package opts

import (
	"strings"

	"github.com/ttd2089/tyers"
)

// An Option represents an CLI option that can be parsed from argv and contain a value or a flag
// that will modify the beahviour of a program.
type Option interface {

	// Parse attempts to consume a specifiation for the target Option from the first value of args.
	//
	// If the first value of args is a valid reference to the Option then Parse will return true
	// and a modified copy of args with the reference to the Option removed.
	//
	// If the reference to the Option is invalid then Parse will return an error.
	//
	// If the fierst value of args does not reference the option then Parse will return false and
	// an unmodified reference to args.
	Parse(args []string) (bool, []string, error)
}

// Parse parses args using the values of options until the first non-option argument is encountered
// or an error occurs.
func Parse(args []string, options []Option) ([]string, error) {
	for len(args) > 0 {
		if isEndOfOptionsDelimiter(args[0]) {
			return args[1:], nil
		}
		if !isOptionRef(args[0]) {
			return args, nil
		}
		parsed, remaining, err := parseNextOption(args, options)
		if err != nil {
			return nil, err
		}
		if !parsed {
			return nil, tyers.Errorf(ErrUnknownOption, "unknown option '%s'", args[0])
		}
		args = remaining
	}
	return []string{}, nil
}

func isEndOfOptionsDelimiter(token string) bool {
	// https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html#tag_12_02
	// Guideline 10
	return token == "--"
}

func isOptionRef(token string) bool {
	return strings.HasPrefix(token, "-")
}

func parseNextOption(args []string, options []Option) (bool, []string, error) {
	for _, opt := range options {
		parsed, remaining, err := opt.Parse(args)
		if err != nil {
			return false, nil, err
		}
		if parsed {
			return true, remaining, nil
		}
	}
	return false, args, nil
}
