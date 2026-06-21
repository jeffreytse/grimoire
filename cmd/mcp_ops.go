package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/jeffreytse/grimoire/internal/compliance"
	"github.com/jeffreytse/grimoire/internal/detect"
	gitops "github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/profiles"
	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/skills"
)

// ── MCP-only output types ─────────────────────────────────────────────────────

type mcpRegistrySetOutput struct {
	Registry string `json:"registry"`
	IsLocal  bool   `json:"is_local"`
}

type mcpPresetListEntry struct {
	Name     string `json:"name"`
	Registry string `json:"registry"`
}

// ── Output types ──────────────────────────────────────────────────────────────

type mcpVersionOutput struct {
	CLI      string `json:"cli"`
	Grimoire string `json:"grimoire,omitempty"`
	Home     string `json:"home"`
}

type mcpInitOutput struct {
	Dir              string   `json:"dir"`
	Profile          string   `json:"profile,omitempty"`
	Skills           []string `json:"skills,omitempty"`
	SkillCount       int      `json:"skill_count,omitempty"`
	HasReport        bool     `json:"has_report"`
	CoveragePct      float64  `json:"coverage_pct,omitempty"`
	ComplianceStatus string   `json:"compliance_status,omitempty"`
}

type mcpSelfUpdateOutput struct {
	CurrentVersion  string `json:"current_version"`
	LatestVersion   string `json:"latest_version"`
	AlreadyUpToDate bool   `json:"already_up_to_date"`
	UpdateAvailable bool   `json:"update_available"`
	Method          string `json:"method"`
	Updated         bool   `json:"updated"`
}

type mcpRegistryRemoveOutput struct {
	Name         string `json:"name"`
	Removed      bool   `json:"removed"`
	CloneDeleted bool   `json:"clone_deleted"`
}

type mcpRegistryUpdateResult struct {
	Name   string `json:"name"`
	Status string `json:"status"` // "ok" | "cloned" | "error"
	Error  string `json:"error,omitempty"`
}

// ── Handlers ──────────────────────────────────────────────────────────────────

func toolGrimoireVersion(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	home := skills.GrimoireHome()
	out := mcpVersionOutput{
		CLI:  strings.TrimPrefix(cliVersion, "v"),
		Home: home,
	}
	if _, err := os.Stat(home); err == nil {
		if state, err := gitops.CurrentState(home); err == nil {
			out.Grimoire = state.Version
		} else {
			out.Grimoire = skills.GrimoireVersion()
		}
	}
	return jsonResult(out)
}

func toolGrimoireInit(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	cwd := getProjectDir()
	dir := filepath.Join(cwd, ".grimoire")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("creating .grimoire/: %v", err)), nil
	}
	cfg := loadExistingInitConfig(cwd, detect.Profile(cwd))
	// Explicit parameters override auto-detected values.
	if v := request.GetString("profile", ""); v != "" {
		cfg.Profile = v
	}
	if v := request.GetFloat("threshold", 0); v > 0 {
		cfg.Threshold = int(v)
	}
	if v := request.GetFloat("max_errors", -1); v >= 0 {
		cfg.MaxErrors = int(v)
	}
	if err := writeSettings(dir, cfg); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	out := mcpInitOutput{Dir: dir, Profile: cfg.Profile}
	if cfg.Profile != "" {
		if p, err := profiles.ResolveWithOptions(cfg.Profile, cwd, resolveOpts(cwd)); err == nil {
			out.SkillCount = len(p.Skills)
			names := make([]string, len(p.Skills))
			for i, sk := range p.Skills {
				names[i] = sk.Name
			}
			out.Skills = names
		}
	}
	if r, err := compliance.Load(resolvedReportPath(cwd)); err == nil {
		out.HasReport = true
		out.CoveragePct = r.Coverage.OverallPct
		out.ComplianceStatus = r.Threshold.Status
	}
	return jsonResult(out)
}

func toolGrimoireSelfUpdate(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	apply := request.GetString("yes", "") == "true"
	out, err := performSelfUpdate(apply)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(out)
}

func toolGrimoireRegistryList(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	entries, err := collectRegistryList()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(entries)
}

func toolGrimoireRegistryUpdate(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	name := request.GetString("name", "")
	results, err := performRegistryUpdate(name)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(results)
}

// ── Self-update ───────────────────────────────────────────────────────────────

func performSelfUpdate(apply bool) (mcpSelfUpdateOutput, error) {
	exePath, err := resolvedExePath()
	if err != nil {
		return mcpSelfUpdateOutput{}, fmt.Errorf("locating binary: %w", err)
	}
	method := detectInstallMethod(exePath)

	rel, err := fetchLatestRelease()
	if err != nil {
		return mcpSelfUpdateOutput{}, fmt.Errorf("fetching release info: %w", err)
	}

	latestTag := strings.TrimPrefix(rel.TagName, "v")
	currentTag := strings.TrimPrefix(cliVersion, "v")

	out := mcpSelfUpdateOutput{
		CurrentVersion: currentTag,
		LatestVersion:  latestTag,
		Method:         method,
	}

	if currentTag != "dev" && currentTag == latestTag {
		out.AlreadyUpToDate = true
		return out, nil
	}

	out.UpdateAvailable = true

	if !apply {
		return out, nil
	}

	switch method {
	case "go":
		pkg := ghModule + "@" + rel.TagName
		var buf bytes.Buffer
		cmd := exec.Command("go", "install", pkg)
		cmd.Stdout = &buf
		cmd.Stderr = &buf
		if err := cmd.Run(); err != nil {
			return out, fmt.Errorf("go install failed: %w", err)
		}
	default:
		if err := performSelfUpdateBinaryPlatform(exePath, rel); err != nil {
			return out, err
		}
	}

	out.Updated = true
	return out, nil
}

func performSelfUpdateBinaryPlatform(exePath string, rel *ghRelease) error {
	assetName := fmt.Sprintf("grimoire-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		assetName += ".exe"
	}

	var downloadURL string
	for _, a := range rel.Assets {
		if a.Name == assetName {
			downloadURL = a.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		return fmt.Errorf("no release asset found for %s/%s (looked for %q)", runtime.GOOS, runtime.GOARCH, assetName)
	}

	if runtime.GOOS == "windows" {
		return downloadFile(downloadURL, exePath+".new")
	}

	tmpPath := exePath + ".tmp"
	if err := downloadFile(downloadURL, tmpPath); err != nil {
		return err
	}
	if err := os.Chmod(tmpPath, 0o755); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("chmod: %w", err)
	}
	if err := os.Rename(tmpPath, exePath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("replacing binary: %w", err)
	}
	return nil
}

// ── Registry ──────────────────────────────────────────────────────────────────

func collectRegistryList() ([]registryListEntry, error) {
	fs, err := settings.LoadGlobal()
	if err != nil {
		return nil, fmt.Errorf("loading settings: %w", err)
	}
	r, _ := settings.Load(getProjectDir())

	coreURL := skills.GrimoireRepoURL()
	coreKind := "core"
	if filepath.IsAbs(coreURL) {
		coreKind = "local"
	}
	coreVersion := skills.GrimoireVersion()
	if fs.Core.Registry != "" {
		_, v := settings.ParseRef(fs.Core.Registry)
		if v != "" {
			coreVersion = v
		}
	}
	officialHome := skills.OfficialRegistryHome()
	entries := []registryListEntry{{
		Name:        settings.DeriveRegistryName(coreURL),
		URL:         coreURL,
		Version:     coreVersion,
		SkillsCount: countSkills(skills.SkillsRoot()),
		Cloned:      dirExists(officialHome),
		Kind:        coreKind,
	}}

	for _, ref := range r.StandardsExtends {
		u, ver := settings.ParseRef(ref)
		name := settings.DeriveRegistryName(u)
		extHome := skills.ExtendsHome(name)
		kind := "extends"
		if filepath.IsAbs(u) {
			kind = "local"
		}
		entries = append(entries, registryListEntry{
			Name:        name,
			URL:         u,
			Version:     ver,
			SkillsCount: countSkills(filepath.Join(extHome, "skills")),
			Cloned:      dirExists(extHome),
			Kind:        kind,
		})
	}
	return entries, nil
}

// performRegistryAdd adds a registry to standards.extends and clones it.
// Mirrors CLI `grimoire registry add` behaviour.
func performRegistryAdd(ref string) (registryListEntry, error) {
	if filepath.IsAbs(ref) {
		if _, err := os.Stat(ref); err != nil {
			return registryListEntry{}, fmt.Errorf("local path %q not found", ref)
		}
		gfs, err := settings.LoadGlobal()
		if err != nil {
			return registryListEntry{}, fmt.Errorf("loading settings: %w", err)
		}
		for _, existing := range gfs.StandardsExtends {
			eu, _ := settings.ParseRef(existing)
			if eu == ref {
				return registryListEntry{Name: ref, URL: ref, Cloned: true, Kind: "local"}, nil
			}
		}
		gfs.StandardsExtends = append(gfs.StandardsExtends, ref)
		if err := settings.SaveGlobal(gfs); err != nil {
			return registryListEntry{}, fmt.Errorf("saving settings: %w", err)
		}
		return registryListEntry{
			Name:        ref,
			URL:         ref,
			SkillsCount: countSkills(filepath.Join(ref, "skills")),
			Cloned:      true,
			Kind:        "local",
		}, nil
	}

	u, _ := settings.ParseRef(ref)
	if !skills.IsGitURL(u) {
		return registryListEntry{}, fmt.Errorf("invalid ref %q — expected owner/repo, git URL, or absolute path", ref)
	}
	name := settings.DeriveRegistryName(u)

	gfs, err := settings.LoadGlobal()
	if err != nil {
		return registryListEntry{}, fmt.Errorf("loading settings: %w", err)
	}
	for _, existing := range gfs.StandardsExtends {
		eu, _ := settings.ParseRef(existing)
		if settings.DeriveRegistryName(eu) == name {
			// already present — just ensure cloned
			dest := skills.ExtendsHome(name)
			if !dirExists(dest) {
				_ = os.MkdirAll(filepath.Dir(dest), 0o755)
				_ = gitops.Clone(u, dest)
			}
			return registryListEntry{
				Name: name, URL: u,
				SkillsCount: countSkills(filepath.Join(dest, "skills")),
				Cloned:      true, Kind: "extends",
			}, nil
		}
	}

	gfs.StandardsExtends = append(gfs.StandardsExtends, ref)
	if err := settings.SaveGlobal(gfs); err != nil {
		return registryListEntry{}, fmt.Errorf("saving settings: %w", err)
	}

	dest := skills.ExtendsHome(name)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return registryListEntry{}, fmt.Errorf("creating dir: %w", err)
	}
	if err := gitops.Clone(u, dest); err != nil {
		return registryListEntry{}, fmt.Errorf("cloning registry: %w", err)
	}

	return registryListEntry{
		Name:        name,
		URL:         u,
		SkillsCount: countSkills(filepath.Join(dest, "skills")),
		Cloned:      true,
		Kind:        "extends",
	}, nil
}

// performRegistryRemove removes a registry from standards.extends by name.
// Mirrors CLI `grimoire registry remove` behaviour.
func performRegistryRemove(name string) (mcpRegistryRemoveOutput, error) {
	gfs, err := settings.LoadGlobal()
	if err != nil {
		return mcpRegistryRemoveOutput{}, fmt.Errorf("loading settings: %w", err)
	}

	var kept []string
	removed := false
	for _, existing := range gfs.StandardsExtends {
		eu, _ := settings.ParseRef(existing)
		if settings.DeriveRegistryName(eu) == name {
			removed = true
			continue
		}
		kept = append(kept, existing)
	}
	if !removed {
		return mcpRegistryRemoveOutput{}, fmt.Errorf("registry %q not found in standards.extends", name)
	}

	gfs.StandardsExtends = kept
	if err := settings.SaveGlobal(gfs); err != nil {
		return mcpRegistryRemoveOutput{}, fmt.Errorf("saving settings: %w", err)
	}
	return mcpRegistryRemoveOutput{Name: name, Removed: true}, nil
}

// performRegistrySet sets core.registry (the official source).
func performRegistrySet(ref string) (mcpRegistrySetOutput, error) {
	u, _ := settings.ParseRef(ref)
	if !skills.IsGitURL(u) && !filepath.IsAbs(u) {
		return mcpRegistrySetOutput{}, fmt.Errorf("invalid ref %q — expected owner/repo[@version], git URL, or absolute path", ref)
	}
	if filepath.IsAbs(u) {
		if _, err := os.Stat(u); err != nil {
			return mcpRegistrySetOutput{}, fmt.Errorf("local path %q not found", u)
		}
	}
	gfs, err := settings.LoadGlobal()
	if err != nil {
		return mcpRegistrySetOutput{}, fmt.Errorf("loading settings: %w", err)
	}
	gfs.Core.Registry = ref
	if err := settings.SaveGlobal(gfs); err != nil {
		return mcpRegistrySetOutput{}, fmt.Errorf("saving settings: %w", err)
	}
	return mcpRegistrySetOutput{Registry: ref, IsLocal: filepath.IsAbs(u)}, nil
}

func performRegistryUpdate(name string) ([]mcpRegistryUpdateResult, error) {
	r, err := settings.Load(getProjectDir())
	if err != nil {
		return nil, fmt.Errorf("loading settings: %w", err)
	}

	var results []mcpRegistryUpdateResult

	if name == "" || name == "official" {
		status, err := updateOneRegistrySilent("official", r)
		res := mcpRegistryUpdateResult{Name: "official", Status: status}
		if err != nil {
			res.Status = "error"
			res.Error = err.Error()
		}
		results = append(results, res)
	}

	if name != "" && name != "official" {
		status, err := updateOneRegistrySilent(name, r)
		res := mcpRegistryUpdateResult{Name: name, Status: status}
		if err != nil {
			res.Status = "error"
			res.Error = err.Error()
		}
		return append(results, res), nil
	}

	for _, ref := range r.StandardsExtends {
		u, _ := settings.ParseRef(ref)
		n := settings.DeriveRegistryName(u)
		status, err := updateOneRegistrySilent(n, r)
		res := mcpRegistryUpdateResult{Name: n, Status: status}
		if err != nil {
			res.Status = "error"
			res.Error = err.Error()
		}
		results = append(results, res)
	}
	return results, nil
}

func updateOneRegistrySilent(name string, r settings.Resolved) (string, error) {
	var refURL, ver string
	if name == "official" {
		refURL = skills.GrimoireRepoURL()
	} else {
		for _, ref := range r.StandardsExtends {
			u, v := settings.ParseRef(ref)
			if settings.DeriveRegistryName(u) == name {
				refURL, ver = u, v
				break
			}
		}
	}
	if refURL == "" {
		return "error", fmt.Errorf("target %q not configured", name)
	}

	// Local path: verify it exists, skip git ops
	if filepath.IsAbs(refURL) {
		if _, err := os.Stat(refURL); err != nil {
			return "error", fmt.Errorf("local registry %q not found", refURL)
		}
		return "ok", nil
	}

	dest := skills.RegistryHome(name)

	if !dirExists(dest) {
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return "error", fmt.Errorf("mkdir: %w", err)
		}
		if err := gitops.Clone(refURL, dest); err != nil {
			return "error", fmt.Errorf("cloning: %w", err)
		}
		if ver != "" {
			if err := gitops.CheckoutTag(dest, ver); err != nil {
				return "error", fmt.Errorf("checkout %s: %w", ver, err)
			}
		}
		return "cloned", nil
	}

	if ver != "" {
		// pinned — ensure correct tag, skip pull
		if err := gitops.CheckoutTag(dest, ver); err != nil {
			return "error", fmt.Errorf("checkout %s: %w", ver, err)
		}
		return "ok", nil
	}

	if err := gitops.Pull(dest); err != nil {
		return "error", fmt.Errorf("pulling: %w", err)
	}
	return "ok", nil
}

// ── New MCP tool handlers ─────────────────────────────────────────────────────

func toolGrimoireRegistryAdd(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	ref := request.GetString("ref", "")
	entry, err := performRegistryAdd(ref)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(entry)
}

func toolGrimoireRegistryRemove(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	name := request.GetString("name", "")
	out, err := performRegistryRemove(name)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(out)
}

func toolGrimoireRegistrySet(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	ref := request.GetString("ref", "")
	out, err := performRegistrySet(ref)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(out)
}

func toolGrimoireRegistryReset(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	gfs, err := settings.LoadGlobal()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	gfs.Core.Registry = ""
	if err := settings.SaveGlobal(gfs); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(map[string]any{"registry": skills.GrimoireRepo, "reset": true})
}

func toolGrimoireRegistryValidate(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	target := request.GetString("target", "")

	var resolvedTarget string
	if target == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		resolvedTarget = cwd
	} else if filepath.IsAbs(target) {
		resolvedTarget = target
	} else {
		u, _ := settings.ParseRef(target)
		name := settings.DeriveRegistryName(u)
		home := skills.RegistryHome(name)
		if !dirExists(home) {
			return mcp.NewToolResultError(fmt.Sprintf("registry %q not installed", target)), nil
		}
		resolvedTarget = home
	}

	type vcheck struct {
		Name   string `json:"name"`
		Status string `json:"status"`
		Detail string `json:"detail"`
	}
	var checks []vcheck
	allOK := true

	check := func(name, status, detail string) {
		checks = append(checks, vcheck{name, status, detail})
		if status == "error" || status == "warn" {
			allOK = false
		}
	}

	hasMarker := false
	for _, marker := range []string{"skills", "profiles", "presets", "settings.toml"} {
		if _, err := os.Stat(filepath.Join(resolvedTarget, marker)); err == nil {
			hasMarker = true
			break
		}
	}
	if !hasMarker {
		check("registry-markers", "error", "no registry markers found")
	} else {
		check("registry-markers", "ok", "registry structure detected")
	}

	skillsDir := filepath.Join(resolvedTarget, "skills")
	if _, err := os.Stat(skillsDir); err == nil {
		allSkills, _ := skills.ListAllSkills(skillsDir)
		if len(allSkills) == 0 {
			check("skills-structure", "warn", "skills/ found but no skills detected")
		} else {
			missing := 0
			for _, sk := range allSkills {
				if sk.Path == "" {
					missing++
				}
			}
			if missing > 0 {
				check("skills-structure", "warn", fmt.Sprintf("%d skills, %d missing SKILL.md", len(allSkills), missing))
			} else {
				check("skills-structure", "ok", fmt.Sprintf("%d skill(s), all have SKILL.md", len(allSkills)))
			}
		}
	} else {
		check("skills-structure", "skip", "no skills/ directory")
	}

	profilesDir := filepath.Join(resolvedTarget, "profiles")
	if _, err := os.Stat(profilesDir); err == nil {
		entries, _ := os.ReadDir(profilesDir)
		total, invalid := 0, 0
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".toml") {
				continue
			}
			total++
			if _, err := settings.ParseFile(filepath.Join(profilesDir, e.Name())); err != nil {
				invalid++
			}
		}
		if total == 0 {
			check("profiles-structure", "warn", "profiles/ found but empty")
		} else if invalid > 0 {
			check("profiles-structure", "error", fmt.Sprintf("%d/%d profile TOML file(s) failed to parse", invalid, total))
		} else {
			check("profiles-structure", "ok", fmt.Sprintf("%d profile(s), all valid TOML", total))
		}
	} else {
		check("profiles-structure", "skip", "no profiles/ directory")
	}

	presetsDir := filepath.Join(resolvedTarget, "presets")
	if _, err := os.Stat(presetsDir); err == nil {
		presets := skills.ListPresets(resolvedTarget)
		if len(presets) == 0 {
			check("presets-structure", "warn", "presets/ found but no presets detected")
		} else {
			invalid := 0
			for _, p := range presets {
				if _, err := settings.ParseFile(filepath.Join(presetsDir, p, "settings.toml")); err != nil {
					invalid++
				}
			}
			if invalid > 0 {
				check("presets-structure", "error", fmt.Sprintf("%d/%d preset settings.toml file(s) failed to parse", invalid, len(presets)))
			} else {
				check("presets-structure", "ok", fmt.Sprintf("%d preset(s), all have valid settings.toml", len(presets)))
			}
		}
	} else {
		check("presets-structure", "skip", "no presets/ directory")
	}

	rootSettings := filepath.Join(resolvedTarget, "settings.toml")
	if _, err := os.Stat(rootSettings); err == nil {
		if _, err := settings.ParseFile(rootSettings); err != nil {
			check("settings-toml", "error", fmt.Sprintf("settings.toml parse error: %v", err))
		} else {
			check("settings-toml", "ok", "settings.toml is valid TOML")
		}
	} else {
		check("settings-toml", "skip", "no settings.toml")
	}

	return jsonResult(map[string]any{
		"target": resolvedTarget,
		"ok":     allOK,
		"checks": checks,
	})
}

func toolGrimoirePresetList(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	var entries []mcpPresetListEntry
	for _, reg := range skills.AllRegistries() {
		for _, name := range skills.ListPresets(reg.Home) {
			entries = append(entries, mcpPresetListEntry{Name: name, Registry: reg.Name})
		}
	}
	return jsonResult(entries)
}
