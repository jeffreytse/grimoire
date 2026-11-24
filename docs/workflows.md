# Workflows

How grimoire's meta workflows route and combine practices.

← [Back to the README](../README.md)

> **MCP mode:** skills are not installed as agent slash commands, so run any workflow below as `/grimoire:skill <name>` or `grimoire:<name> [task]` (e.g. `grimoire:suggest-best-practice my team is losing confidence`). When a workflow invokes another skill or scores installed skills, grimoire tells the agent to use `grimoire_get_skill` / `grimoire_search_skills` for that. The instruction line, server instructions and session hook route every request through `start-best-practice`, and `grimoire_get_skill` takes a `names` list so a plan can load all its skills in one call.

| Your situation                                                                                                | Start here                                |
| ------------------------------------------------------------------------------------------------------------- | ----------------------------------------- |
| Know exactly which skill you need                                                                             | `/skill-name` directly                    |
| Have a problem, unsure which skill                                                                            | `/suggest-best-practice`                  |
| Already have a plan, want gaps checked                                                                        | `/review-best-practice-fit`               |
| Needs 2+ practices coordinated — within one domain or across many — and sub-problems are identifiable upfront | `/plan-best-practice-solution`            |
| Complex problem where sub-problems are opaque and emerge through execution                                    | `/apply-best-practice-tree`               |
| Don't know what practices exist for a topic                                                                   | `/discover-best-practices`                |
| About to start a task — want to catch gaps before you begin                                                   | `/start-best-practice`                    |
| Problem isn't clear yet — need to define it before solving                                                    | `/analyze-best-practice-problem`          |
| Activate a paradigm's best practices (OOP, TDD, etc.)                                                         | `/apply-best-practice-profile`            |
| Align any project or artifact to stated best practice preferences (BPDD)                                      | `/apply-best-practice-driven-development` |
| Check if any artifact aligns with stated best practice preferences                                            | `/check-best-practice-compliance`         |
| Have a specific compliance finding to fix                                                                     | `/fix-best-practice-finding`              |
| Two practices exist — want side-by-side comparison                                                            | `/compare-best-practices`                 |
| Two practices conflict — want to reason through which fits                                                    | `/resolve-best-practice-conflict`         |
| Resolved a conflict — want to save the decision for future sessions                                           | `/pin-best-practice-preference`           |
| Get oriented on grimoire itself — what it is, why, how to use it                                              | `/learn-grimoire`                         |
| Actually internalize a practice — active recall, spaced review, so you can use it without AI                  | `/learn-best-practice`                    |
| Understand what a skill does and why it exists                                                                | `/explain-best-practice`                  |
| Adapt a skill to a different context or constraint                                                            | `/adapt-best-practice`                    |
| Teach a practice to someone else (structured walkthrough)                                                     | `/teach-best-practice`                    |
| Installing/uninstalling skills, upgrading grimoire, checking install health                                   | `/install-grimoire`                       |
| Want an inventory of which best practices your existing work already follows                                  | `/audit-applied-best-practices`           |
| Want to validate a practice profile before using or sharing it                                                | `/review-best-practice-profile`           |
| Want to publish or share a practice profile                                                                   | `/share-best-practice-profile`            |
| Want to create a new practice profile                                                                         | `/write-best-practice-profile`            |
| Authoring a new SKILL.md for grimoire-core or a STANDARD.md-following package                                 | `/write-best-practice-skill`              |
| Reviewing whether a SKILL.md meets STANDARD.md criteria                                                       | `/review-best-practice-skill`             |
| A SKILL.md has review findings, a stale citation, or an outdated referenced tool                              | `/revise-best-practice-skill`             |
| Assessing quality of all skills in a domain/sub-domain (pre-release, post-bulk-contribution)                  | `/audit-best-practice-domain`             |
| A skill's source is outdated or has been superseded — check if it should be deprecated                        | `/deprecate-best-practice-skill`          |
| Adding a new domain or sub-domain to grimoire                                                                 | `/design-best-practice-domain`            |

→ [BPDD guide](./bpdd.md) — cycle, compliance linter, LSP output, suppression, CI integration

### Common paths

**Describe any problem → `suggest-best-practice`**

Describe any situation. `suggest-best-practice` reads it, finds the best-matching skill(s), and routes you there.

<pre>
❯ My senior engineer just quit and the team is losing confidence.

grimoire ▸ finding a practice for team confidence after a senior departure — situation
    described, no approach yet

<b>Options</b>
1. run-employee-offboarding — handle the departure itself: knowledge transfer, how it's
                              announced, and the effect on team morale ★ recommended
2. apply-change-leadership — keep the team steady and give clear direction while things are
                             uncertain
3. plan-team-succession — find the critical knowledge the senior engineer held and who can
                          take it over
4. apply-trust-building — rebuild confidence in you and in the team's direction

Reply 1–4, "all", or "skip".
</pre>

**Multiple practices needed → `plan-best-practice-solution`**

When a problem needs 2+ practices coordinated in sequence — within one domain or across many — `plan-best-practice-solution` decomposes it with MECE methodology, sequences skills by dependency, and recursively plans complex sub-problems.

<pre>
❯ help me plan a migration from a monolith to services over the next year

grimoire ▸ planning a multi-step migration roadmap — several practices need sequencing

<b>Plan</b> — 7 steps
⎿  Step 1  design-greenfield-brownfield-strategy — confirm services are worth it and pick
                                                   the rollout approach
⎿  Step 2  service boundaries — choose one
     2a. design-microservices — split by domain area and find the first services to extract
                                ★ recommended
     2b. apply-domain-driven-design — model the business domain in depth before drawing
                                      boundaries
⎿  Step 3  incremental extraction — nested plan, 2 skills
     3a. apply-strangler-fig-pattern — move traffic from the monolith to new services a
                                       piece at a time
     3b. apply-branch-by-abstraction — swap out tightly coupled parts inside the monolith
⎿  Step 4  data decomposition — nested plan, 2 skills
     4a. apply-parallel-change — keep old and new schemas compatible during the split
     4b. write-database-migration — make each schema move safe and reversible
⎿  Step 5  design-event-driven-architecture — loosely coupled, asynchronous communication
                                              between services
⎿  Step 6  No match  "observability, contract testing, deployment platform" — general
    knowledge
⎿  Step 7  No match  "team ownership and the 12-month roadmap by quarter" — general
    knowledge
Apply this plan? Reply "yes", or pick another choice (e.g. "2b"). Steps that change files or
    run commands ask first.
</pre>

You confirm the plan once. grimoire then works through it, showing where it is (`— step 3a of 7`), and asks again only before a step that changes files or runs commands.

**Catching gaps before you start → `start-best-practice`**

Before any task begins, `start-best-practice` finds the practice that applies and applies it, so you start the work already aligned instead of discovering gaps halfway through. It runs automatically on every request once grimoire is installed. You can also call it explicitly with `grimoire:start-best-practice`. When one practice clearly fits, grimoire applies it straight away:

<pre>
❯ I'm about to write performance reviews for my team.

grimoire ▸ finding a practice for writing performance reviews — you've described a goal

<b>write-performance-review</b>
⎿  Approach  gather evidence, rate against the role, use specific examples, deliver before
    submitting

Before you write
1. Gather evidence from the whole period. Look at 1:1 notes, project outcomes, peer feedback
    and the work itself …
2. Rate each person against their role and level, not against each other …
</pre>

**Unclear problem → one question first**

When the problem itself is unclear, grimoire asks one question to narrow it down before choosing a practice, instead of guessing:

<pre>
❯ My product is losing users and I don't know why.

grimoire ▸ finding a practice for diagnosing user churn — problem described, no approach
                                                          chosen yet

You've told me the symptom (users are leaving) but not yet where it happens or what's behind
    it.
Fixes for churn depend a lot on when people leave, so I'll narrow that down first. …

Where are you seeing the drop: new signups who never stick, long-time users who stop coming
    back, or you can't tell yet?
</pre>

**Opaque complex problem → `apply-best-practice-tree`**

When the sub-problems only become clear as you solve them, `apply-best-practice-tree` applies a practice, looks at what it reveals, and chooses the next practice from there, so the structure of the solution emerges as you go.

**Internalizing a practice → `learn-best-practice`**

When you want to actually own a practice — not just apply it once — `learn-best-practice` runs active recall, spaced-repetition review, and scenario drilling so you can use it without AI in the future.

<pre>
❯ I want to actually learn MECE decomposition, not just use it.

grimoire ▸ teaching mece decomposition — you want to learn the practice, not just apply it

<b>apply-mece</b>
⎿  Mode  learn for durable recall: baseline, short card, practice scenarios, memory aid,
    review schedule

1. Baseline: you explain MECE in your own words so I know where to start.
2. Short card: when to use it, the core steps, and the most common mistake.
3. Practice: I give you a realistic problem, you decompose it without looking anything up,
    and I point out exactly where you have gaps or overlaps.
4. Memory aid: a 3–5 item checklist you can run in your head.
5. Review schedule: quick self-tests tomorrow, in 3 days and in a week, so it lasts.

In your own words, what is MECE decomposition for, and what makes a breakdown MECE?
</pre>
