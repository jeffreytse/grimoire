package cmd

import (
	"bytes"
	"context"
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
// grimoire is active, and the skill grimoire last served when that was
// recent: [GRIMOIRE] or [GRIMOIRE:suggest-best-practice].
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
		runStatusLine(os.Stdin, os.Stdout, flagStatusLineThen, lastSkillFile())
		return nil
	},
}

func init() {
	statusLineCmd.Flags().StringVar(&flagStatusLineThen, "then", "", "status line command to run first; the badge is appended to its output")
	rootCmd.AddCommand(statusLineCmd)
}

func runStatusLine(in io.Reader, out io.Writer, then, lastFile string) {
	raw, _ := io.ReadAll(io.LimitReader(in, 1<<20))
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

// statusBadge renders the badge, naming the last served skill when the file
// at lastFile is recent and safe to show.
func statusBadge(lastFile string) string {
	label := "GRIMOIRE"
	if name := recentLastSkill(lastFile); name != "" {
		label += ":" + name
	}
	return "\033[38;5;141m[" + label + "]\033[0m"
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

// lastSkillFile is where the MCP server records the last skill it served,
// or "" when no cache dir is available.
func lastSkillFile() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "grimoire", "last-skill")
}

// recordLastSkill notes name as the last skill served, for the status line.
// Failures are ignored: the badge is cosmetic.
func recordLastSkill(name string) {
	path := lastSkillFile()
	name = sanitizeSkillName(name)
	if path == "" || name == "" {
		return
	}
	if fi, err := os.Lstat(path); err == nil && fi.Mode()&os.ModeSymlink != 0 {
		return
	}
	if os.MkdirAll(filepath.Dir(path), 0o700) == nil {
		_ = os.WriteFile(path, []byte(name), 0o600)
	}
}
