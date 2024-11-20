package core

import (
	"testing"

	"github.com/matryer/is"
)

func TestIsDebug(t *testing.T) {
	t.Run("returns false if env is unset", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(RUNNER_DEBUG, "")

		result := IsDebug()

		is.Equal(result, false) // should match error
	})

	t.Run("returns true if env is set", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(RUNNER_DEBUG, "1")

		result := IsDebug()

		is.Equal(result, true) // should match error
	})
}
