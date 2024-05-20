package stats

import (
	"fmt"
	"sort"
	"strings"
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
	mu   sync.Mutex
}

func NewTable() Table {
	return Table{data: make(map[TwoStringKey]int)}
}

func (t *Table) AddValue(row string, column string, value int) error {
	key := NewTwoStringKey(row, column)
	_, ok := t.data[key]
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
	headersSlice := sort.StringSlice(headers)
	headersSlice.Sort()
	return headersSlice
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

func (t *Table) charactersForGap(totalChars int, textChars int, gapInBeginning bool) []string {
	if gapInBeginning {
		return []string{strings.Repeat(" ", totalChars-textChars)}
	} else {
		gapCount := totalChars - textChars
		if gapCount%2 == 0 {
			halfGap := strings.Repeat(" ", gapCount/2)
			return []string{
				halfGap,
				halfGap,
			}
		} else {
			firstGap := strings.Repeat(" ", (gapCount/2)+1)
			lastGet := strings.Repeat(" ", gapCount/2)
			return []string{
				firstGap,
				lastGet,
			}
		}
	}
}

func (t *Table) formatRow(firstColumn string, otherColumns []string, firstWidth int, otherWidth int) string {
	s := firstColumn + t.charactersForGap(firstWidth, len(firstColumn), true)[0] + "|"
	for _, column := range otherColumns {
		gaps := t.charactersForGap(otherWidth, len(column), false)
		s += gaps[0] + column + gaps[1] + "|"
	}
	return s
}

func (t *Table) Lines4Print() []string {
	langs := t.RowHeaders()
	maxLen := 0
	for _, lang := range langs {
		if len(lang) > maxLen {
			maxLen = len(lang)
		}
	}
	const firstColumnHeader = "Lang"
	if len(firstColumnHeader) > maxLen {
		maxLen = len(firstColumnHeader)
	}
	maxLen += 1
	days := t.ColumnHeaders()
	lines := make([]string, 0, len(langs)+1)
	lines = append(lines, t.formatRow(firstColumnHeader, days, maxLen, 12))
	for _, lang := range langs {
		values := t.Row(lang)
		strValues := make([]string, 0, len(values))
		for _, day := range days {
			strValues = append(strValues, fmt.Sprintf("%d", values[day]))
		}
		lines = append(lines, t.formatRow(lang, strValues, maxLen, 12))
	}
	return lines
}
