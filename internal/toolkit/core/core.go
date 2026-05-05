package core

import (
	"os"
)

const runnerDebug = "RUNNER_DEBUG"

// IsDebug returns true if the current GitHub Actions step is in debug mode.
func IsDebug() bool {
	return os.Getenv(runnerDebug) == "1"
}
