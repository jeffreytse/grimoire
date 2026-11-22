package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectMCPSetup(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("CODEX_HOME", "")
	t.Setenv("COPILOT_HOME", "")
	proj := t.TempDir()

	for _, ag := range []string{"claude", "codex", "copilot", "opencode", "antigravity"} {
		if s := DetectMCPSetup(ag, "user", proj); s.Active() {
			t.Errorf("%s: active before install: %+v", ag, s)
		}
	}

	// claude: instruction + hook
	if _, err := ConfigureMCPInstruction(MCPInstructionFile("claude", "user", proj)); err != nil {
		t.Fatal(err)
	}
	h, _ := PromptHookFor("claude")
	if _, err := ConfigurePromptHook(PromptHookFile("claude", "user", proj), h, PromptHookCommand("grimoire")); err != nil {
		t.Fatal(err)
	}
	if s := DetectMCPSetup("claude", "user", proj); !s.Instruction || !s.Hook || !s.HookSupported {
		t.Errorf("claude: %+v", s)
	}

	// copilot: instruction only, no hook support
	if _, err := ConfigureMCPInstruction(MCPInstructionFile("copilot", "user", proj)); err != nil {
		t.Fatal(err)
	}
	if s := DetectMCPSetup("copilot", "user", proj); !s.Instruction || s.Hook || s.HookSupported {
		t.Errorf("copilot: %+v", s)
	}

	// antigravity: the session hook stands in for the instruction file
	if _, err := ConfigureAntigravityHooks(AntigravityHooksFile("user", proj), "grimoire"); err != nil {
		t.Fatal(err)
	}
	if s := DetectMCPSetup("antigravity", "user", proj); !s.Instruction || !s.Hook {
		t.Errorf("antigravity: %+v", s)
	}

	// opencode: plugin only (instruction missing) is still MCP mode
	if _, err := ConfigureOpenCodePlugin(OpenCodePluginFile("user", proj), "grimoire"); err != nil {
		t.Fatal(err)
	}
	if s := DetectMCPSetup("opencode", "user", proj); !s.Active() || s.Instruction || !s.Hook {
		t.Errorf("opencode: %+v", s)
	}
}

func TestLinkedSkillCount(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	pkg := filepath.Join(home, ".grimoire", "packages", "p")
	skill := filepath.Join(pkg, "skills", "x")
	dir := SkillsDir("claude")
	for _, d := range []string{skill, dir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(skill, filepath.Join(dir, "x")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(dir, "context7-mcp")); err != nil {
		t.Fatal(err)
	}
	if n := LinkedSkillCount("claude", pkg); n != 1 {
		t.Errorf("LinkedSkillCount = %d; want 1", n)
	}
}
