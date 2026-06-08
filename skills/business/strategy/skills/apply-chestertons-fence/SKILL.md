---
name: apply-chestertons-fence
description: Use when taking over an inherited system, process, or organization — or when evaluating a proposed change to a working system — to establish whether you understand why each element exists before making any modifications
source: "资治通鉴 (Zīzhì Tōngjiàn) 汉纪 (~188 BC) — Cao Can inherits Xiao He's governance system and refuses all changes: "萧何与陛下定天下，法令既明...遵而勿失，不亦可乎"; Chesterton "The Thing" (1929) — the fence principle; Deming "Out of the Crisis" (1986) — tampering with a stable system always makes it worse"
tags: [change-management, organizational-design, inherited-systems, risk-management, governance]
verified: true
---

# Apply Chesterton's Fence

Before modifying any element of an inherited working system, establish why that element exists — because systems that work encode solutions to problems no longer visible, and removing elements without that understanding destroys the solution along with the evidence of the problem.

## Why This Is Best Practice

When Cao Can (曹参) became prime minister of Han (~188 BC), succeeding the legendary Xiao He, officials expected active reform. He changed nothing. When Emperor Hui pressed him on this "laziness," Cao Can asked two questions: Could the Emperor surpass Emperor Gaozu in wisdom? No. Could Cao Can surpass Xiao He? No. "Then," he said, "their laws and institutions are already clear. We have only to maintain them faithfully." 资治通鉴 records the era that followed as one of the most stable in Han history. Sima Guang cited this as evidence that maintaining a working system faithfully outperforms well-intentioned reform of a system you don't fully understand.

G.K. Chesterton articulated the same principle in "The Thing" (1929): "Don't ever take a fence down until you know the reason it was put up." The reformer who removes a fence without understanding it has not improved the system — they have removed a constraint whose purpose they couldn't see.

**Adopted by:** Amazon (two-way door vs. one-way door framework — one-way-door decisions, which are hard to reverse, require deep understanding of why the current state exists before reversal); Google, Netflix, Stripe (engineering principle: "don't refactor code you don't understand" — standard software engineering practice); McKinsey, BCG, Bain (diagnostic before solution is foundational methodology — understand current state thoroughly before prescribing change).

**Impact:** Deming "Out of the Crisis" (1986) mathematically demonstrated that "tampering" — adjusting a stable system without a model of why it works — consistently makes the system perform worse. His funnel experiment is a standard illustration in statistical process control education. Organizations that change working systems without understanding them routinely break what worked and then struggle to diagnose why, because the problem the system was solving is no longer visible once the solution is gone.

**Why best:** The instinct to change things upon arrival — to signal activity, to differentiate from a predecessor, or to apply prior-context knowledge — is nearly universal and nearly always costly in inherited contexts. The disciplined alternative is not passivity but diagnosis: build a model of why each element exists before touching any element.

**Distinct from `apply-premortem`:** apply-premortem surfaces failure modes of a NEW plan before commitment. apply-chestertons-fence applies to EXISTING working systems — the question is not "what could go wrong with our new plan" but "why does the current system work the way it does, and what would break if we changed this element?"

**Distinct from `apply-inversion`:** apply-inversion reasons backward from failure to identify risks in a plan. apply-chestertons-fence is a prerequisite discipline: before any change decision, establish whether you can explain the function of what you're proposing to change.

## Steps

1. **Identify what "working" means.** Before diagnosing the system, define what "working" looks like quantitatively. What outputs does the current system produce? At what rate? With what variance? This establishes the baseline you need to detect whether a change improved or degraded performance.

2. **Map the system before touching it.** List every major element — rule, process, structure, constraint, norm — that you are considering changing. For each element, write one sentence explaining why you believe it exists.

3. **Apply the fence test to each element.** For each element where you cannot write a complete explanation of why it exists: **do not change it yet.** The inability to explain why something is there is the signal that you don't yet understand it — not that it's unnecessary.

4. **Investigate unexplained elements.** Ask people who have worked with the system longest: "Why does this exist? What was the problem it was solving? What happened when it wasn't here?" If no one knows, look for historical incidents or failures that preceded the element's introduction. Often the element was introduced after a costly failure — and its function is to prevent that failure from recurring.

5. **Only change elements you can fully explain.** Prioritize changes where you can articulate: (a) why the current element exists, (b) what problem the change solves, (c) why the alternative is better than the current approach. For elements where you cannot articulate all three: wait.

6. **Change one element at a time in stable systems.** In a working system, changing multiple elements simultaneously makes it impossible to diagnose which change caused a subsequent problem. The discipline of single-variable change is especially important in inherited systems where your model is incomplete.

7. **Document the reason for each change.** Write down — before making the change — why the element existed, why you are changing it, and what you expect to happen. This record prevents future inheritors from treating your changes as Chesterton's fences they don't understand.

## Rules

- "I don't know why this exists" is a reason to investigate, not a reason to remove. An unexplained element has the same probability of being unnecessary as it does of encoding a critical constraint.
- "This looks old/inefficient/overcomplicated" is not a diagnosis. Old systems often encode hard-won solutions. Efficiency and cleanliness are properties you discover through understanding, not justify through appearance.
- Apply the fence test more strictly when taking over from a predecessor who was highly capable. Capable predecessors make fewer arbitrary decisions; unexplained elements from capable predecessors are more likely to encode real solutions.
- Greenfield contexts (building something new from scratch) do not require this discipline. The fence test applies only to inherited working systems — systems that are currently producing an outcome, even if that outcome could be improved.

## Examples

**New executive inheriting a team:** Incoming VP notices the team has a weekly all-hands meeting that seems redundant. Before canceling it, they attend for a month and discover it's the only forum where cross-functional blockers surface early — without it, issues accumulate until they become crises. The meeting survives; the format is improved instead.

**Acquiring a company:** Integration team decides to migrate the acquired company's customer data to the acquirer's CRM system. The acquired CRM has custom fields that seem redundant. Before removing them, the team investigates and finds those fields are used by the support team to flag accounts with unusual contractual terms that require special handling. Removing them would cause standard process to be applied to accounts where it would breach contract. The fields are preserved and migrated.

**Software engineering:** New engineer takes over a codebase with a seemingly redundant validation check on a field that appears to always be populated. They remove it as cleanup. In production, they discover that the check was protecting against a specific edge case in an upstream system that occasionally sends null values for that field — an edge case the original engineer had experienced and fixed. The check is restored after an incident.

**Organizational redesign:** Consulting team recommends removing a middle-management layer as an efficiency measure. Before implementing, a manager with 10 years of tenure explains that the layer was introduced after a product incident where field teams made contradictory commitments to the same customer. The layer's function is coordination of commitments. The redesign is restructured to preserve the coordination function through a different mechanism.

## Common Mistakes

- **Changing before mapping:** Making changes in the first 30 days of arrival, before building a model of why things work the way they do. This is the most common and costly failure pattern in leadership transitions.
- **Treating explanation difficulty as evidence of arbitrariness:** "I asked five people and none of them know why this exists" sometimes means the element is arbitrary — and sometimes means the person who knew has left, taking the institutional knowledge with them. Both are possible; neither is proof.
- **Applying this too broadly as a conservatism principle:** Chesterton's fence applies to working systems. Systems that are not working — producing consistent failures, declining outputs, clear dysfunction — do not benefit from the fence test; they require diagnosis and change. The fence test is not an argument for leaving broken things in place.
- **Changing everything you do understand:** Understanding why something exists does not obligate you to change it. The fence test establishes the minimum standard for change eligibility — it doesn't imply that everything eligible should be changed.
