package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeffreytse/grimoire/internal/compliance"
)

func sandboxHome(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("GRIMOIRE_SYSTEM_CONFIG", filepath.Join(t.TempDir(), "none.toml"))
}

func TestCollectStatus_NoReport(t *testing.T) {
	sandboxHome(t)
	st := collectStatus(t.TempDir())
	if st.LastCheck != nil || st.Result != "" || st.Errors != nil {
		t.Errorf("no report should leave check fields empty: %+v", st)
	}
	if st.StalenessDays != 7 {
		t.Errorf("default staleness = %d, want 7", st.StalenessDays)
	}
}

func TestCollectStatus_WithReport(t *testing.T) {
	sandboxHome(t)
	dir := t.TempDir()
	path := filepath.Join(dir, compliance.DefaultReportPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	report := `{"version":"1","coverage":{"overall_pct":87.5},"threshold":{"required":80,"actual":87.5,"status":"pass"},` +
		`"diagnostics":[{"severity":1,"message":"a"},{"severity":2,"message":"b"},{"severity":2,"message":"c"}]}`
	if err := os.WriteFile(path, []byte(report), 0o644); err != nil {
		t.Fatal(err)
	}
	st := collectStatus(dir)
	if st.LastCheck == nil || st.Stale == nil || *st.Stale {
		t.Fatalf("fresh report should be found and not stale: %+v", st)
	}
	if st.Result != "pass" || st.CoveragePct == nil || *st.CoveragePct != 87.5 {
		t.Errorf("result/coverage = %q/%v", st.Result, st.CoveragePct)
	}
	if st.Errors == nil || *st.Errors != 1 || *st.Warnings != 2 {
		t.Errorf("errors/warnings wrong: %+v", st)
	}
}
