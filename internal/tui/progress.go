package tui

import "fmt"

// Progress is a single-line in-place counter for TTY terminals.
// On non-TTY, Set is a no-op; Done prints msg unconditionally.
type Progress struct {
	label  string
	tty    bool
	active bool
}

// NewProgress creates a Progress with the given label (e.g. "Installing").
func NewProgress(label string) *Progress {
	return &Progress{label: label, tty: IsTTY()}
}

// Set overwrites the current line on TTY with "  label  N/total  detail".
// When total == 0, omits the "/total" part. No-op on non-TTY.
func (p *Progress) Set(n, total int, detail string) {
	if !p.tty {
		return
	}
	if p.active {
		fmt.Print("\r\033[K")
	}
	if total > 0 {
		fmt.Printf("  %s  %d/%d  %s", p.label, n, total, detail)
	} else {
		fmt.Printf("  %s  %d  %s", p.label, n, detail)
	}
	p.active = true
}

// Done clears the progress line on TTY and prints msg (if non-empty).
func (p *Progress) Done(msg string) {
	if p.tty && p.active {
		fmt.Print("\r\033[K")
	}
	if msg != "" {
		fmt.Println(msg)
	}
	p.active = false
}
