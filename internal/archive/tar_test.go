package archive

import (
	"archive/tar"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/matryer/is"

	"github.com/action-stars/ghactl/internal/fileio"
)

func TestUnTar(t *testing.T) {
	tests := []struct {
		name      string
		tarFile   string
		dest      string
		gz        bool
		wantDirs  []string
		wantFiles []string
		wantErr   bool
	}{
		{
			name:    "errors_if_the_file_does_not_exist",
			tarFile: "non-existent-file",
			wantErr: true,
		},
		{
			name:    "errors_if_the_destination_does_not_exist",
			tarFile: "../../testdata/file-and-dir.tar",
			dest:    "non-existent-dir",
			wantErr: true,
		},
		{
			name:      "can_extract_tar_file",
			tarFile:   "../../testdata/file-and-dir.tar",
			gz:        false,
			wantDirs:  []string{"dir"},
			wantFiles: []string{"small.txt"},
		},
		{
			name:      "can_extract_tar_gz_file",
			tarFile:   "../../testdata/file-and-dir.tar.gz",
			gz:        true,
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

			err := UnTar(tt.tarFile, dest, tt.gz)

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

func Test_extractTar(t *testing.T) {
	tests := []struct {
		name      string
		entries   []testTarEntry
		wantDirs  []string
		wantFiles map[string]string
		wantErr   bool
	}{
		{
			name: "extracts_directories_and_files",
			entries: []testTarEntry{
				{name: "dir/", typeflag: tar.TypeDir, mode: 0o755},
				{name: "dir/file.txt", typeflag: tar.TypeReg, mode: 0o644, content: "hello"},
			},
			wantDirs:  []string{"dir"},
			wantFiles: map[string]string{"dir/file.txt": "hello"},
		},
		{
			name: "creates_parent_directories_for_files",
			entries: []testTarEntry{
				{name: "a/b/c.txt", typeflag: tar.TypeReg, mode: 0o644, content: "nested"},
			},
			wantFiles: map[string]string{"a/b/c.txt": "nested"},
		},
		{
			name: "rejects_path_traversal",
			entries: []testTarEntry{
				{name: "../escape.txt", typeflag: tar.TypeReg, mode: 0o644, content: "bad"},
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

			tr := newTestTar(t, tt.entries)
			err = extractTar(root, tr)

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

func Test_extractTarEntry(t *testing.T) {
	tests := []struct {
		name     string
		header   *tar.Header
		content  string
		wantDir  string
		wantFile string
		wantData string
		wantErr  bool
	}{
		{
			name:    "creates_directory",
			header:  &tar.Header{Name: "testdir/", Typeflag: tar.TypeDir, Mode: 0o755},
			wantDir: "testdir",
		},
		{
			name:     "creates_regular_file",
			header:   &tar.Header{Name: "test.txt", Typeflag: tar.TypeReg, Mode: 0o644},
			content:  "content",
			wantFile: "test.txt",
			wantData: "content",
		},
		{
			name:     "creates_parent_directories",
			header:   &tar.Header{Name: "a/b/file.txt", Typeflag: tar.TypeReg, Mode: 0o644},
			content:  "deep",
			wantFile: "a/b/file.txt",
			wantData: "deep",
		},
		{
			name:     "strips_leading_slash",
			header:   &tar.Header{Name: "/absolute.txt", Typeflag: tar.TypeReg, Mode: 0o644},
			content:  "abs",
			wantFile: "absolute.txt",
			wantData: "abs",
		},
		{
			name:    "rejects_path_traversal",
			header:  &tar.Header{Name: "../../etc/passwd", Typeflag: tar.TypeReg, Mode: 0o644},
			content: "bad",
			wantErr: true,
		},
		{
			name:   "skips_empty_name",
			header: &tar.Header{Name: "/", Typeflag: tar.TypeDir, Mode: 0o755},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			dest := t.TempDir()

			root, err := os.OpenRoot(dest)
			is.NoErr(err)
			defer root.Close()

			err = extractTarEntry(root, tt.header, strings.NewReader(tt.content))

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
