package statsV2

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/chytilp/supStats/persistence"
)

type Grid[N persistence.Number] struct {
	languages []string
	months    []string
	values    map[persistence.MyTuple]N
	typeObj   N
}

func NewGrid[N persistence.Number](langs []string, months []string, values map[persistence.MyTuple]N) *Grid[N] {
	grid := Grid[N]{languages: langs, months: months, values: values, typeObj: *new(N)}
	return &grid
}

func (g *Grid[N]) getRow(lang string) map[string]N {
	result := make(map[string]N)
	for _, month := range g.months {
		key := persistence.MyTuple{Language: lang, Month: month}
		value, ok := g.values[key]
		if ok {
			result[month] = value
		}
	}
	return result
}

func (g *Grid[N]) maxLangSize() int {
	size := 0
	for _, lang := range g.languages {
		if len(lang) > size {
			size = utf8.RuneCountInString(lang)
		}
	}
	return size
}

func (g *Grid[N]) formatBySize(value string, maxSize int) string {
	diff := maxSize - len(value)
	uniDiff := len(value) - utf8.RuneCountInString(value)
	if diff < 0 {
		return value + strings.Repeat(" ", maxSize)
	} else {
		if uniDiff == 0 {
			return value + strings.Repeat(" ", diff+2)
		} else {
			return value + strings.Repeat(" ", diff+2+uniDiff)
		}
	}

}

func (g *Grid[N]) getFormat() (string, error) {
	typeName := fmt.Sprintf("%T", g.typeObj)
	switch typeName {
	case "int":
		return "%d", nil
	case "float64":
		return "%.2f", nil
	default:
		return "", fmt.Errorf("unknown type: %s", typeName)
	}
}

func (g *Grid[N]) Print() {
	maxSize := g.maxLangSize()
	// fmt.Println("size: ", maxSize)
	g.getFormat()
	fmt.Printf(g.formatBySize("languages", maxSize))
	months := make([]string, 0)
	for _, month := range g.months {
		fmt.Printf("%s", g.formatBySize(month, 7))
		months = append(months, month)
	}
	fmt.Printf("\n")
	format, err := g.getFormat()
	if err != nil {
		fmt.Printf("Get format error: %v\n", err)
		return
	}
	for _, lang := range g.languages {
		values := g.getRow(lang)
		fmt.Printf("%s", g.formatBySize(lang, maxSize))
		for _, month := range months {
			value, ok := values[month]
			if ok {
				strValue := "NA"
				if value >= 0 {
					strValue = fmt.Sprintf(format, value)
				}
				fmt.Printf("%s", g.formatBySize(strValue, 7))
			} else {
				fmt.Printf("%s", g.formatBySize(fmt.Sprintf(format, ""), 7))
			}

		}
		fmt.Printf("\n")
	}
}
