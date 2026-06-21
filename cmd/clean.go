package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var flagCleanTarget string

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove stale grimoire-managed skills from agent directories",
	RunE:  runClean,
}

func init() {
	cleanCmd.Flags().StringVar(&flagCleanTarget, "target", "", "agent to clean (default: all detected)")
}

func runClean(cmd *cobra.Command, args []string) error {
	targets := resolveTargets(flagCleanTarget)
	total := 0
	for _, ag := range targets {
		dir := agent.SkillsDir(ag)
		n, err := skills.CleanBrokenSymlinks(dir)
		if err != nil {
			fmt.Printf("  %s  %s: %v\n", tui.IconWarn, agent.DisplayName(ag), err)
			continue
		}
		total += n
	}
	if total == 0 {
		fmt.Println("  nothing to clean")
	} else {
		fmt.Printf("  %s  %d broken symlink(s) removed\n", tui.IconOK, total)
	}
	return nil
}

func resolveTargets(target string) []string {
	switch target {
	case "", "auto":
		// Pinned agents in settings take priority over auto-detection
		r, _ := settings.Load(getProjectDir())
		if len(r.Core.Agents) > 0 {
			return r.Core.Agents
		}
		detected := agent.Detected()
		if len(detected) == 0 {
			return []string{"claude"}
		}
		return detected
	case "all":
		return agent.All
	default:
		return []string{target}
	}
}
