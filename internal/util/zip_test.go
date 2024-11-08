package util

import (
	"path"
	"testing"

	"github.com/matryer/is"
)

func TestUnZip(t *testing.T) {
	t.Run("errors if the file does not exist", func(t *testing.T) {
		is := is.New(t)

		err := UnZip("non-existent-file", t.TempDir())

		is.True(err != nil) // should error
	})

	t.Run("errors if the destination does not exist", func(t *testing.T) {
		is := is.New(t)

		err := UnZip("../../testdata/file-and-dir.zip", "non-existent-dir")

		is.True(err != nil) // should error
	})

	t.Run("can extract zip file", func(t *testing.T) {
		is := is.New(t)
		dest := t.TempDir()

		err := UnZip("../../testdata/file-and-dir.zip", dest)

		dirExists, _ := DirExists(path.Join(dest, "dir"))
		fileExists, _ := FileExists(path.Join(dest, "small.txt"))

		is.NoErr(err)       // should not error
		is.True(dirExists)  // dir should exist
		is.True(fileExists) // dir should exist
	})
}
