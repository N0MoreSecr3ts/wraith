package core

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"sync"
)

// cloneGithub will set the clone config and then either do a plain clone if it is going to disk
// or a full clone if going ito memory.
func cloneGithub(cloneConfig *CloneConfiguration) (*git.Repository, string, error) {

	cloneOptions := &git.CloneOptions{
		URL:           *cloneConfig.URL,
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
	apiClient *github.Client
}

// addUser will add a new user to the sess for further scanning and analyzing
func (s *Session) addUser(user *github.User) {
	s.Lock()
	defer s.Unlock()
	h := md5.New()
	_, _ = io.WriteString(h, *user.Login)                     // TODO handle error
	_, _ = io.WriteString(h, strconv.FormatInt(*user.ID, 10)) // TODO handle error
	userMD5 := fmt.Sprintf("%x", h.Sum(nil))

	for _, o := range s.GithubUsers {
		j := md5.New()
		_, _ = io.WriteString(j, *o.Login)                     // TODO handle error
		_, _ = io.WriteString(h, strconv.FormatInt(*o.ID, 10)) // TODO handle error
		sessMD5 := fmt.Sprintf("%x", h.Sum(nil))

		if userMD5 == sessMD5 {
			return
		}
	}
	s.GithubUsers = append(s.GithubUsers, user)
}

// GatherUsers will generate a list of users from github.com that can then be filtered down to a specific target range
func GatherUsers(sess *Session) {
	sess.Out.Important("Gathering users...\n")

	ctx := context.Background()

	var opts github.OrganizationsListOptions
	opts.PerPage = 40
	opts.Since = -1
	for _, o := range sess.UserLogins {
		u, _, err := sess.GithubClient.Users.Get(ctx, o)

		if err != nil {
			sess.Out.Error("Unable to collect user %s: %s\n", u, err)
		}

		// Add the user to the session and increment the user count
		sess.addUser(u)
		sess.Stats.IncrementUsers()
		sess.Out.Debug("Added user %s\n", *u.Login)
	}
}

// ValidateUserInput will check for special characters in the strings and make sure we
// have at least one usr/repo/org to scan
func (s *Session) ValidateUserInput(v *viper.Viper) {

	// Raw user inputs
	s.UserDirtyRepos = v.GetStringSlice("github-repos")
	s.UserDirtyOrgs = v.GetStringSlice("github-orgs")
	s.UserDirtyNames = v.GetStringSlice("github-users")
	s.GithubAccessToken = CheckGithubAPIToken(v.GetString("github-api-token"), s)

	// If no targets are given, fail fast
	if s.UserDirtyRepos == nil && s.UserDirtyOrgs == nil && s.UserDirtyNames == nil {
		s.Out.Error("You must enter either a user, org or repo[s] to scan")
	}

	// validate the input does not contain any scary characters
	exp := regexp.MustCompile(`[A-Za-z0-9,-_]*$`)

	for _, o := range s.UserDirtyOrgs {
		if exp.MatchString(o) {
			s.UserOrgs = append(s.UserOrgs, o)

		}
	}

	for _, r := range s.UserDirtyRepos {
		if exp.MatchString(r) {
			s.UserRepos = append(s.UserRepos, r)

		}
	}

	for _, u := range s.UserDirtyNames {
		if exp.MatchString(u) {
			s.UserLogins = append(s.UserLogins, u)
		}
	}

}

// CheckGithubAPIToken will ensure we have a valid github api token
func CheckGithubAPIToken(t string, sess *Session) string {

	// check to make sure the length is proper
	if len(t) != 40 {
		sess.Out.Error("The token is invalid. Please use a valid Github token\n")
		os.Exit(2)
	}

	// match only letters and numbers and ensure you match 40
	exp1 := regexp.MustCompile(`[A-Za-z0-9\_]{40}`)
	if !exp1.MatchString(t) {
		sess.Out.Error("The token is invalid. Please use a valid Github token\n")
		os.Exit(2)
	}
	return t
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

// TODO Do we thread this?
// getRepositoriesFromOrganization will generate a slice of github repo objects for an org. This has only been tested on github enterprise.
func getRepositoriesFromOrganization(login *string, client *github.Client, scanFork bool, sess *Session) ([]*Repository, error) {
	var allRepos []*Repository
	orgName := *login
	ctx := context.Background()
	opt := &github.RepositoryListByOrgOptions{
		Type: "sources",
	}

	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, orgName, opt)
		if err != nil {
			sess.Out.Error("Error listing repos for the org %s: %s\n", orgName, err)
			return allRepos, err
		}
		for _, repo := range repos {
			// If we don't want to scan forked repos, we can use a flag to set this and the
			// loop instance will stop here and go on to the next repo
			if !sess.ScanFork && repo.GetFork() {
				continue
			}
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

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}

// GatherGithubRepositoriesFromOwner is used gather all the repos associated with a github user
func GatherGithubRepositoriesFromOwner(sess *Session) {
	var allRepos []*Repository
	ctx := context.Background()

	// The defaults should be fine for a tool like this but if you want to customize
	// settings like repo type (public, private, etc) or the amount of results returned
	// per page this is where you do it.
	opt := &github.RepositoryListOptions{}

	// TODO This should be threaded
	for _, ul := range sess.UserLogins {
		for {
			repos, resp, err := sess.GithubClient.Repositories.List(ctx, ul, opt)
			if err != nil {
				sess.Out.Error("Error gathering Github repos from %s: %s\n", ul, err)
			}
			for _, repo := range repos {
				// If we don't want to scan forked repos, we can use a flag to set this and the
				// loop instance will stop here and go on to the next repo
				if !sess.ScanFork && repo.GetFork() {
					continue
				}
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
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
	}

	// FIXME what happens if no repos are recovered

	// If we re only looking for a subset of the repos in an org we do a comparison
	// of the repos gathered for the org and the list of repos that we care about.
	for _, repo := range allRepos {
		// Increment the total number of repos found, regardless if we are cloning them
		sess.Stats.IncrementRepositoriesTotal()
		if sess.UserRepos != nil {
			for _, r := range sess.UserRepos {
				if r == *repo.Name {
					sess.Out.Debug(" Retrieved repository %s from user %s\n", *repo.FullName, *repo.Owner)

					// Add the repo to the sess to be scanned
					sess.AddRepository(repo)
				}
			}
			continue
		}
		sess.Out.Debug(" Retrieved repository %s from user %s\n", *repo.FullName, *repo.Owner)

		// If we are not doing any filtering and simply grabbing all available repos we add the repos
		// to the session to be scanned
		sess.AddRepository(repo)
	}
}

// GetOrganizationMembers will gather all the members of a given organization
func (c githubClient) GetOrganizationMembers(target Owner) ([]*Owner, error) {
	var allMembers []*Owner
	ctx := context.Background()
	opt := &github.ListMembersOptions{}

	var wg sync.WaitGroup
	var mut sync.Mutex

	for {
		members, resp, err := c.apiClient.Organizations.ListMembers(ctx, *target.Login, opt)
		if err != nil {
			return allMembers, err
		}

		wg.Add(1)

		go func() {
			for _, member := range members {
				mut.Lock()
				allMembers = append(allMembers, &Owner{Login: member.Login, ID: member.ID, Type: member.Type})
				mut.Unlock()
			}
			wg.Done()
		}()

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	wg.Wait()

	return allMembers, nil
}

// GatherOrgs will use a client to generate a list of all orgs that the client can see. By default this will include
// orgs that contain both public and private repos
func GatherOrgs(sess *Session) {
	sess.Out.Important("Gathering github organizations...\n")

	ctx := context.Background()

	var orgList []*github.Organization

	var orgID int64

	// Options necessary for enumerating the orgs. These are not client options such as TLS or auth,
	// these are options such as orgs per page.
	var opts github.OrganizationsListOptions

	// How many orgs per page
	opts.PerPage = 40

	// This controls pagination, see below. In order for it to work this gets set to the last org that was found
	// and the next page picks up with the next one in line.
	opts.Since = -1

	// Used to track the orgID's for the sake of pagination.
	tmpOrgID := int64(0)

	// If the user did not specify specific orgs then we grab all the orgs the client can see
	if sess.UserOrgs == nil {

		for opts.Since < tmpOrgID {
			if opts.Since == orgID {
				break
			}
			orgs, _, err := sess.GithubClient.Organizations.ListAll(ctx, &opts)

			if err != nil {
				sess.Out.Error("Error gathering Github orgs: %s\n", err)
			}

			for _, org := range orgs {
				orgList = append(orgList, org)
				orgID = *org.ID
			}

			opts.Since = orgID
			tmpOrgID = orgID + 1
		}
	} else {
		// This will handle orgs passed in via flags
		for _, o := range sess.UserOrgs {
			org, _, err := sess.GithubClient.Organizations.Get(ctx, o)

			if err != nil {
				sess.Out.Error("Error gathering the Github org %s: %s\n", o, err)
			}

			orgList = append(orgList, org)
		}
	}

	// Add the orgs to the list for later enumeration of repos
	for _, org := range orgList {
		sess.addOrganization(org)
		sess.Stats.IncrementOrgs()
		sess.Out.Debug("Added org %s\n", *org.Login)

	}
}

// addOrganization will add a new organization to the session for further scanning and analyzing
func (s *Session) addOrganization(organization *github.Organization) {
	s.Lock()
	defer s.Unlock()
	h := md5.New()
	_, _ = io.WriteString(h, *organization.Login) // TODO handle these errors instead of ignoring them explictly
	_, _ = io.WriteString(h, strconv.FormatInt(*organization.ID, 10))
	orgMD5 := fmt.Sprintf("%x", h.Sum(nil))

	for _, o := range s.Organizations {
		j := md5.New()
		_, _ = io.WriteString(j, *o.Login)
		_, _ = io.WriteString(h, strconv.FormatInt(*o.ID, 10))
		sessMD5 := fmt.Sprintf("%x", h.Sum(nil))

		if orgMD5 == sessMD5 {
			return
		}
	}
	s.Organizations = append(s.Organizations, organization)
}

// GatherGithubOrgRepositories will gather all the repositories for a given org.
func GatherGithubOrgRepositories(sess *Session) {

	// Create a channel for each org in the list
	var ch = make(chan *github.Organization, len(sess.Organizations))
	var wg sync.WaitGroup

	// Calculate the number of threads based on the flag and the number of orgs
	// TODO: implement nice in the threading logic to guard against rate limiting and tripping the
	//  security protections
	var threadNum int
	if len(sess.Organizations) <= 1 {
		threadNum = 1
	} else if len(sess.Organizations) <= sess.Threads {
		threadNum = len(sess.Organizations) - 1
	} else {
		threadNum = sess.Threads
	}

	wg.Add(threadNum)
	sess.Out.Debug("Threads for repository gathering: %d\n", threadNum)

	for i := 0; i < threadNum; i++ {
		go func() {
			for {
				var repos []*Repository
				var err error
				org, ok := <-ch
				if !ok {
					wg.Done()
					return
				}
				// Retrieve all the repos in an org regardless of public/private
				repos, err = getRepositoriesFromOrganization(org.Login, sess.GithubClient, sess.ScanFork, sess)

				if err != nil {
					sess.Out.Error(" Failed to retrieve repositories from %s: %s\n", *org.Login, err)
				}

				// In the case where all the repos are private
				if len(repos) == 0 {
					sess.Out.Debug("No repositories have benn gathered for %s\n", *org.Login)
					continue
				}

				// If we re only looking for a subset of the repos in an org we do a comparison
				// of the repos gathered for the org and the list pf repos that we care about.
				for _, repo := range repos {

					// Increment the total number of repos found even if we are not cloning them
					sess.Stats.IncrementRepositoriesTotal()

					if sess.UserRepos != nil {
						for _, r := range sess.UserRepos {
							if r == *repo.Name {
								sess.Out.Debug(" Retrieved repository %s from org %s\n", *repo.FullName, *org.Login)

								// Add the repo to the sess to be scanned
								sess.AddRepository(repo)
							}
						}
						continue
					}
					sess.Out.Debug(" Retrieved repository %s from org %s\n", *repo.FullName, *org.Login)

					// If we are not doing any filtering and simply grabbing all available repos we add the repos
					// to the session to be scanned
					sess.AddRepository(repo)
				}
			}
		}()
	}
	for _, org := range sess.Organizations {
		ch <- org
	}

	close(ch)
	wg.Wait()
}
