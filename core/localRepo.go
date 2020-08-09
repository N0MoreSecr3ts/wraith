// Package localRepo represents github specific functionality
package core

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"io/ioutil"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// CloneRepository will create either an in memory clone of a given repository or clone to a temp dir.
func CloneLocalRepository(cloneConfig *CloneConfiguration) (*git.Repository, string, error) {

	cloneOptions := &git.CloneOptions{
		URL:           *cloneConfig.Url,
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
