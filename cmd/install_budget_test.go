package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeffreytse/grimoire/internal/skills"
)

// fakeSkills returns n skills whose listing entries each cost 100 characters.
func fakeSkills(n int) []skills.Skill {
	out := make([]skills.Skill, n)
	for i := range out {
		out[i] = skills.Skill{Name: "skill-name", Description: strings.Repeat("d", 82)}
	}
	return out
}

// The listing shows the skills that fit the agent's budget, in install order.
func TestListingFit(t *testing.T) {
	visible, need, budget := listingFit("claude", fakeSkills(200))
	if budget != 8000 || need != 20000 || visible != 80 {
		t.Errorf("visible=%d need=%d budget=%d", visible, need, budget)
	}
	if v, _, b := listingFit("gemini", fakeSkills(200)); b != 0 || v != 0 {
		t.Errorf("unknown budget: visible=%d budget=%d", v, b)
	}
}

// Over budget, an MCP-capable agent is pointed to MCP mode; a skills-only
// agent only to a subset; under budget or unknown, nothing is said.
func TestListingWarning(t *testing.T) {
	msg := listingWarning("claude", 80, 200, 20000, 8000)
	for _, want := range []string{"Claude Code will see about 80 of 200 skills", "~8,000", "~20,000", "grimoire mcp install --target claude", "--domain"} {
		if !strings.Contains(msg, want) {
			t.Errorf("claude warning lacks %q:\n%s", want, msg)
		}
	}
	if msg := listingWarning("openclaw", 80, 200, 20000, 8000); strings.Contains(msg, "mcp install") || !strings.Contains(msg, "--domain") {
		t.Errorf("skills-only agent warning:\n%s", msg)
	}
	if listingWarning("claude", 50, 50, 5000, 8000) != "" || listingWarning("gemini", 0, 200, 20000, 0) != "" {
		t.Error("warned under budget or with an unknown budget")
	}
}

// --no-load-optimize silences the warning.
func TestWarnListingBudgets_NoCap(t *testing.T) {
	var out bytes.Buffer
	warnListingBudgets(&out, []string{"claude"}, fakeSkills(200))
	if !strings.Contains(out.String(), "will see about") {
		t.Fatalf("no warning: %q", out.String())
	}
	orig := flagInstallNoCap
	t.Cleanup(func() { flagInstallNoCap = orig })
	flagInstallNoCap = true
	out.Reset()
	warnListingBudgets(&out, []string{"claude"}, fakeSkills(200))
	if out.Len() != 0 {
		t.Errorf("warned with --no-load-optimize: %q", out.String())
	}
}

func TestThousands(t *testing.T) {
	for n, want := range map[int]string{0: "0", 999: "999", 8000: "8,000", 321723: "321,723", 1234567: "1,234,567"} {
		if got := thousands(n); got != want {
			t.Errorf("%d: %q, want %q", n, got, want)
		}
	}
}

// doctor reads installed skills' descriptions and warns when they overflow
// the agent's skill listing; a small install passes.
func TestListingDoctorDetail(t *testing.T) {
	dir := t.TempDir()
	write := func(name, desc string) {
		d := filepath.Join(dir, name)
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
		body := "---\nname: " + name + "\ndescription: " + desc + "\n---\n\n# x\n"
		if err := os.WriteFile(filepath.Join(d, "SKILL.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("small-skill", "short")
	if d := listingDoctorDetail("claude", dir); d != "" {
		t.Errorf("small install warned: %q", d)
	}
	for i := 0; i < 100; i++ {
		write(fmt.Sprintf("skill-%03d", i), strings.Repeat("d", 200))
	}
	d := listingDoctorDetail("claude", dir)
	if !strings.Contains(d, "of 101 installed skills") || !strings.Contains(d, "grimoire mcp install --target claude") {
		t.Errorf("detail = %q", d)
	}
	if listingDoctorDetail("gemini", dir) != "" {
		t.Error("unknown budget warned")
	}
}

// The install path lists skills without descriptions; the estimate reads them.
func TestWithDescriptions(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "my-skill")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("---\nname: my-skill\ndescription: Use when testing.\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := withDescriptions([]skills.Skill{{Name: "my-skill", Path: dir}, {Name: "kept", Description: "already"}})
	if got[0].Description != "Use when testing." || got[1].Description != "already" {
		t.Errorf("got %+v", got)
	}
}
