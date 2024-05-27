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

type Table struct {
	data map[TwoStringKey]int
	mu   sync.RWMutex
}

func NewTable() Table {
	return Table{data: make(map[TwoStringKey]int)}
}

func (t *Table) AddValue(row string, column string, value int) error {
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

func (t *Table) AddValues(column string, rowValues map[string]int) error {
	for row, value := range rowValues {
		err := t.AddValue(row, column, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Table) Row(row string) map[string]int {
	rowValues := make(map[string]int)
	for key, value := range t.data {
		if key.row == row {
			rowValues[key.column] = value
		}
	}
	return rowValues
}

func (t *Table) sortHeaders(headers []string) []string {
	sort.Strings(headers)
	return headers
}

func (t *Table) RowHeaders() []string {
	headers := make(map[string]bool)
	for key, _ := range t.data {
		_, ok := headers[key.row]
		if !ok {
			headers[key.row] = true
		}
	}

	return t.sortHeaders(maps.Keys(headers))
}

func (t *Table) ColumnHeaders() []string {
	headers := make(map[string]bool)
	for key, _ := range t.data {
		_, ok := headers[key.column]
		if !ok {
			headers[key.column] = true
		}
	}
	return t.sortHeaders(maps.Keys(headers))
}
