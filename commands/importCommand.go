package commands

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/chytilp/supStats/common"
	"github.com/chytilp/supStats/model"
	"github.com/chytilp/supStats/persistence"
	"github.com/chytilp/supStats/request"
)

type ImportCommand struct {
	DB          *sql.DB
	InputFolder string
	Version     int
	dataTable   persistence.SupDataTable
}

type FileImportResult struct {
	Folder   string
	Filename string
	Imported bool
	Error    *error
	DataRows int
}

func NewFileImportResult(filename string, folder string) FileImportResult {
	return FileImportResult{
		Folder:   folder,
		Filename: filename,
	}
}

func NewImportCommand(db *sql.DB, inputFolder string, version int) ImportCommand {
	return ImportCommand{DB: db, InputFolder: inputFolder, Version: version}
}

func (i *ImportCommand) setDataTable(dataTable persistence.SupDataTable) {
	i.dataTable = dataTable
}

func (i *ImportCommand) insertDataItem(item request.Item, filename string) error {
	row := model.NewSupdataRow(item, i.Version, filename)
	err := i.dataTable.InsertRow(row)
	if err != nil {
		return err
	}
	return nil
}

func (i *ImportCommand) insertItem(item request.Item, filename string) error {
	err := i.insertDataItem(item, filename)
	if err != nil {
		return err
	}
	return nil
}

func (i *ImportCommand) insertItemAndChildren(item request.Item, filename string) error {
	err := i.insertItem(item, filename)
	if err != nil {
		return err
	}
	for _, child := range item.Children {
		err = i.insertItem(child, filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *ImportCommand) getFolderFiles() []string {
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
		if !v.IsDir() && common.IsCorrectFileFormat(v.Name(), i.Version) {
			resultFiles = append(resultFiles, v.Name())
		}
	}
	return resultFiles
}

func (i *ImportCommand) Run() []FileImportResult {
	files := i.getFolderFiles()
	results := make([]FileImportResult, len(files))
	for index, file := range files {
		results[index] = i.RunFile(file, i.InputFolder)
	}
	return results
}

func (i *ImportCommand) setResult(result FileImportResult, err error) FileImportResult {
	result.Error = &err
	result.Imported = false
	return result
}

func (i *ImportCommand) RunFile(filename string, folder string) FileImportResult {
	result := NewFileImportResult(filename, folder)
	inputFilePath := folder + "/" + filename
	data, err := request.ReadData(inputFilePath)
	if err != nil {
		return i.setResult(result, err)
	}
	date := data.DateInString()
	i.setDataTable(persistence.SupDataTable{
		DB: i.DB,
	})
	existsPtr, err := i.dataTable.ExistsDate(date)
	if err != nil {
		return i.setResult(result, err)
	}
	if *existsPtr {
		err = fmt.Errorf("date %s already exists in database", date)
		return i.setResult(result, err)
	}
	//---------------------- inserts ------------------
	err = i.insertItemAndChildren(*data.Backend, filename)
	if err != nil {
		return i.setResult(result, err)
	}
	err = i.insertItemAndChildren(*data.Frontend, filename)
	if err != nil {
		return i.setResult(result, err)
	}
	err = i.insertItemAndChildren(*data.Mobile, filename)
	if err != nil {
		return i.setResult(result, err)
	}
	//------------------ stats --------------------
	dRows, err := i.dataTable.GetRows()
	if err != nil {
		fmt.Printf("Data table rows error: %v\n", err)
		result.DataRows = -1
	} else {
		result.DataRows = dRows
	}
	fmt.Printf("data table rows: %d\n", dRows)
	result.Imported = true
	return result
}
