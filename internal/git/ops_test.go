package git

import (
	"os"
	"path/filepath"
	"testing"
)

// ── compareSemver ─────────────────────────────────────────────────────────────

func TestCompareSemver_Equal(t *testing.T) {
	cases := []struct{ a, b string }{
		{"1.0.0", "1.0.0"},
		{"v2.3.1", "v2.3.1"},
		{"0.0.0", "0.0.0"},
	}
	for _, tc := range cases {
		got := compareSemver(tc.a, tc.b)
		if got != 0 {
			t.Errorf("compareSemver(%q, %q) = %d; want 0", tc.a, tc.b, got)
		}
	}
}

func TestCompareSemver_AGreaterThanB(t *testing.T) {
	cases := []struct{ a, b string }{
		{"2.0.0", "1.9.9"},
		{"1.1.0", "1.0.9"},
		{"1.0.1", "1.0.0"},
		{"v10.0.0", "v9.9.9"},
	}
	for _, tc := range cases {
		got := compareSemver(tc.a, tc.b)
		if got != 1 {
			t.Errorf("compareSemver(%q, %q) = %d; want 1", tc.a, tc.b, got)
		}
	}
}

func TestCompareSemver_ALessThanB(t *testing.T) {
	cases := []struct{ a, b string }{
		{"1.0.0", "2.0.0"},
		{"1.0.9", "1.1.0"},
		{"0.9.9", "1.0.0"},
		{"v1.2.3", "v1.2.4"},
	}
	for _, tc := range cases {
		got := compareSemver(tc.a, tc.b)
		if got != -1 {
			t.Errorf("compareSemver(%q, %q) = %d; want -1", tc.a, tc.b, got)
		}
	}
}

func TestCompareSemver_VPrefixStripped(t *testing.T) {
	if compareSemver("v1.2.3", "1.2.3") != 0 {
		t.Error("v-prefix should be stripped before comparison")
	}
	if compareSemver("v2.0.0", "v1.0.0") != 1 {
		t.Error("v2.0.0 should be > v1.0.0")
	}
}

func TestCompareSemver_DifferentPartCounts(t *testing.T) {
	// "1.0" vs "1.0.0" — missing part treated as 0
	if compareSemver("1.0", "1.0.0") != 0 {
		t.Error("1.0 should equal 1.0.0")
	}
	if compareSemver("1.1", "1.0.0") != 1 {
		t.Error("1.1 should be > 1.0.0")
	}
}

func TestCompareSemver_ZeroVersions(t *testing.T) {
	if compareSemver("0.0.0", "0.0.1") != -1 {
		t.Error("0.0.0 should be < 0.0.1")
	}
}

// ── readVersion ───────────────────────────────────────────────────────────────

func TestReadVersion_FileAbsent_ReturnsUnknown(t *testing.T) {
	dir := t.TempDir()
	got := readVersion(dir)
	if got != "unknown" {
		t.Errorf("readVersion with no file = %q; want %q", got, "unknown")
	}
}

func TestReadVersion_FilePresent_ReturnsTrimmed(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "VERSION"), []byte("  1.2.3\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	got := readVersion(dir)
	if got != "1.2.3" {
		t.Errorf("readVersion = %q; want %q", got, "1.2.3")
	}
}

func TestReadVersion_EmptyFile_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "VERSION"), []byte(""), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	got := readVersion(dir)
	if got != "" {
		t.Errorf("readVersion(empty file) = %q; want empty string", got)
	}
}
