package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadFile downloads a file from the internet and saves it to the specified path.
// If the destination is a directory, an error is returned.
// If the destination file exists and overwrite is false, an error is returned.
func DownloadFile(url, dest string, overwrite bool) error {
	exists, err := PathExists(dest)
	if err != nil {
		return err
	}

	if exists {
		// Can ignore error here, as it's already been checked'
		fi, _ := os.Stat(dest)

		if fi.IsDir() {
			return fmt.Errorf("%s is a directory", dest)
		}

		if !overwrite {
			return fmt.Errorf("destination %s already exists", dest)
		}

		err := os.Remove(dest)
		if err != nil {
			return err
		}
	} else {
		dir := filepath.Dir(dest)
		dirExists, err := DirExists(dir)
		if err != nil {
			return err
		}

		if !dirExists {
			return fmt.Errorf("destination directory %s does not exist", dir)
		}
	}

	df, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer df.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(df, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
