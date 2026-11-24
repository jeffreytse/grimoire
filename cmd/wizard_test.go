package cmd

import (
	"strings"
	"testing"
)

// The wizard sets up each MCP agent at the chosen scope when it can, and at
// its own scope otherwise, saying so.
func TestMCPScopeGroups(t *testing.T) {
	groups, notes := mcpScopeGroups([]string{"claude", "codex"}, "project")
	if strings.Join(groups["project"], ",") != "claude" || strings.Join(groups["user"], ",") != "codex" {
		t.Errorf("groups = %v", groups)
	}
	if len(notes) != 1 || !strings.Contains(notes[0], "Codex supports user scope only") {
		t.Errorf("notes = %v", notes)
	}
	if groups, notes := mcpScopeGroups([]string{"claude", "codex"}, "user"); len(groups["user"]) != 2 || len(notes) != 0 {
		t.Errorf("user scope: %v %v", groups, notes)
	}
}
