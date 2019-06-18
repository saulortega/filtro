package filtro

import (
	"fmt"
	"net/http"
	"strings"
)

//
//
//

type Condition struct {
	Clause     string
	Args       []interface{}
	searchType SearchType
	valuetype  ValueType
	column     string
}

func (c *Condition) SearchType() SearchType {
	return c.searchType
}

func (c *Condition) ValueType() ValueType {
	return c.valuetype
}

func (c *Condition) Column() string {
	return c.column
}

//
//
//

type SearchingInstance struct {
	Columnas []Column
}

func NewSearching(gruposColumnas ...[]Column) *SearchingInstance {
	var I = new(SearchingInstance)
	I.Columnas = []Column{}

	for _, GC := range gruposColumnas {
		for _, C := range GC {
			I.Columnas = append(I.Columnas, C)
		}
	}

	return I
}

func (I *SearchingInstance) InString(cols ...string) {
	I.Columnas = append(I.Columnas, InString(cols...)...)
}

func (I *SearchingInstance) InInt64(cols ...string) {
	I.Columnas = append(I.Columnas, InInt64(cols...)...)
}

func (I *SearchingInstance) Boolean(cols ...string) {
	I.Columnas = append(I.Columnas, Boolean(cols...)...)
}

func (I *SearchingInstance) Like(cols ...string) {
	I.Columnas = append(I.Columnas, Like(cols...)...)
}

func (I *SearchingInstance) Parse(r *http.Request) ([]*Condition, error) {
	var condiciones = []*Condition{}
	var err error

	for _, col := range I.Columnas {
		switch col.searchType {
		case SearchTypeIn:
			C := new(Condition)
			if col.valuetype == ValueTypeString {
				C = CndcnInString(r, col)
			} else if col.valuetype == ValueTypeInt64 {
				C, err = CndcnInInt64(r, col)
				if err != nil {
					return condiciones, err
				}
			}
			if C != nil {
				condiciones = append(condiciones, C)
			}
		case SearchTypeBoolean:
			C, err := CndcnBoolean(r, col)
			if err != nil {
				return condiciones, err
			} else if C != nil {
				condiciones = append(condiciones, C)
			}
		case SearchTypeLike:
			C := CndcnLike(r, col)
			if C != nil {
				condiciones = append(condiciones, C)
			}
		}
	}

	return condiciones, nil
}

func (I *SearchingInstance) ParseFormatted(r *http.Request, Oprdr string) (string, []interface{}, error) {
	var Clss = []string{}
	var Args = []interface{}{}

	var Cdcns, err = I.Parse(r)
	if err != nil {
		return "", Args, err
	}

	var likesClauses = []string{}
	var likesArgs = []interface{}{}
	for _, c := range Cdcns {
		switch c.SearchType() {
		case SearchTypeIn, SearchTypeBoolean:
			Clss = append(Clss, c.Clause)
			if c.Args != nil {
				Args = append(Args, c.Args...)
			}
		case SearchTypeLike:
			likesClauses = append(likesClauses, c.Clause)
			likesArgs = append(likesArgs, c.Args...)
		}
	}

	if len(likesClauses) > 0 {
		Clss = append(Clss, fmt.Sprintf(`(%s)`, strings.Join(likesClauses, " OR ")))
		Args = append(Args, likesArgs...)
	}

	var AllClss = strings.Join(Clss, fmt.Sprintf(" %s ", strings.TrimSpace(Oprdr)))

	return AllClss, Args, nil
}

//
//
//

func CndcnInString(r *http.Request, col Column) *Condition {
	var C = new(Condition)
	C.Args = []interface{}{}

	var vlrs = vlrsString(r, col.name)
	if len(vlrs) == 0 {
		return nil
	}

	for _, vlr := range vlrs {
		C.Args = append(C.Args, interface{}(vlr))
	}

	C.Clause = fmt.Sprintf(`%s IN ?`, col.name)
	C.searchType = col.searchType
	C.valuetype = col.valuetype
	C.column = col.name

	return C
}

func CndcnInInt64(r *http.Request, col Column) (*Condition, error) {
	var C = new(Condition)
	C.Args = []interface{}{}

	var vlrs, err = vlrsInt64(r, col.name)
	if err != nil {
		return nil, err
	} else if len(vlrs) == 0 {
		return nil, nil
	}

	for _, vlr := range vlrs {
		C.Args = append(C.Args, interface{}(vlr))
	}

	C.Clause = fmt.Sprintf(`%s IN ?`, col.name)
	C.searchType = col.searchType
	C.valuetype = col.valuetype
	C.column = col.name

	return C, nil
}

func CndcnBoolean(r *http.Request, col Column) (*Condition, error) {
	var vlr, prsnte, err = vlrBool(r, col.name)
	if err != nil {
		return nil, err
	} else if !prsnte {
		return nil, nil
	}

	var cndcn = col.name
	if !vlr {
		cndcn = fmt.Sprintf(`NOT %s`, col.name)
	}

	var C = new(Condition)
	C.Clause = cndcn
	C.searchType = col.searchType
	C.valuetype = col.valuetype
	C.column = col.name

	return C, nil
}

func CndcnLike(r *http.Request, col Column) *Condition {
	var C = new(Condition)
	C.Args = []interface{}{}

	var clmn = normalizeColumn(col.name)
	var wrds = words(r.FormValue(Params.Search))
	if len(wrds) == 0 || len(clmn) == 0 {
		return nil
	}

	var ors = []string{}
	for _, wrd := range wrds {
		ors = append(ors, fmt.Sprintf(`%s LIKE ?`, clmn))
		C.Args = append(C.Args, interface{}(fmt.Sprintf(`%%%s%%`, wrd)))
	}

	C.Clause = fmt.Sprintf(`(%s)`, strings.Join(ors, " OR "))
	C.searchType = col.searchType
	C.valuetype = col.valuetype
	C.column = col.name

	return C
}
