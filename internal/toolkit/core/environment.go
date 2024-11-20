package core

import (
	"fmt"
	"os"
)

// ENV_FILE_LOOKUP is the environment variable containing the path to the GitHub Actions environment variable file
const ENV_FILE_LOOKUP = "GITHUB_ENV"

// ExportVariable writes an environment variable to be persisted for the GitHub Actions workflow.
func ExportVariable(key, value string) error {
	p := os.Getenv(ENV_FILE_LOOKUP)
	if len(p) == 0 {
		return fmt.Errorf("%s is not defined", ENV_FILE_LOOKUP)
	}

	return IssueFileCommand(p, key, value)
}
