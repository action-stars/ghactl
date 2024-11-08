package util

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// FileExists checks if a directory exists.
// It returns true if the path exists and is a file.
// It returns false if the path does not exist.
// It returns an error if the path exists and is not a file.
func FileExists(p string) (bool, error) {
	fi, err := os.Stat(p)
	if err == nil && !fi.IsDir() {
		return true, nil
	} else if err == nil {
		return false, fmt.Errorf("%s is not a file", p)
	} else if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

// CopyFile copies a file from src to dest.
// If the destination file exists and overwrite is false, an error is returned.
func CopyFile(src, dest string, overwrite bool) error {
	sfi, err := os.Stat(src)
	if err != nil {
		return err
	}

	dfi, err := os.Stat(dest)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if err == nil && dfi.IsDir() {
		dest = filepath.Join(dest, filepath.Base(src))
	}

	_, err = os.Stat(dest)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if err == nil {
		if !overwrite {
			return fmt.Errorf("destination %s already exists", dest)
		}

		err := os.Remove(dest)
		if err != nil {
			return err
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
	if err != nil {
		return err
	}

	return nil
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
	if err != nil {
		return err
	}

	return nil
}

// WriteFileString writes a string to a file.
// If the file does not exist, it will be created.
// If the file exists, it will be appended to.
func WriteFileString(name, value string) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(value)
	if err != nil {
		return err
	}

	return nil
}
