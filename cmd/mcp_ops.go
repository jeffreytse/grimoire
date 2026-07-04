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
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/jeffreytse/grimoire/internal/compliance"
	"github.com/jeffreytse/grimoire/internal/config"
	grimctx "github.com/jeffreytse/grimoire/internal/context"
	"github.com/jeffreytse/grimoire/internal/detect"
	gitops "github.com/jeffreytse/grimoire/internal/git"
	"github.com/jeffreytse/grimoire/internal/manifest"
	"github.com/jeffreytse/grimoire/internal/profiles"
	"github.com/jeffreytse/grimoire/internal/skills"
)

// ── MCP-only output types ─────────────────────────────────────────────────────

type mcpPackageSetOutput struct {
	Package string `json:"package"`
	IsLocal bool   `json:"is_local"`
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

type mcpPackageRemoveOutput struct {
	Name         string `json:"name"`
	Removed      bool   `json:"removed"`
	CloneDeleted bool   `json:"clone_deleted"`
}

type mcpPackageUpdateResult struct {
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
	if _, err := scaffoldManifest(manifest.ProjectPath(cwd), filepath.Base(cwd), cfg.Profile, cfg.Threshold, cfg.MaxErrors); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	out := mcpInitOutput{Dir: dir, Profile: cfg.Profile}
	if cfg.Profile != "" {
		if refs, err := profiles.ResolveEffectiveSkills([]string{cfg.Profile}, cwd, nil, nil); err == nil {
			out.SkillCount = len(refs)
			names := make([]string, len(refs))
			for i, sk := range refs {
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

func toolGrimoirePackageList(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	entries, err := collectPackageList()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(entries)
}

func toolGrimoirePackageUpdate(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	name := request.GetString("name", "")
	results, err := performPackageUpdate(name)
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

// ── Package ───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────

func collectPackageList() ([]packageListEntry, error) {
	cfg, err := config.LoadGlobal()
	if err != nil {
		return nil, fmt.Errorf("loading settings: %w", err)
	}

	regs := skills.AllPackages()
	entries := make([]packageListEntry, 0, len(regs))
	for _, reg := range regs {
		var url, ver string
		for _, rd := range cfg.Packages {
			if rd.Name == reg.Name {
				url, ver = config.ParseRef(rd.URL)
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
		entries = append(entries, packageListEntry{
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

// performPackageAdd adds a named package to [[package]] and clones it.
// Mirrors CLI `grimoire package add` behaviour.
// Name is derived from the URL when not explicitly provided.
func performPackageAdd(ref string) (packageListEntry, error) {
	u, _ := config.ParseRef(ref)
	if u == "" {
		u = ref
	}
	if !skills.IsGitURL(u) && !filepath.IsAbs(u) {
		return packageListEntry{}, fmt.Errorf("invalid grimoire-ref %q — expected owner/repo[@version], git URL, or absolute path", ref)
	}

	name := config.DerivePackageName(u)

	gfs, err := config.LoadGlobal()
	if err != nil {
		return packageListEntry{}, fmt.Errorf("loading config: %w", err)
	}

	// Idempotent: if name already exists, just ensure cloned.
	for _, rd := range gfs.Packages {
		if rd.Name != name {
			continue
		}
		home := skills.PackageHome(name)
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
		return packageListEntry{
			Name:        name,
			URL:         u,
			SkillsCount: countSkills(filepath.Join(home, "skills")),
			Cloned:      dirExists(home),
			Kind:        kind,
		}, nil
	}

	rd := config.PackageDef{Name: name, URL: ref, Enabled: true}
	gfs.Packages = append(gfs.Packages, rd)
	if err := config.SaveGlobal(gfs); err != nil {
		return packageListEntry{}, fmt.Errorf("saving settings: %w", err)
	}

	home := skills.PackageHome(name)
	if filepath.IsAbs(u) {
		home = u
		kind := "local"
		return packageListEntry{
			Name:        name,
			URL:         u,
			SkillsCount: countSkills(filepath.Join(home, "skills")),
			Cloned:      true,
			Kind:        kind,
		}, nil
	}

	if err := os.MkdirAll(filepath.Dir(home), 0o755); err != nil {
		return packageListEntry{}, fmt.Errorf("creating dir: %w", err)
	}
	if err := gitops.Clone(u, home); err != nil {
		return packageListEntry{}, fmt.Errorf("cloning package: %w", err)
	}

	return packageListEntry{
		Name:        name,
		URL:         u,
		SkillsCount: countSkills(filepath.Join(home, "skills")),
		Cloned:      true,
		Kind:        "user",
	}, nil
}

// performPackageRemove removes a package from [[package]] by name.
// Mirrors CLI `grimoire package remove` behaviour.
func performPackageRemove(name string) (mcpPackageRemoveOutput, error) {
	gfs, err := config.LoadGlobal()
	if err != nil {
		return mcpPackageRemoveOutput{}, fmt.Errorf("loading settings: %w", err)
	}

	var kept []config.PackageDef
	removed := false
	for _, rd := range gfs.Packages {
		if rd.Name == name {
			removed = true
			continue
		}
		kept = append(kept, rd)
	}
	if !removed {
		return mcpPackageRemoveOutput{}, fmt.Errorf("package %q not found in [[package]]", name)
	}

	gfs.Packages = kept
	if err := config.SaveGlobal(gfs); err != nil {
		return mcpPackageRemoveOutput{}, fmt.Errorf("saving settings: %w", err)
	}
	return mcpPackageRemoveOutput{Name: name, Removed: true}, nil
}

// performPackageSet sets the official package URL via the [[package]] model.
func performPackageSet(ref string) (mcpPackageSetOutput, error) {
	u, _ := config.ParseRef(ref)
	if !skills.IsGitURL(u) && !filepath.IsAbs(u) {
		return mcpPackageSetOutput{}, fmt.Errorf("invalid grimoire-ref %q — expected owner/repo[@version], git URL, or absolute path", ref)
	}
	if filepath.IsAbs(u) {
		if _, err := os.Stat(u); err != nil {
			return mcpPackageSetOutput{}, fmt.Errorf("local path %q not found", u)
		}
	}
	gfs, err := config.LoadGlobal()
	if err != nil {
		return mcpPackageSetOutput{}, fmt.Errorf("loading settings: %w", err)
	}
	for i, rd := range gfs.Packages {
		if rd.Official {
			gfs.Packages[i].URL = ref
			if err := config.SaveGlobal(gfs); err != nil {
				return mcpPackageSetOutput{}, fmt.Errorf("saving settings: %w", err)
			}
			return mcpPackageSetOutput{Package: ref, IsLocal: filepath.IsAbs(u)}, nil
		}
	}
	gfs.Packages = append(gfs.Packages, config.PackageDef{
		Name:     "official",
		URL:      ref,
		Official: true,
		Priority: 100,
		Enabled:  true,
	})
	if err := config.SaveGlobal(gfs); err != nil {
		return mcpPackageSetOutput{}, fmt.Errorf("saving settings: %w", err)
	}
	return mcpPackageSetOutput{Package: ref, IsLocal: filepath.IsAbs(u)}, nil
}

func performPackageUpdate(name string) ([]mcpPackageUpdateResult, error) {
	cfg, err := config.LoadGlobal()
	if err != nil {
		return nil, fmt.Errorf("loading settings: %w", err)
	}

	updateOne := func(n string) mcpPackageUpdateResult {
		status, err := updateOnePackageSilent(n, &cfg)
		res := mcpPackageUpdateResult{Name: n, Status: status}
		if err != nil {
			res.Status = "error"
			res.Error = err.Error()
		}
		return res
	}

	if name != "" {
		return []mcpPackageUpdateResult{updateOne(name)}, nil
	}

	if len(cfg.Packages) == 0 {
		return []mcpPackageUpdateResult{updateOne("official")}, nil
	}

	var results []mcpPackageUpdateResult
	for _, rd := range cfg.Packages {
		if rd.Enabled {
			results = append(results, updateOne(rd.Name))
		}
	}
	return results, nil
}

func updateOnePackageSilent(name string, cfg *config.FileConfig) (string, error) {
	var refURL, ver string
	var dest string

	for _, rd := range cfg.Packages {
		if rd.Name != name {
			continue
		}
		u, v := config.ParseRef(rd.URL)
		if u == "" {
			u = rd.URL
		}
		refURL, ver = u, v
		if filepath.IsAbs(u) {
			dest = u
		} else {
			dest = skills.PackageHome(config.DeriveVersionedName(u, v))
		}
		break
	}

	if refURL == "" {
		if name != skills.OfficialPackageName && name != skills.OfficialPackageDerivedName() {
			return "error", fmt.Errorf("target %q not configured", name)
		}
		refURL = skills.GrimoireRepoURL()
		dest = skills.OfficialPackageHome()
	}
	// Local path: verify it exists, skip git ops
	if filepath.IsAbs(refURL) {
		if _, err := os.Stat(refURL); err != nil {
			return "error", fmt.Errorf("local package %q not found", refURL)
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

func toolGrimoirePackageAdd(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	ref := request.GetString("ref", "")
	entry, err := performPackageAdd(ref)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(entry)
}

func toolGrimoirePackageRemove(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	name := request.GetString("name", "")
	out, err := performPackageRemove(name)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(out)
}

func toolGrimoirePackageEnable(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	name := request.GetString("name", "")
	if err := setPackageEnabled(name, true); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(map[string]any{"name": name, "enabled": true})
}

func toolGrimoirePackageDisable(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	name := request.GetString("name", "")
	if err := setPackageEnabled(name, false); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(map[string]any{"name": name, "enabled": false})
}

func toolGrimoirePackageSet(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	ref := request.GetString("ref", "")
	out, err := performPackageSet(ref)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(out)
}

func toolGrimoirePackageReset(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	gfs, err := config.LoadGlobal()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	var kept []config.PackageDef
	for _, rd := range gfs.Packages {
		if !rd.Official {
			kept = append(kept, rd)
		}
	}
	gfs.Packages = kept
	if err := config.SaveGlobal(gfs); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(map[string]any{"package": skills.GrimoireRepo, "reset": true})
}

func toolGrimoirePackageValidate(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	target := request.GetString("target", "")

	var resolvedTarget string
	switch {
	case target == "":
		cwd, err := os.Getwd()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		resolvedTarget = cwd
	case filepath.IsAbs(target):
		resolvedTarget = target
	default:
		u, ver := config.ParseRef(target)
		name := config.DeriveVersionedName(u, ver)
		home := skills.PackageHome(name)
		if !dirExists(home) {
			return mcp.NewToolResultError(fmt.Sprintf("package %q not installed", target)), nil
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
	for _, marker := range []string{"skills", "profiles", "grimoire.toml"} {
		if _, err := os.Stat(filepath.Join(resolvedTarget, marker)); err == nil {
			hasMarker = true
			break
		}
	}
	if !hasMarker {
		check("package-markers", "error", "no package markers found")
	} else {
		check("package-markers", "ok", "package structure detected")
	}

	skillsDir := filepath.Join(resolvedTarget, "skills")
	if _, err := os.Stat(skillsDir); err == nil {
		allSkills, _ := skills.ListAllSkills(skillsDir)
		if len(allSkills) == 0 {
			check("skills-structure", "warn", "skills/ found but no skills detected")
		} else {
			missing := 0
			for i := range allSkills {
				sk := allSkills[i]
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
			if _, err := config.ParseFile(filepath.Join(profilesDir, e.Name())); err != nil {
				invalid++
			}
		}
		switch {
		case total == 0:
			check("profiles-structure", "warn", "profiles/ found but empty")
		case invalid > 0:
			check("profiles-structure", "error", fmt.Sprintf("%d/%d profile TOML file(s) failed to parse", invalid, total))
		default:
			check("profiles-structure", "ok", fmt.Sprintf("%d profile(s), all valid TOML", total))
		}
	} else {
		check("profiles-structure", "skip", "no profiles/ directory")
	}

	rootConfig := filepath.Join(resolvedTarget, "grimoire.toml")
	if _, err := os.Stat(rootConfig); err == nil {
		if _, err := config.ParseFile(rootConfig); err != nil {
			check("grimoire-toml", "error", fmt.Sprintf("grimoire.toml parse error: %v", err))
		} else {
			check("grimoire-toml", "ok", "grimoire.toml is valid TOML")
		}
	} else {
		check("grimoire-toml", "skip", "no grimoire.toml")
	}

	return jsonResult(map[string]any{
		"target": resolvedTarget,
		"ok":     allOK,
		"checks": checks,
	})
}

// ── Status / Search / Info ────────────────────────────────────────────────────

type mcpStatusOutput struct {
	ProjectDir      string  `json:"project_dir"`
	Profile         string  `json:"profile"`
	SkillCount      int     `json:"skill_count"`
	PackageCount    int     `json:"package_count"`
	HasReport       bool    `json:"has_report"`
	LastCheckAge    string  `json:"last_check_age,omitempty"`
	LastCheckResult string  `json:"last_check_result,omitempty"`
	LastCheckPct    float64 `json:"last_check_pct,omitempty"`
	Stale           bool    `json:"stale"`
	StalenessDays   int     `json:"staleness_days"`
	ErrorCount      int     `json:"error_count"`
	WarningCount    int     `json:"warning_count"`
}

type mcpSkillMatch struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Domain      string   `json:"domain,omitempty"`
	Subdomain   string   `json:"subdomain,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Package     string   `json:"package,omitempty"`
}

type mcpSkillInfo struct {
	Name          string            `json:"name"`
	Description   string            `json:"description,omitempty"`
	Domain        string            `json:"domain,omitempty"`
	Subdomain     string            `json:"subdomain,omitempty"`
	Package       string            `json:"package,omitempty"`
	Version       string            `json:"version,omitempty"`
	Authors       []string          `json:"authors,omitempty"`
	License       string            `json:"license,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
	Compatibility []string          `json:"compatibility,omitempty"`
	Dependencies  map[string]string `json:"dependencies,omitempty"`
	Path          string            `json:"path,omitempty"`
}

func toolGrimoireStatus(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	cwd := getProjectDir()

	ctx := grimctx.Detect(cwd)

	regs := skills.AllSkillsPackages()
	totalSkills := 0
	for _, reg := range regs {
		if all, err := skills.WalkSkills(reg.Root); err == nil {
			totalSkills += len(all)
		}
	}

	cfg, _ := config.Load(cwd)
	stalenessDays := cfg.StalenessDays
	if stalenessDays == 0 {
		stalenessDays = 7
	}

	out := mcpStatusOutput{
		ProjectDir:    cwd,
		Profile:       ctx.Profile,
		SkillCount:    totalSkills,
		PackageCount:  len(regs),
		StalenessDays: stalenessDays,
	}

	reportPath := resolvedReportPath(cwd)
	if fi, err := os.Stat(reportPath); err == nil {
		out.HasReport = true
		age := time.Since(fi.ModTime())
		out.LastCheckAge = formatAge(age)
		out.Stale = age.Hours()/24 > float64(stalenessDays)

		if report, loadErr := compliance.Load(reportPath); loadErr == nil {
			out.LastCheckResult = report.Threshold.Status
			out.LastCheckPct = report.Coverage.OverallPct
			out.ErrorCount = len(filterBySeverity(report.Diagnostics, 1))
			out.WarningCount = len(filterBySeverity(report.Diagnostics, 2))
		}
	}

	return jsonResult(out)
}

func toolGrimoireSearch(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	query := strings.ToLower(strings.TrimSpace(request.GetString("query", "")))

	regs := skills.AllSkillsPackages()
	if len(regs) == 0 {
		return mcp.NewToolResultError("no packages installed — run grimoire_update first"), nil
	}

	all, _, err := skills.ListAllSkillsFromPackages(regs)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	matches := make([]mcpSkillMatch, 0)
	for i := range all {
		if query == "" || matchesQuery(&all[i], query) {
			matches = append(matches, mcpSkillMatch{
				Name:        all[i].Name,
				Description: all[i].Description,
				Domain:      all[i].Domain,
				Subdomain:   all[i].Subdomain,
				Tags:        all[i].Tags,
				Package:     all[i].Package,
			})
		}
	}
	return jsonResult(matches)
}

func toolGrimoireInfo(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	name := strings.TrimSpace(request.GetString("skill", ""))
	if name == "" {
		return mcp.NewToolResultError("skill name is required"), nil
	}

	regs := skills.AllSkillsPackages()
	if len(regs) == 0 {
		return mcp.NewToolResultError("no packages installed — run grimoire_update first"), nil
	}

	_, src, err := resolveSkillFromPackages(regs, name)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	all, err := skills.WalkSkills(src.Root)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	for i := range all {
		sk := all[i]
		if sk.Name == name {
			return jsonResult(mcpSkillInfo{
				Name:          sk.Name,
				Description:   sk.Description,
				Domain:        sk.Domain,
				Subdomain:     sk.Subdomain,
				Package:       sk.Package,
				Version:       sk.Version,
				Authors:       sk.Authors,
				License:       sk.License,
				Tags:          sk.Tags,
				Compatibility: sk.Compatibility,
				Dependencies:  sk.Dependencies,
				Path:          sk.Path,
			})
		}
	}
	return mcp.NewToolResultError(fmt.Sprintf("skill %q not found", name)), nil
}
