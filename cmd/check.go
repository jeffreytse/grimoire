package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/compliance"
	"github.com/jeffreytse/grimoire/internal/settings"
)

const (
	ansiRed    = "\033[31m"
	ansiYellow = "\033[33m"
	ansiGreen  = "\033[32m"
	ansiGray   = "\033[90m"
	ansiReset  = "\033[0m"
)

var (
	flagReport       string
	flagJSON         bool
	flagNoColor      bool
	flagFailOnError  bool
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check compliance against the latest report",
	RunE:  runCheck,
}

func init() {
	checkCmd.Flags().StringVar(&flagReport, "report", "", "path to compliance report (default: .grimoire/reports/compliance-latest.json)")
	checkCmd.Flags().BoolVar(&flagJSON, "json", false, "output raw JSON")
	checkCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "disable ANSI color")
	checkCmd.Flags().BoolVar(&flagFailOnError, "fail-on-error", false, "exit 1 if any error-severity diagnostic exists (no settings.toml required)")
}

func runCheck(cmd *cobra.Command, args []string) error {
	report, err := compliance.Load(flagReport)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2) //nolint:revive // intentional: compliance failure must return exit code 2
	}

	if flagJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(report)
	}

	printSummary(report)

	// Settings-driven thresholds override what the AI wrote in the report,
	// making CI gating deterministic regardless of LLM output variation.
	resolved, _ := settings.Load(".")
	section := resolved.ResolveSection(report.Scope)

	failed := false

	if section.ComplianceThreshold > 0 {
		if report.Coverage.OverallPct < section.ComplianceThreshold {
			fmt.Fprintf(os.Stderr, "\n  threshold (settings): %.0f%% required, got %.1f%%\n",
				section.ComplianceThreshold, report.Coverage.OverallPct)
			failed = true
		}
	} else if report.Threshold.Status != "pass" {
		failed = true
	}

	errorCount := len(filterBySeverity(report.Diagnostics, 1))

	if section.ComplianceThresholdError >= 0 && errorCount > section.ComplianceThresholdError {
		fmt.Fprintf(os.Stderr, "\n  error limit (settings): %d allowed, found %d\n",
			section.ComplianceThresholdError, errorCount)
		failed = true
	}

	if flagFailOnError && errorCount > 0 {
		fmt.Fprintf(os.Stderr, "\n  --fail-on-error: %d error-severity diagnostic(s) found\n", errorCount)
		failed = true
	}

	if failed {
		os.Exit(1) //nolint:revive // intentional: compliance failure must return non-zero exit code
	}
	return nil
}

func printSummary(r *compliance.Report) {
	pass := r.Threshold.Status == "pass"

	statusLabel := colorize(ansiGreen, "PASS")
	if !pass {
		statusLabel = colorize(ansiRed, "FAIL")
	}

	fmt.Printf("Compliance: %.1f%% — %s\n", r.Coverage.OverallPct, statusLabel)
	fmt.Printf("  Practices: %d total, %d passing, %d partial, %d failing\n",
		r.Coverage.Practices.Total,
		r.Coverage.Practices.Passing,
		r.Coverage.Practices.Partial,
		r.Coverage.Practices.Failing,
	)

	errors := filterBySeverity(r.Diagnostics, 1)
	warnings := filterBySeverity(r.Diagnostics, 2)

	if len(errors) > 0 {
		fmt.Println()
		for i := range errors {
			loc := formatLoc(&errors[i])
			fmt.Printf("  %s %s%s\n", colorize(ansiRed, "✗"), errors[i].Message, loc)
		}
	}

	if len(warnings) > 0 {
		fmt.Println()
		for i := range warnings {
			loc := formatLoc(&warnings[i])
			fmt.Printf("  %s %s%s\n", colorize(ansiYellow, "⚠"), warnings[i].Message, loc)
		}
	}

	if pass && len(errors) == 0 && len(warnings) == 0 {
		fmt.Printf("  %s All criteria pass.\n", colorize(ansiGreen, "✓"))
	}

	if len(r.Coverage.Details) > 0 {
		fmt.Println()
		fmt.Println("  Practices:")
		for _, d := range r.Coverage.Details {
			bar := colorize(ansiGreen, "✓")
			if d.Failing > 0 {
				bar = colorize(ansiRed, "✗")
			} else if d.Partial > 0 {
				bar = colorize(ansiYellow, "~")
			}
			fmt.Printf("    %s  %-40s %d/%d  %.0f%%\n",
				bar, d.Name, d.Passing, d.Total, d.CoveragePct)
		}
	}

	fmt.Printf("\nThreshold: %.0f%% required, %.1f%% actual — %s\n",
		r.Threshold.Required,
		r.Threshold.Actual,
		statusLabel,
	)
}

func filterBySeverity(diags []compliance.Diagnostic, severity int) []compliance.Diagnostic {
	var out []compliance.Diagnostic
	for i := range diags {
		if diags[i].Severity == severity {
			out = append(out, diags[i])
		}
	}
	return out
}

func formatLoc(d *compliance.Diagnostic) string {
	if d.URI == "" {
		return ""
	}
	line := d.Range.Start.Line + 1
	return colorize(ansiGray, fmt.Sprintf(" (%s:%d)", d.URI, line))
}

func colorize(code, s string) string {
	if flagNoColor || os.Getenv("NO_COLOR") != "" {
		return s
	}
	return code + s + ansiReset
}
