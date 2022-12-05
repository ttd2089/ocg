package opts

import (
	"errors"
	"testing"
)

func TestNewInvalidOptionValue(t *testing.T) {

	t.Run("errors.Is returns true for ErrInvalidOptionValue", func(t *testing.T) {
		a := NewInvalidOptionValue("a", "b")
		if !errors.Is(a, ErrInvalidOptionValue) {
			t.Errorf("expected true; got false")
		}
	})

}
