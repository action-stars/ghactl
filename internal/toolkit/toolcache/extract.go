package toolcache

import (
	"github.com/action-stars/ghactl/internal/archive"
	"github.com/action-stars/ghactl/internal/toolkit/core"
)

// ExtractTar extracts a tarball into a temporary directory.
// If gz is true, it will decompress the file.
func ExtractTar(tarFile string, gz bool) (string, error) {
	dest, err := core.CreateTempDir()
	if err != nil {
		return "", err
	}

	if err := archive.UnTar(tarFile, dest, gz); err != nil {
		return "", err
	}
	return dest, nil
}

// ExtractZip extracts a zip archive into a temporary directory.
func ExtractZip(zipFile string) (string, error) {
	dest, err := core.CreateTempDir()
	if err != nil {
		return "", err
	}

	if err := archive.UnZip(zipFile, dest); err != nil {
		return "", err
	}
	return dest, nil
}
