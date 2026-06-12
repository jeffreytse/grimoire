---
name: write-changelog
description: Use when writing or updating a CHANGELOG for a software project before a release
source: Keep a Changelog (keepachangelog.com); Conventional Commits specification (conventionalcommits.org)
tags: [changelog, release, documentation, conventional-commits, versioning]
verified: true
---

# Write Changelog

Produce a human-readable changelog that communicates what changed, for whom, and why.

## Why This Is Best Practice

**Adopted by:** Most major open-source projects (Node.js, Vue, Angular, Rails); mandated by Conventional Commits adopters
**Impact:** Changelogs reduce support tickets at release time by giving users a clear upgrade path; Angular's adoption of Conventional Commits automated changelog generation and halved release preparation time.

**Why best:** A changelog is a contract with users. It separates "what the code did" (commit log) from "what users need to know" (changelog). The Keep a Changelog format provides a widely understood, machine-parseable structure.

## Steps

1. **Choose format** — Use Keep a Changelog structure: `## [version] - YYYY-MM-DD` with subsections Added, Changed, Deprecated, Removed, Fixed, Security.
2. **Collect commits since last release** — Run `git log v1.2.0..HEAD --oneline` or parse Conventional Commit messages with a tool (standard-version, release-please, semantic-release).
3. **Translate commits to user impact** — Reframe technical changes as user-facing outcomes. "refactor: extract auth module" → omit; "feat: OAuth login" → Added.
4. **Highlight breaking changes** — Mark with `BREAKING CHANGE:` or a visible banner; explain migration steps.
5. **Write the Unreleased section first** — Keep `## [Unreleased]` at the top during development; move it to a versioned section at release.
6. **Link versions** — Add diff links at the bottom: `[1.3.0]: https://github.com/org/repo/compare/v1.2.0...v1.3.0`.

## Rules

- Write for users, not developers — omit internal refactors unless they affect public APIs.
- Never omit security fixes — always include them under the Security subsection.
- One entry per logical change, not per commit — consolidate related commits.
- Keep entries concise: one sentence describing the change and its user impact.

## Examples

```markdown
## [2.4.0] - 2026-05-01
### Added
- OAuth 2.0 login via Google and GitHub providers.
### Fixed
- Password reset email not sending when username contained special characters.
### Security
- Updated dependency `axios` to patch CVE-2026-12345 (SSRF vulnerability).
```

## Common Mistakes

- **Copying commit messages verbatim** — developer jargon is not user-facing communication.
- **Skipping the Unreleased section** — teams scramble at release to reconstruct what changed.
- **Combining multiple releases in one entry** — makes it impossible to correlate a bug with a version.
