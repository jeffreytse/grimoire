package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/jeffreytse/grimoire/internal/search"
	"github.com/jeffreytse/grimoire/internal/skills"
)

// Calling an exact skill in MCP mode, in three layers (most to least
// deterministic):
//
//  1. The `skill` MCP prompt — user-controlled per the MCP spec; Claude Code
//     shows it as `/grimoire:skill <name>` and injects the skill directly.
//     One prompt with a completable argument, not one prompt per skill.
//  2. Server instructions documenting the `grimoire:<name>` text convention,
//     for clients without prompt support (and `claude -p`).
//  3. resolveSkillOrSuggest: strips `grimoire:`, and answers a miss with
//     "did you mean" suggestions.

// skillRefPrefix is the text-convention prefix users type before a skill name.
const skillRefPrefix = "grimoire:"

// mcpServerInstructions is sent in the initialize result. Clients cut long
// instructions (Claude Code: 4,096 chars), so the critical rule comes first.
const mcpServerInstructions = `grimoire serves professional best-practice skills.
- If the user writes ` + "`grimoire:<skill-name>`" + ` (e.g. ` + "`grimoire:propose-conventional-commit`" + `), call grimoire_get_skill with exactly that name — do not search first — and follow the skill, applying it to the rest of the user's message.
- For other tasks, call grimoire_search_skills to find a relevant skill, then grimoire_get_skill to load it and follow it.
- If a skill name is not found, the error lists close matches; use the right one.
- Use grimoire_list_skills (no arguments) for a domain index.`

// skillCompletionLimit caps completion values (the spec allows up to 100).
const skillCompletionLimit = 20

// mcpServerOptions returns the server options for `grimoire mcp serve`.
func mcpServerOptions() []server.ServerOption {
	return []server.ServerOption{
		server.WithInstructions(mcpServerInstructions),
		server.WithPromptCapabilities(false),
		server.WithCompletions(),
		server.WithPromptCompletionProvider(skillNameCompleter{}),
	}
}

func registerMCPPrompts(s *server.MCPServer) {
	s.AddPrompt(
		mcp.NewPrompt("skill",
			mcp.WithPromptDescription("Load one grimoire skill by exact name and follow it for the current task (e.g. propose-conventional-commit)."),
			mcp.WithArgument("name",
				mcp.ArgumentDescription("Skill name, e.g. propose-conventional-commit"),
				mcp.RequiredArgument(),
			),
		),
		promptGrimoireSkill,
	)
}

func promptGrimoireSkill(_ context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) { //nolint:gocritic
	name, path, err := resolveSkillOrSuggest(request.Params.Arguments["name"])
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(filepath.Join(path, "SKILL.md"))
	if err != nil {
		return nil, fmt.Errorf("reading SKILL.md for %q: %w", name, err)
	}
	text := fmt.Sprintf("Follow this grimoire skill (%s) for the current task.\n\n%s", name, content)
	return mcp.NewGetPromptResult(
		"grimoire skill: "+name,
		[]mcp.PromptMessage{mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(text))},
	), nil
}

// resolveSkillOrSuggest resolves a user-supplied skill reference to its
// directory. It strips a leading "grimoire:" and, when nothing matches,
// returns an error naming the closest skills by BM25 relevance.
func resolveSkillOrSuggest(ref string) (name, path string, err error) {
	name = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(ref), skillRefPrefix))
	if name == "" {
		return "", "", fmt.Errorf("skill name is required")
	}
	regs := skills.AllSkillsPackages()
	if len(regs) == 0 {
		return name, "", fmt.Errorf("no packages installed — run grimoire_update first")
	}
	if p, _, rerr := resolveSkillFromPackages(regs, name); rerr == nil {
		return name, p, nil
	}
	msg := fmt.Sprintf("skill %q not found", name)
	if all, _, lerr := skills.ListAllSkillsFromPackages(regs); lerr == nil {
		query := strings.NewReplacer("-", " ", "_", " ", "/", " ").Replace(name)
		hits := search.BM25Search(query, all, 3)
		if len(hits) > 0 {
			names := make([]string, len(hits))
			for i := range hits {
				names[i] = hits[i].Skill.Name
			}
			msg += " — did you mean: " + strings.Join(names, ", ") + "?"
		}
	}
	return name, "", fmt.Errorf("%s", msg)
}

// skillNameCompleter completes the `name` argument of the `skill` prompt:
// prefix matches first, then substring matches, capped at skillCompletionLimit.
type skillNameCompleter struct{}

func (skillNameCompleter) CompletePromptArgument(_ context.Context, promptName string, arg mcp.CompleteArgument, _ mcp.CompleteContext) (*mcp.Completion, error) { //nolint:gocritic
	if promptName != "skill" || arg.Name != "name" {
		return &mcp.Completion{Values: []string{}}, nil
	}
	all, _, err := skills.ListAllSkillsFromPackagesMeta(skills.AllSkillsPackages())
	if err != nil {
		return nil, err
	}
	typed := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(arg.Value), skillRefPrefix))
	var prefix, contains []string
	seen := map[string]bool{}
	for i := range all {
		n := all[i].Name
		if seen[n] {
			continue
		}
		seen[n] = true
		switch ln := strings.ToLower(n); {
		case strings.HasPrefix(ln, typed):
			prefix = append(prefix, n)
		case strings.Contains(ln, typed):
			contains = append(contains, n)
		}
	}
	sort.Strings(prefix)
	sort.Strings(contains)
	matches := append(prefix, contains...) //nolint:gocritic // prefix is a fresh slice
	total := len(matches)
	if total > skillCompletionLimit {
		matches = matches[:skillCompletionLimit]
	}
	return &mcp.Completion{Values: matches, Total: total, HasMore: total > skillCompletionLimit}, nil
}
