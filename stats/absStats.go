package stats

import (
	"fmt"
	"math"
	"sort"
	"time"
)

type Day string

type Language string

type LanguageRow struct {
	Name   Language
	Values map[Day]int
}

func (l *LanguageRow) AddDay(day Day, value int) error {
	_, ok := l.Values[day]
	if ok {
		return fmt.Errorf("Day %s in values already exists", day)
	}
	l.Values[day] = value
	return nil
}

func (l *LanguageRow) Mean() int {
	sum := 0
	for _, value := range l.Values {
		sum += value
	}
	mean := float64(sum) / float64(len(l.Values))
	return int(math.Round(mean))
}

func (l *LanguageRow) middleIndexes(values []int) []int {
	if len(values)%2 != 0 {
		middle := int(math.Floor(float64(len(values)) / float64(2)))
		return []int{middle}
	} else {
		upperMiddle := len(values) / 2
		lowerMiddle := upperMiddle - 1
		return []int{lowerMiddle, upperMiddle}
	}
}

func (l *LanguageRow) Median() int {
	var values []int = make([]int, len(l.Values))
	index := 0
	for _, value := range l.Values {
		values[index] = value
		index += 1
	}
	sort.Ints(values)
	middleIndexes := l.middleIndexes(values)
	if len(middleIndexes) == 1 {
		return values[middleIndexes[0]]
	}
	return int(math.Round(float64(values[middleIndexes[0]]+values[middleIndexes[1]]) / float64(2)))
}

func (l *LanguageRow) Min() int {
	min := math.MaxInt
	for _, value := range l.Values {
		if value < min {
			min = value
		}
	}
	return min
}

func (l *LanguageRow) Max() int {
	max := 0
	for _, value := range l.Values {
		if value > max {
			max = value
		}
	}
	return max
}

type StatInput struct {
	DateFrom  time.Time
	DateTo    time.Time
	Tech      Technology
	Max       bool
	Min       bool
	Mean      bool
	Median    bool
	languages *map[Language]LanguageRow
	table     *Table
}
