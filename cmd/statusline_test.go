package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// sessionLast returns the last-skill file for session s1 under a temp state
// root, and the status line input naming that session.
func sessionLast(t *testing.T) (root, last, in string) {
	t.Helper()
	root = t.TempDir()
	dir := filepath.Join(root, "s1")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	return root, filepath.Join(dir, "last-skill"), `{"session_id":"s1"}`
}

func TestStatusLine_Badge(t *testing.T) {
	root, last, in := sessionLast(t)

	var out bytes.Buffer
	runStatusLine(strings.NewReader(in), &out, "", root)
	if !strings.Contains(out.String(), "[GRIMOIRE]") {
		t.Errorf("no last skill: got %q", out.String())
	}

	if err := os.WriteFile(last, []byte("Run-Game-Day\x1b[31m\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	runStatusLine(strings.NewReader(in), &out, "", root)
	if !strings.Contains(out.String(), "[GRIMOIRE:run-game-day31m]") || strings.Count(out.String(), "\x1b") != 2 {
		t.Errorf("sanitized skill: got %q", out.String())
	}
	if err := os.WriteFile(last, []byte("suggest-best-practice"), 0o600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	runStatusLine(strings.NewReader(in), &out, "", root)
	if !strings.Contains(out.String(), "[GRIMOIRE:suggest]") || strings.Count(out.String(), "\x1b") != 2 {
		t.Errorf("sanitized skill: got %q", out.String())
	}

	// Stale skill names drop out of the badge.
	old := time.Now().Add(-2 * statusLineFresh)
	if err := os.Chtimes(last, old, old); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	runStatusLine(strings.NewReader(in), &out, "", root)
	if !strings.Contains(out.String(), "[GRIMOIRE]") {
		t.Errorf("stale skill: got %q", out.String())
	}
}

// A symlinked state file is never read: it could point at a secret.
func TestStatusLine_RefusesSymlink(t *testing.T) {
	dir := t.TempDir()
	secret := filepath.Join(dir, "secret")
	if err := os.WriteFile(secret, []byte("topsecret"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "last-skill")
	if err := os.Symlink(secret, link); err != nil {
		t.Skip("symlinks unsupported:", err)
	}
	if got := recentLastSkill(link); got != "" {
		t.Errorf("read through symlink: %q", got)
	}
}

// --then runs the chained status line with the same stdin and appends the badge.
func TestStatusLine_Chains(t *testing.T) {
	var out bytes.Buffer
	runStatusLine(strings.NewReader("IN-MARKER"), &out, "cat; echo ' [OTHER]'", "")
	got := out.String()
	if !strings.HasPrefix(got, "IN-MARKER [OTHER] ") || !strings.HasSuffix(got, "[GRIMOIRE]\x1b[0m") {
		t.Errorf("chained output = %q", got)
	}
}

// The badge shows the workflow and the step skill served after it.
func TestStatusLine_WorkflowAndSkill(t *testing.T) {
	root, last, in := sessionLast(t)
	write := func(p, v string) {
		if err := os.WriteFile(p, []byte(v), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write(workflowFile(last), "plan-best-practice-solution")
	write(last, "design-slo")
	var out bytes.Buffer
	runStatusLine(strings.NewReader(in), &out, "", root)
	if !strings.Contains(out.String(), "[GRIMOIRE:plan › design-slo]") {
		t.Errorf("got %q", out.String())
	}
	write(last, "plan-best-practice-solution")
	out.Reset()
	runStatusLine(strings.NewReader(in), &out, "", root)
	if !strings.Contains(out.String(), "[GRIMOIRE:plan]") {
		t.Errorf("workflow only: got %q", out.String())
	}
}

// A session never shows another session's skill, and input without a
// session id shows the bare badge.
func TestStatusLine_PerSession(t *testing.T) {
	root := t.TempDir()
	recordSkill(root, "s1", "plan-best-practice-solution")
	recordSkill(root, "s1", "design-slo")
	recordSkill(root, "s2", "run-game-day")
	for in, want := range map[string]string{
		`{"session_id":"s1"}`:        "[GRIMOIRE:plan › design-slo]",
		`{"session_id":"s2"}`:        "[GRIMOIRE:run-game-day]",
		`{"session_id":"s3"}`:        "[GRIMOIRE]",
		`{}`:                         "[GRIMOIRE]",
		`{"session_id":"../../etc"}`: "[GRIMOIRE]",
	} {
		var out bytes.Buffer
		runStatusLine(strings.NewReader(in), &out, "", root)
		if !strings.Contains(out.String(), want) {
			t.Errorf("%s: got %q, want %q", in, out.String(), want)
		}
	}
}

// The badge uses grimoire's brand gold (#c9933a): 24-bit where the terminal
// supports it, the nearest 256-color otherwise.
func TestStatusLine_BrandColor(t *testing.T) {
	t.Setenv("COLORTERM", "truecolor")
	if b := statusBadge(""); !strings.HasPrefix(b, "\x1b[1;38;2;201;147;58m[GRIMOIRE]") {
		t.Errorf("truecolor badge = %q", b)
	}
	t.Setenv("COLORTERM", "")
	if b := statusBadge(""); !strings.HasPrefix(b, "\x1b[1;38;5;173m[GRIMOIRE]") {
		t.Errorf("256-color badge = %q", b)
	}
}
