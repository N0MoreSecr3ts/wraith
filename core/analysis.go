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
	//fmt.Println("targets: ",sess.Stats.Targets)//TODO remove me

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
					//sess.Stats.IncrementTargets() // TODO remove me
					//fmt.Println("repos: ", sess.Stats.Repositories)//TODO remove me
					//fmt.Println("repos input: ", len(sess.RepoDirs))//TODO remove me

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

					// Set the url to the relative path of the repo based on the execution path of grover
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

					//sess.Stats.IncrementRepositories() //TODO remove me
					//sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))
					//fmt.Println(len(sess.Repositories))
					//fmt.Println("here") // TODO remove me
				}
			}
			return nil
		})
		if err0 != nil {
			//fmt.Println("err0", err0) // TODO remove me
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
				sess.Stats.IncrementTargets() //TODO why are we incrementing here and not above within the loop
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

// createFinding will create a discrete finding based on a match in a given repo and given commit
//func createFinding(repo Repository,
//	commit object.Commit,
//	change *object.Change,
//	fileSignature FileSignature,
//	contentSignature ContentSignature,
//scanType string) *Finding {

//	finding := &Finding{
//		FilePath:                    GetChangePath(change),
//		Action:                      GetChangeAction(change),
//		FileSignatureDescription:    fileSignature.GetDescription(),
//		FileSignatureComment:        fileSignature.GetComment(),
//		ContentSignatureDescription: contentSignature.GetDescription(),
//		ContentSignatureComment:     contentSignature.GetComment(),
//		RepositoryOwner:             *repo.Owner,
//		RepositoryName:              *repo.Name,
//		CommitHash:                  commit.Hash.String(),
//		CommitMessage:               strings.TrimSpace(commit.Message),
//		CommitAuthor:                commit.Author.String(),
//		CloneUrl:                    *repo.CloneURL,
//	}
//	finding.Initialize(scanType)
//	return finding
//
//}

// matchContent will attempt to match the content of a file such as a password or access token within the file
//func matchContent(sess *Session,
//	matchTarget MatchFile,
//	repo Repository,
//	change *object.Change,
//	commit object.Commit,
//	fileSignature FileSignature,
//	threadId int) {
//
//	content, err := GetChangeContent(change)
//	if err != nil {
//		sess.Out.Error("Error retrieving content in commit %s, change %s:  %s", commit.String(), change.String(), err)
//	}
//	matchTarget.Content = content
//	sess.Out.Debug("[THREAD #%d][%s] Matching content in %s...\n", threadId, *repo.CloneURL, commit.Hash)
//	for _, contentSignature := range sess.Signatures.ContentSignatures {
//		matched, err := contentSignature.Match(matchTarget)
//		if err != nil {
//			sess.Out.Error("Error while performing content match with '%s': %s\n", contentSignature.Description, err)
//		}
//		if !matched {
//			continue
//		}
//		finding := createFinding(repo, commit, change, fileSignature, contentSignature, sess.ScanType)
//		sess.AddFinding(finding)
//	}
//}

// findSecrets will attempt to find secrets in the content of a given file or match file in a given path
//func findSecrets(sess *Session, repo *Repository, commit *object.Commit, changes object.Changes, threadId int) {
//	for _, change := range changes {
//
//		path := GetChangePath(change)
//		matchTarget := newMatchFile(path)
//		if matchTarget.isSkippable(sess) {
//			sess.Out.Debug("[THREAD #%d][%s] Skipping %s\n", threadId, *repo.CloneURL, matchTarget.Path) // TODO implement me
//			continue
//		}
//		sess.Out.Debug("[THREAD #%d][%s] Inspecting file: %s...\n", threadId, *repo.CloneURL, matchTarget.Path) // TODO implement me
//
//		if sess.Mode != 3 {
//			for _, fileSignature := range sess.Signatures.FileSignatures {
//				matched, err := fileSignature.Match(matchTarget)
//				if err != nil {
//					sess.Out.Error(fmt.Sprintf("Error while performing file match: %s\n", err))
//				}
//				if !matched {
//					continue
//				}
//				if sess.Mode == 1 {
//					finding := createFinding(*repo, *commit, change, fileSignature,
//						ContentSignature{Description: "NA"}, sess.ScanType)
//					sess.AddFinding(finding)
//				}
//				if sess.Mode == 2 {
//					matchContent(sess, matchTarget, *repo, change, *commit, fileSignature, threadId)
//				}
//				break
//			}
//			sess.Stats.IncrementFiles()
//		} else {
//			matchContent(sess, matchTarget, *repo, change, *commit, FileSignature{Description: "NA"}, threadId)
//			sess.Stats.IncrementFiles()
//		}
//	}
//}

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
		//fmt.Println("clone path: ", path) // TODO remove me
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
			sess.Stats.IncrementRepositories()
			sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))
			return nil, "", err
		default:
			sess.Out.Error("Error cloning repository %s: %s\n", *repo.CloneURL, err)
			sess.Stats.IncrementRepositories()
			sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))
			return nil, "", err
		}
	}
	sess.Stats.IncrementRepositories()
	sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))
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
		sess.Stats.IncrementRepositories()
		sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))
		return nil, err
	}
	sess.Out.Debug("[THREAD #%d][%s] Number of commits: %d\n", threadId, *repo.CloneURL, len(history))
	return history, err
}

// AnalyzeRepositories will take a given repository, clone it, pull the commit history and use that as a basis for
// scanning for secrets within the repo and based on that output create a finding associated with that repo
//func AnalyzeRepositories(sess *Session) {
//	sess.Stats.Status = StatusAnalyzing
//	var ch = make(chan *Repository, len(sess.Repositories))
//	var wg sync.WaitGroup
//	var threadNum int
//	if len(sess.Repositories) <= 1 {
//		threadNum = 1
//	} else if len(sess.Repositories) <= sess.Threads {
//		threadNum = len(sess.Repositories) - 1
//	} else {
//		threadNum = sess.Threads
//	}
//	wg.Add(threadNum)
//	sess.Out.Debug("Threads for repository analysis: %d\n", threadNum)
//
//	sess.Out.Important("Analyzing %d %s...\n", len(sess.Repositories), Pluralize(len(sess.Repositories), "repository", "repositories"))
//
//	for i := 0; i < threadNum; i++ {
//		go func(tid int) {
//			for {
//				sess.Out.Debug("[THREAD #%d] Requesting new repository to analyze...\n", tid)
//				repo, ok := <-ch
//				if !ok {
//					sess.Out.Debug("[THREAD #%d] No more tasks, marking WaitGroup as done\n", tid)
//					wg.Done()
//					return
//				}
//
//				clone, path, err := cloneRepository(sess, repo, tid)
//				if err != nil {
//					continue
//				}
//
//				history, err := getRepositoryHistory(sess, clone, repo, path, tid)
//				if err != nil {
//					continue
//				}
//
//				for _, commit := range history {
//					sess.Out.Debug("[THREAD #%d][%s] Analyzing commit: %s\n", tid, *repo.CloneURL, commit.Hash)
//					changes, _ := GetChanges(commit, clone)
//					sess.Out.Debug("[THREAD #%d][%s] %s changes in %d\n", tid, *repo.CloneURL, commit.Hash, len(changes))
//
//					findSecrets(sess, repo, commit, changes, tid)
//
//					sess.Stats.IncrementCommits()
//					sess.Out.Debug("[THREAD #%d][%s] Done analyzing changes in %s\n", tid, *repo.CloneURL, commit.Hash)
//				}
//
//				sess.Out.Debug("[THREAD #%d][%s] Done analyzing commits\n", tid, *repo.CloneURL)
//				if sess.InMemClone {
//					os.RemoveAll(path)
//				}
//				sess.Out.Debug("[THREAD #%d][%s] Deleted %s\n", tid, *repo.CloneURL, path)
//				sess.Stats.IncrementRepositories()
//				sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))
//			}
//		}(i)
//	}
//	for _, repo := range sess.Repositories {
//		ch <- repo
//	}
//	close(ch)
//	wg.Wait()
//}
func AnalyzeRepositories(sess *Session) {
	sess.Stats.Status = StatusAnalyzing
	if len(sess.Repositories) == 0 {
		sess.Out.Error("No repositories have been gathered.")
		os.Exit(2)
	}

	var ch = make(chan *Repository, len(sess.Repositories))
	var wg sync.WaitGroup
	//var threadNum int

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
				//repo, ok := <-ch
				//if !ok {
				//	wg.Done()
				//	return
				//}
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


				// This is to set the org for inclusion in the reports. If we are scanning a local repo or filesystem then we default to local.
				// If we are scanning a GHE instance then we set the org based on the cloneURL
				//var orgPath string

				//if *hunt.Organizations[0].Login == "localRepo" {
				//	orgPath = "localRepo"
				//} else if *hunt.Organizations[0].Login == "github.com" {
				//	orgPath = "github.com"
				//} else {
				//	fullPath := strings.Split(*repo.CloneURL, ":")
				//	fullRepoPath := fullPath[1]
				//	repoPath := strings.Split(fullRepoPath, "/")
				//	orgPath = repoPath[0]
				//}

				// Get the hash for the last commit and add it. This is used to determine if any new commits have been made since the last time this was run
				//ref, _ := clone.Head()
				//lastCommit := fmt.Sprint(ref.Hash())
				//report.addCommit(orgPath, *repo.Name, *repo.DefaultBranch, lastCommit)

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

						//fmt.Println("path: ", path)

						// This is the total number of files that we know exist in out path. This does not care about the scan, it is simply the total number of files found
						//sess.Stats.IncrementFilesTotal() TODO implement in stats

						//likelyTestFile := false

						//if !hunt.ScanTests {  TODO implement in scanning for tests
						//	likelyTestFile = isTestFileOrPath(fullFilePath)
						//}

						if fi, err := os.Stat(fullFilePath); err == nil {
							fileSize := fi.Size()

							var mbFileMaxSize int64
							mbFileMaxSize = sess.MaxFileSize * 1024 * 1024 // TODO fully implement this

							// If the file is greater than the max size of a file we want to deal with then ignore it
							if fileSize > mbFileMaxSize {
								// If we are not scanning the file then by definition we are ignoring it
								//sess.Stats.IncrementFilesIgnored() TODO implement in stats
								continue
							}
						}

						// If the file is likely a test then ignore it
						//if likelyTestFile { TODO implement in stats
						//	// If we are not scanning the file then by definition we are ignoring it
						//	hunt.Stats.IncrementFilesIgnored()
						//	continue
						//}

						// If the file matches a file extension or other method that precludes it from a scan
						matchFile := newMatchFile(fullFilePath)
						if matchFile.isSkippable(sess) {
							// If we are not scanning the file then by definition we are ignoring it
							//hunt.Stats.IncrementFilesIgnored()TODO implement in stats
							continue
						}
						sess.Stats.IncrementFiles() // TODO this should be total files
						//fmt.Println("I am the matchFile:",matchFile)

						//if hunt.Debug {
						//	fmt.Println("Scanning ", fullFilePath)
						//} TODO implement in stats

						// We are now finally at the point where we are going to scan a file
						//hunt.Stats.IncrementFilesScanned()TODO implement in stats
						//idx := 1 TODO remove me
						// for each signature that is loaded scan the file as a whole and generate a map of the match and the line number the match was found on
						for _, signature := range Signatures {

							//fmt.Println("I am a sig: ", idx) TODO remove me
							//idx++ TODO remove me
							bMatched, matchMap := signature.ExtractMatch(matchFile)
							//fmt.Println("I done trying be matched") // TODO remove me
							if bMatched {
								//fmt.Println("I am matched") // TODO remove me

								var content string   // this is because file matches are puking
								var genericID string // the generic id used in the finding

								// for every instance of the secret that matched the specific rule create a new finding
								for k, v := range matchMap {

									// Increment the total number of findings found
									//sess.Stats.IncrementFindingsTotal()TODO implement in stats

									// Is the secret known to us already
									//knownSecret := false

									cleanK := strings.SplitAfterN(k, "_", 2)
									if matchMap == nil {
										content = ""
										genericID = *repo.Name + "://" + fPath + "_" + generateGenericID(content)
									} else {
										content = cleanK[1]
										genericID = *repo.Name + "://" + fPath + "_" + generateGenericID(content)

									}

									// destroy the secret if the flag is set
									//if hunt.HideSecrets { TODO implement in stats
									//	content = ""
									//}

									// if the secret, via the id, is already in the triage file we skip it
									//for _, h := range hunt.TriageIDs {
									//	if h == genericID {
									//		// increment the count of findings that are previously known and already accounted for by the triage file
									//		hunt.Stats.IncrementFindingsKnown()
									//		knownSecret = true
									//	}
									//}

									// if the secret is in the triage file do not report it
									//if knownSecret {
									//	continue
									//}

									finding := &Finding{
										Action:          changeAction,
										Comment:         content,
										CommitAuthor:    commit.Author.String(),
										CommitHash:      commit.Hash.String(),
										CommitMessage:   strings.TrimSpace(commit.Message),
										Description:     signature.Description(),
										FilePath:        fPath,
										GroverVersion:   version.AppVersion(),
										LineNumber:      strconv.Itoa(v),
										RepositoryName:  *repo.Name,
										RepositoryOwner: *repo.Owner,
										Ruleid:          signature.Ruleid(),
										//RulesVersion:    hunt.RulesVersion, TODO implement this
										SecretID: genericID,
									}

									// Get a proper uid for the finding
									finding.Initialize(sess.ScanType)

									//secret := &Secret{ // TODO this is all for the db output
									//	FilePath:    fPath,
									//	CommitHash:  commit.Hash.String(),
									//	Description: signature.Description(),
									//	ID:          generateGenericID(content),
									//	RuleID:      signature.Ruleid(),
									//}
									//
									//// Generate a proper uid
									//secret.Initialize()

									// Add the secret to the report for later inclusion in the db or dump to json/csv
									//report.addSecret(orgPath, *repo.Name, *repo.DefaultBranch, secret)

									// Add it to the hunt
									sess.AddFinding(finding)
									sess.Stats.IncrementCommits()
									sess.Out.Debug("[THREAD #%d][%s] Done analyzing changes in %s\n", tid, *repo.CloneURL, commit.Hash)

									dirtyCommit = true

									// print realtime data to stdout
									//realTimeOutput(finding, hunt)

									//if hunt.DBOutput {
									//	writeSecretToDB(orgPath, repo, secret, db, orgMap, repoMap, branchMap)
									//}
								}
								sess.Out.Debug("[THREAD #%d][%s] Done analyzing commits\n", tid, *repo.CloneURL)
								if sess.InMemClone {
									os.RemoveAll(path)
								}
								sess.Out.Debug("[THREAD #%d][%s] Deleted %s\n", tid, *repo.CloneURL, path)
								sess.Stats.IncrementRepositories()
								sess.Stats.UpdateProgress(sess.Stats.Repositories, len(sess.Repositories))
							}
						}
					}
					// Increment the number of commits that were found t be dirty
					if dirtyCommit {
						//	hunt.Stats.IncrementCommitsDirty() TODO implemnt in stats
					}
				}

				os.RemoveAll(path)
				sess.Stats.IncrementRepositories()
				//fmt.Println(len(sess.Repositories)) // TODO remove me
				//sess.Stats.IncrementRepositoriesScanned() TODO implement in stats
				//db.Close()
			}
		}(i)
	}
	for _, repo := range sess.Repositories {
		ch <- repo
	}

	close(ch)
	wg.Wait()

}
