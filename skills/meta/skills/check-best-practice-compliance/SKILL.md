---
name: check-best-practice-compliance
description: Use when the user wants to check whether any artifact or work product aligns with their stated best practice preferences — e.g., "check compliance", "linter for best practices", "are we following our pinned practices?", "check this document against our standards".
source: 'Sadowski et al. (2015) "Tricorder: Building a Program Analysis Ecosystem Inside Google", ICSE; Ford, Parsons, Kua (2017) "Building Evolutionary Architectures"; Open Policy Agent documentation'
tags: [compliance, linter, enforcement, audit, check, preferences, settings, domain-agnostic]
related: [apply-best-practice-driven-development, review-best-practice-fit, pin-best-practice-preference, apply-best-practice-profile]
---

# Check Best Practice Compliance

Run a compliance check against any artifact using the resolved effective preferences from settings.toml. Outputs LSP-compatible JSON (machine-readable, editor-consumable) and an HTML coverage report (human-readable).

## Why This Is Best Practice

**Adopted by:** Linters (ESLint, Checkstyle, Rubocop), architecture validation tools (ArchUnit), and policy-as-code frameworks (OPA, Conftest) encode quality criteria once and run them repeatedly against any artifact. Google's Tricorder runs 30+ analyzers on every CL and blocks 15–20% that would have passed human review; false-positive rate kept below 5% by restricting to well-defined criteria (Sadowski et al., ICSE 2015).
**Impact:** Manual best-practice review is inconsistent — what gets checked depends on who reviews and when. A compliance check applies the same criteria every run and produces a structured, comparable result. Coverage metrics make progress visible: teams that track compliance coverage improve it (Hawthorne effect for quality metrics, confirmed in code review studies at Microsoft Research).
**Why best:** `review-best-practice-fit` is one-shot and subjective. `check-best-practice-compliance` is repeatable, explicit, and machine-readable — it runs against the user's actual stated preferences (resolved from settings), producing LSP-compatible diagnostics that any editor or CI system can consume directly.

Sources: Sadowski et al. (2015) "Tricorder: Building a Program Analysis Ecosystem Inside Google"; Ford, Parsons, Kua (2017) "Building Evolutionary Architectures"; ESLint architecture documentation

## Steps

### 1. Resolve the effective spec (silent)

Apply all settings layers in precedence order (session > project-local > project-shared > global). Expand profiles, apply domain overrides, remove disabled entries. The resolved spec is the check suite — not any individual layer.

---

### 2. Select scope

Default: `[a] Full project`. Infer from context if user specified a file, directory, or "changed only" — no need to ask. Only ask if scope is genuinely ambiguous.

| Keyword in request | Inferred scope |
|--------------------|---------------|
| names a file/component | `[s]` specific artifact |
| names a directory/module | `[r]` region |
| "changed", "diff", "PR", "commit" | `[c]` changed only |
| no specific target | `[a]` full project (default) |

---

### 3. Evaluate each practice

For each practice in the resolved spec, check the target artifact against its observable criteria. Each criterion: pass, fail, partial, or suppressed (via inline annotation).

**False positive suppression** — any finding can be suppressed inline. Suppressed findings are never dropped from output — they appear as `"status": "suppressed"` in JSON.

```
# grimoire-ignore: apply-solid-principles/srp
class LegacyAdapter:  ...  # intentional god class — refactor blocked by contract

# grimoire-ignore-start: apply-low-coupling
... third-party integration block ...
# grimoire-ignore-end
```

**Configuring suppressions:** Findings can be suppressed two ways:
1. Inline annotation: add `# grimoire-ignore: [practice]/[criterion]` on the line before the violation
2. Settings entry: add the practice to the `disabled` array in `.grimoire/settings.toml` for the relevant domain

Suppressed findings always appear in JSON output with `"status": "suppressed"` — they are never silently dropped.

---

### 4. JSON output (primary)

Always written to `.grimoire/reports/compliance-<timestamp>.json` (e.g., `.grimoire/reports/compliance-2026-06-14T143000.json`). Also write a symlink `.grimoire/reports/compliance-latest.json` pointing to the most recent report.

LSP-compatible schema — consumable by editors, CI pipelines, LSP servers, dashboards:

```json
{
  "version": "1",
  "timestamp": "2026-06-09T14:32:00Z",
  "mode": "full",
  "scope": "src/contracts/VendorAgreement.md",
  "spec": {
    "sources": [".grimoire/settings.toml", "~/.config/grimoire/settings.toml"],
    "resolved_from": "project-shared + global"
  },
  "result": "fail",
  "coverage": {
    "overall_pct": 61.1,
    "practices": { "total": 4, "passing": 1, "partial": 1, "failing": 2, "coverage_pct": 37.5 },
    "criteria":  { "total": 18, "passing": 11, "failing": 7, "suppressed": 1, "coverage_pct": 61.1 }
  },
  "threshold": { "required": 80, "actual": 61.1, "status": "fail" },
  "summary": { "error": 2, "warning": 1, "info": 0, "pass": 11, "suppressed": 1 },
  "diagnostics": [
    {
      "uri": "file:///project/src/contracts/VendorAgreement.md",
      "range": { "start": { "line": 42, "character": 0 }, "end": { "line": 58, "character": 0 } },
      "severity": 1,
      "code": "audit-gdpr-compliance/dpa-required",
      "source": "grimoire",
      "message": "No data processing agreement section found — GDPR Art.28 requires one before processing EU customer data",
      "practice": "audit-gdpr-compliance",
      "criterion": "dpa-required",
      "status": "fail",
      "coverage": { "total": 5, "passing": 2, "coverage_pct": 40.0 }
    },
    {
      "uri": "file:///project/src/UserService.ts",
      "range": { "start": { "line": 12, "character": 0 }, "end": { "line": 45, "character": 1 } },
      "severity": 2,
      "code": "apply-solid-principles/srp",
      "source": "grimoire",
      "message": "UserService handles auth, email, and billing (3 concerns) — violates SRP",
      "practice": "apply-solid-principles",
      "criterion": "srp",
      "status": "suppressed",
      "suppressed_by": "inline",
      "coverage": { "total": 5, "passing": 3, "coverage_pct": 60.0 }
    }
  ]
}
```

**LSP severity:** 1 = Error (FAIL) · 2 = Warning (PARTIAL) · 3 = Information · 4 = Hint

**Incremental mode** (`[c] changed only`): output includes `"mode": "incremental"` and `"base_ref": "HEAD~1"` so downstream tools know what was checked.

---

### 5. HTML report (human summary)

Always generated alongside JSON at `.grimoire/reports/compliance-<timestamp>.html`.

Structure:
- Header: artifact, timestamp, spec sources, overall result badge (PASS/FAIL)
- Coverage gauge: overall % with threshold indicator
- Per-practice sections: collapsible, progress bar per practice, criteria checklist
- Suppressed findings: shown separately, not counted in coverage
- Footer: links to settings files, next actions

Color coding: green (pass) · red (fail) · amber (partial) · grey (suppressed).
Suitable for browser viewing or CI artifact upload.

---

### 6. Console summary

```
Compliance report — src/contracts/VendorAgreement.md

  ✓ apply-kiss-principle       100%  (3/3 criteria)
  ✗ audit-gdpr-compliance       40%  (2/5 criteria)
    ERROR  §42–58: No DPA section — GDPR Art.28 requires one
    ERROR  §12:    Liability cap missing — required for processor agreements
  ~ apply-solid-principles      60%  (3/5 criteria, 1 suppressed)
  ✗ apply-domain-driven-design   0%  (0/4 criteria)

Coverage: 61.1%  ·  Threshold: 80%  ·  Result: FAIL
Reports:  .grimoire/reports/compliance-2026-06-09T14-32.json
          .grimoire/reports/compliance-2026-06-09T14-32.html
```

---

### 7. Offer fixes

For each FAIL, offer to invoke the relevant grimoire skill, or run the full BPDD cycle:

```
Fix with:
  [1] audit-gdpr-compliance    — add DPA section and liability clause
  [2] apply-domain-driven-design — extract domain model
  [a] Run full BPDD cycle      → /apply-best-practice-driven-development
```

**Routing:** Auto-invoke only the top-priority FAIL (highest severity, lowest coverage %). For all other FAILs: list them in the output but do not invoke — wait for user to select. User can type a number to fix that specific finding, or [a] to run the full BPDD cycle.

## Common Mistakes

**Checking against a single settings layer instead of the resolved spec.** Always resolve the full precedence chain — project overrides global.

**Treating suppressed findings as passing.** Suppressed findings are excluded from coverage calculation but counted separately. Review suppressions periodically — suppressed ≠ acceptable.

**Setting threshold to 100% on legacy projects.** Start with the current coverage as baseline, improve incrementally. A threshold below current = no gate; set it at current + 10% to enforce improvement without blocking work.

## When NOT to Use

- **For one-time exploration** — use `review-best-practice-fit` to discover what practices apply; use this when you know your spec and want to verify alignment
- **When settings.toml is empty** — set preferences first; this checks against your stated spec, not general best practices
