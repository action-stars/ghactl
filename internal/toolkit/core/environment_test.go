package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestExportVariable(t *testing.T) {
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
			name:  "writes_single_line_entry",
			key:   "key",
			value: "value",
		},
		{
			name:  "writes_multi_line_entry",
			key:   "key",
			value: "hello\n    world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			envFile := tt.envFile
			if !tt.wantErr {
				envFile = filepath.Join(t.TempDir(), "test")
			}
			t.Setenv(envFileLookup, envFile)

			err := ExportVariable(tt.key, tt.value)

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
