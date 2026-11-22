package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readFile(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestConfigureMCPInstruction_CreatesFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "sub", "CLAUDE.md")
	changed, err := ConfigureMCPInstruction(p)
	if err != nil || !changed {
		t.Fatalf("changed=%v err=%v", changed, err)
	}
	want := sectionHeader + "\n" + MCPTriggerLine + "\n"
	if got := readFile(t, p); got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

func TestConfigureMCPInstruction_ReplacesSkillsTriggerKeepsContent(t *testing.T) {
	p := filepath.Join(t.TempDir(), "CLAUDE.md")
	orig := "@RTK.md\n\n" + sectionHeader + "\n" + triggerLine + "\n\n## Other\nkeep me\n"
	if err := os.WriteFile(p, []byte(orig), 0o644); err != nil {
		t.Fatal(err)
	}
	if changed, err := ConfigureMCPInstruction(p); err != nil || !changed {
		t.Fatalf("changed=%v err=%v", changed, err)
	}
	got := readFile(t, p)
	if strings.Contains(got, triggerLine) {
		t.Error("skills trigger should be replaced")
	}
	if strings.Count(got, sectionHeader) != 1 || strings.Count(got, MCPTriggerLine) != 1 {
		t.Errorf("want exactly one grimoire section:\n%s", got)
	}
	for _, keep := range []string{"@RTK.md", "## Other", "keep me"} {
		if !strings.Contains(got, keep) {
			t.Errorf("lost %q:\n%s", keep, got)
		}
	}
	// Idempotent.
	if changed, err := ConfigureMCPInstruction(p); err != nil || changed {
		t.Errorf("second run changed=%v err=%v", changed, err)
	}
}

func TestRemoveMCPInstruction(t *testing.T) {
	dir := t.TempDir()
	if changed, err := RemoveMCPInstruction(filepath.Join(dir, "missing.md")); err != nil || changed {
		t.Errorf("missing file: changed=%v err=%v", changed, err)
	}
	p := filepath.Join(dir, "CLAUDE.md")
	if err := os.WriteFile(p, []byte("@RTK.md\n\n"+sectionHeader+"\n"+MCPTriggerLine+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if changed, err := RemoveMCPInstruction(p); err != nil || !changed {
		t.Fatalf("changed=%v err=%v", changed, err)
	}
	if got := readFile(t, p); got != "@RTK.md\n" {
		t.Errorf("got %q; want %q", got, "@RTK.md\n")
	}
}
