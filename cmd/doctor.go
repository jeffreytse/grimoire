package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

var flagDoctorJSON bool

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run a health check on the grimoire installation",
	RunE:  runDoctor,
}

func init() {
	doctorCmd.Flags().BoolVar(&flagDoctorJSON, "json", false, "output as JSON")
}

type doctorCheck struct {
	Name   string `json:"name"`
	Status string `json:"status"` // "ok" | "warn" | "error" | "skip"
	Detail string `json:"detail,omitempty"`
}

type doctorOutput struct {
	OK     bool          `json:"ok"`
	Checks []doctorCheck `json:"checks"`
}

func runDoctor(cmd *cobra.Command, args []string) error {
	out := collectDoctorChecks()

	if flagDoctorJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	printDoctorHuman(out)
	return nil
}

func collectDoctorChecks() doctorOutput {
	var checks []doctorCheck
	ok := true

	// ── Source ──────────────────────────────────────────────────────────────────
	home := skills.GrimoireHome()
	if _, err := os.Stat(home); err != nil {
		checks = append(checks, doctorCheck{
			Name:   "grimoire-source",
			Status: "error",
			Detail: fmt.Sprintf("grimoire not found at %s — run: grimoire update", home),
		})
		ok = false
	} else {
		state, err := git.CurrentState(home)
		if err != nil {
			ver := skills.GrimoireVersion()
			checks = append(checks, doctorCheck{
				Name:   "grimoire-source",
				Status: "warn",
				Detail: fmt.Sprintf("grimoire %s (no git state at %s)", ver, home),
			})
		} else {
			checks = append(checks, doctorCheck{
				Name:   "grimoire-source",
				Status: "ok",
				Detail: fmt.Sprintf("grimoire %s (commit %s, %s) at %s", state.Version, state.Commit, state.Date, home),
			})
		}
	}

	// ── AI agents ────────────────────────────────────────────────────────────────
	for _, ag := range agent.All {
		label := agent.DisplayName(ag)
		if _, err := exec.LookPath(ag); err != nil {
			checks = append(checks, doctorCheck{
				Name:   "agent-" + ag,
				Status: "skip",
				Detail: label + " not found",
			})
			continue
		}
		ver := agent.Version(ag)
		vs := "detected"
		if ver != "" {
			vs = ver
		}
		dir := agent.SkillsDir(ag)
		if _, err := os.Stat(dir); err != nil {
			checks = append(checks, doctorCheck{
				Name:   "agent-" + ag,
				Status: "warn",
				Detail: fmt.Sprintf("%s %s — no skills installed, run: grimoire install --target %s", label, vs, ag),
			})
			ok = false
			continue
		}
		count := agent.SkillCount(ag)
		broken := agent.BrokenSymlinkCount(ag)
		if broken > 0 {
			checks = append(checks, doctorCheck{
				Name:   "agent-" + ag,
				Status: "warn",
				Detail: fmt.Sprintf("%s %s — %d skills, %d broken, run: grimoire clean --target %s", label, vs, count, broken, ag),
			})
			ok = false
		} else {
			checks = append(checks, doctorCheck{
				Name:   "agent-" + ag,
				Status: "ok",
				Detail: fmt.Sprintf("%s %s — %d grimoire skills", label, vs, count),
			})
		}
		if agent.IsConfigured(ag) {
			checks = append(checks, doctorCheck{
				Name:   "agent-" + ag + "-configured",
				Status: "ok",
				Detail: fmt.Sprintf("%s start-best-practice active", label),
			})
		} else {
			checks = append(checks, doctorCheck{
				Name:   "agent-" + ag + "-configured",
				Status: "warn",
				Detail: fmt.Sprintf("%s start-best-practice not configured — run: grimoire install --target %s", label, ag),
			})
			ok = false
		}
	}

	// Detect-only agents
	detectOnly := []struct{ label, cmd string }{
		{"Cursor", "cursor"},
		{"Windsurf", "windsurf"},
		{"Aider", "aider"},
	}
	for _, item := range detectOnly {
		if _, err := exec.LookPath(item.cmd); err == nil {
			ver := agent.Version(item.cmd)
			vs := "detected"
			if ver != "" {
				vs = ver
			}
			checks = append(checks, doctorCheck{
				Name:   "agent-" + item.cmd,
				Status: "ok",
				Detail: fmt.Sprintf("%s %s", item.label, vs),
			})
		} else {
			checks = append(checks, doctorCheck{
				Name:   "agent-" + item.cmd,
				Status: "skip",
				Detail: item.label + " not found",
			})
		}
	}

	// ── Config files ─────────────────────────────────────────────────────────────
	cwd, _ := os.Getwd()
	home2, _ := os.UserHomeDir()
	cfgPaths := []struct{ path, label string }{
		{filepath.Join(cwd, ".grimoire", "settings.toml"), "project (.grimoire/settings.toml)"},
		{filepath.Join(home2, ".config", "grimoire", "settings.toml"), "global (~/.config/grimoire/settings.toml)"},
		{settings.SystemPath(), "system (" + settings.SystemPath() + ")"},
	}
	for _, c := range cfgPaths {
		if _, err := os.Stat(c.path); err != nil {
			checks = append(checks, doctorCheck{Name: "config-" + filepath.Base(c.path), Status: "skip", Detail: c.label + " not found"})
		} else {
			checks = append(checks, doctorCheck{Name: "config-" + filepath.Base(c.path), Status: "ok", Detail: c.label + " present"})
		}
	}

	sharedPath := filepath.Join(cwd, ".grimoire", "settings.toml")
	if shared, err := settings.ParseFile(sharedPath); err == nil {
		if shared.Core.Home != "" || shared.Core.Source != "" {
			checks = append(checks, doctorCheck{
				Name:   "config-core-in-shared",
				Status: "warn",
				Detail: "[core] home/source in .grimoire/settings.toml — move to global: grimoire config set core.home <path> --global",
			})
			ok = false
		}
	}

	return doctorOutput{OK: ok, Checks: checks}
}

func printDoctorHuman(out doctorOutput) {
	fmt.Println("\nGrimoire health check")
	warnings, errs := 0, 0

	fmt.Println()
	fmt.Println("  Source")
	for _, c := range out.Checks {
		if c.Name != "grimoire-source" {
			continue
		}
		icon, _ := doctorIcon(c.Status)
		fmt.Printf("    %s  %s\n", icon, c.Detail)
		if c.Status == "error" {
			errs++
		} else if c.Status == "warn" {
			warnings++
		}
	}

	fmt.Println()
	fmt.Println("  AI agents")
	for _, c := range out.Checks {
		if len(c.Name) < 6 || c.Name[:6] != "agent-" {
			continue
		}
		icon, _ := doctorIcon(c.Status)
		fmt.Printf("    %s  %s\n", icon, c.Detail)
		if c.Status == "error" {
			errs++
		} else if c.Status == "warn" {
			warnings++
		}
	}

	fmt.Println()
	fmt.Println("  Config")
	for _, c := range out.Checks {
		if len(c.Name) < 7 || c.Name[:7] != "config-" {
			continue
		}
		icon, _ := doctorIcon(c.Status)
		fmt.Printf("    %s  %s\n", icon, c.Detail)
		if c.Status == "error" {
			errs++
		} else if c.Status == "warn" {
			warnings++
		}
	}

	fmt.Println()
	if errs == 0 && warnings == 0 {
		fmt.Printf("  %s  All checks passed.\n\n", tui.IconOK)
	} else {
		parts := ""
		if errs > 0 {
			parts += fmt.Sprintf("%d error(s)", errs)
		}
		if warnings > 0 {
			if parts != "" {
				parts += ", "
			}
			parts += fmt.Sprintf("%d warning(s)", warnings)
		}
		fmt.Printf("  Summary: %s.\n\n", parts)
	}
}

func doctorIcon(status string) (string, string) {
	switch status {
	case "ok":
		return tui.IconOK, "ok"
	case "warn":
		return tui.IconWarn, "warn"
	case "error":
		return tui.IconFail, "error"
	default:
		return tui.IconSkip, "skip"
	}
}
