package compliance

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

//go:embed templates/report.html
var reportTemplateContent string

var reportTmpl = template.Must(template.New("report").Funcs(template.FuncMap{
	"lower": strings.ToLower,
}).Parse(reportTemplateContent))

// reportTemplateData is the data model passed to templates/report.html.
type reportTemplateData struct {
	Scope       string
	Mode        string
	TimeStr     string
	BadgeClass  string
	BadgeText   string
	CovPct      float64
	Required    float64
	PracPassing int
	PracTotal   int
	ErrCount    int
	WarnCount   int
	DiagTotal   int
	Practices   []pracRow
	Errors      []diagRow
	Warnings    []diagRow
	Infos       []diagRow
	CLIVersion  string
	LiveMode    bool
}

type criterionRow struct {
	Name  string
	Count int
}

type pracRow struct {
	IconClass       string
	IconChar        string
	Name            string
	Passing         int
	Total           int
	BarClass        string
	BarPct          float64
	FailingCriteria []criterionRow
	PassingCriteria []criterionRow
}

type codeLine struct {
	Num         int
	Text        string
	Highlighted bool
	// Inline char split — set on start line only when start.char != end.char (same-line range).
	// Zero-length points (char:0,char:0) are left as plain whole-line highlights.
	HasCharHL bool
	PreHl     string
	InHl      string
	PostHl    string
	// InlineCriterion is non-empty on the last highlighted line only; rendered as
	// a GitHub-style inline annotation block immediately after that line.
	InlineCriterion string
}

type diagRow struct {
	Message       string
	Loc           string
	Practice      string
	Criterion     string
	SeverityLabel string
	CodeLines     []codeLine
}

// RenderHTMLReport renders the compliance report HTML from the JSON at jsonPath.
// When liveMode is true the SSE status badge is included.
func RenderHTMLReport(jsonPath, version, projectDir string, liveMode bool) ([]byte, error) {
	r, err := Load(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("loading report for HTML: %w", err)
	}

	errors := filterDiagsBySeverity(r.Diagnostics, 1)
	warnings := filterDiagsBySeverity(r.Diagnostics, 2)
	infos := filterDiagsBySeverity(r.Diagnostics, 3)

	ts := r.Timestamp
	timeStr := ""
	if t, err := time.Parse(time.RFC3339, ts); err == nil {
		timeStr = t.Local().Format("2006-01-02 15:04:05")
	} else if len(ts) >= 19 {
		timeStr = ts[:10] + " " + ts[11:19]
	} else if len(ts) >= 10 {
		timeStr = ts[:10]
	}

	failingCritByPractice := map[string]map[string]int{}
	pracHasFailingDiag := map[string]bool{}
	// Primary: criteria_matrix (natural LLM output).
	passingCritByPractice := map[string][]string{}
	for practice, entry := range r.CriteriaMatrix {
		passingCritByPractice[practice] = entry.Pass
	}
	for i := range r.Diagnostics {
		d := r.Diagnostics[i]
		if d.Severity <= 2 && d.Practice != "" {
			pracHasFailingDiag[d.Practice] = true
			if d.Criterion != "" {
				if failingCritByPractice[d.Practice] == nil {
					failingCritByPractice[d.Practice] = map[string]int{}
				}
				failingCritByPractice[d.Practice][d.Criterion]++
			}
		}
		// Fallback: severity-4 pass diagnostics.
		if d.Severity == 4 && d.Status == "pass" && d.Practice != "" && d.Criterion != "" {
			if _, ok := passingCritByPractice[d.Practice]; !ok {
				passingCritByPractice[d.Practice] = append(passingCritByPractice[d.Practice], d.Criterion)
			}
		}
	}

	// Enrich practice details from criteria_matrix when AI skipped the criteria array.
	for i := range r.Coverage.Details {
		if len(r.Coverage.Details[i].Criteria) == 0 {
			if entry, ok := r.CriteriaMatrix[r.Coverage.Details[i].Name]; ok {
				for _, n := range entry.Pass {
					r.Coverage.Details[i].Criteria = append(r.Coverage.Details[i].Criteria,
						CriterionDetail{Name: n, Status: "pass"})
				}
				for _, n := range entry.Fail {
					r.Coverage.Details[i].Criteria = append(r.Coverage.Details[i].Criteria,
						CriterionDetail{Name: n, Status: "fail"})
				}
			}
		}
	}

	pracRows := buildPracRows(r.Coverage.Details, failingCritByPractice, passingCritByPractice, pracHasFailingDiag)

	covPct := math.Min(r.Coverage.OverallPct, 100)
	pracPassing := r.Coverage.Practices.Passing
	pass := r.Threshold.Status == "pass"
	badgeClass, badgeText := "pass", "PASS"
	if !pass {
		badgeClass, badgeText = "fail", "FAIL"
	}

	data := reportTemplateData{
		Scope:       r.Scope,
		Mode:        r.Mode,
		TimeStr:     timeStr,
		BadgeClass:  badgeClass,
		BadgeText:   badgeText,
		CovPct:      covPct,
		Required:    r.Threshold.Required,
		PracPassing: pracPassing,
		PracTotal:   len(pracRows),
		ErrCount:    len(errors),
		WarnCount:   len(warnings),
		DiagTotal:   len(errors) + len(warnings) + len(infos),
		Practices:   pracRows,
		Errors:      buildDiagRows(errors, projectDir),
		Warnings:    buildDiagRows(warnings, projectDir),
		Infos:       buildDiagRows(infos, projectDir),
		CLIVersion:  version,
		LiveMode:    liveMode,
	}

	var buf bytes.Buffer
	if err := reportTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering HTML report: %w", err)
	}
	return buf.Bytes(), nil
}

// WriteHTMLReport generates and writes the static HTML report (no live badge).
// Returns the HTML file path.
func WriteHTMLReport(jsonPath, version, projectDir string) (string, error) {
	htmlBytes, err := RenderHTMLReport(jsonPath, version, projectDir, false)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(jsonPath)
	htmlPath := filepath.Join(dir, "compliance-latest.html")
	if err := os.WriteFile(htmlPath, htmlBytes, 0o644); err != nil {
		return "", fmt.Errorf("writing HTML report: %w", err)
	}
	return htmlPath, nil
}

// RenderLoadingHTML renders a minimal live-mode page shown before the first check completes.
func RenderLoadingHTML() ([]byte, error) {
	var buf bytes.Buffer
	if err := reportTmpl.Execute(&buf, reportTemplateData{LiveMode: true, BadgeClass: "pass", BadgeText: "..."}); err != nil {
		return nil, fmt.Errorf("rendering loading HTML: %w", err)
	}
	return buf.Bytes(), nil
}

func buildPracRows(details []PracticeDetail, failingCritByPractice map[string]map[string]int, passingCritByPractice map[string][]string, pracHasFailingDiag map[string]bool) []pracRow {
	rows := make([]pracRow, len(details))
	for i, d := range details {
		iconClass, iconChar, barClass := "ok", "✓", "ok"
		if d.Failing > 0 {
			iconClass, iconChar, barClass = "fail", "✗", "fail"
		} else if d.Partial > 0 {
			iconClass, iconChar, barClass = "partial", "~", "partial"
		}
		var failingCriteria, passingCriteria []criterionRow
		if len(d.Criteria) > 0 {
			// criteria_matrix canonical names — authoritative source for both statuses.
			for _, c := range d.Criteria {
				if c.Status == "fail" {
					failingCriteria = append(failingCriteria, criterionRow{Name: c.Name})
				} else {
					passingCriteria = append(passingCriteria, criterionRow{Name: c.Name})
				}
			}
			sort.Slice(failingCriteria, func(a, b int) bool {
				return failingCriteria[a].Name < failingCriteria[b].Name
			})
		} else {
			// Fallback: per-file diagnostic names (may have wording variants for same criterion).
			if critMap := failingCritByPractice[d.Name]; len(critMap) > 0 {
				failingCriteria = make([]criterionRow, 0, len(critMap))
				for name, count := range critMap {
					failingCriteria = append(failingCriteria, criterionRow{Name: name, Count: count})
				}
				// Sort by occurrences desc so the most-representative name comes first.
				sort.Slice(failingCriteria, func(a, b int) bool {
					if failingCriteria[a].Count != failingCriteria[b].Count {
						return failingCriteria[a].Count > failingCriteria[b].Count
					}
					return failingCriteria[a].Name < failingCriteria[b].Name
				})
				// AI summary d.Failing is ground truth; wording variants inflate the list.
				if d.Failing > 0 && len(failingCriteria) > d.Failing {
					failingCriteria = failingCriteria[:d.Failing]
				}
			}
			for _, name := range passingCritByPractice[d.Name] {
				passingCriteria = append(passingCriteria, criterionRow{Name: name})
			}
		}
		// Sync display counts to named criteria when canonical data is available.
		passing, total := d.Passing, d.Total
		if len(d.Criteria) > 0 {
			passing = len(passingCriteria)
			total = len(passingCriteria) + len(failingCriteria)
		}
		// Recompute icon from actual criteria lists (canonical or fallback).
		// d.Failing only counts severity-1; fallback criteria include severity-2 warnings.
		barPct := math.Min(d.CoveragePct, 100)
		if len(failingCriteria) > 0 || len(passingCriteria) > 0 {
			if len(failingCriteria) > 0 {
				iconClass, iconChar, barClass = "fail", "✗", "fail"
			} else {
				iconClass, iconChar, barClass = "ok", "✓", "ok"
			}
			// Bar recalculated only for PATH A (canonical) where passing/total are precise.
			if len(d.Criteria) > 0 && total > 0 {
				barPct = math.Min(float64(passing)/float64(total)*100, 100)
			}
		}
		// Safety net: trust d.Failing and raw diagnostics over criteria lists.
		// Fires when all criteria-based logic shows "ok" but evidence says otherwise.
		if iconClass == "ok" && (d.Failing > 0 || pracHasFailingDiag[d.Name]) {
			iconClass, iconChar, barClass = "fail", "✗", "fail"
			if d.Total > 0 && d.Failing > 0 {
				barPct = math.Min(float64(d.Total-d.Failing)/float64(d.Total)*100, 100)
			}
		}
		rows[i] = pracRow{
			IconClass:       iconClass,
			IconChar:        iconChar,
			Name:            d.Name,
			Passing:         passing,
			Total:           total,
			BarClass:        barClass,
			BarPct:          barPct,
			FailingCriteria: failingCriteria,
			PassingCriteria: passingCriteria,
		}
	}
	return rows
}

var severityLabels = map[int]string{1: "Error", 2: "Warning", 3: "Info", 4: "Hint"}

const snippetCtx = 3

func extractCodeSnippet(d *Diagnostic, projectDir string) []codeLine {
	if d.URI == "" || projectDir == "" {
		return nil
	}
	data, err := os.ReadFile(filepath.Join(projectDir, d.URI))
	if err != nil {
		return nil
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	hlStart := d.Range.Start.Line
	hlEnd := d.Range.End.Line
	if hlEnd < hlStart {
		hlEnd = hlStart
	}
	if hlStart < 0 || hlStart >= len(lines) {
		return nil
	}
	start := hlStart - snippetCtx
	if start < 0 {
		start = 0
	}
	end := hlEnd + snippetCtx
	if end >= len(lines) {
		end = len(lines) - 1
	}
	startChar := d.Range.Start.Character
	endChar := d.Range.End.Character
	// Inline char split only for same-line ranges with a genuine non-zero char span.
	// Zero-length points (startChar==endChar) keep plain whole-line highlight.
	hasCharHL := hlStart == hlEnd && startChar != endChar

	result := make([]codeLine, 0, end-start+1)
	for i := start; i <= end; i++ {
		cl := codeLine{Num: i + 1, Text: lines[i], Highlighted: i >= hlStart && i <= hlEnd}
		if hasCharHL && i == hlStart {
			runes := []rune(lines[i])
			sc, ec := startChar, endChar
			if sc > len(runes) {
				sc = len(runes)
			}
			if ec > len(runes) {
				ec = len(runes)
			}
			if ec < sc {
				ec = sc
			}
			cl.HasCharHL = true
			cl.PreHl = string(runes[:sc])
			cl.InHl = string(runes[sc:ec])
			cl.PostHl = string(runes[ec:])
		}
		result = append(result, cl)
	}
	return result
}

func buildDiagRows(diags []Diagnostic, projectDir string) []diagRow {
	rows := make([]diagRow, len(diags))
	for i := range diags {
		d := diags[i]
		loc := ""
		if d.URI != "" {
			sl, sc := d.Range.Start.Line+1, d.Range.Start.Character
			el, ec := d.Range.End.Line+1, d.Range.End.Character
			switch {
			case sl == el && sc == 0 && ec == 0:
				loc = fmt.Sprintf("%s:%d", d.URI, sl)
			case sl == el:
				loc = fmt.Sprintf("%s:%d:%d-%d", d.URI, sl, sc, ec)
			default:
				loc = fmt.Sprintf("%s:%d-%d", d.URI, sl, el)
			}
		}
		codeLines := extractCodeSnippet(&d, projectDir)
		if d.Criterion != "" {
			for j := len(codeLines) - 1; j >= 0; j-- {
				if codeLines[j].Highlighted {
					codeLines[j].InlineCriterion = d.Criterion
					break
				}
			}
		}
		rows[i] = diagRow{
			Message:       d.Message,
			Loc:           loc,
			Practice:      d.Practice,
			Criterion:     d.Criterion,
			SeverityLabel: severityLabels[d.Severity],
			CodeLines:     codeLines,
		}
	}
	return rows
}

func filterDiagsBySeverity(diags []Diagnostic, severity int) []Diagnostic {
	var out []Diagnostic
	for i := range diags {
		if diags[i].Severity == severity {
			out = append(out, diags[i])
		}
	}
	return out
}
