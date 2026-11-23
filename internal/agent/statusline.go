package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Claude Code status line badge. Claude Code has a single statusLine slot
// ({"statusLine": {"type": "command", "command": ...}} in settings.json), so
// grimoire never takes it over silently:
//
//	no statusLine        set it to `grimoire statusline`
//	someone else's       left alone, unless the user opts in to chaining:
//	                     `grimoire statusline --then '<their command>'`
//	                     runs theirs and appends the grimoire badge
//
// The original command is kept inside the --then argument, so uninstall
// restores it without any other state.

// statusLineArgs follows the grimoire binary in the statusLine command.
const statusLineArgs = "statusline"

// statusLineThen introduces the chained command.
const statusLineThen = " --then "

// StatusLine states reported by StatusLineState.
const (
	StatusLineNone    = "none"     // no statusLine configured
	StatusLineOurs    = "grimoire" // grimoire's badge alone
	StatusLineChained = "chained"  // grimoire's badge after another command
	StatusLineOther   = "other"    // someone else's statusLine
)

// StatusLineFile returns Claude Code's settings file for scope.
func StatusLineFile(scope, projectDir string) string {
	return PromptHookFile("claude", scope, projectDir)
}

// StatusLineCommand is the statusLine command for the grimoire binary bin,
// chaining then when it is not "".
func StatusLineCommand(bin, then string) string {
	c := shellQuote(bin) + " " + statusLineArgs
	if then != "" {
		c += statusLineThen + quoteAlways(then)
	}
	return c
}

// quoteAlways single-quotes s for a POSIX shell.
func quoteAlways(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// unquoteAlways reverses quoteAlways, or reports false when s is not in that
// form.
func unquoteAlways(s string) (string, bool) {
	if len(s) < 2 || s[0] != '\'' || s[len(s)-1] != '\'' {
		return "", false
	}
	inner := strings.ReplaceAll(s[1:len(s)-1], `'\''`, "\x00")
	if strings.Contains(inner, "'") {
		return "", false
	}
	return strings.ReplaceAll(inner, "\x00", "'"), true
}

// parseStatusLine splits a grimoire statusLine command into its chained
// command. ours is false when command is not grimoire's.
func parseStatusLine(command string) (then string, ours bool) {
	c := strings.TrimSpace(command)
	head, rest, chained := strings.Cut(c, statusLineThen)
	head = strings.TrimSpace(head)
	if !strings.HasSuffix(head, " "+statusLineArgs) {
		return "", false
	}
	bin := strings.Trim(strings.TrimSpace(strings.TrimSuffix(head, statusLineArgs)), `'"`)
	if filepath.Base(bin) != "grimoire" {
		return "", false
	}
	if !chained {
		return "", true
	}
	then, ok := unquoteAlways(strings.TrimSpace(rest))
	if !ok {
		return "", false
	}
	return then, true
}

// statusLineCommandOf returns cfg.statusLine.command, whether statusLine is
// set at all, and an error when it has a shape grimoire cannot read.
func statusLineCommandOf(cfg map[string]any, path string) (command string, set bool, err error) {
	v, ok := cfg["statusLine"]
	if !ok || v == nil {
		return "", false, nil
	}
	m, ok := v.(map[string]any)
	if !ok {
		return "", true, fmt.Errorf(`%s: "statusLine" is not an object — fix it by hand`, path)
	}
	c, _ := m["command"].(string)
	return c, true, nil
}

// StatusLineState reports what the statusLine in the settings file at path
// is, and the foreign command when it is someone else's.
func StatusLineState(path string) (state, other string) {
	cfg, err := readSettings(path)
	if err != nil {
		return StatusLineNone, ""
	}
	c, set, err := statusLineCommandOf(cfg, path)
	if err != nil || !set {
		return StatusLineNone, ""
	}
	then, ours := parseStatusLine(c)
	switch {
	case ours && then != "":
		return StatusLineChained, then
	case ours:
		return StatusLineOurs, ""
	}
	return StatusLineOther, c
}

// ConfigureStatusLine puts grimoire's badge in the statusLine of the settings
// file at path. An empty statusLine is set to the badge. Someone else's is
// chained behind the badge only when chain is true; otherwise it is left
// alone and returned as other. Reports whether the file changed.
func ConfigureStatusLine(path, bin string, chain bool) (changed bool, other string, err error) {
	cfg, err := readSettings(path)
	if err != nil {
		return false, "", err
	}
	cur, set, err := statusLineCommandOf(cfg, path)
	if err != nil {
		return false, "", err
	}
	var want string
	switch then, ours := parseStatusLine(cur); {
	case !set || strings.TrimSpace(cur) == "":
		want = StatusLineCommand(bin, "")
	case ours:
		want = StatusLineCommand(bin, then) // refresh the binary path
	case !chain:
		return false, cur, nil
	default:
		want = StatusLineCommand(bin, cur)
	}
	if want == cur {
		return false, "", nil
	}
	sl, _ := cfg["statusLine"].(map[string]any)
	if sl == nil {
		sl = map[string]any{}
	}
	sl["type"] = "command"
	sl["command"] = want
	cfg["statusLine"] = sl
	return true, "", writeSettings(path, cfg)
}

// RemoveStatusLine takes grimoire's badge out of the statusLine of the
// settings file at path, restoring a chained command. Someone else's
// statusLine is left alone. Reports whether the file changed.
func RemoveStatusLine(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}
	cfg, err := readSettings(path)
	if err != nil {
		return false, err
	}
	cur, set, err := statusLineCommandOf(cfg, path)
	if err != nil || !set {
		return false, err
	}
	then, ours := parseStatusLine(cur)
	if !ours {
		return false, nil
	}
	if then == "" {
		delete(cfg, "statusLine")
	} else {
		sl := cfg["statusLine"].(map[string]any)
		sl["command"] = then
	}
	return true, writeSettings(path, cfg)
}
