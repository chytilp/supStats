package stats

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

var vector Vector[int] = NewVector(
	"C--",
	[]int{9, 9, 11, 10, 13},
)

var vector2 Vector[int] = NewVector(
	"J--",
	[]int{9, 9, 11, 10, 13, 12},
)

func TestVectorMin(t *testing.T) {
	min := vector.Min()
	assert.Equal(t, min, 9)
}

func TestVectorMax(t *testing.T) {
	max := vector.Max()
	assert.Equal(t, max, 13)
}

func TestVectorMean(t *testing.T) {
	mean := vector.Mean()
	assert.Equal(t, mean, 10)
}

func TestVectorMedian(t *testing.T) {
	median := vector.Median()
	assert.Equal(t, median, 10)
	median = vector2.Median()
	assert.Equal(t, median, 11)
}
