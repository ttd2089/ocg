package opts

import (
	"errors"
	"fmt"

	"github.com/ttd2089/tyers"
)

// ErrUnknownOption is returned when Parse encounters an option token that can't be parsed by any
// of the supplied Options.
var ErrUnknownOption error = errors.New("ErrUnknownOption")

// ErrInvalidOptionValue is returned when Option.Parse encounters an otherwise valid reference to
// itself with an invalid value for the option.
var ErrInvalidOptionValue error = errors.New("ErrInvalidOptionValue")

// NewInvalidOptionValue returns a new error describing an invalid value for a CLI option. The
// returned value will cause errors.Is to return true when ErrInvalidOptionValue is the target.
func NewInvalidOptionValue(option, value string) error {
	return NewInvalidOptionValueHelpText(option, value, "")
}

// NewInvalidOptionValueHelpText returns a new error describing an invalid value for a CLI option
// and a description of how to fix it. The returned value will cause errors.Is to return true when
// ErrInvalidOptionValue is the target.
func NewInvalidOptionValueHelpText(option, value, helpText string) error {
	return tyers.As(ErrInvalidOptionValue, &invalidOptionValue{
		option:   option,
		value:    value,
		helpText: helpText,
	})
}

type invalidOptionValue struct {
	option   string
	value    string
	helpText string
}

func (e *invalidOptionValue) Error() string {
	if e.helpText != "" {
		return fmt.Sprintf("invalid value '%s' for option '%s': %s", e.value, e.option, e.helpText)
	}
	return fmt.Sprintf("invalid value '%s' for option '%s'", e.value, e.option)
}
