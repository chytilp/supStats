package stats

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestTableDisplayLines4Print(t *testing.T) {
	table := NewTable()
	table.AddValue("java", "2024-05-02", 4)
	table.AddValue("java", "2024-05-03", 5)
	table.AddValue("java", "2024-05-01", 1)
	table.AddValue("c#", "2024-05-01", 2)
	table.AddValue("python", "2024-05-01", 3)
	display := NewDisplay(&table)
	lines := display.Lines4Print()
	assert.Equal(t, lines[0], "Lang   | 2024-05-01 | 2024-05-02 | 2024-05-03 |")
	assert.Equal(t, lines[1], "c#     |      2     |      0     |      0     |")
	assert.Equal(t, lines[2], "java   |      1     |      4     |      5     |")
	assert.Equal(t, lines[3], "python |      3     |      0     |      0     |")
}
