// Package matching contains specific functionality elated to scanning and detecting secrets within the given input.
package core

import (
	"path/filepath"
	"strings"
)

// MatchTarget holds the various parts of a file that will be matched using either regex's or simple pattern matches.
type MatchTarget struct {
	Path      string
	Filename  string
	Extension string
	Content   string
}

// IsSkippable will check the matched file against a list of extensions or paths either
// supplied by the user or set by default
func (f *MatchTarget) IsSkippable(paths []string, exts []string) bool {
	ext := strings.ToLower(f.Extension)
	path := strings.ToLower(f.Path)
	for _, skippableExt := range exts {
		if ext == skippableExt {
			return true
		}
	}
	for _, skippablePath := range paths {
		if strings.Contains(path, skippablePath) {
			return true
		}
	}
	return false
}

// NewMatchTarget splits a filename into its composite pieces so that it may be measured
// and classified for scanning
func NewMatchTarget(path string) MatchTarget {
	_, filename := filepath.Split(path)
	extension := filepath.Ext(path)
	return MatchTarget{
		Path:      path,
		Filename:  filename,
		Extension: extension,
		Content:   "",
	}
}
