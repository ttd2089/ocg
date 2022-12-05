package opts

import (
	"errors"
	"testing"
)

func TestFlagOpt(t *testing.T) {

	underTest := FlagOpt{
		OptionName: OptionName{
			ShortName: 'h',
			LongName:  "help",
		},
	}

	type result struct {
		parsed    bool
		remaining []string
		err       error
	}

	tests := []struct {
		name           string
		startingValue  bool
		input          []string
		expectedResult result
		expectedValue  bool
	}{
		{
			name:          "Parses standalone short name reference",
			startingValue: false,
			input:         []string{"-h"},
			expectedResult: result{
				parsed:    true,
				remaining: []string{},
			},
			expectedValue: true,
		},
		{
			name:          "Parses clustered shortname reference",
			startingValue: false,
			input:         []string{"-hijkl"},
			expectedResult: result{
				parsed:    true,
				remaining: []string{"-ijkl"},
			},
			expectedValue: true,
		}, {
			name:          "Parses standalone longname reference",
			startingValue: false,
			input:         []string{"--help"},
			expectedResult: result{
				parsed:    true,
				remaining: []string{},
			},
			expectedValue: true,
		},
		{
			name:          "Parses longname reference with equals true",
			startingValue: false,
			input:         []string{"--help=true"},
			expectedResult: result{
				parsed:    true,
				remaining: []string{},
			},
			expectedValue: true,
		},
		{
			name:          "Parses longname reference with equals false",
			startingValue: true,
			input:         []string{"--help=false"},
			expectedResult: result{
				parsed:    true,
				remaining: []string{},
			},
			expectedValue: false,
		},
		{
			name:          "Parses longname reference with no- prefix",
			startingValue: true,
			input:         []string{"--no-help"},
			expectedResult: result{
				parsed:    true,
				remaining: []string{},
			},
			expectedValue: false,
		},
		{
			name:          "Parses longname reference with no- prefix and equals true",
			startingValue: true,
			input:         []string{"--no-help=true"},
			expectedResult: result{
				parsed:    true,
				remaining: []string{},
			},
			expectedValue: false,
		},
		{
			name:          "Parses longname reference with no- prefix and equals false",
			startingValue: false,
			input:         []string{"--no-help=false"},
			expectedResult: result{
				parsed:    true,
				remaining: []string{},
			},
			expectedValue: true,
		},
		{
			name:          "Ignores subsequent tokens after shortname reference",
			startingValue: false,
			input:         []string{"-h", "jkl", "mno"},
			expectedResult: result{
				parsed:    true,
				remaining: []string{"jkl", "mno"},
			},
			expectedValue: true,
		},
		{
			name:          "Ignores subsequent tokens after longname reference",
			startingValue: false,
			input:         []string{"--help", "jkl", "mno"},
			expectedResult: result{
				parsed:    true,
				remaining: []string{"jkl", "mno"},
			},
			expectedValue: true,
		},
		{
			name:          "Ignores different shortname reference",
			startingValue: false,
			input:         []string{"-i"},
			expectedResult: result{
				parsed:    false,
				remaining: []string{"-i"},
			},
			expectedValue: false,
		},
		{
			name:          "Ignores different longname reference",
			startingValue: false,
			input:         []string{"--indent"},
			expectedResult: result{
				parsed:    false,
				remaining: []string{"--indent"},
			},
			expectedValue: false,
		},
		{
			name:          "Returns error for equals with non-bool value",
			startingValue: false,
			input:         []string{"--help=7"},
			expectedResult: result{
				err: ErrInvalidOptionValue,
			},
			expectedValue: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			underTest.Value = tt.startingValue

			parsed, remaining, err := underTest.Parse(tt.input)

			if tt.expectedResult.err != nil {
				if !errors.Is(err, tt.expectedResult.err) {
					t.Errorf("expected err='%v'; got '%v'", tt.expectedResult.err, err)
					t.FailNow()
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if parsed != tt.expectedResult.parsed {
				t.Errorf("expected parsed='%t'; got '%t'", tt.expectedResult.parsed, parsed)
			}

			for i := range tt.expectedResult.remaining {
				if remaining[i] != tt.expectedResult.remaining[i] {
					t.Errorf("expected '%+v'; got '%v'", tt.expectedResult.remaining, remaining)
					t.FailNow()
				}
			}
		})
	}
}
