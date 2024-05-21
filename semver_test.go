package semver_test

import (
	"testing"

	"github.com/networkteam/semver"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		version     string
		major       int
		minor       int
		patch       int
		preRelease  string
		build       string
		expectedErr string
	}{
		{"1.0.0", 1, 0, 0, "", "", ""},
		{"1.00.0", 1, 0, 0, "", "", "invalid version core: minor: leading zero is not allowed (at position 2)"},
		{"1.2.3", 1, 2, 3, "", "", ""},
		{"0.1.0", 0, 1, 0, "", "", ""},
		{"1.0.", 1, 0, 0, "", "", "invalid version core: patch: unexpected end of input (at position 4)"},
		{"1.0.0-beta.2", 1, 0, 0, "beta.2", "", ""},
		{"1.0.0-alpha+001", 1, 0, 0, "alpha", "001", ""},
		{"1.0.0+20130313144700", 1, 0, 0, "", "20130313144700", ""},
		{"1.0.0-beta+exp.sha.5114f85", 1, 0, 0, "beta", "exp.sha.5114f85", ""},
		{"1.0.0+21AF26D3----117B344092BD", 1, 0, 0, "", "21AF26D3----117B344092BD", ""},
	}

	for _, test := range tests {
		t.Run(test.version, func(t *testing.T) {
			v, err := semver.ParseVersion(test.version)
			if err != nil {
				if test.expectedErr == "" {
					t.Errorf("Unexpected error: %v", err)
				}
				if err.Error() != test.expectedErr {
					t.Errorf("Expected error %q, got %q", test.expectedErr, err)
				}
				return
			}

			if v.Major != test.major {
				t.Errorf("Expected major version %d, got %d", test.major, v.Major)
			}
			if v.Minor != test.minor {
				t.Errorf("Expected minor version %d, got %d", test.minor, v.Minor)
			}
			if v.Patch != test.patch {
				t.Errorf("Expected patch version %d, got %d", test.patch, v.Patch)
			}
			if v.PreRelease != test.preRelease {
				t.Errorf("Expected pre-release version %q, got %q", test.preRelease, v.PreRelease)
			}
			if v.Build != test.build {
				t.Errorf("Expected build version %q, got %q", test.build, v.Build)
			}
		})
	}
}

func TestBefore(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.0.0", "1.0.1", true},
		{"1.0.0", "2.0.0", true},
		{"1.0.0", "1.1.0", true},
		{"1.0.0-alpha", "1.0.0", true},
		{"1.0.0-alpha", "1.0.0-alpha.1", true},
		{"1.0.0-alpha.1", "1.0.0-alpha.beta", true},
		{"1.0.0-alpha.beta", "1.0.0-beta", true},
		{"1.0.0-beta", "1.0.0-beta.2", true},
		{"1.0.0-beta.2", "1.0.0-beta.11", true},
		{"1.0.0-beta.11", "1.0.0-rc.1", true},
		{"1.0.0-rc.1", "1.0.0", true},
		{"1.0.0+build1", "1.0.0+build2", false}, // Build metadata does not affect precedence
		{"1.0.0+build2", "1.0.0+build1", false}, // Build metadata does not affect precedence
		{"1.0.0", "1.0.0-alpha", false},
		{"2.0.0", "1.0.0", false},
		{"1.1.0", "1.0.0", false},
		{"1.0.1", "1.0.0", false},
		{"1.2.3", "1.3.0-rc.1", true},
		{"1.2.3-rc.0200", "1.2.3-rc.030", true},
	}

	for _, test := range tests {
		t.Run(test.v1+" < "+test.v2, func(t *testing.T) {
			v1, err1 := semver.ParseVersion(test.v1)
			if err1 != nil {
				t.Errorf("Error parsing version %q: %v", test.v1, err1)
				return
			}
			v2, err2 := semver.ParseVersion(test.v2)
			if err2 != nil {
				t.Errorf("Error parsing version %q: %v", test.v2, err2)
				return
			}

			result := v1.Before(v2)
			if result != test.expected {
				t.Errorf("Expected %q.Before(%q) to be %v, got %v", test.v1, test.v2, test.expected, result)
			}
		})
	}
}

func TestEquals(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.0.0", "1.0.0", true},
		{"1.0.0-alpha", "1.0.0-alpha", true},
		{"1.0.0+build1", "1.0.0+build1", true},       // Build metadata is ignored in equality
		{"1.0.0-alpha+001", "1.0.0-alpha+002", true}, // Build metadata is ignored in equality
		{"1.0.0", "1.0.1", false},
		{"1.0.0", "2.0.0", false},
		{"1.0.0", "1.1.0", false},
		{"1.0.0-alpha", "1.0.0-beta", false},
		{"1.0.0-alpha.1", "1.0.0-alpha.2", false},
	}

	for _, test := range tests {
		t.Run(test.v1+" == "+test.v2, func(t *testing.T) {
			v1, err1 := semver.ParseVersion(test.v1)
			if err1 != nil {
				t.Errorf("Error parsing version %q: %v", test.v1, err1)
				return
			}
			v2, err2 := semver.ParseVersion(test.v2)
			if err2 != nil {
				t.Errorf("Error parsing version %q: %v", test.v2, err2)
				return
			}

			result := v1.Equals(v2)
			if result != test.expected {
				t.Errorf("Expected %q.Equals(%q) to be %v, got %v", test.v1, test.v2, test.expected, result)
			}
		})
	}
}
