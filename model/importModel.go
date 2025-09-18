package model

import (
	"fmt"
	"time"

	"github.com/chytilp/supStats/common"
	"github.com/chytilp/supStats/request"
)

type RecordType int

const (
	Category   RecordType = 1
	Technology RecordType = 2
)

type SupdataRow struct {
	Language string
	Type     int
	Count    int
	Version  int
	Date     string
}

func getType(isTechnology bool) int {
	if isTechnology {
		return 2
	} else {
		return 1
	}
}

func NewSupdataRow(item request.Item, version int, filename string) SupdataRow {
	var dayPtr *time.Time
	dayPtr, _ = common.GetFileDate(filename)
	var date time.Time = *dayPtr
	day := fmt.Sprintf("%04d-%02d-%02d 00:00:00", date.Year(), date.Month(), date.Day())
	var isTechnology bool = len(item.Children) == 0
	return SupdataRow{
		Language: item.Name,
		Count:    item.OfferCount,
		Type:     getType(isTechnology),
		Version:  version,
		Date:     day,
	}
}

func NewSupdata25Row(item FileContentItem, version int, itemType string, date time.Time) SupdataRow {
	day := fmt.Sprintf("%04d-%02d-%02d 00:00:00", date.Year(), date.Month(), date.Day())
	var isTechnology bool = (itemType == "technology")
	return SupdataRow{
		Language: item.Name,
		Count:    item.Count,
		Type:     getType(isTechnology),
		Version:  version,
		Date:     day,
	}
}

type IndexdataRow struct {
	Language  string
	IndexType int
	Rating    float32
	Order     int
	Date      string
}

func getIndexType(name string) int {
	if name == "tiobe" {
		return 1
	} else if name == "pypl" {
		return 2
	}
	return 0
}

func NewIndexdataRow(indexName string, lang string, rating float32, order int, month string) IndexdataRow {
	return IndexdataRow{
		Language:  lang,
		IndexType: getIndexType(indexName),
		Rating:    rating,
		Order:     order,
		Date:      month,
	}
}
