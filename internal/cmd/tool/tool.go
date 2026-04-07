package tool

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"runtime"

	"github.com/urfave/cli/v3"

	"github.com/action-stars/ghactl/internal/toolkit/toolcache"
)

// Cmd provides the action logic for tool subcommands.
type Cmd struct{}

// New returns the fully-wired "tool" CLI command tree.
func New() *cli.Command {
	c := &Cmd{}

	return &cli.Command{
		Name:  "tool",
		Usage: "Manage GitHub runner tools.",
		Commands: []*cli.Command{
			c.cacheCommand(),
			c.downloadCommand(),
			c.extractCommand(),
			c.versionCommand(),
		},
	}
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
