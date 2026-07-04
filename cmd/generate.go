package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	flagGenerateDomain string
	flagGenerateForce  bool
)

var generateCmd = &cobra.Command{
	Use:   "generate <type> <name>",
	Short: "Scaffold grimoire artifacts (skill, profile)",
	Args:  cobra.ExactArgs(2),
	RunE:  runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVar(&flagGenerateDomain, "domain", "", "domain directory for the generated skill (e.g. engineering)")
	generateCmd.Flags().BoolVar(&flagGenerateForce, "force", false, "overwrite existing files")
}

func runGenerate(_ *cobra.Command, args []string) error {
	kind := strings.ToLower(args[0])
	name := args[1]

	switch kind {
	case "skill":
		return generateSkill(name)
	case "profile":
		return generateProfile(name)
	default:
		return fmt.Errorf("unknown type %q — supported: skill, profile", kind)
	}
}

func generateSkill(name string) error {
	var dir string
	if flagGenerateDomain != "" {
		dir = filepath.Join("skills", flagGenerateDomain, name)
	} else {
		dir = filepath.Join("skills", name)
	}

	skillMD := filepath.Join(dir, "SKILL.md")
	if !flagGenerateForce {
		if _, err := os.Stat(skillMD); err == nil {
			return fmt.Errorf("%s already exists — use --force to overwrite", skillMD)
		}
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}

	domain := flagGenerateDomain
	if domain == "" {
		domain = "general"
	}

	content := fmt.Sprintf(`---
name: %s
version: 0.1.0
description: ""
authors: []
license: MIT
tags:
  - %s
compatibility:
  - claude
  - opencode
dependencies: {}
---

# %s

<!-- Describe what this skill teaches the AI to do. Be direct and instructional. -->
<!-- The AI reads this full body to understand what good looks like — write freely. -->
`, name, domain, name)

	if err := os.WriteFile(skillMD, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", skillMD, err)
	}

	fmt.Printf("✓ Created %s\n", skillMD)
	fmt.Printf("  Edit the file to add skill instructions.\n")
	fmt.Printf("  Publish by adding the parent directory as a package.\n")
	return nil
}

func generateProfile(name string) error {
	dir := "profiles"
	profileTOML := filepath.Join(dir, name+".toml")
	if !flagGenerateForce {
		if _, err := os.Stat(profileTOML); err == nil {
			return fmt.Errorf("%s already exists — use --force to overwrite", profileTOML)
		}
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}

	content := fmt.Sprintf(`[package]
name = "%s"
version = "0.1.0"
description = ""
tags = []

# extends = ["engineering"]   # inherit skills from another profile

[skills]
# List skills this profile activates. Value is a semver constraint.
# apply-solid-principles = "*"
# apply-dry-principle = "*"
`, name)

	if err := os.WriteFile(profileTOML, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", profileTOML, err)
	}

	fmt.Printf("✓ Created %s\n", profileTOML)
	fmt.Printf("  Add skills under [skills] and commit to your package.\n")
	return nil
}
