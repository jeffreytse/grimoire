---
name: audit-applied-best-practices
description: Use when the user wants to inventory which best practices are already applied in their existing work — any artifact across any domain (codebase, document, design, music, finance, law, science, writing, marketing, and more) — and pin intentional choices as preferences that routing skills will respect.
source: ISO 9001 gap audit standards (ISO 9001:2015 §4.1); Nielsen Norman Group contextual inquiry methodology
tags: [audit, inventory, applied-practices, preference, existing-work, context-awareness, onboarding]
---

# Audit Applied Best Practices

Scan existing work to map which best practices are already in use, then let the user pin intentional choices as preferences that all routing skills respect.

## Why This Is Best Practice

**Adopted by:** ISO 9001 mandates formal context-of-organization audits before any quality initiative — organizations must document current practices before assessing gaps (ISO 9001:2015 §4.1). McKinsey's as-is mapping precedes every strategy engagement. Nielsen Norman Group's contextual inquiry methodology requires observing what users *actually do* before designing interventions.
**Impact:** Skipping the as-is audit is the primary cause of process improvement failure — organizations apply new practices that conflict with established workflows, creating resistance rather than improvement (Kotter, "Leading Change", 1996). ISO 9001 organizations that conduct thorough current-state audits have 60% higher implementation success rates than those that skip this step (BSI Group, 2019).
**Why best:** You cannot improve what you haven't mapped. Without an applied-practice inventory, routing skills suggest changes to things the user deliberately chose — destroying trust. Pinning preferences converts noise into signal: the user tells the system "I know, and I chose this intentionally."

Sources: ISO 9001:2015 §4.1; Kotter (1996) "Leading Change"; BSI Group (2019) ISO certification outcomes study; Nielsen Norman Group contextual inquiry guidelines

## Steps

### Step 1: Detect work type (silent)

Auto-classify from all available signals — files present, user-provided content, and verbal description. Check in this order:

**File-based signals (check first):**

| Signal | Detected type |
|--------|--------------|
| `.git` directory; `package.json`, `go.mod`, `Cargo.toml`, `requirements.txt`, `Gemfile`, `pom.xml` | Codebase |
| `.sketch`, `.fig`, `figma.json`, `.xd`, `.ai`, `.psd`, design token files | Design |
| `*.tex`, `*.bib`, structured academic directory (`sections/`, `figures/`, `references/`) | Science / Academic writing |
| `.flp`, `.als`, `.logic`, `*.mid`, `*.wav` collection, DAW project files | Music |
| `*.psd`, `*.ai`, large `*.svg` / `*.png` collection, illustration project files | Art |
| `*.xlsx`, `*.csv` with financial headers (revenue, expenses, EBITDA), financial model files | Finance |
| `*.docx` / `*.pdf` with legal structure (parties, whereas, definitions, governing law) | Law |
| Screenplay file (`.fdx`, `.fountain`) or content matching INT./EXT. scene headers | Film / Screenwriting |
| `*.recipe`, structured cooking files, ingredient/instruction pattern | Cooking |
| Marketing campaign folder (`campaigns/`, `ads/`, UTM-tagged links, social copy files) | Marketing |
| `*.jpg` / `*.raw` / `*.dng` collection, Lightroom catalog, editing preset files | Photography |
| Academic or research structure (abstract, methodology, results sections) | Science / Writing |
| Single document or paste (essay, article, report, proposal) | Document |
| Multiple heterogeneous file types not matching above | Mixed artifact |

**Content-based signals (if no files):**
If user pastes or describes content, infer type from terminology and structure:
- Code syntax, imports, function definitions → Codebase
- Musical notation, tempo, key, DAW terminology → Music
- Legal terms (indemnify, jurisdiction, whereas) → Law
- Financial terms (P&L, EBITDA, amortization) → Finance
- Script formatting (slug lines, action lines, dialogue) → Film
- Scientific methodology (hypothesis, variables, p-value) → Science
- Anything else → Document or Verbal description

**Verbal fallback:**
If no files and no content signals, detect domain from what the user describes. Ask ONE targeted question: "What kind of work are we auditing? (e.g. codebase, business strategy doc, music project, legal contract)"

For Codebase, Design, Finance, and Music: proceed automatically.
For all others: ask ONE targeted question to confirm scope before proceeding.

### Step 2: Extract signals

Gather evidence of existing practices based on detected type:

| Work type | Where to look |
|-----------|--------------|
| **Codebase** | `package.json` / `Gemfile` / `go.mod` / `requirements.txt` (dependencies → infer testing framework, linter, formatter); CI config (`.github/workflows/`, `.gitlab-ci.yml`); test file structure (co-located vs. `tests/`); `git log --oneline -50` (commit message style); linter configs (`.eslintrc`, `.rubocop.yml`, `pyproject.toml`); folder structure and naming conventions |
| **Design** | Component library presence, design token structure, spacing/color system, accessibility annotations, handoff conventions |
| **Science / Academic** | Citation style (APA/MLA/Chicago), abstract structure, statistical reporting conventions, data organization, reproducibility setup (notebooks, scripts) |
| **Music** | Song structure conventions, mix bus chain, reference track usage, naming/versioning of project files, stem export format |
| **Art** | File naming and versioning, layer organization, color profile, export format and resolution standards |
| **Finance** | Model structure (assumptions tab, outputs tab), formula transparency vs. hardcoded values, version control convention, audit trail |
| **Law** | Clause ordering convention, defined terms usage, governing law, signature block format, amendment process |
| **Film / Screenwriting** | Format compliance (margins, slug lines), scene numbering, revision mark convention, breakdown document structure |
| **Marketing** | UTM naming convention, campaign naming taxonomy, asset versioning, brief format, approval workflow |
| **Photography** | Import/folder naming convention, culling workflow, editing preset usage, export format and naming |
| **Document** | Section headings, citation format, terminology, structural conventions (executive summary, appendices), voice and tense |
| **Mixed / verbal** | Ask one targeted question per ambiguous area — maximum 3 questions total |

### Step 3: Match signals to grimoire skills

Score each detected pattern against installed skills:

```
score = (tag_overlap × 2) + (description_match × 3) + (domain_plausibility × 1)
```

Surface only matches scoring ≥ 0.5. Group by domain.

### Step 4: Report applied practices

Present findings:

```
Applied practices detected in [project/artifact name]:

Domain                     Practice                        Evidence
───────────────────────────────────────────────────────────────────────────
✓ engineering/testing      write-unit-test                 Jest, co-located tests
✓ engineering/development  propose-conventional-commit     detected in git log
✓ engineering/devops       design-ci-pipeline              GitHub Actions workflow
? engineering/architecture                                 no clear pattern found
? engineering/security                                     no clear pattern found
```

### Step 5: Interactive preference pinning

For each ✓ result, ask in sequence:

```
Pin "propose-conventional-commit" as your intentional choice for engineering/development?
[y] pin  [n] skip  [r] pin with reason
```

If user selects `[r]`, ask: "Reason? (e.g. 'required by semantic-release')"

For `?` results, after all `✓` items are processed, offer once:

```
3 domains have no detected practice. Want to specify preferences for any of them?
[y] go through each  [n] skip
```

If yes, for each undetected domain, ask: "Which practice do you use for [domain]? (skill name or describe it)"

### Step 6: Write preferences file

Ask where to save:

```
Save preferences to:
  [0] This session only  → in memory; not written to disk; resets when session ends
  [1] This project only  → <project-root>/.grimoire/preferences.md
  [2] All my projects    → ~/.grimoire/preferences.md
                           (uses ~/.config/grimoire/preferences.md if XDG_CONFIG_HOME is set)
  [3] Both (project + global)
```

Write to selected location(s) in this format:

```markdown
# Grimoire Practice Preferences

<!-- Intentional choices. Routing skills will not suggest alternatives for pinned practices. -->

## engineering/development
- propose-conventional-commit: conventional commits format
  reason: required by semantic-release

## engineering/testing
- write-unit-test: Jest, co-located test files
  reason: team standard, enforced in CI
```

If file already exists at the target path: append new domain sections only. Never silently overwrite existing pins. If a domain conflict exists, ask:
```
You already have [existing-skill] pinned for [domain]. Replace it with [new-skill]? [y/n]
```

After writing, confirm:
```
Preferences saved to [path]. Routing skills will now respect these choices.
```

## Rules

- Never overwrite existing preferences without explicit confirmation
- Pin only what the user confirms as intentional — never auto-pin detected practices
- Precedence order: session > project > global > CLAUDE.md fallback
- Session-level pins are in-memory only — never written to disk, reset when session ends
- Project-level file takes precedence over global for any domain where both exist
- XDG compliance: use `$XDG_CONFIG_HOME/grimoire/preferences.md` if `XDG_CONFIG_HOME` is set, else `~/.grimoire/preferences.md`
- If no `.git` directory and no files are accessible, ask the user to describe their domain and practices verbally — maximum 3 questions
- After writing, always confirm the path and what was saved

## Common Mistakes

**Auto-pinning without asking**: always ask the user to confirm each pin. Detected ≠ intentional.

**Overwriting silently**: existing preference files may have carefully chosen entries. Append only, never overwrite without confirmation.

**Skipping undetected domains**: `?` domains are often the most important — the user may have strong preferences for architecture or security that don't show up in file scans.

**Forgetting to confirm write**: after writing the file, always confirm the path and what was saved.
