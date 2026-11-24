package cmd

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func callSearch(t *testing.T, args map[string]any) (hits []mcpSearchSkillsResult, errText string) {
	t.Helper()
	req := mcp.CallToolRequest{}
	req.Params.Name = "grimoire_search_skills"
	req.Params.Arguments = args
	res, err := toolGrimoireSearchSkills(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	text := res.Content[0].(mcp.TextContent).Text
	if res.IsError {
		return nil, text
	}
	var out []mcpSearchSkillsResult
	if err := json.Unmarshal([]byte(text), &out); err != nil {
		t.Fatalf("not JSON: %v: %s", err, text)
	}
	return out, ""
}

// Search freely: several phrasings in one call, filters, and paging with
// stable ranks.
func TestSearchSkills_QueriesFiltersPaging(t *testing.T) {
	setupListSkillsHome(t)
	all, errText := callSearch(t, map[string]any{"queries": []any{"write unit tests", "commit message", "readme"}, "top_n": 10})
	if errText != "" || len(all) != 3 {
		t.Fatalf("multi-query: %d results, err %q", len(all), errText)
	}
	for i, r := range all {
		if r.Rank != i+1 {
			t.Errorf("rank %d at index %d", r.Rank, i)
		}
	}
	page2, _ := callSearch(t, map[string]any{"queries": []any{"write unit tests", "commit message", "readme"}, "top_n": 1, "offset": 1})
	if len(page2) != 1 || page2[0].Name != all[1].Name || page2[0].Rank != 2 {
		t.Errorf("page 2 = %+v, want %s at rank 2", page2, all[1].Name)
	}
	only, _ := callSearch(t, map[string]any{"query": "write", "domain": "writing"})
	for _, r := range only {
		if r.Domain != "writing" {
			t.Errorf("domain filter leaked %s (%s)", r.Name, r.Domain)
		}
	}
	if _, errText := callSearch(t, map[string]any{}); !strings.Contains(errText, "required") {
		t.Errorf("no query: %q", errText)
	}
	many := make([]any, mcpSearchQueriesMax+1)
	for i := range many {
		many[i] = "x"
	}
	if _, errText := callSearch(t, map[string]any{"queries": many}); !strings.Contains(errText, "at most") {
		t.Errorf("too many queries: %q", errText)
	}
}

func TestRelatedSkills(t *testing.T) {
	for in, want := range map[string]string{
		"---\nname: a\nrelated: [b, \"c\"]\n---\nbody related: [x]": "b,c",
		"---\nname: a\nrelated:\n  - b\n  - c\ntags: [t]\n---\n":    "b,c",
		"---\nname: a\n---\nrelated: [not-front-matter]":            "",
	} {
		if got := strings.Join(relatedSkills(in), ","); got != want {
			t.Errorf("%q: %q, want %q", in, got, want)
		}
	}
}
