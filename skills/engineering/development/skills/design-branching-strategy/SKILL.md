---
name: design-branching-strategy
description: Use when choosing or designing a Git branching model for a team or project
source: Trunk-Based Development (trunkbaseddevelopment.com); "A successful Git branching model" — Vincent Driessen; GitHub Flow docs
tags: [git, branching, trunk-based, gitflow, github-flow, devops]
verified: true
---

# Design Branching Strategy

Select and configure a Git branching model that matches team size, release cadence, and deployment frequency.

## Why This Is Best Practice

**Adopted by:** Google (mono-trunk), Netflix, Etsy (trunk-based); GitFlow used widely in release-gated software
**Impact:** DORA research (Accelerate) shows trunk-based development correlates with elite software delivery performance — 46x more frequent deployments, 440x faster lead time.

**Why best:** The right branching strategy reduces merge conflicts, clarifies code ownership, and aligns version control with deployment processes. Mismatched strategies (e.g., GitFlow for a team deploying daily) create unnecessary overhead.

## Steps

1. **Assess release cadence** — Continuous deployment favors trunk-based or GitHub Flow; scheduled releases favor GitFlow.
2. **Assess team size** — Small teams (<10): GitHub Flow. Medium teams: trunk-based with feature flags. Large/multi-release: GitFlow or release branches.
3. **Define branch types** — Document which branch types exist (main, feature, release, hotfix), their naming convention, and lifetime expectations.
4. **Set merge policies** — Require pull requests, minimum approvals, and passing CI before merge to main.
5. **Configure branch protections** — Enforce rules in GitHub/GitLab: no direct pushes to main, require status checks.
6. **Document the strategy** — Add a one-page CONTRIBUTING.md section so every contributor follows the same model.
7. **Review periodically** — Revisit the strategy as team size or release cadence changes.

## Rules

- Long-lived feature branches (>2 days) accumulate merge debt — use feature flags instead.
- `main` must always be deployable.
- Hotfix branches must be merged to both main and any active release branches.
- Avoid "GitFlow for startups" — the overhead exceeds the benefit at low team sizes.

## Examples

GitHub Flow (continuous deployment):
- Branch from `main` → `feature/user-auth`
- Open PR → review → CI green → squash-merge → deploy.

GitFlow (scheduled releases):
- `develop` integrates features; `release/1.4` is cut for stabilization; hotfixes branch from `main`.

## Common Mistakes

- **Trunk-based without feature flags** — incomplete features ship to production.
- **GitFlow on a CD team** — `develop` diverges from `main`; integration pain increases.
- **No branch naming convention** — branches like `fix` or `test2` are untrackable.

## When NOT to Use

- When working as a solo developer on a personal project with no CI/CD pipeline, the overhead of a formal branching strategy exceeds its coordination benefit; commit directly to main.
- When the repository hosts infrastructure-as-code or configuration files that are always applied from a single source of truth, a multi-branch model introduces drift risk that outweighs its benefits.
- When a team has already converged on a functioning strategy and deployment metrics are healthy, redesigning the branching model creates churn without a concrete problem to solve.
