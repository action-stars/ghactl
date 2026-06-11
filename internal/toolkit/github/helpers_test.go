package github

// func mustNewGitHubClient(t *testing.T, c *http.Client, u string, token *string) *github.Client {
// 	t.Helper()

// 	options := []github.ClientOptionsFunc{
// 		github.WithHTTPClient(c),
// 		github.WithURLs(&u, nil),
// 	}

// 	if token != nil && *token != "" {
// 		options = append(options, github.WithAuthToken(*token))
// 	}

// 	client, err := github.NewClient(options...)
// 	if err != nil {
// 		t.Fatalf("failed to create GitHub client: %v", err)
// 	}

// 	return client
// }
