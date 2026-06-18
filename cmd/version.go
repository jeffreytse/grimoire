package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show grimoire version information",
	RunE:  runVersion,
}

func runVersion(cmd *cobra.Command, args []string) error {
	fmt.Printf("%s  cli:        %s\n", tui.IconOK, cliVersion)

	home := skills.GrimoireHome()
	if _, err := os.Stat(home); err != nil {
		fmt.Printf("%s  grimoire:   not installed (run: grimoire update)\n", tui.IconWarn)
		return nil
	}

	state, err := git.CurrentState(home)
	if err != nil {
		ver := skills.GrimoireVersion()
		fmt.Printf("%s  grimoire:   v%s\n", tui.IconOK, ver)
		return nil
	}

	fmt.Printf("%s  grimoire:   v%s (commit %s, %s)\n",
		tui.IconOK, state.Version, state.Commit, state.Date)
	fmt.Printf("%s  location:   %s\n",
		tui.IconOK, home)
	return nil
}
