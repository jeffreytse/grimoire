package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var flagVersionJSON bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show grimoire version information",
	RunE:  runVersion,
}

func init() {
	versionCmd.Flags().BoolVar(&flagVersionJSON, "json", false, "output as JSON")
}

func runVersion(cmd *cobra.Command, args []string) error {
	home := skills.OfficialPackageHome()

	type versionOut struct {
		CLI      string `json:"cli"`
		Grimoire string `json:"grimoire,omitempty"`
		Home     string `json:"home"`
	}

	out := versionOut{
		CLI:  strings.TrimPrefix(cliVersion, "v"),
		Home: home,
	}

	if _, err := os.Stat(home); err == nil {
		if state, err := git.CurrentState(home); err == nil {
			out.Grimoire = state.Version
		} else {
			out.Grimoire = skills.GrimoireVersion()
		}
	}

	if flagVersionJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	fmt.Printf("%s  cli:        %s\n", tui.IconOK, cliVersion)
	if out.Grimoire == "" {
		fmt.Printf("%s  grimoire:   not installed (run: grimoire update)\n", tui.IconWarn)
		return nil
	}
	fmt.Printf("%s  grimoire:   %s\n", tui.IconOK, out.Grimoire)
	fmt.Printf("%s  location:   %s\n", tui.IconOK, home)
	return nil
}
