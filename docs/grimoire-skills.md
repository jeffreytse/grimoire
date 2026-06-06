# Meta Skills

Meta skills run the grimoire framework itself. Unlike domain skills (which encode best practices in engineering, health, finance, etc.), meta skills are how you find, apply, contribute, and maintain those domain skills.

There are two audiences: **users** who want to apply best practices, and **contributors** who manage the library.

---

## User-facing skills

These three skills cover every user interaction: finding a practice, planning multi-domain work, and evaluating an existing solution.

### suggest-best-practice

**When to use:** Any time a user describes a situation, problem, goal, or question — including when they don't know which domain applies or don't know a best practice exists for their situation.

**What it produces:** Either applies the best-matching skill directly, or presents a ranked list of candidates for the user to choose from. If no installed skill matches, it identifies which domain to install.

**Example invocation:**
```
User: "Our API response times have been spiking for two days and we can't figure out why."
→ suggest-best-practice classifies this as engineering/debugging, routes to the appropriate skill
```

```
User: "I need to write a performance review for a direct report."
→ suggest-best-practice: no installed skill matches — suggests installing a management domain
```

---

### plan-best-practice-solution

**When to use:** When the problem spans multiple domains, requires coordinating several best practices in sequence, or is too complex for a single skill — e.g., launching a startup, handling a compliance audit, building a hiring process from scratch.

**What it produces:** A sequenced plan — sub-problems identified, each matched to the best skill, ordered by dependency. Executes one skill at a time through the plan.

**Example invocation:**
```
User: "We're a 12-person startup about to raise Series A. What do we need to get right?"
→ plan-best-practice-solution decomposes: legal entity structure → cap table clean-up →
  unit economics → pitch narrative → data room — each mapped to a skill, sequenced
```

---

### review-best-practice-fit

**When to use:** When the user already has a solution, plan, design, or approach and wants to know how well it aligns with best practices — what's missing, what's wrong, what to fix first.

**What it produces:** ALIGNED / PARTIAL / MISSING verdict per applicable practice, with a prioritized fix list.

**Example invocation:**
```
User: "Here's our API design doc — does it follow REST best practices?"
→ review-best-practice-fit evaluates against applicable skills:
  ALIGNED: resource naming, HTTP verbs
  PARTIAL: error response format (missing error codes)
  MISSING: pagination contract, versioning strategy
  → fix list ordered by impact
```

---

## Contributor skills

These six skills manage the library: authoring, reviewing, fixing, auditing, retiring, and designing new domains.

### write-best-practice-skill

**When to use:** When authoring a new SKILL.md — whether starting from scratch, adapting existing knowledge, or encoding a domain best practice.

**What it produces:** A complete SKILL.md with all required sections: frontmatter, Why This Is Best Practice (Adopted by / Impact / Why best / Sources), and Steps. Structured to pass `review-best-practice-skill`.

**Example invocation:**
```
User: "I want to add a skill for conducting blameless post-mortems."
→ write-best-practice-skill guides: qualify the practice → name it → write frontmatter →
  find credible sources → write Why section → write Steps → self-review
```

---

### review-best-practice-skill

**When to use:** When evaluating whether a SKILL.md meets grimoire standards — for self-review before submitting a PR, for maintainer PR review, or for auditing existing skills.

**What it produces:** A structured verdict — PASS, NEEDS-REVISION, or REJECT — with a per-criterion findings table. REJECT blocks merge. NEEDS-REVISION lists exactly what to fix.

**Example invocation:**
```
Applying review-best-practice-skill to skills/engineering/development/skills/review-pull-request/SKILL.md
→ verdict table: name ✅, description ✅, tags ⚠️ (missing outcome axis), Impact ⚠️ (vague)
→ Overall: NEEDS-REVISION — 2 required fixes before merge
```

---

### revise-best-practice-skill

**When to use:** When an existing SKILL.md has `review-best-practice-skill` findings to address, a citation has become inaccurate, steps reference an outdated tool, or scope needs adjusting.

**What it produces:** Targeted changes to the specific failing sections only. Passing sections are left untouched.

**Example invocation:**
```
review-best-practice-skill found: Impact line says "significantly improves" — no number or study
→ revise-best-practice-skill: load the Impact finding → find the specific sentence →
  replace with causal argument or cited number → re-run review-best-practice-skill → PASS
```

---

### audit-best-practice-domain

**When to use:** When assessing the quality of all skills in a domain — before a release, after a bulk contribution, when adopting a domain for the first time, or on a weekly maintenance schedule.

**What it produces:** Per-skill verdicts (PASS/NEEDS-REVISION/REJECT) for every SKILL.md in the domain, plus summary counts. REJECT findings trigger deprecation or revision.

**Example invocation:**
```
Applying audit-best-practice-domain to skills/engineering/
→ engineering/development: 8 PASS, 1 NEEDS-REVISION (propose-conventional-commit)
→ engineering/testing: 5 PASS, 0 findings
→ Action items: open issue for propose-conventional-commit revision
```

---

### deprecate-best-practice-skill

**When to use:** When a skill has become outdated — the referenced tool no longer exists at scale, the source institution revised its position, a newer practice has achieved majority top-tier adoption that supersedes it, or the skill fails criteria it once passed.

**What it produces:** A two-step PR process: (1) deprecation notice added to the skill file + `deprecated: true` in marketplace.json; (2) after one release cycle, a removal PR deleting the directory.

**Example invocation:**
```
User: "The skill references a tool that was sunset two years ago."
→ deprecate-best-practice-skill: add notice → point to replacement skill →
  merge deprecation PR → wait one release cycle → open removal PR
```

---

### design-best-practice-domain

**When to use:** When adding a new domain or sub-domain to grimoire — whether starting a brand-new domain (health, law, cooking), adding a new sub-domain to an existing domain, or deciding whether a new sub-domain is needed at all.

**What it produces:** A complete domain scaffold: directory structure, domain plugin.json, sub-domain plugin.json files, marketplace.json entries, README table update, and minimum 2 seed skills — all verified by `audit-best-practice-domain` before PR.

**Example invocation:**
```
User: "I want to add a 'leadership' domain with skills for 1:1s, performance reviews, team structure."
→ design-best-practice-domain: qualify (3+ separable subdomains?) → choose flat vs hierarchical →
  create dirs → write plugin.json → write 2 seed skills → audit-best-practice-domain → open PR
```

---

## Lifecycle: how meta skills connect

The meta skills form a complete workflow. No step is orphaned.

### Finding and applying practices (users)

```
Describe situation
  → suggest-best-practice
       → single skill match: apply it directly
       → multiple matches: ranked list for user to choose
       → no match: identify domain to install
       → multi-domain problem: hand off to plan-best-practice-solution

Existing solution to evaluate
  → review-best-practice-fit
       → ALIGNED / PARTIAL / MISSING per practice
       → prioritized fix list
```

### Contributing a new skill

```
suggest-best-practice            ← check no duplicate exists
  └→ design-best-practice-domain          ← if new domain needed
       └→ write-best-practice-skill       ← author SKILL.md
            └→ review-best-practice-skill ← self-review
                 └→ revise-best-practice-skill   ← fix NEEDS-REVISION findings
                      └→ review-best-practice-skill   ← re-verify → PASS
                           └→ open PR
```

### Maintaining the library

```
audit-best-practice-domain                ← weekly or pre-release
  └→ NEEDS-REVISION found → revise-best-practice-skill → re-run review-best-practice-skill
  └→ REJECT found
       → outdated practice: deprecate-best-practice-skill
       → fixable: revise-best-practice-skill
       → replaced by new domain: design-best-practice-domain
```

---

## Quick reference

| Skill | Audience | Trigger | Output |
|-------|----------|---------|--------|
| `suggest-best-practice` | User | any situation or question | matching skill, ranked list, or install recommendation |
| `plan-best-practice-solution` | User | multi-domain or complex problem | sequenced skill application plan |
| `review-best-practice-fit` | User | existing solution to evaluate | ALIGNED/PARTIAL/MISSING per practice + fix list |
| `write-best-practice-skill` | Contributor | new practice to encode | complete SKILL.md |
| `review-best-practice-skill` | Contributor | SKILL.md to evaluate | PASS/NEEDS-REVISION/REJECT verdict with findings |
| `revise-best-practice-skill` | Contributor | review findings to address | targeted fixes to existing SKILL.md |
| `audit-best-practice-domain` | Contributor | domain to assess | per-skill verdicts + summary counts |
| `deprecate-best-practice-skill` | Contributor | outdated skill | deprecation notice + removal PR |
| `design-best-practice-domain` | Contributor | new domain concept | directory scaffold + plugin files + seed skills |

For contributor workflow detail, see [authoring-skills.md](./authoring-skills.md) and [maintaining.md](./maintaining.md).
