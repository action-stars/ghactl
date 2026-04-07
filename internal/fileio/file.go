package fileio

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// FileExists checks if a file exists.
// It returns true if the path exists and is a file.
// It returns false if the path does not exist.
// It returns an error if the path exists and is not a file.
func FileExists(p string) (bool, error) {
	fi, err := os.Stat(p)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	if fi.IsDir() {
		return false, fmt.Errorf("%s is not a file", p)
	}
	return true, nil
}

// CopyFile copies a file from src to dest.
// If the destination file exists and overwrite is false, an error is returned.
func CopyFile(src, dest string, overwrite bool) error {
	sfi, err := os.Stat(src)
	if err != nil {
		return err
	}

	// If dest is a directory, use the source filename to create the full destination path.
	if dfi, err := os.Stat(dest); err == nil && dfi.IsDir() {
		dest = filepath.Join(dest, filepath.Base(src))
	}

	// Check if we can't overwrite the destination file.
	if !overwrite {
		if _, err := os.Stat(dest); err == nil {
			return fmt.Errorf("destination %s already exists", dest)
		}
	}

	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, sfi.Mode())
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	return err
}

// WriteFile writes bytes to a file.
// If the file does not exist, it will be created.
// If the file exists, it will be appended to.
func WriteFile(name string, value []byte) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(value)
	return err
}
