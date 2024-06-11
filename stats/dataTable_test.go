package stats

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
	"golang.org/x/exp/maps"
)

func TestTableAddValue(t *testing.T) {
	table := NewTable[int]()
	table.AddValue("java", "2024-05-01", 1)
	table.AddValue("java", "2024-05-02", 2)
	table.AddValue("java", "2024-05-03", 3)
	table.AddValue("java", "2024-05-04", 4)
	table.AddValue("java", "2024-05-05", 5)
	daysValues := table.Row("java")
	assert.Equal(t, len(daysValues), 5)
	for i := 1; i <= 5; i++ {
		assert.Equal(t, daysValues[fmt.Sprintf("2024-05-0%d", i)], i)
	}
}

func TestTableColumn(t *testing.T) {
	table := NewTable[int]()
	table.AddValue("java", "2024-05-01", 1)
	table.AddValue("c++", "2024-05-01", 2)
	table.AddValue("python", "2024-05-01", 3)
	table.AddValue("java", "2024-05-02", 2)
	table.AddValue("java", "2024-05-03", 3)
	table.AddValue("java", "2024-05-04", 4)
	table.AddValue("java", "2024-05-05", 5)
	daysValues := table.Column("2024-05-01")
	assert.Equal(t, len(daysValues), 3)
	assert.Equal(t, daysValues["java"], 1)
	assert.Equal(t, daysValues["c++"], 2)
	assert.Equal(t, daysValues["python"], 3)
}

func TestTableCopy(t *testing.T) {
	table := NewTable[int]()
	table.AddValue("java", "2024-05-01", 1)
	table.AddValue("c++", "2024-05-01", 2)
	table.AddValue("python", "2024-05-01", 3)
	table.AddValue("java", "2024-05-02", 2)
	table.AddValue("java", "2024-05-03", 3)
	table.AddValue("java", "2024-05-04", 4)
	table.AddValue("java", "2024-05-05", 5)
	requiredColumns := []string{"2024-05-01", "2024-05-02"}
	newTable := table.Copy(requiredColumns)
	assert.Equal(t, newTable.ColumnHeaders(), requiredColumns)
	values := newTable.Column("2024-05-01")
	assert.Equal(t, len(values), 3)
	rowValues := newTable.Row("java")
	assert.Equal(t, len(rowValues), 2)
	assert.Equal(t, rowValues["2024-05-01"], 1)
	assert.Equal(t, rowValues["2024-05-02"], 2)
}

func TestTableRowHeaders(t *testing.T) {
	table := NewTable[int]()
	table.AddValue("java", "2024-05-01", 1)
	table.AddValue("c#", "2024-05-01", 2)
	table.AddValue("python", "2024-05-01", 3)
	table.AddValue("java", "2024-05-02", 4)
	table.AddValue("java", "2024-05-03", 5)
	langs := table.RowHeaders()
	expected := map[string]bool{
		"c#":     false,
		"java":   false,
		"python": false,
	}
	assert.Equal(t, len(langs), len(expected))
	assert.Equal(t, langs, maps.Keys(expected))
}

func TestTableColumnHeaders(t *testing.T) {
	table := NewTable[int]()
	table.AddValue("java", "2024-05-02", 4)
	table.AddValue("java", "2024-05-03", 5)
	table.AddValue("java", "2024-05-01", 1)
	table.AddValue("c#", "2024-05-01", 2)
	table.AddValue("python", "2024-05-01", 3)
	days := table.ColumnHeaders()
	expected := map[string]bool{
		"2024-05-01": false,
		"2024-05-02": false,
		"2024-05-03": false,
	}
	assert.Equal(t, len(days), len(expected))
	assert.Equal(t, days, maps.Keys(expected))
}

func TestSelectColumns(t *testing.T) {
	columns := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	selected := selectColumns(columns, 2)
	assert.Equal(t, []string{"1", "10"}, selected)
	selected = selectColumns(columns, 3)
	assert.Equal(t, []string{"1", "6", "10"}, selected)
	selected = selectColumns(columns, 4)
	assert.Equal(t, []string{"1", "4", "7", "10"}, selected)
}
