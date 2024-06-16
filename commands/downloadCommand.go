package commands

import (
	"fmt"
	"sync"
	"time"

	"github.com/chytilp/supStats/common"
	"github.com/chytilp/supStats/request"
)

type DownloadCommand struct {
	config *common.Config
}

func NewDownloadCommand(config *common.Config) DownloadCommand {
	return DownloadCommand{config: config}
}

func (d *DownloadCommand) Run() (*string, error) {
	outputData, err := d.downloadDailyData()
	if err != nil {
		return nil, err
	}
	filePath, err := request.MarshalToFile(*outputData, d.config)
	if err != nil {
		fmt.Println("err in MarshalToFile")
		return nil, err
	}
	return filePath, nil
}

func (d *DownloadCommand) downloadDailyData() (*request.OutputData, error) {
	var wg sync.WaitGroup
	wg.Add(3)
	fe, feErr := d.scrapeAndParseData(d.config.FrontendUrl, &wg)
	be, beErr := d.scrapeAndParseData(d.config.BackendUrl, &wg)
	mb, mbErr := d.scrapeAndParseData(d.config.MobileUrl, &wg)
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

func (d *DownloadCommand) scrapeAndParseData(url string, wg *sync.WaitGroup) (*request.ResponseData, error) {
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
