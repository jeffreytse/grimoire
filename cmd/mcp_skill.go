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

// mcpChainRule maps the skills-mode wording inside grimoire skills to MCP
// tools. Workflow skills (suggest-best-practice, plan-best-practice-solution,
// apply-best-practice-tree, start-best-practice, …) were written for skills
// linked into the agent: they "invoke" other skills, "load the skill using the
// Skill tool", and "score all installed skills" from a catalog already in
// context. In MCP mode none of that exists, so this rule travels with every
// skill grimoire hands out.
const mcpChainRule = `grimoire skills run over MCP here; no grimoire skills are installed in the agent itself. When a grimoire skill says to invoke, load, run or apply another skill — "load the skill using the Skill tool", "invoke ` + "`X`" + `", "/X" — call grimoire_get_skill with that skill name and follow it. When it says to scan, score or list installed skills, get candidates with grimoire_search_skills (or grimoire_list_skills for a domain) and only name skills those tools return. Install hints like "/plugin install" mean grimoire_package_add.`

// mcpServerInstructions is sent in the initialize result. Clients cut long
// instructions (Claude Code: 4,096 chars), so the critical rule comes first.
const mcpServerInstructions = `grimoire serves professional best-practice skills.
- If the user writes ` + "`grimoire:<skill-name>`" + ` (e.g. ` + "`grimoire:propose-conventional-commit`" + `), call grimoire_get_skill with exactly that name — do not search first — and follow the skill, applying it to the rest of the user's message.
- For other tasks, call grimoire_search_skills to find a relevant skill, then grimoire_get_skill to load it and follow it.
- For a multi-step or unclear problem, load start-best-practice first. It routes to the workflow that fits: suggest-best-practice (one skill), plan-best-practice-solution (several skills in order), apply-best-practice-tree (break an opaque problem down), discover-best-practices (browse a domain), learn-best-practice (retain a practice).
- To load several skills in one call, pass grimoire_get_skill a names list.
- If a skill name is not found, the error lists close matches; use the right one.
- Use grimoire_list_skills (no arguments) for a domain index.
- ` + mcpChainRule

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
	text := fmt.Sprintf("Follow this grimoire skill (%s) for the current task. %s\n\n%s", name, mcpChainRule, content)
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
