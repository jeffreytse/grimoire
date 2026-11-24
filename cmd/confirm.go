package cmd

import (
	"errors"
	"fmt"

	"github.com/jeffreytse/grimoire/internal/tui"
)

// confirmSelect and confirmIsTTY are tui.RunSelect and tui.IsTTY; tests
// replace them.
var (
	confirmSelect = tui.RunSelect
	confirmIsTTY  = tui.IsTTY
)

// confirm is the one rule every grimoire prompt follows:
//   - yes (--yes/-y) → go ahead
//   - someone at the terminal → ask
//   - no terminal, safe action → go ahead: running the command is the request
//   - no terminal, destructive action → fail, so a script never deletes data
//     unasked and never reports success for work it skipped
//
// It returns false only when the person answered No.
func confirm(title string, yes, destructive bool) (bool, error) {
	if yes {
		return true, nil
	}
	if !confirmIsTTY() {
		if destructive {
			return false, fmt.Errorf("%s — this can't be undone, so it needs confirmation; rerun with --yes", title)
		}
		return true, nil
	}
	chosen, ok := confirmSelect(title, []string{"Yes", "No"})
	return ok && chosen == "Yes", nil
}

// errNeedsTerminal is returned by commands that only work interactively.
var errNeedsTerminal = errors.New("the wizard needs an interactive terminal — to set up grimoire from a script, run: grimoire mcp install")
