package model

import (
	"fmt"
)

type IndexRow struct {
	Lang          string   `json:"lang"`
	Order         int      `json:"order"`
	OrderChange   *string  `json:"change"`
	OrderPrevYear *int     `json:"orderPreviousYear"`
	Rating        float64  `json:"ratingPercent"`
	RatingChange  *float64 `json:"ratingChangePercent"`
}

func (i *IndexRow) String() string {
	var change string = "null"
	if i.OrderChange != nil {
		change = *i.OrderChange
	}
	var orderPrev string = "null"
	if i.OrderPrevYear != nil {
		orderPrev = fmt.Sprintf("%d", *i.OrderPrevYear)
	}
	var ratingChange string = "null"
	if i.RatingChange != nil {
		ratingChange = fmt.Sprintf("%f", *i.RatingChange)
	}
	return fmt.Sprintf("lang: %s, order: %d, orderChange: %s, orderPrev: %s, Rating: %f, RatingChange: %s",
		i.Lang, i.Order, change, orderPrev, i.Rating, ratingChange)
}

func NewTiobeTopRow(lang string, order int, orderChange string, orderPrevYear int, rating float64, ratingChange float64) IndexRow {
	return IndexRow{
		Lang:          lang,
		Order:         order,
		OrderChange:   &orderChange,
		OrderPrevYear: &orderPrevYear,
		Rating:        rating,
		RatingChange:  &ratingChange,
	}
}

func NewTiobeOtherRow(lang string, order int, rating float64) IndexRow {
	return IndexRow{
		Lang:          lang,
		Order:         order,
		OrderChange:   nil,
		OrderPrevYear: nil,
		Rating:        rating,
		RatingChange:  nil,
	}
}

func NewPyplRow(lang string, order int, orderChange string, rating float64, ratingChange float64) IndexRow {
	return IndexRow{
		Lang:          lang,
		Order:         order,
		OrderChange:   &orderChange,
		OrderPrevYear: nil,
		Rating:        rating,
		RatingChange:  &ratingChange,
	}
}
