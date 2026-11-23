package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/compliance"
	"github.com/jeffreytse/grimoire/internal/config"
	"github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/profiles"
	"github.com/jeffreytse/grimoire/internal/rules"
	"github.com/jeffreytse/grimoire/internal/skills"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "MCP server integration",
	Long: `Model Context Protocol (MCP) integration for AI assistants.

  grimoire mcp install         Set up MCP mode: register the server + add the agent instruction
  grimoire mcp uninstall       Undo mcp install
  grimoire mcp serve           Start the MCP server (stdio transport)
  grimoire mcp config          Print MCP server configuration for your AI assistant`,
}

var mcpServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the grimoire MCP server (stdio transport)",
	Long: `Start the grimoire MCP server using stdio transport.

AI assistants that support MCP (Claude Code, Cursor, Windsurf, Cline, etc.)
can connect to grimoire tools natively once configured.

To set up Claude Code (server + CLAUDE.md instruction):
  grimoire mcp install

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
		mcpServerOptions()...,
	)

	registerMCPTools(s)
	registerMCPPrompts(s)

	stdio := server.NewStdioServer(s)
	return stdio.Listen(context.Background(), os.Stdin, os.Stdout)
}

func registerMCPTools(s *server.MCPServer) {
	addTool(s,
		mcp.NewTool("grimoire_context",
			mcp.WithDescription("Get full grimoire environment context: CLI version, grimoire version, installed agents, resolved settings, last compliance report, packages, and structural rule findings. Call this at session start to get deterministic ground truth."),
		),
		toolGrimoireContext,
	)

	addTool(s,
		mcp.NewTool("grimoire_check",
			mcp.WithDescription("Read the cached compliance report (.grimoire/reports/compliance-latest.json). Returns coverage score, per-practice breakdown, diagnostics, and rule engine findings. Does not re-run the check — run `grimoire check` first for a fresh analysis. Use this to assess BPDD compliance state."),
		),
		toolGrimoireCheck,
	)

	addTool(s,
		mcp.NewTool("grimoire_list_skills",
			mcp.WithDescription("Browse available skills. Without arguments, returns a compact index of domains and subdomains with skill counts. Pass a domain (and optionally a subdomain) to list those skills with name, description, tags, and package. For keyword lookup use grimoire_search_skills; to read a skill use grimoire_get_skill."),
			mcp.WithString("domain", mcp.Description("Domain to list (e.g. engineering). Omit for the domain index.")),
			mcp.WithString("subdomain", mcp.Description("Subdomain filter within the domain (e.g. testing)")),
		),
		toolGrimoireListSkills,
	)

	addTool(s,
		mcp.NewTool("grimoire_get_skill",
			mcp.WithDescription("Fetch the full SKILL.md content for a skill by exact name. When the user writes `grimoire:<skill-name>`, call this with that name directly (no search). A leading `grimoire:` is accepted; an unknown name returns close matches. Returns name, full content, and source path. To load several skills at once (e.g. the steps of a plan), pass `names` instead of `name`."),
			mcp.WithString("name", mcp.Description("Skill name (e.g. propose-conventional-commit)")),
			mcp.WithArray("names", mcp.WithStringItems(), mcp.Description("Several skill names to load in one call (max 10); used instead of name")),
		),
		toolGrimoireGetSkill,
	)

	addTool(s,
		mcp.NewTool("grimoire_search_skills",
			mcp.WithDescription("Search skills using BM25 relevance ranking. Handles semantic queries like 'write better commit messages' matching propose-conventional-commit — better than exact substring matching. Returns top matches with scores."),
			mcp.WithString("query", mcp.Required(), mcp.Description("Natural language or keyword search query")),
			mcp.WithNumber("top_n", mcp.Description("Maximum results to return (default: 10, max: 100)")),
		),
		toolGrimoireSearchSkills,
	)

	addTool(s,
		mcp.NewTool("grimoire_doctor",
			mcp.WithDescription("Run a health check on the grimoire installation. Returns status of grimoire source, AI agent installations, and config files. Use this to diagnose setup issues."),
		),
		toolGrimoireDoctor,
	)

	addTool(s,
		mcp.NewTool("grimoire_run_rules",
			mcp.WithDescription("Run the deterministic structural rule engine. Checks skill directory structure (missing SKILL.md, missing frontmatter fields), broken agent symlinks, and grimoire.toml parseability. Returns diagnostics tagged source=grimoire-rules."),
		),
		toolGrimoireRunRules,
	)

	addTool(s,
		mcp.NewTool("grimoire_get_config",
			mcp.WithDescription("Return resolved grimoire config for the current project after merging all layers (project > global > system). Equivalent to grimoire config show --json."),
		),
		toolGrimoireGetConfig,
	)

	addTool(s,
		mcp.NewTool("grimoire_profile_list",
			mcp.WithDescription("List all available grimoire profiles (project .grimoire/profiles/ and user directories) plus profiles active in config."),
		),
		toolGrimoireProfileList,
	)

	addTool(s,
		mcp.NewTool("grimoire_profile_show",
			mcp.WithDescription("Show resolved skills and metadata for a named grimoire profile."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Profile name (without .toml extension)")),
		),
		toolGrimoireProfileShow,
	)

	addTool(s,
		mcp.NewTool("grimoire_profile_init",
			mcp.WithDescription("Create a new profile file at .grimoire/profiles/<name>.toml with a starter template. Returns the relative path of the created file."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Profile name (without .toml extension)")),
		),
		toolGrimoireProfileInit,
	)

	addTool(s,
		mcp.NewTool("grimoire_config_get",
			mcp.WithDescription("Get the resolved value and source file for a grimoire config key (e.g. standards.engineering.compliance-threshold)."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Dotted config key")),
		),
		toolGrimoireConfigGet,
	)

	addTool(s,
		mcp.NewTool("grimoire_config_set",
			mcp.WithDescription("Write a grimoire config key. Level defaults to 'local' for standards.* keys, 'global' for core.* keys."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Dotted config key")),
			mcp.WithString("value", mcp.Required(), mcp.Description("New value")),
			mcp.WithString("level", mcp.Description("Target layer: local | global | system (default: auto)")),
		),
		toolGrimoireConfigSet,
	)

	addTool(s,
		mcp.NewTool("grimoire_config_unset",
			mcp.WithDescription("Clear a grimoire config key from a config file."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Dotted config key")),
			mcp.WithString("level", mcp.Description("Target layer: local | global | system (default: auto)")),
		),
		toolGrimoireConfigUnset,
	)

	addTool(s,
		mcp.NewTool("grimoire_config_add",
			mcp.WithDescription("Append a value to a list-type config key (idempotent). Supports: standards.extends, standards.profiles, standards.<domain>.practices, standards.<domain>.disabled."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Dotted config key (must be a list key)")),
			mcp.WithString("value", mcp.Required(), mcp.Description("Value to append")),
			mcp.WithString("level", mcp.Description("Target layer: local | global | system (default: auto)")),
		),
		toolGrimoireConfigAdd,
	)

	addTool(s,
		mcp.NewTool("grimoire_config_remove",
			mcp.WithDescription("Remove a value from a list-type config key (idempotent, no error if not present). Mirrors grimoire_config_add."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Dotted config key (must be a list key)")),
			mcp.WithString("value", mcp.Required(), mcp.Description("Value to remove")),
			mcp.WithString("level", mcp.Description("Target layer: local | global | system (default: auto)")),
		),
		toolGrimoireConfigRemove,
	)

	addTool(s,
		mcp.NewTool("grimoire_install",
			mcp.WithDescription("Install grimoire skills to AI agent directories. Omit all params to install everything."),
			mcp.WithString("domain", mcp.Description("Domain to install (e.g. engineering)")),
			mcp.WithString("subdomain", mcp.Description("Subdomain filter (requires domain)")),
			mcp.WithString("skill", mcp.Description("Single skill ref: domain/subdomain/name")),
			mcp.WithString("target", mcp.Description("Agent: claude, codex, gemini, antigravity, openclaw, opencode, copilot, pi, dsh, kimi, omp, grok, muse, minimax, all, auto (default: auto)")),
		),
		toolGrimoireInstall,
	)

	addTool(s,
		mcp.NewTool("grimoire_uninstall",
			mcp.WithDescription("Remove grimoire skills from AI agent directories."),
			mcp.WithString("domain", mcp.Description("Domain to uninstall")),
			mcp.WithString("subdomain", mcp.Description("Subdomain filter (requires domain)")),
			mcp.WithString("skill", mcp.Description("Single skill ref: domain/subdomain/name")),
			mcp.WithString("target", mcp.Description("Agent: claude, codex, gemini, antigravity, openclaw, opencode, copilot, pi, dsh, kimi, omp, grok, muse, minimax, all, auto (default: auto)")),
		),
		toolGrimoireUninstall,
	)

	addTool(s,
		mcp.NewTool("grimoire_update",
			mcp.WithDescription("Pull the latest grimoire skills and relink. Clones if not yet installed."),
			mcp.WithBoolean("stable", mcp.Description("Check out latest tagged release instead of HEAD")),
			mcp.WithBoolean("force", mcp.Description("Discard local package modifications and update")),
		),
		toolGrimoireUpdate,
	)

	addTool(s,
		mcp.NewTool("grimoire_clean",
			mcp.WithDescription("Remove stale grimoire-managed skills from agent directories (broken symlinks and stale copy-mode installs)."),
			mcp.WithString("target", mcp.Description("Agent to clean: claude, codex, gemini, antigravity, openclaw, opencode, copilot, pi, dsh, kimi, omp, grok, muse, minimax, all, auto (default: auto)")),
		),
		toolGrimoireClean,
	)

	addTool(s,
		mcp.NewTool("grimoire_version",
			mcp.WithDescription("Show grimoire CLI and skills version information."),
		),
		toolGrimoireVersion,
	)

	addTool(s,
		mcp.NewTool("grimoire_init",
			mcp.WithDescription("Initialize .grimoire/ in the project directory. Creates grimoire.toml with auto-detected profile."),
			mcp.WithString("profile", mcp.Description("Profile to activate (e.g. engineering, writing, design). Defaults to auto-detected.")),
			mcp.WithNumber("threshold", mcp.Description("Compliance threshold percentage (0–100). Default: 80.")),
			mcp.WithNumber("max_errors", mcp.Description("Max allowed compliance errors. Default: 0.")),
		),
		toolGrimoireInit,
	)

	addTool(s,
		mcp.NewTool("grimoire_self_update",
			mcp.WithDescription("Check for or apply updates to the grimoire CLI binary. Default is check-only."),
			mcp.WithBoolean("yes", mcp.Description("Apply the update (default: check only, returns available version)")),
		),
		toolGrimoireSelfUpdate,
	)

	addTool(s,
		mcp.NewTool("grimoire_package_list",
			mcp.WithDescription("List all configured grimoire skill packages with skill counts and clone status."),
		),
		toolGrimoirePackageList,
	)

	addTool(s,
		mcp.NewTool("grimoire_package_add",
			mcp.WithDescription("Add a package to standards.extends and clone it. Accepts owner/repo shorthand, full git URL, or absolute local path."),
			mcp.WithString("ref", mcp.Required(), mcp.Description("Grimoire ref: owner/repo[@version], git URL, or absolute local path")),
		),
		toolGrimoirePackageAdd,
	)

	addTool(s,
		mcp.NewTool("grimoire_package_remove",
			mcp.WithDescription("Remove a package from standards.extends by name (derived owner/repo or absolute path)."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Package name as returned by grimoire_package_list")),
		),
		toolGrimoirePackageRemove,
	)

	addTool(s,
		mcp.NewTool("grimoire_package_enable",
			mcp.WithDescription("Enable a previously disabled package."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Package name as returned by grimoire_package_list")),
		),
		toolGrimoirePackageEnable,
	)

	addTool(s,
		mcp.NewTool("grimoire_package_disable",
			mcp.WithDescription("Disable a package without removing it. Skills from this package will not be loaded."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Package name as returned by grimoire_package_list")),
		),
		toolGrimoirePackageDisable,
	)

	addTool(s,
		mcp.NewTool("grimoire_package_update",
			mcp.WithDescription("Pull the latest skills from all packages or a specific one."),
			mcp.WithString("name", mcp.Description("Package name (omit to update all)")),
			mcp.WithBoolean("force", mcp.Description("Discard local package modifications and update")),
		),
		toolGrimoirePackageUpdate,
	)

	addTool(s,
		mcp.NewTool("grimoire_package_set",
			mcp.WithDescription("Set the official (core) package. Accepts owner/repo shorthand, full git URL, or absolute local path. Run grimoire_update after to apply."),
			mcp.WithString("ref", mcp.Required(), mcp.Description("Grimoire ref: owner/repo[@version], git URL, or absolute local path")),
		),
		toolGrimoirePackageSet,
	)

	addTool(s,
		mcp.NewTool("grimoire_package_reset",
			mcp.WithDescription("Clear core.package and revert to the built-in official default (jeffreytse/grimoire-core)."),
		),
		toolGrimoirePackageReset,
	)

	addTool(s,
		mcp.NewTool("grimoire_package_validate",
			mcp.WithDescription("Validate a package's structure before publishing. Pass a path or installed package name; omit to validate the current directory."),
			mcp.WithString("target", mcp.Description("Absolute path, installed package name (owner/repo), or omit for current directory")),
		),
		toolGrimoirePackageValidate,
	)

	addTool(s,
		mcp.NewTool("grimoire_status",
			mcp.WithDescription("Show project compliance health: active profile, skill count, last check result and age, staleness."),
		),
		toolGrimoireStatus,
	)

	addTool(s,
		mcp.NewTool("grimoire_search",
			mcp.WithDescription("Search installed packages for skills matching a query. Searches name, description, domain, and tags."),
			mcp.WithString("query", mcp.Description("Search query (omit to list all skills)")),
		),
		toolGrimoireSearch,
	)

	addTool(s,
		mcp.NewTool("grimoire_info",
			mcp.WithDescription("Show metadata for an installed skill: description, domain, tags, path, authors, dependencies."),
			mcp.WithString("skill", mcp.Required(), mcp.Description("Skill name (e.g. apply-solid-principles)")),
		),
		toolGrimoireInfo,
	)
}

func toolGrimoireContext(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	home := skills.OfficialPackageHome()
	grimoireVer := ""
	if state, err := git.CurrentState(home); err == nil {
		grimoireVer = state.Version
	} else {
		grimoireVer = skills.GrimoireVersion()
	}

	eng := &rules.Engine{
		SkillsPackages: skills.AllSkillsPackages(),
		ProjectDir:     getProjectDir(),
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
		ConfigSources:    buildConfigSources(),
		Agents:           buildAgentInfos(),
		Config:           buildConfigMap(),
		Compliance:       complianceReport,
		Packages:         buildPackageInfos(),
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
		SkillsPackages: skills.AllSkillsPackages(),
		ProjectDir:     getProjectDir(),
	}
	if found := eng.Run(); len(found) > 0 {
		report.Diagnostics = append(found, report.Diagnostics...)
	}
	return jsonResult(report)
}

// mcpSkillSummary is the per-skill shape grimoire_list_skills returns: metadata
// only. Full SKILL.md bodies for every skill run to megabytes, far beyond what
// an MCP client will pass to a model; grimoire_get_skill loads one on demand.
type mcpSkillSummary struct {
	Name        string   `json:"name"`
	Subdomain   string   `json:"subdomain,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Package     string   `json:"package,omitempty"`
}

type mcpDomainIndex struct {
	Domain     string         `json:"domain"`
	Total      int            `json:"total"`
	Subdomains map[string]int `json:"subdomains,omitempty"`
}

const domainListDescMax = 100

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return strings.TrimSpace(string(r[:n])) + "…"
}

func toolGrimoireListSkills(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	domain := request.GetString("domain", "")
	subdomain := request.GetString("subdomain", "")
	if subdomain != "" && domain == "" {
		return mcp.NewToolResultError("subdomain requires domain"), nil
	}
	// The Meta variant skips frontmatter (no description/tags), so parse fully;
	// only the summary fields below are sent to the client.
	all, _, err := skills.ListAllSkillsFromPackages(skills.AllSkillsPackages())
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if domain == "" {
		byDomain := map[string]*mcpDomainIndex{}
		var order []string
		for i := range all {
			sk := &all[i]
			d := byDomain[sk.Domain]
			if d == nil {
				d = &mcpDomainIndex{Domain: sk.Domain, Subdomains: map[string]int{}}
				byDomain[sk.Domain] = d
				order = append(order, sk.Domain)
			}
			d.Total++
			if sk.Subdomain != "" {
				d.Subdomains[sk.Subdomain]++
			}
		}
		sort.Strings(order)
		index := make([]mcpDomainIndex, 0, len(order))
		for _, name := range order {
			index = append(index, *byDomain[name])
		}
		return jsonResult(map[string]any{
			"total":   len(all),
			"domains": index,
			"hint":    "Pass domain (and optional subdomain) to list skills; grimoire_get_skill loads one.",
		})
	}

	var out []mcpSkillSummary
	for i := range all {
		sk := &all[i]
		if sk.Domain != domain || (subdomain != "" && sk.Subdomain != subdomain) {
			continue
		}
		sum := mcpSkillSummary{
			Name:        sk.Name,
			Subdomain:   sk.Subdomain,
			Description: sk.Description,
			Tags:        sk.Tags,
			Package:     sk.Package,
		}
		if subdomain == "" {
			// A whole domain can hold ~300 skills; keep it within typical MCP
			// output limits (~25k tokens). Narrow by subdomain for full detail.
			sum.Description = truncateRunes(sum.Description, domainListDescMax)
			sum.Tags, sum.Package = nil, ""
		}
		out = append(out, sum)
	}
	if len(out) == 0 {
		scope := domain
		if subdomain != "" {
			scope += "/" + subdomain
		}
		return mcp.NewToolResultError(fmt.Sprintf("no skills in %q — call grimoire_list_skills without arguments for the index", scope)), nil
	}
	return jsonResult(out)
}

func toolGrimoireDoctor(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	return jsonResult(collectDoctorChecks())
}

func toolGrimoireRunRules(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	eng := &rules.Engine{
		SkillsPackages: skills.AllSkillsPackages(),
		ProjectDir:     getProjectDir(),
	}
	return jsonResult(eng.Run())
}

func toolGrimoireGetConfig(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	r, err := config.Load(getProjectDir())
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(configToMap(r))
}

func toolGrimoireProfileList(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	return jsonResult(listProfileEntries(getProjectDir(), skills.OfficialPackageHome()))
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
	if strings.TrimSpace(name) == "" {
		return mcp.NewToolResultError("name is required"), nil
	}
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
	r, err := config.Load(getProjectDir())
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
	fs, err := config.LoadFile(path)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := applyKey(&fs, key, value); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := config.WriteFile(path, fs); err != nil {
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
	fs, err := config.LoadFile(path)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := applyKey(&fs, key, ""); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := config.WriteFile(path, fs); err != nil {
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
	fs, err := config.LoadFile(path)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := appendToKey(&fs, key, value); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := config.WriteFile(path, fs); err != nil {
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
	fs, err := config.LoadFile(path)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := removeFromKey(&fs, key, value); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := config.WriteFile(path, fs); err != nil {
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
		fmt.Print(`Easiest: grimoire mcp install   (registers the server and adds the CLAUDE.md instruction)

Or register with the Claude Code CLI:

  claude mcp add -s user grimoire -- grimoire mcp serve      # all projects (~/.claude.json)
  claude mcp add -s project grimoire -- grimoire mcp serve   # this repo (.mcp.json, commit to share)

Or write .mcp.json at the project root yourself:

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
