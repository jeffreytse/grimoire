package tui

import (
	"fmt"
	"sync"
	"time"
)

// LiveBoard renders a fixed list of labeled rows that animate while in-progress
// and show a final icon when complete. No-op on non-TTY.
type LiveBoard struct {
	mu      sync.Mutex
	rows    []liveRow
	frame   int
	started bool
	stop    chan struct{}
	done    chan struct{}
}

type liveRow struct {
	label  string
	fin    bool
	icon   string
	detail string
}

// NewLiveBoard creates a board with one row per label.
func NewLiveBoard(labels []string) *LiveBoard {
	b := &LiveBoard{
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	for _, l := range labels {
		b.rows = append(b.rows, liveRow{label: l})
	}
	return b
}

// Start prints initial spinner rows and starts the render loop. No-op on non-TTY.
func (b *LiveBoard) Start() {
	if !IsTTY() {
		return
	}
	b.mu.Lock()
	for _, r := range b.rows {
		fmt.Printf("  %s %s\n", StyleCyan.Render(spinnerFrames[0]), r.label)
	}
	b.started = true
	b.mu.Unlock()

	go func() {
		defer close(b.done)
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-b.stop:
				b.render()
				return
			case <-ticker.C:
				b.render()
			}
		}
	}()
}

func (b *LiveBoard) render() {
	b.mu.Lock()
	defer b.mu.Unlock()
	n := len(b.rows)
	if n == 0 {
		return
	}
	b.frame++
	fmt.Printf("\033[%dA", n)
	for _, r := range b.rows {
		fmt.Print("\r\033[K")
		if r.fin {
			if r.detail != "" {
				fmt.Printf("  %s %s  %s\n", r.icon, r.label, r.detail)
			} else {
				fmt.Printf("  %s %s\n", r.icon, r.label)
			}
		} else {
			fmt.Printf("  %s %s\n", StyleCyan.Render(spinnerFrames[b.frame%len(spinnerFrames)]), r.label)
		}
	}
}

// Complete marks row i as done with icon and optional detail. Thread-safe.
func (b *LiveBoard) Complete(i int, icon, detail string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if i < 0 || i >= len(b.rows) {
		return
	}
	b.rows[i].fin = true
	b.rows[i].icon = icon
	b.rows[i].detail = detail
}

// Stop signals the render loop, does a final render, and waits for exit.
// No-op on non-TTY or if Start was never called.
func (b *LiveBoard) Stop() {
	if !IsTTY() || !b.started {
		return
	}
	close(b.stop)
	<-b.done
}
