---
name: apply-fishbone-diagram
description: Use when investigating the root cause of a quality defect, recurring process failure, or operational problem — to map all contributing causes across multiple categories simultaneously before selecting which cause to address.
source: 'Kaoru Ishikawa, "Guide to Quality Control", JUSE Press (1968); ISO 9001:2015 Quality Management Systems; ASQ (American Society for Quality) Body of Knowledge; Juran & Godfrey, "Juran''s Quality Handbook", McGraw-Hill (1999)'
tags: [root-cause, quality, fishbone, ishikawa, problem-analysis, six-sigma, operations]
---

# Apply Fishbone Diagram

Map all potential causes of a problem across 6 standard categories (Machine, Method, Material, Man, Measurement, Environment) in a visual cause-and-effect diagram — surfacing the full landscape of contributing factors before selecting which root cause to investigate, preventing the premature single-cause fixation that makes problems recur.

## Why This Is Best Practice

**Adopted by:** Standard tool in ISO 9001 quality management systems (adopted by 1M+ certified organizations globally). Required technique in Six Sigma DMAIC (Analyze phase) — used at GE, Motorola, Toyota, Boeing, and virtually all Fortune 500 manufacturing and service operations quality programs. Ishikawa's original application was at Kawasaki Steel; the technique is now standard in healthcare quality (Joint Commission hospital accreditation), software engineering retrospectives, and food safety (HACCP protocols). The ASQ (American Society for Quality, 80,000+ members) includes it in the core quality tools body of knowledge.
**Impact:** Ishikawa (1968) documented application at 230+ Japanese factories where structured multi-cause analysis reduced defect recurrence rates by 40–60% compared to single-cause investigation. Toyota Production System documentation credits fishbone analysis as a key tool in achieving defect rates 10× lower than Western competitors in the 1980s (J.D. Power quality studies). The mechanism: single-cause investigation produces fixes that address symptoms; systematic cause mapping across 6 categories identifies systemic causes — which are 3–5× more common in recurring defects than proximate causes (Juran Institute data).
**Why best:** The 5 Whys (a complementary tool) drills deep on one causal chain; the fishbone maps wide across all possible cause categories. The two tools are used in sequence: fishbone first to identify which causal category to investigate, then 5 Whys to drill into that category. Using 5 Whys alone without a fishbone risks drilling deep on the wrong cause — a common failure mode in process improvement. The alternative (brainstorming without structure) produces unorganized lists where systemic causes get overlooked in favor of recent or visible events.

Sources: Ishikawa (1968) "Guide to Quality Control"; ISO 9001:2015; ASQ Body of Knowledge; Juran & Godfrey (1999) "Juran's Quality Handbook"

## Steps

### 1. Define the problem statement — the "head" of the fish

Write a precise, specific problem statement and place it at the right end of a horizontal arrow (the "spine"):

```
                                           ┌─────────────┐
─────────────────────────────────────────→ │   PROBLEM   │
                                           └─────────────┘

Good problem statement:
  ✅ "Defect rate on Assembly Line 3 increased from 0.8% to 3.2% in March"
  ✅ "Customer complaint rate for product X increased 40% in Q1"
  ✅ "Server response time exceeds 2s for 15% of requests since the v2.4 deploy"

❌ Too vague: "Quality is bad" / "System is slow" / "Customers are unhappy"
```

The specificity of the problem statement determines the quality of the analysis. Vague problem = vague causes = ineffective fixes.

### 2. Draw the 6 main cause branches (the "bones")

Draw 6 diagonal lines branching off the spine — the standard 6M categories:

```
         MACHINE          METHOD
            \               \
             \               \
──────────────────────────────────────→ [PROBLEM]
             /               /
            /               /
         MATERIAL          MAN
         
         (also: MEASUREMENT and ENVIRONMENT as additional branches above/below)
```

The 6M categories:
- **Machine** — equipment, tools, technology used in the process
- **Method** — procedures, protocols, standards, processes followed
- **Material** — inputs, raw materials, components, data
- **Man** (People) — human factors: training, skills, attention, communication
- **Measurement** — how outcomes are measured: accuracy, calibration, metrics
- **Environment** — physical conditions: temperature, workspace, external factors

For service/software contexts, adapt to 8P: People, Process, Policy, Procedure, Place, Product, Price, Promotion — or use the 6M as written, relabeling as needed.

### 3. Brainstorm causes for each category

For each of the 6 branches, ask the team: "How could [this category] contribute to [the problem]?"

```
Example — Machine branch for "server response time >2s":
  - Underpowered database server (CPU throttling)
  - Memory leak in application process
  - Network interface at 90% capacity
  - Disk I/O bottleneck from logging

Example — Method branch:
  - No query optimization in new endpoints
  - Synchronous calls where async would work
  - No caching layer for frequent reads
```

Write each cause as a sub-bone branching off the main category line. Add sub-sub-bones for causes of causes.

### 4. Identify the most likely root cause candidates

After populating all 6 branches:

1. **Vote** — ask the team to mark the 3 causes they believe most likely (dot voting or checkmarks)
2. **Evidence check** — for each high-vote cause, ask: "What data would confirm this?" Causes without verifiable evidence are hypotheses, not root causes
3. **Correlation check** — did the problem start when this cause appeared? (If the defect started in March, a cause present since January cannot be the root cause of the March increase)

```
Priority scoring:
  High votes + verifiable data + timing correlation = investigate first
  High votes but no timing correlation = contributing factor, not root cause
  Low votes = deprioritize unless no other candidates
```

### 5. Investigate top candidates — transition to 5 Whys

For each top candidate, apply `apply-five-whys` to drill into the causal chain:

```
Fishbone identified: "No query optimization in new endpoints"
5 Whys:
  Why? → New endpoints were shipped without performance review
  Why? → Performance review was not in the deployment checklist
  Why? → Checklist was created before performance requirements existed
  Why? → Performance SLAs were added after the initial process design
  Root cause → Process documents not updated when requirements changed
```

The fishbone tells you WHERE to look. The 5 Whys tells you WHY it happened.

### 6. Document and share

```
Document format:
  Problem statement: [specific, measurable]
  Fishbone diagram: [attached image or drawn]
  Top 3 root cause candidates: [listed with evidence]
  Selected root cause for action: [one specific cause]
  Action taken: [specific corrective action with owner and date]
  Verification: [how we'll confirm the cause is resolved]
```

## Rules

- One problem per diagram — fishbone for "all our quality problems" produces an unusable diagram; the problem statement must be specific enough that all 6 branches are relevant
- All 6 categories before voting — teams naturally gravitate to familiar categories (usually Method and Man); requiring all 6 forces consideration of Equipment and Measurement, which contain root causes in ~40% of manufacturing defects (ASQ data)
- Evidence required before acting — a cause with many votes but no data to verify it should be confirmed before resources are committed; the diagram is a hypothesis map, not a decision
- Use fishbone before 5 Whys, not instead of — fishbone identifies the right causal category; 5 Whys drills into it; using 5 Whys alone risks drilling into the wrong category

## Common Mistakes

- **Skipping the problem statement definition** — starting with causes before writing a specific problem statement produces a vague diagram where causes from multiple problems get mixed; write the problem statement first, verify it's specific
- **Using fishbone to confirm a pre-existing theory** — teams who already believe "it's a people problem" fill in Man with detail and leave other branches empty; this defeats the purpose; enforce equal brainstorming across all 6 categories
- **Treating the fishbone as the solution** — the diagram identifies candidates; acting on a cause without verifying it is confirmed produces fixes that feel productive but don't resolve the problem; always verify before implementing corrective action
- **Not using it with 5 Whys** — fishbone shows WHAT categories contain the cause; 5 Whys finds WHY the cause exists; using fishbone alone produces surface-level causes ("machine malfunction") without the systemic root ("maintenance schedule was defunded in Q3 budget cuts")
