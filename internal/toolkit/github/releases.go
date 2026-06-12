package github

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v88/github"
)

// ReleaseResolution is the selected release, asset, and normalized version to install.
type ReleaseResolution struct {
	Version   string
	AssetName string
	AssetURL  string
}

// ResolveToolRelease selects a release and matching asset for a tool install.
func ResolveToolRelease(ctx context.Context, token, owner, repo, version, osName, arch string, includePreRelease bool) (*ReleaseResolution, error) {
	if owner == "" {
		return nil, fmt.Errorf("owner is not defined")
	}

	if repo == "" {
		return nil, fmt.Errorf("repo is not defined")
	}

	if arch == "" {
		return nil, fmt.Errorf("arch is not defined")
	}

	if osName == "" {
		return nil, fmt.Errorf("os is not defined")
	}

	var tokenPtr *string
	if token != "" {
		tokenPtr = &token
	}

	client, err := GetClient(tokenPtr)
	if err != nil {
		return nil, err
	}

	releases, err := listReleases(ctx, client, owner, repo)
	if err != nil {
		return nil, err
	}

	release, normalizedVersion, err := selectRelease(releases, version, includePreRelease)
	if err != nil {
		return nil, err
	}

	asset, err := selectAsset(release.Assets, osName, arch)
	if err != nil {
		return nil, err
	}

	assetName := asset.GetName()
	assetURL := asset.GetBrowserDownloadURL()
	if assetName == "" || assetURL == "" {
		return nil, fmt.Errorf("resolved release asset is missing required fields")
	}

	return &ReleaseResolution{
		Version:   normalizedVersion,
		AssetName: assetName,
		AssetURL:  assetURL,
	}, nil
}

func listReleases(ctx context.Context, client *github.Client, owner, repo string) ([]*github.RepositoryRelease, error) {
	all := []*github.RepositoryRelease{}

	for release, err := range client.Repositories.ListReleasesIter(ctx, owner, repo, &github.ListOptions{PerPage: 100}) {
		if err != nil {
			return nil, err
		}

		all = append(all, release)
	}

	return all, nil
}

func selectRelease(releases []*github.RepositoryRelease, version string, includePreRelease bool) (*github.RepositoryRelease, string, error) {
	if len(releases) == 0 {
		return nil, "", fmt.Errorf("no releases found")
	}

	requested := strings.TrimSpace(version)
	if requested == "" {
		requested = "latest"
	}

	if strings.EqualFold(requested, "latest") {
		type candidate struct {
			release *github.RepositoryRelease
			ver     *semver.Version
		}

		candidates := []*candidate{}
		for _, release := range releases {
			if release.GetDraft() {
				continue
			}
			if release.GetPrerelease() && !includePreRelease {
				continue
			}

			normalized := normalizeVersion(release.GetTagName())
			if normalized == "" {
				continue
			}

			ver, err := semver.StrictNewVersion(normalized)
			if err != nil {
				continue
			}

			if ver.Prerelease() != "" && !includePreRelease {
				continue
			}

			candidates = append(candidates, &candidate{release: release, ver: ver})
		}

		if len(candidates) == 0 {
			return nil, "", fmt.Errorf("no matching releases found")
		}

		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].ver.LessThan(candidates[j].ver)
		})

		selected := candidates[len(candidates)-1]
		return selected.release, selected.ver.String(), nil
	}

	normalizedRequested := normalizeVersion(requested)
	if normalizedRequested == "" {
		return nil, "", fmt.Errorf("version is not defined")
	}

	_, err := semver.StrictNewVersion(normalizedRequested)
	if err != nil {
		return nil, "", err
	}

	for _, release := range releases {
		if release.GetDraft() {
			continue
		}

		if normalizeVersion(release.GetTagName()) == normalizedRequested {
			return release, normalizedRequested, nil
		}
	}

	return nil, "", fmt.Errorf("release %s not found", requested)
}

func normalizeVersion(value string) string {
	version := strings.TrimSpace(value)
	version = strings.TrimPrefix(version, "refs/tags/")
	version = strings.TrimPrefix(version, "v")
	return version
}

func selectAsset(assets []*github.ReleaseAsset, osName, arch string) (*github.ReleaseAsset, error) {
	if len(assets) == 0 {
		return nil, fmt.Errorf("no release assets found")
	}

	osTokens := mapOS(osName)
	archTokens := mapArch(arch)

	type candidate struct {
		asset *github.ReleaseAsset
		score int
	}

	candidates := []*candidate{}
	for _, asset := range assets {
		if asset.GetBrowserDownloadURL() == "" {
			continue
		}

		name := strings.ToLower(asset.GetName())
		osScore := tokenScore(name, osTokens)
		archScore := tokenScore(name, archTokens)
		if osScore == 0 || archScore == 0 {
			continue
		}

		score := (osScore * 100) + (archScore * 100) + archiveScore(name)
		candidates = append(candidates, &candidate{asset: asset, score: score})
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no matching release asset found for os=%s arch=%s", osName, arch)
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	best := candidates[0]
	for i := 1; i < len(candidates); i++ {
		if candidates[i].score != best.score {
			break
		}

		if candidates[i].asset.GetName() != best.asset.GetName() {
			return nil, fmt.Errorf("ambiguous release assets for os=%s arch=%s", osName, arch)
		}
	}

	return best.asset, nil
}

func tokenScore(name string, tokens []string) int {
	score := 0
	for _, token := range tokens {
		if strings.Contains(name, token) {
			score++
		}
	}
	return score
}

func archiveScore(name string) int {
	switch {
	case strings.HasSuffix(name, ".tar.gz"):
		return 4
	case strings.HasSuffix(name, ".tgz"):
		return 3
	case strings.HasSuffix(name, ".tar"):
		return 2
	case strings.HasSuffix(name, ".zip"):
		return 1
	default:
		return 0
	}
}

func mapOS(osName string) []string {
	switch osName {
	case "darwin":
		return []string{"darwin", "macos", "mac"}
	case "windows":
		return []string{"windows", "win"}
	case "linux":
		return []string{"linux"}
	default:
		return []string{strings.ToLower(osName)}
	}
}

func mapArch(arch string) []string {
	switch arch {
	case "amd64":
		return []string{"amd64", "x86_64", "x64"}
	case "arm64":
		return []string{"arm64", "aarch64"}
	case "386":
		return []string{"386", "x86", "ia32", "i386"}
	default:
		return []string{strings.ToLower(arch)}
	}
}
