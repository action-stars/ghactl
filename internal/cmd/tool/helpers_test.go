package tool

import (
	"os"
	"path/filepath"
	"testing"
)

func setupToolCache(t *testing.T) string {
	t.Helper()

	abs, err := filepath.Abs("../../../testdata/tool-cache")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("RUNNER_TOOL_CACHE", abs)
	return abs
}

func setupTempDir(t *testing.T) {
	t.Helper()

	d := t.TempDir()
	t.Setenv("RUNNER_TEMP", d)
}

func createSourceDir(t *testing.T) string {
	t.Helper()

	d := t.TempDir()
	err := os.WriteFile(filepath.Join(d, "tool-binary"), []byte("binary"), 0o755)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func createSourceFile(t *testing.T) string {
	t.Helper()

	d := t.TempDir()
	p := filepath.Join(d, "tool-binary")
	err := os.WriteFile(p, []byte("binary"), 0o755)
	if err != nil {
		t.Fatal(err)
	}
	return p
}
