package commands

import (
	"fmt"
	"sync"
	"time"

	"github.com/chytilp/supStats/common"
	"github.com/chytilp/supStats/request"
)

func scrapeAndParseData(url string, wg *sync.WaitGroup) (*request.ResponseData, error) {
	defer wg.Done()
	data, err := request.SendRequest(url)
	if err != nil {
		fmt.Println("err in SendRequest")
		return nil, err
	}
	result, err := request.ParseResponse(data)
	if err != nil {
		fmt.Println("err in ParseResponse")
		return nil, err
	}
	return result, nil
}

func CreateDailyData(config *common.Config) (*request.OutputData, error) {
	var wg sync.WaitGroup
	wg.Add(3)
	fe, feErr := scrapeAndParseData(config.FrontendUrl, &wg)
	be, beErr := scrapeAndParseData(config.BackendUrl, &wg)
	mb, mbErr := scrapeAndParseData(config.MobileUrl, &wg)
	wg.Wait()
	if feErr != nil {
		return nil, feErr
	}
	if beErr != nil {
		return nil, beErr
	}
	if mbErr != nil {
		return nil, mbErr
	}
	outData := request.OutputData{
		Frontend:     &fe.Root[0],
		Backend:      &be.Root[0],
		Mobile:       &mb.Root[0],
		DownloadedAt: time.Now(),
	}
	return &outData, nil
}
