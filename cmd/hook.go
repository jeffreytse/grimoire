package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/search"
	"github.com/jeffreytse/grimoire/internal/skills"
)

// `grimoire hook prompt` is a Claude Code UserPromptSubmit hook. Claude Code
// pipes the hook input JSON to stdin; whatever this prints as
// hookSpecificOutput.additionalContext is added to the model's context for
// that prompt. Unlike the CLAUDE.md instruction, the hook runs on every
// prompt, so the skill lookup cannot be skipped by the model.
//
// The hook must never block or break a prompt: every failure exits 0 with no
// output.

const (
	// hookMaxQuery caps how much of the prompt is searched; long pastes
	// dilute BM25 and slow the hook.
	hookMaxQuery = 2000
	// hookTopN is how many matching skills are suggested.
	hookTopN = 3
	// hookMaxDesc trims each suggested skill's description.
	hookMaxDesc = 200
)

// hookInlineBudget caps inlined skill context below Claude Code's hook
// context limit (10,000 chars), past which the text is moved to a file.
var hookInlineBudget = 9000

// hookMinScore is the BM25 score a skill needs to be suggested. Below it the
// hook stays silent so small talk ("ok", "thanks") costs no tokens.
var hookMinScore = 8.0

// hookMinMatched is how many distinct prompt terms a skill must contain.
var hookMinMatched = 2

// hookRelScore drops suggestions scoring below this fraction of the best one.
var hookRelScore = 0.5

var hookCmd = &cobra.Command{
	Use:    "hook",
	Short:  "Agent hooks that put grimoire skills in front of the model",
	Hidden: true,
}

var hookPromptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Claude Code UserPromptSubmit hook: suggest skills matching the prompt",
	Long: `Reads Claude Code's UserPromptSubmit hook input on stdin and prints
additional context for the model:

  - grimoire:<name> in the prompt  the named skill's SKILL.md, to follow
  - otherwise                      the best-matching skills (BM25), with a
                                   note to load one with grimoire_get_skill

Prints nothing when no skill matches well. Registered by: grimoire mcp install`,
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		runHookPrompt(os.Stdin, os.Stdout)
		return nil
	},
}

func init() {
	hookPromptCmd.Flags().Float64Var(&hookMinScore, "min-score", hookMinScore, "minimum BM25 score for a skill to be suggested")
	hookCmd.AddCommand(hookPromptCmd)
	rootCmd.AddCommand(hookCmd)
}

type hookInput struct {
	Prompt string `json:"prompt"`
}

type hookOutput struct {
	HookSpecificOutput struct {
		HookEventName     string `json:"hookEventName"`
		AdditionalContext string `json:"additionalContext"`
	} `json:"hookSpecificOutput"`
}

func runHookPrompt(in io.Reader, out io.Writer) {
	raw, err := io.ReadAll(io.LimitReader(in, 1<<20))
	if err != nil {
		return
	}
	var input hookInput
	if json.Unmarshal(raw, &input) != nil {
		return
	}
	ctx := hookContext(input.Prompt)
	if ctx == "" {
		return
	}
	var o hookOutput
	o.HookSpecificOutput.HookEventName = "UserPromptSubmit"
	o.HookSpecificOutput.AdditionalContext = ctx
	b, err := json.Marshal(o)
	if err != nil {
		return
	}
	fmt.Fprintln(out, string(b))
}

var skillRefPattern = regexp.MustCompile(`(?:^|[\s(\[` + "`" + `])grimoire:([A-Za-z0-9][A-Za-z0-9._/-]*)`)

// hookContext returns the context to add for prompt, or "" for none.
func hookContext(prompt string) string {
	prompt = strings.TrimSpace(prompt)
	// Slash commands are handled by Claude Code (including /grimoire:skill).
	if prompt == "" || strings.HasPrefix(prompt, "/") {
		return ""
	}
	if refs := explicitSkillRefs(prompt); len(refs) > 0 {
		return explicitSkillContext(refs)
	}
	return suggestedSkillContext(prompt)
}

// explicitSkillRefs returns the distinct skill names written as grimoire:<name>.
func explicitSkillRefs(prompt string) []string {
	var names []string
	seen := map[string]bool{}
	for _, m := range skillRefPattern.FindAllStringSubmatch(prompt, -1) {
		n := strings.TrimRight(m[1], ".,;:!?/")
		if n != "" && !seen[n] {
			seen[n] = true
			names = append(names, n)
		}
	}
	return names
}

// explicitSkillContext inlines each named skill so the model follows it
// without a tool call; an unknown name gets the "did you mean" hint instead.
//
// Claude Code moves hook context over its size cap into a file the model must
// Read back, which costs a turn. A skill that would push the context past
// hookInlineBudget is therefore named for grimoire_get_skill instead of
// inlined: one tool call either way, and the server attaches the chain rule.
func explicitSkillContext(names []string) string {
	var b strings.Builder
	resolved := false
	for _, ref := range names {
		name, path, err := resolveSkillOrSuggest(ref)
		if err != nil {
			msg := err.Error()
			if !strings.HasSuffix(msg, "?") {
				msg += "."
			}
			fmt.Fprintf(&b, "grimoire: %s Tell the user, or load the right skill with grimoire_get_skill.\n\n", msg)
			continue
		}
		content, err := os.ReadFile(filepath.Join(path, "SKILL.md"))
		if err != nil {
			continue
		}
		if !resolved {
			resolved = true
			b.WriteString(mcpChainRule + "\n\n")
		}
		inline := fmt.Sprintf("The user invoked grimoire skill %q. Follow it for this prompt, applied to the rest of the user's message.\n\n%s\n\n", name, strings.TrimSpace(string(content)))
		if b.Len()+len(inline) <= hookInlineBudget {
			b.WriteString(inline)
			continue
		}
		fmt.Fprintf(&b, "The user invoked grimoire skill %q. Before anything else, call grimoire_get_skill with name %q and follow it for this prompt, applied to the rest of the user's message.\n\n", name, name)
	}
	return strings.TrimSpace(b.String())
}

// suggestedSkillContext lists the skills that best match prompt, or "" when
// none scores at least hookMinScore.
func suggestedSkillContext(prompt string) string {
	regs := skills.AllSkillsPackages()
	if len(regs) == 0 {
		return ""
	}
	all, _, err := skills.ListAllSkillsFromPackages(regs)
	if err != nil {
		return ""
	}
	picks := hookPicks(prompt, all)
	if len(picks) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("grimoire: best-practice skills matching this prompt:\n")
	for i := range picks {
		sk := &picks[i].Skill
		desc := strings.Join(strings.Fields(sk.Description), " ")
		if r := []rune(desc); len(r) > hookMaxDesc {
			desc = string(r[:hookMaxDesc]) + "…"
		}
		if desc == "" {
			fmt.Fprintf(&b, "- %s\n", sk.Name)
		} else {
			fmt.Fprintf(&b, "- %s — %s\n", sk.Name, desc)
		}
	}
	b.WriteString("Before acting, if one fits the task, call grimoire_get_skill with its name and follow it. If none fits, ignore this note.")
	return b.String()
}

// hookPicks returns the skills worth suggesting for prompt, best first.
func hookPicks(prompt string, all []skills.Skill) []search.Result {
	query := prompt
	if len(query) > hookMaxQuery {
		query = query[:hookMaxQuery]
	}
	var picks []search.Result
	hits := search.BM25Search(query, all, hookTopN)
	for i := range hits {
		h := &hits[i]
		// One shared word ("time", "plan") is weak evidence unless the
		// prompt has only one topic word to begin with.
		if h.Matched < min(hookMinMatched, h.Terms) {
			continue
		}
		if h.Score < hookMinScore || h.Score < hookRelScore*hits[0].Score {
			continue
		}
		picks = append(picks, *h)
	}
	return picks
}
