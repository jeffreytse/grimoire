package cmd

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func getSkillPrompt(t *testing.T, name string) (*mcp.GetPromptResult, error) {
	t.Helper()
	req := mcp.GetPromptRequest{}
	req.Params.Name = "skill"
	req.Params.Arguments = map[string]string{"name": name}
	return promptGrimoireSkill(context.Background(), req)
}

func promptText(t *testing.T, res *mcp.GetPromptResult) string {
	t.Helper()
	if len(res.Messages) != 1 {
		t.Fatalf("want 1 message, got %d", len(res.Messages))
	}
	tc, ok := res.Messages[0].Content.(mcp.TextContent)
	if !ok {
		t.Fatalf("content is %T, want TextContent", res.Messages[0].Content)
	}
	return tc.Text
}

func TestSkillPrompt_ExactAndPrefixed(t *testing.T) {
	setupListSkillsHome(t)
	for _, in := range []string{"write-unit-test", "grimoire:write-unit-test", "  grimoire:write-unit-test "} {
		res, err := getSkillPrompt(t, in)
		if err != nil {
			t.Fatalf("%q: %v", in, err)
		}
		if text := promptText(t, res); !strings.Contains(text, "BODY-MARKER") || !strings.Contains(text, "write-unit-test") {
			t.Errorf("%q: prompt missing skill content", in)
		}
		if res.Messages[0].Role != mcp.RoleUser {
			t.Errorf("role = %s; want user", res.Messages[0].Role)
		}
	}
}

func TestResolveSkillOrSuggest_Typo(t *testing.T) {
	setupListSkillsHome(t)
	_, _, err := resolveSkillOrSuggest("grimoire:propose-conventional-message")
	if err == nil {
		t.Fatal("typo should not resolve")
	}
	if !strings.Contains(err.Error(), "did you mean") || !strings.Contains(err.Error(), "propose-conventional-commit") {
		t.Errorf("want suggestion of propose-conventional-commit, got %v", err)
	}
	if _, _, err := resolveSkillOrSuggest("  "); err == nil {
		t.Error("empty name should error")
	}
}

func TestGetSkillTool_AcceptsPrefixAndSuggests(t *testing.T) {
	setupListSkillsHome(t)
	call := func(name string) (string, bool) {
		req := mcp.CallToolRequest{}
		req.Params.Arguments = map[string]any{"name": name}
		res, err := toolGrimoireGetSkill(context.Background(), req)
		if err != nil {
			t.Fatal(err)
		}
		tc, ok := res.Content[0].(mcp.TextContent)
		if !ok {
			t.Fatalf("content is %T", res.Content[0])
		}
		return tc.Text, res.IsError
	}
	if out, isErr := call("grimoire:write-readme"); isErr || !strings.Contains(out, "BODY-MARKER") {
		t.Errorf("prefixed name: isErr=%v out=%.80s", isErr, out)
	}
	if out, isErr := call("write-readmee"); !isErr || !strings.Contains(out, "write-readme") {
		t.Errorf("typo: isErr=%v out=%s", isErr, out)
	}
	out, _ := call("write-readme")
	var got mcpGetSkillOutput
	if err := json.Unmarshal([]byte(out), &got); err != nil || got.Note != mcpChainRule {
		t.Errorf("get_skill note missing or wrong: err=%v note=%q", err, got.Note)
	}
}

// Workflow skills say "invoke X" / "load the skill using the Skill tool";
// every way grimoire delivers a skill must say how that maps to MCP tools.
func TestChainRuleDelivered(t *testing.T) {
	if !strings.Contains(mcpServerInstructions, mcpChainRule) {
		t.Error("server instructions lack the chain rule")
	}
	if n := len(mcpServerInstructions); n > 4096 {
		t.Errorf("server instructions are %d chars; Claude Code cuts at 4096", n)
	}
	setupListSkillsHome(t)
	res, err := getSkillPrompt(t, "write-unit-test")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(promptText(t, res), mcpChainRule) {
		t.Error("skill prompt lacks the chain rule")
	}
	if ctx := hookContext("grimoire:write-unit-test"); !strings.Contains(ctx, mcpChainRule) {
		t.Error("hook explicit-skill context lacks the chain rule")
	}
}

func TestSkillNameCompleter(t *testing.T) {
	setupListSkillsHome(t)
	c := skillNameCompleter{}
	got, err := c.CompletePromptArgument(context.Background(), "skill", mcp.CompleteArgument{Name: "name", Value: "write"}, mcp.CompleteContext{})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got.Values, ",") != "write-readme,write-unit-test" {
		t.Errorf("prefix completion = %v", got.Values)
	}
	// Substring matches come after prefix matches; "grimoire:" is ignored.
	got, _ = c.CompletePromptArgument(context.Background(), "skill", mcp.CompleteArgument{Name: "name", Value: "grimoire:commit"}, mcp.CompleteContext{})
	if len(got.Values) != 1 || got.Values[0] != "propose-conventional-commit" {
		t.Errorf("substring completion = %v", got.Values)
	}
	// Other prompts/arguments get nothing.
	got, _ = c.CompletePromptArgument(context.Background(), "other", mcp.CompleteArgument{Name: "name", Value: "w"}, mcp.CompleteContext{})
	if len(got.Values) != 0 {
		t.Errorf("other prompt completion = %v", got.Values)
	}
}

// The initialize result must carry the instructions and advertise the
// prompts and completions capabilities; prompts/list must hold one prompt.
func TestMCPServer_InitializeAdvertisesSkillAccess(t *testing.T) {
	setupListSkillsHome(t)
	s := server.NewMCPServer("grimoire", "test", mcpServerOptions()...)
	registerMCPPrompts(s)

	send := func(msg string, out any) {
		t.Helper()
		b, err := json.Marshal(s.HandleMessage(context.Background(), json.RawMessage(msg)))
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(b, out); err != nil {
			t.Fatal(err)
		}
	}
	var init struct {
		Result struct {
			Instructions string                     `json:"instructions"`
			Capabilities map[string]json.RawMessage `json:"capabilities"`
		} `json:"result"`
	}
	send(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"t","version":"1"}}}`, &init)
	instr := init.Result.Instructions
	if !strings.HasPrefix(instr, "grimoire serves") || !strings.Contains(instr, "grimoire:<skill-name>") {
		t.Errorf("instructions = %q", instr)
	}
	if len(instr) > 4096 {
		t.Errorf("instructions are %d chars; Claude Code cuts at 4096", len(instr))
	}
	if init.Result.Capabilities["prompts"] == nil || init.Result.Capabilities["completions"] == nil {
		t.Errorf("capabilities = %v; want prompts and completions", init.Result.Capabilities)
	}

	var list struct {
		Result struct {
			Prompts []struct {
				Name string `json:"name"`
			} `json:"prompts"`
		} `json:"result"`
	}
	send(`{"jsonrpc":"2.0","id":2,"method":"prompts/list"}`, &list)
	if len(list.Result.Prompts) != 1 || list.Result.Prompts[0].Name != "skill" {
		t.Errorf("prompts/list = %+v; want exactly the skill prompt", list.Result.Prompts)
	}
}
