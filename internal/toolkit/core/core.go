package core

import (
	"os"
)

const RUNNER_DEBUG = "RUNNER_DEBUG"

// IsDebug returns true if the current GitHub Actions step is in debug mode.
func IsDebug() bool {
	return os.Getenv(RUNNER_DEBUG) == "1"
}
