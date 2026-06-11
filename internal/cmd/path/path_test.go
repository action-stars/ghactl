package path

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/urfave/cli/v3"
)

func TestNew_Add(t *testing.T) {
	t.Run("adds_path_entry", func(t *testing.T) {
		is := is.New(t)

		envFile := filepath.Join(t.TempDir(), "github-path")
		t.Setenv("GITHUB_PATH", envFile)
		t.Setenv("PATH", "existing-path")

		buf := new(bytes.Buffer)
		cmd := New()
		cmd.Writer = buf

		err := cmd.Run(context.Background(), []string{"path", "add", "--path", "$HOME/.local/bin"})

		data, readErr := os.ReadFile(envFile)

		is.NoErr(err)                                                     // should not error
		is.NoErr(readErr)                                                 // should not error
		is.True(strings.Contains(string(data), "$HOME/.local/bin"))       // should write path
		is.True(strings.HasPrefix(os.Getenv("PATH"), "$HOME/.local/bin")) // should prepend path
		is.Equal(buf.Len(), 0)                                            // should not output
	})

	t.Run("errors_when_github_path_not_set", func(t *testing.T) {
		is := is.New(t)

		t.Setenv("GITHUB_PATH", "")

		buf := new(bytes.Buffer)
		cmd := New()
		cmd.Writer = buf
		cmd.ExitErrHandler = func(_ context.Context, _ *cli.Command, _ error) {}

		err := cmd.Run(context.Background(), []string{"path", "add", "--path", "$HOME/.local/bin"})

		is.True(err != nil) // should error
	})
}
