package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ── DisplayName ─────────────────────────────────────────────────────────────

func TestDisplayName_KnownAgents(t *testing.T) {
	cases := []struct {
		ag   string
		want string
	}{
		{"claude", "Claude Code"},
		{"codex", "Codex"},
		{"gemini", "Gemini CLI"},
		{"antigravity", "Antigravity"},
		{"openclaw", "OpenClaw"},
		{"opencode", "OpenCode"},
	}
	for _, tc := range cases {
		t.Run(tc.ag, func(t *testing.T) {
			got := DisplayName(tc.ag)
			if got != tc.want {
				t.Errorf("DisplayName(%q) = %q; want %q", tc.ag, got, tc.want)
			}
		})
	}
}

func TestDisplayName_UnknownAgentPassthrough(t *testing.T) {
	got := DisplayName("unknown-agent")
	if got != "unknown-agent" {
		t.Errorf("DisplayName(unknown) = %q; want passthrough %q", got, "unknown-agent")
	}
}

// ── FromDisplayName ─────────────────────────────────────────────────────────

func TestFromDisplayName_KnownDisplayNames(t *testing.T) {
	cases := []struct {
		display string
		want    string
	}{
		{"Claude Code", "claude"},
		{"Codex", "codex"},
		{"Gemini CLI", "gemini"},
		{"Antigravity", "antigravity"},
		{"OpenClaw", "openclaw"},
		{"OpenCode", "opencode"},
	}
	for _, tc := range cases {
		t.Run(tc.display, func(t *testing.T) {
			got := FromDisplayName(tc.display)
			if got != tc.want {
				t.Errorf("FromDisplayName(%q) = %q; want %q", tc.display, got, tc.want)
			}
		})
	}
}

func TestFromDisplayName_UnknownPassthrough(t *testing.T) {
	got := FromDisplayName("SomeOtherTool")
	if got != "SomeOtherTool" {
		t.Errorf("FromDisplayName(unknown) = %q; want passthrough", got)
	}
}

func TestDisplayName_FromDisplayName_Roundtrip(t *testing.T) {
	for _, ag := range All {
		display := DisplayName(ag)
		back := FromDisplayName(display)
		if back != ag {
			t.Errorf("roundtrip failed: %q → %q → %q", ag, display, back)
		}
	}
}

// ── SkillsDir ───────────────────────────────────────────────────────────────

func TestSkillsDir_KnownAgents_ReturnsNonEmpty(t *testing.T) {
	for _, ag := range All {
		got := SkillsDir(ag)
		if got == "" {
			t.Errorf("SkillsDir(%q) returned empty string", ag)
		}
	}
}

func TestSkillsDir_UnknownAgent_ReturnsEmpty(t *testing.T) {
	got := SkillsDir("nonexistent-agent")
	if got != "" {
		t.Errorf("SkillsDir(unknown) = %q; want empty", got)
	}
}

func TestSkillsDir_ClaudeContainsDotClaude(t *testing.T) {
	got := SkillsDir("claude")
	// Must live under ~/.claude/skills
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, ".claude", "skills")
	if got != want {
		t.Errorf("SkillsDir(claude) = %q; want %q", got, want)
	}
}

func TestSkillsDir_CodexContainsDotAgents(t *testing.T) {
	got := SkillsDir("codex")
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, ".agents", "skills")
	if got != want {
		t.Errorf("SkillsDir(codex) = %q; want %q", got, want)
	}
}

func TestSkillsDir_GeminiContainsDotGemini(t *testing.T) {
	got := SkillsDir("gemini")
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, ".gemini", "skills")
	if got != want {
		t.Errorf("SkillsDir(gemini) = %q; want %q", got, want)
	}
}

func TestSkillsDir_OpenCodeUnderDotConfig(t *testing.T) {
	got := SkillsDir("opencode")
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, ".config", "opencode", "skills")
	if got != want {
		t.Errorf("SkillsDir(opencode) = %q; want %q", got, want)
	}
}

func TestSkillsDir_AntigravityUnderGeminiConfig(t *testing.T) {
	got := SkillsDir("antigravity")
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, ".gemini", "config", "skills")
	if got != want {
		t.Errorf("SkillsDir(antigravity) = %q; want %q", got, want)
	}
}

// ── ConfigFile ──────────────────────────────────────────────────────────────

func TestConfigFile_KnownAgents_ReturnsNonEmpty(t *testing.T) {
	for _, ag := range All {
		got := ConfigFile(ag)
		if got == "" {
			t.Errorf("ConfigFile(%q) returned empty string", ag)
		}
	}
}

func TestConfigFile_ClaudeIsCLAUDEmd(t *testing.T) {
	got := ConfigFile("claude")
	if filepath.Base(got) != "CLAUDE.md" {
		t.Errorf("ConfigFile(claude) base = %q; want CLAUDE.md", filepath.Base(got))
	}
}

func TestConfigFile_CodexIsAGENTSmd(t *testing.T) {
	got := ConfigFile("codex")
	if filepath.Base(got) != "AGENTS.md" {
		t.Errorf("ConfigFile(codex) base = %q; want AGENTS.md", filepath.Base(got))
	}
}

func TestConfigFile_AntigravityIsAGENTSmd(t *testing.T) {
	got := ConfigFile("antigravity")
	if filepath.Base(got) != "AGENTS.md" {
		t.Errorf("ConfigFile(antigravity) base = %q; want AGENTS.md", filepath.Base(got))
	}
}

// ── ConfigDir ───────────────────────────────────────────────────────────────

func TestConfigDir_KnownAgents_ReturnsNonEmpty(t *testing.T) {
	for _, ag := range All {
		got := ConfigDir(ag)
		if got == "" {
			t.Errorf("ConfigDir(%q) returned empty string", ag)
		}
	}
}

func TestConfigDir_SharesRootWithSkillsDir(t *testing.T) {
	// ConfigDir and SkillsDir must share the same agent root directory.
	// For most agents ConfigDir is the direct parent of SkillsDir; for
	// openclaw they are siblings under ~/.openclaw, so we check one level up.
	for _, ag := range All {
		cfgDir := ConfigDir(ag)
		skillsDir := SkillsDir(ag)
		cfgRoot := filepath.Dir(cfgDir)
		skillsRoot := filepath.Dir(skillsDir)
		if cfgRoot != skillsRoot && !strings.HasPrefix(skillsDir, cfgDir) {
			t.Errorf("agent %q: ConfigDir %q and SkillsDir %q share no common root", ag, cfgDir, skillsDir)
		}
	}
}

// ── All list ─────────────────────────────────────────────────────────────────

func TestAll_ContainsExpectedAgents(t *testing.T) {
	expected := []string{"claude", "codex", "gemini", "antigravity", "openclaw", "opencode"}
	if len(All) != len(expected) {
		t.Errorf("All has %d agents; want %d", len(All), len(expected))
	}
	byName := map[string]bool{}
	for _, ag := range All {
		byName[ag] = true
	}
	for _, ag := range expected {
		if !byName[ag] {
			t.Errorf("All does not contain %q", ag)
		}
	}
}

// ── SkillCount ───────────────────────────────────────────────────────────────

func TestSkillCount_EmptyDir_ReturnsZero(t *testing.T) {
	dir := t.TempDir()
	// Override HOME so SkillsDir returns our temp dir
	oldHome := os.Getenv("HOME")
	defer func() { _ = os.Setenv("HOME", oldHome) }()

	// Place skills dir at expected location for claude
	skillsDir := filepath.Join(dir, ".claude", "skills")
	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	_ = os.Setenv("HOME", dir)

	count := SkillCount("claude")
	if count != 0 {
		t.Errorf("SkillCount = %d; want 0", count)
	}
}

func TestSkillCount_CountsNonHidden(t *testing.T) {
	dir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer func() { _ = os.Setenv("HOME", oldHome) }()

	skillsDir := filepath.Join(dir, ".claude", "skills")
	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Create 3 visible entries + 1 hidden
	for _, name := range []string{"skill-a", "skill-b", "skill-c"} {
		if err := os.MkdirAll(filepath.Join(skillsDir, name), 0o755); err != nil {
			t.Fatalf("mkdir skill: %v", err)
		}
	}
	if err := os.WriteFile(filepath.Join(skillsDir, ".gitkeep"), nil, 0o644); err != nil {
		t.Fatalf("write hidden: %v", err)
	}

	_ = os.Setenv("HOME", dir)
	count := SkillCount("claude")
	if count != 3 {
		t.Errorf("SkillCount = %d; want 3", count)
	}
}

func TestSkillCount_NonexistentDir_ReturnsZero(t *testing.T) {
	dir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer func() { _ = os.Setenv("HOME", oldHome) }()
	// Do NOT create the skills dir
	_ = os.Setenv("HOME", dir)

	count := SkillCount("claude")
	if count != 0 {
		t.Errorf("SkillCount = %d; want 0 for nonexistent dir", count)
	}
}

// ── BrokenSymlinkCount ───────────────────────────────────────────────────────

func TestBrokenSymlinkCount_NoBroken_ReturnsZero(t *testing.T) {
	dir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer func() { _ = os.Setenv("HOME", oldHome) }()

	skillsDir := filepath.Join(dir, ".claude", "skills")
	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Regular dir — not broken
	if err := os.MkdirAll(filepath.Join(skillsDir, "normal-skill"), 0o755); err != nil {
		t.Fatalf("mkdir skill: %v", err)
	}
	_ = os.Setenv("HOME", dir)

	count := BrokenSymlinkCount("claude")
	if count != 0 {
		t.Errorf("BrokenSymlinkCount = %d; want 0", count)
	}
}
