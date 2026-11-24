package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jeffreytse/grimoire/internal/agent"
	"github.com/jeffreytse/grimoire/internal/skills"
	"github.com/jeffreytse/grimoire/internal/tui"
)

// listingFit estimates how many of sks an agent's skill listing shows: it
// fills the agent's listing budget in install order. need is the characters
// all of them would take. budget is 0 when the agent's budget is unknown.
func listingFit(ag string, sks []skills.Skill) (visible, need, budget int) {
	budget = agent.SkillListingBudget(ag)
	used := 0
	for i := range sks {
		c := agent.SkillListingCost(sks[i].Name, sks[i].Description)
		need += c
		if used+c <= budget {
			used += c
			visible++
		}
	}
	return visible, need, budget
}

// listingWarning returns the warning for an agent whose skill listing cannot
// hold total skills needing need characters, or "" when they fit or the
// budget is unknown. Agents with MCP support are pointed to MCP mode, which
// has no listing limit; every agent is pointed to installing a subset.
func listingWarning(ag string, visible, total, need, budget int) string {
	if budget == 0 || need <= budget {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "\n  %s  %s will see about %d of %d skills — its skill list holds ~%s characters; these need ~%s.\n",
		tui.IconWarn, agent.DisplayName(ag), visible, total, thousands(budget), thousands(need))
	b.WriteString("     The rest stay installed but invisible to its skill discovery.\n")
	if findMCPAgent(ag) != nil {
		fmt.Fprintf(&b, "     Full library, no limit:  grimoire mcp install --target %s\n", ag)
	}
	b.WriteString("     Or a smaller subset:     grimoire install --domain <domain> --subdomain <subdomain>\n")
	return b.String()
}

// warnListingBudgets prints a listing warning for each target the skills
// overflow. --no-load-optimize silences it, as it lifts the load caps.
func warnListingBudgets(w io.Writer, targets []string, sks []skills.Skill) {
	if flagInstallNoCap {
		return
	}
	known := false
	for _, ag := range targets {
		known = known || agent.SkillListingBudget(ag) > 0
	}
	if !known {
		return
	}
	sks = withDescriptions(sks)
	for _, ag := range targets {
		visible, need, budget := listingFit(ag, sks)
		if msg := listingWarning(ag, visible, len(sks), need, budget); msg != "" {
			fmt.Fprint(w, msg)
		}
	}
}

// withDescriptions returns sks with each missing description read from its
// SKILL.md. The install path lists skills without reading any file, for
// speed; the listing estimate needs the descriptions, so it reads them only
// when a target agent has a known budget.
func withDescriptions(sks []skills.Skill) []skills.Skill {
	out := make([]skills.Skill, len(sks))
	for i := range sks {
		out[i] = sks[i]
		if out[i].Description == "" && out[i].Path != "" {
			out[i].Description = skills.Description(out[i].Path)
		}
	}
	return out
}

// thousands formats n with comma separators: 321723 → "321,723".
func thousands(n int) string {
	s := fmt.Sprint(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

// installedSkills lists the skills linked in an agent's skills dir, with
// their descriptions, for estimating its listing cost.
func installedSkills(dir string) []skills.Skill {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []skills.Skill
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") || e.Type().IsRegular() {
			continue
		}
		out = append(out, skills.Skill{Name: e.Name(), Description: skills.Description(filepath.Join(dir, e.Name()))})
	}
	return out
}

// listingDoctorDetail returns doctor's warning for an agent whose installed
// skills overflow its skill listing, or "" when they fit.
func listingDoctorDetail(ag, dir string) string {
	sks := installedSkills(dir)
	visible, need, budget := listingFit(ag, sks)
	if budget == 0 || need <= budget {
		return ""
	}
	d := fmt.Sprintf("%s sees about %d of %d installed skills — its skill list holds ~%s characters; they need ~%s",
		agent.DisplayName(ag), visible, len(sks), thousands(budget), thousands(need))
	if findMCPAgent(ag) != nil {
		d += fmt.Sprintf(" — for the full library use MCP mode: grimoire mcp install --target %s", ag)
	} else {
		d += " — install a smaller subset: grimoire install --domain <domain> --subdomain <subdomain>"
	}
	return d
}
