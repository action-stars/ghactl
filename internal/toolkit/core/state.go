package core

import (
	"fmt"
	"os"
)

// stateFileLookup is the environment variable containing the path to the GitHub Actions state file.
const stateFileLookup = "GITHUB_STATE"

// SaveState saves state for the current action.
// The state can only be retrieved by this action's post job execution.
func SaveState(key, value string) error {
	p := os.Getenv(stateFileLookup)
	if p == "" {
		return fmt.Errorf("%s is not defined", stateFileLookup)
	}

	return IssueFileCommand(p, key, value)
}

// GetState gets the value of a state set by this action's main execution.
func GetState(name string) string {
	return os.Getenv("STATE_" + name)
}
