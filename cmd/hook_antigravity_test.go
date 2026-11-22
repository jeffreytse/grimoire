package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/jeffreytse/grimoire/internal/agent"
)

func agySteps(t *testing.T, out *bytes.Buffer) []agyStep {
	t.Helper()
	if out.Len() == 0 {
		return nil
	}
	var o agyOutput
	if err := json.Unmarshal(out.Bytes(), &o); err != nil {
		t.Fatalf("output is not JSON: %v: %s", err, out.String())
	}
	return o.InjectSteps
}

func TestAgySession_InjectsInstruction(t *testing.T) {
	var out bytes.Buffer
	runAgySession(strings.NewReader(`{"conversationId":"c1"}`), &out)
	steps := agySteps(t, &out)
	if len(steps) != 1 || steps[0].SystemMessage == nil || steps[0].SystemMessage.SystemMessage != agent.MCPTriggerLine {
		t.Fatalf("steps = %+v", steps)
	}
}

// PreInvocation fires before every model call; context is injected once per
// new prompt in a conversation.
func TestAgyPrompt_OncePerPrompt(t *testing.T) {
	setupListSkillsHome(t)
	state := t.TempDir()
	run := func(conv, prompt string) []agyStep {
		in, _ := json.Marshal(map[string]any{"conversationId": conv, "lastUserInput": prompt, "invocationNum": 0})
		var out bytes.Buffer
		runAgyPrompt(bytes.NewReader(in), &out, state)
		return agySteps(t, &out)
	}
	steps := run("conv/1", "use grimoire:write-unit-test")
	if len(steps) != 1 || !strings.Contains(steps[0].EphemeralMessage, "BODY-MARKER") {
		t.Fatalf("first call steps = %+v", steps)
	}
	if steps := run("conv/1", "use grimoire:write-unit-test"); steps != nil {
		t.Errorf("same prompt again: want nothing, got %+v", steps)
	}
	if steps := run("conv-2", "use grimoire:write-unit-test"); len(steps) != 1 {
		t.Errorf("other conversation: want context, got %+v", steps)
	}
	if steps := run("conv/1", "now grimoire:write-readme please"); len(steps) != 1 {
		t.Errorf("new prompt: want context, got %+v", steps)
	}

	var out bytes.Buffer
	runAgyPrompt(strings.NewReader("not json"), &out, state)
	if out.Len() != 0 {
		t.Errorf("bad input: want silence, got %s", out.String())
	}
}
