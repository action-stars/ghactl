package core

import (
	"fmt"
	"os"

	"github.com/action-stars/ghactl/internal/util"
)

// SUMMARY_FILE_LOOKUP is the environment variable containing the path to the GitHub Actions step summary file.
const SUMMARY_FILE_LOOKUP = "GITHUB_STEP_SUMMARY"

// WriteSummary writes a summary to be persisted for the current GitHub Actions workflow step.
func WriteSummary(value string) error {
	p := os.Getenv(SUMMARY_FILE_LOOKUP)
	if len(p) == 0 {
		return fmt.Errorf("%s is not defined", SUMMARY_FILE_LOOKUP)
	}

	return util.WriteFileString(p, value)
}
