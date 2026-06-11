package github

import (
	"testing"

	"github.com/matryer/is"
)

func TestGetClient(t *testing.T) {
	t.Run("returns_a_client_without_auth", func(t *testing.T) {
		is := is.New(t)

		client, err := GetClient(nil)

		is.NoErr(err)
		is.True(client != nil)
	})

	t.Run("returns_a_client_with_an_auth_token", func(t *testing.T) {
		is := is.New(t)
		token := "test-token"

		client, err := GetClient(&token)

		is.NoErr(err)
		is.True(client != nil)
	})
}
