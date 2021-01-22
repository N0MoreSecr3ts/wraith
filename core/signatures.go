package core

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// These are the various items that we are attempting to match against using either regex's or simple pattern matches.
const (
	PartExtension = "extension" // file extension
	PartFilename  = "filename"  // file name
	PartPath      = "path"      // the path to the file
	PartContent   = "content"   // the content of the file
)

// Signatures holds a list of all signatures used during the session
var Signatures []Signature

// SafeFunctionSignatures is a collection of safe function sigs
var SafeFunctionSignatures []SafeFunctionSignature

// loadSignatureSet will read in the defined signatures from an external source
func loadSignatureSet(filename string) (SignatureConfig, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return SignatureConfig{}, err
	}

	var c SignatureConfig
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return SignatureConfig{}, err
	}

	return c, nil
}

// get EntropyInt will calculate the entrophy based upon Shannon Entropy
func getEntropyInt(s string) float64 {
	//Shannon Entropy calculation
	m := map[rune]float64{}
	for _, r := range s {
		m[r]++
	}
	var hm float64
	for _, c := range m {
		hm += c * math.Log2(c)
	}
	l := float64(len(s))
	res := math.Log2(l) - hm/l
	return res
}

// Signature is an expression that we are looking for in a file
type Signature interface {
	Description() string
	Enable() int
	ExtractMatch(file MatchFile, sess *Session, change *object.Change) (bool, map[string]int)
	ConfidenceLevel() int
	Part() string
	SignatureID() string // TODO change id -> ID
}

// SignaturesMetaData is used by updateSignatures to determine if/how to update the signatures
type SignaturesMetaData struct {
	Date    string
	Time    int
	Version string
}

// SafeFunctionSignature holds the information about a safe function, that is used to detect and mitigate false positives
type SafeFunctionSignature struct {
	comment         string
	description     string
	enable          int
	entropy         float64
	match           *regexp.Regexp
	confidenceLevel int
	part            string
	signatureid     string
}

// SimpleSignature holds the information about a simple signature which is used to match a path or filename
type SimpleSignature struct {
	comment         string
	description     string
	enable          int
	entropy         float64
	match           string
	confidenceLevel int
	part            string
	signatureid     string
}

// PatternSignature holds the information about a pattern signature which is a regex used to match content within a file
type PatternSignature struct {
	comment         string
	description     string
	enable          int
	entropy         float64
	match           *regexp.Regexp
	confidenceLevel int
	part            string
	signatureid     string
}

// SignatureDef maps to a signature within the yaml file
type SignatureDef struct {
	Comment         string  `yaml:"comment"`
	Description     string  `yaml:"description"`
	Enable          int     `yaml:"enable"`
	Entropy         float64 `yaml:"entropy"`
	Match           string  `yaml:"match"`
	ConfidenceLevel int     `yaml:"confidence-level"`
	Part            string  `yaml:"part"`
	SignatureID     string  `yaml:"signatureid"`
}

// SignatureConfig holds the base file structure for the signatures file
type SignatureConfig struct {
	Meta                   SignaturesMetaData `yaml:"Meta"`
	PatternSignatures      []SignatureDef     `yaml:"PatternSignatures"`
	SimpleSignatures       []SignatureDef     `yaml:"SimpleSignatures"`
	SafeFunctionSignatures []SignatureDef     `yaml:"SafeFunctionSignatures"`
}

// ExtractMatch will attempt to match a path or file name of the given file
func (s SimpleSignature) ExtractMatch(file MatchFile, sess *Session, change *object.Change) (bool, map[string]int) {
	var haystack *string
	var bResult = false

	// this is empty and could be removed but it here to streamline all the match functions
	var results map[string]int

	switch s.part {
	case PartPath:
		haystack = &file.Path
		bResult = true
	case PartFilename:
		haystack = &file.Filename
		bResult = true
	case PartExtension:
		haystack = &file.Extension
		bResult = true
	default:
		return bResult, results
	}

	return s.match == *haystack, results
}

// Enable sets whether as signature is active or not
func (s SimpleSignature) Enable() int {
	return s.enable
}

// ConfidenceLevel sets the confidence level of the pattern
func (s SimpleSignature) ConfidenceLevel() int {
	return s.confidenceLevel
}

// Part sets the part of the file/path that is matched [ filename content extension ]
func (s SimpleSignature) Part() string {
	return s.part
}

// Description sets the user comment of the signature
func (s SimpleSignature) Description() string {
	return s.description
}

// SignatureID sets the id used to identify the signature. This id is immutable and generated from a has of the signature and is changed with every update to a signature.
func (s SimpleSignature) SignatureID() string {
	return s.signatureid
}

// IsSafeText check against known "safe" (aka not a password) list
func IsSafeText(sMatchString *string) bool {
	bResult := false
	for _, safeSig := range SafeFunctionSignatures {
		if safeSig.match.MatchString(*sMatchString) {
			bResult = true
		}
	}
	return bResult
}

// confirmEntropy will determine correct entrophy of the string and decide if we move forward with the match
func confirmEntropy(thisMatch string, iSessionEntropy float64) bool {
	bResult := false

	iEntropy := getEntropyInt(thisMatch)

	if (iSessionEntropy == 0) || (iEntropy >= iSessionEntropy) {
		if !IsSafeText(&thisMatch) {
			bResult = true
		}
	}

	return bResult
}

// ExtractMatch will try and find a match within the content of the file.
func (s PatternSignature) ExtractMatch(file MatchFile, sess *Session, change *object.Change) (bool, map[string]int) {

	var haystack *string            // this is a pointer to the item we want to match
	var bResult = false             // match result
	results := make(map[string]int) // the secret and the line number in a map

	switch s.part {
	case PartPath:
		haystack = &file.Path
		bResult = s.match.MatchString(*haystack)
	case PartFilename:
		haystack = &file.Filename
		bResult = s.match.MatchString(*haystack)
	case PartExtension:
		haystack = &file.Extension
		bResult = s.match.MatchString(*haystack)
	case PartContent:
		haystack := &file.Path
		if PathExists(*haystack, sess) {
			if _, err := os.Stat(*haystack); err == nil {
				data, err := ioutil.ReadFile(*haystack)
				if err != nil {
					sErrAppend := fmt.Sprintf("ERROR --- Unable to open file for scanning: <%s> \nError Message: <%s>", *haystack, err)
					results[sErrAppend] = 0 // set to zero due to error, we never have a line 0 so we can always ignore that or error on it
					return false, results
				}

				// The regex that we are going to try and match against
				r := s.match

				var contextMatches []string

				// Check to see if there is a match in the data and if so switch to a Findall that
				// will get a slice of all the individual matches. Doing this ahead of time saves us
				// from looping through if it is not necessary.
				if r.Match(data) {
					for _, curRegexMatch := range r.FindAll(data, -1) {
						contextMatches = append(contextMatches, string(curRegexMatch))
					}
					if len(contextMatches) > 0 {
						bResult = true
						for i, curMatch := range contextMatches {

							thisMatch := string(curMatch[:])
							thisMatch = strings.TrimSuffix(thisMatch, "\n")

							bResult = confirmEntropy(thisMatch, s.entropy)

							if bResult {
								linesOfScannedFile := strings.Split(string(data), "\n")

								num := fetchLineNumber(&linesOfScannedFile, thisMatch, 0)
								results[strconv.Itoa(i)+"_"+thisMatch] = num
							}
						}
						return bResult, results
					}
				}

				if sess.ScanType != "localPath" {

					content, err := GetChangeContent(change)
					if err != nil {
						sess.Out.Error("Error retrieving content in commit %s, change %s:  %s\n", "commit.String()", change.String(), err)
					}

					if r.Match([]byte(content)) {
						for _, curRegexMatch := range r.FindAll([]byte(content), -1) {
							contextMatches = append(contextMatches, string(curRegexMatch))
						}
						if len(contextMatches) > 0 {
							bResult = true
							for i, curMatch := range contextMatches {
								thisMatch := string(curMatch[:])
								thisMatch = strings.TrimSuffix(thisMatch, "\n")

								bResult = confirmEntropy(thisMatch, s.entropy)

								if bResult {
									linesOfScannedFile := strings.Split(content, "\n")

									num := fetchLineNumber(&linesOfScannedFile, thisMatch, i)
									results[strconv.Itoa(i)+"_"+thisMatch] = num
								}
							}
							return bResult, results
						}
					}
				}

			}
		}
	default: // TODO We need to do something with this
		return bResult, results
	}
	return bResult, results

}

// fetchLineNumber will read a file line by line and when the match is found, save the line number.
// It manages multiple matches in a file by way of the count and an index
func fetchLineNumber(input *[]string, thisMatch string, idx int) int {
	linesOfScannedFile := *input
	lineNumIndexMap := make(map[int]int)

	count := 0

	for i, line := range linesOfScannedFile {
		if strings.Contains(line, thisMatch) {

			// We need to add 1 here as the index starts at zero so every line number would be line -1 normally
			lineNumIndexMap[count] = i + 1
			count = count + 1
		}
	}
	return lineNumIndexMap[idx]
}

// Enable sets whether as signature is active or not
func (s PatternSignature) Enable() int {
	return s.enable
}

// ConfidenceLevel sets the confidence level of the pattern
func (s PatternSignature) ConfidenceLevel() int {
	return s.confidenceLevel
}

// Part sets the part of the file/path that is matched [ filename content extension ]
func (s PatternSignature) Part() string {
	return s.part
}

// Description sets the user comment of the signature
func (s PatternSignature) Description() string {
	return s.description
}

// SignatureID sets the id used to identify the signature. This id is immutable and generated from a has of the signature and is changed with every update to a signature.
func (s PatternSignature) SignatureID() string {
	return s.signatureid
}

// Enable sets whether as signature is active or not
func (s SafeFunctionSignature) Enable() int {
	return s.enable
}

// ConfidenceLevel sets the confidence level of the pattern
func (s SafeFunctionSignature) ConfidenceLevel() int {
	return s.confidenceLevel
}

// Part sets the part of the file/path that is matched [ filename content extension ]
func (s SafeFunctionSignature) Part() string {
	return s.part
}

// Description sets the user comment of the signature
func (s SafeFunctionSignature) Description() string {
	return s.description
}

// SignatureID sets the id used to identify the signature. This id is immutable and generated from a has of the signature and is changed with every update to a signature.
func (s SafeFunctionSignature) SignatureID() string {
	return s.signatureid
}

// ExtractMatch is a placeholder to ensure min code complexity and allow the reuse of the functions
func (s SafeFunctionSignature) ExtractMatch(file MatchFile, sess *Session, change *object.Change) (bool, map[string]int) {
	var results map[string]int

	return false, results
}

// LoadSignatures will load all known signatures for the various match types into the session
func LoadSignatures(filePath string, mLevel int, sess *Session) []Signature { // TODO we don't need to bring in session here

	// ensure that we have the proper home directory
	filePath = SetHomeDir(filePath, sess)

	c, err := loadSignatureSet(filePath)
	if err != nil {
		sess.Out.Error("Failed to load signatures file %s: %s\n", filePath, err.Error())
		os.Exit(2)
	}

	signaturesMetaData := SignaturesMetaData{
		Version: c.Meta.Version,
		Date:    c.Meta.Date,
		Time:    c.Meta.Time,
	}

	sess.SignatureVersion = signaturesMetaData.Version

	var SimpleSignatures []SimpleSignature
	var PatternSignatures []PatternSignature
	for _, curSig := range c.SimpleSignatures {

		if curSig.Enable > 0 && curSig.ConfidenceLevel >= mLevel {

			var part string
			switch strings.ToLower(curSig.Part) {
			case "partpath":
				part = PartPath
			case "partfilename":
				part = PartFilename
			case "partextension":
				part = PartExtension
			case "partcontent":
				part = PartContent
			default:
				part = PartContent
			}

			SimpleSignatures = append(SimpleSignatures, SimpleSignature{
				curSig.Comment,
				curSig.Description,
				curSig.Enable,
				curSig.Entropy,
				curSig.Match,
				curSig.ConfidenceLevel,
				part,
				curSig.SignatureID,
			})
		}
	}

	for _, curSig := range c.PatternSignatures {
		if curSig.Enable > 0 && curSig.ConfidenceLevel >= mLevel {
			var part string
			switch strings.ToLower(curSig.Part) {
			case "partpath":
				part = PartPath
			case "partfilename":
				part = PartFilename
			case "partextension":
				part = PartExtension
			case "partcontent":
				part = PartContent
			default:
				part = PartContent
			}

			match := regexp.MustCompile(curSig.Match)
			PatternSignatures = append(PatternSignatures, PatternSignature{
				curSig.Comment,
				curSig.Description,
				curSig.Enable,
				curSig.Entropy,
				match,
				curSig.ConfidenceLevel,
				part,
				curSig.SignatureID,
			})
		}
	}
	for _, curSig := range c.SafeFunctionSignatures {
		if curSig.Enable > 0 && curSig.ConfidenceLevel >= mLevel {
			var part string
			switch strings.ToLower(curSig.Part) {
			case "partpath":
				part = PartPath
			case "partfilename":
				part = PartFilename
			case "partextension":
				part = PartExtension
			case "partcontent":
				part = PartContent
			default:
				part = PartContent
			}

			match := regexp.MustCompile(curSig.Match)
			SafeFunctionSignatures = append(SafeFunctionSignatures, SafeFunctionSignature{
				curSig.Comment,
				curSig.Description,
				curSig.Enable,
				curSig.Entropy,
				match,
				curSig.ConfidenceLevel,
				part,
				curSig.SignatureID,
			})
		}
	}

	idx := len(PatternSignatures) + len(SimpleSignatures)

	Signatures := make([]Signature, idx)
	jdx := 0
	for _, v := range SimpleSignatures {
		Signatures[jdx] = v
		jdx++
	}

	for _, v := range PatternSignatures {
		Signatures[jdx] = v
		jdx++
	}

	// TODO are we loading the safe ones somewhere

	return Signatures
}
