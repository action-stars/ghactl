package util

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestDownloadFile(t *testing.T) {
	t.Run("errors if url is invalid", func(t *testing.T) {
		is := is.New(t)
		dest := filepath.Join(t.TempDir(), "test")

		err := DownloadFile("http://non-existent-url", dest, false)

		is.True(err != nil) // should error
	})

	t.Run("downloads file", func(t *testing.T) {
		is := is.New(t)
		dest := filepath.Join(t.TempDir(), "test")

		err := DownloadFile("https://dl.k8s.io/release/stable.txt", dest, false)

		is.NoErr(err) // should not error

		if err == nil {
			destContent, err := os.ReadFile(dest)
			if err != nil {
				t.Fatal(err)
			}

			is.True(strings.HasPrefix(string(destContent), "v")) // should match
		}
	})

	t.Run("errors if dest dir does not exist", func(t *testing.T) {
		is := is.New(t)
		dest := filepath.Join(t.TempDir(), "non-existent-dir", "test")

		err := DownloadFile("https://dl.k8s.io/release/stable.txt", dest, false)

		is.True(err != nil) // should error
	})

	t.Run("errors if dest already exists", func(t *testing.T) {
		is := is.New(t)
		dest := func() string {
			tmp, err := os.CreateTemp(t.TempDir(), "file")
			if err != nil {
				t.Fatal(err)
			}
			defer tmp.Close()
			_, err = tmp.WriteString("dest")
			if err != nil {
				t.Fatal(err)
			}
			return tmp.Name()
		}()

		err := DownloadFile("https://dl.k8s.io/release/stable.txt", dest, false)

		is.True(err != nil) // should error
	})

	t.Run("downloads file if dest already exists with override", func(t *testing.T) {
		is := is.New(t)
		dest := func() string {
			tmp, err := os.CreateTemp(t.TempDir(), "file")
			if err != nil {
				t.Fatal(err)
			}
			defer tmp.Close()
			_, err = tmp.WriteString("dest")
			if err != nil {
				t.Fatal(err)
			}
			return tmp.Name()
		}()

		err := DownloadFile("https://dl.k8s.io/release/stable.txt", dest, true)

		destContent, _ := os.ReadFile(dest)

		is.NoErr(err)                                        // should not error
		is.True(strings.HasPrefix(string(destContent), "v")) // should match
	})
}
