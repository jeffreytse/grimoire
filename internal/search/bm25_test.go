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
