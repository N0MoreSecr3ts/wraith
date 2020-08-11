// Package core represents the core functionality of all commands
package core

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	"wraith/version"

	"github.com/gin-gonic/gin"
)

// These are varios environment variables and tool statuses used in auth and displaying messages
const (
	StatusInitializing = "initializing"
	StatusGathering    = "gathering"
	StatusAnalyzing    = "analyzing"
	StatusFinished     = "finished"
)

// skippableExtensions is an array of extensions that if they match a file that file will be excluded
var defaultIgnoreExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff",
	".tif", ".psd", ".xcf"}

// skippablePathIndicators is an array of directories that will be excluded from all types of scans.
var defaultIgnorePaths = []string{"node_modules/", "vendor/bundle", "vendor/cache", "/proc/"}

var DefaultValues = map[string]interface{}{
	"bind-address":     "127.0.0.1",
	"bind-port":        9393,
	"commit-depth":     0,
	"config-file":      "$HOME/.wraith/config.yaml",
	"debug":            false,
	"github-targets":   "",
	"github-api-token": "0123456789ABCDEFGHIJKLMNOPQRSTUVWXVZabcd",
	"gitlab-targets":   "",
	"gitlab-api-token": "0123456789ABCDEFGHIJ",
	"ignore-extension": "",
	"ignore-path":      "",
	"in-mem-clone":     false,
	"max-file-size":    50,
	"local-dirs":       "",
	"scan-forks":       true,
	"scan-tests":       false,
	"scan-type":        "",
	"silent":           false,
	//"csv":                     false,
	//"db-output":               false,
	//"display-changelog":       false,
	//"json":                    false,
	//"low-priority":            false,
	"match-level":     3,
	"signature-file":  "default_signatures.yml",
	"signatures-path": "$HOME/.wraith/signatures",
	//"signatures-url":               "",
	"scan-dir":                "",
	"scan-file":               "",
	"hide-secrets": false,
	//"test-signatures":              false, // TODO implement this as a bool
}

// Session contains all the necessary values and parameters used during a scan
type Session struct {
	sync.Mutex

	BindAddress       string
	BindPort          int
	Client            IClient `json:"-"`
	CommitDepth       int
	Debug             bool
	Findings          []*Finding
	GithubAccessToken string
	GithubTargets     []string
	GitlabAccessToken string
	GitlabTargets     []string
	HideSecrets       bool
	InMemClone        bool
	MaxFileSize       int64
	NoExpandOrgs      bool
	Out               *Logger `json:"-"`
	LocalDirs         []string
	Repositories      []*Repository
	Router            *gin.Engine `json:"-"`
	SignatureVersion  string
	ScanFork          bool
	ScanTests         bool
	ScanType          string
	Signatures        []*Signature
	Silent            bool
	SkippableExt      []string
	SkippablePath     []string
	Stats             *Stats
	Targets           []*Owner
	Threads           int
	Version           string
	MatchLevel        int
}

// setConfig will set the defaults, and load a config file and environment variables if they are present
func SetConfig() *viper.Viper {

	v := viper.New()

	for key, value := range DefaultValues {
		v.SetDefault(key, value)
	}

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	v.AddConfigPath(home + "/.wraith/")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err == nil {
	}

	v.AutomaticEnv()

	return v
}

// Initialize will set the initial values and options used during a scan session
func (s *Session) Initialize(v *viper.Viper, scanType string) {

	s.BindAddress = v.GetString("bind-address")
	s.BindPort = v.GetInt("bind-port")
	s.CommitDepth = setCommitDepth(v.GetInt("commit-depth"))
	s.Debug = v.GetBool("debug")
	s.GithubAccessToken = v.GetString("github-api-token")
	s.GithubTargets = v.GetStringSlice("github-targets")
	s.GitlabAccessToken = v.GetString("gitlab-api-token")
	s.GitlabTargets = v.GetStringSlice("gitlab-targets")
	s.InMemClone = v.GetBool("in-mem-clone")
	s.MaxFileSize = v.GetInt64("max-file-size")
	s.LocalDirs = v.GetStringSlice("local-dirs")
	s.ScanFork = v.GetBool("scan-forks") //TODO Need to implement
	s.ScanTests = v.GetBool("scan-tests")
	s.ScanType = scanType
	s.Silent = v.GetBool("silent")
	s.Threads = v.GetInt("num-threads")
	s.Version = version.AppVersion()
	//s.CSVOutput = v.GetBool("csv")
	//s.GithubEnterpriseURL = v.GetString("github-enterprise-url")
	//s.GithubURL = v.GetString("github-url")
	//s.JSONOutput = v.GetBool("json")

	s.HideSecrets = v.GetBool("hide-secrets")
	s.MatchLevel = v.GetInt("match-level")

	// add the default directories to the sess if they don't already exist
	for _, e := range defaultIgnorePaths {
		e = strings.TrimSpace(e)
		s.SkippablePath = AppendIfMissing(s.SkippablePath, e)
	}

	// add any additional paths the user requested to exclude to the pre-defined slice
	userIgnorePath := v.GetString("ignore-path")
	if userIgnorePath != "" {
		p := strings.Split(v.GetString("ignore-path"), ",")

		for _, e := range p {
			e = strings.TrimSpace(e)
			s.SkippablePath = AppendIfMissing(s.SkippablePath, e)
		}
	}

	// the default ignorable extensions
	for _, e := range defaultIgnoreExtensions {
		s.SkippableExt = AppendIfMissing(s.SkippableExt, e)
	}

	// add any additional extensions the user requested to ignore
	userIgnoreExtensions := v.GetString("ignore-extension")
	if userIgnoreExtensions != "" {
		e := strings.Split(userIgnoreExtensions, ",")

		for _, f := range e {
			f = strings.TrimSpace(f)
			s.SkippableExt = AppendIfMissing(s.SkippableExt, f)
		}
	}

	s.InitStats()
	s.InitLogger()
	s.InitThreads()
	s.InitAPIClient()

	if !s.Silent {
		s.InitRouter()
	}

	var curSig []Signature
	var combinedSig []Signature
	SignaturesFile := v.GetString("signatures-file")
	if SignaturesFile != "" {
		Signatures := strings.Split(SignaturesFile, ",")

		for _, f := range Signatures {
			f = strings.TrimSpace(f)
			if PathExists(f) {
				curSig = LoadSignatures(f, s.MatchLevel, s)
				combinedSig = append(combinedSig, curSig...)
			}
		}
	} else {
		curSig = LoadSignatures(v.GetString(".")+"default_signatures.yml", s.MatchLevel, s) // TODO implement this
		combinedSig = append(combinedSig, curSig...)
	}

	Signatures = combinedSig
}

// setCommitDepth will set the commit depth to go to during a sess. This is an ugly way of doing it but for the moment it works fine.
func setCommitDepth(c int) int {
	if c == 0 {
		return 9999999999
	}
	return c
}

// Finish is called at the end of a scan session and used to generate discrete data points
// for a given scan session including setting the status of a scan to finished.
func (s *Session) Finish() {
	s.Stats.FinishedAt = time.Now()
	s.Stats.Status = StatusFinished
}

// AddTarget will add a new target to a session to be scanned during that session
func (s *Session) AddTarget(target *Owner) {
	s.Lock()
	defer s.Unlock()
	for _, t := range s.Targets {
		if *target.ID == *t.ID {
			return
		}
	}
	s.Targets = append(s.Targets, target)
	s.Stats.IncrementTargets()
}

// AddRepository will add a given repository to be scanned to a session. This counts as
// the total number of repos that have been gathered during a session.
func (s *Session) AddRepository(repository *Repository) {
	s.Lock()
	defer s.Unlock()
	for _, r := range s.Repositories {
		if *repository.ID == *r.ID {
			return
		}
	}
	s.Repositories = append(s.Repositories, repository)
	s.Stats.IncrementRepositoriesTotal()

}

// AddFinding will add a finding that has been discovered during a session to the list of findings
// for that session
func (s *Session) AddFinding(finding *Finding) {
	s.Lock()
	defer s.Unlock()
	const MaxStrLen = 100
	s.Findings = append(s.Findings, finding)
	s.Stats.IncrementFindingsTotal()
}

// InitStats will set the initial values for a hunt
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

// InitLogger will initialize the logger for the session
func (s *Session) InitLogger() {
	s.Out = &Logger{}
	s.Out.SetDebug(s.Debug)
	s.Out.SetSilent(s.Silent)
}

// InitAPIClient will create a new gitlab or github api client based on the session identifier
func (s *Session) InitAPIClient() {

	switch s.ScanType {
	case "github":
		CheckGithubAPIToken(s.GithubAccessToken)
		s.Client = githubClient.NewClient(githubClient{}, s.GithubAccessToken)
	case "gitlab":
		CheckGitlabAPIToken(s.GitlabAccessToken)
		var err error
		s.Client, err = gitlabClient.NewClient(gitlabClient{}, s.GitlabAccessToken, s.Out)
		if err != nil {
			s.Out.Fatal("Error initializing GitLab client: %s", err)
		}
	default:
		// TODO put something in here when needed
	}
}

// InitThreads will set the correct number of threads based on the commandline flags
func (s *Session) InitThreads() {
	if s.Threads == 0 {
		numCPUs := runtime.NumCPU()
		s.Threads = numCPUs
	}
	runtime.GOMAXPROCS(s.Threads + 2) // thread count + main + web server
}

// InitRouter will configure and start the webserver for graphical output and status messages
func (s *Session) InitRouter() {
	bind := fmt.Sprintf("%s:%d", s.BindAddress, s.BindPort)
	s.Router = NewRouter(s)
	go func(sess *Session) {
		if err := sess.Router.Run(bind); err != nil {
			sess.Out.Fatal("Error when starting web server: %s\n", err)
		}
	}(s)
}

// SaveToFile will save a json representation of the session output to a file
func (s *Session) SaveToFile(location string) error {
	sessionJson, err := json.Marshal(s)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(location, sessionJson, 0644)
	if err != nil {
		return err
	}
	return nil
}

// IncrementTargets will add one to the running target count during the target discovery phase of a session
func (s *Stats) IncrementTargets() {
	s.Lock()
	defer s.Unlock()
	s.Targets++
}

// IncrementRepositories will add one to the running repository count during the target discovery phase of a session
func (s *Stats) IncrementRepositories() {
	s.Lock()
	defer s.Unlock()
	s.Repositories++
}

// IncrementCommits will add one to the running count of commits during the target discovery phase of a session
func (s *Stats) IncrementCommits() {
	s.Lock()
	defer s.Unlock()
	s.Commits++
}

// IncrementFiles will add one to the running count of files during the target discovery phase of a session
func (s *Stats) IncrementFiles() {
	s.Lock()
	defer s.Unlock()
	s.Files++
}

// IncrementFindings will add one to the running count of findings during the target discovery phase of a session
func (s *Stats) IncrementFindings() {
	s.Lock()
	defer s.Unlock()
	s.Findings++
}

// UpdateProgress will update the progress percentage
func (s *Stats) UpdateProgress(current int, total int) {
	s.Lock()
	defer s.Unlock()
	if current >= total {
		s.Progress = 100.0
	} else {
		s.Progress = (float64(current) * float64(100)) / float64(total)
	}
}

// NewSession  is the entry point for starting a new scan session
func NewSession(v *viper.Viper, scanType string) (*Session, error) {
	var session Session

	session.Initialize(v, scanType)

	return &session, nil
}
