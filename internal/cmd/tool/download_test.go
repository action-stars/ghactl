package tool

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestCmd_Download(t *testing.T) {
	c := &Cmd{}

	t.Run("errors_if_URL_is_invalid", func(t *testing.T) {
		is := is.New(t)
		setupTempDir(t)

		_, err := c.Download(t.Context(), "://invalid")

		is.True(err != nil) // should error
	})

	t.Run("errors_if_temp_dir_is_not_defined", func(t *testing.T) {
		is := is.New(t)
		t.Setenv("RUNNER_TEMP", "")

		_, err := c.Download(t.Context(), "http://example.com/file")

		is.True(err != nil) // should error
	})

	t.Run("downloads_file_successfully", func(t *testing.T) {
		is := is.New(t)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprint(w, "file content")
		}))
		defer ts.Close()

		setupTempDir(t)

		p, err := c.Download(t.Context(), ts.URL+"/tool")

		is.NoErr(err)    // should not error
		is.True(p != "") // should return path
	})

	t.Run("errors_on_HTTP_failure", func(t *testing.T) {
		is := is.New(t)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer ts.Close()

		setupTempDir(t)

		_, err := c.Download(t.Context(), ts.URL+"/tool")

		is.True(err != nil) // should error
	})
}
