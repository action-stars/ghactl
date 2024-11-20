package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// UnTar extracts a tar file to a destination.
// If dest is not a directory or doesn't exist, it will return an error.
// If gz is true, it will decompress the file.
func UnTar(tarFile, dest string, gz bool) error {
	destIsDir, err := DirExists(dest)
	if err != nil {
		return err
	}

	if !destIsDir {
		return fmt.Errorf("%s is not a dir", dest)
	}

	fr, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer fr.Close()

	var r io.Reader
	if gz {
		gzr, err := gzip.NewReader(fr)
		if err != nil {
			return err
		}
		defer gzr.Close()
		r = gzr
	} else {
		r = fr
	}

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		switch {
		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			exists, err := DirExists(target)
			if err != nil {
				return err
			}

			if !exists {
				err := os.MkdirAll(target, 0o755)
				if err != nil {
					return err
				}
			}

		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			_, err = io.Copy(f, tr)
			if err != nil {
				return err
			}

			f.Close()
		}
	}
}
