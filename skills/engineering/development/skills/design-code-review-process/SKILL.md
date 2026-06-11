---
name: design-code-review-process
description: Use when establishing or improving a team's code review workflow, when reviews are causing bottlenecks, or when quality issues are escaping review into production
source: Google Engineering Practices Documentation (2019); Cohen et al. "Best Kept Secrets of Peer Code Review" SmartBear (2006); DORA State of DevOps Report (Forsgren et al. 2019)
tags: [development, code-quality, process, team]
verified: true
---

# Design Code Review Process

Establish a code review process that catches defects, spreads knowledge, and maintains velocity without creating bottlenecks.

## Why This Is Best Practice

**Adopted by:** Google (mandatory for all commits), Microsoft, Amazon, Meta — all require code review before merge; codified in Google Engineering Practices public documentation
**Impact:** SmartBear study: code review finds 60% of defects before testing; DORA research shows teams with effective code review deploy 46x more frequently with 5x lower change failure rate; review sessions >60 minutes lose effectiveness
**Why best:** Humans write code with blind spots that fresh eyes catch; code review distributes knowledge, enforces standards, and prevents technical debt accumulation

Sources: Google "Engineering Practices Documentation" (2019, opensource.google); Cohen et al. "Best Kept Secrets of Peer Code Review" SmartBear (2006); Forsgren, Humble & Kim "Accelerate" (2018)

## Steps

1. **Define review scope and SLA** — Specify what requires review (all merges to main, not experimental branches), who must approve (1 reviewer for low-risk, 2 for high-risk), and the SLA (respond within 1 business day, resolve within 2). Without SLAs, reviews become bottlenecks. Google's standard: first response within 1 business day. **Critically: fast *individual responses* matter more than fast overall process.** A reviewer who responds in hours — even to say "I'll look at this fully tomorrow" — eliminates most developer frustration with code review. A reviewer who responds every 2 days while requesting major changes each time produces legitimate complaints about review being too slow, even if their standards are correct. Most complaints about "strict reviewers" disappear when those same reviewers respond quickly.

2. **Set PR size guidelines** — Maximum 400 lines changed per PR (SmartBear data: review effectiveness drops sharply above 400 lines). Enforce this via PR description templates or CI checks. Large PRs should be decomposed into a stack of smaller PRs with a feature flag. Smaller PRs get faster, more thorough reviews.

3. **Write PR descriptions that enable review** — Structure: **First line** = short, focused summary of *what* changed (imperative mood, ≤72 chars); **Body** = *why* this change was made, the problem it solves, the approach chosen and its tradeoffs, any limitations, links to tickets/design docs. The first line is what appears in `git log` — it must stand alone. The body provides the context reviewers need to evaluate the approach rather than just the implementation. A PR without a body receives surface-level review; a PR with a clear body enables reviewers to focus on the right questions.

4. **Distinguish blocking from non-blocking feedback** — Use explicit prefixes: `[blocking]` must be fixed before merge, `[nit]` stylistic preference the author may ignore, `[question]` genuine question (not critique), `[suggestion]` optional improvement. This prevents nit-picking from blocking merges on correctness.

5. **Establish reviewer assignment policy** — Define: who must review (code owner), who may review (anyone), and load balancing (rotate reviewers, maximum 3 PRs assigned simultaneously). Use CODEOWNERS files to auto-assign. Unassigned PRs sit in a queue no one feels responsible for.

6. **Checklist for reviewers** — Review for: correctness (does it do what it claims?), tests (are edge cases covered?), security (injection, auth bypass, secrets in code), performance (N+1 queries, missing indexes), maintainability (naming, complexity, dead code). Not aesthetics (handled by linter).

7. **Automate what automation handles** — Lint, format, style, and test coverage checks must run in CI before human review begins. Reviewers should never comment on formatting that a linter catches. Configure auto-merge for dependabot when CI passes. Humans review logic, not style.

8. **Handle disagreements with a clear escalation path** — When author and reviewer disagree after one round of discussion: author explains rationale; if still disagreed, escalate to tech lead within 24 hours. Unresolved disagreements that block merge for > 2 days indicate a process problem.

9. **Measure review health** — Track: median time-to-first-review, median time-to-merge, average PR size, and number of review rounds per PR. Review health correlates directly with team velocity. Alert when time-to-first-review exceeds SLA. Review these metrics monthly.

10. **Conduct periodic process retrospectives** — Quarterly: ask the team what's working and what's frustrating about the review process. Adjust SLAs, size guidelines, or reviewer rotation based on feedback. Code review process that doesn't evolve becomes a bureaucratic friction point.

## Rules

- Reviewers critique code, not authors; all feedback is about the change, not the person.
- The author is responsible for getting their PR reviewed; don't wait passively, follow up after 24 hours.
- Approval does not eliminate author responsibility; the author owns the change in production regardless of reviewer sign-off.
- PRs over 600 lines are refused and returned for decomposition; no exceptions.

## Common Mistakes

- **Reviewing for style instead of correctness** — misallocates reviewer attention to low-value feedback; automate style with linters.
- **Approving without reading** — rubber-stamp approvals defeat the purpose of review; require evidence of review (specific comments) before auto-merge rules apply.
- **No PR size limit** — 1000-line PRs are reviewed superficially; reviewers scan rather than read, missing logical errors.
- **Conflating code review with architecture review** — architectural concerns should be resolved in design docs before implementation; code review is too late to redesign.

## When NOT to Use

- Hotfixes during active production incidents where time-to-fix outweighs review benefit (review post-incident)
- Internal tooling with a single maintainer and no external consumers, where review overhead exceeds risk
- Prototype branches explicitly marked as throwaway that will be rewritten before merging to main
