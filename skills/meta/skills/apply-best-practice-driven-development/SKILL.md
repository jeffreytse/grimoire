---
name: apply-best-practice-driven-development
description: Use when the user wants to systematically align any project or artifact to their stated best practice preferences — e.g., "apply BPDD", "align this project to our settings", "close the gap between our practices and reality", "enforce our best practices like TDD".
source: 'Beck (2003) "Test-Driven Development: By Example"; Ford, Parsons, Kua (2017) "Building Evolutionary Architectures"; Nagappan et al. (2008) "Realizing Quality Improvement Through Test Driven Development", IEEE'
tags: [bpdd, tdd, compliance, enforcement, cycle, preferences, settings, domain-agnostic]
related: [check-best-practice-compliance, pin-best-practice-preference, apply-best-practice-profile, review-best-practice-fit]
---

# Apply Best Practice Driven Development

Systematically align any project or artifact to its stated best practice preferences using a Red-Green-Refactor cycle driven by `settings.toml` and profiles. Works for codebases, legal documents, business plans, training programs, marketing campaigns — any work product with declared practices.

## Why This Is Best Practice

**Adopted by:** Architecture fitness functions (Ford, Parsons, Kua — "Building Evolutionary Architectures", O'Reilly 2017) encode quality criteria as executable checks and validate continuously — the same inversion as TDD: define the check before the implementation. Netflix, Thoughtworks, and Amazon use fitness functions to enforce architectural constraints at CI time. The same principle applies beyond code: quality management systems (ISO 9001) and legal compliance frameworks define compliance criteria before auditing against them.
**Impact:** Without BPDD, teams apply best practices opportunistically — gaps accumulate and preferences drift silently from reality. TDD reduces defect density 40–80% (Nagappan et al., IEEE 2008) not because tests catch bugs, but because the write-test-first discipline prevents gaps from forming. BPDD applies the same mechanism to best practices across any domain.
**Why best:** `review-best-practice-fit` identifies gaps reactively, after implementation. BPDD inverts this: the preferences declared in settings are the specification; the artifact follows. The cycle continues until artifact and spec agree — leaving no gaps to discover later.

Sources: Beck (2003) "Test-Driven Development: By Example"; Ford, Parsons, Kua (2017) "Building Evolutionary Architectures"; Nagappan et al. (2008) "Realizing Quality Improvement Through Test Driven Development"

## Steps

### 1. Resolve the effective spec (silent)

Read all settings layers in precedence order and produce the single effective spec:

```
session pins
  ↓ overrides
.grimoire/settings.local.toml   (project-personal)
  ↓ overrides
.grimoire/settings.toml         (project-shared)
  ↓ overrides
~/.config/grimoire/settings.toml  (global)
```

- Expand `profiles = [...]` to skill lists (tag query or file)
- Apply domain overrides (`[engineering.architecture]` overrides `[engineering]`)
- Apply `disabled` entries — remove those skills from the spec
- Drop any practice overridden by a higher-precedence layer

---

### 2. Show the resolved spec

Display the effective spec before running — the user must see what they're aligning to:

```
Effective spec — 5 practices

  apply-solid-principles       [profile: oop → .grimoire/settings.toml]
  apply-domain-driven-design   [profile: oop → .grimoire/settings.toml]
  apply-kiss-principle         [engineering.architecture → global]
  apply-pyramid-principle      [writing → .grimoire/settings.toml]
  audit-gdpr-compliance        [law → global]
  ⊘ apply-law-of-demeter       [disabled by .grimoire/settings.toml]
```

Wait for confirmation. If settings are empty, stop and direct user to `pin-best-practice-preference` or `apply-best-practice-profile` first.

---

### 3. Select scope

```
Scope:
  [s] Specific artifact   — a file, document, section, or component
  [r] Region / scope      — a directory, chapter, module, or area
  [c] Changed only        — diff-based, checks only modified parts
  [a] Full project        — everything in scope of the settings
```

---

### 4. Run compliance check

Call `check-best-practice-compliance` against the selected scope using the resolved spec. Read the JSON output to classify each practice:

- **PASS** — already applied
- **PARTIAL** — partially applied, gaps remain
- **FAIL** — not applied or actively violated

```
Compliance baseline

  ✓ apply-kiss-principle          PASS     (3/3 criteria)
  ✗ apply-solid-principles        FAIL     (2/5 criteria, 40% coverage)
  ~ apply-pyramid-principle       PARTIAL  (2/3 criteria, 67% coverage)
  ✗ apply-domain-driven-design    FAIL     (0/4 criteria, 0% coverage)
  ✓ audit-gdpr-compliance         PASS     (5/5 criteria)

Overall: 55.0% criteria coverage  ·  threshold: 80%  ·  STATUS: FAIL
2 passing · 2 failing · 1 partial
```

If all pass and threshold met: codebase already aligns with spec. Report and stop.

---

### 5. Pick next gap

Priority order (explicit, shown to user):

1. Practices that are causing a threshold violation (blocking — fix these first)
2. FAIL before PARTIAL
3. Within same status: settings array order (index 0 = highest priority)
4. User override: "fix X next" accepted at any point in the cycle

```
Next: apply-solid-principles (FAIL, 40% coverage)
  Reason: causing threshold failure + listed before apply-domain-driven-design

  ✗ SRP: UserService handles auth, email, and billing (3 concerns)
  ✗ DIP: direct dependency on MySQLUserRepository
  ✓ OCP, LSP, ISP — already passing
```

---

### 6. Green — implement

Read the `"practice"` field from each failing diagnostic in the JSON report — that value is the grimoire skill to invoke. `"practice": "apply-solid-principles"` → `/apply-solid-principles`; `"practice": "apply-domain-driven-design"` → `/apply-domain-driven-design`; and so on. Re-run `check-best-practice-compliance` for this practice after each change until it passes.

---

### 7. Refactor

Clean up the implementation while keeping the check green. Re-run to confirm no regression.

---

### 8. Commit

Record progress. Return to step 5 for the next failing practice. Continue until all practices pass and threshold is met.

---

### 9. Final report

```
✓ BPDD complete — 5/5 practices aligned  ·  92% criteria coverage

  ✓ apply-solid-principles      (fixed: 3 violations resolved)
  ✓ apply-domain-driven-design  (fixed: domain model extracted)
  ✓ apply-pyramid-principle     (fixed: 1 remaining gap closed)
  ✓ apply-kiss-principle        (was already passing)
  ✓ audit-gdpr-compliance       (was already passing)

Threshold: 92% ≥ 80% ✓
Artifact now aligns with .grimoire/settings.toml
```

## Common Mistakes

**Running BPDD without reviewing the resolved spec first.** Always show step 2 — the user must confirm what spec they're aligning to. A misconfigured settings file means fixing the wrong things.

**Skipping commit between cycles.** Each Green must be committed before moving to the next practice — same discipline as TDD.

**Ignoring suppressed findings.** Suppressed violations (via `# grimoire-ignore`) still appear in the JSON report. Review them periodically — suppressed ≠ resolved.

## When NOT to Use

- **For one-time review** — use `review-best-practice-fit` instead; BPDD is for systematic, committed alignment
- **When settings.toml is empty** — set preferences first with `pin-best-practice-preference` or `apply-best-practice-profile`
