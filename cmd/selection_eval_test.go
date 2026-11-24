package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/jeffreytse/grimoire/internal/skills"
)

// TestSelectionEval measures whether grimoire leads the model to the right
// practice — the rate users feel — against the labeled prompts in
// testdata/hook_eval.json (the same labels as the keyword-candidate eval).
// Opt-in; it runs Claude Code headless and costs API usage:
//
//	GRIMOIRE_SELECTION_EVAL=1 go test ./cmd -run TestSelectionEval -v -timeout 120m
//	GRIMOIRE_DISPLAY_EVAL_MODE=mcp|skills   grimoire mode (default mcp; not both)
//	GRIMOIRE_DISPLAY_EVAL_MODEL=<m>         pin the model when comparing runs
//	GRIMOIRE_DISPLAY_EVAL_LIMIT=N           first N cases only
//	GRIMOIRE_DISPLAY_EVAL_ONLY=2,5,39       only these cases (1-based)
//	GRIMOIRE_DISPLAY_EVAL_OUT=<dir>         keep transcripts and picks
//	GRIMOIRE_SOURCE=<package>               evaluate a draft skills package
//	GRIMOIRE_SELECTION_EVAL_SET=<file>      cases file under testdata (default
//	                                        hook_eval.json; selection_hard.json
//	                                        holds the harder cases)
//
// A case's picks are, in order: the practice the reply recommends with ★
// (suggest-best-practice often recommends and waits rather than loading),
// the skills shown as display blocks (a block names the skill it applies,
// even when the model had its content without a fresh load), then the
// non-workflow skills the run loaded. Workflow skills
// (start-best-practice, suggest-best-practice, …) are routing, not picks.
//
// Rates, with 95% Wilson intervals:
//   - routing recall / precision / accuracy: engaged on labeled prompts, and
//     only on them (engaged = loaded any grimoire skill or recommended one)
//   - top-1 accuracy: the first pick is a labeled skill
//   - hit rate: some pick is a labeled skill
//   - pick precision: share of all picks that are labeled skills
//   - clarify rate: routed replies that ask the user one question before
//     choosing — what the workflows are meant to do when a request is
//     unclear. They have no pick yet, so they are left out of top-1 and hit
//     rate and reported on their own; a rising clarify rate is a signal too.
//
// Every pick missing from its labels is listed: review each, then either
// accept it (add it to "want") or count it as a real miss. Labels written
// for keyword search are often narrower than the set of good practices.
func TestSelectionEval(t *testing.T) {
	if dir := os.Getenv("GRIMOIRE_SELECTION_RESCORE"); dir != "" {
		rescoreSelection(t, dir)
		return
	}
	if os.Getenv("GRIMOIRE_SELECTION_EVAL") == "" {
		t.Skip("set GRIMOIRE_SELECTION_EVAL=1 to run against Claude Code (uses API credits)")
	}
	if _, err := exec.LookPath("claude"); err != nil {
		t.Skip("claude CLI not on PATH")
	}
	file := "hook_eval.json"
	if f := os.Getenv("GRIMOIRE_SELECTION_EVAL_SET"); f != "" {
		file = f
	}
	raw, err := os.ReadFile(filepath.Join("testdata", file))
	if err != nil {
		t.Fatal(err)
	}
	var set struct {
		Cases []selectionCase `json:"cases"`
	}
	if err := json.Unmarshal(raw, &set); err != nil {
		t.Fatal(err)
	}
	if n, err := strconv.Atoi(os.Getenv("GRIMOIRE_DISPLAY_EVAL_LIMIT")); err == nil && n > 0 && n < len(set.Cases) {
		set.Cases = set.Cases[:n]
	}
	if only := os.Getenv("GRIMOIRE_DISPLAY_EVAL_ONLY"); only != "" {
		var picked []selectionCase
		for _, f := range strings.Split(only, ",") {
			if n, err := strconv.Atoi(strings.TrimSpace(f)); err == nil && n >= 1 && n <= len(set.Cases) {
				picked = append(picked, set.Cases[n-1])
			}
		}
		set.Cases = picked
	}
	mode := os.Getenv("GRIMOIRE_DISPLAY_EVAL_MODE")
	if mode == "" {
		mode = "mcp"
	}
	if mode != "mcp" && mode != "skills" {
		t.Fatalf("GRIMOIRE_DISPLAY_EVAL_MODE=%q: want mcp or skills", mode)
	}
	if mode == "skills" && os.Getenv("GRIMOIRE_DISPLAY_EVAL_SUBDOMAINS") == "" {
		// Link every subdomain a label lives in, so skills mode can pick it.
		subs, err := labeledSubdomains(set.Cases)
		if err != nil {
			t.Fatal(err)
		}
		t.Setenv("GRIMOIRE_DISPLAY_EVAL_SUBDOMAINS", strings.Join(subs, ","))
	}

	dir := t.TempDir()
	bin := filepath.Join(dir, "grimoire")
	if out, err := exec.Command("go", "build", "-o", bin, "..").CombinedOutput(); err != nil {
		t.Fatalf("build: %v: %s", err, out)
	}
	run, err := evalRunner(t, dir, bin, mode)
	if err != nil {
		t.Fatalf("%s setup: %v", mode, err)
	}

	var results []selectionResult
	var cost float64
	for i, c := range set.Cases {
		r, err := run(c.Prompt, "brief", "")
		cost += r.cost
		if err != nil {
			t.Logf("ERROR #%d %q: %v", i+1, c.Prompt, err)
			continue
		}
		res := scoreSelection(c, strings.Join(r.texts, "\n"), r.loaded)
		results = append(results, res)
		saveSelection(t, mode, i+1, c, res, r.texts)
	}
	logSelection(t, mode, results)
	t.Logf("%s cases=%d cost=$%.2f", mode, len(results), cost)
}

type selectionCase struct {
	Prompt string   `json:"prompt"`
	Want   []string `json:"want"`
	// Accepted lists picks judged appropriate after a run, beyond the
	// original labels (each with its reason in the case's note). Strict
	// rates use Want only; adjudicated rates count Accepted too.
	Accepted []string `json:"accepted"`
}

type selectionResult struct {
	Case      selectionCase
	Engaged   bool     // loaded or recommended any grimoire skill
	Picks     []string // ★ recommendation, then loaded non-workflow skills
	Clarified bool     // routed, no pick yet, and asked the user a question
}

// routingSkills are grimoire's routing workflows: they find a practice,
// they are not one. Other meta skills (install-grimoire, configure-grimoire,
// explain-best-practice …) are real picks when the request is about them.
var routingSkills = map[string]bool{
	"start-best-practice": true, "suggest-best-practice": true,
	"plan-best-practice-solution": true, "apply-best-practice-tree": true,
	"analyze-best-practice-problem": true, "review-best-practice-fit": true,
	"discover-best-practices": true,
}

// workflowSkill reports whether name is a routing workflow, not a pick.
func workflowSkill(name string) bool {
	return routingSkills[name]
}

// starPick matches the practice a reply recommends with ★.
var starPick = regexp.MustCompile(`★\s*(?:\*\*)?\s*([a-z0-9]+(?:-[a-z0-9]+)+)`)

// stepPick matches a plan step row: ⎿  Step N  skill-name  `domain`.
var stepPick = regexp.MustCompile("(?m)^⎿\\s+Step\\s+\\d+\\s+\\**([a-z0-9]+(?:-[a-z0-9]+)+)\\**\\s+`[a-z0-9-]+/[a-z0-9-]+`")

// listPick matches a ranked option: "2. **skill-name**" (with or without ★).
var listPick = regexp.MustCompile(`(?m)^\s*\d+\.\s*(?:★\s*)?\*\*([a-z0-9]+(?:-[a-z0-9]+)+)\*\*`)

// matchPick matches a skill named in a Match/Pick/Recommended row.
var matchPick = regexp.MustCompile(`(?m)^⎿\s+(?:Match|Pick|Recommended|Recommend)\s+\**([a-z0-9]+(?:-[a-z0-9]+)+)`)

// askedQuestion reports whether a reply ends by asking the user something.
func askedQuestion(reply string) bool {
	lines := strings.Split(strings.TrimSpace(reply), "\n")
	for i := len(lines) - 1; i >= 0 && i >= len(lines)-12; i-- {
		if strings.Contains(lines[i], "?") {
			return true
		}
	}
	return false
}

// blockPick matches a display block header: **skill-name** then its domain.
var blockPick = regexp.MustCompile("(?m)^\\*\\*([a-z0-9]+(?:-[a-z0-9]+)+)\\*\\*[^\\n]*`[a-z0-9-]+/[a-z0-9-]+`")

// scoreSelection derives a case's picks from its reply and loaded skills.
func scoreSelection(c selectionCase, reply string, loaded []string) selectionResult {
	res := selectionResult{Case: c, Engaged: len(loaded) > 0 || displayMarker.MatchString(reply)}
	seen := map[string]bool{}
	add := func(n string) {
		if n != "" && !seen[n] && !workflowSkill(n) {
			seen[n] = true
			res.Picks = append(res.Picks, n)
		}
	}
	if m := starPick.FindStringSubmatch(reply); m != nil {
		add(m[1])
		res.Engaged = true
	}
	for _, re := range []*regexp.Regexp{matchPick, stepPick, listPick} {
		for _, m := range re.FindAllStringSubmatch(reply, -1) {
			add(m[1])
		}
	}
	for _, m := range blockPick.FindAllStringSubmatch(reply, -1) {
		add(m[1])
	}
	for _, n := range loaded {
		add(n)
	}
	res.Clarified = res.Engaged && len(res.Picks) == 0 && askedQuestion(reply)
	return res
}

// selectionRates are the eval's rates as k/n counts.
type selectionRates struct {
	recallK, recallN       int // engaged on labeled prompts
	precisionK, precisionN int // labeled among engaged
	accuracyK, accuracyN   int // right engage/stay-out decision
	top1K, top1N           int
	hitK, hitN             int
	pickK, pickN           int // labeled picks among all picks
	clarifyK, clarifyN     int // asked a question first, among engaged labeled
}

func selectionRatesOf(results []selectionResult) selectionRates {
	var r selectionRates
	for _, res := range results {
		labeled := len(res.Case.Want) > 0
		want := map[string]bool{}
		for _, w := range res.Case.Want {
			want[w] = true
		}
		r.accuracyN++
		if res.Engaged == labeled {
			r.accuracyK++
		}
		if res.Engaged {
			r.precisionN++
			if labeled {
				r.precisionK++
			}
		}
		if !labeled {
			continue
		}
		r.recallN++
		if res.Engaged {
			r.recallK++
			r.clarifyN++
		}
		if res.Clarified {
			r.clarifyK++
			continue // no pick yet: not a miss, not a hit
		}
		r.top1N++
		r.hitN++
		hit := false
		for i, p := range res.Picks {
			r.pickN++
			if want[p] {
				r.pickK++
				hit = true
				if i == 0 {
					r.top1K++
				}
			}
		}
		if hit {
			r.hitK++
		}
	}
	return r
}

func logSelection(t *testing.T, mode string, results []selectionResult) {
	strict := selectionRatesOf(results)
	adjudicated := make([]selectionResult, len(results))
	for i, res := range results {
		res.Case.Want = append(append([]string{}, res.Case.Want...), res.Case.Accepted...)
		adjudicated[i] = res
	}
	adj := selectionRatesOf(adjudicated)
	for _, row := range []struct {
		name string
		k, n int
	}{
		{"top-1 adjudicated", adj.top1K, adj.top1N},
		{"hit adjudicated", adj.hitK, adj.hitN},
		{"pick prec. adjud.", adj.pickK, adj.pickN},
	} {
		if row.n > 0 {
			lo, hi := wilson(row.k, row.n)
			t.Logf("%-6s %-18s %3d/%-3d %5.1f%%  95%% CI %5.1f–%5.1f%%", mode, row.name, row.k, row.n,
				100*float64(row.k)/float64(row.n), 100*lo, 100*hi)
		}
	}
	for _, res := range adjudicated {
		if len(res.Case.Want) == 0 || len(res.Picks) == 0 {
			continue
		}
		first := res.Picks[0]
		ok := false
		for _, w := range res.Case.Want {
			ok = ok || w == first
		}
		if !ok {
			t.Logf("TOP1-MISS %q → %s", res.Case.Prompt, first)
		}
	}
	r := strict
	for _, row := range []struct {
		name string
		k, n int
	}{
		{"routing recall", r.recallK, r.recallN},
		{"routing precision", r.precisionK, r.precisionN},
		{"routing accuracy", r.accuracyK, r.accuracyN},
		{"top-1 accuracy", r.top1K, r.top1N},
		{"hit rate", r.hitK, r.hitN},
		{"pick precision", r.pickK, r.pickN},
		{"clarify rate", r.clarifyK, r.clarifyN},
	} {
		if row.n == 0 {
			continue
		}
		lo, hi := wilson(row.k, row.n)
		t.Logf("%-6s %-18s %3d/%-3d %5.1f%%  95%% CI %5.1f–%5.1f%%", mode, row.name, row.k, row.n,
			100*float64(row.k)/float64(row.n), 100*lo, 100*hi)
	}
	// Picks outside the labels, for review.
	var review []string
	for _, res := range results {
		want := map[string]bool{}
		for _, w := range res.Case.Want {
			want[w] = true
		}
		for _, p := range res.Picks {
			if !want[p] {
				review = append(review, fmt.Sprintf("%q → %s (labels: %s)", res.Case.Prompt, p, strings.Join(res.Case.Want, ", ")))
			}
		}
	}
	for _, line := range review {
		t.Logf("REVIEW %s", line)
	}
}

// wilson is the 95% Wilson score interval for k successes in n trials.
func wilson(k, n int) (lo, hi float64) {
	if n == 0 {
		return 0, 0
	}
	const z = 1.96
	p, nf := float64(k)/float64(n), float64(n)
	d := 1 + z*z/nf
	c := p + z*z/(2*nf)
	r := z * math.Sqrt(p*(1-p)/nf+z*z/(4*nf*nf))
	return (c - r) / d, (c + r) / d
}

// labeledSubdomains lists domain/subdomain for every labeled skill.
func labeledSubdomains(cases []selectionCase) ([]string, error) {
	all, _, err := skills.ListAllSkillsFromPackages(skills.AllSkillsPackages())
	if err != nil {
		return nil, err
	}
	where := map[string]string{}
	for i := range all {
		where[all[i].Name] = all[i].Domain + "/" + all[i].Subdomain
	}
	set := map[string]bool{}
	for _, c := range cases {
		for _, w := range c.Want {
			if s, ok := where[w]; ok {
				set[s] = true
			}
		}
	}
	subs := make([]string, 0, len(set))
	for s := range set {
		subs = append(subs, s)
	}
	sort.Strings(subs)
	return subs, nil
}

// saveSelection keeps a case's picks and reply under GRIMOIRE_DISPLAY_EVAL_OUT.
func saveSelection(t *testing.T, mode string, n int, c selectionCase, res selectionResult, texts []string) {
	dir := os.Getenv("GRIMOIRE_DISPLAY_EVAL_OUT")
	if dir == "" {
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}
	body := fmt.Sprintf("mode: %s\nprompt: %s\nwant: %s\nengaged: %v\npicks: %s\n\n%s\n",
		mode, c.Prompt, strings.Join(c.Want, ", "), res.Engaged, strings.Join(res.Picks, ", "),
		strings.Join(texts, "\n\n---\n\n"))
	if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("selection-%s-%02d.md", mode, n)), []byte(body), 0o644); err != nil {
		t.Logf("transcript: %v", err)
	}
}

// The selection scoring is deterministic and tested offline.
func TestScoreSelection(t *testing.T) {
	launch := selectionCase{Prompt: "launch", Want: []string{"apply-premortem", "design-deployment-strategy"}}
	none := selectionCase{Prompt: "capital of Australia?"}

	r := scoreSelection(launch, "1. ★ **apply-premortem** (recommended)\n2. **run-game-day**",
		[]string{"start-best-practice", "suggest-best-practice", "run-game-day"})
	if !r.Engaged || strings.Join(r.Picks, ",") != "apply-premortem,run-game-day" {
		t.Errorf("picks = %v engaged=%v", r.Picks, r.Engaged)
	}
	if r := scoreSelection(none, "Canberra.", nil); r.Engaged || len(r.Picks) != 0 {
		t.Errorf("silent case: %+v", r)
	}
	if r := scoreSelection(launch, "plain answer", []string{"start-best-practice"}); !r.Engaged || len(r.Picks) != 0 {
		t.Errorf("routed without a pick: %+v", r)
	}
	// A block applied without a fresh load still counts, after the ★ pick;
	// workflow blocks and the route line engage but are not picks.
	block := "start-best-practice ▸ suggest-best-practice  `meta/routing` — x\n\n**suggest-best-practice**  `meta/discovery`\n\n**propose-conventional-commit**  `engineering/development`\n⎿  Type  docs"
	commit := selectionCase{Prompt: "commit", Want: []string{"propose-conventional-commit"}}
	if r := scoreSelection(commit, block, nil); !r.Engaged || strings.Join(r.Picks, ",") != "propose-conventional-commit" {
		t.Errorf("block pick: %+v", r)
	}

	rates := selectionRatesOf([]selectionResult{
		scoreSelection(launch, "★ apply-premortem", nil),                      // top-1 right
		scoreSelection(launch, "★ run-game-day", []string{"apply-premortem"}), // hit, not top-1
		scoreSelection(launch, "answer", nil),                                 // missed
		scoreSelection(none, "★ apply-premortem", nil),                        // false engage
		scoreSelection(none, "Canberra.", nil),                                // right silence
	})
	want := selectionRates{
		recallK: 2, recallN: 3, precisionK: 2, precisionN: 3, accuracyK: 3, accuracyN: 5,
		top1K: 1, top1N: 3, hitK: 2, hitN: 3, pickK: 2, pickN: 3, clarifyN: 2,
	}
	if rates != want {
		t.Errorf("rates = %+v\nwant    %+v", rates, want)
	}
	asked := scoreSelection(launch, "start-best-practice ▸ suggest-best-practice  `meta/routing` — x\n\nWhat is weighing on you most?", nil)
	if !asked.Clarified || len(asked.Picks) != 0 {
		t.Errorf("clarifying reply: %+v", asked)
	}
	matched := scoreSelection(launch, "**suggest-best-practice**  `meta/discovery`\n⎿  Match  apply-premortem  `business/strategy`", nil)
	if strings.Join(matched.Picks, ",") != "apply-premortem" {
		t.Errorf("match row pick: %+v", matched)
	}
	if r := selectionRatesOf([]selectionResult{asked}); r.clarifyK != 1 || r.top1N != 0 || r.recallK != 1 {
		t.Errorf("clarified rates: %+v", r)
	}
	plan := "**plan-best-practice-solution**  `meta/planning`\n⎿  Step 1  design-slo  `engineering/reliability`: x\n⎿  Step 2  write-sre-runbook  `engineering/reliability`: y"
	if r := scoreSelection(launch, plan, nil); strings.Join(r.Picks, ",") != "design-slo,write-sre-runbook" {
		t.Errorf("plan step picks: %+v", r.Picks)
	}
	ranked := "1. ★ **apply-five-whys**: a\n2. **plan-incident-response**: b"
	if r := scoreSelection(launch, ranked, nil); strings.Join(r.Picks, ",") != "apply-five-whys,plan-incident-response" {
		t.Errorf("ranked picks: %+v", r.Picks)
	}
	if lo, hi := wilson(11, 12); lo < 0.64 || lo > 0.65 || hi < 0.98 || hi > 0.99 {
		t.Errorf("wilson(11,12) = %.3f–%.3f", lo, hi)
	}
}

func TestLoadedSkills(t *testing.T) {
	for _, c := range []struct {
		tool, input, want string
	}{
		{mcpGetSkillTool, `{"name":"grimoire:apply-premortem"}`, "apply-premortem"},
		{mcpGetSkillTool, `{"names":["design-slo","run-game-day"]}`, "design-slo,run-game-day"},
		{"Skill", `{"skill":"plugin:write-unit-test","args":"x"}`, "write-unit-test"},
		{"Read", `{"file_path":"x"}`, ""},
	} {
		if got := strings.Join(loadedSkills(c.tool, json.RawMessage(c.input)), ","); got != c.want {
			t.Errorf("%s %s: %q, want %q", c.tool, c.input, got, c.want)
		}
	}
}

// rescoreSelection re-applies the current scorer and labels to transcripts a
// run saved (GRIMOIRE_DISPLAY_EVAL_OUT), free of API cost:
//
//	GRIMOIRE_SELECTION_RESCORE=<dir> [GRIMOIRE_SELECTION_EVAL_SET=…] go test ./cmd -run TestSelectionEval -v
//
// Loaded skills are not in the transcript; the picks the run recorded stand
// in for them.
func rescoreSelection(t *testing.T, dir string) {
	file := "hook_eval.json"
	if f := os.Getenv("GRIMOIRE_SELECTION_EVAL_SET"); f != "" {
		file = f
	}
	raw, err := os.ReadFile(filepath.Join("testdata", file))
	if err != nil {
		t.Fatal(err)
	}
	var set struct {
		Cases []selectionCase `json:"cases"`
	}
	if err := json.Unmarshal(raw, &set); err != nil {
		t.Fatal(err)
	}
	labels := map[string]selectionCase{}
	for _, c := range set.Cases {
		labels[c.Prompt] = c
	}
	files, _ := filepath.Glob(filepath.Join(dir, "selection-*.md"))
	var results []selectionResult
	mode := "rescore"
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		head, body, _ := strings.Cut(string(b), "\n\n")
		fields := map[string]string{}
		for _, line := range strings.Split(head, "\n") {
			if k, v, ok := strings.Cut(line, ": "); ok {
				fields[k] = v
			}
		}
		c, ok := labels[fields["prompt"]]
		if !ok {
			continue
		}
		mode = fields["mode"]
		var loaded []string
		for _, p := range strings.Split(fields["picks"], ", ") {
			if p = strings.TrimSpace(p); p != "" {
				loaded = append(loaded, p)
			}
		}
		if fields["engaged"] == "true" && len(loaded) == 0 {
			loaded = []string{"start-best-practice"} // engaged without a recorded pick
		}
		results = append(results, scoreSelection(c, body, loaded))
	}
	logSelection(t, mode, results)
	t.Logf("%s rescored=%d from %s", mode, len(results), dir)
}
