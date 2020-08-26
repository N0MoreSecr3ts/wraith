// Package matching contains specific functionality elated to scanning and detecting secrets within the given input.
package core

import (
	"crypto/sha1"
	"fmt"
	"io"
)

// Finding is a secret that has been discovered within a target by a discovery method
type Finding struct {
	Action            string
	Comment           string
	CommitAuthor      string
	CommitHash        string
	CommitMessage     string
	CommitUrl         string
	Description       string
	FilePath          string
	FileUrl           string
	WraithVersion     string
	Hash              string
	LineNumber        string
	RepositoryName    string
	RepositoryOwner   string
	RepositoryUrl     string
	Signatureid       string
	SignaturesVersion string
	SecretID          string
}

// setupUrls will set the urls used to search through either github or gitlab for inclusion in the finding data
func (f *Finding) setupUrls(scanType string, enterpriseUrl string) {
	baseUrl := ""
	if len(enterpriseUrl) > 0 {
		baseUrl = enterpriseUrl
	} else if scanType == "github" {
		baseUrl = "https://github.com"
	} else {
		baseUrl = "https://gitlab.com"
	}
	switch scanType {
	case "github":
		f.RepositoryUrl = fmt.Sprintf("%s/%s/%s", baseUrl, f.RepositoryOwner, f.RepositoryName)
		f.FileUrl = fmt.Sprintf("%s/blob/%s/%s", f.RepositoryUrl, f.CommitHash, f.FilePath)
		f.CommitUrl = fmt.Sprintf("%s/commit/%s", f.RepositoryUrl, f.CommitHash)
	case "gitlab":
		results := CleanUrlSpaces(f.RepositoryOwner, f.RepositoryName)
		f.RepositoryUrl = fmt.Sprintf("%s/%s/%s", baseUrl, results[0], results[1])
		f.FileUrl = fmt.Sprintf("%s/blob/%s/%s", f.RepositoryUrl, f.CommitHash, f.FilePath)
		f.CommitUrl = fmt.Sprintf("%s/commit/%s", f.RepositoryUrl, f.CommitHash)
	}

}

// generateID will create an ID for each finding based up the SHA1 of discrete data points associated
// with the finding
func (f *Finding) generateID() {
	h := sha1.New()
	_, err := io.WriteString(h, f.FilePath)
	_, err = io.WriteString(h, f.Action)
	_, err = io.WriteString(h, f.RepositoryOwner)
	_, err = io.WriteString(h, f.RepositoryName)
	_, err = io.WriteString(h, f.CommitHash)
	_, err = io.WriteString(h, f.CommitMessage)
	_, err = io.WriteString(h, f.CommitAuthor)

	if err != nil {
		fmt.Println("Not able to generate finding ID: ", err)
	}
	f.SecretID = fmt.Sprintf("%x", h.Sum(nil))
}

// Initialize will set the urls and create an ID for inclusion within the finding
func (f *Finding) Initialize(scanType string, enterpriseUrl string) {
	f.setupUrls(scanType, enterpriseUrl)
	f.generateID()
}
