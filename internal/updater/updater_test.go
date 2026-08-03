package updater

import (
	"strings"
	"testing"
)

func TestCompareSemver(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.0.1", -1},
		{"1.1.0", "1.0.9", 1},
		{"2.0.0", "1.9.9", 1},

		// Numeric, not lexicographic: this is the comparison that decides
		// whether an update is offered, and "0.9.0" > "0.10.0" as strings.
		{"0.10.0", "0.9.0", 1},
		{"0.9.0", "0.10.0", -1},
		{"0.10.0", "0.10.0", 0},

		// Tag prefixes and partial versions.
		{"v1.2.3", "1.2.3", 0},
		{"1.2", "1.2.0", 0},
		{"1", "1.0.0", 0},
	}

	for _, tt := range tests {
		t.Run(tt.a+" vs "+tt.b, func(t *testing.T) {
			if got := compareSemver(tt.a, tt.b); got != tt.want {
				t.Fatalf("compareSemver(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// Prerelease and build metadata are stripped before comparing, so a release
// candidate compares equal to its release and does not trigger an update.
func TestCompareSemverIgnoresPrereleaseAndBuild(t *testing.T) {
	if got := compareSemver("1.2.3-rc1", "1.2.3"); got != 0 {
		t.Fatalf("compareSemver(1.2.3-rc1, 1.2.3) = %d, want 0", got)
	}
	if got := compareSemver("1.2.3+build9", "1.2.3"); got != 0 {
		t.Fatalf("compareSemver(1.2.3+build9, 1.2.3) = %d, want 0", got)
	}
}

func TestParseSemverNonNumericIsZero(t *testing.T) {
	if got := parseSemver("not.a.version"); got != [3]int{0, 0, 0} {
		t.Fatalf("parseSemver(not.a.version) = %v, want [0 0 0]", got)
	}
	if got := parseSemver(""); got != [3]int{0, 0, 0} {
		t.Fatalf("parseSemver(\"\") = %v, want [0 0 0]", got)
	}
}

// The asset name has to match the file names published by the release
// workflow, or Check reports "no asset in release" on every platform.
func TestPlatformAssetName(t *testing.T) {
	name := platformAssetName()

	if !strings.HasPrefix(name, "ftm-") {
		t.Fatalf("platformAssetName() = %q, want an ftm- prefix", name)
	}

	parts := strings.Split(strings.TrimPrefix(name, "ftm-"), "-")
	if len(parts) != 2 {
		t.Fatalf("platformAssetName() = %q, want ftm-<os>-<arch>", name)
	}
	for _, part := range parts {
		if part == "" {
			t.Fatalf("platformAssetName() = %q has an empty segment", name)
		}
	}
}

func TestOSAndArchAliases(t *testing.T) {
	if got := osAlias(); got == "" {
		t.Fatal("osAlias() is empty")
	}
	if got := archAlias(); got == "" {
		t.Fatal("archAlias() is empty")
	}
}

func TestListAssetNames(t *testing.T) {
	assets := []Asset{{Name: "ftm-linux-x64"}, {Name: "ftm-macos-arm64"}}

	if got := listAssetNames(assets); got != "ftm-linux-x64, ftm-macos-arm64" {
		t.Fatalf("listAssetNames() = %q", got)
	}
	if got := listAssetNames(nil); got != "" {
		t.Fatalf("listAssetNames(nil) = %q, want empty", got)
	}
}
