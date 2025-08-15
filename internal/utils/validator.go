package utils

import (
	"regexp"
)

func IsValidTaskName(input string) bool {
	// Task27: Regex
	match, _ := regexp.MatchString(`[a-zA-Z]+`, input)
	return match
}

func IsValidDate(input string) bool {
	// Task27: Regex
	match, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, input)
	return match
}
