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

type IndexRecord struct {
	Language      string   `json:"lang"`
	Month         string   `json:"-"`
	Order         int      `json:"order"`
	OrderChange   *string  `json:"change"`
	OrderPrevYear *int     `json:"orderPreviousYear,omitempty"`
	Rating        float32  `json:"ratingPercent"`
	RatingChange  *float32 `json:"ratingChangePercent"`
	Type          string   `json:"-"`
}

func (ir *IndexRecord) GetOrderChange() string {
	if ir.IsOrderChangeNil() {
		return "nil"
	} else {
		return *ir.OrderChange
	}
}

func (ir *IndexRecord) IsOrderChangeNil() bool {
	return ir.OrderChange == nil
}

func (ir *IndexRecord) IsOrderPrevYearNil() bool {
	return ir.OrderPrevYear == nil
}

func (ir *IndexRecord) GetOrderPrevYear() string {
	if ir.IsOrderPrevYearNil() {
		return "nil"
	} else {
		return fmt.Sprintf("%d", *ir.OrderPrevYear)
	}
}

func (ir *IndexRecord) IsRatingChangeNil() bool {
	return ir.RatingChange == nil
}

func (ir *IndexRecord) GetRatingChange() string {
	if ir.IsRatingChangeNil() {
		return "nil"
	} else {
		return fmt.Sprintf("%.2f", *ir.RatingChange)
	}
}

func (ir *IndexRecord) String() string {

	return fmt.Sprintf("lang: %s, order: %d, month: %s, type: %s, orderChange: %s, rating: %f, ratingChange: %s, orderPrevYear: %s",
		ir.Language, ir.Order, ir.Month, ir.Type, ir.GetOrderChange(), ir.Rating, ir.GetRatingChange(), ir.GetOrderPrevYear())
}

type MonthIndex []*IndexRecord

func (mi *MonthIndex) GetIndex(index int) *IndexRecord {
	for idx, value := range *mi {
		if idx == index {
			return value
		}
	}
	return nil
}

func (mi *MonthIndex) FillMonthAndType(month string, indexType string) {
	for index := range *mi {
		mi.GetIndex(index).Month = month
		mi.GetIndex(index).Type = indexType
	}
}
