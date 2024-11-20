package github

import (
	"testing"

	"github.com/matryer/is"
)

func TestGetClient(t *testing.T) {
	t.Run("returns a client", func(t *testing.T) {
		is := is.New(t)

		client := GetClient("")

		is.True(client != nil)
	})

	t.Run("returns a client with an auth token", func(t *testing.T) {
		is := is.New(t)

		client := GetClient("token")

		is.True(client != nil)
	})
}
