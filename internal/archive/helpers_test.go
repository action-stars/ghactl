package archive

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"os"
	"testing"
)

type testTarEntry struct {
	name     string
	typeflag byte
	mode     int64
	content  string
}

func newTestTar(t *testing.T, entries []testTarEntry) *tar.Reader {
	t.Helper()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	for _, e := range entries {
		hdr := &tar.Header{
			Name:     e.name,
			Typeflag: e.typeflag,
			Mode:     e.mode,
			Size:     int64(len(e.content)),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}

		if e.content != "" {
			if _, err := tw.Write([]byte(e.content)); err != nil {
				t.Fatal(err)
			}
		}
	}

	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}

	return tar.NewReader(&buf)
}

type testZipEntry struct {
	name    string
	content string
	isDir   bool
}

func newTestZip(t *testing.T, entries []testZipEntry) []*zip.File {
	t.Helper()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	for _, e := range entries {
		if e.isDir {
			hdr := &zip.FileHeader{Name: e.name}
			hdr.SetMode(os.ModeDir | 0o755)

			if _, err := zw.CreateHeader(hdr); err != nil {
				t.Fatal(err)
			}

			continue
		}

		hdr := &zip.FileHeader{Name: e.name}
		hdr.SetMode(0o644)

		w, err := zw.CreateHeader(hdr)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := w.Write([]byte(e.content)); err != nil {
			t.Fatal(err)
		}
	}

	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	data := buf.Bytes()

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}

	return zr.File
}
