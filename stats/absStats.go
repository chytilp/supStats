package stats

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/chytilp/supStats/request"
)

type Technology int

type Day string

type Language string

const (
	Frontend Technology = iota
	Backend
	Mobile
)

type LanguageRow struct {
	Name   Language
	Values map[Day]int
}

func (l *LanguageRow) AddDay(day Day, value int) error {
	_, ok := l.Values[day]
	if ok {
		return fmt.Errorf("Day %s in values already exists", day)
	}
	l.Values[day] = value
	return nil
}

func (l *LanguageRow) Mean() int {
	sum := 0
	for _, value := range l.Values {
		sum += value
	}
	mean := float64(sum) / float64(len(l.Values))
	return int(math.Round(mean))
}

func (l *LanguageRow) middleIndexes(values []int) []int {
	if len(values)%2 != 0 {
		middle := int(math.Floor(float64(len(values)) / float64(2)))
		return []int{middle}
	} else {
		upperMiddle := len(values) / 2
		lowerMiddle := upperMiddle - 1
		return []int{lowerMiddle, upperMiddle}
	}
}

func (l *LanguageRow) Median() int {
	var values []int = make([]int, len(l.Values))
	index := 0
	for _, value := range l.Values {
		values[index] = value
		index += 1
	}
	sort.Ints(values)
	middleIndexes := l.middleIndexes(values)
	if len(middleIndexes) == 1 {
		return values[middleIndexes[0]]
	}
	return int(math.Round(float64(values[middleIndexes[0]]+values[middleIndexes[1]]) / float64(2)))
}

func (l *LanguageRow) Min() int {
	min := math.MaxInt
	for _, value := range l.Values {
		if value < min {
			min = value
		}
	}
	return min
}

func (l *LanguageRow) Max() int {
	max := 0
	for _, value := range l.Values {
		if value > max {
			max = value
		}
	}
	return max
}

type StatInput struct {
	DateFrom  time.Time
	DateTo    time.Time
	Tech      Technology
	Max       bool
	Min       bool
	Mean      bool
	Median    bool
	languages *map[Language]LanguageRow
	table     *Table
}

func (s *StatInput) fileList() (*[]string, error) {
	return nil, nil
}

func (s *StatInput) ReadData() error {
	langMap := map[Language]LanguageRow{}
	s.languages = &langMap
	table := NewTable()
	s.table = &table
	// read files
	files, err := s.fileList()
	if err != nil {
		return err
	}
	// download data from files to table
	var wg sync.WaitGroup
	errChan := make(chan error)
	wg.Add(len(*files))
	for _, file := range *files {
		path := file
		go s.readFile(path, &wg, errChan)
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

func (s *StatInput) readFile(path string, wg *sync.WaitGroup, errChan chan<- error) {
	data, err := request.UnmarshalFromFile[request.OutputData](path)
	if err != nil {
		errChan <- fmt.Errorf("UnmarshalFromFile - error in goroutine processing file: %s, err: %v", path, err)
		return
	}
	var root *request.Item
	switch s.Tech {
	case Backend:
		root = data.Backend
	case Frontend:
		root = data.Frontend
	case Mobile:
		root = data.Mobile
	}
	dayOfData := data.Day()
	err = s.table.AddValue(root.Name, dayOfData, root.OfferCount)
	if err != nil {
		errChan <- fmt.Errorf("AddValue - error in goroutine processing file: %s, err: %v", path, err)
		return
	}
	for _, child := range root.Children {
		err = s.table.AddValue(child.Name, dayOfData, child.OfferCount)
		if err != nil {
			errChan <- fmt.Errorf("AddValue - error in goroutine processing file: %s, err: %v", path, err)
			return
		}
	}
	fmt.Printf("Data from file %s added to table.\n", path)
	wg.Done()
}
