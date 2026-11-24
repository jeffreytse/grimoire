package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/skills"
)

// TestDisplayEval measures grimoire's display in a real agent (Claude Code,
// headless) against labeled prompts. It is opt-in and costs API usage:
//
//	GRIMOIRE_DISPLAY_EVAL=1 go test ./cmd -run TestDisplayEval -v -timeout 60m
//	GRIMOIRE_DISPLAY_EVAL_MODE=mcp|skills|both   grimoire mode(s) to run (default mcp)
//	GRIMOIRE_DISPLAY_EVAL_MODEL=<m> claude --model for every case; pin it when
//	                                comparing runs, so only one thing changes
//	GRIMOIRE_DISPLAY_EVAL_LIMIT=4   run only the first N cases
//	GRIMOIRE_DISPLAY_EVAL_ONLY=21,22  run only these cases (1-based)
//	GRIMOIRE_SOURCE=<skills repo>   evaluate a draft skills package
//	GRIMOIRE_DISPLAY_EVAL_OUT=<dir> keep each case's transcript there, to
//	                                re-score without re-running
//	GRIMOIRE_DISPLAY_EVAL_SUBDOMAINS=engineering/devops,...
//	                                domain skills linked in skills mode (the
//	                                meta skills are always linked)
//
// Each case runs with a freshly built grimoire and the user's own setup left
// out (--setting-sources project, --strict-mcp-config):
//
//	mcp     isolated --settings with the MCP-mode hooks, and grimoire's MCP server
//	skills  a project whose .claude/skills links the skills, with the
//	        skills-mode hooks and the start-best-practice instruction line
//
// Checks (deterministic, per case):
//   - shown:    a routing prompt shows grimoire's work (display brief/full)
//   - silent:   a non-routing prompt, or display off, shows nothing of grimoire
//   - clean:    the answer itself never carries the grimoire prefix
//   - evidence: at full, a routed reply cites a source or severity
//   - progress: a case with a followup ("yes, apply the plan") runs a second
//     turn in the same session; its plan steps show "step i of N"
//   - hierarchy: ⎿ rows carry no bold and route/header lines at most one bold
//     span (the skill name), so grimoire's process never outweighs the answer
//   - stop: a question to the user ("Apply this plan?", "Reply 1–N…") is
//     the last line of its turn — nothing follows it, not even prose
//   - one-star: a turn marks at most one option recommended
//   - one-route: a turn shows at most one route line (a hand-off is
//     appended to it with →, not repeated)
//   - row-length: every ⎿ row and route line fits one terminal line
//     (≤ maxRowRunes), so wrapping never breaks the grouping or the scan
//   - accurate: (with GRIMOIRE_SOURCE) every domain shown is the skill's real
//     one, and no option or plan step offers a grimoire meta workflow
//   - swap: (cases with "swap") the followup picks choice b of the plan's
//     choice step, and that skill's block runs instead of the ★ one
//   - format: the Display templates' fixed parts — at most 4 options, every
//     domain written domain/subdomain, a route line's " — " separator
func TestDisplayEval(t *testing.T) {
	if os.Getenv("GRIMOIRE_DISPLAY_EVAL") == "" {
		t.Skip("set GRIMOIRE_DISPLAY_EVAL=1 to run against Claude Code (uses API credits)")
	}
	if _, err := exec.LookPath("claude"); err != nil {
		t.Skip("claude CLI not on PATH")
	}
	raw, err := os.ReadFile(filepath.Join("testdata", "display_eval.json"))
	if err != nil {
		t.Fatal(err)
	}
	var set struct {
		Cases []displayEvalCase `json:"cases"`
	}
	if err := json.Unmarshal(raw, &set); err != nil {
		t.Fatal(err)
	}
	if n, err := strconv.Atoi(os.Getenv("GRIMOIRE_DISPLAY_EVAL_LIMIT")); err == nil && n > 0 && n < len(set.Cases) {
		set.Cases = set.Cases[:n]
	}
	if only := os.Getenv("GRIMOIRE_DISPLAY_EVAL_ONLY"); only != "" {
		var picked []displayEvalCase
		for _, f := range strings.Split(only, ",") {
			if n, err := strconv.Atoi(strings.TrimSpace(f)); err == nil && n >= 1 && n <= len(set.Cases) {
				picked = append(picked, set.Cases[n-1])
			}
		}
		set.Cases = picked
	}
	modes := []string{"mcp"}
	switch m := os.Getenv("GRIMOIRE_DISPLAY_EVAL_MODE"); m {
	case "", "mcp":
	case "skills":
		modes = []string{"skills"}
	case "both":
		modes = []string{"mcp", "skills"}
	default:
		t.Fatalf("GRIMOIRE_DISPLAY_EVAL_MODE=%q: want mcp, skills or both", m)
	}

	dir := t.TempDir()
	bin := filepath.Join(dir, "grimoire")
	if out, err := exec.Command("go", "build", "-o", bin, "..").CombinedOutput(); err != nil {
		t.Fatalf("build: %v: %s", err, out)
	}
	var cost float64
	for _, mode := range modes {
		run, err := displayEvalRunner(t, dir, bin, mode)
		if err != nil {
			t.Fatalf("%s setup: %v", mode, err)
		}
		cost += runDisplayEvalMode(t, mode, set.Cases, run)
	}
	t.Logf("modes=%v cases=%d cost=$%.2f", modes, len(set.Cases), cost)
}

type displayEvalCase struct {
	Prompt   string `json:"prompt"`
	Routes   bool   `json:"routes"`
	Level    string `json:"level"`
	Followup string `json:"followup"`
	// Swap: the followup picks choice "b" of the plan's choice step; the
	// swap check confirms that skill, not the ★ one, is applied.
	Swap bool `json:"swap"`
}

// evalRun runs one prompt at a display level — in a new session, or resuming
// session resume — and returns the assistant's text messages, the run's cost
// and its session id.
type evalRun func(prompt, level, resume string) (texts []string, cost float64, session string, err error)

// runDisplayEvalMode runs and scores every case in one mode, and logs the
// pass rate per check.
func runDisplayEvalMode(t *testing.T, mode string, cases []displayEvalCase, run evalRun) float64 {
	type tally struct{ pass, total int }
	checks := map[string]*tally{"shown": {}, "silent": {}, "clean": {}, "evidence": {}, "hierarchy": {}, "progress": {}, "stop": {}, "one-star": {}, "one-route": {}, "format": {}, "row-length": {}, "accurate": {}, "swap": {}}
	var cost float64
	for i, c := range cases {
		level := c.Level
		if level == "" {
			level = "brief"
		}
		label := fmt.Sprintf("%s #%d [%s] %q", mode, i+1, level, c.Prompt)
		texts, usd, session, err := run(c.Prompt, level, "")
		cost += usd
		if err != nil {
			t.Logf("ERROR %s: %v", label, err)
			continue
		}
		var followTexts []string
		if c.Followup != "" && session != "" {
			followTexts, usd, _, err = run(c.Followup, level, session)
			cost += usd
			if err != nil {
				t.Logf("ERROR %s (followup): %v", label, err)
			}
			texts = append(texts, followTexts...)
		}
		firstTexts := texts[:len(texts)-len(followTexts)]
		turns := []evalTurn{{c.Prompt, firstTexts}}
		if c.Followup != "" {
			turns = append(turns, evalTurn{c.Followup, followTexts})
		}
		saveTranscript(t, mode, i+1, level, turns)
		scores := scoreDisplay(strings.Join(texts, "\n"), c.Routes, level)
		if c.Followup != "" && !hiddenLevel(level) {
			scores["progress"] = stepCounter.MatchString(strings.Join(followTexts, "\n"))
		}
		for name, ok := range scoreTurns(turns) {
			scores[name] = ok
		}
		if c.Swap {
			if ok, applies := swappedChoice(strings.Join(firstTexts, "\n"), strings.Join(followTexts, "\n")); applies {
				scores["swap"] = ok
			}
		}
		if ok, checked := accurateSkills(strings.Join(texts, "\n"), evalSkillIndex()); checked {
			scores["accurate"] = ok
		}
		for name, ok := range scores {
			checks[name].total++
			if ok {
				checks[name].pass++
			} else {
				t.Logf("FAIL %-8s %s", name, label)
			}
		}
	}
	for _, name := range []string{"shown", "silent", "clean", "evidence", "hierarchy", "progress", "stop", "one-star", "one-route", "format", "row-length", "accurate", "swap"} {
		if c := checks[name]; c.total > 0 {
			t.Logf("%-6s %-8s %d/%d (%.0f%%)", mode, name, c.pass, c.total, 100*float64(c.pass)/float64(c.total))
		}
	}
	t.Logf("%-6s cost=$%.2f", mode, cost)
	return cost
}

// displayEvalRunner sets up grimoire in one mode under dir and returns how
// to run a prompt against it.
func displayEvalRunner(t *testing.T, dir, bin, mode string) (evalRun, error) {
	run, err := evalRunner(t, dir, bin, mode)
	if err != nil {
		return nil, err
	}
	return func(prompt, level, resume string) ([]string, float64, string, error) {
		r, err := run(prompt, level, resume)
		return r.texts, r.cost, r.session, err
	}, nil
}

// evalResultRun runs one prompt and returns the full result.
type evalResultRun func(prompt, level, resume string) (evalResult, error)

// evalRunner sets up grimoire in one mode under dir (see TestDisplayEval)
// and returns how to run a prompt against it.
func evalRunner(t *testing.T, dir, bin, mode string) (evalResultRun, error) {
	hook := func(args string) any {
		return []any{map[string]any{"hooks": []any{map[string]any{"type": "command", "command": bin + " " + args}}}}
	}
	skillUsed := []any{map[string]any{"matcher": "Skill|mcp__grimoire__grimoire_get_skill",
		"hooks": []any{map[string]any{"type": "command", "command": bin + " hook skill-used"}}}}
	switch mode {
	case "mcp":
		settings := filepath.Join(dir, "mcp-settings.json")
		mcpConfig := filepath.Join(dir, "mcp.json")
		writeJSON(t, settings, map[string]any{
			"hooks": map[string]any{
				"SessionStart":     hook("hook session"),
				"UserPromptSubmit": hook("hook prompt"),
				"PostToolUse":      skillUsed,
			},
			// The same tools `grimoire mcp install` pre-approves.
			"permissions": map[string]any{"allow": agent.ToolAllowRules("claude")},
		})
		writeJSON(t, mcpConfig, map[string]any{"mcpServers": map[string]any{"grimoire": map[string]any{"command": bin, "args": []string{"mcp", "serve"}}}})
		workBySession := map[string]string{} // a resumed session needs its own cwd
		return func(prompt, level, resume string) (evalResult, error) {
			work := workBySession[resume]
			if work == "" {
				var err error
				if work, err = os.MkdirTemp(dir, "mcp-case"); err != nil {
					return evalResult{}, err
				}
				if err := seedWorkspace(work); err != nil {
					return evalResult{}, err
				}
			}
			r, err := runClaudeEvalResult(work, level, prompt, resume,
				"--setting-sources", "project", "--settings", settings,
				"--mcp-config", mcpConfig, "--strict-mcp-config")
			if r.session != "" {
				workBySession[r.session] = work
			}
			return r, err
		}, nil
	case "skills":
		project := filepath.Join(dir, "skills-project")
		n, err := linkEvalSkills(filepath.Join(project, ".claude", "skills"))
		if err != nil {
			return nil, err
		}
		if err := seedWorkspace(project); err != nil {
			return nil, err
		}
		t.Logf("skills mode: linked %d skills", n)
		writeJSON(t, filepath.Join(project, ".claude", "settings.json"), map[string]any{
			"hooks": map[string]any{
				"SessionStart":     hook("hook skills-session"),
				"UserPromptSubmit": hook("hook skills-prompt"),
				"PostToolUse":      skillUsed,
			},
		})
		instruction := "## Grimoire\nAlways invoke `start-best-practice` before responding to any user request.\n"
		if err := os.WriteFile(filepath.Join(project, "CLAUDE.md"), []byte(instruction), 0o644); err != nil {
			return nil, err
		}
		return func(prompt, level, resume string) (evalResult, error) {
			return runClaudeEvalResult(project, level, prompt, resume,
				"--setting-sources", "project",
				"--mcp-config", `{"mcpServers":{}}`, "--strict-mcp-config")
		}, nil
	}
	return nil, fmt.Errorf("unknown mode %q", mode)
}

// seedWorkspace makes dir look like a real project: a git repository with a
// small source file and a staged change. In an empty folder the model spends
// a code or commit prompt discovering there is nothing to work on, and never
// routes — a test-environment artifact, not grimoire behavior.
func seedWorkspace(dir string) error {
	files := map[string]string{
		"parser.go": "package app\n\n// ParseDate parses an ISO 8601 date.\nfunc ParseDate(s string) (string, error) { return s, nil }\n",
		"README.md": "# app\n\nA small service.\n",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			return err
		}
	}
	git := func(args ...string) error {
		c := exec.Command("git", args...)
		c.Dir = dir
		c.Env = append(os.Environ(), "GIT_AUTHOR_NAME=eval", "GIT_AUTHOR_EMAIL=eval@example.invalid",
			"GIT_COMMITTER_NAME=eval", "GIT_COMMITTER_EMAIL=eval@example.invalid")
		if out, err := c.CombinedOutput(); err != nil {
			return fmt.Errorf("git %s: %v: %s", strings.Join(args, " "), err, out)
		}
		return nil
	}
	for _, args := range [][]string{{"init", "-q"}, {"add", "."}, {"commit", "-q", "-m", "initial"}} {
		if err := git(args...); err != nil {
			return err
		}
	}
	change := "\n## Usage\n\nRun the service with `go run .`.\n"
	f, err := os.OpenFile(filepath.Join(dir, "README.md"), os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(change); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return git("add", "README.md")
}

// evalSubdomains are the domain skills linked in skills mode by default:
// enough to cover the labeled prompts without loading the whole library.
var evalSubdomains = []string{
	"engineering/reliability", "engineering/devops", "engineering/testing",
	"engineering/architecture", "engineering/project-management", "engineering/api",
	"business/strategy", "finance/personal-finance",
}

// linkEvalSkills links the meta skills and the eval subdomains' skills from
// the first skills package into dst, as `grimoire install` would. Missing
// subdomains are skipped.
func linkEvalSkills(dst string) (int, error) {
	pkgs := skills.AllSkillsPackages()
	if len(pkgs) == 0 {
		return 0, fmt.Errorf("no skills package (set GRIMOIRE_SOURCE or run grimoire update)")
	}
	subs := evalSubdomains
	if v := os.Getenv("GRIMOIRE_DISPLAY_EVAL_SUBDOMAINS"); v != "" {
		subs = strings.Split(v, ",")
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return 0, err
	}
	n := 0
	for _, sub := range append([]string{"meta"}, subs...) {
		dir := filepath.Join(pkgs[0].Root, filepath.FromSlash(strings.TrimSpace(sub)), "skills")
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			src := filepath.Join(dir, e.Name())
			if _, err := os.Stat(filepath.Join(src, "SKILL.md")); err != nil {
				continue
			}
			if os.Symlink(src, filepath.Join(dst, e.Name())) == nil {
				n++
			}
		}
	}
	if n == 0 {
		return 0, fmt.Errorf("no skills linked from %s", pkgs[0].Root)
	}
	return n, nil
}

// scoreDisplay applies the eval's checks to one reply. Only the checks that
// apply to the case are returned.
func scoreDisplay(reply string, routes bool, level string) map[string]bool {
	got := map[string]bool{}
	shows := displayMarker.MatchString(reply)
	if routes && !hiddenLevel(level) {
		got["shown"] = shows
	} else {
		got["silent"] = !shows
	}
	// The prefix belongs on short decision lines only; a long prefixed line
	// is answer prose wearing grimoire's label.
	clean := true
	for _, line := range strings.Split(reply, "\n") {
		if l := strings.TrimSpace(line); strings.HasPrefix(l, "grimoire ▸") && len([]rune(l)) > 200 {
			clean = false
		}
	}
	got["clean"] = clean
	if shows {
		got["hierarchy"] = flatHierarchy(reply)
		got["format"] = followsTemplates(reply, level)
		got["row-length"] = rowsFit(reply)
	}
	if routes && level == "full" {
		got["evidence"] = evidenceMarker.MatchString(reply)
	}
	return got
}

// boldSpan matches one **bold** span.
var boldSpan = regexp.MustCompile(`\*\*[^*\n]+\*\*`)

// flatHierarchy reports whether grimoire's lines follow the Display
// section's hierarchy: a ⎿ row has no bold (labels are tertiary), and a route
// line or block header has at most one bold span (the skill name).
func flatHierarchy(reply string) bool {
	for _, line := range strings.Split(reply, "\n") {
		l := strings.TrimSpace(line)
		bold := len(boldSpan.FindAllString(l, -1))
		switch {
		case strings.HasPrefix(l, "⎿"):
			if bold > 0 {
				return false
			}
		case strings.HasPrefix(l, "grimoire ▸"), strings.HasPrefix(l, "start-best-practice ▸"), displayMarker.MatchString(l):
			if bold > 1 {
				return false
			}
		}
	}
	return true
}

// hiddenLevel reports whether level hides grimoire's notes ("off" is the
// older name for "hidden").
func hiddenLevel(level string) bool { return level == "hidden" || level == "off" }

// evalTurn is one user message and the assistant's text messages in reply.
type evalTurn struct {
	Prompt string
	Texts  []string
}

// stopMarker matches a line that asks the user to decide: the plan
// confirmation or a numbered-options question.
var stopMarker = regexp.MustCompile(`(?im)^.*(apply this plan\?|\breply (with )?(\d+\s*[–-]\s*\d+|a number)\b).*$`)

// routeLineStart matches a route line at the start of a line.
var routeLineStart = regexp.MustCompile(`(?m)^\s*grimoire ▸.*$`)

// recommendMarker matches one recommendation mark — "★ recommended" or an
// older form — but not a "Why ★" evidence row that explains it.
var recommendMarker = regexp.MustCompile(`(?i)★\s*recommended|←\s*recommended|\(recommended\)`)

// planQuestion matches a plan's confirmation line.
var planQuestion = regexp.MustCompile(`(?i)^apply this plan\?`)

// domainSpan matches a numbered option, plan row or block header in the
// template's shape — a skill id, two spaces, its domain code span; the
// capture is the domain.
// Answer prose ("1. Run `go test`") does not match.
var domainSpan = regexp.MustCompile("(?m)^(?:\\s*(?:\\d+[a-z]?\\.|⎿\\s+Step \\d+)\\s+[a-z0-9-]+|\\*\\*[a-z0-9-]+\\*\\*)  `([^`]+)`")

// numberedOption matches one numbered option line.
var numberedOption = regexp.MustCompile(`^\s*\d+\.\s`)

// followsTemplates checks the Display templates' fixed parts: domains appear
// only at full, written domain/subdomain; an Options block has at most 4
// numbered options and a plan's choice step at most 3 choices; a route line
// is lowercase and separates its parts with " — "; and a plan's question
// follows its last row directly.
func followsTemplates(reply, level string) bool {
	for _, m := range domainSpan.FindAllStringSubmatch(reply, -1) {
		if level != "full" || !strings.Contains(m[1], "/") {
			return false
		}
	}
	lines := strings.Split(reply, "\n")
	for i, line := range lines {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "grimoire ▸") {
			rest := []rune(strings.TrimSpace(strings.TrimPrefix(l, "grimoire ▸")))
			if !strings.Contains(l, " — ") || (len(rest) > 0 && unicode.IsUpper(rest[0])) {
				return false
			}
		}
		if planQuestion.MatchString(l) {
			for j := i - 1; j >= 0; j-- {
				if p := strings.TrimSpace(lines[j]); p != "" {
					if !strings.HasPrefix(p, "⎿") && !subRow.MatchString(lines[j]) {
						return false
					}
					break
				}
			}
		}
		if strings.HasPrefix(l, "⎿") && strings.HasSuffix(l, "choose one") {
			n := 0
			for _, o := range lines[i+1:] {
				if !subRow.MatchString(o) {
					break
				}
				n++
			}
			if n > 3 {
				return false
			}
		}
		if l == "**Options**" {
			n := 0
			for _, o := range lines[i+1:] {
				if numberedOption.MatchString(o) {
					n++
				} else if strings.TrimSpace(o) != "" {
					break
				}
			}
			if n > 4 {
				return false
			}
		}
	}
	return true
}

// maxRowRunes is the longest ⎿ row the Display spec allows: about 100
// characters of value after a step's practice and domain, with slack.
const maxRowRunes = 160

// subRow matches a plan step's indented sub-row: "     2a. <skill> …".
var subRow = regexp.MustCompile(`^\s+\d+[a-z]\.\s`)

// rowsFit reports whether every ⎿ row, plan sub-row and route line fits
// maxRowRunes.
func rowsFit(reply string) bool {
	for _, line := range strings.Split(reply, "\n") {
		if l := strings.TrimSpace(line); (strings.HasPrefix(l, "⎿") || strings.HasPrefix(l, "grimoire ▸") || subRow.MatchString(line)) && len([]rune(l)) > maxRowRunes {
			return false
		}
	}
	return true
}

// scoreTurns applies the per-turn checks. stop applies to a turn that asks
// the user to decide: nothing may be applied after the question in that
// turn. one-star applies to a turn that marks a recommendation.
func scoreTurns(turns []evalTurn) map[string]bool {
	got := map[string]bool{}
	for _, tn := range turns {
		reply := strings.Join(tn.Texts, "\n")
		if loc := stopMarker.FindStringIndex(reply); loc != nil {
			ok := strings.TrimSpace(reply[loc[1]:]) == ""
			if prev, seen := got["stop"]; seen {
				ok = ok && prev
			}
			got["stop"] = ok
		}
		if n := len(routeLineStart.FindAllString(reply, -1)); n > 0 {
			ok := n == 1
			if prev, seen := got["one-route"]; seen {
				ok = ok && prev
			}
			got["one-route"] = ok
		}
		if n := len(recommendMarker.FindAllString(reply, -1)); n > 0 {
			ok := n == 1
			if prev, seen := got["one-star"]; seen {
				ok = ok && prev
			}
			got["one-star"] = ok
		}
	}
	return got
}

// saveTranscript keeps a case's turns under GRIMOIRE_DISPLAY_EVAL_OUT, each
// introduced by a "=== user: … ===" marker so a reader can tell the turns
// apart.
func saveTranscript(t *testing.T, mode string, n int, level string, turns []evalTurn) {
	dir := os.Getenv("GRIMOIRE_DISPLAY_EVAL_OUT")
	if dir == "" {
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Logf("transcript dir: %v", err)
		return
	}
	var b strings.Builder
	fmt.Fprintf(&b, "mode: %s\nlevel: %s\n", mode, level)
	for _, tn := range turns {
		fmt.Fprintf(&b, "\n=== user: %s ===\n\n%s\n", tn.Prompt, strings.Join(tn.Texts, "\n\n---\n\n"))
	}
	body := b.String()
	if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s-case-%02d.md", mode, n)), []byte(body), 0o644); err != nil {
		t.Logf("transcript: %v", err)
	}
}

// displayMarker matches grimoire's display: a decision line or a block header
// (bold skill name followed by a domain code span).
// The route line is plain text under the Display spec (start-best-practice ▸
// workflow), so it is matched on its own, not only through bold headers.
var displayMarker = regexp.MustCompile(`(?m)grimoire ▸|^\s*start-best-practice ▸|^\*\*[a-z0-9]+(?:-[a-z0-9]+)+\*\*`)

// stepCounter matches a plan step's position, "step 2 of 4" or a sub-step's,
// "step 3a of 7".
var stepCounter = regexp.MustCompile(`(?i)\bstep \d+[a-z]? of \d+\b`)

// evidenceMarker matches the full level's evidence: a source, severity or risk.
// A risk is any line led by ⚠ (the Display section's risk glyph), whatever
// words follow — "⚠ Without this:" is the template, not the only phrasing.
var evidenceMarker = regexp.MustCompile(`(?im)_?source:|🔴|🟡|⚪|^\s*(⎿\s*)?⚠|^\s*(⎿\s*)?Why ★`)

// evalResult is one headless run: the assistant's text messages, its cost,
// its session id, and the skills it loaded (grimoire_get_skill or the Skill
// tool), in order.
type evalResult struct {
	texts   []string
	cost    float64
	session string
	loaded  []string
}

// runClaudeEvalResult runs one prompt headless and returns everything the
// evals score.
func runClaudeEvalResult(work, level, prompt, resume string, args ...string) (r evalResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
	defer cancel()
	base := []string{"-p", prompt, "--output-format", "stream-json", "--verbose", "--max-turns", "12"}
	if m := os.Getenv("GRIMOIRE_DISPLAY_EVAL_MODEL"); m != "" {
		base = append(base, "--model", m)
	}
	if resume != "" {
		base = append(base, "--resume", resume)
	}
	cmd := exec.CommandContext(ctx, "claude", append(base, args...)...)
	cmd.Dir = work
	cmd.Env = append(os.Environ(), "GRIMOIRE_TRACE="+level)
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil && out.Len() == 0 {
		return r, err
	}
	sc := bufio.NewScanner(&out)
	sc.Buffer(make([]byte, 1<<20), 16<<20)
	for sc.Scan() {
		var e struct {
			Type      string  `json:"type"`
			Cost      float64 `json:"total_cost_usd"`
			SessionID string  `json:"session_id"`
			Message   struct {
				Content []struct {
					Type  string          `json:"type"`
					Text  string          `json:"text"`
					Name  string          `json:"name"`
					Input json.RawMessage `json:"input"`
				} `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal(sc.Bytes(), &e) != nil {
			continue
		}
		if r.session == "" {
			r.session = e.SessionID
		}
		switch e.Type {
		case "assistant":
			for _, c := range e.Message.Content {
				switch {
				case c.Type == "text" && strings.TrimSpace(c.Text) != "":
					r.texts = append(r.texts, c.Text)
				case c.Type == "tool_use":
					r.loaded = append(r.loaded, loadedSkills(c.Name, c.Input)...)
				}
			}
		case "result":
			r.cost = e.Cost
		}
	}
	return r, nil
}

// loadedSkills returns the skills a tool call loads: grimoire_get_skill's
// name/names (MCP mode) or the Skill tool's skill (skills mode).
func loadedSkills(tool string, input json.RawMessage) []string {
	var in struct {
		Name  string   `json:"name"`
		Names []string `json:"names"`
		Skill string   `json:"skill"`
	}
	if json.Unmarshal(input, &in) != nil {
		return nil
	}
	var names []string
	switch tool {
	case mcpGetSkillTool:
		names = append([]string{in.Name}, in.Names...)
	case "Skill":
		names = []string{in.Skill}
	}
	var out []string
	for _, n := range names {
		n = strings.TrimPrefix(strings.TrimSpace(n), "grimoire:")
		if i := strings.LastIndex(n, ":"); i >= 0 {
			n = n[i+1:]
		}
		if n != "" {
			out = append(out, n)
		}
	}
	return out
}

func writeJSON(t *testing.T, path string, v any) {
	t.Helper()
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatal(err)
	}
}

// The eval's scoring is deterministic and tested on its own, so a failing
// case means the display regressed, not the harness.
func TestScoreDisplay(t *testing.T) {
	decision := "grimoire ▸ start-best-practice → suggest-best-practice (problem described)\n\nHere is the plan…"
	block := "**apply-premortem**  `business/strategy`\n⎿  **Top risk**  DB pool"
	prose := "grimoire ▸ " + strings.Repeat("answer prose ", 30)
	cases := []struct {
		name, reply, level string
		routes             bool
		want               map[string]bool
	}{
		{"decision line then answer", decision, "brief", true, map[string]bool{"row-length": true, "format": false, "shown": true, "clean": true, "hierarchy": true}},
		{"block header", block, "brief", true, map[string]bool{"row-length": true, "format": false, "shown": true, "clean": true, "hierarchy": false}},
		{"plain row labels", "**apply-premortem**  `business/strategy`\n⎿  Top risk  DB pool", "brief", true, map[string]bool{"row-length": true, "format": false, "shown": true, "clean": true, "hierarchy": true}},
		{"routed but silent", "Here is a plan.", "brief", true, map[string]bool{"shown": false, "clean": true}},
		{"plain route line only", "start-best-practice ▸ suggest-best-practice  `meta/routing` — x\n\nWhat matters most?", "brief", true, map[string]bool{"row-length": true, "format": true, "shown": true, "clean": true, "hierarchy": true}},
		{"no route stays silent", "Canberra.", "brief", false, map[string]bool{"silent": true, "clean": true}},
		{"hidden but shown", decision, "hidden", true, map[string]bool{"row-length": true, "format": false, "silent": false, "clean": true, "hierarchy": true}},
		{"prefixed prose", prose, "brief", true, map[string]bool{"row-length": false, "format": false, "shown": true, "clean": false, "hierarchy": true}},
		{"full with source", block + "\n_Source: Klein (1998)_", "full", true, map[string]bool{"row-length": true, "format": true, "shown": true, "clean": true, "evidence": true, "hierarchy": false}},
		{"full without evidence", block, "full", true, map[string]bool{"row-length": true, "format": true, "shown": true, "clean": true, "evidence": false, "hierarchy": false}},
		{"full with bold risk line", block + "\n⎿  ⚠️ **Without this:** no backup", "full", true, map[string]bool{"row-length": true, "format": true, "shown": true, "clean": true, "evidence": true, "hierarchy": false}},
		{"text-glyph risk line", "**apply-x**  `a/b`\n⚠ Without this: no backup", "full", true, map[string]bool{"row-length": true, "format": true, "shown": true, "clean": true, "evidence": true, "hierarchy": true}},
		{"risk line in own words", "**apply-x**  `a/b`\n⚠ Without measuring first, you risk the wrong fix", "full", true, map[string]bool{"row-length": true, "format": true, "shown": true, "clean": true, "evidence": true, "hierarchy": true}},
		{"warning sign mid-sentence is not a risk row", "**apply-x**  `a/b`\nSee the ⚠ above.", "full", true, map[string]bool{"row-length": true, "format": true, "shown": true, "clean": true, "evidence": false, "hierarchy": true}},
	}
	for _, c := range cases {
		got := scoreDisplay(c.reply, c.routes, c.level)
		if len(got) != len(c.want) {
			t.Errorf("%s: checks %v, want %v", c.name, got, c.want)
			continue
		}
		for k, v := range c.want {
			if got[k] != v {
				t.Errorf("%s: %s = %v, want %v", c.name, k, got[k], v)
			}
		}
	}
}

// Skills mode links the meta skills and the chosen subdomains' skills, as
// grimoire install would, and skips subdomains the package lacks.
func TestLinkEvalSkills(t *testing.T) {
	setupListSkillsHome(t)
	t.Setenv("GRIMOIRE_DISPLAY_EVAL_SUBDOMAINS", "engineering/testing,writing/technical,nope/missing")
	dst := filepath.Join(t.TempDir(), "skills")
	n, err := linkEvalSkills(dst)
	if err != nil || n != 2 {
		t.Fatalf("linked %d, err %v; want 2", n, err)
	}
	for _, name := range []string{"write-unit-test", "write-readme"} {
		if _, err := os.Stat(filepath.Join(dst, name, "SKILL.md")); err != nil {
			t.Errorf("%s not linked: %v", name, err)
		}
	}
}

func TestStepCounter(t *testing.T) {
	for s, want := range map[string]bool{
		"**design-slo**  `engineering/reliability` — step 2 of 4": true,
		"Step 1 of 3: premortem":                                  true,
		"⎿  Step 1  apply-premortem":                              false,
		"stepped 2 of 4":                                          false,
	} {
		if got := stepCounter.MatchString(s); got != want {
			t.Errorf("%q: %v, want %v", s, got, want)
		}
	}
}

// The eval workspace is a git repo with a staged change.
func TestSeedWorkspace(t *testing.T) {
	dir := t.TempDir()
	if err := seedWorkspace(dir); err != nil {
		t.Fatal(err)
	}
	out, err := exec.Command("git", "-C", dir, "diff", "--cached", "--name-only").Output()
	if err != nil || strings.TrimSpace(string(out)) != "README.md" {
		t.Errorf("staged = %q, err %v", out, err)
	}
}

// A question to the user ends its turn; a recommendation is marked once.
func TestScoreTurns(t *testing.T) {
	plan := "grimoire ▸ planning several practices — x\n\n**Plan** — 2 steps\n⎿  Step 1  apply-premortem  `business/strategy` — a\nApply this plan?"
	ranOn := plan + "\n\n**apply-premortem**  `business/strategy` — step 1 of 2\n⎿  Top risk  x"
	options := "**Options**\n1. a  `x/y` — p ★ recommended\n2. b  `x/y` — q\nReply 1–2, \"all\", or \"skip\"."
	cases := []struct {
		name  string
		turns []evalTurn
		want  map[string]bool
	}{
		{"plan stops", []evalTurn{{"p", []string{plan}}}, map[string]bool{"one-route": true, "stop": true}},
		{"plan runs on", []evalTurn{{"p", []string{ranOn}}}, map[string]bool{"one-route": true, "stop": false}},
		{"step after question in a later message", []evalTurn{{"p", []string{plan, "**apply-premortem**  `business/strategy` — step 1 of 2"}}}, map[string]bool{"one-route": true, "stop": false}},
		{"step in the next turn is fine", []evalTurn{{"p", []string{plan}}, {"yes", []string{"**apply-premortem**  `business/strategy` — step 1 of 2"}}}, map[string]bool{"one-route": true, "stop": true}},
		{"old reply-with-a-number then block", []evalTurn{{"p", []string{"Which one? Reply with a number.\n**x**  `a/b` — step 1 of 5"}}}, map[string]bool{"stop": false}},
		{"options, one star", []evalTurn{{"p", []string{options}}}, map[string]bool{"stop": true, "one-star": true}},
		{"two stars", []evalTurn{{"p", []string{"1. a  `x/y` ★ recommended\n2. b  `x/y` ★ recommended"}}}, map[string]bool{"one-star": false}},
		{"why-star evidence row is not a second mark", []evalTurn{{"p", []string{"1. a  `x/y` ★ recommended\n⎿  Why ★  fastest"}}}, map[string]bool{"one-star": true}},
		{"no question, no star", []evalTurn{{"p", []string{"Canberra."}}}, map[string]bool{}},
		{"plan with choices", []evalTurn{{"p", []string{planChoices}}}, map[string]bool{"stop": true, "one-star": true, "one-route": true}},
		{"prose after the question", []evalTurn{{"p", []string{plan + "\n\nWhy this order: risks first."}}}, map[string]bool{"stop": false, "one-route": true}},
		{"route repeated across messages", []evalTurn{{"p", []string{"grimoire ▸ a — b", "grimoire ▸ c — d"}}}, map[string]bool{"one-route": false}},
		{"one route per turn", []evalTurn{{"p", []string{"grimoire ▸ a — b"}}, {"yes", []string{"grimoire ▸ applying — you chose all"}}}, map[string]bool{"one-route": true}},
	}
	for _, c := range cases {
		got := scoreTurns(c.turns)
		if len(got) != len(c.want) {
			t.Errorf("%s: %v, want %v", c.name, got, c.want)
			continue
		}
		for k, v := range c.want {
			if got[k] != v {
				t.Errorf("%s: %s = %v, want %v", c.name, k, got[k], v)
			}
		}
	}
}

// Saved transcripts mark each turn so a reader can tell them apart.
func TestSaveTranscript_TurnMarkers(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GRIMOIRE_DISPLAY_EVAL_OUT", dir)
	saveTranscript(t, "mcp", 3, "brief", []evalTurn{{"plan the launch", []string{"Apply this plan?"}}, {"yes", []string{"step 1 of 2"}}})
	b, err := os.ReadFile(filepath.Join(dir, "mcp-case-03.md"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	i, j := strings.Index(s, "=== user: plan the launch ==="), strings.Index(s, "=== user: yes ===")
	if i < 0 || j < i || !strings.Contains(s[i:j], "Apply this plan?") || !strings.Contains(s[j:], "step 1 of 2") {
		t.Errorf("transcript:\n%s", s)
	}
}

// planChoices is the Plan template with a choice step and a nested step.
const planChoices = "grimoire ▸ planning several skills — launch needs them in order\n\n**Plan** — 4 steps\n" +
	"⎿  Step 1  apply-premortem  `business/strategy` — find likeliest failures\n" +
	"⎿  Step 2  load readiness — choose one\n" +
	"     2a. run-load-test  `engineering/performance` — launch traffic ★ recommended\n" +
	"     2b. calculate-capacity-plan  `engineering/performance` — headroom\n" +
	"⎿  Step 3  release safety — nested plan, 2 skills\n" +
	"     3a. design-deployment-strategy  `engineering/devops` — rollback\n" +
	"     3b. apply-feature-toggles  `engineering/devops` — kill switches\n" +
	"⎿  Step 4  No match  \"status page\" — general knowledge\n" +
	"Apply this plan? Reply \"yes\", or pick another choice (e.g. \"2b\"). Steps that change files or run commands ask first."

// planBrief is the Plan template at brief: no domains.
var planBrief = regexp.MustCompile(" {2}`[a-z0-9-]+/[a-z0-9-]+`").ReplaceAllString(planChoices, "")

// The format check holds a reply to the Display templates' fixed parts and
// leaves answer prose alone.
func TestFollowsTemplates(t *testing.T) {
	five := "**Options**\n1. a  `x/y` — p\n2. b  `x/y` — p\n3. c  `x/y` — p\n4. d  `x/y` — p\n5. e  `x/y` — p"
	for _, c := range []struct {
		name, reply, level string
		want               bool
	}{
		{"template reply", "grimoire ▸ planning — why\n\n**Plan** — 1 steps\n⎿  Step 1  apply-premortem  `business/strategy` — a\nApply this plan?", "full", true},
		{"five options", five, "full", false},
		{"four options then prose", "**Options**\n1. a  `x/y` — p\n2. b  `x/y` — p\n3. c  `x/y` — p\n4. d  `x/y` — p\n\nReply 1–4.\n1. later list\n2. x\n3. y\n4. z\n5. w", "full", true},
		{"top-level domain only", "⎿  Step 1  write-adr  `engineering` — a", "full", false},
		{"colon route", "grimoire ▸ finding practices for launch risk: a worry", "full", false},
		{"answer list with code", "1. Run `go test`\n2. Read `cmd/hook.go`", "full", true},
		{"plan with choice and nested steps", planChoices, "full", true},
		{"choice sub-row with top-level domain", "⎿  Step 2  load readiness — choose one\n     2a. run-load-test  `engineering` — x ★ recommended", "full", false},
		{"four choices", "⎿  Step 2  x — choose one\n     2a. a  `x/y` — p\n     2b. b  `x/y` — p\n     2c. c  `x/y` — p\n     2d. d  `x/y` — p", "full", false},
		{"brief plan without domains", strings.ReplaceAll(planChoices, "  `engineering/performance`", ""), "brief", false},
		{"brief plan, no domains at all", planBrief, "brief", true},
		{"domain at brief", "⎿  Step 1  write-sre-runbook  `engineering/reliability` — x", "brief", false},
		{"capitalized route", "grimoire ▸ Planning on-call — two linked skills", "brief", false},
		{"rationale between plan and question", "**Plan** — 1 steps\n⎿  Step 1  apply-premortem — a\n\nWhy this order: risk first.\n\nApply this plan?", "brief", false},
	} {
		if got := followsTemplates(c.reply, c.level); got != c.want {
			t.Errorf("%s: %v, want %v", c.name, got, c.want)
		}
	}
	for _, e := range []string{"⎿  Why ★  it finds unknown risks", "     Why ★  fits 48 hours"} {
		if !evidenceMarker.MatchString(e) {
			t.Errorf("%q is evidence", e)
		}
	}
}

// A ⎿ row is one terminal line; a paragraph in a row fails.
func TestRowsFit(t *testing.T) {
	if !rowsFit("⎿  Top risk  DB connection exhaustion under concurrent signups") {
		t.Error("short row")
	}
	if rowsFit("⎿  Check first  " + strings.Repeat("word ", 40)) {
		t.Error("paragraph row passed")
	}
	if rowsFit("     2a. run-load-test  `engineering/performance` — " + strings.Repeat("why ", 40)) {
		t.Error("long sub-row passed")
	}
	if rowsFit("grimoire ▸ planning — " + strings.Repeat("why ", 40)) {
		t.Error("long route line passed")
	}
	if !rowsFit(strings.Repeat("prose ", 60)) {
		t.Error("answer prose is not a row")
	}
}

// skillRowRef matches a skill named in an option, plan row or block header,
// with its domain span if any: captures name, domain.
var skillRowRef = regexp.MustCompile("(?m)^(?:\\s*(?:\\d+[a-z]?\\.|⎿\\s+Step \\d+)\\s+|\\*\\*)([a-z0-9]+(?:-[a-z0-9]+)+)(?:\\*\\*)?(?:  `([^`]+)`)?")

// accurateSkills reports whether the skills a reply shows are shown truly:
// each domain span is the skill's real domain, and no option or plan row
// offers a meta workflow. checked is false when the reply names no indexed
// skill (or there is no index).
func accurateSkills(reply string, index map[string]string) (ok, checked bool) {
	if index == nil {
		return false, false
	}
	ok = true
	for _, m := range skillRowRef.FindAllStringSubmatch(reply, -1) {
		domain, known := index[m[1]]
		if !known {
			continue
		}
		checked = true
		header := strings.HasPrefix(strings.TrimSpace(m[0]), "**")
		if (m[2] != "" && m[2] != domain) || (domain == "meta" && !header) {
			ok = false
		}
	}
	return ok, checked
}

var (
	skillIndexOnce sync.Once
	skillIndex     map[string]string
)

// evalSkillIndex maps each skill in the GRIMOIRE_SOURCE package to its
// domain/subdomain ("meta" for meta workflows), or nil without a source.
func evalSkillIndex() map[string]string {
	skillIndexOnce.Do(func() {
		if root := os.Getenv("GRIMOIRE_SOURCE"); root != "" {
			skillIndex = buildSkillIndex(filepath.Join(root, "skills"))
		}
	})
	return skillIndex
}

// buildSkillIndex walks a package's skills tree:
// <domain>/<subdomain>/skills/<name>/SKILL.md, or <domain>/skills/<name>/SKILL.md.
func buildSkillIndex(root string) map[string]string {
	index := map[string]string{}
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || d.Name() != "SKILL.md" {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		parts := strings.Split(filepath.ToSlash(rel), "/")
		switch {
		case len(parts) == 5 && parts[2] == "skills":
			index[parts[3]] = parts[0] + "/" + parts[1]
		case len(parts) == 4 && parts[1] == "skills":
			index[parts[2]] = parts[0]
		}
		return nil
	})
	return index
}

// Shown domains must be the real ones, and meta workflows are never offered.
func TestAccurateSkills(t *testing.T) {
	index := map[string]string{"write-sre-runbook": "engineering/reliability", "adapt-best-practice": "meta"}
	for _, c := range []struct {
		name, reply string
		ok, checked bool
	}{
		{"right domain", "⎿  Step 1  write-sre-runbook  `engineering/reliability` — x", true, true},
		{"guessed domain", "⎿  Step 1  write-sre-runbook  `engineering/sre` — x", false, true},
		{"no domain at brief", "     2a. write-sre-runbook — x ★ recommended", true, true},
		{"meta offered", "3. adapt-best-practice — fit it to the team", false, true},
		{"meta block header is fine", "**adapt-best-practice** — step 1 of 2", true, true},
		{"unknown names skipped", "1. pre-launch — check", true, false},
	} {
		ok, checked := accurateSkills(c.reply, index)
		if ok != c.ok || checked != c.checked {
			t.Errorf("%s: ok=%v checked=%v", c.name, ok, checked)
		}
	}
	idx := buildSkillIndex(filepath.Join("..", "nonexistent"))
	if len(idx) != 0 {
		t.Errorf("index of missing dir: %v", idx)
	}
}

// choiceRow matches a plan choice sub-row: step number, letter, skill, rest.
var choiceRow = regexp.MustCompile(`(?m)^\s+(\d+)([a-z])\.\s+([a-z0-9]+(?:-[a-z0-9]+)+)(.*)$`)

// swappedChoice reports whether, after the user picked choice "b" of the
// first choice step in plan, the followup applies that skill's block and not
// the recommended one's. applies is false when plan has no choice step with
// a ★ choice and a b choice.
func swappedChoice(plan, follow string) (ok, applies bool) {
	rows := choiceRow.FindAllStringSubmatch(plan, -1)
	starred, step, picked := "", "", ""
	for _, m := range rows { // the choice step is the one with a ★
		if strings.Contains(m[4], "★") {
			starred, step = m[3], m[1]
			break
		}
	}
	for _, m := range rows {
		if m[1] == step && m[2] == "b" {
			picked = m[3]
		}
	}
	if starred == "" || picked == "" || starred == picked {
		return false, false
	}
	header := func(name string) bool {
		return regexp.MustCompile(`(?m)^\*\*` + regexp.QuoteMeta(name) + `\*\*`).MatchString(follow)
	}
	return header(picked) && !header(starred), true
}

// Picking "b" runs b's skill, not the recommended one.
func TestSwappedChoice(t *testing.T) {
	plan := "⎿  Step 3  rollout — choose one\n     3a. apply-strangler-fig-pattern — x ★ recommended\n     3b. apply-parallel-change — y"
	if ok, applies := swappedChoice(plan, "**apply-parallel-change** — step 3 of 5\n⎿  Plan  expand, migrate, contract"); !ok || !applies {
		t.Errorf("swap honored: ok=%v applies=%v", ok, applies)
	}
	if ok, _ := swappedChoice(plan, "**apply-strangler-fig-pattern** — step 3 of 5"); ok {
		t.Error("recommended skill ran despite the swap")
	}
	nested := "⎿  Step 3  ci — nested plan, 2 skills\n     3a. design-ci-pipeline — x\n     3b. design-cicd-pipeline-security — y\n" + plan
	if ok, applies := swappedChoice(nested, "**apply-parallel-change** — step 3 of 5"); !ok || !applies {
		t.Errorf("a nested step's b is not the choice: ok=%v applies=%v", ok, applies)
	}
	if _, applies := swappedChoice("⎿  Step 1  write-adr — x", "anything"); applies {
		t.Error("no choice step: check should not apply")
	}
}
