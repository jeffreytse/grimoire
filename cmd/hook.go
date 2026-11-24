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

// `grimoire hook prompt` is a prompt hook for Claude Code and Codex
// (UserPromptSubmit) and Gemini CLI (BeforeAgent). The agent pipes the hook
// input JSON to stdin; whatever this prints as
// hookSpecificOutput.additionalContext is added to the model's context for
// that prompt. Unlike the instruction file, the hook runs on every prompt, so
// the skill lookup cannot be skipped by the model.
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

// hookInlineBudget caps inlined skill context below the agents' hook context
// limits (Claude Code: 10,000 chars; Codex: 2,500 tokens by default), past
// which the text is moved to a file.
var hookInlineBudget = 9000

// hookMinScore is the BM25F score a skill needs to be suggested. Below it the
// hook lists no candidates. Tuned with TestHookEval (cmd/testdata/hook_eval.json);
// retune there if the scoring changes.
var hookMinScore = 9.0

// hookRouteReminder is the per-prompt nudge onto grimoire's router. The
// router itself is injected once per session by `grimoire hook session`;
// re-asserting it on every prompt is what keeps the model from judging a task
// "simple" and answering from its own knowledge. In skills mode the
// instruction runs start-best-practice on every task; this is MCP mode's
// equivalent.
const hookRouteReminder = "grimoire: before responding, apply the start-best-practice routing (injected at session start; if it is not in your context, load it with grimoire_get_skill). If it routes to a workflow, follow that skill: it asks clarifying questions when needed and loads the specific best practices. If nothing routes, proceed normally."

// routerSkill is grimoire's entry-point workflow router.
const routerSkill = "start-best-practice"

// routerSections are the parts of the router SKILL.md the session hook
// inlines: the routing table and its rules. The rationale and examples are
// left to grimoire_get_skill.
var routerSections = []string{"Steps", "Rules", "Display"}

// hookSmallTalk are replies that never start a new task, so the hook stays
// silent for them (start-best-practice passes them through anyway).
var hookSmallTalk = map[string]bool{
	"ok": true, "okay": true, "k": true, "thanks": true, "thank you": true, "thx": true, "ty": true,
	"yes": true, "y": true, "no": true, "n": true, "yep": true, "nope": true, "sure": true,
	"continue": true, "go on": true, "go ahead": true, "proceed": true, "next": true, "done": true,
	"looks good": true, "lgtm": true, "got it": true, "great": true, "nice": true, "cool": true,
	"perfect": true, "good": true, "do it": true, "ship it": true, "sounds good": true,
}

// hookMinWords is the prompt length (words) below which a prompt is treated
// as a follow-up rather than a new request.
var hookMinWords = 3

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
	Short: "Agent prompt hook: suggest skills matching the prompt",
	Long: `Reads a prompt hook's input on stdin (Claude Code and Codex
UserPromptSubmit, Gemini CLI BeforeAgent) and prints
additional context for the model:

  - grimoire:<name> in the prompt  the named skill's SKILL.md, to follow
  - otherwise                      a reminder to apply the start-best-practice
                                   routing, plus candidate skills (BM25)

Prints nothing for slash commands and small talk ("ok", "thanks").
Registered by: grimoire mcp install`,
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		runHookPrompt(os.Stdin, os.Stdout)
		return nil
	},
}

var hookSessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Agent session-start hook: inject grimoire's workflow router",
	Long: `Reads a session-start hook's input on stdin (Claude Code, Codex and
Gemini CLI SessionStart) and prints additional context for the model: the
start-best-practice routing table and rules, read from the installed skills
package, so every request in the session is routed to the grimoire workflow
that fits. Registered by: grimoire mcp install`,
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		runHookSession(os.Stdin, os.Stdout)
		return nil
	},
}

func init() {
	hookPromptCmd.Flags().Float64Var(&hookMinScore, "min-score", hookMinScore, "minimum BM25 score for a skill to be suggested")
	hookCmd.AddCommand(hookPromptCmd, hookSessionCmd)
	rootCmd.AddCommand(hookCmd)
}

type hookInput struct {
	Prompt string `json:"prompt"`
	// HookEventName is echoed back: Claude Code and Codex send
	// "UserPromptSubmit", Gemini CLI sends "BeforeAgent".
	HookEventName string `json:"hook_event_name"`
}

type hookOutput struct {
	// SystemMessage is shown to the user, not the model: grimoire's status
	// line ("grimoire ▸ …"). Empty when trace is off.
	SystemMessage      string `json:"systemMessage,omitempty"`
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
	l := displayLevelSetting()
	ctx, msg := hookTurn(input.Prompt, l)
	writeHookOutput(out, input.HookEventName, "UserPromptSubmit", ctx, withOnboardingTip(l, ctx, msg))
}

// runHookSession prints the router context for a session-start hook.
func runHookSession(in io.Reader, out io.Writer) {
	raw, err := io.ReadAll(io.LimitReader(in, 1<<20))
	if err != nil {
		return
	}
	var input hookInput
	// Some agents send no input to session hooks; that is fine.
	_ = json.Unmarshal(raw, &input)
	ctx, msg := sessionTurn(displayLevelSetting())
	writeHookOutput(out, input.HookEventName, "SessionStart", ctx, msg)
}

// writeHookOutput prints ctx as the hook's additional context and msg as the
// status line shown to the user, echoing the agent's event name (or fallback
// when it sent none). Empty ctx prints nothing.
func writeHookOutput(out io.Writer, event, fallback, ctx, msg string) {
	if ctx == "" {
		return
	}
	var o hookOutput
	o.SystemMessage = msg
	o.HookSpecificOutput.HookEventName = event
	if event == "" {
		o.HookSpecificOutput.HookEventName = fallback
	}
	o.HookSpecificOutput.AdditionalContext = ctx
	b, err := json.Marshal(o)
	if err != nil {
		return
	}
	fmt.Fprintln(out, string(b))
}

// sessionContext is the router context for a session at the brief level.
func sessionContext() string {
	ctx, _ := sessionTurn(displayBrief)
	return ctx
}

// sessionTurn returns the router injected once per session — the chain rule
// and the start-best-practice routing, inlined from the installed package so
// it stays in step with the skills library — and the status line for the
// user. Without the package it falls back to a directive to load the router.
//
// The status line states facts: the router and display level in effect, or
// what is wrong when the skills package is missing.
func sessionTurn(l displayLevel) (ctx, msg string) {
	head := "grimoire is active for this session. Before responding to each new user request, apply the start-best-practice routing below and follow the workflow it routes to; load routed skills with grimoire_get_skill.\n\n" + mcpChainRule + "\n\n" + displayContext(l)
	fallback := head + "\n\nThe router is not inlined here: call grimoire_get_skill with name \"" + routerSkill + "\" once, at the first request of the session, and follow it."
	// Not a problem — the model loads the router itself — so full only.
	fallbackMsg := ""
	if l == displayFull {
		fallbackMsg = statusLine(l, "router not inlined · the model loads "+routerSkill+" on demand")
	}
	_, path, err := resolveSkillOrSuggest(routerSkill)
	if err != nil {
		return fallback, statusLine(l, "no skills package with "+routerSkill+" — run `grimoire update`")
	}
	raw, err := os.ReadFile(filepath.Join(path, "SKILL.md"))
	if err != nil {
		return fallback, fallbackMsg
	}
	body := markdownSections(string(raw), routerSections...)
	if body == "" {
		return fallback, fallbackMsg
	}
	ctx = head + "\n\n# start-best-practice (router)\n\n" + body
	if len(ctx) > hookInlineBudget {
		return fallback, fallbackMsg
	}
	n := countRoutes(body)
	unit := "routes"
	if n == 1 {
		unit = "route"
	}
	return ctx, readyStatus(l, "", fmt.Sprintf("router: %s, %d %s", routerSkill, n, unit))
}

// countRoutes counts the routing table rows that invoke a workflow.
func countRoutes(body string) int {
	n := 0
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "|") && strings.Contains(line, "invoke `") {
			n++
		}
	}
	return n
}

// markdownSections returns the level-2 sections of md whose headings are in
// names, in document order, each with its heading.
func markdownSections(md string, names ...string) string {
	want := map[string]bool{}
	for _, n := range names {
		want[n] = true
	}
	var b strings.Builder
	keep, inFence := false, false
	for _, line := range strings.Split(md, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inFence = !inFence
		}
		if !inFence && strings.HasPrefix(line, "## ") {
			keep = want[strings.TrimSpace(strings.TrimPrefix(line, "## "))]
		}
		if keep {
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
	return strings.TrimSpace(b.String())
}

// isSmallTalk reports whether prompt is an acknowledgement or short
// follow-up rather than a new request.
func isSmallTalk(prompt string) bool {
	norm := strings.ToLower(strings.Trim(strings.TrimSpace(prompt), ".!?,;: "))
	return hookSmallTalk[norm] || len(strings.Fields(norm)) < hookMinWords
}

var skillRefPattern = regexp.MustCompile(`(?:^|[\s(\[` + "`" + `])grimoire:([A-Za-z0-9][A-Za-z0-9._/-]*)`)

// hookContext returns the context to add for prompt (brief level), or "" for none.
func hookContext(prompt string) string {
	ctx, _ := hookTurn(prompt, displayBrief)
	return ctx
}

// hookTurn returns the context to add for prompt and the status line for the
// user, or "" for none. The status line states only what the hook knows: a
// skill it loaded, or (full level) the unverified keyword candidates it
// passed to the model. It never claims a route — the model shows that.
func hookTurn(prompt string, l displayLevel) (ctx, msg string) {
	prompt = strings.TrimSpace(prompt)
	// Slash commands are handled by the agent (including /grimoire:skill).
	if prompt == "" || strings.HasPrefix(prompt, "/") {
		return "", ""
	}
	if refs := explicitSkillRefs(prompt); len(refs) > 0 {
		ctx, loaded, unknown := explicitSkillContext(refs)
		var parts []string
		if len(loaded) > 0 {
			parts = append(parts, "loading skill "+strings.Join(loaded, ", "))
		}
		if len(unknown) > 0 {
			parts = append(parts, "unknown skill "+strings.Join(unknown, ", "))
		}
		return ctx, statusLine(l, strings.Join(parts, " · "))
	}
	if isSmallTalk(prompt) {
		return "", ""
	}
	reminder := hookRouteReminder + " " + displayContext(l)
	if cands, names := suggestedSkillContext(prompt); cands != "" {
		return reminder + "\n" + cands, candidatesStatus(l, names)
	}
	return reminder, ""
}

// candidatesStatus shows the keyword candidates at the full level only,
// labeled as unverified: at brief they are noise, and often off-topic.
func candidatesStatus(l displayLevel, names []string) string {
	if l != displayFull {
		return ""
	}
	return statusLine(l, "keyword candidates (unverified): "+strings.Join(names, ", "))
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
func explicitSkillContext(names []string) (ctx string, loaded, unknown []string) {
	var b strings.Builder
	resolved := false
	for _, ref := range names {
		name, path, err := resolveSkillOrSuggest(ref)
		if err != nil {
			unknown = append(unknown, ref)
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
		loaded = append(loaded, name)
		inline := fmt.Sprintf("The user invoked grimoire skill %q. Follow it for this prompt, applied to the rest of the user's message.\n\n%s\n\n", name, strings.TrimSpace(string(content)))
		if b.Len()+len(inline) <= hookInlineBudget {
			b.WriteString(inline)
			continue
		}
		fmt.Fprintf(&b, "The user invoked grimoire skill %q. Before anything else, call grimoire_get_skill with name %q and follow it for this prompt, applied to the rest of the user's message.\n\n", name, name)
	}
	return strings.TrimSpace(b.String()), loaded, unknown
}

// suggestedSkillContext lists the skills that best match prompt as
// candidates for the routed workflow, or "" when none scores at least
// hookMinScore. Keyword matches can be off-topic, so they are offered as
// unverified rather than as the answer.
func suggestedSkillContext(prompt string) (ctx string, names []string) {
	regs := skills.AllSkillsPackages()
	if len(regs) == 0 {
		return "", nil
	}
	all, _, err := skills.ListAllSkillsFromPackages(regs)
	if err != nil {
		return "", nil
	}
	picks := hookPicks(prompt, all)
	if len(picks) == 0 {
		return "", nil
	}

	var b strings.Builder
	b.WriteString("grimoire: candidate skills by keyword match (unverified; the routed workflow confirms or discards them):\n")
	for i := range picks {
		sk := &picks[i].Skill
		names = append(names, sk.Name)
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
	return strings.TrimRight(b.String(), "\n"), names
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
