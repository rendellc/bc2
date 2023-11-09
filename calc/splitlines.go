package calc

import (
	"strings"
)


func SplitLines(script string) []string {
	return strings.Split(strings.ReplaceAll(script, "\r\n", "\n"), "\n")
}

