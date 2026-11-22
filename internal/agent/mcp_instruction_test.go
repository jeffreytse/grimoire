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

// A file written by an older grimoire (previous MCPTriggerLine) is upgraded in
// place on configure, and still cleaned up on uninstall.
func TestMCPInstruction_UpgradesAndRemovesLegacyLine(t *testing.T) {
	legacy := legacyMCPTriggerLines[0]
	if strings.Contains(legacy, "grimoire:<skill-name>") || !strings.Contains(MCPTriggerLine, "grimoire:<skill-name>") {
		t.Fatal("current line must add the grimoire:<skill-name> convention")
	}
	p := filepath.Join(t.TempDir(), "AGENTS.md")
	orig := "# Notes\n\n" + sectionHeader + "\n" + legacy + "\n"
	if err := os.WriteFile(p, []byte(orig), 0o644); err != nil {
		t.Fatal(err)
	}
	if changed, err := ConfigureMCPInstruction(p); err != nil || !changed {
		t.Fatalf("upgrade: changed=%v err=%v", changed, err)
	}
	got := readFile(t, p)
	if want := "# Notes\n\n" + sectionHeader + "\n" + MCPTriggerLine + "\n"; got != want {
		t.Errorf("upgraded file = %q; want %q", got, want)
	}
	if changed, err := ConfigureMCPInstruction(p); err != nil || changed {
		t.Errorf("re-run after upgrade: changed=%v err=%v", changed, err)
	}

	if err := os.WriteFile(p, []byte(orig), 0o644); err != nil {
		t.Fatal(err)
	}
	if changed, err := RemoveMCPInstruction(p); err != nil || !changed {
		t.Fatalf("remove legacy: changed=%v err=%v", changed, err)
	}
	if got := readFile(t, p); got != "# Notes\n" {
		t.Errorf("after remove = %q", got)
	}
}

// Every earlier trigger line must differ from the current one and be
// replaced on reinstall, so users upgrading get the routing clause.
func TestMCPTriggerLine_RoutesAndUpgradesLegacy(t *testing.T) {
	if !strings.Contains(MCPTriggerLine, "start-best-practice") {
		t.Error("trigger line does not route through start-best-practice")
	}
	for _, old := range legacyMCPTriggerLines {
		if old == MCPTriggerLine {
			t.Fatal("legacy list contains the current trigger line")
		}
		p := filepath.Join(t.TempDir(), "AGENTS.md")
		if err := os.WriteFile(p, []byte(sectionHeader+"\n"+old+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if changed, err := ConfigureMCPInstruction(p); err != nil || !changed {
			t.Fatalf("configure: changed=%v err=%v", changed, err)
		}
		b, _ := os.ReadFile(p)
		if string(b) != sectionHeader+"\n"+MCPTriggerLine+"\n" {
			t.Errorf("legacy line not replaced: %q", b)
		}
	}
}
