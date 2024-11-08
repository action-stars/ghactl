package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestWriteSummary(t *testing.T) {
	t.Run("errors if file env variable is not defined", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(SUMMARY_FILE_LOOKUP, "")

		err := WriteSummary("### Hello world! :rocket:")

		is.True(err != nil) // should error
	})

	t.Run("writes summary", func(t *testing.T) {
		is := is.New(t)
		p := filepath.Join(t.TempDir(), "test")
		t.Setenv(SUMMARY_FILE_LOOKUP, p)
		summary := "### Hello world! :rocket:\n"

		err := WriteSummary(summary)

		data, _ := os.ReadFile(p)

		is.NoErr(err)                   // should not error
		is.Equal(string(data), summary) // should match
	})
}
