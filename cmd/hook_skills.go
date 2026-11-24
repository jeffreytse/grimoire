package cmd

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Skills-mode hooks, registered by `grimoire install`. There the skills are
// linked into the agent and its instruction file already runs
// start-best-practice, so these hooks carry only what MCP mode gets from its
// server: grimoire's display format, the per-prompt reminder, status lines
// for the user, and (Claude Code) the skill in use for the status bar badge.

// skillsSessionHead tells the model how grimoire runs in skills mode.
const skillsSessionHead = "grimoire skills are installed in this agent. Before responding to each new user request, invoke the start-best-practice skill and follow the workflow it routes to; invoke the skills it names the way this agent loads skills."

// skillsRouteReminder is skills mode's per-prompt nudge onto the router.
const skillsRouteReminder = "grimoire: before responding, invoke the start-best-practice skill and follow the workflow it routes to: it asks clarifying questions when needed and loads the specific best practices. If nothing routes, proceed normally."

var hookSkillsSessionCmd = &cobra.Command{
	Use:           "skills-session",
	Short:         "Skills-mode session hook: grimoire's display format",
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		runSkillsSession(os.Stdin, os.Stdout)
		return nil
	},
}

var hookSkillsPromptCmd = &cobra.Command{
	Use:           "skills-prompt",
	Short:         "Skills-mode prompt hook: routing reminder and status line",
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		runSkillsPrompt(os.Stdin, os.Stdout)
		return nil
	},
}

var hookSkillUsedCmd = &cobra.Command{
	Use:           "skill-used",
	Short:         "Claude Code PostToolUse hook on Skill: record the skill for the status line",
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		runSkillUsed(os.Stdin)
		return nil
	},
}

func init() {
	hookCmd.AddCommand(hookSkillsSessionCmd, hookSkillsPromptCmd, hookSkillUsedCmd)
}

func runSkillsSession(in io.Reader, out io.Writer) {
	raw, err := io.ReadAll(io.LimitReader(in, 1<<20))
	if err != nil {
		return
	}
	var input hookInput
	_ = json.Unmarshal(raw, &input)
	l := displayLevelSetting()
	if active, src := grimoireActive(input.SessionID); !active {
		writeHookOutput(out, input.HookEventName, "SessionStart", inactiveRule, inactiveSessionStatus(l, src))
		return
	}
	ctx, msg := skillsSessionTurn(l)
	writeHookOutput(out, input.HookEventName, "SessionStart", ctx, msg)
}

// skillsSessionTurn returns the session context and status line in skills mode.
func skillsSessionTurn(l displayLevel) (ctx, msg string) {
	return skillsSessionHead + "\n\n" + displayContext(l),
		readyStatus(l, "skills mode", "")
}

func runSkillsPrompt(in io.Reader, out io.Writer) {
	raw, err := io.ReadAll(io.LimitReader(in, 1<<20))
	if err != nil {
		return
	}
	var input hookInput
	if json.Unmarshal(raw, &input) != nil {
		return
	}
	ctx, msg := promptTurn(input.Prompt, input.SessionID, displayLevelSetting(), true, skillsHookTurn)
	writeHookOutput(out, input.HookEventName, "UserPromptSubmit", ctx, msg)
}

// skillsHookTurn is hookTurn for skills mode: the skills are already in the
// agent, so it neither inlines skills nor names MCP tools.
func skillsHookTurn(prompt string, l displayLevel) (ctx, msg string) {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" || strings.HasPrefix(prompt, "/") || isSmallTalk(prompt) {
		return "", ""
	}
	reminder := skillsRouteReminder + " " + displayContext(l)
	if cands, names := suggestedSkillContext(prompt); cands != "" {
		return reminder + "\n" + cands, candidatesStatus(l, names)
	}
	return reminder, ""
}

// skillUsedInput is the part of Claude Code's PostToolUse input we read.
type skillUsedInput struct {
	SessionID string `json:"session_id"`
	ToolName  string `json:"tool_name"`
	ToolInput struct {
		Skill string   `json:"skill"`
		Name  string   `json:"name"`
		Names []string `json:"names"`
	} `json:"tool_input"`
}

// mcpGetSkillTool is grimoire_get_skill as Claude Code names it.
const mcpGetSkillTool = "mcp__grimoire__grimoire_get_skill"

// runSkillUsed records, per session, the skill Claude Code just loaded:
// through its Skill tool (skills mode) or grimoire_get_skill (MCP mode).
// Plugin skills come as "plugin:name"; only the name is kept.
func runSkillUsed(in io.Reader) {
	raw, err := io.ReadAll(io.LimitReader(in, 1<<20))
	if err != nil {
		return
	}
	var input skillUsedInput
	if json.Unmarshal(raw, &input) != nil {
		return
	}
	var names []string
	switch input.ToolName {
	case "Skill":
		names = []string{input.ToolInput.Skill}
	case mcpGetSkillTool:
		names = append([]string{input.ToolInput.Name}, input.ToolInput.Names...)
	default:
		return
	}
	root := grimoireStateRoot()
	for _, name := range names {
		name = strings.TrimPrefix(strings.TrimSpace(name), "grimoire:")
		if i := strings.LastIndex(name, ":"); i >= 0 {
			name = name[i+1:]
		}
		recordSkill(root, input.SessionID, name)
	}
}
