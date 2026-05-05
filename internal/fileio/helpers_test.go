package fileio

import (
	"os"
	"testing"
)

func createTempFile(t *testing.T, content string) string {
	t.Helper()

	tmp, err := os.CreateTemp(t.TempDir(), "file")
	if err != nil {
		t.Fatal(err)
	}
	defer tmp.Close()

	if _, err := tmp.WriteString(content); err != nil {
		t.Fatal(err)
	}

	return tmp.Name()
}
