package core

import (
	"fmt"
	"os"
)

// OUTPUT_FILE_LOOKUP is the environment variable containing the path to the GitHub Actions output file
const OUTPUT_FILE_LOOKUP = "GITHUB_OUTPUT"

// SetOutput writes an output to be persisted for the current GitHub Actions workflow step.
func SetOutput(key, value string) error {
	p := os.Getenv(OUTPUT_FILE_LOOKUP)
	if len(p) == 0 {
		return fmt.Errorf("%s is not defined", OUTPUT_FILE_LOOKUP)
	}

	return IssueFileCommand(p, key, value)
}
