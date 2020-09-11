// Package matching contains specific functionality elated to scanning and detecting secrets within the given input.
package core

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"time"
)

// Finding is a secret that has been discovered within a target by a discovery method
type Finding struct {
	Action            string
	Content           string
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
func (f *Finding) setupUrls(sess *Session) {
	baseUrl := ""
	if len(sess.EnterpriseURL) > 0 && sess.EnterpriseScan {
		baseUrl = sess.EnterpriseURL
	} else if sess.ScanType == "github" {
		baseUrl = "https://github.com"
	} else {
		baseUrl = "https://gitlab.com"
	}
	switch sess.ScanType {
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
// with the finding.
func generateID() string {
	h := sha1.New()
	source := rand.NewSource(time.Now().UnixNano())
	randNum := rand.New(source)

	_, err := io.WriteString(h, fmt.Sprintf("%x", randNum.Intn(10000000000)))

	if err != nil {
		fmt.Println("Not able to generate finding ID: ", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Initialize will set the urls and create an ID for inclusion within the finding
func (f *Finding) Initialize(sess *Session) {
	f.setupUrls(sess)
}
