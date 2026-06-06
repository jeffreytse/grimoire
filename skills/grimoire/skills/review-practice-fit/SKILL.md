---
name: review-practice-fit
description: Use when the user already has a solution, plan, approach, or design and wants to know how well it aligns with best practices — including gaps, what's missing, and what to fix.
source: McKinsey as-is/to-be gap analysis methodology, Google Engineering Practices design review, ISO 9001 gap audit standards
tags: [gap-analysis, solution-review, practice-alignment, quality-audit, practitioner, decision-maker, solution-improvement, practice-compliance]
---

# Review Practice Fit

Evaluate an existing solution against applicable best practices, identify gaps, and produce a prioritized fix list.

## Why This Is Best Practice

**Adopted by:** McKinsey and BCG use structured as-is/to-be gap analysis before every engagement to identify which best practices a client's current approach violates — it is the primary tool for diagnosing why organizations underperform their industry peers. Google's Engineering Practices mandate design reviews against explicit quality criteria before large features ship. ISO 9001 certification requires formal gap audits comparing current processes against the standard before submission.
**Impact:** Google's structured design review process reduces post-ship defects by ~50% and catches architecture problems that survive code review (Google Engineering Practices documentation). ISO 9001 mandates formal gap audits before certification precisely because self-assessment consistently misses systemic gaps — organizations routinely discover critical non-conformances only during external audits that internal reviews missed (ISO 9001:2015 §9.2 internal audit requirements). McKinsey's structured gap analysis prevents clients from investing in solutions that miss the dimensions that actually drive performance — the most expensive mistake in strategy.
**Why best:** Self-assessment without an external standard is systematically optimistic — practitioners overweight what they did and underweight what they omitted. A structured comparison against explicit best-practice criteria catches omissions (the invisible gaps), not just flaws (the visible ones). Ad-hoc feedback ("this looks good, but maybe add X") is alternative — it finds surface problems only and produces no prioritized action plan.

Sources: Google Engineering Practices; ISO 9001:2015 §9.2 internal audit requirements; McKinsey Problem Solving methodology

## Steps

### 1. Extract the solution

From the user's description, identify:

| Element | Extract |
|---------|---------|
| **What** | What is the solution, plan, approach, or design? |
| **Domain(s)** | Which fields does it operate in? |
| **Goal** | What problem is it trying to solve? |
| **Constraints** | Any stated limitations (time, budget, team size, technology)? |
| **Maturity** | Is this a draft, in-progress, or already deployed? |

If the solution description is too vague to evaluate, ask ONE targeted question:
```
To review this properly, I need to understand [specific missing element].
Can you describe [that element] in more detail?
```

### 2. Identify applicable practices

Score candidate practices using the `suggest-practice` scoring model:

```
score = (tag_overlap × 2) + (description_match × 3) + (domain_plausibility × 1)
```

Select all practices scoring ≥ 0.4. Cap at 7 practices — if more qualify, take the 7 highest-scoring.

If no practice scores ≥ 0.4: state "No installed skills closely match this solution's domain. Install relevant domain skills first."

### 3. Evaluate fit for each practice

For each applicable practice, evaluate the solution against the practice's core criteria:

**ALIGNED** — solution demonstrably follows the practice's key steps and principles
**PARTIAL** — some elements present, but one or more critical criteria are missing or weak
**MISSING** — practice not addressed at all

For each PARTIAL or MISSING verdict, extract:
- What the solution currently does (or doesn't do)
- Which specific criterion from the practice it violates or omits
- The concrete consequence of this gap (what goes wrong without it)

### 4. Prioritize gaps

Classify each gap by impact:

| Priority | When |
|----------|------|
| 🔴 **Critical** | Violates a core principle of the practice; high risk of failure, harm, or waste |
| 🟡 **Significant** | Reduces effectiveness meaningfully; workaround exists but at cost |
| ⚪ **Minor** | Polish or optimization; solution works without it |

Order: Critical → Significant → Minor within the report.

### 5. Produce the gap report

```
## Practice Fit Review

Solution: [one-sentence description of what was evaluated]

---

### [practice-name] — ALIGNED / PARTIAL / MISSING
Domain: [domain/subdomain]

✓ [What the solution gets right — be specific]
✗ [What's missing or weak — cite the specific criterion]
→ Fix: [concrete action, not advice — what exactly to do]

### [practice-name] — ALIGNED / PARTIAL / MISSING
...

---

### Priority gaps

🔴 Critical
- [gap]: [consequence] → [fix]

🟡 Significant
- [gap]: [consequence] → [fix]

⚪ Minor
- [gap]: [consequence] → [fix]

---

### Verdict
[STRONG / ADEQUATE / NEEDS WORK / REBUILD]
[1–2 sentences: overall assessment and single most important action]
```

**Verdict scale:**
- **STRONG**: ≥ 80% of practices ALIGNED, no Critical gaps
- **ADEQUATE**: no Critical gaps, some Significant gaps
- **NEEDS WORK**: 1+ Critical gaps, core structure is sound
- **REBUILD**: 2+ Critical gaps across different practices, fundamental approach is flawed

### 6. Offer to apply fixes

After the report, ask:
```
Want me to apply any of the practices above to help close these gaps?
```

If yes, load the relevant skill and follow its steps to guide the improvement.

## Rules

- Never fabricate practice criteria — evaluate only against the actual steps and rules in each matched skill
- If a practice is not installed, name it and give the install command instead of skipping it
- ALIGNED does not mean perfect — state what qualifies and what could still improve
- Fix recommendations must be concrete and actionable — not "consider improving X"
- If the solution is already deployed (not a draft), flag Critical gaps as urgent risk, not just improvement opportunities
- Do not soften verdicts — a MISSING is a MISSING, even if the overall solution is otherwise strong

## Examples

> Skill names in examples are illustrative — actual matches depend on what domains are installed. If a skill is not installed, `review-practice-fit` names it and gives the install command.

**Example 1 — Engineering architecture review**
> "Our API: REST endpoints, JWT auth, PostgreSQL, deployed on Heroku, no rate limiting, logs to console only."

Matches: `design-api-architecture`, `review-security-posture`, `design-observability`
- `design-api-architecture`: PARTIAL — REST ✓, stateless auth ✓, no versioning ✗, no pagination standard ✗
- `review-security-posture`: MISSING — no rate limiting, no input validation mentioned, JWT secret management unknown
- `design-observability`: MISSING — console logs only, no structured logging, no alerting, no tracing

🔴 Critical: No rate limiting → DoS exposure → add rate limiter at gateway before next deploy
🔴 Critical: No structured logging → incidents uninvestigable → switch to structured JSON logs with correlation IDs

---

**Example 2 — Business plan review**
> "Startup plan: build a mobile app, charge $9.99/month, target college students, raise seed round."

Matches: `design-business-model`, `calculate-unit-economics`
- `design-business-model`: PARTIAL — revenue model ✓, no customer segment validation, no competitive moat stated
- `calculate-unit-economics`: MISSING — no LTV/CAC calculation, no payback period, no cohort assumptions

🔴 Critical: No unit economics → seed investors will reject without LTV/CAC → calculate before fundraising

---

**Example 3 — Strong fit**
> "Code review process: async PR reviews, two approvers required, automated linting and tests must pass, comments must cite a reason, author resolves all comments before merge."

Matches: `review-pull-request`
- ALIGNED — two-approver gate ✓, automated checks ✓, reasoned feedback ✓, resolution required ✓

⚪ Minor: No stated SLA for review turnaround — can cause blocked PRs

Verdict: STRONG — process follows the practice. One minor improvement: add a 24hr review SLA.

## Common Mistakes

**Softening verdicts**: a Critical gap on a deployed system is a risk, not a suggestion. State it clearly.

**Evaluating against invented criteria**: only evaluate against what the matched skill's actual steps say. Don't add your own criteria.

**Skipping ALIGNED items**: report what works too — it confirms the user's judgment and anchors the gaps in context.

**Generic fixes**: "improve your security posture" is not a fix. "Add rate limiting of 100 req/min per IP at the API gateway" is a fix.
