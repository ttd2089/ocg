package main

import (
	"fmt"
	"os"

	"github.com/ttd2089/ocg/pkg/gitrepos"
)

func main() {
	repoIter := gitrepos.NewIter()
	repos, err := repoIter.Iterate(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to enumerate git repositories: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("found %d git repositories:\n", len(repos))
	for _, repo := range repos {
		fmt.Printf("- %s\n", repo)
	}
}
