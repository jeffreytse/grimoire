package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/compliance"
	"github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/profiles"
	"github.com/jeffreytse/grimoire/internal/rules"
	"github.com/jeffreytse/grimoire/internal/settings"
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
  grimoire mcp config --target claude

To pin a specific project directory (recommended for MCP):
  grimoire --project-dir /path/to/project mcp serve
  GRIMOIRE_PROJECT_DIR=/path/to/project grimoire mcp serve`,
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

	s.AddTool(
		mcp.NewTool("grimoire_get_settings",
			mcp.WithDescription("Return resolved grimoire settings for the current project after merging all layers (project > global > system). Equivalent to grimoire settings --json."),
		),
		toolGrimoireGetSettings,
	)

	s.AddTool(
		mcp.NewTool("grimoire_profile_list",
			mcp.WithDescription("List all available grimoire profiles (project .grimoire/profiles/ and user directories) plus profiles active in settings."),
		),
		toolGrimoireProfileList,
	)

	s.AddTool(
		mcp.NewTool("grimoire_profile_show",
			mcp.WithDescription("Show resolved skills and metadata for a named grimoire profile."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Profile name (without .toml extension)")),
		),
		toolGrimoireProfileShow,
	)

	s.AddTool(
		mcp.NewTool("grimoire_profile_init",
			mcp.WithDescription("Create a new profile file at .grimoire/profiles/<name>.toml with a starter template. Returns the relative path of the created file."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Profile name (without .toml extension)")),
		),
		toolGrimoireProfileInit,
	)

	s.AddTool(
		mcp.NewTool("grimoire_config_get",
			mcp.WithDescription("Get the resolved value and source file for a grimoire config key (e.g. standards.engineering.compliance-threshold)."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Dotted config key")),
		),
		toolGrimoireConfigGet,
	)

	s.AddTool(
		mcp.NewTool("grimoire_config_set",
			mcp.WithDescription("Write a grimoire config key. Level defaults to 'local' for standards.* keys, 'global' for core.* keys."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Dotted config key")),
			mcp.WithString("value", mcp.Required(), mcp.Description("New value")),
			mcp.WithString("level", mcp.Description("Target layer: local | global | system (default: auto)")),
		),
		toolGrimoireConfigSet,
	)

	s.AddTool(
		mcp.NewTool("grimoire_config_unset",
			mcp.WithDescription("Clear a grimoire config key from a settings file."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Dotted config key")),
			mcp.WithString("level", mcp.Description("Target layer: local | global | system (default: auto)")),
		),
		toolGrimoireConfigUnset,
	)

	s.AddTool(
		mcp.NewTool("grimoire_config_add",
			mcp.WithDescription("Append a value to a list-type config key (idempotent). Supports: standards.extends, standards.profiles, standards.<domain>.practices, standards.<domain>.disabled."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Dotted config key (must be a list key)")),
			mcp.WithString("value", mcp.Required(), mcp.Description("Value to append")),
			mcp.WithString("level", mcp.Description("Target layer: local | global | system (default: auto)")),
		),
		toolGrimoireConfigAdd,
	)

	s.AddTool(
		mcp.NewTool("grimoire_config_remove",
			mcp.WithDescription("Remove a value from a list-type config key (idempotent, no error if not present). Mirrors grimoire_config_add."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Dotted config key (must be a list key)")),
			mcp.WithString("value", mcp.Required(), mcp.Description("Value to remove")),
			mcp.WithString("level", mcp.Description("Target layer: local | global | system (default: auto)")),
		),
		toolGrimoireConfigRemove,
	)

	s.AddTool(
		mcp.NewTool("grimoire_install",
			mcp.WithDescription("Install grimoire skills to AI agent directories. Omit all params to install everything."),
			mcp.WithString("domain", mcp.Description("Domain to install (e.g. engineering)")),
			mcp.WithString("subdomain", mcp.Description("Subdomain filter (requires domain)")),
			mcp.WithString("skill", mcp.Description("Single skill ref: domain/subdomain/name")),
			mcp.WithString("target", mcp.Description("Agent: claude, codex, gemini, all, auto (default: auto)")),
		),
		toolGrimoireInstall,
	)

	s.AddTool(
		mcp.NewTool("grimoire_uninstall",
			mcp.WithDescription("Remove grimoire skills from AI agent directories."),
			mcp.WithString("domain", mcp.Description("Domain to uninstall")),
			mcp.WithString("subdomain", mcp.Description("Subdomain filter (requires domain)")),
			mcp.WithString("skill", mcp.Description("Single skill ref: domain/subdomain/name")),
			mcp.WithString("target", mcp.Description("Agent: claude, codex, gemini, all, auto (default: auto)")),
		),
		toolGrimoireUninstall,
	)

	s.AddTool(
		mcp.NewTool("grimoire_update",
			mcp.WithDescription("Pull the latest grimoire skills and relink. Clones if not yet installed."),
			mcp.WithString("stable", mcp.Description("Set 'true' to check out latest tagged release instead of HEAD")),
		),
		toolGrimoireUpdate,
	)

	s.AddTool(
		mcp.NewTool("grimoire_clean",
			mcp.WithDescription("Remove stale grimoire-managed skills from agent directories (broken symlinks and stale copy-mode installs)."),
			mcp.WithString("target", mcp.Description("Agent to clean: claude, codex, gemini, all, auto (default: auto)")),
		),
		toolGrimoireClean,
	)

	s.AddTool(
		mcp.NewTool("grimoire_version",
			mcp.WithDescription("Show grimoire CLI and skills version information."),
		),
		toolGrimoireVersion,
	)

	s.AddTool(
		mcp.NewTool("grimoire_init",
			mcp.WithDescription("Initialize .grimoire/ in the project directory. Creates settings.toml with auto-detected profile."),
			mcp.WithString("profile", mcp.Description("Profile to activate (e.g. engineering, writing, design). Defaults to auto-detected.")),
			mcp.WithNumber("threshold", mcp.Description("Compliance threshold percentage (0–100). Default: 80.")),
			mcp.WithNumber("max_errors", mcp.Description("Max allowed compliance errors. Default: 0.")),
		),
		toolGrimoireInit,
	)

	s.AddTool(
		mcp.NewTool("grimoire_self_update",
			mcp.WithDescription("Check for or apply updates to the grimoire CLI binary. Default is check-only."),
			mcp.WithString("yes", mcp.Description("Set 'true' to apply the update (default: check only, returns available version)")),
		),
		toolGrimoireSelfUpdate,
	)

	s.AddTool(
		mcp.NewTool("grimoire_registry_list",
			mcp.WithDescription("List all configured grimoire skill registries with skill counts and clone status."),
		),
		toolGrimoireRegistryList,
	)

	s.AddTool(
		mcp.NewTool("grimoire_registry_add",
			mcp.WithDescription("Add a registry to standards.extends and clone it. Accepts owner/repo shorthand, full git URL, or absolute local path."),
			mcp.WithString("ref", mcp.Required(), mcp.Description("Registry ref: owner/repo, git URL, or absolute local path")),
		),
		toolGrimoireRegistryAdd,
	)

	s.AddTool(
		mcp.NewTool("grimoire_registry_remove",
			mcp.WithDescription("Remove a registry from standards.extends by name (derived owner/repo or absolute path)."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Registry name as returned by grimoire_registry_list")),
		),
		toolGrimoireRegistryRemove,
	)

	s.AddTool(
		mcp.NewTool("grimoire_registry_update",
			mcp.WithDescription("Pull the latest skills from all registries or a specific one."),
			mcp.WithString("name", mcp.Description("Registry name (omit to update all)")),
		),
		toolGrimoireRegistryUpdate,
	)

	s.AddTool(
		mcp.NewTool("grimoire_registry_set",
			mcp.WithDescription("Set the official (core) registry. Accepts owner/repo shorthand, full git URL, or absolute local path. Run grimoire_update after to apply."),
			mcp.WithString("ref", mcp.Required(), mcp.Description("Registry ref: owner/repo[@version], git URL, or absolute local path")),
		),
		toolGrimoireRegistrySet,
	)

	s.AddTool(
		mcp.NewTool("grimoire_registry_reset",
			mcp.WithDescription("Clear core.registry and revert to the built-in official default (jeffreytse/grimoire-hub)."),
		),
		toolGrimoireRegistryReset,
	)

	s.AddTool(
		mcp.NewTool("grimoire_registry_validate",
			mcp.WithDescription("Validate a registry's structure before publishing. Pass a path or installed registry name; omit to validate the current directory."),
			mcp.WithString("target", mcp.Description("Absolute path, installed registry name (owner/repo), or omit for current directory")),
		),
		toolGrimoireRegistryValidate,
	)

	s.AddTool(
		mcp.NewTool("grimoire_preset_list",
			mcp.WithDescription("List all available presets from all installed registries."),
		),
		toolGrimoirePresetList,
	)
}

func toolGrimoireContext(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	home := skills.OfficialRegistryHome()
	grimoireVer := ""
	if state, err := git.CurrentState(home); err == nil {
		grimoireVer = state.Version
	} else {
		grimoireVer = skills.GrimoireVersion()
	}

	eng := &rules.Engine{
		SkillsSources: skills.AllSkillsSources(),
		ProjectDir:    getProjectDir(),
	}

	var complianceReport *compliance.Report
	if r, err := compliance.Load(resolvedReportPath(getProjectDir())); err == nil {
		complianceReport = r
	}

	out := contextOutput{
		CLIVersion:       strings.TrimPrefix(cliVersion, "v"),
		GrimoireVersion:  grimoireVer,
		GrimoireHome:     home,
		ProfileDirs:      buildProfileDirs(home),
		ResolvedProfiles: buildResolvedProfiles(),
		ProfileSources:   buildProfileSources(),
		DomainSections:   buildDomainSections(),
		SettingsSources:  buildSettingsSources(),
		Agents:           buildAgentInfos(),
		Settings:         buildSettingsMap(),
		Compliance:       complianceReport,
		Registries:       buildRegistryInfos(),
		RuleFindings:     eng.Run(),
	}
	return jsonResult(out)
}

func toolGrimoireCheck(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	report, err := compliance.Load(resolvedReportPath(getProjectDir()))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	eng := &rules.Engine{
		SkillsSources: skills.AllSkillsSources(),
		ProjectDir:    getProjectDir(),
	}
	if found := eng.Run(); len(found) > 0 {
		report.Diagnostics = append(found, report.Diagnostics...)
	}
	return jsonResult(report)
}

func toolGrimoireListSkills(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	all, err := skills.ListAllSkillsFromSources(skills.AllSkillsSources())
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(all)
}

func toolGrimoireDoctor(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	return jsonResult(collectDoctorChecks())
}

func toolGrimoireRunRules(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	eng := &rules.Engine{
		SkillsSources: skills.AllSkillsSources(),
		ProjectDir:    getProjectDir(),
	}
	return jsonResult(eng.Run())
}

func toolGrimoireGetSettings(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	r, err := settings.Load(getProjectDir())
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(settingsToMap(r))
}

func toolGrimoireProfileList(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	return jsonResult(listProfileEntries(getProjectDir(), skills.OfficialRegistryHome()))
}

func toolGrimoireProfileShow(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	name := request.GetString("name", "")
	cwd := getProjectDir()
	p, err := profiles.ResolveWithOptions(name, cwd, resolveOpts(cwd))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(p)
}

func toolGrimoireProfileInit(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	name := request.GetString("name", "")
	dir := filepath.Join(getProjectDir(), ".grimoire", "profiles")
	path := filepath.Join(dir, name+".toml")
	if _, err := os.Stat(path); err == nil {
		return mcp.NewToolResultError("profile already exists: " + path), nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := os.WriteFile(path, []byte(buildProfileTemplate(name)), 0o644); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	rel, _ := filepath.Rel(getProjectDir(), path)
	return jsonResult(map[string]any{"path": rel})
}

func toolGrimoireConfigGet(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	key := request.GetString("key", "")
	r, err := settings.Load(getProjectDir())
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	val, src, err := getKeyResolved(&r, key)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(map[string]any{"value": val, "source": src})
}

func toolGrimoireConfigSet(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	key := request.GetString("key", "")
	value := request.GetString("value", "")
	level := request.GetString("level", "")
	path, err := levelToFilePath(key, level)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	fs, err := settings.LoadFile(path)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := applyKey(&fs, key, value); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := settings.WriteFile(path, fs); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(map[string]any{"key": key, "value": value, "path": path})
}

func toolGrimoireConfigUnset(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	key := request.GetString("key", "")
	level := request.GetString("level", "")
	path, err := levelToFilePath(key, level)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	fs, err := settings.LoadFile(path)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := applyKey(&fs, key, ""); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := settings.WriteFile(path, fs); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(map[string]any{"key": key, "path": path})
}

func toolGrimoireConfigAdd(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	key := request.GetString("key", "")
	value := request.GetString("value", "")
	level := request.GetString("level", "")
	path, err := levelToFilePath(key, level)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	fs, err := settings.LoadFile(path)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := appendToKey(&fs, key, value); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := settings.WriteFile(path, fs); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(map[string]any{"key": key, "value": value, "path": path})
}

func toolGrimoireConfigRemove(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	key := request.GetString("key", "")
	value := request.GetString("value", "")
	level := request.GetString("level", "")
	path, err := levelToFilePath(key, level)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	fs, err := settings.LoadFile(path)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := removeFromKey(&fs, key, value); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := settings.WriteFile(path, fs); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(map[string]any{"key": key, "value": value, "path": path})
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
