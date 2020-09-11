package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"

	"gopkg.in/src-d/go-git.v4/storage/memory"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

//<<<<<<< HEAD
// CloneRepository will crete either an in memory clone of a given repository or clone to a temp dir.
func cloneGithub(cloneConfig *CloneConfiguration) (*git.Repository, string, error) {
	////=======
	//// CloneRepository will create either an in memory clone of a given repository or clone to a temp dir.
	//func CloneGithubRepository(cloneConfig *CloneConfiguration) (*git.Repository, string, error) {
	//>>>>>>> 33e8672995d58dbbbca9fe5a6d5e56505d77f933

	cloneOptions := &git.CloneOptions{
		URL:           *cloneConfig.Url,
		Depth:         *cloneConfig.Depth,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", *cloneConfig.Branch)),
		SingleBranch:  true,
		Tags:          git.NoTags,
		Auth: &http.BasicAuth{
			Username: "doesn't matter",
			Password: *cloneConfig.Token,
		},
	}

	var repository *git.Repository
	var err error
	var dir string

	if !*cloneConfig.InMemClone {
		dir, err = ioutil.TempDir("", "wraith")
		if err != nil {
			return nil, "", err
		}
		repository, err = git.PlainClone(dir, false, cloneOptions)
	} else {
		repository, err = git.Clone(memory.NewStorage(), nil, cloneOptions)
	}
	if err != nil {
		return nil, dir, err
	}
	return repository, dir, nil
}

// Client holds a github api client instance
type githubClient struct {
	apiClient *github.Client // TODO put this into the session struct
}

// TODO make this a single function
// CheckAPIToken will ensure we have a valid github api token
func CheckGithubAPIToken(t string, sess *Session) string {

	// check to make sure the length is proper
	if len(t) != 40 {
		sess.Out.Error("The token is invalid. Please use a valid Github token")
		os.Exit(2)
	}

	// match only letters and numbers and ensure you match 40
	exp1 := regexp.MustCompile(`^[A-Za-z0-9]{40}`)
	if !exp1.MatchString(t) {
		sess.Out.Error("The token is invalid. Please use a valid Github token")
		os.Exit(2)
	}
	return t
}

// InitGithubClient will create a new github client of the type given by the input string. Currently Enterprise and github.com are supported
func (s *Session) InitGitClient() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s.GithubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	// TODO need to make this a switch
	if s.ScanType == "github-enterprise" {

		if s.GithubEnterpriseURL != "" {

			_, err := url.Parse(s.GithubEnterpriseURL)
			if err != nil {
				s.Out.Error("Unable to parse --github-enterprise-url: <%s>", s.GithubEnterpriseURL)
			}
		}
		s.GithubClient, _ = github.NewEnterpriseClient(s.GithubEnterpriseURL, "", tc)
	}

	if s.ScanType == "github" {
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
		s.Client, err = gitlabClient.NewClient(gitlabClient{}, s.GitlabAccessToken, s.Out) // TODO set this to some sort of consistent client, look to github for ideas
		if err != nil {
			s.Out.Fatal("Error initializing GitLab client: %s", err)
		}
	}
}

// GetUserOrganization is used to enumerate the owner in a given org
func (c githubClient) GetUserOrganization(login string) (*Owner, error) {
	ctx := context.Background()
	user, _, err := c.apiClient.Users.Get(ctx, login)
	if err != nil {
		return nil, err
	}
	return &Owner{
		Login:     user.Login,
		ID:        user.ID,
		Type:      user.Type,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
		URL:       user.HTMLURL,
		Company:   user.Company,
		Blog:      user.Blog,
		Location:  user.Location,
		Email:     user.Email,
		Bio:       user.Bio,
	}, nil
}

// GetRepositoriesFromOwner is used gather all the repos associated with the org owner or other user
func (c githubClient) GetRepositoriesFromOwner(target Owner) ([]*Repository, error) {
	var allRepos []*Repository
	ctx := context.Background()
	opt := &github.RepositoryListOptions{
		Type: "sources",
	}

	for {
		repos, resp, err := c.apiClient.Repositories.List(ctx, *target.Login, opt)
		if err != nil {
			return allRepos, err
		}
		for _, repo := range repos {
			if !*repo.Fork {
				r := Repository{
					Owner:         repo.Owner.Login,
					ID:            repo.ID,
					Name:          repo.Name,
					FullName:      repo.FullName,
					CloneURL:      repo.CloneURL,
					URL:           repo.HTMLURL,
					DefaultBranch: repo.DefaultBranch,
					Description:   repo.Description,
					Homepage:      repo.Homepage,
				}
				allRepos = append(allRepos, &r)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil
}

// GetOrganizationMembers will gather all the members of a given organization
func (c githubClient) GetOrganizationMembers(target Owner) ([]*Owner, error) {
	var allMembers []*Owner
	ctx := context.Background()
	opt := &github.ListMembersOptions{}
	for {
		members, resp, err := c.apiClient.Organizations.ListMembers(ctx, *target.Login, opt)
		if err != nil {
			return allMembers, err
		}
		for _, member := range members {
			allMembers = append(allMembers, &Owner{Login: member.Login, ID: member.ID, Type: member.Type})
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allMembers, nil
}
