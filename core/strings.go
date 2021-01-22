package core

import (
	"fmt"
	"regexp"
	"strings"
)

// Pluralize will take in a count and if the count is != 1 it will return the singular of the word.
func Pluralize(count int, singular string, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

// TruncateString will take an integer and cut a string at that length and append an ellipsis to it.
func TruncateString(str string, maxLength int) string {

	// match a carriage return or newline character and use that as a delimiter
	// https://regex101.com/r/gb6pcj/2
	var NewlineRegex = regexp.MustCompile(`\r?\n`)

	str = NewlineRegex.ReplaceAllString(str, " ")
	str = strings.TrimSpace(str)
	if len(str) > maxLength {
		str = fmt.Sprintf("%s...", str[0:maxLength])
	}
	return str
}
