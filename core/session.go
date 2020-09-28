// Package core represents the core functionality of all commands
package core

import (
	"encoding/json"
	//"context"
	"fmt"
	"github.com/google/go-github/github"
	//"golang.org/x/oauth2"
	"io/ioutil"
	//"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	"wraith/version"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

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
	"enterprise-scan":	false,
	"enterprise-url":	  "",
	"github-targets":   "",
	"github-api-token": "0123456789ABCDEFGHIJKLMNOPQRSTUVWXVZabcd",
	"gitlab-targets":   "",
	"gitlab-api-token": "0123456789ABCDEFGHIJ",
	"ignore-extension": "",
	"ignore-path":      "",
	"in-mem-clone":     false,
	"max-file-size":    50,
	"num-threads":      0,
	"local-dirs":       nil,
	"local-files":      nil,
	"scan-forks":       true,
	"scan-tests":       false,
	"scan-type":        "",
	"silent":           false,
	"CSVOutput":        false,
	"JSONOutput":       false,
	"match-level":      3,
	"output-dir":		    "./",
	"output-prefix":	  "wraith",
	"signature-file":   "$HOME/.wraith/signatures/default.yml",
	"signature-path":   "$HOME/.wraith/signatures/",
	"signature-url":    "",
	"scan-dir":         "",
	"scan-file":        "",
	"hide-secrets":     false,
	"github-url":       "https://api.github.com",
	"gitlab-url":       "", // TODO set the default
	"rules-url":        "git@example.com:foo/bar.git",
	"github-orgs":      "",
	"github-repos":     "",
	"github-users":     "",
}

// Session contains all the necessary values and parameters used during a scan
type Session struct {
	sync.Mutex

	BindAddress       string
	BindPort          int
	Client            IClient `json:"-"`
	CommitDepth       int
	CSVOutput         bool
	Debug             bool
	ExpandOrgs        bool
	Findings          []*Finding
	GithubAccessToken string
	EnterpriseScan    bool
	EnterpriseURL     string
  Organizations     []*github.Organization
	GithubClient        *github.Client `json:"-"`
	GithubTargets     []string
	GitlabAccessToken string
	GitlabTargets     []string
  GithubUsers         []*github.User
	HideSecrets       bool
	InMemClone        bool
	JSONOutput        bool
	MaxFileSize       int64
	NoExpandOrgs      bool
	Out               *Logger `json:"-"`
	OutputDir         string
	OutputPrefix      string
	LocalDirs         []string
	LocalFiles        []string
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
  GithubURL           string
	GitlabURL           string
	UserDirtyNames      string
	UserDirtyOrgs       string
	UserDirtyRepos      string
	UserLogins          []string
	UserOrgs            []string
	UserRepos           []string

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
	s.CSVOutput = v.GetBool("csv")
	s.Debug = v.GetBool("debug")
	s.ExpandOrgs = v.GetBool("expaand-orgs")
	s.GithubEnterpriseURL = v.GetString("github-enterprise-url")
	s.GithubAccessToken = v.GetString("github-api-token")
	s.GithubTargets = v.GetStringSlice("github-targets")
	s.GitlabAccessToken = v.GetString("gitlab-api-token")
	s.GitlabTargets = v.GetStringSlice("gitlab-targets")
	s.HideSecrets = v.GetBool("hide-secrets")
	s.InMemClone = v.GetBool("in-mem-clone")
	s.JSONOutput = v.GetBool("json")
	s.LocalDirs = v.GetStringSlice("local-dirs")
	s.MaxFileSize = v.GetInt64("max-file-size")
	s.MatchLevel = v.GetInt("match-level")
	s.OutputDir = v.GetString("output-dir")
	s.OutputPrefix = v.GetString("output-prefix")
	s.ScanFork = v.GetBool("scan-forks") //TODO Need to implement
	s.ScanTests = v.GetBool("scan-tests")
	s.ScanType = scanType
	s.Silent = v.GetBool("silent")
	s.Threads = v.GetInt("num-threads")
	s.Version = version.AppVersion()
	v.GetStringSlice("scan-dir")
	v.GetStringSlice("scan-file")

	// add the default directories to the sess if they don't already exist
	for _, e := range defaultIgnorePaths {
		e = strings.TrimSpace(e)
		s.SkippablePath = AppendIfMissing(s.SkippablePath, e)
	}

	// add any additional paths the user requested to exclude to the pre-defined slice
	userIgnorePath := v.GetString("ignore-path")
	if userIgnorePath != "" {
		p := strings.Split(v.GetString("ignore-path"), ",") // TODO make slice

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
		e := strings.Split(userIgnoreExtensions, ",") // TODO make slice

		for _, f := range e {
			f = strings.TrimSpace(f)
			s.SkippableExt = AppendIfMissing(s.SkippableExt, f)
		}
	}

	s.InitStats()
	s.InitLogger()
	s.InitThreads()

	if !s.Silent {
		s.InitRouter()
	}

	var curSig []Signature
	var combinedSig []Signature
	SignaturesFile := v.GetString("signature-file")
	if SignaturesFile != "" {
		Signatures := strings.Split(SignaturesFile, ",") // TODO make slice

		for _, f := range Signatures {
			f = strings.TrimSpace(f)
			h := SetHomeDir(f)
			if PathExists(h, s) {
				curSig = LoadSignatures(h, s.MatchLevel, s)
				combinedSig = append(combinedSig, curSig...)
			}
		}
	} // TODO need to catch this error here
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

// InitGithubClient will create a new github client of the type given by the input string. Currently Enterprise and github.com are supported
//func (s *Session) InitAPIClient() {
//	ctx := context.Background()
//	ts := oauth2.StaticTokenSource(
//		&oauth2.Token{AccessToken: s.GithubAccessToken},
//	)
//	tc := oauth2.NewClient(ctx, ts)
//
//	if s.ScanType == "github-enterprise" {
//
//		if s.GithubEnterpriseURL != "" {
//
//			_, err := url.Parse(s.GithubEnterpriseURL)
//			if err != nil {
//				s.Out.Error("Unable to parse --github-enterprise-url: <%s>", s.GithubEnterpriseURL)
//			}
//		}
//		s.GithubClient, _ = github.NewEnterpriseClient(s.GithubEnterpriseURL, "", tc)
//	}
//
//	if t == "github" {
//		if s.GithubURL != "" {
//			_, err := url.Parse(s.GithubURL)
//			if err != nil {
//				s.Out.Error("Unable to parse --github-url: <%s>", s.GithubURL)
//			}
//		}
//		s.GithubClient = github.NewClient(tc)
//	}
//}

// InitAPIClient will create a new gitlab or github api client based on the session identifier
//func (s *Session) InitAPIClient() {
//
//	switch s.ScanType {
//	case "github":
//		CheckGithubAPIToken(s.GithubAccessToken, s)
//		s.Client = githubClient.NewClient(githubClient{}, s)
//	case "gitlab":
//		CheckGitlabAPIToken(s.GitlabAccessToken, s)
//		var err error
//		s.Client, err = gitlabClient.NewClient(gitlabClient{}, s.GitlabAccessToken, s.Out)
//		if err != nil {
//			s.Out.Fatal("Error initializing GitLab client: %s", err)
//		}
//	default:
//		// TODO put something in here when needed
//	}
//}

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
func NewSession(v *viper.Viper, scanType string) *Session { // TODO refactor out this function
	var session Session

	session.Initialize(v, scanType)

	return &session
}
