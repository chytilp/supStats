package stats

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"golang.org/x/exp/maps"
)

type TwoStringKey struct {
	row    string
	column string
}

func NewTwoStringKey(row string, column string) TwoStringKey {
	return TwoStringKey{row: row, column: column}
}

type Table[T Number] struct {
	data          map[TwoStringKey]T
	mu            sync.RWMutex
	decimalPLaces int
}

func NewTable[T Number](optional ...int) Table[T] {
	decimalPLaces := 0
	if len(optional) > 0 {
		decimalPLaces = optional[0]
	}
	return Table[T]{data: make(map[TwoStringKey]T), decimalPLaces: decimalPLaces}
}

func (t *Table[T]) Copy(columns []string) Table[T] {
	newTable := NewTable[T](t.decimalPLaces)
	for _, column := range columns {
		values := t.Column(column)
		_ = newTable.AddValues(column, values)
	}
	return newTable
}

func (t *Table[T]) AddValue(row string, column string, value T) error {
	key := NewTwoStringKey(row, column)
	t.mu.Lock()
	_, ok := t.data[key]
	t.mu.Unlock()
	if ok {
		return fmt.Errorf("value row: %s, column: %s", row, column)
	}

	t.mu.Lock()
	t.data[key] = value
	t.mu.Unlock()
	return nil
}

func (t *Table[T]) AddValues(column string, rowValues map[string]T) error {
	for row, value := range rowValues {
		err := t.AddValue(row, column, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Table[T]) Row(row string) map[string]T {
	rowValues := make(map[string]T)
	for key, value := range t.data {
		if key.row == row {
			rowValues[key.column] = value
		}
	}
	return rowValues
}

func (t *Table[T]) Column(column string) map[string]T {
	columnValues := make(map[string]T)
	for key, value := range t.data {
		if key.column == column {
			columnValues[key.row] = value
		}
	}
	return columnValues
}

func (t *Table[T]) sortHeaders(headers []string) []string {
	sort.Strings(headers)
	return headers
}

func (t *Table[T]) RowHeaders() []string {
	headers := make(map[string]bool)
	for key := range t.data {
		_, ok := headers[key.row]
		if !ok {
			headers[key.row] = true
		}
	}

	return t.sortHeaders(maps.Keys(headers))
}

func (t *Table[T]) ColumnHeaders() []string {
	headers := make(map[string]bool)
	for key := range t.data {
		_, ok := headers[key.column]
		if !ok {
			headers[key.column] = true
		}
	}
	return t.sortHeaders(maps.Keys(headers))
}

func rateCombination(dist1 int, dist2 int, dist3 int) int {
	return int(math.Abs(float64(dist1-dist2)) + math.Abs(float64(dist2-dist3)) + math.Abs(float64(dist1-dist3)))
}

func select3Items(data []string) []string {
	middleIndex := len(data) / 2
	return []string{data[0], data[middleIndex], data[len(data)-1]}
}

func select4Items(data []string) []string {
	startIndex := 0
	endIndex := len(data) - 1
	currentFromStart := startIndex
	currentFromEnd := endIndex
	result := [...]int{startIndex, 0, 0, endIndex}
	minRate := 100
	var dist1, dist2, dist3 int
	for currentFromStart <= currentFromEnd {
		currentFromStart += 1
		currentFromEnd -= 1
		dist1 = currentFromStart - startIndex
		dist2 = currentFromEnd - currentFromStart
		dist3 = endIndex - currentFromEnd
		rate := rateCombination(dist1, dist2, dist3)
		if rate < minRate {
			minRate = rate
			result[1] = currentFromStart
			result[2] = currentFromEnd
		}
	}
	return []string{data[result[0]], data[result[1]], data[result[2]], data[result[3]]}
}

func selectColumns(data []string, count int) []string {
	if count == 2 {
		return []string{data[0], data[len(data)-1]}
	}
	if count == 3 {
		return select3Items(data)
	}
	return select4Items(data)
}

func TableWithSelectedColumns[T Number](table *Table[T], columnCount int) (*Table[T], error) {
	if columnCount > 4 {
		return nil, fmt.Errorf("select more than 4 column not implemented, count: %d", columnCount)
	}
	columns := table.ColumnHeaders()
	filteredColumns := selectColumns(columns, columnCount)
	newTable := table.Copy(filteredColumns)
	return &newTable, nil
}
