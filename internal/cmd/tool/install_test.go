package tool

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"

	toolgithub "github.com/action-stars/ghactl/internal/toolkit/github"
)

func TestNew_Install(t *testing.T) {
	t.Run("installs_and_caches_tool_with_defaults", func(t *testing.T) {
		is := is.New(t)
		t.Setenv("RUNNER_TOOL_CACHE", t.TempDir())
		setupTempDir(t)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("binary-content"))
		}))
		defer ts.Close()

		c := &Cmd{
			releaseResolver: func(_ context.Context, token, owner, repo, toolName, version, osName, arch string, includePreRelease bool) (*toolgithub.ReleaseResolution, error) {
				is.Equal(token, "")
				is.Equal(owner, "sharkdp")
				is.Equal(repo, "bat")
				is.Equal(toolName, "bat")
				is.Equal(version, "latest")
				is.True(osName != "")
				is.True(arch != "")
				is.Equal(includePreRelease, false)

				return &toolgithub.ReleaseResolution{
					Version:   "1.2.3",
					AssetName: "bat-1.2.3-linux-x64",
					AssetURL:  ts.URL + "/bat",
				}, nil
			},
		}

		installedPath, installErr := c.Install(context.Background(), InstallOptions{Owner: "sharkdp", Repo: "bat", Version: "latest"})
		is.NoErr(installErr)
		is.True(installedPath != "")

		cachedPath, findErr := (&Cmd{}).CacheFind("bat", "", "1.2.3")
		is.NoErr(findErr)
		is.Equal(installedPath, cachedPath)
	})

	t.Run("short_circuits_when_resolved_version_is_already_cached", func(t *testing.T) {
		is := is.New(t)
		t.Setenv("RUNNER_TOOL_CACHE", t.TempDir())
		setupTempDir(t)

		source := createSourceFile(t)
		preCachedPath, cacheErr := (&Cmd{}).CacheFile(source, "bat", "bat", "1.2.3", "")
		is.NoErr(cacheErr)

		c := &Cmd{
			releaseResolver: func(_ context.Context, _, _, _, _, _, _, _ string, _ bool) (*toolgithub.ReleaseResolution, error) {
				return &toolgithub.ReleaseResolution{
					Version:   "1.2.3",
					AssetName: "bat-1.2.3-linux-x64",
					AssetURL:  "http://127.0.0.1:0/should-not-download",
				}, nil
			},
		}

		installedPath, err := c.Install(context.Background(), InstallOptions{Owner: "sharkdp", Repo: "bat", Version: "latest"})

		is.NoErr(err)
		is.Equal(installedPath, preCachedPath)
	})

	t.Run("uses_explicit_name_for_cache_path", func(t *testing.T) {
		is := is.New(t)
		t.Setenv("RUNNER_TOOL_CACHE", t.TempDir())
		setupTempDir(t)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("binary-content"))
		}))
		defer ts.Close()

		c := &Cmd{
			releaseResolver: func(_ context.Context, _, _, _, _, version, _, _ string, _ bool) (*toolgithub.ReleaseResolution, error) {
				is.Equal(version, "1.2.3")
				return &toolgithub.ReleaseResolution{
					Version:   "1.2.3",
					AssetName: "bat-1.2.3-linux-x64",
					AssetURL:  ts.URL + "/bat",
				}, nil
			},
		}

		installedPath, err := c.Install(context.Background(), InstallOptions{Owner: "sharkdp", Repo: "bat", Name: "custom-bat", Version: "1.2.3"})

		is.NoErr(err)
		is.True(strings.Contains(filepath.ToSlash(installedPath), "/custom-bat/1.2.3/"))
	})
}

func TestCmd_Install_validation(t *testing.T) {
	c := &Cmd{}

	t.Run("errors_when_owner_missing", func(t *testing.T) {
		is := is.New(t)

		_, err := c.Install(t.Context(), InstallOptions{Repo: "bat"})

		is.True(err != nil)
		is.Equal(err.Error(), "owner is not defined")
	})

	t.Run("errors_when_repo_missing", func(t *testing.T) {
		is := is.New(t)

		_, err := c.Install(t.Context(), InstallOptions{Owner: "sharkdp"})

		is.True(err != nil)
		is.Equal(err.Error(), "repo is not defined")
	})

	t.Run("returns_resolver_error", func(t *testing.T) {
		is := is.New(t)
		t.Setenv("RUNNER_TOOL_CACHE", t.TempDir())

		c := &Cmd{
			releaseResolver: func(_ context.Context, _, _, _, _, _, _, _ string, _ bool) (*toolgithub.ReleaseResolution, error) {
				return nil, fmt.Errorf("boom")
			},
		}

		_, err := c.Install(t.Context(), InstallOptions{Owner: "sharkdp", Repo: "bat", Version: "latest", Arch: "amd64", OS: "linux"})

		is.True(err != nil)
		is.Equal(err.Error(), "boom")
	})
}
