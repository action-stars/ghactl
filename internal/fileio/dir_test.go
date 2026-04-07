package fileio

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestDirExists(t *testing.T) {
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
			path: "non-existent-directory",
			want: false,
		},
		{
			name: "returns_true_if_the_path_is_a_dir",
			path: dir,
			want: true,
		},
		{
			name:    "errors_if_the_path_is_a_file",
			path:    file.Name(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			got, err := DirExists(tt.path)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)          // should not error
			is.Equal(got, tt.want) // existence check
		})
	}
}
