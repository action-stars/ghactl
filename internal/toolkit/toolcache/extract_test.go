package toolcache

import (
	"testing"

	"github.com/action-stars/ghactl/internal/toolkit/core"
	"github.com/action-stars/ghactl/internal/util"
	"github.com/matryer/is"
)

func TestExtractTar(t *testing.T) {
	t.Run("errors if temp dir is not defined", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(core.RUNNER_TEMP_LOOKUP, "")

		_, err := ExtractTar("../../../testdata/file-and-dir.tar", false)

		is.True(err != nil) // should error
	})

	t.Run("errors if cannot extract from invalid tar file", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(core.RUNNER_TEMP_LOOKUP, t.TempDir())

		_, err := ExtractTar("non-existent-file", false)

		is.True(err != nil) // should error
	})

	t.Run("extracts tar file", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(core.RUNNER_TEMP_LOOKUP, t.TempDir())

		p, err := ExtractTar("../../../testdata/file-and-dir.tar", false)

		dirExists, _ := util.DirExists(p)

		is.NoErr(err)      // should not error
		is.True(dirExists) // extracted dir should exists
	})

	t.Run("extracts tar gz file", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(core.RUNNER_TEMP_LOOKUP, t.TempDir())

		p, err := ExtractTar("../../../testdata/file-and-dir.tar.gz", true)

		dirExists, _ := util.DirExists(p)

		is.NoErr(err)      // should not error
		is.True(dirExists) // extracted dir should exists
	})
}

func TestExtractZip(t *testing.T) {
	t.Run("errors if temp dir is not defined", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(core.RUNNER_TEMP_LOOKUP, "")

		_, err := ExtractZip("../../../testdata/file-and-dir.zip")

		is.True(err != nil) // should error
	})

	t.Run("errors if cannot extract from invalid zip file", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(core.RUNNER_TEMP_LOOKUP, t.TempDir())

		_, err := ExtractZip("non-existent-file")

		is.True(err != nil) // should error
	})

	t.Run("extracts zip file", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(core.RUNNER_TEMP_LOOKUP, t.TempDir())

		p, err := ExtractZip("../../../testdata/file-and-dir.zip")

		dirExists, _ := util.DirExists(p)

		is.NoErr(err)      // should not error
		is.True(dirExists) // extracted dir should exists
	})
}
