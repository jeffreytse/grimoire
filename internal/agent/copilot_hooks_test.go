package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopilotHooks_ConfigureAndRemove(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hooks", "grimoire.json")
	if changed, err := ConfigureCopilotHooks(path, "/opt/bin/grimoire"); err != nil || !changed {
		t.Fatalf("configure: changed=%v err=%v", changed, err)
	}
	if changed, err := ConfigureCopilotHooks(path, "/opt/bin/grimoire"); err != nil || changed {
		t.Errorf("re-configure: changed=%v err=%v; want no change", changed, err)
	}
	b, _ := os.ReadFile(path)
	for _, s := range []string{`"version": 1`, `"SessionStart"`, `"UserPromptSubmit"`,
		`/opt/bin/grimoire hook session --format flat`, `/opt/bin/grimoire hook prompt --format flat`, `"timeoutSec": 10`} {
		if !strings.Contains(string(b), s) {
			t.Errorf("hook file lacks %s:\n%s", s, b)
		}
	}
	if !HasCopilotHooks(path) {
		t.Error("HasCopilotHooks = false after configure")
	}
	// A new binary path replaces grimoire's own file.
	if changed, err := ConfigureCopilotHooks(path, "grimoire"); err != nil || !changed {
		t.Errorf("new binary: changed=%v err=%v", changed, err)
	}
	if changed, err := RemoveCopilotHooks(path); err != nil || !changed {
		t.Fatalf("remove: changed=%v err=%v", changed, err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("hook file not deleted")
	}
}

// A hook file with someone else's hooks is never overwritten or deleted.
func TestCopilotHooks_LeavesForeignFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "grimoire.json")
	mine := `{"version":1,"hooks":{"preToolUse":[{"type":"command","bash":"./check.sh"}]}}`
	if err := os.WriteFile(path, []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ConfigureCopilotHooks(path, "grimoire"); err == nil {
		t.Error("want error for a file grimoire did not write")
	}
	if changed, err := RemoveCopilotHooks(path); err != nil || changed {
		t.Errorf("remove foreign: changed=%v err=%v", changed, err)
	}
	if b, _ := os.ReadFile(path); string(b) != mine {
		t.Error("foreign hook file modified")
	}
}

// The --format flag doesn't stop grimoire from recognizing its own hooks.
func TestIsGrimoireHook_WithFormatFlag(t *testing.T) {
	for _, c := range []string{"grimoire hook prompt --format flat", "/x/grimoire hook session --format text", "grimoire hook prompt"} {
		if !isGrimoireHook(c) {
			t.Errorf("%q not recognized", c)
		}
	}
	if isGrimoireHook("other hook prompt --format flat") {
		t.Error("non-grimoire binary recognized")
	}
}
