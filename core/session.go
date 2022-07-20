// Package core represents the core functionality of all commands
package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/N0MoreSecr3ts/wraith/version"
	"github.com/google/go-github/github"
	"github.com/mitchellh/go-homedir"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// WraithConfig holds the configuration data the commands
var WraithConfig *viper.Viper

// These are various environment variables and tool statuses used in auth and displaying messages
const (
	StatusInitializing = "initializing"
	StatusGathering    = "gathering"
	StatusAnalyzing    = "analyzing"
	StatusFinished     = "finished"
)

// defaultIgnoreExtensions is an array of extensions that if they match a file that file will be excluded
var defaultIgnoreExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff",
	".tif", ".psd", ".xcf"}

// defaultIgnorePaths is an array of directories that will be excluded from all types of scans.
var defaultIgnorePaths = []string{"node_modules/", "vendor/bundle", "vendor/cache", "/proc/"}

// DefaultValues is a map of all flag default values and other mutable variables
var DefaultValues = map[string]interface{}{
	"bind-address":                "127.0.0.1",
	"bind-port":                   9393,
	"commit-depth":                -1,
	"config-file":                 "$HOME/.wraith/config.yaml",
	"csv":                         false,
	"debug":                       false,
	"add-org-members":             false,
	"github-enterprise-url":       "",
	"github-api-token":            "",
	"github-enterprise-api-token": "",
	"gitlab-targets":              nil,
	"gitlab-api-token":            "",
	"ignore-extension":            nil,
	"ignore-path":                 nil,
	"in-mem-clone":                false,
	"json":                        false,
	"max-file-size":               10,
	"num-threads":                 -1,
	"local-paths":                 nil,
	"scan-forks":                  false,
	"scan-tests":                  false,
	"scan-type":                   "",
	"silent":                      false,
	"confidence-level":            3,
	"signature-file":              "$HOME/.wraith/signatures/default.yaml",
	"signature-path":              "$HOME/.wraith/signatures/",
	"scan-dir":                    nil,
	"scan-file":                   nil,
	"hide-secrets":                false,
	"github-url":                  "https://api.github.com",
	//"gitlab-url":                 "", // TODO set the default
	"rules-url":               "",
	"signatures-path":         "$HOME/.wraith/signatures/",
	"signatures-url":          "https://github.com/N0MoreSecr3ts/wraith-signatures",
	"signatures-version":      "",
	"test-signatures":         false,
	"github-enterprise-orgs":  nil,
	"github-enterprise-repos": nil,
	"github-orgs":             nil,
	"github-repos":            nil,
	"github-users":            nil,
	"web-server":              false,
}

// Session contains all the necessary values and parameters used during a scan
type Session struct {
	sync.Mutex

	BindAddress         string
	BindPort            int
	Client              IClient `json:"-"`
	CommitDepth         int
	ConfidenceLevel     int
	CSVOutput           bool
	Debug               bool
	ExpandOrgs          bool
	Findings            []*Finding
	GithubAccessToken   string
	GithubClient        *github.Client `json:"-"`
	GithubEnterpriseURL string
	GithubURL           string
	GitlabAccessToken   string
	GitlabTargets       []string
	GitlabURL           string
	GithubUsers         []*github.User
	HideSecrets         bool
	InMemClone          bool
	JSONOutput          bool
	LocalPaths          []string
	MaxFileSize         int64
	Organizations       []*github.Organization
	Out                 *Logger `json:"-"`
	Repositories        []*Repository
	Router              *gin.Engine `json:"-"`
	SignatureVersion    string
	ScanFork            bool
	ScanTests           bool
	ScanType            string
	Signatures          []*Signature
	Silent              bool
	SkippableExt        []string
	SkippablePath       []string
	Stats               *Stats
	Targets             []*Owner
	Threads             int
	UserDirtyNames      []string
	UserDirtyOrgs       []string
	UserDirtyRepos      []string
	UserLogins          []string
	UserOrgs            []string
	UserRepos           []string
	WebServer           bool
	WraithVersion       string
}

// githubRepository is the holds the necessary fields in a simpler structure
type githubRepository struct {
	Owner         *string
	ID            *int64
	Name          *string
	FullName      *string
	CloneURL      *string
	URL           *string
	DefaultBranch *string
	Description   *string
	Homepage      *string
}

// SetConfig will set the defaults, and load a config file and environment variables if they are present
func SetConfig() {
	for key, value := range DefaultValues {
		viper.SetDefault(key, value)
	}

	configFile := viper.GetString("config-file")

	if configFile != DefaultValues["config-file"] {
		viper.SetConfigFile(configFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home + "/.wraith/")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
	}

	viper.AutomaticEnv()

	WraithConfig = viper.GetViper()
}

// Initialize will set the initial values and options used during a scan session
func (s *Session) Initialize(scanType string) {

	s.BindAddress = WraithConfig.GetString("bind-address")
	s.BindPort = WraithConfig.GetInt("bind-port")
	s.CommitDepth = setCommitDepth(WraithConfig.GetFloat64("commit-depth"))
	s.CSVOutput = WraithConfig.GetBool("csv")
	s.Debug = WraithConfig.GetBool("debug")
	s.ExpandOrgs = WraithConfig.GetBool("expand-orgs")
	s.GithubEnterpriseURL = WraithConfig.GetString("github-enterprise-url")
	s.GithubAccessToken = WraithConfig.GetString("github-api-token")
	s.GitlabAccessToken = WraithConfig.GetString("gitlab-api-token")
	s.GitlabTargets = WraithConfig.GetStringSlice("gitlab-targets")
	s.HideSecrets = WraithConfig.GetBool("hide-secrets")
	s.InMemClone = WraithConfig.GetBool("in-mem-clone")
	s.JSONOutput = WraithConfig.GetBool("json")
	s.MaxFileSize = WraithConfig.GetInt64("max-file-size")
	s.ConfidenceLevel = WraithConfig.GetInt("confidence-level")
	s.ScanFork = WraithConfig.GetBool("scan-forks")
	s.ScanTests = WraithConfig.GetBool("scan-tests")
	s.ScanType = scanType
	s.Silent = WraithConfig.GetBool("silent")
	s.Threads = WraithConfig.GetInt("num-threads")
	s.WraithVersion = version.AppVersion()
	s.WebServer = WraithConfig.GetBool("web-server")

	if s.ScanType == "localGit" {
		s.LocalPaths = WraithConfig.GetStringSlice("local-repos")
	} else if s.ScanType == "localPath" {
		s.LocalPaths = WraithConfig.GetStringSlice("local-paths")
	}

	// Add the default directories to the sess if they don't already exist
	for _, e := range defaultIgnorePaths {
		e = strings.TrimSpace(e)
		s.SkippablePath = AppendIfMissing(s.SkippablePath, e)
	}

	// add any additional paths the user requested to exclude to the pre-defined slice
	for _, e := range WraithConfig.GetStringSlice("ignore-path") {
		e = strings.TrimSpace(e)
		s.SkippablePath = AppendIfMissing(s.SkippablePath, e)
	}

	// the default ignorable extensions
	for _, e := range defaultIgnoreExtensions {
		s.SkippableExt = AppendIfMissing(s.SkippableExt, e)
	}

	// add any additional extensions the user requested to ignore
	for _, f := range WraithConfig.GetStringSlice("ignore-extension") {
		f = strings.TrimSpace(f)
		s.SkippableExt = AppendIfMissing(s.SkippableExt, f)
	}

	s.InitStats()
	s.InitLogger()
	s.InitThreads()

	if !s.Silent && s.WebServer {
		s.InitRouter()
	}

	var curSig []Signature
	var combinedSig []Signature

	// TODO need to catch this error here
	for _, f := range WraithConfig.GetStringSlice("signature-file") {
		f = strings.TrimSpace(f)
		h := SetHomeDir(f, s)
		if PathExists(h, s) {
			curSig = LoadSignatures(h, s.ConfidenceLevel, s)
			combinedSig = append(combinedSig, curSig...)
		}
	}
	Signatures = combinedSig
}

// setCommitDepth will set the commit depth for the current session. This is an ugly way of doing it
// but for the moment it works fine.
// TODO dynamically acquire the commit depth of a given repo
func setCommitDepth(c float64) int {
	if c == -1 {
		return 9999999999
	}
	return int(c)
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
	sessionJSON, err := json.Marshal(s)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(location, sessionJSON, 0644)
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

// IncrementCommitsTotal will add one to the running count of commits during the target discovery phase of a session
func (s *Stats) IncrementCommitsTotal() {
	s.Lock()
	defer s.Unlock()
	s.CommitsTotal++
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
	//s.Lock() TODO REMOVE ME
	//defer s.Unlock() TODO REMOVE ME
	if current >= total {
		s.Progress = 100.0
	} else {
		s.Progress = (float64(current) * float64(100)) / float64(total)
	}
}

// NewSession  is the entry point for starting a new scan session
func NewSession(scanType string) *Session {
	var session Session

	session.Initialize(scanType)

	return &session
}
