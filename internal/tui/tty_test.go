package tui

import (
	"os"
	"testing"
)

// CI jobs often run with stdin from /dev/null, which is a character device
// but not a terminal; IsTTY must say no so nothing waits on a prompt.
func TestIsTTY_DevNullIsNotATerminal(t *testing.T) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		t.Skip(err)
	}
	defer func() { _ = f.Close() }()
	orig := os.Stdin
	os.Stdin = f
	defer func() { os.Stdin = orig }()
	if IsTTY() {
		t.Error("stdin from /dev/null reported as a terminal")
	}
}
