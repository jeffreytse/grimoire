---
name: deprecate-skill
description: Use when a grimoire skill has become outdated — the referenced tool no longer exists at scale, the source institution revised its position, a newer practice has achieved majority top-tier adoption that supersedes it, or the skill fails review-skill criteria it once passed.
source: npm deprecation guidelines, IETF RFC obsolescence process (BCP 9), Semantic Versioning (semver.org)
tags: [skill-maintenance, staleness, deprecation, outdated-practice, maintainer, knowledge-freshness]
verified: true
---

# Deprecate Skill

Mark an outdated grimoire skill for removal and, where possible, point users to the
practice that supersedes it.

## Why This Is Best Practice

**Adopted by:** Explicit deprecation with replacement pointers is standard across all
major knowledge systems — npm packages use `deprecated` metadata with migration notes,
MDN Web Docs marks obsolete APIs with explicit "Use X instead" banners, and the Python
Package Index surfaces deprecation notices prominently to prevent installation of
abandoned packages.
**Impact:** Without deprecation, outdated skills get served to users indefinitely.
Research on documentation debt (Aghajani et al., 2019, MSR) found 50–60% of surveyed
developers encountered bugs caused by stale documentation; explicit deprecation markers
reduce this by surfacing decay at the point of use.
**Why best:** Silent removal leaves users with broken install paths and no migration
guidance. Deprecation-with-pointer gives users a safe off-ramp and maintains trust in
the repo as a reliable source.

Sources: npm deprecation model, MDN Web Docs obsolescence guidelines, Aghajani et al.
2019 (MSR documentation debt), IETF RFC obsolescence process (BCP 9)

## Steps

### 1. Confirm the skill meets staleness criteria

A skill qualifies for deprecation if ANY of the following are true:

- The source institution has revised its position and the skill no longer reflects it
- A newer practice has achieved majority top-tier adoption that directly supersedes this one
- The tool or API referenced in the skill no longer exists at scale (abandoned, renamed,
  or replaced by something with 10× adoption)
- The skill now fails `review-skill` criteria it once passed (e.g., the "Adopted by"
  claim is no longer accurate)

If none apply, the skill needs updating — not deprecating. Use `write-skill` to
revise it instead.

### 2. Identify the replacement (if one exists)

| Situation | Action |
|-----------|--------|
| A newer grimoire skill supersedes this one | Point to that skill |
| The practice is superseded but no grimoire skill covers the new approach yet | Note the superseding practice and flag for a new skill to be written |
| The practice was wrong / never qualified | No replacement — just remove |

### 3. Add a deprecation notice to the skill file

At the very top of the SKILL.md body (after frontmatter, before the title), add:

```markdown
> **Deprecated:** This skill is outdated. [Reason in one sentence.]
> Use [`replacement-skill-name`](../../replacement-path/SKILL.md) instead.
```

If no replacement exists:
```markdown
> **Deprecated:** This skill is outdated. [Reason in one sentence.]
> No direct replacement — see [domain]/[subdomain] for related skills.
```

Do NOT delete the skill file in this step — deprecation precedes removal to give
users time to migrate.

### 4. Update marketplace.json

Find the entry for this skill's domain/sub-domain and add a `deprecated` flag:

```json
{
  "name": "grimoire-engineering-some-subdomain",
  "source": { ... },
  "description": "...",
  "deprecated": true,
  "deprecation_note": "Superseded by grimoire-meta-replacement-skill"
}
```

### 5. Open a PR

PR title: `deprecate(<domain>/<skill-name>): <one-line reason>`

PR description must include:
- What changed that made the skill stale (source revision, tool abandonment, etc.)
- Link to the superseding practice or skill
- Whether a new skill needs to be written to cover the replacement practice

### 6. Removal (maintainers only)

Maintainers remove deprecated skills in the next release cycle after the PR merges.
Removal steps:
1. Delete the skill directory
2. Remove from `marketplace.json`
3. Remove from the domain's `.claude-plugin/plugin.json`
4. Tag the release with a breaking-change note if the skill had significant install volume

## Rules

- Never deprecate without a stated reason — "Deprecated" with no explanation is worse
  than nothing
- Always check if a replacement skill should be written before deprecating
- Deprecation PR and removal are two separate steps — never delete in the same PR as
  the deprecation notice
- If a skill is wrong but not stale (e.g., always had incorrect "Adopted by" claims),
  use `review-skill` + a fix PR instead of deprecation
- Only maintainers perform step 6 (removal)

## Common Mistakes

**Removing instead of deprecating**: deleting a skill file in one PR gives users no
warning and breaks install paths. Always deprecate first, remove in the next cycle.

**Deprecating because you disagree**: deprecation is for staleness, not preference.
If the practice is still majority top-tier adopted, open a revision PR instead.

**No replacement pointer**: "Deprecated" with no guidance forces users to start from
scratch. Always name the superseding practice, even if no grimoire skill covers it yet.
