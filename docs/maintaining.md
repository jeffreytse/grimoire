# Maintaining grimoire

This guide covers the maintainer workflow: reviewing PRs, running domain audits, managing the deprecation lifecycle, and releasing new versions. Contributors write and submit skills; maintainers gate quality and merge.

---

## Maintainer responsibilities

| Task | Who |
|------|-----|
| Review skill PRs (apply review-skill) | Maintainers |
| Run periodic domain audits (audit-domain) | Maintainers |
| Execute deprecation removal | Maintainers only |
| Write skills, open PRs | Contributors |
| Self-review before submitting | Contributors |

The separation is intentional: contributors cannot merge their own skills. All skill PRs pass through `review-skill` before merge.

---

## PR review workflow

For every incoming skill PR:

### 1. Apply review-skill

Invoke `review-skill` on the skill file:

```
Applying review-skill to skills/<domain>/<subdomain>/skills/<name>/SKILL.md
```

`review-skill` produces a structured verdict with per-criterion findings.

### 2. Act on the verdict

**PASS** → merge. No further review needed.

**NEEDS-REVISION** → comment on the PR with the specific findings from `review-skill`. Ask the contributor to invoke `revise-skill` to address each one. Example comment:

```
review-skill findings:

- Impact: "significantly improves" — needs a cited number or named study
- tags: missing outcome axis — add a tag describing the result the user gets

Please use `revise-skill` to address these before we merge.
```

**REJECT** → a required field is missing or a core criterion fails. Close the PR with a clear explanation of which criterion failed and what it needs. Contributor can re-open after a full fix.

### 3. One skill per PR

If a PR contains multiple skills, ask the contributor to split it before reviewing. Multi-skill PRs are harder to review accurately and harder to revert cleanly.

### 4. New domain or subdomain PRs

If the PR adds a new domain or subdomain:

- Verify the directory structure matches the pattern in [architecture.md](./architecture.md)
- Verify plugin.json follows the naming convention: `grimoire-<domain>` / `grimoire-<domain>-<subdomain>`
- Verify marketplace.json has entries for the new domain and each subdomain
- Verify at least 2 seed skills are included and all pass review-skill
- Apply `audit-domain` to the new domain to confirm all seed skills PASS

---

## Domain audits

### Weekly audit

Run `audit-domain` on each domain that has received contributions in the past week:

```
Applying audit-domain to skills/engineering/
Applying audit-domain to skills/health/
...
```

`audit-domain` produces a summary: Total / PASS / NEEDS-REVISION / REJECT per domain.

### Acting on findings

**NEEDS-REVISION findings** → open a GitHub issue for each affected skill, assigned to the original contributor if identifiable. Label: `needs-revision`. Link to the specific review-skill findings.

**REJECT findings** → the skill should not be in the library. Either:
- Open a revision PR yourself (if the fix is straightforward — e.g., vague "Adopted by" can be named)
- Trigger the deprecation lifecycle if the practice never qualified

**Tracking quality over time** → record domain audit counts in a CHANGELOG entry each release. A domain improving from "3 PASS, 2 NEEDS-REVISION" to "5 PASS" is a meaningful quality signal.

---

## Deprecation lifecycle

Deprecation is a two-step process: mark first, remove after one release cycle.

### Step 1: Deprecation PR

Invoke `deprecate-skill` to guide the process. The PR must:

1. Add a deprecation notice at the top of the skill's SKILL.md body (after frontmatter, before the title):

```markdown
> **Deprecated:** This skill is outdated. [Reason in one sentence.]
> Use [`replacement-skill-name`](../../replacement-path/SKILL.md) instead.
```

2. Add a `deprecated: true` flag to the skill's marketplace.json entry

PR title format: `deprecate(<domain>/<skill-name>): <one-line reason>`

### Step 2: Removal PR (maintainers only)

After the deprecation PR merges and one release cycle passes:

1. Delete the skill directory
2. Remove the entry from `.claude-plugin/marketplace.json`
3. Remove the path from the domain's `.claude-plugin/plugin.json` `skills` array
4. If the skill had significant install volume, tag the release with a breaking-change note

**Never delete in the same PR as the deprecation notice** — users need time to migrate.

---

## Release process

grimoire uses semantic versioning at the domain level (each domain plugin.json has its own version).

### Version bump rules

| Change type | Version bump |
|------------|-------------|
| Bug fix in skill content, citation correction | Patch (0.0.x) |
| New skill added to domain | Minor (0.x.0) |
| Skill removed (deprecation removal) | Minor (0.x.0) |
| Subdomain restructure, renamed paths | Major (x.0.0) |

### Release steps

1. Bump `version` in affected domain plugin.json files
2. Update `.claude-plugin/marketplace.json` if any `source.path` values changed
3. Add a CHANGELOG entry (see below)
4. Tag the release on GitHub: `git tag v<version>`
5. Push the tag: `git push origin v<version>`

### CHANGELOG format

```markdown
## [0.2.0] - YYYY-MM-DD

### Added
- engineering/development: review-pull-request

### Fixed
- engineering/development: propose-conventional-commit — corrected Adopted by claim

### Deprecated
- engineering/development: old-skill-name — superseded by new-skill-name

### Removed
- engineering/testing: stale-skill — deprecated in 0.1.0
```

---

## Handling contributor disputes

**"My skill should PASS but review-skill rejected it"**
→ The criteria are in STANDARD.md and are not negotiable per-contributor. If the criterion itself is wrong, open a PR to revise STANDARD.md and run `revise-skill` on all meta skills that reference the changed criterion.

**"This practice is contested — why can't I encode my position?"**
→ If the split is 50/50 among top-tier professionals, no skill can encode either position as a best practice. If there's a clear majority position, encode it. Refer to STANDARD.md §Contested practices.

**"I want to add a skill for an emerging practice"**
→ Emerging practices not yet proven at scale don't qualify. grimoire encodes practices with demonstrated impact, not trends. Revisit when majority top-tier adoption is established.
