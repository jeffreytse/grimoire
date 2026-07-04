package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var infoCmd = &cobra.Command{
	Use:   "info <skill>",
	Short: "Show metadata for an installed skill",
	Args:  cobra.ExactArgs(1),
	RunE:  runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func runInfo(_ *cobra.Command, args []string) error {
	skillName := args[0]
	regs := skills.AllSkillsPackages()
	if len(regs) == 0 {
		return fmt.Errorf("no packages installed — run: grimoire update")
	}

	_, src, err := resolveSkillFromPackages(regs, skillName)
	if err != nil {
		return err
	}

	// WalkSkills to find the matching skill with its full metadata
	all, walkErr := skills.WalkSkills(src.Root)
	if walkErr != nil {
		return walkErr
	}

	var found *skills.Skill
	for i := range all {
		if all[i].Name == skillName {
			found = &all[i]
			break
		}
	}
	if found == nil {
		return fmt.Errorf("skill %q not found", skillName)
	}

	printSkillInfo(found)
	return nil
}

func printSkillInfo(sk *skills.Skill) {
	fmt.Printf("%s", tui.StyleBold.Render(sk.Name))
	if sk.Version != "" {
		fmt.Printf("  v%s", sk.Version)
	}
	fmt.Println()

	if sk.Description != "" {
		fmt.Printf("  %s\n", sk.Description)
	}
	fmt.Println()

	if sk.Domain != "" {
		domain := sk.Domain
		if sk.Subdomain != "" {
			domain += "/" + sk.Subdomain
		}
		fmt.Printf("  domain:        %s\n", domain)
	}
	if sk.Package != "" {
		fmt.Printf("  package:       %s\n", sk.Package)
	}
	if len(sk.Authors) > 0 {
		fmt.Printf("  authors:       %s\n", strings.Join(sk.Authors, ", "))
	}
	if sk.License != "" {
		fmt.Printf("  license:       %s\n", sk.License)
	}
	if len(sk.Compatibility) > 0 {
		fmt.Printf("  compatibility: %s\n", strings.Join(sk.Compatibility, ", "))
	}
	if len(sk.Tags) > 0 {
		fmt.Printf("  tags:          %s\n", strings.Join(sk.Tags, ", "))
	}
	if len(sk.Dependencies) > 0 {
		fmt.Printf("  dependencies:\n")
		for dep, constraint := range sk.Dependencies {
			fmt.Printf("    %s = %q\n", dep, constraint)
		}
	}
	if sk.Path != "" {
		fmt.Printf("  path:          %s\n", sk.Path)
	}
	fmt.Println()
}
