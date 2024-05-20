package convertor

import (
	"encoding/json"
	"io/ioutil"

	"github.com/chytilp/supStats/request"
)

func ReadDataFile(datafilePath string) (*[]OldItem, error) {
	data, err := ioutil.ReadFile(datafilePath)
	if err != nil {
		return nil, err
	}
	var items []OldItem
	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, err
	}
	return &items, nil
}

func Transform(oldItems []OldItem) (*request.OutputData, error) {
	return nil, nil
}
