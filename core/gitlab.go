package core

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"sync"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// CloneRepository will create either an in memory clone of a given repository or clone to a temp dir.
func cloneGitlab(cloneConfig *CloneConfiguration) (*git.Repository, string, error) {

	cloneOptions := &git.CloneOptions{
		URL:           *cloneConfig.URL,
		Depth:         *cloneConfig.Depth,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", *cloneConfig.Branch)),
		SingleBranch:  true,
		Tags:          git.NoTags,
		Auth: &http.BasicAuth{
			Username: *cloneConfig.Username,
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

// Client holds a gitlab api client instance
type gitlabClient struct {
	apiClient *gitlab.Client
	logger    *Logger
}

// NewClient creates a gitlab api client instance using a token
func (c gitlabClient) NewClient(token string, logger *Logger) (gitlabClient, error) {
	var err error
	c.apiClient, err = gitlab.NewClient(token)
	if err != nil {
		return gitlabClient{}, err
	}
	c.apiClient.UserAgent = UserAgent
	c.logger = logger
	return c, nil
}

// CheckGitlabAPIToken will ensure we have a valid github api token
func CheckGitlabAPIToken(t string, sess *Session) string {

	// check to make sure the length is proper
	if len(t) != 20 {
		sess.Out.Error("Gitlab token is invalid\n")
		os.Exit(2)
	}

	return t
}

// GetUserOrganization is used to enumerate the owner in a given org
func (c gitlabClient) GetUserOrganization(login string) (*Owner, error) {
	emptyString := gitlab.String("")
	org, orgErr := c.getOrganization(login)
	if orgErr != nil {
		user, userErr := c.getUser(login)
		if userErr != nil {
			return nil, userErr
		}
		id := int64(user.ID)
		return &Owner{
			Login:     gitlab.String(user.Username),
			ID:        &id,
			Type:      gitlab.String(TargetTypeUser),
			Name:      gitlab.String(user.Name),
			AvatarURL: gitlab.String(user.AvatarURL),
			URL:       gitlab.String(user.WebsiteURL),
			Company:   gitlab.String(user.Organization),
			Blog:      emptyString,
			Location:  emptyString,
			Email:     gitlab.String(user.PublicEmail),
			Bio:       gitlab.String(user.Bio),
		}, nil
	}
	id := int64(org.ID)
	return &Owner{
		Login:     gitlab.String(org.Name),
		ID:        &id,
		Type:      gitlab.String(TargetTypeOrganization),
		Name:      gitlab.String(org.Name),
		AvatarURL: gitlab.String(org.AvatarURL),
		URL:       gitlab.String(org.WebURL),
		Company:   gitlab.String(org.FullName),
		Blog:      emptyString,
		Location:  emptyString,
		Email:     emptyString,
		Bio:       gitlab.String(org.Description),
	}, nil

}

// GetOrganizationMembers will gather all the members of a given organization
func (c gitlabClient) GetOrganizationMembers(target Owner) ([]*Owner, error) {
	var allMembers []*Owner
	opt := &gitlab.ListGroupMembersOptions{}
	sID := strconv.FormatInt(*target.ID, 10) //safely downcast an int64 to an int
	for {
		members, resp, err := c.apiClient.Groups.ListAllGroupMembers(sID, opt)
		if err != nil {
			return nil, err
		}
		for _, member := range members {
			id := int64(member.ID)
			allMembers = append(allMembers,
				&Owner{
					Login: gitlab.String(member.Username),
					ID:    &id,
					Type:  gitlab.String(TargetTypeUser)})
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allMembers, nil
}

// GetRepositoriesFromOwner is used gather all the repos associated with the org owner or other user
func (c gitlabClient) GetRepositoriesFromOwner(target Owner) ([]*Repository, error) {
	var allProjects []*Repository
	id := int(*target.ID)
	if *target.Type == TargetTypeUser {
		userProjects, err := c.getUserProjects(id)
		if err != nil {
			return nil, err
		}
		for _, project := range userProjects {
			allProjects = append(allProjects, project)
		}
	} else {
		groupProjects, err := c.getGroupProjects(target)
		if err != nil {
			return nil, err
		}
		for _, project := range groupProjects {
			allProjects = append(allProjects, project)
		}
	}
	return allProjects, nil
}

// getUser will get the necessary info from a given user
func (c gitlabClient) getUser(login string) (*gitlab.User, error) {
	users, _, err := c.apiClient.Users.ListUsers(&gitlab.ListUsersOptions{Username: gitlab.String(login)})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("No GitLab %s or %s %s was found.  If you are targeting a GitLab group, be sure to"+
			" use an ID in place of a name.",
			strings.ToLower(TargetTypeUser),
			strings.ToLower(TargetTypeOrganization),
			login)
	}
	return users[0], err
}

// getOrganization will get the necessary info from an org
func (c gitlabClient) getOrganization(login string) (*gitlab.Group, error) {
	id, err := strconv.Atoi(login)
	if err != nil {
		return nil, err
	}
	org, _, err := c.apiClient.Groups.GetGroup(id)
	if err != nil {
		return nil, err
	}
	return org, err
}

// getUserProjects will gather the projects associated with a given user
func (c gitlabClient) getUserProjects(id int) ([]*Repository, error) {
	var allUserProjects []*Repository
	listUserProjectsOps := &gitlab.ListProjectsOptions{}

	var wg sync.WaitGroup
	var mut sync.Mutex

	for {
		projects, response, err := c.apiClient.Projects.ListUserProjects(id, listUserProjectsOps)
		if err != nil {
			return nil, err
		}

		wg.Add(1)

		go func() {
			for _, project := range projects {
				//don't capture forks
				if project.ForkedFromProject == nil {
					id := int64(project.ID)
					p := Repository{
						Owner:         gitlab.String(project.Owner.Username),
						ID:            &id,
						Name:          gitlab.String(project.Name),
						FullName:      gitlab.String(project.NameWithNamespace),
						CloneURL:      gitlab.String(project.HTTPURLToRepo),
						URL:           gitlab.String(project.WebURL),
						DefaultBranch: gitlab.String(project.DefaultBranch),
						Description:   gitlab.String(project.Description),
						Homepage:      gitlab.String(project.WebURL),
					}
					mut.Lock()
					allUserProjects = append(allUserProjects, &p)
					mut.Unlock()
				}
			}
			wg.Done()
		}()

		if response.NextPage == 0 {
			break
		}
		listUserProjectsOps.Page = response.NextPage
	}
	wg.Wait()

	return allUserProjects, nil
}

// getGroupProjects will gather the projects associated with a given group
func (c gitlabClient) getGroupProjects(target Owner) ([]*Repository, error) {
	var allGroupProjects []*Repository
	listGroupProjectsOps := &gitlab.ListGroupProjectsOptions{}
	id := strconv.FormatInt(*target.ID, 10)

	var wg sync.WaitGroup
	var mut sync.Mutex

	for {
		projects, response, err := c.apiClient.Groups.ListGroupProjects(id, listGroupProjectsOps)
		if err != nil {
			return nil, err
		}

		wg.Add(1)

		go func() {
			for _, project := range projects {
				//don't capture forks
				if project.ForkedFromProject == nil {
					id := int64(project.ID)
					p := Repository{
						Owner:         gitlab.String(project.Namespace.FullPath),
						ID:            &id,
						Name:          gitlab.String(project.Name),
						FullName:      gitlab.String(project.NameWithNamespace),
						CloneURL:      gitlab.String(project.HTTPURLToRepo),
						URL:           gitlab.String(project.WebURL),
						DefaultBranch: gitlab.String(project.DefaultBranch),
						Description:   gitlab.String(project.Description),
						Homepage:      gitlab.String(project.WebURL),
					}
					mut.Lock()
					allGroupProjects = append(allGroupProjects, &p)
					mut.Unlock()
				}
			}
			wg.Done()
		}()

		if response.NextPage == 0 {
			break
		}
		listGroupProjectsOps.Page = response.NextPage
	}
	wg.Wait()

	return allGroupProjects, nil
}

// GetRepositoriesFromOwner is used gather all the repos associated with the org owner or other user.
// This is only used by the gitlab client. The github client use a github specific function.
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
