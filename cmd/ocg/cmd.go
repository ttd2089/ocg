package main

import "github.com/ttd2089/shgit"

type appContext struct {
	wd     string
	gitCLI shgit.CLI
}

type cmd interface {
	run(args []string) int
}
