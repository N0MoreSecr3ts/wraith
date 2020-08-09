//// Package matching contains specific functionality elated to scanning and detecting secrets within the given input.
package core

//
//import (
//	"errors"
//	"fmt"
//	"regexp"
//)
//
//// FileSignatureType is a breakdown of the various parts associated with a file
//type FileSignatureType struct {
//	Extension string
//	Filename  string
//	Path      string
//}
//
//// fileSignatureTypes stores the common names used to reference the types within the codebase
//var fileSignatureTypes = FileSignatureType{
//	Extension: "extension",
//	Filename:  "filename",
//	Path:      "path",
//}
//
//// FileSignature holds various values associated with a specific signature used to find a secret.
//type FileSignature struct {
//	Part        string
//	MatchOn     string
//	Description string
//	Comment     string
//}
//
//// Match will attempt to match the path or a given part of the file name. This is used to match specific files such as
//// private keys or specific token files.
//func (f FileSignature) Match(target MatchTarget) (bool, error) {
//	var haystack *string
//	switch f.Part {
//	case fileSignatureTypes.Path:
//		haystack = &target.Path
//	case fileSignatureTypes.Filename:
//		haystack = &target.Filename
//	case fileSignatureTypes.Extension:
//		haystack = &target.Extension
//	default:
//		return false, errors.New(fmt.Sprintf("Unrecognized 'Part' parameter: %s\n", f.Part))
//	}
//	return regexp.MatchString(f.MatchOn, *haystack)
//}
//
//// GetDescription will return the description of the signature
//func (f FileSignature) GetDescription() string {
//	return f.Description
//}
//
//// GetComment will return the comment of the signature
//func (f FileSignature) GetComment() string {
//	return f.Comment
//}
