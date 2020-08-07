// Package matching contains specific functionality elated to scanning and detecting secrets within the given input.
package matching

import "regexp"

type ContentSignature struct {
	MatchOn     string
	Description string
	Comment     string
}

// Match will attempt a match based on the content of the file. This is used to match items like
// access tokens or passwords within a given file.
func (c ContentSignature) Match(target MatchTarget) (bool, error) {
	return regexp.MatchString(c.MatchOn, target.Content)
}

// GetDescription will return the description of the signature
func (c ContentSignature) GetDescription() string {
	return c.Description
}

// GetComment will return the comment of the signature
func (c ContentSignature) GetComment() string {
	return c.Comment
}
