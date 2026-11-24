package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

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
}

// evalRun runs one prompt at a display level — in a new session, or resuming
// session resume — and returns the assistant's text messages, the run's cost
// and its session id.
type evalRun func(prompt, level, resume string) (texts []string, cost float64, session string, err error)

// runDisplayEvalMode runs and scores every case in one mode, and logs the
// pass rate per check.
func runDisplayEvalMode(t *testing.T, mode string, cases []displayEvalCase, run evalRun) float64 {
	type tally struct{ pass, total int }
	checks := map[string]*tally{"shown": {}, "silent": {}, "clean": {}, "evidence": {}, "hierarchy": {}, "progress": {}}
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
		saveTranscript(t, mode, i+1, level, c.Prompt, texts)
		scores := scoreDisplay(strings.Join(texts, "\n"), c.Routes, level)
		if c.Followup != "" && level != "off" {
			scores["progress"] = stepCounter.MatchString(strings.Join(followTexts, "\n"))
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
	for _, name := range []string{"shown", "silent", "clean", "evidence", "hierarchy", "progress"} {
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
		f.Close()
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
	if routes && level != "off" {
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

// saveTranscript keeps a case's reply under GRIMOIRE_DISPLAY_EVAL_OUT.
func saveTranscript(t *testing.T, mode string, n int, level, prompt string, texts []string) {
	dir := os.Getenv("GRIMOIRE_DISPLAY_EVAL_OUT")
	if dir == "" {
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Logf("transcript dir: %v", err)
		return
	}
	body := fmt.Sprintf("mode: %s\nlevel: %s\nprompt: %s\n\n%s\n", mode, level, prompt, strings.Join(texts, "\n\n---\n\n"))
	if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s-case-%02d.md", mode, n)), []byte(body), 0o644); err != nil {
		t.Logf("transcript: %v", err)
	}
}

// displayMarker matches grimoire's display: a decision line or a block header
// (bold skill name followed by a domain code span).
// The route line is plain text under the Display spec (start-best-practice ▸
// workflow), so it is matched on its own, not only through bold headers.
var displayMarker = regexp.MustCompile("(?m)grimoire ▸|^\\s*start-best-practice ▸|^\\*\\*[a-z0-9-]+\\*\\*[^\\n]*`[a-z0-9-]+/[a-z0-9-]+`")

// stepCounter matches a plan step's position, "step 2 of 4".
var stepCounter = regexp.MustCompile(`(?i)\bstep \d+ of \d+\b`)

// evidenceMarker matches the full level's evidence: a source, severity or risk.
// A risk is any line led by ⚠ (the Display section's risk glyph), whatever
// words follow — "⚠ Without this:" is the template, not the only phrasing.
var evidenceMarker = regexp.MustCompile(`(?im)_?source:|🔴|🟡|⚪|^\s*(⎿\s*)?⚠`)

// runClaudeEval runs one prompt headless in work with extra claude args —
// resuming session resume when set — and returns the assistant's text
// messages, the run's cost and its session id.
func runClaudeEval(work, level, prompt, resume string, args ...string) (texts []string, cost float64, session string, err error) {
	r, err := runClaudeEvalResult(work, level, prompt, resume, args...)
	return r.texts, r.cost, r.session, err
}

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
		{"decision line then answer", decision, "brief", true, map[string]bool{"shown": true, "clean": true, "hierarchy": true}},
		{"block header", block, "brief", true, map[string]bool{"shown": true, "clean": true, "hierarchy": false}},
		{"plain row labels", "**apply-premortem**  `business/strategy`\n⎿  Top risk  DB pool", "brief", true, map[string]bool{"shown": true, "clean": true, "hierarchy": true}},
		{"routed but silent", "Here is a plan.", "brief", true, map[string]bool{"shown": false, "clean": true}},
		{"plain route line only", "start-best-practice ▸ suggest-best-practice  `meta/routing` — x\n\nWhat matters most?", "brief", true, map[string]bool{"shown": true, "clean": true, "hierarchy": true}},
		{"no route stays silent", "Canberra.", "brief", false, map[string]bool{"silent": true, "clean": true}},
		{"off but shown", decision, "off", true, map[string]bool{"silent": false, "clean": true, "hierarchy": true}},
		{"prefixed prose", prose, "brief", true, map[string]bool{"shown": true, "clean": false, "hierarchy": true}},
		{"full with source", block + "\n_Source: Klein (1998)_", "full", true, map[string]bool{"shown": true, "clean": true, "evidence": true, "hierarchy": false}},
		{"full without evidence", block, "full", true, map[string]bool{"shown": true, "clean": true, "evidence": false, "hierarchy": false}},
		{"full with bold risk line", block + "\n⎿  ⚠️ **Without this:** no backup", "full", true, map[string]bool{"shown": true, "clean": true, "evidence": true, "hierarchy": false}},
		{"text-glyph risk line", "**x**  `a/b`\n⚠ Without this: no backup", "full", true, map[string]bool{"shown": true, "clean": true, "evidence": true, "hierarchy": true}},
		{"risk line in own words", "**x**  `a/b`\n⚠ Without measuring first, you risk the wrong fix", "full", true, map[string]bool{"shown": true, "clean": true, "evidence": true, "hierarchy": true}},
		{"warning sign mid-sentence is not a risk row", "**x**  `a/b`\nSee the ⚠ above.", "full", true, map[string]bool{"shown": true, "clean": true, "evidence": false, "hierarchy": true}},
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
