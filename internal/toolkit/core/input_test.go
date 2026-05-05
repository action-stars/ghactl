package core

import (
	"testing"

	"github.com/matryer/is"
)

func TestGetInput(t *testing.T) {
	boolPtr := func(b bool) *bool { return &b }

	tests := []struct {
		name    string
		envKey  string
		envVal  string
		input   string
		opts    InputOptions
		want    string
		wantErr bool
	}{
		{
			name:   "gets_non-required_input",
			envKey: "INPUT_MY_INPUT",
			envVal: "val",
			input:  "my input",
			opts:   InputOptions{},
			want:   "val",
		},
		{
			name:   "gets_required_input",
			envKey: "INPUT_MY_INPUT",
			envVal: "val",
			input:  "my input",
			opts:   InputOptions{Required: true},
			want:   "val",
		},
		{
			name:    "errors_on_missing_required_input",
			envKey:  "INPUT_MISSING",
			envVal:  "",
			input:   "missing",
			opts:    InputOptions{Required: true},
			wantErr: true,
		},
		{
			name:   "returns_empty_for_missing_non-required_input",
			envKey: "INPUT_MISSING",
			envVal: "",
			input:  "missing",
			opts:   InputOptions{},
			want:   "",
		},
		{
			name:   "is_case_insensitive",
			envKey: "INPUT_MY_INPUT",
			envVal: "val",
			input:  "My InPuT",
			opts:   InputOptions{},
			want:   "val",
		},
		{
			name:   "handles_spaces_in_name",
			envKey: "INPUT_MULTIPLE_SPACES_VARIABLE",
			envVal: "I have multiple spaces",
			input:  "multiple spaces variable",
			opts:   InputOptions{},
			want:   "I have multiple spaces",
		},
		{
			name:   "trims_whitespace_by_default",
			envKey: "INPUT_WITH_TRAILING_WHITESPACE",
			envVal: "  some val  ",
			input:  "with trailing whitespace",
			opts:   InputOptions{},
			want:   "some val",
		},
		{
			name:   "trims_whitespace_when_option_is_true",
			envKey: "INPUT_WITH_TRAILING_WHITESPACE",
			envVal: "  some val  ",
			input:  "with trailing whitespace",
			opts:   InputOptions{TrimWhitespace: boolPtr(true)},
			want:   "some val",
		},
		{
			name:   "does_not_trim_whitespace_when_option_is_false",
			envKey: "INPUT_WITH_TRAILING_WHITESPACE",
			envVal: "  some val  ",
			input:  "with trailing whitespace",
			opts:   InputOptions{TrimWhitespace: boolPtr(false)},
			want:   "  some val  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			t.Setenv(tt.envKey, tt.envVal)

			result, err := GetInput(tt.input, tt.opts)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)             // should not error
			is.Equal(result, tt.want) // should match expected
		})
	}
}

func TestGetMultilineInput(t *testing.T) {
	boolPtr := func(b bool) *bool { return &b }

	tests := []struct {
		name   string
		envKey string
		envVal string
		input  string
		opts   InputOptions
		want   []string
	}{
		{
			name:   "gets_multiline_input",
			envKey: "INPUT_MY_INPUT_LIST",
			envVal: "val1\nval2\nval3",
			input:  "my input list",
			opts:   InputOptions{},
			want:   []string{"val1", "val2", "val3"},
		},
		{
			name:   "trims_whitespace_by_default",
			envKey: "INPUT_LIST_WITH_TRAILING_WHITESPACE",
			envVal: "  val1  \n  val2  \n  ",
			input:  "list with trailing whitespace",
			opts:   InputOptions{},
			want:   []string{"val1", "val2"},
		},
		{
			name:   "does_not_trim_whitespace_when_option_is_false",
			envKey: "INPUT_LIST_WITH_TRAILING_WHITESPACE",
			envVal: "  val1  \n  val2  \n  ",
			input:  "list with trailing whitespace",
			opts:   InputOptions{TrimWhitespace: boolPtr(false)},
			want:   []string{"  val1  ", "  val2  ", "  "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			t.Setenv(tt.envKey, tt.envVal)

			result, err := GetMultilineInput(tt.input, tt.opts)

			is.NoErr(err)             // should not error
			is.Equal(result, tt.want) // should match expected
		})
	}
}

func TestGetBooleanInput(t *testing.T) {
	tests := []struct {
		name    string
		envKey  string
		envVal  string
		input   string
		opts    InputOptions
		want    bool
		wantErr bool
	}{
		{
			name:   "gets_true_for_true",
			envKey: "INPUT_BOOL",
			envVal: "true",
			input:  "bool",
			want:   true,
		},
		{
			name:   "gets_true_for_True",
			envKey: "INPUT_BOOL",
			envVal: "True",
			input:  "bool",
			want:   true,
		},
		{
			name:   "gets_true_for_TRUE",
			envKey: "INPUT_BOOL",
			envVal: "TRUE",
			input:  "bool",
			want:   true,
		},
		{
			name:   "gets_false_for_false",
			envKey: "INPUT_BOOL",
			envVal: "false",
			input:  "bool",
			want:   false,
		},
		{
			name:   "gets_false_for_False",
			envKey: "INPUT_BOOL",
			envVal: "False",
			input:  "bool",
			want:   false,
		},
		{
			name:   "gets_false_for_FALSE",
			envKey: "INPUT_BOOL",
			envVal: "FALSE",
			input:  "bool",
			want:   false,
		},
		{
			name:    "errors_on_invalid_boolean",
			envKey:  "INPUT_WRONG_BOOLEAN",
			envVal:  "wrong",
			input:   "wrong boolean",
			wantErr: true,
		},
		{
			name:    "errors_on_required_missing_boolean",
			envKey:  "INPUT_MISSING_BOOL",
			envVal:  "",
			input:   "missing bool",
			opts:    InputOptions{Required: true},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			t.Setenv(tt.envKey, tt.envVal)

			result, err := GetBooleanInput(tt.input, tt.opts)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)             // should not error
			is.Equal(result, tt.want) // should match expected
		})
	}
}
