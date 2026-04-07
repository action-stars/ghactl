package core

import (
	"fmt"
	"os"
)

// pathFileLookup is the environment variable containing the path to the GitHub Actions runner PATH file.
const pathFileLookup = "GITHUB_PATH"

// AddPath writes a path entry to be persisted for the current GitHub Actions workflow run.
// It also prepends the path to the current process PATH.
func AddPath(value string) error {
	p := os.Getenv(pathFileLookup)
	if p == "" {
		return fmt.Errorf("%s is not defined", pathFileLookup)
	}

	if err := IssueFileCommand(p, value, ""); err != nil {
		return err
	}

	return os.Setenv("PATH", value+string(os.PathListSeparator)+os.Getenv("PATH"))
}
