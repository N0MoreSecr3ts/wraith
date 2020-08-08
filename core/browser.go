// Package common contains functionality not critical to the core project but still essential.
package core

// TODO refactor out the common package

import (
	"fmt"
	"strings"
)

// UserAgent set the browser user agent when required.
var UserAgent = fmt.Sprintf("%s v%s", Name, Version)

// CleanUrlSpaces will take a string and replace any spaces with dashes so that is may be used in a url.
func CleanUrlSpaces(dirtyStrings ...string) []string {
	var result []string
	for _, s := range dirtyStrings {
		result = append(result, strings.ReplaceAll(s, " ", "-"))
	}
	return result
}
