package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestSaveState(t *testing.T) {
	tests := []struct {
		name    string
		envFile string
		key     string
		value   string
		wantErr bool
	}{
		{
			name:    "errors_if_file_env_variable_is_not_defined",
			envFile: "",
			key:     "key",
			value:   "value",
			wantErr: true,
		},
		{
			name:  "writes_state_entry",
			key:   "my state",
			value: "out val",
		},
		{
			name:  "writes_multi_line_state_entry",
			key:   "my state",
			value: "hello\nworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			envFile := tt.envFile
			if !tt.wantErr {
				envFile = filepath.Join(t.TempDir(), "test")
			}
			t.Setenv(stateFileLookup, envFile)

			err := SaveState(tt.key, tt.value)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			data, _ := os.ReadFile(envFile)

			is.NoErr(err)                                     // should not error
			is.True(strings.Contains(string(data), tt.key))   // should contain key
			is.True(strings.Contains(string(data), tt.value)) // should contain value
		})
	}
}

func TestGetState(t *testing.T) {
	tests := []struct {
		name  string
		env   string
		value string
		key   string
		want  string
	}{
		{
			name:  "gets_state_value",
			env:   "STATE_TEST_1",
			value: "state_val",
			key:   "TEST_1",
			want:  "state_val",
		},
		{
			name: "returns_empty_string_for_undefined_state",
			key:  "UNDEFINED_STATE",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			if tt.env != "" {
				t.Setenv(tt.env, tt.value)
			}

			result := GetState(tt.key)

			is.Equal(result, tt.want) // should match
		})
	}
}
