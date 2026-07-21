package toolcache

import (
	"path/filepath"
	"testing"

	"github.com/matryer/is"

	"github.com/action-stars/ghactl/internal/fileio"
)

func TestExtractTar(t *testing.T) {
	tests := []struct {
		name      string
		tempUnset bool
		file      string
		gz        bool
		wantErr   bool
	}{
		{
			name:      "errors_if_temp_dir_is_not_defined",
			tempUnset: true,
			file:      "../../../testdata/file-and-dir.tar",
			wantErr:   true,
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

			temp := ""
			if tt.tempUnset {
				temp = ""
			} else {
				temp = t.TempDir()
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
		name      string
		tempUnset bool
		file      string
		wantErr   bool
	}{
		{
			name:      "errors_if_temp_dir_is_not_defined",
			tempUnset: true,
			file:      "../../../testdata/file-and-dir.zip",
			wantErr:   true,
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

			temp := ""
			if tt.tempUnset {
				temp = ""
			} else {
				temp = t.TempDir()
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

func TestResolveToolDirectory(t *testing.T) {
	t.Parallel()

	t.Run("returns_path_for_empty_directory", func(t *testing.T) {
		t.Parallel()

		is := is.New(t)
		root := t.TempDir()

		result, err := ResolveToolDirectory(root)

		is.NoErr(err)
		is.Equal(result, root)
	})

	t.Run("returns_path_with_multiple_items", func(t *testing.T) {
		t.Parallel()

		is := is.New(t)
		root := t.TempDir()
		mustCreateTestFile(t, filepath.Join(root, "file1"), "content")
		mustCreateTestFile(t, filepath.Join(root, "file2"), "content")
		mustCreateTestDir(t, filepath.Join(root, "dir1"))

		result, err := ResolveToolDirectory(root)

		is.NoErr(err)
		is.Equal(result, root)
	})

	t.Run("steps_into_single_nested_directory", func(t *testing.T) {
		t.Parallel()

		is := is.New(t)
		root := t.TempDir()
		nested := filepath.Join(root, "tool-v1.0.0")
		mustCreateTestDir(t, nested)
		mustCreateTestFile(t, filepath.Join(nested, "tool.exe"), "")

		result, err := ResolveToolDirectory(root)

		is.NoErr(err)
		is.Equal(result, nested)
	})

	t.Run("steps_into_multiple_nested_directories", func(t *testing.T) {
		t.Parallel()

		is := is.New(t)
		root := t.TempDir()
		nested := filepath.Join(root, "tool-v1.0.0", "bin", "amd64")
		mustCreateTestDir(t, nested)
		mustCreateTestFile(t, filepath.Join(nested, "tool.exe"), "")

		result, err := ResolveToolDirectory(root)

		is.NoErr(err)
		is.Equal(result, nested)
	})

	t.Run("steps_into_nested_directory_with_bin_subdirectory", func(t *testing.T) {
		t.Parallel()

		is := is.New(t)
		root := t.TempDir()
		nested := filepath.Join(root, "jsonschema-v1.0.0")
		mustCreateTestDir(t, nested)
		bin := filepath.Join(nested, "bin")
		mustCreateTestDir(t, bin)
		mustCreateTestFile(t, filepath.Join(bin, "jsonschema"), "")
		mustCreateTestDir(t, filepath.Join(nested, "lib"))

		result, err := ResolveToolDirectory(root)

		is.NoErr(err)
		is.Equal(result, bin)
	})

	t.Run("steps_into_bin_if_multiple_items_at_root", func(t *testing.T) {
		t.Parallel()

		is := is.New(t)
		root := t.TempDir()
		bin := filepath.Join(root, "bin")
		mustCreateTestDir(t, bin)
		mustCreateTestDir(t, filepath.Join(root, "lib"))
		mustCreateTestFile(t, filepath.Join(root, "file.txt"), "")

		result, err := ResolveToolDirectory(root)

		is.NoErr(err)
		is.Equal(result, bin)
	})

	t.Run("handles_non_existent_path_error", func(t *testing.T) {
		t.Parallel()

		is := is.New(t)
		root := t.TempDir()
		nonExistentPath := filepath.Join(root, "does-not-exist")

		_, err := ResolveToolDirectory(nonExistentPath)

		is.True(err != nil)
	})

	t.Run("handles_file_path_error", func(t *testing.T) {
		t.Parallel()

		is := is.New(t)
		root := t.TempDir()
		filePath := filepath.Join(root, "file.txt")
		mustCreateTestFile(t, filePath, "content")

		_, err := ResolveToolDirectory(filePath)

		is.True(err != nil)
	})
}
