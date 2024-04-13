package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TrimAllSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func CapitilizeString(s string) string {
	return cases.Title(language.English, cases.Compact).String(s)
}

func EventDetailsStringCleaner(m map[string]string) map[string]string {

	for k, v := range m {
		m[k] = CapitilizeString(TrimAllSpace(v))
	}

	return m
}
