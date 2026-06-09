# Best Practice Driven Development

Systematically align any project or artifact to its stated best practice preferences — using a Red-Green-Refactor cycle driven by `settings.toml`.

---

## What It Is

BPDD inverts the usual workflow: declare what "good" looks like first, then bring the artifact into alignment. The spec is `settings.toml` plus any active profiles — no separate test files, no checklist to maintain. When the compliance check passes, the artifact matches the spec.

It works for any work product: a codebase, a legal contract, a business plan, a training program, a marketing campaign. The cycle is identical regardless of domain.

---

## The Cycle

```
1. Red      — run compliance check; identify which practices FAIL or are PARTIAL
2. Green    — invoke the relevant grimoire skill; fix until the check passes
3. Refactor — clean up the implementation while keeping the check green
4. Commit   — record progress; return to step 1 for the next gap
```

Invoke with `/apply-best-practice-driven-development`. The skill resolves the effective spec from all settings layers, shows what it's aligning to, then drives the cycle until all practices pass.

---

## Grimoire as a Linter

`check-best-practice-compliance` works like ESLint or Rubocop — but for best practices across any domain.

Encode quality criteria once in `settings.toml`. Run against any artifact, any number of times. Same criteria every run: no inconsistency from who reviewed or when. Gaps that would survive human review get caught by the check.

```toml
# .grimoire/settings.toml
profiles = ["oop"]

[engineering.architecture]
practices = ["apply-solid-principles", "apply-law-of-demeter"]
compliance-threshold = 80
```

Invoke with `/check-best-practice-compliance`. Scope options: specific artifact, region, changed only (for CI), or full project.

---

## Running a Check

Console output:

```
Compliance report — src/contracts/VendorAgreement.md

  ✓ apply-kiss-principle          100%  (3/3 criteria)
  ✗ audit-gdpr-compliance          40%  (2/5 criteria)
    ERROR  §42–58: No DPA section — GDPR Art.28 requires one
    ERROR  §12:    Liability cap missing — required for processor agreements
  ~ apply-solid-principles         60%  (3/5 criteria, 1 suppressed)
  ✗ apply-domain-driven-design      0%  (0/4 criteria)

Coverage: 61.1%  ·  Threshold: 80%  ·  Result: FAIL
Reports:  .grimoire/reports/compliance-2026-06-09T14-32.json
          .grimoire/reports/compliance-2026-06-09T14-32.html
```

**Output files** — always written alongside the console summary:

| File | Format | Use |
|------|--------|-----|
| `compliance-<timestamp>.json` | LSP-compatible JSON | editors, CI, dashboards, other skills |
| `compliance-latest.json` | same, symlinked | always points to most recent run |
| `compliance-<timestamp>.html` | HTML | browser or CI artifact upload |

**Coverage** — `overall_pct` = criteria passing / criteria total. Set a threshold to gate CI:

```toml
[engineering]
compliance-threshold = 80        # fail if overall coverage drops below 80%
compliance-threshold-error = 0   # fail if any error-severity violations remain
```

---

## LSP-Compatible Output

The JSON report follows the LSP Diagnostic schema — consumable by any editor, CI pipeline, LSP server, or downstream grimoire skill without transformation.

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
  "diagnostics": [
    {
      "uri": "file:///project/src/contracts/VendorAgreement.md",
      "range": { "start": { "line": 42, "character": 0 }, "end": { "line": 58, "character": 0 } },
      "severity": 1,
      "code": "audit-gdpr-compliance/dpa-required",
      "source": "grimoire",
      "message": "No data processing agreement section found — GDPR Art.28 requires one",
      "practice": "audit-gdpr-compliance",
      "criterion": "dpa-required",
      "status": "fail"
    }
  ]
}
```

**Severity:** `1 = Error (FAIL)` · `2 = Warning (PARTIAL)` · `3 = Information` · `4 = Hint`

`uri` + `range` locate any text artifact — a legal clause is a text file; `line`/`character` still points to the problematic section. A future grimoire LSP server can read `compliance-latest.json` and publish diagnostics directly to any editor.

**CI integration** — check `"result"` or `"threshold.status"` in the JSON to gate a pipeline:

```bash
result=$(jq -r '.threshold.status' .grimoire/reports/compliance-latest.json)
[ "$result" = "pass" ] || exit 1
```

---

## Suppressing False Positives

Suppress a single finding inline:

```
# grimoire-ignore: apply-solid-principles/srp
class LegacyAdapter:
    ...  # intentional god class — refactor blocked by vendor contract
```

Suppress a block:

```
# grimoire-ignore-start: apply-low-coupling
... third-party integration block ...
# grimoire-ignore-end
```

Suppressed findings are never silently dropped. They appear as `"status": "suppressed"` in the JSON output and in a separate section of the HTML report. Coverage counts suppressed findings separately — suppressed ≠ passing.

---

## Incremental Mode

Scope `[c] Changed only` checks only modified lines or sections — faster for CI and pre-commit hooks.

```json
{ "mode": "incremental", "base_ref": "HEAD~1", ... }
```

`base_ref` is included in the JSON so downstream tools know exactly what was checked.

---

## Settings Reference

| Key | Where | Description |
|-----|-------|-------------|
| `profiles` | top-level | Activate skill bundles by name or tag |
| `practices` | `[domain]` section | Explicit ordered skill list for a domain |
| `compliance-threshold` | `[domain]` section | Minimum overall coverage % to pass |
| `compliance-threshold-error` | `[domain]` section | Max allowed error-severity violations |

See [profiles.md](./profiles.md) for how `profiles` resolves and how to author custom profile files.

---

## When to Use Which Skill

| Task | Skill |
|------|-------|
| One-off check against current spec | `check-best-practice-compliance` |
| Fix one specific compliance finding (targeted, location-aware) | `fix-best-practice-finding` |
| Full systematic alignment until all practices pass | `apply-best-practice-driven-development` |
| Discover which practices apply to an artifact | `review-best-practice-fit` |
| Change or add a preference | `pin-best-practice-preference` |
| Activate a paradigm's practices in bulk | `apply-best-practice-profile` |
