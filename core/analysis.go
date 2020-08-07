// Package core represents the core functionality of all commands
package core

import (
	"crypto/sha1"
	"fmt"
	"gitrob/common"
	"gitrob/github"
	"gitrob/gitlab"
	"gitrob/localRepo"
	"gitrob/matching"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// PrintSessionStats will print the performance and sessions stats to stdout at the conclusion of a session scan
func PrintSessionStats(sess *Session) {

	sess.Out.Important("\n--------Results--------\n")
	sess.Out.Important("\n")
	sess.Out.Important("-------Findings------\n")
	sess.Out.Info("Total Findings......: %d\n", sess.Stats.Findings)
	sess.Out.Important("\n")
	sess.Out.Important("--------Files--------\n")
	sess.Out.Info("Total Files.........: %d\n", sess.Stats.Files)
	sess.Out.Info("Files Scanned.......: %d\n", -1) // TODO implement skipping files and tests
	sess.Out.Info("Files Ignored.......: %d\n", -1) // TODO implement skipping files and tests
	sess.Out.Important("\n")
	sess.Out.Important("---------SCM---------\n")
	sess.Out.Info("Orgs................: %d\n", -1) // TODO need to implement
	sess.Out.Info("Users...............: %d\n", -1) // TODO need to implement
	sess.Out.Info("Repos Found.........: %d\n", sess.Stats.Repositories)
	sess.Out.Info("Repos Cloned........: %d\n", -1) // TODO need to implement
	sess.Out.Info("Repos Scanned.......: %d\n", -1) // TODO need to implement
	sess.Out.Info("Commits Scanned.....: %d\n", sess.Stats.Commits)
	sess.Out.Info("Commits Dirty.......: %d\n", -1) // TODO need to implement
	sess.Out.Important("\n")
	sess.Out.Important("-------General-------\n")
	sess.Out.Info("Grover Version......: %s\n", sess.Version)
	sess.Out.Info("Rules Version.......: %d\n", -1) // TODO need to implement
	sess.Out.Info("Elapsed Time........: %s\n\n", time.Since(sess.Stats.StartedAt))
}

// GatherTargets will enumerate github orgs and members and add them to the running target list of a session
func GatherTargets(sess *Session) {
	sess.Stats.Status = StatusGathering
	sess.Out.Important("Gathering targets...\n")

	var targets []string

	switch sess.ScanType {
	case "github":
		targets = sess.GithubTargets
	case "gitlab":
		targets = sess.GitlabTargets
	}

	for _, loginOption := range targets {
		target, err := sess.Client.GetUserOrganization(loginOption)
		if err != nil || target == nil {
			sess.Out.Error(" Error retrieving information on %s: %s\n", loginOption, err)
			continue
		}
		sess.Out.Debug("%s (ID: %d) type: %s\n", *target.Login, *target.ID, *target.Type)
		sess.AddTarget(target)
		if sess.NoExpandOrgs == false && *target.Type == common.TargetTypeOrganization {
			sess.Out.Debug("Gathering members of %s (ID: %d)...\n", *target.Login, *target.ID)
			members, err := sess.Client.GetOrganizationMembers(*target)
			if err != nil {
				sess.Out.Error(" Error retrieving members of %s: %s\n", *target.Login, err)
				continue
			}
			for _, member := range members {
				sess.Out.Debug("Adding organization member %s (ID: %d) to targets\n", *member.Login, *member.ID)
				sess.AddTarget(member)
			}
		}
	}
}

// GatherLocalRepositories will grab all the local repos from the user input and generate a repository
// object, putting dummy or generated values in where necessary
func GatherLocalRepositories(sess *Session) {

	for _, pth := range sess.RepoDirs {

		if !common.PathExists(pth) {
			sess.Out.Error("\n[*] <%s> does not exist! Quitting.\n", pth)
			os.Exit(1)
		}

		// Gather all paths in the tree
		err0 := filepath.Walk(pth, func(path string, f os.FileInfo, err1 error) error {
			if err1 != nil {
				fmt.Println(err1) // TODO use the error logging capability here
				return nil
			}

			// If it is a directory then move forward
			if f.IsDir() {

				// If there is a .git directory then we have a repo
				if filepath.Ext(path) == ".git" { // TODO Should we reverse this to ! to make the code cleaner

					parent, _ := filepath.Split(path)

					gitProjName, _ := filepath.Split(parent)

					openRepo, err2 := git.PlainOpen(parent)
					if err2 != nil {

						return nil
					}

					ref, err3 := openRepo.Head()
					if err3 != nil {
						fmt.Println("err3: ", err3)
						return nil
					}

					// Get the name of the branch we are working on
					s := ref.Strings()
					branchPath := fmt.Sprintf("%s", s[0])
					branchPathParts := strings.Split(branchPath, string("refs/heads/"))
					branchName := branchPathParts[len(branchPathParts)-1]
					pBranchName := &branchName

					commit, _ := openRepo.CommitObject(ref.Hash())
					var commitHash = commit.Hash[:]

					// TODO make this a generic function at some point
					// Generate a uid for the repo
					h := sha1.New()
					repoID := fmt.Sprintf("%x", h.Sum(commitHash))

					intRepoID, _ := strconv.ParseInt(repoID, 10, 64)
					var pRepoID *int64
					pRepoID = &intRepoID

					// Set the url to the relative path of the repo based on the execution path of grover
					pRepoURL := &parent

					// This is used to id the owner, fullname, and description of the repo. It is ugly but effective. It is the relative path to the repo, for example ../foo
					pGitProjName := &gitProjName

					// The project name is simply the parent directory in the case of a local scan with all other path bits removed for example ../foo -> foo.
					projectPathParts := strings.Split(*pGitProjName, string(os.PathSeparator))
					pProjectName := &projectPathParts[len(projectPathParts)-2]

					sessR := common.Repository{
						Owner:         pGitProjName,
						ID:            pRepoID,
						Name:          pProjectName,
						FullName:      pGitProjName,
						CloneURL:      pRepoURL,
						URL:           pRepoURL,
						DefaultBranch: pBranchName,
						Description:   pGitProjName,
						Homepage:      pRepoURL,
					}

					// Add the repo to the sess to be cloned and scanned
					sess.AddRepository(&sessR)

					sess.Stats.IncrementTargets()
					fmt.Println(len(sess.Repositories))
					fmt.Println("here")
				}
			}
			return nil
		})
		if err0 != nil {
			fmt.Println("err0", err0)
		}
	}
}

// Gather Repositories will gather all repositories associated with a given target during a scan session.
// This is done using threads, whose count is set via commandline flag. Care much be taken to avoid rate
// limiting associated with suspected DOS attacks.
func GatherRepositories(sess *Session) {
	var ch = make(chan *common.Owner, len(sess.Targets))
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
				sess.Stats.IncrementTargets()
				sess.Out.Info(" Retrieved %d %s from %s\n", len(repos), common.Pluralize(len(repos), "repository", "repositories"), *target.Login)
			}
		}()
	}

	for _, target := range sess.Targets {
		ch <- target
	}
	close(ch)
	wg.Wait()
}

// createFinding will create a discrete finding based on a match in a given repo and given commit
func createFinding(repo common.Repository,
	commit object.Commit,
	change *object.Change,
	fileSignature matching.FileSignature,
	contentSignature matching.ContentSignature,
	scanType string) *matching.Finding {

	finding := &matching.Finding{
		FilePath:                    common.GetChangePath(change),
		Action:                      common.GetChangeAction(change),
		FileSignatureDescription:    fileSignature.GetDescription(),
		FileSignatureComment:        fileSignature.GetComment(),
		ContentSignatureDescription: contentSignature.GetDescription(),
		ContentSignatureComment:     contentSignature.GetComment(),
		RepositoryOwner:             *repo.Owner,
		RepositoryName:              *repo.Name,
		CommitHash:                  commit.Hash.String(),
		CommitMessage:               strings.TrimSpace(commit.Message),
		CommitAuthor:                commit.Author.String(),
		CloneUrl:                    *repo.CloneURL,
	}
	finding.Initialize(scanType)
	return finding

}

// matchContent will attempt to match the content of a file such as a password or access token within the file
func matchContent(sess *Session,
	matchTarget matching.MatchTarget,
	repo common.Repository,
	change *object.Change,
	commit object.Commit,
	fileSignature matching.FileSignature,
	threadId int) {

	content, err := common.GetChangeContent(change)
	if err != nil {
		sess.Out.Error("Error retrieving content in commit %s, change %s:  %s", commit.String(), change.String(), err)
	}
	matchTarget.Content = content
	sess.Out.Debug("[THREAD #%d][%s] Matching content in %s...\n", threadId, *repo.CloneURL, commit.Hash)
	for _, contentSignature := range sess.Signatures.ContentSignatures {
		matched, err := contentSignature.Match(matchTarget)
		if err != nil {
			sess.Out.Error("Error while performing content match with '%s': %s\n", contentSignature.Description, err)
		}
		if !matched {
			continue
		}
		finding := createFinding(repo, commit, change, fileSignature, contentSignature, sess.ScanType)
		sess.AddFinding(finding)
	}
}

// findSecrets will attempt to find secrets in the content of a given file or match file in a given path
func findSecrets(sess *Session, repo *common.Repository, commit *object.Commit, changes object.Changes, threadId int) {
	for _, change := range changes {

		path := common.GetChangePath(change)
		matchTarget := matching.NewMatchTarget(path)
		if matchTarget.IsSkippable(sess.SkippablePath, sess.SkippableExt) {
			sess.Out.Debug("[THREAD #%d][%s] Skipping %s\n", threadId, *repo.CloneURL, matchTarget.Path)
			continue
		}
		sess.Out.Debug("[THREAD #%d][%s] Inspecting file: %s...\n", threadId, *repo.CloneURL, matchTarget.Path)

		if sess.Mode != 3 {
			for _, fileSignature := range sess.Signatures.FileSignatures {
				matched, err := fileSignature.Match(matchTarget)
				if err != nil {
					sess.Out.Error(fmt.Sprintf("Error while performing file match: %s\n", err))
				}
				if !matched {
					continue
				}
				if sess.Mode == 1 {
					finding := createFinding(*repo, *commit, change, fileSignature,
						matching.ContentSignature{Description: "NA"}, sess.ScanType)
					sess.AddFinding(finding)
				}
				if sess.Mode == 2 {
					matchContent(sess, matchTarget, *repo, change, *commit, fileSignature, threadId)
				}
				break
			}
			sess.Stats.IncrementFiles()
		} else {
			matchContent(sess, matchTarget, *repo, change, *commit, matching.FileSignature{Description: "NA"}, threadId)
			sess.Stats.IncrementFiles()
		}
	}
}

// cloneRepository will clone a given repository based upon a configured set or options a user provides
func cloneRepository(sess *Session, repo *common.Repository, threadId int) (*git.Repository, string, error) {
	sess.Out.Debug("[THREAD #%d][%s] Cloning repository...\n", threadId, *repo.CloneURL)

	cloneConfig := common.CloneConfiguration{
		Url:        repo.CloneURL,
		Branch:     repo.DefaultBranch,
		Depth:      &sess.CommitDepth,
		Token:      &sess.GitlabAccessToken, // TODO Is this need since we already have a client?
		InMemClone: &sess.InMemClone,
	}

	var clone *git.Repository
	var path string
	var err error

	switch sess.ScanType {
	case "github":
		clone, path, err = github.CloneRepository(&cloneConfig)
	case "gitlab":
		userName := "oauth2"
		cloneConfig.Username = &userName
		clone, path, err = gitlab.CloneRepository(&cloneConfig)
	case "localGit":
		clone, path, err = localRepo.CloneRepository(&cloneConfig)

	}
	if err != nil {
		if err.Error() != "remote repository is empty" {
			sess.Out.Error("Error cloning repository %s: %s\n", *repo.CloneURL, err)
		}
		sess.Stats.IncrementRepositories()
		sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))
		return nil, "", err
	}
	sess.Out.Debug("[THREAD #%d][%s] Cloned repository to: %s\n", threadId, *repo.CloneURL, path)
	return clone, path, err
}

// getRepositoryHistory will attempt to get the commit history of a given repo and if successful increment the repo
// count and update the progress.
func getRepositoryHistory(sess *Session, clone *git.Repository, repo *common.Repository, path string, threadId int) ([]*object.Commit, error) {
	history, err := common.GetRepositoryHistory(clone)
	if err != nil {
		sess.Out.Error("[THREAD #%d][%s] Error getting commit history: %s\n", threadId, *repo.CloneURL, err)
		if sess.InMemClone {
			os.RemoveAll(path)
		}
		sess.Stats.IncrementRepositories()
		sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))
		return nil, err
	}
	sess.Out.Debug("[THREAD #%d][%s] Number of commits: %d\n", threadId, *repo.CloneURL, len(history))
	return history, err
}

// AnalyzeRepositories will take a given repository, clone it, pull the commit history and use that as a basis for
// scanning for secrets within the repo and based on that output create a finding associated with that repo
func AnalyzeRepositories(sess *Session) {
	sess.Stats.Status = StatusAnalyzing
	var ch = make(chan *common.Repository, len(sess.Repositories))
	var wg sync.WaitGroup
	var threadNum int
	if len(sess.Repositories) <= 1 {
		threadNum = 1
	} else if len(sess.Repositories) <= sess.Threads {
		threadNum = len(sess.Repositories) - 1
	} else {
		threadNum = sess.Threads
	}
	wg.Add(threadNum)
	sess.Out.Debug("Threads for repository analysis: %d\n", threadNum)

	sess.Out.Important("Analyzing %d %s...\n", len(sess.Repositories), common.Pluralize(len(sess.Repositories), "repository", "repositories"))

	for i := 0; i < threadNum; i++ {
		go func(tid int) {
			for {
				sess.Out.Debug("[THREAD #%d] Requesting new repository to analyze...\n", tid)
				repo, ok := <-ch
				if !ok {
					sess.Out.Debug("[THREAD #%d] No more tasks, marking WaitGroup as done\n", tid)
					wg.Done()
					return
				}

				clone, path, err := cloneRepository(sess, repo, tid)
				if err != nil {
					continue
				}

				history, err := getRepositoryHistory(sess, clone, repo, path, tid)
				if err != nil {
					continue
				}

				for _, commit := range history {
					sess.Out.Debug("[THREAD #%d][%s] Analyzing commit: %s\n", tid, *repo.CloneURL, commit.Hash)
					changes, _ := common.GetChanges(commit, clone)
					sess.Out.Debug("[THREAD #%d][%s] %s changes in %d\n", tid, *repo.CloneURL, commit.Hash, len(changes))

					findSecrets(sess, repo, commit, changes, tid)

					sess.Stats.IncrementCommits()
					sess.Out.Debug("[THREAD #%d][%s] Done analyzing changes in %s\n", tid, *repo.CloneURL, commit.Hash)
				}

				sess.Out.Debug("[THREAD #%d][%s] Done analyzing commits\n", tid, *repo.CloneURL)
				if sess.InMemClone {
					os.RemoveAll(path)
				}
				sess.Out.Debug("[THREAD #%d][%s] Deleted %s\n", tid, *repo.CloneURL, path)
				sess.Stats.IncrementRepositories()
				sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))
			}
		}(i)
	}
	for _, repo := range sess.Repositories {
		ch <- repo
	}
	close(ch)
	wg.Wait()
}
