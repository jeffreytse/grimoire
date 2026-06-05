---
name: design-slo-sla-framework
description: Use when establishing reliability targets for services, defining error budgets for engineering teams, or negotiating service commitments with customers
source: Google SRE Workbook (2018) Ch. 2; Beyer et al. "Site Reliability Engineering" (2016); Treynor "The Calculus of Service Availability" Queue (2017)
tags: [reliability, sre, slo, observability]
verified: true
---

# Design SLO/SLA Framework

Define Service Level Objectives (SLOs) that balance reliability with feature velocity, and Service Level Agreements (SLAs) that set customer expectations with contractual consequences.

## Why This Is Best Practice

**Adopted by:** Google (SRE origin), Atlassian, Spotify, Cloudflare, AWS (every managed service has published SLAs) — the SRE discipline is built on this foundation
**Impact:** Google SRE book: teams using error budgets ship 2x more features with 2x fewer incidents; SLOs replace "keep the lights on" with measurable, negotiable reliability targets; error budgets make reliability vs. velocity trade-offs explicit and data-driven
**Why best:** Without SLOs, reliability is subjective; "the system is unreliable" vs. "we missed our 99.9% availability SLO by 0.2% this month" — the latter drives precise conversations and improvement

Sources: Beyer et al. "Site Reliability Engineering" O'Reilly (2016); Murphy et al. "The Site Reliability Workbook" O'Reilly (2018); Treynor "The Calculus of Service Availability" ACM Queue (2017)

## Steps

1. **Define SLIs (Service Level Indicators)** — An SLI is a quantitative measure of the service's behavior from the user's perspective. Common SLIs: availability (successful requests / total requests), latency (% of requests completing < threshold), throughput (requests processed per second), error rate (failed requests / total requests). Choose SLIs that reflect what users care about, not what's easy to measure.

2. **Set SLO targets** — An SLO is a target range for an SLI over a time window. "99.9% of requests succeed over a rolling 30-day window." The target must be aspirational but achievable. Start with historical performance and set the SLO at the 10th percentile of recent 30-day windows. Tighten over time as reliability improves.

3. **Calculate error budgets** — Error budget = 1 - SLO. For 99.9% availability: error budget = 0.1% of requests = 43.2 minutes of downtime per 30-day window. The error budget is the allowable unreliability. Track budget consumption daily. When the error budget is exhausted, reliability work takes priority over feature development.

4. **Define the SLO measurement window** — Rolling windows (last 30 days) are more stable than calendar windows and avoid the "reset on the 1st" gaming problem. Use a 28-day or 30-day rolling window for most SLOs. Weekly SLOs are useful for fast feedback but noisy; use as leading indicators, not commitments.

5. **Choose SLO tiers for service criticality** — Tier 1 (user-facing critical): 99.9% or 99.95% availability, p99 < 200 ms. Tier 2 (internal business logic): 99.5%, p99 < 500 ms. Tier 3 (batch, background): 99.0%, no latency SLO. Not all services need the same reliability target; higher SLOs cost proportionally more to maintain.

6. **Implement SLO monitoring** — Measure SLIs continuously. Use Prometheus + Grafana, Datadog SLO tracking, or Google Cloud Monitoring. Alert at 2% and 5% of monthly error budget consumed in the last hour (alert on burn rate, not remaining budget). Burn rate alerting provides time to act before the budget is exhausted.

7. **Define SLAs (Service Level Agreements)** — An SLA is a commercial agreement with defined consequences for missing the SLO (service credits, refunds). SLAs must be weaker than SLOs: if your SLO is 99.9%, your SLA should commit to 99.5% — the gap is your margin for SLO misses before SLA consequences trigger. Involve legal and product in SLA definition.

8. **Establish error budget policy** — Write and publish: when the error budget is > 50% remaining: normal development velocity. When 25-50% remaining: reliability review required before major deployments. When < 25% remaining: freeze feature deployments; focus on reliability. When exhausted: full reliability incident review; no feature work until budget is replenished.

9. **Conduct SLO reviews** — Monthly: review SLO compliance, error budget consumption, and trend. Quarterly: review whether SLO targets are appropriate — too easy targets provide false confidence; too hard targets paralyze teams. Annual: renegotiate SLAs based on current SLO track record.

10. **Make SLOs visible** — Publish SLO status on an internal dashboard accessible to all engineers and stakeholders. SLO transparency creates shared accountability. Consider publishing SLO status externally on a status page (Statuspage, Cachet) to set user expectations proactively.

## Rules

- SLOs must be set based on user needs, not technical limitations; "this is the best we can do" is not a user-centered SLO.
- The error budget policy must have organizational authority; a VP overriding the error budget freeze defeats the system.
- SLAs must always be weaker than SLOs; promising customers what you haven't achieved internally guarantees SLA violations.
- Burn rate alerting is mandatory; alerting only when the budget is nearly exhausted gives no time to remediate.

## Common Mistakes

- **Setting 100% SLO targets** — 100% is unachievable; it eliminates the error budget and forces all incidents to be crises; 99.9% is the standard for high-availability user-facing services.
- **Measuring from infrastructure metrics instead of user perspective** — server uptime ≠ user-experienced availability; measure from the user's perspective (successful API responses, successful page loads).
- **No error budget policy enforcement** — SLOs without consequences for consumption are aspirational metrics, not reliability management; the error budget policy gives SLOs teeth.
- **Too many SLOs** — 20+ SLOs per service creates alert fatigue; focus on 2-3 SLIs per service that most directly reflect user experience.

## When NOT to Use

- Internal developer tools with no external SLAs and where developers can tolerate occasional downtime
- Experimental or prototype services not yet serving production user traffic
- Batch systems where availability and latency SLOs are meaningless; use throughput and completion-time SLOs instead
