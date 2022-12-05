package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ttd2089/shgit"
	"github.com/ttd2089/tyers"
)

// IsRepo returns a bool indicating whether the given path points to a git repository.
func IsRepo(path string) (bool, error) {
	info, err := os.Stat(filepath.Join(path, ".git"))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to determine if '%s' is a Git repo: %w", path, err)
	}
	return info.IsDir(), nil
}

// A Repo represents a Git repository on the local file system.
type Repo interface {

	// Name returns the basename of the repository directory.
	Name() string

	// Path returns the absolute path of the respository directory.
	Path() string

	LocalBranches() ([]LocalBranch, error)
}

// NewRepo returns a Repo representing the given path.
func NewRepo(absPath string, gitCLI shgit.CLI) (Repo, error) {
	if !filepath.IsAbs(absPath) {
		return nil, errors.New("NewRepo: absPath must be absolute")
	}
	isRepo, err := IsRepo(absPath)
	if err != nil {
		return nil, err
	}
	if !isRepo {
		return nil, tyers.Errorf(ErrNotAGitRepo, "'%s' is not a Git repo", absPath)
	}
	return &repo{
		path:   absPath,
		gitCLI: gitCLI,
	}, nil
}

type repo struct {
	path   string
	gitCLI shgit.CLI
}

func (r *repo) Name() string {
	return filepath.Base(r.path)
}

func (r *repo) Path() string {
	return r.path
}

func (r *repo) LocalBranches() ([]LocalBranch, error) {
	output, err := r.gitCLI.Run(
		"-C",
		r.path,
		"for-each-ref",
		"--format=%(refname) %(objectname) %(upstream:short)")
	if err != nil {
		return nil, fmt.Errorf("failed to get branches in repo '%s': %v", r.Path(), err)
	}

	lines := tokenizeLines(output)

	remotes := map[string]*Branch{}
	for _, tokens := range lines {
		if len(tokens) < 2 || len(tokens) > 3 {
			line := strings.Join(tokens, " ")
			return nil, fmt.Errorf("repo.LocalBranches(): unexpected output from git command: %s", line)
		}
		if !strings.HasPrefix(tokens[0], "refs/remotes/") {
			continue
		}
		name := strings.TrimPrefix(tokens[0], "refs/remotes/")
		remotes[name] = &Branch{
			Name: name,
			SHA:  tokens[1],
		}
	}

	locals := make([]LocalBranch, 0, (len(lines) - len(remotes)))
	for _, tokens := range lines {
		if !strings.HasPrefix(tokens[0], "refs/heads/") {
			continue
		}
		var tracking *Branch
		if len(tokens) == 3 {
			tracking, _ = remotes[tokens[2]]
		}
		name := strings.TrimPrefix(tokens[0], "refs/heads/")
		locals = append(locals, LocalBranch{
			Branch: Branch{
				Name: name,
				SHA:  tokens[1],
			},
			Tracking: tracking,
		})
	}

	return locals, nil
}

func tokenizeLines(s string) [][]string {
	lines := strings.Split(s, "\n")
	tokens := make([][]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tokens = append(tokens, strings.Split(line, " "))
	}
	return tokens
}
