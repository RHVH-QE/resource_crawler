package crawlers

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// Rhvh4xImg is
type Rhvh4xImg struct {
	BuildName  string `bson:"_id"`
	URLs       []string
	Downloaded bool
}

// Rhvh4xImgCrawler is
type Rhvh4xImgCrawler struct {
	Conf
	BrewRootURL string
}

// NewRhvh4xImgCrawler is
func NewRhvh4xImgCrawler(conf map[string]string) *Rhvh4xImgCrawler {
	ret := Rhvh4xImgCrawler{}
	ret.CrawlerName = conf["crawler_name"]
	ret.CollectionName = conf["col_name"]
	ret.AllowDomains = conf["domains"]
	ret.StartURL = conf["start_url"]
	ret.BrewRootURL = conf["brew_root_url"]
	return &ret
}

func (r Rhvh4xImgCrawler) getBuildName(url string) (buildName string) {
	ret := strings.Split(url, "/")
	buildName = ret[len(ret)-1]
	buildName = strings.Replace(buildName, ".x86_64.liveimg.squashfs", "", 1)
	return buildName
}

// ColName is
func (r Rhvh4xImgCrawler) ColName() string {
	return r.CollectionName
}

// Crawl is start point
func (r Rhvh4xImgCrawler) Crawl() (result interface{}, err error) {
	log.Infof("crawl %s is starting", r.CrawlerName)

	var res []Rhvh4xImg

	c := colly.NewCollector(
		colly.AllowedDomains(r.AllowDomains),
		colly.CacheDir(fmt.Sprintf("/tmp/%s_cache", r.CrawlerName)),
	)
	c.SetRequestTimeout(time.Second * 30)
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	d := c.Clone()

	c.OnHTML(`tr[class*=row-]`, func(e *colly.HTMLElement) {
		var buildURL string
		e.ForEach("td", func(i int, e *colly.HTMLElement) {
			switch i {
			case 0:
				buildURL = e.ChildAttr("a", "href")
			case 3:
				if e.Attr("class") == "complete" {
					err = d.Visit(fmt.Sprintf("%s%s", r.BrewRootURL, buildURL))
					CheckError(err)
				}
			}
		})
	})

	d.OnHTML("table", func(e *colly.HTMLElement) {
		var rhvh4xImg Rhvh4xImg
		urls := e.ChildAttrs(`a[href*="liveimg.squashfs"],a[href*="host-image-update"]`, "href")
		if len(urls) == 2 {
			rhvh4xImg.BuildName = r.getBuildName(urls[1])
			rhvh4xImg.URLs = urls
			rhvh4xImg.Downloaded = false
			res = append(res, rhvh4xImg)
		}
	})

	err = c.Visit(r.StartURL)
	if err != nil {
		log.Error(err)
		return res, err
	}

	return res, nil
}
