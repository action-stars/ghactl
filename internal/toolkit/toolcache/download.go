package toolcache

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/action-stars/ghactl/internal/toolkit/core"
)

// DownloadTool downloads a tool from a URL and saves it to a temporary directory.
func DownloadTool(ctx context.Context, logger any, url url.URL) (string, error) {
	d, err := core.GetTempDir()
	if err != nil {
		return "", err
	}

	client := retryablehttp.NewClient()
	client.Logger = logger

	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	dest, err := os.CreateTemp(d, "tool-*")
	if err != nil {
		return "", err
	}
	defer dest.Close()

	if _, err := dest.ReadFrom(resp.Body); err != nil {
		return "", err
	}

	return dest.Name(), nil
}
