package github

import "github.com/google/go-github/v66/github"

// GetClient returns a GitHub client.
// If a token is provided, it will be used for authentication.
func GetClient(token string) *github.Client {
	client := github.NewClient(nil)

	if len(token) > 0 {
		client = client.WithAuthToken(token)
	}

	return client
}
