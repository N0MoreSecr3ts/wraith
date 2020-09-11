// Package core represents the core functionality of all commands
package core

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"wraith/version"
)

// GatherTargets will enumerate git targets adding them to a running target list. This will set the targets based
// on the scan type set within the cmd package.
func GatherTargets(sess *Session) {
	sess.Stats.Status = StatusGathering
	sess.Out.Important("Gathering targets...\n")

	var targets []string

	// Based on the type of scan, set in the cmd package, we set a generic
	// variable to the specific targets
	switch sess.ScanType {
	case "github":
		targets = sess.GithubTargets
	case "gitlab":
		targets = sess.GitlabTargets
	}

	// For each target that the user provided, we use the client set in the session
	// initialization to enumerate the target. There are flag that be used here to
	// decide if forks are followed the scope of a target can be increased a lot. This
	// could be useful as some developers may keep secrets in their forks, yet purge
	// them before creating a pull request. Developers may also keep a specific environment
	// file within there repo that is not set to be ignored so they can more easily develop
	// on multiple boxes or collaborate with multiple people.
	for _, loginOption := range targets {
		target, err := sess.Client.GetUserOrganization(loginOption)
		if err != nil || target == nil {
			sess.Out.Error(" Error retrieving information on %s: %s\n", loginOption, err)
			continue
		}
		sess.Out.Debug("%s (ID: %d) type: %s\n", *target.Login, *target.ID, *target.Type)
		sess.AddTarget(target)
		// If forking is false AND the target type is an Organization as set above in GetUserOrganization
		if sess.ExpandOrgs == true && *target.Type == TargetTypeOrganization {
			sess.Out.Debug("Gathering members of %s (ID: %d)...\n", *target.Login, *target.ID)
			members, err := sess.Client.GetOrganizationMembers(*target)
			if err != nil {
				sess.Out.Error(" Error retrieving members of %s: %s\n", *target.Login, err)
				continue
			}
			// Add organization members gathered above to the target list
			// TODO Do we want to spider this out at some point to enumerate all members of an org?
			for _, member := range members {
				sess.Out.Debug("Adding organization member %s (ID: %d) to targets\n", *member.Login, *member.ID)
				sess.AddTarget(member)
			}
		}
	}
}

// GatherLocalRepositories will grab all the local repos from the user input and generate a repository
// object, putting dummy or generated values in where necessary
//func GatherLocalRepositories(sess *Session) {
//
//	// This is the number of targets as we don't do forks or anything else.
//	// It will contain directorys, that will then be added to the repo count
//	// if they contain a .git directory
//	sess.Stats.Targets = len(sess.LocalDirs)
//	sess.Stats.Status = StatusGathering
//	sess.Out.Important("Gathering Local Repositories...\n")
//
//	for _, pth := range sess.LocalDirs {
//
//		if !PathExists(pth, sess) {
//			sess.Out.Error("\n[*] <%s> does not exist! Quitting.\n", pth)
//			os.Exit(1)
//		}
//
//		// Gather all paths in the tree
//		err0 := filepath.Walk(pth, func(path string, f os.FileInfo, err1 error) error {
//			if err1 != nil {
//				sess.Out.Error("Failed to enumerate the path: %s\n", err1.Error())
//				return nil
//			}
//
//			// If it is a directory then move forward
//			if f.IsDir() {
//
//				// If there is a .git directory then we have a repo
//				if filepath.Ext(path) == ".git" { // TODO Should we reverse this to ! to make the code cleaner
//
//					parent, _ := filepath.Split(path)
//
//					gitProjName, _ := filepath.Split(parent)
//
//					openRepo, err2 := git.PlainOpen(parent)
//					if err2 != nil {
//
//						return nil
//					}
//
//					ref, err3 := openRepo.Head()
//					if err3 != nil {
//						sess.Out.Error("Failed to open the repo HEAD: %s\n", err3.Error())
//						return nil
//					}
//
//					// Get the name of the branch we are working on
//					s := ref.Strings()
//					branchPath := fmt.Sprintf("%s", s[0])
//					branchPathParts := strings.Split(branchPath, string("refs/heads/"))
//					branchName := branchPathParts[len(branchPathParts)-1]
//					pBranchName := &branchName
//
//					commit, _ := openRepo.CommitObject(ref.Hash())
//					var commitHash = commit.Hash[:]
//
//					// TODO make this a generic function at some point
//					// Generate a uid for the repo
//					h := sha1.New()
//					repoID := fmt.Sprintf("%x", h.Sum(commitHash))
//
//					intRepoID, _ := strconv.ParseInt(repoID, 10, 64)
//					var pRepoID *int64
//					pRepoID = &intRepoID
//
//					// Set the url to the relative path of the repo based on the execution path of wraith
//					pRepoURL := &parent
//
//					// This is used to id the owner, fullname, and description of the repo. It is ugly but effective. It is the relative path to the repo, for example ../foo
//					pGitProjName := &gitProjName
//
//					// The project name is simply the parent directory in the case of a local scan with all other path bits removed for example ../foo -> foo.
//					projectPathParts := strings.Split(*pGitProjName, string(os.PathSeparator))
//					pProjectName := &projectPathParts[len(projectPathParts)-2]
//
//					sessR := Repository{
//						Owner:         pGitProjName,
//						ID:            pRepoID,
//						Name:          pProjectName,
//						FullName:      pGitProjName,
//						CloneURL:      pRepoURL,
//						URL:           pRepoURL,
//						DefaultBranch: pBranchName,
//						Description:   pGitProjName,
//						Homepage:      pRepoURL,
//					}
//
//					// Add the repo to the sess to be cloned and scanned
//					sess.AddRepository(&sessR)
//				}
//			}
//			return nil
//		})
//		if err0 != nil {
//		}
//	}
//}

// Gather Repositories will gather all repositories associated with a given target during a scan session.
// This is done using threads, whose count is set via commandline flag. Care much be taken to avoid rate
// limiting associated with suspected DOS attacks.
//func GatherRepositories(sess *Session) {
//	var ch = make(chan *Owner, len(sess.Targets))
//	sess.Out.Debug("Number of targets: %d\n", len(sess.Targets))
//	var wg sync.WaitGroup
//	var threadNum int
//	if len(sess.Targets) == 1 {
//		threadNum = 1
//	} else if len(sess.Targets) <= sess.Threads {
//		threadNum = len(sess.Targets) - 1
//	} else {
//		threadNum = sess.Threads
//	}
//	wg.Add(threadNum)
//	sess.Out.Debug("Threads for repository gathering: %d\n", threadNum)
//	for i := 0; i < threadNum; i++ {
//		go func() {
//			for {
//				target, ok := <-ch
//				if !ok {
//					wg.Done()
//					return
//				}
//				repos, err := sess.Client.GetRepositoriesFromOwner(*target)
//				if err != nil {
//					sess.Out.Error(" Failed to retrieve repositories from %s: %s\n", *target.Login, err)
//				}
//				if len(repos) == 0 {
//					continue
//				}
//				for _, repo := range repos {
//					sess.Out.Debug(" Retrieved repository: %s\n", *repo.CloneURL)
//					sess.AddRepository(repo)
//				}
//				sess.Out.Info(" Retrieved %d %s from %s\n", len(repos), Pluralize(len(repos), "repository", "repositories"), *target.Login)
//			}
//		}()
//	}
//
//	for _, target := range sess.Targets {
//		ch <- target
//	}
//	close(ch)
//	wg.Wait()
//}
//
////	sess.Out.Debug("[THREAD #%d][%s] Skipping %s\n", threadId, *repo.CloneURL, matchTarget.Path) // TODO implement me
////
////sess.Out.Debug("[THREAD #%d][%s] Inspecting file: %s...\n", threadId, *repo.CloneURL, matchTarget.Path) // TODO implement me
////
////			sess.Out.Error(fmt.Sprintf("Error while performing file match: %s\n", err))

// cloneRepository will clone a given repository based upon a configured set or options a user provides
//func cloneRepository(sess *Session, repo *Repository, threadId int) (*git.Repository, string, error) {
//	sess.Out.Debug("[THREAD #%d][%s] Cloning repository...\n", threadId, *repo.CloneURL)
//
//	var clone *git.Repository
//	var path string
//	var err error
//
//	switch sess.ScanType {
//	case "github":
//		cloneConfig := CloneConfiguration{
//			Url:        repo.CloneURL,
//			Branch:     repo.DefaultBranch,
//			Depth:      &sess.CommitDepth,
//			InMemClone: &sess.InMemClone,
//			Token:      &sess.GithubAccessToken,
//		}
//		clone, path, err = CloneGithubRepository(&cloneConfig)
//	case "gitlab":
//		userName := "oauth2"
//		cloneConfig := CloneConfiguration{
//			Url:        repo.CloneURL,
//			Branch:     repo.DefaultBranch,
//			Depth:      &sess.CommitDepth,
//			Token:      &sess.GitlabAccessToken, // TODO Is this need since we already have a client?
//			InMemClone: &sess.InMemClone,
//			Username:   &userName,
//		}
//		clone, path, err = CloneGitlabRepository(&cloneConfig)
//	case "localGit":
//		cloneConfig := CloneConfiguration{
//			Url:        repo.CloneURL,
//			Branch:     repo.DefaultBranch,
//			Depth:      &sess.CommitDepth,
//			InMemClone: &sess.InMemClone,
//		}
//		clone, path, err = CloneLocalRepository(&cloneConfig)
//
//	}
//	if err != nil {
//		switch err.Error() {
//		case "remote repository is empty":
//			sess.Out.Error("Repository %s is empty: %s\n", *repo.CloneURL, err)
//			sess.Stats.IncrementRepositoriesCloned()
//			//sess.Stats.UpdateProgress(sess.Stats.RepositoriesCloned, len(sess.Repositories))
//			return nil, "", err
//		default:
//			sess.Out.Error("Error cloning repository %s: %s\n", *repo.CloneURL, err)
//			//sess.Stats.UpdateProgress(sess.Stats.RepositoriesCloned, len(sess.Repositories))
//			return nil, "", err
//		}
//	}
//	sess.Stats.IncrementRepositoriesCloned()
//	//sess.Stats.UpdateProgress(sess.Stats.RepositoriesCloned, len(sess.Repositories))
//	sess.Out.Debug("[THREAD #%d][%s] Cloned repository to: %s\n", threadId, *repo.CloneURL, path)
//	return clone, path, err
//}

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

//>>>>>>> 33e8672995d58dbbbca9fe5a6d5e56505d77f933
func AnalyzeRepositories(sess *Session) {
	sess.Stats.Status = StatusAnalyzing
	if len(sess.Repositories) == 0 {
		sess.Out.Error("No repositories have been gathered.")
		os.Exit(2)
	}

	var ch = make(chan *Repository, len(sess.Repositories))
	var wg sync.WaitGroup

	// Calculate the number of threads based on the flag and the number of repos
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

				// Clone the repository from the remote source or if a local repo from the path
				// The path variable is returning the path that the clone was done to. The repo is cloned directly
				// there.
				clone, path, err := cloneRepository(sess, repo, tid)
				if err != nil {
					if err.Error() != "repository is empty" {
						sess.Out.Error("Error cloning repository %s: %s\n", *repo.FullName, err)
					}
					continue
				}

				// Get the full commit history for the repo
				history, err := GetRepositoryHistory(clone)
				if err != nil {
					sess.Out.Error("[THREAD #%d][%s] Error getting commit history: %s\n", tid, *repo.CloneURL, err)
					if sess.InMemClone {
						err := os.RemoveAll(path)
						sess.Out.Error("[THREAD #%d][%s] Error removing path from memory: %s\n", tid, *repo.CloneURL, err)
					} else {
						err := os.RemoveAll(path)
						sess.Out.Error("[THREAD #%d][%s] Error removing path from disk: %s\n", tid, *repo.CloneURL, err)
					}
					continue
				}
				//sess.Stats.IncrementRepositories()
				//sess.Stats.UpdateProgress(sess.Stats.RepositoriesCloned, len(sess.Repositories))
				sess.Out.Debug("[THREAD #%d][%s] Number of commits: %d\n", tid, *repo.CloneURL, len(history))

				// For every commit in the history we want to look through it for any changes
				// there is a known bug in here related to files that have changed paths from the most
				// recent path. The does not do a fetch per history so if a file changes paths from
				// the current one it will throw a file not found error. You can see this by turning
				// on debugging.
				for _, commit := range history {

					sess.Out.Debug("[THREAD #%d][%s] Analyzing commit: %s\n", tid, *repo.CloneURL, commit.Hash)

					// Increment the total number of commits. This needs to be used in conjunction with
					// the total number of commits scanned as a commit may have issues and not be scanned once
					// it is found.
					sess.Stats.IncrementCommits()

					// This will be used to increment the dirty commit stat if any matches are found. A dirty commit
					// means that a secret was found in that commit. This provides an easier way to manually to look
					// through the commit history of a given repo.
					dirtyCommit := false

					// TODO what is this actually doing here?
					changes, _ := GetChanges(commit, clone)
					sess.Out.Debug("[THREAD #%d][%s] %s changes in %d\n", tid, *repo.CloneURL, commit.Hash, len(changes))

					for _, change := range changes {

						// TODO Is this need for the finding object, why are we saving this?
						changeAction := GetChangeAction(change)

						// TODO Add an example of the output from this function
						fPath := GetChangePath(change)

						// TODO Add an example of this
						fullFilePath := path + "/" + fPath

						sess.Stats.IncrementFilesTotal()

						// required as that is a map of interfaces.
						scanTests := DefaultValues["scan-tests"]
						likelyTestFile := scanTests.(bool)

						// If we do not want to scan tests we run some checks to see if the file in
						// question is a test file. This will return a true if it is a test file.
						if !sess.ScanTests {
							likelyTestFile = isTestFileOrPath(fullFilePath)
						}

						// If the file is likely a test then ignore it. By default this is currently
						// set to false which means we do NOT want to scan tests. This means that we
						// check above and if this returns true because it is likely a test file, we
						// increment the ignored file count and pass through scanning the file and content.
						if likelyTestFile {
							// If we are not scanning the file then by definition we are ignoring it
							sess.Stats.IncrementFilesIgnored()
							sess.Out.Debug("%s is a test file and being ignored\n", fPath)

							continue
						}

						// Check the file size of the file. If it is greater than the default size then
						// then we increment the ignored file count and pass on through.
						if IsMaxFileSize(fullFilePath, sess) {

							sess.Stats.IncrementFilesIgnored()
							sess.Out.Debug("%s is too large and being ignored\n", fPath)

							continue
						}

						// Break a file name up into its composite pieces including the extension and base name
						matchFile := newMatchFile(fullFilePath)

						// If the file extension matches an extension or other criteria that precludes
						// it from a scan we increment the ignored files count and pass on through.
						if matchFile.isSkippable(sess) {
							sess.Stats.IncrementFilesIgnored()
							sess.Out.Debug("%s is skippable and being ignored\n", fPath)

							continue
						}

						// The total number of files that were evaluated
						sess.Stats.IncrementFilesTotal()

						// We are now finally at the point where we are going to scan a file so we implement
						// that count.
						sess.Stats.IncrementFilesScanned()

						// for each signature that is loaded scan the file as a whole and generate a map of
						// the match and the line number the match was found on
						for _, signature := range Signatures {

							bMatched, matchMap := signature.ExtractMatch(matchFile, sess, change)
							if bMatched {

								// Incremented the count of files that contain secrets
								sess.Stats.IncrementFilesDirty()

								// content will hold the secret found within the target
								var content string

								// For every instance of the secret that matched the specific signatures
								// create a new finding. Thi will produce dupes as the file may exist
								// in multiple commits.
								for k, v := range matchMap {

									// This sets the content for the finding, in this case the actual secret
									// is the content. This can be removed and hidden via a commandline flag.
									cleanK := strings.SplitAfterN(k, "_", 2)
									if matchMap == nil {
										content = ""
									} else {
										content = cleanK[1]

									}

									// Destroy the secret by zeroing the content if the flag is set
									if sess.HideSecrets {
										content = ""
									}

									// Create a new instance of a finding and set the necessary fields.
									finding := &Finding{
										Action:            changeAction,
										Content:           content,
										CommitAuthor:      commit.Author.String(),
										CommitHash:        commit.Hash.String(),
										CommitMessage:     strings.TrimSpace(commit.Message),
										Description:       signature.Description(),
										FilePath:          fPath,
										WraithVersion:     version.AppVersion(),
										LineNumber:        strconv.Itoa(v),
										RepositoryName:    *repo.Name,
										RepositoryOwner:   *repo.Owner,
										Signatureid:       signature.Signatureid(),
										SignaturesVersion: sess.SignatureVersion,
										SecretID:          generateID(),
									}
									// Set the urls for the finding
									finding.Initialize(sess)
									//fNew := true

									//for _, f := range sess.Findings { // TODO this is for de-duping if needed
									//	if f.CommitHash == finding.CommitHash && f.SecretID == finding.SecretID && f.Description == finding.Description {
									//		fNew = false
									//		continue
									//	}
									//}

									if true {
										// Add it to the session
										sess.AddFinding(finding)
										sess.Out.Debug("[THREAD #%d][%s] Done analyzing changes in %s\n", tid, *repo.CloneURL, commit.Hash)

										dirtyCommit = true

										// Print realtime data to stdout
										realTimeOutput(finding, sess)
									}

								}
								sess.Out.Debug("[THREAD #%d][%s] Done analyzing commits\n", tid, *repo.CloneURL)
								if sess.InMemClone {
									err = os.RemoveAll(path)
									if err != nil {
										sess.Out.Error("Could not remove path from memory: %s", err.Error())
									}
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

				err = os.RemoveAll(path)
				if err != nil {
					sess.Out.Error("Could not remove path from disk: %s", err.Error())
				}
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
