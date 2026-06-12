package toolcache

import (
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
//   - Single root directory (steps into it)
//   - 'bin' subdirectory (steps into it if it exists after entering a single nested directory)
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
	finalPath := extractedPath

	// Check if directory exists
	exists, err := fileio.DirExists(finalPath)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", nil
	}

	// Check for single root directory
	entries, err := os.ReadDir(finalPath)
	if err != nil {
		return "", err
	}

	steppedIntoNested := false
	if len(entries) == 1 && entries[0].IsDir() {
		finalPath = filepath.Join(finalPath, entries[0].Name())
		steppedIntoNested = true
	}

	// Check for bin subdirectory only after stepping into a nested directory
	if steppedIntoNested {
		binPath := filepath.Join(finalPath, "bin")
		if exists, err := fileio.DirExists(binPath); err == nil && exists {
			finalPath = binPath
		}
	}

	return finalPath, nil
}
