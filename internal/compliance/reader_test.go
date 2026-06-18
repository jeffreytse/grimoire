package compliance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// writeReport writes a Report as JSON to path (creating parent dirs).
func writeReport(t *testing.T, path string, r *Report) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("writeReport mkdir: %v", err)
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("writeReport marshal: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("writeReport write: %v", err)
	}
}

// minimalReport returns a populated but minimal Report for use in tests.
func minimalReport() Report {
	return Report{
		Version:   "1",
		Timestamp: "2026-06-17T00:00:00Z",
		Mode:      "full",
		Result:    "pass",
		Coverage: Coverage{
			OverallPct: 85.0,
			Practices: PracticeSummary{
				Total:       10,
				Passing:     8,
				Partial:     1,
				Failing:     1,
				CoveragePct: 85.0,
			},
		},
		Threshold: Threshold{
			Required: 80.0,
			Actual:   85.0,
			Status:   "pass",
		},
	}
}

// ── Load with explicit path ──────────────────────────────────────────────────

func TestLoad_ValidReport_ReturnsReport(t *testing.T) {
	dir := t.TempDir()
	reportPath := filepath.Join(dir, "compliance-latest.json")
	want := minimalReport()
	writeReport(t, reportPath, &want)

	got, err := Load(reportPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Version != want.Version {
		t.Errorf("Version = %q; want %q", got.Version, want.Version)
	}
	if got.Coverage.OverallPct != want.Coverage.OverallPct {
		t.Errorf("OverallPct = %v; want %v", got.Coverage.OverallPct, want.Coverage.OverallPct)
	}
	if got.Threshold.Status != want.Threshold.Status {
		t.Errorf("Threshold.Status = %q; want %q", got.Threshold.Status, want.Threshold.Status)
	}
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no-such-file.json")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_MalformedJSON_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("{not valid json"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

// ── Diagnostics ──────────────────────────────────────────────────────────────

func TestLoad_Diagnostics_Preserved(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.json")
	r := minimalReport()
	r.Diagnostics = []Diagnostic{
		{
			URI:      "main.go",
			Severity: 1,
			Message:  "missing error check",
			Range: DiagRange{
				Start: Position{Line: 10, Character: 4},
				End:   Position{Line: 10, Character: 20},
			},
		},
		{
			URI:      "util.go",
			Severity: 2,
			Message:  "consider using constants",
		},
	}
	writeReport(t, path, &r)

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got.Diagnostics) != 2 {
		t.Fatalf("expected 2 diagnostics, got %d", len(got.Diagnostics))
	}
	if got.Diagnostics[0].Severity != 1 {
		t.Errorf("Severity = %d; want 1", got.Diagnostics[0].Severity)
	}
	if got.Diagnostics[0].Range.Start.Line != 10 {
		t.Errorf("Range.Start.Line = %d; want 10", got.Diagnostics[0].Range.Start.Line)
	}
	if got.Diagnostics[1].Severity != 2 {
		t.Errorf("second Severity = %d; want 2", got.Diagnostics[1].Severity)
	}
}

// ── Load with empty path (uses DefaultReportPath) ────────────────────────────

func TestLoad_EmptyPath_UsesDefault_MissingReturnsError(t *testing.T) {
	// Change working directory to a temp dir so DefaultReportPath doesn't exist.
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	defer func() { _ = os.Chdir(orig) }()

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("Chdir: %v", err)
	}

	_, err = Load("") // should try DefaultReportPath which doesn't exist
	if err == nil {
		t.Fatal("expected error when default report path does not exist")
	}
}

func TestLoad_EmptyPath_UsesDefault_WhenPresent(t *testing.T) {
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	defer func() { _ = os.Chdir(orig) }()

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("Chdir: %v", err)
	}

	r := minimalReport()
	writeReport(t, DefaultReportPath, &r)

	got, err := Load("")
	if err != nil {
		t.Fatalf("Load with empty path: %v", err)
	}
	if got.Threshold.Status != "pass" {
		t.Errorf("Threshold.Status = %q; want pass", got.Threshold.Status)
	}
}

// ── PracticeSummary fields ───────────────────────────────────────────────────

func TestLoad_PracticeSummaryFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.json")
	r := minimalReport()
	r.Coverage.Practices = PracticeSummary{
		Total:       20,
		Passing:     15,
		Partial:     3,
		Failing:     2,
		CoveragePct: 75.0,
	}
	writeReport(t, path, &r)

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	ps := got.Coverage.Practices
	if ps.Total != 20 {
		t.Errorf("Total = %d; want 20", ps.Total)
	}
	if ps.Passing != 15 {
		t.Errorf("Passing = %d; want 15", ps.Passing)
	}
	if ps.Partial != 3 {
		t.Errorf("Partial = %d; want 3", ps.Partial)
	}
	if ps.Failing != 2 {
		t.Errorf("Failing = %d; want 2", ps.Failing)
	}
	if ps.CoveragePct != 75.0 {
		t.Errorf("CoveragePct = %v; want 75.0", ps.CoveragePct)
	}
}

// ── Edge cases ───────────────────────────────────────────────────────────────

func TestLoad_EmptyDiagnostics_ReturnsNilSlice(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.json")
	r := minimalReport()
	// no Diagnostics field → should unmarshal as nil/empty
	writeReport(t, path, &r)

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	// nil or empty both acceptable — the key test is no panic
	_ = got.Diagnostics
}

func TestLoad_EmptyFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for empty file")
	}
}

func TestLoad_FailThreshold(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fail.json")
	r := minimalReport()
	r.Threshold.Status = "fail"
	r.Coverage.OverallPct = 50.0
	writeReport(t, path, &r)

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Threshold.Status != "fail" {
		t.Errorf("Threshold.Status = %q; want fail", got.Threshold.Status)
	}
}
