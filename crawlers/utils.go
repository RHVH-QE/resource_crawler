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
