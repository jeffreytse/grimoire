package agent

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAntigravityHooks_ConfigureAndRemove(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hooks.json")
	seed := `{"my-linter":{"PostToolUse":[{"matcher":"run_command","hooks":[{"command":"lint.sh"}]}]}}`
	if err := os.WriteFile(path, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	if changed, err := ConfigureAntigravityHooks(path, "/opt/my tools/grimoire"); err != nil || !changed {
		t.Fatalf("configure: changed=%v err=%v", changed, err)
	}
	cfg := readJSON(t, path)
	if cfg["my-linter"] == nil {
		t.Errorf("other hook lost: %v", cfg)
	}
	g, ok := cfg["grimoire"].(map[string]any)
	if !ok || g["SessionStart"] == nil || g["PreInvocation"] == nil {
		t.Fatalf("grimoire entry = %v", cfg["grimoire"])
	}
	raw, _ := os.ReadFile(path)
	if !strings.Contains(string(raw), `'/opt/my tools/grimoire' hook antigravity-prompt`) {
		t.Errorf("binary path not quoted: %s", raw)
	}
	if changed, err := ConfigureAntigravityHooks(path, "/opt/my tools/grimoire"); err != nil || changed {
		t.Errorf("re-configure: changed=%v err=%v; want no change", changed, err)
	}
	if changed, err := RemoveAntigravityHooks(path); err != nil || !changed {
		t.Fatalf("remove: changed=%v err=%v", changed, err)
	}
	if cfg := readJSON(t, path); cfg["grimoire"] != nil || cfg["my-linter"] == nil {
		t.Errorf("after remove: %v", cfg)
	}
}

func TestOpenCodePlugin_ConfigureAndRemove(t *testing.T) {
	path := filepath.Join(t.TempDir(), "plugins", "grimoire.js")
	if changed, err := ConfigureOpenCodePlugin(path, "grimoire"); err != nil || !changed {
		t.Fatalf("configure: changed=%v err=%v", changed, err)
	}
	if changed, err := ConfigureOpenCodePlugin(path, "grimoire"); err != nil || changed {
		t.Errorf("re-configure: changed=%v err=%v; want no change", changed, err)
	}
	b, _ := os.ReadFile(path)
	for _, s := range []string{`const BIN = "grimoire"`, `"chat.message"`, `"experimental.chat.system.transform"`} {
		if !strings.Contains(string(b), s) {
			t.Errorf("plugin lacks %s", s)
		}
	}
	// The generated plugin must at least parse as an ES module.
	if node, err := exec.LookPath("node"); err == nil {
		mjs := filepath.Join(t.TempDir(), "grimoire.mjs")
		if err := os.WriteFile(mjs, b, 0o644); err != nil {
			t.Fatal(err)
		}
		if out, err := exec.Command(node, "--check", mjs).CombinedOutput(); err != nil {
			t.Errorf("plugin does not parse: %v: %s", err, out)
		}
	}
	if changed, err := RemoveOpenCodePlugin(path); err != nil || !changed {
		t.Fatalf("remove: changed=%v err=%v", changed, err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("plugin not deleted")
	}
}

// A user's own plugin file with the same name is never overwritten or deleted.
func TestOpenCodePlugin_LeavesForeignFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "grimoire.js")
	if err := os.WriteFile(path, []byte("export const Mine = async () => ({})\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ConfigureOpenCodePlugin(path, "grimoire"); err == nil {
		t.Error("want error for a file grimoire did not write")
	}
	if changed, err := RemoveOpenCodePlugin(path); err != nil || changed {
		t.Errorf("remove foreign: changed=%v err=%v", changed, err)
	}
	if b, _ := os.ReadFile(path); !strings.Contains(string(b), "Mine") {
		t.Error("foreign plugin modified")
	}
}
