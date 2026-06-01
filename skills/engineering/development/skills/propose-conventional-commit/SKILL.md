---
name: propose-conventional-commit
description: Use when the user asks to commit or wants a git commit message drafted — for any staged changes, whether or not they know the Conventional Commits format.
source: Conventional Commits specification (conventionalcommits.org), Angular commit guidelines, Semantic Release
tags: [commit-messages, conventional-commits, git, developer, changelog, semantic-versioning]
verified: true
---

# Propose Conventional Commit

Inspect staged changes and recent context, draft a Conventional Commits message, wait for approval.

## Why This Is Best Practice

**Adopted by:** Standard at Google, Meta, Stripe, and most teams with structured
engineering workflows. Conventional Commits format adopted by Angular, Vue, Electron,
and thousands of OSS projects as the dominant commit message standard.
**Impact:** Semantic Release (used by Angular, Electron, Vue) generates changelogs and
bumps semver fully automatically from Conventional Commits — eliminating 30–60 min of
manual release note writing per release. Type-prefixed commits also enable git bisect
efficiency: engineers can skip `style`/`chore` commits when bisecting regressions.
**Why best:** Freeform messages ("fix stuff", "wip") degrade under pressure and are
unsearchable. Gitmoji uses emoji prefixes — visually distinct but no tooling support
for automated changelogs. Conventional Commits has the broadest tooling ecosystem
(Semantic Release, commitlint, changelogen, release-please) and is the only format
with a published specification (conventionalcommits.org).

Sources: Conventional Commits specification, Angular commit message guidelines,
Semantic Release documentation, Google Engineering Practices

## Steps

### 1. Gather context

Run in parallel:
```bash
git diff --cached --stat
git diff --cached
git log --oneline -10
```

### 2. Draft message

- Derive type from nature of change: `feat`, `fix`, `refactor`, `style`, `test`, `build`, `docs`, `chore`
- Subject line: imperative mood, ≤72 chars, no period
- Body (optional): include only if the *why* is non-obvious — a hidden constraint, specific bug, or non-obvious trade-off
- **No `Co-Authored-By` trailer** — author is git user only

### 3. Present for approval

Show the proposed message as a plain fenced code block:

```
<type>(<optional scope>): <subject>

<optional body>
```

Then stop. Do NOT run `git commit` until the user explicitly approves (says "commit", "looks good", "yes", or sends the message back).

### 4. On approval

Run:
```bash
git commit -m "$(cat <<'EOF'
<approved message>
EOF
)"
```

Then report the commit hash.

## Rules

- Never commit without explicit user approval in this turn
- Never add `Co-Authored-By` lines
- If nothing is staged, say so and stop
- If the diff is ambiguous (e.g., multiple unrelated changes), note it and suggest splitting
