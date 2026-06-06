# Tag Vocabulary

Tags are how `suggest-best-practice` routes any situation to the right skill. Every skill requires tags covering all four axes. Inconsistent tags silently break routing — two skills covering the same role tagged differently won't both surface when they should.

This document defines accepted values per axis and explains how to choose.

---

## The four axes

| Axis | Question it answers | Example |
|------|-------------------|---------|
| **Problem** | What problem or situation is this skill for? | `commit-messages`, `deprecation`, `gap-analysis` |
| **Method** | What technique or tool does it apply? | `checklist`, `mece`, `semantic-versioning` |
| **Role** | Who uses this skill? | `developer`, `maintainer`, `reviewer` |
| **Outcome** | What does the user get? | `defect-prevention`, `consistency`, `domain-release` |

All four axes must be covered. A skill with only problem and method tags won't route to the right audience (missing role) and won't match queries framed as outcomes (missing outcome).

The total tag count should be **3–8**. Under 3 usually means an axis is missing. Over 8 usually means the skill covers too many concepts — consider splitting.

---

## Problem keywords

Use these to tag what the skill addresses. Pick the most specific match; add a broader category if needed.

### Commits and version control
| Tag | Meaning |
|-----|---------|
| `commit-messages` | Writing or evaluating commit message content |
| `conventional-commits` | Conventional Commits specification specifically |
| `changelog` | Maintaining or generating changelogs |
| `semantic-versioning` | Version numbering decisions |
| `git` | Git workflows and operations |

### Skill library management
| Tag | Meaning |
|-----|---------|
| `skill-authoring` | Writing a new SKILL.md |
| `skill-review` | Evaluating a SKILL.md |
| `skill-revision` | Fixing findings in an existing SKILL.md |
| `skill-discovery` | Finding the right skill for a situation |
| `skill-maintenance` | Keeping skills accurate and current |
| `skill-quality` | Assessing overall quality of skills |
| `deprecation` | Retiring an outdated skill or practice |
| `staleness` | Detecting skills or references that are out of date |
| `knowledge-freshness` | Ensuring knowledge stays current over time |
| `outdated-practice` | Practices that have been superseded |

### Domain and library structure
| Tag | Meaning |
|-----|---------|
| `domain-creation` | Adding a new domain to the library |
| `domain-health` | Quality state of an entire domain |
| `knowledge-architecture` | How knowledge is structured and organized |
| `knowledge-governance` | Rules and processes for managing a knowledge base |

### Practice evaluation and gap analysis
| Tag | Meaning |
|-----|---------|
| `gap-analysis` | Identifying what's missing from an approach |
| `practice-alignment` | How well an approach matches best practices |
| `solution-review` | Evaluating an existing solution or design |
| `quality-audit` | Batch assessment of quality across many items |
| `quality-debt` | Accumulated quality problems that need addressing |
| `quality-gate` | A checkpoint that blocks progress until quality criteria pass |

### Problem structure
| Tag | Meaning |
|-----|---------|
| `multi-domain` | Problems spanning more than one domain |
| `cross-domain` | Skills that apply across multiple domains |
| `problem-decomposition` | Breaking complex problems into sub-problems |
| `naming-inconsistency` | Inconsistent naming in code, docs, or taxonomy |
| `structured-authoring` | Writing structured documents like specs or proposals |

---

## Method tags

Use these to tag the technique or approach the skill applies.

| Tag | Meaning |
|-----|---------|
| `auto-classify` | Automatic categorization of input to the right domain |
| `checklist` | Evaluation via a fixed list of criteria |
| `mece` | Mutually Exclusive, Collectively Exhaustive decomposition |
| `problem-routing` | Directing a problem to the right skill or domain |
| `situation-analysis` | Analyzing a user's context before recommending |
| `skill-orchestration` | Coordinating multiple skills in sequence |
| `structured-solution` | Producing output in a defined structured format |

---

## Role tags

Use the role that best describes who invokes this skill. Pick the primary audience; add a second if genuinely two distinct audiences exist.

| Tag | Meaning |
|-----|---------|
| `developer` | Software engineer writing or reviewing code |
| `architect` | Technical decision-maker designing systems |
| `reviewer` | Person evaluating someone else's work (code, skills, plans) |
| `maintainer` | Person responsible for maintaining a library or codebase |
| `contributor` | Person contributing to an open-source project |
| `domain-expert` | Subject matter expert in a specific field |
| `practitioner` | Generic professional practitioner in any domain |
| `decision-maker` | Person making a strategic or resource decision |
| `problem-solver` | General — anyone trying to solve a problem across domains |

**Rule:** Don't use `practitioner` when a more specific role applies. `practitioner` is a fallback for domain skills where the audience can't be narrowed further (e.g., a health skill for anyone in medicine). For meta skills and engineering skills, use a specific role.

---

## Outcome tags

Use these to tag what the user achieves by applying the skill.

| Tag | Meaning |
|-----|---------|
| `consistency` | Uniform results across multiple applications |
| `contribution-quality` | Higher quality of contributions to a shared library |
| `defect-prevention` | Fewer bugs or errors in output |
| `documentation-standards` | Documentation that meets a defined quality bar |
| `domain-release` | A new domain ready to publish |
| `knowledge-governance` | A governed, maintained knowledge base |
| `onboarding-speed` | Faster time-to-productivity for new contributors |
| `practice-compliance` | Output that conforms to established best practices |
| `practice-recommendation` | A concrete recommendation for what practice to apply |
| `repo-health` | A codebase or skill library in good structural health |
| `review-pass` | A skill or artifact that passes a quality review |
| `searchability` | Content that can be found and routed accurately |
| `solution-improvement` | An existing solution improved against best practices |

---

## Choosing tags: decision rules

**Step 1 — Problem axis**
Ask: "What situation would cause someone to use this skill?" Pick 1–2 tags that match that situation.

**Step 2 — Method axis**
Ask: "What technique does this skill apply?" Pick 1 tag. If none fit, the method may be implicit in the problem — that's acceptable, but check the vocabulary first.

**Step 3 — Role axis**
Ask: "Who is this skill for — by job or responsibility?" Pick 1 tag. If two genuinely distinct audiences exist (e.g., both `developer` and `reviewer`), use both.

**Step 4 — Outcome axis**
Ask: "What does the user walk away with?" Pick 1–2 outcome tags. This is the most commonly missing axis.

**Step 5 — Count check**
If total tags < 3: an axis is missing. If total tags > 8: the skill covers too many concepts — split it.

---

## Common mistakes

**Missing the outcome axis**
```yaml
# Bad — no outcome
tags: [skill-review, checklist, reviewer]

# Good — outcome added
tags: [skill-review, checklist, reviewer, review-pass]
```

**Using a role too broad**
```yaml
# Bad — practitioner is too generic for an engineering skill
tags: [commit-messages, conventional-commits, practitioner, consistency]

# Good
tags: [commit-messages, conventional-commits, developer, consistency]
```

**Tags that describe the same thing twice**
```yaml
# Bad — skill-review and quality-gate both say "evaluation checkpoint"
tags: [skill-review, quality-gate, reviewer, review-pass]

# Good — pick the most specific
tags: [skill-review, checklist, reviewer, review-pass]
```

**Inventing tags with no equivalent in this vocabulary**
Before adding a new tag, check whether an existing tag covers it. Tags outside the vocabulary don't improve routing — they're noise. If a genuine gap exists, add the tag to this document in the same PR as the skill.

---

## Adding new tags

If no existing tag fits your skill, add one — but follow these rules:

1. **Check the vocabulary first.** A synonym of an existing tag does not qualify.
2. **Place it in the correct axis section.** If you're unsure which axis it belongs to, the tag may be combining two concepts — split it.
3. **Write a one-line definition.** Every tag in this document has a definition. Add yours.
4. **Include it in the same PR as the skill.** Don't add tags in isolation; they need a skill to validate their meaning.
