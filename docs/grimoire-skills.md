# Meta Skills

Meta skills run the grimoire framework itself. Unlike domain skills (which encode best practices in engineering, health, finance, etc.), meta skills are how you find, apply, contribute, and maintain those domain skills.

There are two audiences: **users** who want to apply best practices, and **contributors** who manage the library.

---

## User-facing skills

These skills cover every stage of the user journey: defining the problem, finding the right practice, planning multi-domain work, evaluating existing solutions, and going deeper on any practice.

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

### analyze-problem

**When to use:** When the user's problem isn't clearly defined — goal is unclear, scope is unstated, or what's described looks like a symptom rather than a root cause.

**What it produces:** A clarification loop (one question at a time), then a structured problem space map (problem statement / scope / constraints / root cause) and 2–4 possible solution routes before handing off to solution skills.

**Example invocation:**
```
User: "My team isn't performing well."
→ analyze-problem asks: "What outcome are you trying to achieve — better output, faster delivery, or fewer conflicts?"
→ After answers: Problem statement, scope, routes A/B/C, handoff to suggest-best-practice
```

---

### discover-best-practices

**When to use:** When the user mentions a domain, role, or field without a specific problem — or wants to know what best practices exist before encountering a problem.

**What it produces:** A grouped list of available practices by subdomain, framed as "what goes wrong without each one," with 1–3 "most commonly discovered too late" highlights.

**Example invocation:**
```
User: "What best practices exist for contract negotiation?"
→ discover-best-practices lists practices by subdomain, surfaces the 2 most commonly skipped
→ Offers to apply any or learn more
```

---

### intercept-best-practice

**When to use:** When the user signals they are about to begin a task — "I'm going to write X", "starting to build Y", "planning to do Z" — before they have started.

**What it produces:** A proactive surface of the most relevant practice for that task, with the key gaps it would prevent, and an offer to apply it before they begin.

**Example invocation:**
```
User: "I'm about to refactor our auth module."
→ intercept-best-practice: "Before you start: review-security-posture applies here.
   Common gaps at this stage: token expiry strategy, rate limiting, audit logging.
   Apply now, or continue without? [y/n]"
```

---

### apply-best-practice-tree

**When to use:** When a problem is complex but stays within one domain — requiring recursive decomposition into sub-problems, each matched to the best installed skill.

**What it produces:** A recursive skill tree — the top-level problem matched to a skill, which surfaces sub-problems, each matched and applied in turn until the full problem is covered.

**Example invocation:**
```
User: "Our production incident response is a mess."
→ apply-best-practice-tree: top match design-incident-response-process →
  sub-problems: Detection (apply-alerting-strategy), Escalation (design-on-call-rotation),
  Communication (write-status-page-update), Post-incident (apply-five-whys)
→ Apply each sub-skill in order
```

---

### audit-applied-best-practices

**When to use:** When the user has existing work — code, a document, a design, a plan — and wants to know which best practices were applied, which are missing, and what gaps to address.

**What it produces:** An audit of the existing work against applicable practices (APPLIED / PARTIAL / MISSING), with a prioritized list of gaps to close.

**Example invocation:**
```
User: "Here's our current incident response runbook."
→ audit-applied-best-practices: identifies domain (engineering/reliability),
  scores runbook against applicable practices, surfaces MISSING: blameless framing,
  contributing factors section, PARTIAL: timeline structure
→ Fix list ordered by impact
```

---

### compare-best-practices

**When to use:** When the user needs to choose between two or more practices that could apply to the same problem — or when suggest-best-practice surfaces multiple close matches.

**What it produces:** A side-by-side table (approach / evidence / best for / trade-off / adoption cost) and a clear recommendation with reasoning.

**Example invocation:**
```
User: "Should I use Kanban or Scrum for my 4-person team?"
→ compare-best-practices: side-by-side table, then:
  Recommendation: Kanban
  Reason: 4-person teams with mixed work types benefit from flow-based prioritization;
  Scrum's sprint ceremonies add overhead that benefits larger teams more.
```

---

### explain-best-practice

**When to use:** When the user wants to understand WHY a practice works, not just apply it — or asks "what's the evidence for X?", "why do we do Y?".

**What it produces:** A structured explanation with six sections: the problem, the origin, the evidence, the mechanism, failure modes, and common misconceptions — plus a real-world failure example.

**Example invocation:**
```
User: "Why do code reviews matter? Can't we just rely on tests?"
→ explain-best-practice: covers problem (bugs tests miss), origin (Fagan 1976),
  evidence (Google 50% defect reduction), mechanism (second reader catches assumptions),
  failure modes (review fatigue, rubber-stamping), misconceptions (it's about finding bugs)
```

---

### adapt-best-practice

**When to use:** When the user wants to apply a practice but has constraints that differ from the canonical case — team size, regulation, budget, maturity, or tech stack.

**What it produces:** An adapted version of the practice with each step labeled ✓ unchanged / adapted: / ✗ skipped, plus a trade-off summary. Core steps are never silently removed.

**Example invocation:**
```
User: "Apply design-api-architecture, but we're a 2-person startup — we don't have time for full API governance."
→ adapt-best-practice: classifies each step Core/Adjustable/Optional,
  skips Optional governance steps, adapts versioning to lightweight URL convention,
  warns: "Rate limiting is a Core step — skipping it is a risk."
```

---

### teach-best-practice

**When to use:** When the user wants to explain or advocate for a practice with others — a team, manager, client, or stakeholder.

**What it produces:** Audience-tailored talking points, a written brief, or a slide outline — with a hook framed around what the audience cares about losing, anticipated objections, and concrete responses.

**Example invocation:**
```
User: "Help me convince my manager to adopt blameless post-mortems."
→ teach-best-practice: detects audience (non-technical manager, cares about risk/cost),
  produces: hook (without this, incidents repeat and cost 3× more), what it is (one sentence),
  why it works (incentivizes reporting over blame), evidence (Google SRE data),
  objections: "We don't have time" → response
```

---

### pin-best-practice-preference

**When to use:** When the user wants to save a practice preference for a domain or subdomain — so future sessions automatically apply it without re-matching.

**What it produces:** A saved preference at session level, project level (`<root>/.grimoire/preferences.md`), or global level (`~/.grimoire/preferences.md`), depending on user choice.

**Example invocation:**
```
User: "Always use Kanban for project management, not Scrum."
→ pin-best-practice-preference: saves preference at user-chosen scope,
  future suggest-best-practice calls skip matching for this subdomain
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
Browse before a problem
  → discover-best-practices
       → grouped list by subdomain with gap framing
       → highlights most commonly discovered too late
       → offer to apply or learn more

About to start a task
  → intercept-best-practice
       → proactive match before work begins
       → key gaps identified
       → apply now or continue

Problem isn't clear yet
  → analyze-problem
       → clarification loop (one question at a time)
       → problem statement + space map + possible routes
       → hands off to suggest-best-practice or plan-best-practice-solution

Describe situation
  → suggest-best-practice (or analyze-problem if vague)
       → single skill match: apply it directly
       → multiple matches: ranked list with ★ recommendation
       → no match: identify domain to install
       → multi-domain problem: hand off to plan-best-practice-solution

Complex single-domain problem
  → apply-best-practice-tree
       → recursive decomposition into sub-problems
       → each sub-problem matched and applied in turn

Existing solution to evaluate
  → review-best-practice-fit
       → ALIGNED / PARTIAL / MISSING per practice
       → prioritized fix list

Audit what was applied to existing work
  → audit-applied-best-practices
       → identifies which practices were applied, which are missing
       → fix list ordered by impact

Choosing between practices
  → compare-best-practices
       → side-by-side table: approach / evidence / best for / trade-offs
       → clear recommendation with reasoning

Deep understanding of a practice
  → explain-best-practice
       → problem → origin → evidence → mechanism → failure modes → misconceptions

Applying under constraints
  → adapt-best-practice
       → Core / Adjustable / Optional classification per step
       → adapted version with labeled changes

Sharing with others
  → teach-best-practice
       → audience analysis → tailored hook → talking points / brief / slide outline

Save a preference
  → pin-best-practice-preference
       → session / project / global scope
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
| `analyze-problem` | User | problem isn't clearly defined | problem statement + space map + possible routes |
| `discover-best-practices` | User | user mentions a domain, no specific problem | grouped practice list with gap framing |
| `intercept-best-practice` | User | about to begin a task | proactive practice match + apply offer |
| `suggest-best-practice` | User | any situation or question | matching skill, ranked list, or install recommendation |
| `plan-best-practice-solution` | User | multi-domain or complex problem | sequenced skill application plan |
| `apply-best-practice-tree` | User | complex single-domain problem | recursive skill tree applied in sequence |
| `review-best-practice-fit` | User | existing solution to evaluate | ALIGNED/PARTIAL/MISSING per practice + fix list |
| `audit-applied-best-practices` | User | existing work to audit | APPLIED/PARTIAL/MISSING per practice + fix list |
| `compare-best-practices` | User | multiple practices could apply | side-by-side table + recommendation |
| `explain-best-practice` | User | wants to understand WHY a practice works | 6-section explanation + failure example |
| `adapt-best-practice` | User | practice has constraint conflicts | adapted version with labeled step changes |
| `teach-best-practice` | User | wants to share practice with others | audience-tailored talking points/brief/outline |
| `pin-best-practice-preference` | User | wants to save a practice preference | saved preference at session/project/global scope |
| `write-best-practice-skill` | Contributor | new practice to encode | complete SKILL.md |
| `review-best-practice-skill` | Contributor | SKILL.md to evaluate | PASS/NEEDS-REVISION/REJECT verdict with findings |
| `revise-best-practice-skill` | Contributor | review findings to address | targeted fixes to existing SKILL.md |
| `audit-best-practice-domain` | Contributor | domain to assess | per-skill verdicts + summary counts |
| `deprecate-best-practice-skill` | Contributor | outdated skill | deprecation notice + removal PR |
| `design-best-practice-domain` | Contributor | new domain concept | directory scaffold + plugin files + seed skills |

For contributor workflow detail, see [authoring-skills.md](./authoring-skills.md) and [maintaining.md](./maintaining.md).
