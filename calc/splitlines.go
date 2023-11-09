package calc

import (
	"bufio"
	"strings"
)


func SplitLines(script string) []string {
	lines := []string{}
	scanner := bufio.NewScanner(strings.NewReader(script))

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

