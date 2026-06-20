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
	"github.com/jeffreytse/grimoire/internal/settings"
	"github.com/jeffreytse/grimoire/internal/skills"
)

// ── Output types ──────────────────────────────────────────────────────────────

type mcpVersionOutput struct {
	CLI      string `json:"cli"`
	Grimoire string `json:"grimoire,omitempty"`
	Home     string `json:"home"`
}

type mcpInitOutput struct {
	Dir              string  `json:"dir"`
	Profile          string  `json:"profile,omitempty"`
	HasReport        bool    `json:"has_report"`
	CoveragePct      float64 `json:"coverage_pct,omitempty"`
	ComplianceStatus string  `json:"compliance_status,omitempty"`
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

func toolGrimoireInit(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	cwd := getProjectDir()
	dir := filepath.Join(cwd, ".grimoire")
	if _, err := os.Stat(dir); err == nil {
		return mcp.NewToolResultError(".grimoire/ already exists — edit .grimoire/settings.toml directly or call grimoire_config_set"), nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("creating .grimoire/: %v", err)), nil
	}
	profile := detect.Profile(cwd)
	if err := writeSettings(dir, profile); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := writeGitignore(dir); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	out := mcpInitOutput{Dir: dir, Profile: profile}
	if r, err := compliance.Load(filepath.Join(cwd, compliance.DefaultReportPath)); err == nil {
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

func toolGrimoireRegistryAdd(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) { //nolint:gocritic
	name := request.GetString("name", "")
	url := request.GetString("url", "")
	entry, err := performRegistryAdd(name, url)
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

	officialRoot := skills.SkillsRoot()
	officialCloned := false
	if _, err := os.Stat(officialRoot); err == nil {
		officialCloned = true
	}
	entries := []registryListEntry{{
		Name:        "official",
		URL:         skills.GrimoireRepoURL(),
		SkillsCount: countSkills(officialRoot),
		Cloned:      officialCloned,
		Version:     skills.GrimoireVersion(),
	}}

	for _, name := range sortedRegistryNames(fs.Registries) {
		rc := fs.Registries[name]
		regHome := skills.RegistryHome(name)
		regRoot := regHome + "/skills"
		cloned := false
		if _, err := os.Stat(regHome); err == nil {
			cloned = true
		}
		entries = append(entries, registryListEntry{
			Name:        name,
			URL:         rc.URL,
			SkillsCount: countSkills(regRoot),
			Cloned:      cloned,
		})
	}
	return entries, nil
}

func performRegistryAdd(name, url string) (registryListEntry, error) {
	if name == skills.OfficialRegistryName {
		return registryListEntry{}, fmt.Errorf("cannot add registry named %q — use grimoire_config_set with key core.source", skills.OfficialRegistryName)
	}
	if !skills.IsGitURL(url) {
		return registryListEntry{}, fmt.Errorf("url must be a git URL (https://, git://, git@): %s", url)
	}

	fs, err := settings.LoadGlobal()
	if err != nil {
		return registryListEntry{}, fmt.Errorf("loading settings: %w", err)
	}
	if fs.Registries == nil {
		fs.Registries = make(map[string]settings.RegistryConfig)
	}
	if _, exists := fs.Registries[name]; exists {
		return registryListEntry{}, fmt.Errorf("registry %q already exists — remove it first with grimoire_registry_remove", name)
	}

	dest := skills.RegistryHome(name)
	if err := gitops.Clone(url, dest); err != nil {
		return registryListEntry{}, fmt.Errorf("cloning registry: %w", err)
	}

	fs.Registries[name] = settings.RegistryConfig{URL: url}
	if err := settings.SaveGlobal(fs); err != nil {
		return registryListEntry{}, fmt.Errorf("saving settings: %w", err)
	}

	return registryListEntry{
		Name:        name,
		URL:         url,
		SkillsCount: countSkills(dest + "/skills"),
		Cloned:      true,
	}, nil
}

func performRegistryRemove(name string) (mcpRegistryRemoveOutput, error) {
	if name == skills.OfficialRegistryName {
		return mcpRegistryRemoveOutput{}, fmt.Errorf("cannot remove the official registry")
	}

	fs, err := settings.LoadGlobal()
	if err != nil {
		return mcpRegistryRemoveOutput{}, fmt.Errorf("loading settings: %w", err)
	}
	if _, exists := fs.Registries[name]; !exists {
		return mcpRegistryRemoveOutput{}, fmt.Errorf("registry %q not found", name)
	}

	delete(fs.Registries, name)
	if err := settings.SaveGlobal(fs); err != nil {
		return mcpRegistryRemoveOutput{}, fmt.Errorf("saving settings: %w", err)
	}

	out := mcpRegistryRemoveOutput{Name: name, Removed: true}
	regHome := skills.RegistryHome(name)
	if _, err := os.Stat(regHome); err == nil {
		if err := os.RemoveAll(regHome); err == nil {
			out.CloneDeleted = true
		}
	}
	return out, nil
}

func performRegistryUpdate(name string) ([]mcpRegistryUpdateResult, error) {
	fs, err := settings.LoadGlobal()
	if err != nil {
		return nil, fmt.Errorf("loading settings: %w", err)
	}

	var results []mcpRegistryUpdateResult

	if name != "" {
		status, err := updateOneRegistrySilent(name, fs)
		r := mcpRegistryUpdateResult{Name: name, Status: status}
		if err != nil {
			r.Status = "error"
			r.Error = err.Error()
		}
		return append(results, r), nil
	}

	// update all: official first, then custom
	for _, n := range append([]string{skills.OfficialRegistryName}, sortedRegistryNames(fs.Registries)...) {
		status, err := updateOneRegistrySilent(n, fs)
		r := mcpRegistryUpdateResult{Name: n, Status: status}
		if err != nil {
			r.Status = "error"
			r.Error = err.Error()
		}
		results = append(results, r)
	}
	return results, nil
}

func updateOneRegistrySilent(name string, fs settings.FileSettings) (string, error) {
	var url string
	if name == skills.OfficialRegistryName {
		url = skills.GrimoireRepoURL()
	} else {
		rc, ok := fs.Registries[name]
		if !ok {
			return "error", fmt.Errorf("registry %q not configured", name)
		}
		url = rc.URL
	}

	dest := skills.RegistryHome(name)
	if _, err := os.Stat(dest); err != nil {
		if err := gitops.Clone(url, dest); err != nil {
			return "error", fmt.Errorf("cloning: %w", err)
		}
		return "cloned", nil
	}

	if err := gitops.Pull(dest); err != nil {
		return "error", fmt.Errorf("pulling: %w", err)
	}
	return "ok", nil
}
