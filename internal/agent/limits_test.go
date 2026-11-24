package agent

import "testing"

func TestSkillListingBudget(t *testing.T) {
	for ag, want := range map[string]int{"claude": 8000, "codex": 8000, "gemini": 0, "openclaw": 0} {
		if got := SkillListingBudget(ag); got != want {
			t.Errorf("%s: %d, want %d", ag, got, want)
		}
	}
	if c := SkillListingCost("abc", "12345"); c != 16 {
		t.Errorf("cost = %d, want 16", c)
	}
}
