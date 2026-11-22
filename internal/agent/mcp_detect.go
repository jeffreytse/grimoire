package agent

import (
	"bytes"
	"os"
	"strings"
)

// MCPSetup is what `grimoire mcp install` left in place for one agent and
// scope, as read back from the agent's own files.
type MCPSetup struct {
	// Instruction: the grimoire instruction line is in the agent's
	// instruction file (for Antigravity: injected by its SessionStart hook).
	Instruction bool
	// Hook: grimoire's prompt hook is registered.
	Hook bool
	// HookSupported is false for agents that cannot run a prompt hook.
	HookSupported bool
}

// Active reports whether grimoire is set up in MCP mode here at all.
func (s MCPSetup) Active() bool { return s.Instruction || s.Hook }

// DetectMCPSetup reads ag's files for scope ("user" or "project").
func DetectMCPSetup(ag, scope, projectDir string) MCPSetup {
	var s MCPSetup
	if p := MCPInstructionFile(ag, scope, projectDir); p != "" {
		if b, err := os.ReadFile(p); err == nil && strings.Contains(string(b), MCPTriggerLine) {
			s.Instruction = true
		}
	}
	switch ag {
	case "opencode":
		s.HookSupported = true
		if b, err := os.ReadFile(OpenCodePluginFile(scope, projectDir)); err == nil {
			s.Hook = bytes.HasPrefix(b, []byte(openCodePluginMarker))
		}
	case "antigravity":
		s.HookSupported = true
		if cfg, err := readSettings(AntigravityHooksFile(scope, projectDir)); err == nil {
			_, s.Hook = cfg[agyHookName]
		}
		s.Instruction = s.Instruction || s.Hook
	default:
		h, ok := PromptHookFor(ag)
		if !ok {
			return s
		}
		s.HookSupported = true
		path := PromptHookFile(ag, scope, projectDir)
		cfg, err := readSettings(path)
		if err != nil {
			return s
		}
		groups, err := hookGroups(cfg, path, h.Event)
		if err != nil {
			return s
		}
		for _, g := range groups {
			for _, hk := range groupHooks(g) {
				if isGrimoireHook(hookCommand(hk)) {
					s.Hook = true
				}
			}
		}
	}
	return s
}

// LinkedSkillCount counts the symlinks in ag's skills dir that point into
// root (a grimoire package home), i.e. skills linked by `grimoire install`.
func LinkedSkillCount(ag, root string) int {
	dir := SkillsDir(ag)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	prefix := strings.TrimRight(root, "/") + "/"
	n := 0
	for _, e := range entries {
		if e.Type()&os.ModeSymlink == 0 {
			continue
		}
		if t, err := os.Readlink(dir + "/" + e.Name()); err == nil && strings.HasPrefix(t, prefix) {
			n++
		}
	}
	return n
}
