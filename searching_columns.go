package filtro

import (
	"strings"
)

type ValueType int

const (
	ValueTypeString ValueType = iota + 1
	ValueTypeInt64
	ValueTypeBool
)

type SearchType int

const (
	SearchTypeIn SearchType = iota + 1
	SearchTypeLike
	SearchTypeBoolean
)

type Column struct {
	searchType SearchType
	valuetype  ValueType
	name       string
}

func InString(cols ...string) []Column {
	return columnsType(SearchTypeIn, ValueTypeString, cols...)
}

func InInt64(cols ...string) []Column {
	return columnsType(SearchTypeIn, ValueTypeInt64, cols...)
}

func Boolean(cols ...string) []Column {
	return columnsType(SearchTypeBoolean, ValueTypeBool, cols...)
}

func Like(cols ...string) []Column {
	return columnsType(SearchTypeLike, ValueTypeString, cols...)
}

func columnsType(sType SearchType, vType ValueType, cols ...string) []Column {
	var Cols = []Column{}

	for _, c := range cols {
		col := strings.TrimSpace(c)
		if col == "" {
			continue
		}

		Col := Column{sType, vType, col}
		Cols = append(Cols, Col)
	}

	return Cols
}
