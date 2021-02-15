// Package core represents the core functionality of all commands
package core

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
	"net/http"
	"net/url"
	"sync"
)

// Set easier names to refer to
const (
	TargetTypeUser         = "User"
	TargetTypeOrganization = "Organization"
)

// GithubRepository holds the necessary information for a repository,
// this data is specific to Github.
type GithubRepository struct {
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

// CloneConfiguration holds the configurations for cloning a repo
type CloneConfiguration struct {
	InMemClone *bool
	URL        *string
	Username   *string
	Token      *string
	Branch     *string
	Depth      *int
}

// Owner holds the info that we want for a repo owner
type Owner struct {
	Login     *string
	ID        *int64
	Type      *string
	Name      *string
	AvatarURL *string
	URL       *string
	Company   *string
	Blog      *string
	Location  *string
	Email     *string
	Bio       *string
}

// Repository holds the info we want for a repo itself
type Repository struct {
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

// EmptyTreeCommit is a dummy commit id used as a placeholder and for testing
const (
	EmptyTreeCommitID = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"
)

// GetParentCommit will get the parent commit from a specific point. If the current commit
// has no parents then it will create a dummy commit.
func getParentCommit(commit *object.Commit, repo *git.Repository) (*object.Commit, error) {
	if commit.NumParents() == 0 {
		parentCommit, err := repo.CommitObject(plumbing.NewHash(EmptyTreeCommitID))
		if err != nil {
			return nil, err
		}
		return parentCommit, nil
	}
	parentCommit, err := commit.Parents().Next()
	if err != nil {
		return nil, err
	}
	return parentCommit, nil
}

// GetRepositoryHistory gets the commit history of a repository
func GetRepositoryHistory(repository *git.Repository) ([]*object.Commit, error) {
	var commits []*object.Commit
	ref, err := repository.Head()
	if err != nil {
		return nil, err
	}
	cIter, err := repository.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}
	_ = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})
	return commits, nil
}

// GetChanges will get the changes between to specific commits. It grabs the parent commit of
// the one being passed and uses that to fetch the tree for that commit. If no commit is found,
// it will create a fake on. It then takes that parent tree along with the tree for the commit
// passed in and does a diff
func GetChanges(commit *object.Commit, repo *git.Repository) (object.Changes, error) {
	parentCommit, err := getParentCommit(commit, repo)
	if err != nil {
		return nil, err
	}

	commitTree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	parentCommitTree, err := parentCommit.Tree()
	if err != nil {
		return nil, err
	}

	changes, err := object.DiffTree(parentCommitTree, commitTree)
	if err != nil {
		return nil, err
	}
	return changes, nil
}

// GetChangeAction returns a more condensed and user friendly action for further reference
func GetChangeAction(change *object.Change) string {
	action, err := change.Action()
	if err != nil {
		return "Unknown"
	}
	switch action {
	case merkletrie.Insert:
		return "Insert"
	case merkletrie.Modify:
		return "Modify"
	case merkletrie.Delete:
		return "Delete"
	default:
		return "Unknown"
	}
}

// GetChangePath will set the action of the commit for further action
func GetChangePath(change *object.Change) string {
	action, err := change.Action()
	if err != nil {
		return change.To.Name
	}

	if action == merkletrie.Delete {
		return change.From.Name
	}

	return change.To.Name
}

// GetChangeContent will get the contents of a git change or patch.
func GetChangeContent(change *object.Change) (result string, contentError error) {
	//temporary response to:  https://github.com/sergi/go-diff/issues/89
	defer func() {
		if err := recover(); err != nil {
			contentError = errors.New("panic occurred while retrieving change content: ")
		}
	}()
	patch, err := change.Patch()
	if err != nil {
		return "", err
	}
	for _, filePatch := range patch.FilePatches() {
		if filePatch.IsBinary() {
			continue
		}
		for _, chunk := range filePatch.Chunks() {
			result += chunk.Content()
		}
	}
	return result, nil
}

// GatherGitlabRepositories will gather all repositories associated with a given target during a scan session.
// This is done using threads, whose count is set via commandline flag. Care much be taken to avoid rate
// limiting associated with suspected DOS attacks.
func GatherGitlabRepositories(sess *Session) {
	var ch = make(chan *Owner, len(sess.Targets))
	sess.Out.Debug("Number of targets: %d\n", len(sess.Targets))
	var wg sync.WaitGroup
	var threadNum int
	if len(sess.Targets) == 1 {
		threadNum = 1
	} else if len(sess.Targets) <= sess.Threads {
		threadNum = len(sess.Targets) - 1
	} else {
		threadNum = sess.Threads
	}
	wg.Add(threadNum)
	sess.Out.Debug("Threads for repository gathering: %d\n", threadNum)
	for i := 0; i < threadNum; i++ {
		go func() {
			for {
				target, ok := <-ch
				if !ok {
					wg.Done()
					return
				}
				repos, err := sess.Client.GetRepositoriesFromOwner(*target)
				if err != nil {
					sess.Out.Error(" Failed to retrieve repositories from %s: %s\n", *target.Login, err)
				}
				if len(repos) == 0 {
					continue
				}
				for _, repo := range repos {
					sess.Out.Debug(" Retrieved repository: %s\n", *repo.CloneURL)
					sess.AddRepository(repo)
				}
				sess.Out.Info(" Retrieved %d %s from %s\n", len(repos), Pluralize(len(repos), "repository", "repositories"), *target.Login)
			}
		}()
	}

	for _, target := range sess.Targets {
		ch <- target
	}
	close(ch)
	wg.Wait()
}

// InitGitClient will create a new github client of the type given by the input string. Currently Enterprise and github.com are supported
func (s *Session) InitGitClient() {

	// TODO need to make this a switch
	if s.ScanType == "github-enterprise" {

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		sslcli := &http.Client{Transport: tr}

		ctx := context.Background()
		ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: s.GithubAccessToken},
		)
		tc := oauth2.NewClient(ctx, ts)

		if s.GithubEnterpriseURL != "" {

			_, err := url.Parse(s.GithubEnterpriseURL)
			if err != nil {
				s.Out.Error("Unable to parse --github-enterprise-url: <%s>", s.GithubEnterpriseURL)
			}
		}
		s.GithubClient, _ = github.NewEnterpriseClient(s.GithubEnterpriseURL, "", tc)
	}

	if s.ScanType == "github" {

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		sslcli := &http.Client{Transport: tr}

		ctx := context.Background()
		ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: s.GithubAccessToken},
		)
		tc := oauth2.NewClient(ctx, ts)

		if s.GithubURL != "" {
			_, err := url.Parse(s.GithubURL)
			if err != nil {
				s.Out.Error("Unable to parse --github-url: <%s>", s.GithubURL)
			}
		}
		s.GithubClient = github.NewClient(tc)
	}

	if s.ScanType == "gitlab" { // TODO need to refactor all this
		CheckGitlabAPIToken(s.GitlabAccessToken, s) // TODO move this out
		var err error
		// TODO need to add in the bits to parse the url here as well
		// TODO set this to some sort of consistent client, look to github for ideas
		s.Client, err = gitlabClient.NewClient(gitlabClient{}, s.GitlabAccessToken, s.Out)
		if err != nil {
			s.Out.Fatal("Error initializing GitLab client: %s", err)
		}
	}
}

// cloneRepository will clone a given repository based upon a configured set or options a user provides.
// This is a catchall for all different types of repos and create a single entry point for cloning a repo.
func cloneRepository(sess *Session, repo *Repository, threadID int) (*git.Repository, string, error) {
	sess.Out.Debug("[THREAD #%d][%s] Cloning repository...\n", threadID, *repo.CloneURL)

	var clone *git.Repository
	var path string
	var err error

	switch sess.ScanType {
	case "github":
		cloneConfig := CloneConfiguration{
			URL:        repo.CloneURL,
			Branch:     repo.DefaultBranch,
			Depth:      &sess.CommitDepth,
			InMemClone: &sess.InMemClone,
			Token:      &sess.GithubAccessToken,
		}
		// Clone a github repo
		clone, path, err = cloneGithub(&cloneConfig)

	case "github-enterprise":
		cloneConfig := CloneConfiguration{
			URL:        repo.CloneURL,
			Branch:     repo.DefaultBranch,
			Depth:      &sess.CommitDepth,
			InMemClone: &sess.InMemClone,
			Token:      &sess.GithubAccessToken,
		}
		// Clone a github repo
		clone, path, err = cloneGithub(&cloneConfig)

	case "gitlab":
		userName := "oauth2"
		cloneConfig := CloneConfiguration{
			URL:        repo.CloneURL,
			Branch:     repo.DefaultBranch,
			Depth:      &sess.CommitDepth,
			Token:      &sess.GitlabAccessToken, // TODO Is this need since we already have a client?
			InMemClone: &sess.InMemClone,
			Username:   &userName,
		}
		// Clone a gitlab repo
		clone, path, err = cloneGitlab(&cloneConfig)
	case "localGit":
		cloneConfig := CloneConfiguration{
			URL:        repo.CloneURL,
			Branch:     repo.DefaultBranch,
			Depth:      &sess.CommitDepth,
			InMemClone: &sess.InMemClone,
		}
		// Clone a local repo
		clone, path, err = cloneLocal(&cloneConfig)

	}
	if err != nil {
		switch err.Error() {
		case "remote repository is empty":
			sess.Out.Error("Repository %s is empty: %s\n", *repo.CloneURL, err)
			sess.Stats.IncrementRepositoriesCloned()
			return nil, "", err
		default:
			sess.Out.Error("Error cloning repository %s: %s\n", *repo.CloneURL, err)
			return nil, "", err
		}
	}
	sess.Stats.IncrementRepositoriesCloned()
	sess.Out.Debug("[THREAD #%d][%s] Cloned repository to: %s\n", threadID, *repo.CloneURL, path)
	return clone, path, err
}
