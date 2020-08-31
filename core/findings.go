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
	Fields            []string
}

// setupUrls will set the urls used to search through either github or gitlab for inclusion in the finding data
func (f *Finding) setupUrls(sess *Session) {
	baseUrl := ""
	if len(sess.EnterpriseURL) > 0 {
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

func (f *Finding) getFieldNames() []string {
	return f.Fields
}

func (f *Finding) getValues() []string {
	values := *new([]string)
	values = append(values, f.Action)
	values = append(values, f.Comment)
	values = append(values, f.CommitAuthor)
	values = append(values, f.CommitHash)
	values = append(values, f.CommitMessage)
	values = append(values, f.CommitUrl)
	values = append(values, f.Description)
	values = append(values, f.FilePath)
	values = append(values, f.FileUrl)
	values = append(values, f.WraithVersion)
	values = append(values, f.Hash)
	values = append(values, f.LineNumber)
	values = append(values, f.RepositoryName)
	values = append(values, f.RepositoryOwner)
	values = append(values, f.RepositoryUrl)
	values = append(values, f.Signatureid)
	values = append(values, f.SignaturesVersion)
	values = append(values, f.SecretID)
	return values
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
func (f *Finding) Initialize(sess *Session) {
	f.setupUrls(sess)
	f.generateID()
	f.Fields = []string{"Action","Comment","CommitAuthor","CommitHash","CommitMessage","CommitUrl","Description","FilePath",
		"FileUrl","WraithVersion","Hash","LineNumber","RepositoryName","RepositoryOwner","RepositoryUrl","Signatureid",
		"SignaturesVersion","SecretID"}
}
