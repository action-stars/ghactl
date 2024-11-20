package util

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

func TestFileExists(t *testing.T) {
	t.Run("returns false if the path does not exist", func(t *testing.T) {
		is := is.New(t)

		exists, err := FileExists("non-existent-file")

		is.NoErr(err)    // should not error
		is.True(!exists) // should not exists
	})

	t.Run("returns true if the path is a file", func(t *testing.T) {
		is := is.New(t)
		tmp, err := os.CreateTemp(t.TempDir(), "file")
		if err != nil {
			t.Fatal(err)
		}
		tmp.Close()

		exists, err := FileExists(tmp.Name())

		is.NoErr(err)   // should not error
		is.True(exists) // should exist
	})

	t.Run("errors if the path is a directory", func(t *testing.T) {
		is := is.New(t)

		_, err := FileExists(t.TempDir())

		is.True(err != nil) // should error
	})
}

func TestCopyFile(t *testing.T) {
	t.Run("errors if the source path does not exist", func(t *testing.T) {
		is := is.New(t)

		err := CopyFile("non-existent-file", t.TempDir(), false)

		is.True(err != nil) // should error
	})

	t.Run("errors if the destination path exist", func(t *testing.T) {
		is := is.New(t)
		content := "src"
		src := func() string {
			tmp, err := os.CreateTemp(t.TempDir(), "file")
			if err != nil {
				t.Fatal(err)
			}
			defer tmp.Close()
			_, err = tmp.WriteString(content)
			if err != nil {
				t.Fatal(err)
			}
			return tmp.Name()
		}()
		dest := func() string {
			tmp, err := os.CreateTemp(t.TempDir(), "file")
			if err != nil {
				t.Fatal(err)
			}
			defer tmp.Close()
			_, err = tmp.WriteString("dest")
			if err != nil {
				t.Fatal(err)
			}
			return tmp.Name()
		}()

		err := CopyFile(src, dest, false)

		is.True(err != nil) // should error
	})

	t.Run("copies to empty dir", func(t *testing.T) {
		is := is.New(t)
		content := "src"
		src := func() string {
			tmp, err := os.CreateTemp(t.TempDir(), "file")
			if err != nil {
				t.Fatal(err)
			}
			defer tmp.Close()
			_, err = tmp.WriteString(content)
			if err != nil {
				t.Fatal(err)
			}
			return tmp.Name()
		}()
		dest := t.TempDir()

		err := CopyFile(src, dest, false)

		destContent, _ := os.ReadFile(filepath.Join(dest, filepath.Base(src)))

		is.NoErr(err)                          // should not error
		is.Equal(string(destContent), content) // should match
	})

	t.Run("copies to empty path", func(t *testing.T) {
		is := is.New(t)
		content := "src"
		src := func() string {
			tmp, err := os.CreateTemp(t.TempDir(), "file")
			if err != nil {
				t.Fatal(err)
			}
			defer tmp.Close()
			_, err = tmp.WriteString(content)
			if err != nil {
				t.Fatal(err)
			}
			return tmp.Name()
		}()
		dest := filepath.Join(t.TempDir(), "file")

		err := CopyFile(src, dest, false)

		destContent, _ := os.ReadFile(dest)

		is.NoErr(err)                          // should not error
		is.Equal(string(destContent), content) // should match
	})

	t.Run("copies to existing path with override", func(t *testing.T) {
		is := is.New(t)
		content := "src"
		src := func() string {
			tmp, err := os.CreateTemp(t.TempDir(), "file")
			if err != nil {
				t.Fatal(err)
			}
			defer tmp.Close()
			_, err = tmp.WriteString(content)
			if err != nil {
				t.Fatal(err)
			}
			return tmp.Name()
		}()
		dest := func() string {
			tmp, err := os.CreateTemp(t.TempDir(), "file")
			if err != nil {
				t.Fatal(err)
			}
			defer tmp.Close()
			_, err = tmp.WriteString("dest")
			if err != nil {
				t.Fatal(err)
			}
			return tmp.Name()
		}()

		err := CopyFile(src, dest, true)

		destContent, _ := os.ReadFile(dest)

		is.NoErr(err)                          // should not error
		is.Equal(string(destContent), content) // should match
	})
}

func TestWriteFile(t *testing.T) {
	t.Run("errors if path is invalid", func(t *testing.T) {
		is := is.New(t)

		err := WriteFile("/invalid-file-path", []byte("value"))

		is.True(err != nil) // should error
	})

	t.Run("can write a value", func(t *testing.T) {
		is := is.New(t)
		name := filepath.Join(t.TempDir(), "test")
		value := "value\n"

		defer func() { _ = os.Remove(name) }()

		err := WriteFile(name, []byte(value))

		data, _ := os.ReadFile(name)

		is.NoErr(err)                 // should not error
		is.Equal(string(data), value) // should match
	})

	t.Run("can append value", func(t *testing.T) {
		is := is.New(t)
		name := filepath.Join(t.TempDir(), "test")
		value1 := "hello\n"
		value2 := "world\n"

		defer func() { _ = os.Remove(name) }()

		err := WriteFile(name, []byte(value1))
		err2 := WriteFile(name, []byte(value2))

		data, _ := os.ReadFile(name)

		is.NoErr(err)                                               // should not error
		is.NoErr(err2)                                              // should not error
		is.Equal(string(data), fmt.Sprintf("%s%s", value1, value2)) // should match
	})
}

func TestWriteFileString(t *testing.T) {
	t.Run("errors if path is invalid", func(t *testing.T) {
		is := is.New(t)

		err := WriteFileString("/invalid-file-path", "value")

		is.True(err != nil) // should error
	})

	t.Run("can write a value", func(t *testing.T) {
		is := is.New(t)
		name := filepath.Join(t.TempDir(), "test")
		value := "value\n"

		defer func() { _ = os.Remove(name) }()

		err := WriteFileString(name, value)

		data, _ := os.ReadFile(name)

		is.NoErr(err)                 // should not error
		is.Equal(string(data), value) // should match
	})

	t.Run("can append value", func(t *testing.T) {
		is := is.New(t)
		name := filepath.Join(t.TempDir(), "test")
		value1 := "hello\n"
		value2 := "world\n"

		defer func() { _ = os.Remove(name) }()

		err := WriteFileString(name, value1)
		err2 := WriteFileString(name, value2)

		data, _ := os.ReadFile(name)

		is.NoErr(err)                                               // should not error
		is.NoErr(err2)                                              // should not error
		is.Equal(string(data), fmt.Sprintf("%s%s", value1, value2)) // should match
	})
}
