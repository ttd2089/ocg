package gitrepos

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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
	cmd := exec.Command("git", "for-each-ref", "--format=%(refname:short)", "refs/heads/")
	cmd.Dir = r.path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		branches := strings.Split(stdout.String(), "\n")
		return branches[:len(branches)-1], nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		exitCode := exitErr.ExitCode()
		if exitCode == 128 {
			errMsg := fmt.Sprintf("'%s' is not a git repository", r.path)
			return nil, newTypedError(ErrNotAGitRepo, errMsg)
		}
		msg := strings.Split(stderr.String(), "\n")[0]
		return nil, fmt.Errorf("git command returned exit code '%d': %s", exitCode, msg)
	}
	return nil, fmt.Errorf("failed to execute git command: %w", err)
}
