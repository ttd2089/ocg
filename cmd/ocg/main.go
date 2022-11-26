package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 || os.Args[1] != "list" {
		fmt.Fprintf(os.Stderr, "usage: ocg list\n")
		os.Exit(1)
	}
	cmd := newListCmd()
	os.Exit(cmd.run(os.Args[2:]))
}
