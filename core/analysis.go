// Package core represents the core functionality of all commands
package core

import (
	"crypto/sha1"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"wraith/version"
)

// PrintSessionStats will print the performance and sessions stats to stdout at the conclusion of a session scan
func PrintSessionStats(sess *Session) {

	sess.Out.Important("\n--------Results--------\n")
	sess.Out.Important("\n")
	sess.Out.Important("-------Findings------\n")
	sess.Out.Info("Total Findings......: %d\n", sess.Stats.Findings)
	sess.Out.Important("\n")
	sess.Out.Important("--------Files--------\n")
	sess.Out.Info("Total Files.........: %d\n", sess.Stats.FilesTotal)
	sess.Out.Info("Files Scanned.......: %d\n", sess.Stats.FilesScanned)
	sess.Out.Info("Files Ignored.......: %d\n", sess.Stats.FilesIgnored)
	sess.Out.Info("Files Dirty.........: %d\n", sess.Stats.FilesDirty)
	sess.Out.Important("\n")
	sess.Out.Important("---------SCM---------\n")
	sess.Out.Info("Repos Found.........: %d\n", sess.Stats.RepositoriesTotal)
	sess.Out.Info("Repos Cloned........: %d\n", sess.Stats.RepositoriesCloned)
	sess.Out.Info("Repos Scanned.......: %d\n", sess.Stats.RepositoriesScanned)
	sess.Out.Info("Commits Scanned.....: %d\n", sess.Stats.Commits)
	sess.Out.Info("Commits Dirty.......: %d\n", sess.Stats.CommitsDirty)
	sess.Out.Important("\n")
	sess.Out.Important("-------General-------\n")
	sess.Out.Info("Wraith Version......: %s\n", sess.Version)
	sess.Out.Info("Rules Version.......: %s\n", sess.RulesVersion)
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
		if sess.NoExpandOrgs == false && *target.Type == TargetTypeOrganization {
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

	// This is the number of targets as we don't do forks or anything else.
	// It will contain directorys, that will then be added to the repo count
	// if they contain a .git directory
	sess.Stats.Targets = len(sess.LocalDirs)

	for _, pth := range sess.LocalDirs {

		if !PathExists(pth) {
			sess.Out.Error("\n[*] <%s> does not exist! Quitting.\n", pth)
			os.Exit(1)
		}

		// Gather all paths in the tree
		err0 := filepath.Walk(pth, func(path string, f os.FileInfo, err1 error) error {
			if err1 != nil {
				//fmt.Println(err1) // TODO use the error logging capability here
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
						//fmt.Println("err3: ", err3) //TODO remove me
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

					// Set the url to the relative path of the repo based on the execution path of wraith
					pRepoURL := &parent

					// This is used to id the owner, fullname, and description of the repo. It is ugly but effective. It is the relative path to the repo, for example ../foo
					pGitProjName := &gitProjName

					// The project name is simply the parent directory in the case of a local scan with all other path bits removed for example ../foo -> foo.
					projectPathParts := strings.Split(*pGitProjName, string(os.PathSeparator))
					pProjectName := &projectPathParts[len(projectPathParts)-2]

					sessR := Repository{
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
				}
			}
			return nil
		})
		if err0 != nil {
		}
	}
}

// Gather Repositories will gather all repositories associated with a given target during a scan session.
// This is done using threads, whose count is set via commandline flag. Care much be taken to avoid rate
// limiting associated with suspected DOS attacks.
func GatherRepositories(sess *Session) {
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

//	sess.Out.Debug("[THREAD #%d][%s] Skipping %s\n", threadId, *repo.CloneURL, matchTarget.Path) // TODO implement me
//
//sess.Out.Debug("[THREAD #%d][%s] Inspecting file: %s...\n", threadId, *repo.CloneURL, matchTarget.Path) // TODO implement me
//
//			sess.Out.Error(fmt.Sprintf("Error while performing file match: %s\n", err))

// cloneRepository will clone a given repository based upon a configured set or options a user provides
func cloneRepository(sess *Session, repo *Repository, threadId int) (*git.Repository, string, error) {
	sess.Out.Debug("[THREAD #%d][%s] Cloning repository...\n", threadId, *repo.CloneURL)

	var clone *git.Repository
	var path string
	var err error

	switch sess.ScanType {
	case "github":
		cloneConfig := CloneConfiguration{
			Url:        repo.CloneURL,
			Branch:     repo.DefaultBranch,
			Depth:      &sess.CommitDepth,
			InMemClone: &sess.InMemClone,
		}
		clone, path, err = CloneGithubRepository(&cloneConfig)
	case "gitlab":
		userName := "oauth2"
		cloneConfig := CloneConfiguration{
			Url:        repo.CloneURL,
			Branch:     repo.DefaultBranch,
			Depth:      &sess.CommitDepth,
			Token:      &sess.GitlabAccessToken, // TODO Is this need since we already have a client?
			InMemClone: &sess.InMemClone,
			Username:   &userName,
		}
		clone, path, err = CloneGitlabRepository(&cloneConfig)
	case "localGit":
		cloneConfig := CloneConfiguration{
			Url:        repo.CloneURL,
			Branch:     repo.DefaultBranch,
			Depth:      &sess.CommitDepth,
			InMemClone: &sess.InMemClone,
		}
		clone, path, err = CloneLocalRepository(&cloneConfig)

	}
	if err != nil {
		switch err.Error() {
		case "remote repository is empty":
			sess.Out.Error("Repository %s is empty: %s\n", *repo.CloneURL, err)
			sess.Stats.IncrementRepositoriesCloned()
			//sess.Stats.UpdateProgress(sess.Stats.RepositoriesCloned, len(sess.Repositories))
			return nil, "", err
		default:
			sess.Out.Error("Error cloning repository %s: %s\n", *repo.CloneURL, err)
			//sess.Stats.UpdateProgress(sess.Stats.RepositoriesCloned, len(sess.Repositories))
			return nil, "", err
		}
	}
	sess.Stats.IncrementRepositoriesCloned()
	//sess.Stats.UpdateProgress(sess.Stats.RepositoriesCloned, len(sess.Repositories))
	sess.Out.Debug("[THREAD #%d][%s] Cloned repository to: %s\n", threadId, *repo.CloneURL, path)
	return clone, path, err
}

// getRepositoryHistory will attempt to get the commit history of a given repo and if successful increment the repo
// count and update the progress.
func getRepositoryHistory(sess *Session, clone *git.Repository, repo *Repository, path string, threadId int) ([]*object.Commit, error) {
	history, err := GetRepositoryHistory(clone)
	if err != nil {
		sess.Out.Error("[THREAD #%d][%s] Error getting commit history: %s\n", threadId, *repo.CloneURL, err)
		if sess.InMemClone {
			os.RemoveAll(path)
		}
		//sess.Stats.IncrementRepositories()
		//sess.Stats.UpdateProgress(sess.Stats.RepositoriesCloned, len(sess.Repositories))
		return nil, err
	}
	sess.Out.Debug("[THREAD #%d][%s] Number of commits: %d\n", threadId, *repo.CloneURL, len(history))
	return history, err
}

//sess.Out.Debug("Threads for repository analysis: %d\n", threadNum)
//sess.Out.Important("Analyzing %d %s...\n", len(sess.Repositories), Pluralize(len(sess.Repositories), "repository", "repositories"))
//				sess.Out.Debug("[THREAD #%d] No more tasks, marking WaitGroup as done\n", tid)

//					sess.Out.Debug("[THREAD #%d][%s] Analyzing commit: %s\n", tid, *repo.CloneURL, commit.Hash)
//					sess.Out.Debug("[THREAD #%d][%s] %s changes in %d\n", tid, *repo.CloneURL, commit.Hash, len(changes))
//
//					sess.Out.Debug("[THREAD #%d][%s] Done analyzing changes in %s\n", tid, *repo.CloneURL, commit.Hash)
//
//				sess.Out.Debug("[THREAD #%d][%s] Done analyzing commits\n", tid, *repo.CloneURL)
//				sess.Out.Debug("[THREAD #%d][%s] Deleted %s\n", tid, *repo.CloneURL, path)

func AnalyzeRepositories(sess *Session) {
	sess.Stats.Status = StatusAnalyzing
	if len(sess.Repositories) == 0 {
		sess.Out.Error("No repositories have been gathered.")
		os.Exit(2)
	}

	var ch = make(chan *Repository, len(sess.Repositories))
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

	sess.Out.Important("Analyzing %d %s...\n", len(sess.Repositories), Pluralize(len(sess.Repositories), "repository", "repositories"))

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

				// Clone the repository from the remote source or if local from the path
				clone, path, err := cloneRepository(sess, repo, tid)
				if err != nil {
					if err.Error() != "remote repository is empty" {
						sess.Out.Error("Error cloning repository %s: %s\n", *repo.FullName, err)
					}
					continue
				}

				// Get the commit history for the repo
				history, err := GetRepositoryHistory(clone)
				if err != nil {
					sess.Out.Error("[THREAD #%d][%s] Error getting commit history: %s\n", tid, *repo.FullName, err)
					os.RemoveAll(path)
					continue
				}

				for _, commit := range history {
					sess.Out.Debug("[THREAD #%d][%s] Analyzing commit: %s\n", tid, *repo.CloneURL, commit.Hash)

					// Increment the total number of commits scanned
					sess.Stats.IncrementCommits()
					//sess.Stats.IncrementCommitsScanned() // TODO implement in stats

					// This will be used to increment the dirty commit stat if any matches are found
					dirtyCommit := false

					changes, _ := GetChanges(commit, clone)
					sess.Out.Debug("[THREAD #%d][%s] %s changes in %d\n", tid, *repo.CloneURL, commit.Hash, len(changes))

					for _, change := range changes {

						changeAction := GetChangeAction(change)
						fPath := GetChangePath(change)
						fullFilePath := path + "/" + fPath

						sess.Stats.IncrementFilesTotal()

						likelyTestFile := false

						if !sess.ScanTests {
							likelyTestFile = isTestFileOrPath(fullFilePath)
						}

						// If the file is likely a test then ignore it
						if likelyTestFile {
							// If we are not scanning the file then by definition we are ignoring it
							sess.Stats.IncrementFilesIgnored()
							continue
						}

						if fi, err := os.Stat(fullFilePath); err == nil {
							fileSize := fi.Size()

							var mbFileMaxSize int64
							mbFileMaxSize = sess.MaxFileSize * 1024 * 1024

							// If the file is greater than the max size of a file we want to deal with then ignore it
							if fileSize > mbFileMaxSize {
								// If we are not scanning the file then by definition we are ignoring it
								sess.Stats.IncrementFilesIgnored()
								continue
							}
						}

						// If the file matches a file extension or other method that precludes it from a scan
						matchFile := newMatchFile(fullFilePath)
						if matchFile.isSkippable(sess) {
							// If we are not scanning the file then by definition we are ignoring it
							sess.Stats.IncrementFilesIgnored()
							continue
						}
						sess.Stats.IncrementFilesTotal()

						// We are now finally at the point where we are going to scan a file
						sess.Stats.IncrementFilesScanned()

						// for each signature that is loaded scan the file as a whole and generate a map of the match and the line number the match was found on
						for _, signature := range Signatures {

							bMatched, matchMap := signature.ExtractMatch(matchFile)
							if bMatched {

								sess.Stats.IncrementFilesDirty()

								var content string   // this is because file matches are puking
								var genericID string // the generic id used in the finding

								// for every instance of the secret that matched the specific rule create a new finding
								for k, v := range matchMap {

									cleanK := strings.SplitAfterN(k, "_", 2)
									if matchMap == nil {
										content = ""
										genericID = *repo.Name + "://" + fPath + "_" + generateGenericID(content)
									} else {
										content = cleanK[1]
										genericID = *repo.Name + "://" + fPath + "_" + generateGenericID(content)

									}

									// destroy the secret if the flag is set
									if sess.HideSecrets {
										content = ""
									}

									finding := &Finding{
										Action:          changeAction,
										Comment:         content,
										CommitAuthor:    commit.Author.String(),
										CommitHash:      commit.Hash.String(),
										CommitMessage:   strings.TrimSpace(commit.Message),
										Description:     signature.Description(),
										FilePath:        fPath,
										WraithVersion:   version.AppVersion(),
										LineNumber:      strconv.Itoa(v),
										RepositoryName:  *repo.Name,
										RepositoryOwner: *repo.Owner,
										Ruleid:          signature.Ruleid(),
										RulesVersion:    sess.RulesVersion,
										SecretID:        genericID,
									}

									// Get a proper uid for the finding
									finding.Initialize(sess.ScanType)

									// Add it to the hunt
									sess.AddFinding(finding)
									sess.Stats.IncrementCommits()
									sess.Out.Debug("[THREAD #%d][%s] Done analyzing changes in %s\n", tid, *repo.CloneURL, commit.Hash)

									dirtyCommit = true

									//print realtime data to stdout
									realTimeOutput(finding, sess)

								}
								sess.Out.Debug("[THREAD #%d][%s] Done analyzing commits\n", tid, *repo.CloneURL)
								if sess.InMemClone {
									os.RemoveAll(path)
								}
								sess.Out.Debug("[THREAD #%d][%s] Deleted %s\n", tid, *repo.CloneURL, path)
								//sess.Stats.IncrementRepositoriesScanned()
								//sess.Stats.UpdateProgress(sess.Stats.RepositoriesScanned, len(sess.Repositories))
							}
						}
					}
					// Increment the number of commits that were found t be dirty
					if dirtyCommit {
						sess.Stats.IncrementCommitsDirty()
					}
				}

				os.RemoveAll(path)
				sess.Stats.IncrementRepositoriesScanned()
			}
		}(i)
	}
	for _, repo := range sess.Repositories {
		ch <- repo
	}

	close(ch)
	wg.Wait()

}
