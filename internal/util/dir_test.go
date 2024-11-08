package util

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestDirExists(t *testing.T) {
	t.Run("returns false if the path does not exist", func(t *testing.T) {
		is := is.New(t)

		exists, err := DirExists("non-existent-directory")

		is.NoErr(err)    // should not error
		is.True(!exists) // should not exists
	})

	t.Run("returns true if the path is a dir", func(t *testing.T) {
		is := is.New(t)

		exists, err := DirExists(t.TempDir())

		is.NoErr(err)   // should not error
		is.True(exists) // should exist
	})

	t.Run("errors if the path is a file", func(t *testing.T) {
		is := is.New(t)
		tmp, err := os.CreateTemp(t.TempDir(), "file")
		if err != nil {
			t.Fatal(err)
		}
		tmp.Close()

		_, err = DirExists(tmp.Name())

		is.True(err != nil) // should error
	})
}
