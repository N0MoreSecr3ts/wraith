// Package core represents the core functionality of all commands
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
	Action           string
	Content          string
	CommitAuthor     string
	CommitHash       string
	CommitMessage    string
	CommitURL        string
	Description      string
	FilePath         string
	FileURL          string
	WraithVersion    string
	Hash             string
	LineNumber       string
	RepositoryName   string
	RepositoryOwner  string
	RepositoryURL    string
	SignatureID      string
	signatureVersion string
	SecretID         string
}

// setupUrls will set the urls used to search through either github or gitlab for inclusion in the finding data
func (f *Finding) setupUrls(sess *Session) {
	baseURL := ""
	if sess.ScanType == "github-enterprise" {
		baseURL = sess.GithubEnterpriseURL
	} else if sess.ScanType == "github" {
		baseURL = "https://github.com"
	} else {
		baseURL = "https://gitlab.com"
	}
	switch sess.ScanType {
	case "github":
		f.RepositoryURL = fmt.Sprintf("%s/%s/%s", baseURL, f.RepositoryOwner, f.RepositoryName)
		f.FileURL = fmt.Sprintf("%s/blob/%s/%s", f.RepositoryURL, f.CommitHash, f.FilePath)
		f.CommitURL = fmt.Sprintf("%s/commit/%s", f.RepositoryURL, f.CommitHash)
	case "gitlab":
		results := CleanURLSpaces(f.RepositoryOwner, f.RepositoryName)
		f.RepositoryURL = fmt.Sprintf("%s/%s/%s", baseURL, results[0], results[1])
		f.FileURL = fmt.Sprintf("%s/blob/%s/%s", f.RepositoryURL, f.CommitHash, f.FilePath)
		f.CommitURL = fmt.Sprintf("%s/commit/%s", f.RepositoryURL, f.CommitHash)
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
