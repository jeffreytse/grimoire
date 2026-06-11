---
name: apply-small-cls
description: Use when preparing a pull request or planning a feature branch — to keep changes small, focused, and independently reviewable so reviews are faster, more thorough, and easier to roll back.
source: 'Google Engineering Practices "Small CLs" (google.github.io/eng-practices/review/developer/small-cls.html); Google Engineering Practices "Separate Out Refactorings"'
tags: [code-review, pull-requests, git, collaboration, developer-experience, ci-cd]
verified: true
---

# Keep CLs Small

Submit changes that do exactly one thing and are as small as possible — small CLs are reviewed faster, reviewed more thoroughly, less likely to introduce bugs, easier to roll back, and cause fewer merge conflicts.

## Why This Is Best Practice

**Adopted by:** Google's eng-practices guide mandates small CLs (changelists = pull requests) as a core discipline. The principle is independently described at Stripe, Shopify, and most high-velocity engineering teams. Trunk-based development (adopted by Google, Meta, Netflix) depends on small, frequent merges to keep the main branch stable.

**Impact:** Large CLs are reviewed slowly and superficially. Reviewers faced with 800-line diffs switch to approval mode — they check that the overall approach makes sense but miss individual bugs and design issues. Small CLs (under 200 lines of substantive change) receive thorough line-by-line review. A Google internal study found that CLs over 400 lines receive significantly less thorough review than CLs under 200 lines.

**Why best:** A small CL does exactly one thing. "One thing" means the reviewer needs to hold only one context in their head, the change can be reverted without reverting unrelated work, and the diff in `git log` communicates a clear intent. Large CLs conflate multiple concerns, making review, revert, and bisection all harder.

Sources: Google Engineering Practices; Winters, Manshreck & Wright *Software Engineering at Google* (2020); Trunk-Based Development (trunkbaseddevelopment.com)

## Steps

### Step 1: Define "small" — and aim for it

A CL is small when:
- It does **one self-contained thing** — one feature, one bug fix, one refactor, one cleanup
- A reviewer can review it thoroughly in **under 20 minutes**
- It is typically **under 200 lines** of substantive change (excluding generated code, test fixtures, and mechanical renames)

Lines of code is a proxy, not the definition. A 400-line CL that adds a single well-defined feature with tests is better than three 50-line CLs tangled together. "Self-contained and focused" is the rule; line count is the signal.

### Step 2: Separate refactoring from feature work

Never combine a refactor with a new feature or bug fix in the same CL. These are independent concerns and combining them makes review much harder.

```
# BAD — one CL does two things
CL: "Refactor PaymentService and add discount code support"
- Renames PaymentService → BillingService across 40 files
- Adds discount_code field to charge() method

# GOOD — two separate CLs in sequence
CL 1: "Rename PaymentService to BillingService" (refactor only, no behavior change)
CL 2: "Add discount code support to BillingService" (feature only)
```

Even small cleanups (renaming a local variable, removing a dead import) should be in a separate CL if they'd obscure the diff of the actual change.

### Step 3: Split large features horizontally (by layer) or vertically (by sub-feature)

When a feature is inherently large, split it into a sequence of independent CLs.

**Horizontal split (by layer):**
```
CL 1: Add data model + migration  (no business logic yet)
CL 2: Add repository layer        (data access, depends on CL 1)
CL 3: Add service layer           (business logic, depends on CL 2)
CL 4: Add API endpoint            (controller, depends on CL 3)
CL 5: Add frontend component      (depends on CL 4)
```
Each CL is reviewable on its own. Reviewers see one layer at a time.

**Vertical split (by sub-feature):**
```
CL 1: Add multiplication operator  (full stack: model + service + API + UI)
CL 2: Add division operator        (independent, parallel track)
```
Each CL delivers a complete sub-feature and can be reviewed independently.

### Step 4: Keep tests in the same CL as the code they test

Tests belong in the same CL as the production code they cover. Separate "add feature" and "add tests for feature" CLs leave the codebase in a state where the feature exists but is untested — that defeats the purpose of tests as a verification gate.

Exception: adding tests for *existing untested code* is a valid standalone CL — it improves code health without changing behavior.

### Step 5: Write a CL description that makes the scope clear

A focused CL is easy to describe in one sentence. If you can't write a clear first-line description, the CL is probably doing too many things.

```
# CL that's too broad — description is vague or compound
"Various payment fixes and refactoring"
"Clean up user module and add email verification"

# CL that's well-scoped — description is specific
"Fix race condition in PaymentService when two charges submit simultaneously"
"Add email verification step to user registration flow"
```

### Step 6: When large CLs are acceptable

Two situations where a large CL is justified:

1. **File deletion:** Deleting an entire file counts as one logical change regardless of line count — reviewers don't need to read every deleted line.

2. **Automated mechanical changes:** A CL generated by a refactoring tool (e.g., rename across the codebase, code format migration) can be large if the reviewer only needs to verify the intent, not inspect each individual change. Include the tool command in the description so reviewers know what to verify.

```
# Good description for a large mechanical CL
"Run gofmt across the entire repo (no behavior changes)
 Command: gofmt -w ./..."
```

## Rules

- One CL = one self-contained thing. If the description requires "and", it's two CLs.
- Refactoring and feature work are always separate CLs.
- Tests belong in the same CL as the code they cover.
- If a CL will take a reviewer more than 20 minutes, find a way to split it.
- Large CLs are acceptable only for file deletions and verified automated changes.

## Common Mistakes

**"I'll just put it all in one PR to save time":** Large PRs don't save time — they shift the cost to the reviewer, who reviews them slowly, misses bugs, and resents the process. Small PRs get reviewed faster and more thoroughly; the total elapsed time to merge is usually shorter.

**Combining refactor with feature:** The refactor changes lines 1-400, the feature changes lines 401-450. The reviewer can't tell which lines are behavior changes and which are pure renames. Split them; the feature-only CL review takes 10 minutes instead of 45.

**Separating tests from code:** "I'll add the tests in the next PR" means the tests often never arrive, and the code exists in a verified-untested state in production until they do.

**Treating "large" as a line count rather than a focus count:** 500 lines of migration SQL that does one thing is better than 3 CLs of 50 lines each that are logically entangled. The metric is focus, not size.
