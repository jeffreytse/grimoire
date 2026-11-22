package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readJSON(t *testing.T, path string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	m := map[string]any{}
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	return m
}

func promptCommands(t *testing.T, path string) []string {
	t.Helper()
	groups, err := hookGroups(readJSON(t, path), path)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, g := range groups {
		for _, h := range groupHooks(g) {
			out = append(out, hookCommand(h))
		}
	}
	return out
}

func TestConfigureClaudeHook_KeepsOtherSettingsAndIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	seed := `{"model":"opus","hooks":{"UserPromptSubmit":[{"hooks":[{"type":"command","command":"caveman-hook"}]}],"Stop":[{"hooks":[{"type":"command","command":"x"}]}]}}`
	if err := os.WriteFile(path, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := ClaudeHookCommand("grimoire")
	changed, err := ConfigureClaudeHook(path, cmd)
	if err != nil || !changed {
		t.Fatalf("first configure: changed=%v err=%v", changed, err)
	}
	if got := promptCommands(t, path); len(got) != 2 || got[0] != "caveman-hook" || got[1] != "grimoire hook prompt" {
		t.Fatalf("UserPromptSubmit commands = %v", got)
	}
	cfg := readJSON(t, path)
	if cfg["model"] != "opus" || cfg["hooks"].(map[string]any)["Stop"] == nil {
		t.Errorf("other settings lost: %v", cfg)
	}

	if changed, err := ConfigureClaudeHook(path, cmd); err != nil || changed {
		t.Fatalf("re-configure: changed=%v err=%v; want no change", changed, err)
	}
}

func TestConfigureClaudeHook_ReplacesOldBinaryPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if _, err := ConfigureClaudeHook(path, ClaudeHookCommand("/opt/old/grimoire")); err != nil {
		t.Fatal(err)
	}
	if _, err := ConfigureClaudeHook(path, ClaudeHookCommand("grimoire")); err != nil {
		t.Fatal(err)
	}
	if got := promptCommands(t, path); len(got) != 1 || got[0] != "grimoire hook prompt" {
		t.Fatalf("commands = %v; want only the new one", got)
	}
}

func TestRemoveClaudeHook(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if changed, err := RemoveClaudeHook(path); err != nil || changed {
		t.Fatalf("missing file: changed=%v err=%v", changed, err)
	}
	seed := `{"hooks":{"UserPromptSubmit":[{"hooks":[{"type":"command","command":"caveman-hook"},{"type":"command","command":"'/a b/grimoire' hook prompt"}]}]}}`
	if err := os.WriteFile(path, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	if changed, err := RemoveClaudeHook(path); err != nil || !changed {
		t.Fatalf("remove: changed=%v err=%v", changed, err)
	}
	if got := promptCommands(t, path); len(got) != 1 || got[0] != "caveman-hook" {
		t.Fatalf("commands = %v; want caveman-hook only", got)
	}

	// Removing the last hook drops the empty "hooks" object.
	only := filepath.Join(t.TempDir(), "settings.json")
	if _, err := ConfigureClaudeHook(only, ClaudeHookCommand("grimoire")); err != nil {
		t.Fatal(err)
	}
	if _, err := RemoveClaudeHook(only); err != nil {
		t.Fatal(err)
	}
	if cfg := readJSON(t, only); cfg["hooks"] != nil {
		t.Errorf("empty hooks left behind: %v", cfg)
	}
}

func TestConfigureClaudeHook_RefusesMalformed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	for _, bad := range []string{`{"hooks":[]}`, `{"hooks":{"UserPromptSubmit":{}}}`, `not json`} {
		if err := os.WriteFile(path, []byte(bad), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := ConfigureClaudeHook(path, "grimoire hook prompt"); err == nil {
			t.Errorf("%s: want error", bad)
		}
		if b, _ := os.ReadFile(path); string(b) != bad {
			t.Errorf("%s: file was rewritten to %s", bad, b)
		}
	}
}

func TestClaudeHookCommandAndMatch(t *testing.T) {
	if got := ClaudeHookCommand("/Users/me/My Tools/grimoire"); got != "'/Users/me/My Tools/grimoire' hook prompt" {
		t.Errorf("quoted command = %q", got)
	}
	for cmd, want := range map[string]bool{
		"grimoire hook prompt":                    true,
		"/usr/local/bin/grimoire hook prompt":     true,
		"'/a b/grimoire' hook prompt":             true,
		"other hook prompt":                       false,
		"grimoire hook prompt --min-score 5 junk": false,
		strings.Repeat(" ", 3):                    false,
	} {
		if got := isGrimoireHook(cmd); got != want {
			t.Errorf("isGrimoireHook(%q) = %v; want %v", cmd, got, want)
		}
	}
}
