package archive

import (
	"os"
	"path"
	"testing"

	"github.com/matryer/is"

	"github.com/action-stars/ghactl/internal/fileio"
)

func TestUnZip(t *testing.T) {
	tests := []struct {
		name      string
		zipFile   string
		dest      string
		wantDirs  []string
		wantFiles []string
		wantErr   bool
	}{
		{
			name:    "errors_if_the_file_does_not_exist",
			zipFile: "non-existent-file",
			wantErr: true,
		},
		{
			name:    "errors_if_the_destination_does_not_exist",
			zipFile: "../../testdata/file-and-dir.zip",
			dest:    "non-existent-dir",
			wantErr: true,
		},
		{
			name:      "can_extract_zip_file",
			zipFile:   "../../testdata/file-and-dir.zip",
			wantDirs:  []string{"dir"},
			wantFiles: []string{"small.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			dest := tt.dest
			if dest == "" {
				dest = t.TempDir()
			}

			err := UnZip(tt.zipFile, dest)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err) // should not error

			for _, d := range tt.wantDirs {
				dirExists, _ := fileio.DirExists(path.Join(dest, d))
				is.True(dirExists) // dir should exist
			}

			for _, f := range tt.wantFiles {
				fileExists, _ := fileio.FileExists(path.Join(dest, f))
				is.True(fileExists) // file should exist
			}
		})
	}
}

func Test_extractZip(t *testing.T) {
	tests := []struct {
		name      string
		entries   []testZipEntry
		wantDirs  []string
		wantFiles map[string]string
		wantErr   bool
	}{
		{
			name: "extracts_directories_and_files",
			entries: []testZipEntry{
				{name: "dir/", isDir: true},
				{name: "dir/file.txt", content: "hello"},
			},
			wantDirs:  []string{"dir"},
			wantFiles: map[string]string{"dir/file.txt": "hello"},
		},
		{
			name: "creates_parent_directories_for_files",
			entries: []testZipEntry{
				{name: "a/b/c.txt", content: "nested"},
			},
			wantFiles: map[string]string{"a/b/c.txt": "nested"},
		},
		{
			name: "rejects_path_traversal",
			entries: []testZipEntry{
				{name: "../escape.txt", content: "bad"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			dest := t.TempDir()

			root, err := os.OpenRoot(dest)
			is.NoErr(err) // should open root
			defer root.Close()

			files := newTestZip(t, tt.entries)
			err = extractZip(root, files)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err) // should not error

			for _, d := range tt.wantDirs {
				dirExists, _ := fileio.DirExists(path.Join(dest, d))
				is.True(dirExists) // dir should exist
			}

			for f, want := range tt.wantFiles {
				data, err := os.ReadFile(path.Join(dest, f))
				is.NoErr(err)                // should read file
				is.Equal(string(data), want) // content should match
			}
		})
	}
}

func Test_extractZipFile(t *testing.T) {
	tests := []struct {
		name     string
		entry    testZipEntry
		wantDir  string
		wantFile string
		wantData string
		wantErr  bool
	}{
		{
			name:    "creates_directory",
			entry:   testZipEntry{name: "testdir/", isDir: true},
			wantDir: "testdir",
		},
		{
			name:     "creates_regular_file",
			entry:    testZipEntry{name: "test.txt", content: "content"},
			wantFile: "test.txt",
			wantData: "content",
		},
		{
			name:     "creates_parent_directories",
			entry:    testZipEntry{name: "a/b/file.txt", content: "deep"},
			wantFile: "a/b/file.txt",
			wantData: "deep",
		},
		{
			name:    "rejects_path_traversal",
			entry:   testZipEntry{name: "../escape.txt", content: "bad"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			dest := t.TempDir()

			root, err := os.OpenRoot(dest)
			is.NoErr(err)
			defer root.Close()

			files := newTestZip(t, []testZipEntry{tt.entry})
			err = extractZipFile(root, files[0])

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err) // should not error

			if tt.wantDir != "" {
				dirExists, _ := fileio.DirExists(path.Join(dest, tt.wantDir))
				is.True(dirExists) // directory should exist
			}

			if tt.wantFile != "" {
				data, err := os.ReadFile(path.Join(dest, tt.wantFile))
				is.NoErr(err)                       // should read file
				is.Equal(string(data), tt.wantData) // content should match
			}
		})
	}
}
