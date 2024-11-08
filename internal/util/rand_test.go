package util

import (
	"testing"

	"github.com/matryer/is"
)

func TestGenerateRandomString(t *testing.T) {
	t.Run("has the correct length", func(t *testing.T) {
		is := is.New(t)

		length := 10
		randomString := GenerateRandomString(length)

		is.Equal(len(randomString), length) // should have a known size
	})

	t.Run("does not match", func(t *testing.T) {
		is := is.New(t)

		length := 10
		randomStringA := GenerateRandomString(length)
		randomStringB := GenerateRandomString(length)

		is.True(randomStringA != randomStringB) // should be different
	})
}
