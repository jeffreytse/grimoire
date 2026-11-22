package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Claude Code hook that runs `grimoire hook prompt` on every user prompt.
// The instruction line in CLAUDE.md is advisory: the model decides whether a
// task "needs" a skill and often skips the lookup. A UserPromptSubmit hook
// runs on every prompt regardless, so matching skills are always put in front
// of the model.

// ClaudeHookEvent is the Claude Code hook event grimoire registers under.
const ClaudeHookEvent = "UserPromptSubmit"

// claudeHookArgs is what follows the grimoire binary in the hook command.
const claudeHookArgs = "hook prompt"

// claudeHookTimeout (seconds) bounds the hook; a slow search must never stall
// the user's prompt.
const claudeHookTimeout = 10

// ClaudeSettingsFile returns Claude Code's settings.json for scope.
func ClaudeSettingsFile(scope, projectDir string) string {
	if scope == "project" {
		return filepath.Join(projectDir, ".claude", "settings.json")
	}
	return filepath.Join(mustHome(), ".claude", "settings.json")
}

// ClaudeHookCommand is the hook command line for the grimoire binary bin.
func ClaudeHookCommand(bin string) string {
	if strings.ContainsAny(bin, " \t'\"") {
		bin = "'" + strings.ReplaceAll(bin, "'", `'\''`) + "'"
	}
	return bin + " " + claudeHookArgs
}

// isGrimoireHook reports whether command is a grimoire prompt hook, whatever
// binary path it was registered with.
func isGrimoireHook(command string) bool {
	c := strings.TrimSpace(command)
	if !strings.HasSuffix(c, " "+claudeHookArgs) {
		return false
	}
	bin := strings.Trim(strings.TrimSpace(strings.TrimSuffix(c, claudeHookArgs)), `'"`)
	return filepath.Base(bin) == "grimoire"
}

// ConfigureClaudeHook makes the settings file at path run command on every
// prompt, replacing any older grimoire hook (e.g. a different binary path).
// Other settings and hooks are kept. Reports whether the file changed.
func ConfigureClaudeHook(path, command string) (bool, error) {
	cfg, err := readSettings(path)
	if err != nil {
		return false, err
	}
	groups, err := hookGroups(cfg, path)
	if err != nil {
		return false, err
	}

	// Already present exactly once with this command: nothing to do.
	count, exact := 0, false
	for _, g := range groups {
		for _, h := range groupHooks(g) {
			if isGrimoireHook(hookCommand(h)) {
				count++
				exact = exact || hookCommand(h) == command
			}
		}
	}
	if count == 1 && exact {
		return false, nil
	}

	groups = withoutGrimoireHooks(groups)
	groups = append(groups, map[string]any{
		"hooks": []any{map[string]any{
			"type":    "command",
			"command": command,
			"timeout": claudeHookTimeout,
		}},
	})
	setHookGroups(cfg, groups)
	return true, writeSettings(path, cfg)
}

// RemoveClaudeHook drops grimoire's prompt hook from the settings file at
// path, leaving every other hook. Reports whether the file changed.
func RemoveClaudeHook(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}
	cfg, err := readSettings(path)
	if err != nil {
		return false, err
	}
	groups, err := hookGroups(cfg, path)
	if err != nil {
		return false, err
	}
	kept := withoutGrimoireHooks(groups)
	if len(kept) == len(groups) && countHooks(kept) == countHooks(groups) {
		return false, nil
	}
	setHookGroups(cfg, kept)
	return true, writeSettings(path, cfg)
}

func readSettings(path string) (map[string]any, error) {
	cfg := map[string]any{}
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(string(raw)) == "" {
		return cfg, nil
	}
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return cfg, nil
}

func writeSettings(path string, cfg map[string]any) error {
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

// hookGroups returns cfg.hooks.UserPromptSubmit. A value of the wrong shape is
// an error: grimoire does not overwrite settings it cannot read.
func hookGroups(cfg map[string]any, path string) ([]any, error) {
	hv, ok := cfg["hooks"]
	if !ok {
		return nil, nil
	}
	hooks, ok := hv.(map[string]any)
	if !ok {
		return nil, fmt.Errorf(`%s: "hooks" is not an object — fix it by hand`, path)
	}
	ev, ok := hooks[ClaudeHookEvent]
	if !ok {
		return nil, nil
	}
	groups, ok := ev.([]any)
	if !ok {
		return nil, fmt.Errorf(`%s: "hooks.%s" is not an array — fix it by hand`, path, ClaudeHookEvent)
	}
	return groups, nil
}

func setHookGroups(cfg map[string]any, groups []any) {
	hooks, ok := cfg["hooks"].(map[string]any)
	if !ok {
		hooks = map[string]any{}
	}
	if len(groups) == 0 {
		delete(hooks, ClaudeHookEvent)
	} else {
		hooks[ClaudeHookEvent] = groups
	}
	if len(hooks) == 0 {
		delete(cfg, "hooks")
		return
	}
	cfg["hooks"] = hooks
}

func groupHooks(g any) []any {
	m, ok := g.(map[string]any)
	if !ok {
		return nil
	}
	if hs, ok := m["hooks"].([]any); ok {
		return hs
	}
	return nil
}

func hookCommand(h any) string {
	m, ok := h.(map[string]any)
	if !ok {
		return ""
	}
	if c, ok := m["command"].(string); ok {
		return c
	}
	return ""
}

func countHooks(groups []any) int {
	n := 0
	for _, g := range groups {
		n += len(groupHooks(g))
	}
	return n
}

// withoutGrimoireHooks removes grimoire hooks from each group and drops groups
// left empty by that. Groups of unknown shape are kept untouched.
func withoutGrimoireHooks(groups []any) []any {
	out := make([]any, 0, len(groups))
	for _, g := range groups {
		m, ok := g.(map[string]any)
		if !ok {
			out = append(out, g)
			continue
		}
		hs := groupHooks(g)
		kept := make([]any, 0, len(hs))
		for _, h := range hs {
			if !isGrimoireHook(hookCommand(h)) {
				kept = append(kept, h)
			}
		}
		if len(kept) == len(hs) {
			out = append(out, g)
			continue
		}
		if len(kept) == 0 {
			continue
		}
		cp := make(map[string]any, len(m))
		for k, v := range m {
			cp[k] = v
		}
		cp["hooks"] = kept
		out = append(out, cp)
	}
	return out
}
