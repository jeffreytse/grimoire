package cmd

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCP tool titles and behavior hints. Clients show the title in place of the
// raw tool name, and use the hints to group read-only lookups and to decide
// what needs approval. mcp-go's defaults mark every tool as a destructive,
// open-world write, which overstates every skill lookup.

type toolEffect int

const (
	toolRead        toolEffect = iota // reads local state only
	toolWrite                         // changes config or files, additively
	toolDestructive                   // removes or resets something
)

type toolMeta struct {
	title  string
	effect toolEffect
	// network: the tool reaches outside the machine (git clone/pull,
	// release downloads).
	network bool
}

var mcpToolMeta = map[string]toolMeta{
	"grimoire_search_skills":    {"Search grimoire skills", toolRead, false},
	"grimoire_get_skill":        {"Load grimoire skill", toolRead, false},
	"grimoire_list_skills":      {"Browse grimoire skills", toolRead, false},
	"grimoire_search":           {"Search grimoire packages", toolRead, false},
	"grimoire_info":             {"Show grimoire skill info", toolRead, false},
	"grimoire_context":          {"Show grimoire context", toolRead, false},
	"grimoire_check":            {"Read grimoire compliance report", toolRead, false},
	"grimoire_status":           {"Show grimoire status", toolRead, false},
	"grimoire_doctor":           {"Check grimoire health", toolRead, false},
	"grimoire_run_rules":        {"Run grimoire structural rules", toolRead, false},
	"grimoire_get_config":       {"Show grimoire config", toolRead, false},
	"grimoire_config_get":       {"Get grimoire config value", toolRead, false},
	"grimoire_profile_list":     {"List grimoire profiles", toolRead, false},
	"grimoire_profile_show":     {"Show grimoire profile", toolRead, false},
	"grimoire_package_list":     {"List grimoire packages", toolRead, false},
	"grimoire_package_validate": {"Validate grimoire package", toolRead, false},
	"grimoire_version":          {"Show grimoire version", toolRead, false},

	"grimoire_profile_init":   {"Create grimoire profile", toolWrite, false},
	"grimoire_config_set":     {"Set grimoire config value", toolWrite, false},
	"grimoire_config_add":     {"Add to grimoire config list", toolWrite, false},
	"grimoire_init":           {"Initialize grimoire in project", toolWrite, false},
	"grimoire_package_enable": {"Enable grimoire package", toolWrite, false},
	"grimoire_package_set":    {"Set grimoire package option", toolWrite, false},
	"grimoire_install":        {"Install grimoire skills", toolWrite, true},
	"grimoire_update":         {"Update grimoire skills", toolWrite, true},
	"grimoire_package_add":    {"Add grimoire package", toolWrite, true},
	"grimoire_package_update": {"Update grimoire package", toolWrite, true},
	"grimoire_self_update":    {"Update grimoire CLI", toolWrite, true},

	"grimoire_config_unset":    {"Unset grimoire config value", toolDestructive, false},
	"grimoire_config_remove":   {"Remove from grimoire config list", toolDestructive, false},
	"grimoire_uninstall":       {"Uninstall grimoire skills", toolDestructive, false},
	"grimoire_clean":           {"Clean stale grimoire skills", toolDestructive, false},
	"grimoire_package_remove":  {"Remove grimoire package", toolDestructive, false},
	"grimoire_package_disable": {"Disable grimoire package", toolDestructive, false},
	"grimoire_package_reset":   {"Reset grimoire package", toolDestructive, false},
}

// annotateTool sets t's title and hints from mcpToolMeta. Unknown tools keep
// mcp-go's conservative defaults.
func annotateTool(t mcp.Tool) mcp.Tool { //nolint:gocritic // mcp.Tool is passed by value throughout mcp-go
	m, ok := mcpToolMeta[t.Name]
	if !ok {
		return t
	}
	t.Title = m.title
	t.Annotations.Title = m.title
	t.Annotations.ReadOnlyHint = mcp.ToBoolPtr(m.effect == toolRead)
	t.Annotations.DestructiveHint = mcp.ToBoolPtr(m.effect == toolDestructive)
	t.Annotations.IdempotentHint = mcp.ToBoolPtr(m.effect != toolWrite || !m.network)
	t.Annotations.OpenWorldHint = mcp.ToBoolPtr(m.network)
	return t
}

// addTool registers t, annotated, with handler h.
func addTool(s *server.MCPServer, t mcp.Tool, h server.ToolHandlerFunc) { //nolint:gocritic // see annotateTool
	s.AddTool(annotateTool(t), h)
}
