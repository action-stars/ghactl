package toolcache

import (
	"testing"

	"github.com/matryer/is"
)

func TestCheckVersion(t *testing.T) {
	t.Run("errors if version is empty", func(t *testing.T) {
		is := is.New(t)

		_, err := CheckVersion("", "*")

		is.True(err != nil) // should error
	})

	t.Run("errors if version is invalid", func(t *testing.T) {
		is := is.New(t)

		_, err := CheckVersion("a.b.c", "*")

		is.True(err != nil) // should error
	})

	t.Run("errors if version spec is empty", func(t *testing.T) {
		is := is.New(t)

		_, err := CheckVersion("1.0.0", "")

		is.True(err != nil) // should error
	})

	t.Run("errors if version spec is invalid", func(t *testing.T) {
		is := is.New(t)

		_, err := CheckVersion("1.0.0", "test")

		is.True(err != nil) // should error
	})

	t.Run("matches explicit version", func(t *testing.T) {
		is := is.New(t)

		valid, err := CheckVersion("1.0.0", "1.0.0")

		is.NoErr(err)  // should not error
		is.True(valid) // should be valid
	})

	t.Run("blocks invalid explicit version", func(t *testing.T) {
		is := is.New(t)

		valid, err := CheckVersion("1.0.0", "1.0.1")

		is.NoErr(err)   // should not error
		is.True(!valid) // should not be valid
	})

	t.Run("matches patch version constraint", func(t *testing.T) {
		is := is.New(t)

		valid, err := CheckVersion("1.0.1", "~1.0.0")

		is.NoErr(err)  // should not error
		is.True(valid) // should be valid
	})

	t.Run("blocks invalid patch version constraint", func(t *testing.T) {
		is := is.New(t)

		valid, err := CheckVersion("1.0.1", "~1.1.0")

		is.NoErr(err)   // should not error
		is.True(!valid) // should not be valid
	})

	t.Run("matches minor version constraint", func(t *testing.T) {
		is := is.New(t)

		valid, err := CheckVersion("1.1.1", "^1.0.0")

		is.NoErr(err)  // should not error
		is.True(valid) // should be valid
	})

	t.Run("blocks invalid minor version constraint", func(t *testing.T) {
		is := is.New(t)

		valid, err := CheckVersion("1.1.1", "^2.0.0")

		is.NoErr(err)   // should not error
		is.True(!valid) // should not be valid
	})

	t.Run("matches greater than version constraint", func(t *testing.T) {
		is := is.New(t)

		valid, err := CheckVersion("1.1.1", ">=1.0.0")

		is.NoErr(err)  // should not error
		is.True(valid) // should be valid
	})

	t.Run("blocks invalid greater than version constraint", func(t *testing.T) {
		is := is.New(t)

		valid, err := CheckVersion("1.1.1", ">=1.2.0")

		is.NoErr(err)   // should not error
		is.True(!valid) // should not be valid
	})
}
