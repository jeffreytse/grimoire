---
name: apply-jobs-to-be-done
description: Use when designing a product, evaluating product-market fit, or investigating why customers switch to or away from a solution — to understand what functional, social, and emotional progress customers are trying to make, independent of your existing product framing.
source: 'Christensen, Hall, Dillon & Duncan, "Know Your Customers'' Jobs to Be Done", Harvard Business Review (2016); Christensen, "The Innovator''s Dilemma", Harvard Business School Press (1997); Moesta & Spiek, "Demand-Side Sales 101" (2020); Ulwick, "Jobs to Be Done: Theory to Practice", IDEA Bite Press (2016)'
tags: [jobs-to-be-done, product-market-fit, customer-research, innovation, strategy, jtbd, product]
---

# Apply Jobs to Be Done

Identify the specific progress a customer is trying to make in a given situation — the functional outcome they need, the social identity they want to project, and the emotional state they seek — then evaluate your product against the job it's hired to do rather than the features it provides.

## Why This Is Best Practice

**Adopted by:** Developed at Harvard Business School by Clayton Christensen and refined at the Clayton Christensen Institute. Taught in HBS MBA core curriculum and at MIT Sloan, Stanford GSB, and INSEAD. Applied formally at Apple (Bob Moesta worked with Jobs on the original Mac and iPod), Intercom (publicly documented JTBD adoption), Basecamp (Jason Fried and DHH cite it explicitly), and Procter & Gamble (Lafley's transformation cited JTBD as core to product strategy). Bob Moesta and Chris Spiek's "Demand-Side Sales 101" and Tony Ulwick's "Outcome-Driven Innovation" are the practitioner implementations adopted across product management communities.
**Impact:** Christensen (2016) HBS case study corpus: companies that designed products around jobs rather than customer segments consistently outperformed segment-focused competitors in product-market fit metrics. Ulwick (2016) documents Outcome-Driven Innovation (ODI) applications: Bosch power tools redesigned using JTBD methodology achieved 35% market share increase; Arm & Hammer baking soda expanded into 11 new product categories by identifying new jobs the product could be hired for. The mechanism: features get commoditized; jobs are stable — "get from A to B quickly" has been the job for horses, cars, and ride-sharing apps; companies that own the job survive disruption.
**Why best:** Traditional segmentation (demographics, psychographics) describes who buys but not why they switch. JTBD answers the switching question: what situation prompted the hire, what progress was sought, what alternatives were considered and why they failed. This is the information that drives product decisions. The alternative (feature-based roadmapping from user surveys) produces locally optimal improvements ("make the UI faster") but misses the structural question: "is this the right product for this job, or are we optimizing the wrong thing?"

Sources: Christensen et al. (2016) Harvard Business Review; Christensen (1997); Moesta & Spiek (2020); Ulwick (2016)

## Steps

### 1. Identify the unit of analysis — one job per session

A "job" is a specific situation in which a person seeks progress:

```
Job structure: [Situation] → [Motivation] → [Outcome]
"When [situation], I want to [progress], so I can [outcome]"

Examples:
  "When I need to get 8 coffees to a meeting, I want to carry them safely,
   so I can arrive looking competent."
  "When I'm learning to code at night after work, I want to make progress
   without losing my place, so I can feel I'm building toward a career change."
  "When my team is remote and async, I want to coordinate without meetings,
   so I can respect everyone's deep work time."
```

Avoid defining the job by reference to your product ("the job of using our app"). Define it by the customer's situation and desired progress.

### 2. Conduct switch interviews — the highest-signal data source

The most revealing JTBD data comes from interviewing recent purchasers or recent churners about the moment of switching:

```
Switch interview structure (Bob Moesta's timeline method):

1. First thought: "When did you first think about making this change?"
   → Surfaces the trigger event (the push from the old solution)

2. Passive looking: "What did you start doing once that thought occurred?"
   → Surfaces the problem they were trying to solve, not what they bought

3. Active looking: "When did you start actively searching for solutions?"
   → Identifies the forcing event (why they acted when they did)

4. Decision: "Walk me through the day you decided."
   → Reveals the hire moment, the alternatives considered, and the anxiety

5. First use: "What happened when you first used it?"
   → Reveals whether the job was satisfied
```

Key: ask about chronology, not opinions. "Walk me through the day you bought it" produces facts. "What do you value in a product?" produces rationalizations.

### 3. Map the four forces of switching

Every switch from one solution to another involves four forces:

```
Forces toward new solution:
  PUSH: frustrations with current solution ("the old way was too slow")
  PULL: attraction of new solution ("I heard this one is much simpler")

Forces against switching:
  ANXIETY: "What if the new thing doesn't work?"
  HABIT: "I know how the current solution works, even if imperfectly"
```

A switch only happens when Push + Pull > Anxiety + Habit. In interviews, identify all four forces. If your product has strong Pull but the Anxiety is not addressed, customers will not switch even when attracted.

### 4. Define the job in three dimensions

After 5–8 switch interviews, synthesize the job across three dimensions:

```
Functional dimension:
  What task must get done? What outcome is measurable?
  "Transfer 50 contacts from my phone to a CRM in under 10 minutes"

Social dimension:
  How does the person want to be seen by others?
  "Look organized and professional to my manager"
  "Be the person who brings useful tools to the team"

Emotional dimension:
  How does the person want to feel?
  "Not anxious that important contacts will be lost"
  "Confident that my system is under control"
```

Products that only satisfy the functional dimension compete on efficiency. Products that satisfy all three dimensions create loyalty. The emotional and social dimensions are where switching costs are built.

### 5. Evaluate your product against the job

With the job defined, assess your product honestly:

```
For each dimension of the job:
  Functional: Does the product complete the task reliably? Where does it fail?
  Social: Does using this product make the customer look good to others?
  Emotional: Does the product reduce anxiety or create it?

Hire/fire test:
  "Is our product hired for this job?" 
  → Check: are customers actually using it for this job, or a different one?
  
  "Is our product at risk of being fired?"
  → Check: what competing solution could do this job better?
  → Check: are there struggles we're ignoring because they seem minor?
```

### 6. Apply to product decisions

JTBD changes which product questions are worth asking:

```
Instead of:                         Ask instead:
"What features do users want?"   → "What progress are they trying to make?"
"What's our target demographic?" → "What job are they hiring for?"
"Why did they churn?"            → "What job did we fail to do, or what job did they hire someone else for?"
"What's our competitive moat?"   → "Do we own the job, or do we own features that can be copied?"
```

## Rules

- Define the job independent of your product — "the job of using our app" is not a job; "the job of staying connected with distributed teammates without synchronous meetings" is a job; the former traps you in feature-optimization, the latter shows you the full competitive landscape
- Interview recent switchers, not long-term users — long-term users have rationalized their choice; recent switchers (bought in last 90 days or cancelled in last 90 days) still remember the situation that triggered the switch, which is the high-signal data
- Capture all four forces — a product roadmap that only addresses Push (customer frustrations) misses Anxiety; if customers are attracted to your product but not converting, usually the Anxiety force is unaddressed, not the Pull
- One job per analysis session — mixing multiple jobs in a single diagram produces contradictory insights; if interview data reveals two distinct jobs, run separate analyses

## Common Mistakes

- **Defining the job by product feature** — "the job of getting directions" is not a job; "the job of getting to an unfamiliar location confidently without looking incompetent" is; the feature-framed job produces incremental improvements; the progress-framed job shows substitution risk (GPS apps replacing printed maps, which replaced asking locals)
- **Conducting surveys instead of interviews** — JTBD data comes from narrative interviews about specific switch moments; surveys produce averages that obscure the qualitative texture of the job; "what features matter to you?" on a scale of 1–5 tells you nothing about why customers switch
- **Ignoring the social and emotional dimensions** — teams default to functional jobs because they're measurable; social and emotional dimensions are harder to quantify but are where differentiation lives; Apple's success with iPod was not functional (mp3 players already existed) — it was social and emotional
- **Confusing job with use case** — a use case is what the product does; a job is the progress the customer seeks; Milkshake example (Christensen): the milkshake's use case was "breakfast food"; the job was "make my long boring commute bearable and filling" — understanding the job led to a thicker shake that took longer to consume, not a healthier one
