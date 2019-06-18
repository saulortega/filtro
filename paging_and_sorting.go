package filtro

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Si SortingColumns es nil se asumirá cualquier parámetro recibido,
// lo cual es susceptible a inyección SQL.
// Si SortingColumns está vacío se ignorará el parámetro de ordenamiento.
type PagingAndSortingInstance struct {
	SortingColumns []string
	RowsPerPage    string
	Descending     string
	SortBy         string
	Page           string
}

// Datos de paginación.
type PagingAndSorting struct {
	Order  string
	Limit  int
	Offset int
}

// Crea una nueva instancia con los parámetros predeterminados.
func NewPagingAndSorting(cols ...string) *PagingAndSortingInstance {
	// Si se quiere establecer nil, debe hacerse explícitamente.
	if cols == nil {
		cols = []string{}
	}

	var I = new(PagingAndSortingInstance)
	I.SortingColumns = cols
	I.RowsPerPage = Params.RowsPerPage
	I.Descending = Params.Descending
	I.SortBy = Params.SortBy
	I.Page = Params.Page
	return I
}

// Obtiene los datos de paginación de la solicitud web.
func (I *PagingAndSortingInstance) Parse(r *http.Request) (*PagingAndSorting, error) {
	var P = new(PagingAndSorting)
	var err error

	lmt := r.FormValue(I.RowsPerPage)
	if lmt == "" {
		return nil, nil
	}

	P.Limit, err = strconv.Atoi(lmt)
	if err != nil {
		return nil, err
	} else if P.Limit < 1 {
		return nil, errors.New("Wrong limit.")
	}

	pge := r.FormValue(I.Page)
	if len(pge) > 0 {
		pgn, err := strconv.Atoi(pge)
		if err != nil {
			return nil, err
		}
		if pgn == 0 {
			pgn = 1
		}
		if pgn < 1 {
			return nil, errors.New("Wrong page.")
		}

		P.Offset = (P.Limit * pgn) - P.Limit
	}

	srt := r.FormValue(I.SortBy)
	if len(srt) > 0 {
		sort := I.SortingColumns == nil || in(srt, I.SortingColumns)
		if sort {
			P.Order = oriArr(srt, I.SortingColumns)

			var desc bool
			des := r.FormValue(I.Descending)
			if des != "" {
				desc, err = strconv.ParseBool(des)
				if err != nil {
					return nil, err
				}
			}

			if desc {
				P.Order += " DESC"
			} else {
				P.Order += " ASC"
			}
		}
	}

	return P, nil
}

func (I *PagingAndSortingInstance) ParseFormatted(r *http.Request, maxLimit int, colOrdnPrdtmnda ...string) (string, error) {
	var PAS string

	var Pgncn, err = I.Parse(r)
	if err != nil {
		return "", err
	}

	if Pgncn != nil {
		if len(Pgncn.Order) > 0 {
			PAS = fmt.Sprintf("ORDER BY %s", Pgncn.Order)
		} else if len(colOrdnPrdtmnda) == 1 {
			PAS = fmt.Sprintf("ORDER BY %s", colOrdnPrdtmnda[0])
		}

		if Pgncn.Limit > 0 && Pgncn.Limit < maxLimit {
			PAS = fmt.Sprintf("%s LIMIT %v", PAS, Pgncn.Limit)
		} else {
			PAS = fmt.Sprintf("%s LIMIT %v", PAS, maxLimit)
		}

		if Pgncn.Offset > 0 {
			PAS = fmt.Sprintf("%s OFFSET %v", PAS, Pgncn.Offset)
		}
	}

	if PAS == "" && len(colOrdnPrdtmnda) == 1 {
		PAS = fmt.Sprintf("ORDER BY %s", colOrdnPrdtmnda[0])
	}

	if Pgncn == nil {
		PAS = fmt.Sprintf("%s LIMIT %v", PAS, maxLimit)
	}

	PAS = strings.TrimSpace(PAS)

	return PAS, nil
}

//
//
//

func in(col string, arr []string) bool {
	for _, c := range arr {
		if colVlr(c) == col {
			return true
		}
	}

	return false
}

func oriArr(col string, arr []string) string {
	if arr == nil {
		return col
	}

	for _, c := range arr {
		if colVlr(c) == col {
			return c
		}
	}

	return col
}
