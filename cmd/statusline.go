package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// `grimoire statusline` prints a Claude Code status line badge showing that
// grimoire is active, and while it is working, the workflow and the skill it
// last served: [GRIMOIRE], [GRIMOIRE:plan], [GRIMOIRE:plan › design-slo].
// The color is grimoire's brand gold from assets/icon-book.svg.
//
// Claude Code has one statusLine slot. When another tool already owns it,
// `grimoire mcp install --statusline` chains it: --then runs that command
// with the same stdin and the badge is appended to its output.

var flagStatusLineThen string

// statusLineFresh is how long the last served skill stays in the badge.
var statusLineFresh = 30 * time.Minute

// statusLineThenTimeout bounds the chained command; the status line must not hang.
const statusLineThenTimeout = 3 * time.Second

// lastSkillMax caps how much of the last-skill file is read.
const lastSkillMax = 64

var statusLineCmd = &cobra.Command{
	Use:   "statusline",
	Short: "Claude Code status line badge for grimoire",
	Long: `Prints a status line badge for Claude Code: [GRIMOIRE], or
[GRIMOIRE:<skill>] while a skill grimoire served is recent.

With --then, runs that command first (with the same stdin) and appends the
badge to its output, so grimoire can share the status line with another tool.
Registered by: grimoire mcp install`,
	Args:          cobra.NoArgs,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		runStatusLine(os.Stdin, os.Stdout, flagStatusLineThen, grimoireStateRoot())
		return nil
	},
}

func init() {
	statusLineCmd.Flags().StringVar(&flagStatusLineThen, "then", "", "status line command to run first; the badge is appended to its output")
	rootCmd.AddCommand(statusLineCmd)
}

// runStatusLine prints the chained status line (if any) and the badge for
// the session named in Claude Code's status line input. State lives under
// stateRoot, per session, so concurrent sessions never show each other's skill.
func runStatusLine(in io.Reader, out io.Writer, then, stateRoot string) {
	raw, _ := io.ReadAll(io.LimitReader(in, 1<<20))
	var input struct {
		SessionID string `json:"session_id"`
	}
	_ = json.Unmarshal(raw, &input)
	lastFile := ""
	if dir := sessionStateDir(stateRoot, input.SessionID); dir != "" {
		lastFile = filepath.Join(dir, "last-skill")
	}
	var parts []string
	if then != "" {
		if s := runChainedStatusLine(then, raw); s != "" {
			parts = append(parts, s)
		}
	}
	parts = append(parts, statusBadge(lastFile))
	fmt.Fprint(out, strings.Join(parts, " "))
}

// runChainedStatusLine runs command with stdin and returns its output without
// trailing newlines, or "" on failure.
func runChainedStatusLine(command string, stdin []byte) string {
	ctx, cancel := context.WithTimeout(context.Background(), statusLineThenTimeout)
	defer cancel()
	c := exec.CommandContext(ctx, "sh", "-c", command)
	c.Stdin = bytes.NewReader(stdin)
	var buf bytes.Buffer
	c.Stdout = &buf
	_ = c.Run()
	return strings.TrimRight(buf.String(), "\r\n")
}

// workflowShort names grimoire's workflow skills in the badge.
var workflowShort = map[string]string{
	"start-best-practice":           "route",
	"suggest-best-practice":         "suggest",
	"plan-best-practice-solution":   "plan",
	"apply-best-practice-tree":      "tree",
	"analyze-best-practice-problem": "analyze",
	"review-best-practice-fit":      "review",
	"discover-best-practices":       "discover",
}

// statusBadge renders the badge: the recent workflow and the last skill
// served after it, when the state files are recent and safe to show.
func statusBadge(lastFile string) string {
	label := "GRIMOIRE"
	wf := workflowShort[recentLastSkill(workflowFile(lastFile))]
	skill := recentLastSkill(lastFile)
	if short, isWorkflow := workflowShort[skill]; isWorkflow {
		wf, skill = short, "" // the last skill served is the workflow itself
	}
	switch {
	case wf != "" && skill != "":
		label += ":" + wf + " › " + skill
	case wf != "":
		label += ":" + wf
	case skill != "":
		label += ":" + skill
	}
	return badgeColor() + "[" + label + "]\033[0m"
}

// brandGold is grimoire's brand color, #c9933a (assets/icon-book.svg).
const (
	brandGoldTrue = "\033[1;38;2;201;147;58m" // 24-bit
	brandGold256  = "\033[1;38;5;173m"        // nearest xterm-256 color, #d7875f
)

// badgeColor is the brand gold in 24-bit color where the terminal says it
// supports it (COLORTERM), else the nearest 256-color.
func badgeColor() string {
	switch strings.ToLower(os.Getenv("COLORTERM")) {
	case "truecolor", "24bit":
		return brandGoldTrue
	}
	return brandGold256
}

// workflowFile holds the last workflow skill served, next to lastFile.
func workflowFile(lastFile string) string {
	if lastFile == "" {
		return ""
	}
	return filepath.Join(filepath.Dir(lastFile), "last-workflow")
}

// recentLastSkill reads the last served skill name. The file is refused when
// it is a symlink, and its content is cut to lastSkillMax bytes of
// [a-z0-9-], so it can never put terminal escapes on the status line.
func recentLastSkill(path string) string {
	if path == "" {
		return ""
	}
	fi, err := os.Lstat(path)
	if err != nil || !fi.Mode().IsRegular() || time.Since(fi.ModTime()) > statusLineFresh {
		return ""
	}
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	raw, _ := io.ReadAll(io.LimitReader(f, lastSkillMax))
	return sanitizeSkillName(string(raw))
}

func sanitizeSkillName(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// grimoireStateRoot is where grimoire keeps per-session status state, or ""
// when no cache dir is available.
func grimoireStateRoot() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "grimoire", "sessions")
}

// sessionStateDir is the state dir for one agent session, or "" when either
// part is missing. The session id is reduced to a safe file name.
func sessionStateDir(root, sessionID string) string {
	id := agyStateName(strings.TrimSpace(sessionID))
	if root == "" || id == "" || len(id) > 128 {
		return ""
	}
	return filepath.Join(root, id)
}

// sessionStateTTL is how long a session's state is kept after its last write.
const sessionStateTTL = 48 * time.Hour

// recordSkill notes name as the skill in use for a session, and as its
// workflow when name is one of grimoire's workflows. Failures are ignored:
// the badge is cosmetic.
func recordSkill(root, sessionID, name string) {
	dir := sessionStateDir(root, sessionID)
	name = sanitizeSkillName(name)
	if dir == "" || name == "" {
		return
	}
	if fi, err := os.Lstat(dir); err == nil && fi.Mode()&os.ModeSymlink != 0 {
		return
	}
	if os.MkdirAll(dir, 0o700) != nil {
		return
	}
	write := func(file string) {
		p := filepath.Join(dir, file)
		if fi, err := os.Lstat(p); err == nil && fi.Mode()&os.ModeSymlink != 0 {
			return
		}
		_ = os.WriteFile(p, []byte(name), 0o600)
	}
	if _, isWorkflow := workflowShort[name]; isWorkflow {
		write("last-workflow")
	}
	write("last-skill")
	pruneSessionState(root)
}

// pruneSessionState removes session dirs untouched for sessionStateTTL.
func pruneSessionState(root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}
	for _, e := range entries {
		info, err := e.Info()
		if err != nil || !e.IsDir() || time.Since(info.ModTime()) < sessionStateTTL {
			continue
		}
		_ = os.RemoveAll(filepath.Join(root, e.Name()))
	}
}
