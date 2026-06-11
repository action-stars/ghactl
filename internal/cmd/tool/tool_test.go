package tool

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/urfave/cli/v3"
)

func TestNew_CacheGet(t *testing.T) {
	t.Run("outputs_tool_cache_directory", func(t *testing.T) {
		is := is.New(t)
		tc := setupToolCache(t)

		buf := new(bytes.Buffer)
		cmd := New()
		cmd.Writer = buf

		err := cmd.Run(context.Background(), []string{"tool", "cache", "get"})

		is.NoErr(err)                                 // should not error
		is.Equal(strings.TrimSpace(buf.String()), tc) // should output cache dir
	})

	t.Run("errors_when_env_not_set", func(t *testing.T) {
		is := is.New(t)
		t.Setenv("RUNNER_TOOL_CACHE", "")

		buf := new(bytes.Buffer)
		cmd := New()
		cmd.Writer = buf
		cmd.ExitErrHandler = func(_ context.Context, _ *cli.Command, _ error) {}

		err := cmd.Run(context.Background(), []string{"tool", "cache", "get"})

		is.True(err != nil) // should error
	})
}

func TestNew_CacheFind(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		t.Run("outputs_found_tool_versions", func(t *testing.T) {
			is := is.New(t)
			setupToolCache(t)

			buf := new(bytes.Buffer)
			cmd := New()
			cmd.Writer = buf

			err := cmd.Run(context.Background(), []string{"tool", "cache", "find", "--all", "--name", "test-tool", "--arch", "amd64"})

			is.NoErr(err)          // should not error
			is.True(buf.Len() > 0) // should output versions
		})

		t.Run("produces_no_output_for_missing_tool", func(t *testing.T) {
			is := is.New(t)
			setupToolCache(t)

			buf := new(bytes.Buffer)
			cmd := New()
			cmd.Writer = buf

			err := cmd.Run(context.Background(), []string{"tool", "cache", "find", "--all", "--name", "nonexistent", "--arch", "amd64"})

			is.NoErr(err)          // should not error
			is.Equal(buf.Len(), 0) // should produce no output
		})

		t.Run("filters_versions_with_version_spec", func(t *testing.T) {
			is := is.New(t)
			setupToolCache(t)

			buf := new(bytes.Buffer)
			cmd := New()
			cmd.Writer = buf

			err := cmd.Run(context.Background(), []string{"tool", "cache", "find", "--all", "--name", "test-tool", "--arch", "amd64", "--version", "^1.0.0"})

			is.NoErr(err) // should not error

			lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
			is.Equal(lines, []string{"1.0.0", "1.0.1", "1.2.0"}) // should include only matching versions
		})

		t.Run("errors_for_invalid_version_spec", func(t *testing.T) {
			is := is.New(t)
			setupToolCache(t)

			buf := new(bytes.Buffer)
			cmd := New()
			cmd.Writer = buf
			cmd.ExitErrHandler = func(_ context.Context, _ *cli.Command, _ error) {}

			err := cmd.Run(context.Background(), []string{"tool", "cache", "find", "--all", "--name", "test-tool", "--arch", "amd64", "--version", "bad"})

			is.True(err != nil) // should error
		})
	})

	t.Run("outputs_found_tool_path", func(t *testing.T) {
		is := is.New(t)
		setupToolCache(t)

		buf := new(bytes.Buffer)
		cmd := New()
		cmd.Writer = buf

		err := cmd.Run(context.Background(), []string{"tool", "cache", "find", "--name", "test-tool", "--arch", "amd64", "--version", "^1.0.0"})

		is.NoErr(err)          // should not error
		is.True(buf.Len() > 0) // should output path
	})

	t.Run("produces_no_output_for_non-matching_version", func(t *testing.T) {
		is := is.New(t)
		setupToolCache(t)

		buf := new(bytes.Buffer)
		cmd := New()
		cmd.Writer = buf

		err := cmd.Run(context.Background(), []string{"tool", "cache", "find", "--name", "test-tool", "--arch", "amd64", "--version", "^99.0.0"})

		is.NoErr(err)          // should not error
		is.Equal(buf.Len(), 0) // should produce no output
	})
}

func TestNew_Download(t *testing.T) {
	t.Run("outputs_downloaded_file_path", func(t *testing.T) {
		is := is.New(t)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprint(w, "tool-binary")
		}))
		defer ts.Close()

		setupTempDir(t)

		buf := new(bytes.Buffer)
		cmd := New()
		cmd.Writer = buf

		err := cmd.Run(context.Background(), []string{"tool", "download", "--url", ts.URL + "/tool"})

		is.NoErr(err)          // should not error
		is.True(buf.Len() > 0) // should output path
	})
}

func TestNew_ExtractTar(t *testing.T) {
	t.Run("outputs_extracted_directory_path", func(t *testing.T) {
		is := is.New(t)
		setupTempDir(t)

		tarPath, _ := filepath.Abs("../../../testdata/file-and-dir.tar")

		buf := new(bytes.Buffer)
		cmd := New()
		cmd.Writer = buf

		err := cmd.Run(context.Background(), []string{"tool", "extract", "tar", "--path", tarPath})

		is.NoErr(err)          // should not error
		is.True(buf.Len() > 0) // should output path
	})
}

func TestNew_ExtractTgz(t *testing.T) {
	t.Run("outputs_extracted_directory_path", func(t *testing.T) {
		is := is.New(t)
		setupTempDir(t)

		tgzPath, _ := filepath.Abs("../../../testdata/file-and-dir.tar.gz")

		buf := new(bytes.Buffer)
		cmd := New()
		cmd.Writer = buf

		err := cmd.Run(context.Background(), []string{"tool", "extract", "tgz", "--path", tgzPath})

		is.NoErr(err)          // should not error
		is.True(buf.Len() > 0) // should output path
	})
}

func TestNew_ExtractZip(t *testing.T) {
	t.Run("outputs_extracted_directory_path", func(t *testing.T) {
		is := is.New(t)
		setupTempDir(t)

		zipPath, _ := filepath.Abs("../../../testdata/file-and-dir.zip")

		buf := new(bytes.Buffer)
		cmd := New()
		cmd.Writer = buf

		err := cmd.Run(context.Background(), []string{"tool", "extract", "zip", "--path", zipPath})

		is.NoErr(err)          // should not error
		is.True(buf.Len() > 0) // should output path
	})
}

func TestNew_VersionCheck(t *testing.T) {
	t.Run("outputs_true_for_matching_version", func(t *testing.T) {
		is := is.New(t)

		buf := new(bytes.Buffer)
		cmd := New()
		cmd.Writer = buf

		err := cmd.Run(context.Background(), []string{"tool", "version", "check", "--version", "1.0.0", "--version-spec", "^1.0.0"})

		is.NoErr(err)                                     // should not error
		is.Equal(strings.TrimSpace(buf.String()), "true") // should output true
	})

	t.Run("outputs_false_for_non-matching_version", func(t *testing.T) {
		is := is.New(t)

		buf := new(bytes.Buffer)
		cmd := New()
		cmd.Writer = buf

		err := cmd.Run(context.Background(), []string{"tool", "version", "check", "--version", "2.0.0", "--version-spec", "^1.0.0"})

		is.NoErr(err)                                      // should not error
		is.Equal(strings.TrimSpace(buf.String()), "false") // should output false
	})

	t.Run("errors_for_invalid_version", func(t *testing.T) {
		is := is.New(t)

		buf := new(bytes.Buffer)
		cmd := New()
		cmd.Writer = buf
		cmd.ExitErrHandler = func(_ context.Context, _ *cli.Command, _ error) {}

		err := cmd.Run(context.Background(), []string{"tool", "version", "check", "--version", "bad", "--version-spec", "^1.0.0"})

		is.True(err != nil) // should error
	})
}
