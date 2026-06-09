---
name: fix-best-practice-finding
description: Use when the user wants to fix a specific compliance finding from a check-best-practice-compliance report — e.g., "fix this finding", "resolve violation #3", "fix the SRP violation in UserService", "close this linter error".
source: LSP code action protocol (Microsoft, 2016); ESLint --fix design; VS Code Quick Fix API
tags: [compliance, linter, fix, finding, bpdd, enforcement, meta, grimoire]
related: [check-best-practice-compliance, apply-best-practice-driven-development, review-best-practice-fit]
---

# Fix Best Practice Finding

Fix one specific compliance finding from a `check-best-practice-compliance` report — targeted, location-aware, verified.

## Why This Is Best Practice

**Adopted by:** LSP Code Action protocol (used by VS Code, JetBrains, Neovim) defines the same pattern: a diagnostic identifies the problem at a location; a code action fixes it at that location. ESLint `--fix` applies the same model — run the linter, get findings, fix individual rules. The separation between detection (compliance check) and remediation (fix action) is foundational to every production linter toolchain.
**Impact:** Without targeted fix, users must either run the full BPDD cycle (overkill for one finding) or invoke a domain skill with no location context (the skill applies the practice generally, not to the specific violation). Targeted fix closes the gap: the compliance JSON provides the exact `uri`, `range`, `criterion`, and `message` — passing that context to the domain skill produces a precise, verifiable fix.
**Why best:** `apply-best-practice-driven-development` is systematic but heavy — it processes all practices in priority order. This skill is surgical: one finding, one fix, one verification. Same relationship as `eslint --fix --rule srp` vs `eslint --fix`.

Sources: LSP specification §3.16 (Code Action); ESLint architecture documentation; VS Code Language Server Extension Guide

## Steps

### 1. Identify the finding

Accept input in any form:

- JSON diagnostic object pasted from `compliance-latest.json`
- Finding code (e.g., `apply-solid-principles/srp`)
- Natural language description ("fix the SRP violation in UserService.ts around line 12")
- Range reference from editor (file + line/character)

If multiple findings match (e.g., same practice, multiple locations), list them and ask the user to select one:

```
Multiple findings match apply-solid-principles/srp:

  [1] src/UserService.ts §12–45    — handles auth, email, and billing (3 concerns)
  [2] src/OrderService.ts §8–62    — handles payment, inventory, and shipping (3 concerns)

Which finding to fix?
```

---

### 2. Show the finding

Confirm what will be fixed before acting:

```
Finding:  apply-solid-principles/srp
Location: src/UserService.ts §12–45
Severity: Error
Message:  UserService handles auth, email, and billing (3 concerns) — violates SRP
```

---

### 3. Invoke the domain skill with context

Read `"practice"` from the finding — that value is the grimoire skill to invoke. Pass the following as targeted context:

- **Location** — `uri` + `range` (the exact artifact and lines to fix)
- **Criterion** — `criterion` field (the specific rule being violated)
- **Message** — `message` field (the specific problem description)

The domain skill receives this context and focuses only on the identified violation — not a general application of the practice across the whole artifact.

Example: for `"practice": "apply-solid-principles"` with `"criterion": "srp"` at `src/UserService.ts §12–45`, invoke `/apply-solid-principles` with the context: "Fix SRP violation at §12–45: UserService handles auth, email, and billing. Extract concerns into separate services."

---

### 4. Re-run targeted compliance check

Re-run `check-best-practice-compliance` scoped to:
- The affected artifact (`[s] Specific artifact`)
- Filtered to the affected practice

Confirm the specific finding is resolved — no longer appears in diagnostics with matching `uri` + `code`.

If still present: the domain skill partially addressed it. Show what changed, continue fixing.

---

### 5. Report

```
✓ Fixed: apply-solid-principles/srp in src/UserService.ts §12–45
  UserService now delegates auth to AuthService, email to NotificationService

Remaining findings in this artifact: 1
  ✗ apply-solid-principles/dip — §8: direct dependency on MySQLUserRepository

Fix next? [y] apply-solid-principles/dip  [n] done  [a] run full BPDD cycle
```

If all findings in the artifact are resolved, report final coverage and offer to run the full project check.

## Common Mistakes

**Fixing at the wrong scope.** The fix must target the exact `uri` + `range` from the finding — not the whole file, class, or module. Passing location context to the domain skill is required for a precise fix.

**Marking fixed before re-check.** Always re-run the compliance check after the fix. The domain skill may have addressed part of the violation while introducing a different one. The check is the source of truth — not the AI's judgment.

**Fixing suppressed findings.** If a finding has `"status": "suppressed"`, it was intentionally ignored. Do not fix suppressed findings unless the user explicitly asks to remove the suppression.

## When NOT to Use

- **To fix all findings at once** — use `apply-best-practice-driven-development` (full BPDD cycle with priority ordering)
- **To explore what practices apply to an artifact** — use `review-best-practice-fit` first
- **When no compliance report exists** — run `check-best-practice-compliance` first to generate findings
