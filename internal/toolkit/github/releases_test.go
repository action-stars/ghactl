package github

import (
	"regexp"
	"testing"

	"github.com/google/go-github/v88/github"
	"github.com/matryer/is"
)

func Test_normalizeVersion(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		input string
		want  string
	}{
		{name: "plain_version", input: "1.2.3", want: "1.2.3"},
		{name: "v_prefixed", input: "v1.2.3", want: "1.2.3"},
		{name: "refs_tags_prefixed", input: "refs/tags/v1.2.3", want: "1.2.3"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			got := normalizeVersion(tt.input)

			is.Equal(got, tt.want)
		})
	}
}

func Test_selectRelease(t *testing.T) {
	t.Parallel()

	releases := []*github.RepositoryRelease{
		{TagName: new("v1.1.0")},
		{TagName: new("v1.2.0")},
		{TagName: new("v2.0.0-rc.1"), Prerelease: new(true)},
	}

	for _, tt := range []struct {
		name              string
		version           string
		includePrerelease bool
		wantTagName       string
		wantVersion       string
		wantErr           bool
	}{
		{
			name:              "selects_latest_stable_by_default",
			version:           "latest",
			includePrerelease: false,
			wantTagName:       "v1.2.0",
			wantVersion:       "1.2.0",
		},
		{
			name:              "selects_latest_including_prerelease_when_enabled",
			version:           "latest",
			includePrerelease: true,
			wantTagName:       "v2.0.0-rc.1",
			wantVersion:       "2.0.0-rc.1",
		},
		{
			name:              "matches_exact_without_v_prefix",
			version:           "1.1.0",
			includePrerelease: false,
			wantTagName:       "v1.1.0",
			wantVersion:       "1.1.0",
		},
		{
			name:              "matches_exact_with_v_prefix",
			version:           "v1.1.0",
			includePrerelease: false,
			wantTagName:       "v1.1.0",
			wantVersion:       "1.1.0",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			release, version, err := selectRelease(releases, tt.version, tt.includePrerelease)

			is.NoErr(err)
			is.Equal(release.GetTagName(), tt.wantTagName)
			is.Equal(version, tt.wantVersion)
		})
	}
}

func Test_selectAsset(t *testing.T) {
	t.Parallel()

	defaultAssets := []*github.ReleaseAsset{
		{Name: new("tool_v1.2.3_linux_x64.tar.gz"), BrowserDownloadURL: new("https://example.com/tar.gz")},
		{Name: new("tool_v1.2.3_linux_x64.zip"), BrowserDownloadURL: new("https://example.com/zip")},
		{Name: new("tool_v1.2.3_darwin_x64.tar.gz"), BrowserDownloadURL: new("https://example.com/darwin")},
		{Name: new("tool_v1.2.3.tar.gz"), BrowserDownloadURL: new("https://example.com/none")},
	}

	for _, tt := range []struct {
		name     string
		assets   []*github.ReleaseAsset
		toolName string
		toolRepo string
		goos     string
		goarch   string
		want     string
		wantErr  string
	}{
		{
			name:     "selects_best_matching_asset",
			assets:   defaultAssets,
			toolName: "tool",
			toolRepo: "toolrepo",
			goos:     "linux",
			goarch:   "amd64",
			want:     "tool_v1.2.3_linux_x64.tar.gz",
		},
		{
			name:     "fallback_to_name_match_when_os_arch_not_found",
			assets:   defaultAssets,
			toolName: "tool",
			toolRepo: "toolrepo",
			goos:     "windows",
			goarch:   "arm64",
			want:     "tool_v1.2.3.tar.gz",
		},
		{
			name:     "errors_when_no_name_match",
			assets:   defaultAssets,
			toolName: "unknown",
			toolRepo: "toolrepo",
			goos:     "windows",
			goarch:   "arm64",
			wantErr:  "no matching release asset found",
		},
		{
			name: "prefers_non_musl_over_musl_when_both_match",
			assets: []*github.ReleaseAsset{
				{Name: new("mdq-linux-x64-musl.tar.gz"), BrowserDownloadURL: new("https://example.com/musl")},
				{Name: new("mdq-linux-x64.tar.gz"), BrowserDownloadURL: new("https://example.com/glibc")},
			},
			toolName: "mdq",
			toolRepo: "mdq",
			goos:     "linux",
			goarch:   "amd64",
			want:     "mdq-linux-x64.tar.gz",
		},
		{
			name: "can_find_non_archive_asset",
			assets: []*github.ReleaseAsset{
				{Name: new("hadolint-linux-arm64"), BrowserDownloadURL: new("https://example.com/hadolint-linux-arm64")},
				{Name: new("hadolint-linux-x86_64"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64")},
				{Name: new("hadolint-windows-x86_64.exe"), BrowserDownloadURL: new("https://example.com/hadolint-windows-x86_64.exe")},
			},
			toolName: "hadolint",
			toolRepo: "hadolint",
			goos:     "linux",
			goarch:   "amd64",
			want:     "hadolint-linux-x86_64",
		},
		{
			name: "can_find_windows_exe",
			assets: []*github.ReleaseAsset{
				{Name: new("hadolint-linux-arm64"), BrowserDownloadURL: new("https://example.com/hadolint-linux-arm64")},
				{Name: new("hadolint-linux-x86_64"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64")},
				{Name: new("hadolint-windows-x86_64.exe"), BrowserDownloadURL: new("https://example.com/hadolint-windows-x86_64.exe")},
			},
			toolName: "hadolint",
			toolRepo: "hadolint",
			goos:     "windows",
			goarch:   "amd64",
			want:     "hadolint-windows-x86_64.exe",
		},
		{
			name: "ignores_checksum_assets_when_selecting",
			assets: []*github.ReleaseAsset{
				{Name: new("hadolint-linux-x86_64"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64")},
				{Name: new("hadolint-linux-x86_64.sha1"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64.sha1")},
				{Name: new("hadolint-linux-x86_64.sha256"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64.sha256")},
				{Name: new("hadolint-linux-x86_64.sha512"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64.sha512")},
				{Name: new("hadolint-linux-x86_64.sha256sum"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64.sha256sum")},
				{Name: new("hadolint-linux-x86_64.sha512sum"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64.sha512sum")},
				{Name: new("hadolint-linux-x86_64.sig"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64.sig")},
				{Name: new("hadolint-linux-x86_64.asc"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64.asc")},
			},
			toolName: "hadolint",
			toolRepo: "hadolint",
			goos:     "linux",
			goarch:   "amd64",
			want:     "hadolint-linux-x86_64",
		},
		{
			name: "errors_when_only_sidecar_assets_are_present",
			assets: []*github.ReleaseAsset{
				{Name: new("hadolint-linux-x86_64.sha256"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64.sha256")},
				{Name: new("hadolint-linux-x86_64.sha512"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64.sha512")},
				{Name: new("hadolint-linux-x86_64.sha256sum"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64.sha256sum")},
				{Name: new("hadolint-linux-x86_64.sha512sum"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64.sha512sum")},
				{Name: new("hadolint-linux-x86_64.sig"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64.sig")},
				{Name: new("hadolint-linux-x86_64.asc"), BrowserDownloadURL: new("https://example.com/hadolint-linux-x86_64.asc")},
			},
			toolName: "hadolint",
			toolRepo: "hadolint",
			goos:     "linux",
			goarch:   "amd64",
			wantErr:  "no matching release asset found",
		},
		{
			name: "can_find_self",
			assets: []*github.ReleaseAsset{
				{Name: new("checksums.txt"), BrowserDownloadURL: new("https://example.com/checksums.txt")},
				{Name: new("checksums.txt.sbom.json.bundle"), BrowserDownloadURL: new("https://example.com/checksums.txt.sbom.json.bundle")},
				{Name: new("ghactl_0.0.11_linux_amd64.tar.gz"), BrowserDownloadURL: new("https://example.com/ghactl_0.0.11_linux_amd64.tar.gz")},
				{Name: new("ghactl_0.0.11_linux_arm64.tar.gz"), BrowserDownloadURL: new("https://example.com/ghactl_0.0.11_linux_arm64.tar.gz")},
				{Name: new("ghactl_0.0.11_windows_amd64.tar.gz"), BrowserDownloadURL: new("https://example.com/ghactl_0.0.11_windows_amd64.tar.gz")},
				{Name: new("ghactl_0.0.11_windows_arm64.tar.gz"), BrowserDownloadURL: new("https://example.com/ghactl_0.0.11_windows_arm64.tar.gz")},
			},
			toolName: "ghactl",
			toolRepo: "ghactl",
			goos:     "linux",
			goarch:   "amd64",
			want:     "ghactl_0.0.11_linux_amd64.tar.gz",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			is := is.New(t)

			asset, err := selectAsset(tt.assets, tt.toolName, tt.toolRepo, tt.goos, tt.goarch)

			if tt.wantErr != "" {
				is.True(regexp.MustCompile(regexp.QuoteMeta(tt.wantErr)).MatchString(err.Error())) // Expect error to match
				return
			}

			is.NoErr(err)                      // Expect no error
			is.True(asset != nil)              // Expect asset to be non-nil
			is.Equal(asset.GetName(), tt.want) // Expect asset name to match
		})
	}
}
