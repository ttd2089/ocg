package opts

import (
	"fmt"
	"strconv"
	"strings"
)

// An OptionName represents the names an option may be addressed by.
type OptionName struct {

	// The long name of the option.
	LongName string

	// The short name of the option.
	ShortName rune
}

// A FlagOpt represents an option that contains a bool value.
type FlagOpt struct {

	// The name(s) of the flag.
	OptionName

	// The value of the flag.
	Value bool
}

func (f *FlagOpt) Parse(args []string) (bool, []string, error) {
	if len(args) == 0 {
		return false, args, nil
	}
	parsed, remaining, err := f.parseShort(args)
	if err != nil {
		return false, nil, err
	}
	if parsed {
		return true, remaining, nil
	}
	return f.parseLong(args, true)
}

func (f *FlagOpt) parseShort(args []string) (bool, []string, error) {
	if f.ShortName == 0 {
		return false, args, nil
	}
	shortNameRef := fmt.Sprintf("-%c", f.ShortName)
	if !strings.HasPrefix(args[0], shortNameRef) {
		return false, args, nil
	}
	f.Value = true
	if args[0] == shortNameRef {
		return true, args[1:], nil
	}
	remaining := append([]string{strings.Replace(args[0], shortNameRef, "-", 1)}, args[1:]...)
	return true, remaining, nil
}

func (f *FlagOpt) parseLong(args []string, supportNoPrefix bool) (bool, []string, error) {
	if f.LongName == "" {
		return false, args, nil
	}
	longNameRef := fmt.Sprintf("--%s", f.LongName)
	if args[0] == longNameRef {
		f.Value = true
		return true, args[1:], nil
	}
	refWithEquals := fmt.Sprintf("%s=", longNameRef)
	if strings.HasPrefix(args[0], refWithEquals) {
		optArg := strings.TrimPrefix(args[0], refWithEquals)
		b, err := strconv.ParseBool(optArg)
		if err != nil {
			return false, nil, NewInvalidOptionValueHelpText(f.LongName, optArg, "must be true or false")
		}
		f.Value = b
		return true, args[1:], nil
	}
	if supportNoPrefix {
		return f.parseNoPrefix(args)
	}
	return false, args, nil
}

func (f *FlagOpt) parseNoPrefix(args []string) (bool, []string, error) {
	noOpt := &FlagOpt{
		OptionName: OptionName{
			LongName: fmt.Sprintf("no-%s", f.LongName),
		},
		Value: true,
	}
	parsed, remaining, err := noOpt.parseLong(args, false)
	f.Value = !noOpt.Value
	return parsed, remaining, err
}
