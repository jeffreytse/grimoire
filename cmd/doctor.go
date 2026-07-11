package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/config"
	"github.com/jeffreytse/grimoire/internal/git"
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
	home := skills.OfficialPackageHome()
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
		if _, err := exec.LookPath(agent.Binary(ag)); err != nil {
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
		if loadCap := agent.SkillLoadCap(ag); loadCap > 0 && count > loadCap {
			checks = append(checks, doctorCheck{
				Name:   "agent-" + ag + "-skill-cap",
				Status: "warn",
				Detail: fmt.Sprintf(
					"%s: %d skills installed but agent only loads %d by default — "+
						"create or edit ~/.openclaw/openclaw.json:\n"+
						`      { "skills": { "limits": { "maxSkillsLoadedPerSource": %d, "maxSkillsInPrompt": %d } } }`,
					label, count, loadCap, count, count),
			})
			ok = false
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
	cwd := getProjectDir()
	cfgPaths := []struct{ path, label string }{
		{config.ProjectPath(cwd), "project (grimoire.toml)"},
		{config.GlobalPath(), "global (" + config.GlobalPath() + ")"},
		{config.SystemPath(), "system (" + config.SystemPath() + ")"},
	}
	for _, c := range cfgPaths {
		if _, err := os.Stat(c.path); err != nil {
			checks = append(checks, doctorCheck{Name: "config-" + filepath.Base(c.path), Status: "skip", Detail: c.label + " not found"})
		} else {
			checks = append(checks, doctorCheck{Name: "config-" + filepath.Base(c.path), Status: "ok", Detail: c.label + " present"})
		}
	}

	if shared, err := config.ParseFile(config.ProjectPath(cwd)); err == nil {
		if shared.Core.Home != "" {
			checks = append(checks, doctorCheck{
				Name:   "config-core-in-project",
				Status: "warn",
				Detail: "[core] section in grimoire.toml is ignored — [core] is user-level; move to global: grimoire config set core.home <path> --global",
			})
			ok = false
		}
	}

	// ── Package config ────────────────────────────────────────────────────────────
	if cfg, err := config.LoadGlobal(); err == nil {
		officialCount := 0
		for _, rd := range cfg.Packages {
			if rd.Official {
				officialCount++
			}
		}
		if officialCount > 1 {
			checks = append(checks, doctorCheck{
				Name:   "package-multiple-official",
				Status: "warn",
				Detail: fmt.Sprintf("%d packages have official=true — only the first is used; run: grimoire package list to review", officialCount),
			})
			ok = false
		}
	}

	// ── Migration: old unversioned package dirs ──────────────────────────────────
	if old := scanUnversionedPackageDirs(skills.PackagesRoot()); len(old) > 0 {
		checks = append(checks, doctorCheck{
			Name:   "packages-unversioned",
			Status: "warn",
			Detail: fmt.Sprintf("%d old-style unversioned package dir(s) found (%s) — run: grimoire install to re-clone to versioned paths, then remove manually",
				len(old), strings.Join(old, ", ")),
		})
		ok = false
	}

	return doctorOutput{OK: ok, Checks: checks}
}

// scanUnversionedPackageDirs walks PackagesRoot and finds repo-level dirs that
// lack a "@version" suffix — these are from the pre-versioned-path layout.
func scanUnversionedPackageDirs(root string) []string {
	var found []string
	// Structure: <root>/<host>/<owner>/<repo>[@version]
	// Walk two levels deep, collect leaf dirs without '@'.
	hosts, _ := os.ReadDir(root)
	for _, h := range hosts {
		if !h.IsDir() {
			continue
		}
		owners, _ := os.ReadDir(filepath.Join(root, h.Name()))
		for _, o := range owners {
			if !o.IsDir() {
				continue
			}
			repos, _ := os.ReadDir(filepath.Join(root, h.Name(), o.Name()))
			for _, r := range repos {
				if !r.IsDir() {
					continue
				}
				if !strings.Contains(r.Name(), "@") {
					found = append(found, h.Name()+"/"+o.Name()+"/"+r.Name())
				}
			}
		}
	}
	return found
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
		switch c.Status {
		case "error":
			errs++
		case "warn":
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
		switch c.Status {
		case "error":
			errs++
		case "warn":
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
		switch c.Status {
		case "error":
			errs++
		case "warn":
			warnings++
		}
	}

	hasExtends := false
	for _, c := range out.Checks {
		if strings.HasPrefix(c.Name, "extends-") {
			hasExtends = true
			break
		}
	}
	if hasExtends {
		fmt.Println()
		fmt.Println("  Packages")
		for _, c := range out.Checks {
			if !strings.HasPrefix(c.Name, "extends-") {
				continue
			}
			icon, _ := doctorIcon(c.Status)
			fmt.Printf("    %s  %s\n", icon, c.Detail)
			switch c.Status {
			case "error":
				errs++
			case "warn":
				warnings++
			}
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

func doctorIcon(status string) (icon, label string) {
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
