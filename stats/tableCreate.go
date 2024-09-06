package stats

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
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
	Unknown
)

func TechnologyFromString(stringTechnology string) Technology {
	switch stringTechnology {
	case "fe":
		return Frontend
	case "be":
		return Backend
	case "mb":
		return Mobile
	default:
		return Unknown
	}
}

type TableCreate[T Number] struct {
	DateFrom time.Time
	DateTo   time.Time
	Tech     Technology
	Config   *common.Config
	table    *Table[T]
}

func (t *TableCreate[T]) folderList() []string {
	months := []string{}
	first := time.Date(t.DateFrom.Year(), t.DateFrom.Month(), 1, 0, 0, 0, 0, time.Local)
	last := time.Date(t.DateTo.Year(), t.DateTo.Month(), 1, 0, 0, 0, 0, time.Local)
	tmp := first
	for tmp == last || tmp.Before(last) {
		months = append(months, common.GetFolder(tmp))
		tmp = tmp.AddDate(0, 1, 0)
	}
	return months
}

func (t *TableCreate[T]) includeFile(filename string, firstDay int, lastDay int) bool {
	parts := strings.Split(filename, ".")
	parts2 := strings.Split(parts[0], "_")
	num, _ := strconv.Atoi(parts2[3])
	return num >= firstDay && num <= lastDay
}

func (t *TableCreate[T]) folderFiles(absFolder string, firstDay int, lastDay int) (*[]string, error) {
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
		if !v.IsDir() && t.correctFileFormat(v.Name()) && t.includeFile(v.Name(), firstDay, lastDay) {
			absFiles = append(absFiles, v.Name())
		}
	}
	return &absFiles, nil
}

func (t *TableCreate[T]) daysInMonth(folder string) int {
	parts := strings.Split(folder, "-")
	year, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])
	m := time.Month(month + 1)
	dt := time.Date(year, m, 1, 0, 0, 0, 0, time.UTC)
	dt2 := dt.AddDate(0, 0, -1)
	return dt2.Day()
}

func (t *TableCreate[T]) correctFileFormat(fileName string) bool {
	r, _ := regexp.Compile(`^data_([0-9]{4})_([0-9]{2})_([0-9]{2})\.json$`)
	return r.MatchString(fileName)
}

func (t *TableCreate[T]) fileList() (*[]string, error) {
	allFiles := []string{}
	var absFolder string
	var err error
	var files *[]string
	var firstDay int = 1
	var lastDay int
	folderList := t.folderList()
	var lastIndex int = len(folderList) - 1
	for index, folderName := range folderList {
		absFolder = t.Config.DataFolder + "/" + folderName
		if index == 0 {
			firstDay = t.DateFrom.Day()
		} else {
			firstDay = 1
		}

		if index == lastIndex {
			lastDay = t.DateTo.Day()
		} else {
			lastDay = t.daysInMonth(folderName)
		}
		files, err = t.folderFiles(absFolder, firstDay, lastDay)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		for _, file := range *files {
			allFiles = append(allFiles, absFolder+file)
		}
	}
	sort.Strings(allFiles)
	return &allFiles, nil
}

func (t *TableCreate[T]) ReadData(inclParent bool) error {
	table := NewTable[T]()
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

func (t *TableCreate[T]) readFile(path string, inclParent bool, wg *sync.WaitGroup, errChan chan<- error) {
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
		err = t.table.AddValue(root.Name, dayOfData, T(root.OfferCount))
		if err != nil {
			errChan <- fmt.Errorf("AddValue - error in goroutine processing file: %s, err: %v", path, err)
			return
		}
	}
	for _, child := range root.Children {
		err = t.table.AddValue(child.Name, dayOfData, T(child.OfferCount))
		if err != nil {
			errChan <- fmt.Errorf("AddValue - error in goroutine processing file: %s, err: %v", path, err)
			return
		}
	}
	fmt.Printf("Data from file %s added to table.\n", path)
	wg.Done()
}

func (t *TableCreate[T]) Table() *Table[T] {
	return t.table
}
