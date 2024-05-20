package main

import (
	"fmt"
	"log"

	"github.com/chytilp/supStats/commands"
	"github.com/chytilp/supStats/common"
	"github.com/chytilp/supStats/request"
)

func main() {
	fmt.Println("Start app")
	config := common.GetConfig()
	outputData, err := commands.CreateDailyData(config)
	if err != nil {
		fmt.Println("err in CreateDailyData")
		log.Fatalln(err.Error())
	}
	err = request.MarshalToFile(*outputData)
	if err != nil {
		fmt.Println("err in MarshalToFile")
		log.Fatalln(err.Error())
	}
	fmt.Println("Saved")
}
