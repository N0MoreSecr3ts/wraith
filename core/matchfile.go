package core

import (
	//"gopkg.in/src-d/go-git.v4/plumbing/object"
	"path/filepath"
	//"strconv"
	"strings"
	//"fmt"
	//"github.com/N0MoreSecr3ts/wraith/version"
)

// MatchFile holds the various parts of a file that will be matched using either regex's or simple pattern matches.
type MatchFile struct {
	Path      string
	Filename  string
	Extension string
}

// newMatchFile will generate a match object by dissecting a filename
func newMatchFile(path string) MatchFile {
	_, filename := filepath.Split(path)
	extension := filepath.Ext(path)
	return MatchFile{
		Path:      path,
		Filename:  filename,
		Extension: extension,
	}
}

// isSkippable will check the matched file against a list of extensions or paths either supplied by the user or set by default
func (f *MatchFile) isSkippable(sess *Session) bool {
	ext := strings.ToLower(f.Extension)
	path := strings.ToLower(f.Path)
	for _, skippableExt := range sess.SkippableExt {
		if ext == skippableExt {
			return true
		}
	}
	for _, skippablePath := range sess.SkippablePath {
		if strings.Contains(path, skippablePath) {
			return true
		}
	}
	return false
}
