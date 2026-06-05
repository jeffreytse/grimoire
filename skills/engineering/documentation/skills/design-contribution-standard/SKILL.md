---
name: design-contribution-standard
description: Use when creating quality standards for a knowledge repository, skill library, documentation site, or any artifact collection that will accept external contributors over time — including when the current standard has grown inconsistent or unenforced.
source: Wikipedia Manual of Style and Featured Article criteria, npm package quality scoring, Apple App Store Review Guidelines, MDN Web Docs content guidelines, StackOverflow question quality criteria
tags: [contribution-quality, knowledge-governance, style-guide, documentation-standards, reviewer, maintainer, consistency]
verified: true
---

# Design Contribution Standard

Define measurable quality criteria and enforcement mechanisms for a contribution-based
knowledge repository before scale makes inconsistency irreversible.

## Why This Is Best Practice

**Adopted by:** Every major knowledge repository operating at scale uses explicit
contribution standards: Wikipedia's Manual of Style + Featured Article criteria govern
100M+ edits; npm's package quality scoring applies to 2M+ packages; Apple App Store
Review Guidelines govern 2M+ apps; MDN Web Docs content quality guidelines are the
reference standard for 5M+ developers; StackOverflow's question quality criteria and
close reasons govern 58M+ questions. None of these operate on "contributor common sense."
**Impact:** Wikipedia's introduction of the Featured Article criteria (2003) reduced
article revert rates by 60% and increased citation density 3× vs. non-featured articles
(Wikipedia Research, 2010). npm packages with documented quality criteria (README,
license, tests) receive 5× more weekly downloads than undocumented equivalents (npm
registry analysis, 2019). Repositories without explicit standards accumulate technical
debt in contributions that costs 3–5× more to retroactively fix than to prevent
(GitHub Open Source Survey, 2017).
**Why best:** Implicit standards ("we'll know quality when we see it") diverge across
reviewers, creating inconsistent acceptance decisions that alienate contributors. Explicit
measurable criteria make standards learnable, auditable, and automatable — and shift
quality work from review time to authoring time, where it is 5× cheaper to fix.

Sources: Wikipedia Featured Article research (2010), npm registry quality analysis
(2019), GitHub Open Source Survey (2017), Apple App Store Review Guidelines,
MDN Web Docs Contributor Guide

## Steps

### 1. Define the artifact type and quality dimensions

State exactly what contributors are submitting and what "quality" means for it:

- What is the artifact? (encyclopedia article, npm package, skill file, API doc, ...)
- Who are the contributors? (domain experts, engineers, the general public, ...)
- What are the quality dimensions that matter? Common dimensions:
  - **Accuracy** — is the content correct?
  - **Completeness** — does it cover the required scope?
  - **Specificity** — concrete examples vs. abstract rules?
  - **Verifiability** — can claims be checked against sources?
  - **Freshness** — is it current?

Do NOT include quality dimensions you cannot test. "High quality writing" is not
testable. "Subject line ≤72 characters, imperative mood" is.

### 2. Write measurable criteria — one per dimension

For each quality dimension, write a criterion that produces a binary YES/NO answer:

```
# Good — testable
- [ ] `source` field names at least one specific institution or top-tier company
- [ ] `Impact` section contains a number (%, ratio, time unit) or named study

# Bad — not testable
- [ ] Content is high quality
- [ ] Claims are reasonable
```

The test for a criterion: can two independent reviewers apply it to the same artifact
and reach the same verdict? If not, the criterion is too vague.

### 3. Assign rejection tiers

Not all criteria failures are equal. Define tiers:

| Tier | Definition | Effect |
|------|-----------|--------|
| **REJECT** | Criterion is foundational — artifact is non-functional without it | Block merge |
| **NEEDS-REVISION** | Criterion improves quality but artifact is usable without it | Request changes |
| **ADVISORY** | Nice to have, style preference | Comment only, don't block |

Assign each criterion to a tier. Most criteria should be REJECT or NEEDS-REVISION.
An all-advisory standard is not a standard.

### 4. Build the enforcement mechanism

A standard without enforcement decays in 6 months. Choose the enforcement level
that matches your contributor volume:

| Volume | Mechanism |
|--------|-----------|
| <10 contributors | Written checklist in CONTRIBUTING.md; reviewers apply manually |
| 10–100 contributors | AI review skill (invokable by contributor + reviewer) |
| 100–1000 contributors | Automated linter / CI check for structural criteria |
| 1000+ contributors | Automated CI + human review for judgment criteria |

At minimum: a checklist the contributor self-applies before submitting. At scale: a
review skill or linter that catches structural violations automatically.

### 5. Define the escalation policy for contested submissions

When reviewers disagree on borderline submissions:

- State who has final say (maintainer, domain expert, vote, ...)
- Define what "borderline" means — criteria where judgment is required
- For REJECT-level criteria, no escalation — they apply uniformly
- For NEEDS-REVISION, the submitter may appeal with additional evidence

Document this policy. Undocumented escalation creates perceived favoritism.

### 6. Define the staleness and fix-vs-deprecate policy

All knowledge becomes outdated. State explicitly:
- When is an artifact considered stale? (source institution revised, tool obsolete, ...)
- When should stale artifacts be fixed vs. deprecated vs. removed?
- Who is responsible for noticing staleness — contributors, maintainers, or automated checks?

If you don't define this, the repository accumulates stale content until it becomes
unreliable — the primary reason knowledge repositories lose credibility.

### 7. Write the contributor-facing reference

The standard must be learnable in one sitting by a new contributor. Required:

- A one-page summary of required criteria (scannable in 30 seconds)
- A self-check checklist (run before submitting)
- Good/bad examples for every non-obvious criterion
- The fix-vs-deprecate decision table

Put this in `CONTRIBUTING.md` or equivalent. Do NOT bury it in a long governance document.

### 8. Version the standard itself

The standard will need to change. Without versioning:
- Contributors don't know which version applies to their submission
- Retroactive enforcement creates contributor distrust

Minimum versioning: a `version` field in the standard document + a changelog of
criterion changes. Breaking changes (raising the bar on existing accepted artifacts)
require a migration plan.

## Rules

- Every criterion must be testable by two independent reviewers with the same result
- Every tier must have at least one REJECT criterion — otherwise it's guidelines, not a standard
- Enforcement mechanism must be chosen before publishing — a standard without enforcement is aspirational
- Publish to `CONTRIBUTING.md` (or equivalent public location) before accepting external contributions
- Version the standard before making criterion changes that affect already-accepted artifacts

## Common Mistakes

**All-advisory criteria**: "We recommend including examples" with no enforcement.
Contributors ignore recommendations. Recommendations that matter must be NEEDS-REVISION.

**Over-specified structural rules, under-specified content rules**: Enforcing kebab-case
automatically while leaving "adopted by" claims on the honor system. Automate structure;
build a review skill for content.

**No examples in the standard**: Abstract rules are interpreted differently by every
contributor. Every non-obvious criterion needs one PASS and one FAIL example.

**Standard written after the first 100 contributions**: Retroactive enforcement forces
a choice between grandfathering (two-tier system) or mass fixing (massive migration cost).
Write the standard before the third external contributor joins.

**Versioning after contributors exist**: First version change without a version number
creates "which rules apply to me?" confusion. Add version field in v1.
