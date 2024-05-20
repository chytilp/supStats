package stats

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

var row LanguageRow = LanguageRow{
	Name: "C--",
	Values: map[Day]int{
		"2024-05-01": 9,
		"2024-05-02": 9,
		"2024-05-03": 11,
		"2024-05-04": 10,
		"2024-05-05": 13,
	}}

var row2 LanguageRow = LanguageRow{
	Name: "J--",
	Values: map[Day]int{
		"2024-05-01": 9,
		"2024-05-02": 9,
		"2024-05-03": 11,
		"2024-05-04": 10,
		"2024-05-05": 13,
		"2024-05-06": 12,
	}}

func TestLanguageRowMin(t *testing.T) {
	min := row.Min()
	assert.Equal(t, min, 9)
}

func TestLanguageRowMax(t *testing.T) {
	max := row.Max()
	assert.Equal(t, max, 13)
}

func TestLanguageRowMean(t *testing.T) {
	mean := row.Mean()
	assert.Equal(t, mean, 10)
}

func TestLanguageRowMedian(t *testing.T) {
	median := row.Median()
	assert.Equal(t, median, 10)
	median = row2.Median()
	assert.Equal(t, median, 11)
}
