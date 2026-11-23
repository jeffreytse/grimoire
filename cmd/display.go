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

// displayOffRule tells the model to keep grimoire silent.
const displayOffRule = "grimoire display is off: do not show grimoire blocks or `grimoire ▸` lines; route and load skills silently."

// displayContext is the display instruction for the model at level l. When
// the installed router skill defines the format (its Display section), this
// only names the level. Older skill packages get a one-line-per-decision
// fallback rather than a second copy of the format.
func displayContext(l displayLevel) string {
	switch {
	case !l.on():
		return displayOffRule
	case routerHasDisplay():
		return fmt.Sprintf("grimoire display: %s. Show grimoire's work in the format of start-best-practice's Display section, at the %s level. The user can say \"show grimoire details\" (full) or \"quiet grimoire\" (off) for the rest of the session.", l, l)
	}
	return fmt.Sprintf("grimoire display: %s. Tell the user each grimoire decision in one line starting with `grimoire ▸`: the workflow routed to and why, then the skills applied and why. Use the prefix only on those lines, never on the answer. Show nothing for grimoire when nothing routes.", l)
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
