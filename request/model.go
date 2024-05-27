package request

import (
	"time"
)

type Item struct {
	Identifier    string `json:"identifier"`
	UrlIdentifier string `json:"urlIdentifier"`
	Name          string `json:"name"`
	Children      []Item `json:"children"`
	OfferCount    int    `json:"offer_count"`
}

type ResponseData struct {
	Root []Item `json:"initiallySelectedOptions"`
}

type OutputData struct {
	Frontend     *Item     `json:"frontend"`
	Backend      *Item     `json:"backend"`
	Mobile       *Item     `json:"mobile"`
	DownloadedAt time.Time `json:"downloaded"`
}

func (o *OutputData) Day() string {
	filename := GetFileName(o.DownloadedAt) // data_2024_05_20.json
	return filename[5:15]
}
