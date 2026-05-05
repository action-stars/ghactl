package toolcache

import (
	"testing"

	"github.com/matryer/is"

	"github.com/action-stars/ghactl/internal/fileio"
)

func TestExtractTar(t *testing.T) {
	tests := []struct {
		name    string
		temp    string
		file    string
		gz      bool
		wantErr bool
	}{
		{
			name:    "errors_if_temp_dir_is_not_defined",
			temp:    "UNSET",
			file:    "../../../testdata/file-and-dir.tar",
			wantErr: true,
		},
		{
			name:    "errors_if_cannot_extract_from_invalid_tar_file",
			file:    "non-existent-file",
			wantErr: true,
		},
		{
			name: "extracts_tar_file",
			file: "../../../testdata/file-and-dir.tar",
		},
		{
			name: "extracts_tar_gz_file",
			file: "../../../testdata/file-and-dir.tar.gz",
			gz:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			temp := tt.temp
			switch temp {
			case "":
				temp = t.TempDir()
			case "UNSET":
				temp = ""
			}
			t.Setenv("RUNNER_TEMP", temp)
			t.Setenv("TEMPDIR", "")
			t.Setenv("TMP", "")

			p, err := ExtractTar(tt.file, tt.gz)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err) // should not error

			dirExists, _ := fileio.DirExists(p)
			is.True(dirExists) // extracted dir should exist
		})
	}
}

func TestExtractZip(t *testing.T) {
	tests := []struct {
		name    string
		temp    string
		file    string
		wantErr bool
	}{
		{
			name:    "errors_if_temp_dir_is_not_defined",
			temp:    "UNSET",
			file:    "../../../testdata/file-and-dir.zip",
			wantErr: true,
		},
		{
			name:    "errors_if_cannot_extract_from_invalid_zip_file",
			file:    "non-existent-file",
			wantErr: true,
		},
		{
			name: "extracts_zip_file",
			file: "../../../testdata/file-and-dir.zip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			temp := tt.temp
			switch temp {
			case "":
				temp = t.TempDir()
			case "UNSET":
				temp = ""
			}
			t.Setenv("RUNNER_TEMP", temp)
			t.Setenv("TEMPDIR", "")
			t.Setenv("TMP", "")

			p, err := ExtractZip(tt.file)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err) // should not error

			dirExists, _ := fileio.DirExists(p)
			is.True(dirExists) // extracted dir should exist
		})
	}
}
