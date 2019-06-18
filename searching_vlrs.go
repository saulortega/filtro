package filtro

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func vlrsString(r *http.Request, col string) []string {
	var vlrs = []string{}
	col = colVlr(col)

	if _, exte := r.Form[col]; !exte {
		col = col + "[]"
	}

	for _, v := range r.Form[col] {
		val := strings.TrimSpace(v)
		if val == "" {
			continue
		}

		vlrs = append(vlrs, val)
	}

	return vlrs
}

func vlrsInt64(r *http.Request, col string) ([]int64, error) {
	var vlrsStr = vlrsString(r, col)
	var vlrs = []int64{}

	for _, v := range vlrsStr {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return vlrs, errors.New("Invalid values for «" + col + "»")
		}

		vlrs = append(vlrs, i)
	}

	return vlrs, nil
}

// El segundo Bool indica si hay un valor presente.
// Se debe evaluar primero el error, y luego el segundo bool.
func vlrBool(r *http.Request, col string) (bool, bool, error) {
	var val = r.FormValue(colVlr(col))
	if val == "" {
		return false, false, nil
	}

	var vlr, err = strconv.ParseBool(val)
	if err != nil {
		return false, true, errors.New("Invalid values for «" + col + "»")
	}

	return vlr, true, nil
}

func colVlr(C string) string {
	var cc = strings.Split(C, ".")
	return cc[len(cc)-1]
}