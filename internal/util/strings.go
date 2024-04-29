package util

import "golang.org/x/exp/constraints"

func Pluralize[T constraints.Integer](singular, plural string, count T) string {
	if count == 1 {
		return singular
	}
	return plural
}
