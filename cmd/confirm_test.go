package cmd

import (
	"strings"
	"testing"
)

func stubConfirm(t *testing.T, tty bool, answer string) *bool {
	t.Helper()
	origSel, origTTY := confirmSelect, confirmIsTTY
	t.Cleanup(func() { confirmSelect, confirmIsTTY = origSel, origTTY })
	asked := false
	confirmIsTTY = func() bool { return tty }
	confirmSelect = func(string, []string) (string, bool) { asked = true; return answer, true }
	return &asked
}

func TestConfirm(t *testing.T) {
	cases := []struct {
		name                   string
		tty, yes, destructive  bool
		answer                 string
		want, wantAsk, wantErr bool
	}{
		{"--yes skips the prompt", true, true, true, "No", true, false, false},
		{"terminal: Yes proceeds", true, false, false, "Yes", true, true, false},
		{"terminal: No cancels", true, false, true, "No", false, true, false},
		{"no terminal, safe: proceeds", false, false, false, "", true, false, false},
		{"no terminal, destructive: fails", false, false, true, "", false, false, true},
		{"no terminal, destructive, --yes: proceeds", false, true, true, "", true, false, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			asked := stubConfirm(t, c.tty, c.answer)
			got, err := confirm("Replace it?", c.yes, c.destructive)
			if got != c.want || *asked != c.wantAsk || (err != nil) != c.wantErr {
				t.Errorf("got %v (asked %v, err %v); want %v (asked %v, err %v)",
					got, *asked, err, c.want, c.wantAsk, c.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), "--yes") {
				t.Errorf("error should say how to proceed: %v", err)
			}
		})
	}
}

// update without a terminal goes ahead rather than "cancelling" with exit 0.
func TestConfirmUpdate_NoTerminalProceeds(t *testing.T) {
	asked := stubConfirm(t, false, "No")
	orig := flagUpdateYes
	t.Cleanup(func() { flagUpdateYes = orig })
	flagUpdateYes = false
	if !confirmUpdate("Pull now?") || *asked {
		t.Error("update without a terminal should proceed without prompting")
	}
}

// Replacing the skill library deletes it, so a script must pass --yes.
func TestInstallFrom_ReplaceNeedsYesWithoutTerminal(t *testing.T) {
	setupListSkillsHome(t) // official package home exists
	stubConfirm(t, false, "")
	orig := flagInstallYes
	t.Cleanup(func() { flagInstallYes = orig })
	flagInstallYes = false
	_, err := resolveAndPersistSource("https://github.com/example/skills.git")
	if err == nil || !strings.Contains(err.Error(), "--yes") {
		t.Fatalf("want a --yes error, got %v", err)
	}
}

func TestWizard_NeedsTerminal(t *testing.T) {
	stubConfirm(t, false, "")
	if err := wizardCmd.RunE(wizardCmd, nil); err != errNeedsTerminal {
		t.Errorf("wizard without a terminal = %v, want errNeedsTerminal", err)
	}
}

func TestCleanAndUninstall_RejectUnknownTarget(t *testing.T) {
	origC, origU := flagCleanTarget, flagUninstallTarget
	t.Cleanup(func() { flagCleanTarget, flagUninstallTarget = origC, origU })
	flagCleanTarget, flagUninstallTarget = "bogus", "bogus"
	if err := runClean(nil, nil); err == nil {
		t.Error("clean --target bogus should error")
	}
	if err := runUninstall(nil, nil); err == nil {
		t.Error("uninstall --target bogus should error")
	}
}
