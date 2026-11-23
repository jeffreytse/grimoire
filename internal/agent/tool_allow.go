package agent

import (
	"fmt"
	"os"
	"path/filepath"
)

// MCP tools that only read skills. `grimoire mcp install` pre-approves these
// where the agent has a permission allow list, so a skill lookup never stops
// on an approval prompt (and works in headless runs). Tools that change
// config or packages are left to the agent's normal approval flow.
var mcpReadOnlyTools = []string{
	"grimoire_search_skills",
	"grimoire_get_skill",
	"grimoire_list_skills",
	"grimoire_context",
	"grimoire_info",
}

// mcpServerName is the name grimoire registers its MCP server under.
const mcpServerName = "grimoire"

// ToolAllowFile returns the settings file holding ag's permission allow list
// for scope, or "" when grimoire does not pre-approve tools for ag:
//
//	claude       permissions.allow  mcp__grimoire__<tool>   .claude/settings.json
//	antigravity  permissions.allow  mcp(grimoire/<tool>)    ~/.gemini/antigravity-cli/settings.json
func ToolAllowFile(ag, scope, projectDir string) string {
	switch ag {
	case "claude":
		return PromptHookFile(ag, scope, projectDir)
	case "antigravity":
		if scope == "user" {
			return filepath.Join(mustHome(), ".gemini", "antigravity-cli", "settings.json")
		}
	}
	return ""
}

// ToolAllowRules returns the allow-list entries for grimoire's read-only
// tools in ag's syntax.
func ToolAllowRules(ag string) []string {
	var rules []string
	for _, t := range mcpReadOnlyTools {
		switch ag {
		case "claude":
			rules = append(rules, "mcp__"+mcpServerName+"__"+t)
		case "antigravity":
			rules = append(rules, "mcp("+mcpServerName+"/"+t+")")
		}
	}
	return rules
}

// ConfigureToolAllow adds rules to permissions.allow in the settings file at
// path, keeping every other entry and setting. Reports whether it changed.
func ConfigureToolAllow(path string, rules []string) (bool, error) {
	cfg, err := readSettings(path)
	if err != nil {
		return false, err
	}
	perms, allow, err := allowList(cfg, path)
	if err != nil {
		return false, err
	}
	have := map[string]bool{}
	for _, v := range allow {
		if s, ok := v.(string); ok {
			have[s] = true
		}
	}
	changed := false
	for _, r := range rules {
		if !have[r] {
			allow = append(allow, r)
			changed = true
		}
	}
	if !changed {
		return false, nil
	}
	perms["allow"] = allow
	cfg["permissions"] = perms
	return true, writeSettings(path, cfg)
}

// RemoveToolAllow drops rules from permissions.allow in the settings file at
// path. Reports whether it changed.
func RemoveToolAllow(path string, rules []string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}
	cfg, err := readSettings(path)
	if err != nil {
		return false, err
	}
	perms, allow, err := allowList(cfg, path)
	if err != nil {
		return false, err
	}
	drop := map[string]bool{}
	for _, r := range rules {
		drop[r] = true
	}
	kept := make([]any, 0, len(allow))
	for _, v := range allow {
		if s, ok := v.(string); ok && drop[s] {
			continue
		}
		kept = append(kept, v)
	}
	if len(kept) == len(allow) {
		return false, nil
	}
	perms["allow"] = kept
	cfg["permissions"] = perms
	return true, writeSettings(path, cfg)
}

// allowList returns cfg.permissions and its allow array. A value of the wrong
// shape is an error: grimoire does not overwrite settings it cannot read.
func allowList(cfg map[string]any, path string) (perms map[string]any, allow []any, err error) {
	perms = map[string]any{}
	if pv, ok := cfg["permissions"]; ok {
		p, ok := pv.(map[string]any)
		if !ok {
			return nil, nil, fmt.Errorf(`%s: "permissions" is not an object — fix it by hand`, path)
		}
		perms = p
	}
	if av, ok := perms["allow"]; ok {
		a, ok := av.([]any)
		if !ok {
			return nil, nil, fmt.Errorf(`%s: "permissions.allow" is not an array — fix it by hand`, path)
		}
		allow = a
	}
	return perms, allow, nil
}
