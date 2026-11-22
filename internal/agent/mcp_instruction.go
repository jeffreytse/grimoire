package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MCPTriggerLine tells the agent to consult grimoire through its MCP tools.
// It replaces triggerLine (which invokes the start-best-practice skill) when
// grimoire is used in MCP mode, where no skills are linked into the agent.
//
// It must live in the agent's instruction file (e.g. CLAUDE.md): in testing,
// Claude Code followed it from CLAUDE.md but ignored the same text passed via
// --append-system-prompt.
const MCPTriggerLine = "For any task, first use the grimoire MCP tools: `grimoire_search_skills` to find a relevant best-practice skill, then `grimoire_get_skill` to load and follow it."

// ConfigureMCPInstruction makes path contain the grimoire section with the MCP
// trigger, replacing a skills-mode trigger if present. It creates the file
// (and parent dirs) when missing. Reports whether the file changed.
func ConfigureMCPInstruction(path string) (bool, error) {
	raw, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}
	content := string(raw)
	if strings.Contains(content, MCPTriggerLine) && !strings.Contains(content, triggerLine) {
		return false, nil
	}

	lines := stripGrimoireLines(content)
	body := strings.TrimRight(strings.Join(lines, "\n"), "\n")
	if body != "" {
		body += "\n\n"
	}
	body += sectionHeader + "\n" + MCPTriggerLine + "\n"

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, fmt.Errorf("creating %s: %w", filepath.Dir(path), err)
	}
	return true, os.WriteFile(path, []byte(body), 0o644)
}

// RemoveMCPInstruction removes the grimoire section (MCP or skills trigger)
// from path. A missing file is not an error. Reports whether the file changed.
func RemoveMCPInstruction(path string) (bool, error) {
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	content := string(raw)
	if !strings.Contains(content, MCPTriggerLine) && !strings.Contains(content, triggerLine) {
		return false, nil
	}
	body := strings.TrimRight(strings.Join(stripGrimoireLines(content), "\n"), "\n")
	if body != "" {
		body += "\n"
	}
	return true, os.WriteFile(path, []byte(body), 0o644)
}

// stripGrimoireLines drops grimoire's header and trigger lines, keeping all
// other content in order.
func stripGrimoireLines(content string) []string {
	var out []string
	for _, l := range strings.Split(content, "\n") {
		switch strings.TrimSpace(l) {
		case sectionHeader, triggerLine, MCPTriggerLine:
			continue
		}
		out = append(out, l)
	}
	return out
}

// MCPInstructionFile returns the instruction file an agent reads for the
// given scope ("user" or "project"), or "" when that scope has none.
// Paths follow each agent's documented location:
//
//	claude   ~/.claude/CLAUDE.md                         ./CLAUDE.md
//	codex    ${CODEX_HOME:-~/.codex}/AGENTS.md            (user only)
//	gemini   ~/.gemini/GEMINI.md                          ./GEMINI.md
//	copilot  ${COPILOT_HOME:-~/.copilot}/copilot-instructions.md
//	                                                      ./.github/copilot-instructions.md
//	opencode ${XDG_CONFIG_HOME:-~/.config}/opencode/AGENTS.md  ./AGENTS.md
//	pi       <pi agent dir>/AGENTS.md                   ./AGENTS.md
//	kimi     ${KIMI_CODE_HOME:-~/.kimi-code}/AGENTS.md   ./AGENTS.md
//	omp      <omp agent dir>/AGENTS.md                  ./.omp/AGENTS.md
func MCPInstructionFile(ag, scope, projectDir string) string {
	home := mustHome()
	if scope == "project" {
		switch ag {
		case "claude":
			return filepath.Join(projectDir, "CLAUDE.md")
		case "gemini":
			return filepath.Join(projectDir, "GEMINI.md")
		case "copilot":
			return filepath.Join(projectDir, ".github", "copilot-instructions.md")
		case "opencode", "pi", "kimi":
			return filepath.Join(projectDir, "AGENTS.md")
		case "omp":
			return filepath.Join(projectDir, ".omp", "AGENTS.md")
		}
		return ""
	}
	switch ag {
	case "claude":
		return filepath.Join(home, ".claude", "CLAUDE.md")
	case "codex":
		return filepath.Join(envDir("CODEX_HOME", filepath.Join(home, ".codex")), "AGENTS.md")
	case "gemini":
		return filepath.Join(home, ".gemini", "GEMINI.md")
	case "copilot":
		return filepath.Join(copilotHome(home), "copilot-instructions.md")
	case "opencode":
		return filepath.Join(OpenCodeConfigDir(), "AGENTS.md")
	case "pi":
		return filepath.Join(piAgentDir(home), "AGENTS.md")
	case "kimi":
		return filepath.Join(kimiHome(home), "AGENTS.md")
	case "omp":
		return filepath.Join(ompAgentDir(home), "AGENTS.md")
	}
	// antigravity: agy reads GEMINI.md/AGENTS.md only from the working
	// directory up to the repo root — there is no user-level file.
	return ""
}

// MCPServersFile returns the "mcpServers" JSON config an agent reads for the
// scope, or "" when the agent has none for it. All use the common shape
// {"mcpServers": {"<name>": {"command": ..., "args": [...]}}}.
//
//	pi           <pi agent dir>/mcp.json          ./.pi/mcp.json
//	kimi         ${KIMI_CODE_HOME:-~/.kimi-code}/mcp.json   ./.kimi-code/mcp.json
//	omp          <omp agent dir>/mcp.json         ./.omp/mcp.json
//	antigravity  ~/.gemini/config/mcp_config.json (user only)
func MCPServersFile(ag, scope, projectDir string) string {
	home := mustHome()
	if scope == "project" {
		switch ag {
		case "pi":
			return filepath.Join(projectDir, ".pi", "mcp.json")
		case "kimi":
			return filepath.Join(projectDir, ".kimi-code", "mcp.json")
		case "omp":
			return filepath.Join(projectDir, ".omp", "mcp.json")
		}
		return ""
	}
	switch ag {
	case "pi":
		return filepath.Join(piAgentDir(home), "mcp.json")
	case "kimi":
		return filepath.Join(kimiHome(home), "mcp.json")
	case "omp":
		return filepath.Join(ompAgentDir(home), "mcp.json")
	case "antigravity":
		return filepath.Join(home, ".gemini", "config", "mcp_config.json")
	}
	return ""
}

// OpenCodeConfigDir is OpenCode's global config dir
// (${XDG_CONFIG_HOME:-~/.config}/opencode).
func OpenCodeConfigDir() string {
	return filepath.Join(envDir("XDG_CONFIG_HOME", filepath.Join(mustHome(), ".config")), "opencode")
}
