package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jeffreytse/grimoire/internal/agent"
)

// TestOutcomeEval measures whether grimoire makes the agent's work better
// (docs/outcome-eval-plan.md): each task runs twice with the same model,
// workspace and permissions — arm A is Claude Code alone, arm B is Claude
// Code with grimoire in MCP mode as installed — and a blind judge scores each
// final result against a rubric written from the practice's original source.
//
//	GRIMOIRE_OUTCOME_EVAL=1 GRIMOIRE_SOURCE=<skills repo> \
//	  go test ./cmd -run TestOutcomeEval -v -timeout 120m
//
// Options:
//
//	GRIMOIRE_OUTCOME_EVAL_MODEL=<m>  agent model for both arms (default claude-opus-5-5)
//	GRIMOIRE_OUTCOME_JUDGE_MODEL=<m> judge model (default claude-sonnet-5-5: a
//	                                 different model from the one under test)
//	GRIMOIRE_OUTCOME_EVAL_ONLY=a,b   task ids to run
//	GRIMOIRE_OUTCOME_EVAL_RUNS=3     runs per task per arm (default 1)
//	GRIMOIRE_OUTCOME_EVAL_OUT=<dir>  keep transcripts, judged outputs and verdicts
//	GRIMOIRE_OUTCOME_RESCORE=<dir>   re-judge a kept run without re-running agents
//	GRIMOIRE_OUTCOME_REPORT=<dir>    recompute the report from a kept run's verdicts
//	                                 (no API calls)
//	GRIMOIRE_OUTCOME_NEUTRALIZE=0    skip the neutral-style rewrite of replies
//	                                 (on by default; see neutralize)
func TestOutcomeEval(t *testing.T) {
	if dir := os.Getenv("GRIMOIRE_OUTCOME_REPORT"); dir != "" {
		logOutcome(t, loadVerdicts(t, dir))
		return
	}
	rescore := os.Getenv("GRIMOIRE_OUTCOME_RESCORE")
	if os.Getenv("GRIMOIRE_OUTCOME_EVAL") == "" && rescore == "" {
		t.Skip("set GRIMOIRE_OUTCOME_EVAL=1 to run against Claude Code (uses API credits)")
	}
	if _, err := exec.LookPath("claude"); err != nil {
		t.Skip("claude CLI not on PATH")
	}
	tasks := loadOutcomeTasks(t)
	if only := os.Getenv("GRIMOIRE_OUTCOME_EVAL_ONLY"); only != "" {
		keep := map[string]bool{}
		for _, id := range strings.Split(only, ",") {
			keep[strings.TrimSpace(id)] = true
		}
		var picked []outcomeTask
		for i := range tasks {
			if keep[tasks[i].ID] {
				picked = append(picked, tasks[i])
			}
		}
		tasks = picked
	}
	model := envOr("GRIMOIRE_OUTCOME_EVAL_MODEL", "claude-opus-5-5")
	judgeModel := envOr("GRIMOIRE_OUTCOME_JUDGE_MODEL", "claude-sonnet-5-5")
	runs := 1
	if n, err := strconv.Atoi(os.Getenv("GRIMOIRE_OUTCOME_EVAL_RUNS")); err == nil && n > 0 {
		runs = n
	}
	out := os.Getenv("GRIMOIRE_OUTCOME_EVAL_OUT")
	if rescore != "" {
		out = rescore
	}

	var arms map[string]outcomeArm
	if rescore == "" {
		dir := t.TempDir()
		bin := filepath.Join(dir, "grimoire")
		if b, err := exec.Command("go", "build", "-o", bin, "..").CombinedOutput(); err != nil {
			t.Fatalf("build: %v: %s", err, b)
		}
		arms = outcomeArms(t, dir, bin, model)
	}

	var results []outcomeResult
	var cost float64
	fallbacks := 0
	for i := range tasks {
		task := &tasks[i]
		for run := 1; run <= runs; run++ {
			key := fmt.Sprintf("%s-run%d", task.ID, run)
			outputs := map[string]string{}
			for _, name := range []string{"A", "B"} {
				path := filepath.Join(out, key+"-"+name+".md")
				if rescore != "" {
					b, err := os.ReadFile(path)
					if err != nil {
						t.Fatalf("rescore: %v", err)
					}
					outputs[name] = reblind(string(b))
					continue
				}
				o, usd, err := arms[name](task)
				cost += usd
				if err != nil {
					t.Logf("ERROR %s arm %s: %v", key, name, err)
				}
				outputs[name] = o
				if out != "" {
					writeOutcomeFile(t, path, o)
				}
			}
			if os.Getenv("GRIMOIRE_OUTCOME_NEUTRALIZE") != "0" {
				for _, name := range []string{"A", "B"} {
					cache := ""
					if out != "" {
						cache = filepath.Join(out, key+"-"+name+".neutral.md")
					}
					n, usd, kept := neutralize(outputs[name], judgeModel, cache)
					cost += usd
					if !kept {
						fallbacks++
						t.Logf("NOTE %s arm %s: rewrite changed content; judging the plain blinded text", key, name)
					}
					outputs[name] = n
				}
			}
			v, usd, err := judgePair(task, outputs["A"], outputs["B"], judgeModel)
			cost += usd
			if err != nil {
				t.Logf("ERROR %s judge: %v", key, err)
				continue
			}
			if out != "" {
				b, _ := json.MarshalIndent(v, "", "  ")
				writeOutcomeFile(t, filepath.Join(out, key+"-verdict.json"), string(b))
			}
			results = append(results, outcomeResult{task: task.ID, field: task.Field, verdict: v})
			t.Logf("%-28s A %.0f%%  B %.0f%%  better=%s  order-disagreements=%d", key, 100*v.ARate, 100*v.BRate, v.Better, v.Disagree)
		}
	}
	logOutcome(t, results)
	t.Logf("neutral rewrites rejected by the fidelity check: %d", fallbacks)
	t.Logf("cost=$%.2f (agents %s, judge %s)", cost, model, judgeModel)
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// outcomeTask is one eval task (cmd/testdata/outcome_eval.json).
type outcomeTask struct {
	ID       string            `json:"id"`
	Field    string            `json:"field"`
	Practice string            `json:"practice"`
	Prompt   string            `json:"prompt"`
	Files    map[string]string `json:"files"`
	Staged   map[string]string `json:"staged"`
	Exec     string            `json:"exec"`
	Source   string            `json:"source"`
	Rubric   []string          `json:"rubric"`
}

func loadOutcomeTasks(t *testing.T) []outcomeTask {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", "outcome_eval.json"))
	if err != nil {
		t.Fatal(err)
	}
	var set struct {
		Tasks []outcomeTask `json:"tasks"`
	}
	if err := json.Unmarshal(raw, &set); err != nil {
		t.Fatal(err)
	}
	return set.Tasks
}

// outcomeArm runs a task and returns the judged output: the final turn's
// reply plus every file the agent created or changed.
type outcomeArm func(task *outcomeTask) (output string, cost float64, err error)

// outcomeConfirm is the one reply the harness gives when an arm ends its turn
// with a question or plan — the same for both arms.
const outcomeConfirm = "Yes — go ahead with your recommendation."

// outcomeArms returns arm A (Claude Code alone) and arm B (with grimoire in
// MCP mode, configured as `grimoire mcp install` configures it). Both arms
// use the same model, a fresh seeded workspace and acceptEdits permissions.
func outcomeArms(t *testing.T, dir, bin, model string) map[string]outcomeArm {
	hook := func(args string) any {
		return []any{map[string]any{"hooks": []any{map[string]any{"type": "command", "command": bin + " " + args}}}}
	}
	settings := filepath.Join(dir, "outcome-settings.json")
	mcpConfig := filepath.Join(dir, "outcome-mcp.json")
	writeJSON(t, settings, map[string]any{
		"hooks": map[string]any{
			"SessionStart":     hook("hook session"),
			"UserPromptSubmit": hook("hook prompt"),
		},
		"permissions": map[string]any{"allow": agent.ToolAllowRules("claude")},
	})
	writeJSON(t, mcpConfig, map[string]any{"mcpServers": map[string]any{"grimoire": map[string]any{"command": bin, "args": []string{"mcp", "serve"}}}})
	common := []string{"--model", model, "--permission-mode", "acceptEdits", "--setting-sources", "project"}
	armArgs := map[string][]string{
		"A": append(append([]string{}, common...), "--mcp-config", `{"mcpServers":{}}`, "--strict-mcp-config"),
		"B": append(append([]string{}, common...), "--settings", settings, "--mcp-config", mcpConfig, "--strict-mcp-config"),
	}
	arms := map[string]outcomeArm{}
	for name, args := range armArgs {
		args := args
		arms[name] = func(task *outcomeTask) (string, float64, error) {
			work, err := os.MkdirTemp(dir, task.ID+"-"+name)
			if err != nil {
				return "", 0, err
			}
			if err := seedOutcomeWorkspace(work, task); err != nil {
				return "", 0, err
			}
			r, err := runClaudeEvalResult(work, "brief", task.Prompt, "", args...)
			cost, final := r.cost, r.texts
			if err == nil && endsWithQuestion(strings.Join(r.texts, "\n")) && r.session != "" {
				r2, err2 := runClaudeEvalResult(work, "brief", outcomeConfirm, r.session, args...)
				cost += r2.cost
				if err2 == nil {
					final = r2.texts
				}
			}
			return outcomeOutput(final, workspaceChanges(work), execCheck(work, task.Exec)), cost, err
		}
	}
	return arms
}

// seedOutcomeWorkspace writes a task's files into a git repository, commits
// them, then writes and stages the task's staged changes.
func seedOutcomeWorkspace(dir string, task *outcomeTask) error {
	write := func(files map[string]string) error {
		for name, body := range files {
			p := filepath.Join(dir, filepath.FromSlash(name))
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
				return err
			}
		}
		return nil
	}
	if err := write(task.Files); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, ".gitkeep"), nil, 0o644); err != nil {
		return err
	}
	for _, args := range [][]string{{"init", "-q"}, {"add", "."}, {"commit", "-q", "-m", "initial"}} {
		if err := outcomeGit(dir, args...); err != nil {
			return err
		}
	}
	if len(task.Staged) == 0 {
		return nil
	}
	if err := write(task.Staged); err != nil {
		return err
	}
	return outcomeGit(dir, "add", ".")
}

func outcomeGit(dir string, args ...string) error {
	c := exec.Command("git", args...)
	c.Dir = dir
	c.Env = append(os.Environ(), "GIT_AUTHOR_NAME=eval", "GIT_AUTHOR_EMAIL=eval@example.invalid",
		"GIT_COMMITTER_NAME=eval", "GIT_COMMITTER_EMAIL=eval@example.invalid")
	if out, err := c.CombinedOutput(); err != nil {
		return fmt.Errorf("git %s: %v: %s", strings.Join(args, " "), err, out)
	}
	return nil
}

// endsWithQuestion reports whether a reply hands the turn back to the user:
// its last line asks a question or is a plan or options prompt.
func endsWithQuestion(reply string) bool {
	lines := strings.Split(strings.TrimSpace(reply), "\n")
	last := strings.TrimSpace(lines[len(lines)-1])
	return strings.HasSuffix(last, "?") || stopMarker.MatchString(last)
}

// Output limits keep a judged output within the judge's context.
const (
	outcomeFileMax  = 6000
	outcomeTotalMax = 24000
)

// workspaceChanges lists the files the agent created or changed since the
// seeded commit, with their contents (staged task changes excluded).
func workspaceChanges(dir string) string {
	c := exec.Command("git", "status", "--porcelain", "--untracked-files=all")
	c.Dir = dir
	out, err := c.Output()
	if err != nil {
		return ""
	}
	var b strings.Builder
	for _, line := range strings.Split(strings.TrimRight(string(out), "\n"), "\n") {
		if len(line) < 4 || line[1] == ' ' && line[0] != '?' {
			continue // staged only (the task's own staged change)
		}
		name := strings.TrimSpace(line[3:])
		body, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		content := string(body)
		if len(content) > outcomeFileMax {
			content = content[:outcomeFileMax] + "\n…(truncated)"
		}
		fmt.Fprintf(&b, "\n--- file: %s ---\n%s\n", name, content)
		if b.Len() > outcomeTotalMax {
			break
		}
	}
	return b.String()
}

// execCheck runs a task's execution grader in the workspace and reports the
// result as a line the judge treats as fact. "go-test": do the tests compile
// and pass?
func execCheck(dir, kind string) string {
	if kind != "go-test" {
		return ""
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	c := exec.CommandContext(ctx, "go", "test", "./...")
	c.Dir = dir
	out, err := c.CombinedOutput()
	result := "PASS"
	if err != nil {
		result = "FAIL"
	}
	tail := strings.TrimSpace(string(out))
	if len(tail) > 1500 {
		tail = tail[len(tail)-1500:]
	}
	return fmt.Sprintf("[execution check: `go test ./...` %s]\n%s", result, tail)
}

// displayLine matches lines of grimoire's display that would reveal which
// arm wrote an output: route lines, practice block headers, plan and options
// headings, and their questions.
var displayLine = regexp.MustCompile(`(?m)^\s*(grimoire ▸.*|\*\*[a-z0-9]+(?:-[a-z0-9]+)+\*\*[^\n]*|\*\*(Plan|Options)\*\*.*|Apply this plan\?.*|Reply \d+\s*[–-]\s*\d+.*)$`)

// displayRow matches a display row, "⎿  <Label>  <value>".
var displayRow = regexp.MustCompile(`(?m)^(\s*)⎿\s+(\S[^\n]*?)\s{2,}(\S[^\n]*)$`)

// blind removes the structure that tells the arms apart, the same way for
// both outputs: grimoire's display lines go, display rows become plain
// "Label: value" lines, and recommendation and step markers are dropped.
func blind(reply string) string {
	reply = displayLine.ReplaceAllString(reply, "")
	reply = displayRow.ReplaceAllString(reply, "${1}${2}: ${3}")
	reply = strings.ReplaceAll(reply, "⎿", "")
	// The Outcome template's "✓ <result>" lines read as a plain list.
	reply = regexp.MustCompile(`(?m)^(\s*)✓\s+`).ReplaceAllString(reply, "${1}- ")
	reply = regexp.MustCompile(`\s*★ recommended`).ReplaceAllString(reply, "")
	reply = regexp.MustCompile(`(?i)\s*— step \d+[a-z]? of \d+`).ReplaceAllString(reply, "")
	return regexp.MustCompile(`\n{3,}`).ReplaceAllString(strings.TrimSpace(reply), "\n\n")
}

// legacyRow matches a row as the first harness version saved it: "⎿"
// replaced by "-", which left "-  Label  value" (two spaces after the dash,
// unlike a markdown list item).
var legacyRow = regexp.MustCompile(`(?m)^(\s*)- {2}(\S[^\n]*?)\s{2,}(\S[^\n]*)$`)

// reblind applies the current blinding to a saved output, including rows
// saved by the first harness version.
func reblind(saved string) string {
	return blind(legacyRow.ReplaceAllString(saved, "${1}${2}: ${3}"))
}

// filesMarker starts the part of an output that is the agent's work product
// (files it created or changed, and execution checks) — never rewritten.
const filesMarker = "\n\n=== files created or changed ==="

// neutralize restates an output's reply in one fixed house style, so the two
// arms can't be told apart by presentation; the judge then compares content.
// Files and execution checks are left exact. A rewrite that alters a code
// block or drops a number is rejected and the blinded original is kept
// (kept=false). Rewrites are cached at cache when it is set.
func neutralize(output, model, cache string) (neutral string, cost float64, kept bool) {
	if cache != "" {
		if b, err := os.ReadFile(cache); err == nil {
			return string(b), 0, true
		}
	}
	reply, rest := output, ""
	if i := strings.Index(output, filesMarker); i >= 0 {
		reply, rest = output[:i], output[i:]
	}
	if strings.TrimSpace(reply) == "" {
		return output, 0, true
	}
	dir, err := os.MkdirTemp("", "grimoire-neutral")
	if err != nil {
		return output, 0, false
	}
	defer func() { _ = os.RemoveAll(dir) }()
	r, err := runClaudeEvalResult(dir, "hidden", neutralPrompt+reply, "", "--model", model, "--max-turns", "1",
		"--setting-sources", "project", "--mcp-config", `{"mcpServers":{}}`, "--strict-mcp-config", "--tools", "")
	rewritten := strings.TrimSpace(strings.Join(r.texts, "\n"))
	if err != nil || rewritten == "" || !faithful(reply, rewritten) {
		return output, r.cost, false
	}
	neutral = rewritten + rest
	if cache != "" {
		_ = os.WriteFile(cache, []byte(neutral), 0o644)
	}
	return neutral, r.cost, true
}

const neutralPrompt = `Rewrite the assistant reply below into a fixed neutral house style, so that replies from different assistants cannot be told apart by their formatting.

Style: plain markdown. "## " headings only where the original has a section. Lists use "- ". No bold or italics, no emoji, no decorative symbols or glyphs, no tables unless the original's content is tabular. Plain sentences.

Content: keep every fact, number, name, recommendation, caveat, question and step, in the same order. Keep every fenced code block byte-for-byte. Do not add, remove, correct or improve anything, and do not comment on the reply. Same language as the original.

Output only the rewritten reply.

REPLY:
`

// fencedBlock matches a fenced code block.
var fencedBlock = regexp.MustCompile("(?s)```[^\n]*\n.*?```")

// numberToken matches a number as written: 3,120 · 24.9% · $1.8M · 0.035.
var numberToken = regexp.MustCompile(`\d[\d,]*(?:\.\d+)?`)

// faithful reports whether a rewrite kept the original's code blocks
// verbatim and every number it states.
func faithful(original, rewritten string) bool {
	for _, block := range fencedBlock.FindAllString(original, -1) {
		body := block[strings.Index(block, "\n")+1:]
		if !strings.Contains(rewritten, strings.TrimSuffix(body, "```")) {
			return false
		}
	}
	have := map[string]bool{}
	for _, n := range numberToken.FindAllString(rewritten, -1) {
		have[strings.ReplaceAll(n, ",", "")] = true
	}
	for _, n := range numberToken.FindAllString(original, -1) {
		if !have[strings.ReplaceAll(n, ",", "")] {
			return false
		}
	}
	return true
}

// outcomeOutput assembles what the judge sees, blinded the same way for both
// arms: the final reply, then the changed files and any execution check.
func outcomeOutput(texts []string, changes, execResult string) string {
	reply := blind(strings.Join(texts, "\n\n"))
	var b strings.Builder
	b.WriteString(reply)
	if changes != "" {
		b.WriteString("\n\n=== files created or changed ===")
		b.WriteString(changes)
	}
	if execResult != "" {
		b.WriteString("\n\n" + execResult)
	}
	return b.String()
}

// outcomeCriterion is the judge's verdict on one rubric criterion.
type outcomeCriterion struct {
	Criterion int    `json:"criterion"`
	Met       bool   `json:"met"`
	Evidence  string `json:"evidence"`
}

type outcomeCriteria []outcomeCriterion

func (cs outcomeCriteria) met() int {
	n := 0
	for _, c := range cs {
		if c.Met {
			n++
		}
	}
	return n
}

func (cs outcomeCriteria) rate() float64 {
	if len(cs) == 0 {
		return 0
	}
	return float64(cs.met()) / float64(len(cs))
}

// outcomeVerdict is one task's judged result, mapped back to the arms.
type outcomeVerdict struct {
	A        outcomeCriteria `json:"a"`
	B        outcomeCriteria `json:"b"`
	Better   string          `json:"better"` // "A", "B" or "tie"
	Why      string          `json:"why"`
	Rejected int             `json:"evidence_rejected"` // "met" verdicts whose quote was not in the output
	AFirst   bool            `json:"a_shown_first"`
	// Guess is the judge's guess at which arm used a best-practice plugin:
	// "A", "B" or "unsure". Guesses near chance mean the blinding holds.
	Guess string `json:"guess"`
}

// outcomePair is a task's result, judged twice — once in each presentation
// order — so a judge's preference for the first or second output cancels
// out. A criterion's score is the mean of the two judgments; a preference
// counts only when both orders agree.
type outcomePair struct {
	Rounds   [2]outcomeVerdict `json:"rounds"`
	ARate    float64           `json:"a_rate"`
	BRate    float64           `json:"b_rate"`
	Better   string            `json:"better"`   // "A", "B" or "tie"
	Disagree int               `json:"disagree"` // criteria the two orders judged differently
	Rejected int               `json:"evidence_rejected"`
}

type outcomeResult struct {
	task, field string
	verdict     outcomePair
	guesses     []string // every round's plugin guess, after averageRuns
}

// judgePair judges a task's two outputs in both presentation orders and
// combines the rounds.
func judgePair(task *outcomeTask, a, b, model string) (outcomePair, float64, error) {
	var p outcomePair
	var cost float64
	for i, aFirst := range []bool{true, false} {
		v, usd, err := judgeOutcome(task, a, b, model, aFirst)
		cost += usd
		if err != nil {
			return p, cost, err
		}
		p.Rounds[i] = v
	}
	return combineRounds(&p.Rounds), cost, nil
}

// combineRounds averages two judgments of the same pair.
func combineRounds(rounds *[2]outcomeVerdict) outcomePair {
	p := outcomePair{Rounds: *rounds}
	r1, r2 := rounds[0], rounds[1]
	p.ARate = (r1.A.rate() + r2.A.rate()) / 2
	p.BRate = (r1.B.rate() + r2.B.rate()) / 2
	p.Better = "tie"
	if r1.Better == r2.Better {
		p.Better = r1.Better
	}
	for _, pair := range [][2]outcomeCriteria{{r1.A, r2.A}, {r1.B, r2.B}} {
		for i := range pair[0] {
			if i < len(pair[1]) && pair[0][i].Met != pair[1][i].Met {
				p.Disagree++
			}
		}
	}
	p.Rejected = r1.Rejected + r2.Rejected
	return p
}

// judgeOutcome asks a different model to score both outputs, shown in random
// order as X and Y, against the task's rubric. Every "met" must quote the
// output; a quote that isn't there turns the verdict into "not met".
func judgeOutcome(task *outcomeTask, a, b, model string, aFirst bool) (outcomeVerdict, float64, error) {
	x, y := b, a
	if aFirst {
		x, y = a, b
	}
	prompt := judgePrompt(task, x, y)
	dir, err := os.MkdirTemp("", "grimoire-judge")
	if err != nil {
		return outcomeVerdict{}, 0, err
	}
	defer func() { _ = os.RemoveAll(dir) }()
	r, err := runClaudeEvalResult(dir, "hidden", prompt, "", "--model", model, "--max-turns", "1",
		"--setting-sources", "project", "--mcp-config", `{"mcpServers":{}}`, "--strict-mcp-config", "--tools", "")
	if err != nil {
		return outcomeVerdict{}, r.cost, err
	}
	raw, err := parseJudge(strings.Join(r.texts, "\n"))
	if err != nil {
		return outcomeVerdict{}, r.cost, err
	}
	v := mapJudge(&raw, aFirst, len(task.Rubric))
	v.Rejected = checkEvidence(v.A, a) + checkEvidence(v.B, b)
	return v, r.cost, nil
}

func judgePrompt(task *outcomeTask, x, y string) string {
	var b strings.Builder
	b.WriteString("You are grading two AI assistants' work on the same task. Judge only what is in each output. Be strict and consistent; do not reward length or formatting for its own sake.\n\n")
	fmt.Fprintf(&b, "TASK given to both:\n%s\n\n", task.Prompt)
	if len(task.Files) > 0 {
		b.WriteString("Files the assistants were given:\n")
		names := make([]string, 0, len(task.Files))
		for n := range task.Files {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, n := range names {
			fmt.Fprintf(&b, "--- %s ---\n%s\n", n, task.Files[n])
		}
		for n, body := range task.Staged {
			fmt.Fprintf(&b, "--- %s (staged change) ---\n%s\n", n, body)
		}
		b.WriteString("\n")
	}
	fmt.Fprintf(&b, "RUBRIC (from %s). For each criterion, decide met or not met for each output:\n", task.Source)
	for i, c := range task.Rubric {
		fmt.Fprintf(&b, "%d. %s\n", i+1, c)
	}
	fmt.Fprintf(&b, "\n=== OUTPUT X ===\n%s\n=== END X ===\n\n=== OUTPUT Y ===\n%s\n=== END Y ===\n\n", x, y)
	b.WriteString(`For every "met", copy a short exact quote (5-25 words) from that output as evidence; if you cannot quote it, it is not met. For a criterion met by something being absent (e.g. "no sleeps or randomness"), write "ABSENCE" as the evidence. An "[execution check]" line is a fact from running the code.
Then say which output better serves the user overall: "X", "Y" or "tie".
Finally, one of the two assistants used a plugin that changes how it presents its work. Judging only by presentation and formatting, not by quality, guess which: "X", "Y", or "unsure" if you can't tell. This guess does not affect the grading.
Reply with only this JSON:
{"x":[{"criterion":1,"met":true,"evidence":"..."}],"y":[{"criterion":1,"met":false,"evidence":""}],"better":"X","why":"one sentence","plugin_guess":"unsure"}`)
	return b.String()
}

// judgeRaw is the judge's JSON in presentation order (X, Y).
type judgeRaw struct {
	X      outcomeCriteria `json:"x"`
	Y      outcomeCriteria `json:"y"`
	Better string          `json:"better"`
	Why    string          `json:"why"`
	Guess  string          `json:"plugin_guess"`
}

// parseJudge extracts the judge's JSON object from its reply.
func parseJudge(reply string) (judgeRaw, error) {
	var raw judgeRaw
	i := strings.Index(reply, "{")
	if i < 0 {
		return raw, fmt.Errorf("no JSON in judge reply: %.200s", reply)
	}
	// Decode the first JSON value only: evidence can hold braces (YAML's
	// ${{ }}), and the judge sometimes adds text after the object.
	if err := json.NewDecoder(strings.NewReader(reply[i:])).Decode(&raw); err != nil {
		return raw, fmt.Errorf("judge JSON: %w: %.200s", err, reply[i:])
	}
	return raw, nil
}

// mapJudge maps presentation order back to arms, keeping one verdict per
// rubric criterion (a missing criterion counts as not met).
func mapJudge(raw *judgeRaw, aFirst bool, n int) outcomeVerdict {
	norm := func(cs outcomeCriteria) outcomeCriteria {
		out := make(outcomeCriteria, n)
		for i := range out {
			out[i] = outcomeCriterion{Criterion: i + 1}
		}
		for _, c := range cs {
			if c.Criterion >= 1 && c.Criterion <= n {
				out[c.Criterion-1] = c
			}
		}
		return out
	}
	v := outcomeVerdict{Why: raw.Why, AFirst: aFirst}
	better := strings.ToUpper(strings.TrimSpace(raw.Better))
	if aFirst {
		v.A, v.B = norm(raw.X), norm(raw.Y)
		v.Better = map[string]string{"X": "A", "Y": "B"}[better]
	} else {
		v.A, v.B = norm(raw.Y), norm(raw.X)
		v.Better = map[string]string{"X": "B", "Y": "A"}[better]
	}
	if v.Better == "" {
		v.Better = "tie"
	}
	guess := strings.ToUpper(strings.TrimSpace(raw.Guess))
	if aFirst {
		v.Guess = map[string]string{"X": "A", "Y": "B"}[guess]
	} else {
		v.Guess = map[string]string{"X": "B", "Y": "A"}[guess]
	}
	if v.Guess == "" {
		v.Guess = "unsure"
	}
	return v
}

// quoteNorm normalizes text for evidence matching: lowercase, markdown and
// punctuation markers dropped, whitespace collapsed.
var quoteStrip = regexp.MustCompile("[*_`#>|]+")

func quoteNorm(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(quoteStrip.ReplaceAllString(s, " "))), " ")
}

// quoteSplit separates the fragments of a quote that cites several places at
// once ("a … b", "a | b").
var quoteSplit = regexp.MustCompile(`\.\.\.+|…| \| `)

// checkEvidence turns every "met" whose quote is not in the output into
// "not met", and returns how many it changed. A quote may cite several
// places, joined by "..." or " | "; each fragment must be in the output.
func checkEvidence(cs outcomeCriteria, output string) int {
	out, rejected := quoteNorm(output), 0
	for i := range cs {
		if !cs[i].Met || quoteFound(cs[i].Evidence, out) {
			continue
		}
		cs[i].Met = false
		rejected++
	}
	return rejected
}

// absenceEvidence is what the judge writes for a criterion met by something
// being absent ("no sleeps or randomness"), which has nothing to quote.
const absenceEvidence = "ABSENCE"

// quoteToken matches the words a quote is checked by.
var quoteToken = regexp.MustCompile(`[a-z0-9]{2,}`)

// quoteFound reports whether every fragment of a quote is in the normalized
// output — verbatim, or with at least 90% of its distinct words present (the
// judge stitches code rows from different lines). Fragments shorter than 4
// characters (stray separators) are ignored.
func quoteFound(evidence, normalizedOutput string) bool {
	if strings.TrimSpace(evidence) == absenceEvidence {
		return true
	}
	words := map[string]bool{}
	for _, w := range quoteToken.FindAllString(normalizedOutput, -1) {
		words[w] = true
	}
	found := 0
	for _, frag := range quoteSplit.Split(evidence, -1) {
		q := quoteNorm(frag)
		if len(q) < 4 {
			continue
		}
		if !strings.Contains(normalizedOutput, q) && !mostWordsPresent(q, words) {
			return false
		}
		found++
	}
	return found > 0
}

// mostWordsPresent reports whether at least 90% of a fragment's distinct
// words (3 or more of them) appear in the output.
func mostWordsPresent(fragment string, words map[string]bool) bool {
	distinct := map[string]bool{}
	for _, w := range quoteToken.FindAllString(fragment, -1) {
		distinct[w] = true
	}
	if len(distinct) < 3 {
		return false
	}
	hit := 0
	for w := range distinct {
		if words[w] {
			hit++
		}
	}
	return float64(hit) >= 0.9*float64(len(distinct))
}

// logOutcome reports the paired difference in criteria met (B − A) with a
// bootstrap 95% CI over tasks, the pairwise win rate, and both by field.
func logOutcome(t *testing.T, results []outcomeResult) {
	if len(results) == 0 {
		t.Log("no results")
		return
	}
	report := func(label string, runs []outcomeResult) {
		rs := averageRuns(runs)
		var diffs []float64
		var a, b float64
		wins, losses, ties, rejected, disagree, guessRight, guessWrong, guessUnsure := 0, 0, 0, 0, 0, 0, 0, 0
		for i := range rs {
			v := &rs[i].verdict
			a += v.ARate
			b += v.BRate
			diffs = append(diffs, v.BRate-v.ARate)
			rejected += v.Rejected
			disagree += v.Disagree
			guesses := rs[i].guesses
			if guesses == nil {
				for _, r := range v.Rounds {
					guesses = append(guesses, r.Guess)
				}
			}
			for _, g := range guesses {
				switch g {
				case "B":
					guessRight++
				case "A":
					guessWrong++
				default:
					guessUnsure++
				}
			}
			switch v.Better {
			case "B":
				wins++
			case "A":
				losses++
			default:
				ties++
			}
		}
		n := float64(len(rs))
		lo, hi := bootstrapMeanCI(diffs, 10000, 1)
		wlo, whi := wilson(wins, wins+losses)
		t.Logf("%-12s tasks=%d runs=%d  criteria met: A %.0f%%  B %.0f%%  diff %+.1f pts [95%% CI %+.1f, %+.1f]  B preferred %d/%d (%.0f%%, CI %.0f–%.0f%%), ties %d",
			label, len(rs), len(runs), 100*a/n, 100*b/n, 100*(b-a)/n, 100*lo, 100*hi,
			wins, wins+losses, pct(wins, wins+losses), 100*wlo, 100*whi, ties)
		t.Logf("%-12s checks: order disagreements %d criteria, evidence rejected %d, blinding guess right %d / wrong %d / unsure %d (near chance = blind)",
			label, disagree, rejected, guessRight, guessWrong, guessUnsure)
	}
	report("all", results)
	byField := map[string][]outcomeResult{}
	var fields []string
	for i := range results {
		f := results[i].field
		if byField[f] == nil {
			fields = append(fields, f)
		}
		byField[f] = append(byField[f], results[i])
	}
	sort.Strings(fields)
	for _, f := range fields {
		report(f, byField[f])
	}
}

// averageRuns collapses repeated runs of a task into one result — the mean
// rates, and the majority preference — so the confidence interval is over
// tasks, not over correlated runs of the same task. Round-level checks
// (guesses, disagreements, rejections) are kept in full.
func averageRuns(runs []outcomeResult) []outcomeResult {
	var order []string
	byTask := map[string][]outcomeResult{}
	for i := range runs {
		id := runs[i].task
		if byTask[id] == nil {
			order = append(order, id)
		}
		byTask[id] = append(byTask[id], runs[i])
	}
	out := make([]outcomeResult, 0, len(order))
	for _, id := range order {
		rs := byTask[id]
		agg := outcomeResult{task: id, field: rs[0].field}
		votes := map[string]int{}
		for i := range rs {
			v := &rs[i].verdict
			agg.verdict.ARate += v.ARate / float64(len(rs))
			agg.verdict.BRate += v.BRate / float64(len(rs))
			agg.verdict.Disagree += v.Disagree
			agg.verdict.Rejected += v.Rejected
			votes[v.Better]++
		}
		agg.verdict.Rounds = rs[0].verdict.Rounds
		agg.verdict.Better = "tie"
		switch {
		case votes["B"] > votes["A"]:
			agg.verdict.Better = "B"
		case votes["A"] > votes["B"]:
			agg.verdict.Better = "A"
		}
		// Keep every round's guess: one synthetic round per run would lose them.
		agg.guesses = nil
		for i := range rs {
			for _, r := range rs[i].verdict.Rounds {
				agg.guesses = append(agg.guesses, r.Guess)
			}
		}
		out = append(out, agg)
	}
	return out
}

// loadVerdicts reads a kept run's verdict files back into results.
func loadVerdicts(t *testing.T, dir string) []outcomeResult {
	t.Helper()
	fields := map[string]string{}
	tasks := loadOutcomeTasks(t)
	for i := range tasks {
		fields[tasks[i].ID] = tasks[i].Field
	}
	files, _ := filepath.Glob(filepath.Join(dir, "*-verdict.json"))
	sort.Strings(files)
	var out []outcomeResult
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		var v outcomePair
		if err := json.Unmarshal(b, &v); err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		id := strings.TrimSuffix(filepath.Base(f), "-verdict.json")
		if i := strings.LastIndex(id, "-run"); i > 0 {
			id = id[:i]
		}
		out = append(out, outcomeResult{task: id, field: fields[id], verdict: v})
	}
	return out
}

func pct(k, n int) float64 {
	if n == 0 {
		return 0
	}
	return 100 * float64(k) / float64(n)
}

// bootstrapMeanCI returns the percentile bootstrap 95% CI of the mean of xs.
func bootstrapMeanCI(xs []float64, resamples int, seed int64) (lo, hi float64) {
	if len(xs) == 0 {
		return 0, 0
	}
	rng := rand.New(rand.NewSource(seed)) //nolint:gosec // reproducible resampling
	means := make([]float64, resamples)
	for r := range means {
		var sum float64
		for range xs {
			sum += xs[rng.Intn(len(xs))]
		}
		means[r] = sum / float64(len(xs))
	}
	sort.Float64s(means)
	return means[int(0.025*float64(resamples))], means[int(0.975*float64(resamples))-1]
}

func writeOutcomeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Logf("out dir: %v", err)
		return
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Logf("write %s: %v", path, err)
	}
}

// ── offline tests of the harness (no API) ────────────────────────────────

// Every task's fixture builds, its staged change is staged, and its rubric
// is non-empty — before any money is spent.
func TestOutcomeTasks_Fixtures(t *testing.T) {
	tasks := loadOutcomeTasks(t)
	if len(tasks) == 0 {
		t.Fatal("no tasks")
	}
	seen := map[string]bool{}
	for i := range tasks {
		task := &tasks[i]
		if task.ID == "" || seen[task.ID] || len(task.Rubric) < 4 || task.Source == "" || task.Practice == "" {
			t.Errorf("task %d (%s): id, practice, source and 4+ rubric criteria required", i, task.ID)
		}
		seen[task.ID] = true
		dir := t.TempDir()
		if err := seedOutcomeWorkspace(dir, task); err != nil {
			t.Fatalf("%s: seed: %v", task.ID, err)
		}
		if _, ok := task.Files["go.mod"]; ok {
			c := exec.Command("go", "vet", "./...")
			c.Dir = dir
			if out, err := c.CombinedOutput(); err != nil {
				t.Errorf("%s: fixture doesn't build: %v: %s", task.ID, err, out)
			}
		}
		if changes := workspaceChanges(dir); changes != "" {
			t.Errorf("%s: fresh workspace reports changes: %s", task.ID, changes)
		}
	}
}

// The judge never sees grimoire's display, and agent edits are included.
func TestOutcomeOutput_Blinded(t *testing.T) {
	texts := []string{"grimoire ▸ applying a practice — why\n\n**write-sre-runbook**\n⎿  Goal  a runbook\n\nHere is the runbook."}
	got := outcomeOutput(texts, "\n--- file: RUNBOOK.md ---\nbody\n", "")
	if strings.Contains(got, "grimoire") || strings.Contains(got, "**write-sre-runbook**") || strings.Contains(got, "⎿") {
		t.Errorf("not blinded:\n%s", got)
	}
	if !strings.Contains(got, "Goal: a runbook") || !strings.Contains(got, "Here is the runbook.") || !strings.Contains(got, "RUNBOOK.md") {
		t.Errorf("content lost:\n%s", got)
	}
}

// New and modified files are reported; the task's own staged change is not.
func TestWorkspaceChanges(t *testing.T) {
	dir := t.TempDir()
	task := &outcomeTask{Files: map[string]string{"a.go": "package a\n"}, Staged: map[string]string{"b.go": "package a\n"}}
	if err := seedOutcomeWorkspace(dir, task); err != nil {
		t.Fatal(err)
	}
	if got := workspaceChanges(dir); got != "" {
		t.Fatalf("staged task change reported: %s", got)
	}
	_ = os.WriteFile(filepath.Join(dir, "a_test.go"), []byte("package a\n// new\n"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\n// edited\n"), 0o644)
	got := workspaceChanges(dir)
	if !strings.Contains(got, "a_test.go") || !strings.Contains(got, "// edited") || strings.Contains(got, "b.go") {
		t.Errorf("changes = %s", got)
	}
}

// Presentation order is undone, missing criteria count as unmet, and an
// unquotable "met" is rejected.
func TestJudgeMapping(t *testing.T) {
	raw, err := parseJudge("Sure:\n```json\n" + `{"x":[{"criterion":1,"met":true,"evidence":"Hello, guest"}],"y":[{"criterion":2,"met":true,"evidence":"made up words"}],"better":"Y","why":"w"}` + "\n```")
	if err != nil {
		t.Fatal(err)
	}
	v := mapJudge(&raw, false, 2) // Y was A
	if v.Better != "A" || !v.B[0].Met || v.A[0].Met || !v.A[1].Met {
		t.Fatalf("mapping: %+v", v)
	}
	if r := checkEvidence(v.B, "returns **Hello,   guest** for nil"); r != 0 || !v.B[0].Met {
		t.Errorf("real quote rejected: %d %+v", r, v.B)
	}
	if r := checkEvidence(v.A, "something else entirely"); r != 1 || v.A[1].Met {
		t.Errorf("fake quote kept: %d %+v", r, v.A)
	}
	stitched := outcomeCriteria{{Criterion: 1, Met: true, Evidence: `{"2d4h", 3120}, {"90m", 90}`}, {Criterion: 2, Met: true, Evidence: absenceEvidence}, {Criterion: 3, Met: true, Evidence: "invented claim about rollback timing windows"}}
	if r := checkEvidence(stitched, "cases := []struct{}{\n\t{\"90m\", 90},\n\t{\"4h\", 240},\n\t{\"2d4h\", 3120},\n}"); r != 1 || !stitched[0].Met || !stitched[1].Met || stitched[2].Met {
		t.Errorf("stitched/absence/invented: rejected %d %+v", r, stitched)
	}
	if raw, err := parseJudge(`{"x":[],"y":[],"better":"tie","why":"uses ${{ secrets.TOKEN }}"} Commentary after.`); err != nil || raw.Better != "tie" {
		t.Errorf("trailing text / braces: %v %+v", err, raw)
	}
	multi := outcomeCriteria{{Criterion: 1, Met: true, Evidence: `{"2d4h", 3120} ... {"4h", 240}`}, {Criterion: 2, Met: true, Evidence: "real part | invented part"}}
	if r := checkEvidence(multi, `cases := {"2d4h", 3120}, {"90m", 90}, {"4h", 240} and a real part`); r != 1 || !multi[0].Met || multi[1].Met {
		t.Errorf("multi-fragment quotes: rejected %d %+v", r, multi)
	}
	if v := mapJudge(&judgeRaw{Better: "maybe"}, true, 1); v.Better != "tie" || len(v.A) != 1 {
		t.Errorf("unknown better: %+v", v)
	}
}

func TestEndsWithQuestion(t *testing.T) {
	for reply, want := range map[string]bool{
		"Plan…\nApply this plan? Steps that change files or run commands ask first.": true,
		"**Options**\n1. a\nReply 1–2, \"all\", or \"skip\".":                        true,
		"Want me to go deeper?":         true,
		"Here is the runbook.\nDone.":   false,
		"Is it fast? Yes.\nAll done.\n": false,
	} {
		if got := endsWithQuestion(reply); got != want {
			t.Errorf("%q: %v", reply, got)
		}
	}
}

func TestBootstrapMeanCI(t *testing.T) {
	lo, hi := bootstrapMeanCI([]float64{0.1, 0.2, 0.3, 0.2, 0.1, 0.3}, 2000, 1)
	if lo <= 0 || hi >= 0.31 || lo >= hi {
		t.Errorf("CI = [%v, %v]", lo, hi)
	}
	if lo, hi := bootstrapMeanCI([]float64{0.5, 0.5}, 100, 1); lo != 0.5 || hi != 0.5 {
		t.Errorf("constant: [%v, %v]", lo, hi)
	}
}

// Both arms get the same normalization: plan and options scaffolding goes,
// rows become plain lines, markers are dropped.
func TestBlind(t *testing.T) {
	if got := blind("**Runbook written**\n✓ `runbooks/x.md` (not committed yet)"); strings.Contains(got, "✓") || !strings.Contains(got, "- `runbooks/x.md`") {
		t.Errorf("outcome line kept: %q", got)
	}
	in := "**Plan** — 2 steps\n⎿  Step 1  write-adr — record it\n     2a. design-ci-pipeline — move CI ★ recommended\nApply this plan? Steps that change files ask first.\n\n**design-ci-pipeline** — step 2a of 5\n⎿  Approach  run both side by side\nReply 1–2, \"all\", or \"skip\"."
	got := blind(in)
	for _, gone := range []string{"**Plan**", "Apply this plan?", "★", "step 2a of 5", "⎿", "Reply 1–2", "**design-ci-pipeline**"} {
		if strings.Contains(got, gone) {
			t.Errorf("%q survived:\n%s", gone, got)
		}
	}
	if !strings.Contains(got, "Approach: run both side by side") || !strings.Contains(got, "Step 1: write-adr — record it") {
		t.Errorf("content lost:\n%s", got)
	}
	plain := "Here is the plan.\n\n1. Do this\n2. Then that"
	if blind(plain) != plain {
		t.Errorf("plain output changed:\n%s", blind(plain))
	}
}

// Two rounds average per criterion; a preference counts only when both
// orders agree; disagreements are counted.
func TestCombineRounds(t *testing.T) {
	met := func(bs ...bool) outcomeCriteria {
		cs := make(outcomeCriteria, len(bs))
		for i, b := range bs {
			cs[i] = outcomeCriterion{Criterion: i + 1, Met: b}
		}
		return cs
	}
	p := combineRounds(&[2]outcomeVerdict{
		{A: met(true, false), B: met(true, true), Better: "B"},
		{A: met(false, false), B: met(true, true), Better: "B"},
	})
	if p.ARate != 0.25 || p.BRate != 1 || p.Better != "B" || p.Disagree != 1 {
		t.Errorf("agreeing rounds: %+v", p)
	}
	if p := combineRounds(&[2]outcomeVerdict{{A: met(true), B: met(true), Better: "A"}, {A: met(true), B: met(true), Better: "B"}}); p.Better != "tie" {
		t.Errorf("order-dependent preference not a tie: %+v", p)
	}
}

// The plugin guess maps back through the presentation order.
func TestJudgeGuess(t *testing.T) {
	if v := mapJudge(&judgeRaw{Guess: "X"}, true, 1); v.Guess != "A" {
		t.Errorf("X with A first: %s", v.Guess)
	}
	if v := mapJudge(&judgeRaw{Guess: "x"}, false, 1); v.Guess != "B" {
		t.Errorf("X with B first: %s", v.Guess)
	}
	if v := mapJudge(&judgeRaw{}, true, 1); v.Guess != "unsure" {
		t.Errorf("missing guess: %s", v.Guess)
	}
}

// Outputs saved by the first harness version are re-blinded on rescore.
func TestReblind(t *testing.T) {
	got := reblind("-  Change  `Greeting` now handles a nil user\n\n- a normal list item\n- another  item with  spaces")
	if !strings.Contains(got, "Change: `Greeting` now handles a nil user") {
		t.Errorf("legacy row kept:\n%s", got)
	}
	if !strings.Contains(got, "- a normal list item") {
		t.Errorf("list item changed:\n%s", got)
	}
}

// A rewrite must keep code blocks verbatim and every number.
func TestFaithful(t *testing.T) {
	orig := "**Target:** $1,800,000 at 4%.\n\n```yaml\non: [push]\n```\n⎿  Step 2  pay the card first"
	if !faithful(orig, "## Target\n\nThe target is $1800000 at 4%.\n\n```yaml\non: [push]\n```\n\n- Step 2: pay the card first") {
		t.Error("faithful rewrite rejected")
	}
	if faithful(orig, "The target is $1,800,000 at 4%.\n\n```yaml\non: [push, pull_request]\n```") {
		t.Error("changed code block accepted")
	}
	if faithful(orig, "The target is $1,800,000.\n\n```yaml\non: [push]\n```\nStep 2: pay the card first") {
		t.Error("dropped number (4) accepted")
	}
}

// The neutral rewrite never touches the agent's files or execution checks,
// and a cached rewrite is reused.
func TestNeutralize_Cache(t *testing.T) {
	cache := filepath.Join(t.TempDir(), "x.neutral.md")
	if err := os.WriteFile(cache, []byte("cached"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got, cost, kept := neutralize("anything"+filesMarker+"\n--- file: a.go ---", "m", cache); got != "cached" || cost != 0 || !kept {
		t.Errorf("cache not used: %q %v %v", got, cost, kept)
	}
	if got, _, kept := neutralize(filesMarker+"\n--- file: a.go ---\nbody", "m", ""); !kept || !strings.Contains(got, "a.go") {
		t.Errorf("files-only output changed: %q", got)
	}
}

// Repeated runs of a task count once: the CI is over tasks.
func TestAverageRuns(t *testing.T) {
	run := func(task string, a, b float64, better string) outcomeResult {
		return outcomeResult{task: task, field: "f", verdict: outcomePair{ARate: a, BRate: b, Better: better, Rounds: [2]outcomeVerdict{{Guess: "B"}, {Guess: "unsure"}}}}
	}
	got := averageRuns([]outcomeResult{run("t1", 0.5, 1, "B"), run("t1", 0.7, 0.8, "B"), run("t1", 0.6, 0.9, "A"), run("t2", 1, 1, "tie")})
	if len(got) != 2 || got[0].task != "t1" || got[0].verdict.Better != "B" || len(got[0].guesses) != 6 {
		t.Fatalf("averaged: %+v", got)
	}
	if d := got[0].verdict.ARate - 0.6; d > 1e-9 || d < -1e-9 {
		t.Errorf("A rate = %v, want 0.6", got[0].verdict.ARate)
	}
}
