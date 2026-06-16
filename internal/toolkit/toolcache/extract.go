package toolcache

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/action-stars/ghactl/internal/archive"
	"github.com/action-stars/ghactl/internal/fileio"
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

// ResolveToolDirectory navigates nested directories to find the actual tool location.
// It handles:
//   - Single nested directories
//   - 'bin' subdirectory after evaluating single nested directories
//
// For example, if an archive extracts to:
//
//	extracted/
//	├── jsonschema-v1.0.0/
//	│   ├── bin/
//	│   │   └── jsonschema
//	│   └── lib/
//
// This will resolve to: extracted/jsonschema-v1.0.0/bin.
func ResolveToolDirectory(extractedPath string) (string, error) {
	exists, err := fileio.DirExists(extractedPath)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("extracted path does not exist: %s", extractedPath)
	}

	path := extractedPath
	for {
		entries, err := os.ReadDir(path)
		if err != nil {
			return "", err
		}

		if len(entries) != 1 || !entries[0].IsDir() {
			break
		}

		path = filepath.Join(path, entries[0].Name())
	}

	binPath := filepath.Join(path, "bin")
	if exists, err := fileio.DirExists(binPath); err == nil && exists {
		path = binPath
	}

	return path, nil
}
