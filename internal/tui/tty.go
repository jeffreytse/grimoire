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
