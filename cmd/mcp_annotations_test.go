package cmd

import (
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

// Every registered tool has a title and hints; a new tool without an entry in
// mcpToolMeta would fall back to "destructive, open-world write".
func TestMCPTools_Annotated(t *testing.T) {
	s := server.NewMCPServer("grimoire", "test")
	registerMCPTools(s)
	tools := s.ListTools()
	if len(tools) == 0 {
		t.Fatal("no tools registered")
	}
	for name, st := range tools {
		m, ok := mcpToolMeta[name]
		if !ok {
			t.Errorf("%s: no entry in mcpToolMeta", name)
			continue
		}
		a := st.Tool.Annotations
		if st.Tool.Title != m.title || a.Title != m.title {
			t.Errorf("%s: title = %q / %q", name, st.Tool.Title, a.Title)
		}
		if a.ReadOnlyHint == nil || *a.ReadOnlyHint != (m.effect == toolRead) {
			t.Errorf("%s: readOnlyHint wrong", name)
		}
	}
	for _, name := range []string{"grimoire_search_skills", "grimoire_get_skill", "grimoire_list_skills"} {
		a := tools[name].Tool.Annotations
		if !*a.ReadOnlyHint || *a.DestructiveHint || *a.OpenWorldHint {
			t.Errorf("%s: want read-only, non-destructive, closed-world", name)
		}
	}
	for name := range mcpToolMeta {
		if _, ok := tools[name]; !ok {
			t.Errorf("mcpToolMeta has %s, which is not registered", name)
		}
	}
}

// The routing directive rides on the skill tools' descriptions: the one
// channel every MCP client shows the model (Copilot CLI drops server
// instructions by default, and several agents have no hooks).
func TestMCPTools_DescriptionsCarryRouteDirective(t *testing.T) {
	s := server.NewMCPServer("grimoire", "test")
	registerMCPTools(s)
	for _, name := range []string{"grimoire_get_skill", "grimoire_search_skills"} {
		if d := s.GetTool(name).Tool.Description; !strings.HasPrefix(d, mcpRouteDirective) {
			t.Errorf("%s description = %q", name, d)
		}
	}
}
