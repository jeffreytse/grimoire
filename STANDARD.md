# The Grimoire Skill Standard

**Version:** 1.0
**Status:** Active
**Authors:** grimoire contributors
**Scope:** Defines the quality requirements for AI agent skills across all domains.
**License:** [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/) — freely adoptable by any AI agent skill project.

This standard is maintained by the [grimoire project](https://github.com/jeffreytse/grimoire). Other AI agent skill libraries are invited to adopt it. See [Adopting this standard](#adopting-this-standard) below.

---

## Purpose

A grimoire skill exists for one reason: to promote a practice people wouldn't otherwise know to apply, and to guide them through applying it correctly.

Most people solving a problem don't know which expert standard governs it. LLMs that know those standards won't enforce them without explicit guidance. A skill closes both gaps: it tells the AI **when** to apply a practice (trigger condition in `description`) and **how** to apply it (structured steps from a verified source).

Skills are not documentation. They are instructions that change what an AI does.

---

## Changelog

| Version | Change |
|---------|--------|
| 1.0 | Initial release — naming standard, 4-axis tags, skill lifecycle model, machine-readable schema (`schema/skill.schema.json`), conformance test suite (`schema/tests/`) |

A grimoire skill encodes a **best practice**: a technique with strong, demonstrated
impact that is adopted by most top-tier companies or leading professionals in the domain.

Not expert opinion. Not an interesting technique. Not one company's approach.
A practice proven at scale, used by the best.

The reference implementation is `skills/engineering/development/skills/propose-conventional-commit/SKILL.md`.

## Framework Quality Basis

This framework is an implementation of the `design-contribution-standard` methodology
(see `skills/engineering/documentation/skills/design-contribution-standard/`), adapted for AI
agent skill libraries. Its quality claim is derived from that methodology, which is
adopted at scale by Wikipedia (100M+ edits), npm (2M+ packages), Apple App Store
(2M+ apps), and MDN Web Docs (5M+ developers).

The framework does not claim to be a majority-adopted pattern in AI skill library
design — that category is new. It claims to correctly apply a majority-adopted
contribution-standards methodology to a new domain. Reviewers can verify this claim
by running `review-best-practice-skill` on any meta skill in `meta/`.

---

## File Location

```
skills/<domain>/<subdomain>/skills/<skill-name>/SKILL.md
```

Examples:
```
skills/engineering/development/skills/code-review/SKILL.md
skills/health/fitness/skills/design-training-program/SKILL.md
skills/finance/investing/skills/dcf-valuation/SKILL.md
skills/law/contracts/skills/review-saas-agreement/SKILL.md
```

---

## SKILL.md Format

### Frontmatter (required)

```yaml
---
name: verb-first-kebab-case
description: Use when <triggering conditions and context>.
source: <Author/Org, "Title", Year — or "Widely adopted at [Company1, Company2, ...]">
tags: [problem-keyword-1, problem-keyword-2, problem-keyword-3]
related: [companion-skill-name, next-skill-in-sequence]  # optional
practitioner: true  # optional — marks practitioner-contributed skill (see Practitioner path)
---
```

**`name`**

Pattern: `<verb>-<subject>[-<qualifier>]`

- **verb** — imperative present tense, from the approved tier below
- **subject** — the specific thing being acted on (not a generic category word)
- **qualifier** — optional; only add when `verb-subject` collides across domains

**Approved verbs:**

| Verb | Use for |
|------|---------|
| `propose-` | Draft an artifact for human approval |
| `write-` | Author a document, message, or content |
| `review-` | Evaluate quality of something that exists |
| `audit-` | Batch evaluation across many items |
| `design-` | Architect a system, plan, or program |
| `calculate-` | Compute a numeric value or formula |
| `diagnose-` | Identify the root cause of a problem |
| `optimize-` | Improve a measured metric |
| `suggest-` | Recommend options for user selection |
| `deprecate-` | Retire an outdated artifact |
| `plan-` | Create a structured sequence of actions |
| `negotiate-` | Handle a back-and-forth agreement process |
| `apply-` | Apply a technique, method, or framework to a situation |
| `run-` | Facilitate or execute a meeting, session, or process |

Verbs not in this table are allowed if none of the above fit — but vague verbs below are always rejected:

| Reject | Problem | Use instead |
|--------|---------|-------------|
| `do-` | Says nothing | The actual action verb |
| `handle-` | Too vague | `diagnose-`, `review-`, `resolve-` |
| `manage-` | Too vague | `plan-`, `design-`, `audit-` |
| `improve-` | Doesn't say what improves | `optimize-query-latency`, `reduce-churn` |
| `get-` | Ambiguous (fetch? compute?) | `calculate-`, `fetch-`, `extract-` |
| `use-` | Reads like a tutorial | The action the skill actually performs |
| `help-` | Too generic | The actual action verb |

**Subject specificity:**
- If "which kind?" is a valid follow-up, the subject is too generic → reject
- ❌ Too generic: `commit`, `code`, `contract`, `document`, `test`
- ✅ Specific: `conventional-commit`, `pull-request`, `saas-contract`, `api-docs`, `unit-test`

**Abbreviations:**
- Use the industry-recognized abbreviation when it is more commonly searched than the
  full term: `dcf`, `sql`, `okr`, `api`, `roi`, `kpi`, `cta`, `crm`
- Spell out when the abbreviation is domain-internal or ambiguous across domains
- Test: would a practitioner outside the domain recognize it? If yes, abbreviate.

**Qualifier rules:**
- Only add when `verb-subject` collides across domains
- `review-contract` → `review-saas-contract`, `review-employment-contract`
- `design-program` → `design-training-program`, `design-onboarding-program`
- Do NOT add qualifiers preemptively

**Format:**
- Kebab-case only
- ≤50 characters; ideal 2–4 words (verb + 1–2 word subject + optional qualifier)
  - ✅ `calculate-macros` (2), `review-saas-contract` (3), `design-training-program` (3)
  - ❌ `diagnose-slow-database-query-latency-regression` (≤50 chars but 6 words)
- No `skill-`, `best-practice-`, `guide-` prefix — redundant in a skill library
- No noun-first: ~~`macro-calculation`~~, ~~`contract-review`~~, ~~`training-program-design`~~

**Examples:**
```
❌ handle-deployment-failure   → rejected verb
✅ diagnose-deployment-failure

❌ manage-technical-debt       → rejected verb
✅ audit-technical-debt

❌ improve-performance         → rejected verb + generic subject
✅ optimize-query-latency

❌ code-review                 → noun-first
✅ review-pull-request

❌ commit                      → no verb
✅ propose-conventional-commit

❌ contract-review-best-practice-skill       → noun-first, redundant suffix
✅ review-saas-contract
```

**`description`**
- Must start with "Use when"
- Describes WHEN to use the skill — triggering conditions, symptoms, context
- Must NOT summarize the skill's content or steps
- Max 500 characters

**`source`**
- Required — makes every skill's origin auditable by a third party
- Format: `Author/Org, "Title", Year` — enough for a reader to locate the original
- Examples:
  - `Google, "Google Engineering Practices", 2023`
  - `WHO, "Global Guidelines on Hand Hygiene in Health Care", 2009`
  - `CFA Institute, "Standards of Practice Handbook", 12th ed. 2022`
  - `ABA, "Model Rules of Professional Conduct", 2023`
  - `ISDA, "Master Agreement", 2002`
  - `McKinsey & Company, "Problem Solving", internal methodology`
  - `BCG, "Strategic Analysis Frameworks"`
  - `Bacchelli & Bird, "Expectations, Outcomes, and Challenges of Modern Code Review", ICSE 2013`
  - `Widely adopted at Google, Netflix, Stripe, Airbnb` — acceptable when no single document exists

**`tags`** (required)
- Used by `suggest-best-practice` to match skills to user situations — cover all four axes:
  - **Problem**: what problem does this skill solve? (`code-quality`, `muscle-gain`, `debt-reduction`)
  - **Tool/method**: what technique or tool is involved? (`git`, `progressive-overload`, `dcf`, `sql`)
  - **Role/context**: who uses this / in what context? (`developer`, `athlete`, `startup`, `manager`)
  - **Outcome**: what result does the user get? (`defect-reduction`, `strength-gain`, `cost-savings`)
- Not domain names (`health`, `engineering` — those are captured by file path)
- 3–8 tags, lowercase kebab-case
- Examples by domain:
  - Engineering: `code-quality`, `git`, `pull-request`, `developer`, `defect-reduction`, `test-flakiness`
  - Health: `muscle-gain`, `progressive-overload`, `athlete`, `strength-gain`, `injury-prevention`
  - Finance: `retirement-planning`, `dcf`, `investor`, `portfolio-risk`, `tax-optimization`
  - Law: `saas-contract`, `negotiation`, `founder`, `ip-protection`, `liability-reduction`
  - Business: `team-performance`, `okr`, `manager`, `cost-reduction`, `hiring`

```yaml
# Bad — summarizes content
description: Use when committing — inspects staged files, drafts conventional commit message, presents for approval.

# Good — triggering conditions only
description: Use when the user asks to commit, wants a commit message, or invokes /propose-commit.
```

**`related`** (optional)
- Companion skills commonly used before, after, or alongside this one
- Use when skills form a natural sequence or reinforce each other
- Examples:
  - `apply-first-principles` → `related: [design-solution, evaluate-tradeoffs]`
  - `diagnose-deployment-failure` → `related: [write-postmortem, design-runbook]`
- Do NOT list skills just because they share a domain — only list genuine sequences

### Content structure

**Required:**
- `# Title` — Title Case version of the skill name
- One-sentence purpose statement immediately after the title

| Section | Required? |
|---------|-----------|
| `## Why This Is Best Practice` | **Required** |
| `## Steps` or `## Core Pattern` | **Required** |
| `## Rules` | Recommended |
| `## Examples` | Recommended |
| `## Common Mistakes` | **Recommended** |
| `## When NOT to Use` | **Recommended** |

### Why This Is Best Practice (required)

Every skill must argue its own case. This section proves the skill belongs in grimoire.

Required content:
- **Adopted by**: name the most top-tier companies or credentialed professionals — not vague "many companies"
- **Impact**: measurable outcome with evidence — defect reduction, time saved, performance gain, risk reduction
- **Why best**: why this approach over alternatives
- **Sources**: institution, engineering blog, clinical guideline, or standard body

```markdown
## Why This Is Best Practice

**Adopted by:** Google, Microsoft, Meta, Stripe, and virtually all software companies
with >50 engineers. Codified in Google's Engineering Practices documentation.
**Impact:** Structured code review reduces post-ship defects by ~50% (Google internal
data). #1 defect-detection technique per Microsoft Research (Bacchelli & Bird, 2013).
**Why best:** Catches logic errors, design issues, and security flaws before merge —
earlier than testing, cheaper than post-production fixes, faster than formal inspection.

Sources: Google Engineering Practices, Microsoft Research ICSE 2013
```

```markdown
## Why This Is Best Practice

**Adopted by:** Virtually all elite strength coaches; foundational in NSCA Essentials
of Strength Training and ACSM Exercise Guidelines for all competitive athletic programs.
**Impact:** Progressive overload produces 2–3× greater strength gains vs. fixed-load
training over 12 weeks [Ralston et al., 2017, systematic review].
**Why best:** Only training method with consistent long-term adaptation evidence;
fixed-volume alternatives plateau within 4–6 weeks.

Sources: NSCA CSCS Exam Content, ACSM Position Stand on Resistance Training
```

---

## 5 Quality Criteria

Every skill must pass all five.

### 1. Actionable

The reader can DO something immediately. Steps are concrete and commandable.

```
❌ "Good tests should be maintainable, reliable, and independent."
✅ "Step 1: Write the assertion before writing the implementation.
   Step 2: Run the test — it must fail. If it passes, it tests nothing."
```

### 2. Scoped

One skill = one concept. If a skill covers two separable things, split it.

```
❌ code-review-and-refactoring (two different activities)
✅ code-review  +  refactor-safely  (two skills)

❌ nutrition-and-exercise-planning (two disciplines)
✅ calculate-macros  +  design-training-program  (two skills)
```

### 3. Industry-proven

Adopted by **most** top-tier companies or credentialed professionals in the domain,
with demonstrated strong impact. Majority adoption at the top — not a niche technique,
not one company's approach, not an emerging trend.

A practice must also be **falsifiable** — there must exist possible evidence that could
prove it wrong. If no conceivable outcome could disprove the claim, it is a belief or
philosophy, not a practice. Unfalsifiable claims are always rejected regardless of which
qualification path is used.

```
❌ "Drink enough water throughout the day." (generic, not industry standard)
✅ "For endurance athletes, sodium supplementation during efforts over 90 min
   prevents hyponatremia [ACSM/NSCA position stand]. Target 500–1000mg/hr."

❌ "Review contracts carefully before signing." (generic)
✅ "Uncapped liability clauses are the highest-risk SaaS provision per ABA
   commercial practice standards. Check: unlimited indemnification scope,
   no mutual fee-based cap, consequential damages not excluded."
```

Qualifying sources by domain:

| Domain | Top-tier sources |
|--------|-----------------|
| Engineering | Google, Netflix, AWS, Stripe, Airbnb engineering practices |
| Finance | CFA Institute, Bridgewater, Goldman Sachs, JPMorgan frameworks |
| Health | WHO, Mayo Clinic, NIH, ACSM, NSCA clinical guidelines |
| Law | ABA standards, BigLaw practices, ISDA/NVCA model agreements |
| Marketing | P&G brand management, HubSpot, Ogilvy advertising principles |
| Business | McKinsey/BCG/Bain frameworks, HBS/INSEAD core curriculum |
| Design | Apple HIG, Material Design, Nielsen Norman Group standards |
| Sports | NSCA, USA Weightlifting, Olympic coaching methodology |
| Psychology | APA clinical guidelines, CBT evidence base (Beck Institute) |
| Education | Bloom's Taxonomy, Hattie's Visible Learning, Sweller's CLT |

**Excluded from this criterion:**
- Fame or follower count without verifiable organizational outcomes
- "Interesting techniques" adopted by only a few top-tier companies
- Marginal optimizations when stronger proven practices exist
- One company's internal convention not adopted elsewhere
- Emerging practices not yet proven at scale
- Abstract philosophies that cannot be expressed as actionable steps
- Large or popular community adoption without top-tier endorsement — community size is not a quality signal

#### Visionary practitioner exception

A practice also qualifies if it meets **all three**:

1. **Championed by a proven visionary** — a founder, CEO, or domain leader with verifiable outsized organizational outcomes at global scale (e.g., Musk at SpaceX/Tesla, Jobs at Apple, Dalio at Bridgewater, Bezos at Amazon, Buffett at Berkshire Hathaway)
2. **Outcomes are verifiable** — the organization's results can be cited as evidence (cost reduction %, market leadership, fund returns, etc.)
3. **Specific and actionable** — a concrete step or method, not a general philosophy or mindset

This path is not a shortcut — it requires the same specificity and actionability as majority-adoption practices.

```
❌ "Reason from first principles" — fails criterion 3 (too abstract to act on)
✅ "Question every requirement before optimizing it; then optimize; then automate" — Musk's manufacturing rule, verifiable: SpaceX reduced launch cost by ~90%

❌ "Think different" — fails criterion 3 (philosophy, not method)
✅ "Show the product in 30 seconds or it doesn't ship" — Jobs' product clarity rule, verifiable: iPhone became most valuable product ever
```

**Disqualifiers:**
- Motivational speakers, influencers, or public figures whose fame is not tied to verifiable organizational results
- Practices tied to a single domain outcome that does not generalize beyond that org
- Personal opinions of visionaries outside their domain of proven outcomes (e.g., Musk on nutrition)

**Source format for visionary practices:**
```
source: Elon Musk manufacturing principles (SpaceX) — verified: ~90% launch cost reduction vs. incumbents
source: Steve Jobs product simplicity (Apple) — verified: iPhone became highest-revenue product in history
source: Ray Dalio radical transparency (Bridgewater) — verified: world's largest hedge fund by AUM
```

#### Contested practices

If credible top-tier professionals actively debate a practice, it is **not automatically disqualified** — but requires explicit handling:

- State the majority position and its evidence
- Acknowledge the minority position and its strongest argument
- Skill encodes the majority-adopted approach, not a verdict

If the split is roughly 50/50 among top-tier (genuine controversy with no consensus), exclude it — no consensus means no best practice.

```
❌ "Always use microservices architecture" — debated at the top tier
✅ "Microservices adoption criteria" — encodes when to use vs. monolith with evidence from both sides
```

#### Emerging practices

A practice qualifies as *emerging* (not yet best practice, but on the trajectory) if it meets **all four**:

1. **Adopted by several top-tier orgs** — at least 3–5 named top-tier companies or credentialed professionals with positive early results (not majority, but not isolated either)
2. **Early evidence of impact** — measurable outcome exists even if not at full scale: pilot data, early adopter case studies, or peer-reviewed pre-print research
3. **Specific and actionable** — same bar as best practices; no abstract philosophies
4. **Not contested at the top** — if top-tier professionals are actively debating it, wait for consensus

Emerging skills **must** carry `emerging: true` in frontmatter and a status note in `## Why This Is Best Practice`:

```markdown
**Status:** Emerging — adopted by [org1, org2, org3] with early evidence. Not yet majority top-tier adoption. Review for promotion or deprecation by [YYYY].
```

**Auto-deprecation rule:** If an emerging skill has not achieved majority top-tier adoption within 2 years of addition, run `audit-best-practice-domain` and either promote it (upgrade to best practice) or deprecate it.

```
❌ "AI agents will replace all software engineers" — contested, speculative, no evidence
❌ "Use LLMs for all code review" — interesting but unproven at scale
✅ "Structured prompt templates (XML tags) for LLM instructions" — adopted at Anthropic, OpenAI, Google DeepMind; early evidence of output reliability improvement
✅ "Continuous deployment with feature flags" — adopted at Netflix, Etsy, GitHub before it became mainstream; reduced deployment risk demonstrated in early case studies
```

#### Community standard path

A practice also qualifies if it is codified by a **recognized standards body** and
meets **all three**:

1. **Body is authoritative in its domain** — an established international, national,
   or professional standards organization with a formal review process (see approved
   bodies below)
2. **Standard is current** — not superseded, withdrawn, or in draft-only status
3. **Specific and actionable** — a concrete requirement or procedure from the standard,
   not a general goal or aspiration

**Approved bodies by domain:**

| Domain | Bodies |
|--------|--------|
| Security | NIST, OWASP, ISO/IEC 27000-series, CIS |
| Networking / Web | IETF (RFCs), W3C, IEEE |
| Accessibility | W3C WCAG |
| Engineering | IEEE, ANSI, ISO |
| Health / Medicine | WHO, NIH, CDC, ACSM, NSCA |
| Law | ABA, ISDA, NVCA model agreements |
| Finance | CFA Institute, FASB, IFRS Foundation |
| Psychology | APA, NICE guidelines |

Bodies not in this table are allowed if they meet criterion 1 — but self-published
"community standards" (GitHub orgs, individual blogs) are always rejected.

**Source format for community standard practices:**
```
source: OWASP Top 10 (2021) — A03:2021 Injection
source: NIST SP 800-53 Rev 5 — AC-2 Account Management
source: W3C WCAG 2.2 — Success Criterion 1.4.3 Contrast (Minimum)
source: IETF RFC 9110 — HTTP Semantics, Section 9.3.1
```

#### Practitioner path

A practice also qualifies if contributed by a **verified domain practitioner** and
meets **all four**:

1. **Credentials are verifiable** — contributor names a specific certification,
   competition record, years of recognized practice, or publication. Self-declared
   mastery is not sufficient.
2. **Practice is specific and actionable** — concrete steps, not a philosophy or mindset
3. **Practice is falsifiable** — a conceivable outcome could disprove it; unfalsifiable
   claims are rejected on this path exactly as on all others
4. **Scoped to practitioner's domain of expertise** — a master saddler contributes
   leatherwork skills, not nutrition advice

Practitioner skills **must** carry `practitioner: true` in frontmatter. They are
displayed in a separate README section and are not mixed with Active/Stable skills.

**Source format for practitioner contributions:**
```
source: Jane Smith, Master Saddler, LCSB-certified, 30yr practice
source: Dr. John Doe, APA-licensed CBT therapist, 15yr specialization
source: Park Jae-won, 5× national archery champion, KAA Level 3 coach
```

**Promotion path:** If a practitioner skill later achieves majority top-tier adoption,
remove `practitioner: true` and promote to Active via the normal lifecycle.

---

### 4. Specific over general

Concrete examples beat abstract rules. Name the tool, cite the number, show the command.

```
❌ "Use appropriate data structures for your use case."
✅ "O(1) lookup by key → hash map. Sorted iteration with O(log n) insert →
   balanced BST. FIFO queue with O(1) both ends → deque."

❌ "Diversify your investment portfolio."
✅ "At accumulation stage (20+ years to retirement): 90% equities (split 60%
   domestic total market, 30% international total market), 10% bonds.
   Rebalance annually or when any allocation drifts >5% from target."
```

### 5. Complete

Covers failure modes, edge cases, and the non-obvious constraints — not just the happy path.

```
❌ Steps that only work when everything goes right
✅ Steps that include: "If X fails, do Y. If the output is ambiguous, ask
   for clarification before proceeding."
```

---

## User Agency: Multiple Matching Practices

Any skill that scores, matches, or routes to best practices MUST follow this rule:

**When multiple practices could apply, always present a ranked list with a recommendation and let the user decide. Never silently pick one.**

### Decision rule

| Result | Condition | Required behavior |
|--------|-----------|------------------|
| **Sole clear match** | 1 practice scores ≥ 0.7 AND second place < 0.4 | Apply directly — no list needed |
| **Multiple candidates** | 2+ practices score ≥ 0.4 | Present ranked list, mark ★ recommendation, wait for user choice |
| **No match** | All practices score < 0.4 | State no match found; ask one clarifying question or flag for manual research |

### List format (required when multiple candidates)

```
Multiple best practices apply. Recommended: [top-practice-name]

1. ★ [top-practice] — [one sentence: what problem it solves]  ← recommended
   Domain: [domain/subdomain]

2. [second-practice] — [one sentence: what problem it solves]
   Domain: [domain/subdomain]

3. [third-practice] — [one sentence]
   Domain: [domain/subdomain]

Which would you like to apply? (Enter number, or press Enter for ★ recommended)
```

### Rules

- ★ marks the highest-scoring match. If two score equally, ★ goes to the one whose `Use when...` description is the closest literal match to the user's input.
- Max 5 options in a ranked list — drop lower-scoring matches beyond 5.
- Never hallucinate practice names — only list practices that exist in installed domains.
- If a practice is not installed, include the install command alongside its entry.
- This rule applies to every skill that performs scoring/matching: `suggest-best-practice`, `intercept-best-practice`, `review-best-practice-fit`, `plan-best-practice-solution`, `apply-best-practice-tree`, and any future skill that routes to other skills.

---

## Domain Safety Standards

Skills in regulated or sensitive domains carry additional obligations.

### Health / Medicine

- Cite the evidence level when making factual claims:
  - `[RCT]` — randomized controlled trial evidence
  - `[SR]` — systematic review / meta-analysis
  - `[Consensus]` — expert consensus / clinical guidelines
  - `[Practical]` — widely used practice without strong RCT support
- Never prescribe specific medications, diagnose conditions, or give personalized medical advice
- Required footer on every health skill:

```markdown
> For personal health decisions, consult a qualified healthcare provider.
```

### Law

- State the jurisdiction scope explicitly: "US law", "EU GDPR", "general contract principles (common law)"
- Don't give advice on specific legal situations
- Required footer on every law skill:

```markdown
> This is educational content, not legal advice. Consult qualified legal counsel for your situation.
```

### Finance / Investing

- Note where past performance does not predict future results
- Distinguish between general principles and personalized financial advice
- Required footer on every finance skill:

```markdown
> This is educational content, not financial advice. Consult a licensed financial advisor for personal decisions.
```

### Psychology / Mental Health

- Don't diagnose conditions or prescribe treatment protocols
- Don't recommend specific medications
- Required footer on every mental health skill:

```markdown
> For mental health concerns, consult a qualified mental health professional.
```

---

## What NOT to Include

| Exclude | Why |
|---------|-----|
| Generic advice any non-expert already knows | Adds no value ("stay hydrated", "test your code") |
| Two separable concepts in one file | Split into two skills |
| Background theory that belongs in a textbook | Link to authoritative source instead |
| Style preferences without expert justification | Opinions aren't skills |
| Content that's already well-covered by documentation | Link; don't duplicate |

---

## Size

- **Target: 50–300 lines**
- Under 50 lines: probably too shallow — add more concrete steps or examples
- Over 300 lines: probably two skills — split it

---

## Skill Freshness

Skills become outdated when:
- The source institution revises its position
- A newer practice achieves majority top-tier adoption that supersedes this one
- The tools or context referenced no longer exist at scale

### Fix vs. Deprecate

| Situation | Action |
|-----------|--------|
| Source institution revised its position | Fix — update Why section and sources |
| Tool referenced is outdated but practice is sound | Fix — update tool reference only |
| A newer practice supersedes this one at majority top-tier | Deprecate — point to replacement |
| Skill never qualified (Adopted by was always inaccurate) | Deprecate — no replacement needed |
| Skill is correct but scope was wrong (too broad or narrow) | Fix — adjust scope, re-run review-best-practice-skill |
| 50/50 controversy has resolved to a clear consensus | Fix — update to reflect new consensus |

When submitting a PR that supersedes an existing skill, mark the old skill for
deprecation in the PR description. Maintainers remove deprecated skills in the
next release cycle.

---

## Skill Lifecycle

Every skill moves through a defined lifecycle. The current state is declared in frontmatter.

| State | Frontmatter flag | Meaning |
|-------|-----------------|---------|
| **Proposed** | *(no file — PR under review)* | Awaiting merge |
| **Emerging** | `emerging: true` | Early adopters, 2-year promotion window |
| **Active** | *(none — default)* | Majority adoption, cited evidence |
| **Stable** | `stable: true` | 5+ years proven, rarely needs updates |
| **Deprecated** | `deprecated: true` + `deprecated_by: skill-name` | Superseded or retired |
| **Practitioner** | `practitioner: true` | Expert-contributed; no org-scale adoption evidence required |

```
PROPOSED → EMERGING    → ACTIVE → STABLE → DEPRECATED
         ↘ PRACTITIONER ↗
(PR state)  (tagged)     (default) (tagged)   (tagged)
```

### Promotion criteria

**Emerging → Active:** Remove `emerging: true`. Requires: 2+ years of adoption, ≥10 top-tier organizations, impact evidence, maintainer sign-off.

**Active → Stable:** Add `stable: true`. Requires: 5+ years uncontested, broad cross-industry adoption, content unchanged over 12+ months.

**→ Deprecated:** Add `deprecated: true` and `deprecated_by: <replacement-skill-name>` (or `deprecated_by: none` if retired with no successor). Skill file is **kept, never deleted** — so links remain valid.

**Emerging → Deprecated:** 2-year window elapsed with no promotion. Add `deprecated: true` + `deprecated_by: none` and note the reason in the Why section.

### Conflicting states

The following combinations are invalid and rejected by the validator:

- `stable: true` + `emerging: true` — a skill cannot be both unproven and proven
- `deprecated: true` + `emerging: true` — a retired skill is not also in trial
- `deprecated: true` without `deprecated_by:` — deprecation must name the successor or `none`
- `practitioner: true` + `stable: true` — practitioner skills have no adoption evidence; stable requires 5+ years broad cross-industry adoption

---

## Verified Tier

Skills marked `verified: true` in frontmatter receive a `✓ Verified` badge in the README featured section.

**Criteria — all three must be met:**

1. **Attributed** — `source` names a specific institution, company, or published methodology (not vague "widely adopted")
2. **Tested** — contributor has used this practice in production or a real engagement; noted in PR description
3. **Reviewed** — a maintainer has run `review-best-practice-skill` and approved without major findings

**How to request verification:**
- Add `verified: true` to your skill's frontmatter
- In your PR, write one sentence: "I have used this in production at [context]."
- A maintainer will run `review-best-practice-skill` and either approve or return findings via `revise-best-practice-skill`

**Revoking verification:**
- If `audit-best-practice-domain` or `review-best-practice-skill` finds the skill outdated or inaccurate, `verified` is removed pending a fix
- Run `revise-best-practice-skill` to restore it

---

## Machine-Readable Schema

The Grimoire Skill Standard is formally specified as a JSON Schema:

```
schema/skill.schema.json
```

Any tool that validates SKILL.md frontmatter should conform to this schema. IDE plugins, GitHub Actions, and third-party validators can reference the schema directly to validate frontmatter without parsing the prose standard.

**Conformance test suite:** `schema/tests/` contains canonical SKILL.md fixtures — valid and invalid. A validator implementation is conformant when it produces the correct result for every fixture:

- `schema/tests/valid/` — must PASS (exit 0)
- `schema/tests/invalid/` — must FAIL (exit 1)

The reference implementation is `scripts/validate-skill.sh`. Run the full suite with:

```bash
bash scripts/test-schema.sh
```

---

## Adopting This Standard

Any AI agent skill library may adopt this standard. To declare compliance:

1. Reference this standard in your contributing guide: `Skills must comply with the Grimoire Skill Standard v1.0`
2. Use the [Review Checklist](#review-checklist) to gate contributions
3. Open an [Adopt the Standard](https://github.com/jeffreytse/grimoire/issues/new?template=adopt-standard.yml) issue to be listed in [ADOPTERS.md](./ADOPTERS.md)

Adopters gain: a proven quality framework, a shared contributor vocabulary, and cross-project skill compatibility.

**Automated enforcement:** The reference CI configuration is at `.github/workflows/validate.yml`. Copy it to your repo to enforce the standard on every pull request automatically.

**Governance:** Changes to this standard follow the process defined in [GOVERNANCE.md](./GOVERNANCE.md), including a 7-day discussion period and maintainer approval.

---

## Review Checklist

Before submitting a skill, verify:

- [ ] `name` passes naming standard (see `**\`name\`**` section): verb-first, specific subject, 2–4 words, abbreviation policy applied, no rejected verbs
- [ ] `description` starts with "Use when"
- [ ] `description` describes triggering conditions ONLY — no content summary
- [ ] Title is h1, matches name in Title Case
- [ ] One-sentence purpose statement after title
- [ ] `## Why This Is Best Practice` section present
- [ ] Section names specific top-tier companies or credentialed professionals (not vague "many")
- [ ] Section states measurable impact with evidence
- [ ] Section explains why this approach over alternatives
- [ ] Steps are concrete and immediately actionable
- [ ] Scoped to one concept
- [ ] Industry-proven — passes at least one qualification path: (a) majority top-tier company/professional adoption, (b) visionary practitioner with verifiable outcomes, (c) recognized standards body (NIST, OWASP, IETF, W3C, IEEE, WHO, ABA, etc.), or (d) verified domain practitioner (`practitioner: true`)
- [ ] If `practitioner: true`, `source` names specific verifiable credentials (certification, competition record, publication, or years in recognized context)
- [ ] If `practitioner: true`, skill is falsifiable and actionable — practitioner path cannot substitute for a philosophy or ideology
- [ ] `tags` present with 3–8 tags covering all 4 axes: problem keyword, tool/method, role/context, outcome
- [ ] `source` field present and cites credible institution or top-tier adopters
- [ ] Practice has strong impact — not a marginal optimization
- [ ] Specific: names tools, cites numbers, shows commands
- [ ] Covers edge cases and failure modes
- [ ] Safety footer present (health / law / finance / mental-health)
- [ ] Not superseded by a newer practice with equal or broader top-tier adoption
- [ ] 50–300 lines
- [ ] Practice is falsifiable — conceivable evidence could disprove it (not a philosophy or belief)
- [ ] Lifecycle flags consistent — no conflicting combinations (`stable` + `emerging`, `deprecated` + `emerging`, `deprecated` without `deprecated_by`)
- [ ] If `related:` present, lists genuine sequences only (not domain neighbors)
- [ ] If health domain, evidence-tier tags present on factual claims (`[RCT]`, `[SR]`, `[Consensus]`, `[Practical]`)
- [ ] If `emerging: true`, `## Why This Is Best Practice` includes status note with named orgs and review date
