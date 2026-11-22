package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeffreytse/grimoire/internal/skills"
)

// TestHookEval measures the prompt hook against labeled prompts using the
// skills package installed on this machine. It is opt-in because the result
// depends on that package:
//
//	GRIMOIRE_HOOK_EVAL=1 go test ./cmd -run TestHookEval -v
//
// Metrics:
//   - top1:   the first suggestion is an acceptable skill
//   - hit:    some suggestion is an acceptable skill
//   - silent: positive prompt got no suggestion (a miss, but a cheap one)
//   - noise:  prompt that should stay silent got a suggestion
func TestHookEval(t *testing.T) {
	if os.Getenv("GRIMOIRE_HOOK_EVAL") == "" {
		t.Skip("set GRIMOIRE_HOOK_EVAL=1 to run against the installed skills package")
	}
	raw, err := os.ReadFile(filepath.Join("testdata", "hook_eval.json"))
	if err != nil {
		t.Fatal(err)
	}
	var set struct {
		Cases []struct {
			Prompt string   `json:"prompt"`
			Want   []string `json:"want"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(raw, &set); err != nil {
		t.Fatal(err)
	}
	all, _, err := skills.ListAllSkillsFromPackages(skills.AllSkillsPackages())
	if err != nil || len(all) == 0 {
		t.Fatalf("no installed skills: %v", err)
	}
	known := map[string]bool{}
	for i := range all {
		known[all[i].Name] = true
	}

	var pos, neg, top1, hit, silent, noise int
	for _, c := range set.Cases {
		for _, w := range c.Want {
			if !known[w] {
				t.Errorf("label %q (prompt %q) is not an installed skill", w, c.Prompt)
			}
		}
		picks := hookPicks(c.Prompt, all)
		names := make([]string, len(picks))
		for i := range picks {
			names[i] = picks[i].Skill.Name
		}
		got := strings.Join(names, ", ")

		if len(c.Want) == 0 {
			neg++
			if len(picks) > 0 {
				noise++
				t.Logf("NOISE  %-60q -> %s", c.Prompt, got)
			}
			continue
		}
		pos++
		want := map[string]bool{}
		for _, w := range c.Want {
			want[w] = true
		}
		switch {
		case len(picks) == 0:
			silent++
			t.Logf("SILENT %-60q", c.Prompt)
		case want[names[0]]:
			top1++
			hit++
		default:
			found := false
			for _, n := range names {
				found = found || want[n]
			}
			if found {
				hit++
				t.Logf("HIT@N  %-60q -> %s", c.Prompt, got)
			} else {
				t.Logf("WRONG  %-60q -> %s", c.Prompt, got)
			}
		}
	}
	pct := func(n, d int) float64 { return 100 * float64(n) / float64(max(d, 1)) }
	t.Logf("positives=%d top1=%.0f%% hit=%.0f%% silent=%.0f%% wrong=%.0f%% | negatives=%d noise=%.0f%%",
		pos, pct(top1, pos), pct(hit, pos), pct(silent, pos), pct(pos-hit-silent, pos), neg, pct(noise, neg))
}
