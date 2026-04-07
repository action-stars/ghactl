package core

import (
	"fmt"
	"os"
)

// outputFileLookup is the environment variable containing the path to the GitHub Actions output file.
const outputFileLookup = "GITHUB_OUTPUT"

// SetOutput writes an output to be persisted for the current GitHub Actions workflow step.
func SetOutput(key, value string) error {
	p := os.Getenv(outputFileLookup)
	if p == "" {
		return fmt.Errorf("%s is not defined", outputFileLookup)
	}

	return IssueFileCommand(p, key, value)
}
