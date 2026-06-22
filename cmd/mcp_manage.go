package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/jeffreytse/grimoire/internal/agent"
	gitops "github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/skills"
)

// ── Output types ──────────────────────────────────────────────────────────────

type mcpInstallOutput struct {
	Installed map[string]int `json:"installed"`
	Skipped   map[string]int `json:"skipped"`
	Errors    []string       `json:"errors,omitempty"`
}

type mcpUninstallOutput struct {
	Removed map[string]int `json:"removed"`
	Errors  []string       `json:"errors,omitempty"`
}

type mcpUpdateOutput struct {
	AlreadyUpToDate bool   `json:"already_up_to_date"`
	OldVersion      string `json:"old_version,omitempty"`
	NewVersion      string `json:"new_version,omitempty"`
	OldCommit       string `json:"old_commit,omitempty"`
	NewCommit       string `json:"new_commit,omitempty"`
	SkillsAdded     []string `json:"skills_added,omitempty"`
	SkillsUpdated   []string `json:"skills_updated,omitempty"`
	ProfilesAdded   []string `json:"profiles_added,omitempty"`
	ProfilesUpdated []string `json:"profiles_updated,omitempty"`
	PresetsAdded    []string `json:"presets_added,omitempty"`
	PresetsUpdated  []string `json:"presets_updated,omitempty"`
}

type mcpCleanOutput struct {
	Removed map[string]int `json:"removed"`
	Total   int            `json:"total"`
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func toolGrimoireInstall(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	domain := request.GetString("domain", "")
	subdomain := request.GetString("subdomain", "")
	skill := request.GetString("skill", "")
	target := request.GetString("target", "")
	out, err := performInstall(domain, subdomain, skill, target)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(out)
}

func toolGrimoireUninstall(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	domain := request.GetString("domain", "")
	subdomain := request.GetString("subdomain", "")
	skill := request.GetString("skill", "")
	target := request.GetString("target", "")
	out, err := performUninstall(domain, subdomain, skill, target)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(out)
}

func toolGrimoireUpdate(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	stable := request.GetString("stable", "") == "true"
	out, err := performUpdate(stable)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(out)
}

func toolGrimoireClean(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	target := request.GetString("target", "")
	out, err := performClean(target)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(out)
}

// ── Install ───────────────────────────────────────────────────────────────────

func performInstall(domain, subdomain, skill, target string) (mcpInstallOutput, error) {
	sources := skills.AllSkillsSources()
	if len(sources) == 0 {
		return mcpInstallOutput{}, fmt.Errorf("skills not found — run grimoire_update first")
	}
	targets := resolveTargets(target)
	perAgent := make(map[string]int)
	var errs []string

	switch {
	case skill != "":
		skillPath, _, err := resolveSkillFromSources(sources, skill)
		if err != nil {
			return mcpInstallOutput{}, err
		}
		for _, ag := range targets {
			ok, err := skills.InstallSkill(skillPath, agent.SkillsDir(ag), true)
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", ag, err))
			} else if ok {
				perAgent[ag]++
			}
		}

	case domain != "":
		for _, src := range sources {
			for _, ag := range targets {
				n, agErrs := installDomainSilent(src.Root, domain, subdomain, ag)
				perAgent[ag] += n
				errs = append(errs, agErrs...)
			}
		}

	default:
		all, err := skills.ListAllSkillsFromSources(sources)
		if err != nil {
			return mcpInstallOutput{}, err
		}
		for _, sk := range all {
			for _, ag := range targets {
				ok, err := skills.InstallSkill(sk.Path, agent.SkillsDir(ag), true)
				if err != nil {
					errs = append(errs, err.Error())
				} else if ok {
					perAgent[ag]++
				}
			}
		}
	}

	for _, ag := range targets {
		_, _ = skills.CleanBrokenSymlinks(agent.SkillsDir(ag))
		_ = agent.ConfigureAgentMD(ag)
	}

	installed := make(map[string]int)
	skipped := make(map[string]int)
	for _, ag := range targets {
		if perAgent[ag] > 0 {
			installed[agent.DisplayName(ag)] = perAgent[ag]
		} else {
			skipped[agent.DisplayName(ag)] = 0
		}
	}
	return mcpInstallOutput{Installed: installed, Skipped: skipped, Errors: errs}, nil
}

func installDomainSilent(root, domain, subdomain, ag string) (count int, errs []string) {
	domainDir := root + "/" + domain
	if _, err := os.Stat(domainDir); err != nil {
		return 0, []string{fmt.Sprintf("domain not found: %s", domain)}
	}
	destDir := agent.SkillsDir(ag)

	if skills.IsNested(domainDir) {
		subs, err := skills.ListSubdomains(domainDir)
		if err != nil {
			return 0, []string{err.Error()}
		}
		for _, sub := range subs {
			if subdomain != "" && sub.Name != subdomain {
				continue
			}
			skillList, err := skills.ListSkillsInDir(sub.Path, domain, sub.Name)
			if err != nil {
				continue
			}
			for _, sk := range skillList {
				ok, err := skills.InstallSkill(sk.Path, destDir, true)
				if err != nil {
					errs = append(errs, err.Error())
				} else if ok {
					count++
				}
			}
		}
	} else {
		skillList, err := skills.ListSkillsInDir(domainDir, domain, "")
		if err != nil {
			return 0, []string{err.Error()}
		}
		for _, sk := range skillList {
			ok, err := skills.InstallSkill(sk.Path, destDir, true)
			if err != nil {
				errs = append(errs, err.Error())
			} else if ok {
				count++
			}
		}
	}
	return count, errs
}

// ── Uninstall ─────────────────────────────────────────────────────────────────

func performUninstall(domain, subdomain, skill, target string) (mcpUninstallOutput, error) {
	root := skills.SkillsRoot()
	targets := resolveTargets(target)
	removed := make(map[string]int)
	var errs []string

	switch {
	case skill != "":
		skillPath, err := skills.ResolveSkillPath(root, skill)
		if err != nil {
			skillPath = skill
		}
		name := skillNameFromPath(skillPath, skill)
		for _, ag := range targets {
			ok, err := skills.UninstallSkill(name, agent.SkillsDir(ag))
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", ag, err))
			} else if ok {
				removed[agent.DisplayName(ag)]++
			}
		}

	case domain != "":
		for _, ag := range targets {
			n, agErrs := uninstallDomainSilent(root, domain, subdomain, ag)
			removed[agent.DisplayName(ag)] += n
			errs = append(errs, agErrs...)
		}

	default:
		domains, err := skills.ListDomains(root)
		if err != nil {
			return mcpUninstallOutput{}, err
		}
		for _, d := range domains {
			for _, ag := range targets {
				n, agErrs := uninstallDomainSilent(root, d.Name, "", ag)
				removed[agent.DisplayName(ag)] += n
				errs = append(errs, agErrs...)
			}
		}
	}

	for _, ag := range targets {
		_, _ = skills.CleanBrokenSymlinks(agent.SkillsDir(ag))
		if agent.SkillCount(ag) == 0 {
			_ = agent.RemoveAgentMDConfig(ag)
		}
	}
	return mcpUninstallOutput{Removed: removed, Errors: errs}, nil
}

func uninstallDomainSilent(root, domain, subdomain, ag string) (count int, errs []string) {
	domainDir := root + "/" + domain
	if _, err := os.Stat(domainDir); err != nil {
		return 0, nil
	}
	destDir := agent.SkillsDir(ag)

	if skills.IsNested(domainDir) {
		subs, err := skills.ListSubdomains(domainDir)
		if err != nil {
			return 0, []string{err.Error()}
		}
		for _, sub := range subs {
			if subdomain != "" && sub.Name != subdomain {
				continue
			}
			skillList, err := skills.ListSkillsInDir(sub.Path, domain, sub.Name)
			if err != nil {
				continue
			}
			for _, sk := range skillList {
				ok, err := skills.UninstallSkill(sk.Name, destDir)
				if err != nil {
					errs = append(errs, err.Error())
				} else if ok {
					count++
				}
			}
		}
	} else {
		skillList, err := skills.ListSkillsInDir(domainDir, domain, "")
		if err != nil {
			return 0, []string{err.Error()}
		}
		for _, sk := range skillList {
			ok, err := skills.UninstallSkill(sk.Name, destDir)
			if err != nil {
				errs = append(errs, err.Error())
			} else if ok {
				count++
			}
		}
	}
	return count, errs
}

// ── Update ────────────────────────────────────────────────────────────────────

func performUpdate(stable bool) (mcpUpdateOutput, error) {
	home := skills.OfficialRegistryHome()
	url := skills.GrimoireRepoURL()

	// Local registry: skip all git ops
	if filepath.IsAbs(url) {
		if _, err := os.Stat(home); err != nil {
			return mcpUpdateOutput{}, fmt.Errorf("local registry %q not found", home)
		}
		return mcpUpdateOutput{AlreadyUpToDate: true}, nil
	}

	if _, err := os.Stat(home); err != nil {
		if err := os.MkdirAll(filepath.Dir(home), 0o755); err != nil {
			return mcpUpdateOutput{}, fmt.Errorf("creating dir: %w", err)
		}
		if err := gitops.Clone(url, home); err != nil {
			return mcpUpdateOutput{}, fmt.Errorf("cloning grimoire: %w", err)
		}
		state, _ := gitops.CurrentState(home)
		return mcpUpdateOutput{NewVersion: state.Version, NewCommit: state.Commit}, nil
	}

	current, err := gitops.CurrentState(home)
	if err != nil {
		return mcpUpdateOutput{}, fmt.Errorf("reading current state: %w", err)
	}

	if stable {
		_ = gitops.FetchTags(home) // non-fatal
		latest, err := gitops.LatestTag(home)
		if err != nil {
			return mcpUpdateOutput{}, fmt.Errorf("finding latest tag: %w", err)
		}
		tagState, err := gitops.TagState(home, latest)
		if err != nil {
			return mcpUpdateOutput{}, fmt.Errorf("reading tag state: %w", err)
		}
		if tagState.Commit == current.Commit {
			return mcpUpdateOutput{
				AlreadyUpToDate: true,
				OldVersion:      current.Version,
				OldCommit:       current.Commit,
			}, nil
		}
		if err := gitops.CheckoutTag(home, latest); err != nil {
			return mcpUpdateOutput{}, fmt.Errorf("checking out %s: %w", latest, err)
		}
		ch, _ := gitops.RegistryChangesSince(home, current.Commit)
		relinkNewSkills(home, current.Commit)
		return mcpUpdateOutput{
			OldVersion: current.Version, OldCommit: current.Commit,
			NewVersion: tagState.Version, NewCommit: tagState.Commit,
			SkillsAdded: ch.SkillsAdded, SkillsUpdated: ch.SkillsUpdated,
			ProfilesAdded: ch.ProfilesAdded, ProfilesUpdated: ch.ProfilesUpdated,
			PresetsAdded: ch.PresetsAdded, PresetsUpdated: ch.PresetsUpdated,
		}, nil
	}

	upToDate, _, _, err := gitops.IsUpToDate(home)
	if err != nil {
		return mcpUpdateOutput{}, fmt.Errorf("checking upstream: %w", err)
	}
	if upToDate {
		return mcpUpdateOutput{
			AlreadyUpToDate: true,
			OldVersion:      current.Version,
			OldCommit:       current.Commit,
		}, nil
	}
	if err := gitops.PullWithForceFallback(home); err != nil {
		return mcpUpdateOutput{}, fmt.Errorf("updating: %w", err)
	}
	newState, _ := gitops.CurrentState(home)
	ch, _ := gitops.RegistryChangesSince(home, current.Commit)
	relinkNewSkills(home, current.Commit)
	updateCustomRegistries()
	return mcpUpdateOutput{
		OldVersion: current.Version, OldCommit: current.Commit,
		NewVersion: newState.Version, NewCommit: newState.Commit,
		SkillsAdded: ch.SkillsAdded, SkillsUpdated: ch.SkillsUpdated,
		ProfilesAdded: ch.ProfilesAdded, ProfilesUpdated: ch.ProfilesUpdated,
		PresetsAdded: ch.PresetsAdded, PresetsUpdated: ch.PresetsUpdated,
	}, nil
}

// ── Clean ─────────────────────────────────────────────────────────────────────

func performClean(target string) (mcpCleanOutput, error) {
	targets := resolveTargets(target)
	removed := make(map[string]int)
	for _, ag := range targets {
		n, err := skills.CleanBrokenSymlinks(agent.SkillsDir(ag))
		if err == nil {
			removed[agent.DisplayName(ag)] = n
		}
	}
	total := 0
	for _, n := range removed {
		total += n
	}
	return mcpCleanOutput{Removed: removed, Total: total}, nil
}
