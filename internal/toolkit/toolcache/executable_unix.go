//go:build !windows

package toolcache

import (
	"os"
)

// ensureExecutable sets the user executable bit on Unix-like systems.
func ensureExecutable(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	mode := fi.Mode()
	if mode&0o100 != 0 {
		return nil
	}

	return os.Chmod(path, mode|0o100)
}
