package core

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"github.com/google/go-github/github"
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

// TODO can we clean this up at all
// validateGHEInput will ensure that the user supplied input for orgs and repos is valid and not malicious.
func ValidateGHInput(s *Session) {
	// Ex. bar, 3ar, bar_foo, baz-foo
	// This regex is case insensitive
	exp1 := regexp.MustCompile(`[A-Za-z0-9,-_]*$`)

	if !exp1.MatchString(strings.TrimSpace(s.UserDirtyOrgs)) {
		fmt.Println("The orgs are in an unsupported format. Please use a coma separated list with no whitespace or special characters")
		os.Exit(2)
	}
	if !exp1.MatchString(strings.TrimSpace(s.UserDirtyRepos)) {
		fmt.Println("The repos are in an unsupported format. Please use a coma separated list with no whitespace or special characters")
		os.Exit(2)
	}
	if !exp1.MatchString(strings.TrimSpace(s.UserDirtyNames)) {
		fmt.Println("The users are in an unsupported format. Please use a coma separated list with no whitespace or special characters")
		os.Exit(2)
	}

	orgs := strings.Split(s.UserDirtyOrgs, ",")
	repos := strings.Split(s.UserDirtyRepos, ",")
	users := strings.Split(s.UserDirtyNames, ",")

	for _, o := range orgs {
		if o != "" && o != "." {
			s.UserOrgs = append(s.UserOrgs, o)
		}
	}

	for _, r := range repos {
		if r != "" && r != "." {
			s.UserRepos = append(s.UserRepos, r)
		}
	}

	for _, u := range users {
		if u != "" && u != "." {
			s.UserLogins = append(s.UserLogins, u)
		}
	}
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

	var wg sync.WaitGroup
	var mut sync.Mutex

	for {
		repos, resp, err := c.apiClient.Repositories.List(ctx, *target.Login, opt)
		if err != nil {
			return allRepos, err
		}

		wg.Add(1)

		go func() {
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
					mut.Lock()
					allRepos = append(allRepos, &r)
					mut.Unlock()
				}
			}
			wg.Done()
		}()

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}
	wg.Wait()

	return allRepos, nil
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

	// List of orgs that we find
	var orgList []*github.Organization

	var orgID int64

	// Options necessary for enumerating the orgs. These are not client options such as TLS or auth,
	// these are options such as orgs per page.
	var opts github.OrganizationsListOptions

	// How many orgs per page
	opts.PerPage = 40

	// this controls pagination, see below. In order for it to work this gets set to the last org that was found
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

			if err != nil { // TODO add better error checking here
				fmt.Println(err)
				os.Exit(99)
			}

			for _, org := range orgs {
				orgList = append(orgList, org)
				orgID = *org.ID
			}

			opts.Since = orgID
			tmpOrgID = orgID + 1
		}
	} else {
		// This will handle user specificed orgs //TODO can we do something better here
		for _, o := range sess.UserOrgs {
			org, _, err := sess.GithubClient.Organizations.Get(ctx, o)

			if err != nil {
				fmt.Println(err)
				os.Exit(99)
			}

			orgList = append(orgList, org)
		}
	}
	for _, org := range orgList {

		// Add the orgs to the list for later enumeration of repos
		sess.addOrganization(org)
	}

	//fmt.Println(sess.Organizations) // TODO remove me

	// Set the total count of orgs scanned
	sess.Stats.Organizations = len(orgList)
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

// GatherGithubRepositories will walk a tree and create a repo object for each repository found. After the
// object is completed is will increment the total number of repositories by 1.
func GatherGithubRepositories(sess *Session) {

	// Create a channel for each org in the list
	var ch = make(chan *github.Organization, len(sess.Organizations))
	var wg sync.WaitGroup
	//var orgName *string

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

	fmt.Println("I am here") // TODO Remove me

	wg.Add(threadNum)
	sess.Out.Debug("Threads for repository gathering: %d\n", threadNum)

	for i := 0; i < threadNum; i++ {
		go func() {
			for {
				fmt.Println("I am in the first for loop")
				var repos []*Repository
				var err error
				org, ok := <-ch
				if !ok {
					wg.Done()
					return
				}
				// Retrieve all the repos in an org regardless of public/private
				repos, err = getRepositoriesFromOrganization(org.Login, sess.GithubClient, sess.ScanFork)
				fmt.Println("I am after we gather the org repos")

				if err != nil {
					sess.Out.Error(" Failed to retrieve repositories from %s: %s\n", *org.Login, err)
				}

				// In the case where all the repos are private
				if len(repos) == 0 {
					fmt.Println("I am at zero")
					continue
				}
				// This is for all the repos we could see in a given org.
				for _, repo := range repos {
					if len(sess.UserRepos) >= 1 && sess.UserRepos[0] != "" {
						for _, r := range sess.UserRepos {
							if r == *repo.Name {

								sess.Out.Debug(" Retrieved repository: %s\n", *repo.FullName)
								// Add the repo to the sess to be scanned

								sess.AddRepository(repo)
								// Increment the total count of repos found, regardless if it gets cloned or scanned
								sess.Stats.IncrementRepositoriesTotal()
							}
						}
						continue
					}
					sess.Out.Debug(" Retrieved repository: %s\n", *repo.FullName)
					// Add the repo to the sess to be scanned

					sess.AddRepository(repo)
					// Increment the total count of repos found, regardless if it gets cloned or scanned
					sess.Stats.IncrementRepositoriesTotal()
				}
			}
		}()
	}
	for _, org := range sess.Organizations {
		ch <- org
	}

	//for idx, gheRepo := range sess.Organizations {
	//	// This associates an org with all the repos under it
	//	orgName = sess.Organizations[idx].Login
	//	ch <- gheRepo
	//}
	close(ch)
	wg.Wait()
}
