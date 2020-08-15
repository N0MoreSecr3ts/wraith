// Package matching contains specific functionality elated to scanning and detecting secrets within the given input.
package core

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// this are the various items that we are attempting to match against using either regex's or simple pattern matches.
const (
	PartExtension = "extension" // file extension
	PartFilename  = "filename"  // file name
	PartPath      = "path"      // the path to the file
	PartContent   = "content"   // the content of the file
)

// Signatures holds a list of all signatures used during the session
var Signatures []Signature

// SafeFunctionSignatures is a collection of safe function sigs
var SafeFunctionSignatures = []SafeFunctionSignature{}

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

// generateGenericID will return an id with sufficient enough entropy to be usable for larger scale sessionss
func generateGenericID(val1 string) string {
	id := sha1.New()

	str := val1

	io.WriteString(id, str)

	j := id.Sum(nil)

	encodedStr := hex.EncodeToString(j)

	return encodedStr
}

// Signature is an expression that we are looking for in a file
type Signature interface {
	Description() string
	Enable() int
	ExtractMatch(file MatchFile, sess *Session) (bool, map[string]int)
	MatchLevel() int
	Part() string
	Signatureid() string // TODO change id -> ID
}

// SignaturesMetaData is used by updateSignatures to determine if/how to update the signatures
type SignaturesMetaData struct {
	Date    string
	Time    int
	Version string
}

// SafeFunctionSignature holds the information about a safe function, that is used to detect and mitigate false positives
type SafeFunctionSignature struct {
	comment     string
	description string
	enable      int
	entropy     float64
	match       *regexp.Regexp
	matchLevel  int
	part        string
	signatureid string
}

// SimpleSignature holds the information about a simple signature which is used to match a path or filename
type SimpleSignature struct {
	comment     string
	description string
	enable      int
	entropy     float64
	match       string
	matchLevel  int
	part        string
	signatureid string
}

// PatternSignature holds the information about a pattern signature which is a regex used to match content within a file
type PatternSignature struct {
	comment     string
	description string
	enable      int
	entropy     float64
	match       *regexp.Regexp
	matchLevel  int
	part        string
	signatureid string
}

// SignatureDef maps to a signature within the yaml file
type SignatureDef struct {
	Comment     string  `yaml:"comment"`
	Description string  `yaml:"description"`
	Enable      int     `yaml:"enable"`
	Entropy     float64 `yaml:"entropy"`
	Match       string  `yaml:"match"`
	MatchLevel  int     `yaml:"match-level"`
	Part        string  `yaml:"part"`
	Signatureid string  `yaml:"signatureid"`
}

// SignatureConfig holds the base file structure for the signatures file
type SignatureConfig struct {
	Meta                   SignaturesMetaData `yaml:"Meta"`
	PatternSignatures      []SignatureDef     `yaml:"PatternSignatures"`
	SimpleSignatures       []SignatureDef     `yaml:"SimpleSignatures"`
	SafeFunctionSignatures []SignatureDef     `yaml:"SafeFunctionSignatures"`
}

// ExtractMatch will attempt to match a path or file name of the given file
func (s SimpleSignature) ExtractMatch(file MatchFile, sess *Session) (bool, map[string]int) {
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

// MatchLevel sets the confidence level of the pattern
func (s SimpleSignature) MatchLevel() int {
	return s.matchLevel
}

// Part sets the part of the file/path that is matched [ filename content extension ]
func (s SimpleSignature) Part() string {
	return s.part
}

// Description sets the user comment of the signature
func (s SimpleSignature) Description() string {
	return s.description
}

// Sugnatureid sets the id used to identify the signature. This id is immutable and generated from a has of the signature and is changed with every update to a signature.
func (s SimpleSignature) Signatureid() string {
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
func (s PatternSignature) ExtractMatch(file MatchFile, sess *Session) (bool, map[string]int) {

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
		//fmt.Println("I am here at the content") // TODO remove me
		haystack := &file.Path
		fmt.Println(&file.Path)
		fmt.Println("extracting a content match")
		if PathExists(*haystack, sess) {
			if _, err := os.Stat(*haystack); err == nil {
				data, err := ioutil.ReadFile(*haystack)
				if err != nil {
					sErrAppend := fmt.Sprintf("ERROR --- Unable to open file for scanning: <%s> \nError Message: <%s>", *haystack, err)
					results[sErrAppend] = 0 // set to zero due to error, we never have a line 0 so we can always ignore that or error on it
					return false, results
				}

				r := s.match // this is the regex that we are going to try and match against

				var contextMatches []string
				if r.Match(data) {
					for _, curRegexMatch := range r.FindAll(data, -1) {
						contextMatches = append(contextMatches, string(curRegexMatch))
					}
				}

				if len(contextMatches) > 0 {
					bResult = true
					for i, curMatch := range contextMatches {
						thisMatch := string(curMatch[:])
						thisMatch = strings.TrimSuffix(thisMatch, "\n")

						bResult = confirmEntropy(thisMatch, s.entropy)

						if bResult {
							linesOfScannedFile := strings.Split(string(data), "\n")
							//linesOfScannedFile = linesOfScannedFile[:len(linesOfScannedFile)] // TODO Is this needed?

							num := fetchLineNumber(&linesOfScannedFile, thisMatch, i)
							results[strconv.Itoa(i)+"_"+thisMatch] = num
						}
					}
					return bResult, results
				}
			}
		}
	default:
		return bResult, results
	}
	return bResult, results
}

// fetchLineNumber will read a file in line by line and when the match is found, save the line number. It manages multiple matches in a file by way of the count and an index
func fetchLineNumber(input *[]string, thisMatch string, idx int) int {
	linesOfScannedFile := *input
	lineNumIndexMap := make(map[int]int)

	count := 0

	for i, line := range linesOfScannedFile {
		if strings.Contains(line, thisMatch) {

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

// MatchLevel sets the confidence level of the pattern
func (s PatternSignature) MatchLevel() int {
	return s.matchLevel
}

// Part sets the part of the file/path that is matched [ filename content extension ]
func (s PatternSignature) Part() string {
	return s.part
}

// Description sets the user comment of the signature
func (s PatternSignature) Description() string {
	return s.description
}

// Signatureid sets the id used to identify the signature. This id is immutable and generated from a has of the signature and is changed with every update to a signature.
func (s PatternSignature) Signatureid() string {
	return s.signatureid
}

// Enable sets whether as signature is active or not
func (s SafeFunctionSignature) Enable() int {
	return s.enable
}

// MatchLevel sets the confidence level of the pattern
func (s SafeFunctionSignature) MatchLevel() int {
	return s.matchLevel
}

// Part sets the part of the file/path that is matched [ filename content extension ]
func (s SafeFunctionSignature) Part() string {
	return s.part
}

// Description sets the user comment of the signature
func (s SafeFunctionSignature) Description() string {
	return s.description
}

// Signatureid sets the id used to identify the signature. This id is immutable and generated from a has of the signature and is changed with every update to a signature.
func (s SafeFunctionSignature) Signatureid() string {
	return s.signatureid
}

// ExtractMatch is a placeholder to ensure min code complexity and allow the reuse of the functions
func (s SafeFunctionSignature) ExtractMatch(file MatchFile, sess *Session) (bool, map[string]int) {
	var results map[string]int

	return false, results
}

// LoadSignatures will load all known signatures for the various match types into the session
func LoadSignatures(filePath string, mLevel int, sess *Session) []Signature { // TODO we don't need to bring in session here

	// ensure that we have the proper home directory
	filePath = SetHomeDir(filePath)

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

	SimpleSignatures := []SimpleSignature{}   // TODO change this variable name
	PatternSignatures := []PatternSignature{} // TODO change this variable name
	for _, curSig := range c.SimpleSignatures {

		if curSig.Enable > 0 && curSig.MatchLevel >= mLevel {

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
				curSig.MatchLevel,
				part,
				curSig.Signatureid,
			})
		}
	}

	for _, curSig := range c.PatternSignatures {
		if curSig.Enable > 0 && curSig.MatchLevel >= mLevel {
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
				curSig.MatchLevel,
				part,
				curSig.Signatureid,
			})
		}
	}
	for _, curSig := range c.SafeFunctionSignatures {
		if curSig.Enable > 0 && curSig.MatchLevel >= mLevel {
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
				curSig.MatchLevel,
				part,
				curSig.Signatureid,
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
