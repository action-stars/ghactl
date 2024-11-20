package util

import (
	"errors"
	"io/fs"
	"os"
)

// PathExists checks if a path exists.
func PathExists(p string) (bool, error) {
	_, err := os.Stat(p)
	if err == nil {
		return true, nil
	} else if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}
