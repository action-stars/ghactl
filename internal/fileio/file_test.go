package fileio

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestFileExists(t *testing.T) {
	dir := t.TempDir()

	file, err := os.CreateTemp(dir, "file")
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	tests := []struct {
		name    string
		path    string
		want    bool
		wantErr bool
	}{
		{
			name: "returns_false_if_the_path_does_not_exist",
			path: "non-existent-file",
			want: false,
		},
		{
			name: "returns_true_if_the_path_is_a_file",
			path: file.Name(),
			want: true,
		},
		{
			name:    "errors_if_the_path_is_a_directory",
			path:    dir,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			got, err := FileExists(tt.path)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)          // should not error
			is.Equal(got, tt.want) // existence check
		})
	}
}

func TestCopyFile(t *testing.T) {
	src := createTempFile(t, "src")
	existingDest := createTempFile(t, "dest")
	emptyDir := t.TempDir()
	emptyPath := filepath.Join(t.TempDir(), "file")
	overwriteDest := createTempFile(t, "dest")

	tests := []struct {
		name      string
		src       string
		dest      string
		overwrite bool
		wantErr   bool
		wantData  string
	}{
		{
			name:    "errors_if_the_source_path_does_not_exist",
			src:     "non-existent-file",
			dest:    t.TempDir(),
			wantErr: true,
		},
		{
			name:    "errors_if_the_destination_path_exists",
			src:     src,
			dest:    existingDest,
			wantErr: true,
		},
		{
			name:     "copies_to_empty_dir",
			src:      src,
			dest:     emptyDir,
			wantData: "src",
		},
		{
			name:     "copies_to_empty_path",
			src:      src,
			dest:     emptyPath,
			wantData: "src",
		},
		{
			name:      "copies_to_existing_path_with_override",
			src:       src,
			dest:      overwriteDest,
			overwrite: true,
			wantData:  "src",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			err := CopyFile(tt.src, tt.dest, tt.overwrite)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err) // should not error

			// If dest is a directory, the file is copied into it.
			dest := tt.dest
			if fi, statErr := os.Stat(dest); statErr == nil && fi.IsDir() {
				dest = filepath.Join(dest, filepath.Base(tt.src))
			}

			destContent, _ := os.ReadFile(dest)
			is.Equal(string(destContent), tt.wantData) // content should match
		})
	}
}

func TestWriteFile(t *testing.T) {
	validPath := filepath.Join(t.TempDir(), "test")
	appendPath := filepath.Join(t.TempDir(), "test")
	_ = WriteFile(appendPath, []byte("value\n"))

	tests := []struct {
		name     string
		path     string
		value    string
		wantData string
		wantErr  bool
	}{
		{
			name:    "errors_if_path_is_invalid",
			path:    "/invalid-file-path",
			value:   "value",
			wantErr: true,
		},
		{
			name:     "can_write_a_value",
			path:     validPath,
			value:    "value\n",
			wantData: "value\n",
		},
		{
			name:     "can_append_value",
			path:     appendPath,
			value:    "value\n",
			wantData: "value\nvalue\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			err := WriteFile(tt.path, []byte(tt.value))

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err) // should not error

			data, _ := os.ReadFile(tt.path)
			is.Equal(string(data), tt.wantData) // should match
		})
	}
}
