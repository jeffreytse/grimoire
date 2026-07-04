package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/jeffreytse/grimoire/internal/compliance"
	"github.com/jeffreytse/grimoire/internal/tui"
)

// liveEvent is a typed SSE message sent to all connected browser tabs.
type liveEvent struct {
	Event string // SSE named event: "status" | "reload"
	Data  string // raw JSON payload
}

// runLiveCheck runs an initial compliance check, starts an HTTP server that serves
// the HTML report with SSE-based auto-refresh and in-page status badge, then watches
// for file changes and re-runs the check on each save.
func runLiveCheck(projectDir string) error {
	var subsMu sync.Mutex
	var subs []chan liveEvent
	var lastStatus liveEvent // replayed to new subscribers so they see current state

	broadcast := func(ev liveEvent) {
		subsMu.Lock()
		defer subsMu.Unlock()
		if ev.Event == "status" {
			lastStatus = ev
		}
		for _, ch := range subs {
			select {
			case ch <- ev:
			default:
			}
		}
	}

	broadcastStatus := func(status string, files []string, elapsed string) {
		filesJSON, _ := json.Marshal(files)
		broadcast(liveEvent{
			Event: "status",
			Data:  fmt.Sprintf(`{"status":%q,"files":%s,"elapsed":%q}`, status, filesJSON, elapsed),
		})
	}

	broadcastReload := func() {
		broadcast(liveEvent{Event: "reload", Data: `{}`})
	}

	subscribe := func() (chan liveEvent, func()) {
		ch := make(chan liveEvent, 2)
		subsMu.Lock()
		subs = append(subs, ch)
		last := lastStatus
		subsMu.Unlock()
		if last.Event != "" {
			ch <- last // replay current state on connect
		}
		return ch, func() {
			subsMu.Lock()
			for i, s := range subs {
				if s == ch {
					subs = append(subs[:i], subs[i+1:]...)
					break
				}
			}
			subsMu.Unlock()
		}
	}

	jsonPath := filepath.Join(projectDir, ".grimoire", "reports", "compliance-latest.json")

	mux := http.NewServeMux()

	// Serve the HTML report rendered fresh from JSON on each request.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html, err := compliance.RenderHTMLReport(jsonPath, cliVersion, projectDir, true)
		if err != nil {
			// JSON not yet written — show minimal loading page with live badge.
			html, _ = compliance.RenderLoadingHTML()
			if html == nil {
				http.Error(w, "Check still running…", http.StatusServiceUnavailable)
				return
			}
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(html)
	})

	// SSE endpoint — each tab subscribes and receives typed named events.
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		ch, unsub := subscribe()
		defer unsub()
		for {
			select {
			case <-r.Context().Done():
				return
			case ev := <-ch:
				_, _ = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Event, ev.Data)
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}
		}
	})

	host := flagHost
	if host == "" {
		host = "0.0.0.0"
	}
	addr := fmt.Sprintf("%s:%d", host, flagPort)
	srv := &http.Server{Addr: addr, Handler: mux}
	srvErr := make(chan error, 1)
	go func() { srvErr <- srv.ListenAndServe() }()
	select {
	case err := <-srvErr:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("live server: %w (try --port or --host to change address)", err)
		}
	case <-time.After(150 * time.Millisecond):
		// server started successfully
	}

	stateFile := filepath.Join(os.TempDir(), fmt.Sprintf("grimoire-live-%d", flagPort))
	shouldOpen := true
	if info, err := os.Stat(stateFile); err == nil && time.Since(info.ModTime()) < 5*time.Minute {
		shouldOpen = false
	}
	_ = os.WriteFile(stateFile, nil, 0o600)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	broadcastShutdown := func() {
		broadcast(liveEvent{Event: "shutdown", Data: `{}`})
		time.Sleep(300 * time.Millisecond)
	}

	displayHost := flagHost
	if displayHost == "" || displayHost == "0.0.0.0" {
		displayHost = "localhost"
	}
	url := fmt.Sprintf("http://%s:%d", displayHost, flagPort)

	// clearLive clears the screen and reprints the sticky header so each check
	// cycle replaces the previous output rather than accumulating as logs.
	clearLive := func() {
		fmt.Print("\033[2J\033[H")
		if tui.IsTTY() {
			fmt.Printf("  \033[1;36mGRIMOIRE\033[0m  \033[2m%s\033[0m  \033[32mlive\033[0m\n\n", cliVersion)
			fmt.Printf("  \033[36m➜\033[0m  Report:   \033[0;4m%s\033[0m\n", url)
			fmt.Printf("  \033[2m➜  press ctrl+c to stop\033[0m\n\n")
		} else {
			fmt.Printf("  GRIMOIRE  %s  live\n\n  Report: %s\n  press ctrl+c to stop\n\n", cliVersion, url)
		}
	}

	liveCtx, liveCancel := context.WithCancel(context.Background())
	defer liveCancel()

	clearLive()
	fmt.Printf("[%s] ── initial check\n\n", time.Now().Format("15:04:05"))
	if shouldOpen {
		openBrowser(url)
	}

	broadcastStatus("analyzing", nil, "")
	t0 := time.Now()
	ran, _ := runIndependentCheck(liveCtx, projectDir)
	elapsed := time.Since(t0).Round(time.Millisecond).String()
	broadcastStatus("done", nil, elapsed)
	broadcastReload() // always: transition browser from loading page to report
	if ran {
		fmt.Printf("\n[%s] ✓ done in %s\n", time.Now().Format("15:04:05"), elapsed)
	} else {
		fmt.Printf("\n[%s] · unchanged\n", time.Now().Format("15:04:05"))
	}

	// File watcher — reuse helpers from watch.go (same package).
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("creating watcher: %w", err)
	}
	defer func() { _ = watcher.Close() }()

	if err := watchAddRecursive(watcher, projectDir); err != nil {
		return fmt.Errorf("watching %s: %w", projectDir, err)
	}

	const debounceDelay = 500 * time.Millisecond
	var debounce *time.Timer
	var mu sync.Mutex
	var checkMu sync.Mutex // serializes concurrent check runs
	pending := make(map[string]bool)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if watchShouldIgnore(event.Name) {
				continue
			}
			if event.Has(fsnotify.Create) {
				if fi, statErr := os.Stat(event.Name); statErr == nil && fi.IsDir() {
					_ = watcher.Add(event.Name)
				}
			}
			rel, _ := filepath.Rel(projectDir, event.Name)
			rel = filepath.ToSlash(rel)
			if shouldSkip(rel, false) {
				continue
			}
			mu.Lock()
			pending[rel] = true
			mu.Unlock()

			if debounce != nil {
				debounce.Stop()
			}
			debounce = time.AfterFunc(debounceDelay, func() {
				mu.Lock()
				var files []string
				for f := range pending {
					files = append(files, f)
				}
				pending = make(map[string]bool)
				mu.Unlock()

				if len(files) == 0 {
					return
				}

				// Cancel any in-flight check before acquiring the mutex — kills the
				// subprocess fast so checkMu.Lock() unblocks almost immediately.
				liveCancel()
				liveCtx, liveCancel = context.WithCancel(context.Background())

				checkMu.Lock()
				defer checkMu.Unlock()

				label := files[0]
				if len(files) > 1 {
					label = fmt.Sprintf("%d files", len(files))
				}
				ts := time.Now().Format("15:04:05")

				for _, f := range files {
					if !isGrimoireConfig(projectDir, filepath.Join(projectDir, f)) {
						continue
					}
					clearLive()
					fmt.Printf("[%s] ── config changed · full re-check\n\n", ts)
					broadcastStatus("analyzing", files, "")
					t0 := time.Now()
					ran, _ := runIndependentCheck(liveCtx, projectDir)
					if liveCtx.Err() != nil {
						return
					}
					elapsed := time.Since(t0).Round(time.Millisecond).String()
					broadcastStatus("done", files, elapsed)
					if ran {
						fmt.Printf("\n[%s] ✓ done in %s\n", time.Now().Format("15:04:05"), elapsed)
						broadcastReload()
					} else {
						fmt.Printf("\n[%s] · unchanged\n", time.Now().Format("15:04:05"))
					}
					return
				}

				clearLive()
				fmt.Printf("[%s] ── %s changed · analyzing…\n\n", ts, label)
				broadcastStatus("analyzing", files, "")
				t0 := time.Now()
				ran, _ := runIndependentCheck(liveCtx, projectDir, files...)
				if liveCtx.Err() != nil {
					return
				}
				elapsed := time.Since(t0).Round(time.Millisecond).String()
				broadcastStatus("done", files, elapsed)
				if ran {
					fmt.Printf("\n[%s] ✓ done in %s\n", time.Now().Format("15:04:05"), elapsed)
					broadcastReload()
				} else {
					fmt.Printf("\n[%s] · unchanged\n", time.Now().Format("15:04:05"))
				}
			})

		case watchErr, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(os.Stderr, "watcher: %v\n", watchErr)

		case <-sigCh:
			broadcastShutdown()
			shutCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			_ = srv.Shutdown(shutCtx)
			cancel()
			return nil
		}
	}
}
