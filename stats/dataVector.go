package stats

import (
	"math"
	"sort"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type Vector[T Number] struct {
	Label         string
	values        []T
	decimalPLaces int
}

func NewVector[T Number](label string, values []T, optionalParams ...int) Vector[T] {
	decimalPlaces := 0
	if len(optionalParams) > 0 {
		decimalPlaces = optionalParams[0]
	}
	return Vector[T]{Label: label, values: values, decimalPLaces: decimalPlaces}
}

func sortSlice[T Number](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

func (v *Vector[T]) Min() T {
	var min T = v.values[0]
	for _, value := range v.values[1:] {
		if value < min {
			min = value
		}
	}
	return min
}

func (v *Vector[T]) Max() T {
	var max T = v.values[0]
	for _, value := range v.values[1:] {
		if value > max {
			max = value
		}
	}
	return max
}

func (v *Vector[T]) Mean() T {
	var sum T = T(0)
	for _, value := range v.values {
		sum += value
	}
	mean := float64(sum) / float64(len(v.values))
	placesConstant := int(math.Pow10(v.decimalPLaces))
	roundedMean := math.Round(mean*float64(placesConstant)) / float64(placesConstant)
	return T(roundedMean)
}

func (v *Vector[T]) middleIndexes(values []T) []int {
	if len(values)%2 != 0 {
		middle := int(math.Floor(float64(len(values)) / float64(2)))
		return []int{middle}
	} else {
		upperMiddle := len(values) / 2
		lowerMiddle := upperMiddle - 1
		return []int{lowerMiddle, upperMiddle}
	}
}

func (v *Vector[T]) Median() T {
	var values []T = make([]T, len(v.values))
	index := 0
	for _, value := range v.values {
		values[index] = value
		index += 1
	}
	sortSlice(values)
	middleIndexes := v.middleIndexes(values)
	if len(middleIndexes) == 1 {
		return values[middleIndexes[0]]
	}
	avg := float64(values[middleIndexes[0]]+values[middleIndexes[1]]) / float64(2)
	placesConstant := int(math.Pow10(v.decimalPLaces))
	roundedAvg := math.Round(avg*float64(placesConstant)) / float64(placesConstant)
	return T(roundedAvg)
}
