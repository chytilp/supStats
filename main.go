package main

import (
	"fmt"
	"log"
	"time"

	//"github.com/chytilp/supStats/commands"
	"github.com/chytilp/supStats/common"
	//"github.com/chytilp/supStats/request"
	"github.com/chytilp/supStats/stats"
)

func main() {
	fmt.Println("Start app")
	config := common.GetConfig()
	/*outputData, err := commands.CreateDailyData(config)
	if err != nil {
		fmt.Println("err in CreateDailyData")
		log.Fatalln(err.Error())
	}
	err = request.MarshalToFile(*outputData, config)
	if err != nil {
		fmt.Println("err in MarshalToFile")
		log.Fatalln(err.Error())
	}
	fmt.Println("Saved")*/
	tableCreate := stats.TableCreate{
		DateFrom: time.Date(2024, time.May, 1, 0, 0, 0, 0, time.Local),
		DateTo:   time.Date(2024, time.May, 31, 0, 0, 0, 0, time.Local),
		Tech:     stats.Frontend,
		Config:   config,
	}
	err := tableCreate.ReadData(true)
	if err != nil {
		fmt.Println("err in ReadData")
		log.Fatalln(err.Error())
	}
	display := stats.NewDisplay(tableCreate.Table())
	lines := display.Lines4Print()
	for _, line := range lines {
		fmt.Println(line)
	}
}
