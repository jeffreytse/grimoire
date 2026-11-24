package search

import (
	"strings"
	"testing"

	"github.com/jeffreytse/grimoire/internal/skills"
)

func TestTokenize(t *testing.T) {
	cases := map[string]string{
		"propose-conventional-commit":                   "propose conventional commit",
		"how does Claude trigger grimoire's workflows?": "claude trigger grimoire workflow",
		"please help me write tests":                    "write test",
		"renovating vs renovation":                      "renov vs renov",
		"class status analysis policies":                "class status analysis policy",
		"a I ok thanks":                                 "",
	}
	for in, want := range cases {
		if got := strings.Join(tokenize(in), " "); got != want {
			t.Errorf("tokenize(%q) = %q; want %q", in, got, want)
		}
	}
}

func TestBM25Search_MatchedTerms(t *testing.T) {
	corpus := []skills.Skill{
		{Name: "propose-conventional-commit", Description: "Draft a commit message"},
		{Name: "design-morning-routine", Description: "Plan a morning routine"},
	}
	hits := BM25Search("write a commit message please", corpus, 0)
	if len(hits) != 1 || hits[0].Skill.Name != "propose-conventional-commit" {
		t.Fatalf("hits = %+v", hits)
	}
	if hits[0].Matched != 2 || hits[0].Terms != 3 {
		t.Errorf("Matched/Terms = %d/%d; want 2/3 (write, commit, message)", hits[0].Matched, hits[0].Terms)
	}
	if got := BM25Search("thanks, ok", corpus, 0); got != nil {
		t.Errorf("stopword-only query should return nil, got %+v", got)
	}
}

// A term in a skill's name outranks the same term in passing in another
// skill's description (BM25F field weights).
func TestBM25Search_NameOutweighsDescription(t *testing.T) {
	corpus := []skills.Skill{
		{Name: "plan-race-meals", Description: "load test"},
		{Name: "run-api-load-test-suite", Description: "Use when checking capacity under real production traffic"},
		{Name: "write-readme", Description: "Use when writing a README"},
	}
	hits := BM25Search("load test", corpus, 0)
	if len(hits) != 2 || hits[0].Skill.Name != "run-api-load-test-suite" || hits[0].Score <= hits[1].Score {
		t.Fatalf("hits = %+v", hits)
	}
}

// Bare numbers are quantities, not topics; mixed tokens like "oauth2" stay.
func TestTokenize_DropsBareNumbers(t *testing.T) {
	if got := strings.Join(tokenize("launch in 48 hours with oauth2 and s3"), " "); got != "launch hour oauth2 s3" {
		t.Errorf("tokenize = %q", got)
	}
}

// A skill whose problem vocabulary is only in its body is found by the
// hybrid search, while the field match stays on top.
func TestHybridSearch_BodyRecallFieldPrecision(t *testing.T) {
	corpus := []skills.Skill{
		{Name: "design-deployment-strategy", Description: "Use when releasing software", Body: "Plan the go-live: rollback, canary, blue-green."},
		{Name: "plan-go-live-party", Description: "Use when planning a go-live celebration", Tags: []string{"go-live"}},
		{Name: "write-readme", Description: "Use when writing a README"},
	}
	hits := HybridSearch([]string{"go-live rollback"}, corpus, 0, Options{})
	if len(hits) < 2 {
		t.Fatalf("hits = %+v", hits)
	}
	found := map[string]bool{}
	for _, h := range hits {
		found[h.Skill.Name] = true
	}
	if !found["design-deployment-strategy"] {
		t.Error("body-only match not found")
	}
	if len(BM25Search("rollback", corpus, 0)) != 0 {
		t.Error("field-only search should not see body words")
	}
}

// Several phrasings of one need fuse into one ranking; filters apply.
func TestHybridSearch_MultiQueryAndFilters(t *testing.T) {
	corpus := []skills.Skill{
		{Name: "diagnose-sleep-issue", Domain: "health", Subdomain: "sleep", Tags: []string{"insomnia"}, Description: "Use when you cannot sleep"},
		{Name: "design-sleep-hygiene-protocol", Domain: "health", Subdomain: "sleep", Description: "Use when building better sleep habits"},
		{Name: "apply-sleep-mode-power", Domain: "engineering", Subdomain: "embedded", Description: "Use when putting a device in sleep mode"},
	}
	hits := HybridSearch([]string{"can't fall asleep", "insomnia"}, corpus, 0, Options{Domain: "health"})
	if len(hits) == 0 || hits[0].Skill.Name != "diagnose-sleep-issue" {
		t.Fatalf("hits = %+v", hits)
	}
	for _, h := range hits {
		if h.Skill.Domain != "health" {
			t.Errorf("filter leaked %s", h.Skill.Name)
		}
	}
	if got := HybridSearch([]string{"sleep"}, corpus, 0, Options{Tags: []string{"insomnia"}}); len(got) != 1 {
		t.Errorf("tag filter: %d hits", len(got))
	}
}

// "one" is filler in chat but names topics like one-on-ones, so it is kept
// ("down" was tried and reverted: it lowered hybrid top-1 and MRR); short
// slash acronyms stay one token.
func TestTokenize_MeaningfulWordsAndAcronyms(t *testing.T) {
	for in, want := range map[string]string{
		"prepare for a one on one": "prepare one one",
		"run an A/B test":          "run ab test",
		"set up CI/CD and I/O":     "set up cicd io",
		"and/or either way":        "andor either",
	} {
		if got := strings.Join(tokenize(in), " "); got != want {
			t.Errorf("tokenize(%q) = %q; want %q", in, got, want)
		}
	}
}
