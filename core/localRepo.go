package core

import (
	"crypto/sha1"
	"fmt"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// cloneLocal will create either an in memory clone of a given repository or clone to a temp dir.
func cloneLocal(cloneConfig *CloneConfiguration) (*git.Repository, string, error) {

	cloneOptions := &git.CloneOptions{
		URL:           *cloneConfig.URL,
		Depth:         *cloneConfig.Depth,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", *cloneConfig.Branch)),
		SingleBranch:  true,
		Tags:          git.NoTags,
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

// GatherLocalRepositories will grab all the local repos from the user input and generate a repository
// object, putting dummy or generated values in where necessary.
func GatherLocalRepositories(sess *Session) {

	// This is the number of targets as we don't do forks or anything else.
	// It will contain directorys, that will then be added to the repo count
	// if they contain a .git directory
	sess.Stats.Targets = len(sess.LocalPaths)
	sess.Stats.Status = StatusGathering
	sess.Out.Important("Gathering Local Repositories...\n")

	for _, pth := range sess.LocalPaths {

		if !PathExists(pth, sess) {
			sess.Out.Error("\n[*] <%s> does not exist! Quitting.\n", pth)
			os.Exit(1)
		}

		// Gather all paths in the tree
		err0 := filepath.Walk(pth, func(path string, f os.FileInfo, err1 error) error {
			if err1 != nil {
				sess.Out.Error("Failed to enumerate the path: %s\n", err1.Error())
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
						sess.Out.Error("Failed to open the repo HEAD: %s\n", err3.Error())
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
