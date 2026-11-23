package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestStatusLine_Badge(t *testing.T) {
	dir := t.TempDir()
	last := filepath.Join(dir, "last-skill")

	var out bytes.Buffer
	runStatusLine(strings.NewReader("{}"), &out, "", last)
	if !strings.Contains(out.String(), "[GRIMOIRE]") {
		t.Errorf("no last skill: got %q", out.String())
	}

	if err := os.WriteFile(last, []byte("Suggest-Best-Practice\x1b[31m\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	runStatusLine(strings.NewReader("{}"), &out, "", last)
	if !strings.Contains(out.String(), "[GRIMOIRE:suggest-best-practice31m]") || strings.Count(out.String(), "\x1b") != 2 {
		t.Errorf("sanitized skill: got %q", out.String())
	}

	// Stale skill names drop out of the badge.
	old := time.Now().Add(-2 * statusLineFresh)
	if err := os.Chtimes(last, old, old); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	runStatusLine(strings.NewReader("{}"), &out, "", last)
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
