package github

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v88/github"
)

var (
	// osIndicators maps OS names to their file name indicators.
	osIndicators = map[string][]string{
		"darwin":  {"darwin", "macos", "mac"},
		"windows": {"windows", "win"},
		"linux":   {"linux"},
	}

	// archIndicators maps arch names to their file name indicators.
	archIndicators = map[string][]string{
		"amd64": {"amd64", "x86_64", "x64", "64bit"},
		"arm64": {"arm64", "aarch64"},
		"386":   {"386", "x86", "ia32", "i386", "32bit"},
	}

	// allOSIndicators is a flattened list of all known OS indicators.
	allOSIndicators = flatten(osIndicators)

	// allArchIndicators is a flattened list of all known arch indicators.
	allArchIndicators = flatten(archIndicators)
)

// ReleaseResolution is the selected release, asset, and normalized version to install.
type ReleaseResolution struct {
	Version   string
	AssetName string
	AssetURL  string
}

// ResolveToolRelease selects a release and matching asset for a tool install.
// Assets are matched primarily by tool name (exact or substring match preferred),
// with optional fallback to repo name. OS and arch indicators are optional but
// preferred when present. Returns an error if no matching asset is found.
func ResolveToolRelease(ctx context.Context, token, owner, repo, toolName, version, osName, arch string, includePreRelease bool) (*ReleaseResolution, error) {
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

	asset, err := selectAsset(release.Assets, toolName, repo, osName, arch)
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

// selectAsset selects the best matching release asset based on tool name and OS/arch.
// Scoring priority:
// 1. Tool name match (exact, substring, or fallback to repo name)
// 2. OS and arch specificity (optional but preferred)
// 3. Archive format preference (tar.gz > tgz > tar > zip)
// 4. Variant preference (non-musl preferred over musl, negative scoring allows fallback while penalizing undesirable variants)
// Assets matching by name are accepted even without OS/arch specificity.
func selectAsset(assets []*github.ReleaseAsset, toolName, repo, osName, arch string) (*github.ReleaseAsset, error) {
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

		assetName := strings.ToLower(asset.GetName())

		// Score based on name match first (tool name preferred, then repo name)
		nameScore := scoreNameMatch(assetName, toolName, repo)
		if nameScore == 0 {
			continue
		}

		// Score OS and arch matches (0 if no match, can be multiple matches)
		osScore := tokenScore(assetName, osTokens)
		archScore := tokenScore(assetName, archTokens)

		// All name-matched assets are candidates; scoring determines ranking.
		// Higher scores for better matches, negative scoring penalizes undesirable variants.
		score := (nameScore * 100000) + (osScore * 10000) + (archScore * 10000) + archiveScore(assetName) + variantScore(assetName)
		candidates = append(candidates, &candidate{asset: asset, score: score})
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no matching release asset found for tool=%s repo=%s os=%s arch=%s", toolName, repo, osName, arch)
	}

	// Filter: reject assets with explicit OS/arch indicators that don't match our requested OS/arch.
	// Assets without explicit indicators are generic fallbacks and always acceptable.
	filtered := []*candidate{}
candidateLoop:
	for _, c := range candidates {
		assetName := strings.ToLower(c.asset.GetName())

		// If asset has an OS indicator, it must match our requested OS
		if hasOSIndicator(assetName) && tokenScore(assetName, osTokens) == 0 {
			continue candidateLoop
		}

		// If asset has an arch indicator, it must match our requested arch
		if hasArchIndicator(assetName) && tokenScore(assetName, archTokens) == 0 {
			continue candidateLoop
		}

		// Asset passed all checks: accept it
		filtered = append(filtered, c)
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no matching release asset found for tool=%s repo=%s os=%s arch=%s", toolName, repo, osName, arch)
	}

	candidates = filtered

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	best := candidates[0]
	for i := 1; i < len(candidates); i++ {
		if candidates[i].score != best.score {
			break
		}

		if candidates[i].asset.GetName() != best.asset.GetName() {
			return nil, fmt.Errorf("ambiguous release assets for tool=%s repo=%s os=%s arch=%s", toolName, repo, osName, arch)
		}
	}

	return best.asset, nil
}

// scoreNameMatch returns a score indicating how well the asset name matches the tool.
// Scoring: 3 for exact match on tool name, 2 for tool name substring,
// 1 for repo name match, 0 for no match.
func scoreNameMatch(assetName, toolName, repo string) int {
	// Prefer exact match on tool name, then repo name
	if assetName == toolName {
		return 3 // Exact match on tool name
	}
	if strings.Contains(assetName, toolName) {
		return 2 // Tool name is a substring (e.g., tool name in larger filename)
	}
	if repo != "" && repo != toolName {
		if assetName == repo {
			return 1 // Exact match on repo name (fallback)
		}
		if strings.Contains(assetName, repo) {
			return 1 // Repo name is a substring (fallback)
		}
	}
	return 0 // No match
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

func variantScore(name string) int {
	if strings.Contains(name, "musl") {
		return -1
	}

	return 0
}

// hasOSIndicator checks if name contains any known OS indicator.
func hasOSIndicator(name string) bool {
	for _, indicator := range allOSIndicators {
		if strings.Contains(name, indicator) {
			return true
		}
	}
	return false
}

// hasArchIndicator checks if name contains any known arch indicator.
func hasArchIndicator(name string) bool {
	for _, indicator := range allArchIndicators {
		if strings.Contains(name, indicator) {
			return true
		}
	}
	return false
}

func mapOS(osName string) []string {
	if tokens, ok := osIndicators[osName]; ok {
		return tokens
	}
	return []string{strings.ToLower(osName)}
}

func mapArch(arch string) []string {
	if tokens, ok := archIndicators[arch]; ok {
		return tokens
	}
	return []string{strings.ToLower(arch)}
}

// flatten extracts all values from a map and returns them as a single slice.
func flatten(m map[string][]string) []string {
	var result []string
	for _, tokens := range m {
		result = append(result, tokens...)
	}
	return result
}
