---
name: audit-error-budget
description: Use when reviewing error budget consumption, setting burn rate alerts, or making reliability investment decisions based on budget status
source: "Site Reliability Engineering (Google, 2016) Ch. 3; 'Alerting on SLOs' (Google SRE Workbook, 2018 Ch. 5); Alex Hidalgo 'Implementing SLOs' (2020)"
tags: [sre, error-budget, burn-rate, reliability, alerting, slo, operations]
verified: true
---

# Audit Error Budget

Analyze error budget consumption to decide whether to invest in reliability work or continue shipping features.

## Why This Is Best Practice

**Adopted by:** Google SRE teams; Spotify engineering; Atlassian reliability organization; formalized in Google's SRE Workbook as the standard method for balancing reliability against feature velocity

**Impact:** Google's original SRE book documents error budget policies eliminating the feature/reliability conflict by making trade-offs data-driven; organizations that implement error budget policies report reduced escalation frequency and faster incident resolution prioritization

**Why best:** Error budget is the mathematical expression of acceptable unreliability. Consuming it is not inherently bad — it means features shipped. Consuming it too fast is bad — it means users suffered unexpectedly. The audit identifies whether consumption is within the expected rate and triggers defined policy responses when it is not.

## Steps

1. **Calculate current budget remaining** — `budget_remaining = (current_error_rate - SLO_target) × window_length`; express in minutes of downtime equivalent and percentage of budget consumed
2. **Calculate burn rate** — `burn_rate = actual_error_rate / error_rate_allowed_by_SLO`; a burn rate of 1.0 is neutral; above 1.0 depletes budget; below 1.0 recovers it
3. **Classify burn rate severity** — fast burn (>14.4x over 1 hour): page immediately; slow burn (>1x over 3 days): ticket required; under 1x: healthy
4. **Identify consumption source** — attribute budget consumption to deployments, infrastructure changes, or traffic spikes; use deployment markers in your metrics system
5. **Apply the error budget policy** — if >50% consumed with >14 days remaining: reliability work enters the sprint; if 100% consumed: feature releases freeze until budget recovers
6. **Forecast end-of-window position** — project current burn rate to end of measurement window; if trajectory shows overrun, intervene now
7. **Document and present findings** — monthly error budget review should include: budget consumed, top three consuming incidents, policy actions taken, and trend vs. prior periods

## Rules

- An error budget review is mandatory input to sprint planning — reliability investment must be data-driven, not reactive
- Budget consumption attributed to planned maintenance must still count against the budget unless users were notified in advance
- Burn rate alerts must use multi-window detection (1-hour + 6-hour) to avoid both alert fatigue and missed slow burns
- Never reset the error budget manually mid-window to avoid policy consequences

## Examples

**Scenario:** SLO is 99.9% over 28 days. After 14 days, error rate is 0.3% (3× the allowed 0.1%). Burn rate = 3.0. Budget consumed: 60% with 50% of window remaining. **Policy trigger:** reliability work enters next sprint; no new feature releases until burn rate drops below 1.0.

## Common Mistakes

- Reviewing error budget only after incidents: by the time you review, the policy consequence has already been delayed
- Attributing all consumption to unavoidable causes: this removes the feedback loop that improves reliability
- Setting burn rate alerts that only trigger at 100% consumption: at that point, users have already suffered the full impact

## When NOT to Use

- No SLO has been defined for the service yet — error budget auditing requires a target SLO as the denominator; run `design-slo` first to establish the baseline before attempting an audit.
- The service is in alpha or private beta with fewer than 100 users — statistical noise in error rates at low traffic volumes makes burn rate calculations meaningless and will trigger false policy responses.
- The team is in the middle of a major incident with active user impact — auditing historical budget consumption during an active outage delays mitigation; complete incident response first and run the audit during the post-incident review.
