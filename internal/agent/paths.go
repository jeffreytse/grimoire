package agent

import (
	"fmt"
	"os"
	"path/filepath"
)

// All lists agents that receive installed skills (have a managed skills directory).
var All = []string{"claude", "codex", "gemini", "antigravity", "openclaw", "opencode"}

// CheckAgents is the default resolution order for `grimoire check --independent`.
// Separate from All — includes "copilot" (maps to "gh copilot"), which cannot receive
// skills but can perform compliance checks.
var CheckAgents = []string{"claude", "gemini", "codex", "antigravity", "copilot", "opencode", "openclaw"}

func mustHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("grimoire: cannot determine home directory: %v\nSet $HOME or run as a user with a home directory.", err))
	}
	return home
}

func SkillsDir(ag string) string {
	home := mustHome()
	switch ag {
	case "claude":
		return filepath.Join(home, ".claude", "skills")
	case "codex":
		return filepath.Join(home, ".agents", "skills")
	case "gemini":
		return filepath.Join(home, ".gemini", "skills")
	case "antigravity":
		return filepath.Join(home, ".gemini", "antigravity-cli", "skills")
	case "openclaw":
		return filepath.Join(home, ".openclaw", "skills")
	case "opencode":
		return filepath.Join(home, ".config", "opencode", "skills")
	}
	return ""
}

// ProjectSkillsDir returns the project-relative skills directory for an agent.
// Mirrors SkillsDir but rooted at projectDir instead of the user home.
func ProjectSkillsDir(ag, projectDir string) string {
	switch ag {
	case "claude":
		return filepath.Join(projectDir, ".claude", "skills")
	case "codex":
		return filepath.Join(projectDir, ".agents", "skills")
	case "gemini":
		return filepath.Join(projectDir, ".gemini", "skills")
	case "antigravity":
		return filepath.Join(projectDir, ".agents", "skills")
	case "openclaw":
		return filepath.Join(projectDir, ".openclaw", "skills")
	case "opencode":
		return filepath.Join(projectDir, ".opencode", "skills")
	}
	return ""
}

func ConfigFile(ag string) string {
	home := mustHome()
	switch ag {
	case "claude":
		return filepath.Join(home, ".claude", "CLAUDE.md")
	case "codex":
		return filepath.Join(home, ".agents", "AGENTS.md")
	case "gemini":
		return filepath.Join(home, ".gemini", "GEMINI.md")
	case "antigravity":
		return filepath.Join(home, ".gemini", "AGENTS.md")
	case "openclaw":
		return filepath.Join(home, ".openclaw", "workspace", "AGENTS.md")
	case "opencode":
		return filepath.Join(home, ".config", "opencode", "AGENTS.md")
	}
	return ""
}

func ConfigDir(ag string) string {
	home := mustHome()
	switch ag {
	case "claude":
		return filepath.Join(home, ".claude")
	case "codex":
		return filepath.Join(home, ".agents")
	case "gemini":
		return filepath.Join(home, ".gemini")
	case "antigravity":
		return filepath.Join(home, ".gemini")
	case "openclaw":
		return filepath.Join(home, ".openclaw", "workspace")
	case "opencode":
		return filepath.Join(home, ".config", "opencode")
	}
	return ""
}

func DisplayName(ag string) string {
	switch ag {
	case "claude":
		return "Claude Code"
	case "codex":
		return "Codex"
	case "gemini":
		return "Gemini CLI"
	case "antigravity":
		return "Antigravity"
	case "openclaw":
		return "OpenClaw"
	case "opencode":
		return "OpenCode"
	}
	return ag
}

func FromDisplayName(display string) string {
	switch display {
	case "Claude Code":
		return "claude"
	case "Codex":
		return "codex"
	case "Gemini CLI":
		return "gemini"
	case "Antigravity":
		return "antigravity"
	case "OpenClaw":
		return "openclaw"
	case "OpenCode":
		return "opencode"
	}
	return display
}
