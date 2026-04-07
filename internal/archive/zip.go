package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// UnZip extracts a zip file to a destination directory.
func UnZip(zipFile, dest string) error {
	root, err := os.OpenRoot(dest)
	if err != nil {
		return fmt.Errorf("opening destination: %w", err)
	}
	defer root.Close()

	zr, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("opening zip file: %w", err)
	}
	defer zr.Close()

	return extractZip(root, zr.File)
}

// extractZip extracts zip files into root.
func extractZip(root *os.Root, files []*zip.File) error {
	for _, zf := range files {
		if err := extractZipFile(root, zf); err != nil {
			return err
		}
	}
	return nil
}

// extractZipFile extracts a single zip file entry into root.
func extractZipFile(root *os.Root, zf *zip.File) error {
	name := strings.TrimLeft(zf.Name, "/")
	if name == "" {
		return nil
	}

	if zf.FileInfo().IsDir() {
		if err := root.MkdirAll(name, 0o755); err != nil {
			return fmt.Errorf("creating directory %s: %w", name, err)
		}
		return nil
	}

	if dir := path.Dir(name); dir != "." {
		if err := root.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("creating parent directory for %s: %w", name, err)
		}
	}

	f, err := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zf.Mode()&0o777)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", name, err)
	}

	zaf, err := zf.Open()
	if err != nil {
		f.Close()
		return fmt.Errorf("opening zip entry %s: %w", name, err)
	}

	_, copyErr := io.Copy(f, zaf)
	closeErr := f.Close()
	zaf.Close()

	if copyErr != nil {
		return fmt.Errorf("writing file %s: %w", name, copyErr)
	}
	if closeErr != nil {
		return fmt.Errorf("closing file %s: %w", name, closeErr)
	}

	return nil
}
