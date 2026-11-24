# Best Practice Driven Development

Systematically align any project or artifact to its stated best practice preferences — using a Red-Green-Refactor cycle driven by `grimoire.toml`.

← [Back to the README](../README.md)

---

## What It Is

BPDD inverts the usual workflow: declare what "good" looks like first, then bring the artifact into alignment. The spec is `grimoire.toml` plus any active profiles — no separate test files, no checklist to maintain. When the compliance check passes, the artifact matches the spec.

It works for any work product: a codebase, a legal contract, a business plan, a training program, a marketing campaign. The cycle is identical regardless of domain.

**The body-as-rubric insight:** Every skill file is both a guide and an evaluator. The same SKILL.md body that tells an AI how to apply a practice is what grimoire reads when checking compliance — no separate criteria schema, no drift between "what we say is good" and "what we check for." Update a skill and both the guidance and the evaluation standard update in one step.

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

Encode quality criteria once in `grimoire.toml`. Run against any artifact, any number of times. Same criteria every run: no inconsistency from who reviewed or when. Gaps that would survive human review get caught by the check.

```toml
# grimoire.toml
[standards]
profiles = ["oop"]

[standards.engineering.architecture]
practices = ["apply-solid-principles", "apply-law-of-demeter"]
compliance-threshold = 80
```

Invoke with `/check-best-practice-compliance`. Scope options: specific artifact, region, changed only (for CI), or full project.

---

## Two Modes

`grimoire check` has two modes — choose based on how you want to run the AI:

**Independent mode (default)** — grimoire selects an installed local AI agent (claude, gemini, codex, copilot) or a configured API provider and runs the check itself. No prior AI session needed.

```bash
grimoire check               # auto-selects local CLI or API provider
grimoire check --via claude  # force a specific local agent
```

**Report mode** — an AI session skill (`/check-best-practice-compliance`) generates the JSON report; the CLI reads it and enforces thresholds. Useful when you want the AI to produce a rich report with precise file locations and then gate CI on the result.

```bash
# In AI session:
/check-best-practice-compliance   → writes .grimoire/reports/compliance-latest.json

# In terminal / CI:
grimoire check --from-report      # reads JSON, enforces thresholds, exits 0/1
```

Both modes respect the same `compliance-threshold` and `compliance-threshold-error` settings and support `--ci` / `--junit` output flags.

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
[standards.engineering]
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
    "sources": ["grimoire.toml", "~/.config/grimoire/grimoire.toml"],
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

> **Note:** Suppression markers are evaluated at **report generation time** — when `/check-best-practice-compliance` runs in an AI session (report mode). The CLI (`grimoire check --from-report`) only reads and enforces the already-generated report; it does not re-scan source files for ignore markers.

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
| `profiles` | `[standards]` | Activate skill bundles by name or tag |
| `practices` | `[standards.domain]` section | Explicit ordered skill list for a domain |
| `compliance-threshold` | `[standards.domain]` section | Minimum overall coverage % to pass |
| `compliance-threshold-error` | `[standards.domain]` section | Max allowed error-severity violations |

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

## Grimoire vs Rule Engines

ESLint, Rubocop, Checkstyle: deterministic. Same code → same finding. They enforce codifiable rules — syntax, style, known anti-patterns.

Grimoire: AI judgment. Covers non-codifiable criteria — architecture intent, legal risk, business logic coherence. Two runs may differ, the same way two senior engineers flag different issues in review. This is a feature, not a flaw.

Important findings are stable. Suppress noise with `# grimoire-ignore`. Use thresholds to gate on coverage %, not finding-level exactness.

**Use both.** Linters for syntax rules, Grimoire for judgment.

## CI Integration

Use the `grimoire` CLI to gate PRs on compliance.

**One-time setup:**

```bash
# Install (one step, no runtime dependency)
curl -fsSL https://github.com/jeffreytse/grimoire/releases/latest/download/grimoire-linux-amd64 \
  -o /usr/local/bin/grimoire && chmod +x /usr/local/bin/grimoire

# Initialize project (creates grimoire.toml — commit this)
grimoire init
```

**Independent mode** — grimoire runs the AI check itself in CI (requires a local CLI or API key in the environment):

```bash
grimoire check                                         # auto-selects AI; exits 0/1
grimoire check --ci                                    # + GitHub Actions annotations
grimoire check --junit .grimoire/reports/junit.xml     # + JUnit XML for test reporters
```

`--ci` auto-enables when `GITHUB_ACTIONS=true` — no flag needed in a GHA workflow.

**Report mode** — AI session generates the report, CLI enforces thresholds:

```bash
# After /check-best-practice-compliance in an AI session:
grimoire check --from-report --ci --junit .grimoire/reports/junit.xml
```

**GitHub Actions example:**

```yaml
# .github/workflows/compliance.yml
- name: Run grimoire compliance check
  run: grimoire check --from-report --ci --junit .grimoire/reports/junit.xml
  env:
    ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}   # only needed for independent mode

- name: Upload JUnit report
  uses: actions/upload-artifact@v4
  with:
    name: junit-compliance
    path: .grimoire/reports/junit.xml
```

**`--project-dir`** — pin the project root when CI runs from a different working directory:

```bash
grimoire --project-dir /workspace/myproject check
# or:
GRIMOIRE_PROJECT_DIR=/workspace/myproject grimoire check
```

**Local dev companion** — `grimoire watch` re-runs the check on every file save, giving the same continuous feedback loop as a linter:

```bash
grimoire watch            # re-checks on save using active profile
grimoire watch --via gemini
```

**MCP alternative** — if grimoire is connected via MCP, the AI can call `grimoire_check` directly instead of running the CLI. `grimoire_context` at session start gives the full environment snapshot including current compliance status.


## Linting for practices

Declare which practices your project follows. `grimoire check` enforces them — CI gate, watch mode, or on-demand.

```bash
grimoire check                              # run compliance check against declared practices
grimoire check --live                       # watch mode — re-check on every file save
grimoire check --live --port 8080           # custom port (default 7890)
grimoire check --live --host 127.0.0.1      # localhost-only (default: all interfaces)
grimoire check --junit report.xml           # JUnit XML output for CI systems
grimoire check --scope changed              # check only changed files (incremental)
```

Sample output:

```
  ✓  propose-conventional-commit   100%  all 4 criteria passing
  ✓  apply-solid-principles         88%  7/8 criteria passing
  ✗  review-pull-request            62%  below threshold (80%)
     └─ missing: security-checklist, breaking-change-annotation

2 passed · 1 failed · exit 1
```

Gate any PR: `grimoire check --junit report.xml` exits non-zero when coverage drops below the configured threshold. No other tool checks whether your project _follows_ its declared practices — only whether the config file is valid.

## Quick reference: commands, output and thresholds

grimoire is a linter for best practices — same model as ESLint for code style:

```
# Code style linter:             │  # Best-practice linter:
npm install eslint               │  grimoire mcp install --target all
eslint .                         │  grimoire check
# watch mode:                    │  # watch mode:
eslint --watch                   │  grimoire check --live
```

Declare which practices you require in `grimoire.toml`, install the packages that encode them, then run `grimoire check` against any artifact — a codebase, a legal contract, a business plan, a training program. Same criteria every run. Gaps that survive human review get caught by the check.

**The cycle** — same inversion as TDD: declare what "good" looks like first, then bring the artifact into alignment.

```
1. Red      — run compliance check; identify which practices FAIL or are PARTIAL
2. Green    — invoke the relevant grimoire skill; fix until the check passes
3. Refactor — clean up while keeping the check green
4. Commit   — record progress; repeat for next gap
```

Run `/apply-best-practice-driven-development` to drive the full cycle. Run `/check-best-practice-compliance` for a one-off check.

**Two modes** — choose based on how you want to run the AI:

- **Independent mode** (default) — grimoire runs the AI check itself; no prior AI session needed
- **Report mode** — an AI session skill (`/check-best-practice-compliance`) generates the JSON report; `grimoire check --from-report` reads it and enforces thresholds

```bash
grimoire init                        # one-time project setup
grimoire check                       # independent mode: auto-selects local CLI or API provider
grimoire check --via claude          # force a specific local agent
grimoire check --live                       # file watcher + browser report — like eslint --watch
grimoire check --live --port 8080           # custom port (default 7890)
grimoire check --live --host 127.0.0.1     # localhost-only (default: all interfaces)
grimoire check --ci                  # + GitHub Actions annotations
grimoire check --junit report.xml    # + JUnit XML for CI reporters
grimoire check --from-report         # report mode: reads .grimoire/reports/compliance-latest.json
```

**Output** — always written to `.grimoire/reports/`:

| File                          | Format              | Use                                            |
| ----------------------------- | ------------------- | ---------------------------------------------- |
| `compliance-<timestamp>.json` | LSP-compatible JSON | editors, CI pipelines, LSP servers, dashboards |
| `compliance-latest.json`      | JSON (symlink)      | Always points to most recent run — use in CI   |
| `compliance-<timestamp>.html` | HTML                | browser or CI artifact upload                  |

The JSON follows the LSP Diagnostic schema — `uri` + `range` locate any finding in any text artifact, not just code.

**Coverage thresholds** — set in `grimoire.toml`, enforced on every check:

```toml
[standards.engineering]
compliance-threshold = 80        # fail if overall criteria coverage < 80%
compliance-threshold-error = 0   # fail if any error-severity violations remain
```

Use `/fix-best-practice-finding` to fix one specific compliance finding — targeted, location-aware, verified. Use `/apply-best-practice-driven-development` to fix everything systematically.

→ [Full BPDD guide](./bpdd.md) — cycle, linter, LSP schema, false positive suppression, incremental mode
