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
	groups, err := hookGroups(readJSON(t, path), path, "UserPromptSubmit")
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

func TestConfigurePromptHook_KeepsOtherSettingsAndIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	seed := `{"model":"opus","hooks":{"UserPromptSubmit":[{"hooks":[{"type":"command","command":"caveman-hook"}]}],"Stop":[{"hooks":[{"type":"command","command":"x"}]}]}}`
	if err := os.WriteFile(path, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := PromptHookCommand("grimoire")
	changed, err := ConfigurePromptHook(path, promptHooks["claude"], cmd)
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

	if changed, err := ConfigurePromptHook(path, promptHooks["claude"], cmd); err != nil || changed {
		t.Fatalf("re-configure: changed=%v err=%v; want no change", changed, err)
	}
}

func TestConfigurePromptHook_ReplacesOldBinaryPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if _, err := ConfigurePromptHook(path, promptHooks["claude"], PromptHookCommand("/opt/old/grimoire")); err != nil {
		t.Fatal(err)
	}
	if _, err := ConfigurePromptHook(path, promptHooks["claude"], PromptHookCommand("grimoire")); err != nil {
		t.Fatal(err)
	}
	if got := promptCommands(t, path); len(got) != 1 || got[0] != "grimoire hook prompt" {
		t.Fatalf("commands = %v; want only the new one", got)
	}
}

func TestRemovePromptHook(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if changed, err := RemovePromptHook(path, promptHooks["claude"]); err != nil || changed {
		t.Fatalf("missing file: changed=%v err=%v", changed, err)
	}
	seed := `{"hooks":{"UserPromptSubmit":[{"hooks":[{"type":"command","command":"caveman-hook"},{"type":"command","command":"'/a b/grimoire' hook prompt"}]}]}}`
	if err := os.WriteFile(path, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	if changed, err := RemovePromptHook(path, promptHooks["claude"]); err != nil || !changed {
		t.Fatalf("remove: changed=%v err=%v", changed, err)
	}
	if got := promptCommands(t, path); len(got) != 1 || got[0] != "caveman-hook" {
		t.Fatalf("commands = %v; want caveman-hook only", got)
	}

	// Removing the last hook drops the empty "hooks" object.
	only := filepath.Join(t.TempDir(), "settings.json")
	if _, err := ConfigurePromptHook(only, promptHooks["claude"], PromptHookCommand("grimoire")); err != nil {
		t.Fatal(err)
	}
	if _, err := RemovePromptHook(only, promptHooks["claude"]); err != nil {
		t.Fatal(err)
	}
	if cfg := readJSON(t, only); cfg["hooks"] != nil {
		t.Errorf("empty hooks left behind: %v", cfg)
	}
}

func TestConfigurePromptHook_RefusesMalformed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	for _, bad := range []string{`{"hooks":[]}`, `{"hooks":{"UserPromptSubmit":{}}}`, `not json`} {
		if err := os.WriteFile(path, []byte(bad), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := ConfigurePromptHook(path, promptHooks["claude"], "grimoire hook prompt"); err == nil {
			t.Errorf("%s: want error", bad)
		}
		if b, _ := os.ReadFile(path); string(b) != bad {
			t.Errorf("%s: file was rewritten to %s", bad, b)
		}
	}
}

func TestPromptHookCommandAndMatch(t *testing.T) {
	if got := PromptHookCommand("/Users/me/My Tools/grimoire"); got != "'/Users/me/My Tools/grimoire' hook prompt" {
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

func TestConfigurePromptHook_GeminiUsesBeforeAgentAndMillis(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	seed := `{"security":{"auth":{"selectedType":"oauth-personal"}},"hooks":{"BeforeAgent":[{"matcher":"*","hooks":[{"name":"memory","type":"command","command":"node mem.js"}]}]}}`
	if err := os.WriteFile(path, []byte(seed), 0o644); err != nil {
		t.Fatal(err)
	}
	h, _ := PromptHookFor("gemini")
	if changed, err := ConfigurePromptHook(path, h, PromptHookCommand("grimoire")); err != nil || !changed {
		t.Fatalf("configure: changed=%v err=%v", changed, err)
	}
	cfg := readJSON(t, path)
	if cfg["security"] == nil {
		t.Errorf("other settings lost: %v", cfg)
	}
	groups, err := hookGroups(cfg, path, "BeforeAgent")
	if err != nil || len(groups) != 2 {
		t.Fatalf("BeforeAgent groups = %v, err %v", groups, err)
	}
	entry, ok := groupHooks(groups[1])[0].(map[string]any)
	if !ok {
		t.Fatalf("gemini hook entry is not an object: %v", groups[1])
	}
	if entry["command"] != "grimoire hook prompt" || entry["name"] != "grimoire" || entry["timeout"] != float64(10000) {
		t.Errorf("gemini hook entry = %v", entry)
	}
	if _, ok := cfg["hooks"].(map[string]any)["UserPromptSubmit"]; ok {
		t.Error("gemini hook written under UserPromptSubmit")
	}

	if changed, err := RemovePromptHook(path, h); err != nil || !changed {
		t.Fatalf("remove: changed=%v err=%v", changed, err)
	}
	if groups, _ := hookGroups(readJSON(t, path), path, "BeforeAgent"); len(groups) != 1 {
		t.Errorf("want only the memory hook left, got %v", groups)
	}
}

func TestPromptHookFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CODEX_HOME", filepath.Join(home, "cx"))
	proj := filepath.Join(home, "proj")
	for _, tc := range []struct{ ag, scope, want string }{
		{"claude", "user", filepath.Join(home, ".claude", "settings.json")},
		{"claude", "project", filepath.Join(proj, ".claude", "settings.json")},
		{"codex", "user", filepath.Join(home, "cx", "hooks.json")},
		{"codex", "project", filepath.Join(proj, ".codex", "hooks.json")},
		{"gemini", "user", filepath.Join(home, ".gemini", "settings.json")},
		{"gemini", "project", filepath.Join(proj, ".gemini", "settings.json")},
		{"copilot", "user", ""},
	} {
		if got := PromptHookFile(tc.ag, tc.scope, proj); got != tc.want {
			t.Errorf("PromptHookFile(%s, %s) = %q; want %q", tc.ag, tc.scope, got, tc.want)
		}
	}
}
