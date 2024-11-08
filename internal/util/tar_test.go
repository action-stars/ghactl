package util

import (
	"path"
	"testing"

	"github.com/matryer/is"
)

func TestUnTar(t *testing.T) {
	t.Run("errors if the file does not exist", func(t *testing.T) {
		is := is.New(t)

		err := UnTar("non-existent-file", t.TempDir(), false)

		is.True(err != nil) // should error
	})

	t.Run("errors if the destination does not exist", func(t *testing.T) {
		is := is.New(t)

		err := UnTar("../../testdata/file-and-dir.tar", "non-existent-dir", false)

		is.True(err != nil) // should error
	})

	t.Run("can extract tar file", func(t *testing.T) {
		is := is.New(t)
		dest := t.TempDir()

		err := UnTar("../../testdata/file-and-dir.tar", dest, false)

		dirExists, _ := DirExists(path.Join(dest, "dir"))
		fileExists, _ := FileExists(path.Join(dest, "small.txt"))

		is.NoErr(err)       // should not error
		is.True(dirExists)  // dir should exist
		is.True(fileExists) // dir should exist
	})

	t.Run("can extract tar gz file", func(t *testing.T) {
		is := is.New(t)
		dest := t.TempDir()

		err := UnTar("../../testdata/file-and-dir.tar.gz", dest, true)

		dirExists, _ := DirExists(path.Join(dest, "dir"))
		fileExists, _ := FileExists(path.Join(dest, "small.txt"))

		is.NoErr(err)       // should not error
		is.True(dirExists)  // dir should exist
		is.True(fileExists) // dir should exist
	})
}
