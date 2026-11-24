package agent

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestOMPExtension_ConfigureAndRemove(t *testing.T) {
	path := filepath.Join(t.TempDir(), "extensions", "grimoire.ts")
	if changed, err := ConfigureOMPExtension(path, "grimoire"); err != nil || !changed {
		t.Fatalf("configure: changed=%v err=%v", changed, err)
	}
	if changed, err := ConfigureOMPExtension(path, "grimoire"); err != nil || changed {
		t.Errorf("re-configure: changed=%v err=%v", changed, err)
	}
	b, _ := os.ReadFile(path)
	for _, s := range []string{ompExtensionMarker, `from "@oh-my-pi/pi-coding-agent"`, `event.reason !== "prompt"`, `systemPrompt: [...event.systemPrompt, ...extra]`} {
		if !strings.Contains(string(b), s) {
			t.Errorf("extension lacks %s", s)
		}
	}
	if node, err := exec.LookPath("node"); err == nil {
		mts := filepath.Join(t.TempDir(), "grimoire.mts")
		if err := os.WriteFile(mts, b, 0o644); err != nil {
			t.Fatal(err)
		}
		if out, err := exec.Command(node, "--experimental-strip-types", "--check", mts).CombinedOutput(); err != nil {
			t.Errorf("extension does not parse: %v: %s", err, out)
		}
	}
	if changed, err := RemoveOMPExtension(path); err != nil || !changed {
		t.Fatalf("remove: changed=%v err=%v", changed, err)
	}
}

func TestOMPExtensionFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("PI_CODING_AGENT_DIR", "")
	t.Setenv("PI_CONFIG_DIR", "")
	proj := t.TempDir()
	if got, want := OMPExtensionFile("user", proj), filepath.Join(home, ".omp", "agent", "extensions", "grimoire.ts"); got != want {
		t.Errorf("user = %s, want %s", got, want)
	}
	if got, want := OMPExtensionFile("project", proj), filepath.Join(proj, ".omp", "extensions", "grimoire.ts"); got != want {
		t.Errorf("project = %s, want %s", got, want)
	}
}

const kimiUserConfig = `# my settings
default_model = "kimi-k2"

[[hooks]]
event = "PreToolUse"
matcher = "Bash"
command = "node ~/.kimi-code/hooks/check.mjs"
timeout = 5

[mcp]
enabled = true
`

func TestKimiHook_ConfigureKeepsUserConfigAndRemoves(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(kimiUserConfig), 0o644); err != nil {
		t.Fatal(err)
	}
	if changed, err := ConfigureKimiHook(path, "/opt/bin/grimoire"); err != nil || !changed {
		t.Fatalf("configure: changed=%v err=%v", changed, err)
	}
	if changed, err := ConfigureKimiHook(path, "/opt/bin/grimoire"); err != nil || changed {
		t.Errorf("re-configure: changed=%v err=%v", changed, err)
	}
	b, _ := os.ReadFile(path)
	if !strings.HasPrefix(string(b), kimiUserConfig) {
		t.Errorf("user config changed:\n%s", b)
	}
	var cfg struct {
		Hooks []map[string]any `toml:"hooks"`
	}
	if err := toml.Unmarshal(b, &cfg); err != nil {
		t.Fatalf("config no longer parses: %v\n%s", err, b)
	}
	if len(cfg.Hooks) != 2 {
		t.Fatalf("want the user's hook plus grimoire's, got %v", cfg.Hooks)
	}
	g := cfg.Hooks[1]
	if g["event"] != "UserPromptSubmit" || g["command"] != "/opt/bin/grimoire hook prompt --format text" || len(g) != 3 {
		t.Errorf("grimoire hook = %v (Kimi allows only event, matcher, command, timeout)", g)
	}
	if !HasKimiHook(path) {
		t.Error("HasKimiHook = false")
	}
	// A new binary path replaces the entry rather than adding a second one.
	if _, err := ConfigureKimiHook(path, "grimoire"); err != nil {
		t.Fatal(err)
	}
	b, _ = os.ReadFile(path)
	if strings.Count(string(b), kimiHookMarker) != 1 {
		t.Errorf("want one grimoire entry:\n%s", b)
	}
	if changed, err := RemoveKimiHook(path); err != nil || !changed {
		t.Fatalf("remove: changed=%v err=%v", changed, err)
	}
	b, _ = os.ReadFile(path)
	if string(b) != kimiUserConfig {
		t.Errorf("remove should restore the user's file exactly:\n%q", b)
	}
}

func TestKimiHook_NewFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "kimi", "config.toml")
	if _, err := ConfigureKimiHook(path, "grimoire"); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(path)
	var probe map[string]any
	if err := toml.Unmarshal(b, &probe); err != nil {
		t.Fatalf("new config does not parse: %v\n%s", err, b)
	}
}
