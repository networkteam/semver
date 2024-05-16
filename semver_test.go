package semver_test

import (
	"testing"

	"github.com/networkteam/semver"
)

// TestBefore is a table-driven test for the Before method of the Version struct.
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
	}

	for _, test := range tests {
		v1, err1 := semver.ParseVersion(test.v1)
		if err1 != nil {
			t.Errorf("Error parsing version %s: %v", test.v1, err1)
			continue
		}
		v2, err2 := semver.ParseVersion(test.v2)
		if err2 != nil {
			t.Errorf("Error parsing version %s: %v", test.v2, err2)
			continue
		}

		result := v1.Before(v2)
		if result != test.expected {
			t.Errorf("Expected %v.Before(%v) to be %v, got %v", test.v1, test.v2, test.expected, result)
		}
	}
}

// TestEquals is a table-driven test for the Equals method of the Version struct.
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
		v1, err1 := semver.ParseVersion(test.v1)
		if err1 != nil {
			t.Errorf("Error parsing version %s: %v", test.v1, err1)
			continue
		}
		v2, err2 := semver.ParseVersion(test.v2)
		if err2 != nil {
			t.Errorf("Error parsing version %s: %v", test.v2, err2)
			continue
		}

		result := v1.Equals(v2)
		if result != test.expected {
			t.Errorf("Expected %v.Equals(%v) to be %v, got %v", test.v1, test.v2, test.expected, result)
		}
	}
}
