package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

// MCP mode: the agent reaches grimoire skills through `grimoire mcp serve`
// instead of skills linked into its skills directory. `grimoire mcp install`
// registers the server with each agent and writes the instruction that makes
// the agent consult it.

var (
	flagMCPSetupTarget      string
	flagMCPSetupScope       string
	flagMCPSetupNoConfigure bool
	flagMCPSetupNoHook      bool
	flagMCPSetupNoAllow     bool
)

// mcpAgent describes how to set up MCP mode for one agent. Registration
// goes through the agent's own CLI (add/remove), or through its config file
// when the CLI has no non-interactive add (fileRegister/fileUnregister).
// Each entry was checked against the real CLI: register, re-register, list.
type mcpAgent struct {
	key, display, bin string
	scopes            []string
	add               func(scope, serve string) []string
	remove            func(scope string) []string
	existsMarker      string // CLI output meaning "already registered" on add
	fileRegister      func(scope, serve string) (bool, error)
	fileUnregister    func(scope string) (bool, error)
	verify            string
}

var mcpAgents = []mcpAgent{
	{
		key: "claude", display: "Claude Code", bin: "claude", scopes: []string{"user", "project"},
		add: func(scope, serve string) []string {
			return []string{"mcp", "add", "-s", scope, "grimoire", "--", serve, "mcp", "serve"}
		},
		remove:       func(scope string) []string { return []string{"mcp", "remove", "-s", scope, "grimoire"} },
		existsMarker: "already exists",
		verify:       "/mcp in a new session",
	},
	{
		key: "codex", display: "Codex", bin: "codex", scopes: []string{"user"},
		add:    func(_, serve string) []string { return []string{"mcp", "add", "grimoire", "--", serve, "mcp", "serve"} },
		remove: func(string) []string { return []string{"mcp", "remove", "grimoire"} },
		verify: "codex mcp list",
	},
	{
		key: "gemini", display: "Gemini CLI", bin: "gemini", scopes: []string{"user", "project"},
		add: func(scope, serve string) []string {
			return []string{"mcp", "add", "-s", scope, "grimoire", serve, "mcp", "serve"}
		},
		remove: func(scope string) []string { return []string{"mcp", "remove", "-s", scope, "grimoire"} },
		// Gemini keeps all MCP servers disabled in folders the user hasn't trusted.
		verify: "gemini mcp list (servers stay disabled until you trust the folder in Gemini)",
	},
	{
		key: "copilot", display: "GitHub Copilot CLI", bin: "copilot", scopes: []string{"user"},
		add:          func(_, serve string) []string { return []string{"mcp", "add", "grimoire", "--", serve, "mcp", "serve"} },
		remove:       func(string) []string { return []string{"mcp", "remove", "grimoire"} },
		existsMarker: "already exists",
		verify:       "copilot mcp list",
	},
	{
		key: "opencode", display: "OpenCode", bin: "opencode", scopes: []string{"user", "project"},
		fileRegister:   registerOpenCodeMCP,
		fileUnregister: unregisterOpenCodeMCP,
		verify:         "opencode mcp list",
	},
	// Pi, Kimi, OMP, Antigravity: documented "mcpServers" JSON files. Their
	// CLIs are not needed for registration (bin is only used by --target all).
	fileAgent("pi", "Pi", "pi", []string{"user", "project"}, commandArgsEntry, "pi mcp list"),
	fileAgent("kimi", "Kimi Code", "kimi", []string{"user", "project"}, commandArgsEntry, "/mcp in a new kimi session"),
	fileAgent("omp", "OMP", "omp", []string{"user", "project"}, func(serve string) map[string]any {
		e := commandArgsEntry(serve)
		e["type"] = "stdio"
		return e
	}, "/mcp in a new omp session"),
	fileAgent("antigravity", "Antigravity", "agy", []string{"user"}, commandArgsEntry, "/mcp in a new agy session"),
}

func fileAgent(key, display, bin string, scopes []string, entry func(string) map[string]any, verify string) mcpAgent {
	f := mcpServersFile(key, entry)
	return mcpAgent{
		key: key, display: display, bin: bin, scopes: scopes,
		fileRegister: f.register, fileUnregister: f.unregister, verify: verify,
	}
}

func mcpAgentKeys() []string {
	keys := make([]string, len(mcpAgents))
	for i := range mcpAgents {
		keys[i] = mcpAgents[i].key
	}
	return keys
}

func findMCPAgent(key string) *mcpAgent {
	for i := range mcpAgents {
		if mcpAgents[i].key == key {
			return &mcpAgents[i]
		}
	}
	return nil
}

func (a *mcpAgent) supports(scope string) bool {
	for _, s := range a.scopes {
		if s == scope {
			return true
		}
	}
	return false
}

var mcpInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Set up grimoire in MCP mode for AI agents",
	Long: `Register the grimoire MCP server with AI agents and tell them to use it.

  1. Fetches the grimoire skills package if it is not present yet.
  2. Registers "grimoire mcp serve" with each agent (its own CLI, or its
     config file for OpenCode).
  3. Writes the grimoire instruction to the agent's instruction file,
     replacing a skills-mode start-best-practice trigger. Re-running is safe:
     nothing is duplicated.
  4. Adds a prompt hook that runs on every prompt and puts the best-matching
     skills (or the skill named as grimoire:<name>) in front of the model, so
     the lookup happens even when the model would skip it:
       claude, codex  UserPromptSubmit in .claude/settings.json, .codex/hooks.json
       gemini         BeforeAgent in .gemini/settings.json
       opencode       plugin plugins/grimoire.js
       antigravity    SessionStart (injects the instruction) + PreInvocation
                      in ~/.gemini/config/hooks.json
     Copilot CLI drops command-hook output, so it gets none. Codex runs the
     hook only after you trust it in /hooks. --no-hook skips it.
  5. Claude Code, Antigravity: pre-approves the read-only skill tools
     (grimoire_search_skills, grimoire_get_skill, grimoire_list_skills) in
     permissions.allow, so lookups never wait on an approval prompt and work
     in headless runs. Tools that change config or packages still ask.
     --no-allow skips it.

Agents and scopes:
  claude    user, project   (instructions: CLAUDE.md)
  codex     user            (~/.codex/AGENTS.md)
  gemini    user, project   (GEMINI.md)
  copilot   user            (~/.copilot/copilot-instructions.md)
  opencode  user, project   (opencode.json + AGENTS.md)
  pi        user, project   (mcp.json + AGENTS.md)
  kimi      user, project   (mcp.json + AGENTS.md)
  omp       user, project   (mcp.json + AGENTS.md)
  antigravity user          (~/.gemini/config/mcp_config.json; no instruction
                            file — the SessionStart hook injects it)

No skills are linked into the agent's skills directory. Start a new agent
session afterwards.`,
	Example: `  grimoire mcp install                          # Claude Code, all projects
  grimoire mcp install --target codex,gemini
  grimoire mcp install --target all             # every supported agent on PATH
  grimoire mcp install --scope project          # this repo only`,
	RunE: runMCPInstall,
}

var mcpUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove grimoire MCP mode from AI agents",
	Long: `Unregister the grimoire MCP server from each agent and remove the grimoire
instruction from its instruction file. The skills package and config are kept.`,
	RunE: runMCPUninstall,
}

func init() {
	for _, c := range []*cobra.Command{mcpInstallCmd, mcpUninstallCmd} {
		c.Flags().StringVar(&flagMCPSetupTarget, "target", "claude",
			"agents, comma-separated: "+strings.Join(mcpAgentKeys(), ", ")+", or all (those on PATH)")
		c.Flags().StringVar(&flagMCPSetupScope, "scope", "user", `"user" (all projects) or "project" (current project)`)
		mcpCmd.AddCommand(c)
	}
	mcpInstallCmd.Flags().BoolVar(&flagMCPSetupNoConfigure, "no-configure", false, "do not write the grimoire instruction or hook")
	mcpInstallCmd.Flags().BoolVar(&flagMCPSetupNoHook, "no-hook", false, "do not add the prompt hook (Claude Code, Codex, Gemini CLI)")
	mcpInstallCmd.Flags().BoolVar(&flagMCPSetupNoAllow, "no-allow", false, "do not pre-approve grimoire's read-only skill tools (Claude Code, Antigravity)")
}

// mcpRunCmd runs an external command in dir and returns combined output.
// Tests replace it.
var mcpRunCmd = func(dir, name string, args ...string) (string, error) {
	c := exec.Command(name, args...)
	c.Dir = dir
	out, err := c.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// mcpLookPath is exec.LookPath; tests replace it.
var mcpLookPath = exec.LookPath

// resolveMCPAgents validates flags and returns the agents to act on.
// "all" selects every supported agent whose CLI is on PATH and that supports
// the scope; naming an agent explicitly that lacks the scope is an error.
func resolveMCPAgents() ([]*mcpAgent, error) {
	if flagMCPSetupScope != "user" && flagMCPSetupScope != "project" {
		return nil, fmt.Errorf(`--scope must be "user" or "project", got %q`, flagMCPSetupScope)
	}
	var out []*mcpAgent
	if flagMCPSetupTarget == "all" {
		for i := range mcpAgents {
			a := &mcpAgents[i]
			if _, err := mcpLookPath(a.bin); err == nil && a.supports(flagMCPSetupScope) {
				out = append(out, a)
			}
		}
		if len(out) == 0 {
			return nil, fmt.Errorf("no supported agent found on PATH (supported: %s)", strings.Join(mcpAgentKeys(), ", "))
		}
		return out, nil
	}
	seen := map[string]bool{}
	for _, key := range strings.Split(flagMCPSetupTarget, ",") {
		key = strings.TrimSpace(key)
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		a := findMCPAgent(key)
		if a == nil {
			return nil, fmt.Errorf("unsupported --target %q — supported: %s, all (other MCP clients: grimoire mcp config --target <client>)",
				key, strings.Join(mcpAgentKeys(), ", "))
		}
		if !a.supports(flagMCPSetupScope) {
			return nil, fmt.Errorf("%s supports --scope %s only", a.display, strings.Join(a.scopes, ", "))
		}
		out = append(out, a)
	}
	if len(out) == 0 {
		return nil, errors.New("--target is empty")
	}
	return out, nil
}

// grimoireServeCommand is the command agents should launch: plain
// "grimoire" when it resolves on PATH, else this binary's absolute path.
func grimoireServeCommand() string {
	if _, err := mcpLookPath("grimoire"); err == nil {
		return "grimoire"
	}
	if exe, err := os.Executable(); err == nil {
		if abs, err := filepath.EvalSymlinks(exe); err == nil {
			return abs
		}
		return exe
	}
	return "grimoire"
}

func runMCPInstall(_ *cobra.Command, _ []string) error {
	targets, err := resolveMCPAgents()
	if err != nil {
		return err
	}
	serve := grimoireServeCommand()

	if len(skills.AllSkillsPackages()) == 0 {
		fmt.Println("  Fetching grimoire skills…")
		if _, err := performUpdate(false, false); err != nil {
			return fmt.Errorf("fetching skills: %w", err)
		}
	}
	fmt.Printf("  %s  skills package ready\n\n", tui.IconOK)

	var failed []string
	for _, a := range targets {
		fmt.Printf("  %s\n", tui.StyleBold.Render(a.display))
		if err := installMCPAgent(a, serve); err != nil {
			fmt.Fprintf(os.Stderr, "  %s  %v\n", tui.IconWarn, err)
			failed = append(failed, a.key)
			continue
		}
		fmt.Printf("     verify: %s\n", a.verify)
	}
	fmt.Println("\n  Start a new agent session to pick up the grimoire server and instruction.")
	if len(failed) > 0 {
		return fmt.Errorf("setup incomplete for: %s", strings.Join(failed, ", "))
	}
	return nil
}

func installMCPAgent(a *mcpAgent, serve string) error {
	scope := flagMCPSetupScope
	switch {
	case a.fileRegister != nil:
		changed, err := a.fileRegister(scope, serve)
		if err != nil {
			return err
		}
		reportChange(changed, fmt.Sprintf("registered MCP server (%s scope)", scope), "MCP server already registered")
	default:
		args := a.add(scope, serve)
		manual := a.bin + " " + strings.Join(args, " ")
		if _, err := mcpLookPath(a.bin); err != nil {
			return fmt.Errorf("%s CLI not found in PATH — install it, then rerun, or register manually: %s", a.bin, manual)
		}
		out, err := mcpRunCmd(getProjectDir(), a.bin, args...)
		switch {
		case err == nil:
			fmt.Printf("  %s  registered MCP server (%s scope)\n", tui.IconOK, scope)
		case a.existsMarker != "" && strings.Contains(out, a.existsMarker):
			fmt.Printf("  %s  MCP server already registered (%s scope)\n", tui.IconSkip, scope)
		case a.key == "gemini" && strings.Contains(out, "untrusted workspace"):
			// Gemini refuses project settings until the user trusts the folder;
			// grimoire must not bypass that (GEMINI_CLI_TRUST_WORKSPACE).
			return errors.New("gemini won't change settings in an untrusted folder — open gemini here and trust the folder, then rerun (or use --scope user)")
		default:
			return fmt.Errorf("registering with %s failed: %v: %s — register manually: %s", a.bin, err, out, manual)
		}
	}

	if flagMCPSetupNoConfigure {
		return nil
	}
	path := agent.MCPInstructionFile(a.key, scope, getProjectDir())
	switch {
	case path != "":
		changed, err := agent.ConfigureMCPInstruction(path)
		if err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
		reportChange(changed, "added grimoire instruction to "+path, "grimoire instruction already in "+path)
	case a.key == "antigravity" && !flagMCPSetupNoHook:
		// No instruction file: the SessionStart hook injects the line.
	default:
		fmt.Printf("  %s  %s has no %s-level instruction file — add this line to a project AGENTS.md yourself:\n       %s\n",
			tui.IconSkip, a.display, scope, agent.MCPTriggerLine)
	}

	if !flagMCPSetupNoAllow {
		if err := installToolAllow(a.key, scope); err != nil {
			return err
		}
	}
	if flagMCPSetupNoHook {
		return nil
	}
	return installPromptHook(a.key, scope, serve)
}

// installPromptHook adds the hook that puts matching skills in front of the
// model on every prompt, in whatever form the agent supports. Agents without
// a usable prompt hook (Copilot CLI drops command-hook output) are skipped.
func installPromptHook(ag, scope, serve string) error {
	dir := getProjectDir()
	switch ag {
	case "opencode":
		p := agent.OpenCodePluginFile(scope, dir)
		changed, err := agent.ConfigureOpenCodePlugin(p, serve)
		if err != nil {
			return fmt.Errorf("adding prompt hook plugin: %w", err)
		}
		reportChange(changed, "added prompt hook plugin "+p, "prompt hook plugin already at "+p)
	case "antigravity":
		p := agent.AntigravityHooksFile(scope, dir)
		changed, err := agent.ConfigureAntigravityHooks(p, serve)
		if err != nil {
			return fmt.Errorf("adding hooks to %s: %w", p, err)
		}
		reportChange(changed, "added session + prompt hooks to "+p+" (they inject the grimoire instruction)",
			"session + prompt hooks already in "+p)
	default:
		h, ok := agent.PromptHookFor(ag)
		if !ok {
			return nil
		}
		settings := agent.PromptHookFile(ag, scope, dir)
		changed, err := agent.ConfigurePromptHook(settings, h, agent.PromptHookCommand(serve))
		if err != nil {
			return fmt.Errorf("adding prompt hook to %s: %w", settings, err)
		}
		reportChange(changed, "added prompt hook to "+settings, "prompt hook already in "+settings)
		if changed && h.Trust != "" {
			fmt.Printf("  %s  %s\n", tui.IconWarn, h.Trust)
		}
	}
	return nil
}

// removePromptHook undoes installPromptHook. Errors are reported, not fatal.
func removePromptHook(ag, scope string) {
	dir := getProjectDir()
	var (
		changed bool
		err     error
		where   string
	)
	switch ag {
	case "opencode":
		where = agent.OpenCodePluginFile(scope, dir)
		changed, err = agent.RemoveOpenCodePlugin(where)
	case "antigravity":
		where = agent.AntigravityHooksFile(scope, dir)
		changed, err = agent.RemoveAntigravityHooks(where)
	default:
		h, ok := agent.PromptHookFor(ag)
		if !ok {
			return
		}
		where = agent.PromptHookFile(ag, scope, dir)
		changed, err = agent.RemovePromptHook(where, h)
	}
	switch {
	case err != nil:
		fmt.Fprintf(os.Stderr, "  %s  %v\n", tui.IconWarn, err)
	case changed:
		fmt.Printf("  %s  removed prompt hook from %s\n", tui.IconOK, where)
	}
}

func reportChange(changed bool, did, already string) {
	if changed {
		fmt.Printf("  %s  %s\n", tui.IconOK, did)
	} else {
		fmt.Printf("  %s  %s\n", tui.IconSkip, already)
	}
}

func runMCPUninstall(_ *cobra.Command, _ []string) error {
	targets, err := resolveMCPAgents()
	if err != nil {
		return err
	}
	scope := flagMCPSetupScope
	for _, a := range targets {
		fmt.Printf("  %s\n", tui.StyleBold.Render(a.display))
		switch {
		case a.fileUnregister != nil:
			changed, err := a.fileUnregister(scope)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  %s  %v\n", tui.IconWarn, err)
			} else if changed {
				fmt.Printf("  %s  unregistered MCP server (%s scope)\n", tui.IconOK, scope)
			}
		default:
			if _, err := mcpLookPath(a.bin); err != nil {
				fmt.Printf("  %s  %s CLI not found; skipped server removal\n", tui.IconSkip, a.bin)
				break
			}
			if out, err := mcpRunCmd(getProjectDir(), a.bin, a.remove(scope)...); err != nil {
				fmt.Printf("  %s  MCP server not removed (%s scope): %s\n", tui.IconSkip, scope, out)
			} else {
				fmt.Printf("  %s  unregistered MCP server (%s scope)\n", tui.IconOK, scope)
			}
		}
		removePromptHook(a.key, scope)
		removeToolAllow(a.key, scope)
		path := agent.MCPInstructionFile(a.key, scope, getProjectDir())
		if path == "" {
			continue
		}
		changed, err := agent.RemoveMCPInstruction(path)
		if err != nil {
			return fmt.Errorf("updating %s: %w", path, err)
		}
		if changed {
			fmt.Printf("  %s  removed grimoire instruction from %s\n", tui.IconOK, path)
		}
	}
	return nil
}

// ── Config-file registration ─────────────────────────────────────────────────
// For agents without a non-interactive `mcp add`, grimoire edits their JSON
// config directly, touching only the "grimoire" entry. JSONC files are left
// alone (rewriting them would drop comments); the user gets the entry to paste.

// jsonServerFile describes one agent's MCP config file: where it lives, which
// key holds the server map, and the entry for grimoire.
type jsonServerFile struct {
	path   func(scope string) string
	mapKey string
	entry  func(serve string) map[string]any
	seed   map[string]any // initial content for a new file
}

func loadJSONConfig(path string, seed map[string]any, pasteHint string) (map[string]any, error) {
	raw, err := os.ReadFile(path)
	// An agent may create its config as an empty file (Antigravity does);
	// treat that like a missing one.
	empty := err == nil && strings.TrimSpace(string(raw)) == ""
	if os.IsNotExist(err) || empty {
		jsonc := strings.TrimSuffix(path, ".json") + ".jsonc"
		if _, statErr := os.Stat(jsonc); statErr == nil && !empty {
			return nil, fmt.Errorf("found %s — add this yourself (grimoire won't rewrite JSONC comments): %s", jsonc, pasteHint)
		}
		cfg := map[string]any{}
		for k, v := range seed {
			cfg[k] = v
		}
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}
	cfg := map[string]any{}
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return cfg, nil
}

func writeJSONConfig(path string, cfg map[string]any) error {
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

// serverMap returns cfg[key] as a map (empty when absent). A non-object value
// is an error rather than something to overwrite.
func serverMap(cfg map[string]any, key, path string) (map[string]any, error) {
	v, present := cfg[key]
	if !present {
		return map[string]any{}, nil
	}
	servers, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf(`%s: %q is not an object — fix it by hand`, path, key)
	}
	return servers, nil
}

func (f jsonServerFile) register(scope, serve string) (bool, error) {
	path := f.path(scope)
	want := f.entry(serve)
	hint, _ := json.Marshal(map[string]any{f.mapKey: map[string]any{"grimoire": want}})
	cfg, err := loadJSONConfig(path, f.seed, string(hint))
	if err != nil {
		return false, err
	}
	servers, err := serverMap(cfg, f.mapKey, path)
	if err != nil {
		return false, err
	}
	// Compare through JSON so a file-loaded entry ([]any) matches the wanted one.
	have, _ := json.Marshal(servers["grimoire"])
	wantJSON, _ := json.Marshal(want)
	if servers["grimoire"] != nil && bytes.Equal(have, wantJSON) {
		return false, nil
	}
	servers["grimoire"] = want
	cfg[f.mapKey] = servers
	return true, writeJSONConfig(path, cfg)
}

func (f jsonServerFile) unregister(scope string) (bool, error) {
	path := f.path(scope)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}
	cfg, err := loadJSONConfig(path, nil, "")
	if err != nil {
		return false, err
	}
	servers, err := serverMap(cfg, f.mapKey, path)
	if err != nil {
		return false, err
	}
	if _, ok := servers["grimoire"]; !ok {
		return false, nil
	}
	delete(servers, "grimoire")
	if len(servers) == 0 {
		delete(cfg, f.mapKey)
	}
	return true, writeJSONConfig(path, cfg)
}

func commandArgsEntry(serve string) map[string]any {
	return map[string]any{"command": serve, "args": []string{"mcp", "serve"}}
}

// mcpServersFile builds the jsonServerFile for agents using the common
// {"mcpServers": {...}} shape (paths from agent.MCPServersFile).
func mcpServersFile(ag string, entry func(string) map[string]any) jsonServerFile {
	return jsonServerFile{
		path:   func(scope string) string { return agent.MCPServersFile(ag, scope, getProjectDir()) },
		mapKey: "mcpServers",
		entry:  entry,
	}
}

// OpenCode: `opencode mcp add` is interactive, so grimoire edits opencode.json:
//
//	"mcp": { "grimoire": { "type": "local", "command": [...], "enabled": true } }
//
// Checked with `opencode mcp list` (shows grimoire connected).
var openCodeFile = jsonServerFile{
	path: func(scope string) string {
		if scope == "project" {
			return filepath.Join(getProjectDir(), "opencode.json")
		}
		return filepath.Join(agent.OpenCodeConfigDir(), "opencode.json")
	},
	mapKey: "mcp",
	entry: func(serve string) map[string]any {
		return map[string]any{"type": "local", "command": []string{serve, "mcp", "serve"}, "enabled": true}
	},
	seed: map[string]any{"$schema": "https://opencode.ai/config.json"},
}

func registerOpenCodeMCP(scope, serve string) (bool, error) {
	return openCodeFile.register(scope, serve)
}
func unregisterOpenCodeMCP(scope string) (bool, error) { return openCodeFile.unregister(scope) }

// installToolAllow pre-approves grimoire's read-only skill tools in the
// agent's permission allow list, so a skill lookup never waits on an
// approval prompt (headless runs refuse unapproved MCP tools outright).
func installToolAllow(ag, scope string) error {
	p := agent.ToolAllowFile(ag, scope, getProjectDir())
	if p == "" {
		return nil
	}
	changed, err := agent.ConfigureToolAllow(p, agent.ToolAllowRules(ag))
	if err != nil {
		return fmt.Errorf("pre-approving grimoire skill tools in %s: %w", p, err)
	}
	reportChange(changed, "pre-approved read-only grimoire skill tools in "+p,
		"read-only grimoire skill tools already approved in "+p)
	return nil
}

// removeToolAllow undoes installToolAllow. Errors are reported, not fatal.
func removeToolAllow(ag, scope string) {
	p := agent.ToolAllowFile(ag, scope, getProjectDir())
	if p == "" {
		return
	}
	changed, err := agent.RemoveToolAllow(p, agent.ToolAllowRules(ag))
	switch {
	case err != nil:
		fmt.Fprintf(os.Stderr, "  %s  %v\n", tui.IconWarn, err)
	case changed:
		fmt.Printf("  %s  removed grimoire tool approvals from %s\n", tui.IconOK, p)
	}
}
