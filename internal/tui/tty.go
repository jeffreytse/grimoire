package tui

import (
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	xterm "github.com/charmbracelet/x/term"
)

// ttyOpts returns ProgramOptions for bubbletea programs.
// Uses stderr for output so stdout stays clean for piping.
func ttyOpts() []tea.ProgramOption {
	return []tea.ProgramOption{tea.WithOutput(os.Stderr)}
}

// runProgram runs p and restores the terminal state afterward.
// bubbletea v1.x leaves ONLCR disabled on exit in some configurations,
// corrupting all subsequent shell prompt newlines. We open /dev/tty directly
// rather than using os.Stdin.Fd() — when invoked via `go run` or `make`,
// stdin may be a pipe (not the terminal), causing GetState/IsTerminal to
// silently fail and skip the entire restore path.
func runProgram(p *tea.Program) (tea.Model, error) {
	tty, _ := os.OpenFile("/dev/tty", os.O_RDWR, 0)

	fd := os.Stdin.Fd()
	if tty != nil {
		fd = tty.Fd()
		defer func() { _ = tty.Close() }()
	}
	state, _ := xterm.GetState(fd)

	final, err := p.Run()

	if state != nil {
		_ = xterm.Restore(fd, state)
	}
	sane := exec.Command("stty", "sane")
	if tty != nil {
		sane.Stdin = tty
	} else {
		sane.Stdin = os.Stdin
	}
	_ = sane.Run()

	return final, err
}

// IsTTY reports whether stdin is an interactive terminal.
func IsTTY() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
