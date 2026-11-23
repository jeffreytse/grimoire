package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runHook(t *testing.T, prompt string) string {
	t.Helper()
	ctx, _ := runHookFull(t, prompt)
	return ctx
}

// runHookFull runs the prompt hook and returns its context and status line.
func runHookFull(t *testing.T, prompt string) (ctx, msg string) {
	t.Helper()
	in, _ := json.Marshal(map[string]string{"prompt": prompt, "hook_event_name": "UserPromptSubmit"})
	var out bytes.Buffer
	runHookPrompt(bytes.NewReader(in), &out)
	if out.Len() == 0 {
		return "", ""
	}
	var o hookOutput
	if err := json.Unmarshal(out.Bytes(), &o); err != nil {
		t.Fatalf("hook output is not JSON: %v: %s", err, out.String())
	}
	if o.HookSpecificOutput.HookEventName != "UserPromptSubmit" {
		t.Errorf("hookEventName = %q", o.HookSpecificOutput.HookEventName)
	}
	return o.HookSpecificOutput.AdditionalContext, o.SystemMessage
}

func TestHookPrompt_SuggestsMatchingSkill(t *testing.T) {
	setupListSkillsHome(t)
	orig := hookMinScore
	hookMinScore = 0.1 // tiny fixture corpus scores low
	t.Cleanup(func() { hookMinScore = orig })

	ctx := runHook(t, "write a conventional commit message")
	if !strings.Contains(ctx, "propose-conventional-commit") || !strings.Contains(ctx, "grimoire_get_skill") {
		t.Errorf("context = %q", ctx)
	}
}

func TestHookPrompt_SilentWhenNoGoodMatch(t *testing.T) {
	setupListSkillsHome(t)
	for _, p := range []string{"ok", "", "   ", "/caveman lite", "/grimoire:skill write-unit-test"} {
		if ctx := runHook(t, p); ctx != "" {
			t.Errorf("%q: want no context, got %q", p, ctx)
		}
	}
}

func TestHookPrompt_ExplicitSkillInlinesContent(t *testing.T) {
	setupListSkillsHome(t)
	ctx := runHook(t, "please use grimoire:write-unit-test, for parser.go")
	if !strings.Contains(ctx, "BODY-MARKER") || !strings.Contains(ctx, `"write-unit-test"`) {
		t.Errorf("context = %q", ctx)
	}

	ctx = runHook(t, "grimoire:propose-conventional-message now")
	if !strings.Contains(ctx, "did you mean") || strings.Contains(ctx, "?.") {
		t.Errorf("typo context = %q", ctx)
	}
}

// A skill too big for Claude Code's hook context is named for
// grimoire_get_skill rather than inlined (inlined text over the cap is moved
// to a file the model has to Read back).
func TestHookPrompt_LargeSkillPointsToGetSkill(t *testing.T) {
	setupListSkillsHome(t)
	orig := hookInlineBudget
	t.Cleanup(func() { hookInlineBudget = orig })

	hookInlineBudget = len(mcpChainRule) + 10 // nothing fits after the rule
	ctx := runHook(t, "grimoire:write-unit-test for parser.go")
	if strings.Contains(ctx, "BODY-MARKER") {
		t.Error("large skill was inlined")
	}
	if !strings.Contains(ctx, `call grimoire_get_skill with name "write-unit-test"`) || !strings.Contains(ctx, mcpChainRule) {
		t.Errorf("context = %q", ctx)
	}

	// Budget is shared: the first skill inlines, the second is pointed to.
	hookInlineBudget = orig
	ctx = runHook(t, "grimoire:write-unit-test and grimoire:write-readme")
	if !strings.Contains(ctx, "BODY-MARKER") || strings.Count(ctx, mcpChainRule) != 1 {
		t.Errorf("two-skill context: want inlined content and one chain rule, got %q", ctx)
	}
	if n := len(ctx); n > hookInlineBudget {
		t.Errorf("context is %d bytes, over budget %d", n, hookInlineBudget)
	}
}

func TestHookPrompt_BadInputIsSilent(t *testing.T) {
	for _, in := range []string{"", "not json", `{"prompt":42}`} {
		var out bytes.Buffer
		runHookPrompt(strings.NewReader(in), &out)
		if out.Len() != 0 {
			t.Errorf("%q: want no output, got %q", in, out.String())
		}
	}
}

func TestExplicitSkillRefs(t *testing.T) {
	got := explicitSkillRefs("grimoire:a-b, then (grimoire:c) and grimoire:a-b. not xgrimoire:d")
	if strings.Join(got, ",") != "a-b,c" {
		t.Errorf("refs = %v", got)
	}
}

// Gemini CLI sends BeforeAgent; the reply must carry the same event name or
// Gemini drops the context. A missing event name falls back to Claude's.
func TestHookPrompt_EchoesEventName(t *testing.T) {
	setupListSkillsHome(t)
	for event, want := range map[string]string{"BeforeAgent": "BeforeAgent", "": "UserPromptSubmit"} {
		in, _ := json.Marshal(map[string]string{"prompt": "use grimoire:write-unit-test", "hook_event_name": event})
		var out bytes.Buffer
		runHookPrompt(bytes.NewReader(in), &out)
		var o hookOutput
		if err := json.Unmarshal(out.Bytes(), &o); err != nil {
			t.Fatalf("%q: output is not JSON: %v: %s", event, err, out.String())
		}
		if o.HookSpecificOutput.HookEventName != want {
			t.Errorf("event %q: hookEventName = %q; want %q", event, o.HookSpecificOutput.HookEventName, want)
		}
	}
}

// Every substantive prompt gets the start-best-practice reminder, long or
// short, matches or not; small talk and follow-ups get nothing.
func TestHookPrompt_RouteReminderOnEveryRequest(t *testing.T) {
	setupListSkillsHome(t)
	for _, p := range []string{
		"migrate the billing service",
		"We're launching our SaaS in 48 hours. I'm terrified something will break. What do we do?",
	} {
		ctx := runHook(t, p)
		if !strings.HasPrefix(ctx, hookRouteReminder) {
			t.Errorf("%q: want route reminder first, got %q", p, ctx)
		}
		if strings.Contains(ctx, "ignore this note") {
			t.Errorf("%q: reminder must not offer an opt-out, got %q", p, ctx)
		}
	}
	for _, p := range []string{"ok", "Thanks!", "looks good", "go ahead", "fix it"} {
		if ctx := runHook(t, p); ctx != "" {
			t.Errorf("%q: want silence, got %q", p, ctx)
		}
	}
}

// Keyword matches are offered as unverified candidates after the reminder.
func TestHookPrompt_CandidatesFollowReminder(t *testing.T) {
	setupListSkillsHome(t)
	orig := hookMinScore
	hookMinScore = 0.1
	t.Cleanup(func() { hookMinScore = orig })

	ctx := runHook(t, "write a conventional commit message")
	r := strings.Index(ctx, hookRouteReminder)
	c := strings.Index(ctx, "propose-conventional-commit")
	if r != 0 || c < 0 || !strings.Contains(ctx, "unverified") {
		t.Errorf("context = %q", ctx)
	}
}

func runSessionHook(t *testing.T, event string) string {
	t.Helper()
	in, _ := json.Marshal(map[string]string{"hook_event_name": event, "source": "startup"})
	var out bytes.Buffer
	runHookSession(bytes.NewReader(in), &out)
	var o hookOutput
	if err := json.Unmarshal(out.Bytes(), &o); err != nil {
		t.Fatalf("session hook output is not JSON: %v: %s", err, out.String())
	}
	want := event
	if want == "" {
		want = "SessionStart"
	}
	if o.HookSpecificOutput.HookEventName != want {
		t.Errorf("hookEventName = %q; want %q", o.HookSpecificOutput.HookEventName, want)
	}
	return o.HookSpecificOutput.AdditionalContext
}

const routerFixture = `---
name: start-best-practice
description: Use when responding to any user request
---

# Start Best Practice

## Why This Is Best Practice

WHY-MARKER

## Steps

### Step 0: Route

| Intent signal | Route |
|---|---|
| Describes a problem | invoke ` + "`suggest-best-practice`" + ` |

` + "```" + `
## not a heading inside a fence
` + "```" + `

## Rules

- RULES-MARKER

## Common Mistakes

MISTAKES-MARKER
`

// The session hook inlines the router's Steps and Rules, nothing else.
func TestHookSession_InlinesRouter(t *testing.T) {
	setupListSkillsHome(t)
	dir := filepath.Join(os.Getenv("GRIMOIRE_SOURCE"), "skills", "meta", "skills", "start-best-practice")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(routerFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx := runSessionHook(t, "SessionStart")
	for _, want := range []string{"suggest-best-practice", "RULES-MARKER", "not a heading inside a fence", mcpChainRule} {
		if !strings.Contains(ctx, want) {
			t.Errorf("missing %q in %q", want, ctx)
		}
	}
	for _, bad := range []string{"WHY-MARKER", "MISTAKES-MARKER"} {
		if strings.Contains(ctx, bad) {
			t.Errorf("unexpected %q in %q", bad, ctx)
		}
	}
	if len(ctx) > hookInlineBudget {
		t.Errorf("context is %d bytes, over budget %d", len(ctx), hookInlineBudget)
	}
	if got := runSessionHook(t, "BeforeAgent"); got == "" {
		t.Error("event name not echoed")
	}
}

// Without the router installed (or over budget) the session hook still tells
// the model to load it.
func TestHookSession_FallsBackToDirective(t *testing.T) {
	setupListSkillsHome(t)
	ctx := runSessionHook(t, "")
	if !strings.Contains(ctx, `grimoire_get_skill with name "start-best-practice"`) {
		t.Errorf("context = %q", ctx)
	}
}

// Status lines state only facts the hook knows: a skill it loaded, or at the
// full level the unverified keyword candidates. Brief shows no per-prompt
// line — the model's own blocks show the route.
func TestHookPrompt_StatusLineFactsOnly(t *testing.T) {
	setupListSkillsHome(t)
	orig := hookMinScore
	hookMinScore = 0.1
	t.Cleanup(func() { hookMinScore = orig })

	ctx, msg := runHookFull(t, "write a conventional commit message")
	if msg != "" || !strings.Contains(ctx, "propose-conventional-commit") || !strings.Contains(ctx, "grimoire display: brief") {
		t.Errorf("brief: ctx=%q msg=%q", ctx, msg)
	}
	t.Setenv("GRIMOIRE_TRACE", "full")
	if _, msg := runHookFull(t, "write a conventional commit message"); !strings.HasPrefix(msg, "grimoire ▸ keyword candidates (unverified): ") {
		t.Errorf("full candidates status = %q", msg)
	}
	t.Setenv("GRIMOIRE_TRACE", "")
	if _, msg := runHookFull(t, "use grimoire:write-unit-test and grimoire:nope-skill"); msg != "grimoire ▸ loading skill write-unit-test · unknown skill nope-skill" {
		t.Errorf("explicit status = %q", msg)
	}
	if ctx, msg := runHookFull(t, "thanks!"); ctx != "" || msg != "" {
		t.Errorf("small talk: ctx=%q msg=%q", ctx, msg)
	}
}

// Display levels come from GRIMOIRE_TRACE, then core.trace; "on" reads as
// brief; off silences status lines and tells the model to stay silent.
func TestDisplayLevelSetting(t *testing.T) {
	setupListSkillsHome(t)
	for v, want := range map[string]displayLevel{"off": displayOff, "brief": displayBrief, "full": displayFull, "on": displayBrief, "bogus": displayBrief} {
		t.Setenv("GRIMOIRE_TRACE", v)
		if got := displayLevelSetting(); got != want {
			t.Errorf("GRIMOIRE_TRACE=%s: level %q, want %q", v, got, want)
		}
	}
	t.Setenv("GRIMOIRE_TRACE", "off")
	ctx, msg := runHookFull(t, "migrate the billing service")
	if msg != "" || !strings.Contains(ctx, displayOffRule) {
		t.Errorf("off: ctx=%q msg=%q", ctx, msg)
	}

	t.Setenv("GRIMOIRE_TRACE", "")
	cfg := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "grimoire")
	if err := os.MkdirAll(cfg, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cfg, "grimoire.toml"), []byte("[core]\ntrace = \"full\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := displayLevelSetting(); got != displayFull {
		t.Errorf("core.trace = full in %s: got %q", cfg, got)
	}
}

func TestHookSession_StatusLine(t *testing.T) {
	setupListSkillsHome(t)
	if _, msg := sessionTurn(displayBrief); msg != "grimoire ▸ no skills package with start-best-practice — run `grimoire update`" {
		t.Errorf("missing-router status = %q", msg)
	}
	dir := filepath.Join(os.Getenv("GRIMOIRE_SOURCE"), "skills", "meta", "skills", "start-best-practice")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(routerFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx, msg := sessionTurn(displayBrief)
	if msg != "grimoire ▸ ready · start-best-practice (1 route) · display brief" || !strings.Contains(ctx, "grimoire display: brief") {
		t.Errorf("status = %q", msg)
	}
	if ctx, msg := sessionTurn(displayOff); msg != "" || !strings.Contains(ctx, displayOffRule) {
		t.Errorf("off: msg = %q", msg)
	}
}

// A router skill with a Display section is the single source of the format:
// the session hook inlines it and the CLI only names the level.
func TestHookSession_RouterDisplaySection(t *testing.T) {
	setupListSkillsHome(t)
	dir := filepath.Join(os.Getenv("GRIMOIRE_SOURCE"), "skills", "meta", "skills", "start-best-practice")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if got := displayContext(displayBrief); !strings.Contains(got, "one line starting with `grimoire ▸`") {
		t.Errorf("fallback without Display = %q", got)
	}
	withDisplay := routerFixture + "\n## Display\n\nDISPLAY-MARKER\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(withDisplay), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx, _ := sessionTurn(displayFull)
	if !strings.Contains(ctx, "DISPLAY-MARKER") || !strings.Contains(ctx, "Display section, at the full level") {
		t.Errorf("ctx = %q", ctx)
	}
	if n := skillNote(); !strings.Contains(n, "Display section, at the brief level") {
		t.Errorf("note = %q", n)
	}
}
