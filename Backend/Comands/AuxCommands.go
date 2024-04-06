package Comands

import (
	"strings"
)

func Compare(a string, b string) bool {
	if strings.ToUpper(a) == strings.ToUpper(b) {
		return true
	}
	return false
}
