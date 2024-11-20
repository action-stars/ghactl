package core

import (
	"fmt"
	"os"

	"github.com/action-stars/ghactl/internal/util"
)

// RUNNER_TEMP_LOOKUP is the environment variable used to find the runner temp directory.
const RUNNER_TEMP_LOOKUP = "RUNNER_TEMP"

// GetTempDirectory returns the path to the runner temp directory.
func GetTempDirectory() (string, error) {
	p := os.Getenv(RUNNER_TEMP_LOOKUP)
	if len(p) == 0 {
		return "", fmt.Errorf("%s is not defined", RUNNER_TEMP_LOOKUP)
	}

	exists, err := util.DirExists(p)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", fmt.Errorf("dir %s does not exist", p)
	}

	return p, nil
}
