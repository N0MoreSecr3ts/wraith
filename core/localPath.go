package core

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/N0MoreSecr3ts/wraith/version"

	"golang.org/x/sync/errgroup"
)

// Search will walk the path or a given directory and append each viable path to an array
func Search(ctx context.Context, root string, skippablePath []string, sess *Session) ([]string, error) {
	sess.Out.Important("Enumerating Paths\n")
	g, ctx := errgroup.WithContext(ctx)
	paths := make(chan string, 20)

	// get all the paths within a tree
	g.Go(func() error {
		defer close(paths)

		return filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
			// This will check against the combined list of directories that we want to exclude
			// There is the stock list that we pre-defined and then user have the ability to add to this list via the commandline
			for _, p := range skippablePath {
				if strings.HasPrefix(path, p) {
					return nil
				}
			}

			if os.IsPermission(err) {
				return nil
			}
			if !fi.Mode().IsRegular() {
				return nil
			}

			select {
			case paths <- path:
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil
		})
	})

	var m []string
	for r := range paths {
		m = append(m, r)
	}
	return m, g.Wait()
}

// DoFileScan with create a match object and then test for various criteria necessary in order to determine if it should be scanned. This includes if it should be skipped due to a default or user supplied extension, if it matches a test regex, or is in a protected directory or is itself protected. This will only run when doing scanLocalPath.
func DoFileScan(filename string, sess *Session) {

	// Set default values for all pre-requisites for a file scan
	likelyTestFile := false

	// This is the total number of files that we know exist in out path. This does not care about the scan, it is simply the total number of files found
	sess.Stats.IncrementFilesTotal()

	matchFile := newMatchFile(filename)
	if matchFile.isSkippable(sess) {
		sess.Out.Debug("%s is listed as skippable and is being ignored\n", filename)
		sess.Stats.IncrementFilesIgnored()
		return
	}

	// If we are not scanning tests then drop all files that match common test file patterns
	// If we do not want to scan any test files or paths we check for them and then exclude them if they are found
	// The default is to not scan test files or common test paths
	if !sess.ScanTests {
		likelyTestFile = isTestFileOrPath(filename)
	}

	if likelyTestFile {
		// We want to know how many files have been ignored
		sess.Stats.IncrementFilesIgnored()
		sess.Out.Debug("%s is a test file and being ignored\n", filename)
		return
	}

	// Check the file size of the file. If it is greater than the default size then
	// then we increment the ignored file count and pass on through.
	val, msg := IsMaxFileSize(filename, sess)
	if val {

		sess.Stats.IncrementFilesIgnored()
		sess.Out.Debug("%s %s\n", filename, msg)

		return
	}

	if sess.Debug {
		// Print the filename of every file being scanned
		sess.Out.Debug("Analyzing %s\n", filename)
	}

	// Increment the number of files scanned
	sess.Stats.IncrementFilesScanned()

	// Scan the file for know signatures
	for _, signature := range Signatures {
		bMatched, matchMap := signature.ExtractMatch(matchFile, sess, nil)

		var content string // this is because file matches are puking

		// for every instance of the secret that matched the specific rule create a new finding
		for k, v := range matchMap {

			// Increment the total number of findings
			sess.Stats.IncrementFindingsTotal()

			cleanK := strings.SplitAfterN(k, "_", 2)

			content = cleanK[1]

			// destroy the secret if the flag is set
			if sess.HideSecrets {
				content = ""
			}

			if bMatched {
				newFinding := &Finding{
					FilePath:         filename,
					Action:           `File Scan`,
					Description:      signature.Description(),
					SignatureID:      signature.SignatureID(),
					Content:          content,
					RepositoryOwner:  `not-a-repo`,
					RepositoryName:   `not-a-repo`,
					CommitHash:       ``,
					CommitMessage:    ``,
					CommitAuthor:     ``,
					LineNumber:       strconv.Itoa(v),
					SecretID:         generateID(),
					WraithVersion:    version.AppVersion(),
					signatureVersion: sess.SignatureVersion,
				}

				// Add a new finding and increment the total
				newFinding.Initialize(sess)
				sess.AddFinding(newFinding)

				// print the current finding to stdout
				realTimeOutput(newFinding, sess)
			}
		}
	}
}

// ScanDir will scan a directory for all the files and then kick a file scan on each of them
func ScanDir(path string, sess *Session) {

	ctx, cf := context.WithTimeout(context.Background(), 3600*time.Second)
	defer cf()

	// get an slice of of all paths
	files, err1 := Search(ctx, path, sess.SkippablePath, sess)
	if err1 != nil {
		sess.Out.Error("There is an error scanning %s: %s\n", path, err1.Error())
	}

	maxThreads := 100
	sem := make(chan struct{}, maxThreads)

	var wg sync.WaitGroup

	wg.Add(len(files))
	for _, file := range files {
		p := file
		sem <- struct{}{}
		go func() {
			defer wg.Done()

			// scan the specific file if it is found to be a valid candidate
			DoFileScan(p, sess)
			<-sem
		}()
	}

	wg.Wait()
}

// CheckArgs will ensure that both a directory and file are not defined at the same time
func CheckArgs(sFile []string, sDir []string, sess *Session) {
	if sFile != nil && sDir != nil {
		sess.Out.Error("You cannot set both scan-file and scan-dir at the same time\n")
		os.Exit(1)
	}

	if sFile == nil && sDir == nil {
		sess.Out.Error("You must set either a path or file to scan\n")
		os.Exit(1)
	}
}
