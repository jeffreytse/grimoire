package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jeffreytse/grimoire/internal/config"
)

func TestParseToggle(t *testing.T) {
	for p, want := range map[string]string{
		"grimoire off":                    "off",
		"Grimoire off.":                   "off",
		"turn off grimoire":               "off",
		"please disable grimoire mode":    "off",
		"quiet grimoire":                  "off",
		"grimoire on":                     "on",
		"enable grimoire":                 "on",
		"Activate grimoire!":              "on",
		"resume grimoire":                 "on",
		"is grimoire off by default?":     "",
		"turn off grimoire logging":       "",
		"why did grimoire pick that?":     "",
		"we should enable grimoire later": "",
	} {
		on, off := parseToggle(p)
		got := ""
		if on {
			got = "on"
		} else if off {
			got = "off"
		}
		if got != want {
			t.Errorf("%q: %q, want %q", p, got, want)
		}
	}
}

// runPromptHook runs the MCP-mode prompt hook for a session.
func runPromptHook(t *testing.T, session, prompt string) (ctx, msg string) {
	t.Helper()
	in, _ := json.Marshal(map[string]string{"prompt": prompt, "hook_event_name": "UserPromptSubmit", "session_id": session})
	var out bytes.Buffer
	runHookPrompt(bytes.NewReader(in), &out)
	if out.Len() == 0 {
		return "", ""
	}
	var o hookOutput
	if err := json.Unmarshal(out.Bytes(), &o); err != nil {
		t.Fatalf("hook output is not JSON: %v: %s", err, out.String())
	}
	return o.HookSpecificOutput.AdditionalContext, o.SystemMessage
}

// "grimoire off" turns grimoire off for that session only; other sessions
// keep routing; a named skill still loads; "grimoire on" restores routing.
func TestSessionToggle(t *testing.T) {
	setupListSkillsHome(t)

	ctx, msg := runPromptHook(t, "s1", "grimoire off")
	if !strings.Contains(ctx, "turned grimoire off") || msg != toggleStatus(false) {
		t.Fatalf("off toggle: ctx=%q msg=%q", ctx, msg)
	}
	if ctx, msg := runPromptHook(t, "s1", "migrate the billing service"); ctx != inactiveRule || msg != "" {
		t.Errorf("inactive request: ctx=%q msg=%q", ctx, msg)
	}
	if ctx, _ := runPromptHook(t, "s2", "migrate the billing service"); !strings.HasPrefix(ctx, hookRouteReminder) {
		t.Errorf("other session lost routing: %q", ctx)
	}
	if ctx, _ := runPromptHook(t, "s1", "use grimoire:write-unit-test for parser.go"); !strings.Contains(ctx, "BODY-MARKER") || !strings.HasSuffix(ctx, inactiveRule) {
		t.Errorf("named skill while inactive: %q", ctx)
	}
	if ctx, _ := runPromptHook(t, "s1", "thanks"); ctx != "" {
		t.Errorf("small talk while inactive: %q", ctx)
	}

	if _, msg := runPromptHook(t, "s1", "grimoire on"); msg != toggleStatus(true) {
		t.Fatalf("on toggle: %q", msg)
	}
	if ctx, _ := runPromptHook(t, "s1", "migrate the billing service"); !strings.HasPrefix(ctx, hookRouteReminder) {
		t.Errorf("routing not restored: %q", ctx)
	}
}

// A toggle without a session id cannot stick, and says so.
func TestSessionToggle_NoSession(t *testing.T) {
	setupListSkillsHome(t)
	ctx, msg := runPromptHook(t, "", "grimoire off")
	if !strings.Contains(ctx, "could not be saved") || !strings.Contains(msg, "could not change grimoire") {
		t.Errorf("ctx=%q msg=%q", ctx, msg)
	}
}

// Precedence: GRIMOIRE_ACTIVE, then the session's choice, then core.active.
func TestGrimoireActivePrecedence(t *testing.T) {
	setupListSkillsHome(t)
	if a, src := grimoireActive("s1"); !a || src != activeFromDefault {
		t.Errorf("default: %v %s", a, src)
	}
	if err := setGlobalActive(false); err != nil {
		t.Fatal(err)
	}
	if a, src := grimoireActive("s1"); a || src != activeFromConfig {
		t.Errorf("after disable: %v %s", a, src)
	}
	if err := setSessionActive("s1", true); err != nil {
		t.Fatal(err)
	}
	if a, src := grimoireActive("s1"); !a || src != activeFromSession {
		t.Errorf("session on over config off: %v %s", a, src)
	}
	t.Setenv("GRIMOIRE_ACTIVE", "0")
	if a, src := grimoireActive("s1"); a || src != activeFromEnv {
		t.Errorf("env off over session on: %v %s", a, src)
	}
	t.Setenv("GRIMOIRE_ACTIVE", "")
	if err := setGlobalActive(true); err != nil {
		t.Fatal(err)
	}
	if a, _ := grimoireActive("s2"); !a {
		t.Error("enable did not restore")
	}
}

// With grimoire disabled, the session hook sends the inactive rule instead
// of the router and tells the user why; skills mode does the same.
func TestSessionHook_Inactive(t *testing.T) {
	setupListSkillsHome(t)
	if err := setGlobalActive(false); err != nil {
		t.Fatal(err)
	}
	for name, run := range map[string]func(*bytes.Reader, *bytes.Buffer){
		"mcp":    func(in *bytes.Reader, out *bytes.Buffer) { runHookSession(in, out) },
		"skills": func(in *bytes.Reader, out *bytes.Buffer) { runSkillsSession(in, out) },
	} {
		in, _ := json.Marshal(map[string]string{"hook_event_name": "SessionStart", "session_id": "s1"})
		var out bytes.Buffer
		run(bytes.NewReader(in), &out)
		var o hookOutput
		if err := json.Unmarshal(out.Bytes(), &o); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if o.HookSpecificOutput.AdditionalContext != inactiveRule || !strings.Contains(o.SystemMessage, "off for all sessions") {
			t.Errorf("%s: ctx=%q msg=%q", name, o.HookSpecificOutput.AdditionalContext, o.SystemMessage)
		}
	}
}

// Agents without hooks learn grimoire is off from the MCP server.
func TestMCP_Inactive(t *testing.T) {
	setupListSkillsHome(t)
	if !strings.HasPrefix(routeDirective(), mcpRouteDirective) {
		t.Error("active directive changed")
	}
	t.Setenv("GRIMOIRE_ACTIVE", "false")
	if instr := mcpServerInstructions(displayBrief); !strings.HasPrefix(instr, mcpInactiveInstruction) {
		t.Errorf("instructions = %q", instr)
	}
	if routeDirective() != mcpInactiveDirective {
		t.Errorf("directive = %q", routeDirective())
	}
	if n := skillNote(); !strings.Contains(n, "grimoire is inactive") {
		t.Errorf("note = %q", n)
	}
}

// The badge shows "off" for a session that turned grimoire off.
func TestStatusLine_Off(t *testing.T) {
	setupListSkillsHome(t)
	if err := setSessionActive("s1", false); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	runStatusLine(strings.NewReader(`{"session_id":"s1"}`), &out, "", grimoireStateRoot())
	if out.String() != "[GRIMOIRE off]" {
		t.Errorf("off badge = %q", out.String())
	}
}

// grimoire disable / enable write core.active to the global config.
func TestEnableDisable_WriteConfig(t *testing.T) {
	setupListSkillsHome(t)
	if err := setGlobalActive(false); err != nil {
		t.Fatal(err)
	}
	fs, err := config.LoadFile(config.GlobalPath())
	if err != nil || fs.Core.Active != "false" {
		t.Fatalf("after disable: active=%q err=%v", fs.Core.Active, err)
	}
	if b, _ := os.ReadFile(config.GlobalPath()); !strings.Contains(string(b), "active = false") {
		t.Errorf("config file:\n%s", b)
	}
	var cfs config.FileConfig
	if err := applyKey(&cfs, "core.active", "maybe"); err == nil {
		t.Error("core.active = maybe: want error")
	}
	if err := applyKey(&cfs, "core.active", "0"); err != nil || cfs.Core.Active != "false" {
		t.Errorf("core.active = 0: %q %v", cfs.Core.Active, err)
	}
}

// Antigravity toggles per conversation too, and — having no user-visible
// hook output — never uses up the one-time tip.
func TestAgy_ToggleAndTip(t *testing.T) {
	setupListSkillsHome(t)
	dir, _ := os.UserCacheDir()
	tip := filepath.Join(dir, "grimoire", "tip-shown")
	if err := os.Remove(tip); err != nil {
		t.Fatal(err)
	}
	run := func(prompt string) string {
		in, _ := json.Marshal(map[string]any{"conversationId": "c1", "lastUserInput": prompt})
		var out bytes.Buffer
		runAgyPrompt(bytes.NewReader(in), &out, "")
		steps := agySteps(t, &out)
		if len(steps) == 0 {
			return ""
		}
		return steps[0].EphemeralMessage
	}
	if ctx := run("migrate the billing service"); !strings.HasPrefix(ctx, hookRouteReminder) {
		t.Errorf("active: %q", ctx)
	}
	if _, err := os.Stat(tip); err == nil {
		t.Error("Antigravity used up the one-time tip")
	}
	if ctx := run("grimoire off"); !strings.Contains(ctx, "turned grimoire off") {
		t.Errorf("toggle: %q", ctx)
	}
	if ctx := run("migrate the billing service again"); ctx != inactiveRule {
		t.Errorf("inactive: %q", ctx)
	}
}
