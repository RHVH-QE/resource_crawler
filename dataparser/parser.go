package dataparser

import (
	"github.com/globalsign/mgo"
	log "github.com/sirupsen/logrus"

	cs "github.com/dracher/resource_crawler/crawlers"
)

// CrawledDatabase is
type CrawledDatabase struct {
	DB *mgo.Database
}

// NewCrawledDatabase is
func NewCrawledDatabase(db *mgo.Database) *CrawledDatabase {
	return &CrawledDatabase{db}
}

// SaveCrawledData is
func (cd CrawledDatabase) SaveCrawledData(colName string, data interface{}) {
	bk := cd.DB.C(colName).Bulk()
	bk.Unordered()
	defer bk.Run()

	switch v := data.(type) {
	case []cs.Rhvh4xISO:
		for _, i := range v {
			bk.Insert(i)
		}
	case []cs.Rhvh4xImg:
		for _, i := range v {
			bk.Insert(i)
		}
	case []cs.RhevmBuild:
		for _, i := range v {
			bk.Insert(i)
		}
	default:
		log.Error("don't know about type %T!\n", v)
	}
}
