package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ttd2089/ocg/internal/opts"
	"github.com/ttd2089/shgit"
)

// go build -ldflags="-X 'main.OCGVersion=<version>'"
var OCGVersion = "0.0.0.dev"

var ocgHelpText []string = []string{
	"usage: ocg [<option>...] <command> [<cmd-option>...] [<arg>...]",
	"",
	"options:",
	"  -h, --help       Invokes the help command",
	"  -v, --version    Invokes the version command",
	"",
	"commands:",
	"  list       List git repositories and their statuses",
	"  help       Print help text",
	"  version    Print OCG version information",
}

type ocgOptions struct {
	help    opts.FlagOpt
	version opts.FlagOpt
}

func main() {

	appCtx, err := getAppContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n\n", err)
		os.Exit(127)
	}

	ocgOpts, args, err := parseOptions(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n\n", err)
		help(os.Stderr)
		os.Exit(1)
	}

	if ocgOpts.help.Value {
		args = append([]string{"help"}, args...)
	} else if ocgOpts.version.Value {
		args = append([]string{"version"}, args...)
	} else if len(args) == 0 {
		args = []string{"help"}
	}

	var command cmd

	switch args[0] {
	case "list":
		command = newListCmd(appCtx)
	case "version":
		version()
		return
	case "help":
		help(os.Stdout)
		return
	default:
		fmt.Fprintf(os.Stderr, "ocg: unknown command '%s'\n\n", os.Args[1])
		help(os.Stderr)
		os.Exit(1)
	}

	os.Exit(command.run(os.Args[2:]))
}

func parseOptions(args []string) (ocgOptions, []string, error) {

	ocgOpts := ocgOptions{}

	ocgOpts.help = opts.FlagOpt{
		OptionName: opts.OptionName{
			LongName:  "help",
			ShortName: 'h',
		},
	}

	ocgOpts.version = opts.FlagOpt{
		OptionName: opts.OptionName{
			LongName:  "version",
			ShortName: 'v',
		},
	}

	remaining, err := opts.Parse(
		args,
		[]opts.Option{
			&ocgOpts.help,
			&ocgOpts.version,
		})

	return ocgOpts, remaining, err
}

func getAppContext() (appCtx appContext, err error) {

	wd, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("failed to determine working directory: %w", err)
		return
	}
	appCtx.wd = wd

	appCtx.gitCLI = shgit.NewCLI()

	return
}

func help(w io.Writer) {
	fmt.Fprintf(w, "%s", strings.Join(ocgHelpText, "\n"))
}

func version() {
	fmt.Printf("ocg version %s\n", OCGVersion)
}
