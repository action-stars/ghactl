package core

import (
	"fmt"
	"os"

	"github.com/action-stars/ghactl/internal/fileio"
)

// summaryFileLookup is the environment variable containing the path to the GitHub Actions step summary file.
const summaryFileLookup = "GITHUB_STEP_SUMMARY"

// WriteSummary writes a summary to be persisted for the current GitHub Actions workflow step.
func WriteSummary(value string) error {
	p := os.Getenv(summaryFileLookup)
	if p == "" {
		return fmt.Errorf("%s is not defined", summaryFileLookup)
	}

	return fileio.WriteFile(p, []byte(value))
}
