package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jeffreytse/grimoire/internal/tui"
)

const (
	ghReleaseAPI = "https://api.github.com/repos/jeffreytse/grimoire/releases/latest"
	ghModule     = "github.com/jeffreytse/grimoire"
)

type ghRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []ghAsset `json:"assets"`
}

type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

var flagSelfUpdateYes bool

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update the grimoire CLI binary to the latest release",
	RunE:  runSelfUpdate,
}

func init() {
	selfUpdateCmd.Flags().BoolVar(&flagSelfUpdateYes, "yes", false, "skip confirmation prompt")
}

func runSelfUpdate(cmd *cobra.Command, args []string) error {
	exePath, err := resolvedExePath()
	if err != nil {
		return fmt.Errorf("locating current binary: %w", err)
	}

	method := detectInstallMethod(exePath)
	switch method {
	case "go":
		fmt.Printf("  Installed via:  go install (%s)\n\n", exePath)
	default:
		fmt.Printf("  Installed via:  binary (%s)\n\n", exePath)
	}

	fmt.Println("Checking for latest release...")
	rel, err := fetchLatestRelease()
	if err != nil {
		return fmt.Errorf("fetching release info: %w", err)
	}

	latestTag := strings.TrimPrefix(rel.TagName, "v")
	currentTag := strings.TrimPrefix(cliVersion, "v")

	if currentTag != "dev" && currentTag == latestTag {
		fmt.Printf("  %s  Already up to date. (%s)\n", tui.IconOK, rel.TagName)
		return nil
	}

	if currentTag == "dev" {
		fmt.Printf("  Current:  dev build\n")
	} else {
		fmt.Printf("  Current:  v%s\n", currentTag)
	}
	fmt.Printf("  Latest:   %s\n\n", rel.TagName)

	if !flagSelfUpdateYes {
		chosen, ok := tui.RunSelect("Update to "+rel.TagName+"?", []string{"Yes", "No"})
		if !ok || chosen == "No" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	switch method {
	case "go":
		return selfUpdateViaGoInstall(rel.TagName)
	default:
		return selfUpdateBinary(exePath, rel)
	}
}

func resolvedExePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(exe)
}

func detectInstallMethod(exePath string) string {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		if out, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
			goPath = strings.TrimSpace(string(out))
		}
	}
	if goPath != "" && filepath.Dir(exePath) == filepath.Join(goPath, "bin") {
		return "go"
	}
	return "binary"
}

func fetchLatestRelease() (*ghRelease, error) {
	req, err := http.NewRequest(http.MethodGet, ghReleaseAPI, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var rel ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}
	return &rel, nil
}

func selfUpdateViaGoInstall(tag string) error {
	pkg := ghModule + "@" + tag
	fmt.Printf("Running: go install %s\n", pkg)
	cmd := exec.Command("go", "install", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go install failed: %w", err)
	}
	fmt.Printf("  %s  Updated to %s\n", tui.IconOK, tag)
	return nil
}

func selfUpdateBinary(exePath string, rel *ghRelease) error {
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

	fmt.Printf("Downloading %s...\n", assetName)

	if runtime.GOOS == "windows" {
		return selfUpdateBinaryWindows(exePath, downloadURL, rel.TagName)
	}
	return selfUpdateBinaryUnix(exePath, downloadURL, rel.TagName)
}

func selfUpdateBinaryUnix(exePath, downloadURL, tag string) error {
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

	fmt.Printf("  %s  Replaced %s (%s → %s)\n", tui.IconOK, exePath, cliVersion, tag)
	return nil
}

func selfUpdateBinaryWindows(exePath, downloadURL, tag string) error {
	newPath := exePath + ".new"

	if err := downloadFile(downloadURL, newPath); err != nil {
		return err
	}

	fmt.Printf("  %s  Downloaded to: %s\n", tui.IconOK, newPath)
	fmt.Printf("\n  Windows cannot replace a running binary.\n")
	fmt.Printf("  Close this terminal, then run:\n\n")
	fmt.Printf("    move \"%s\" \"%s\"\n\n", newPath, exePath)
	return nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url) //nolint:gosec // URL constructed from known GitHub release patterns
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned %d", resp.StatusCode)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		_ = f.Close()
		_ = os.Remove(dest)
		return fmt.Errorf("writing download: %w", err)
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(dest)
		return fmt.Errorf("closing download: %w", err)
	}
	return nil
}
