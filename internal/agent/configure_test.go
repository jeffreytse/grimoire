package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupAgentDir creates $tmp/.claude/ and returns (home, cfgDir, cfgFile).
func setupAgentDir(t *testing.T) (home, cfgDir, cfgFile string) {
	t.Helper()
	home = t.TempDir()
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() { _ = os.Setenv("HOME", oldHome) })
	_ = os.Setenv("HOME", home)

	cfgDir = filepath.Join(home, ".claude")
	cfgFile = filepath.Join(cfgDir, "CLAUDE.md")
	return home, cfgDir, cfgFile
}

func writeCfg(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func readCfg(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	return string(data)
}

// ── IsConfigured ─────────────────────────────────────────────────────────────

func TestIsConfigured_FileAbsent_ReturnsFalse(t *testing.T) {
	_, _, _ = setupAgentDir(t)
	// Do NOT create any file
	if IsConfigured("claude") {
		t.Error("IsConfigured = true; want false when file absent")
	}
}

func TestIsConfigured_FileExistsWithoutTrigger_ReturnsFalse(t *testing.T) {
	_, _, cfgFile := setupAgentDir(t)
	writeCfg(t, cfgFile, "# My config\nsome other content\n")
	if IsConfigured("claude") {
		t.Error("IsConfigured = true; want false without trigger")
	}
}

func TestIsConfigured_FileExistsWithTrigger_ReturnsTrue(t *testing.T) {
	_, _, cfgFile := setupAgentDir(t)
	writeCfg(t, cfgFile, "# Config\n"+triggerLine+"\n")
	if !IsConfigured("claude") {
		t.Error("IsConfigured = false; want true when trigger present")
	}
}

// ── ConfigureAgentMD ─────────────────────────────────────────────────────────

func TestConfigureAgentMD_AgentDirAbsent_DoesNothing(t *testing.T) {
	home, _, _ := setupAgentDir(t)
	// cfgDir (~/.claude) does NOT exist
	if err := ConfigureAgentMD("claude"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// No file should be created
	cfgFile := filepath.Join(home, ".claude", "CLAUDE.md")
	if _, err := os.Stat(cfgFile); err == nil {
		t.Error("config file created; want nothing created when agent dir absent")
	}
}

func TestConfigureAgentMD_AppendsSection_WhenNotPresent(t *testing.T) {
	_, cfgDir, cfgFile := setupAgentDir(t)
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	writeCfg(t, cfgFile, "# Existing content\n")

	if err := ConfigureAgentMD("claude"); err != nil {
		t.Fatalf("ConfigureAgentMD: %v", err)
	}

	got := readCfg(t, cfgFile)
	if !strings.Contains(got, sectionHeader) {
		t.Errorf("section header %q not found in output:\n%s", sectionHeader, got)
	}
	if !strings.Contains(got, triggerLine) {
		t.Errorf("trigger line not found in output:\n%s", got)
	}
	if !strings.Contains(got, "# Existing content") {
		t.Error("existing content was overwritten")
	}
}

func TestConfigureAgentMD_CreatesFileIfAbsent(t *testing.T) {
	_, cfgDir, cfgFile := setupAgentDir(t)
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// cfgFile does NOT exist

	if err := ConfigureAgentMD("claude"); err != nil {
		t.Fatalf("ConfigureAgentMD: %v", err)
	}

	got := readCfg(t, cfgFile)
	if !strings.Contains(got, triggerLine) {
		t.Errorf("trigger line not found; got:\n%s", got)
	}
}

func TestConfigureAgentMD_Idempotent(t *testing.T) {
	_, cfgDir, cfgFile := setupAgentDir(t)
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err := ConfigureAgentMD("claude"); err != nil {
		t.Fatalf("first call: %v", err)
	}
	first := readCfg(t, cfgFile)

	if err := ConfigureAgentMD("claude"); err != nil {
		t.Fatalf("second call: %v", err)
	}
	second := readCfg(t, cfgFile)

	if first != second {
		t.Errorf("second call modified file:\nbefore: %q\nafter:  %q", first, second)
	}
}

// ── RemoveAgentMDConfig ──────────────────────────────────────────────────────

func TestRemoveAgentMDConfig_NotConfigured_DoesNothing(t *testing.T) {
	_, _, cfgFile := setupAgentDir(t)
	original := "# My rules\ndo something\n"
	writeCfg(t, cfgFile, original)

	if err := RemoveAgentMDConfig("claude"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := readCfg(t, cfgFile)
	if got != original {
		t.Errorf("file modified; got %q; want %q", got, original)
	}
}

func TestRemoveAgentMDConfig_RemovesTriggerLines(t *testing.T) {
	_, _, cfgFile := setupAgentDir(t)
	writeCfg(t, cfgFile, "# Before\n\n"+sectionHeader+"\n"+triggerLine+"\n")

	if err := RemoveAgentMDConfig("claude"); err != nil {
		t.Fatalf("RemoveAgentMDConfig: %v", err)
	}

	got := readCfg(t, cfgFile)
	if strings.Contains(got, sectionHeader) {
		t.Errorf("section header still present after removal:\n%s", got)
	}
	if strings.Contains(got, triggerLine) {
		t.Errorf("trigger line still present after removal:\n%s", got)
	}
	if !strings.Contains(got, "# Before") {
		t.Error("pre-existing content was removed")
	}
}

func TestRemoveAgentMDConfig_EmptyAfterRemoval_WritesEmpty(t *testing.T) {
	_, _, cfgFile := setupAgentDir(t)
	writeCfg(t, cfgFile, "\n"+sectionHeader+"\n"+triggerLine+"\n")

	if err := RemoveAgentMDConfig("claude"); err != nil {
		t.Fatalf("RemoveAgentMDConfig: %v", err)
	}

	got := readCfg(t, cfgFile)
	if strings.TrimSpace(got) != "" {
		t.Errorf("file not empty after removing all content; got %q", got)
	}
}

func TestRemoveAgentMDConfig_FileAbsent_DoesNothing(t *testing.T) {
	setupAgentDir(t)
	// No file created — should not error
	if err := RemoveAgentMDConfig("claude"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemoveAgentMDConfig_RoundtripWithConfigure(t *testing.T) {
	_, cfgDir, cfgFile := setupAgentDir(t)
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	writeCfg(t, cfgFile, "# Preamble\n")
	original := readCfg(t, cfgFile)

	if err := ConfigureAgentMD("claude"); err != nil {
		t.Fatalf("configure: %v", err)
	}
	if err := RemoveAgentMDConfig("claude"); err != nil {
		t.Fatalf("remove: %v", err)
	}

	got := readCfg(t, cfgFile)
	if got != original {
		t.Errorf("file differs after configure+remove roundtrip:\ngot:  %q\nwant: %q", got, original)
	}
}
