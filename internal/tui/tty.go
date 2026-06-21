package tui

import (
	"io"
	"os"
)

// selectOutput returns the writer for TUI programs.
// Uses stderr so stdout stays clean for piping.
func selectOutput() io.Writer {
	return os.Stderr
}

// IsTTY reports whether stdin is an interactive terminal.
func IsTTY() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
