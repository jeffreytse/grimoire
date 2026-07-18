package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var (
	flagValidateJSON       bool
	flagValidateStrict     bool
	flagValidateTestSchema string
	flagValidateFix        bool
	flagValidateVia        string
	flagValidatePreferAPI  bool
	flagValidateNoDuplicates bool
)

var validateCmd = &cobra.Command{
	Use:   "validate [<skill-path>...]",
	Short: "Validate skill files against the grimoire STANDARD",
	Long: `Validate SKILL.md files for conformance with the grimoire STANDARD.

Each path can be a SKILL.md file, a directory, or a glob pattern.
Defaults to the current working directory when no paths are given.

Exit codes: 0 = all pass, 1 = errors found (or warnings with --strict).`,
	RunE: runValidate,
}

func init() {
	validateCmd.Flags().BoolVar(&flagValidateJSON, "json", false, "output as JSON")
	validateCmd.Flags().BoolVar(&flagValidateStrict, "strict", false, "treat warnings as errors (non-zero exit)")
	validateCmd.Flags().StringVar(&flagValidateTestSchema, "test-schema", "", "run conformance test against a schema/tests/ directory")
	validateCmd.Flags().BoolVar(&flagValidateFix, "fix", false, "attempt to fix validation errors via AI (requires local AI agent or API key)")
	validateCmd.Flags().StringVar(&flagValidateVia, "via", "", "force a specific local AI agent for --fix (claude, gemini, codex, …)")
	validateCmd.Flags().BoolVar(&flagValidatePreferAPI, "prefer-api", false, "use API provider instead of local AI agent for --fix")
	validateCmd.Flags().BoolVar(&flagValidateNoDuplicates, "no-duplicates", false, "skip semantic near-duplicate detection")
}

type validateFinding struct {
	Level   string `json:"level"` // "error" | "warn"
	Message string `json:"message"`
}

type validateResult struct {
	Path              string            `json:"path"`
	Name              string            `json:"name"`
	Status            string            `json:"status"` // "pass" | "fail" | "warn"
	Findings          []validateFinding `json:"findings"`
	DuplicateReviewed bool              `json:"-"` // from frontmatter; used by --duplicates pass
}

type validateOutput struct {
	OK     bool             `json:"ok"`
	Skills []validateResult `json:"skills"`
}

func runValidate(cmd *cobra.Command, args []string) error {
	if flagValidateTestSchema != "" {
		return runTestSchema(cmd, flagValidateTestSchema)
	}

	projectDir := getProjectDir()

	var skillMDPaths []string
	seen := map[string]bool{}
	addPath := func(p string) {
		abs, err := filepath.Abs(p)
		if err != nil {
			abs = p
		}
		if !seen[abs] {
			seen[abs] = true
			skillMDPaths = append(skillMDPaths, abs)
		}
	}

	if len(args) == 0 {
		// Default to CWD — mirrors how go vet, eslint, etc. behave.
		args = []string{"."}
	}

	for _, arg := range args {
		if isGlobPattern(arg) {
			found, err := resolveGlobPattern(arg)
			if err != nil {
				return err
			}
			if len(found) == 0 {
				fmt.Fprintf(os.Stderr, "warn: no SKILL.md files matched %q\n", arg)
			}
			for _, p := range found {
				addPath(p)
			}
		} else {
			p, err := resolveSkillMDPath(arg)
			if err != nil {
				// Fallback: if arg is a directory, walk it recursively.
				info, serr := os.Stat(arg)
				if serr != nil || !info.IsDir() {
					return err
				}
				found, werr := skills.WalkSkillFiles(arg)
				if werr != nil {
					return fmt.Errorf("walking %s: %w", arg, werr)
				}
				if len(found) == 0 {
					fmt.Fprintf(os.Stderr, "warn: no SKILL.md files found under %q\n", arg)
				}
				for _, fp := range found {
					addPath(fp)
				}
			} else {
				addPath(p)
			}
		}
	}

	if len(skillMDPaths) == 0 {
		fmt.Printf("%s  no SKILL.md files found\n", tui.IconSkip)
		return nil
	}

	if !flagValidateJSON {
		fmt.Printf("Validating %d skill(s)...\n\n", len(skillMDPaths))
	}

	results := make([]validateResult, len(skillMDPaths))
	concurrency := runtime.GOMAXPROCS(0) * 2
	if concurrency > 32 {
		concurrency = 32
	}
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	for i, p := range skillMDPaths {
		wg.Add(1)
		go func(i int, p string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			results[i] = validateSkill(p)
		}(i, p)
	}
	wg.Wait()

	// Duplicate name check — flag all occurrences, including when 3+ share the same name.
	if len(skillMDPaths) > 1 {
		nameToIdxs := map[string][]int{}
		for i, r := range results {
			if r.Name != "" {
				nameToIdxs[r.Name] = append(nameToIdxs[r.Name], i)
			}
		}
		for _, idxs := range nameToIdxs {
			if len(idxs) < 2 {
				continue
			}
			for _, i := range idxs {
				dup := validateFinding{"error", fmt.Sprintf("name: '%s' duplicated at %d locations", results[i].Name, len(idxs))}
				results[i].Findings = append(results[i].Findings, dup)
				results[i].Status = "fail"
			}
		}
	}

	// ── Semantic near-duplicate check ─────────────────────────────────────────
	if !flagValidateNoDuplicates {
		corpusRoot := ""
		if len(skillMDPaths) > 0 {
			// Infer from skill paths: first 'skills/' component in path is the corpus root.
			// Falls back to installed package root when skills are not inside a known structure.
			corpusRoot = inferCorpusRoot(skillMDPaths[0])
		}
		if corpusRoot == "" {
			corpusRoot = skills.SkillsRoot()
		}
		if abs, err := filepath.Abs(corpusRoot); err == nil {
			corpusRoot = abs
		}
		allPaths, _ := skills.WalkSkillFiles(corpusRoot)

		if !flagValidateJSON && len(allPaths) > 0 {
			fmt.Printf("Checking duplicates against %d-skill corpus...\n\n", len(allPaths))
		}

		// Build corpus: precompute tag/desc sets once per entry to avoid
		// per-comparison allocations in the O(n×m) scoring loop.
		type dupEntry struct {
			name    string
			path    string
			tagSet  map[string]bool
			descSet map[string]bool
			domain  string
		}
		corpus := make([]dupEntry, len(allPaths))
		var cwg sync.WaitGroup
		csem := make(chan struct{}, concurrency)
		for ci, cp := range allPaths {
			cwg.Add(1)
			go func(ci int, cp string) {
				defer cwg.Done()
				csem <- struct{}{}
				defer func() { <-csem }()
				sk, _ := skills.ParseSkillFile(cp)
				ts := make(map[string]bool, len(sk.Tags))
				for _, t := range sk.Tags {
					ts[strings.ToLower(t)] = true
				}
				ds := make(map[string]bool)
				for _, tok := range dupTokenize(sk.Description) {
					ds[tok] = true
				}
				corpus[ci] = dupEntry{sk.Name, sk.Path, ts, ds, skillDomain(sk.Path)}
			}(ci, cp)
		}
		cwg.Wait()

		// Parallelize the outer loop: each goroutine scores one validated skill
		// against the full corpus and appends findings to results[i].
		type nearMatch struct {
			name  string
			score float64
		}
		var dupWg sync.WaitGroup
		dupSem := make(chan struct{}, concurrency)
		for i := range results {
			dupWg.Add(1)
			go func(i int) {
				defer dupWg.Done()
				dupSem <- struct{}{}
				defer func() { <-dupSem }()

				validatedSk, _ := skills.ParseSkillFile(skillMDPaths[i])
				validatedDir := filepath.Dir(skillMDPaths[i])
				vTagSet := make(map[string]bool, len(validatedSk.Tags))
				for _, t := range validatedSk.Tags {
					vTagSet[strings.ToLower(t)] = true
				}
				vDescSet := make(map[string]bool)
				for _, tok := range dupTokenize(validatedSk.Description) {
					vDescSet[tok] = true
				}
				vDomain := skillDomain(validatedDir)

				validatedName := results[i].Name
				var matches []nearMatch
				for _, e := range corpus {
					if e.path == validatedDir || e.name == validatedName {
						continue
					}
					score := dupScoreSets(vTagSet, vDescSet, vDomain, e.tagSet, e.descSet, e.domain)
					if score >= 0.5 {
						matches = append(matches, nearMatch{e.name, score})
					}
				}
				sort.Slice(matches, func(a, b int) bool { return matches[a].score > matches[b].score })
				if len(matches) > 3 {
					matches = matches[:3]
				}
				for _, m := range matches {
					var f validateFinding
					if m.score >= 0.7 {
						if results[i].DuplicateReviewed {
							f = validateFinding{"warn", fmt.Sprintf("near-duplicate (reviewed): '%s' (score=%.2f)", m.name, m.score)}
						} else {
							f = validateFinding{"error", fmt.Sprintf("near-duplicate: '%s' (score=%.2f) — extend existing skill or add duplicate-reviewed: true", m.name, m.score)}
							results[i].Status = "fail"
						}
					} else {
						f = validateFinding{"warn", fmt.Sprintf("possible near-duplicate: '%s' (score=%.2f)", m.name, m.score)}
					}
					results[i].Findings = append(results[i].Findings, f)
					if results[i].Status == "pass" {
						results[i].Status = "warn"
					}
				}
			}(i)
		}
		dupWg.Wait()
	}

	// ── AI fix pass ───────────────────────────────────────────────────────────
	if flagValidateFix {
		type fixCandidate struct {
			idx    int
			result validateResult
		}
		var candidates []fixCandidate
		for i, r := range results {
			if r.Status == "pass" {
				continue
			}
			// Warn-only skills are acceptable — skip unless --strict treats them as errors.
			if r.Status == "warn" && !flagValidateStrict {
				continue
			}
			candidates = append(candidates, fixCandidate{i, r})
		}

		if !flagValidateJSON {
			fmt.Printf("%s\n\n", tui.StyleBold.Render("Running AI fix pass..."))
		}

		if len(candidates) > 0 {
			goCtx := context.Background()
			cwd, _ := os.Getwd()
			ex := resolveExecutorFor(projectDir, flagValidateVia, flagValidatePreferAPI)

			// Concurrency cap: fewer goroutines for local CLI (process spawn overhead).
			fixConcurrency := 15
			if ex.Kind == execLocalCLI {
				fixConcurrency = 5
			}

			type fixWork struct {
				updated validateResult
				buf     strings.Builder
				fixed   bool
				skipped bool
			}
			work := make([]fixWork, len(candidates))

			var mu sync.Mutex
			sem := make(chan struct{}, fixConcurrency)
			var wg sync.WaitGroup
			for j := range candidates {
				wg.Add(1)
				go func(j int) {
					defer wg.Done()
					sem <- struct{}{}
					defer func() { <-sem }()

					cand := candidates[j]
					w := &work[j]
					r := cand.result

					// Flush this skill's output block as soon as the goroutine exits,
					// regardless of which return path is taken.
					defer func() {
						if !flagValidateJSON {
							mu.Lock()
							fmt.Print(w.buf.String())
							mu.Unlock()
						}
					}()

					rel, err := filepath.Rel(cwd, r.Path)
					if err != nil {
						rel = r.Path
					}
					fmt.Fprintf(&w.buf, "  %s %s\n", tui.StyleCyan.Render("→"), tui.StyleBold.Render(r.Name))
					fmt.Fprintf(&w.buf, "    %s\n", tui.StyleDim.Render(rel))
					for _, f := range r.Findings {
						if f.Level == "warn" {
							fmt.Fprintf(&w.buf, "    %s %s\n", tui.IconWarn, f.Message)
						} else {
							fmt.Fprintf(&w.buf, "    %s %s\n", tui.IconFail, f.Message)
						}
					}

					// ── deterministic pre-fix ────────────────────────────────
					if raw, readErr := os.ReadFile(r.Path); readErr == nil {
						if fixed, changed := deterministicFix(string(raw), r.Findings); changed {
							if writeErr := os.WriteFile(r.Path, []byte(fixed), 0o644); writeErr == nil {
								r = validateSkill(r.Path)
								fmt.Fprintf(&w.buf, "    %s %s\n", tui.StyleCyan.Render("⚡"), tui.StyleDim.Render("deterministic fixes applied"))
								if r.Status == "pass" {
									w.updated = r
									w.fixed = true
									fmt.Fprintf(&w.buf, "    %s %s\n\n", tui.IconOK, tui.StyleGreen.Render("fixed (deterministic)"))
									return
								}
							}
						}
					}

					// ── AI fix ───────────────────────────────────────────────
					if ex.Kind == execPrint {
						fmt.Fprintf(&w.buf, "    %s %s\n\n", tui.IconFail, tui.StyleRed.Render("skipped: no AI executor — install a local agent or set an API key"))
						w.updated = r
						w.skipped = true
						return
					}
					aiFindings := fixableFindings(r.Findings)
					if len(aiFindings) == 0 {
						fmt.Fprintf(&w.buf, "    %s %s\n\n", tui.IconSkip, tui.StyleDim.Render("no AI-fixable findings — manual intervention required"))
						w.updated = r
						w.skipped = true
						return
					}
					rToFix := r
					rToFix.Findings = aiFindings
					if err := fixSkillWithAI(goCtx, rToFix, projectDir, &ex); err != nil {
						fmt.Fprintf(&w.buf, "    %s %s\n\n", tui.IconFail, tui.StyleRed.Render("skipped: "+err.Error()))
						w.updated = r
						w.skipped = true
						return
					}
					updated := validateSkill(r.Path)
					w.updated = updated
					w.fixed = true
					if updated.Status == "pass" {
						fmt.Fprintf(&w.buf, "    %s %s\n\n", tui.IconOK, tui.StyleGreen.Render("fixed"))
					} else {
						fmt.Fprintf(&w.buf, "    %s %s\n\n", tui.IconWarn, tui.StyleYellow.Render(fmt.Sprintf("%d finding(s) remain", len(updated.Findings))))
					}

				}(j)
			}
			wg.Wait()

			// Apply results in order; output was already printed live inside each goroutine.
			fixedCount, skippedCount := 0, 0
			for j, cand := range candidates {
				w := &work[j]
				results[cand.idx] = w.updated
				if w.fixed {
					fixedCount++
				}
				if w.skipped {
					skippedCount++
				}
			}
			if !flagValidateJSON {
				fmt.Printf("%s %d fixed, %d skipped\n\n",
					tui.StyleBold.Render("Fix summary:"), fixedCount, skippedCount)
			}
		}
	}

	out := validateOutput{OK: true}
	passed, failed, warnOnly := 0, 0, 0
	for _, r := range results {
		switch r.Status {
		case "fail":
			out.OK = false
			failed++
		case "warn":
			warnOnly++
			if flagValidateStrict {
				out.OK = false
			}
		default:
			passed++
		}
		out.Skills = append(out.Skills, r)
	}

	if flagValidateJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			return fmt.Errorf("writing JSON output: %w", err)
		}
	} else {
		printValidateHuman(results)
		fmt.Printf("Summary: %d passed", passed)
		if failed > 0 {
			fmt.Printf(", %d failed", failed)
		}
		if warnOnly > 0 {
			fmt.Printf(", %d warning-only", warnOnly)
		}
		fmt.Println()
	}

	if !out.OK {
		// Suppress cobra's "Error: ..." print — we've already shown our summary.
		cmd.Root().SilenceErrors = true
		return fmt.Errorf("validation failed")
	}
	return nil
}

func printValidateHuman(results []validateResult) {
	for _, r := range results {
		switch r.Status {
		case "pass":
			fmt.Printf("%s  %s\n", tui.IconOK, r.Name)
		case "fail":
			errs, warns := 0, 0
			for _, f := range r.Findings {
				if f.Level == "error" {
					errs++
				} else {
					warns++
				}
			}
			label := fmt.Sprintf("%d error(s)", errs)
			if warns > 0 {
				label += fmt.Sprintf(", %d warning(s)", warns)
			}
			fmt.Printf("%s  %s (%s)\n", tui.IconFail, r.Name, label)
			for _, f := range r.Findings {
				tag := "[FAIL]"
				if f.Level == "warn" {
					tag = "[WARN]"
				}
				fmt.Printf("   %s %s\n", tag, f.Message)
			}
		case "warn":
			fmt.Printf("%s  %s (%d warning(s))\n", tui.IconWarn, r.Name, len(r.Findings))
			for _, f := range r.Findings {
				fmt.Printf("   [WARN] %s\n", f.Message)
			}
		}
	}
	fmt.Println()
}

// validateSkill runs all STANDARD checks against a single SKILL.md path.
func validateSkill(skillMDPath string) validateResult {
	name := filepath.Base(filepath.Dir(skillMDPath))
	result := validateResult{Path: skillMDPath, Name: name, Status: "pass"}

	raw, err := os.ReadFile(skillMDPath)
	if err != nil {
		result.Status = "fail"
		result.Findings = append(result.Findings, validateFinding{"error", "cannot read file: " + err.Error()})
		return result
	}
	content := string(raw)
	lowerContent := strings.ToLower(content)
	lineCount := strings.Count(content, "\n") + 1

	fail := func(msg string) {
		result.Findings = append(result.Findings, validateFinding{"error", msg})
		result.Status = "fail"
	}
	warn := func(msg string) {
		result.Findings = append(result.Findings, validateFinding{"warn", msg})
		if result.Status == "pass" {
			result.Status = "warn"
		}
	}

	// Check frontmatter structure (cheap string ops) then YAML validity via single parse.
	yamlOK := true
	if !strings.HasPrefix(content, "---") {
		fail("frontmatter: file must start with '---'")
		yamlOK = false
	} else if !strings.Contains(content[3:], "\n---") {
		fail("frontmatter: missing closing '---' delimiter")
		yamlOK = false
	}

	// Parse from already-read content — single YAML parse for both validity check and field access.
	sk, yamlErr := skills.ParseSkillFromContent(content, filepath.Dir(skillMDPath))
	if yamlOK && yamlErr != nil {
		fail("frontmatter: invalid YAML (check for unquoted ':' in field values) — " + yamlErr.Error())
		yamlOK = false
	}
	if sk.Name != "" {
		result.Name = sk.Name
	}
	result.DuplicateReviewed = sk.DuplicateReviewed

	// Field-presence and content checks only make sense when YAML parsed cleanly.
	if !yamlOK {
		// Skip to size check — YAML error already reported above.
		if lineCount < 50 {
			fail(fmt.Sprintf("size: only %d lines — minimum is 50", lineCount))
		} else if lineCount > 300 {
			fail(fmt.Sprintf("size: %d lines — exceeds 300-line limit, consider splitting", lineCount))
		}
		return result
	}

	// ── name ──────────────────────────────────────────────────────────────────
	if sk.Name == "" {
		fail("name: field missing")
	} else {
		if !validateNameRe.MatchString(sk.Name) {
			fail("name: must be lowercase kebab-case (e.g. 'review-pull-request')")
		}
		if len(sk.Name) > 50 {
			fail(fmt.Sprintf("name: '%s' exceeds 50 characters (%d)", sk.Name, len(sk.Name)))
		}
		if !isMetaSkill(skillMDPath) {
			for _, bad := range validateRejectedPrefixes {
				if strings.HasPrefix(sk.Name, bad) {
					fail(fmt.Sprintf("name: prefix '%s' is redundant in a skill library", bad))
				}
			}
			firstVerb := strings.SplitN(sk.Name, "-", 2)[0]
			if validateRejectedVerbs[firstVerb] {
				fail(fmt.Sprintf("name: verb '%s' is a rejected verb (too vague)", firstVerb))
			} else if !validateApprovedVerbs[firstVerb] {
				warn(fmt.Sprintf("name: verb '%s' not in approved verb list", firstVerb))
			}
		}
	}

	// ── description ───────────────────────────────────────────────────────────
	if sk.Description == "" {
		fail("description: field missing")
	} else {
		if !strings.HasPrefix(sk.Description, "Use when") {
			fail("description: must start with 'Use when'")
		}
		if len(sk.Description) > 500 {
			fail(fmt.Sprintf("description: exceeds 500 characters (%d)", len(sk.Description)))
		}
	}

	// ── source ────────────────────────────────────────────────────────────────
	if strings.TrimSpace(sk.Source) == "" {
		fail("source: field missing")
	}

	// ── tags ──────────────────────────────────────────────────────────────────
	switch {
	case len(sk.Tags) == 0:
		fail("tags: field missing")
	case len(sk.Tags) < 3:
		fail(fmt.Sprintf("tags: need at least 3 tags, found %d", len(sk.Tags)))
	case len(sk.Tags) > 8:
		warn(fmt.Sprintf("tags: %d tags — consider trimming to ≤8", len(sk.Tags)))
	}
	for _, t := range sk.Tags {
		if !validateTagRe.MatchString(t) {
			switch {
			case t != strings.ToLower(t):
				fail(fmt.Sprintf("tags: '%s' must be lowercase (use only a-z, 0-9, hyphens)", t))
			default:
				fail(fmt.Sprintf("tags: '%s' must contain only lowercase letters, digits, hyphens", t))
			}
		}
	}

	// ── related ───────────────────────────────────────────────────────────────
	for _, r := range sk.Related {
		if !validateNameRe.MatchString(r) {
			warn(fmt.Sprintf("related: '%s' is not a valid skill name (must be kebab-case)", r))
		}
	}

	// ── body sections ─────────────────────────────────────────────────────────
	body := sk.Body
	if !strings.HasPrefix(body, "# ") && !strings.Contains(body, "\n# ") {
		fail("body: missing '# Title' h1 heading")
	}
	if !strings.Contains(body, "## Why This Is Best Practice") {
		fail("body: missing '## Why This Is Best Practice' section")
	}
	if !strings.Contains(body, "## Steps") && !strings.Contains(body, "## Core Pattern") {
		fail("body: missing '## Steps' or '## Core Pattern' section")
	}

	// ── emerging status marker ─────────────────────────────────────────────────
	if sk.Emerging && !strings.Contains(body, "**Status:** Emerging") {
		fail("emerging: true but missing '**Status:** Emerging' in Why section")
	}

	// ── lifecycle conflicts ───────────────────────────────────────────────────
	if sk.Practitioner && sk.Stable {
		fail("lifecycle: practitioner: true and stable: true are mutually exclusive")
	}
	if sk.Stable && sk.Emerging {
		fail("lifecycle: stable: true and emerging: true are mutually exclusive")
	}
	if sk.Deprecated && sk.Emerging {
		fail("lifecycle: deprecated: true and emerging: true are mutually exclusive")
	}
	if sk.Deprecated && sk.DeprecatedBy == "" {
		fail("lifecycle: deprecated: true requires deprecated_by: to be set")
	}
	if sk.DeprecatedBy != "" && !sk.Deprecated {
		warn("lifecycle: deprecated_by: is set but deprecated: true is missing")
	}

	// ── Why section content (non-emerging only) ───────────────────────────────
	if !sk.Emerging {
		if !strings.Contains(body, "**Adopted by:**") {
			fail("body: '## Why This Is Best Practice' missing '**Adopted by:**' line")
		}
		if !strings.Contains(body, "**Impact:**") {
			fail("body: '## Why This Is Best Practice' missing '**Impact:**' line")
		}
		if !strings.Contains(body, "**Why best:**") {
			fail("body: '## Why This Is Best Practice' missing '**Why best:**' line")
		}
	}

	// ── size ──────────────────────────────────────────────────────────────────
	if lineCount < 50 {
		fail(fmt.Sprintf("size: only %d lines — minimum is 50", lineCount))
	} else if lineCount > 300 {
		fail(fmt.Sprintf("size: %d lines — exceeds 300-line limit, consider splitting", lineCount))
	}

	// ── domain safety footer ──────────────────────────────────────────────────
	switch skillDomain(skillMDPath) {
	case "health":
		if !strings.Contains(lowerContent, "healthcare provider") && !strings.Contains(lowerContent, "medical advice") {
			fail("safety: health skill missing required disclaimer (must mention 'healthcare provider' or 'medical advice')")
		}
	case "law":
		if !strings.Contains(lowerContent, "legal advice") && !strings.Contains(lowerContent, "legal counsel") {
			fail("safety: law skill missing required disclaimer (must mention 'legal advice' or 'legal counsel')")
		}
	case "finance":
		if !strings.Contains(lowerContent, "financial advice") && !strings.Contains(lowerContent, "financial advisor") {
			fail("safety: finance skill missing required disclaimer (must mention 'financial advice' or 'financial advisor')")
		}
	case "psychology":
		if !strings.Contains(lowerContent, "mental health professional") {
			fail("safety: psychology skill missing required disclaimer (must mention 'mental health professional')")
		}
	}

	// ── title matches skill name (Title Case) ─────────────────────────────────
	if h1, ok := extractH1(body); ok {
		title := strings.TrimPrefix(h1, "# ")
		if sk.Name != "" {
			nameWords := strings.Split(sk.Name, "-")
			titleWords := normalizeH1Words(title)

			nameSet := make(map[string]bool, len(nameWords))
			for _, w := range nameWords {
				nameSet[w] = true
			}
			var filtered []string
			for _, w := range titleWords {
				if !h1Stopwords[w] || nameSet[w] {
					filtered = append(filtered, w)
				}
			}

			if !h1TitleMatchesName(nameWords, filtered) {
				if len(nameWords) != len(filtered) {
					warn(fmt.Sprintf("body: h1 '%s' word count doesn't match skill name '%s'", h1, sk.Name))
				} else {
					warn(fmt.Sprintf("body: h1 '%s' doesn't match skill name '%s' (Title Case expected)", h1, sk.Name))
				}
			}
		}
	}

	// ── purpose statement after title ────────────────────────────────────────
	if !hasPurposeStatement(body) {
		warn("body: missing one-sentence purpose statement after '# Title' heading")
	}

	return result
}

// h1Stopwords lists English articles and prepositions that may appear in an H1
// title but are omitted from the kebab-case skill name (e.g. "Apply the Exposure
// Triangle" for "apply-exposure-triangle"). A stopword is only filtered when it is
// absent from the skill name itself, so "proof-of-concept" still matches correctly.
var h1Stopwords = map[string]bool{
	"a": true, "an": true, "the": true,
	"of": true, "in": true, "for": true, "on": true,
	"with": true, "at": true, "by": true, "from": true,
	"to": true, "and": true, "or": true, "but": true, "vs": true,
}

// versionZeroRe matches version tokens like "2.0" or "1.00" where the minor part is
// all zeros. normalizeH1Words uses it to strip the minor suffix so "OAuth 2.0" →
// tokens ["oauth","2"] rather than ["oauth","20"], enabling matching against "oauth2".
var versionZeroRe = regexp.MustCompile(`^\d+\.0+$`)

var (
	validateNameRe           = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)+$`)
	validateTagRe            = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)
	validateRejectedPrefixes = []string{"skill-", "best-practice-", "guide-"}

	// Deterministic fix regexes (compiled once, reused across all skills).
	deterministicNameRe        = regexp.MustCompile(`(?m)^(name:\s*)(\S+)(.*)$`)
	deterministicYAMLQuotingRe = regexp.MustCompile(`(?m)^((description|source):\s*)(.+)$`)
	deterministicTagsRe        = regexp.MustCompile(`(?m)^(tags:\s*\[)([^\]]+)(\])`)
)

var validateRejectedVerbs = map[string]bool{
	"do": true, "handle": true, "manage": true, "improve": true,
	"set": true, "get": true, "use": true, "help": true,
}

var validateApprovedVerbs = map[string]bool{
	"propose": true, "write": true, "review": true, "audit": true,
	"design": true, "calculate": true, "diagnose": true, "optimize": true,
	"suggest": true, "deprecate": true, "plan": true, "negotiate": true,
	"apply": true, "prevent": true, "profile": true, "validate": true,
	"run": true, "refactor": true, "build": true, "delegate": true,
	"give": true, "resolve": true, "bisect": true, "triage": true,
	"configure": true, "fix": true, "share": true, "teach": true,
	"explain": true, "compare": true, "adapt": true, "analyze": true,
	"discover": true, "install": true, "check": true, "pin": true,
	"start": true, "with": true, "create": true,
}

// isGlobPattern reports whether s contains any glob metacharacter.
func isGlobPattern(s string) bool {
	return strings.ContainsAny(s, "*?{[")
}

// resolveGlobPattern expands a glob pattern into matching SKILL.md paths.
// The pattern is resolved relative to cwd.
func resolveGlobPattern(pattern string) ([]string, error) {
	absPattern := pattern
	if !filepath.IsAbs(pattern) {
		cwd, _ := os.Getwd()
		absPattern = filepath.Join(cwd, pattern)
	}
	absPattern = filepath.ToSlash(absPattern)

	root := globWalkRoot(absPattern)
	found, err := skills.WalkSkillFiles(root)
	if err != nil {
		return nil, fmt.Errorf("walking %s: %w", root, err)
	}

	var matched []string
	for _, p := range found {
		if ok, _ := doublestar.Match(absPattern, filepath.ToSlash(p)); ok {
			matched = append(matched, p)
		}
	}
	return matched, nil
}

// globWalkRoot returns the longest path prefix before the first glob component.
func globWalkRoot(absSlashPattern string) string {
	parts := strings.Split(absSlashPattern, "/")
	clean := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.ContainsAny(p, "*?{[") {
			break
		}
		clean = append(clean, p)
	}
	if len(clean) == 0 {
		return "/"
	}
	return strings.Join(clean, "/")
}

// resolveSkillMDPath normalizes a user-supplied path to a SKILL.md file path.
func resolveSkillMDPath(arg string) (string, error) {
	if filepath.Base(arg) == "SKILL.md" {
		if _, err := os.Stat(arg); err != nil {
			return "", fmt.Errorf("%s: %w", arg, err)
		}
		return arg, nil
	}
	p := filepath.Join(arg, "SKILL.md")
	if _, err := os.Stat(p); err != nil {
		return "", fmt.Errorf("%s: no SKILL.md found", arg)
	}
	return p, nil
}

// isMetaSkill returns true when the skill path includes a 'meta' directory component.
func isMetaSkill(path string) bool {
	for _, part := range strings.Split(filepath.ToSlash(path), "/") {
		if part == "meta" {
			return true
		}
	}
	return false
}

// normalizeH1Words splits an h1 title string by spaces and hyphens, strips
// non-alphanumeric chars (e.g. apostrophes in "Don't"), and returns lowercase
// tokens. This handles hyphenated compounds like "Show-Don't-Tell" correctly.
func normalizeH1Words(title string) []string {
	var words []string
	for _, part := range strings.FieldsFunc(title, func(r rune) bool { return r == ' ' || r == '-' }) {
		if versionZeroRe.MatchString(part) {
			part = part[:strings.IndexByte(part, '.')]
		}
		var clean strings.Builder
		for _, r := range part {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				_, _ = clean.WriteRune(r)
			}
		}
		if w := clean.String(); w != "" {
			words = append(words, strings.ToLower(w))
		}
	}
	return words
}

// h1TitleMatchesName returns true when nameWords aligns with titleWords using greedy
// bidirectional compound matching: up to 3 adjacent name words may collapse into 1 title
// token (e.g. ["vo2","max"] → "vo2max"), and up to 3 adjacent title tokens may collapse
// into 1 name word (e.g. ["oauth","2"] → "oauth2" for "OAuth 2.0" matching "oauth2").
func h1TitleMatchesName(nameWords, titleWords []string) bool {
	ni, ti := 0, 0
	for ni < len(nameWords) && ti < len(titleWords) {
		matched := false
	outer:
		for nMerge := 1; nMerge <= 3 && ni+nMerge <= len(nameWords); nMerge++ {
			nc := strings.Join(nameWords[ni:ni+nMerge], "")
			for tMerge := 1; tMerge <= 3 && ti+tMerge <= len(titleWords); tMerge++ {
				if nc == strings.Join(titleWords[ti:ti+tMerge], "") {
					ni += nMerge
					ti += tMerge
					matched = true
					break outer
				}
			}
		}
		if !matched {
			return false
		}
	}
	return ni == len(nameWords) && ti == len(titleWords)
}

// skillDomain returns the top-level domain segment from an absolute skill path.
// It looks for the segment immediately after the first 'skills' path component.
func skillDomain(absPath string) string {
	parts := strings.Split(filepath.ToSlash(absPath), "/")
	for i, p := range parts {
		if p == "skills" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// inferCorpusRoot finds the outermost 'skills' directory in an absolute skill
// path — that is the package-level skills root used as the duplicate corpus.
// Returns "" when no 'skills' component is found.
func inferCorpusRoot(skillMDPath string) string {
	parts := strings.Split(filepath.ToSlash(skillMDPath), "/")
	for i, p := range parts {
		if p == "skills" {
			return filepath.FromSlash(strings.Join(parts[:i+1], "/"))
		}
	}
	return ""
}

// ── Near-duplicate detection helpers ─────────────────────────────────────────

var dupStopwords = map[string]bool{
	"use": true, "when": true, "a": true, "an": true, "the": true,
	"is": true, "are": true, "for": true, "to": true, "in": true,
	"of": true, "and": true, "or": true, "with": true, "that": true,
	"this": true, "as": true, "by": true, "on": true, "at": true,
	"from": true, "be": true, "been": true, "being": true, "have": true,
	"has": true, "had": true, "do": true, "does": true, "did": true,
	"will": true, "would": true, "could": true, "should": true, "may": true,
	"might": true, "it": true, "its": true, "not": true, "but": true,
	"if": true, "you": true, "your": true,
}

var dupNonWordRe = regexp.MustCompile(`\W+`)

// dupTokenize lowercases, splits on non-word chars, and drops tokens ≤2 chars or in dupStopwords.
func dupTokenize(s string) []string {
	raw := dupNonWordRe.Split(strings.ToLower(s), -1)
	out := raw[:0]
	for _, w := range raw {
		if len(w) > 2 && !dupStopwords[w] {
			out = append(out, w)
		}
	}
	return out
}

// dupJaccard computes intersection/union for two string slices (case-insensitive, set-based).
func dupJaccard(a, b []string) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	setA := make(map[string]bool, len(a))
	for _, v := range a {
		setA[strings.ToLower(v)] = true
	}
	setB := make(map[string]bool, len(b))
	for _, v := range b {
		setB[strings.ToLower(v)] = true
	}
	union := make(map[string]bool, len(setA)+len(setB))
	for k := range setA {
		union[k] = true
	}
	for k := range setB {
		union[k] = true
	}
	intersection := 0
	for k := range setA {
		if setB[k] {
			intersection++
		}
	}
	return float64(intersection) / float64(len(union))
}

// dupScore computes the weighted similarity between two skills.
// Formula (matches check-duplicates.sh): (tagJaccard×2 + descTokenJaccard×3 + sameDomain×0.5) / 5.5
func dupScore(a, b skills.Skill) float64 {
	t := dupJaccard(a.Tags, b.Tags)
	d := dupJaccard(dupTokenize(a.Description), dupTokenize(b.Description))
	dom := 0.0
	if skillDomain(a.Path) == skillDomain(b.Path) {
		dom = 0.5
	}
	return (t*2 + d*3 + dom) / 5.5
}

// jaccardSets computes intersection/union on precomputed bool sets — no allocation.
func jaccardSets(a, b map[string]bool) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	intersection := 0
	for k := range a {
		if b[k] {
			intersection++
		}
	}
	union := len(a) + len(b) - intersection
	return float64(intersection) / float64(union)
}

// dupScoreSets is the hot-path scorer that operates on precomputed sets.
func dupScoreSets(aTagSet, aDescSet map[string]bool, aDomain string, bTagSet, bDescSet map[string]bool, bDomain string) float64 {
	t := jaccardSets(aTagSet, bTagSet)
	d := jaccardSets(aDescSet, bDescSet)
	dom := 0.0
	if aDomain == bDomain {
		dom = 0.5
	}
	return (t*2 + d*3 + dom) / 5.5
}

// extractH1 returns the first h1 line from the body and whether one was found.
func extractH1(body string) (string, bool) {
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "# ") {
			return line, true
		}
	}
	return "", false
}

// hasPurposeStatement returns true when a non-heading, non-empty line appears
// within 5 non-empty lines after the first h1 heading.
func hasPurposeStatement(body string) bool {
	lines := strings.Split(body, "\n")
	inH1 := false
	checked := 0
	for _, l := range lines {
		if strings.HasPrefix(l, "# ") {
			inH1 = true
			continue
		}
		if !inH1 {
			continue
		}
		if strings.TrimSpace(l) == "" {
			continue
		}
		checked++
		if !strings.HasPrefix(l, "#") {
			return true
		}
		if checked > 5 {
			break
		}
	}
	return false
}

// fixableFindings returns only the findings that an AI can plausibly correct.
// Excludes duplicate-name errors (require human file rename) and size violations
// (require human editorial judgment to truncate or expand content).
func fixableFindings(findings []validateFinding) []validateFinding {
	var out []validateFinding
	for _, f := range findings {
		if strings.Contains(f.Message, "' duplicated at ") {
			continue
		}
		if strings.HasPrefix(f.Message, "size: ") {
			continue
		}
		out = append(out, f)
	}
	return out
}

// deterministicFix applies rule-based transforms that need no AI judgment.
// Returns the fixed content and whether any change was made.
// Transforms operate only on the YAML frontmatter; body is untouched.
func deterministicFix(content string, findings []validateFinding) (string, bool) {
	if !strings.HasPrefix(content, "---") {
		return content, false
	}
	rest := content[3:]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return content, false
	}
	frontmatter := rest[:end]
	suffix := rest[end:]

	var fixName, fixDesc, fixTags bool
	for _, f := range findings {
		switch {
		case strings.Contains(f.Message, "name: must be lowercase kebab-case"):
			fixName = true
		case strings.Contains(f.Message, "frontmatter: invalid YAML") &&
			strings.Contains(f.Message, "mapping values are not allowed"):
			fixDesc = true
		case strings.HasPrefix(f.Message, "tags: '") &&
			strings.Contains(f.Message, "must be lowercase"):
			fixTags = true
		}
	}

	changed := false
	if fixName {
		if fm, ok := deterministicFixName(frontmatter); ok {
			frontmatter, changed = fm, true
		}
	}
	if fixDesc {
		if fm, ok := deterministicFixYAMLQuoting(frontmatter); ok {
			frontmatter, changed = fm, true
		}
	}
	if fixTags {
		if fm, ok := deterministicFixTags(frontmatter); ok {
			frontmatter, changed = fm, true
		}
	}
	if !changed {
		return content, false
	}
	return "---" + frontmatter + suffix, true
}

func deterministicFixName(frontmatter string) (string, bool) {
	orig := frontmatter
	result := deterministicNameRe.ReplaceAllStringFunc(frontmatter, func(match string) string {
		parts := deterministicNameRe.FindStringSubmatch(match)
		return parts[1] + strings.ToLower(parts[2]) + parts[3]
	})
	return result, result != orig
}

func deterministicFixYAMLQuoting(frontmatter string) (string, bool) {
	orig := frontmatter
	result := deterministicYAMLQuotingRe.ReplaceAllStringFunc(frontmatter, func(match string) string {
		parts := deterministicYAMLQuotingRe.FindStringSubmatch(match)
		val := parts[3]
		if strings.HasPrefix(val, `"`) || strings.HasPrefix(val, `'`) || !strings.Contains(val, ":") {
			return match // already quoted, or no colon — nothing to do
		}
		escaped := strings.ReplaceAll(val, `"`, `\"`)
		return parts[1] + `"` + escaped + `"`
	})
	return result, result != orig
}

func deterministicFixTags(frontmatter string) (string, bool) {
	orig := frontmatter
	result := deterministicTagsRe.ReplaceAllStringFunc(frontmatter, func(match string) string {
		parts := deterministicTagsRe.FindStringSubmatch(match)
		items := strings.Split(parts[2], ",")
		lowered := make([]string, len(items))
		for i, item := range items {
			lowered[i] = strings.ToLower(strings.TrimSpace(item))
		}
		return parts[1] + strings.Join(lowered, ", ") + parts[3]
	})
	return result, result != orig
}

// fixSkillWithAI calls the given AI executor to repair a failing SKILL.md.
// It reads the file, builds a targeted prompt from the findings, dispatches to the AI,
// extracts the corrected SKILL.md from the response, and writes it back.
// Callers must resolve the executor before calling and must not pass execPrint.
func fixSkillWithAI(goCtx context.Context, result validateResult, projectDir string, ex *executorSpec) error {
	content, err := os.ReadFile(result.Path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	var fb strings.Builder
	for _, f := range result.Findings {
		fmt.Fprintf(&fb, "- [%s] %s\n", strings.ToUpper(f.Level), f.Message)
	}

	const system = "You are grimoire, a SKILL.md editor. " +
		"Output ONLY the corrected SKILL.md content. " +
		"Start your output with exactly '---' (the YAML frontmatter delimiter). " +
		"Do not add any prose, explanation, or markdown code fences around the output."
	userMsg := "Fix the following SKILL.md file. Apply ONLY the fixes required for the findings listed below.\n" +
		"Do not alter any content that is not mentioned in the findings.\n\n" +
		"Findings to fix:\n" + fb.String() +
		"\nCurrent SKILL.md content:\n" + string(content)

	var raw string
	switch ex.Kind {
	case execLocalCLI:
		raw, err = callViaLocalAgent(goCtx, ex.Agent, projectDir, system+"\n\n"+userMsg)
	case execAPI:
		if ex.Provider.Format == "anthropic" {
			raw, err = callAnthropicAPI(ex.Provider.APIKey, ex.Provider.Model, ex.Provider.MaxTokens, system, userMsg)
		} else {
			raw, err = callOpenAICompatible(ex.Provider, system, userMsg)
		}
	}
	if err != nil {
		return fmt.Errorf("AI call failed: %w", err)
	}

	fixed, ok := extractSkillMDContent(raw)
	if !ok {
		return fmt.Errorf("AI response did not contain valid SKILL.md content (no '---' delimiter found)")
	}
	if err := os.WriteFile(result.Path, []byte(fixed), 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

// extractSkillMDContent strips leading markdown fences from an AI response and
// returns the SKILL.md content starting from the first '---' delimiter.
func extractSkillMDContent(raw string) (string, bool) {
	// Strip leading markdown code fences (```markdown, ```, etc.)
	s := strings.TrimSpace(raw)
	if strings.HasPrefix(s, "```") {
		end := strings.Index(s[3:], "\n")
		if end >= 0 {
			s = s[3+end+1:]
		}
		if idx := strings.LastIndex(s, "```"); idx >= 0 {
			s = strings.TrimSpace(s[:idx])
		}
	}
	idx := strings.Index(s, "---")
	if idx < 0 {
		return "", false
	}
	return s[idx:], true
}

// runTestSchema validates each fixture in <dir>/valid/ (must pass) and <dir>/invalid/
// (must fail) individually, reporting conformance per fixture.
func runTestSchema(cmd *cobra.Command, dir string) error {
	type fixtureResult struct {
		name string
		want string // "pass" or "fail"
		got  string // "pass" or "fail"
	}
	var results []fixtureResult

	runFixtures := func(subdir, want string) error {
		entries, err := os.ReadDir(filepath.Join(dir, subdir))
		if err != nil {
			return fmt.Errorf("reading %s: %w", filepath.Join(dir, subdir), err)
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			skillMD := filepath.Join(dir, subdir, e.Name(), "SKILL.md")
			if _, err := os.Stat(skillMD); err != nil {
				continue
			}
			r := validateSkill(skillMD)
			got := "pass"
			if r.Status == "fail" {
				got = "fail"
			}
			results = append(results, fixtureResult{subdir + "/" + e.Name(), want, got})
		}
		return nil
	}

	if err := runFixtures("valid", "pass"); err != nil {
		return err
	}
	if err := runFixtures("invalid", "fail"); err != nil {
		return err
	}

	passed, failed := 0, 0
	for _, r := range results {
		if r.want == r.got {
			fmt.Printf("%s  %s\n", tui.IconOK, r.name)
			passed++
		} else {
			fmt.Printf("%s  %s — expected %s, got %s\n", tui.IconFail, r.name, r.want, r.got)
			failed++
		}
	}
	fmt.Printf("\nConformance: %d/%d fixtures correct\n", passed, passed+failed)
	if failed > 0 {
		cmd.Root().SilenceErrors = true
		return fmt.Errorf("conformance test failed: %d fixture(s) incorrect", failed)
	}
	return nil
}
