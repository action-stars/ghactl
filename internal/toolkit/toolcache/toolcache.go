package toolcache

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/Masterminds/semver/v3"

	"github.com/action-stars/ghactl/internal/fileio"
)

// runnerToolCacheLookup is the environment variable used to find the GitHub Actions runner tool cache directory.
const runnerToolCacheLookup = "RUNNER_TOOL_CACHE"

// GetToolCacheDirectory returns the path to the GitHub Actions runner tool cache directory.
func GetToolCacheDirectory() (string, error) {
	d := os.Getenv(runnerToolCacheLookup)
	if d == "" {
		return "", fmt.Errorf("%s is not defined", runnerToolCacheLookup)
	}

	exists, err := fileio.DirExists(d)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", fmt.Errorf("dir %s does not exist", d)
	}

	return d, nil
}

// FindAllToolVersions returns all versions of a tool in the GitHub Actions runner tool cache.
func FindAllToolVersions(tool, arch string) ([]string, error) {
	d, err := GetToolCacheDirectory()
	if err != nil {
		return nil, err
	}

	if tool == "" {
		return nil, fmt.Errorf("tool is not defined")
	}

	if arch == "" {
		return nil, fmt.Errorf("arch is not defined")
	}

	toolPath := filepath.Join(d, tool)
	exists, err := fileio.DirExists(toolPath)
	if err != nil {
		return nil, err
	}

	versions := []string{}

	if !exists {
		return versions, nil
	}

	items, err := os.ReadDir(toolPath)
	if err != nil {
		return nil, err
	}

	nodeArch := getNodeArch(arch)

	for _, item := range items {
		if item.IsDir() {
			_, err := semver.StrictNewVersion(item.Name())
			if err == nil {
				complete, _ := fileio.FileExists(getMarkerPath(filepath.Join(toolPath, item.Name(), nodeArch)))
				if complete {
					versions = append(versions, item.Name())
				}
			}
		}
	}

	return versions, nil
}

// FindTool finds a tool in the GitHub Actions runner tool cache that matches the version constraint.
// If the versionSpec isn't an explicit version then it will be evaluated as a semver constraint.
// Will return the path to the tool or an empty string if no tool is found.
func FindTool(tool, arch, versionSpec string) (string, error) {
	d, err := GetToolCacheDirectory()
	if err != nil {
		return "", err
	}

	if tool == "" {
		return "", fmt.Errorf("tool is not defined")
	}

	if arch == "" {
		return "", fmt.Errorf("arch is not defined")
	}

	if versionSpec == "" {
		return "", fmt.Errorf("versionSpec is not defined")
	}

	toolPath := filepath.Join(d, tool)
	exists, err := fileio.DirExists(toolPath)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", nil
	}

	items, err := os.ReadDir(toolPath)
	if err != nil {
		return "", err
	}

	c, err := semver.NewConstraint(versionSpec)
	if err != nil {
		return "", err
	}

	vs := []*semver.Version{}

	nodeArch := getNodeArch(arch)

	for _, item := range items {
		if item.IsDir() {
			v, err := semver.StrictNewVersion(item.Name())
			if err == nil {
				complete, _ := fileio.FileExists(getMarkerPath(filepath.Join(toolPath, item.Name(), nodeArch)))
				if complete {
					if c.Check(v) {
						vs = append(vs, v)
					}
				}
			}
		}
	}

	if len(vs) == 0 {
		return "", nil
	}

	if len(vs) > 1 {
		sort.Sort(semver.Collection(vs))
	}

	return filepath.Join(toolPath, vs[len(vs)-1].String(), nodeArch), nil
}

// CacheDir caches a tool dir into the GitHub Actions runner tool cache.
// If addPath is true, the tool will be added to the PATH.
func CacheDir(source, tool, version, arch string) (string, error) {
	sourceExists, err := fileio.DirExists(source)
	if err != nil {
		return "", err
	}
	if !sourceExists {
		return "", fmt.Errorf("source %s does not exist", source)
	}

	nodeArch := getNodeArch(arch)

	toolPath, err := createToolPath(tool, version, nodeArch)
	if err != nil {
		return "", err
	}

	if err := os.CopyFS(toolPath, os.DirFS(source)); err != nil {
		return "", err
	}

	markerPath := getMarkerPath(toolPath)
	marker, err := os.Create(markerPath)
	if err != nil {
		return "", err
	}
	marker.Close()

	return toolPath, nil
}

// CacheFile caches a tool file into the GitHub Actions runner tool cache.
// If addPath is true, the tool will be added to the PATH.
func CacheFile(source, targetName, tool, version, arch string) (string, error) {
	sourceExists, err := fileio.FileExists(source)
	if err != nil {
		return "", err
	}
	if !sourceExists {
		return "", fmt.Errorf("source %s does not exist", source)
	}

	if targetName == "" {
		return "", fmt.Errorf("targetName is not defined")
	}

	nodeArch := getNodeArch(arch)

	toolPath, err := createToolPath(tool, version, nodeArch)
	if err != nil {
		return "", err
	}

	if err := fileio.CopyFile(source, filepath.Join(toolPath, targetName), false); err != nil {
		return "", err
	}

	markerPath := getMarkerPath(toolPath)
	marker, err := os.Create(markerPath)
	if err != nil {
		return "", err
	}
	marker.Close()

	return toolPath, nil
}

// createToolPath creates a path to a specific version and architecture of a tool in the GitHub Actions runner tool cache.
func createToolPath(tool, version, arch string) (string, error) {
	cacheDir, err := GetToolCacheDirectory()
	if err != nil {
		return "", err
	}

	if tool == "" {
		return "", fmt.Errorf("tool is not defined")
	}

	if version == "" {
		return "", fmt.Errorf("version is not defined")
	}

	if arch == "" {
		return "", fmt.Errorf("arch is not defined")
	}

	toolPath := filepath.Join(cacheDir, tool, version, arch)
	markerPath := getMarkerPath(toolPath)

	if err := os.RemoveAll(toolPath); err != nil {
		return "", err
	}

	if err := os.Remove(markerPath); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return "", err
	}

	if err := os.MkdirAll(toolPath, 0o755); err != nil {
		return "", err
	}

	return toolPath, nil
}

// getMarkerPath returns the path to the marker file for a specific version and architecture of a tool in the GitHub Actions runner tool cache.
func getMarkerPath(toolPath string) string {
	return fmt.Sprintf("%s.complete", toolPath)
}

// getNodeArch returns the architecture in Node.js format.
func getNodeArch(arch string) string {
	switch arch {
	case "amd64":
		return "x64"
	case "386":
		return "ia32"
	default:
		return arch
	}
}
