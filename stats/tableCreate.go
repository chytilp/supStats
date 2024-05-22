package stats

import (
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/chytilp/supStats/common"
	"github.com/chytilp/supStats/request"
)

type Technology int

const (
	Frontend Technology = iota
	Backend
	Mobile
)

type TableCreate struct {
	DateFrom time.Time
	DateTo   time.Time
	Tech     Technology
	Config   *common.Config
	table    *Table
}

func (t *TableCreate) folderList() []string {
	months := []string{}
	first := time.Date(t.DateFrom.Year(), t.DateFrom.Month(), 1, 0, 0, 0, 0, nil)
	last := time.Date(t.DateTo.Year(), t.DateTo.Month(), 1, 0, 0, 0, 0, nil)

	for tmp := first; tmp.Before(last); tmp = tmp.AddDate(0, 1, 0) {
		months = append(months, request.GetFolder(tmp))
	}
	return months
}

func (t *TableCreate) folderFiles(folderName string) (*[]string, error) {
	absFolder := t.Config.DataFolder + folderName
	folder, err := os.Open(absFolder)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	files, err := folder.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	absFiles := []string{}
	for _, v := range files {
		if !v.IsDir() && t.correctFileFormat(v.Name()) {
			absFiles = append(absFiles, v.Name())
		}
	}
	return &absFiles, nil
}

func (t *TableCreate) correctFileFormat(fileName string) bool {
	r, _ := regexp.Compile(`^data_([0-9]{4})_([0-9]{2})_([0-9]{2})\.json$`)
	return r.MatchString(fileName)
}

func (t *TableCreate) fileList() (*[]string, error) {
	allFiles := []string{}
	var absFolder string
	var err error
	var files *[]string
	for _, folderName := range t.folderList() {
		absFolder = t.Config.DataFolder + folderName
		files, err = t.folderFiles(absFolder)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		for _, file := range *files {
			allFiles = append(allFiles, absFolder+file)
		}

	}
	return &allFiles, nil
}

func (t *TableCreate) ReadData(inclParent bool) error {
	table := NewTable()
	t.table = &table
	// read files
	files, err := t.fileList()
	if err != nil {
		return err
	}
	// download data from files to table
	var wg sync.WaitGroup
	errChan := make(chan error)
	wg.Add(len(*files))
	for _, file := range *files {
		path := file
		go t.readFile(path, inclParent, &wg, errChan)
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TableCreate) readFile(path string, inclParent bool, wg *sync.WaitGroup, errChan chan<- error) {
	data, err := request.UnmarshalFromFile[request.OutputData](path)
	if err != nil {
		errChan <- fmt.Errorf("UnmarshalFromFile - error in goroutine processing file: %s, err: %v", path, err)
		return
	}
	var root *request.Item
	switch t.Tech {
	case Backend:
		root = data.Backend
	case Frontend:
		root = data.Frontend
	case Mobile:
		root = data.Mobile
	}
	dayOfData := data.Day()
	if inclParent {
		err = t.table.AddValue(root.Name, dayOfData, root.OfferCount)
		if err != nil {
			errChan <- fmt.Errorf("AddValue - error in goroutine processing file: %s, err: %v", path, err)
			return
		}
	}
	for _, child := range root.Children {
		err = t.table.AddValue(child.Name, dayOfData, child.OfferCount)
		if err != nil {
			errChan <- fmt.Errorf("AddValue - error in goroutine processing file: %s, err: %v", path, err)
			return
		}
	}
	fmt.Printf("Data from file %s added to table.\n", path)
	wg.Done()
}
