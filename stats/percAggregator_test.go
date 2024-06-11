package stats

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestAggregatePercent(t *testing.T) {
	table := NewTable[int]()
	table.AddValue("c++", "2024-05-01", 5)
	table.AddValue("java", "2024-05-01", 7)
	table.AddValue("c#", "2024-05-01", 6)
	table.AddValue("python", "2024-05-01", 20)
	aggregator := PercentAggregator[int]{table: table}
	newTable := aggregator.Aggregate()
	assert.Equal(t, newTable.ColumnHeaders(), []string{"2024-05-01"})
	rows := newTable.RowHeaders()
	assert.Equal(t, len(rows), 4)
	assert.Equal(t, newTable.Row("c++")["2024-05-01"], 13.16)
	assert.Equal(t, newTable.Row("c#")["2024-05-01"], 15.79)
	assert.Equal(t, newTable.Row("java")["2024-05-01"], 18.42)
	assert.Equal(t, newTable.Row("python")["2024-05-01"], 52.63)
}
