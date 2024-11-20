package util

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// DirExists checks if a directory exists.
// It returns true if the path exists and is a directory.
// It returns false if the path does not exist.
// It returns an error if the path exists and is not a directory.
func DirExists(p string) (bool, error) {
	fi, err := os.Stat(p)
	if err == nil && fi.IsDir() {
		return true, nil
	} else if err == nil {
		return false, fmt.Errorf("%s is not a directory", p)
	} else if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}
