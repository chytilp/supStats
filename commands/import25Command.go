package commands

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/chytilp/supStats/model"
	"github.com/chytilp/supStats/persistence"
	"github.com/chytilp/supStats/request"
)

type Import25Command struct {
	DB          *sql.DB
	InputFolder string
	Version     int
	dataTable   persistence.SupDataTable
}

func NewImport25Command(db *sql.DB, inputFolder string) Import25Command {
	return Import25Command{DB: db, InputFolder: inputFolder, Version: 25}
}

func (i *Import25Command) Run() []FileImportResult {
	files := getFolderFiles(i.InputFolder, i.Version)
	results := make([]FileImportResult, len(files))
	for index, file := range files {
		results[index] = i.RunFile(file, i.InputFolder)
	}
	return results
}

func (i *Import25Command) insertItem(item model.FileContentItem, itemType string, date time.Time) error {
	row := model.NewSupdata25Row(item, i.Version, itemType, date)
	err := i.dataTable.InsertRow(row)
	if err != nil {
		return err
	}
	return nil
}

func (i *Import25Command) RunFile(filename string, folder string) FileImportResult {
	result := NewFileImportResult(filename, folder)
	i.dataTable = persistence.SupDataTable{DB: i.DB}
	originalRows, err := i.dataTable.GetRows(i.Version)
	if err != nil {
		return result.SetErrorResult(err)
	}
	inputFilePath := folder + "/" + filename
	obj, err := request.UnmarshalFromFile[model.FileContent](inputFilePath)
	if err != nil {
		return result.SetErrorResult(err)
	}
	var date string = obj.DateInString()
	existsPtr, err := i.dataTable.ExistsDate(date, i.Version)
	if err != nil {
		return result.SetErrorResult(err)
	}
	if *existsPtr {
		err = fmt.Errorf("date %s already exists in database table supdata (version: %d)", date, i.Version)
		return result.SetErrorResult(err)
	}
	//---------------------- inserts ------------------
	//-- categories --
	for _, categoryItem := range obj.Categories {
		err = i.insertItem(categoryItem, "category", obj.DownloadedAt)
		if err != nil {
			return result.SetErrorResult(err)
		}
	}
	//-- technologies --
	for _, technologyItem := range obj.Technologies {
		err = i.insertItem(technologyItem, "technology", obj.DownloadedAt)
		if err != nil {
			return result.SetErrorResult(err)
		}
	}
	//------------------ stats --------------------
	dRows, err := i.dataTable.GetRows(i.Version)
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
