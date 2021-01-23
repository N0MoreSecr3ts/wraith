package core

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

// PathExists will check if a path exists or not and is used to validate user input
func PathExists(path string, sess *Session) bool {
	_, err := os.Stat(path)

	if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOSPC {
		return false
	}

	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		sess.Out.Debug("Path does not exist: %s\n", err.Error())
		return false
	}

	return true
}

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
func SetHomeDir(h string, sess *Session) string {

	if strings.Contains(h, "$HOME") {
		home, err := homedir.Dir()
		if err != nil {
			sess.Out.Error(err.Error())
			os.Exit(2)
		}

		h = strings.Replace(h, "$HOME", home, -1)
	}

	if strings.Contains(h, "~") {
		home, err := homedir.Dir()
		if err != nil {
			sess.Out.Error(err.Error())
			os.Exit(2)
		}
		h = strings.Replace(h, "~", home, -1)
	}
	return h
}

// realTimeOutput will print out the current finding to stdout if all conditions are met
func realTimeOutput(finding *Finding, sess *Session) {
	if !sess.Silent && !sess.CSVOutput && !sess.JSONOutput {

		sess.Out.Warn(" %s\n", strings.ToUpper(finding.Description))
		sess.Out.Info("  SignatureID..........: %s\n", finding.SignatureID)
		sess.Out.Info("  Repo.................: %s\n", finding.RepositoryName)
		sess.Out.Info("  File Path............: %s\n", finding.FilePath)
		sess.Out.Info("  Line Number..........: %s\n", finding.LineNumber)
		sess.Out.Info("  Message..............: %s\n", TruncateString(finding.CommitMessage, 100))
		sess.Out.Info("  Commit Hash..........: %s\n", TruncateString(finding.CommitHash, 100))
		sess.Out.Info("  Author...............: %s\n", finding.CommitAuthor)
		sess.Out.Info("  SecretID.............: %v\n", finding.SecretID)
		sess.Out.Info("  Wraith Version.......: %s\n", finding.WraithVersion)
		sess.Out.Info("  Signatures Version...: %v\n", finding.signatureVersion)
		if len(finding.Content) > 0 {
			issues := "\n\t" + finding.Content
			sess.Out.Info("  Issues..........: %s\n", issues)
		}

		sess.Out.Info(" ------------------------------------------------\n\n")
	}
}

// IsMaxFileSize will determine if the file size is under the max limit set by maxFileSize
func IsMaxFileSize(filename string, sess *Session) (bool, string) {

	fi, err := os.Stat(filename)

	if err != nil {
		return true, "does not exist"
	}

	fileSize := fi.Size()
	var maxFileSize int64
	maxFileSize = sess.MaxFileSize * 1024 * 1024

	if fileSize > maxFileSize {
		return true, "is too large"
	}
	return false, ""
}

// isTestFileOrPath will run various regex's against a target to determine if it is a test file or contained in a test directory.
func isTestFileOrPath(fullPath string) bool {
	fName := filepath.Base(fullPath)

	// If the directory contains "test"
	// Ex. foo/test/bar
	r := regexp.MustCompile(`(?i)[/\\]test?[/\\]`)
	if r.MatchString(fullPath) {
		return true
	}

	// If the directory starts with test, the leading slash gets dropped by default
	// Ex. test/foo/bar
	r = regexp.MustCompile(`(?i)test?[/\\]`)
	if r.MatchString(fullPath) {
		return true
	}

	// If the directory path starts with a different root but has the word test in it somewhere
	// Ex. foo/test-secrets/bar
	r = regexp.MustCompile(`/test.*/`)
	if r.MatchString(fullPath) {
		return true
	}

	// A the word Test is in the string, case sensitive
	// Ex. ghTestlk
	// Ex. Testllfhe
	// Ex. Test
	r = regexp.MustCompile(`Test`)
	if r.MatchString(fName) {
		return true
	}

	// A file has a suffix of _test
	// Golang uses this as the default test file naming convention
	//Ex. foo_test.go
	r = regexp.MustCompile(`(?i)_test`)
	if r.MatchString(fName) {
		return true
	}

	// If the pattern _test_ is in the string
	// Ex. foo_test_baz
	r = regexp.MustCompile(`(?i)_test?_`)
	if r.MatchString(fName) {
		return true
	}

	return false
}
