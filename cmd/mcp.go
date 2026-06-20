package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/compliance"
	"github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/rules"
	"github.com/jeffreytse/grimoire/internal/skills"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "MCP server integration",
	Long: `Model Context Protocol (MCP) integration for AI assistants.

  grimoire mcp serve           Start the MCP server (stdio transport)
  grimoire mcp config          Print MCP server configuration for your AI assistant`,
}

var mcpServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the grimoire MCP server (stdio transport)",
	Long: `Start the grimoire MCP server using stdio transport.

AI assistants that support MCP (Claude Code, Cursor, Windsurf, Cline, etc.)
can connect to grimoire tools natively once configured.

To configure Claude Code, add to ~/.claude/mcp.json:
  grimoire mcp config --target claude`,
	RunE: runMCPServe,
}

var flagMCPConfigTarget string

var mcpConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Print MCP configuration snippet for an AI assistant",
	RunE:  runMCPConfig,
}

func init() {
	mcpConfigCmd.Flags().StringVar(&flagMCPConfigTarget, "target", "claude", "AI assistant to configure (claude, cursor, windsurf, cline)")
	mcpCmd.AddCommand(mcpServeCmd)
	mcpCmd.AddCommand(mcpConfigCmd)
}

func runMCPServe(_ *cobra.Command, _ []string) error {
	s := server.NewMCPServer(
		"grimoire",
		strings.TrimPrefix(cliVersion, "v"),
	)

	registerMCPTools(s)

	stdio := server.NewStdioServer(s)
	return stdio.Listen(context.Background(), os.Stdin, os.Stdout)
}

func registerMCPTools(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("grimoire_context",
			mcp.WithDescription("Get full grimoire environment context: CLI version, grimoire version, installed agents, resolved settings, last compliance report, registries, and structural rule findings. Call this at session start to get deterministic ground truth."),
		),
		toolGrimoireContext,
	)

	s.AddTool(
		mcp.NewTool("grimoire_check",
			mcp.WithDescription("Run the compliance check against the latest compliance report (.grimoire/reports/compliance-latest.json). Returns coverage score, per-practice breakdown, diagnostics, and rule engine findings. Use this to assess BPDD compliance state."),
		),
		toolGrimoireCheck,
	)

	s.AddTool(
		mcp.NewTool("grimoire_list_skills",
			mcp.WithDescription("List all available skills across all registries. Returns skill name, domain, subdomain, tags, and registry. Use this to discover what practices are available before applying BPDD."),
		),
		toolGrimoireListSkills,
	)

	s.AddTool(
		mcp.NewTool("grimoire_settings",
			mcp.WithDescription("Get the resolved grimoire settings for the current project. Shows the effective configuration after merging all layers (env vars > project personal > project shared > global). Use this to understand what standards and profiles are active."),
		),
		toolGrimoireSettings,
	)

	s.AddTool(
		mcp.NewTool("grimoire_doctor",
			mcp.WithDescription("Run a health check on the grimoire installation. Returns status of grimoire source, AI agent installations, and config files. Use this to diagnose setup issues."),
		),
		toolGrimoireDoctor,
	)

	s.AddTool(
		mcp.NewTool("grimoire_run_rules",
			mcp.WithDescription("Run the deterministic structural rule engine. Checks skill directory structure (missing SKILL.md, missing frontmatter fields), broken agent symlinks, and settings.toml parseability. Returns diagnostics tagged source=grimoire-rules."),
		),
		toolGrimoireRunRules,
	)
}

func toolGrimoireContext(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	home := skills.GrimoireHome()
	grimoireVer := ""
	if state, err := git.CurrentState(home); err == nil {
		grimoireVer = state.Version
	} else {
		grimoireVer = skills.GrimoireVersion()
	}

	eng := &rules.Engine{
		SkillsSources: skills.AllSkillsSources(),
		ProjectDir:    ".",
	}

	var complianceReport *compliance.Report
	if r, err := compliance.Load(""); err == nil {
		complianceReport = r
	}

	out := contextOutput{
		CLIVersion:      strings.TrimPrefix(cliVersion, "v"),
		GrimoireVersion: grimoireVer,
		GrimoireHome:    home,
		Agents:          buildAgentInfos(),
		Settings:        buildSettingsMap(),
		Compliance:      complianceReport,
		Registries:      buildRegistryInfos(),
		RuleFindings:    eng.Run(),
	}
	return jsonResult(out)
}

func toolGrimoireCheck(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	report, err := compliance.Load("")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	eng := &rules.Engine{
		SkillsSources: skills.AllSkillsSources(),
		ProjectDir:    ".",
	}
	if found := eng.Run(); len(found) > 0 {
		report.Diagnostics = append(found, report.Diagnostics...)
	}
	return jsonResult(report)
}

func toolGrimoireListSkills(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	all, err := skills.ListAllSkillsFromSources(skills.AllSkillsSources())
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(all)
}

func toolGrimoireSettings(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return jsonResult(buildSettingsMap())
}

func toolGrimoireDoctor(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return jsonResult(collectDoctorChecks())
}

func toolGrimoireRunRules(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	eng := &rules.Engine{
		SkillsSources: skills.AllSkillsSources(),
		ProjectDir:    ".",
	}
	return jsonResult(eng.Run())
}

func jsonResult(v any) (*mcp.CallToolResult, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("marshal error: %v", err)), nil
	}
	return mcp.NewToolResultText(string(b)), nil
}

func runMCPConfig(_ *cobra.Command, _ []string) error {
	switch flagMCPConfigTarget {
	case "claude":
		fmt.Print(`Add to ~/.claude/mcp.json (global) or .claude/mcp.json (project):

{
  "mcpServers": {
    "grimoire": {
      "command": "grimoire",
      "args": ["mcp", "serve"]
    }
  }
}
`)
	case "cursor":
		fmt.Print(`Add to .cursor/mcp.json:

{
  "mcpServers": {
    "grimoire": {
      "command": "grimoire",
      "args": ["mcp", "serve"]
    }
  }
}
`)
	case "windsurf":
		fmt.Print(`Add to ~/.codeium/windsurf/mcp_config.json:

{
  "mcpServers": {
    "grimoire": {
      "command": "grimoire",
      "args": ["mcp", "serve"]
    }
  }
}
`)
	case "cline":
		fmt.Print(`In Cline settings → MCP Servers → Add Server:

Command: grimoire
Args: ["mcp", "serve"]
`)
	default:
		return fmt.Errorf("unknown target %q — supported: claude, cursor, windsurf, cline", flagMCPConfigTarget)
	}
	return nil
}
