package tool

import (
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestCmd_ExtractTar(t *testing.T) {
	c := &Cmd{}

	t.Run("errors_if_temp_dir_is_not_defined", func(t *testing.T) {
		is := is.New(t)
		t.Setenv("RUNNER_TEMP", "")

		_, err := c.ExtractTar("testdata/file.tar", false)

		is.True(err != nil) // should error
	})

	t.Run("errors_on_invalid_tar_file", func(t *testing.T) {
		is := is.New(t)
		setupTempDir(t)

		_, err := c.ExtractTar("non-existent-file", false)

		is.True(err != nil) // should error
	})

	t.Run("extracts_tar_file", func(t *testing.T) {
		is := is.New(t)
		setupTempDir(t)

		tarPath, _ := filepath.Abs("../../../testdata/file-and-dir.tar")
		p, err := c.ExtractTar(tarPath, false)

		is.NoErr(err)    // should not error
		is.True(p != "") // should return path
	})

	t.Run("extracts_tar.gz_file", func(t *testing.T) {
		is := is.New(t)
		setupTempDir(t)

		tarPath, _ := filepath.Abs("../../../testdata/file-and-dir.tar.gz")
		p, err := c.ExtractTar(tarPath, true)

		is.NoErr(err)    // should not error
		is.True(p != "") // should return path
	})
}

func TestCmd_ExtractZip(t *testing.T) {
	c := &Cmd{}

	t.Run("errors_if_temp_dir_is_not_defined", func(t *testing.T) {
		is := is.New(t)
		t.Setenv("RUNNER_TEMP", "")

		_, err := c.ExtractZip("testdata/file.zip")

		is.True(err != nil) // should error
	})

	t.Run("errors_on_invalid_zip_file", func(t *testing.T) {
		is := is.New(t)
		setupTempDir(t)

		_, err := c.ExtractZip("non-existent-file")

		is.True(err != nil) // should error
	})

	t.Run("extracts_zip_file", func(t *testing.T) {
		is := is.New(t)
		setupTempDir(t)

		zipPath, _ := filepath.Abs("../../../testdata/file-and-dir.zip")
		p, err := c.ExtractZip(zipPath)

		is.NoErr(err)    // should not error
		is.True(p != "") // should return path
	})
}
