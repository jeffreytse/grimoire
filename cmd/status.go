package cmd

import (
	"encoding/json"
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

var flagStatusJSON bool

func init() {
	statusCmd.Flags().BoolVar(&flagStatusJSON, "json", false, "output as JSON")
}

// statusReport is the project health that `grimoire status` shows, and what
// `--json` prints.
type statusReport struct {
	Project       string     `json:"project"`
	Profile       string     `json:"profile"`
	Skills        int        `json:"skills"`
	Packages      int        `json:"packages"`
	LastCheck     *time.Time `json:"last_check,omitempty"`
	Result        string     `json:"result,omitempty"` // "pass" | "fail"
	CoveragePct   *float64   `json:"coverage_pct,omitempty"`
	Stale         *bool      `json:"stale,omitempty"`
	StalenessDays int        `json:"staleness_days"`
	Errors        *int       `json:"errors,omitempty"`
	Warnings      *int       `json:"warnings,omitempty"`
}

func collectStatus(projectDir string) statusReport {
	st := statusReport{Project: projectDir, Profile: grimctx.Detect(projectDir).Profile}
	if abs, err := filepath.Abs(projectDir); err == nil {
		st.Project = abs
	}
	regs := skills.AllSkillsPackages()
	st.Packages = len(regs)
	for _, reg := range regs {
		if all, err := skills.WalkSkills(reg.Root); err == nil {
			st.Skills += len(all)
		}
	}
	cfg, _ := config.Load(projectDir)
	st.StalenessDays = cfg.StalenessDays
	if st.StalenessDays == 0 {
		st.StalenessDays = 7
	}
	reportPath := resolvedReportPath(projectDir)
	fi, err := os.Stat(reportPath)
	if err != nil {
		return st
	}
	mod := fi.ModTime()
	stale := time.Since(mod).Hours()/24 > float64(st.StalenessDays)
	st.LastCheck, st.Stale = &mod, &stale
	if report, loadErr := compliance.Load(reportPath); loadErr == nil {
		pct := report.Coverage.OverallPct
		st.CoveragePct = &pct
		st.Result = "fail"
		if report.Threshold.Status == "pass" {
			st.Result = "pass"
		}
		errs := len(filterBySeverity(report.Diagnostics, 1))
		warns := len(filterBySeverity(report.Diagnostics, 2))
		st.Errors, st.Warnings = &errs, &warns
	}
	return st
}

func runStatus(_ *cobra.Command, _ []string) error {
	st := collectStatus(getProjectDir())
	if flagStatusJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(st)
	}

	profileLabel := st.Profile
	if profileLabel == "" {
		profileLabel = colorize(ansiYellow, "(none)")
	}
	skillsLabel := fmt.Sprintf("%d installed", st.Skills)
	if st.Packages > 0 {
		skillsLabel += fmt.Sprintf(" across %d %s", st.Packages, plural("package", "packages", st.Packages))
	}
	checkLabel := colorize(ansiGray, "(no report — run `grimoire check`)")
	stalenessLabel := colorize(ansiYellow, "unknown")
	if st.LastCheck != nil {
		age := time.Since(*st.LastCheck)
		checkLabel = formatAge(age)
		if st.CoveragePct != nil {
			pct := fmt.Sprintf("%.1f%%", *st.CoveragePct)
			if st.Result == "pass" {
				checkLabel = fmt.Sprintf("%s — %s (%s)", checkLabel, colorize(ansiGreen, "PASS"), pct)
			} else {
				checkLabel = fmt.Sprintf("%s — %s (%s)", checkLabel, colorize(ansiRed, "FAIL"), pct)
			}
		}
		ageDays := age.Hours() / 24
		if *st.Stale {
			stalenessLabel = colorize(ansiRed, fmt.Sprintf("STALE (%.0f days, threshold: %d)", ageDays, st.StalenessDays))
		} else {
			stalenessLabel = colorize(ansiGreen, fmt.Sprintf("OK (%.0f days, threshold: %d)", ageDays, st.StalenessDays))
		}
	}

	label := func(key, val string) {
		fmt.Printf("  %-14s %s\n", key+":", val)
	}
	fmt.Printf("\n%s\n\n", tui.StyleBold.Render("grimoire status"))
	label("Project", st.Project)
	label("Profile", profileLabel)
	label("Skills", skillsLabel)
	label("Last check", checkLabel)
	label("Staleness", stalenessLabel)

	// Show counts from last report diagnostics.
	if st.Errors != nil {
		if *st.Errors+*st.Warnings > 0 {
			parts := []string{}
			if *st.Errors > 0 {
				parts = append(parts, colorize(ansiRed, fmt.Sprintf("%d error", *st.Errors))+pluralSuffix(*st.Errors))
			}
			if *st.Warnings > 0 {
				parts = append(parts, colorize(ansiYellow, fmt.Sprintf("%d warning", *st.Warnings))+pluralSuffix(*st.Warnings))
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
