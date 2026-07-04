package cmd

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/compliance"
	"github.com/jeffreytse/grimoire/internal/config"
	grimctx "github.com/jeffreytse/grimoire/internal/context"
	"github.com/jeffreytse/grimoire/internal/profiles"
	"github.com/jeffreytse/grimoire/internal/rules"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

const (
	ansiRed    = "\033[31m"
	ansiYellow = "\033[33m"
	ansiGreen  = "\033[32m"
	ansiGray   = "\033[90m"
	ansiReset  = "\033[0m"
)

var (
	flagReport         string
	flagFromReport     bool
	flagJSON           bool
	flagNoColor        bool
	flagFailOnError    bool
	flagVia            string
	flagCI             bool
	flagJUnit          string
	flagScope          string   // "auto" | "full" | "changed"
	flagCheckBatchSize int      // max KB of file content per AI call
	flagExclude        []string // glob patterns — files to exclude from analysis
	flagPreferAPI      bool     // skip local CLIs, use API directly (temperature=0)
	flagLive           bool     // start live HTTP server with SSE auto-refresh
	flagPort           int      // port for --live server
)

// skillRubric holds a skill name, its SKILL.md body, and explicit criteria for prompting.
type skillRubric struct {
	Name     string
	Body     string
	Criteria []string // from SKILL.md frontmatter criteria: — empty when not defined
}

var checkCmd = &cobra.Command{
	Use:   "check [files/dirs/globs...]",
	Short: "Run independent AI compliance check",
	Long: `Run an independent AI compliance check against the active profile's skills.

Uses an installed local AI agent (claude, gemini, codex, copilot, etc.) when
available, otherwise falls back to a configured or auto-detected API provider.

To read a pre-generated compliance report instead, use --from-report.`,
	RunE: runCheck,
}

func init() {
	checkCmd.Flags().BoolVar(&flagFromReport, "from-report", false, "read pre-generated compliance report instead of running AI check")
	checkCmd.Flags().StringVar(&flagReport, "report", "", "path to compliance report — implies --from-report (default: .grimoire/reports/compliance-latest.json)")
	checkCmd.Flags().BoolVar(&flagJSON, "json", false, "output raw JSON (report mode only)")
	checkCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "disable ANSI color")
	checkCmd.Flags().BoolVar(&flagFailOnError, "fail-on-error", false, "exit 1 if any error-severity diagnostic exists (report mode only)")
	checkCmd.Flags().StringVar(&flagVia, "via", "", "force a specific local AI agent (claude, gemini, codex, copilot, opencode, openclaw)")
	checkCmd.Flags().BoolVar(&flagCI, "ci", false, "CI output mode: GitHub Actions annotations + structured exit codes (auto-enabled when GITHUB_ACTIONS=true)")
	checkCmd.Flags().StringVar(&flagJUnit, "junit", "", "write JUnit XML report to path (e.g. .grimoire/reports/junit.xml)")
	checkCmd.Flags().StringVar(&flagScope, "scope", "auto",
		`file scope: "auto" (default: checksum-based incremental), "full" (all files), "changed" (git diff HEAD)`)
	checkCmd.Flags().IntVar(&flagCheckBatchSize, "batch-size", 256, "max file content per AI call in KB (default 256)")
	checkCmd.Flags().StringArrayVar(&flagExclude, "exclude", nil, `glob patterns for files to exclude from analysis, e.g. --exclude 'vendor/**' --exclude '**/*.pb.go'`)
	checkCmd.Flags().BoolVar(&flagPreferAPI, "prefer-api", false, "skip local AI CLI and route through API directly (sets temperature=0 for deterministic results)")
	checkCmd.Flags().BoolVar(&flagLive, "live", false, "start live server — serve HTML report and auto-refresh browser on file change")
	checkCmd.Flags().IntVar(&flagPort, "port", 7890, "port for --live server")
	checkCmd.Flags().BoolVar(&flagNoGitignore, "no-gitignore", false, "disable .gitignore-based file filtering")
}

func runCheck(cmd *cobra.Command, args []string) error {
	if flagFromReport || cmd.Flags().Changed("report") {
		return runReportCheck(getProjectDir())
	}
	projectDir := getProjectDir()
	initGitignoreMatcher(projectDir)
	initExcludePatterns(projectDir, flagExclude)
	if flagLive {
		return runLiveCheck(projectDir)
	}
	if len(args) > 0 {
		files, err := expandCheckArgs(projectDir, args, flagExclude)
		if err != nil {
			return err
		}
		_, err = runIndependentCheck(context.Background(), projectDir, files...)
		return err
	}
	_, err := runIndependentCheck(context.Background(), projectDir)
	return err
}

// expandCheckArgs resolves each arg (file, directory, or doublestar glob) to a
// deduplicated, relativized list of checkable project files, respecting --exclude.
func expandCheckArgs(projectDir string, args, excludePatterns []string) ([]string, error) {
	seen := make(map[string]bool)
	var files []string
	for _, arg := range args {
		hasGlob := strings.ContainsAny(arg, "*?[")
		if hasGlob {
			matches, err := doublestar.Glob(os.DirFS(projectDir), arg)
			if err != nil {
				return nil, fmt.Errorf("invalid glob %q: %w", arg, err)
			}
			if len(matches) == 0 {
				fmt.Fprintf(os.Stderr, "  warn: no files matched %q\n", arg)
			}
			for _, m := range matches {
				abs := filepath.Join(projectDir, m)
				if fi, _ := os.Stat(abs); fi != nil && fi.IsDir() {
					walkFilesInto(projectDir, abs, excludePatterns, seen, &files)
				} else if !seen[m] && !isSkippableFile(m) && !matchesExcludePatterns(m, excludePatterns) {
					seen[m] = true
					files = append(files, m)
				}
			}
		} else {
			abs := arg
			if !filepath.IsAbs(arg) {
				abs = filepath.Join(projectDir, arg)
			}
			info, err := os.Stat(abs)
			if err != nil {
				return nil, fmt.Errorf("no such file or directory: %s", arg)
			}
			if info.IsDir() {
				walkFilesInto(projectDir, abs, excludePatterns, seen, &files)
			} else {
				rel, _ := filepath.Rel(projectDir, abs)
				if !seen[rel] && !matchesExcludePatterns(rel, excludePatterns) {
					seen[rel] = true
					files = append(files, rel)
				}
			}
		}
	}
	return files, nil
}

// walkFilesInto enumerates checkable files under dir into out, skipping dot dirs,
// vendor, node_modules, skippable files, and any exclude patterns.
func walkFilesInto(projectDir, dir string, excludePatterns []string, seen map[string]bool, out *[]string) {
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		if d.IsDir() {
			if strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			rel, _ := filepath.Rel(projectDir, path)
			if shouldSkip(filepath.ToSlash(rel), true) {
				return filepath.SkipDir
			}
			return nil
		}
		rel, _ := filepath.Rel(projectDir, path)
		if !seen[rel] && !isSkippableFile(rel) && !matchesExcludePatterns(rel, excludePatterns) {
			seen[rel] = true
			*out = append(*out, rel)
		}
		return nil
	})
}

// runReportCheck reads a pre-generated compliance report and enforces thresholds.
func runReportCheck(projectDir string) error {
	reportPath := flagReport
	if reportPath == "" {
		reportPath = resolvedReportPath(projectDir)
	}
	report, err := compliance.Load(reportPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2) //nolint:revive // intentional: compliance failure must return exit code 2
	}

	eng := &rules.Engine{
		SkillsPackages: skills.AllSkillsPackages(),
		ProjectDir:     projectDir,
	}
	if found := eng.Run(); len(found) > 0 {
		report.Diagnostics = append(found, report.Diagnostics...)
	}

	if flagJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(report)
	}

	ciMode := flagCI || os.Getenv("GITHUB_ACTIONS") == "true"

	printSummary(report, "full", nil)

	if htmlPath, htmlErr := compliance.WriteHTMLReport(reportPath, cliVersion, projectDir); htmlErr == nil {
		fmt.Printf("  HTML:   %s\n", htmlPath)
	}

	if ciMode {
		emitGHAAnnotations(report)
	}

	if flagJUnit != "" {
		if err := writeJUnitXML(report, flagJUnit); err != nil {
			fmt.Fprintf(os.Stderr, "  warn: writing JUnit XML: %v\n", err)
		} else {
			fmt.Printf("  JUnit XML written to %s\n", flagJUnit)
		}
	}

	// Settings-driven thresholds override what the AI wrote in the report,
	// making CI gating deterministic regardless of LLM output variation.
	resolved, _ := config.Load(projectDir)
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

func printSummary(r *compliance.Report, mode string, filesToCheck []string) {
	pass := r.Threshold.Status == "pass"

	statusLabel := colorize(ansiGreen, "PASS")
	if !pass {
		statusLabel = colorize(ansiRed, "FAIL")
	}

	if mode == "incremental" && len(filesToCheck) > 0 {
		checked := make(map[string]bool, len(filesToCheck))
		for _, f := range filesToCheck {
			checked[f] = true
		}
		var cur []compliance.Diagnostic
		for i := range r.Diagnostics {
			d := r.Diagnostics[i]
			if checked[d.URI] {
				cur = append(cur, d)
			}
		}
		errs := filterBySeverity(cur, 1)
		warns := filterBySeverity(cur, 2)

		fmt.Printf("Compliance: %.1f%% — %s\n", r.Coverage.OverallPct, statusLabel)
		if len(errs) > 0 {
			fmt.Println()
			for i := range errs {
				fmt.Printf("  %s %s%s\n", colorize(ansiRed, "✗"), errs[i].Message, formatLoc(&errs[i]))
			}
		}
		if len(warns) > 0 {
			fmt.Println()
			for i := range warns {
				fmt.Printf("  %s %s%s\n", colorize(ansiYellow, "⚠"), warns[i].Message, formatLoc(&warns[i]))
			}
		}
		if len(errs) == 0 && len(warns) == 0 {
			fmt.Printf("  %s All criteria pass.\n", colorize(ansiGreen, "✓"))
		}
		fmt.Printf("\nThreshold: %.0f%% required, %.1f%% actual — %s\n",
			r.Threshold.Required, r.Threshold.Actual, statusLabel)
		return
	}

	fmt.Printf("Compliance: %.1f%% — %s\n", r.Coverage.OverallPct, statusLabel)
	fmt.Printf("  Practices: %d total, %d passing, %d partial, %d failing\n",
		r.Coverage.Practices.Total,
		r.Coverage.Practices.Passing,
		r.Coverage.Practices.Partial,
		r.Coverage.Practices.Failing,
	)

	errs := filterBySeverity(r.Diagnostics, 1)
	warnings := filterBySeverity(r.Diagnostics, 2)

	if len(errs) > 0 {
		fmt.Println()
		for i := range errs {
			loc := formatLoc(&errs[i])
			fmt.Printf("  %s %s%s\n", colorize(ansiRed, "✗"), errs[i].Message, loc)
		}
	}

	if len(warnings) > 0 {
		fmt.Println()
		for i := range warnings {
			loc := formatLoc(&warnings[i])
			fmt.Printf("  %s %s%s\n", colorize(ansiYellow, "⚠"), warnings[i].Message, loc)
		}
	}

	if pass && len(errs) == 0 && len(warnings) == 0 {
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

func computeSummary(r *compliance.Report) compliance.ReportSummary {
	s := compliance.ReportSummary{Pass: r.Coverage.Practices.Passing}
	for i := range r.Diagnostics {
		switch r.Diagnostics[i].Severity {
		case 1:
			s.Errors++
		case 2:
			s.Warnings++
		case 3:
			s.Info++
		}
		if r.Diagnostics[i].Status == "suppressed" {
			s.Suppressed++
		}
	}
	return s
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

func resolvedReportPath(cwd string) string {
	r, _ := config.Load(cwd)
	if r.ReportPath != "" {
		if filepath.IsAbs(r.ReportPath) {
			return r.ReportPath
		}
		return filepath.Join(cwd, r.ReportPath)
	}
	return filepath.Join(cwd, compliance.DefaultReportPath)
}

// ── CI output ────────────────────────────────────────────────────────────────

// emitGHAAnnotations writes GitHub Actions workflow commands for each diagnostic.
// Errors emit ::error, warnings emit ::warning.
// Format: ::error file={uri},line={line},title={title}::{message}.
func emitGHAAnnotations(r *compliance.Report) {
	for i := range r.Diagnostics {
		d := &r.Diagnostics[i]
		line := d.Range.Start.Line + 1
		title := d.Code
		if title == "" {
			title = d.Practice
		}
		props := fmt.Sprintf("file=%s,line=%d,title=%s", d.URI, line, title)
		switch d.Severity {
		case 1:
			fmt.Printf("::error %s::%s\n", props, d.Message)
		case 2:
			fmt.Printf("::warning %s::%s\n", props, d.Message)
		default:
			fmt.Printf("::notice %s::%s\n", props, d.Message)
		}
	}
	// Emit overall result as a GHA step summary notice.
	status := "PASS"
	if r.Threshold.Status != "pass" {
		status = "FAIL"
	}
	fmt.Printf("::notice title=Grimoire Compliance::%.1f%% — %s (%d practices)\n",
		r.Coverage.OverallPct, status, r.Coverage.Practices.Total)
}

// writeJUnitXML writes a JUnit XML report to path. Each practice maps to a testcase;
// failing/partial practices are reported as failures.
func writeJUnitXML(r *compliance.Report, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	failures := 0
	for i := range r.Diagnostics {
		if r.Diagnostics[i].Severity == 1 {
			failures++
		}
	}

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	fmt.Fprintf(&sb, `<testsuites name="grimoire-compliance" tests="%d" failures="%d">`,
		r.Coverage.Practices.Total, r.Coverage.Practices.Failing)
	sb.WriteString("\n")
	fmt.Fprintf(&sb, `  <testsuite name="%s" tests="%d" failures="%d" timestamp="%s">`, //nolint:gocritic
		xmlEscape(r.Scope), r.Coverage.Practices.Total, r.Coverage.Practices.Failing, r.Timestamp)
	sb.WriteString("\n")

	// One testcase per practice detail.
	for _, d := range r.Coverage.Details {
		fmt.Fprintf(&sb, `    <testcase name="%s" classname="grimoire">`, xmlEscape(d.Name)) //nolint:gocritic
		if d.Failing > 0 {
			sb.WriteString("\n")
			fmt.Fprintf(&sb, `      <failure message="%.0f%% coverage (%d/%d criteria passing)"/>`,
				d.CoveragePct, d.Passing, d.Total)
			sb.WriteString("\n    </testcase>\n")
		} else {
			sb.WriteString("</testcase>\n")
		}
	}

	// If no practice details, emit one testcase for overall result.
	if len(r.Coverage.Details) == 0 {
		sb.WriteString(`    <testcase name="overall" classname="grimoire">`)
		if r.Threshold.Status != "pass" {
			sb.WriteString("\n")
			fmt.Fprintf(&sb, `      <failure message="%.1f%% overall (threshold %.0f%%)"/>`,
				r.Coverage.OverallPct, r.Threshold.Required)
			sb.WriteString("\n    </testcase>\n")
		} else {
			sb.WriteString("</testcase>\n")
		}
	}

	sb.WriteString("  </testsuite>\n</testsuites>\n")
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}

// ── Provider resolution ──────────────────────────────────────────────────────

// builtinProvider holds default settings for a named LLM provider.
type builtinProvider struct {
	baseURL   string
	apiKeyEnv string
	model     string
	format    string // "anthropic" or "openai"
}

var builtinProviders = map[string]builtinProvider{
	"anthropic":  {"https://api.anthropic.com", "ANTHROPIC_API_KEY", "claude-haiku-4-5-20251001", "anthropic"},
	"openai":     {"https://api.openai.com/v1", "OPENAI_API_KEY", "gpt-4o-mini", "openai"},
	"openrouter": {"https://openrouter.ai/api/v1", "OPENROUTER_API_KEY", "anthropic/claude-3.5-haiku", "openai"},
	"grok":       {"https://api.x.ai/v1", "XAI_API_KEY", "grok-3-mini", "openai"},
	"ollama":     {"http://localhost:11434/v1", "", "llama3.2", "openai"},
	"groq":       {"https://api.groq.com/openai/v1", "GROQ_API_KEY", "llama-3.1-8b-instant", "openai"},
}

// autoDetectOrder is the provider priority when no provider is configured.
var autoDetectOrder = []string{"anthropic", "openai", "openrouter", "grok", "groq", "ollama"}

// resolvedProvider is a fully-hydrated provider ready to call.
type resolvedProvider struct {
	BaseURL   string
	APIKey    string
	Model     string
	Format    string
	MaxTokens int // resolved default: 8192
}

// resolveProvider merges user config with built-in defaults.
func resolveProvider(cfg *config.LLMProviderConfig) (resolvedProvider, bool) {
	builtin, known := builtinProviders[cfg.Name]

	baseURL := cfg.BaseURL
	if baseURL == "" && known {
		baseURL = builtin.baseURL
	}
	keyEnv := cfg.APIKeyEnv
	if keyEnv == "" && known {
		keyEnv = builtin.apiKeyEnv
	}
	model := cfg.Model
	if model == "" && known {
		model = builtin.model
	}
	format := cfg.Format
	if format == "" && known {
		format = builtin.format
	}

	if baseURL == "" {
		return resolvedProvider{}, false
	}
	apiKey := ""
	if keyEnv != "" {
		apiKey = os.Getenv(keyEnv)
		if apiKey == "" && cfg.Name != "ollama" {
			return resolvedProvider{}, false // key required but not set
		}
	}
	maxTokens := cfg.MaxTokens
	if maxTokens == 0 {
		maxTokens = 8192
	}
	return resolvedProvider{baseURL, apiKey, model, format, maxTokens}, true
}

// ── Executor resolution ──────────────────────────────────────────────────────

type executorKind int

const (
	execLocalCLI executorKind = iota
	execAPI
	execPrint
)

type executorSpec struct {
	Kind     executorKind
	Agent    string           // for execLocalCLI
	Provider resolvedProvider // for execAPI
}

// resolveExecutor picks the best executor given the --via flag and project config.
func resolveExecutor(projectDir, via string) executorSpec {
	// 1. --via flag: force a specific local CLI.
	if via != "" {
		return executorSpec{Kind: execLocalCLI, Agent: via}
	}

	cfg, _ := config.Load(projectDir)

	// 2. Local CLIs: intersect configured order with detected agents.
	// Skipped when --prefer-api is set — routes directly to API for temperature=0.
	if !flagPreferAPI {
		order := cfg.Core.CheckAgents
		if len(order) == 0 {
			order = agent.CheckAgents
		}
		detected := make(map[string]struct{})
		for _, ag := range agent.DetectedCheckAgents() {
			detected[ag] = struct{}{}
		}
		for _, ag := range order {
			if _, ok := detected[ag]; ok {
				return executorSpec{Kind: execLocalCLI, Agent: ag}
			}
		}
	}

	// 3. Configured API provider.
	if cfg.Core.CheckProvider.Name != "" {
		if p, ok := resolveProvider(&cfg.Core.CheckProvider); ok {
			return executorSpec{Kind: execAPI, Provider: p}
		}
	}

	// 4. Auto-detect known API keys.
	for _, name := range autoDetectOrder {
		builtin := builtinProviders[name]
		if builtin.apiKeyEnv == "" {
			continue // skip keyless providers (ollama) in auto-detect
		}
		if key := os.Getenv(builtin.apiKeyEnv); key != "" {
			p := resolvedProvider{builtin.baseURL, key, builtin.model, builtin.format, 8192}
			return executorSpec{Kind: execAPI, Provider: p}
		}
	}

	return executorSpec{Kind: execPrint}
}

// ── runIndependentCheck ──────────────────────────────────────────────────────

// runIndependentCheck implements `grimoire check`. Accepts optional override file list
// for watch-triggered incremental checks.
func runIndependentCheck(goCtx context.Context, projectDir string, overrideFiles ...string) (bool, error) {
	ctx := grimctx.Detect(projectDir)
	if ctx.Profile == "" {
		return false, fmt.Errorf("no active profile detected — set [standards] profiles in grimoire.toml")
	}

	opts := resolveOpts(projectDir)
	p, err := profiles.ResolveWithOptions(ctx.Profile, projectDir, opts)
	if err != nil {
		return false, fmt.Errorf("resolving profile %q: %w", ctx.Profile, err)
	}

	regs := skills.AllSkillsPackages()
	allSkillMap := make(map[string]skills.Skill)
	for _, reg := range regs {
		all, walkErr := skills.WalkSkills(reg.Root)
		if walkErr != nil {
			continue
		}
		for i := range all {
			sk := all[i]
			if _, exists := allSkillMap[sk.Name]; !exists {
				allSkillMap[sk.Name] = sk
			}
		}
	}

	var rubrics []skillRubric
	for _, ref := range p.Skills {
		sk, ok := allSkillMap[ref.Name]
		if !ok {
			continue
		}
		if strings.TrimSpace(sk.Body) == "" {
			continue
		}
		rubrics = append(rubrics, skillRubric{sk.Name, sk.Body, sk.Criteria})
	}

	if len(rubrics) == 0 {
		fmt.Printf("  %s  No skill content found for active profile %q.\n", colorize(ansiYellow, "⚠"), ctx.Profile)
		return false, nil
	}

	cfg, _ := config.Load(projectDir)
	section := cfg.ResolveSection(ctx.Profile)
	threshold := section.ComplianceThreshold
	if threshold == 0 {
		threshold = 80
	}

	rubricHash := computeRubricHash(rubrics)

	// ── File collection ───────────────────────────────────────────────────────
	if !flagLive {
		fmt.Println("  Collecting files…")
	}

	var filesToCheck []string
	var fileInfos map[string]os.FileInfo
	var effectiveMode string
	var cache *compliance.CheckCache

	bustRubricCache := func(c *compliance.CheckCache) {
		if c.RubricHash != "" && c.RubricHash != rubricHash {
			if !flagLive {
				fmt.Printf("  Skills updated — running full scan\n")
			}
			c.Files = make(map[string]compliance.FileCacheEntry)
			c.PracticeTotals = nil
			c.PracticeCriteria = nil
		}
		c.RubricHash = rubricHash
	}

	if len(overrideFiles) > 0 {
		cache, _ = compliance.LoadCache(projectDir)
		bustRubricCache(cache)
		filesToCheck, effectiveMode = resolveWatchFiles(overrideFiles, projectDir, cache)
	} else {
		switch flagScope {
		case "full":
			filesToCheck, fileInfos = allProjectFiles(projectDir)
			effectiveMode = "full"
		case "changed":
			filesToCheck = gitChangedFiles(projectDir)
			effectiveMode = "incremental"
		default: // "auto"
			cache, _ = compliance.LoadCache(projectDir)
			bustRubricCache(cache)
			filesToCheck, fileInfos, effectiveMode = autoScopeFiles(projectDir, cache)
		}
	}

	if effectiveMode == "skip" {
		_ = compliance.TouchTimestamp(projectDir)
		return false, runReportCheck(projectDir)
	}
	if len(filesToCheck) == 0 {
		if effectiveMode == "full" {
			fmt.Printf("  %s  No project files found.\n", colorize(ansiYellow, "⚠"))
		}
		// Regenerate HTML from existing report so layout/theme always stays current.
		reportPath := resolvedReportPath(projectDir)
		if htmlPath, htmlErr := compliance.WriteHTMLReport(reportPath, cliVersion, projectDir); htmlErr == nil {
			if !flagLive {
				fmt.Printf("  HTML:   %s\n", htmlPath)
			}
		}
		return false, nil
	}

	ex := resolveExecutor(projectDir, flagVia)

	if !flagLive {
		switch ex.Kind {
		case execLocalCLI:
			fmt.Printf("  Agent:  %s (local CLI — temperature uncontrolled; use --prefer-api for determinism)\n", ex.Agent)
		case execAPI:
			fmt.Printf("  Agent:  %s API (temperature=0)\n", ex.Provider.Model)
		}
	}

	// ── No executor: print rubric + guidance ─────────────────────────────────
	if ex.Kind == execPrint {
		label := "Project files"
		if effectiveMode == "incremental" {
			label = "Changed files"
		}
		var rb strings.Builder
		rb.WriteString(label + ":\n")
		for _, f := range filesToCheck {
			fmt.Fprintf(&rb, "  %s\n", f)
		}
		rb.WriteString("\n")
		for _, r := range rubrics {
			fmt.Fprintf(&rb, "=== Skill: %s ===\n", r.Name)
			if len(r.Criteria) > 0 {
				fmt.Fprintf(&rb, "Evaluate EXACTLY these %d criteria (use these exact names in criteria_matrix):\n", len(r.Criteria))
				for i, c := range r.Criteria {
					fmt.Fprintf(&rb, "  %d. %s\n", i+1, c)
				}
				rb.WriteString("\n")
			}
			rb.WriteString(r.Body)
			rb.WriteString("\n\n")
		}
		fmt.Printf("Independent check rubric (profile: %s, %d skill(s)):\n\n", ctx.Profile, len(rubrics))
		fmt.Println(rb.String())
		detected := agent.DetectedCheckAgents()
		if len(detected) > 0 {
			fmt.Printf("Detected local agents: %s\n", strings.Join(detected, ", "))
			fmt.Printf("Run: grimoire check  (auto-selects first detected agent)\n")
		} else {
			fmt.Printf("No local AI agents detected. Options:\n")
			fmt.Printf("  • Install claude, gemini, or codex CLI\n")
			fmt.Printf("  • Set ANTHROPIC_API_KEY / OPENAI_API_KEY / OPENROUTER_API_KEY\n")
			fmt.Printf("  • Configure [core.check-provider] in grimoire.toml\n")
		}
		return false, nil
	}

	agentLabel := ex.Agent
	if ex.Kind == execAPI {
		agentLabel = ex.Provider.Model
	}
	if !flagLive {
		fmt.Printf("  Running check via %s (profile: %s)…\n", agentLabel, ctx.Profile)
	}

	// Build skills portion shared across all batches.
	var skillsStr strings.Builder
	for _, r := range rubrics {
		fmt.Fprintf(&skillsStr, "=== Skill: %s ===\n", r.Name)
		if len(r.Criteria) > 0 {
			fmt.Fprintf(&skillsStr, "Evaluate EXACTLY these %d criteria (use these exact names in criteria_matrix):\n", len(r.Criteria))
			for i, c := range r.Criteria {
				fmt.Fprintf(&skillsStr, "  %d. %s\n", i+1, c)
			}
			skillsStr.WriteString("\n")
		}
		skillsStr.WriteString(r.Body)
		skillsStr.WriteString("\n\n")
	}

	label := "Project files"
	if effectiveMode == "incremental" {
		label = "Changed files"
	}

	// ── Pre-filter: skip binary/asset files, excluded patterns, and out-of-scope files ──
	extHints := extractExtHints(rubrics)
	var preFiltered []string
	var preSkipped int
	for _, f := range filesToCheck {
		if isKnownBinaryExt(f) {
			preSkipped++
			continue
		}
		if len(globalCheckExclude) > 0 && matchesExcludePatterns(f, globalCheckExclude) {
			preSkipped++
			continue
		}
		if len(extHints) > 0 {
			ext := strings.ToLower(filepath.Ext(f))
			found := false
			for _, h := range extHints {
				if h == ext {
					found = true
					break
				}
			}
			if !found {
				preSkipped++
				continue
			}
		}
		preFiltered = append(preFiltered, f)
	}
	if preSkipped > 0 && !flagLive {
		fmt.Printf("  Skipping %d file(s) — binary, asset, or outside skill scope\n", preSkipped)
	}
	filesToCheck = preFiltered
	if len(filesToCheck) == 0 {
		if !flagLive {
			fmt.Printf("  %s  No matching files to analyze.\n", colorize(ansiYellow, "⚠"))
		}
		reportPath := resolvedReportPath(projectDir)
		if htmlPath, htmlErr := compliance.WriteHTMLReport(reportPath, cliVersion, projectDir); htmlErr == nil {
			if !flagLive {
				fmt.Printf("  HTML:   %s\n", htmlPath)
			}
		}
		return false, nil
	}
	var displayBytes int
	for _, f := range filesToCheck {
		if info, ok := fileInfos[f]; ok {
			displayBytes += int(info.Size())
		} else if info, err := os.Stat(filepath.Join(projectDir, f)); err == nil {
			displayBytes += int(info.Size())
		}
	}
	if !flagLive {
		fmt.Printf("  %s  %d files, %s\n", tui.IconDone, len(filesToCheck), formatBytes(displayBytes))
	}

	// ── Per-batch parallel execution ─────────────────────────────────────────
	type fileResult struct {
		file             string
		diags            []compliance.Diagnostic
		practiceTotals   map[string]int
		practiceCriteria map[string][]compliance.CriterionDetail
		err              error
	}

	parallelism := 5
	if ex.Kind == execAPI {
		parallelism = 15
	}
	start := time.Now()

	fileIndex := make(map[string]int, len(filesToCheck))
	for i, f := range filesToCheck {
		fileIndex[f] = i
	}
	fmt.Println()
	board := tui.NewLiveBoard(filesToCheck)
	board.Start()

	maxBatch := flagCheckBatchSize * 1024
	batches := batchFilesBySize(filesToCheck, fileInfos, projectDir, maxBatch)

	sem := make(chan struct{}, parallelism)
	resultCh := make(chan fileResult, len(filesToCheck))

	var wg sync.WaitGroup
	var printMu sync.Mutex
	var batchErrMu sync.Mutex
	var batchErrors []string
	for _, batch := range batches {
		batch := batch
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			// Read each file individually; skip binary content per file.
			var batchFiles []string
			var combinedContents strings.Builder
			for _, f := range batch {
				data, err := os.ReadFile(filepath.Join(projectDir, f))
				if err != nil {
					continue
				}
				if bytes.Contains(data, []byte{0}) {
					if tui.IsTTY() {
						board.Complete(fileIndex[f], colorize(ansiGray, "-"), "")
					}
					resultCh <- fileResult{file: f}
					continue
				}
				fmt.Fprintf(&combinedContents, "=== %s ===\n%s\n\n", f, string(data))
				batchFiles = append(batchFiles, f)
			}
			if len(batchFiles) == 0 {
				return
			}

			var rb strings.Builder
			rb.WriteString(label + ":\n")
			for _, f := range batchFiles {
				fmt.Fprintf(&rb, "  %s\n", f)
			}
			rb.WriteString("\n")
			rb.WriteString(skillsStr.String())

			result, callErr := callBatch(goCtx, &ex, projectDir, ctx.Profile, rb.String(), combinedContents.String(), threshold)
			if callErr != nil {
				batchErrMu.Lock()
				batchErrors = append(batchErrors, callErr.Error())
				batchErrMu.Unlock()
			}

			// Parse once; distribute diagnostics to per-file results by URI.
			var parsedDiags []compliance.Diagnostic
			var practiceTotals map[string]int
			var practiceCriteria map[string][]compliance.CriterionDetail
			if callErr == nil {
				var rep compliance.Report
				if json.Unmarshal([]byte(extractJSON(result)), &rep) == nil {
					parsedDiags = rep.Diagnostics
					if len(rep.Coverage.Details) > 0 {
						practiceTotals = make(map[string]int, len(rep.Coverage.Details))
						practiceCriteria = make(map[string][]compliance.CriterionDetail, len(rep.Coverage.Details))
						for _, d := range rep.Coverage.Details {
							practiceTotals[d.Name] = d.Total
							if len(d.Criteria) > 0 {
								practiceCriteria[d.Name] = d.Criteria
							}
						}
					}
					// Primary: criteria_matrix (natural LLM output — categorized name lists).
					for practice, entry := range rep.CriteriaMatrix {
						if practiceCriteria == nil {
							practiceCriteria = map[string][]compliance.CriterionDetail{}
						}
						for _, n := range entry.Pass {
							practiceCriteria[practice] = append(practiceCriteria[practice],
								compliance.CriterionDetail{Name: n, Status: "pass"})
						}
						for _, n := range entry.Fail {
							practiceCriteria[practice] = append(practiceCriteria[practice],
								compliance.CriterionDetail{Name: n, Status: "fail"})
						}
					}
					// Fallback: severity-4 pass diagnostics.
					for i := range parsedDiags {
						d := parsedDiags[i]
						if d.Severity == 4 && d.Status == "pass" && d.Practice != "" && d.Criterion != "" {
							if practiceCriteria == nil {
								practiceCriteria = map[string][]compliance.CriterionDetail{}
							}
							practiceCriteria[d.Practice] = append(practiceCriteria[d.Practice],
								compliance.CriterionDetail{Name: d.Criterion, Status: "pass"})
						}
					}
				}
			}

			for _, f := range batchFiles {
				var fileDiags []compliance.Diagnostic
				for i := range parsedDiags {
					d := parsedDiags[i]
					if d.URI == f {
						fileDiags = append(fileDiags, d)
					}
				}
				fr := fileResult{file: f, diags: fileDiags, practiceTotals: practiceTotals, practiceCriteria: practiceCriteria, err: callErr}
				nErr := len(filterBySeverity(fr.diags, 1))
				nWarn := len(filterBySeverity(fr.diags, 2))
				icon := colorize(ansiGreen, "✓")
				detail := ""
				switch {
				case fr.err != nil:
					icon = colorize(ansiRed, "✗")
					detail = "(failed)"
				case nErr > 0 && nWarn > 0:
					icon = colorize(ansiRed, "✗")
					detail = fmt.Sprintf("%d errors, %d warnings", nErr, nWarn)
				case nErr > 0:
					icon = colorize(ansiRed, "✗")
					detail = fmt.Sprintf("%d errors", nErr)
				case nWarn > 0:
					icon = colorize(ansiYellow, "~")
					detail = fmt.Sprintf("%d warnings", nWarn)
				}
				if tui.IsTTY() {
					board.Complete(fileIndex[f], icon, detail)
				} else {
					printMu.Lock()
					if detail != "" {
						fmt.Printf("  %s %s  %s\n", icon, f, detail)
					} else {
						fmt.Printf("  %s %s\n", icon, f)
					}
					printMu.Unlock()
				}
				resultCh <- fr
			}
		}()
	}
	go func() { wg.Wait(); close(resultCh) }()

	var allDiags []compliance.Diagnostic
	var practiceTotals map[string]int
	var practiceCriteria map[string][]compliance.CriterionDetail

	var batchErrCount int
	for r := range resultCh {
		allDiags = append(allDiags, r.diags...)
		if r.err != nil {
			batchErrCount++
		}
		if practiceTotals == nil && r.practiceTotals != nil {
			practiceTotals = r.practiceTotals
		}
		if practiceCriteria == nil && r.practiceCriteria != nil {
			practiceCriteria = r.practiceCriteria
		}
	}
	board.Stop()
	fmt.Println()

	for _, e := range batchErrors {
		fmt.Fprintf(os.Stderr, "  Agent error: %s\n", e)
	}
	if batchErrCount > 0 && practiceTotals == nil && len(allDiags) == 0 {
		return true, fmt.Errorf("all AI batch calls failed (%d file(s)) — see errors above", batchErrCount)
	}

	// Derive totals from criteria_matrix entries when AI omitted coverage.details.
	if practiceTotals == nil && len(practiceCriteria) > 0 {
		practiceTotals = make(map[string]int, len(practiceCriteria))
		for practice, crit := range practiceCriteria {
			practiceTotals[practice] = len(crit)
		}
	}

	// Deduplicate and merge all per-file results.
	allDiags = dedupDiagnostics(allDiags)

	merged := compliance.Report{
		Version:     "1",
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Mode:        effectiveMode,
		Scope:       ".",
		Diagnostics: allDiags,
	}
	if effectiveMode == "incremental" {
		merged.Scope = "changed"
		merged.Git.BaseRef = "HEAD"
	}
	if practiceTotals != nil {
		merged.Coverage = recomputeCoverage(allDiags, practiceTotals, practiceCriteria)
		merged.Threshold = compliance.Threshold{
			Required: threshold,
			Actual:   merged.Coverage.OverallPct,
		}
		if merged.Coverage.OverallPct >= threshold {
			merged.Threshold.Status = "pass"
			merged.Result = "pass"
		} else {
			merged.Threshold.Status = "fail"
			merged.Result = "fail"
		}
	}

	mergedJSON, marshalErr := json.Marshal(merged)
	if marshalErr != nil {
		return true, fmt.Errorf("merging results: %w", marshalErr)
	}
	return true, handleCheckResult(string(mergedJSON), projectDir, effectiveMode, cache, filesToCheck, start)
}

// autoScopeFiles determines which files need checking based on checksums.
// Returns (nil, nil, "skip") when nothing changed.
func autoScopeFiles(dir string, cache *compliance.CheckCache) (files []string, infos map[string]os.FileInfo, mode string) {
	if cache.ConfigHash == "" {
		f, i := allProjectFiles(dir)
		return f, i, "full"
	}
	cfgHash, _ := compliance.ConfigHash(dir)
	if cfgHash != cache.ConfigHash {
		fmt.Printf("  Config changed — running full scan\n")
		cache.Files = make(map[string]compliance.FileCacheEntry)
		f, i := allProjectFiles(dir)
		return f, i, "full"
	}
	all, allInfos := allProjectFiles(dir)

	// Evict cache entries for deleted files.
	current := make(map[string]bool, len(all))
	for _, f := range all {
		current[f] = true
	}
	for f := range cache.Files {
		if !current[f] {
			delete(cache.Files, f)
		}
	}

	// Fast-path: use precomputed FileInfo — no re-stat needed.
	// Files that fail the fast-path are hashed in parallel.
	type hashJob struct{ file string }
	var toHash []hashJob
	var changed []string
	changedInfos := make(map[string]os.FileInfo)
	for _, f := range all {
		info, ok := allInfos[f]
		if !ok {
			continue
		}
		if entry, ok := cache.Files[f]; ok &&
			entry.Mtime == info.ModTime().UnixNano() &&
			entry.Size == info.Size() {
			continue // unchanged — skip SHA256
		}
		toHash = append(toHash, hashJob{f})
	}

	if len(toHash) > 0 {
		type hashResult struct {
			file string
			hash string
		}
		hashSem := make(chan struct{}, 16)
		hashCh := make(chan hashResult, len(toHash))
		var hashWG sync.WaitGroup
		for _, job := range toHash {
			job := job
			hashWG.Add(1)
			go func() {
				defer hashWG.Done()
				hashSem <- struct{}{}
				defer func() { <-hashSem }()
				h, err := compliance.FileHash(filepath.Join(dir, job.file))
				if err == nil {
					hashCh <- hashResult{job.file, h}
				}
			}()
		}
		go func() { hashWG.Wait(); close(hashCh) }()
		for r := range hashCh {
			if entry, ok := cache.Files[r.file]; !ok || entry.Hash != r.hash {
				changed = append(changed, r.file)
				if info, ok := allInfos[r.file]; ok {
					changedInfos[r.file] = info
				}
			}
		}
	}

	if len(changed) == 0 {
		if !flagLive {
			fmt.Printf("  %s  No changes since last check — skipping.\n", colorize(ansiGreen, "✓"))
		}
		return nil, nil, "skip"
	}
	if !flagLive {
		fmt.Printf("  Incremental: %d/%d file(s) changed\n", len(changed), len(all))
		const maxShow = 10
		for i, f := range changed {
			if i >= maxShow {
				fmt.Printf("  ↳ … +%d more\n", len(changed)-maxShow)
				break
			}
			fmt.Printf("  ↳ %s\n", f)
		}
	}
	return changed, changedInfos, "incremental"
}

// resolveWatchFiles determines which watch-triggered files need re-checking.
func resolveWatchFiles(triggered []string, dir string, cache *compliance.CheckCache) (files []string, reason string) {
	cfgHash, _ := compliance.ConfigHash(dir)
	if cfgHash != cache.ConfigHash {
		cache.Files = make(map[string]compliance.FileCacheEntry)
		files, _ := allProjectFiles(dir)
		return files, "full"
	}
	// SHA256 check: skip files whose content hasn't changed since last AI check.
	var toCheck []string
	for _, f := range triggered {
		entry, cached := cache.Files[f]
		if !cached {
			toCheck = append(toCheck, f)
			continue
		}
		h, err := compliance.FileHash(filepath.Join(dir, f))
		if err != nil || h != entry.Hash {
			toCheck = append(toCheck, f)
		}
	}
	if len(toCheck) == 0 {
		return nil, "skip"
	}
	return toCheck, "incremental"
}

// allProjectFiles returns all project files and their FileInfo, using git ls-files with WalkDir fallback.
func allProjectFiles(dir string) (files []string, infos map[string]os.FileInfo) {
	if files, infos := gitAllTrackedFiles(dir); len(files) > 0 {
		return files, infos
	}
	infos = make(map[string]os.FileInfo)
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		if d.IsDir() {
			if strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			rel, _ := filepath.Rel(dir, path)
			if shouldSkip(filepath.ToSlash(rel), true) {
				return filepath.SkipDir
			}
			return nil
		}
		fi, statErr := d.Info()
		if statErr != nil || fi.Size() == 0 {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		if isSkippableFile(rel) {
			return nil
		}
		files = append(files, rel)
		infos[rel] = fi
		return nil
	})
	return files, infos
}

// isGrimoireConfig returns true for user-editable config files.
// Excludes .grimoire/cache.json and .grimoire/reports/ to prevent watch loop.
func isGrimoireConfig(projectDir, absPath string) bool {
	rel, err := filepath.Rel(projectDir, absPath)
	if err != nil {
		return false
	}
	rel = filepath.ToSlash(rel)
	if rel == "grimoire.toml" {
		return true
	}
	if strings.HasPrefix(rel, ".grimoire/") {
		return !strings.HasPrefix(rel, ".grimoire/reports/") &&
			rel != ".grimoire/cache.json"
	}
	return false
}

// recomputeCoverage rebuilds Coverage from merged diagnostics + cached criteria totals.
// Primary source: criteria_matrix fail/pass entries (authoritative).
// Fallback: severity<=2 diagnostics for practices without criteria data.
func recomputeCoverage(diags []compliance.Diagnostic, totals map[string]int, criteria map[string][]compliance.CriterionDetail) compliance.Coverage {
	practiceHasCritData := make(map[string]bool)
	failingByCrit := make(map[string]int)
	passingByCrit := make(map[string]int)
	for practice, crit := range criteria {
		if len(crit) > 0 {
			practiceHasCritData[practice] = true
		}
		for _, c := range crit {
			if c.Status == "fail" {
				failingByCrit[practice]++
			} else {
				passingByCrit[practice]++
			}
		}
	}
	failingByDiag := make(map[string]int)
	for i := range diags {
		d := diags[i]
		if d.Severity <= 2 && d.Practice != "" && !practiceHasCritData[d.Practice] {
			failingByDiag[d.Practice]++
		}
	}
	var details []compliance.PracticeDetail
	totalCriteria, passingCriteria := 0, 0
	pracPassing, pracFailing := 0, 0
	for practice, total := range totals {
		var passing, failing int
		if practiceHasCritData[practice] {
			failing = failingByCrit[practice]
			passing = passingByCrit[practice]
			total = passing + failing
		} else {
			failing = failingByDiag[practice]
			passing = total - failing
			if passing < 0 {
				passing = 0
			}
		}
		pct := 0.0
		if total > 0 {
			pct = float64(passing) / float64(total) * 100
		}
		details = append(details, compliance.PracticeDetail{
			Name:        practice,
			Total:       total,
			Passing:     passing,
			Failing:     failing,
			CoveragePct: pct,
			Criteria:    criteria[practice],
		})
		totalCriteria += total
		passingCriteria += passing
		if failing > 0 {
			pracFailing++
		} else {
			pracPassing++
		}
	}
	overallPct := 0.0
	if totalCriteria > 0 {
		overallPct = float64(passingCriteria) / float64(totalCriteria) * 100
	}
	pracTotal := len(totals)
	pracCovPct := 0.0
	if pracTotal > 0 {
		pracCovPct = float64(pracPassing) / float64(pracTotal) * 100
	}
	return compliance.Coverage{
		OverallPct: overallPct,
		Practices: compliance.PracticeSummary{
			Total:       pracTotal,
			Passing:     pracPassing,
			Failing:     pracFailing,
			CoveragePct: pracCovPct,
		},
		Criteria: compliance.CriteriaSummary{
			Total:       totalCriteria,
			Passing:     passingCriteria,
			Failing:     totalCriteria - passingCriteria,
			CoveragePct: overallPct,
		},
		Details: details,
	}
}

// gitAllTrackedFiles returns all git-tracked non-skippable files and their FileInfo.
func gitAllTrackedFiles(dir string) (files []string, infos map[string]os.FileInfo) {
	cmd := exec.Command("git", "ls-files")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, nil
	}
	infos = make(map[string]os.FileInfo)
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if isSkippableFile(line) {
			continue
		}
		info, statErr := os.Stat(filepath.Join(dir, line))
		if statErr != nil || info.Size() == 0 {
			continue
		}
		files = append(files, line)
		infos[line] = info
	}
	return files, infos
}

// gitChangedFiles returns paths of files modified since the last commit.
func gitChangedFiles(dir string) []string {
	cmd := exec.Command("git", "diff", "--name-only", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		// Try unstaged changes too.
		cmd2 := exec.Command("git", "diff", "--name-only")
		cmd2.Dir = dir
		out, err = cmd2.Output()
		if err != nil {
			return nil
		}
	}
	var files []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}
	return files
}

var extHintRe = regexp.MustCompile(`\*\.(\w+)`)

// isKnownBinaryExt returns true for file extensions that are always binary.
// No I/O — used in the serial pre-filter where reads would be a bottleneck.
func isKnownBinaryExt(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico", ".webp",
		".woff", ".woff2", ".ttf", ".eot", ".otf",
		".mp3", ".mp4", ".wav", ".ogg", ".flac",
		".pdf", ".zip", ".tar", ".gz", ".bz2", ".xz", ".7z",
		".exe", ".dll", ".so", ".dylib", ".bin", ".o", ".a":
		return true
	}
	return false
}

// isAnalyzableFile returns false for binary files and known media/asset types.
// Reads at most 512 bytes to detect binary content.
func isAnalyzableFile(path string) bool {
	if isKnownBinaryExt(path) {
		return false
	}
	f, err := os.Open(path)
	if err != nil {
		return true
	}
	defer func() { _ = f.Close() }()
	var buf [512]byte
	n, _ := f.Read(buf[:])
	return !bytes.Contains(buf[:n], []byte{0})
}

// extractExtHints scans skill bodies for explicit file extension patterns (e.g. *.go).
// Returns nil when no patterns are found — no extension restriction applies.
func extractExtHints(rubrics []skillRubric) []string {
	seen := make(map[string]bool)
	for _, r := range rubrics {
		for _, m := range extHintRe.FindAllStringSubmatch(r.Body, -1) {
			seen["."+strings.ToLower(m[1])] = true
		}
	}
	if len(seen) == 0 {
		return nil
	}
	out := make([]string, 0, len(seen))
	for e := range seen {
		out = append(out, e)
	}
	return out
}

// matchesExcludePatterns reports whether path matches any of the given glob patterns.
// Patterns follow doublestar semantics (** supported). Special cases:
//   - no glob chars: treated as directory prefix (e.g. "vendor" → "vendor/**")
//   - glob with no /: also tried as basename match (e.g. "*.pb.go" matches "internal/foo.pb.go")
func matchesExcludePatterns(path string, patterns []string) bool {
	for _, pat := range patterns {
		hasGlob := strings.ContainsAny(pat, "*?[")
		if !hasGlob {
			// Plain name — match as exact file OR directory prefix.
			plain := strings.TrimSuffix(pat, "/")
			if ok, _ := doublestar.Match(plain+"/**", path); ok {
				return true
			}
			if ok, _ := doublestar.Match(plain, path); ok {
				return true
			}
			continue
		}
		if ok, _ := doublestar.Match(pat, path); ok {
			return true
		}
		// Pattern has glob but no slash — also try as basename match.
		if !strings.Contains(pat, "/") {
			if ok, _ := filepath.Match(pat, filepath.Base(path)); ok {
				return true
			}
		}
	}
	return false
}

// isSkippableFile returns true for files that add byte count but no compliance signal:
// lock files, generated code, minified assets, source maps, build artifacts.
func isSkippableFile(rel string) bool {
	base := filepath.Base(rel)
	switch base {
	case "package-lock.json", "yarn.lock", "pnpm-lock.yaml",
		"Gemfile.lock", "Cargo.lock", "poetry.lock", "composer.lock",
		"go.sum", "flake.lock":
		return true
	}
	ext := strings.ToLower(filepath.Ext(rel))
	switch ext {
	case ".lock", ".snap", ".map":
		return true
	}
	if strings.HasSuffix(rel, ".min.js") || strings.HasSuffix(rel, ".min.css") ||
		strings.HasSuffix(rel, "_generated.go") || strings.HasSuffix(rel, ".pb.go") {
		return true
	}
	return shouldSkip(filepath.ToSlash(rel), false)
}

// batchFilesBySize groups files into batches not exceeding maxBatchBytes each.
// Uses precomputed FileInfo to avoid re-statting; falls back to os.Stat for any missing entries.
// A single file exceeding the limit gets its own batch rather than being dropped.
func batchFilesBySize(files []string, infos map[string]os.FileInfo, dir string, maxBatchBytes int) [][]string {
	var batches [][]string
	var cur []string
	var curSize int
	for _, f := range files {
		var sz int
		if info, ok := infos[f]; ok {
			sz = int(info.Size())
		} else if info, err := os.Stat(filepath.Join(dir, f)); err == nil {
			sz = int(info.Size())
		}
		if sz == 0 {
			continue
		}
		if len(cur) > 0 && curSize+sz > maxBatchBytes {
			batches = append(batches, cur)
			cur, curSize = nil, 0
		}
		cur = append(cur, f)
		curSize += sz
	}
	if len(cur) > 0 {
		batches = append(batches, cur)
	}
	return batches
}

// callBatch invokes the configured AI executor with the pre-built rubric + file contents.
func callBatch(goCtx context.Context, ex *executorSpec, dir, profile, fullRubric, fileContents string, threshold float64) (string, error) {
	prompt := buildLocalPrompt(profile, fullRubric, fileContents, threshold)
	switch ex.Kind {
	case execLocalCLI:
		return callViaLocalAgent(goCtx, ex.Agent, dir, prompt)
	case execAPI:
		const system = "You are grimoire, an independent compliance checker. Output ONLY a single JSON object matching the compliance report schema in the user message. No prose, no markdown, no explanation — JSON only."
		if ex.Provider.Format == "anthropic" {
			return callAnthropicAPI(ex.Provider.APIKey, ex.Provider.Model, ex.Provider.MaxTokens, prompt, "")
		}
		return callOpenAICompatible(ex.Provider, system, prompt)
	default:
		return "", fmt.Errorf("no AI executor configured")
	}
}

// formatBytes formats a byte count as a human-readable string.
func formatBytes(b int) string {
	switch {
	case b < 1024:
		return fmt.Sprintf("%d B", b)
	case b < 1024*1024:
		return fmt.Sprintf("%.0f KB", float64(b)/1024)
	default:
		return fmt.Sprintf("%.1f MB", float64(b)/(1024*1024))
	}
}

// computeRubricHash returns a SHA256 digest of all active skill names, bodies, and criteria.
// Used to detect when skills change (via grimoire update) and bust the file cache.
func computeRubricHash(rubrics []skillRubric) string {
	h := sha256.New()
	for _, r := range rubrics {
		_, _ = h.Write([]byte(r.Name))
		_, _ = h.Write([]byte(r.Body))
		for _, c := range r.Criteria {
			_, _ = h.Write([]byte(c))
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}

// dedupDiagnostics removes duplicate diagnostics keyed by (URI, start line, code).
func dedupDiagnostics(diags []compliance.Diagnostic) []compliance.Diagnostic {
	type key struct {
		uri  string
		line int
		code string
	}
	seen := make(map[key]bool, len(diags))
	var out []compliance.Diagnostic
	for i := range diags {
		d := diags[i]
		k := key{d.URI, d.Range.Start.Line, d.Code}
		if !seen[k] {
			seen[k] = true
			out = append(out, d)
		}
	}
	return out
}

// callAnthropicAPI calls the Anthropic /v1/messages endpoint.
func callAnthropicAPI(apiKey, model string, maxTokens int, rubric, fileContents string) (string, error) {
	userMsg := rubric
	if fileContents != "" {
		userMsg += "\n\nFile contents to evaluate:\n" + fileContents
	}
	body := map[string]any{
		"model":       model,
		"max_tokens":  maxTokens,
		"temperature": 0,
		"system":      "You are grimoire, an independent compliance checker. Output ONLY a single JSON object matching the compliance report schema in the user message. No prose, no markdown, no explanation — JSON only.",
		"messages":    []map[string]any{{"role": "user", "content": userMsg}},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Error != nil {
		return "", fmt.Errorf("API error: %s", result.Error.Message)
	}
	if len(result.Content) == 0 {
		return "", fmt.Errorf("empty API response")
	}
	return result.Content[0].Text, nil
}

// callOpenAICompatible calls any OpenAI-compatible /v1/chat/completions endpoint.
// If p.APIKey is empty the Authorization header is omitted (e.g. for Ollama).
func callOpenAICompatible(p resolvedProvider, system, userMsg string) (string, error) {
	body := map[string]any{
		"model":       p.Model,
		"max_tokens":  p.MaxTokens,
		"temperature": 0,
		"messages": []map[string]any{
			{"role": "system", "content": system},
			{"role": "user", "content": userMsg},
		},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	url := strings.TrimRight(p.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	if p.APIKey != "" {
		req.Header.Set("authorization", "Bearer "+p.APIKey)
	}
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Error != nil {
		return "", fmt.Errorf("API error: %s", result.Error.Message)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("empty API response")
	}
	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

// callViaLocalAgent pipes prompt to an installed AI CLI and returns its output.
// dir sets the working directory so the CLI's tools resolve paths against the project root.
func callViaLocalAgent(goCtx context.Context, ag, dir, prompt string) (string, error) {
	// Scale timeout with prompt size: base 3m + 1m per 20 KB, capped at 15m.
	promptKB := len(prompt) / 1024
	timeout := 3*time.Minute + time.Duration(promptKB/20+1)*time.Minute
	if timeout > 15*time.Minute {
		timeout = 15 * time.Minute
	}
	ctx, cancel := context.WithTimeout(goCtx, timeout)
	defer cancel()

	var cmd *exec.Cmd
	switch ag {
	case "claude":
		cmd = exec.CommandContext(ctx, "claude", "-p")
	case "copilot":
		cmd = exec.CommandContext(ctx, "gh", "copilot", "suggest", "-t", "shell")
	default:
		cmd = exec.CommandContext(ctx, ag)
	}
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(prompt)
	out, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("%s: timed out after %s", ag, timeout.Round(time.Second))
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && len(exitErr.Stderr) > 0 {
			return "", fmt.Errorf("%s: %s", ag, strings.TrimSpace(string(exitErr.Stderr)))
		}
		return "", fmt.Errorf("%s: %w", ag, err)
	}
	return strings.TrimSpace(string(out)), nil
}

// buildLocalPrompt builds a single prompt string for piping to a local AI CLI.
// The prompt instructs the AI to return a JSON compliance report in the compliance.Report schema.
func buildLocalPrompt(profile, rubric, fileContents string, threshold float64) string {
	var b strings.Builder
	b.WriteString("You are grimoire, an independent compliance checker.\n")
	b.WriteString("Output ONLY a single valid JSON object. No other text, no markdown code blocks, no explanation.\n\n")
	b.WriteString("Derive ALL criteria from each skill's rubric. Do not limit to a fixed count — output as many criteria as the skill defines.\n\n")
	b.WriteString("Required JSON schema (replace placeholder values with your evaluation):\n")
	b.WriteString(`{"version":"1","timestamp":"2026-01-01T00:00:00Z","mode":"independent","scope":"<profile>","result":"pass","coverage":{"overall_pct":87.5,"practices":{"total":1,"passing":1,"partial":0,"failing":0,"coverage_pct":100},"criteria":{"total":6,"passing":4,"failing":2,"suppressed":0,"coverage_pct":66.7},"details":[{"name":"apply-solid-principles","total":6,"passing":4,"partial":0,"failing":2,"coverage_pct":66.7}]},"threshold":{"required":80,"actual":87.5,"status":"pass"},"criteria_matrix":{"apply-solid-principles":{"pass":["no god objects","single responsibility per module","dependency injection used","interfaces over concrete types"],"fail":["all public functions have tests","no circular dependencies"]}},"diagnostics":[{"uri":"src/main.go","range":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}},"severity":2,"code":"missing-test","source":"grimoire","message":"Function Foo has no tests","practice":"apply-solid-principles","criterion":"all public functions have tests","status":"fail"}]}`)
	b.WriteString("\n\n")
	fmt.Fprintf(&b, "Profile: %s\n", profile)
	fmt.Fprintf(&b, "Threshold required: %.0f%%\n\n", threshold)
	b.WriteString(rubric)
	if fileContents != "" {
		b.WriteString("\nFile contents to evaluate:\n")
		b.WriteString(fileContents)
	}
	return b.String()
}

// handleCheckResult parses the AI response, merges cache, writes the report, and enforces thresholds.
// Falls back to printing freeform text on parse failure.
func handleCheckResult(result, projectDir, mode string, cache *compliance.CheckCache, filesToCheck []string, start time.Time) error {
	jsonStr := extractJSON(result)
	var report compliance.Report
	if parseErr := json.Unmarshal([]byte(jsonStr), &report); parseErr == nil &&
		(report.Coverage.Practices.Total > 0 || len(report.Diagnostics) > 0 ||
			report.Result == "pass" || report.Result == "fail" ||
			(cache != nil && len(cache.PracticeTotals) > 0)) {
		report.Timestamp = time.Now().UTC().Format(time.RFC3339)
		if report.Version == "" {
			report.Version = "1"
		}
		if report.Mode == "" {
			report.Mode = mode
		}
		if mode == "incremental" {
			report.Scope = "changed"
		} else {
			report.Scope = "."
		}
		if report.Mode == "incremental" && report.Git.BaseRef == "" {
			report.Git.BaseRef = "HEAD"
		}

		// Merge cached diagnostics for files not re-checked this run.
		// Use processedFiles (filesToCheck) as exclusion set, not just checkedURIs:
		// deleted files produce no diagnostics so they'd be absent from checkedURIs,
		// but their stale cached diagnostics must not be merged either.
		if cache != nil && mode == "incremental" {
			processedFiles := make(map[string]bool, len(filesToCheck))
			for _, f := range filesToCheck {
				processedFiles[f] = true
			}
			checkedURIs := make(map[string]bool)
			for i := range report.Diagnostics {
				d := report.Diagnostics[i]
				checkedURIs[d.URI] = true
			}
			for file, entry := range cache.Files {
				if !checkedURIs[file] && !processedFiles[file] {
					report.Diagnostics = append(report.Diagnostics, entry.Diagnostics...)
				}
			}
		}

		// On full scan, store practice criteria totals for coverage recomputation on later incremental runs.
		if cache != nil && mode == "full" && len(report.Coverage.Details) > 0 {
			cache.PracticeTotals = make(map[string]int, len(report.Coverage.Details))
			cache.PracticeCriteria = make(map[string][]compliance.CriterionDetail, len(report.Coverage.Details))
			for _, d := range report.Coverage.Details {
				cache.PracticeTotals[d.Name] = d.Total
				if len(d.Criteria) > 0 {
					cache.PracticeCriteria[d.Name] = d.Criteria
				}
			}
			// Primary: criteria_matrix from AI response.
			for practice, entry := range report.CriteriaMatrix {
				cache.PracticeCriteria[practice] = nil
				for _, n := range entry.Pass {
					cache.PracticeCriteria[practice] = append(cache.PracticeCriteria[practice],
						compliance.CriterionDetail{Name: n, Status: "pass"})
				}
				for _, n := range entry.Fail {
					cache.PracticeCriteria[practice] = append(cache.PracticeCriteria[practice],
						compliance.CriterionDetail{Name: n, Status: "fail"})
				}
			}
			// Fallback: severity-4 pass diagnostics.
			for i := range report.Diagnostics {
				d := report.Diagnostics[i]
				if d.Severity == 4 && d.Status == "pass" && d.Practice != "" && d.Criterion != "" {
					cache.PracticeCriteria[d.Practice] = append(cache.PracticeCriteria[d.Practice],
						compliance.CriterionDetail{Name: d.Criterion, Status: "pass"})
				}
			}
		}

		// Recompute coverage from actual diagnostics + cached criteria totals.
		// Runs for both full and incremental — cache.PracticeTotals is populated above for full scan.
		if cache != nil && len(cache.PracticeTotals) > 0 {
			report.Coverage = recomputeCoverage(report.Diagnostics, cache.PracticeTotals, cache.PracticeCriteria)
			report.Threshold.Actual = report.Coverage.OverallPct
			if report.Threshold.Required > 0 {
				if report.Coverage.OverallPct >= report.Threshold.Required {
					report.Threshold.Status = "pass"
					report.Result = "pass"
				} else {
					report.Threshold.Status = "fail"
					report.Result = "fail"
				}
			}
		}

		report.Summary = computeSummary(&report)

		// Sort diagnostics for deterministic JSON and HTML output regardless of merge order.
		sort.Slice(report.Diagnostics, func(i, j int) bool {
			a, b := &report.Diagnostics[i], &report.Diagnostics[j]
			if a.URI != b.URI {
				return a.URI < b.URI
			}
			if a.Range.Start.Line != b.Range.Start.Line {
				return a.Range.Start.Line < b.Range.Start.Line
			}
			if a.Range.Start.Character != b.Range.Start.Character {
				return a.Range.Start.Character < b.Range.Start.Character
			}
			if a.Severity != b.Severity {
				return a.Severity < b.Severity
			}
			return a.Code < b.Code
		})

		reportPath, writeErr := compliance.WriteReport(&report, projectDir)
		if writeErr != nil {
			fmt.Fprintf(os.Stderr, "  warn: writing report: %v\n", writeErr)
		} else {
			if !flagLive {
				fmt.Printf("  Report: %s\n", reportPath)
			}
			if htmlPath, htmlErr := compliance.WriteHTMLReport(reportPath, cliVersion, projectDir); htmlErr == nil {
				if !flagLive {
					fmt.Printf("  HTML:   %s\n", htmlPath)
				}
			}
		}

		// Update cache entries for files just checked.
		if cache != nil {
			cfgHash, _ := compliance.ConfigHash(projectDir)
			cache.ConfigHash = cfgHash
			byFile := make(map[string][]compliance.Diagnostic)
			for i := range report.Diagnostics {
				d := report.Diagnostics[i]
				byFile[d.URI] = append(byFile[d.URI], d)
			}
			for _, f := range filesToCheck {
				absPath := filepath.Join(projectDir, f)
				info, statErr := os.Stat(absPath)
				if statErr != nil {
					delete(cache.Files, f) // file deleted — evict stale entry
					continue
				}
				hash, hashErr := compliance.FileHash(absPath)
				if hashErr != nil {
					delete(cache.Files, f)
					continue
				}
				cache.Files[f] = compliance.FileCacheEntry{
					Hash:        hash,
					Mtime:       info.ModTime().UnixNano(),
					Size:        info.Size(),
					Diagnostics: byFile[f],
				}
			}
			_ = compliance.SaveCache(cache, projectDir)
		}

		// Apply structural rules engine.
		eng := &rules.Engine{
			SkillsPackages: skills.AllSkillsPackages(),
			ProjectDir:     projectDir,
		}
		if found := eng.Run(); len(found) > 0 {
			report.Diagnostics = append(found, report.Diagnostics...)
		}

		if flagJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(report)
		}

		var doneDiags []compliance.Diagnostic
		if mode == "incremental" && len(filesToCheck) > 0 {
			checkedSet := make(map[string]bool, len(filesToCheck))
			for _, f := range filesToCheck {
				checkedSet[f] = true
			}
			for i := range report.Diagnostics {
				d := report.Diagnostics[i]
				if checkedSet[d.URI] {
					doneDiags = append(doneDiags, d)
				}
			}
		} else {
			doneDiags = report.Diagnostics
		}
		nErrors := len(filterBySeverity(doneDiags, 1))
		nWarnings := len(filterBySeverity(doneDiags, 2))
		elapsed := time.Since(start)
		doneIcon := colorize(ansiGreen, "✓")
		if nErrors > 0 {
			doneIcon = colorize(ansiRed, "✗")
		}
		fmt.Printf("  %s Done in %.1fs — %d errors, %d warnings\n\n", doneIcon, elapsed.Seconds(), nErrors, nWarnings)

		ciMode := flagCI || os.Getenv("GITHUB_ACTIONS") == "true"
		printSummary(&report, mode, filesToCheck)
		if ciMode {
			emitGHAAnnotations(&report)
		}
		if flagJUnit != "" {
			if err := writeJUnitXML(&report, flagJUnit); err != nil {
				fmt.Fprintf(os.Stderr, "  warn: writing JUnit XML: %v\n", err)
			} else {
				fmt.Printf("  JUnit XML written to %s\n", flagJUnit)
			}
		}

		// Threshold enforcement — return error (not os.Exit) so grimoire watch can continue.
		resolved, _ := config.Load(projectDir)
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
			return fmt.Errorf("compliance threshold not met")
		}
		return nil
	}
	// Fallback: AI returned freeform text or an empty/invalid evaluation.
	fmt.Println(result)
	fmt.Fprintln(os.Stderr, "\n  (AI returned no evaluation data — run again to retry)")
	return nil
}

// extractJSON strips markdown code fences and returns the outermost {...} block.
func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.Index(s, "```"); i >= 0 {
		s = s[i+3:]
		if j := strings.Index(s, "\n"); j >= 0 {
			s = s[j+1:]
		}
		if k := strings.Index(s, "```"); k >= 0 {
			s = s[:k]
		}
		s = strings.TrimSpace(s)
	}
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}
