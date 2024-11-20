package toolcache

import (
	"path/filepath"

	"github.com/action-stars/ghactl/internal/toolkit/core"
	"github.com/action-stars/ghactl/internal/util"
)

// DownloadTool downloads a tool from a URL and saves it to a temporary directory.
func DownloadTool(url string) (string, error) {
	d, err := core.GetTempDirectory()
	if err != nil {
		return "", err
	}

	dest := filepath.Join(d, util.GenerateRandomString(16))

	err = util.DownloadFile(url, dest, false)
	if err != nil {
		return "", err
	}

	return dest, nil
}
