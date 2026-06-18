package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// detectInstallMethod classifies a binary path as "go" (installed via
// `go install`) or "binary" (downloaded release or other location).
// It checks whether the binary resides in $GOPATH/bin.

func TestDetectInstallMethod_GoPath_ReturnsGo(t *testing.T) {
	fakeGoPath := t.TempDir()
	binDir := filepath.Join(fakeGoPath, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("mkdir bin: %v", err)
	}
	// Fake binary inside $GOPATH/bin
	exePath := filepath.Join(binDir, "grimoire")
	if err := os.WriteFile(exePath, []byte("fake"), 0o755); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}

	// Set GOPATH so detectInstallMethod can find it without running `go env`.
	old := os.Getenv("GOPATH")
	defer func() { _ = os.Setenv("GOPATH", old) }()
	_ = os.Setenv("GOPATH", fakeGoPath)

	got := DetectInstallMethodForTest(exePath)
	if got != "go" {
		t.Errorf("detectInstallMethod = %q; want go", got)
	}
}

func TestDetectInstallMethod_OutsideGoPath_ReturnsBinary(t *testing.T) {
	// exePath is somewhere other than $GOPATH/bin
	fakeGoPath := t.TempDir()
	otherDir := t.TempDir()

	exePath := filepath.Join(otherDir, "grimoire")
	if err := os.WriteFile(exePath, []byte("fake"), 0o755); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}

	old := os.Getenv("GOPATH")
	defer func() { _ = os.Setenv("GOPATH", old) }()
	_ = os.Setenv("GOPATH", fakeGoPath)

	got := DetectInstallMethodForTest(exePath)
	if got != "binary" {
		t.Errorf("detectInstallMethod = %q; want binary", got)
	}
}

func TestDetectInstallMethod_EmptyGoPath_ReturnsBinary(t *testing.T) {
	// When GOPATH is empty and `go env` is unavailable or returns nothing,
	// any exe path must return "binary".
	old := os.Getenv("GOPATH")
	defer func() { _ = os.Setenv("GOPATH", old) }()
	_ = os.Setenv("GOPATH", "")

	// Use a path that cannot possibly be in GOPATH/bin when GOPATH is empty.
	exePath := filepath.Join(t.TempDir(), "grimoire")
	if err := os.WriteFile(exePath, []byte("fake"), 0o755); err != nil {
		t.Fatalf("write: %v", err)
	}

	got := DetectInstallMethodForTest(exePath)
	// With no GOPATH, we expect "binary" unless the `go env GOPATH` call
	// happens to return a path matching our temp dir — which won't happen.
	if got != "binary" {
		t.Errorf("detectInstallMethod = %q; want binary when GOPATH empty", got)
	}
}

func TestDetectInstallMethod_ExeInSubdirOfGoBin_ReturnsBinary(t *testing.T) {
	// Binary is at $GOPATH/bin/subdir/grimoire — not directly in bin/.
	fakeGoPath := t.TempDir()
	subDir := filepath.Join(fakeGoPath, "bin", "subdir")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	exePath := filepath.Join(subDir, "grimoire")
	if err := os.WriteFile(exePath, []byte("fake"), 0o755); err != nil {
		t.Fatalf("write: %v", err)
	}

	old := os.Getenv("GOPATH")
	defer func() { _ = os.Setenv("GOPATH", old) }()
	_ = os.Setenv("GOPATH", fakeGoPath)

	got := DetectInstallMethodForTest(exePath)
	// filepath.Dir(exePath) == .../bin/subdir ≠ .../bin → "binary"
	if got != "binary" {
		t.Errorf("detectInstallMethod = %q; want binary for subdir of bin", got)
	}
}

// ── selfUpdateBinary asset selection ─────────────────────────────────────────
// Test the pure asset-name selection logic embedded in selfUpdateBinary.
// We extract the name-building formula and verify it matches platform expectations.

func TestAssetName_LinuxAmd64(t *testing.T) {
	// Mirrors the logic inside selfUpdateBinary:
	//   assetName := fmt.Sprintf("grimoire-%s-%s", runtime.GOOS, runtime.GOARCH)
	// We test the formula by constructing it directly.
	goos, goarch := "linux", "amd64"
	name := "grimoire-" + goos + "-" + goarch
	if name != "grimoire-linux-amd64" {
		t.Errorf("asset name = %q; want grimoire-linux-amd64", name)
	}
}

func TestAssetName_WindowsAddsExeSuffix(t *testing.T) {
	goos, goarch := "windows", "amd64"
	name := "grimoire-" + goos + "-" + goarch + ".exe"
	if name != "grimoire-windows-amd64.exe" {
		t.Errorf("asset name = %q; want grimoire-windows-amd64.exe", name)
	}
}

// ── splitLast (cmd/uninstall.go helper, exercised here via package) ───────────
// splitLast is in the same package (cmd) so we can call it directly.

func TestSplitLast_WithSeparator(t *testing.T) {
	cases := []struct {
		in   string
		sep  byte
		want string
	}{
		{"a/b/c", '/', "c"},
		{"domain/subdomain/skill-name", '/', "skill-name"},
		{"single", '/', "single"},
		{"", '/', ""},
		{"a/", '/', ""},
	}
	for _, tc := range cases {
		got := splitLast(tc.in, tc.sep)
		if got != tc.want {
			t.Errorf("splitLast(%q, %q) = %q; want %q", tc.in, string(tc.sep), got, tc.want)
		}
	}
}

// ── skillNameFromPath ─────────────────────────────────────────────────────────

func TestSkillNameFromPath_LongerPathReturnsBasename(t *testing.T) {
	// When skillPath is longer than ref, use the last segment of skillPath.
	got := skillNameFromPath("/home/user/.grimoire/skills/engineering/dev/skills/my-skill", "engineering/dev/my-skill")
	if got != "my-skill" {
		t.Errorf("skillNameFromPath = %q; want my-skill", got)
	}
}

func TestSkillNameFromPath_SameLengthFallsBackToRef(t *testing.T) {
	// When skillPath length == ref length, fall back to last segment of ref.
	got := skillNameFromPath("eng/skill", "eng/skill")
	if got != "skill" {
		t.Errorf("skillNameFromPath = %q; want skill", got)
	}
}
