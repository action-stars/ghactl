package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestWriteSummary(t *testing.T) {
	tests := []struct {
		name    string
		envFile string
		value   string
		wantErr bool
	}{
		{
			name:    "errors_if_file_env_variable_is_not_defined",
			envFile: "",
			value:   "### Hello world! :rocket:",
			wantErr: true,
		},
		{
			name:  "writes_summary",
			value: "### Hello world! :rocket:\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			envFile := tt.envFile
			if !tt.wantErr {
				envFile = filepath.Join(t.TempDir(), "test")
			}
			t.Setenv(summaryFileLookup, envFile)

			err := WriteSummary(tt.value)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			data, _ := os.ReadFile(envFile)

			is.NoErr(err)                    // should not error
			is.Equal(string(data), tt.value) // should match
		})
	}
}
