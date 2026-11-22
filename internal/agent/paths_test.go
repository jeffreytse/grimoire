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
		{"copilot", "GitHub Copilot CLI"},
		{"pi", "Pi"},
		{"dsh", "DeepSeek Harness"},
		{"kimi", "Kimi Code"},
		{"omp", "OMP"},
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
		{"GitHub Copilot CLI", "copilot"},
		{"Pi", "pi"},
		{"DeepSeek Harness", "dsh"},
		{"Kimi Code", "kimi"},
		{"OMP", "omp"},
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

func TestSkillsDir_AntigravityUnderAntigravityCLI(t *testing.T) {
	got := SkillsDir("antigravity")
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, ".gemini", "antigravity-cli", "skills")
	if got != want {
		t.Errorf("SkillsDir(antigravity) = %q; want %q", got, want)
	}
}

// ── ConfigFile ──────────────────────────────────────────────────────────────

func TestConfigFile_KnownAgents_ReturnsNonEmpty(t *testing.T) {
	for _, ag := range All {
		if IsPluginAgent(ag) {
			continue // plugin bundles have no instructions file
		}
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
		if IsPluginAgent(ag) {
			continue // plugin bundles have no instructions file
		}
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
	expected := []string{"claude", "codex", "gemini", "antigravity", "openclaw", "opencode", "copilot", "pi", "dsh", "kimi", "omp", "grok", "muse", "minimax", "agent"}
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

// ── Copilot / Pi / DSH paths ────────────────────────────────────────────────

func TestNewAgentPaths_Defaults(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("COPILOT_HOME", "")
	t.Setenv("PI_CODING_AGENT_DIR", "")
	t.Setenv("KIMI_CODE_HOME", "")
	t.Setenv("OMP_PROFILE", "")
	cases := []struct {
		ag, skills, cfgFile, project string
	}{
		{"copilot", filepath.Join(home, ".copilot", "skills"), filepath.Join(home, ".copilot", "copilot-instructions.md"), filepath.Join("p", ".github", "skills")},
		{"pi", filepath.Join(home, ".pi", "agent", "skills"), filepath.Join(home, ".pi", "agent", "AGENTS.md"), filepath.Join("p", ".pi", "skills")},
		{"dsh", filepath.Join(home, ".agents", "skills"), filepath.Join(home, ".agents", "AGENTS.md"), filepath.Join("p", ".agents", "skills")},
		{"kimi", filepath.Join(home, ".kimi-code", "skills"), filepath.Join(home, ".kimi-code", "AGENTS.md"), filepath.Join("p", ".kimi-code", "skills")},
		{"omp", filepath.Join(home, ".omp", "agent", "skills"), filepath.Join(home, ".omp", "agent", "AGENTS.md"), filepath.Join("p", ".omp", "skills")},
	}
	for _, tc := range cases {
		t.Run(tc.ag, func(t *testing.T) {
			if got := SkillsDir(tc.ag); got != tc.skills {
				t.Errorf("SkillsDir = %q; want %q", got, tc.skills)
			}
			if got := ConfigFile(tc.ag); got != tc.cfgFile {
				t.Errorf("ConfigFile = %q; want %q", got, tc.cfgFile)
			}
			if got := ProjectSkillsDir(tc.ag, "p"); got != tc.project {
				t.Errorf("ProjectSkillsDir = %q; want %q", got, tc.project)
			}
			if !strings.HasPrefix(SkillsDir(tc.ag), ConfigDir(tc.ag)) {
				t.Errorf("SkillsDir %q not under ConfigDir %q", SkillsDir(tc.ag), ConfigDir(tc.ag))
			}
		})
	}
}

func TestNewAgentPaths_EnvOverrides(t *testing.T) {
	cp, piDir := t.TempDir(), t.TempDir()
	t.Setenv("COPILOT_HOME", cp)
	t.Setenv("PI_CODING_AGENT_DIR", piDir)
	if got, want := SkillsDir("copilot"), filepath.Join(cp, "skills"); got != want {
		t.Errorf("copilot SkillsDir = %q; want %q", got, want)
	}
	if got, want := ConfigFile("pi"), filepath.Join(piDir, "AGENTS.md"); got != want {
		t.Errorf("pi ConfigFile = %q; want %q", got, want)
	}
}

func TestOMPAgentDir_Profiles(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("PI_CODING_AGENT_DIR", "")
	t.Setenv("PI_PROFILE", "legacy")
	t.Setenv("OMP_PROFILE", "work")
	if got, want := SkillsDir("omp"), filepath.Join(home, ".omp", "profiles", "work", "agent", "skills"); got != want {
		t.Errorf("OMP_PROFILE: got %q; want %q", got, want)
	}
	// Explicitly-empty OMP_PROFILE selects the default profile, not PI_PROFILE.
	t.Setenv("OMP_PROFILE", "")
	if got, want := SkillsDir("omp"), filepath.Join(home, ".omp", "agent", "skills"); got != want {
		t.Errorf("empty OMP_PROFILE: got %q; want %q", got, want)
	}
	shared := t.TempDir()
	t.Setenv("PI_CODING_AGENT_DIR", shared)
	if got, want := SkillsDir("omp"), filepath.Join(shared, "skills"); got != want {
		t.Errorf("PI_CODING_AGENT_DIR: got %q; want %q", got, want)
	}
	kimi := t.TempDir()
	t.Setenv("KIMI_CODE_HOME", kimi)
	if got, want := ConfigFile("kimi"), filepath.Join(kimi, "AGENTS.md"); got != want {
		t.Errorf("KIMI_CODE_HOME: got %q; want %q", got, want)
	}
}

// ── Plugin agents ───────────────────────────────────────────────────────────

func TestPluginAgents_GlobalOnlyBundlePaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("MINIMAX_DATA_DIR", "")
	cases := map[string]string{
		"minimax": filepath.Join(home, ".minimax", "plugins", "grimoire", "skills"),
		"grok":    filepath.Join(home, ".grimoire", "agent-plugins", "grok", "grimoire", "skills"),
		"muse":    filepath.Join(home, ".grimoire", "agent-plugins", "muse", "grimoire", "skills"),
	}
	for ag, want := range cases {
		if !IsPluginAgent(ag) {
			t.Errorf("IsPluginAgent(%q) = false", ag)
		}
		if got := SkillsDir(ag); got != want {
			t.Errorf("SkillsDir(%q) = %q; want %q", ag, got, want)
		}
		if got := ProjectSkillsDir(ag, "p"); got != "" {
			t.Errorf("ProjectSkillsDir(%q) = %q; want empty (global-only)", ag, got)
		}
	}
	if IsPluginAgent("claude") {
		t.Error("IsPluginAgent(claude) = true")
	}
}

func TestFinalizePlugin_WritesManifestsAndRegisters(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("MINIMAX_DATA_DIR", "")
	bin := t.TempDir()
	for _, b := range []string{"grok", "muse"} {
		if err := os.WriteFile(filepath.Join(bin, b), []byte("#!/bin/sh\n"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", bin)

	for _, ag := range []string{"minimax", "grok", "muse"} {
		for _, n := range []string{"write-readme", "apply-dry-principle"} {
			d := filepath.Join(SkillsDir(ag), n)
			if err := os.MkdirAll(d, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(d, "SKILL.md"), []byte("---\nname: "+n+"\n---\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		var calls []string
		run := func(name string, args ...string) error {
			calls = append(calls, name+" "+strings.Join(args, " "))
			return nil
		}
		if err := FinalizePlugin(ag, run); err != nil {
			t.Fatalf("FinalizePlugin(%s): %v", ag, err)
		}
		bundle := PluginBundleDir(ag)
		manifest := map[string]string{
			"minimax": ".claude-plugin/plugin.json",
			"grok":    ".grok-plugin/plugin.json",
			"muse":    ".muse-plugin/plugin.json",
		}[ag]
		b, err := os.ReadFile(filepath.Join(bundle, manifest))
		if err != nil {
			t.Fatalf("%s manifest: %v", ag, err)
		}
		if !strings.Contains(string(b), `"name": "grimoire"`) {
			t.Errorf("%s manifest missing name: %s", ag, b)
		}
		switch ag {
		case "minimax":
			if len(calls) != 0 {
				t.Errorf("minimax should not run commands; got %v", calls)
			}
		case "grok":
			want := []string{
				"grok plugin uninstall grimoire --keep-data",
				"grok plugin install " + bundle + " --trust",
				"grok plugin enable grimoire",
			}
			if strings.Join(calls, "|") != strings.Join(want, "|") {
				t.Errorf("grok calls = %v; want %v", calls, want)
			}
		case "muse":
			if !strings.Contains(string(b), `"path": "skills/apply-dry-principle/SKILL.md"`) {
				t.Errorf("muse manifest missing skill list: %s", b)
			}
			if len(calls) != 1 || calls[0] != "muse plugins install "+bundle {
				t.Errorf("muse calls = %v", calls)
			}
		}
	}
}
