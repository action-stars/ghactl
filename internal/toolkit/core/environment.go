package core

import (
	"fmt"
	"os"
)

// envFileLookup is the environment variable containing the path to the GitHub Actions environment variable file.
const envFileLookup = "GITHUB_ENV"

// ExportVariable writes an environment variable to be persisted for the GitHub Actions workflow.
// It also sets the variable in the current process environment.
func ExportVariable(key, value string) error {
	p := os.Getenv(envFileLookup)
	if p == "" {
		return fmt.Errorf("%s is not defined", envFileLookup)
	}

	if err := os.Setenv(key, value); err != nil {
		return err
	}

	return IssueFileCommand(p, key, value)
}
