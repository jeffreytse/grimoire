package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func runHook(t *testing.T, prompt string) string {
	t.Helper()
	in, _ := json.Marshal(map[string]string{"prompt": prompt, "hook_event_name": "UserPromptSubmit"})
	var out bytes.Buffer
	runHookPrompt(bytes.NewReader(in), &out)
	if out.Len() == 0 {
		return ""
	}
	var o hookOutput
	if err := json.Unmarshal(out.Bytes(), &o); err != nil {
		t.Fatalf("hook output is not JSON: %v: %s", err, out.String())
	}
	if o.HookSpecificOutput.HookEventName != "UserPromptSubmit" {
		t.Errorf("hookEventName = %q", o.HookSpecificOutput.HookEventName)
	}
	return o.HookSpecificOutput.AdditionalContext
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
