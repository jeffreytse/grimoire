package cmd

import (
	"fmt"
	"os"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/tui"
)

// Skills-mode agent setup. Besides the start-best-practice instruction line,
// `grimoire install` registers the skills-mode hooks (grimoire's display
// format, routing reminder and status lines) and, for Claude Code, the
// [GRIMOIRE] status line badge, at user scope like the instruction line.

var (
	flagInstallNoHook     bool
	flagInstallStatusLine bool
	flagInstallNoStatus   bool
)

func init() {
	installCmd.Flags().BoolVar(&flagInstallNoHook, "no-hook", false, "do not add grimoire's display and prompt hooks (Claude Code, Codex, Gemini CLI)")
	installCmd.Flags().BoolVar(&flagInstallStatusLine, "statusline", false, "chain an existing Claude Code status line so the [GRIMOIRE] badge is shown after it")
	installCmd.Flags().BoolVar(&flagInstallNoStatus, "no-statusline", false, "do not add the [GRIMOIRE] badge to the Claude Code status line")
}

// configureSkillsAgent writes ag's instruction line and, unless turned off,
// its skills-mode hooks and status line badge. Failures are warnings: the
// skills are installed either way.
func configureSkillsAgent(ag string) {
	if err := agent.ConfigureAgentMD(ag); err != nil {
		fmt.Fprintf(os.Stderr, "  warn: configuring %s: %v\n", ag, err)
	}
	if _, err := os.Stat(agent.ConfigDir(ag)); err != nil {
		return // agent not set up on this machine
	}
	bin := grimoireServeCommand()
	if !flagInstallNoHook {
		if err := installSkillsHooks(ag, bin); err != nil {
			fmt.Fprintf(os.Stderr, "  warn: %s hooks: %v\n", ag, err)
		}
	}
	if !flagInstallNoStatus {
		if err := installStatusLineFor(ag, "user", bin, flagInstallStatusLine); err != nil {
			fmt.Fprintf(os.Stderr, "  warn: %v\n", err)
		}
	}
}

// installSkillsHooks registers the skills-mode session and prompt hooks and,
// for Claude Code, the skill-used hook. An MCP-mode setup's hooks on the same
// events are replaced: one grimoire mode per agent.
func installSkillsHooks(ag, bin string) error {
	h, ok := agent.PromptHookFor(ag)
	if !ok {
		return nil
	}
	settings := agent.PromptHookFile(ag, "user", getProjectDir())
	changed, err := agent.ConfigurePromptHook(settings, h, agent.SkillsPromptHookCommand(bin))
	if err != nil {
		return err
	}
	sc, err := agent.ConfigureSessionHook(settings, h, agent.SkillsSessionHookCommand(bin))
	if err != nil {
		return err
	}
	uc, err := agent.ConfigureSkillUsedHook(settings, ag, agent.SkillUsedHookCommand(bin))
	if err != nil {
		return err
	}
	if changed || sc || uc {
		fmt.Printf("  %s  added grimoire display hooks to %s\n", tui.IconOK, settings)
		if h.Trust != "" {
			fmt.Printf("  %s  %s\n", tui.IconWarn, h.Trust)
		}
	}
	return nil
}

// removeSkillsHooks undoes installSkillsHooks, leaving MCP-mode hooks alone.
func removeSkillsHooks(ag string) {
	h, ok := agent.PromptHookFor(ag)
	if !ok {
		return
	}
	settings := agent.PromptHookFile(ag, "user", getProjectDir())
	changed, err := agent.RemoveSkillsHooks(settings, h)
	switch {
	case err != nil:
		fmt.Fprintf(os.Stderr, "  %s  %v\n", tui.IconWarn, err)
	case changed:
		fmt.Printf("  %s  removed grimoire display hooks from %s\n", tui.IconOK, settings)
	}
}
