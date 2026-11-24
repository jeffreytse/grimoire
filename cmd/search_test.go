package cmd

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/jeffreytse/grimoire/internal/skills"
)

func runSearchJSON(t *testing.T, args ...string) []string {
	t.Helper()
	origJSON, origTop := flagSearchJSON, flagSearchTop
	t.Cleanup(func() { flagSearchJSON, flagSearchTop = origJSON, origTop })
	flagSearchJSON, flagSearchTop = true, 10

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = w
	runErr := runSearch(nil, args)
	os.Stdout = stdout
	_ = w.Close()
	out, _ := io.ReadAll(r)
	if runErr != nil {
		t.Fatal(runErr)
	}
	var got []skills.Skill
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, out)
	}
	names := make([]string, len(got))
	for i := range got {
		names[i] = got[i].Name
	}
	return names
}

// The CLI ranks like the MCP tool, so it finds word forms a substring
// match misses: "commits" never appears in the skill, "committing" does.
func TestSearch_MatchesWordForms(t *testing.T) {
	setupListSkillsHome(t)
	got := runSearchJSON(t, "commits")
	if len(got) == 0 || got[0] != "propose-conventional-commit" {
		t.Fatalf("search commits = %v, want propose-conventional-commit first", got)
	}
}

// grimoire search and the agent's grimoire_search_skills return the same
// skills in the same order.
func TestSearch_SameRankingAsMCP(t *testing.T) {
	setupListSkillsHome(t)
	cli := runSearchJSON(t, "writing tests")

	req := mcp.CallToolRequest{}
	req.Params.Name = "grimoire_search_skills"
	req.Params.Arguments = map[string]any{"query": "writing tests"}
	res, err := toolGrimoireSearchSkills(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	var hits []struct {
		Name string `json:"name"`
	}
	text := res.Content[0].(mcp.TextContent).Text
	if err := json.Unmarshal([]byte(text), &hits); err != nil {
		t.Fatalf("MCP result not JSON: %v\n%s", err, text)
	}
	if len(hits) != len(cli) {
		t.Fatalf("CLI %v vs MCP %v", cli, hits)
	}
	for i := range hits {
		if hits[i].Name != cli[i] {
			t.Errorf("rank %d: CLI %s, MCP %s", i+1, cli[i], hits[i].Name)
		}
	}
}

func TestSearch_NoQueryListsAll(t *testing.T) {
	setupListSkillsHome(t)
	if got := runSearchJSON(t); len(got) != 3 {
		t.Errorf("no query should list all 3 skills, got %v", got)
	}
}
