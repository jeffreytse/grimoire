package tui

import (
	"fmt"
	"time"
)

// NewSpinner starts an animated braille spinner on the current line showing msg.
// Returns a stop func that clears the line. No-op on non-TTY.
func NewSpinner(msg string) func() {
	if !IsTTY() {
		return func() {}
	}
	ch := make(chan struct{})
	go func() {
		frame := 0
		for {
			select {
			case <-ch:
				fmt.Print("\r\033[K")
				return
			case <-time.After(80 * time.Millisecond):
				fmt.Printf("\r\033[K  %s %s", StyleCyan.Render(spinnerFrames[frame%len(spinnerFrames)]), msg)
				frame++
			}
		}
	}()
	return func() { close(ch) }
}
