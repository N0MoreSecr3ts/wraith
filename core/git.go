// Package common contains functionality not critical to the core project but still essential.
package core

// TODO refactor out the common package

import (
	"errors"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
)

// Set easier names to refer to
const (
	TargetTypeUser         = "User"
	TargetTypeOrganization = "Organization"
)

// CloneConfiguration holds the configurations for cloning a repo
type CloneConfiguration struct {
	InMemClone *bool
	Url        *string
	Username   *string
	Token      *string
	Branch     *string
	Depth      *int
}

// Owner holds the info that we want for a repo owner
type Owner struct {
	Login     *string
	ID        *int64
	Type      *string
	Name      *string
	AvatarURL *string
	URL       *string
	Company   *string
	Blog      *string
	Location  *string
	Email     *string
	Bio       *string
}

// Repository holds the info we want for a repo itself
type Repository struct {
	Owner         *string
	ID            *int64
	Name          *string
	FullName      *string
	CloneURL      *string
	URL           *string
	DefaultBranch *string
	Description   *string
	Homepage      *string
}

// EmptyTreeCommit is a dummy commit id used as a placeholder and for testing
const (
	EmptyTreeCommitId = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"
)

// GetParentCommit will get the parent commit from a specific point. If the current commit
// has no parents then it will create a dummy commit.
func getParentCommit(commit *object.Commit, repo *git.Repository) (*object.Commit, error) {
	if commit.NumParents() == 0 {
		parentCommit, err := repo.CommitObject(plumbing.NewHash(EmptyTreeCommitId))
		if err != nil {
			return nil, err
		}
		return parentCommit, nil
	}
	parentCommit, err := commit.Parents().Next()
	if err != nil {
		return nil, err
	}
	return parentCommit, nil
}

// GetRepositoryHistory gets the commit history of a repository
func GetRepositoryHistory(repository *git.Repository) ([]*object.Commit, error) {
	var commits []*object.Commit
	ref, err := repository.Head()
	if err != nil {
		return nil, err
	}
	cIter, err := repository.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}
	cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})
	return commits, nil
}

// GetChanges will get the changes between to specific commits
func GetChanges(commit *object.Commit, repo *git.Repository) (object.Changes, error) {
	parentCommit, err := getParentCommit(commit, repo)
	if err != nil {
		return nil, err
	}

	commitTree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	parentCommitTree, err := parentCommit.Tree()
	if err != nil {
		return nil, err
	}

	changes, err := object.DiffTree(parentCommitTree, commitTree)
	if err != nil {
		return nil, err
	}
	return changes, nil
}

func GetChangeAction(change *object.Change) string {
	action, err := change.Action()
	if err != nil {
		return "Unknown"
	}
	switch action {
	case merkletrie.Insert:
		return "Insert"
	case merkletrie.Modify:
		return "Modify"
	case merkletrie.Delete:
		return "Delete"
	default:
		return "Unknown"
	}
}

// GetChangeAction will set the action of the commit to something that is more easily readable.
func GetChangePath(change *object.Change) string {
	action, err := change.Action()
	if err != nil {
		return change.To.Name
	}

	if action == merkletrie.Delete {
		fmt.Println(change.From.Name)
		return change.From.Name
	} else {
		fmt.Println(change.To.Name)
		return change.To.Name
	}
}

// GetChangeContent will get the contents of a git change or patch.
func GetChangeContent(change *object.Change) (result string, contentError error) {
	//temporary response to:  https://github.com/sergi/go-diff/issues/89
	// TODO Where possible switch to libgit2 https://github.com/libgit2/git2go
	defer func() {
		if err := recover(); err != nil {
			contentError = errors.New(fmt.Sprintf("Panic occurred while retrieving change content: %s", err))
		}
	}()
	patch, err := change.Patch()
	if err != nil {
		return "", err
	}
	for _, filePatch := range patch.FilePatches() {
		if filePatch.IsBinary() {
			continue
		}
		for _, chunk := range filePatch.Chunks() {
			result += chunk.Content()
		}
	}
	return result, nil
}
