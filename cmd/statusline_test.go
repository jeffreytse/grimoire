package cmd

import (
	"bytes"
	"math"
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
	t.Setenv("NO_COLOR", "1") // label logic; colors are tested in TestBadgeColors
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
	if !strings.Contains(out.String(), "[GRIMOIRE:run-game-day31m]") || strings.Count(out.String(), "\x1b") != 0 {
		t.Errorf("sanitized skill: got %q", out.String())
	}
	if err := os.WriteFile(last, []byte("suggest-best-practice"), 0o600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	runStatusLine(strings.NewReader(in), &out, "", root)
	if !strings.Contains(out.String(), "[GRIMOIRE:suggest]") || strings.Count(out.String(), "\x1b") != 0 {
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
	t.Setenv("NO_COLOR", "1") // label logic; colors are tested in TestBadgeColors
	var out bytes.Buffer
	runStatusLine(strings.NewReader("IN-MARKER"), &out, "cat; echo ' [OTHER]'", "")
	got := out.String()
	if !strings.HasPrefix(got, "IN-MARKER [OTHER] ") || !strings.HasSuffix(got, "[GRIMOIRE]") {
		t.Errorf("chained output = %q", got)
	}
}

// The badge shows the workflow and the step skill served after it.
func TestStatusLine_WorkflowAndSkill(t *testing.T) {
	t.Setenv("NO_COLOR", "1") // label logic; colors are tested in TestBadgeColors
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
	t.Setenv("NO_COLOR", "1") // label logic; colors are tested in TestBadgeColors
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

// The badge is an ink-on-gold chip in grimoire's brand colors, 24-bit where
// the terminal supports it, 256-color otherwise, and plain text under
// NO_COLOR or TERM=dumb.
func TestBadgeColors(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("COLORTERM", "truecolor")
	if b := renderBadge("GRIMOIRE"); b != "\x1b[1;38;2;26;13;6;48;2;201;147;58m GRIMOIRE \x1b[0m" {
		t.Errorf("truecolor badge = %q", b)
	}
	t.Setenv("COLORTERM", "")
	if b := renderBadge("GRIMOIRE"); b != "\x1b[1;38;5;16;48;5;173m GRIMOIRE \x1b[0m" {
		t.Errorf("256-color badge = %q", b)
	}
	t.Setenv("TERM", "dumb")
	if b := renderBadge("GRIMOIRE"); b != "[GRIMOIRE]" {
		t.Errorf("TERM=dumb badge = %q", b)
	}
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "1")
	if b := renderBadge("GRIMOIRE"); b != "[GRIMOIRE]" {
		t.Errorf("NO_COLOR badge = %q", b)
	}
}

// The badge's text/background pairs meet WCAG AA for text (4.5:1), whatever
// the terminal's own theme.
func TestBadgeContrast(t *testing.T) {
	xterm256 := map[int]rgb{173: {0xd7, 0x87, 0x5f}, 16: {0, 0, 0}}
	for name, pair := range map[string][2]rgb{
		"24-bit":    {brandInk, brandGold},
		"256-color": {xterm256[brandInk256], xterm256[brandGold256]},
	} {
		if r := contrastRatio(pair[0], pair[1]); r < 4.5 {
			t.Errorf("%s badge contrast %.2f:1, want >= 4.5:1", name, r)
		}
	}
}

// contrastRatio is the WCAG 2 contrast ratio of two colors.
func contrastRatio(a, b rgb) float64 {
	lum := func(c rgb) float64 {
		ch := func(v uint8) float64 {
			x := float64(v) / 255
			if x <= 0.03928 {
				return x / 12.92
			}
			return math.Pow((x+0.055)/1.055, 2.4)
		}
		return 0.2126*ch(c.r) + 0.7152*ch(c.g) + 0.0722*ch(c.b)
	}
	l1, l2 := lum(a), lum(b)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}
