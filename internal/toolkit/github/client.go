package github

import "github.com/google/go-github/v88/github"

// GetClient returns a GitHub client.
// If a token is provided, it will be used for authentication.
func GetClient(token *string) (*github.Client, error) {
	options := []github.ClientOptionsFunc{}

	if token != nil && *token != "" {
		options = append(options, github.WithAuthToken(*token))
	}

	return github.NewClient(options...)
}
