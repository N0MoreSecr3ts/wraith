package core

import (
	"sync"
	"time"
)

//type Stats struct {
//	sync.Mutex
//
//	//StartedAt    time.Time
//	//FinishedAt   time.Time
//	//Status       string
//	//Progress     float64
//	//Targets      int
//	//Repositories int
//	//Commits      int
//	//Files        int
//	//Findings     int
//}

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
	FindingsTotal       int       // The total number of findings. There can be more than one finding per file and more than one finding of the same type in a file
	Users               int       // Github users
	Targets             int       // The number of dirs, people, orgs, etc on the command line or config file (what do you want wraith to enumerate on)
	Repositories        int       // This will point to Repositories Scanned
	Commits             int       // This will point to commits scanned
	Findings            int       // This will point to findings total
	Files               int       // This will point to FilesScanned
}

// IncrementFilesTotal will bump the count of files that have been discovered. This does not reflect
// if the file was scanned/skipped. It is simply a count of files that were found.
func (s *Stats) IncrementFilesTotal() {
	s.Lock()
	defer s.Unlock()
	s.FilesTotal++
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
