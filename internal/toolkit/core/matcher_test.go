package core

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestAddMatcher(t *testing.T) {
	t.Run("errors if the file does not exist", func(t *testing.T) {
		is := is.New(t)

		var b bytes.Buffer
		err := AddMatcher(&b, "non-existent-file")

		is.True(err != nil) // should error
	})

	t.Run("writes the add matcher command", func(t *testing.T) {
		is := is.New(t)
		p := func() string {
			tmp, err := os.CreateTemp(t.TempDir(), "problem-matcher.json")
			if err != nil {
				t.Fatal(err)
			}
			tmp.Close()
			return tmp.Name()
		}()

		var b bytes.Buffer
		err := AddMatcher(&b, p)

		is.NoErr(err)                                                     // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s::%s\n", AddMatcherCmd, p)) // should be equal
	})
}

func TestRemoveMatcher(t *testing.T) {
	t.Run("writes the remove matcher command", func(t *testing.T) {
		is := is.New(t)
		owner := "ownerx"

		var b bytes.Buffer
		err := RemoveMatcher(&b, owner)

		is.NoErr(err)                                                                   // should not error
		is.Equal(b.String(), fmt.Sprintf("::%s owner=%s::\n", RemoveMatcherCmd, owner)) // should be equal
	})
}
