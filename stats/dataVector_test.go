package stats

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

var vector Vector[int] = NewVector(
	"C--",
	[]int{9, 9, 11, 10, 13},
)
var vectorF Vector[float64] = NewVector(
	"Go",
	[]float64{3.12, 2.63, 1.71, 9.13, 7.98},
	2,
)

var vectorF2 Vector[float64] = NewVector(
	"Go",
	[]float64{3.12, 2.63, 1.71, 9.13, 7.98, 5.39},
	2,
)

var vector2 Vector[int] = NewVector(
	"J--",
	[]int{9, 9, 11, 10, 13, 12},
)

func TestVectorMin(t *testing.T) {
	min := vector.Min()
	assert.Equal(t, min, 9)
	minF := vectorF.Min()
	assert.Equal(t, minF, 1.71)
}

func TestVectorMax(t *testing.T) {
	max := vector.Max()
	assert.Equal(t, max, 13)
	maxF := vectorF.Max()
	assert.Equal(t, maxF, 9.13)
}

func TestVectorMean(t *testing.T) {
	mean := vector.Mean()
	assert.Equal(t, mean, 10)
	meanF := vectorF.Mean()
	assert.Equal(t, meanF, 4.91)
}

func TestVectorMedian(t *testing.T) {
	median := vector.Median()
	assert.Equal(t, median, 10)
	median = vector2.Median()
	assert.Equal(t, median, 11)
	medianF := vectorF.Median()
	assert.Equal(t, medianF, 3.12)
	medianF = vectorF2.Median()
	assert.Equal(t, medianF, 4.26)
}
