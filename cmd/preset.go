package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var presetCmd = &cobra.Command{
	Use:   "preset",
	Short: "Manage grimoire presets",
	Long: `Presets are starter configurations for common project types.

  grimoire preset list   List available presets from all installed registries`,
}

var flagPresetListJSON bool

var presetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available presets from all installed registries",
	Args:  cobra.NoArgs,
	RunE:  runPresetList,
}

func init() {
	presetListCmd.Flags().BoolVar(&flagPresetListJSON, "json", false, "output as JSON")
	presetCmd.AddCommand(presetListCmd)
}

type presetListEntry struct {
	Name     string `json:"name"`
	Registry string `json:"registry"`
}

func runPresetList(cmd *cobra.Command, args []string) error {
	regs := skills.AllRegistries()

	var entries []presetListEntry
	for _, reg := range regs {
		for _, name := range skills.ListPresets(reg.Home) {
			entries = append(entries, presetListEntry{Name: name, Registry: reg.Name})
		}
	}

	if flagPresetListJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	}

	if len(entries) == 0 {
		fmt.Println("No presets found — run: grimoire update")
		return nil
	}

	for _, e := range entries {
		tag := tui.StyleDim.Render("[" + e.Registry + "]")
		fmt.Printf("  %s  %-30s %s\n", tui.IconOK, e.Name, tag)
	}
	fmt.Printf("\n  apply: grimoire init --preset <name>\n\n")
	return nil
}
