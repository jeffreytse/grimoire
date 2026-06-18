package compliance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const DefaultReportPath = ".grimoire/reports/compliance-latest.json"

// Load reads and parses a compliance report from path.
// If path is empty, DefaultReportPath is used.
func Load(path string) (*Report, error) {
	if path == "" {
		path = DefaultReportPath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no report at %s — run /check-best-practice-compliance first", filepath.Clean(path))
		}
		return nil, fmt.Errorf("reading report: %w", err)
	}

	var r Report
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("parsing report: %w", err)
	}
	return &r, nil
}
