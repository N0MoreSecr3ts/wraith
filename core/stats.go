package core

import (
	"sync"
	"time"
)

// Stats hold various runtime statistics used for perf data as well generating various reports
type Stats struct { // TODO alpha sort this
	sync.Mutex

	StartedAt           time.Time // The time we started the scan
	FinishedAt          time.Time // The time we finished the scan
	Status              string    // The running status of a scan for the web interface
	Progress            float64   // The running progress for the bar on the web interface
	RepositoriesTotal   int       // The toatal number of repos discovered
	RepositoriesScanned int       // The total number of repos scanned (not excluded, errors, empty)
	RepositoriesCloned  int       // The total number of repos cloned (excludes errors and excluded, includes empty)
	Organizations       int       // The number of github orgs
	CommitsScanned      int       // The number of commits scanned in a repo
	CommitsDirty        int       // The number of commits in a repo found to have secrets
	FilesScanned        int       // The number of files actually scanned
	FilesIgnored        int       // The number of files ignored (tests, extensions, paths)
	FilesTotal          int       // The total number of files that were processed
	FilesDirty          int
	FindingsTotal       int // The total number of findings. There can be more than one finding per file and more than one finding of the same type in a file
	Users               int // Github users
	Targets             int // The number of dirs, people, orgs, etc on the command line or config file (what do you want wraith to enumerate on)
	Repositories        int // This will point to Repositories Scanned
	Commits             int // This will point to commits scanned
	Findings            int // This will point to findings total
	Files               int // This will point to FilesScanned
}

// IncrementFilesTotal will bump the count of files that have been discovered. This does not reflect
// if the file was scanned/skipped. It is simply a count of files that were found.
func (s *Stats) IncrementFilesTotal() {
	s.Lock()
	defer s.Unlock()
	s.FilesTotal++
	s.Files++
}

// IncrementFilesTotal will bump the count of files that have been discovered. This does not reflect
// if the file was scanned/skipped. It is simply a count of files that were found.
func (s *Stats) IncrementFilesDirty() {
	s.Lock()
	defer s.Unlock()
	s.FilesDirty++
}

// IncrementFilesScanned will bump the count of files that have been scanned successfully.
func (s *Stats) IncrementFilesScanned() {
	s.Lock()
	defer s.Unlock()
	s.FilesScanned++
}

// IncrementFilesIgnored will bump the number of files that have been ignored for various reasons.
func (s *Stats) IncrementFilesIgnored() {
	s.Lock()
	defer s.Unlock()
	s.FilesIgnored++
}

// IncrementFindingsTotal will bump the total number of findings that have been matched. This does
// exclude any other documented criteria.
func (s *Stats) IncrementFindingsTotal() {
	s.Lock()
	defer s.Unlock()
	s.FindingsTotal++
	s.Findings++
}

// IncrementRepositoriesTotal will bump the total number of repositories that have been discovered.
// This will include empty ones as well as those that had errors
func (s *Stats) IncrementRepositoriesTotal() {
	s.Lock()
	defer s.Unlock()
	s.RepositoriesTotal++
	s.Repositories++
}

// IncrementRepositoriesCloned will bump the number of repositories that have been cloned with errors but may be empty
func (s *Stats) IncrementRepositoriesCloned() {
	s.Lock()
	defer s.Unlock()
	s.RepositoriesCloned++
}

// IncrementRepositoriesScanned will bump the total number of repositories that have been scanned and are not empty
func (s *Stats) IncrementRepositoriesScanned() {
	s.Lock()
	defer s.Unlock()
	s.RepositoriesScanned++
}

// IncrementCommitsScanned will bump the number of commits that have been scanned.
// This is scan wide and not on a per repo/org basis
func (s *Stats) IncrementCommitsScanned() {
	s.Lock()
	defer s.Unlock()
	s.CommitsScanned++
}

// IncrementCommitsDirty will bump the number of commits that have been found to be dirty,
// as in they contain one of more findings
func (s *Stats) IncrementCommitsDirty() {
	s.Lock()
	defer s.Unlock()
	s.CommitsDirty++
}

// InitStats will set the initial values for a session
func (s *Session) InitStats() {
	if s.Stats != nil {
		return
	}
	s.Stats = &Stats{
		FilesIgnored:  0,
		FilesScanned:  0,
		FindingsTotal: 0,
		Organizations: 0,
		Progress:      0.0,
		StartedAt:     time.Now(),
		Status:        StatusFinished,
		Users:         0,
		Targets:       0,
		Repositories:  0,
		Commits:       0,
		Findings:      0,
		Files:         0,
	}
}

/// PrintSessionStats will print the performance and sessions stats to stdout at the conclusion of a session scan
func PrintSessionStats(sess *Session) {

	sess.Out.Important("\n--------Results--------\n")
	sess.Out.Important("\n")
	sess.Out.Important("-------Findings------\n")
	sess.Out.Info("Total Findings......: %d\n", sess.Stats.Findings)
	sess.Out.Important("\n")
	sess.Out.Important("--------Files--------\n")
	sess.Out.Info("Total Files.........: %d\n", sess.Stats.FilesTotal)
	sess.Out.Info("Files Scanned.......: %d\n", sess.Stats.FilesScanned)
	sess.Out.Info("Files Ignored.......: %d\n", sess.Stats.FilesIgnored)
	sess.Out.Info("Files Dirty.........: %d\n", sess.Stats.FilesDirty)
	sess.Out.Important("\n")
	sess.Out.Important("---------SCM---------\n")
	sess.Out.Info("Repos Found.........: %d\n", sess.Stats.RepositoriesTotal)
	sess.Out.Info("Repos Cloned........: %d\n", sess.Stats.RepositoriesCloned)
	sess.Out.Info("Repos Scanned.......: %d\n", sess.Stats.RepositoriesScanned)
	sess.Out.Info("Commits Scanned.....: %d\n", sess.Stats.Commits)
	sess.Out.Info("Commits Dirty.......: %d\n", sess.Stats.CommitsDirty)
	sess.Out.Important("\n")
	sess.Out.Important("-------General-------\n")
	sess.Out.Info("Wraith Version......: %s\n", sess.Version)
	sess.Out.Info("Signatures Version..: %s\n", sess.SignatureVersion)
	sess.Out.Info("Elapsed Time........: %s\n\n", time.Since(sess.Stats.StartedAt))
}
