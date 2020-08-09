// Package matching contains specific functionality elated to scanning and detecting secrets within the given input.
package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

// Signatures holds a list of all signatures used during the session
type Signatures struct {
	FileSignatures    []FileSignature
	ContentSignatures []ContentSignature
}

// loadSignatures will check for a signature file to exist in a known location and load it if it is found.
func (s *Signatures) loadSignatures(path string) error {
	if !FileExists(path) {
		return errors.New(fmt.Sprintf("Missing signature file: %s.\n", path))
	}
	data, readError := ioutil.ReadFile(path)
	if readError != nil {
		return readError
	}
	if unmarshalError := json.Unmarshal(data, &s); unmarshalError != nil {
		return unmarshalError
	}
	return nil
}

// Load will load all known signatures for the various match types into the session
func (s *Signatures) Load(mode int) error {
	var e error
	if mode != 3 {
		e = s.loadSignatures("./rules/filesignatures.json")
		if e != nil {
			return e
		}
	}
	if mode != 1 {
		//source:  https://github.com/dxa4481/truffleHogRegexes/blob/master/truffleHogRegexes/regexes.json
		e = s.loadSignatures("./rules/contentsignatures.json")
		if e != nil {
			return e
		}
	}
	return nil
}
