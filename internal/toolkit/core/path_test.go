package core

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/matryer/is"
)

func TestAddPath(t *testing.T) {
	for _, tt := range []struct {
		name         string
		envFile      string
		existingPath string
		value        string
		wantErr      bool
	}{
		{
			name:    "errors_if_file_env_variable_is_not_set",
			envFile: "",
			value:   "$HOME/.local/bin",
			wantErr: true,
		},
		{
			name:    "writes_the_path",
			envFile: filepath.Join(t.TempDir(), "test"),
			value:   "$HOME/.local/bin",
		},
		{
			name:         "writes_the_path_with_existing_content",
			envFile:      filepath.Join(t.TempDir(), "test"),
			existingPath: "/xxx/yyy/bin",
			value:        "$HOME/.local/bin",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			if tt.envFile != "" && tt.existingPath != "" {
				if err := writeFile(tt.envFile, []byte(tt.existingPath+"\n")); err != nil {
					t.Fatalf("failed to create env file: %v", err)
				}
			}

			t.Setenv(pathFileLookup, tt.envFile)

			err := AddPath(tt.value)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			got, _ := os.ReadFile(tt.envFile)

			is.NoErr(err)                                                                      // should not error
			is.True(regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(tt.value) + `$`).Match(got)) // should contain value on its own line
			if tt.existingPath != "" {
				is.True(regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(tt.existingPath) + `$`).Match(got)) // should not overwrite existing content
			}
		})
	}
}
