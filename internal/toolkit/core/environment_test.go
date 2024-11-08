package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestExportVariable(t *testing.T) {
	t.Run("errors if file env variable is not defined", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(ENV_FILE_LOOKUP, "")

		err := ExportVariable("key", "value")

		is.True(err != nil) // should error
	})

	t.Run("writes single line entry", func(t *testing.T) {
		is := is.New(t)
		p := filepath.Join(t.TempDir(), "test")
		t.Setenv(ENV_FILE_LOOKUP, p)
		k := "key"
		v := "value"

		err := ExportVariable(k, v)

		data, _ := os.ReadFile(p)

		is.NoErr(err)                                        // should not error
		is.Equal(string(data), fmt.Sprintf("%s=%s\n", k, v)) // should match
	})

	t.Run("writes multi line entry", func(t *testing.T) {
		is := is.New(t)
		p := filepath.Join(t.TempDir(), "test")
		t.Setenv(ENV_FILE_LOOKUP, p)
		k := "key"
		v := `hello
    world`

		err := ExportVariable(k, v)

		data, _ := os.ReadFile(p)

		is.NoErr(err)                              // should not error
		is.True(strings.Contains(string(data), k)) // should contain key
		is.True(strings.Contains(string(data), v)) // should contain value
	})
}
