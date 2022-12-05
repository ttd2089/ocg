package git

// A Branch represents a git branch.
type Branch struct {

	// Name is the name of the branch excluding /refs/heads or /refs/remotes.
	Name string

	// SHA is the hash of the commit at the tip of the branch.
	SHA string
}

// A LocalBranch represents a git branch in a repository on the local file system.
type LocalBranch struct {
	Branch

	// RemoteBranch is the Branch that a LocalBranch is tracking.
	Tracking *Branch
}
