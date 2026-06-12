---
name: audit-process-efficiency
description: Use when identifying waste, bottlenecks, or inefficiencies in a repeatable business or operational workflow
source: Taiichi Ohno "Toyota Production System" (1978); Mike Rother "Learning to See" (1998) — Value Stream Mapping
tags: [process-efficiency, lean, value-stream-mapping, waste-reduction, operations, continuous-improvement]
verified: true
---

# Audit Process Efficiency

Systematically identify waste and bottlenecks in a workflow using lean value stream mapping techniques.

## Why This Is Best Practice

**Adopted by:** Toyota, Amazon operations, General Electric Six Sigma programs, NHS Lean transformation
**Impact:** Value stream mapping studies show 30–50% lead-time reduction in manufacturing and 20–40% in knowledge-work processes when waste is systematically eliminated (Lean Enterprise Institute, 2021)

**Why best:** Most process inefficiency is invisible to participants because they are inside it. A structured audit externalizes the process, quantifies non-value-adding time, and produces a prioritized improvement backlog rather than a list of complaints.

## Steps

1. **Define the scope** — Choose one process with a clear start event and end outcome. Example: "customer submits support ticket → ticket resolved and closed." Do not audit everything at once.
2. **Walk the process** — Follow one unit of work (a ticket, a document, an order) through every step from start to finish. Observe in person or via screen share; do not rely on documented procedures alone.
3. **Map the current state** — Draw each step as a box. Between boxes, record: wait time, queue depth, and who hands off to whom. Include decision points (if/else branches). Use sticky notes or a digital canvas.
4. **Classify each step** — Label every step: Value-Adding (VA) = customer would pay for it; Business Non-Value-Adding (BNVA) = necessary but wasteful (compliance, approvals); Non-Value-Adding (NVA) = pure waste to eliminate.
5. **Calculate the efficiency ratio** — Total VA time ÷ Total lead time. A ratio below 20% is common in knowledge work and signals major improvement opportunity.
6. **Identify the constraint** — Find the step with the longest cycle time or largest queue. Per Theory of Constraints, improving any step except the constraint does not improve throughput.
7. **Design the future state** — Redesign the process to eliminate NVA steps, reduce handoffs, and unblock the constraint. Aim for a 25–50% reduction in lead time. Document the future state map alongside the current state.
8. **Build the improvement backlog** — Convert each change into a concrete action item: owner, due date, success metric. Prioritize by impact on the constraint.

## Steps

## Rules

- Audit the actual process, not the documented process — they are almost always different.
- Quantify everything: time, queue depth, error rate, cost. Unquantified waste cannot be prioritized.
- Fix the constraint first; optimizing non-constraints is waste of improvement effort.
- Include the people doing the work in the redesign — they hold knowledge no audit will surface.
- Repeat the audit 90 days after implementing changes to verify improvement and find the new constraint.

## Examples

An HR team auditing new-employee onboarding discovered 14 days of the 21-day process were waiting time (background check vendor wait, IT provisioning queue, manager approval delays) — the actual work took 7 hours. Eliminating two approval hand-offs and pre-provisioning IT access cut lead time to 9 days without adding headcount.

## Common Mistakes

- **Auditing the org chart instead of the work** — Mapping departments instead of following the work unit produces a political map, not an efficiency map.
- **No quantification** — "This step feels slow" without cycle time data produces no defensible prioritization.
- **Improving everywhere equally** — Spreading improvement effort across all steps ignores the constraint; overall throughput does not improve.

## When NOT to Use

- Do not audit a process that runs fewer than ten times per month — the improvement ROI rarely justifies the mapping effort when volume is too low to amortize the change.
- Do not apply lean value stream mapping to creative or knowledge-discovery work (research, design exploration, strategy) where variability and iteration are the intended output, not waste.
- Do not conduct an efficiency audit immediately after a major system migration or reorganization; the process is in flux and the current-state map will be obsolete before improvements can be implemented.
