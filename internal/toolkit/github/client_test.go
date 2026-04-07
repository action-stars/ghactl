package github

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/matryer/is"
)

func TestGetClient(t *testing.T) {
	t.Run("returns_a_client_without_auth", func(t *testing.T) {
		is := is.New(t)

		client := GetClient("")

		is.True(client != nil)
	})

	t.Run("returns_a_client_with_an_auth_token", func(t *testing.T) {
		is := is.New(t)

		client := GetClient("test-token")

		is.True(client != nil)
	})

	t.Run("authenticated_client_sends_auth_header", func(t *testing.T) {
		is := is.New(t)
		var authHeader string
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader = r.Header.Get("Authorization")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[]`))
		}))
		defer ts.Close()

		client := GetClient("test-token")
		client.BaseURL, _ = url.Parse(ts.URL + "/")

		_, _, err := client.Repositories.ListByUser(t.Context(), "test", nil)
		is.NoErr(err)                             // should not error
		is.Equal(authHeader, "Bearer test-token") // should send auth header
	})

	t.Run("unauthenticated_client_does_not_send_auth_header", func(t *testing.T) {
		is := is.New(t)
		var authHeader string
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader = r.Header.Get("Authorization")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[]`))
		}))
		defer ts.Close()

		client := GetClient("")
		client.BaseURL, _ = url.Parse(ts.URL + "/")

		_, _, err := client.Repositories.ListByUser(t.Context(), "test", nil)
		is.NoErr(err)            // should not error
		is.Equal(authHeader, "") // should not send auth header
	})
}
