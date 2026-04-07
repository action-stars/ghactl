package toolcache

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/matryer/is"

	"github.com/action-stars/ghactl/internal/fileio"
)

func TestDownloadTool(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/missing" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Fprint(w, "file content")
	}))
	defer ts.Close()

	tests := []struct {
		name     string
		temp     string
		url      string
		wantFile bool
		wantErr  bool
	}{
		{
			name:    "errors_if_temp_dir_is_not_defined",
			temp:    "",
			url:     "http://example.com/file",
			wantErr: true,
		},
		{
			name:     "downloads_file_to_temp_dir",
			url:      ts.URL + "/file",
			wantFile: true,
		},
		{
			name:    "errors_on_non-2xx_status",
			url:     ts.URL + "/missing",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			temp := tt.temp
			if temp == "" && !tt.wantErr || tt.wantErr && tt.url != "http://example.com/file" {
				temp = t.TempDir()
			}
			t.Setenv("RUNNER_TEMP", temp)
			t.Setenv("TEMPDIR", "")
			t.Setenv("TMP", "")

			u, _ := url.Parse(tt.url)
			p, err := DownloadTool(t.Context(), nil, *u)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err) // should not error

			if tt.wantFile {
				fileExists, _ := fileio.FileExists(p)
				is.True(fileExists) // should exist
			}
		})
	}
}
