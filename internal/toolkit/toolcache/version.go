package toolcache

import "github.com/Masterminds/semver/v3"

// CheckVersion checks if a version meets a version specification.
func CheckVersion(version, versionSpec string) (bool, error) {
	v, err := semver.StrictNewVersion(version)
	if err != nil {
		return false, err
	}

	c, err := semver.NewConstraint(versionSpec)
	if err != nil {
		return false, err
	}

	return c.Check(v), nil
}
