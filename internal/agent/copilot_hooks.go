package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

// ── GitHub Copilot CLI ───────────────────────────────────────────────────────
//
// Copilot CLI loads every *.json in $COPILOT_HOME/hooks (user) and
// .github/hooks (repository, once the folder is trusted):
//
//	{"version": 1, "hooks": {"<event>": [{"type": "command", "command": "…", "timeoutSec": 10}]}}
//
// Under the Claude Code event names, hooks get Claude Code's input fields
// (session_id, prompt), which grimoire's hooks already read. Copilot adds
// only a flat {"additionalContext"} from stdout to the model's context and
// ignores hookSpecificOutput, so the commands run with --format flat
// (checked with Copilot CLI 1.0.114: SessionStart and UserPromptSubmit
// context both reach the model in that shape, and not in Claude Code's).
//
// grimoire owns the whole file grimoire.json; other hook files are untouched.

const copilotHookTimeout = 10 // seconds

// CopilotHooksFile returns grimoire's Copilot hook file for scope.
func CopilotHooksFile(scope, projectDir string) string {
	if scope == "project" {
		return filepath.Join(projectDir, ".github", "hooks", "grimoire.json")
	}
	return filepath.Join(copilotHome(mustHome()), "hooks", "grimoire.json")
}

func copilotHooks(bin string) map[string]any {
	entry := func(args string) []any {
		return []any{map[string]any{
			"type":       "command",
			"command":    shellQuote(bin) + " " + args + " --format flat",
			"timeoutSec": copilotHookTimeout,
		}}
	}
	return map[string]any{
		"version": 1,
		"hooks": map[string]any{
			"SessionStart":     entry(sessionHookArgs),
			"UserPromptSubmit": entry(promptHookArgs),
		},
	}
}

// ConfigureCopilotHooks writes grimoire's hook file at path. A file there that
// holds anything but grimoire's hooks is left alone (error). Reports whether
// the file changed.
func ConfigureCopilotHooks(path, bin string) (bool, error) {
	want := copilotHooks(bin)
	cur, err := readSettings(path)
	if err != nil {
		return false, err
	}
	if len(cur) > 0 {
		var wantCmp map[string]any
		b, _ := json.Marshal(want)
		_ = json.Unmarshal(b, &wantCmp)
		if reflect.DeepEqual(cur, wantCmp) {
			return false, nil
		}
		if !onlyGrimoireCopilotHooks(cur) {
			return false, fmt.Errorf("%s exists and was not written by grimoire — remove or rename it, then rerun", path)
		}
	}
	return true, writeSettings(path, want)
}

// RemoveCopilotHooks deletes grimoire's hook file at path if it holds only
// grimoire's hooks. Reports whether it was removed.
func RemoveCopilotHooks(path string) (bool, error) {
	cur, err := readSettings(path)
	if err != nil {
		return false, err
	}
	if len(cur) == 0 || !onlyGrimoireCopilotHooks(cur) {
		return false, nil
	}
	return true, os.Remove(path)
}

// HasCopilotHooks reports whether grimoire's Copilot hook file is in place.
func HasCopilotHooks(path string) bool {
	cur, err := readSettings(path)
	return err == nil && len(cur) > 0 && onlyGrimoireCopilotHooks(cur)
}

// onlyGrimoireCopilotHooks reports whether every hook in a Copilot hook file
// is a grimoire hook.
func onlyGrimoireCopilotHooks(cfg map[string]any) bool {
	hooks, ok := cfg["hooks"].(map[string]any)
	if !ok || len(hooks) == 0 {
		return false
	}
	for _, v := range hooks {
		list, ok := v.([]any)
		if !ok {
			return false
		}
		for _, h := range list {
			if !isGrimoireHook(hookCommand(h)) {
				return false
			}
		}
	}
	return true
}
