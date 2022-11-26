package gitrepos

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

// An Iterator searches for and returns lists of Git repositories on the local file system.
type Iterator interface {

	// Iterate returns a list containing the absolute path to all Git repositories under the
	// specified root.
	Iterate(root string) ([]Repo, error)
}

// NewIter returns an implementation of Iterator.
func NewIter() Iterator {
	return &iterImpl{}
}

type iterImpl struct{}

// Iterate returns a list containing the absolute path to all Git repositories  under the specified
// root.
func (i *iterImpl) Iterate(root string) ([]Repo, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	var repos []Repo = nil
	walk := func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return err
		}
		isRepo, err := IsRepo(path)
		if err != nil {
			return fmt.Errorf("failed to determine if '%s' is a git repo: %w", path, err)
		}
		if !isRepo {
			return nil
		}
		repo, err := NewRepo(path)
		if err != nil {
			return fmt.Errorf("failed to get repo status for '%s': %w", path, err)
		}
		repos = append(repos, repo)
		return filepath.SkipDir
	}
	if err := filepath.WalkDir(absRoot, walk); err != nil {
		return nil, err
	}
	return repos, nil
}
