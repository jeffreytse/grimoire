package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Re-run compliance check whenever project files change",
	Long: `Watch the project directory for file changes and re-run the AI compliance
check automatically on each save. Uses the same executor resolution as 'grimoire check'.

Press Ctrl+C to stop.`,
	RunE: runWatch,
}

func init() {
	watchCmd.Flags().StringVar(&flagVia, "via", "", "force a specific local AI agent (claude, gemini, codex, copilot, opencode, openclaw)")
	watchCmd.Flags().BoolVar(&flagNoColor, "no-color", false, "disable ANSI color")
	watchCmd.Flags().BoolVar(&flagNoGitignore, "no-gitignore", false, "disable .gitignore-based file filtering")
}

func runWatch(_ *cobra.Command, _ []string) error {
	projectDir := getProjectDir()
	initGitignoreMatcher(projectDir)
	initExcludePatterns(projectDir, nil)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("creating watcher: %w", err)
	}
	defer func() { _ = watcher.Close() }()

	if err := watchAddRecursive(watcher, projectDir); err != nil {
		return fmt.Errorf("watching %s: %w", projectDir, err)
	}

	fmt.Printf("Watching %s — press Ctrl+C to stop\n\n", projectDir)

	_, _ = runIndependentCheck(context.Background(), projectDir)

	const debounceDelay = 500 * time.Millisecond
	var debounce *time.Timer
	var mu sync.Mutex
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
			rel, relErr := filepath.Rel(projectDir, event.Name)
			if relErr != nil {
				rel = filepath.Base(event.Name)
			}
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

				for _, f := range files {
					if isGrimoireConfig(projectDir, filepath.Join(projectDir, f)) {
						fmt.Printf("\n── Config changed — full re-check…\n\n")
						_, _ = runIndependentCheck(context.Background(), projectDir)
						return
					}
				}
				label := files[0]
				if len(files) > 1 {
					label = fmt.Sprintf("%d files", len(files))
				}
				fmt.Printf("\n── %s changed — re-checking…\n\n", label)
				_, _ = runIndependentCheck(context.Background(), projectDir, files...)
			})
		case watchErr, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(os.Stderr, "watcher: %v\n", watchErr)
		}
	}
}

// watchAddRecursive adds dir and all subdirectories to the watcher,
// skipping hidden dirs, vendor, and node_modules.
func watchAddRecursive(w *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		name := d.Name()
		if strings.HasPrefix(name, ".") {
			return filepath.SkipDir
		}
		rel, _ := filepath.Rel(root, path)
		if rel != "." && shouldSkip(filepath.ToSlash(rel), true) {
			return filepath.SkipDir
		}
		return w.Add(path)
	})
}

// watchShouldIgnore returns true for hidden-segment paths (dot-dirs like .git, .grimoire).
// Per-event gitignore+exclude filtering is done at the pending-accumulation point via shouldSkip.
func watchShouldIgnore(name string) bool {
	for _, seg := range strings.Split(filepath.ToSlash(name), "/") {
		if strings.HasPrefix(seg, ".") {
			return true
		}
	}
	return false
}
