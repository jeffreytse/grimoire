package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeffreytse/grimoire/internal/config"
)

// grimoire's display: how much of its work grimoire shows (progressive
// disclosure). The format itself is defined once, in the start-best-practice
// skill's Display section; the CLI only says which level applies and shows
// facts it knows in status lines.
//
//	off    nothing: route and load skills silently
//	brief  (default) the decisions: route, plan, each skill applied, outcome
//	full   brief plus the evidence: candidates and why, severity, sources
//
// Set with core.trace (grimoire config set core.trace brief) or, for one run,
// GRIMOIRE_TRACE. "on" is read as brief for configs written before levels.

type displayLevel string

const (
	displayOff   displayLevel = "off"
	displayBrief displayLevel = "brief"
	displayFull  displayLevel = "full"
)

// on reports whether grimoire shows anything at this level.
func (l displayLevel) on() bool { return l != displayOff }

// parseDisplayLevel reads a core.trace / GRIMOIRE_TRACE value; ok is false
// for anything that is not a level.
func parseDisplayLevel(v string) (l displayLevel, ok bool) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "off", "0", "false":
		return displayOff, true
	case "brief", "on", "1", "true":
		return displayBrief, true
	case "full":
		return displayFull, true
	}
	return displayBrief, false
}

// displayLevelSetting is the level in effect: GRIMOIRE_TRACE, then
// core.trace, then brief. A config that cannot be read counts as brief.
func displayLevelSetting() displayLevel {
	if v := os.Getenv("GRIMOIRE_TRACE"); strings.TrimSpace(v) != "" {
		l, _ := parseDisplayLevel(v)
		return l
	}
	r, err := config.Load(getProjectDir())
	if err != nil || r.Core.Trace == "" {
		return displayBrief
	}
	l, _ := parseDisplayLevel(r.Core.Trace)
	return l
}

// displayPlain reports whether grimoire's display uses words instead of
// decorative glyphs (⎿ ★ ✓ and emoji), which screen readers announce by name:
// GRIMOIRE_DISPLAY_STYLE, then core.display-style, and plain whenever
// NO_COLOR is set.
func displayPlain() bool {
	if v := strings.ToLower(strings.TrimSpace(os.Getenv("GRIMOIRE_DISPLAY_STYLE"))); v != "" {
		return v == "plain"
	}
	if os.Getenv("NO_COLOR") != "" {
		return true
	}
	r, err := config.Load(getProjectDir())
	return err == nil && r.Core.DisplayStyle == "plain"
}

// displayPlainRule asks for the plain style.
const displayPlainRule = " Style: plain — no decorative glyphs or emoji; use words instead (\"Step 1:\", \"Risk (high):\", \"Done:\")."

// displayOffRule tells the model to keep grimoire silent.
const displayOffRule = "grimoire display is off: do not show grimoire blocks or `grimoire ▸` lines; route and load skills silently."

// displayContext is the display instruction for the model at level l. When
// the installed router skill defines the format (its Display section), this
// only names the level. Older skill packages get a one-line-per-decision
// fallback rather than a second copy of the format.
func displayContext(l displayLevel) string {
	if !l.on() {
		return displayOffRule
	}
	style := ""
	if displayPlain() {
		style = displayPlainRule
	}
	if routerHasDisplay() {
		return fmt.Sprintf("grimoire display: %s. Show grimoire's work in the format of start-best-practice's Display section, at the %s level. The user can say \"show grimoire details\" (full) or \"quiet grimoire\" (off) for the rest of the session.%s", l, l, style)
	}
	return fmt.Sprintf("grimoire display: %s. Tell the user each grimoire decision in one short line starting with `grimoire ▸`: the workflow routed to and why, then the skills applied and why. Use the prefix only on those lines, never on the answer. Show nothing for grimoire when nothing routes.%s", l, style)
}

// readyStatus is the session-start status line at the full level: that
// grimoire is ready, the view in effect, how to change it, and detail. At
// brief it is empty — silent on success. Agents render hook messages as
// warnings (Claude Code's systemMessage, Codex's Warning event), so a routine
// "ready" there would wear the warning role and dull real warnings; grimoire
// shows it is on through the badge and its first route line instead.
func readyStatus(l displayLevel, mode, detail string) string {
	ready := "ready"
	if mode != "" {
		ready += " (" + mode + ")"
	}
	if l != displayFull {
		return ""
	}
	line := ready + " · full view — say \"quiet grimoire\" to hide"
	if detail != "" {
		line += " · " + detail
	}
	return statusLine(l, line)
}

// routerHasDisplay reports whether the installed start-best-practice skill
// defines grimoire's display format. Read on each call so a skills update
// applies without restarting the MCP server.
func routerHasDisplay() bool {
	_, path, err := resolveSkillOrSuggest(routerSkill)
	if err != nil {
		return false
	}
	raw, err := os.ReadFile(filepath.Join(path, "SKILL.md"))
	return err == nil && markdownSections(string(raw), "Display") != ""
}

// statusPrefix starts every grimoire status line.
const statusPrefix = "grimoire ▸ "

// statusLine is a status line for the user: facts grimoire knows, never a
// guess about what the model will do. Empty when display is off or there is
// nothing to say.
func statusLine(l displayLevel, text string) string {
	if !l.on() || text == "" {
		return ""
	}
	return statusPrefix + text
}

// onboardingTip is shown once, on the first request grimoire acts on after
// install: the controls a user would otherwise only find in the README. It
// is a contextual reveal — taught when grimoire first does something — and,
// shown once, it does not wear down the warning role the way a per-session
// line would. Recorded under the user cache dir; a cleared cache shows it
// again, which is harmless.
const onboardingTip = `tip: say "show grimoire details" for evidence, "quiet grimoire" to hide it, "skip grimoire" to skip one request`

// takeOnboardingTip returns the tip status line the first time it is called
// (per machine) at a visible level, and "" afterwards.
func takeOnboardingTip(l displayLevel) string {
	if !l.on() {
		return ""
	}
	dir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	marker := filepath.Join(dir, "grimoire", "tip-shown")
	if _, err := os.Lstat(marker); err == nil {
		return ""
	}
	if os.MkdirAll(filepath.Dir(marker), 0o700) != nil || os.WriteFile(marker, nil, 0o600) != nil {
		return "" // cannot remember it: stay quiet rather than repeat it
	}
	return statusLine(l, onboardingTip)
}

// withOnboardingTip adds the one-time tip to a prompt hook's output when the
// hook is adding context (grimoire is acting on the request).
func withOnboardingTip(l displayLevel, ctx, msg string) string {
	if ctx == "" {
		return msg
	}
	tip := takeOnboardingTip(l)
	switch {
	case tip == "":
		return msg
	case msg == "":
		return tip
	}
	return msg + " · " + strings.TrimPrefix(tip, statusPrefix)
}

// preferencesContext summarizes the pinned practice preferences from the
// merged grimoire.toml layers, so workflows that check preferences first
// (suggest-best-practice, review-best-practice-fit …) need no tool call.
// In skills mode that call would be a shell command needing the user's
// permission; in MCP mode it saves a grimoire_context round trip.
func preferencesContext() string {
	r, err := config.Load(getProjectDir())
	if err != nil {
		return ""
	}
	var pinned []string
	for _, key := range r.SectionKeys() {
		if ps := r.ResolveSection(key).Practices; len(ps) > 0 {
			pinned = append(pinned, key+": "+strings.Join(ps, ", "))
		}
	}
	if len(pinned) == 0 {
		return "grimoire preferences (already resolved from grimoire.toml): no practices pinned — skip the pinned-preference check; do not read config files for it."
	}
	return "grimoire preferences (already resolved from grimoire.toml; use these, do not read config files): " + strings.Join(pinned, "; ") + "."
}
