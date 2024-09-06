package request

import (
	"strings"
	"time"

	"github.com/chytilp/supStats/common"
)

type Item struct {
	Identifier    string `json:"identifier"`
	UrlIdentifier string `json:"urlIdentifier"`
	Name          string `json:"name"`
	Children      []Item `json:"children"`
	OfferCount    int    `json:"offer_count"`
}

func (i *Item) Empty() bool {
	if i.Identifier == "" && i.Name == "" && i.OfferCount == 0 {
		return true
	}
	return false
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
	filename := common.GetFileName(o.DownloadedAt) // data_2024_05_20.json
	return filename[5:15]
}

func (o *OutputData) FindItem(name string) Item {
	feResult := o.findInBranch(name, o.Frontend)
	if feResult != nil {
		return *feResult
	}
	beResult := o.findInBranch(name, o.Backend)
	if beResult != nil {
		return *beResult
	}
	mbResult := o.findInBranch(name, o.Mobile)
	if mbResult != nil {
		return *mbResult
	}
	return Item{}
}

func (o *OutputData) findInBranch(name string, branch *Item) *Item {
	if o.NamesAreSame(name, branch.Name) {
		return branch
	}
	for _, child := range branch.Children {
		if o.NamesAreSame(name, child.Name) {
			return &child
		}
	}
	return nil
}

func (o *OutputData) NamesAreSame(name string, itemName string) bool {
	return strings.EqualFold(name, itemName)
}
