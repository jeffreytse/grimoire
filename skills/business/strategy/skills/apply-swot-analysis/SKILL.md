---
name: apply-swot-analysis
description: Use when assessing a business, product, team, or initiative before a strategic decision — to systematically surface internal capabilities and external conditions, then derive concrete strategic options rather than an unactionable list.
source: 'Albert Humphrey, Stanford Research Institute SWOT framework (1960s, widely attributed); Weihrich, "The TOWS Matrix — A Tool for Situational Analysis", Long Range Planning (1982); Porter, "Competitive Strategy", Harvard Business School Press (1980); McKinsey Strategic Analysis frameworks; BCG competitive analysis methodology'
tags: [swot, strategic-analysis, decision-making, competitive-analysis, strategy, situational-analysis, business-planning]
---

# Apply SWOT Analysis

Evaluate Strengths, Weaknesses, Opportunities, and Threats across the correct axis (internal vs. external), then cross them to produce four specific strategy types — SO, WO, ST, WT — translating a 2×2 diagnostic into concrete strategic options.

## Why This Is Best Practice

**Adopted by:** Standard in MBA core curriculum at Harvard Business School, Stanford GSB, Wharton, INSEAD, and virtually every accredited business program worldwide. Applied as the initial situational framing step at McKinsey, BCG, and Bain before strategy design. Embedded in military planning doctrine (US Army Strategic Planning Process), government policy analysis (OECD strategic review templates), and startup accelerator curricula (Y Combinator, Techstars). Albert Humphrey's original framework from Stanford Research Institute has been reproduced in every major strategy textbook published since 1980, including Porter's "Competitive Strategy" and Mintzberg's "Strategy Safari."
**Impact:** Weihrich (1982) Long Range Planning documented the TOWS cross-matrix extension: companies that derived explicit SO/WO/ST/WT strategies from their SWOT inventories produced 2–3× more actionable strategy options versus those who stopped at the 4-quadrant list. McKinsey case documentation: 80%+ of strategy engagements begin with situational analysis that includes a SWOT or equivalent; the framework's consistent adoption comes from its ability to structure diverse stakeholder input into a shared diagnosis before any strategic debate begins. The mechanism: unstructured strategy discussions produce opinion-based debates; SWOT creates a shared factual inventory that moves debate from "what should we do?" to "given these specific facts, what options do we have?"
**Why best:** The alternative — jumping directly to strategy generation without situational inventory — is the most common failure mode in strategic planning. Teams that skip situational analysis debate solutions to problems they haven't verified, miss threats they haven't named, and build on assumed strengths they haven't validated. SWOT's value is not the 4-quadrant diagram — it is the discipline of completing all four quadrants before any strategic option is proposed, and the TOWS cross-matrix that derives strategy from the intersection of facts rather than from intuition.

Sources: Weihrich (1982) Long Range Planning; Porter (1980) Competitive Strategy; McKinsey strategic analysis methodology; Humphrey SRI framework (1960s)

## Steps

### 1. Define the scope — what exactly is being analyzed?

Before filling any quadrant, state the unit of analysis precisely:

```
Scope examples:
  ✅ "Company X's product line Y in the SMB market segment for Q3 strategic planning"
  ✅ "Our mobile app feature set vs. Competitor A for the upcoming roadmap review"
  ✅ "The engineering team's capacity to support a Series A growth plan"

❌ Too broad: "Our company" — SWOT for "the company" produces entries that belong
   to different business units, mixed timeframes, and different competitive contexts
```

Write the scope as a single sentence before opening the 2×2. Every entry you add will be evaluated against this scope.

### 2. Enforce the internal/external axis — the most common failure point

```
INTERNAL (under your control):          EXTERNAL (outside your control):
  Strengths   │  Weaknesses               Opportunities  │  Threats
──────────────┼──────────────           ─────────────────┼────────────
              │                                          │
```

The single most common SWOT error is misclassifying items:

```
Misclassification examples:
  ❌ "Growing market demand" → listed under Strengths
     Correct: market demand is external → Opportunities
  
  ❌ "Competitor launched a new product" → listed under Weaknesses
     Correct: competitor action is external → Threats

  ❌ "We don't have funding yet" → listed under Threats
     Correct: funding gap is internal → Weaknesses

  ❌ "Regulatory change might affect us" → listed under Weaknesses
     Correct: external regulation → Threats (or Opportunities if favorable)
```

**Test for each entry:**
- "Can we change this directly through our own decisions?" → **Internal** (S or W)
- "Does this exist in the environment regardless of our decisions?" → **External** (O or T)

### 3. Fill each quadrant with evidence-based specifics

For each quadrant, aim for 5–8 specific items. Generic items produce generic strategies.

```
Strengths (internal positive):
  ❌ "Strong team" — what specifically? Why is it a competitive advantage?
  ✅ "3× faster deployment cycle than industry average due to CI/CD pipeline"
  ✅ "Exclusive distribution agreement with [Partner] through 2027"
  ✅ "Proprietary dataset of 10M+ user interactions unavailable to competitors"

Weaknesses (internal negative):
  ❌ "Limited resources" — every company has limited resources
  ✅ "No dedicated sales team — growth depends entirely on inbound channels"
  ✅ "Single customer accounts for 40% of ARR — concentration risk"
  ✅ "Engineering team has no ML expertise — required for roadmap items H2"

Opportunities (external positive):
  ❌ "AI is growing" — what specific opportunity does this create for us?
  ✅ "Healthcare sector mandate for digital records by 2026 creates forced adoption"
  ✅ "Lead competitor raised prices 20% — their SMB customers are now evaluable"
  ✅ "Remote work normalization expanded total addressable market by 3×"

Threats (external negative):
  ❌ "Competition" — what specific threat from which competitor?
  ✅ "Google entering our category with $0 freemium offering in Q2"
  ✅ "EU Digital Markets Act requires interoperability — eliminates lock-in advantage"
  ✅ "Key talent segment now receiving 40% higher offers from Series B companies"
```

**Watch for wishful thinking in Opportunities:**
An Opportunity must be external and real — not something you hope will happen or a strength framed as external. If you can control it, it's a Strength. If it requires a market condition that hasn't been validated, it is an assumption, not an opportunity.

### 4. Cross the quadrants — the TOWS strategy matrix

This is the step most teams skip. The SWOT list without the cross-matrix is a diagnostic, not a strategy. Derive one actionable option from each intersection:

```
                   │  STRENGTHS (S)        │  WEAKNESSES (W)
───────────────────┼───────────────────────┼──────────────────────────
OPPORTUNITIES (O)  │  SO Strategy:         │  WO Strategy:
                   │  Use strengths to     │  Address weaknesses to
                   │  pursue opportunities │  capture opportunities
───────────────────┼───────────────────────┼──────────────────────────
THREATS (T)        │  ST Strategy:         │  WT Strategy:
                   │  Use strengths to     │  Minimize weaknesses to
                   │  mitigate threats     │  avoid threats
```

For each cell, write one specific strategy option using the SWOT items you populated:

```
SO example:
  S: "Exclusive distribution partner through 2027"
  O: "Healthcare mandate creating forced adoption in 2026"
  → "Activate distribution partner's healthcare vertical for mandate-driven pipeline
     in Q4 — target 50 accounts before the mandate deadline"

WO example:
  W: "No ML expertise on engineering team"
  O: "AI tooling now commoditized through APIs"
  → "Use pre-built ML APIs (not custom models) to ship AI features in H2 without
     hiring — removes capability gap without headcount"

ST example:
  S: "3× faster deployment than industry average"
  T: "Google entering with $0 freemium offering"
  → "Compete on iteration speed: ship 1 major feature per week vs. Google's slower
     enterprise cycle — position as 'the responsive alternative'"

WT example:
  W: "Single customer = 40% ARR concentration"
  T: "Key talent receiving 40% higher offers"
  → "Diversify top-3 customer concentration to <25% each before attempting
     Series A — de-risks both investor concern and key talent retention pitch"
```

### 5. Prioritize — reduce to 3 actionable options

A SWOT that produces 20 strategy options is not useful. Force prioritization:

```
Scoring each option on two dimensions (1–3 scale each):
  Impact:    3 = strategic, changes competitive position
             2 = meaningful, improves a measurable metric
             1 = incremental, marginal improvement

  Feasibility: 3 = can start this quarter with current resources
               2 = requires one new resource or partnership
               1 = requires fundamental capability change

Priority = Impact × Feasibility (max 9)
```

Select the **top 3** strategy options. Assign an owner and a 90-day next action to each.

```
Output format:
  Priority 1: [SO/WO/ST/WT] "[Strategy name]"
              Score: Impact 3 × Feasibility 3 = 9
              Owner: [name/role]
              Next action: [specific action within 90 days]

  Priority 2: ...
  Priority 3: ...
```

## Rules

- Complete all four quadrants before proposing strategy — proposing SO strategies with an empty Threats quadrant skips risk analysis; require full completion
- Internal vs. external axis is non-negotiable — if an entry doesn't pass the "can we control this directly?" test, move it; misclassified entries produce misaligned strategy
- Evidence required for each entry — every item must be a verifiable fact or a stated assumption (if assumption, label it as one); entries that are intuitions without basis belong in a separate "to validate" list
- The TOWS cross-matrix is mandatory — stopping at the 4-quadrant inventory is the diagnostic without the prescription; always complete all four SO/WO/ST/WT cells

## Common Mistakes

- **Mixing internal and external** — the most common mistake; teams list competitor actions under Weaknesses ("our competitor is better") and market trends under Strengths ("AI is growing" ← this is external); enforce the axis test on every entry before proceeding
- **Wishful thinking as Opportunities** — "market will grow 10×" is a forecast, not a verified opportunity; label unverified market predictions as assumptions and validate before building strategy on them; real Opportunities are observable current conditions (regulation already passed, competitor already raised prices, market already shifted)
- **Stopping at the list** — a completed 2×2 diagram that produces no strategy options has done half the job; the TOWS cross-matrix is what converts the diagnosis into decisions; without it, SWOT is a documentation exercise
- **Failure to prioritize** — a SWOT session that ends with 30 strategic options and no ranking has not made a decision; force scoring and select 3 maximum; more options = less execution; the value of SWOT is narrowing focus, not expanding it
- **Scope creep across business units** — when the scope is "the company" and participants work in different units, Strengths from one unit appear alongside Threats that apply to a different unit; a mixed-scope SWOT produces contradictory strategy; one scope per session
