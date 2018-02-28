package main

import (
	"os"
	"runtime"
	"sync"

	"github.com/globalsign/mgo"
	"github.com/spf13/viper"

	cs "github.com/dracher/resource_crawler/crawlers"
	"github.com/dracher/resource_crawler/dataparser"
	log "github.com/sirupsen/logrus"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	if viper.GetBool("debug") {
		log.Warn("DEBUG MODE IS ON")
		log.SetLevel(log.DebugLevel)
		// do some debugging
		return
	}

	crawlerLogFile, _ := os.Create(viper.GetString("log_file"))
	defer crawlerLogFile.Close()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(crawlerLogFile)

	session, err := mgo.Dial(viper.GetString("database.mongourl"))
	if err != nil {
		log.Panic(err)
	}
	defer session.Close()
	db := dataparser.NewCrawledDatabase(session.DB(viper.GetString("database.dbname")))

	crawlers := []cs.Crawler{
		cs.NewRhvh4xISOCrawler(viper.GetStringMapString("crawlers.rhvh4x_iso")),
		cs.NewRhvh4xImgCrawler(viper.GetStringMapString("crawlers.rhvh4x_img")),
		cs.NewRhevmBuildCrawler(viper.GetStringMapString("crawlers.rhevm_build")),
	}

	var wg sync.WaitGroup
	limit := make(chan uint, viper.GetInt("crawl_conn_limit"))

	for _, crawler := range crawlers {
		limit <- 1
		wg.Add(1)
		go func(c cs.Crawler) {
			data, err := c.Crawl()
			db.SaveCrawledData(c.ColName(), data)
			if err != nil {
				log.Error(err)
			}
			log.Warnf("%s is finished", c.ColName())
			wg.Done()
			<-limit
		}(crawler)
	}
	wg.Wait()
}
