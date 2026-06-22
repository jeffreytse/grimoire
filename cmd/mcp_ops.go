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
	cfg, err := settings.LoadGlobal()
	if err != nil {
		return nil, fmt.Errorf("loading settings: %w", err)
	}

	regs := skills.AllRegistries()
	entries := make([]registryListEntry, 0, len(regs))
	for _, reg := range regs {
		var url, ver string
		for _, rd := range cfg.Registries {
			if rd.Name == reg.Name {
				url, ver = settings.ParseRef(rd.URL)
				if url == "" {
					url = rd.URL
				}
				break
			}
		}
		if url == "" {
			url = skills.GrimoireRepoURL()
		}
		kind := "user"
		if reg.Official {
			kind = "official"
		}
		if filepath.IsAbs(url) {
			kind = "local"
		}
		if ver == "" && reg.Official {
			ver = skills.GrimoireVersion()
		}
		entries = append(entries, registryListEntry{
			Name:        reg.Name,
			URL:         url,
			Version:     ver,
			SkillsCount: countSkills(filepath.Join(reg.Home, "skills")),
			Cloned:      dirExists(reg.Home),
			Kind:        kind,
		})
	}
	return entries, nil
}

// performRegistryAdd adds a named registry to [[registry]] and clones it.
// Mirrors CLI `grimoire registry add` behaviour.
// Name is derived from the URL when not explicitly provided.
func performRegistryAdd(ref string) (registryListEntry, error) {
	u, _ := settings.ParseRef(ref)
	if u == "" {
		u = ref
	}
	if !skills.IsGitURL(u) && !filepath.IsAbs(u) {
		return registryListEntry{}, fmt.Errorf("invalid ref %q — expected owner/repo, git URL, or absolute path", ref)
	}

	name := settings.DeriveRegistryName(u)

	gfs, err := settings.LoadGlobal()
	if err != nil {
		return registryListEntry{}, fmt.Errorf("loading settings: %w", err)
	}

	// Idempotent: if name already exists, just ensure cloned.
	for _, rd := range gfs.Registries {
		if rd.Name == name {
			home := filepath.Join(skills.RegistriesRoot(), name)
			if filepath.IsAbs(u) {
				home = u
			}
			if !dirExists(home) && !filepath.IsAbs(u) {
				_ = os.MkdirAll(filepath.Dir(home), 0o755)
				_ = gitops.Clone(u, home)
			}
			kind := "user"
			if filepath.IsAbs(u) {
				kind = "local"
			}
			return registryListEntry{
				Name:        name,
				URL:         u,
				SkillsCount: countSkills(filepath.Join(home, "skills")),
				Cloned:      dirExists(home),
				Kind:        kind,
			}, nil
		}
	}

	rd := settings.RegistryDef{Name: name, URL: ref, Enabled: true}
	gfs.Registries = append(gfs.Registries, rd)
	if err := settings.SaveGlobal(gfs); err != nil {
		return registryListEntry{}, fmt.Errorf("saving settings: %w", err)
	}

	home := filepath.Join(skills.RegistriesRoot(), name)
	if filepath.IsAbs(u) {
		home = u
		kind := "local"
		return registryListEntry{
			Name:        name,
			URL:         u,
			SkillsCount: countSkills(filepath.Join(home, "skills")),
			Cloned:      true,
			Kind:        kind,
		}, nil
	}

	if err := os.MkdirAll(filepath.Dir(home), 0o755); err != nil {
		return registryListEntry{}, fmt.Errorf("creating dir: %w", err)
	}
	if err := gitops.Clone(u, home); err != nil {
		return registryListEntry{}, fmt.Errorf("cloning registry: %w", err)
	}

	return registryListEntry{
		Name:        name,
		URL:         u,
		SkillsCount: countSkills(filepath.Join(home, "skills")),
		Cloned:      true,
		Kind:        "user",
	}, nil
}

// performRegistryRemove removes a registry from [[registry]] by name.
// Mirrors CLI `grimoire registry remove` behaviour.
func performRegistryRemove(name string) (mcpRegistryRemoveOutput, error) {
	gfs, err := settings.LoadGlobal()
	if err != nil {
		return mcpRegistryRemoveOutput{}, fmt.Errorf("loading settings: %w", err)
	}

	var kept []settings.RegistryDef
	removed := false
	for _, rd := range gfs.Registries {
		if rd.Name == name {
			removed = true
			continue
		}
		kept = append(kept, rd)
	}
	if !removed {
		return mcpRegistryRemoveOutput{}, fmt.Errorf("registry %q not found in [[registry]]", name)
	}

	gfs.Registries = kept
	if err := settings.SaveGlobal(gfs); err != nil {
		return mcpRegistryRemoveOutput{}, fmt.Errorf("saving settings: %w", err)
	}
	return mcpRegistryRemoveOutput{Name: name, Removed: true}, nil
}

// performRegistrySet sets the official registry URL via the [[registry]] model.
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
	for i, rd := range gfs.Registries {
		if rd.Official {
			gfs.Registries[i].URL = ref
			if err := settings.SaveGlobal(gfs); err != nil {
				return mcpRegistrySetOutput{}, fmt.Errorf("saving settings: %w", err)
			}
			return mcpRegistrySetOutput{Registry: ref, IsLocal: filepath.IsAbs(u)}, nil
		}
	}
	gfs.Registries = append(gfs.Registries, settings.RegistryDef{
		Name:     "official",
		URL:      ref,
		Official: true,
		Priority: 100,
		Enabled:  true,
	})
	if err := settings.SaveGlobal(gfs); err != nil {
		return mcpRegistrySetOutput{}, fmt.Errorf("saving settings: %w", err)
	}
	return mcpRegistrySetOutput{Registry: ref, IsLocal: filepath.IsAbs(u)}, nil
}

func performRegistryUpdate(name string) ([]mcpRegistryUpdateResult, error) {
	cfg, err := settings.LoadGlobal()
	if err != nil {
		return nil, fmt.Errorf("loading settings: %w", err)
	}

	updateOne := func(n string) mcpRegistryUpdateResult {
		status, err := updateOneRegistrySilent(n, cfg)
		res := mcpRegistryUpdateResult{Name: n, Status: status}
		if err != nil {
			res.Status = "error"
			res.Error = err.Error()
		}
		return res
	}

	if name != "" {
		return []mcpRegistryUpdateResult{updateOne(name)}, nil
	}

	if len(cfg.Registries) == 0 {
		return []mcpRegistryUpdateResult{updateOne("official")}, nil
	}

	var results []mcpRegistryUpdateResult
	for _, rd := range cfg.Registries {
		if rd.Enabled {
			results = append(results, updateOne(rd.Name))
		}
	}
	return results, nil
}

func updateOneRegistrySilent(name string, cfg settings.FileSettings) (string, error) {
	var refURL, ver string
	var dest string

	for _, rd := range cfg.Registries {
		if rd.Name == name {
			u, v := settings.ParseRef(rd.URL)
			if u == "" {
				u = rd.URL
			}
			refURL, ver = u, v
			if filepath.IsAbs(u) {
				dest = u
			} else {
				dest = skills.RegistryHome(name)
			}
			break
		}
	}

	if refURL == "" {
		if name == skills.OfficialRegistryName {
			refURL = skills.GrimoireRepoURL()
			dest = skills.OfficialRegistryHome()
		} else {
			return "error", fmt.Errorf("target %q not configured", name)
		}
	}
	// Local path: verify it exists, skip git ops
	if filepath.IsAbs(refURL) {
		if _, err := os.Stat(refURL); err != nil {
			return "error", fmt.Errorf("local registry %q not found", refURL)
		}
		return "ok", nil
	}

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
	var kept []settings.RegistryDef
	for _, rd := range gfs.Registries {
		if !rd.Official {
			kept = append(kept, rd)
		}
	}
	gfs.Registries = kept
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
