// Package common contains functionality not critical to the core project but still essential.
package core

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"strings"
	"syscall"
)

// pathExists will check if a path exists or not and is used to validate user input
func PathExists(path string) bool {
	_, err := os.Stat(path)

	if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOSPC {
		return false
	}

	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		//fmt.Println(err) // TODO remove me
		return false
	}

	return true
}

// TODO refactor out the common package

// FileExists will check for the existence of a file and return a bool depending
// on if it exists in a given path or not.
func FileExists(path string) bool {
	// TODO catch the error
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// AppendIfMissing will check a slice for a value before appending it
func AppendIfMissing(slice []string, s string) []string {
	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}
	return append(slice, s)
}

// SetHomeDir will set the correct homedir.
func SetHomeDir(h string) string {

	if strings.Contains(h, "$HOME") {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		h = strings.Replace(h, "$HOME", home, -1)
	}

	if strings.Contains(h, "~/") {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		h = strings.Replace(h, "~/", home, -1)
	}
	return h
}
