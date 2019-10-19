package chglog

import (
	"errors"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var errReachedToCommit = errors.New("reached to commit")

// GitRepo open a GitRepo to use to build the changelog from
func GitRepo(gitPath string) (*git.Repository, error) {
	return git.PlainOpen(gitPath)
}

// GitHashFotTag return the git sha for a particular tag
func GitHashFotTag(gitRepo *git.Repository, tagName string) (hash plumbing.Hash, err error) {
	var ref *plumbing.Reference
	ref, err = gitRepo.Tag(tagName)
	if err == git.ErrTagNotFound {
		ref, err = gitRepo.Tag("v" + tagName)
	}
	if err != nil {
		return plumbing.ZeroHash, err
	}
	return ref.Hash(), nil
}

// CommitsBetween return the list of commits between to commits
func CommitsBetween(gitRepo *git.Repository, start, end plumbing.Hash) (commits []*object.Commit, err error) {
	var (
		commitIter object.CommitIter
	)
	commitIter, err = gitRepo.Log(&git.LogOptions{From: end})
	defer commitIter.Close()
	err = commitIter.ForEach(func(c *object.Commit) error {
		// If no previous tag is found then from and to are equal
		if end == start {
			return nil
		}
		if c.Hash == start {
			return errReachedToCommit
		}
		commits = append(commits, c)

		return nil
	})

	if err != nil && err != errReachedToCommit {
		return nil, err
	}

	return commits, nil
}
