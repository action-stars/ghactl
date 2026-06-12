package core

import (
	"fmt"
	"os"
)

// pathFileLookup is the environment variable containing the path to the GitHub Actions runner PATH file.
const pathFileLookup = "GITHUB_PATH"

// AddPath writes a path entry to be persisted for the current GitHub Actions workflow run.
func AddPath(p string) error {
	fp, ok := os.LookupEnv(pathFileLookup)
	if !ok || fp == "" {
		return fmt.Errorf("%s is not defined", pathFileLookup)
	}

	return writeFile(fp, []byte(p+"\n"))
}
