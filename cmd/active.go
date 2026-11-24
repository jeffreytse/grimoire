package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/config"
	"github.com/jeffreytse/grimoire/internal/tui"
)

// Whether grimoire is active: routing every request through
// start-best-practice (active, the default), or staying out of the way
// (inactive) — no router, no routing reminders; a skill named with
// grimoire:<name> still loads. This is not the display level: display off
// hides grimoire's output but it still routes.
//
// Precedence, first set wins:
//
//	GRIMOIRE_ACTIVE=0|1          one run (CI, scripts)
//	session state                "grimoire off" / "grimoire on" said in a session
//	core.active                  grimoire disable / grimoire enable
//	active                       default

// activeSource names where the active state came from.
type activeSource string

const (
	activeFromDefault activeSource = "default"
	activeFromConfig  activeSource = "config"
	activeFromSession activeSource = "session"
	activeFromEnv     activeSource = "env"
)

// grimoireActive reports whether grimoire is active for session (which may
// be "" when the agent sends none), and where that came from.
func grimoireActive(session string) (bool, activeSource) {
	if v := strings.TrimSpace(os.Getenv("GRIMOIRE_ACTIVE")); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b, activeFromEnv
		}
	}
	if b, ok := sessionActive(session); ok {
		return b, activeFromSession
	}
	if r, err := config.Load(getProjectDir()); err == nil && r.Core.Active != "" {
		if b, err := strconv.ParseBool(r.Core.Active); err == nil {
			return b, activeFromConfig
		}
	}
	return true, activeFromDefault
}

// sessionActiveFile holds a session's own on/off choice.
func sessionActiveFile(session string) string {
	dir := sessionStateDir(grimoireStateRoot(), session)
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, "active")
}

// sessionActive reads a session's own on/off choice; ok is false when the
// session made none.
func sessionActive(session string) (active, ok bool) {
	p := sessionActiveFile(session)
	if p == "" {
		return false, false
	}
	if fi, err := os.Lstat(p); err != nil || !fi.Mode().IsRegular() {
		return false, false
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return false, false
	}
	switch strings.TrimSpace(string(b)) {
	case "on":
		return true, true
	case "off":
		return false, true
	}
	return false, false
}

// setSessionActive records a session's on/off choice. Failures are reported
// so the hook can tell the user the toggle did not stick.
func setSessionActive(session string, active bool) error {
	p := sessionActiveFile(session)
	if p == "" {
		return fmt.Errorf("no session id from the agent")
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}
	if fi, err := os.Lstat(p); err == nil && fi.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("%s is a symlink", p)
	}
	v := "off"
	if active {
		v = "on"
	}
	return os.WriteFile(p, []byte(v), 0o600)
}

// Phrases that turn grimoire on or off for the session. A prompt that is
// mostly the phrase toggles; a prompt that merely mentions grimoire in
// passing ("is grimoire off by default?") does not. "quiet grimoire" is not
// one of them: it hides grimoire's notes (display off) and routing goes on.
var (
	toggleOff = regexp.MustCompile(`(?i)^\s*(?:please\s+)?(?:(?:turn|switch)\s+off|disable|deactivate|stop|pause)\s+grimoire(?:\s+mode)?\s*[.!]*\s*$|^\s*grimoire\s+(?:off|inactive|disable|stop)\s*[.!]*\s*$`)
	toggleOn  = regexp.MustCompile(`(?i)^\s*(?:please\s+)?(?:(?:turn|switch)\s+on|enable|activate|start|resume)\s+grimoire(?:\s+mode)?\s*[.!]*\s*$|^\s*grimoire\s+(?:on|active|enable|start)\s*[.!]*\s*$`)
)

// parseToggle reports whether prompt asks to turn grimoire on or off.
func parseToggle(prompt string) (on, off bool) {
	p := strings.TrimSpace(prompt)
	switch {
	case toggleOff.MatchString(p):
		return false, true
	case toggleOn.MatchString(p):
		return true, false
	}
	return false, false
}

// inactiveRule tells the model grimoire is inactive. It overrides the
// instruction line and tool descriptions, which say to route every request.
const inactiveRule = "grimoire is inactive in this session, at the user's request: do not route requests through start-best-practice, do not load grimoire skills on your own, and ignore any grimoire router or routing reminder above. Load a grimoire skill only when the user names one (grimoire:<name>). The user can say \"grimoire on\" to resume."

// toggleContext tells the model about a toggle the user just made, so it
// acknowledges it instead of treating the phrase as a task.
func toggleContext(active bool) string {
	if active {
		return "The user just turned grimoire on for this session. Acknowledge in one short line; from the next request, route requests through start-best-practice again."
	}
	return "The user just turned grimoire off for this session. Acknowledge in one short line. " + inactiveRule
}

// handleToggle applies an on/off phrase in prompt to session. ok is false
// when prompt is not a toggle.
func handleToggle(prompt, session string) (ctx, msg string, ok bool) {
	on, off := parseToggle(prompt)
	if !on && !off {
		return "", "", false
	}
	if err := setSessionActive(session, on); err != nil {
		cmd := "grimoire disable"
		if on {
			cmd = "grimoire enable"
		}
		return fmt.Sprintf("The user asked to turn grimoire %s for this session, but it could not be saved (%v). Tell them, and that %#q changes it for all sessions.", onOff(on), err, cmd),
			statusPrefix + "could not change grimoire for this session — " + err.Error(), true
	}
	// Confirmed by the model's own one-line acknowledgement (and the badge),
	// not a status line: agents render those as warnings, and a toggle the
	// user just made is no warning.
	return toggleContext(on), "", true
}

// helpPhrase matches a request for grimoire's controls.
var helpPhrase = regexp.MustCompile(`(?i)^\s*(?:grimoire\s+help|help\s+grimoire|grimoire\s+commands|what\s+can\s+i\s+say\s+to\s+grimoire)\s*[?.!]*\s*$`)

// grimoireHelp is the list of controls, answered on "grimoire help".
const grimoireHelp = `The user asked for grimoire's controls. Reply with exactly this list, nothing else:
grimoire controls — say any of these:
- "show grimoire details" — show evidence (candidates, risks, sources)
- "quiet grimoire" — hide grimoire's notes; it keeps routing
- "plain grimoire" — words instead of symbols (screen readers)
- "skip grimoire" — skip it for this request
- "grimoire off" / "grimoire on" — stop or resume routing in this session
For all sessions: grimoire disable / grimoire enable, and grimoire config set core.trace hidden|brief|full.`

func onOff(b bool) string {
	if b {
		return "on"
	}
	return "off"
}

// promptTurn is the prompt hook's answer, shared by every agent's prompt
// hook: an on/off phrase is applied first; then, when active, turn gives
// the routing context; when inactive, only an explicit grimoire:<name> is
// honored and the inactive rule keeps the model from routing. showsStatus is
// false for agents whose hooks cannot show the user a status line
// (Antigravity), so the one-time tip is not used up unseen there.
func promptTurn(prompt, session string, l displayLevel, showsStatus bool, turn func(string, displayLevel) (string, string)) (ctx, msg string) {
	if helpPhrase.MatchString(prompt) {
		return grimoireHelp, ""
	}
	if ctx, msg, ok := handleToggle(prompt, session); ok {
		if !l.on() {
			msg = ""
		}
		return ctx, msg
	}
	if active, _ := grimoireActive(session); active {
		ctx, msg = turn(prompt, l)
		if showsStatus {
			msg = withOnboardingTip(l, ctx, msg)
		}
		return ctx, msg
	}
	p := strings.TrimSpace(prompt)
	if p == "" || strings.HasPrefix(p, "/") || isSmallTalk(p) {
		return "", ""
	}
	if refs := explicitSkillRefs(p); len(refs) > 0 {
		ctx, _ := turn(prompt, l) // loads the named skills
		return ctx + "\n\n" + inactiveRule, ""
	}
	return inactiveRule, ""
}

// inactiveSessionStatus is the session-start status line when grimoire is
// inactive. Grimoire not working is state worth knowing (most agents have no
// badge), but agents render hook messages as warnings, and repeating a
// choice the user made on purpose at every session start is nagging — so at
// brief it appears at most once a day; at full, every session.
func inactiveSessionStatus(l displayLevel, src activeSource) string {
	if l == displayBrief && !takeDailyNotice("off-notice") {
		return ""
	}
	where := "for this session"
	switch src {
	case activeFromConfig:
		where = "for all sessions (grimoire disable)"
	case activeFromEnv:
		where = "(GRIMOIRE_ACTIVE)"
	}
	return statusLine(l, "off "+where+` — say "grimoire on" to resume`)
}

// ── grimoire enable / disable ────────────────────────────────────────────────

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Turn grimoire on: route every request through its best-practice workflows",
	Long: `Turns grimoire on for all sessions (core.active = true, the default).
New agent sessions pick it up. A session can still turn grimoire off for
itself by saying "grimoire off".`,
	Args: cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error { return setGlobalActive(true) },
}

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Turn grimoire off: stop routing requests (named skills still load)",
	Long: `Turns grimoire off for all sessions (core.active = false): no router,
no routing reminders. A skill you name with grimoire:<name> still loads, and
nothing is uninstalled. A session can turn grimoire back on for itself by
saying "grimoire on". Undo with: grimoire enable`,
	Args: cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error { return setGlobalActive(false) },
}

func init() {
	rootCmd.AddCommand(enableCmd, disableCmd)
}

func setGlobalActive(active bool) error {
	path := config.GlobalPath()
	fs, err := config.LoadFile(path)
	if err != nil {
		return fmt.Errorf("loading %s: %w", path, err)
	}
	if err := applyKey(&fs, "core.active", strconv.FormatBool(active)); err != nil {
		return err
	}
	if err := config.WriteFile(path, fs); err != nil {
		return fmt.Errorf("saving %s: %w", path, err)
	}
	if active {
		fmt.Printf("%s  grimoire is on — new sessions route requests through its workflows\n", tui.IconOK)
	} else {
		fmt.Printf("%s  grimoire is off — new sessions won't route requests; grimoire:<name> still loads a skill\n", tui.IconOK)
		fmt.Printf("   undo: grimoire enable · one session: say \"grimoire on\" in it\n")
	}
	fmt.Printf("   saved to %s (core.active)\n", path)
	if v := os.Getenv("GRIMOIRE_ACTIVE"); v != "" {
		fmt.Printf("%s  GRIMOIRE_ACTIVE=%s is set and overrides this in this shell\n", tui.IconWarn, v)
	}
	return nil
}

// takeDailyNotice reports whether the notice named name is due today, and
// records it as shown. Unrecordable means not due: better quiet than
// repeated.
func takeDailyNotice(name string) bool {
	dir, err := os.UserCacheDir()
	if err != nil {
		return false
	}
	p := filepath.Join(dir, "grimoire", name)
	today := time.Now().Format("2006-01-02")
	if fi, err := os.Lstat(p); err == nil {
		if !fi.Mode().IsRegular() {
			return false
		}
		if b, err := os.ReadFile(p); err == nil && strings.TrimSpace(string(b)) == today {
			return false
		}
	}
	if os.MkdirAll(filepath.Dir(p), 0o700) != nil || os.WriteFile(p, []byte(today), 0o600) != nil {
		return false
	}
	return true
}
