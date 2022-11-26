package gitrepos

import "errors"

// ErrNotAGitRepo is returned when a git-related operation is requested against a directory that is
// not a git repository.
var ErrNotAGitRepo error = errors.New("the target directory is not a git repository")

// ErrFailedGitCommand is returned when ocg failed to execute a git command process.
var ErrFailedGitCommand error = errors.New("failed to execute git command")

type typedError struct {
	errType error
	msg     string
}

func newTypedError(errType error, msg string) error {
	return &typedError{
		errType: errType,
		msg:     msg,
	}
}

func (t *typedError) Error() string {
	return t.msg
}

func (t *typedError) Unwrap() error {
	return t.errType
}
