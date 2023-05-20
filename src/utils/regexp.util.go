package util

import (
	"regexp"
	"strings"
)

func SnakeCase(s string) string {
	r := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	snake := r.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}
