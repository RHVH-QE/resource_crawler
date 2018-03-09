package crawlers

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// Rhvh4xISO is
type Rhvh4xISO struct {
	BuildName  string `bson:"_id"`
	URLs       []string
	Downloaded bool
}

// Rhvh4xISOCrawler is
type Rhvh4xISOCrawler struct {
	Conf
	IsoURLTpl string
}

// ColName is
func (r Rhvh4xISOCrawler) ColName() string {
	return r.CollectionName
}

// Crawl is
func (r Rhvh4xISOCrawler) Crawl() (result interface{}, err error) {
	log.Infof("crawl %s is starting", r.CrawlerName)
	c := colly.NewCollector(
		colly.AllowedDomains(r.AllowDomains),
		// colly.CacheDir(fmt.Sprintf("/tmp/%s_cache", r.CrawlerName)),
	)
	c.SetRequestTimeout(time.Second * 30)

	var res []Rhvh4xISO

	pattern := regexp.MustCompile(`RHVH-4.[1-9]`)
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if pattern.MatchString(url) {
			status := fetchSimple(fmt.Sprintf("%s%s%s", r.StartURL, url, "STATUS"))
			if status == "FINISHED" {
				isoName := url[:len(url)-1]
				isoURL := fmt.Sprintf("%s%s%s", r.StartURL, url, fmt.Sprintf(r.IsoURLTpl, isoName))
				rhvh42ISO := Rhvh4xISO{}
				rhvh42ISO.BuildName = isoName
				rhvh42ISO.URLs = append(rhvh42ISO.URLs, isoURL)
				rhvh42ISO.Downloaded = false
				res = append(res, rhvh42ISO)
			}
		}
	})
	err = c.Visit(r.StartURL)
	if err != nil {
		log.Error(err)
		return res, err
	}
	return res, nil
}
