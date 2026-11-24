package cmd

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// scopedProject makes a project dir with the given grimoire.toml (none when
// toml is ""), points grimoire at it, and adds a meta workflow to the test
// package.
func scopedProject(t *testing.T, toml string) string {
	t.Helper()
	setupListSkillsHome(t)
	dir := t.TempDir()
	t.Setenv("GRIMOIRE_PROJECT_DIR", dir)
	if toml != "" {
		writeProjectToml(t, dir, toml)
	}
	meta := filepath.Join(os.Getenv("GRIMOIRE_SOURCE"), "skills", "meta", "skills", "suggest-best-practice")
	if err := os.MkdirAll(meta, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(meta, "SKILL.md"), []byte("---\nname: suggest-best-practice\ndescription: Use when suggesting a practice to commit to\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	scopeMu.Lock()
	scopeCache = map[string]cachedScope{}
	scopeMu.Unlock()
	return dir
}

func writeProjectToml(t *testing.T, dir, toml string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "grimoire.toml"), []byte(toml), 0o644); err != nil {
		t.Fatal(err)
	}
}

// callTool runs an MCP tool handler and returns its text content and whether
// it reported an error.
func callTool(t *testing.T, h func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error), args map[string]any) (string, bool) {
	t.Helper()
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	res, err := h(context.Background(), req)
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

// searchNames runs a search and returns the skill names in rank order, which
// ones are flagged project_standard, and the scope note.
func searchNames(t *testing.T, args map[string]any) (names []string, standard map[string]bool, note string) {
	t.Helper()
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	res, err := toolGrimoireSearchSkills(context.Background(), req)
	if err != nil || res.IsError {
		t.Fatalf("search: %v %+v", err, res)
	}
	var hits []mcpSearchSkillsResult
	if err := json.Unmarshal([]byte(res.Content[0].(mcp.TextContent).Text), &hits); err != nil {
		t.Fatal(err)
	}
	standard = map[string]bool{}
	for i := range hits {
		names = append(names, hits[i].Name)
		if hits[i].ProjectStandard {
			standard[hits[i].Name] = true
		}
	}
	if len(res.Content) > 1 {
		note = res.Content[1].(mcp.TextContent).Text
	}
	return names, standard, note
}

const depsWriteUnitTest = "[dependencies]\nwrite-unit-test = \"*\"\n"

// Without a declaration, MCP mode serves the whole library, unranked.
func TestProjectScope_NoManifest(t *testing.T) {
	scopedProject(t, "")
	if s := projectScope(getProjectDir()); s.hasDeclaration() {
		t.Fatalf("no grimoire.toml declared something: %+v", s)
	}
	names, std, note := searchNames(t, map[string]any{"query": "commit"})
	if !slices.Contains(names, "propose-conventional-commit") || len(std) != 0 || note != "" {
		t.Errorf("names=%v standard=%v note=%q", names, std, note)
	}
}

// By default a declaration prefers its practices: the whole library stays
// searchable, declared ones are flagged and lead comparable matches.
func TestProjectScope_Prefer(t *testing.T) {
	scopedProject(t, depsWriteUnitTest)
	s := projectScope(getProjectDir())
	if !s.hasDeclaration() || s.strict || s.size() != 1 || !s.isDeclared("write-unit-test") {
		t.Fatalf("scope = %+v", s)
	}

	names, _, note := searchNames(t, map[string]any{"query": "commit"})
	if !slices.Contains(names, "propose-conventional-commit") {
		t.Errorf("undeclared skill hidden in prefer mode: %v", names)
	}
	if len(names) > 0 && names[0] == "write-unit-test" {
		t.Errorf("irrelevant declared skill outranked the match: %v", names)
	}
	if !strings.Contains(note, "Ranked by this project's grimoire.toml: its 1 declared practices come first") || !strings.Contains(note, `scope "project"`) {
		t.Errorf("note = %q", note)
	}

	names, std, _ := searchNames(t, map[string]any{"query": "writing unit tests"})
	if len(names) == 0 || names[0] != "write-unit-test" || !std["write-unit-test"] {
		t.Errorf("declared skill not first and flagged: %v %v", names, std)
	}

	list, _ := callTool(t, toolGrimoireListSkills, map[string]any{"domain": "engineering", "subdomain": "testing"})
	if !regexp.MustCompile(`"project_standard":\s*true`).MatchString(list) {
		t.Errorf("list lacks the flag: %s", list)
	}

	lowerHookThresholds(t)
	ctx, _ := suggestedSkillContext("write unit tests then commit the change")
	if !strings.Contains(ctx, "write-unit-test (project standard)") || !strings.Contains(ctx, "propose-conventional-commit") {
		t.Errorf("hook candidates: %s", ctx)
	}
	if strings.Index(ctx, "write-unit-test") > strings.Index(ctx, "propose-conventional-commit") {
		t.Errorf("declared candidate not first: %s", ctx)
	}

	// scope "project" asks for the declared practices only; "all" drops the preference.
	if names, _, _ := searchNames(t, map[string]any{"query": "commit", "scope": "project"}); slices.Contains(names, "propose-conventional-commit") {
		t.Errorf(`scope "project" returned an undeclared skill: %v`, names)
	}
	if _, std, note := searchNames(t, map[string]any{"query": "writing unit tests", "scope": "all"}); len(std) != 0 || note != "" {
		t.Errorf(`scope "all" still ranked: %v %q`, std, note)
	}
}

// [standards] scope = "strict" limits MCP mode to the declared practices.
func TestProjectScope_Strict(t *testing.T) {
	scopedProject(t, "[standards]\nscope = \"strict\"\n\n"+depsWriteUnitTest)
	if s := projectScope(getProjectDir()); !s.restricted() {
		t.Fatalf("not strict: %+v", s)
	}
	names, _, note := searchNames(t, map[string]any{"query": "commit"})
	if slices.Contains(names, "propose-conventional-commit") || slices.Contains(names, "suggest-best-practice") {
		t.Errorf("strict search returned undeclared or meta skills: %v", names)
	}
	if !strings.Contains(note, "Limited to this project's 1 declared practices") || !strings.Contains(note, `scope "all"`) {
		t.Errorf("note = %q", note)
	}
	if names, _, _ := searchNames(t, map[string]any{"query": "commit", "scope": "all"}); !slices.Contains(names, "propose-conventional-commit") {
		t.Errorf(`scope "all": %v`, names)
	}
	list, _ := callTool(t, toolGrimoireListSkills, nil)
	if strings.Contains(list, `"writing"`) || !strings.Contains(list, `"engineering"`) {
		t.Errorf("strict index: %s", list)
	}
	lowerHookThresholds(t)
	if ctx, _ := suggestedSkillContext("write unit tests then commit the change"); strings.Contains(ctx, "propose-conventional-commit") || !strings.Contains(ctx, "write-unit-test") {
		t.Errorf("strict hook candidates: %s", ctx)
	}
}

// A skill asked for by name always loads; an undeclared one carries a note
// worded for the mode. Meta workflows never do.
func TestProjectScope_GetSkillByName(t *testing.T) {
	for _, c := range []struct{ toml, want string }{
		{depsWriteUnitTest, "not one of this project's declared practices"},
		{"[standards]\nscope = \"strict\"\n\n" + depsWriteUnitTest, "not in this project's grimoire.toml"},
	} {
		scopedProject(t, c.toml)
		out, isErr := callTool(t, toolGrimoireGetSkill, map[string]any{"name": "propose-conventional-commit"})
		if isErr || !strings.Contains(out, "BODY-MARKER") || !strings.Contains(out, c.want) {
			t.Errorf("want %q: err=%v %s", c.want, isErr, out[:min(len(out), 300)])
		}
		for _, name := range []string{"write-unit-test", "suggest-best-practice"} {
			if out, _ := callTool(t, toolGrimoireGetSkill, map[string]any{"name": name}); strings.Contains(out, "this project's") {
				t.Errorf("%s got a note", name)
			}
		}
		batch, _ := callTool(t, toolGrimoireGetSkill, map[string]any{"names": []any{"write-unit-test", "propose-conventional-commit"}})
		if strings.Count(batch, c.want) != 1 {
			t.Errorf("batch notes: %s", batch[:min(len(batch), 400)])
		}
	}
}

// A dependency package that is not on disk is reported, not fetched.
func TestProjectScope_MissingPackage(t *testing.T) {
	scopedProject(t, depsWriteUnitTest+"\"acme/missing-skills\" = \"*\"\n")
	s := projectScope(getProjectDir())
	if len(s.missing) != 1 || !strings.Contains(s.missing[0], "acme/missing-skills") {
		t.Fatalf("missing = %v", s.missing)
	}
	if st := scopeStatus(displayBrief); !strings.Contains(st, "this project needs") || !strings.Contains(st, "grimoire update") {
		t.Errorf("status = %q", st)
	}
	if _, _, note := searchNames(t, map[string]any{"query": "tests"}); !strings.Contains(note, "Missing packages") {
		t.Errorf("note = %q", note)
	}
}

// The declaration shows in the status line at full only.
func TestScopeStatus_Levels(t *testing.T) {
	scopedProject(t, depsWriteUnitTest)
	if st := scopeStatus(displayBrief); st != "" {
		t.Errorf("brief: %q", st)
	}
	if st := scopeStatus(displayFull); !strings.Contains(st, "project standards: 1 declared practices (preferred)") {
		t.Errorf("full: %q", st)
	}
}

// The cached scope follows edits to grimoire.toml.
func TestProjectScope_CacheInvalidation(t *testing.T) {
	dir := scopedProject(t, depsWriteUnitTest)
	if s := projectScope(dir); s.size() != 1 {
		t.Fatalf("first: %+v", s)
	}
	writeProjectToml(t, dir, depsWriteUnitTest+"propose-conventional-commit = \"*\"\n")
	touchLater(t, filepath.Join(dir, "grimoire.toml"))
	if s := projectScope(dir); s.size() != 2 {
		t.Errorf("after edit: %+v", s)
	}
}

// Dependencies, profile skills and pinned practices together are the
// declaration (a union: each names practices the project uses).
func TestProjectScope_DeclarationSources(t *testing.T) {
	dir := scopedProject(t, "[standards]\nprofiles = [\"my-team\"]\n\n[standards.writing]\npractices = [\"write-readme\"]\n\n"+depsWriteUnitTest)
	pdir := filepath.Join(dir, ".grimoire", "profiles")
	if err := os.MkdirAll(pdir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pdir, "my-team.toml"), []byte("name = \"my-team\"\n\n[[skills]]\nname = \"propose-conventional-commit\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	touchLater(t, filepath.Join(dir, "grimoire.toml"))
	s := projectScope(dir)
	for _, name := range []string{"write-unit-test", "propose-conventional-commit", "write-readme"} {
		if !s.isDeclared(name) {
			t.Errorf("%s not declared: %+v", name, s.declared)
		}
	}
}

// touchLater moves a file's mtime forward so the scope cache sees a change.
func touchLater(t *testing.T, path string) {
	t.Helper()
	later := time.Now().Add(2 * time.Second)
	if err := os.Chtimes(path, later, later); err != nil {
		t.Fatal(err)
	}
}

// lowerHookThresholds lets the small test package produce keyword
// candidates, which the real thresholds tune for a 1,000+ skill library.
func lowerHookThresholds(t *testing.T) {
	t.Helper()
	score, matched, rel := hookMinScore, hookMinMatched, hookRelScore
	t.Cleanup(func() { hookMinScore, hookMinMatched, hookRelScore = score, matched, rel })
	hookMinScore, hookMinMatched, hookRelScore = 0, 1, 0
}
