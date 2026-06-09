---
name: write-best-practice-profile
description: Use when the user wants to create a new practice profile — e.g., "create a profile for my team", "write an OOP profile", "make a backend-defaults profile".
source: Grimoire project conventions; XDG Base Directory Specification (freedesktop.org)
tags: [profile, skill-authoring, contributor, team, preferences]
related: [review-best-practice-profile, apply-best-practice-profile, pin-best-practice-preference]
---

# Write Best Practice Profile

Create a `.grimoire/profiles/<name>.toml` file that bundles a curated set of skills for a paradigm, team, or project context.

## Why This Is Best Practice

**Adopted by:** Every major toolchain uses guided authoring for shared configuration — ESLint's `--init`, `create-react-app`, `nx generate`, and VS Code's profile export all guide users through structured configuration creation rather than expecting manual file editing.
**Impact:** Manual TOML authoring leads to invalid skill names, missing required fields, and profiles that fail silently at activation time. A guided creation flow catches these errors before the file is written.
**Why best:** A profile is only as good as its skill list. Guided creation ensures each skill name is validated against installed skills at write time — not discovered as a broken reference at activation time.

Sources: XDG Base Directory Specification (freedesktop.org); Grimoire `docs/profiles.md`

## Steps

### 1. Gather profile identity

Ask:
- **Name** (slug, lowercase, hyphens): e.g. `my-team`, `backend-defaults`, `frontend-strict`
- **Description** (one sentence): when should someone activate this profile?

```
Profile name (slug): my-team
Description: Backend team defaults — DDD, SOLID, no framework shortcuts.
```

---

### 2. Collect skills

Two paths — offer both, user chooses:

**A. Manual list** — user names skills directly:
```
Add skills (enter names one by one, blank to finish):
  > apply-solid-principles
  > apply-domain-driven-design
  > apply-low-coupling
  >
```

**B. Tag-assisted** — run tag query, user selects subset:
```
Start from a tag (e.g. "oop", "tdd", "functional")? [tag or skip]:
  > oop

Found 9 skills tagged "oop":
  [x] apply-solid-principles
  [x] apply-law-of-demeter
  [ ] apply-composition-over-inheritance
  [x] apply-information-expert
  ...
Select skills to include (toggle with number, confirm with Enter)
```

---

### 3. Validate skill names

For each collected skill name, check it exists in installed grimoire skills.

- **Found**: mark ✓, include
- **Not found**: warn — ask user to confirm inclusion anyway (may be a skill they plan to install later)

```
Validating skills...
  ✓ apply-solid-principles       — found
  ✓ apply-domain-driven-design   — found
  ⚠ apply-twelve-factor-app      — not installed. Include anyway? [y/n]
```

---

### 4. Choose save location

```
Save to:
  [p] This project   → .grimoire/profiles/my-team.toml  (commit to repo to share with team)
  [u] My user level  → ~/.grimoire/profiles/my-team.toml (personal, applies to all projects)
```

---

### 5. Write the file

```toml
name = "my-team"
description = "Backend team defaults — DDD, SOLID, no framework shortcuts."

[[skills]]
name = "apply-solid-principles"

[[skills]]
name = "apply-domain-driven-design"

[[skills]]
name = "apply-low-coupling"
```

Confirm: `✓ Written to .grimoire/profiles/my-team.toml`

---

### 6. Next steps

```
Activate now?     profiles = ["my-team"] in .grimoire/settings.toml
Validate first?   /review-best-practice-profile my-team
Share publicly?   /share-best-practice-profile my-team
```

## Common Mistakes

**Inventing skill names.** Always validate against installed skills. A profile with a typo in a skill name silently activates nothing for that entry.

**One-size profiles.** Profiles should be scoped — a `backend-api` profile is more useful than a `company-wide` profile that covers too many contexts.
