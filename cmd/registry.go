package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	gitops "github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Manage grimoire skill registries",
	Long: `Add, remove, list, and update grimoire skill registries.

A registry is a git repository that follows the grimoire-skills directory layout.
Multiple registries are searched in priority order: official first, then custom registries.

  grimoire registry list              show all configured registries
  grimoire registry add <name> <url>  add a new registry and clone it
  grimoire registry remove <name>     remove a registry
  grimoire registry update [<name>]   pull latest skills from registries`,
}

var flagRegistryListJSON bool

var registryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured registries",
	RunE:  runRegistryList,
}

var registryAddCmd = &cobra.Command{
	Use:   "add <name> <url>",
	Short: "Add a new registry and clone it",
	Args:  cobra.ExactArgs(2),
	RunE:  runRegistryAdd,
}

var registryRemoveCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm"},
	Short:   "Remove a registry",
	Args:    cobra.ExactArgs(1),
	RunE:    runRegistryRemove,
}

var registryUpdateCmd = &cobra.Command{
	Use:   "update [<name>]",
	Short: "Pull latest skills from all or a specific registry",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runRegistryUpdate,
}

func init() {
	registryListCmd.Flags().BoolVar(&flagRegistryListJSON, "json", false, "output as JSON")
	registryCmd.AddCommand(registryListCmd)
	registryCmd.AddCommand(registryAddCmd)
	registryCmd.AddCommand(registryRemoveCmd)
	registryCmd.AddCommand(registryUpdateCmd)
}

type registryListEntry struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	SkillsCount int    `json:"skills_count"`
	Cloned      bool   `json:"cloned"`
	Version     string `json:"version,omitempty"`
}

func runRegistryList(cmd *cobra.Command, args []string) error {
	fs, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}

	officialURL := skills.GrimoireRepoURL()
	officialRoot := skills.SkillsRoot()
	officialCloned := false
	if _, err := os.Stat(officialRoot); err == nil {
		officialCloned = true
	}
	entries := []registryListEntry{{
		Name:        "official",
		URL:         officialURL,
		SkillsCount: countSkills(officialRoot),
		Cloned:      officialCloned,
		Version:     skills.GrimoireVersion(),
	}}

	for _, name := range sortedRegistryNames(fs.Registries) {
		rc := fs.Registries[name]
		regHome := skills.RegistryHome(name)
		regRoot := regHome + "/skills"
		cloned := false
		if _, err := os.Stat(regHome); err == nil {
			cloned = true
		}
		entries = append(entries, registryListEntry{
			Name:        name,
			URL:         rc.URL,
			SkillsCount: countSkills(regRoot),
			Cloned:      cloned,
		})
	}

	if flagRegistryListJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	}

	for _, e := range entries {
		icon := tui.IconOK
		if !e.Cloned {
			icon = tui.IconWarn
		}
		fmt.Printf("  %s  %-12s %s\n", icon, e.Name, e.URL)
		if e.Cloned {
			detail := fmt.Sprintf("%d skills", e.SkillsCount)
			if e.Version != "" {
				detail += ", version " + e.Version
			}
			fmt.Printf("         %s\n\n", detail)
		} else {
			fmt.Printf("         not cloned yet — run: grimoire registry update %s\n\n", e.Name)
		}
	}

	if len(fs.Registries) == 0 {
		fmt.Printf("  %s  no custom registries — add one with: grimoire registry add <name> <url>\n", tui.StyleDim.Render("i"))
	}
	return nil
}

func runRegistryAdd(cmd *cobra.Command, args []string) error {
	name, url := args[0], args[1]

	if name == skills.OfficialRegistryName {
		return fmt.Errorf("cannot add registry named %q — use: grimoire config set source <url>", skills.OfficialRegistryName)
	}
	if !skills.IsGitURL(url) {
		return fmt.Errorf("url must be a git URL (https://, git://, git@): %s", url)
	}

	fs, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}
	if fs.Registries == nil {
		fs.Registries = make(map[string]settings.RegistryConfig)
	}
	if _, exists := fs.Registries[name]; exists {
		return fmt.Errorf("registry %q already exists — remove it first with: grimoire registry remove %s", name, name)
	}

	dest := skills.RegistryHome(name)
	fmt.Printf("Cloning %s → %s...\n", url, dest)
	if err := gitops.Clone(url, dest); err != nil {
		return fmt.Errorf("cloning registry: %w", err)
	}

	fs.Registries[name] = settings.RegistryConfig{URL: url}
	if err := settings.SaveGlobal(fs); err != nil {
		return fmt.Errorf("saving settings: %w", err)
	}

	count := countSkills(dest + "/skills")
	fmt.Printf("%s  registry %q added — %d skills\n", tui.IconOK, name, count)
	return nil
}

func runRegistryRemove(cmd *cobra.Command, args []string) error {
	name := args[0]

	if name == skills.OfficialRegistryName {
		return fmt.Errorf("cannot remove the official registry")
	}

	fs, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}
	if _, exists := fs.Registries[name]; !exists {
		return fmt.Errorf("registry %q not found", name)
	}

	delete(fs.Registries, name)
	if err := settings.SaveGlobal(fs); err != nil {
		return fmt.Errorf("saving settings: %w", err)
	}

	regHome := skills.RegistryHome(name)
	if _, err := os.Stat(regHome); err == nil {
		chosen, ok := tui.RunSelect(fmt.Sprintf("Delete local clone at %s?", regHome), []string{"Yes", "No"})
		if ok && chosen == "Yes" {
			if err := os.RemoveAll(regHome); err != nil {
				fmt.Fprintf(os.Stderr, "  warn: deleting %s: %v\n", regHome, err)
			} else {
				fmt.Printf("%s  deleted %s\n", tui.IconOK, regHome)
			}
		}
	}

	fmt.Printf("%s  registry %q removed\n", tui.IconOK, name)
	return nil
}

func runRegistryUpdate(cmd *cobra.Command, args []string) error {
	fs, err := settings.LoadGlobal()
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}

	if len(args) == 1 {
		return updateOneRegistry(args[0], fs)
	}

	// update all
	names := sortedRegistryNames(fs.Registries)
	if err := updateOneRegistry(skills.OfficialRegistryName, fs); err != nil {
		fmt.Fprintf(os.Stderr, "  warn: official: %v\n", err)
	}
	for _, name := range names {
		if err := updateOneRegistry(name, fs); err != nil {
			fmt.Fprintf(os.Stderr, "  warn: %s: %v\n", name, err)
		}
	}
	return nil
}

func updateOneRegistry(name string, fs settings.FileSettings) error {
	var url string
	if name == skills.OfficialRegistryName {
		url = skills.GrimoireRepoURL()
	} else {
		rc, ok := fs.Registries[name]
		if !ok {
			return fmt.Errorf("registry %q not configured", name)
		}
		url = rc.URL
	}

	dest := skills.RegistryHome(name)
	if _, err := os.Stat(dest); err != nil {
		fmt.Printf("  cloning %s → %s...\n", name, dest)
		if err := gitops.Clone(url, dest); err != nil {
			return fmt.Errorf("cloning: %w", err)
		}
		fmt.Printf("  %s  %s cloned\n", tui.IconOK, name)
		return nil
	}

	fmt.Printf("  pulling %s...\n", name)
	if err := gitops.Pull(dest); err != nil {
		return fmt.Errorf("pulling: %w", err)
	}
	fmt.Printf("  %s  %s up to date\n", tui.IconOK, name)
	return nil
}

func countSkills(skillsRoot string) int {
	all, _ := skills.ListAllSkills(skillsRoot)
	return len(all)
}

func sortedRegistryNames(registries map[string]settings.RegistryConfig) []string {
	names := make([]string, 0, len(registries))
	for name := range registries {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
