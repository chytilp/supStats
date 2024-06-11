package stats

import (
	"math"
)

type PercentAggregator[T Number] struct {
	table Table[T]
}

func (p *PercentAggregator[T]) sum(values map[string]T) T {
	var sum T
	for _, value := range values {
		sum += value
	}
	return sum
}

func (p *PercentAggregator[T]) aggregateColumn(column string) map[string]float64 {
	columnValues := p.table.Column(column)
	sum := p.sum(columnValues)
	result := map[string]float64{}
	for key, value := range columnValues {
		perc := (float64(value) / float64(sum)) * 100.0
		perc = math.Round(perc*100) / 100
		result[key] = perc
	}
	return result
}

func (p *PercentAggregator[T]) Aggregate() Table[float64] {
	newTable := NewTable[float64]()
	for _, column := range p.table.ColumnHeaders() {
		values := p.aggregateColumn(column)
		_ = newTable.AddValues(column, values)
	}
	return newTable
}
