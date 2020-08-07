// Package gitlab represents github specific functionality
package gitlab

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"os"
	"regexp"
	"strconv"
	"strings"
	"wraith/common"
)

// Client holds a gitlab api client instance
type Client struct {
	apiClient *gitlab.Client
	logger    *common.Logger
}

// NewClient creates a gitlab api client instance using a token
func (c Client) NewClient(token string, logger *common.Logger) (Client, error) {
	var err error
	c.apiClient, err = gitlab.NewClient(token)
	if err != nil {
		return Client{}, err
	}
	c.apiClient.UserAgent = common.UserAgent
	c.logger = logger
	return c, nil
}

// CheckAPIToken will ensure we have a valid github api token
func CheckAPIToken(t string) {

	// check to make sure the length is proper
	if len(t) != 20 {
		fmt.Println("The token is invalid. Please use a valid Gitlab token")
		os.Exit(2)
	}

	// match only letters and numbers and ensure you match 40
	exp1 := regexp.MustCompile(`^[A-Za-z0-9]{20}`)
	if !exp1.MatchString(t) {
		fmt.Println("The token is invalid. Please use a valid Gitlab token")
		os.Exit(2)
	}
	//return t
}

// GetUserOrganization is used to enumerate the owner in a given org
func (c Client) GetUserOrganization(login string) (*common.Owner, error) {
	emptyString := gitlab.String("")
	org, orgErr := c.getOrganization(login)
	if orgErr != nil {
		user, userErr := c.getUser(login)
		if userErr != nil {
			return nil, userErr
		}
		id := int64(user.ID)
		return &common.Owner{
			Login:     gitlab.String(user.Username),
			ID:        &id,
			Type:      gitlab.String(common.TargetTypeUser),
			Name:      gitlab.String(user.Name),
			AvatarURL: gitlab.String(user.AvatarURL),
			URL:       gitlab.String(user.WebsiteURL),
			Company:   gitlab.String(user.Organization),
			Blog:      emptyString,
			Location:  emptyString,
			Email:     gitlab.String(user.PublicEmail),
			Bio:       gitlab.String(user.Bio),
		}, nil
	} else {
		id := int64(org.ID)
		return &common.Owner{
			Login:     gitlab.String(org.Name),
			ID:        &id,
			Type:      gitlab.String(common.TargetTypeOrganization),
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
}

// GetOrganizationMembers will gather all the members of a given organization
func (c Client) GetOrganizationMembers(target common.Owner) ([]*common.Owner, error) {
	var allMembers []*common.Owner
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
				&common.Owner{
					Login: gitlab.String(member.Username),
					ID:    &id,
					Type:  gitlab.String(common.TargetTypeUser)})
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allMembers, nil
}

// GetRepositoriesFromOwner is used gather all the repos associated with the org owner or other user
func (c Client) GetRepositoriesFromOwner(target common.Owner) ([]*common.Repository, error) {
	var allProjects []*common.Repository
	id := int(*target.ID)
	if *target.Type == common.TargetTypeUser {
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
func (c Client) getUser(login string) (*gitlab.User, error) {
	users, _, err := c.apiClient.Users.ListUsers(&gitlab.ListUsersOptions{Username: gitlab.String(login)})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("No GitLab %s or %s %s was found.  If you are targeting a GitLab group, be sure to"+
			" use an ID in place of a name.",
			strings.ToLower(common.TargetTypeUser),
			strings.ToLower(common.TargetTypeOrganization),
			login)
	}
	return users[0], err
}

// getOrganization will get the necessary info from an org
func (c Client) getOrganization(login string) (*gitlab.Group, error) {
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
func (c Client) getUserProjects(id int) ([]*common.Repository, error) {
	var allUserProjects []*common.Repository
	listUserProjectsOps := &gitlab.ListProjectsOptions{}
	for {
		projects, response, err := c.apiClient.Projects.ListUserProjects(id, listUserProjectsOps)
		if err != nil {
			return nil, err
		}
		for _, project := range projects {
			//don't capture forks
			if project.ForkedFromProject == nil {
				id := int64(project.ID)
				p := common.Repository{
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
				allUserProjects = append(allUserProjects, &p)
			}
		}
		if response.NextPage == 0 {
			break
		}
		listUserProjectsOps.Page = response.NextPage
	}
	return allUserProjects, nil
}

// getGroupProjects will gather the projects associated with a given group
func (c Client) getGroupProjects(target common.Owner) ([]*common.Repository, error) {
	var allGroupProjects []*common.Repository
	listGroupProjectsOps := &gitlab.ListGroupProjectsOptions{}
	id := strconv.FormatInt(*target.ID, 10)
	for {
		projects, response, err := c.apiClient.Groups.ListGroupProjects(id, listGroupProjectsOps)
		if err != nil {
			return nil, err
		}
		for _, project := range projects {
			//don't capture forks
			if project.ForkedFromProject == nil {
				id := int64(project.ID)
				p := common.Repository{
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
				allGroupProjects = append(allGroupProjects, &p)
			}
		}
		if response.NextPage == 0 {
			break
		}
		listGroupProjectsOps.Page = response.NextPage
	}
	return allGroupProjects, nil
}
