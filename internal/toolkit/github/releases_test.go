package github

import (
	"testing"

	"github.com/google/go-github/v88/github"
	"github.com/matryer/is"
)

func Test_normalizeVersion(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "plain_version", input: "1.2.3", want: "1.2.3"},
		{name: "v_prefixed", input: "v1.2.3", want: "1.2.3"},
		{name: "refs_tags_prefixed", input: "refs/tags/v1.2.3", want: "1.2.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			got := normalizeVersion(tt.input)

			is.Equal(got, tt.want)
		})
	}
}

func Test_selectRelease(t *testing.T) {
	releases := []*github.RepositoryRelease{
		{TagName: new("v1.1.0")},
		{TagName: new("v1.2.0")},
		{TagName: new("v2.0.0-rc.1"), Prerelease: new(true)},
	}

	t.Run("selects_latest_stable_by_default", func(t *testing.T) {
		is := is.New(t)

		release, version, err := selectRelease(releases, "latest", false)

		is.NoErr(err)
		is.Equal(release.GetTagName(), "v1.2.0")
		is.Equal(version, "1.2.0")
	})

	t.Run("selects_latest_including_prerelease_when_enabled", func(t *testing.T) {
		is := is.New(t)

		release, version, err := selectRelease(releases, "latest", true)

		is.NoErr(err)
		is.Equal(release.GetTagName(), "v2.0.0-rc.1")
		is.Equal(version, "2.0.0-rc.1")
	})

	t.Run("matches_exact_without_v_prefix", func(t *testing.T) {
		is := is.New(t)

		release, version, err := selectRelease(releases, "1.1.0", false)

		is.NoErr(err)
		is.Equal(release.GetTagName(), "v1.1.0")
		is.Equal(version, "1.1.0")
	})

	t.Run("matches_exact_with_v_prefix", func(t *testing.T) {
		is := is.New(t)

		release, version, err := selectRelease(releases, "v1.1.0", false)

		is.NoErr(err)
		is.Equal(release.GetTagName(), "v1.1.0")
		is.Equal(version, "1.1.0")
	})
}

func Test_selectAsset(t *testing.T) {
	assets := []*github.ReleaseAsset{
		{Name: new("tool_v1.2.3_linux_x64.tar.gz"), BrowserDownloadURL: new("https://example.com/tar.gz")},
		{Name: new("tool_v1.2.3_linux_x64.zip"), BrowserDownloadURL: new("https://example.com/zip")},
		{Name: new("tool_v1.2.3_darwin_x64.tar.gz"), BrowserDownloadURL: new("https://example.com/darwin")},
	}

	t.Run("selects_best_matching_asset", func(t *testing.T) {
		is := is.New(t)

		asset, err := selectAsset(assets, "linux", "amd64")

		is.NoErr(err)
		is.Equal(asset.GetName(), "tool_v1.2.3_linux_x64.tar.gz")
	})

	t.Run("errors_when_no_match", func(t *testing.T) {
		is := is.New(t)

		_, err := selectAsset(assets, "windows", "arm64")

		is.True(err != nil)
	})
}
