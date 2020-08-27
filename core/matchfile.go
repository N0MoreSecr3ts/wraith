package core

import (
	//"gopkg.in/src-d/go-git.v4/plumbing/object"
	"path/filepath"
	//"strconv"
	"strings"
	//"fmt"
	//"wraith/version"
)

// MatchFile holds the various parts of a file that will be matched using either regex's or simple pattern matches.
type MatchFile struct {
	Path      string
	Filename  string
	Extension string
}

// newMatchFile will generate a match object by dissecting a filename
func newMatchFile(path string) MatchFile {
	_, filename := filepath.Split(path)
	extension := filepath.Ext(path)
	return MatchFile{
		Path:      path,
		Filename:  filename,
		Extension: extension,
	}
}

// isSkippable will check the matched file against a list of extensions or paths either supplied by the user or set by default
func (f *MatchFile) isSkippable(sess *Session) bool {
	ext := strings.ToLower(f.Extension)
	path := strings.ToLower(f.Path)
	for _, skippableExt := range sess.SkippableExt {
		if ext == skippableExt {
			return true
		}
	}
	for _, skippablePath := range sess.SkippablePath {
		if strings.Contains(path, skippablePath) {
			return true
		}
	}
	return false
}
//
//func findMatch(sess *Session, signature Signature, change *object.Change, file MatchFile )  Finding{
//	bMatched, matchMap := signature.ExtractMatch(file, sess, change)
//	if bMatched {
//
//		// Incremented the count of files that contain secrets
//		sess.Stats.IncrementFilesDirty()
//
//		var content string   // this is because file matches are puking
//		var genericID string // the generic id used in the finding
//
//		// For every instance of the secret that matched the specific signatures
//		// create a new finding. Thi will produce dupes as the file may exist
//		// in multiple commits.
//		for k, v := range matchMap {
//
//			// This sets the content for the finding, in this case the actual secret
//			// is the content. This can be removed and hidden via a commandline flag.
//			cleanK := strings.SplitAfterN(k, "_", 2)
//			if matchMap == nil {
//				content = ""
//				//genericID = *repo.Name + "://" + fPath + "_" + generateGenericID(content) // TODO remove me
//			} else {
//				content = cleanK[1]
//				//genericID = *repo.Name + "://" + fPath + "_" + generateGenericID(content) // TODO remove me
//
//			}
//			genericID = *repo.Name + "://" + fPath + "_" + generateGenericID(content)
//			fmt.Println("Content:") // TODO Remove me
//			fmt.Println(cleanK)     // TODO Remove me
//			fmt.Println()           // TODO Remove me
//
//			// Destroy the secret by zeroing the content if the flag is set
//			if sess.HideSecrets {
//				content = ""
//			}
//
//			// TODO Is this need for the finding object, why are we saving this?
//			changeAction := GetChangeAction(change)
//			// Create a new instance of a finding and set the necessary fields.
//			finding := &Finding{
//				Action:            changeAction,
//				Content:           content,
//				CommitAuthor:      commit.Author.String(),
//				CommitHash:        commit.Hash.String(),
//				CommitMessage:     strings.TrimSpace(commit.Message),
//				Description:       signature.Description(),
//				FilePath:          fPath,
//				WraithVersion:     version.AppVersion(),
//				LineNumber:        strconv.Itoa(v),
//				RepositoryName:    *repo.Name,
//				RepositoryOwner:   *repo.Owner,
//				Signatureid:       signature.Signatureid(),
//				SignaturesVersion: sess.SignatureVersion,
//				SecretID:          genericID,
//			}
//
//			// Get a proper uid for the finding and setup the urls
//			finding.Initialize(sess.ScanType)
//
//			//if fNew {
//				// Add it to the session
//				sess.AddFinding(finding)
//				sess.Stats.IncrementCommits()
//				sess.Out.Debug("[THREAD #%d][%s] Done analyzing changes in %s\n", tid, *repo.CloneURL, commit.Hash)
//
//			// This will be used to increment the dirty commit stat if any matches are found. A dirty commit
//			// means that a secret was found in that commit. This provides an easier way to manually to look
//			// through the commit history of a given repo.
//			sess.Stats.IncrementCommitsDirty()
//
//				//dirtyCommit = true
//
//				//print realtime data to stdout
//				realTimeOutput(finding, sess)
//			//}
//		}
//	}
//	return finding
//}
