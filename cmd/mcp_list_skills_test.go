package cmd

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// setupListSkillsHome builds an isolated official package with a body-heavy
// skill so tests can assert bodies are never returned.
func setupListSkillsHome(t *testing.T) {
	t.Helper()
	pkg := t.TempDir()
	t.Setenv("GRIMOIRE_SOURCE", pkg) // absolute path = official package home
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("HOME", t.TempDir())
	markTipShown(t) // most tests check status lines without the one-time tip
	mk := func(rel, name, desc string) {
		dir := filepath.Join(pkg, "skills", filepath.FromSlash(rel), name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		body := "---\nname: " + name + "\ndescription: " + desc + "\ntags: [a, b]\n---\n\nBODY-MARKER " + strings.Repeat("x", 2000)
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mk("engineering/testing/skills", "write-unit-test", "Use when writing unit tests")
	mk("engineering/development/skills", "propose-conventional-commit", "Use when committing")
	mk("writing/technical/skills", "write-readme", "Use when writing a README")
}

func callListSkills(t *testing.T, args map[string]any) (string, bool) {
	t.Helper()
	req := mcp.CallToolRequest{}
	req.Params.Name = "grimoire_list_skills"
	req.Params.Arguments = args
	res, err := toolGrimoireListSkills(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	var sb strings.Builder
	for _, c := range res.Content {
		if tc, ok := c.(mcp.TextContent); ok {
			sb.WriteString(tc.Text)
		}
	}
	return sb.String(), res.IsError
}

func TestMCPListSkills_IndexByDefault(t *testing.T) {
	setupListSkillsHome(t)
	out, isErr := callListSkills(t, nil)
	if isErr {
		t.Fatalf("unexpected error: %s", out)
	}
	if strings.Contains(out, "BODY-MARKER") {
		t.Fatal("index must not include skill bodies")
	}
	var got struct {
		Total   int `json:"total"`
		Domains []struct {
			Domain     string         `json:"domain"`
			Total      int            `json:"total"`
			Subdomains map[string]int `json:"subdomains"`
		} `json:"domains"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("bad JSON: %v\n%s", err, out)
	}
	if got.Total != 3 || len(got.Domains) != 2 {
		t.Fatalf("total=%d domains=%d; want 3, 2\n%s", got.Total, len(got.Domains), out)
	}
	if d := got.Domains[0]; d.Domain != "engineering" || d.Total != 2 || d.Subdomains["testing"] != 1 {
		t.Errorf("engineering index = %+v", d)
	}
}

func TestMCPListSkills_DomainFilterReturnsMetadataOnly(t *testing.T) {
	setupListSkillsHome(t)
	out, isErr := callListSkills(t, map[string]any{"domain": "engineering", "subdomain": "testing"})
	if isErr {
		t.Fatalf("unexpected error: %s", out)
	}
	if strings.Contains(out, "BODY-MARKER") {
		t.Fatal("filtered list must not include skill bodies")
	}
	var got []mcpSkillSummary
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("bad JSON: %v\n%s", err, out)
	}
	if len(got) != 1 || got[0].Name != "write-unit-test" || got[0].Description != "Use when writing unit tests" {
		t.Fatalf("got %+v", got)
	}
}

func TestMCPListSkills_WholeDomainIsCompact(t *testing.T) {
	setupListSkillsHome(t)
	out, isErr := callListSkills(t, map[string]any{"domain": "engineering"})
	if isErr {
		t.Fatalf("unexpected error: %s", out)
	}
	var got []mcpSkillSummary
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("bad JSON: %v\n%s", err, out)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 engineering skills, got %d", len(got))
	}
	for _, s := range got {
		if s.Tags != nil || s.Package != "" {
			t.Errorf("domain-wide listing should omit tags/package: %+v", s)
		}
	}
}

func TestTruncateRunes(t *testing.T) {
	if got := truncateRunes("short", 10); got != "short" {
		t.Errorf("got %q", got)
	}
	if got := truncateRunes("héllo wörld", 5); got != "héllo…" {
		t.Errorf("got %q", got)
	}
}

func TestMCPListSkills_Errors(t *testing.T) {
	setupListSkillsHome(t)
	if out, isErr := callListSkills(t, map[string]any{"subdomain": "testing"}); !isErr {
		t.Errorf("subdomain without domain should error; got %s", out)
	}
	if out, isErr := callListSkills(t, map[string]any{"domain": "nope"}); !isErr || !strings.Contains(out, "without arguments") {
		t.Errorf("unknown domain should error with a hint; got %v %s", isErr, out)
	}
}
