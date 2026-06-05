---
name: apply-trunk-based-development
description: Use when teams want to increase deployment frequency, reduce merge conflicts, or adopt continuous integration and delivery practices
source: Forsgren, Humble & Kim "Accelerate" (2018); DORA research identifying trunk-based development as elite performance predictor; Paul Hammant trunkbaseddevelopment.com
tags: [development, ci-cd, git, devops]
verified: true
---

# Apply Trunk-Based Development

Integrate code to a single shared branch (trunk/main) continuously — at minimum daily — using short-lived feature branches and feature flags to decouple deployment from release.

## Why This Is Best Practice

**Adopted by:** Google (single monorepo trunk), Facebook, Netflix, Amazon — DORA research found trunk-based development is one of the strongest predictors of elite delivery performance
**Impact:** DORA "Accelerate" (2018): trunk-based teams deploy 46x more frequently with 5x lower change failure rates; reduces integration overhead by 70% vs. long-lived branches; merge conflicts drop to near-zero
**Why best:** Long-lived branches defer integration pain, create merge hell, and hide bugs until the worst possible time; trunk-based development makes integration continuous and cheap

Sources: Forsgren, Humble & Kim "Accelerate" IT Revolution Press (2018); DORA State of DevOps Report 2019; Hammant "Trunk Based Development" (trunkbaseddevelopment.com)

## Steps

1. **Establish main as the permanent trunk** — Protect main: require CI passing before merge, require at least one approval, disallow force-push. This is the single source of truth. All deployments come from trunk. Tag releases from trunk; never create release branches.

2. **Enforce short-lived branches** — Feature branches must be merged within 1-2 days. Branches older than 2 days are a smell; branches older than 5 days are a problem. Enforce via CI: fail the pipeline on branches with no merge activity for 3 days. Long-lived branches reintroduce the integration problems trunk-based development solves.

3. **Implement feature flags** — Decouple deployment from release using feature flags (LaunchDarkly, Unleash, environment variables). New features are deployed to production but hidden behind flags. This allows integration of incomplete features without exposing them to users. Flags also enable instant rollback without redeployment.

4. **Commit to trunk at least daily** — Every engineer integrates their changes to trunk at least once per day. Daily integration is the core practice; less frequent integration is not trunk-based development. If work-in-progress is incomplete, it goes behind a feature flag.

5. **Maintain a passing CI pipeline on trunk at all times** — Trunk must be green. Every commit triggers CI. If a commit breaks trunk, fixing it is the highest priority — above new feature work. Use a team agreement: "if you break trunk, you drop everything and fix it." A broken trunk blocks all deployments.

6. **Write branch-by-abstraction for large refactors** — For changes that touch many files over multiple days: introduce a new abstraction, implement it in parallel with the old, gradually migrate callers, delete the old implementation. This allows incremental integration of large changes without a long-lived branch.

7. **Keep PRs small** — Each PR should be < 400 lines and represent one logical change. Break features into a sequence of small, safe, independently reviewable PRs. A large feature is a series of small integrations, not a large merge at the end. Stacked PRs (PR2 based on PR1) are acceptable.

8. **Use ship/show/ask classification for PR review** — Ship: trivial changes merged without review (typos, dependency bumps with no behavior change). Show: changes merged, then shown to team for async comment. Ask: changes that need discussion before merge. Reduces review bottleneck without sacrificing quality on high-stakes changes.

9. **Automate deployment from trunk** — Every merge to trunk that passes CI should automatically deploy to a staging environment. Production deployment should be a one-step process (merge + manual trigger or automatic). Manual deployment steps are bottlenecks and sources of error.

10. **Manage release cadence with tags, not branches** — Releases are tagged commits on trunk (v1.2.3). Cherry-picks to a release branch are an anti-pattern — fix forward on trunk and release a patch. Release branches that diverge from trunk accumulate debt that never merges back.

## Rules

- Trunk must always be deployable; never commit code that knowingly breaks tests or the build.
- Feature flags are not technical debt when used correctly; they are a deployment safety mechanism; retire them after full rollout.
- No branch lives beyond 2 days without a compelling, reviewed exception; long-lived branches are a symptom of integration avoidance.
- Breaking trunk is a team emergency, not an individual's problem.

## Common Mistakes

- **Merging long branches on Fridays** — integrating large branches at the end of the week before a weekend creates emergency situations; merge small, merge often, merge early in the week.
- **Skipping feature flags** — integrating incomplete features without flags breaks trunk-is-deployable; flags are mandatory for in-progress work.
- **Treating trunk-based as "no branches"** — short-lived branches (< 2 days) are fine and recommended for code review; trunk-based means frequent integration, not no branching.
- **No CI enforcement** — trunk-based development without automated tests is just merging chaos; CI is the enforcement mechanism that makes it safe.

## When NOT to Use

- Open source projects with external contributors who need long-lived forks (use a fork model, not trunk-based)
- Regulated environments requiring formal release approval processes that take > 2 days (adapt with release trains)
- Teams with < 50% automated test coverage where trunk breakage is common and hard to detect
