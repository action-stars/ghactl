package core

import (
	"fmt"
	"os"

	"github.com/action-stars/ghactl/internal/fileio"
)

const (
	// linuxTempLookup is the environment variable used to find the temp directory on Linux.
	linuxTempLookup = "TEMPDIR"

	// windowsTempLookup is the environment variable used to find the temp directory on Windows.
	windowsTempLookup = "TMP"

	// runnerTempLookup is the environment variable used to find the runner temp directory.
	runnerTempLookup = "RUNNER_TEMP"
)

// GetTempDir returns the path to the temporary directory.
// It checks the environment variables in the following order:
// 1. RUNNER_TEMP (for GitHub Actions runner)
// 2. TEMPDIR (for Linux)
// 3. TMP (for Windows)
// If none of these environment variables are set or point to a valid directory, an error is returned.
func GetTempDir() (string, error) {
	for _, lookup := range []string{runnerTempLookup, linuxTempLookup, windowsTempLookup} {
		p := os.Getenv(lookup)
		if p == "" {
			continue
		}

		exists, err := fileio.DirExists(p)
		if err != nil {
			return "", err
		}

		if exists {
			return p, nil
		}
	}

	return "", fmt.Errorf("no temp directory found")
}

// CreateTempDir creates a temporary directory and returns its path.
func CreateTempDir() (string, error) {
	d, err := GetTempDir()
	if err != nil {
		return "", err
	}

	return os.MkdirTemp(d, "")
}
