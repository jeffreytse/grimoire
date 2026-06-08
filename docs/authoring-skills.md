# Authoring Skills

This guide walks through the full lifecycle of writing a grimoire skill — from deciding whether a practice qualifies to opening the PR. Each step maps to a meta skill you can invoke.

For the formal quality checklist, see [STANDARD.md](../STANDARD.md). For the skill template, see [SKILL_TEMPLATE.md](../SKILL_TEMPLATE.md).

---

## The contribution lifecycle

```
suggest-best-practice          ← check no duplicate exists
  └→ design-best-practice-domain        ← if adding a new domain first
       └→ write-best-practice-skill     ← author the SKILL.md
            └→ review-best-practice-skill    ← self-review before PR
                 └→ revise-best-practice-skill   ← fix findings
                      └→ review-best-practice-skill   ← re-verify
                           └→ open PR
```

Invoke each meta skill to get step-by-step guidance for that stage.

---

## Step 1: Qualify the practice

Before writing a single line, confirm all four gates. A skill that fails any gate wastes reviewer time.

**Gate 1 — Majority adoption**
The practice is used by MOST top-tier companies or credentialed professionals in the domain — not a few notable ones, not your team's approach. "Most" means you could name 5+ organizations without straining.

```
❌ "We use this at our company"
❌ "A few top firms use this"
✅ "Google, Microsoft, Meta, Stripe — virtually all software companies >50 engineers"
✅ "Codified in NSCA CSCS exam content, used by all certified strength coaches"
```

**Gate 2 — Measurable impact**
There must be a number or a named study. "Helps" or "improves" without data fails this gate.

```
❌ "Code review improves code quality"
✅ "Code review reduces post-ship defects by ~50% (Google internal data)"
✅ "Progressive overload produces 2–3× greater strength gains vs. fixed-load training [Ralston et al., 2017, SR]"
```

**Gate 3 — No duplicate**
Run `suggest-best-practice` with your topic first. If an existing skill already covers it, extend that skill via [revise-best-practice-skill](../skills/meta/skills/revise-best-practice-skill/SKILL.md) instead of creating a duplicate.

**Gate 4 — No genuine controversy**
If credible top-tier professionals are split ~50/50, there is no consensus and therefore no best practice. Encode the majority position; if genuinely split, skip it.

---

## Step 2: Check for existing domain

If your skill belongs under an existing domain and subdomain, skip to Step 3.

If you need a new domain or subdomain, invoke `design-best-practice-domain` first. It covers:
- Qualifying the domain (3+ separable subdomains required for a top-level domain)
- Directory structure and plugin.json templates
- Marketplace.json entries
- Minimum seed skills before PR

---

## Step 3: Name the skill

Pattern: `<verb>-<subject>[-<qualifier>]`

**Approved verbs:**

| Verb | Use for |
|------|---------|
| `propose-` | Draft an artifact for human approval |
| `write-` | Author a document or content |
| `review-` | Evaluate quality of something |
| `audit-` | Batch evaluation across many items |
| `design-` | Architect a system, plan, or program |
| `calculate-` | Compute a numeric value |
| `diagnose-` | Identify root cause of a problem |
| `optimize-` | Improve a measured metric |
| `suggest-` | Recommend options for selection |
| `deprecate-` | Retire an outdated artifact |
| `plan-` | Create a structured sequence of actions |
| `negotiate-` | Handle a back-and-forth agreement |

Verbs not in this table are allowed if none of the above fit. These are always rejected: `do-`, `handle-`, `manage-`, `improve-`, `get-`, `use-`, `help-`.

**Subject specificity test:** ask "which kind?" — if it's a valid follow-up question, the subject is too generic.

```
❌ review-contract      → which kind of contract?
✅ review-saas-contract
✅ review-employment-contract

❌ design-program       → which kind of program?
✅ design-training-program
✅ design-onboarding-program
```

**Format rules:**
- Kebab-case only
- 2–4 words (verb + 1–2 word subject + optional qualifier)
- ≤50 characters
- No `skill-`, `guide-`, `best-practice-` prefix — redundant in a skill library
- No noun-first: ~~`contract-review`~~, ~~`training-program-design`~~

---

## Step 4: Write the frontmatter

```yaml
---
name: verb-first-kebab-case
description: Use when <triggering conditions — NOT a summary of the steps>.
source: <institution, standard body, or named top-tier companies>
tags: [problem-keyword, tool-method, role-context, outcome]
---
```

**`description`** — the most common mistake here is summarizing what the skill does. It must describe WHEN to use it, not what it does.

```yaml
# Wrong — summarizes content
description: Use when committing — inspects staged files, drafts conventional commit message, presents for approval.

# Correct — triggering conditions only
description: Use when the user asks to commit, wants a commit message, or invokes /propose-commit.
```

**`tags`** — must cover all four axes. These are the signals `suggest-best-practice` uses to route problems.

```yaml
# Bad — only problem keywords, missing role and outcome
tags: [code-quality, pull-request, review]

# Good — all four axes
tags: [code-quality, pull-request, developer, defect-reduction]
#       ^ problem    ^ tool/method  ^ role     ^ outcome
```

---

## Step 5: Write "Why This Is Best Practice"

This is the most scrutinized section. It must prove the skill belongs in grimoire.

```markdown
## Why This Is Best Practice

**Adopted by:** [specific companies or institutions — not "many companies"]
**Impact:** [measurable outcome with cited evidence — not "significantly better"]
**Why best:** [why this over the named alternative — name the alternative]

Sources: [verifiable citation]
```

**Adopted by** — name at least two specific organizations or one named institution. Claim majority adoption, not just notable examples.

**Impact** — the most common failure mode. Every rejected skill has a vague Impact line. Rules:
- Must contain a number (%, ratio, time unit) OR a named study
- The outcome must be causal, not just correlated
- If you're not certain a number is in a specific document, remove the number and keep the causal argument instead

```
❌ "Significantly improves code quality"
❌ "Studies show 60% improvement" (which studies?)
✅ "Reduces post-ship defects by ~50% (Google internal data, Engineering Practices documentation)"
✅ "Teams that skip problem definition produce solutions targeting the wrong problem, requiring full redesign (IDEO Human-Centered Design Field Guide, 2015)"
```

**Why best** — must name at least one alternative and explain why this approach wins over it.

```
❌ "This is the best way to handle this situation"
✅ "Checklist-based review — vs. ad-hoc review where criteria drift reviewer to reviewer — produces consistent verdicts every time. Ad-hoc review finds surface problems only."
```

---

## Step 6: Write Steps

Steps must be immediately executable — the reader can follow them right now.

**Concrete over abstract:**
```
❌ "Review your changes before committing"
✅ "Run `git diff --cached` to see all staged changes"
```

**Complete — include failure modes:**
```
❌ "If the build passes, deploy"
✅ "If the build passes, deploy. If it fails, run `npm run build -- --verbose` to see
   which module failed. If the error is in a test dependency, check that NODE_ENV=test
   is not set in the deployment environment."
```

**Scoped — if you're wondering whether to split, split:**
```
❌ "Optimize database queries and add caching" (two concepts)
✅ optimize-query-latency  +  design-cache-strategy  (two skills)
```

**Size:** 50–300 lines total. Under 50 lines = too shallow, add more concrete steps or examples. Over 300 lines = split into two skills.

---

## Step 7: Add domain safety footer (if applicable)

| Domain | Required footer |
|--------|----------------|
| Health / Medicine | `> For personal health decisions, consult a qualified healthcare provider.` |
| Law | `> This is educational content, not legal advice. Consult qualified legal counsel for your situation.` |
| Finance / Investing | `> This is educational content, not financial advice. Consult a licensed financial advisor for personal decisions.` |
| Psychology / Mental Health | `> For mental health concerns, consult a qualified mental health professional.` |

---

## Step 8: Self-review with review-best-practice-skill

Before opening a PR, invoke `review-best-practice-skill` on your file:

```
Applying review-best-practice-skill to skills/<domain>/<subdomain>/skills/<name>/SKILL.md
```

`review-best-practice-skill` runs the full [STANDARD.md](../STANDARD.md) checklist and produces a verdict:

- **PASS** — ready to open a PR
- **NEEDS-REVISION** — specific fixes required before merge
- **REJECT** — a required field is missing or a core criterion fails; fix and re-run

Fix every NEEDS-REVISION and REJECT finding. Use `revise-best-practice-skill` to apply changes to specific sections without touching passing content. Re-run `review-best-practice-skill` after each round of fixes.

A skill that fails `review-best-practice-skill` will be rejected in the PR review. Run it yourself first.

---

## Step 9: Open a PR

**PR title format:**
```
feat(<domain>/<subdomain>): add <skill-name>
```

Example: `feat(engineering/development): add review-pull-request`

**PR description:** copy the self-check checklist from [CONTRIBUTING.md](../CONTRIBUTING.md) and mark each item.

**One skill per PR** — easier to review, faster to merge. If you're adding multiple skills, open separate PRs.

If your skill introduces a new subdomain, follow the sub-domain setup steps in [CONTRIBUTING.md](../CONTRIBUTING.md) and include those changes in the same PR.

---

## Common failure modes

| Failure | Fix |
|---------|-----|
| "Adopted by" says "many companies" | Name Google, Netflix, Mayo Clinic, CFA Institute — specific |
| "Impact" says "significantly improves" | Find the study or the number; if you can't, keep the causal argument without a number |
| Steps describe concepts instead of actions | Rewrite each step as a command the reader can execute immediately |
| Skill covers two separable concepts | Split into two skills; run `write-best-practice-skill` for each |
| `description` summarizes what the skill does | Replace with triggering conditions only — starts with "Use when" |
| Tags don't cover all four axes | Add tags until problem / tool / role / outcome are all represented |
