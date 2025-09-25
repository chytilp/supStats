package commands

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chytilp/supStats/common"
	"github.com/chytilp/supStats/model"
	"github.com/chytilp/supStats/persistence"
	"github.com/chytilp/supStats/request"
)

type ImportIndexCommand struct {
	DB          *sql.DB
	InputFolder string
	dataTable   persistence.IndexDataTable
}

func NewImportIndexCommand(db *sql.DB, inputFolder string) ImportIndexCommand {
	return ImportIndexCommand{DB: db, InputFolder: inputFolder}
}

func (i *ImportIndexCommand) getFolderFiles() []string {
	resultFiles := []string{}
	folder, err := os.Open(i.InputFolder)
	if err != nil {
		fmt.Println(err)
		return resultFiles
	}
	files, err := folder.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return resultFiles
	}

	for _, v := range files {
		if !v.IsDir() && common.IsCorrectIndexFile(v.Name()) {
			resultFiles = append(resultFiles, v.Name())
		}
	}
	return resultFiles
}

func (i *ImportIndexCommand) getFileMonthAndType(fileName string) (string, string, error) {
	fileParts := strings.Split(fileName, ".")
	parts := strings.Split(fileParts[0], "_")
	if len(parts) < 3 {
		return "", "", fmt.Errorf("wrong filename format: %s", fileName)
	}
	year, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", "", err
	}
	month, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", "", err
	}
	day := 1
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	output := fmt.Sprintf("%04d-%02d-%02d 00:00:00", date.Year(), date.Month(), date.Day())
	return output, parts[0], nil
}

func (i *ImportIndexCommand) Run() []FileImportResult {
	files := i.getFolderFiles()
	results := make([]FileImportResult, len(files))
	for index, file := range files {
		results[index] = i.RunFile(file, i.InputFolder)
	}
	return results
}

func (i *ImportIndexCommand) RunFile(filename string, folder string) FileImportResult {
	wholePath := folder + "/" + filename
	result := NewFileImportResult(filename, folder)
	i.dataTable = persistence.IndexDataTable{DB: i.DB}
	month, indexType, err := i.getFileMonthAndType(wholePath)
	if err != nil {
		return result.SetErrorResult(err)
	}
	originalRows, err := i.dataTable.GetRows(indexType)
	if err != nil {
		return result.SetErrorResult(err)
	}
	obj, err := request.UnmarshalFromFile[model.MonthIndex](wholePath)
	if err != nil {
		return result.SetErrorResult(err)
	}
	obj.FillMonthAndType(month, indexType)
	count := len(*obj)
	var record model.IndexRecord
	for j := 0; j < count; j++ {
		record = *obj.GetIndex(j)
		err = i.dataTable.InsertRow(record)
		if err != nil {
			return result.SetErrorResult(err)
		}
	}
	// -- stats --
	dRows, err := i.dataTable.GetRows(indexType)
	if err != nil {
		fmt.Printf("Data table rows error: %v\n", err)
		result.DataRows = -1
	} else {
		result.DataRows = dRows - originalRows
	}
	fmt.Printf("data table new rows: %d\n", dRows-originalRows)
	result.Imported = true
	return result
}
