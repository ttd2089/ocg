package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ttd2089/ocg/pkg/gitrepos"
)

func newListCmd() cmd {
	return &listCmd{}
}

type listCmd struct{}

func (l *listCmd) run(_ []string) int {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to determine working directory: %s", err)
		return 1
	}
	repoIter := gitrepos.NewIter()
	repos, err := repoIter.Iterate(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to enumerate git repositories: %v\n", err)
		return 1
	}
	fmt.Println("repos:")
	for _, repo := range repos {
		l.printRepo(repo, wd)
	}
	return 0
}

func (l *listCmd) printRepo(repo gitrepos.Repo, wd string) {
	repoPath, _ := filepath.Rel(wd, repo.Path())
	fmt.Printf("- name: %s\n", repo.Name())
	fmt.Printf("  path: %s\n", repoPath)
	l.printBranches(repo)
}

func (_ *listCmd) printBranches(repo gitrepos.Repo) {
	branches, err := repo.BranchNames()
	if err != nil {
		fmt.Printf("  branches: [] # error: %s\n", err)
	} else {
		fmt.Printf("  branches:\n")
		for _, branch := range branches {
			fmt.Printf("  - %s\n", branch)
		}
	}
}
