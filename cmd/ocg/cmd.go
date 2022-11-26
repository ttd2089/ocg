package main

type cmd interface {
	run(args []string) int
}
