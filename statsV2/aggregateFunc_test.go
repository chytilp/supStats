package statsV2

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestAggregateMin(t *testing.T) {
	values := []int{1, 5, 3, 7, 6, 11}
	min_val, err := Min(values)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, min_val, 1.0)
}

func TestAggregateMax(t *testing.T) {
	values := []int{1, 5, 3, 7, 6, 11}
	max_val, err := Max(values)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, max_val, 11.0)
}

func TestAggregateMean(t *testing.T) {
	values := []int{1, 5, 3, 7, 6, 11}
	avg_val, err := Mean(values)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, avg_val, 5.5)
}

func TestAggregateMedian(t *testing.T) {
	values := []int{1, 5, 3, 7, 6, 11}
	med_val, err := Median(values)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, med_val, 5.5)
}
