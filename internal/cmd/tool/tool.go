package tool

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/action-stars/ghactl/internal/toolkit/core"
	"github.com/action-stars/ghactl/internal/toolkit/github"
	"github.com/action-stars/ghactl/internal/toolkit/toolcache"
)

// Cmd provides the action logic for tool subcommands.
type Cmd struct {
	releaseResolver func(ctx context.Context, token, owner, repo, toolName, version, osName, arch string, includePreRelease bool) (*github.ReleaseResolution, error)
}

// New returns the fully-wired "tool" CLI command tree.
func New() *cli.Command {
	c := &Cmd{}

	return &cli.Command{
		Name:  "tool",
		Usage: "Manage GitHub runner tools.",
		Commands: []*cli.Command{
			c.cacheCommand(),
			c.installCommand(),
			c.downloadCommand(),
			c.extractCommand(),
			c.versionCommand(),
		},
	}
}

func (c *Cmd) resolveRelease(ctx context.Context, token, owner, repo, toolName, version, osName, arch string, includePreRelease bool) (*github.ReleaseResolution, error) {
	if c.releaseResolver != nil {
		return c.releaseResolver(ctx, token, owner, repo, toolName, version, osName, arch, includePreRelease)
	}
	return github.ResolveToolRelease(ctx, token, owner, repo, toolName, version, osName, arch, includePreRelease)
}

// CacheGet returns the tool cache directory path.
func (c *Cmd) CacheGet() (string, error) {
	return toolcache.GetToolCacheDirectory()
}

// CacheFindAll returns all cached paths for a tool.
func (c *Cmd) CacheFindAll(tool, arch string) ([]string, error) {
	if arch == "" {
		arch = runtime.GOARCH
	}
	return toolcache.FindAllToolVersions(tool, arch)
}

// CacheFind returns the path to a cached tool matching the version spec.
func (c *Cmd) CacheFind(tool, arch, versionSpec string) (string, error) {
	if arch == "" {
		arch = runtime.GOARCH
	}
	if versionSpec == "" {
		versionSpec = "*"
	}
	return toolcache.FindTool(tool, arch, versionSpec)
}

// CacheDir caches a directory as a tool in the runner tool cache.
func (c *Cmd) CacheDir(source, tool, version, arch string) (string, error) {
	if arch == "" {
		arch = runtime.GOARCH
	}
	return toolcache.CacheDir(source, tool, version, arch)
}

// CacheFile caches a file as a tool in the runner tool cache.
func (c *Cmd) CacheFile(source, targetName, tool, version, arch string) (string, error) {
	if arch == "" {
		arch = runtime.GOARCH
	}
	if targetName == "" {
		targetName = tool
	}
	return toolcache.CacheFile(source, targetName, tool, version, arch)
}

// Download downloads a tool from a URL to a temporary directory.
func (c *Cmd) Download(ctx context.Context, rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return toolcache.DownloadTool(ctx, slog.Default(), *u)
}

// ExtractTar extracts a tar archive to a temporary directory.
func (c *Cmd) ExtractTar(path string, gz bool) (string, error) {
	return toolcache.ExtractTar(path, gz)
}

// ExtractZip extracts a zip archive to a temporary directory.
func (c *Cmd) ExtractZip(path string) (string, error) {
	return toolcache.ExtractZip(path)
}

// VersionCheck checks if a version satisfies a version spec.
func (c *Cmd) VersionCheck(version, versionSpec string) (bool, error) {
	return toolcache.CheckVersion(version, versionSpec)
}

// InstallOptions is the set of options used to install a tool.
type InstallOptions struct {
	Name              string
	Owner             string
	Repo              string
	Version           string
	Arch              string
	OS                string
	IncludePreRelease bool
	Token             string
	AddToPath         bool
}

// Install resolves, downloads, and caches a tool release. It returns the cached tool path.
func (c *Cmd) Install(ctx context.Context, options InstallOptions) (string, error) {
	if options.Owner == "" {
		return "", fmt.Errorf("owner is not defined")
	}

	if options.Repo == "" {
		return "", fmt.Errorf("repo is not defined")
	}

	name := options.Name
	if name == "" {
		name = options.Repo
	}

	arch := defaultArch(options.Arch)
	version := strings.TrimSpace(options.Version)
	if version == "" {
		version = "latest"
	}

	osName := options.OS
	if osName == "" {
		osName = runtime.GOOS
	}

	resolution, err := c.resolveRelease(ctx, options.Token, options.Owner, options.Repo, name, version, osName, arch, options.IncludePreRelease)
	if err != nil {
		return "", err
	}

	cachedPath, err := c.CacheFind(name, arch, resolution.Version)
	if err != nil {
		return "", err
	}

	if cachedPath != "" {
		if options.AddToPath {
			return cachedPath, core.AddPath(cachedPath)
		}
		return cachedPath, nil
	}

	downloadPath, err := c.Download(ctx, resolution.AssetURL)
	if err != nil {
		return "", err
	}

	assetName := strings.ToLower(resolution.AssetName)

	switch {
	case strings.HasSuffix(assetName, ".tar.gz") || strings.HasSuffix(assetName, ".tgz"):
		extractedPath, err := c.ExtractTar(downloadPath, true)
		if err != nil {
			return "", err
		}

		resolvedPath, err := toolcache.ResolveToolDirectory(extractedPath)
		if err != nil {
			return "", err
		}

		cachedPath, err = c.CacheDir(resolvedPath, name, resolution.Version, arch)
		if err != nil {
			return "", err
		}
	case strings.HasSuffix(assetName, ".tar"):
		extractedPath, err := c.ExtractTar(downloadPath, false)
		if err != nil {
			return "", err
		}

		resolvedPath, err := toolcache.ResolveToolDirectory(extractedPath)
		if err != nil {
			return "", err
		}

		cachedPath, err = c.CacheDir(resolvedPath, name, resolution.Version, arch)
		if err != nil {
			return "", err
		}
	case strings.HasSuffix(assetName, ".zip"):
		extractedPath, err := c.ExtractZip(downloadPath)
		if err != nil {
			return "", err
		}

		resolvedPath, err := toolcache.ResolveToolDirectory(extractedPath)
		if err != nil {
			return "", err
		}

		cachedPath, err = c.CacheDir(resolvedPath, name, resolution.Version, arch)
		if err != nil {
			return "", err
		}
	default:
		targetName := filepath.Base(name)
		if targetName == "." || targetName == string(filepath.Separator) {
			targetName = name
		}

		cachedPath, err = c.CacheFile(downloadPath, targetName, name, resolution.Version, arch)
		if err != nil {
			return "", err
		}
	}

	if options.AddToPath {
		return cachedPath, core.AddPath(cachedPath)
	}

	return cachedPath, nil
}

func writeOutput(cmd *cli.Command, v any) error {
	_, err := fmt.Fprintln(cmd.Root().Writer, v)
	if err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func exitErr(err error) error {
	return cli.Exit(err, 1)
}

func defaultArch(arch string) string {
	if arch == "" {
		return runtime.GOARCH
	}
	return arch
}

func archFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "arch",
		Usage: "Architecture of the tool.",
	}
}

func osFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "os",
		Usage: "Operating system of the tool.",
	}
}

func toolNameFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "name",
		Usage:    "Name of the tool.",
		Required: true,
	}
}

func versionFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "version",
		Usage:    "Version of the tool.",
		Required: true,
	}
}
