package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/compliance"
	"github.com/jeffreytse/grimoire/internal/detect"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize .grimoire/ in the current project",
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	dir := filepath.Join(cwd, ".grimoire")
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf(".grimoire/ already exists — edit .grimoire/settings.toml directly or run /configure-grimoire")
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating .grimoire/: %w", err)
	}

	profile := detect.Profile(cwd)
	if err := writeSettings(dir, profile); err != nil {
		return err
	}
	if err := writeGitignore(dir); err != nil {
		return err
	}

	fmt.Println("✓ Grimoire initialized.")
	printInitReport(cwd)
	return nil
}

func printInitReport(projectDir string) {
	reportPath := filepath.Join(projectDir, compliance.DefaultReportPath)
	report, err := compliance.Load(reportPath)
	if err != nil {
		// No report yet — show next-step guidance
		fmt.Println()
		fmt.Println("  Next steps:")
		fmt.Printf("    1. Ask your AI to run %s\n", colorize(ansiGreen, "/check-best-practice-compliance"))
		fmt.Printf("    2. Run %s to see your compliance score\n", colorize(ansiGreen, "grimoire check"))
		return
	}

	pass := report.Threshold.Status == "pass"
	statusLabel := colorize(ansiGreen, "PASS")
	if !pass {
		statusLabel = colorize(ansiRed, "FAIL")
	}

	fmt.Println()
	fmt.Printf("  Existing report: %.1f%% — %s\n", report.Coverage.OverallPct, statusLabel)

	errors := filterBySeverity(report.Diagnostics, 1)
	warnings := filterBySeverity(report.Diagnostics, 2)
	shown := 0
	for i := range errors {
		if shown >= 3 {
			break
		}
		loc := formatLoc(&errors[i])
		fmt.Printf("    %s %s%s\n", colorize(ansiRed, "✗"), errors[i].Message, loc)
		shown++
	}
	for i := range warnings {
		if shown >= 3 {
			break
		}
		loc := formatLoc(&warnings[i])
		fmt.Printf("    %s %s%s\n", colorize(ansiYellow, "⚠"), warnings[i].Message, loc)
		shown++
	}
	total := len(errors) + len(warnings)
	if total > 3 {
		fmt.Printf("    %s\n", colorize(ansiGray, fmt.Sprintf("… and %d more — run `grimoire check` for full report", total-3)))
	}

	fmt.Println()
	fmt.Printf("  Fix first finding: ask your AI to run %s\n", colorize(ansiGreen, "/fix-best-practice-finding"))
}

func writeSettings(dir, profile string) error {
	var profileLine string
	if profile != "" {
		profileLine = fmt.Sprintf("profiles = [%q]", profile)
	} else {
		profileLine = `# profiles = ["engineering"]   # uncomment and set your profile`
	}

	content := fmt.Sprintf(`# Grimoire settings
# Docs: https://github.com/jeffreytse/grimoire/blob/main/docs/settings.md

[core]
# home = "~/.grimoire"     # override clone destination
# source = "https://..."   # override skills repository

[standards]
%s

# Domain standards example:
# [standards.engineering]
# practices = ["apply-solid-principles", "apply-kiss-principle"]
# compliance-threshold = 80
# compliance-threshold-error = 0
`, profileLine)

	path := filepath.Join(dir, "settings.toml")
	return os.WriteFile(path, []byte(content), 0o644)
}

func writeGitignore(dir string) error {
	content := "settings.local.toml\n"
	path := filepath.Join(dir, ".gitignore")
	return os.WriteFile(path, []byte(content), 0o644)
}
