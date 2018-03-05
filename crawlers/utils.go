package crawlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Crawler is
type Crawler interface {
	Crawl() (interface{}, error) // start point to start crawling
	ColName() string             // the mongo collection name which the data will save in
}

// Conf holds basic crawler fields
type Conf struct {
	CrawlerName    string
	CollectionName string
	AllowDomains   string
	StartURL       string
}

// NewCrawler is
func NewCrawler(crawlerName string, conf map[string]interface{}) (c Crawler) {
	switch crawlerName {
	case "rhvh4x_img":
		ret := Rhvh4xImgCrawler{}
		ret.CrawlerName = conf["crawler_name"].(string)
		ret.CollectionName = conf["col_name"].(string)
		ret.AllowDomains = conf["domains"].(string)
		ret.StartURL = conf["start_url"].(string)
		ret.BrewRootURL = conf["brew_root_url"].(string)
		return &ret
	case "rhvh4x_iso":
		ret := Rhvh4xISOCrawler{}
		ret.CrawlerName = conf["crawler_name"].(string)
		ret.CollectionName = conf["col_name"].(string)
		ret.AllowDomains = conf["domains"].(string)
		ret.StartURL = conf["start_url"].(string)
		ret.IsoURLTpl = conf["iso_url_tpl"].(string)
		return &ret
	case "rhevm_build":
		ret := RhevmBuildCrawler{}
		ret.CrawlerName = conf["crawler_name"].(string)
		ret.CollectionName = conf["col_name"].(string)
		ret.AllowDomains = conf["domains"].(string)
		ret.StartURL = conf["start_url"].(string)
		ret.RpmURL = conf["rpm_url"].(string)
		return &ret
	default:
		log.Errorf("%s is unknown", crawlerName)
		return nil
	}
}

func fetchSimple(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%s", body))
}

// CheckError is
func CheckError(err error) {
	if err != nil {
		log.Error(err)
	}
}
