package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run a health check on the grimoire installation",
	RunE:  runDoctor,
}

func runDoctor(cmd *cobra.Command, args []string) error {
	warnings, errs := 0, 0

	fmt.Println("\nGrimoire health check")

	// ── Source ──────────────────────────────────────────────────────────────────
	fmt.Println()
	fmt.Println("  Source")
	home := skills.GrimoireHome()
	if _, err := os.Stat(home); err != nil {
		fmt.Printf("    %s  grimoire not found at %s — run: grimoire update\n", tui.IconFail, home)
		errs++
	} else {
		state, err := git.CurrentState(home)
		if err != nil {
			ver := skills.GrimoireVersion()
			fmt.Printf("    %s  grimoire:   v%s\n", tui.IconOK, ver)
			fmt.Printf("    %s  git repo:   not found at %s\n", tui.IconFail, home)
			errs++
		} else {
			fmt.Printf("    %s  grimoire:   v%s (commit %s, %s)\n",
				tui.IconOK, state.Version, state.Commit, state.Date)
			fmt.Printf("    %s  location:   %s\n", tui.IconOK, home)
		}
	}

	// ── AI agents ────────────────────────────────────────────────────────────────
	fmt.Println()
	fmt.Println("  AI agents")

	for _, ag := range agent.All {
		label := agent.DisplayName(ag)
		_, err := exec.LookPath(ag)
		if err != nil {
			fmt.Printf("    %s  %-16s not found\n", tui.IconSkip, label)
			continue
		}
		ver := agent.Version(ag)
		vs := "detected"
		if ver != "" {
			vs = "v" + ver
		}
		dir := agent.SkillsDir(ag)
		if _, err := os.Stat(dir); err != nil {
			fmt.Printf("    %s  %-16s %-12s (no skills installed — run: grimoire install --target %s)\n",
				tui.IconWarn, label, vs, ag)
			warnings++
			continue
		}
		count := agent.SkillCount(ag)
		broken := agent.BrokenSymlinkCount(ag)
		if broken > 0 {
			fmt.Printf("    %s  %-16s %-12s %d skills, %d broken → run: grimoire clean --target %s\n",
				tui.IconWarn, label, vs, count, broken, ag)
			warnings++
		} else {
			fmt.Printf("    %s  %-16s %-12s %d grimoire skills\n",
				tui.IconOK, label, vs, count)
		}
		cfgFile := agent.ConfigFile(ag)
		if agent.IsConfigured(ag) {
			fmt.Printf("    %s  %-16s %-12s start-best-practice active → %s\n",
				tui.IconOK, "", "", cfgFile)
		} else {
			fmt.Printf("    %s  %-16s %-12s start-best-practice not configured → run: grimoire install --target %s\n",
				tui.IconWarn, "", "", ag)
			warnings++
		}
	}

	// Detect-only agents
	detectOnly := []struct{ label, cmd string }{
		{"Cursor", "cursor"},
		{"Windsurf", "windsurf"},
		{"Aider", "aider"},
	}
	for _, item := range detectOnly {
		if _, err := exec.LookPath(item.cmd); err == nil {
			ver := agent.Version(item.cmd)
			vs := "detected"
			if ver != "" {
				vs = "v" + ver
			}
			fmt.Printf("    %s  %-16s %s\n", tui.IconOK, item.label, vs)
		} else {
			fmt.Printf("    %s  %-16s not found\n", tui.IconSkip, item.label)
		}
	}

	// ── Config files ─────────────────────────────────────────────────────────────
	fmt.Println()
	fmt.Println("  Config")

	cwd, _ := os.Getwd()
	home2, _ := os.UserHomeDir()
	cfgPaths := []struct{ path, label string }{
		{filepath.Join(cwd, ".grimoire", "settings.local.toml"), "project personal (.grimoire/settings.local.toml)"},
		{filepath.Join(cwd, ".grimoire", "settings.toml"), "project shared (.grimoire/settings.toml)"},
		{filepath.Join(home2, ".config", "grimoire", "settings.toml"), "global (~/.config/grimoire/settings.toml)"},
	}
	hasAnyCfg := false
	for _, c := range cfgPaths {
		if _, err := os.Stat(c.path); err != nil {
			fmt.Printf("    %s  %s — not found\n", tui.IconSkip, c.label)
			continue
		}
		hasAnyCfg = true
		fmt.Printf("    %s  %s — present\n", tui.IconOK, c.label)
	}
	if !hasAnyCfg {
		fmt.Printf("    %s  no settings files found (grimoire uses defaults)\n", tui.IconSkip)
	}

	// ── Summary ──────────────────────────────────────────────────────────────────
	fmt.Println()
	if errs == 0 && warnings == 0 {
		fmt.Printf("  %s  All checks passed.\n\n", tui.IconOK)
	} else {
		parts := ""
		if errs > 0 {
			parts += fmt.Sprintf("%d error(s)", errs)
		}
		if warnings > 0 {
			if parts != "" {
				parts += ", "
			}
			parts += fmt.Sprintf("%d warning(s)", warnings)
		}
		fmt.Printf("  Summary: %s.\n\n", parts)
	}
	return nil
}
