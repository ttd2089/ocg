package gitrepos

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ttd2089/ocg/pkg/shellout"
)

// IsRepo returns a bool indicating whether the given path points to a git repository.
func IsRepo(path string) (bool, error) {
	info, err := os.Stat(filepath.Join(path, ".git"))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// A Repo represents a Git repository on the local file system.
type Repo interface {

	// Name returns the basename of the repository directory.
	Name() string

	// Path returns the absolute path of the respository directory.
	Path() string

	// BranchNames returns the names of all local branches.
	BranchNames() ([]string, error)
}

// NewRepo returns a Repo representing the given path.
func NewRepo(path string) (Repo, error) {
	isRepo, err := IsRepo(path)
	if err != nil {
		return nil, err
	}
	if !isRepo {
		return nil, ErrNotAGitRepo
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return &repoImpl{
		path: absPath,
	}, nil
}

type repoImpl struct {
	path string
}

func (r *repoImpl) Name() string {
	return filepath.Base(r.path)
}

func (r *repoImpl) Path() string {
	return r.path
}

func (r *repoImpl) BranchNames() ([]string, error) {
	result, err := shellout.Run(
		"git", "-C", r.path, "for-each-ref", "--format=%(refname:short)", "refs/heads/")
	if err != nil {
		return nil, fmt.Errorf("failed to execute git command: %w", err)
	}
	if result.ExitCode != 0 {
		msg := strings.Split(result.Stderr.String(), "\n")[0]
		return nil, fmt.Errorf("git command returned exit code '%d': %s", result.ExitCode, msg)
	}
	branches := strings.Split(result.Stdout.String(), "\n")
	return branches[:len(branches)-1], nil
}
