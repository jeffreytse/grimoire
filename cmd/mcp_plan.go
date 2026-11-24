package cmd

import (
	"fmt"
	"strings"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/tui"
)

// mcpPlan lists what `grimoire mcp install` (install) or `grimoire mcp
// uninstall` would do for one agent: each command it would run and each file
// it would edit, following the same flags as the real run. Nothing is read
// or written, so it can't say whether a step is already done.
func mcpPlan(a *mcpAgent, scope, serve string, install bool) []string {
	dir := getProjectDir()
	var steps []string
	switch {
	case a.configPath != nil:
		steps = append(steps, "edit "+a.configPath(scope)+" (grimoire server entry)")
	case install:
		steps = append(steps, "run "+a.bin+" "+strings.Join(a.add(scope, serve), " "))
	default:
		steps = append(steps, "run "+a.bin+" "+strings.Join(a.remove(scope), " "))
	}
	// uninstall always cleans up everything install might have added.
	configure := !install || !flagMCPSetupNoConfigure
	if !configure {
		return steps
	}
	if p := agent.MCPInstructionFile(a.key, scope, dir); p != "" {
		steps = append(steps, "edit "+p+" (grimoire instruction)")
	}
	if !install || !flagMCPSetupNoAllow {
		if p := agent.ToolAllowFile(a.key, scope, dir); p != "" {
			steps = append(steps, "edit "+p+" (pre-approved read-only tools)")
		}
	}
	if a.key == "claude" && (!install || !flagMCPSetupNoStatus) {
		steps = append(steps, "edit "+agent.StatusLineFile(scope, dir)+" (status line badge)")
	}
	if install && flagMCPSetupNoHook {
		return steps
	}
	switch a.key {
	case "opencode":
		steps = append(steps, "write "+agent.OpenCodePluginFile(scope, dir)+" (prompt hook plugin)")
	case "pi":
		steps = append(steps, "write "+agent.PiExtensionFile(scope, dir)+" (prompt hook extension)")
	case "copilot":
		steps = append(steps, "write "+agent.CopilotHooksFile(scope, dir)+" (session and prompt hooks)")
	case "omp":
		steps = append(steps, "write "+agent.OMPExtensionFile(scope, dir)+" (prompt hook extension)")
	case "kimi":
		if scope == "user" {
			steps = append(steps, "edit "+agent.KimiConfigFile()+" (prompt hook)")
		}
	case "antigravity":
		steps = append(steps, "edit "+agent.AntigravityHooksFile(scope, dir)+" (session and prompt hooks)")
	default:
		if _, ok := agent.PromptHookFor(a.key); ok {
			steps = append(steps, "edit "+agent.PromptHookFile(a.key, scope, dir)+" (prompt and session hooks)")
		}
	}
	return steps
}

func printMCPPlan(targets []*mcpAgent, serve string, install bool) {
	verb := "set up"
	if !install {
		verb = "remove"
	}
	fmt.Printf("  Dry run — nothing is changed. grimoire mcp %s would %s:\n\n",
		map[bool]string{true: "install", false: "uninstall"}[install], verb)
	for _, a := range targets {
		fmt.Printf("  %s (%s scope)\n", tui.StyleBold.Render(a.display), flagMCPSetupScope)
		for _, step := range mcpPlan(a, flagMCPSetupScope, serve, install) {
			fmt.Printf("     %s\n", step)
		}
	}
	if install {
		fmt.Println("\n  It also downloads the skill library first if it isn't there yet.")
	}
}
