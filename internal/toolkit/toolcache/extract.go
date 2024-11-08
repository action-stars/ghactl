package toolcache

import (
	"os"
	"path/filepath"

	"github.com/action-stars/ghactl/internal/toolkit/core"
	"github.com/action-stars/ghactl/internal/util"
)

// ExtractTar extracts a tarball into a temporary directory.
// If gz is true, it will decompress the file.
func ExtractTar(tarFile string, gz bool) (string, error) {
	d, err := core.GetTempDirectory()
	if err != nil {
		return "", err
	}

	dest := filepath.Join(d, util.GenerateRandomString(16))
	err = os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return "", err
	}

	err = util.UnTar(tarFile, dest, gz)
	if err != nil {
		return "", err
	}
	return dest, nil
}

// ExtractZip extracts a zip archive into a temporary directory.
func ExtractZip(zipFile string) (string, error) {
	d, err := core.GetTempDirectory()
	if err != nil {
		return "", err
	}

	dest := filepath.Join(d, util.GenerateRandomString(16))
	err = os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return "", err
	}

	err = util.UnZip(zipFile, dest)
	if err != nil {
		return "", err
	}
	return dest, nil
}
