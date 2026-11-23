package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Prompt hook that runs `grimoire hook prompt` on every user prompt.
// The instruction line in CLAUDE.md / AGENTS.md / GEMINI.md is advisory: the
// model decides whether a task "needs" a skill and often skips the lookup. A
// prompt hook runs on every prompt regardless, so matching skills are always
// put in front of the model.
//
// Claude Code, Codex and Gemini CLI share one hook protocol: JSON with
// "prompt" and "hook_event_name" on stdin, and
// {"hookSpecificOutput": {"hookEventName", "additionalContext"}} on stdout.
// They also share the settings shape
// {"hooks": {"<event>": [{"hooks": [{"type": "command", "command", "timeout"}]}]}}.
// Only the event name, the timeout unit and the file differ.

// PromptHook describes how one agent registers grimoire's prompt hook.
type PromptHook struct {
	// Event is the hook event under "hooks" in the settings file.
	Event string
	// SessionEvent is the session-start event that injects grimoire's
	// router once per session (and again after /clear or compaction where
	// the agent fires it then).
	SessionEvent string
	// Timeout bounds the hook, in the agent's unit; a slow search must
	// never stall the user's prompt.
	Timeout int
	// Name, when set, is written as the hook's "name" (Gemini shows it in
	// /hooks).
	Name string
	// Trust is a note for the user when the agent will not run the hook
	// until they approve it, or "".
	Trust string
}

// promptHooks are the agents with a prompt hook grimoire can register.
//
//	claude  UserPromptSubmit  timeout s   ~/.claude/settings.json   ./.claude/settings.json
//	codex   UserPromptSubmit  timeout s   $CODEX_HOME/hooks.json    ./.codex/hooks.json
//	gemini  BeforeAgent       timeout ms  ~/.gemini/settings.json   ./.gemini/settings.json
var promptHooks = map[string]PromptHook{
	"claude": {Event: "UserPromptSubmit", SessionEvent: "SessionStart", Timeout: 10},
	"codex": {Event: "UserPromptSubmit", SessionEvent: "SessionStart", Timeout: 10,
		Trust: "Codex skips new hooks until you trust them — run /hooks in codex and trust it"},
	"gemini": {Event: "BeforeAgent", SessionEvent: "SessionStart", Timeout: 10000, Name: "grimoire"},
}

// PromptHookFor returns ag's prompt hook, or false when grimoire has none
// for that agent.
func PromptHookFor(ag string) (PromptHook, bool) {
	h, ok := promptHooks[ag]
	return h, ok
}

// PromptHookFile returns the settings file holding ag's hooks for scope, or
// "" when ag has no prompt hook.
func PromptHookFile(ag, scope, projectDir string) string {
	if _, ok := promptHooks[ag]; !ok {
		return ""
	}
	if scope == "project" {
		switch ag {
		case "codex":
			return filepath.Join(projectDir, ".codex", "hooks.json")
		default:
			return filepath.Join(projectDir, "."+ag, "settings.json")
		}
	}
	home := mustHome()
	switch ag {
	case "codex":
		return filepath.Join(envDir("CODEX_HOME", filepath.Join(home, ".codex")), "hooks.json")
	default:
		return filepath.Join(home, "."+ag, "settings.json")
	}
}

// promptHookArgs and sessionHookArgs follow the grimoire binary in the hook
// commands.
const (
	promptHookArgs  = "hook prompt"
	sessionHookArgs = "hook session"

	// Skills mode (grimoire install): skills are linked into the agent, so
	// these hooks deliver only grimoire's display and status lines.
	skillsPromptHookArgs  = "hook skills-prompt"
	skillsSessionHookArgs = "hook skills-session"
	skillUsedHookArgs     = "hook skill-used"
)

// grimoireHookArgs are all grimoire hook commands, after the binary.
var grimoireHookArgs = []string{promptHookArgs, sessionHookArgs, skillsPromptHookArgs, skillsSessionHookArgs, skillUsedHookArgs}

// skillsHookArgs are the skills-mode hook commands. The skill-used hook is
// shared with MCP mode, so it is removed separately (RemoveSkillUsedHook).
var skillsHookArgs = []string{skillsPromptHookArgs, skillsSessionHookArgs}

// SkillsPromptHookCommand is the skills-mode prompt hook command for bin.
func SkillsPromptHookCommand(bin string) string {
	return shellQuote(bin) + " " + skillsPromptHookArgs
}

// SkillsSessionHookCommand is the skills-mode session hook command for bin.
func SkillsSessionHookCommand(bin string) string {
	return shellQuote(bin) + " " + skillsSessionHookArgs
}

// SkillUsedHookCommand is the hook that records each skill the agent
// invokes, for the status line badge.
func SkillUsedHookCommand(bin string) string {
	return shellQuote(bin) + " " + skillUsedHookArgs
}

// skillUsedEvent and skillUsedMatcher register SkillUsedHookCommand: Claude
// Code runs PostToolUse hooks after its Skill tool (skills mode) and after
// grimoire_get_skill (MCP mode), with the session id the status line keys
// its badge on. Other agents have no status line.
const (
	skillUsedEvent   = "PostToolUse"
	skillUsedMatcher = "Skill|mcp__grimoire__grimoire_get_skill"
)

// PromptHookCommand is the prompt hook command line for the grimoire binary bin.
func PromptHookCommand(bin string) string {
	return shellQuote(bin) + " " + promptHookArgs
}

// SessionHookCommand is the session hook command line for the grimoire binary bin.
func SessionHookCommand(bin string) string {
	return shellQuote(bin) + " " + sessionHookArgs
}

// shellQuote single-quotes bin for a POSIX shell when it needs it.
func shellQuote(bin string) string {
	if strings.ContainsAny(bin, " \t'\"") {
		return "'" + strings.ReplaceAll(bin, "'", `'\''`) + "'"
	}
	return bin
}

// isGrimoireHook reports whether command is a grimoire prompt or session
// hook, whatever binary path it was registered with.
func isGrimoireHook(command string) bool {
	return hookHasArgs(command, grimoireHookArgs)
}

// isSkillsHook reports whether command is one of grimoire's skills-mode hooks.
func isSkillsHook(command string) bool {
	return hookHasArgs(command, skillsHookArgs)
}

// hookHasArgs reports whether command runs the grimoire binary with one of
// argsList, whatever binary path it was registered with.
func hookHasArgs(command string, argsList []string) bool {
	c := strings.TrimSpace(command)
	for _, args := range argsList {
		if !strings.HasSuffix(c, " "+args) {
			continue
		}
		bin := strings.Trim(strings.TrimSpace(strings.TrimSuffix(c, args)), `'"`)
		if filepath.Base(bin) == "grimoire" {
			return true
		}
	}
	return false
}

// ConfigurePromptHook makes the settings file at path run command on every
// prompt under h.Event, replacing any older grimoire hook (e.g. a different
// binary path). Other settings and hooks are kept. Reports whether the file
// changed.
func ConfigurePromptHook(path string, h PromptHook, command string) (bool, error) {
	return configureEventHook(path, h, h.Event, command)
}

// ConfigureSessionHook makes the settings file at path run command at session
// start under h.SessionEvent, like ConfigurePromptHook. It does nothing when
// the agent has no session event.
func ConfigureSessionHook(path string, h PromptHook, command string) (bool, error) {
	if h.SessionEvent == "" {
		return false, nil
	}
	return configureEventHook(path, h, h.SessionEvent, command)
}

func configureEventHook(path string, h PromptHook, event, command string) (bool, error) {
	return configureMatchedHook(path, h, event, "", command)
}

// ConfigureSkillUsedHook registers the skill-used hook (Claude Code's
// PostToolUse on the Skill tool) in the settings file at path. It does
// nothing for other agents.
func ConfigureSkillUsedHook(path, ag, command string) (bool, error) {
	h, ok := promptHooks[ag]
	if !ok || ag != "claude" {
		return false, nil
	}
	return configureMatchedHook(path, h, skillUsedEvent, skillUsedMatcher, command)
}

// configureMatchedHook makes the settings file at path run command under
// event (with matcher, when set), replacing any older grimoire hook there.
func configureMatchedHook(path string, h PromptHook, event, matcher, command string) (bool, error) {
	cfg, err := readSettings(path)
	if err != nil {
		return false, err
	}
	groups, err := hookGroups(cfg, path, event)
	if err != nil {
		return false, err
	}

	// Already present exactly once with this command: nothing to do.
	count, exact := 0, false
	for _, g := range groups {
		for _, hk := range groupHooks(g) {
			if isGrimoireHook(hookCommand(hk)) {
				count++
				exact = exact || hookCommand(hk) == command
			}
		}
	}
	if count == 1 && exact {
		return false, nil
	}

	entry := map[string]any{
		"type":    "command",
		"command": command,
		"timeout": h.Timeout,
	}
	if h.Name != "" {
		entry["name"] = h.Name
	}
	groups = withoutGrimoireHooks(groups)
	group := map[string]any{"hooks": []any{entry}}
	if matcher != "" {
		group["matcher"] = matcher
	}
	groups = append(groups, group)
	setHookGroups(cfg, event, groups)
	return true, writeSettings(path, cfg)
}

// RemovePromptHook drops grimoire's prompt and session hooks from the
// settings file at path, leaving every other hook. Reports whether the file
// changed.
func RemovePromptHook(path string, h PromptHook) (bool, error) {
	return removeHooks(path, []string{h.Event, h.SessionEvent, skillUsedEvent}, isGrimoireHook)
}

// RemoveSkillsHooks drops only grimoire's skills-mode hooks from the settings
// file at path, so an MCP-mode setup on the same agent keeps its hooks.
func RemoveSkillsHooks(path string, h PromptHook) (bool, error) {
	return removeHooks(path, []string{h.Event, h.SessionEvent, skillUsedEvent}, isSkillsHook)
}

// RemoveSkillUsedHook drops grimoire's skill-used hook from the settings file
// at path.
func RemoveSkillUsedHook(path string) (bool, error) {
	return removeHooks(path, []string{skillUsedEvent}, func(c string) bool {
		return hookHasArgs(c, []string{skillUsedHookArgs})
	})
}

// removeHooks drops the hooks matching drop from events in the settings file
// at path, leaving every other hook. Reports whether the file changed.
func removeHooks(path string, events []string, drop func(string) bool) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}
	cfg, err := readSettings(path)
	if err != nil {
		return false, err
	}
	changed := false
	for _, event := range events {
		if event == "" {
			continue
		}
		groups, err := hookGroups(cfg, path, event)
		if err != nil {
			return false, err
		}
		kept := withoutHooks(groups, drop)
		if len(kept) == len(groups) && countHooks(kept) == countHooks(groups) {
			continue
		}
		setHookGroups(cfg, event, kept)
		changed = true
	}
	if !changed {
		return false, nil
	}
	return true, writeSettings(path, cfg)
}

// HasSessionHook reports whether the settings file at path has grimoire's
// session hook under h.SessionEvent.
func HasSessionHook(path string, h PromptHook) bool {
	if h.SessionEvent == "" {
		return false
	}
	cfg, err := readSettings(path)
	if err != nil {
		return false
	}
	groups, err := hookGroups(cfg, path, h.SessionEvent)
	if err != nil {
		return false
	}
	for _, g := range groups {
		for _, hk := range groupHooks(g) {
			if strings.HasSuffix(strings.TrimSpace(hookCommand(hk)), " "+sessionHookArgs) && isGrimoireHook(hookCommand(hk)) {
				return true
			}
		}
	}
	return false
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

// hookGroups returns cfg.hooks.<event>. A value of the wrong shape is
// an error: grimoire does not overwrite settings it cannot read.
func hookGroups(cfg map[string]any, path, event string) ([]any, error) {
	hv, ok := cfg["hooks"]
	if !ok {
		return nil, nil
	}
	hooks, ok := hv.(map[string]any)
	if !ok {
		return nil, fmt.Errorf(`%s: "hooks" is not an object — fix it by hand`, path)
	}
	ev, ok := hooks[event]
	if !ok {
		return nil, nil
	}
	groups, ok := ev.([]any)
	if !ok {
		return nil, fmt.Errorf(`%s: "hooks.%s" is not an array — fix it by hand`, path, event)
	}
	return groups, nil
}

func setHookGroups(cfg map[string]any, event string, groups []any) {
	hooks, ok := cfg["hooks"].(map[string]any)
	if !ok {
		hooks = map[string]any{}
	}
	if len(groups) == 0 {
		delete(hooks, event)
	} else {
		hooks[event] = groups
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
	return withoutHooks(groups, isGrimoireHook)
}

// withoutHooks removes the hooks whose command matches drop from each group
// and drops groups left empty by that. Groups of unknown shape are kept.
func withoutHooks(groups []any, drop func(string) bool) []any {
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
			if !drop(hookCommand(h)) {
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
