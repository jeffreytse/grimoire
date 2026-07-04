package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/compliance"
	"github.com/jeffreytse/grimoire/internal/config"
	grimctx "github.com/jeffreytse/grimoire/internal/context"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show project compliance health at a glance",
	Long: `Show the current compliance status of the project: active profile,
installed skills, last check result, and staleness.`,
	RunE: runStatus,
}

func runStatus(_ *cobra.Command, _ []string) error {
	projectDir := getProjectDir()

	// ── Profile ──────────────────────────────────────────────────────────────
	ctx := grimctx.Detect(projectDir)
	profileLabel := ctx.Profile
	if profileLabel == "" {
		profileLabel = colorize(ansiYellow, "(none)")
	}

	// ── Skills count ─────────────────────────────────────────────────────────
	regs := skills.AllSkillsPackages()
	totalSkills := 0
	for _, reg := range regs {
		if all, err := skills.WalkSkills(reg.Root); err == nil {
			totalSkills += len(all)
		}
	}
	skillsLabel := fmt.Sprintf("%d installed", totalSkills)
	if n := len(regs); n > 0 {
		skillsLabel += fmt.Sprintf(" across %d %s", n, plural("package", "packages", n))
	}

	// ── Last check ───────────────────────────────────────────────────────────
	reportPath := resolvedReportPath(projectDir)
	checkLabel := colorize(ansiGray, "(no report — run `grimoire check`)")
	stalenessLabel := colorize(ansiYellow, "unknown")
	cfg, _ := config.Load(projectDir)
	stalenessDays := cfg.StalenessDays
	if stalenessDays == 0 {
		stalenessDays = 7
	}

	if fi, err := os.Stat(reportPath); err == nil {
		age := time.Since(fi.ModTime())
		ageStr := formatAge(age)

		report, loadErr := compliance.Load(reportPath)
		if loadErr == nil {
			pct := fmt.Sprintf("%.1f%%", report.Coverage.OverallPct)
			if report.Threshold.Status == "pass" {
				checkLabel = fmt.Sprintf("%s — %s (%s)", ageStr, colorize(ansiGreen, "PASS"), pct)
			} else {
				checkLabel = fmt.Sprintf("%s — %s (%s)", ageStr, colorize(ansiRed, "FAIL"), pct)
			}
		} else {
			checkLabel = ageStr
		}

		ageDays := age.Hours() / 24
		if ageDays > float64(stalenessDays) {
			stalenessLabel = colorize(ansiRed, fmt.Sprintf("STALE (%.0f days, threshold: %d)", ageDays, stalenessDays))
		} else {
			stalenessLabel = colorize(ansiGreen, fmt.Sprintf("OK (%.0f days, threshold: %d)", ageDays, stalenessDays))
		}
	}

	// ── Report ───────────────────────────────────────────────────────────────
	label := func(key, val string) {
		fmt.Printf("  %-14s %s\n", key+":", val)
	}

	fmt.Printf("\n%s\n\n", tui.StyleBold.Render("grimoire status"))

	dir := projectDir
	if abs, err := filepath.Abs(projectDir); err == nil {
		dir = abs
	}
	label("Project", dir)
	label("Profile", profileLabel)
	label("Skills", skillsLabel)
	label("Last check", checkLabel)
	label("Staleness", stalenessLabel)

	// Show counts from last report diagnostics.
	if report, err := compliance.Load(reportPath); err == nil {
		errors := filterBySeverity(report.Diagnostics, 1)
		warnings := filterBySeverity(report.Diagnostics, 2)
		if len(errors)+len(warnings) > 0 {
			parts := []string{}
			if len(errors) > 0 {
				parts = append(parts, colorize(ansiRed, fmt.Sprintf("%d error", len(errors)))+pluralSuffix(len(errors)))
			}
			if len(warnings) > 0 {
				parts = append(parts, colorize(ansiYellow, fmt.Sprintf("%d warning", len(warnings)))+pluralSuffix(len(warnings)))
			}
			label("Findings", strings.Join(parts, ", "))
		} else {
			label("Findings", colorize(ansiGreen, "none"))
		}
	}

	fmt.Println()
	return nil
}

func formatAge(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

func plural(singular, pluralForm string, n int) string {
	if n == 1 {
		return singular
	}
	return pluralForm
}

func pluralSuffix(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
