package core

import (
	"fmt"
	"os"

	"github.com/action-stars/ghactl/internal/util"
)

// PATH_FILE_LOOKUP is the environment variable containing the path to the GitHub Actions runner PATH file.
const PATH_FILE_LOOKUP = "GITHUB_PATH"

// AddPath writes a path entry to be persisted for the current GitHub Actions workflow run.
func AddPath(value string) error {
	p := os.Getenv(PATH_FILE_LOOKUP)
	if len(p) == 0 {
		return fmt.Errorf("%s is not defined", PATH_FILE_LOOKUP)
	}

	return util.WriteFileString(p, value)
}
