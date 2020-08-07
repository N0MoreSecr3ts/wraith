// Package github represents github specific functionality
package github

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"wraith/common"
)

// Client holds a github api client instance
type Client struct {
	apiClient *github.Client
}

// CheckAPIToken will ensure we have a valid github api token
func CheckAPIToken(t string) {

	// check to make sure the length is proper
	if len(t) != 40 {
		fmt.Println("The token is invalid. Please use a valid Github token")
		os.Exit(2)
	}

	// match only letters and numbers and ensure you match 40
	exp1 := regexp.MustCompile(`^[A-Za-z0-9]{40}`)
	if !exp1.MatchString(t) {
		fmt.Println("The token is invalid. Please use a valid Github token")
		os.Exit(2)
	}
}

// NewClient creates a github api client instance using oauth2 credentials
func (c Client) NewClient(token string) (apiClient Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	c.apiClient = github.NewClient(tc)
	c.apiClient.UserAgent = common.UserAgent
	return c
}

// GetUserOrganization is used to enumerate the owner in a given org
func (c Client) GetUserOrganization(login string) (*common.Owner, error) {
	ctx := context.Background()
	user, _, err := c.apiClient.Users.Get(ctx, login)
	if err != nil {
		return nil, err
	}
	return &common.Owner{
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
func (c Client) GetRepositoriesFromOwner(target common.Owner) ([]*common.Repository, error) {
	var allRepos []*common.Repository
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
				r := common.Repository{
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
func (c Client) GetOrganizationMembers(target common.Owner) ([]*common.Owner, error) {
	var allMembers []*common.Owner
	ctx := context.Background()
	opt := &github.ListMembersOptions{}
	for {
		members, resp, err := c.apiClient.Organizations.ListMembers(ctx, *target.Login, opt)
		if err != nil {
			return allMembers, err
		}
		for _, member := range members {
			allMembers = append(allMembers, &common.Owner{Login: member.Login, ID: member.ID, Type: member.Type})
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allMembers, nil
}
