// Package matching contains specific functionality elated to scanning and detecting secrets within the given input.
package core

import (
	"crypto/sha1"
	"fmt"
	"io"
)

// Finding is a secret that has been discovered within a target by a discovery method
type Finding struct {
	Id                          string
	FilePath                    string
	Action                      string
	FileSignatureDescription    string
	FileSignatureComment        string
	ContentSignatureDescription string
	ContentSignatureComment     string
	RepositoryOwner             string
	RepositoryName              string
	CommitHash                  string
	CommitMessage               string
	CommitAuthor                string
	FileUrl                     string
	CommitUrl                   string
	RepositoryUrl               string
	CloneUrl                    string
}

// setupUrls will set the urls used to search through either github or gitlab for inclusion in the finding data
func (f *Finding) setupUrls(scanType string) {
	switch scanType {
	case "github":
		f.RepositoryUrl = fmt.Sprintf("https://github.com/%s/%s", f.RepositoryOwner, f.RepositoryName)
		f.FileUrl = fmt.Sprintf("%s/blob/%s/%s", f.RepositoryUrl, f.CommitHash, f.FilePath)
		f.CommitUrl = fmt.Sprintf("%s/commit/%s", f.RepositoryUrl, f.CommitHash)
	case "gitlab":
		results := CleanUrlSpaces(f.RepositoryOwner, f.RepositoryName)
		f.RepositoryUrl = fmt.Sprintf("https://gitlab.com/%s/%s", results[0], results[1])
		f.FileUrl = fmt.Sprintf("%s/blob/%s/%s", f.RepositoryUrl, f.CommitHash, f.FilePath)
		f.CommitUrl = fmt.Sprintf("%s/commit/%s", f.RepositoryUrl, f.CommitHash)
	}

}

// generateID will create an ID for each finding based up the SHA1 of discrete data points associated
// with the finding
func (f *Finding) generateID() {
	h := sha1.New()
	io.WriteString(h, f.FilePath)
	io.WriteString(h, f.Action)
	io.WriteString(h, f.RepositoryOwner)
	io.WriteString(h, f.RepositoryName)
	io.WriteString(h, f.CommitHash)
	io.WriteString(h, f.CommitMessage)
	io.WriteString(h, f.CommitAuthor)
	f.Id = fmt.Sprintf("%x", h.Sum(nil))
}

// Initialize will set the urls and create an ID for inclusion within the finding
func (f *Finding) Initialize(scanType string) {
	f.setupUrls(scanType)
	f.generateID()
}
