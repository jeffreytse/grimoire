package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeffreytse/grimoire/internal/agent"
)

// fakeClaude records claude invocations; registered controls `mcp get`.
type fakeClaude struct {
	calls      []string
	registered bool
}

func (f *fakeClaude) run(_, name string, args ...string) (string, error) {
	call := name + " " + strings.Join(args, " ")
	f.calls = append(f.calls, call)
	if strings.HasPrefix(call, "claude mcp add") {
		if f.registered {
			return "MCP server grimoire already exists in user config", errors.New("exit 1")
		}
		f.registered = true
	}
	return "", nil
}

// setupMCPSetupEnv isolates HOME/XDG, provides a skills package (so no clone
// happens), stubs PATH lookups, and resets flags.
func setupMCPSetupEnv(t *testing.T, claudeOnPath bool) (f *fakeClaude, home string) {
	t.Helper()
	home = t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	pkg := t.TempDir()
	t.Setenv("GRIMOIRE_SOURCE", pkg)
	if err := os.MkdirAll(filepath.Join(pkg, "skills", "x", "skills", "s"), 0o755); err != nil {
		t.Fatal(err)
	}

	f = &fakeClaude{}
	origRun, origLook := mcpRunCmd, mcpLookPath
	mcpRunCmd = f.run
	mcpLookPath = func(name string) (string, error) {
		if name == "claude" && !claudeOnPath {
			return "", errors.New("not found")
		}
		return "/usr/bin/" + name, nil
	}
	t.Cleanup(func() { mcpRunCmd, mcpLookPath = origRun, origLook })
	flagMCPSetupTarget, flagMCPSetupScope, flagMCPSetupNoConfigure = "claude", "user", false
	flagMCPSetupNoHook, flagMCPSetupNoAllow = false, false
	return f, home
}

func TestMCPInstall_AddsAndRemovesClaudePromptHook(t *testing.T) {
	_, home := setupMCPSetupEnv(t, true)
	settings := filepath.Join(home, ".claude", "settings.json")
	if err := runMCPInstall(nil, nil); err != nil {
		t.Fatal(err)
	}
	if err := runMCPInstall(nil, nil); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(settings)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(b), "grimoire hook prompt"); n != 1 {
		t.Fatalf("want hook once after two installs, got %d: %s", n, b)
	}
	if err := runMCPUninstall(nil, nil); err != nil {
		t.Fatal(err)
	}
	b, _ = os.ReadFile(settings)
	if strings.Contains(string(b), "hook prompt") {
		t.Errorf("hook not removed: %s", b)
	}
}

func TestMCPInstall_NoHookFlag(t *testing.T) {
	_, home := setupMCPSetupEnv(t, true)
	flagMCPSetupNoHook = true
	if err := runMCPInstall(nil, nil); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(home, ".claude", "settings.json")); strings.Contains(string(b), "hook prompt") {
		t.Errorf("hook written despite --no-hook: %s", b)
	}
}

func TestMCPInstall_RegistersAndWritesInstruction(t *testing.T) {
	f, home := setupMCPSetupEnv(t, true)
	if err := runMCPInstall(nil, nil); err != nil {
		t.Fatal(err)
	}
	want := "claude mcp add -s user grimoire -- grimoire mcp serve"
	if len(f.calls) != 1 || f.calls[0] != want {
		t.Fatalf("calls = %v; want %q", f.calls, want)
	}
	b, err := os.ReadFile(filepath.Join(home, ".claude", "CLAUDE.md"))
	if err != nil || !strings.Contains(string(b), agent.MCPTriggerLine) {
		t.Fatalf("CLAUDE.md missing MCP instruction: %v %q", err, b)
	}

	// Second run: "already exists" from claude is treated as success.
	if err := runMCPInstall(nil, nil); err != nil {
		t.Fatalf("re-run should succeed when already registered: %v", err)
	}
}

func TestMCPInstall_ProjectScope(t *testing.T) {
	f, _ := setupMCPSetupEnv(t, true)
	proj := t.TempDir()
	origDir := flagProjectDir
	flagProjectDir = proj
	t.Cleanup(func() { flagProjectDir = origDir })
	flagMCPSetupScope = "project"
	if err := runMCPInstall(nil, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(f.calls, "|"), "claude mcp add -s project grimoire") {
		t.Errorf("calls = %v", f.calls)
	}
	if _, err := os.Stat(filepath.Join(proj, "CLAUDE.md")); err != nil {
		t.Errorf("project CLAUDE.md not written: %v", err)
	}
}

func TestMCPInstall_NoClaudeCLI(t *testing.T) {
	setupMCPSetupEnv(t, false)
	if err := runMCPInstall(nil, nil); err == nil || !strings.Contains(err.Error(), "setup incomplete for: claude") {
		t.Fatalf("want incomplete-setup error, got %v", err)
	}
	err := installMCPAgent(findMCPAgent("claude"), "grimoire")
	if err == nil || !strings.Contains(err.Error(), "claude mcp add -s user grimoire -- grimoire mcp serve") {
		t.Fatalf("want error with manual command, got %v", err)
	}
}

// Each CLI-registered agent must run exactly the add command that was checked
// against the real CLI, and write its documented instruction file.
func TestMCPInstall_PerAgentCommandsAndInstructionFiles(t *testing.T) {
	cases := []struct {
		target, wantCall, instr string
	}{
		{"codex", "codex mcp add grimoire -- grimoire mcp serve", ".codex/AGENTS.md"},
		{"gemini", "gemini mcp add -s user grimoire grimoire mcp serve", ".gemini/GEMINI.md"},
		{"copilot", "copilot mcp add grimoire -- grimoire mcp serve", ".copilot/copilot-instructions.md"},
	}
	for _, tc := range cases {
		t.Run(tc.target, func(t *testing.T) {
			f, home := setupMCPSetupEnv(t, true)
			t.Setenv("CODEX_HOME", "")
			t.Setenv("COPILOT_HOME", "")
			flagMCPSetupTarget = tc.target
			if err := runMCPInstall(nil, nil); err != nil {
				t.Fatal(err)
			}
			if len(f.calls) != 1 || f.calls[0] != tc.wantCall {
				t.Errorf("calls = %v; want [%s]", f.calls, tc.wantCall)
			}
			b, err := os.ReadFile(filepath.Join(home, tc.instr))
			if err != nil || !strings.Contains(string(b), agent.MCPTriggerLine) {
				t.Errorf("%s missing instruction: %v", tc.instr, err)
			}
		})
	}
}

func TestMCPInstall_OpenCodeConfigFile(t *testing.T) {
	f, home := setupMCPSetupEnv(t, true)
	flagMCPSetupTarget = "opencode"
	cfgPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatal(err)
	}
	// Existing user settings must survive.
	if err := os.WriteFile(cfgPath, []byte(`{"theme":"dark","mcp":{"other":{"type":"local","command":["x"]}}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runMCPInstall(nil, nil); err != nil {
		t.Fatal(err)
	}
	if len(f.calls) != 0 {
		t.Errorf("opencode should not run a CLI; calls = %v", f.calls)
	}
	var cfg map[string]any
	b, _ := os.ReadFile(cfgPath)
	if err := json.Unmarshal(b, &cfg); err != nil {
		t.Fatal(err)
	}
	mcpCfg, ok := cfg["mcp"].(map[string]any)
	if !ok || cfg["theme"] != "dark" || mcpCfg["other"] == nil {
		t.Fatalf("existing settings lost: %s", b)
	}
	g, ok := mcpCfg["grimoire"].(map[string]any)
	if !ok {
		t.Fatalf("no grimoire entry: %s", b)
	}
	enabled, okEnabled := g["enabled"].(bool)
	cmdArgs, okCmd := g["command"].([]any)
	if g["type"] != "local" || !okEnabled || !enabled || !okCmd || len(cmdArgs) != 3 {
		t.Errorf("grimoire entry = %v", g)
	}
	if _, err := os.Stat(filepath.Join(home, ".config", "opencode", "AGENTS.md")); err != nil {
		t.Errorf("AGENTS.md not written: %v", err)
	}

	// Idempotent: second run leaves the file byte-identical.
	if changed, err := registerOpenCodeMCP("user", "grimoire"); err != nil || changed {
		t.Errorf("re-register changed=%v err=%v", changed, err)
	}

	// Uninstall removes only grimoire.
	if err := runMCPUninstall(nil, nil); err != nil {
		t.Fatal(err)
	}
	b, _ = os.ReadFile(cfgPath)
	if strings.Contains(string(b), "grimoire") || !strings.Contains(string(b), "other") {
		t.Errorf("uninstall should drop only grimoire: %s", b)
	}
}

func TestMCPInstall_OpenCodeRefusesJSONC(t *testing.T) {
	_, home := setupMCPSetupEnv(t, true)
	dir := filepath.Join(home, ".config", "opencode")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "opencode.jsonc"), []byte("// mine\n{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := registerOpenCodeMCP("user", "grimoire"); err == nil || !strings.Contains(err.Error(), "opencode.jsonc") {
		t.Errorf("want JSONC refusal, got %v", err)
	}
}

func TestMCPInstall_GeminiUntrustedFolderHint(t *testing.T) {
	setupMCPSetupEnv(t, true)
	mcpRunCmd = func(_, _ string, _ ...string) (string, error) {
		return "Error: Cannot modify settings in an untrusted workspace. To enable this, ...", errors.New("exit 1")
	}
	flagMCPSetupScope = "project"
	err := installMCPAgent(findMCPAgent("gemini"), "grimoire")
	if err == nil || !strings.Contains(err.Error(), "trust the folder") {
		t.Errorf("want trust hint, got %v", err)
	}
}

func TestResolveMCPAgents(t *testing.T) {
	setupMCPSetupEnv(t, true)
	flagMCPSetupTarget = "codex, gemini,codex"
	got, err := resolveMCPAgents()
	if err != nil || len(got) != 2 || got[0].key != "codex" || got[1].key != "gemini" {
		t.Errorf("list: got %v err %v", got, err)
	}
	// Explicit target without the scope is an error; "all" skips it instead.
	flagMCPSetupTarget, flagMCPSetupScope = "codex", "project"
	if _, err := resolveMCPAgents(); err == nil {
		t.Error("codex --scope project should error")
	}
	flagMCPSetupTarget = "all"
	got, err = resolveMCPAgents()
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range got {
		if !a.supports("project") {
			t.Errorf("all/project included %s", a.key)
		}
	}
	if len(got) != 6 {
		t.Errorf("all/project = %d agents; want claude, gemini, opencode, pi, kimi, omp", len(got))
	}
}

func TestMCPInstall_RejectsBadFlags(t *testing.T) {
	setupMCPSetupEnv(t, true)
	flagMCPSetupTarget = "cursor"
	if err := runMCPInstall(nil, nil); err == nil {
		t.Error("unsupported target should error")
	}
	flagMCPSetupTarget, flagMCPSetupScope = "claude", "global"
	if err := runMCPInstall(nil, nil); err == nil {
		t.Error("bad scope should error")
	}
}

func TestMCPUninstall(t *testing.T) {
	f, home := setupMCPSetupEnv(t, true)
	if err := runMCPInstall(nil, nil); err != nil {
		t.Fatal(err)
	}
	f.calls = nil
	if err := runMCPUninstall(nil, nil); err != nil {
		t.Fatal(err)
	}
	if len(f.calls) != 1 || f.calls[0] != "claude mcp remove -s user grimoire" {
		t.Errorf("calls = %v", f.calls)
	}
	b, _ := os.ReadFile(filepath.Join(home, ".claude", "CLAUDE.md"))
	if strings.Contains(string(b), agent.MCPTriggerLine) {
		t.Errorf("instruction not removed: %q", b)
	}
}

// File-registered agents write the documented "mcpServers" file and
// instruction file, touch nothing else, and are idempotent.
func TestMCPInstall_MCPServersFileAgents(t *testing.T) {
	cases := []struct {
		target, cfg, instr string
		wantType           bool
	}{
		{"pi", ".pi/agent/mcp.json", ".pi/agent/AGENTS.md", false},
		{"kimi", ".kimi-code/mcp.json", ".kimi-code/AGENTS.md", false},
		{"omp", ".omp/agent/mcp.json", ".omp/agent/AGENTS.md", true},
		{"antigravity", ".gemini/config/mcp_config.json", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.target, func(t *testing.T) {
			f, home := setupMCPSetupEnv(t, true)
			for _, k := range []string{"PI_CODING_AGENT_DIR", "KIMI_CODE_HOME", "OMP_PROFILE", "PI_PROFILE"} {
				t.Setenv(k, "")
			}
			cfgPath := filepath.Join(home, tc.cfg)
			if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(cfgPath, []byte(`{"mcpServers":{"other":{"command":"x"}},"keep":1}`), 0o644); err != nil {
				t.Fatal(err)
			}
			flagMCPSetupTarget = tc.target
			if err := runMCPInstall(nil, nil); err != nil {
				t.Fatal(err)
			}
			if len(f.calls) != 0 {
				t.Errorf("file agent ran a CLI: %v", f.calls)
			}
			var cfg struct {
				MCPServers map[string]map[string]any `json:"mcpServers"`
				Keep       int                       `json:"keep"`
			}
			b, _ := os.ReadFile(cfgPath)
			if err := json.Unmarshal(b, &cfg); err != nil {
				t.Fatal(err)
			}
			g := cfg.MCPServers["grimoire"]
			if cfg.Keep != 1 || cfg.MCPServers["other"] == nil || g["command"] != "grimoire" {
				t.Errorf("config = %s", b)
			}
			if _, hasType := g["type"]; hasType != tc.wantType {
				t.Errorf("type field present=%v; want %v", hasType, tc.wantType)
			}
			if tc.instr != "" {
				ib, err := os.ReadFile(filepath.Join(home, tc.instr))
				if err != nil || !strings.Contains(string(ib), agent.MCPTriggerLine) {
					t.Errorf("instruction missing in %s: %v", tc.instr, err)
				}
			}
			before := string(b)
			if err := runMCPInstall(nil, nil); err != nil {
				t.Fatal(err)
			}
			if b2, _ := os.ReadFile(cfgPath); string(b2) != before {
				t.Errorf("rerun changed config:\n%s\n%s", before, b2)
			}
			if err := runMCPUninstall(nil, nil); err != nil {
				t.Fatal(err)
			}
			b3, _ := os.ReadFile(cfgPath)
			if strings.Contains(string(b3), "grimoire") || !strings.Contains(string(b3), "other") {
				t.Errorf("uninstall: %s", b3)
			}
		})
	}
}

func TestMCPInstall_AddsPromptHookForCodexAndGemini(t *testing.T) {
	for _, tc := range []struct{ target, file, event string }{
		{"codex", filepath.Join(".codex", "hooks.json"), "UserPromptSubmit"},
		{"gemini", filepath.Join(".gemini", "settings.json"), "BeforeAgent"},
	} {
		t.Run(tc.target, func(t *testing.T) {
			_, home := setupMCPSetupEnv(t, true)
			t.Setenv("CODEX_HOME", "")
			flagMCPSetupTarget = tc.target
			if err := runMCPInstall(nil, nil); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(home, tc.file)
			b, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var cfg struct {
				Hooks map[string][]struct {
					Hooks []struct{ Command string } `json:"hooks"`
				} `json:"hooks"`
			}
			if err := json.Unmarshal(b, &cfg); err != nil {
				t.Fatal(err)
			}
			groups := cfg.Hooks[tc.event]
			if len(groups) != 1 || len(groups[0].Hooks) != 1 || !strings.HasSuffix(groups[0].Hooks[0].Command, "grimoire hook prompt") {
				t.Fatalf("%s hooks = %s", tc.event, b)
			}
			if err := runMCPUninstall(nil, nil); err != nil {
				t.Fatal(err)
			}
			if b, _ := os.ReadFile(path); strings.Contains(string(b), "hook prompt") {
				t.Errorf("hook not removed: %s", b)
			}
		})
	}
}

// Antigravity creates mcp_config.json as an empty file; registering must
// treat it as {} instead of failing to parse it.
func TestMCPInstall_EmptyMCPServersFile(t *testing.T) {
	_, home := setupMCPSetupEnv(t, true)
	path := filepath.Join(home, ".gemini", "config", "mcp_config.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	flagMCPSetupTarget = "antigravity"
	if err := runMCPInstall(nil, nil); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var cfg struct {
		MCPServers map[string]any `json:"mcpServers"`
	}
	if err := json.Unmarshal(b, &cfg); err != nil || cfg.MCPServers["grimoire"] == nil {
		t.Fatalf("grimoire not registered in %s: %s (err %v)", path, b, err)
	}
}

func TestMCPInstall_OpenCodePluginAndAntigravityHooks(t *testing.T) {
	for _, tc := range []struct{ target, file, marker string }{
		{"opencode", filepath.Join(".config", "opencode", "plugins", "grimoire.js"), "chat.message"},
		{"antigravity", filepath.Join(".gemini", "config", "hooks.json"), "hook antigravity-session"},
	} {
		t.Run(tc.target, func(t *testing.T) {
			_, home := setupMCPSetupEnv(t, true)
			flagMCPSetupTarget = tc.target
			if err := runMCPInstall(nil, nil); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(home, tc.file)
			if b, err := os.ReadFile(path); err != nil || !strings.Contains(string(b), tc.marker) {
				t.Fatalf("%s: err=%v content=%s", path, err, b)
			}
			if err := runMCPUninstall(nil, nil); err != nil {
				t.Fatal(err)
			}
			if b, _ := os.ReadFile(path); strings.Contains(string(b), tc.marker) {
				t.Errorf("hook not removed: %s", b)
			}
		})
	}
}

// Install pre-approves only grimoire's read-only skill tools; uninstall
// removes exactly those entries; --no-allow skips it.
func TestMCPInstall_PreApprovesReadOnlyTools(t *testing.T) {
	for _, tc := range []struct{ target, file, rule string }{
		{"claude", filepath.Join(".claude", "settings.json"), "mcp__grimoire__grimoire_get_skill"},
		{"antigravity", filepath.Join(".gemini", "antigravity-cli", "settings.json"), "mcp(grimoire/grimoire_get_skill)"},
	} {
		t.Run(tc.target, func(t *testing.T) {
			_, home := setupMCPSetupEnv(t, true)
			path := filepath.Join(home, tc.file)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(`{"permissions":{"allow":["command(git log)"]}}`), 0o644); err != nil {
				t.Fatal(err)
			}
			flagMCPSetupTarget = tc.target
			if err := runMCPInstall(nil, nil); err != nil {
				t.Fatal(err)
			}
			b, _ := os.ReadFile(path)
			if !strings.Contains(string(b), tc.rule) || !strings.Contains(string(b), "command(git log)") {
				t.Fatalf("allow list = %s", b)
			}
			if strings.Contains(string(b), "grimoire_config_set") || strings.Contains(string(b), "grimoire_package_add") {
				t.Errorf("mutating tool pre-approved: %s", b)
			}
			if err := runMCPUninstall(nil, nil); err != nil {
				t.Fatal(err)
			}
			b, _ = os.ReadFile(path)
			if strings.Contains(string(b), tc.rule) || !strings.Contains(string(b), "command(git log)") {
				t.Errorf("after uninstall = %s", b)
			}

			flagMCPSetupNoAllow = true
			if err := runMCPInstall(nil, nil); err != nil {
				t.Fatal(err)
			}
			// Match the quoted allow entry: the skill-used hook's matcher also
			// names grimoire_get_skill.
			if b, _ := os.ReadFile(path); strings.Contains(string(b), `"`+tc.rule+`"`) {
				t.Errorf("--no-allow still pre-approved: %s", b)
			}
		})
	}
}

// Claude Desktop: registered in its config file with an absolute command
// (desktop apps don't inherit the shell's PATH); rerun is a no-op, uninstall
// removes only grimoire, and --target all leaves it out (no CLI to detect).
func TestMCPInstall_ClaudeDesktop(t *testing.T) {
	f, _ := setupMCPSetupEnv(t, true)
	cfgPath := agent.ClaudeDesktopConfig()
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cfgPath, []byte(`{"mcpServers":{"other":{"command":"x"}},"keep":1}`), 0o644); err != nil {
		t.Fatal(err)
	}
	flagMCPSetupTarget = "claude-desktop"
	if err := runMCPInstall(nil, nil); err != nil {
		t.Fatal(err)
	}
	if len(f.calls) != 0 {
		t.Errorf("ran a CLI: %v", f.calls)
	}
	var cfg struct {
		MCPServers map[string]map[string]any `json:"mcpServers"`
		Keep       int                       `json:"keep"`
	}
	b, _ := os.ReadFile(cfgPath)
	if err := json.Unmarshal(b, &cfg); err != nil {
		t.Fatal(err)
	}
	cmd, ok := cfg.MCPServers["grimoire"]["command"].(string)
	if !ok || !filepath.IsAbs(cmd) || cfg.Keep != 1 || cfg.MCPServers["other"] == nil {
		t.Errorf("config = %s", b)
	}
	before := string(b)
	if err := runMCPInstall(nil, nil); err != nil {
		t.Fatal(err)
	}
	if b2, _ := os.ReadFile(cfgPath); string(b2) != before {
		t.Errorf("rerun changed config")
	}
	if err := runMCPUninstall(nil, nil); err != nil {
		t.Fatal(err)
	}
	if b3, _ := os.ReadFile(cfgPath); strings.Contains(string(b3), "grimoire") || !strings.Contains(string(b3), "other") {
		t.Errorf("uninstall: %s", b3)
	}

	flagMCPSetupTarget = "all"
	agents, err := resolveMCPAgents()
	if err == nil {
		for _, a := range agents {
			if a.key == "claude-desktop" {
				t.Error("--target all included Claude Desktop")
			}
		}
	}
}
