package filtro

import (
	"fmt"
	"regexp"
	"strings"
)

// Normalizar texto de una columna de BD para mayor precisión en búsqueda.
func normalizeColumn(col string) string {
	return fmt.Sprintf("LOWER(TRANSLATE(%s, 'ÁÉÍÓÚÜáéíóúüÀÈÌÒÙàèìòùÑ', 'aeiouuaeiouuaeiouaeiouñ'))", col)
}

// Normalizar texto para comparación de búsqueda.
func normalizeText(t string) string {
	t = strings.ToLower(t)
	t = strings.TrimSpace(t)
	t = regexp.MustCompile(`\s+`).ReplaceAllString(t, " ")
	t = regexp.MustCompile("á").ReplaceAllString(t, "a")
	t = regexp.MustCompile("é").ReplaceAllString(t, "e")
	t = regexp.MustCompile("í").ReplaceAllString(t, "i")
	t = regexp.MustCompile("ó").ReplaceAllString(t, "o")
	t = regexp.MustCompile("ú").ReplaceAllString(t, "u")
	t = regexp.MustCompile("ü").ReplaceAllString(t, "u")
	t = regexp.MustCompile("à").ReplaceAllString(t, "a")
	t = regexp.MustCompile("è").ReplaceAllString(t, "e")
	t = regexp.MustCompile("ì").ReplaceAllString(t, "i")
	t = regexp.MustCompile("ò").ReplaceAllString(t, "o")
	t = regexp.MustCompile("ù").ReplaceAllString(t, "u")
	return t
}

func words(t string) []string {
	var wrds = []string{}
	t = normalizeText(t)

	for _, p := range strings.Split(t, " ") {
		if len(p) <= 2 || p == "las" || p == "los" || p == "les" || p == "una" || p == "por" {
			continue
		}
		wrds = append(wrds, p)
	}

	return wrds
}
