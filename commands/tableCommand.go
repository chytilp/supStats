package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chytilp/supStats/common"
	"github.com/chytilp/supStats/stats"
)

type TableCommand struct {
	config            *common.Config
	techType          stats.Technology
	fromDate          time.Time
	toDate            time.Time
	columnCount       int
	aggragatedColumns bool
}

func parseDate(stringDate string) time.Time {
	parts := strings.Split(stringDate, "-")
	year, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])
	day, _ := strconv.Atoi(parts[2])
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

func NewTableCommand(config *common.Config, techTYpe stats.Technology, fromTo string, columns int,
	aggColumns bool) TableCommand {
	parts := strings.Split(fromTo, "->")
	fromDate := parseDate(parts[0])
	toDate := parseDate(parts[1])
	return TableCommand{
		config:            config,
		techType:          techTYpe,
		fromDate:          fromDate,
		toDate:            toDate,
		columnCount:       columns,
		aggragatedColumns: aggColumns,
	}
}

func (t *TableCommand) Run() error {
	tableCreate := stats.TableCreate[int]{
		DateFrom: t.fromDate,
		DateTo:   t.toDate,
		Tech:     t.techType,
		Config:   t.config,
	}
	err := tableCreate.ReadData(true)
	if err != nil {
		fmt.Println("err in ReadData")
		return err
	}
	display := stats.NewDisplay(tableCreate.Table())
	lines := display.Lines4Print()
	for _, line := range lines {
		fmt.Println(line)
	}
	return nil
}
