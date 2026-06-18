package compliance

// Report is the top-level structure of compliance-latest.json.
type Report struct {
	Version     string       `json:"version"`
	Timestamp   string       `json:"timestamp"`
	Mode        string       `json:"mode"`
	Scope       string       `json:"scope"`
	Result      string       `json:"result"`
	Coverage    Coverage     `json:"coverage"`
	Threshold   Threshold    `json:"threshold"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type Coverage struct {
	OverallPct float64         `json:"overall_pct"`
	Practices  PracticeSummary `json:"practices"`
	Criteria   CriteriaSummary `json:"criteria"`
}

type PracticeSummary struct {
	Total       int     `json:"total"`
	Passing     int     `json:"passing"`
	Partial     int     `json:"partial"`
	Failing     int     `json:"failing"`
	CoveragePct float64 `json:"coverage_pct"`
}

type CriteriaSummary struct {
	Total       int     `json:"total"`
	Passing     int     `json:"passing"`
	Failing     int     `json:"failing"`
	Suppressed  int     `json:"suppressed"`
	CoveragePct float64 `json:"coverage_pct"`
}

type Threshold struct {
	Required float64 `json:"required"`
	Actual   float64 `json:"actual"`
	Status   string  `json:"status"`
}

// Diagnostic follows the LSP Diagnostic schema.
// Severity: 1=Error, 2=Warning, 3=Information, 4=Hint.
type Diagnostic struct {
	URI       string    `json:"uri"`
	Range     DiagRange `json:"range"`
	Severity  int       `json:"severity"`
	Code      string    `json:"code"`
	Source    string    `json:"source"`
	Message   string    `json:"message"`
	Practice  string    `json:"practice"`
	Criterion string    `json:"criterion"`
	Status    string    `json:"status"`
}

type DiagRange struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}
