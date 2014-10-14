package main

import (
	"strings"
)

func isEmptyString(s string) bool {
	return len(s) == 0 || len(strings.TrimSpace(s)) == 0
}
