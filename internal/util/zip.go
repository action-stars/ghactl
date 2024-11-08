package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// UnZip unzips a zip file to a destination.
// If dest is not a directory or doesn't exist, it will return an error.
func UnZip(zipFile, dest string) error {
	destIsDir, err := DirExists(dest)
	if err != nil {
		return err
	}

	if !destIsDir {
		return fmt.Errorf("%s is not a dir", dest)
	}

	cleanDest := fmt.Sprintf("%s%c", filepath.Clean(dest), os.PathSeparator)

	zr, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zr.Close()

	for _, zf := range zr.File {
		filePath := filepath.Join(dest, zf.Name)

		if !strings.HasPrefix(filePath, cleanDest) {
			return fmt.Errorf("invalid file path %s", filePath)
		}

		if zf.FileInfo().IsDir() {
			err := os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil {
			return err
		}

		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zf.Mode())
		if err != nil {
			return err
		}

		zaf, err := zf.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(f, zaf)
		if err != nil {
			return err
		}

		f.Close()
		zaf.Close()
	}

	return nil
}
