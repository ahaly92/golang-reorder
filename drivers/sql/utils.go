package sql

import "strings"

func QuoteString(input string) string {
	return strings.Replace(input, "'", "''", -1)
}
