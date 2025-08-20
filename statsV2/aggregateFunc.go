package statsV2

import (
	"errors"
	"math"
	"sort"
)

func Min(values []int) (float64, error) {
	if len(values) == 0 {
		return 0.0, errors.New("empty collection on input Min func")
	}
	min_val := values[0]
	for i := 1; i < len(values); i++ {
		if min_val > values[i] {
			min_val = values[i]
		}
	}
	return float64(min_val), nil
}

func Max(values []int) (float64, error) {
	if len(values) == 0 {
		return 0.0, errors.New("empty collection on input Max func")
	}
	max_val := values[0]
	for i := 1; i < len(values); i++ {
		if max_val < values[i] {
			max_val = values[i]
		}
	}
	return float64(max_val), nil
}

func Mean(values []int) (float64, error) {
	if len(values) == 0 {
		return 0.0, errors.New("empty collection on input Mean func")
	}
	sum := 0
	for _, item := range values {
		sum += item
	}
	mean := float64(sum) / float64(len(values))
	return mean, nil
}

func middleIndexes(values []int) []int {
	if len(values)%2 != 0 {
		middle := int(math.Floor(float64(len(values)) / float64(2)))
		return []int{middle}
	} else {
		upperMiddle := len(values) / 2
		lowerMiddle := upperMiddle - 1
		return []int{lowerMiddle, upperMiddle}
	}
}

func Median(values []int) (float64, error) {
	sort.Ints(values)
	middleIndexes := middleIndexes(values)
	if len(middleIndexes) == 1 {
		return float64(values[middleIndexes[0]]), nil
	}
	average := float64(values[middleIndexes[0]]+values[middleIndexes[1]]) / float64(2)
	return average, nil
}
