---
name: apply-early-intervention
description: Use when you detect early warning signals — technical debt accumulating, customer relationship cooling, team friction forming, process degrading — to intervene while the problem is still small and cheap, because intervention cost increases nonlinearly with delay
source: "道德经 第六十三章 (Laozi, ~6th–4th century BC) — 图难于其易，为大于其细；天下难事，必作于易；天下大事，必作于细; IBM Systems Sciences Institute defect cost research (1981) — defect cost in production is 100× defect cost in requirements; DevSecOps shift-left principle; Toyota jidoka (自働化) — stop-and-fix-immediately; Deming PDCA cycle; Total Productive Maintenance (TPM)"
tags: [strategy, risk-management, technical-debt, preventive-maintenance, early-warning]
verified: true
---

# Apply Early Intervention

Tackle difficult problems while they are still small — because the cost to resolve a problem increases nonlinearly with delay, and the signals that make early intervention possible are available long before the crisis makes it unavoidable.

## Why This Is Best Practice

道德经 Chapter 63 (Laozi, ~6th–4th century BC):

> 图难于其易，为大于其细。天下难事，必作于易；天下大事，必作于细。

"Tackle the difficult while it is still easy; accomplish the great while it is still small. All difficult things in the world begin as easy ones; all great things begin as small ones."

**Why best:** Laozi's observation is causal, not motivational: the window in which a problem is easy to solve exists before it is visible as a crisis. Once a problem becomes obviously critical, the easy solution is no longer available — only expensive, disruptive, high-risk interventions remain. The principle is not "fix problems early because it's virtuous" but "fix them early because the cheap path closes as the problem grows."

**IBM Systems Sciences Institute defect research (1981):** The foundational empirical study of intervention cost vs. delay in software development. Finding: the cost to fix a defect identified in the requirements phase is $1; the same defect identified in production costs $100. The nonlinear cost curve — not a linear increase — is the key finding. A 100× cost multiplier is not recoverable through better late-stage execution; it changes the economics of intervention entirely. Replicated in multiple studies; forms the foundation of the "cost of quality" framework used at IBM, Boeing, and throughout software and manufacturing.

**DevSecOps "shift left" (standard practice, 2010s–present):** The industry-wide movement to move security and quality testing earlier in the software development pipeline. The finding: security vulnerabilities found during code review cost 6× less to fix than those found in testing; vulnerabilities found in production cost 30× more than those found in development. Adopted as standard practice at Google, Netflix, Stripe, Amazon, and throughout the FAANG-tier engineering ecosystem. NIST Software Security in Supply Chains report (2022) explicitly mandates shift-left as a supply chain security control.

**Toyota jidoka (自働化, "autonomation"):** One of the two pillars of the Toyota Production System. When a defect is detected, the machine or worker stops the line immediately to fix it — rather than passing the defect downstream. Toyota's principle: a defect caught at its source takes minutes to fix; the same defect, passed through 10 downstream processes, takes days to diagnose and fix. Stopping the line to fix immediately is faster and cheaper than processing to completion and fixing at end-of-line. Adopted globally in Lean manufacturing; standard in automotive, electronics, and high-precision manufacturing.

**Deming PDCA cycle (W. Edwards Deming, "Out of the Crisis", 1982):** Plan-Do-Check-Act as a continuous improvement cycle, with the emphasis that the "Check" phase identifies problems while they are still contained within the cycle — before they compound across cycles. Deming's explicit argument: problems that cross the boundary between cycles become structural; problems caught within a cycle remain solvable. Adopted globally; standard in ISO 9001 quality management systems, healthcare (Institute for Healthcare Improvement), and manufacturing.

**Total Productive Maintenance (TPM, Japan Institute of Plant Maintenance, 1971):** The framework for preventive equipment maintenance. Core principle: scheduled small maintenance interventions prevent catastrophic equipment failures. The economics: a preventive maintenance inspection costs 1/10th the cost of an emergency breakdown repair, not counting production downtime. Standard in aerospace (Boeing, Airbus), automotive (Toyota, BMW), and energy (GE, Siemens).

**Why distinct from `apply-pdca`:** apply-pdca is a cycle for continuous improvement of known processes. apply-early-intervention addresses a different trigger: detecting and acting on early warning signals BEFORE a problem becomes a crisis, in any domain, using cost-of-delay reasoning rather than improvement cycles.

**Adopted by:** IBM (1981 defect cost research, foundational to cost-of-quality frameworks used at IBM and Boeing); Google, Netflix, Stripe, Amazon (DevSecOps shift-left as standard engineering practice); Toyota (jidoka — stop-and-fix-immediately, one of the two pillars of the Toyota Production System); Boeing and Airbus (Total Productive Maintenance, TPM); ISO 9001 quality management systems globally (Deming PDCA); NIST Software Security in Supply Chains report (2022, mandating shift-left as a supply chain security control).

**Impact:** IBM's research established a 100× cost multiplier for production-stage defect correction versus requirements-stage correction — not a linear increase but a fundamental change in the economics of intervention. DevSecOps shift-left adoption at major engineering organizations reduced security vulnerability remediation costs by 6–30× depending on detection stage. Toyota's jidoka principle enabled a defect rate 1/10th that of American automakers and inventory 1/3rd their levels by the 1980s, validated across the global Lean manufacturing adoption that followed.

## Steps

1. **Identify early warning signal categories for the current domain.** Problems have predictable signal sequences: weak early signals → stronger signals → obvious crisis. Map the signal sequence for your context:
   - Technical: test coverage declining, build times increasing, error rates creeping, architectural inconsistency accumulating
   - Organizational: team velocity declining, 1:1 feedback changing tone, attrition signals in language ("when I leave")
   - Customer: support ticket volume increasing, NPS declining before churn, usage patterns shifting
   - Financial: margin compression, CAC increasing, LTV declining
   The goal is to catch signals at the first stage, not the second or third.

2. **Establish monitoring for early signals.** Early signals are invisible without explicit monitoring. Design surveillance for the signal categories you identified:
   - Automated: CI/CD metrics, error rate dashboards, business KPI alerts
   - Manual: regular sampling (reading 10 support tickets per week, 1:1 sentiment tracking)
   - Process: scheduled debt audits, architecture reviews, customer health checks
   The monitoring cadence should be faster than the problem compounds. If a problem doubles in severity weekly, monthly monitoring is too slow.

3. **Quantify the cost of delay before each identified problem.** For each early-stage problem detected, estimate the cost to fix now vs. cost to fix in 30 days, 90 days, 1 year. Use the nonlinear cost model: most problems grow faster than linearly because they compound (technical debt makes future changes slower; team friction reduces throughput; customer dissatisfaction spreads). The cost-of-delay calculation provides the business case for early intervention over deferred action.

4. **Intervene at the smallest effective unit.** Early intervention should be proportionate to the current problem size, not the projected future size. A 10-line architectural inconsistency warrants a 30-minute refactor, not a 2-week architecture overhaul. Over-engineering the intervention at the early stage wastes resources; the principle is "smallest effective intervention at earliest detection."

5. **Document the intervention and its trigger signal.** Record: what signal was detected, when, at what threshold, and what intervention was applied. This creates a learning record that improves signal detection for future cycles. Over 3–6 months, the intervention log reveals which signal categories are highest value to monitor and which interventions are most cost-effective.

6. **Set the threshold below "obvious problem."** The most common failure mode is waiting until the problem is undeniable before acting. Define the intervention threshold explicitly: "If X metric exceeds Y, or if I observe Z signal, I will intervene in the next sprint/week/cycle." The threshold should be uncomfortable — it should feel slightly premature, because the point is to act before it feels obviously necessary.

## Rules

- Intervention cost is nonlinear, not linear. A problem that costs 1 to fix today does not cost 2 in a week — it may cost 10 or 100, depending on the domain. When evaluating whether to intervene now or defer, apply nonlinear rather than linear cost projection.
- Early signals are ambiguous by design. The fact that an early signal might not indicate a real problem is not a reason to defer. The cost of investigating a false positive is small; the cost of ignoring a true positive compounds.
- Do not confuse small problems with unimportant problems. A small problem is not necessarily low-priority — it is low-cost to fix now and high-cost to fix later. Priority is a function of both urgency and cost-of-delay.
- Schedule intervention capacity, not just response capacity. If your process only has capacity for reactive crisis response, early intervention is impossible regardless of how well you detect signals. Reserve a fraction of each cycle for proactive intervention (standard in engineering: "20% capacity for tech debt").

## Examples

**Technical debt:** Static analysis shows 12 functions with cyclomatic complexity above threshold. Current impact: none visible. Cost now: 4 hours of refactoring. Cost in 6 months (after 3 features built on top): 3-day refactor plus regression testing. Early intervention: schedule refactor this sprint.

**Customer relationship:** A key account's monthly active user count dropped from 400 to 310 over 6 weeks. No support tickets filed; no explicit complaint. Early signal: usage decline precedes churn by 3–6 months in this product's data. Intervention now: proactive customer success outreach, health check call. Intervention after churn notice: contract renegotiation under pressure, 6-month sales cycle to replace.

**Team friction:** Two engineers are producing inconsistent code review comments for each other — reviews are slower than average, tone is more formal. No explicit conflict. Early signal: review friction precedes team split or attrition in 60–90 days. Intervention: manager has a direct conversation, clarifies expectations, mediates the underlying disagreement. Intervention after attrition: recruiting costs, 3-month ramp-up for replacement.

**Security vulnerability:** Dependency audit shows a library two versions behind with a known medium-severity CVE. No active exploit in the wild. Cost now: 2 hours to update dependency, run tests. Cost if exploited in production: incident response, customer notification, potential regulatory reporting, reputational damage. Intervention: update scheduled immediately.

## Common Mistakes

- **Waiting for the problem to "prove itself":** The instinct to wait until a signal is confirmed as a real problem delays intervention until the cheap window has closed. Early signals are by definition ambiguous — that is the nature of early-stage problems.
- **Applying future-size interventions to present-size problems:** Detecting an early signal and responding with a large intervention designed for the fully-developed crisis wastes resources and creates organizational fatigue. Match intervention size to current problem size.
- **Monitoring lagging indicators only:** Revenue decline, customer churn, and production incidents are lagging indicators — they confirm the problem has fully developed. Monitor leading indicators (usage trends, error rate trends, team velocity trends) that precede the lagging outcomes by weeks or months.
- **No reserved capacity for proactive work:** If every cycle is fully loaded with reactive work, early intervention competes with crisis response and loses. The fix is structural: reserve a dedicated fraction of each cycle for proactive work before problems are visible.
