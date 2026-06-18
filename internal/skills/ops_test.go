package skills

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// makeSkillDir creates a minimal skill directory (with SKILL.md) in parent.
func makeSkillDir(t *testing.T, parent, name string) string {
	t.Helper()
	dir := filepath.Join(parent, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("makeSkillDir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# "+name), 0o644); err != nil {
		t.Fatalf("makeSkillDir write SKILL.md: %v", err)
	}
	return dir
}

// ── InstallSkill ────────────────────────────────────────────────────────────

func TestInstallSkill_Symlink_CreatesLink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks require elevated privileges on Windows")
	}
	src := makeSkillDir(t, t.TempDir(), "my-skill")
	dest := t.TempDir()

	installed, err := InstallSkill(src, dest, true)
	if err != nil {
		t.Fatalf("InstallSkill: %v", err)
	}
	if !installed {
		t.Fatal("expected installed=true for a new symlink")
	}

	link := filepath.Join(dest, "my-skill")
	target, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("Readlink: %v", err)
	}
	if target != src {
		t.Errorf("symlink target = %q; want %q", target, src)
	}
}

func TestInstallSkill_Symlink_IdempotentWhenCorrect(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks require elevated privileges on Windows")
	}
	src := makeSkillDir(t, t.TempDir(), "skill-a")
	dest := t.TempDir()

	// First install
	if _, err := InstallSkill(src, dest, true); err != nil {
		t.Fatalf("first install: %v", err)
	}
	// Second install — same src, same dest
	installed, err := InstallSkill(src, dest, true)
	if err != nil {
		t.Fatalf("second install: %v", err)
	}
	if installed {
		t.Fatal("expected installed=false when symlink already points to correct target")
	}
}

func TestInstallSkill_Symlink_ReplacesWrongTarget(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks require elevated privileges on Windows")
	}
	base := t.TempDir()
	oldSrc := makeSkillDir(t, base, "old-skill")
	newSrc := makeSkillDir(t, base, "new-skill")
	dest := t.TempDir()

	// Manually plant a symlink pointing at oldSrc, using the name newSrc would get
	link := filepath.Join(dest, "new-skill")
	if err := os.Symlink(oldSrc, link); err != nil {
		t.Fatalf("setup symlink: %v", err)
	}

	installed, err := InstallSkill(newSrc, dest, true)
	if err != nil {
		t.Fatalf("InstallSkill: %v", err)
	}
	if !installed {
		t.Fatal("expected installed=true when replacing wrong-target symlink")
	}

	got, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("Readlink after replace: %v", err)
	}
	if got != newSrc {
		t.Errorf("after replace: link target = %q; want %q", got, newSrc)
	}
}

func TestInstallSkill_Symlink_ReplacesBrokenSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks require elevated privileges on Windows")
	}
	base := t.TempDir()
	src := makeSkillDir(t, base, "real-skill")
	dest := t.TempDir()

	// Plant a broken symlink (pointing nowhere) under the same name
	link := filepath.Join(dest, "real-skill")
	if err := os.Symlink("/nonexistent/path/that/does/not/exist", link); err != nil {
		t.Fatalf("setup broken symlink: %v", err)
	}

	installed, err := InstallSkill(src, dest, true)
	if err != nil {
		t.Fatalf("InstallSkill: %v", err)
	}
	if !installed {
		t.Fatal("expected installed=true after replacing broken symlink")
	}

	got, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("Readlink: %v", err)
	}
	if got != src {
		t.Errorf("link target = %q; want %q", got, src)
	}
}

func TestInstallSkill_Copy_CopiesFiles(t *testing.T) {
	base := t.TempDir()
	src := makeSkillDir(t, base, "copy-skill")
	// Add an extra file inside the skill dir
	if err := os.WriteFile(filepath.Join(src, "README.md"), []byte("readme"), 0o644); err != nil {
		t.Fatalf("write README: %v", err)
	}
	dest := t.TempDir()

	installed, err := InstallSkill(src, dest, false)
	if err != nil {
		t.Fatalf("InstallSkill copy: %v", err)
	}
	if !installed {
		t.Fatal("expected installed=true for copy mode")
	}

	skillDir := filepath.Join(dest, "copy-skill")
	if _, err := os.Stat(filepath.Join(skillDir, "SKILL.md")); err != nil {
		t.Errorf("SKILL.md not found in copy dest: %v", err)
	}
	if _, err := os.Stat(filepath.Join(skillDir, "README.md")); err != nil {
		t.Errorf("README.md not found in copy dest: %v", err)
	}
}

func TestInstallSkill_Copy_CreatesDestDir(t *testing.T) {
	base := t.TempDir()
	src := makeSkillDir(t, base, "skill-x")
	// dest does not exist yet
	dest := filepath.Join(t.TempDir(), "nonexistent", "nested")

	_, err := InstallSkill(src, dest, false)
	if err != nil {
		t.Fatalf("InstallSkill copy into new dir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "skill-x", "SKILL.md")); err != nil {
		t.Errorf("SKILL.md not found after creating dest: %v", err)
	}
}

// ── UninstallSkill ──────────────────────────────────────────────────────────

func TestUninstallSkill_RemovesInstalledSkill(t *testing.T) {
	base := t.TempDir()
	src := makeSkillDir(t, base, "to-remove")
	dest := t.TempDir()

	if _, err := InstallSkill(src, dest, false); err != nil {
		t.Fatalf("setup install: %v", err)
	}

	removed, err := UninstallSkill("to-remove", dest)
	if err != nil {
		t.Fatalf("UninstallSkill: %v", err)
	}
	if !removed {
		t.Fatal("expected removed=true")
	}
	if _, err := os.Stat(filepath.Join(dest, "to-remove")); !os.IsNotExist(err) {
		t.Error("skill dir still exists after uninstall")
	}
}

func TestUninstallSkill_NopWhenNotInstalled(t *testing.T) {
	dest := t.TempDir()
	removed, err := UninstallSkill("ghost-skill", dest)
	if err != nil {
		t.Fatalf("UninstallSkill: %v", err)
	}
	if removed {
		t.Fatal("expected removed=false for non-existent skill")
	}
}

func TestUninstallSkill_RemovesSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks require elevated privileges on Windows")
	}
	base := t.TempDir()
	src := makeSkillDir(t, base, "sym-skill")
	dest := t.TempDir()

	if _, err := InstallSkill(src, dest, true); err != nil {
		t.Fatalf("setup symlink install: %v", err)
	}

	removed, err := UninstallSkill("sym-skill", dest)
	if err != nil {
		t.Fatalf("UninstallSkill: %v", err)
	}
	if !removed {
		t.Fatal("expected removed=true")
	}
	if _, err := os.Lstat(filepath.Join(dest, "sym-skill")); !os.IsNotExist(err) {
		t.Error("symlink still exists after uninstall")
	}
}

// ── CleanBrokenSymlinks ─────────────────────────────────────────────────────

func TestCleanBrokenSymlinks_RemovesBrokenOnly(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks require elevated privileges on Windows")
	}
	dir := t.TempDir()

	// Good symlink — points at a real file
	realFile := filepath.Join(t.TempDir(), "real.md")
	if err := os.WriteFile(realFile, []byte("ok"), 0o644); err != nil {
		t.Fatalf("write real file: %v", err)
	}
	if err := os.Symlink(realFile, filepath.Join(dir, "good")); err != nil {
		t.Fatalf("good symlink: %v", err)
	}

	// Broken symlink
	if err := os.Symlink("/no/such/path", filepath.Join(dir, "broken")); err != nil {
		t.Fatalf("broken symlink: %v", err)
	}

	// Hidden file — should not be touched
	if err := os.WriteFile(filepath.Join(dir, ".hidden"), []byte("x"), 0o644); err != nil {
		t.Fatalf("hidden file: %v", err)
	}

	n, err := CleanBrokenSymlinks(dir)
	if err != nil {
		t.Fatalf("CleanBrokenSymlinks: %v", err)
	}
	if n != 1 {
		t.Errorf("removed count = %d; want 1", n)
	}

	if _, err := os.Lstat(filepath.Join(dir, "good")); err != nil {
		t.Error("good symlink was incorrectly removed")
	}
	if _, err := os.Lstat(filepath.Join(dir, "broken")); !os.IsNotExist(err) {
		t.Error("broken symlink was not removed")
	}
	if _, err := os.Lstat(filepath.Join(dir, ".hidden")); err != nil {
		t.Error("hidden file was removed but should not be")
	}
}

func TestCleanBrokenSymlinks_NonexistentDirReturnsZero(t *testing.T) {
	n, err := CleanBrokenSymlinks(filepath.Join(t.TempDir(), "doesnotexist"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}
