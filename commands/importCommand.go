package commands

import (
	"database/sql"
	"fmt"

	"github.com/chytilp/dbForSupTest/model"
	"github.com/chytilp/dbForSupTest/persistence"
)

type ImportCommand struct {
	DB            *sql.DB
	InputfilePath string
	InputFolder   string
	languageTable persistence.LanguageTable
	dataTable     persistence.DataTable
}

type FileImportResult struct {
	Folder       string
	Filename     string
	Imported     bool
	Error        *error
	LanguageRows int
	DataRows     int
}

func NewFileImportResult(filename string, folder string) FileImportResult {
	return FileImportResult{
		Folder:   folder,
		Filename: filename,
	}
}

func NewImportCommand(db *sql.DB, inputfilePath string) ImportCommand {
	return ImportCommand{DB: db, InputfilePath: inputfilePath}
}

func (i *ImportCommand) setLanguageTable(languageTable persistence.LanguageTable) {
	i.languageTable = languageTable
}

func (i *ImportCommand) setDataTable(dataTable persistence.DataTable) {
	i.dataTable = dataTable
}

func (i *ImportCommand) insertLanguage(item model.Item, parentId *int) (int, error) {
	idPtr, err := i.languageTable.GetLanguageId(item.Identifier)
	if err != nil {
		return 0, err
	}
	if *idPtr > 0 {
		return *idPtr, nil
	}
	row := model.NewLanguageRow(item, parentId)
	newId, err := i.languageTable.InsertLanguage(row)
	if err != nil {
		return 0, err
	}
	return newId, err
}

func (i *ImportCommand) insertDataItem(item model.Item, languageId int, date string) error {
	row := model.NewDataRow(item, languageId, date)
	err := i.dataTable.InsertDataItem(row)
	if err != nil {
		return err
	}
	return nil
}

func (i *ImportCommand) insertItem(item model.Item, parentId *int, date string) (int, error) {
	languageId, err := i.insertLanguage(item, parentId)
	if err != nil {
		return 0, err
	}
	err = i.insertDataItem(item, languageId, date)
	if err != nil {
		return 0, err
	}
	return languageId, nil
}

func (i *ImportCommand) insertItemAndChildren(item model.Item, parentId *int, date string) error {
	id, err := i.insertItem(item, parentId, date)
	if err != nil {
		return err
	}
	for _, child := range item.Children {
		_, err = i.insertItem(child, &id, date)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *ImportCommand) getFolderFiles() []string {
	// TODO:
	files := make([]string, 0)
	return files
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
	data, err := model.ReadData(inputFilePath)
	if err != nil {
		return i.setResult(result, err)
	}
	date := data.DateInString()
	i.setLanguageTable(persistence.LanguageTable{
		DB: i.DB,
	})
	i.setDataTable(persistence.DataTable{
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
	err = i.insertItemAndChildren(*data.Backend, nil, date)
	if err != nil {
		return i.setResult(result, err)
	}
	err = i.insertItemAndChildren(*data.Frontend, nil, date)
	if err != nil {
		return i.setResult(result, err)
	}
	err = i.insertItemAndChildren(*data.Mobile, nil, date)
	if err != nil {
		return i.setResult(result, err)
	}
	//------------------ stats --------------------
	lRows, err := i.languageTable.GetRows()
	if err != nil {
		fmt.Printf("Language table rows error: %v\n", err)
		result.LanguageRows = -1
	} else {
		result.LanguageRows = lRows
	}
	fmt.Printf("Language table rows: %d\n", lRows)
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
