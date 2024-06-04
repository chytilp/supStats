package stats

import (
	"fmt"
	"strings"
)

type TableDisplay[T Number] struct {
	table *Table[T]
}

func NewDisplay[T Number](table *Table[T]) TableDisplay[T] {
	return TableDisplay[T]{table: table}
}

func (t *TableDisplay[T]) Lines4Print() []string {
	langs := t.table.RowHeaders()
	maxLen := 0
	for _, lang := range langs {
		if len(lang) > maxLen {
			maxLen = len(lang)
		}
	}
	const firstColumnHeader = "Lang"
	if len(firstColumnHeader) > maxLen {
		maxLen = len(firstColumnHeader)
	}
	maxLen += 1
	days := t.table.ColumnHeaders()
	lines := make([]string, 0, len(langs)+1)
	const otherWidth = 12
	lines = append(lines, t.formatRow(firstColumnHeader, days, maxLen, otherWidth))
	for _, lang := range langs {
		values := t.table.Row(lang)
		strValues := make([]string, 0, len(values))
		for _, day := range days {
			strValues = append(strValues, fmt.Sprintf("%v", values[day]))
		}
		lines = append(lines, t.formatRow(lang, strValues, maxLen, otherWidth))
	}
	return lines
}

func (t *TableDisplay[T]) charactersForGap(totalChars int, textChars int, gapInBeginning bool) []string {
	if gapInBeginning {
		return []string{strings.Repeat(" ", totalChars-textChars)}
	} else {
		gapCount := totalChars - textChars
		if gapCount%2 == 0 {
			halfGap := strings.Repeat(" ", gapCount/2)
			return []string{
				halfGap,
				halfGap,
			}
		} else {
			firstGap := strings.Repeat(" ", (gapCount/2)+1)
			lastGet := strings.Repeat(" ", gapCount/2)
			return []string{
				firstGap,
				lastGet,
			}
		}
	}
}

func (t *TableDisplay[T]) formatRow(firstColumn string, otherColumns []string, firstWidth int, otherWidth int) string {
	s := firstColumn + t.charactersForGap(firstWidth, len(firstColumn), true)[0] + "|"
	for _, column := range otherColumns {
		gaps := t.charactersForGap(otherWidth, len(column), false)
		s += gaps[0] + column + gaps[1] + "|"
	}
	return s
}
