---
name: apply-gitflow
description: Use when a team has scheduled releases, maintains multiple supported versions simultaneously, or operates in a regulated environment requiring formal release approval gates.
source: Driessen "A Successful Git Branching Model" (nvie.com, 2010); git-flow CLI (Petersen, 2010); DORA State of DevOps Report 2019
tags: [git, branching, gitflow, release-management, versioning, scheduled-releases, regulated]
---

# Apply GitFlow

Implement the GitFlow branching model — two permanent branches (`main`, `develop`) plus short-lived `feature/`, `release/`, and `hotfix/` branches — to manage scheduled releases with clean separation between in-progress work and production code.

## Why This Is Best Practice

**Adopted by:** GitFlow became the de facto standard for scheduled-release software after Driessen's 2010 blog post; git-flow CLI (Petersen) ships natively in Homebrew, apt, and most git clients; widely used in enterprise Java, PHP, and regulated-industry software where batch QA gates and formal release approvals are required
**Impact:** GitFlow's branch model provides a provable audit trail — every production change traces to either a tagged release merge or a named hotfix — a requirement in SOC 2, ISO 27001, and FDA 21 CFR Part 11 regulated software environments; separating `develop` from `main` eliminates the "I deployed unreleased features by accident" class of incident common in single-branch models
**Why best:** GitHub Flow has no release stabilization concept — a hotfix on a busy `main` can inadvertently ship unreleased features; trunk-based development requires feature-flag infrastructure that teams in regulated or batch-release environments often cannot adopt; GitFlow's release branch gives teams a dedicated stabilization window with a hard boundary against new feature inclusion

Sources: Driessen "A Successful Git Branching Model" (nvie.com, 2010); Forsgren, Humble & Kim "Accelerate" (IT Revolution, 2018) §2 — identifies GitFlow as appropriate for release-gated environments; git-flow CLI (github.com/nvie/gitflow)

## Steps

### 1. Initialize the two permanent branches

```bash
git checkout -b develop main   # if develop does not yet exist
git push -u origin develop
```

Configure branch protections on both:
- `main`: no direct push, require CI pass, require at least 1 approval
- `develop`: require CI pass, require at least 1 approval

`main` = production. `develop` = integration target for all feature work. Never commit directly to either.

### 2. Feature development

```bash
# Start
git checkout -b feature/PROJ-42-user-auth develop

# Finish — open PR targeting develop, squash-merge on approval
git checkout develop
git merge --squash feature/PROJ-42-user-auth
git commit -m "feat(auth): add user authentication"
git branch -d feature/PROJ-42-user-auth
```

**Rules for feature branches:**
- Branch from `develop` only — never from `main` or a release branch
- PR target is `develop` — never `main`
- Delete after merge
- No feature branch lives beyond 2 days without an active PR; use feature flags for longer work

### 3. Release branch

Cut a release branch when `develop` has all features for the next version and the team is ready for QA:

```bash
# Cut
git checkout -b release/2.1.0 develop
```

On the release branch, only bug fixes are allowed — **no new features**. This is the stabilization window.

```bash
# Bug fix discovered during QA — commit directly to release branch
git checkout release/2.1.0
git commit -m "fix(auth): handle null token on logout"
```

When QA passes:

```bash
# Merge to main and tag
git checkout main
git merge --no-ff release/2.1.0
git tag -a v2.1.0 -m "Release 2.1.0"

# Merge back to develop (captures QA bug fixes)
git checkout develop
git merge --no-ff release/2.1.0

# Delete
git branch -d release/2.1.0
git push origin --delete release/2.1.0
git push origin main develop --tags
```

### 4. Hotfix branch

Production-only emergency fix — branches from `main`, not `develop`:

```bash
# Cut from main
git checkout -b hotfix/auth-bypass main

# Fix and commit
git commit -m "fix(auth): prevent token reuse after logout"
```

When the fix is ready:

```bash
# Merge to main and tag
git checkout main
git merge --no-ff hotfix/auth-bypass
git tag -a v2.0.1 -m "Hotfix 2.0.1"

# Merge to develop (critical: keeps develop in sync)
git checkout develop
git merge --no-ff hotfix/auth-bypass

# If a release branch is currently active, merge there too
# git checkout release/2.1.0 && git merge --no-ff hotfix/auth-bypass

# Delete
git branch -d hotfix/auth-bypass
git push origin main develop --tags
```

### 5. Branch naming conventions

| Branch | Pattern | Example |
|--------|---------|---------|
| Feature | `feature/<ticket>-<description>` | `feature/PROJ-42-user-auth` |
| Release | `release/<semver>` | `release/2.1.0` |
| Hotfix | `hotfix/<description>` | `hotfix/auth-bypass` |

Never use bare `fix/`, `bugfix/`, or `patch/` prefixes at top level — route through `feature/` (if targeting `develop`) or `hotfix/` (if targeting `main`).

### 6. Document in CONTRIBUTING.md

Add a one-section summary to CONTRIBUTING.md covering:
- Branch types and their purpose
- Where to branch from and where to merge to
- The release cut process (who cuts, on what trigger)
- Hotfix escalation path

Every contributor must be able to answer "where do I branch from?" from CONTRIBUTING.md alone.

## Rules

- `main` is production only — its history must be a sequence of tagged release merges and hotfix merges, nothing else
- `develop` is the integration branch — all feature PRs target `develop`
- Release branches are feature-freeze — no new features, only bug fixes
- Hotfix merges must land on **both** `main` and `develop` (and any active `release/*` branch) — skipping `develop` causes the bug to reappear in the next release
- Every merge to `main` must be tagged — untagged `main` commits are a sign GitFlow is not being followed
- Feature branches that outlive 2 days without a PR are a smell — break the work smaller or use a feature flag

## Common Mistakes

- **Feature PR targeting `main`** — bypasses `develop`; unreleased features can reach production in the next hotfix
- **Forgetting to merge release branch back to `develop`** — QA bug fixes are lost; they reappear in the next release
- **Hotfix merged to `main` only** — the fixed bug resurfaces when `develop` is next released
- **Cutting a release branch from `main` instead of `develop`** — the release contains only previous release content, not current features
- **No branch protection on `develop`** — direct commits bypass CI and code review; `develop` becomes unreliable as an integration target

## When NOT to Use

- When the team deploys more than once a week — GitFlow's batch pipeline adds overhead that blocks fast delivery; use trunk-based development or GitHub Flow instead
- When the team is ≤5 engineers with no formal QA gate — the release branch stabilization window adds process without proportional benefit; GitHub Flow suffices
- When the codebase has no CI pipeline — GitFlow assumes CI guards both `main` and `develop`; without it, the model provides structure but no safety
- When maintaining only one production version at a time with no patch support requirements — the multi-branch model adds complexity that trunk-based with tags handles more simply
