package compliance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DefaultReportDir is the directory where compliance reports are written.
const DefaultReportDir = ".grimoire/reports"

// WriteReport writes r to <projectDir>/.grimoire/reports/compliance-<ts>.json
// and updates compliance-latest.json to point to it.
// Returns the path of the timestamped report file.
func WriteReport(r *Report, projectDir string) (string, error) {
	dir := filepath.Join(projectDir, DefaultReportDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("creating report dir: %w", err)
	}

	ts := time.Now().UTC().Format("20060102T150405Z")
	tsPath := filepath.Join(dir, "compliance-"+ts+".json")
	latestPath := filepath.Join(dir, "compliance-latest.json")

	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshalling report: %w", err)
	}

	if err := os.WriteFile(tsPath, data, 0o644); err != nil {
		return "", fmt.Errorf("writing report: %w", err)
	}

	// Update compliance-latest.json: symlink preferred, file copy as Windows fallback.
	_ = os.Remove(latestPath)
	if err := os.Symlink(tsPath, latestPath); err != nil {
		_ = os.WriteFile(latestPath, data, 0o644)
	}

	return tsPath, nil
}

// TouchTimestamp updates the Timestamp field of compliance-latest.json to now.
// Called when a check cycle is skipped (cache hit) so the report still shows
// the current verification time without re-running the AI.
func TouchTimestamp(projectDir string) error {
	latestPath := filepath.Join(projectDir, DefaultReportDir, "compliance-latest.json")
	data, err := os.ReadFile(latestPath)
	if err != nil {
		return nil // no report yet — nothing to touch
	}
	var r Report
	if err := json.Unmarshal(data, &r); err != nil {
		return nil // malformed — skip silently
	}
	r.Timestamp = time.Now().UTC().Format(time.RFC3339)
	updated, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	// os.WriteFile follows the symlink, updating the underlying timestamped file.
	return os.WriteFile(latestPath, updated, 0o644)
}
