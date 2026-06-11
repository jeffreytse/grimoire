---
name: apply-design-thinking
description: Use when solving a complex, human-centered problem where the right solution is unknown — to move through structured phases of deep user empathy, problem reframing, divergent ideation, and rapid prototype testing before committing to a solution.
source: 'IDEO, "The Field Guide to Human-Centered Design" (2015); Hasso Plattner Institute of Design at Stanford (d.school), "An Introduction to Design Thinking Process Guide" (2010); Brown, "Design Thinking", Harvard Business Review (2008); Cross, "Design Thinking: Understanding How Designers Think and Work", Berg (2011)'
tags: [design-thinking, human-centered-design, ideation, prototyping, empathy, problem-framing, innovation]
---

# Apply Design Thinking

Solve ill-defined problems through five phases — Empathize (deeply understand users), Define (reframe the problem from user insights), Ideate (generate divergent solutions before converging), Prototype (build fast low-fidelity versions), and Test (learn from real users) — ensuring you solve the right problem before optimizing the solution.

## Why This Is Best Practice

**Adopted by:** Stanford d.school teaches Design Thinking to 600+ students per year and has trained 50,000+ alumni; IDEO has applied it at Apple, Procter & Gamble, Mayo Clinic, and GE with documented outcomes. Adopted as a formal innovation methodology at IBM (trained 100,000+ employees 2012–2016, produced $18.6M ROI per dollar invested per IBM study), SAP, PepsiCo, and Airbnb (founders cite design thinking explicitly in company origin story). Taught in 700+ universities worldwide including Harvard, Yale, MIT, Cambridge, and INSEAD. The World Economic Forum lists it among the top 10 skills of the future.
**Impact:** IBM Global Design study (2018): products developed using design thinking came to market 2× faster with 33% fewer defects. IDEO's work on the first Apple mouse (1980) produced a product 1/10th the cost of Xerox's prototype using empathy-driven prototyping. Dam & Siang (Interaction Design Foundation 2020): organizations applying design thinking report 28% higher revenue growth and 35% higher return on investment vs industry benchmarks. The mechanism: 85% of product failures are failures of problem definition, not solution execution (Cooper, Edgett & Kleinschmidt 2004 product innovation research) — Design Thinking front-loads problem definition, preventing expensive late-stage pivots.
**Why best:** Traditional engineering approaches define requirements upfront, then execute. This works for well-defined problems. Design Thinking is optimized for ill-defined problems — where the right solution depends on understanding users in ways that can't be captured in initial requirements. The alternative (requirements → build → test) produces products that solve the stated requirement but miss the actual human need. Design Thinking explicitly includes a "reframe the problem" step (Define phase), which is absent from traditional product development but is what turns "better horse" thinking into automobile thinking.

Sources: IDEO Field Guide (2015); Stanford d.school Process Guide (2010); Brown (2008) HBR; IBM Design Thinking study (2018); Cooper et al. (2004) product innovation research

## Steps

### Phase 1 — Empathize (understand users before defining solutions)

Immerse yourself in the user's world before forming any solution hypotheses:

```
Empathy methods:
  Observation: watch users in their actual environment doing the actual task
               (not in a lab, not in a focus group)
  Interviews:  1-on-1 open-ended interviews focused on stories, not opinions
               "Tell me about the last time you..." not "What do you want?"
  Shadowing:   follow a user through their full day or workflow
  Immersion:   experience the problem yourself if possible
```

**What to capture:**
- What users do (behaviors, workarounds, tools they use)
- What users say (their language for the problem, their stated goals)
- What users think (frustrations, mental models, assumptions)
- What users feel (emotional reactions at each step)

**What to avoid:**
- Asking "what do you want?" — users describe solutions they know, not needs
- Starting with a solution in mind — suspend hypotheses during empathy phase
- Group interviews (people conform to each other's answers)

### Phase 2 — Define (reframe the problem from user insights)

Synthesize empathy data into a precise human-centered problem statement:

```
Point of View (POV) statement format:
  "[User] needs [need] because [insight]"

Examples:
  ✅ "A first-year nurse needs to quickly identify which patient needs
      immediate attention because she can't yet read subtle cues that
      experienced nurses recognize automatically"
  ✅ "A remote team manager needs to sense team morale without scheduled
      check-ins because async communication strips the emotional signals
      she relies on"

❌ Solution-embedded (wrong):
  "A nurse needs a better alert system" → this is a solution, not a need
  "A manager needs a Slack plugin" → solution, not need
```

**How Might We (HMW) questions:** Convert your POV statement into ideation prompts:
- "How might we help the nurse quickly triage patient urgency?"
- "How might we help the manager detect team emotional state passively?"

HMW questions are narrow enough to be useful (not "how might we fix healthcare") but broad enough to invite creative solutions.

### Phase 3 — Ideate (diverge before converging)

Generate as many solutions as possible before evaluating any of them:

```
Rules during ideation:
  ✅ Defer judgment — no "that won't work" during generation
  ✅ Build on others' ideas — "yes, and..." not "but..."
  ✅ Go for volume — aim for 50+ ideas before filtering
  ✅ Encourage wild ideas — the absurd often contains a useful seed
  ❌ No evaluating during generation — evaluation kills divergent thinking

Ideation techniques:
  Brainstorming: 30-minute time-boxed session, all ideas captured
  SCAMPER: Substitute, Combine, Adapt, Modify, Put to other use, Eliminate, Reverse
  Worst possible idea: generate deliberately terrible solutions, then invert them
  Analogies: "How would Amazon/Apple/a hospital solve this?"
```

After ideation, converge: vote on most promising ideas (dot voting), then select 1–3 concepts to prototype.

### Phase 4 — Prototype (build to think, not to ship)

Build the minimum artifact needed to test your assumptions:

```
Prototype fidelity by question type:
  Question: "Does the concept make sense?"
  → Paper prototype: sketches, sticky notes, hand-drawn wireframes (1–2 hours)

  Question: "Can users complete the core flow?"
  → Clickable mockup: Figma, InVision (half day)

  Question: "Does this physical form factor work?"
  → Cardboard/foam model (1 day)

  Question: "Will users pay for this?"
  → Landing page with waitlist or fake-door test (1 day)
```

**Prototype to learn, not to present:** A prototype that takes 3 weeks to build is not a prototype — it's a commitment. If you spend more than 1–2 days on a prototype before testing, you've over-invested in an untested assumption.

Mark every prototype as a test artifact: "This is not the product. We are testing [specific assumption]."

### Phase 5 — Test (learn, not validate)

Put the prototype in front of real users and observe:

```
Testing rules:
  5 users minimum: Nielsen's research shows 5 users reveal 85% of usability problems
  Observe, don't explain: when a user is confused, resist explaining; note the confusion
  Ask "What are you thinking?" not "Is this clear?" — open-ended, not leading
  Test the assumption, not the execution: "Does this concept solve the problem?"
                                          not "Is this design beautiful?"

After each test session:
  What surprised us?
  What did users do that we didn't expect?
  What assumption was confirmed or invalidated?
  What do we do next? (iterate, prototype differently, or go back to Define)
```

### Iterate — the loop is the process

Design Thinking is not a linear 5-step process. After testing:

```
Test revealed wrong problem definition? → Go back to Define
Test revealed users don't understand concept? → Iterate prototype
Test revealed the need doesn't exist? → Go back to Empathize
Test revealed the concept works? → Increase fidelity, test again
```

## Rules

- Complete Empathize before Define — skipping empathy and writing a problem statement from assumptions is the most common failure; the Define phase output is only as good as the Empathize phase input
- Diverge before converging — judgment during ideation kills ideas before they can be useful; enforce a strict "no evaluation" rule during brainstorming; only evaluate after the generation phase is closed
- Prototype low-fidelity first — high-fidelity prototypes make users reluctant to criticize; paper prototypes produce more honest feedback; start rough and increase fidelity only after the concept is validated
- Test with 5 real users — showing prototypes to colleagues or executives is not testing; they know too much about the product context; test with people who represent actual users

## Common Mistakes

- **Skipping empathy and starting with ideation** — the most common mistake; teams jump to "how might we solve X?" before verifying that X is the right problem; this produces creative solutions to the wrong problem
- **Using focus groups instead of observation** — groups produce socially acceptable answers, not honest behavior; observing one user in their actual environment is worth more than 20 people in a focus group
- **Building high-fidelity prototypes before testing the concept** — high fidelity before validation wastes weeks on an unproven assumption; if the concept doesn't work, the polish is irrelevant
- **Treating testing as validation, not learning** — testing designed to confirm a preferred solution produces confirmation bias; testing designed to find what's wrong with your assumptions produces actionable insights
