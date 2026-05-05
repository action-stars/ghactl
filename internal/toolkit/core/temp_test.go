package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestGetTempDir(t *testing.T) {
	validDir := t.TempDir()
	nonExistentDir := filepath.Join(t.TempDir(), "test")

	tests := []struct {
		name    string
		env     string
		wantErr bool
	}{
		{
			name:    "errors_if_temp_directory_env_variable_is_not_defined",
			env:     "",
			wantErr: true,
		},
		{
			name:    "errors_if_temp_directory_does_not_exist",
			env:     nonExistentDir,
			wantErr: true,
		},
		{
			name: "returns_the_temp_directory",
			env:  validDir,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			t.Setenv(runnerTempLookup, tt.env)
			t.Setenv(linuxTempLookup, "")
			t.Setenv(windowsTempLookup, "")

			result, err := GetTempDir()

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)            // should not error
			is.Equal(result, tt.env) // should match
		})
	}
}

func TestCreateTempDir(t *testing.T) {
	validDir := t.TempDir()

	tests := []struct {
		name    string
		env     string
		wantErr bool
	}{
		{
			name:    "errors_if_no_temp_directory_is_configured",
			env:     "",
			wantErr: true,
		},
		{
			name: "creates_a_temporary_directory",
			env:  validDir,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			t.Setenv(runnerTempLookup, tt.env)
			t.Setenv(linuxTempLookup, "")
			t.Setenv(windowsTempLookup, "")

			result, err := CreateTempDir()

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)                              // should not error
			is.True(strings.HasPrefix(result, tt.env)) // should be inside temp dir

			info, err := os.Stat(result)
			is.NoErr(err)         // directory should exist
			is.True(info.IsDir()) // should be a directory
		})
	}
}
