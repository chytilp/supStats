package model

import (
	"github.com/chytilp/supStats/request"
)

type LanguageRow struct {
	Name     string
	Path     string
	ParentId *int
}

func NewLanguageRow(item request.Item, parentId *int) LanguageRow {
	return LanguageRow{
		Name:     item.Name,
		Path:     item.Identifier,
		ParentId: parentId,
	}
}

type DataRow struct {
	LanguageId int
	Count      int
	Date       string
}

func NewDataRow(item request.Item, languageId int, date string) DataRow {
	return DataRow{
		LanguageId: languageId,
		Count:      item.OfferCount,
		Date:       date,
	}
}
