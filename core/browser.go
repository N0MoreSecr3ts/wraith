// Package core represents the core functionality of all commands
package core

import (
	"fmt"
	"strings"
)

// UserAgent set the browser user agent when required.
var UserAgent = fmt.Sprintf("%s v%s", Name, Version)

// CleanURLSpaces will take a string and replace any spaces with dashes so that is may be used in a url.
func CleanURLSpaces(dirtyStrings ...string) []string {
	var result []string
	for _, s := range dirtyStrings {
		result = append(result, strings.ReplaceAll(s, " ", "-"))
	}
	return result
}
