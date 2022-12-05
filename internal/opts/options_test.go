package opts

import (
	"errors"
	"testing"
)

func TestParse(t *testing.T) {

	t.Run("Returns errors from options", func(t *testing.T) {
		expectedErr := errors.New("error from option")
		options := []Option{
			&mockOption{
				err: expectedErr,
			},
		}
		_, err := Parse([]string{"--foo"}, options)
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected '%v'; got '%v'", expectedErr, err)
		}
	})

	t.Run("Returns ErrUnknownOption for reference to unknown option", func(t *testing.T) {
		options := []Option{&mockOption{}}
		_, err := Parse([]string{"--foo"}, options)
		if !errors.Is(err, ErrUnknownOption) {
			t.Errorf("expected '%v'; got '%v'", ErrUnknownOption, err)
		}
	})

	t.Run("Stops parsing at first non-option token", func(t *testing.T) {
		expectedRemaining := []string{"foo", "bar"}
		mockOpt := &mockOption{}
		options := []Option{mockOpt}
		actualRemaining, err := Parse(expectedRemaining, options)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			t.FailNow()
		}
		if mockOpt.calls != 0 {
			t.Errorf("expected 0 call to option.Parse(); got %d", mockOpt.calls)
		}
		for i := range expectedRemaining {
			if actualRemaining[i] != expectedRemaining[i] {
				t.Errorf("expected '%+v'; got '%v'", expectedRemaining, actualRemaining)
				t.FailNow()
			}
		}
	})

	t.Run("Stops parsing at first non-option token", func(t *testing.T) {
		args := []string{"--", "bar"}
		expectedRemaining := []string{"bar"}
		mockOpt := &mockOption{}
		options := []Option{mockOpt}
		actualRemaining, err := Parse(args, options)
		if err != nil {
			t.Errorf("unexpected error: %v\n", err)
			t.FailNow()
		}
		if mockOpt.calls != 0 {
			t.Errorf("expected 0 call to option.Parse(); got %d", mockOpt.calls)
		}
		for i := range expectedRemaining {
			if actualRemaining[i] != expectedRemaining[i] {
				t.Errorf("expected '%+v'; got '%v'", expectedRemaining, actualRemaining)
				t.FailNow()
			}
		}
	})

	t.Run("Continues parsing using the remaining values returned from Option.Parse", func(t *testing.T) {
		// This test relies on the fact that encountering ["bar"] will cause Parse() to stop
		// parsing to demonstrate that Parse() will parse the reamining values returned from
		// Option.Parse after it's called.
		args := []string{"--foo"}
		expectedRemaining := []string{"bar"}
		mockOpt := &mockOption{parsed: true, remaining: []string{"bar"}}
		options := []Option{mockOpt}
		actualRemaining, err := Parse(args, options)
		if err != nil {
			t.Errorf("unexpected error: %v\n", err)
			t.FailNow()
		}
		if mockOpt.calls != 1 {
			t.Errorf("expected 1 call to option.Parse(); got %d", mockOpt.calls)
		}
		for i := range expectedRemaining {
			if actualRemaining[i] != expectedRemaining[i] {
				t.Errorf("expected '%+v'; got '%v'", expectedRemaining, actualRemaining)
				t.FailNow()
			}
		}
	})
}

type mockOption struct {
	parsed    bool
	remaining []string
	err       error
	calls     int
}

func (m *mockOption) Parse(args []string) (bool, []string, error) {
	m.calls += 1
	return m.parsed, m.remaining, m.err
}
