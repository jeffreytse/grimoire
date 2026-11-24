package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSkillsHooks_DisplayAndStatus(t *testing.T) {
	setupListSkillsHome(t)
	ctx, msg := skillsSessionTurn(displayBrief)
	if !strings.Contains(ctx, "grimoire display: brief") || !strings.Contains(ctx, "invoke the start-best-practice skill") || strings.Contains(ctx, "grimoire_get_skill") {
		t.Errorf("session ctx = %q", ctx)
	}
	if msg != "" {
		t.Errorf("brief session msg = %q, want silent on success", msg)
	}
	if _, msg := skillsSessionTurn(displayFull); msg != `grimoire ▸ ready (skills mode) · full view — say "quiet grimoire" to hide` {
		t.Errorf("full session msg = %q", msg)
	}
	if ctx, msg := skillsSessionTurn(displayOff); msg != "" || !strings.Contains(ctx, displayOffRule) {
		t.Errorf("off: ctx=%q msg=%q", ctx, msg)
	}

	ctx, msg = skillsHookTurn("migrate the billing service", displayBrief)
	if !strings.HasPrefix(ctx, skillsRouteReminder) || strings.Contains(ctx, "grimoire_get_skill") || msg != "" {
		t.Errorf("prompt ctx=%q msg=%q", ctx, msg)
	}
	if ctx, msg := skillsHookTurn("thanks!", displayBrief); ctx != "" || msg != "" {
		t.Errorf("small talk: %q %q", ctx, msg)
	}
}

func TestSkillUsed_RecordsSkillPerSession(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	runSkillUsed(strings.NewReader(`{"session_id":"a","tool_name":"Skill","tool_input":{"skill":"grimoire:plan-best-practice-solution"}}`))
	runSkillUsed(strings.NewReader(`{"session_id":"a","tool_name":"mcp__grimoire__grimoire_get_skill","tool_input":{"names":["apply-premortem","design-slo"]}}`))
	runSkillUsed(strings.NewReader(`{"session_id":"b","tool_name":"mcp__grimoire__grimoire_get_skill","tool_input":{"name":"suggest-best-practice"}}`))
	runSkillUsed(strings.NewReader(`{"session_id":"a","tool_name":"Read","tool_input":{"skill":"ignored"}}`))
	root := grimoireStateRoot()
	for sid, want := range map[string]string{"a": "[GRIMOIRE:plan › design-slo]", "b": "[GRIMOIRE:suggest]"} {
		var out strings.Builder
		runStatusLine(strings.NewReader(`{"session_id":"`+sid+`"}`), &out, "", root)
		if !strings.Contains(out.String(), want) {
			t.Errorf("session %s: badge %q, want %q", sid, out.String(), want)
		}
	}
}

// Session context carries the resolved preferences, so workflows need no
// shell command (skills mode) or tool call (MCP mode) to read them.
func TestSessionContext_Preferences(t *testing.T) {
	setupListSkillsHome(t)
	ctx, _ := skillsSessionTurn(displayBrief)
	if !strings.Contains(ctx, "no practices pinned") {
		t.Errorf("no pins: %q", ctx)
	}
	cfg := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "grimoire")
	if err := os.MkdirAll(cfg, 0o755); err != nil {
		t.Fatal(err)
	}
	toml := "[standards.engineering]\npractices = [\"Google Engineering Practices\"]\n"
	if err := os.WriteFile(filepath.Join(cfg, "grimoire.toml"), []byte(toml), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx, _ = skillsSessionTurn(displayBrief)
	if !strings.Contains(ctx, "engineering: Google Engineering Practices") {
		t.Errorf("pinned: %q", ctx)
	}
	if mcp, _ := sessionTurn(displayBrief); !strings.Contains(mcp, "engineering: Google Engineering Practices") {
		t.Errorf("mcp session lacks preferences")
	}
}
