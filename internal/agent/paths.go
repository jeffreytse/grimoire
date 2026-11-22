package agent

import (
	"fmt"
	"os"
	"path/filepath"
)

// All lists agents that receive installed skills (have a managed skills directory).
var All = []string{"claude", "codex", "gemini", "antigravity", "openclaw", "opencode", "copilot", "pi", "dsh", "kimi", "omp", "grok", "muse", "minimax", "agent"}

// CheckAgents is the default resolution order for `grimoire check --independent`.
// Separate from All — "copilot" here maps to the "gh copilot" extension, which
// performs compliance checks; installs target the standalone Copilot CLI.
var CheckAgents = []string{"claude", "gemini", "codex", "antigravity", "copilot", "opencode", "openclaw"}

// envDir returns $name when set and non-empty, else fallback.
func envDir(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}

// copilotHome is the GitHub Copilot CLI config root ($COPILOT_HOME or ~/.copilot).
func copilotHome(home string) string {
	return envDir("COPILOT_HOME", filepath.Join(home, ".copilot"))
}

// piAgentDir is Pi's user config root ($PI_CODING_AGENT_DIR or ~/.pi/agent).
func piAgentDir(home string) string {
	return envDir("PI_CODING_AGENT_DIR", filepath.Join(home, ".pi", "agent"))
}

// kimiHome is Kimi Code's data root ($KIMI_CODE_HOME or ~/.kimi-code).
func kimiHome(home string) string {
	return envDir("KIMI_CODE_HOME", filepath.Join(home, ".kimi-code"))
}

// ompAgentDir is OMP's active agent dir: $PI_CODING_AGENT_DIR, else the
// profile dir for $OMP_PROFILE (falling back to $PI_PROFILE only when
// OMP_PROFILE is unset), else ~/.omp/agent.
func ompAgentDir(home string) string {
	if v := os.Getenv("PI_CODING_AGENT_DIR"); v != "" {
		return v
	}
	profile, ok := os.LookupEnv("OMP_PROFILE")
	if !ok {
		profile = os.Getenv("PI_PROFILE")
	}
	if profile != "" {
		return filepath.Join(home, ".omp", "profiles", profile, "agent")
	}
	return filepath.Join(home, ".omp", "agent")
}

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
	case "copilot":
		return filepath.Join(copilotHome(home), "skills")
	case "pi":
		return filepath.Join(piAgentDir(home), "skills")
	case "kimi":
		return filepath.Join(kimiHome(home), "skills")
	case "omp":
		return filepath.Join(ompAgentDir(home), "skills")
	case "dsh", "agent":
		return filepath.Join(home, ".agents", "skills")
	case "grok", "muse", "minimax":
		return filepath.Join(PluginBundleDir(ag), "skills")
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
		return filepath.Join(projectDir, ".gemini", "antigravity-cli", "skills")
	case "openclaw":
		return filepath.Join(projectDir, ".openclaw", "skills")
	case "opencode":
		return filepath.Join(projectDir, ".opencode", "skills")
	case "copilot":
		return filepath.Join(projectDir, ".github", "skills")
	case "pi":
		return filepath.Join(projectDir, ".pi", "skills")
	case "kimi":
		return filepath.Join(projectDir, ".kimi-code", "skills")
	case "omp":
		return filepath.Join(projectDir, ".omp", "skills")
	case "dsh", "agent":
		return filepath.Join(projectDir, ".agents", "skills")
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
	case "copilot":
		return filepath.Join(copilotHome(home), "copilot-instructions.md")
	case "pi":
		return filepath.Join(piAgentDir(home), "AGENTS.md")
	case "kimi":
		return filepath.Join(kimiHome(home), "AGENTS.md")
	case "omp":
		return filepath.Join(ompAgentDir(home), "AGENTS.md")
	case "dsh", "agent":
		return filepath.Join(home, ".agents", "AGENTS.md")
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
	case "copilot":
		return copilotHome(home)
	case "pi":
		return piAgentDir(home)
	case "kimi":
		return kimiHome(home)
	case "omp":
		return ompAgentDir(home)
	case "dsh", "agent":
		return filepath.Join(home, ".agents")
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
	case "copilot":
		return "GitHub Copilot CLI"
	case "pi":
		return "Pi"
	case "dsh":
		return "DeepSeek Harness"
	case "kimi":
		return "Kimi Code"
	case "omp":
		return "OMP"
	case "grok":
		return "Grok Build"
	case "muse":
		return "Muse"
	case "minimax":
		return "MiniMax Code"
	case "agent":
		return "Generic Agent"
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
	case "GitHub Copilot CLI":
		return "copilot"
	case "Pi":
		return "pi"
	case "DeepSeek Harness":
		return "dsh"
	case "Kimi Code":
		return "kimi"
	case "OMP":
		return "omp"
	case "Grok Build":
		return "grok"
	case "Muse":
		return "muse"
	case "MiniMax Code":
		return "minimax"
	case "Generic Agent":
		return "agent"
	}
	return display
}
