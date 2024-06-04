package stats

import (
	"fmt"
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

func (t *Table[T]) sortHeaders(headers []string) []string {
	sort.Strings(headers)
	return headers
}

func (t *Table[T]) RowHeaders() []string {
	headers := make(map[string]bool)
	for key, _ := range t.data {
		_, ok := headers[key.row]
		if !ok {
			headers[key.row] = true
		}
	}

	return t.sortHeaders(maps.Keys(headers))
}

func (t *Table[T]) ColumnHeaders() []string {
	headers := make(map[string]bool)
	for key, _ := range t.data {
		_, ok := headers[key.column]
		if !ok {
			headers[key.column] = true
		}
	}
	return t.sortHeaders(maps.Keys(headers))
}
