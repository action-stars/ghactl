package toolcache

import (
	"path/filepath"
	"testing"

	"github.com/matryer/is"

	"github.com/action-stars/ghactl/internal/fileio"
)

func TestGetToolCacheDirectory(t *testing.T) {
	validDir := t.TempDir()

	tests := []struct {
		name    string
		tc      string
		want    string
		wantErr bool
	}{
		{
			name:    "errors_if_tool_cache_dir_env_variable_is_not_defined",
			tc:      "",
			wantErr: true,
		},
		{
			name: "returns_the_tool_cache_dir",
			tc:   validDir,
			want: validDir,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			t.Setenv(runnerToolCacheLookup, tt.tc)

			result, err := GetToolCacheDirectory()

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)             // should not error
			is.Equal(result, tt.want) // should match
		})
	}
}

func TestFindAllToolVersions(t *testing.T) {
	tests := []struct {
		name    string
		tc      string
		tool    string
		arch    string
		want    []string
		wantErr bool
	}{
		{
			name:    "errors_if_tool_cache_dir_env_variable_is_not_defined",
			tc:      "",
			tool:    "test-tool",
			arch:    "amd64",
			wantErr: true,
		},
		{
			name:    "errors_if_tool_is_not_defined",
			tc:      "../../../testdata/tool-cache",
			tool:    "",
			arch:    "amd64",
			wantErr: true,
		},
		{
			name:    "errors_if_arch_is_not_defined",
			tc:      "../../../testdata/tool-cache",
			tool:    "test-tool",
			arch:    "",
			wantErr: true,
		},
		{
			name: "returns_an_empty_list_if_tool_not_cached",
			tc:   "../../../testdata/tool-cache",
			tool: "test",
			arch: "amd64",
			want: []string{},
		},
		{
			name: "returns_empty_list_if_tool_is_not_cached_with_marker",
			tc:   "../../../testdata/tool-cache",
			tool: "test-tool2",
			arch: "amd64",
			want: []string{},
		},
		{
			name: "returns_tool_versions_with_completed_markers",
			tc:   "../../../testdata/tool-cache",
			tool: "test-tool",
			arch: "amd64",
			want: []string{"1.0.0", "1.0.1", "1.2.0", "2.0.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			t.Setenv(runnerToolCacheLookup, tt.tc)

			result, err := FindAllToolVersions(tt.tool, tt.arch)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)             // should not error
			is.Equal(result, tt.want) // should match
		})
	}
}

func TestFindTool(t *testing.T) {
	tc := "../../../testdata/tool-cache"

	tests := []struct {
		name        string
		tc          string
		tool        string
		arch        string
		versionSpec string
		want        string
		wantErr     bool
	}{
		{
			name:        "errors_if_tool_cache_dir_env_variable_is_not_defined",
			tc:          "",
			tool:        "test-tool",
			arch:        "amd64",
			versionSpec: "1.0.0",
			wantErr:     true,
		},
		{
			name:        "errors_if_tool_is_not_defined",
			tc:          tc,
			tool:        "",
			arch:        "amd64",
			versionSpec: "1.0.0",
			wantErr:     true,
		},
		{
			name:        "errors_if_arch_is_not_defined",
			tc:          tc,
			tool:        "test-tool",
			arch:        "",
			versionSpec: "1.0.0",
			wantErr:     true,
		},
		{
			name:        "errors_if_version_spec_is_not_defined",
			tc:          tc,
			tool:        "test-tool",
			arch:        "amd64",
			versionSpec: "",
			wantErr:     true,
		},
		{
			name:        "returns_an_empty_string_if_tool_is_not_cached",
			tc:          tc,
			tool:        "test",
			arch:        "amd64",
			versionSpec: "1.0.0",
			want:        "",
		},
		{
			name:        "returns_an_empty_string_if_tool_is_not_cached_with_a_marker",
			tc:          tc,
			tool:        "test-tool2",
			arch:        "amd64",
			versionSpec: "1.0.0",
			want:        "",
		},
		{
			name:        "returns_explicit_version_match",
			tc:          tc,
			tool:        "test-tool",
			arch:        "amd64",
			versionSpec: "1.0.0",
			want:        filepath.Join(tc, "test-tool", "1.0.0", getNodeArch("amd64")),
		},
		{
			name:        "returns_patch_version_match",
			tc:          tc,
			tool:        "test-tool",
			arch:        "amd64",
			versionSpec: "~1.0.0",
			want:        filepath.Join(tc, "test-tool", "1.0.1", getNodeArch("amd64")),
		},
		{
			name:        "returns_minor_version_match",
			tc:          tc,
			tool:        "test-tool",
			arch:        "amd64",
			versionSpec: "^1.0.0",
			want:        filepath.Join(tc, "test-tool", "1.2.0", getNodeArch("amd64")),
		},
		{
			name:        "returns_major_version_match",
			tc:          tc,
			tool:        "test-tool",
			arch:        "amd64",
			versionSpec: ">=1.0.0",
			want:        filepath.Join(tc, "test-tool", "2.0.0", getNodeArch("amd64")),
		},
		{
			name:        "returns_version_match_ignoring_v_prefix",
			tc:          tc,
			tool:        "test-tool",
			arch:        "amd64",
			versionSpec: "v1.0.0",
			want:        filepath.Join(tc, "test-tool", "1.0.0", getNodeArch("amd64")),
		},
		{
			name:        "returns_version_match_with_wildcard",
			tc:          tc,
			tool:        "test-tool",
			arch:        "amd64",
			versionSpec: "*",
			want:        filepath.Join(tc, "test-tool", "2.0.0", getNodeArch("amd64")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			t.Setenv(runnerToolCacheLookup, tt.tc)

			tp, err := FindTool(tt.tool, tt.arch, tt.versionSpec)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)         // should not error
			is.Equal(tp, tt.want) // should match expected
		})
	}
}

func TestCacheDir(t *testing.T) {
	tests := []struct {
		name    string
		tc      string
		source  string
		tool    string
		version string
		arch    string
		wantErr bool
	}{
		{
			name:    "errors_if_the_source_path_does_not_exist",
			source:  "non-existent-path",
			tool:    "test-tool",
			version: "1.0.0",
			arch:    "amd64",
			wantErr: true,
		},
		{
			name:    "errors_if_the_tool_is_not_defined",
			tool:    "",
			version: "1.0.0",
			arch:    "amd64",
			wantErr: true,
		},
		{
			name:    "errors_if_the_version_spec_is_not_defined",
			tool:    "test-tool",
			version: "",
			arch:    "amd64",
			wantErr: true,
		},
		{
			name:    "errors_if_the_arch_is_not_defined",
			tool:    "test-tool",
			version: "1.0.0",
			arch:    "",
			wantErr: true,
		},
		{
			name:    "errors_if_the_tool_cache_path_is_not_defined",
			tc:      "",
			source:  "non-existent-path",
			tool:    "test-tool",
			version: "1.0.0",
			arch:    "amd64",
			wantErr: true,
		},
		{
			name:    "caches_tool_dir",
			source:  "../../../testdata/test-tool",
			tool:    "test-tool",
			version: "1.0.0",
			arch:    "amd64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			tc := tt.tc
			if tc == "" && !tt.wantErr {
				tc = t.TempDir()
			}
			t.Setenv(runnerToolCacheLookup, tc)

			source := tt.source
			if source == "" {
				source = t.TempDir()
			}

			p, err := CacheDir(source, tt.tool, tt.version, tt.arch)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)                                                             // should not error
			is.Equal(p, filepath.Join(tc, tt.tool, tt.version, getNodeArch(tt.arch))) // should be equal

			toolExists, _ := fileio.FileExists(filepath.Join(p, "test-tool"))
			markerExists, _ := fileio.FileExists(getMarkerPath(p))

			is.True(toolExists)   // tool should exist
			is.True(markerExists) // marker file should exist
		})
	}
}

func TestCacheFile(t *testing.T) {
	tests := []struct {
		name       string
		tc         string
		source     string
		targetName string
		tool       string
		version    string
		arch       string
		wantErr    bool
	}{
		{
			name:       "errors_if_the_source_path_does_not_exist",
			source:     "non-existent-path",
			targetName: "test-tool",
			tool:       "test-tool",
			version:    "1.0.0",
			arch:       "amd64",
			wantErr:    true,
		},
		{
			name:       "errors_if_the_target_name_is_not_defined",
			targetName: "",
			tool:       "test-tool",
			version:    "1.0.0",
			arch:       "amd64",
			wantErr:    true,
		},
		{
			name:       "errors_if_the_tool_is_not_defined",
			targetName: "test-tool",
			tool:       "",
			version:    "1.0.0",
			arch:       "amd64",
			wantErr:    true,
		},
		{
			name:       "errors_if_the_version_spec_is_not_defined",
			targetName: "test-tool",
			tool:       "test-tool",
			version:    "",
			arch:       "amd64",
			wantErr:    true,
		},
		{
			name:       "errors_if_the_arch_is_not_defined",
			targetName: "test-tool",
			tool:       "test-tool",
			version:    "1.0.0",
			arch:       "",
			wantErr:    true,
		},
		{
			name:       "errors_if_the_tool_cache_path_is_not_defined",
			tc:         "",
			source:     "non-existent-path",
			targetName: "test-tool",
			tool:       "test-tool",
			version:    "1.0.0",
			arch:       "amd64",
			wantErr:    true,
		},
		{
			name:       "caches_tool_file",
			source:     "../../../testdata/test-tool/test-tool",
			targetName: "test-tool",
			tool:       "test-tool",
			version:    "1.0.0",
			arch:       "amd64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			tc := tt.tc
			if tc == "" && !tt.wantErr {
				tc = t.TempDir()
			}
			t.Setenv(runnerToolCacheLookup, tc)

			source := tt.source
			if source == "" {
				source = t.TempDir()
			}

			p, err := CacheFile(source, tt.targetName, tt.tool, tt.version, tt.arch)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)                                                             // should not error
			is.Equal(p, filepath.Join(tc, tt.tool, tt.version, getNodeArch(tt.arch))) // should be equal

			toolExists, _ := fileio.FileExists(filepath.Join(p, tt.targetName))
			markerExists, _ := fileio.FileExists(getMarkerPath(p))

			is.True(toolExists)   // tool should exist
			is.True(markerExists) // marker file should exist
		})
	}
}

func Test_createToolPath(t *testing.T) {
	tests := []struct {
		name    string
		tc      string
		tool    string
		version string
		arch    string
		wantErr bool
	}{
		{
			name:    "errors_if_the_tool_cache_path_is_not_defined",
			tc:      "UNSET",
			tool:    "tool",
			version: "1.0.0",
			arch:    "amd64",
			wantErr: true,
		},
		{
			name:    "errors_if_the_tool_is_not_defined",
			tool:    "",
			version: "1.0.0",
			arch:    "amd64",
			wantErr: true,
		},
		{
			name:    "errors_if_the_version_spec_is_not_defined",
			tool:    "tool",
			version: "",
			arch:    "amd64",
			wantErr: true,
		},
		{
			name:    "errors_if_the_arch_is_not_defined",
			tool:    "tool",
			version: "1.0.0",
			arch:    "",
			wantErr: true,
		},
		{
			name:    "creates_a_tool_path",
			tool:    "tool",
			version: "1.0.0",
			arch:    "amd64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			tc := tt.tc
			switch tc {
			case "":
				tc = t.TempDir()
			case "UNSET":
				tc = ""
			}
			t.Setenv(runnerToolCacheLookup, tc)

			vtp, err := createToolPath(tt.tool, tt.version, tt.arch)

			if tt.wantErr {
				is.True(err != nil) // should error
				return
			}

			is.NoErr(err)                                                  // should not error
			is.Equal(vtp, filepath.Join(tc, tt.tool, tt.version, tt.arch)) // should be equal

			dirExists, _ := fileio.DirExists(vtp)
			markerExists, _ := fileio.FileExists(getMarkerPath(vtp))

			is.True(dirExists)     // dir should exist
			is.True(!markerExists) // marker file should not exist
		})
	}
}

func Test_getMarkerPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "returns_the_marker_path",
			path: "/tmp/tool/1.0.0/x64",
			want: "/tmp/tool/1.0.0/x64.complete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			mp := getMarkerPath(tt.path)

			is.Equal(mp, tt.want) // should be equal
		})
	}
}

func Test_getNodeArch(t *testing.T) {
	tests := []struct {
		name string
		arch string
		want string
	}{
		{
			name: "returns_correct_amd64_value",
			arch: "amd64",
			want: "x64",
		},
		{
			name: "returns_correct_386_value",
			arch: "386",
			want: "ia32",
		},
		{
			name: "returns_correct_arm64_value",
			arch: "arm64",
			want: "arm64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			na := getNodeArch(tt.arch)

			is.Equal(na, tt.want) // should match expected
		})
	}
}
