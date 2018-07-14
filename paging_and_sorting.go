package filtro

import (
	"errors"
	"net/http"
	"strconv"
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
			P.Order = srt

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

//
//
//

func in(col string, arr []string) bool {
	for _, c := range arr {
		if c == col {
			return true
		}
	}

	return false
}
