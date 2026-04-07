package core

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestAddMatcher(t *testing.T) {
	matcherFile := func() string {
		tmp, err := os.CreateTemp(t.TempDir(), "problem-matcher.json")
		if err != nil {
			t.Fatal(err)
		}
		tmp.Close()
		return tmp.Name()
	}()

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "errors_if_the_file_does_not_exist",
			path:    "non-existent-file",
			wantErr: true,
		},
		{
			name: "writes_the_add_matcher_command",
			path: matcherFile,
			want: fmt.Sprintf("::%s::%s\n", AddMatcherCmd, matcherFile),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			var b bytes.Buffer
			err := AddMatcher(&b, tt.path)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)                 // should not error
			is.Equal(b.String(), tt.want) // should be equal
		})
	}
}

func TestRemoveMatcher(t *testing.T) {
	tests := []struct {
		name  string
		owner string
		want  string
	}{
		{
			name:  "writes_the_remove_matcher_command",
			owner: "ownerx",
			want:  fmt.Sprintf("::%s owner=%s::\n", RemoveMatcherCmd, "ownerx"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			var b bytes.Buffer
			err := RemoveMatcher(&b, tt.owner)

			is.NoErr(err)                 // should not error
			is.Equal(b.String(), tt.want) // should be equal
		})
	}
}
