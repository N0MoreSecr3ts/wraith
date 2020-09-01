package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/src-d/go-git.v4/storage/memory"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// CloneRepository will create either an in memory clone of a given repository or clone to a temp dir.
func CloneGithubRepository(cloneConfig *CloneConfiguration) (*git.Repository, string, error) {

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
		dir, err = ioutil.TempDir("", "wraith") //TODO need to remove this when we are done with it
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
	apiClient *github.Client
}

// TODO make this a single function
// CheckAPIToken will ensure we have a valid github api token
func CheckGithubAPIToken(t string, sess *Session) {

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
}

// NewClient creates a github api client instance using oauth2 credentials
func (c githubClient) NewClient(sess *Session) (apiClient githubClient) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: sess.GithubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	// NewEnterpriseClient creates a github api client for enterprise instances which will use basic auth
	if len(sess.EnterpriseURL) != 0 && sess.EnterpriseScan {
		baseUrl := fmt.Sprintf("%s/api/v3", sess.EnterpriseURL)
		uploadUrl := fmt.Sprintf("%s/api/uploads", sess.EnterpriseURL)
		c.apiClient, _ = github.NewEnterpriseClient(baseUrl, uploadUrl, tc)
	} else {
		c.apiClient = github.NewClient(tc)
	}
	c.apiClient.UserAgent = UserAgent
	return c
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
