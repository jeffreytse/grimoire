package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestWorkflowEval checks that grimoire's meta workflows are reached the way
// start-best-practice's routing table says — the behavior skills mode gets
// from installed skills. Opt-in; runs Claude Code headless (API credits):
//
//	GRIMOIRE_WORKFLOW_EVAL=1 go test ./cmd -run TestWorkflowEval -v
//
// It honors the display eval's GRIMOIRE_DISPLAY_EVAL_MODE, _MODEL, _OUT and
// GRIMOIRE_SOURCE. A case passes when one of its workflows was loaded
// (grimoire_get_skill / Skill) or named on the route line.
func TestWorkflowEval(t *testing.T) {
	if os.Getenv("GRIMOIRE_WORKFLOW_EVAL") == "" {
		t.Skip("set GRIMOIRE_WORKFLOW_EVAL=1 to run against Claude Code (uses API credits)")
	}
	if _, err := exec.LookPath("claude"); err != nil {
		t.Skip("claude CLI not on PATH")
	}
	raw, err := os.ReadFile(filepath.Join("testdata", "workflow_eval.json"))
	if err != nil {
		t.Fatal(err)
	}
	var set struct {
		Cases []struct {
			Prompt   string   `json:"prompt"`
			Workflow []string `json:"workflow"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(raw, &set); err != nil {
		t.Fatal(err)
	}
	if only := os.Getenv("GRIMOIRE_DISPLAY_EVAL_ONLY"); only != "" {
		cases := set.Cases[:0:0]
		for _, f := range strings.Split(only, ",") {
			var n int
			if _, err := fmt.Sscanf(strings.TrimSpace(f), "%d", &n); err == nil && n >= 1 && n <= len(set.Cases) {
				cases = append(cases, set.Cases[n-1])
			}
		}
		set.Cases = cases
	}
	mode := os.Getenv("GRIMOIRE_DISPLAY_EVAL_MODE")
	if mode == "" {
		mode = "mcp"
	}
	dir := t.TempDir()
	bin := filepath.Join(dir, "grimoire")
	if out, err := exec.Command("go", "build", "-o", bin, "..").CombinedOutput(); err != nil {
		t.Fatalf("build: %v: %s", err, out)
	}
	run, err := evalRunner(t, dir, bin, mode)
	if err != nil {
		t.Fatal(err)
	}
	pass, cost := 0, 0.0
	for i, c := range set.Cases {
		r, err := run(c.Prompt, "brief", "")
		cost += r.cost
		if err != nil {
			t.Logf("ERROR %q: %v", c.Prompt, err)
			continue
		}
		reached := reachedWorkflows(strings.Join(r.texts, "\n"), r.loaded)
		ok := false
		for _, w := range c.Workflow {
			ok = ok || reached[w]
		}
		if ok {
			pass++
		}
		var got []string
		for w := range reached {
			got = append(got, w)
		}
		t.Logf("%-4s #%d %q → %s (want one of %s)", map[bool]string{true: "PASS", false: "FAIL"}[ok], i+1, c.Prompt,
			strings.Join(got, ","), strings.Join(c.Workflow, ","))
		if dir := os.Getenv("GRIMOIRE_DISPLAY_EVAL_OUT"); dir != "" {
			_ = os.MkdirAll(dir, 0o755)
			_ = os.WriteFile(filepath.Join(dir, fmt.Sprintf("workflow-%s-%02d.md", mode, i+1)),
				[]byte(c.Prompt+"\n\n"+strings.Join(r.texts, "\n\n---\n\n")), 0o644)
		}
	}
	t.Logf("%s workflows reached %d/%d  cost=$%.2f", mode, pass, len(set.Cases), cost)
}

// routeLine matches the route line's target: start-best-practice ▸ <workflow>.
var routeLine = regexp.MustCompile(`start-best-practice\**\s*(?:▸|→)\s*\**([a-z0-9]+(?:-[a-z0-9]+)+)`)

// reachedWorkflows is the set of skills a run loaded or routed to.
func reachedWorkflows(reply string, loaded []string) map[string]bool {
	got := map[string]bool{}
	for _, n := range loaded {
		got[n] = true
	}
	for _, m := range routeLine.FindAllStringSubmatch(reply, -1) {
		got[m[1]] = true
	}
	return got
}

func TestReachedWorkflows(t *testing.T) {
	got := reachedWorkflows("start-best-practice ▸ discover-best-practices  `meta/routing` — x", []string{"start-best-practice"})
	if !got["discover-best-practices"] || !got["start-best-practice"] || len(got) != 2 {
		t.Errorf("got %v", got)
	}
	if got := reachedWorkflows("**start-best-practice** → **review-best-practice-fit**  `meta/routing`", nil); !got["review-best-practice-fit"] {
		t.Errorf("old-style route line: %v", got)
	}
}
