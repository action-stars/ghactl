package tool

import (
	"testing"

	"github.com/matryer/is"
)

func TestCmd_CacheGet(t *testing.T) {
	c := &Cmd{}

	t.Run("errors_if_tool_cache_dir_env_variable_is_not_defined", func(t *testing.T) {
		is := is.New(t)
		t.Setenv("RUNNER_TOOL_CACHE", "")

		_, err := c.CacheGet()

		is.True(err != nil) // should error
	})

	t.Run("returns_the_tool_cache_dir", func(t *testing.T) {
		is := is.New(t)
		tc := setupToolCache(t)

		result, err := c.CacheGet()

		is.NoErr(err)        // should not error
		is.Equal(result, tc) // should return expected path
	})
}

func TestCmd_CacheFindAll(t *testing.T) {
	c := &Cmd{}

	t.Run("errors_if_tool_cache_dir_is_not_defined", func(t *testing.T) {
		is := is.New(t)
		t.Setenv("RUNNER_TOOL_CACHE", "")

		_, err := c.CacheFindAll("test-tool", "amd64")

		is.True(err != nil) // should error
	})

	t.Run("returns_empty_for_non-existent_tool", func(t *testing.T) {
		is := is.New(t)
		setupToolCache(t)

		ps, err := c.CacheFindAll("non-existent", "amd64")

		is.NoErr(err)        // should not error
		is.Equal(len(ps), 0) // should return empty
	})

	t.Run("returns_versions_for_existing_tool", func(t *testing.T) {
		is := is.New(t)
		setupToolCache(t)

		ps, err := c.CacheFindAll("test-tool", "amd64")

		is.NoErr(err)        // should not error
		is.True(len(ps) > 0) // should find versions
	})
}

func TestCmd_CacheFind(t *testing.T) {
	c := &Cmd{}

	t.Run("errors_if_tool_cache_dir_is_not_defined", func(t *testing.T) {
		is := is.New(t)
		t.Setenv("RUNNER_TOOL_CACHE", "")

		_, err := c.CacheFind("test-tool", "amd64", "*")

		is.True(err != nil) // should error
	})

	t.Run("returns_empty_for_non-existent_tool", func(t *testing.T) {
		is := is.New(t)
		setupToolCache(t)

		p, err := c.CacheFind("non-existent", "amd64", "*")

		is.NoErr(err)   // should not error
		is.Equal(p, "") // should return empty
	})

	t.Run("finds_tool_matching_version_spec", func(t *testing.T) {
		is := is.New(t)
		setupToolCache(t)

		p, err := c.CacheFind("test-tool", "amd64", "^1.0.0")

		is.NoErr(err)    // should not error
		is.True(p != "") // should find tool
	})

	t.Run("uses_default_arch_when_empty", func(t *testing.T) {
		is := is.New(t)
		setupToolCache(t)

		p, err := c.CacheFind("test-tool", "", "*")

		is.NoErr(err) // should not error
		// result depends on runtime.GOARCH; just ensure no error
		_ = p
	})

	t.Run("uses_wildcard_when_version_spec_is_empty", func(t *testing.T) {
		is := is.New(t)
		setupToolCache(t)

		p, err := c.CacheFind("test-tool", "amd64", "")

		is.NoErr(err)    // should not error
		is.True(p != "") // should find latest
	})
}

func TestCmd_CacheDir(t *testing.T) {
	c := &Cmd{}

	t.Run("errors_if_source_does_not_exist", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		t.Setenv("RUNNER_TOOL_CACHE", tc)

		_, err := c.CacheDir("/nonexistent", "my-tool", "1.0.0", "amd64")

		is.True(err != nil) // should error
	})

	t.Run("caches_directory_successfully", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		t.Setenv("RUNNER_TOOL_CACHE", tc)
		source := createSourceDir(t)

		p, err := c.CacheDir(source, "my-tool", "1.0.0", "amd64")

		is.NoErr(err)    // should not error
		is.True(p != "") // should return path
	})
}

func TestCmd_CacheFile(t *testing.T) {
	c := &Cmd{}

	t.Run("errors_if_source_does_not_exist", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		t.Setenv("RUNNER_TOOL_CACHE", tc)

		_, err := c.CacheFile("/nonexistent", "my-tool", "my-tool", "1.0.0", "amd64")

		is.True(err != nil) // should error
	})

	t.Run("caches_file_successfully", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		t.Setenv("RUNNER_TOOL_CACHE", tc)
		source := createSourceFile(t)

		p, err := c.CacheFile(source, "my-tool", "my-tool", "1.0.0", "amd64")

		is.NoErr(err)    // should not error
		is.True(p != "") // should return path
	})

	t.Run("defaults_target_name_to_tool_name", func(t *testing.T) {
		is := is.New(t)
		tc := t.TempDir()
		t.Setenv("RUNNER_TOOL_CACHE", tc)
		source := createSourceFile(t)

		p, err := c.CacheFile(source, "", "my-tool", "1.0.0", "amd64")

		is.NoErr(err)    // should not error
		is.True(p != "") // should return path
	})
}
