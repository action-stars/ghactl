package toolcache

import (
	"testing"

	"github.com/action-stars/ghactl/internal/toolkit/core"
	"github.com/action-stars/ghactl/internal/util"
	"github.com/matryer/is"
)

func TestDownloadTool(t *testing.T) {
	t.Run("errors if temp dir is not defined", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(core.RUNNER_TEMP_LOOKUP, "")

		_, err := DownloadTool("https://dl.k8s.io/release/stable.txt")

		is.True(err != nil) // should error
	})

	t.Run("downloads file to temp dir", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(core.RUNNER_TEMP_LOOKUP, t.TempDir())

		p, err := DownloadTool("https://dl.k8s.io/release/stable.txt")

		is.NoErr(err) // should not error

		if err == nil {
			fileExists, err := util.FileExists(p)
			if err != nil {
				t.Error(err)
			}

			is.True(fileExists) // should exist
		}
	})
}
