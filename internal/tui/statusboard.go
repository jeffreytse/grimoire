package tui

import (
	"fmt"
	"sync"
	"time"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type lineState int

const (
	stateWaiting lineState = iota
	stateUpdating
	stateDone
)

// StatusBoard prints per-registry status lines and updates them in-place on a TTY.
// On non-TTY (pipe/CI), prints simple "updating X..." lines upfront and is a no-op thereafter.
type StatusBoard struct {
	mu     sync.Mutex
	names  []string
	states []lineState
	n      int
	isTTY  bool
}

// NewStatusBoard prints N status lines (waiting state) and returns the board.
func NewStatusBoard(names []string) *StatusBoard {
	sb := &StatusBoard{
		names:  names,
		states: make([]lineState, len(names)),
		n:      len(names),
		isTTY:  IsTTY(),
	}
	waitIcon := StyleDim.Render("○")
	for _, name := range names {
		if sb.isTTY {
			fmt.Printf("  %s  %s\n", waitIcon, StyleDim.Render(name))
		} else {
			fmt.Printf("  updating %s...\n", name)
		}
	}
	return sb
}

// StartSpinner launches the braille animation ticker for lines in "updating" state.
// Returns a stop func — call it after wg.Wait() before flushing results.
func (sb *StatusBoard) StartSpinner() func() {
	if !sb.isTTY {
		return func() {}
	}
	stop := make(chan struct{})
	go func() {
		frame := 0
		for {
			select {
			case <-stop:
				return
			case <-time.After(80 * time.Millisecond):
				sb.tick(spinnerFrames[frame%len(spinnerFrames)])
				frame++
			}
		}
	}()
	return func() { close(stop) }
}

func (sb *StatusBoard) tick(frame string) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	icon := StyleCyan.Render(frame)
	for i, state := range sb.states {
		if state == stateUpdating {
			sb.rewriteLocked(i, icon, sb.names[i])
		}
	}
}

// SetUpdating marks line i as actively updating. Call after the goroutine acquires its semaphore slot.
func (sb *StatusBoard) SetUpdating(i int) {
	if !sb.isTTY {
		return
	}
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.states[i] = stateUpdating
	sb.rewriteLocked(i, StyleCyan.Render(spinnerFrames[0]), sb.names[i])
}

// SetDone marks line i as done. Pass IconDone or IconError.
func (sb *StatusBoard) SetDone(i int, icon string) {
	if !sb.isTTY {
		return
	}
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.states[i] = stateDone
	sb.rewriteLocked(i, icon, sb.names[i])
}

// Finish prints a blank separator line after the status block.
func (sb *StatusBoard) Finish() { fmt.Println() }

// rewriteLocked rewrites status line i in-place using ANSI cursor movement.
// Caller must hold sb.mu. Cursor is assumed to be one line below the last status line.
func (sb *StatusBoard) rewriteLocked(i int, icon, name string) {
	up := sb.n - i
	label := fmt.Sprintf("  %s  %s", icon, name)
	fmt.Printf("\033[%dA\r\033[2K%s\033[%dB\r", up, label, up)
}
