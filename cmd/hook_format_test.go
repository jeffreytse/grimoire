package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestWriteHookOutput_Formats(t *testing.T) {
	orig := hookFormat
	t.Cleanup(func() { hookFormat = orig })

	hookFormat = "flat"
	var b bytes.Buffer
	writeHookOutput(&b, "UserPromptSubmit", "UserPromptSubmit", "ROUTER", "grimoire ▸ ready")
	var flat map[string]any
	if err := json.Unmarshal(b.Bytes(), &flat); err != nil {
		t.Fatalf("flat not JSON: %v: %s", err, b.String())
	}
	if flat["additionalContext"] != "ROUTER" || len(flat) != 1 {
		t.Errorf("flat = %v, want only additionalContext", flat)
	}

	hookFormat = "text"
	b.Reset()
	writeHookOutput(&b, "", "UserPromptSubmit", "ROUTER", "status")
	if strings.TrimSpace(b.String()) != "ROUTER" {
		t.Errorf("text = %q, want the context only", b.String())
	}

	hookFormat = "claude"
	b.Reset()
	writeHookOutput(&b, "", "SessionStart", "ROUTER", "status")
	var cl hookOutput
	if err := json.Unmarshal(b.Bytes(), &cl); err != nil || cl.HookSpecificOutput.AdditionalContext != "ROUTER" || cl.HookSpecificOutput.HookEventName != "SessionStart" {
		t.Errorf("claude = %s (%v)", b.String(), err)
	}
}
