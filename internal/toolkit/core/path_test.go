package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestAddPath(t *testing.T) {
	tests := []struct {
		name    string
		envFile string
		value   string
		wantErr bool
	}{
		{
			name:    "errors_if_file_env_variable_is_not_defined",
			envFile: "",
			value:   "$HOME/.local/bin",
			wantErr: true,
		},
		{
			name:  "writes_the_path_and_prepends_to_PATH",
			value: "$HOME/.local/bin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			envFile := tt.envFile
			if !tt.wantErr {
				envFile = filepath.Join(t.TempDir(), "test")
			}
			t.Setenv(pathFileLookup, envFile)
			t.Setenv("PATH", "existing-path")

			err := AddPath(tt.value)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			data, _ := os.ReadFile(envFile)

			is.NoErr(err)                                           // should not error
			is.True(strings.Contains(string(data), tt.value))       // should contain value
			is.True(strings.HasPrefix(os.Getenv("PATH"), tt.value)) // should prepend to PATH
		})
	}
}
