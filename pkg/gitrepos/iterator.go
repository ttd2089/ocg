package gitrepos

import (
	"io/fs"
	"path/filepath"
)

// An Iterator searches for and returns lists of Git repositories on the local
// file system.
type Iterator interface {

	// Iterate returns a list containing the absolute path to all Git
	// repositories under the specified root.
	Iterate(root string) ([]string, error)
}

// NewIter returns an implementation of Iterator.
func NewIter() Iterator {
	return &iterImpl{}
}

type iterImpl struct{}

// Iterate returns a list containing the absolute path to all Git repositories
// under the specified root.
func (i *iterImpl) Iterate(root string) ([]string, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	var repos []string = nil
	walk := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == ".git" {
			repo, _ := filepath.Split(path)
			repos = append(repos, repo)
		}
		return nil
	}
	if err := filepath.WalkDir(absRoot, walk); err != nil {
		return nil, err
	}
	return repos, nil
}
