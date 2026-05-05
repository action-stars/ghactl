package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// UnTar extracts a tar file to a destination directory.
// If gz is true, the file is decompressed with gzip.
func UnTar(tarFile, dest string, gz bool) error {
	root, err := os.OpenRoot(dest)
	if err != nil {
		return fmt.Errorf("opening destination: %w", err)
	}
	defer root.Close()

	fr, err := os.Open(tarFile)
	if err != nil {
		return fmt.Errorf("opening tar file: %w", err)
	}
	defer fr.Close()

	if !gz {
		return extractTar(root, tar.NewReader(fr))
	}

	gzr, err := gzip.NewReader(fr)
	if err != nil {
		return fmt.Errorf("creating gzip reader: %w", err)
	}
	defer gzr.Close()

	return extractTar(root, tar.NewReader(gzr))
}

// extractTar reads entries from a tar reader and extracts them into root.
func extractTar(root *os.Root, tr *tar.Reader) error {
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("reading tar entry: %w", err)
		}
		if header == nil {
			continue
		}

		if err := extractTarEntry(root, header, tr); err != nil {
			return err
		}
	}
}

// extractTarEntry extracts a single tar entry into root.
func extractTarEntry(root *os.Root, header *tar.Header, r io.Reader) error {
	name := strings.TrimLeft(header.Name, "/")
	if name == "" {
		return nil
	}

	switch header.Typeflag {
	case tar.TypeDir:
		if err := root.MkdirAll(name, 0o755); err != nil {
			return fmt.Errorf("creating directory %s: %w", name, err)
		}

	case tar.TypeReg:
		if dir := path.Dir(name); dir != "." {
			if err := root.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("creating parent directory for %s: %w", name, err)
			}
		}

		f, err := root.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode)&0o777)
		if err != nil {
			return fmt.Errorf("creating file %s: %w", name, err)
		}

		_, copyErr := io.Copy(f, r)
		closeErr := f.Close()
		if copyErr != nil {
			return fmt.Errorf("writing file %s: %w", name, copyErr)
		}
		if closeErr != nil {
			return fmt.Errorf("closing file %s: %w", name, closeErr)
		}
	}

	return nil
}
