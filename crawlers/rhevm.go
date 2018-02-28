package crawlers

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// RhevmBuild is
type RhevmBuild struct {
	BuildName  string `bson:"_id"`
	URLs       []string
	Downloaded bool
}

// RhevmBuildCrawler is
type RhevmBuildCrawler struct {
	Conf
	RpmURL string
}

// NewRhevmBuildCrawler to init rhvh42 crawler
func NewRhevmBuildCrawler(conf map[string]string) *RhevmBuildCrawler {
	ret := RhevmBuildCrawler{}
	ret.CrawlerName = conf["crawler_name"]
	ret.CollectionName = conf["col_name"]
	ret.AllowDomains = conf["domains"]
	ret.StartURL = conf["start_url"]
	ret.RpmURL = conf["rpm_url"]
	return &ret
}

// ColName is
func (r RhevmBuildCrawler) ColName() string {
	return r.CollectionName
}

func (r RhevmBuildCrawler) parse(data []byte) (buildName string, ready bool) {
	ready = bytes.Contains(data, []byte("Ready"))
	if !ready {
		log.Debugf("Not ready, %s", data)
		return buildName, ready
	}
	log.Debugf("Build is Ready: %s", data)
	pattern := regexp.MustCompile(`\[rhv-.+\]`)
	ret := pattern.Find(data)

	return fmt.Sprintf("%s", bytes.Trim(ret, "[]")), ready
}

// Crawl is
func (r RhevmBuildCrawler) Crawl() (result interface{}, err error) {
	log.Infof("crawl %s is starting", r.CrawlerName)
	var res []RhevmBuild

	c := colly.NewCollector(
		colly.AllowedDomains(r.AllowDomains),
		colly.CacheDir(fmt.Sprintf("/tmp/%s_cache", r.CrawlerName)),
	)
	c.SetRequestTimeout(time.Second * 30)
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.OnResponse(func(req *colly.Response) {
		if name, ready := r.parse(req.Body); ready {
			rhevmBuild := RhevmBuild{
				BuildName:  name,
				URLs:       []string{r.RpmURL},
				Downloaded: false,
			}
			res = append(res, rhevmBuild)
		}
	})

	err = c.Visit(r.StartURL)
	if err != nil {
		log.Error(err)
		return result, err
	}

	return res, nil
}
