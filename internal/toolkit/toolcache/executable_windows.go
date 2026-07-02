//go:build windows

package toolcache

// ensureExecutable is a no-op on Windows.
func ensureExecutable(path string) error {
	return nil
}
