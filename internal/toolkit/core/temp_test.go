package core

import (
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestGetTempDirectory(t *testing.T) {
	t.Run("errors if temp directory env variable is not defined", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(RUNNER_TEMP_LOOKUP, "")

		_, err := GetTempDirectory()

		is.True(err != nil) // should error
	})

	t.Run("errors if temp directory does not exist", func(t *testing.T) {
		is := is.New(t)
		p := filepath.Join(t.TempDir(), "test")
		t.Setenv(RUNNER_TEMP_LOOKUP, p)

		_, err := GetTempDirectory()

		is.True(err != nil) // should error
	})

	t.Run("returns the temp directory", func(t *testing.T) {
		is := is.New(t)
		p := t.TempDir()
		t.Setenv(RUNNER_TEMP_LOOKUP, p)

		result, err := GetTempDirectory()

		is.NoErr(err)       // should not error
		is.Equal(result, p) // should match
	})
}
