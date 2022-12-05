package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ttd2089/ocg/internal/git"
	"github.com/ttd2089/ocg/internal/opts"
)

var listHelpText []string = []string{
	"usage: ocg list [<option>...] [<dir>]",
	"",
	"arguments:",
	"  dir    The directory to list (defaults to the current directory)",
	"",
	"options:",
	"  -h, --help    Print help text",
}

func newListCmd(appCtx appContext) cmd {
	return &listCmd{
		helpOpt: opts.FlagOpt{
			OptionName: opts.OptionName{
				LongName:  "help",
				ShortName: 'h',
			},
		},
		appCtx: appCtx,
	}
}

type listCmd struct {
	helpOpt opts.FlagOpt
	appCtx  appContext
}

func (l *listCmd) run(args []string) int {

	args, err := l.parseOptions(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n\n", err)
		l.help(os.Stderr)
		return 1
	}

	if len(args) > 1 {
		l.help(os.Stderr)
		return 1
	}

	if l.helpOpt.Value {
		l.help(os.Stdout)
		return 0
	}

	dir := l.appCtx.wd
	if len(args) == 1 {
		if filepath.IsAbs(args[0]) {
			dir = args[0]
		} else {
			dir = filepath.Join(l.appCtx.wd, args[0])
		}
	}

	repos, err := l.getRepos(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	output := new(bytes.Buffer)
	fmt.Fprintf(output, "repos:\n")
	for _, repo := range repos {
		err := l.printRepo(output, repo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			return 1
		}
	}

	io.Copy(os.Stdout, output)
	return 0
}

func (l *listCmd) parseOptions(args []string) ([]string, error) {
	return opts.Parse(
		args,
		[]opts.Option{
			&l.helpOpt,
		})
}

func (l *listCmd) getRepos(dir string) ([]git.Repo, error) {
	absRoot, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	var repos []git.Repo = nil
	walk := func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return err
		}
		repo, err := git.NewRepo(path, l.appCtx.gitCLI)
		if errors.Is(err, git.ErrNotAGitRepo) {
			return nil
		}
		if err != nil {
			return err
		}
		repos = append(repos, repo)
		return filepath.SkipDir
	}
	if err := filepath.WalkDir(absRoot, walk); err != nil {
		return nil, err
	}
	return repos, nil
}

func (l *listCmd) printRepo(w io.Writer, repo git.Repo) error {
	fmt.Fprintf(w, "- name: %s\n", repo.Name())
	fmt.Fprintf(w, "  path: %s\n", repo.Path())
	return l.printBranches(w, repo)
}

func (_ *listCmd) printBranches(w io.Writer, repo git.Repo) error {
	branches, err := repo.LocalBranches()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "  branches:\n")
	for _, branch := range branches {
		fmt.Fprintf(w, "  - name: %s\n", branch.Name)
		fmt.Fprintf(w, "    sha: %s\n", branch.SHA)
		if branch.Tracking != nil {
			fmt.Fprintf(w, "    remote:\n")
			fmt.Fprintf(w, "      name: %s\n", branch.Tracking.Name)
			fmt.Fprintf(w, "      sha: %s\n", branch.Tracking.SHA)
		}
	}
	return nil
}

func (_ *listCmd) help(w io.Writer) {
	fmt.Fprintf(w, "%s", strings.Join(listHelpText, "\n"))
}
