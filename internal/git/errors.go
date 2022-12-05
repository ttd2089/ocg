package git

import "errors"

// ErrNotAGitRepo is returned when a git-related operation is requested against a directory that is
// not a git repository.
var ErrNotAGitRepo error = errors.New("ErrNotAGitRepo")

// ErrFailedGitCommand is returned when ocg failed to execute a git command process.
var ErrFailedGitCommand error = errors.New("ErrFailedGitCommand")
