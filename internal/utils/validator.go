package utils

import (
	"regexp"
)

func IsValidTaskName(input string) bool {
	// Task27: Regex
	match, _ := regexp.MatchString(`[a-zA-Z]+`, input)
	return match
}
