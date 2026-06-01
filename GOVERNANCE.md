# Governance

This document defines how the Grimoire Skill Standard and the grimoire project
are governed — who can do what, and how decisions are made.

---

## Roles

| Role | Who | Authority |
|------|-----|-----------|
| **Maintainer** | Project founders and appointed contributors | Merge PRs, release versions, resolve disputes, amend the standard |
| **Contributor** | Anyone who opens a PR | Add or fix skills; propose standard changes |
| **Adopter** | Projects listed in [ADOPTERS.md](./ADOPTERS.md) | Influence standard direction; listed in adopter registry |

**Current maintainers:** @jeffreytse

---

## Contributing Skills

See [CONTRIBUTING.md](./CONTRIBUTING.md). Skill PRs require:
- CI passing (`validate-skill.sh --all` and `test-schema.sh`)
- At least one maintainer approval

---

## Changing the Standard

[STANDARD.md](./STANDARD.md) changes require a higher bar than skill PRs — the
standard governs all skills and all adopters.

### Process

1. **Open an issue** labeled `standard-change` describing the proposed amendment,
   the problem it solves, and any impact on existing skills.
2. **Discussion period** — 7 days minimum for community feedback before a PR is
   opened. Adopters are encouraged to weigh in.
3. **Open a PR** referencing the issue. Include:
   - The exact change to STANDARD.md
   - Impact assessment: does this break any existing skills? Does `validate-skill.sh --all` still pass?
   - Update to the conformance test suite if the change affects validation rules
4. **Review** — requires 2 maintainer approvals (or 1 if only 1 maintainer exists).
5. **Merge** — maintainer merges and bumps the version per the version policy below.

### What requires a standard-change issue

| Requires standard-change process | Does NOT require |
|-----------------------------------|-----------------|
| Adding or removing a required frontmatter field | Fixing a typo in STANDARD.md |
| Changing validation rules in `validate-skill.sh` | Clarifying existing prose |
| Adding a new lifecycle state | Adding an example or table row |
| Changing the approved verb list | Fixing a broken link |
| Modifying domain safety requirements | Updating the changelog |

---

## Version Policy

The standard uses [Semantic Versioning](https://semver.org):

| Change type | Version bump | Example |
|-------------|-------------|---------|
| **Major** (breaking) | `X.0` | Removing a required field; changing a validation rule that breaks existing skills |
| **Minor** (additive) | `x.Y` | Adding a new optional field; adding a new lifecycle state; adding domain safety rules |
| **Patch** (clarifying) | `x.y.Z` | Prose clarification; new examples; fixing ambiguity without changing validation behavior |

Breaking changes (major bumps) require a migration guide in the PR describing how
existing skills and validator implementations should update.

---

## Decision-Making

grimoire uses **rough consensus**: a decision is made when no maintainer has a
principled objection. Agreement does not require unanimity.

For standard changes, if consensus cannot be reached, the lead maintainer has final
say after a 14-day deliberation period.

---

## Conflict Resolution

1. Discuss in the relevant issue or PR.
2. If unresolved after 14 days, tag a maintainer for a binding decision.
3. If a maintainer is a party to the conflict, another maintainer decides.

---

## Amendments to This Document

Changes to GOVERNANCE.md follow the same process as standard changes above, with
the same 7-day discussion period and 2-maintainer approval requirement.
