package semver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version represents a parsed semantic version.
type Version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
	Build      string
}

// ParseVersion parses a semantic version string and returns a Version struct or an error if the version is invalid.
func ParseVersion(version string) (*Version, error) {
	// Define the regex pattern for a valid semantic version based on the provided BNF grammar.
	var semverRegex = `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)` +
		`(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?` +
		`(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`

	re := regexp.MustCompile(semverRegex)
	matches := re.FindStringSubmatch(version)

	if matches == nil {
		return nil, errors.New("invalid version")
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])
	preRelease := matches[4]
	build := matches[5]

	return &Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		PreRelease: preRelease,
		Build:      build,
	}, nil
}

// String returns the string representation of the Version.
func (v *Version) String() string {
	version := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.PreRelease != "" {
		version += "-" + v.PreRelease
	}
	if v.Build != "" {
		version += "+" + v.Build
	}
	return version
}

// compareIdentifiers compares two identifiers according to SemVer rules.
func compareIdentifiers(a, b string) int {
	aIsNumeric := isNumeric(a)
	bIsNumeric := isNumeric(b)

	if aIsNumeric && bIsNumeric {
		aNum, _ := strconv.Atoi(a)
		bNum, _ := strconv.Atoi(b)
		return compareInts(aNum, bNum)
	}

	if aIsNumeric {
		return -1
	}

	if bIsNumeric {
		return 1
	}

	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

// isNumeric checks if a string consists only of digits.
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// compareInts compares two integers.
func compareInts(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

// comparePreRelease compares two pre-release versions according to SemVer rules.
func comparePreRelease(a, b string) int {
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		return 1
	}
	if b == "" {
		return -1
	}

	aIdentifiers := strings.Split(a, ".")
	bIdentifiers := strings.Split(b, ".")
	maxLen := len(aIdentifiers)
	if len(bIdentifiers) > maxLen {
		maxLen = len(bIdentifiers)
	}

	for i := 0; i < maxLen; i++ {
		if i >= len(aIdentifiers) {
			return -1
		}
		if i >= len(bIdentifiers) {
			return 1
		}
		result := compareIdentifiers(aIdentifiers[i], bIdentifiers[i])
		if result != 0 {
			return result
		}
	}

	return 0
}

// Equals determines if this version is equal to the provided version.
func (v *Version) Equals(other *Version) bool {
	return v.Major == other.Major &&
		v.Minor == other.Minor &&
		v.Patch == other.Patch &&
		v.PreRelease == other.PreRelease
}

// Before determines if this version is before the provided version.
func (v *Version) Before(other *Version) bool {
	if v.Major != other.Major {
		return v.Major < other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor < other.Minor
	}
	if v.Patch != other.Patch {
		return v.Patch < other.Patch
	}
	return comparePreRelease(v.PreRelease, other.PreRelease) == -1
}
