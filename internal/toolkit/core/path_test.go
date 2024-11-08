package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestAddPath(t *testing.T) {
	t.Run("errors if file env variable is not defined", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(PATH_FILE_LOOKUP, "")

		err := AddPath("$HOME/.local/bin")

		is.True(err != nil) // should error
	})

	t.Run("writes the path", func(t *testing.T) {
		is := is.New(t)
		p := filepath.Join(t.TempDir(), "test")
		t.Setenv(PATH_FILE_LOOKUP, p)
		value := "$HOME/.local/bin"

		err := AddPath(value)

		data, _ := os.ReadFile(p)

		is.NoErr(err)                               // should not error
		is.Equal(string(data), fmt.Sprintln(value)) // should be equal
	})
}
