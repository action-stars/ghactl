package toolcache

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/action-stars/ghactl/internal/util"
	"github.com/matryer/is"
)

func TestGetToolCacheDirectory(t *testing.T) {
	t.Run("errors if tool cache dir env variable is not defined", func(t *testing.T) {
		is := is.New(t)
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, "")

		_, err := GetToolCacheDirectory()

		is.True(err != nil) // should error
	})

	t.Run("returns the tool cache dir", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		result, err := GetToolCacheDirectory()

		is.NoErr(err)        // should not error
		is.Equal(result, tc) // should match
	})
}

func TestFindAllToolVersions(t *testing.T) {
	t.Run("errors if tool cache dir env variable is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := ""
		tool := "test-tool"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := FindAllToolVersions(tool, arch)

		is.True(err != nil) // should error
	})

	t.Run("errors if tool is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := ""
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := FindAllToolVersions(tool, arch)

		is.True(err != nil) // should error
	})

	t.Run("errors if arch is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := "test-tool"
		arch := ""
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := FindAllToolVersions(tool, arch)

		is.True(err != nil) // should error
	})

	t.Run("returns an empty list if tool not cached", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := "test"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		result, err := FindAllToolVersions(tool, arch)

		is.NoErr(err)                // should not error
		is.Equal(result, []string{}) // should match
	})

	t.Run("returns empty list if tool is not cached with marker", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := "test-tool2"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		result, err := FindAllToolVersions(tool, arch)

		is.NoErr(err)                // should not error
		is.Equal(result, []string{}) // should match
	})

	t.Run("returns tool versions with completed markers", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := "test-tool"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		result, err := FindAllToolVersions(tool, arch)

		is.NoErr(err)                                                  // should not error
		is.Equal(result, []string{"1.0.0", "1.0.1", "1.2.0", "2.0.0"}) // should match
	})
}

func TestFindTool(t *testing.T) {
	t.Run("errors if tool cache dir env variable is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := ""
		tool := "test-tool"
		arch := "amd64"
		versionSpec := "1.0.0"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := FindTool(tool, arch, versionSpec)

		is.True(err != nil) // should error
	})

	t.Run("errors if tool is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := ""
		arch := "amd64"
		versionSpec := "1.0.0"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := FindTool(tool, arch, versionSpec)

		is.True(err != nil) // should error
	})

	t.Run("errors if arch is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := "test-tool"
		arch := ""
		versionSpec := "1.0.0"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := FindTool(tool, arch, versionSpec)

		is.True(err != nil) // should error
	})

	t.Run("errors if version spec is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := "test-tool"
		arch := "amd64"
		versionSpec := ""
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := FindTool(tool, arch, versionSpec)

		is.True(err != nil) // should error
	})

	t.Run("returns an empty string if tool is not cached", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := "test"
		arch := "amd64"
		versionSpec := "1.0.0"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		tp, err := FindTool(tool, arch, versionSpec)

		is.NoErr(err)    // should not error
		is.Equal(tp, "") // should be empty
	})

	t.Run("returns an empty string if tool is not cached with a marker", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := "test-tool2"
		arch := "amd64"
		versionSpec := "1.0.0"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		tp, err := FindTool(tool, arch, versionSpec)

		is.NoErr(err)    // should not error
		is.Equal(tp, "") // should be empty
	})

	t.Run("returns explicit version match", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := "test-tool"
		arch := "amd64"
		versionSpec := "1.0.0"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		tp, err := FindTool(tool, arch, versionSpec)

		is.NoErr(err)                                            // should not error
		is.Equal(tp, filepath.Join(tc, tool, versionSpec, arch)) // should be equal
	})

	t.Run("returns patch version match", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := "test-tool"
		arch := "amd64"
		versionSpec := "~1.0.0"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		tp, err := FindTool(tool, arch, versionSpec)

		is.NoErr(err)                                        // should not error
		is.Equal(tp, filepath.Join(tc, tool, "1.0.1", arch)) // should be equal
	})

	t.Run("returns minor version match", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := "test-tool"
		arch := "amd64"
		versionSpec := "^1.0.0"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		tp, err := FindTool(tool, arch, versionSpec)

		is.NoErr(err)                                        // should not error
		is.Equal(tp, filepath.Join(tc, tool, "1.2.0", arch)) // should be equal
	})

	t.Run("returns major version match", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := "test-tool"
		arch := "amd64"
		versionSpec := ">=1.0.0"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		tp, err := FindTool(tool, arch, versionSpec)

		is.NoErr(err)                                        // should not error
		is.Equal(tp, filepath.Join(tc, tool, "2.0.0", arch)) // should be equal
	})

	t.Run("returns version match ignoring v prefix", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := "test-tool"
		arch := "amd64"
		versionSpec := "v1.0.0"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		tp, err := FindTool(tool, arch, versionSpec)

		is.NoErr(err)                                        // should not error
		is.Equal(tp, filepath.Join(tc, tool, "1.0.0", arch)) // should be equal
	})

	t.Run("returns version match with wildcard", func(t *testing.T) {
		is := is.New(t)
		tc := "../../../testdata/tool-cache"
		tool := "test-tool"
		arch := "amd64"
		versionSpec := "*"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		tp, err := FindTool(tool, arch, versionSpec)

		is.NoErr(err)                                        // should not error
		is.Equal(tp, filepath.Join(tc, tool, "2.0.0", arch)) // should be equal
	})
}

func TestCacheDir(t *testing.T) {
	t.Run("errors if the source path does not exist", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		tool := "test-tool"
		version := "1.0.0"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := CacheDir("non-existent-path", tool, version, arch)

		is.True(err != nil) // should error
	})

	t.Run("errors if the tool is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		tool := ""
		version := "1.0.0"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := CacheDir(t.TempDir(), tool, version, arch)

		is.True(err != nil) // should error
	})

	t.Run("errors if the version spec is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		tool := "test-tool"
		version := ""
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := CacheDir(t.TempDir(), tool, version, arch)

		is.True(err != nil) // should error
	})

	t.Run("errors if the arch is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		tool := "test-tool"
		version := "1.0.0"
		arch := ""
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := CacheDir(t.TempDir(), tool, version, arch)

		is.True(err != nil) // should error
	})

	t.Run("errors if the tool cache path is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := ""
		tool := "test-tool"
		version := "1.0.0"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := CacheDir(t.TempDir(), tool, version, arch)

		is.True(err != nil) // should error
	})

	t.Run("caches tool dir", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		tool := "test-tool"
		version := "1.0.0"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		p, err := CacheDir("../../../testdata/test-tool", tool, version, arch)

		toolExists, _ := util.FileExists(filepath.Join(p, "test-tool"))
		markerExists, _ := util.FileExists(getMarkerPath(p))

		is.NoErr(err)                                                    // should not error
		is.Equal(p, filepath.Join(tc, tool, version, getNodeArch(arch))) // should be equal
		is.True(toolExists)                                              // tool should exist
		is.True(markerExists)                                            // marker file should exist
	})
}

func TestCacheFile(t *testing.T) {
	t.Run("errors if the source path does not exist", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		targetName := "test-tool"
		tool := "test-tool"
		version := "1.0.0"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := CacheFile("non-existent-path", targetName, tool, version, arch)

		is.True(err != nil) // should error
	})

	t.Run("errors if the target name is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		targetName := ""
		tool := "test-tool"
		version := "1.0.0"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := CacheFile(t.TempDir(), targetName, tool, version, arch)

		is.True(err != nil) // should error
	})

	t.Run("errors if the tool is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		targetName := "test-tool"
		tool := ""
		version := "1.0.0"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := CacheFile(t.TempDir(), targetName, tool, version, arch)

		is.True(err != nil) // should error
	})

	t.Run("errors if the version spec is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		targetName := "test-tool"
		tool := "test-tool"
		version := ""
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := CacheFile(t.TempDir(), targetName, tool, version, arch)

		is.True(err != nil) // should error
	})

	t.Run("errors if the arch is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		targetName := "test-tool"
		tool := "test-tool"
		version := "1.0.0"
		arch := ""
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := CacheFile(t.TempDir(), targetName, tool, version, arch)

		is.True(err != nil) // should error
	})

	t.Run("errors if the tool cache path is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := ""
		targetName := "test-tool"
		tool := "test-tool"
		version := "1.0.0"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := CacheFile(t.TempDir(), targetName, tool, version, arch)

		is.True(err != nil) // should error
	})

	t.Run("caches tool file", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		targetName := "test-tool"
		tool := "test-tool"
		version := "1.0.0"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		p, err := CacheFile("../../../testdata/test-tool/test-tool", targetName, tool, version, arch)

		toolExists, _ := util.FileExists(filepath.Join(p, "test-tool"))
		markerExists, _ := util.FileExists(getMarkerPath(p))

		is.NoErr(err)                                                    // should not error
		is.Equal(p, filepath.Join(tc, tool, version, getNodeArch(arch))) // should be equal
		is.True(toolExists)                                              // tool should exist
		is.True(markerExists)                                            // marker file should exist
	})
}

func Test_createToolPath(t *testing.T) {
	t.Run("errors if the tool cache path is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := ""
		tool := "tool"
		version := "1.0.0"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := createToolPath(tool, version, arch)

		is.True(err != nil) // should error
	})

	t.Run("errors if the tool is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		tool := ""
		version := "1.0.0"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := createToolPath(tool, version, arch)

		is.True(err != nil) // should error
	})

	t.Run("errors if the version spec is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		tool := "tool"
		version := ""
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := createToolPath(tool, version, arch)

		is.True(err != nil) // should error
	})

	t.Run("errors if the arch is not defined", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		tool := "tool"
		version := "1.0.0"
		arch := ""
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		_, err := createToolPath(tool, version, arch)

		is.True(err != nil) // should error
	})

	t.Run("creates a tool path", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		tool := "tool"
		version := "1.0.0"
		arch := "amd64"
		t.Setenv(RUNNER_TOOL_CACHE_LOOKUP, tc)

		vtp, err := createToolPath(tool, version, arch)

		dirExists, _ := util.DirExists(vtp)
		markerExists, _ := util.FileExists(getMarkerPath(vtp))

		is.NoErr(err)                                                      // should not error
		is.Equal(vtp, filepath.Join(tc, tool, version, getNodeArch(arch))) // should be equal
		is.True(dirExists)                                                 // dir should exist
		is.True(!markerExists)                                             // marker file should not exist
	})
}

func Test_getMarkerPath(t *testing.T) {
	t.Run("returns the marker path", func(t *testing.T) {
		is := is.New(t)
		tp := t.TempDir()

		mp := getMarkerPath(tp)

		is.Equal(mp, fmt.Sprintf("%s.complete", tp)) // should be equal
	})
}

func Test_getNodeArch(t *testing.T) {
	t.Run("returns correct amd64 value", func(t *testing.T) {
		is := is.New(t)

		na := getNodeArch("amd64")

		is.Equal(na, "x64") // should be equal
	})

	t.Run("returns correct 386 value", func(t *testing.T) {
		is := is.New(t)

		na := getNodeArch("386")

		is.Equal(na, "ia32") // should be equal
	})

	t.Run("returns correct arm64 value", func(t *testing.T) {
		is := is.New(t)

		na := getNodeArch("arm64")

		is.Equal(na, "arm64") // should be equal
	})
}
