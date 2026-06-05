---
name: calculate-capacity-plan
description: Use when provisioning infrastructure for a new system, preparing for anticipated traffic growth, or ensuring headroom before a major product launch
source: Google SRE Workbook capacity planning chapter; Allspaw "The Art of Capacity Planning" (O'Reilly 2008); NIST cloud scalability guidelines
tags: [performance, capacity, infrastructure, reliability]
verified: true
---

# Calculate Capacity Plan

Model current resource utilization and traffic growth to determine when and how much additional infrastructure to provision to maintain SLO compliance.

## Why This Is Best Practice

**Adopted by:** Google SRE (capacity planning is a core SRE function), Netflix (capacity headroom model), AWS (Auto Scaling Group sizing methodology)
**Impact:** Allspaw (2008): organizations without capacity planning experience 2-3 avoidable capacity-related outages per year; over-provisioning without capacity plans wastes 30-40% of cloud spend; under-provisioning causes SLO violations
**Why best:** Infrastructure takes time to provision; without a forward-looking model, you react to capacity issues after they cause outages rather than before

Sources: Allspaw "The Art of Capacity Planning" O'Reilly (2008); Murphy et al. "Site Reliability Workbook" O'Reilly (2018) Ch. 17; Google SRE capacity planning practices

## Steps

1. **Measure current resource utilization** — Collect 90-day p50 and p95 utilization for each resource dimension: CPU, memory, disk I/O, network bandwidth, database connections, and cache hit rate. Use percentiles, not averages — averages mask peaks. P95 utilization is your effective utilization baseline.

2. **Identify the bottleneck resource** — From your utilization data, identify which resource reaches saturation first as traffic increases. Common bottlenecks: CPU (compute-bound services), database connections (connection pool exhaustion), memory (high-memory workloads), I/O (storage-bound services). All capacity modeling focuses on the bottleneck resource.

3. **Establish a traffic growth model** — Analyze 90-day traffic trends (requests per second, active users, data volume). Fit a growth model: linear (constant traffic addition), exponential (percentage growth per period), or step-function (event-driven spikes). Get business context: planned marketing campaigns, product launches, seasonal peaks. Your growth model drives the capacity timeline.

4. **Calculate resource-to-traffic ratio** — Determine: at current traffic, what resource utilization exists? Example: 1000 RPS → 40% CPU on 8 vCPUs. Compute the ratio: 1000 RPS / 40% = 2500 RPS before CPU saturation. Apply a safety buffer: plan capacity for 70% max utilization (30% headroom for traffic spikes and failover).

5. **Project when capacity will be exhausted** — Apply the growth model to the safe capacity threshold. Example: 2500 RPS safe limit, current 1000 RPS, growing 15% per month. Months to exhaustion = log(2500/1000) / log(1.15) ≈ 6.5 months. Subtract your infrastructure provisioning lead time (cloud: 0 days for on-demand, 2-4 weeks for reserved procurement).

6. **Model for N+1 and multi-region redundancy** — Capacity must support N+1 node failure. If you have 3 web servers and one fails, the remaining 2 must handle 100% of traffic. Plan capacity at N+1 utilization: 2/3 of normal utilization per server (67%), so 3 servers run at 67%, not 90%. Multi-region: each region must handle full traffic load if the peer region fails.

7. **Account for efficiency initiatives** — If planned optimizations will reduce resource consumption (caching, query optimization, algorithm improvement), model their impact. Example: adding read replicas expected to reduce primary DB load by 40%. Apply these reductions conservatively (50% of the planned gain) to account for implementation delays.

8. **Calculate the procurement/provisioning timeline** — Determine lead times: on-demand cloud instances (minutes), reserved instances (1-3 days to activate), physical hardware (6-12 weeks), network capacity increases (2-4 weeks for ISP circuits). The capacity plan output is: "We need to provision X resources by date Y or we will breach our SLO."

9. **Build a capacity model spreadsheet** — Columns: date, traffic (RPS/users), resource utilization (by dimension), % of safe capacity, headroom remaining. Rows: monthly projections for 12 months. Include scenario analysis: baseline, +50% growth spike, major launch event. Share with engineering leadership monthly.

10. **Review and update quarterly** — Traffic growth rates change; efficiency initiatives complete or slip; product roadmaps shift. Review the capacity model quarterly against actuals. Adjust the growth model based on observed trends. Capacity plans that aren't updated against actuals drift from reality within 3 months.

## Rules

- Never plan for > 70% utilization on any bottleneck resource; 30% headroom is the minimum for safe operations (handles 1.4x traffic spikes without degradation).
- Always include N+1 redundancy in capacity calculations; single-server headroom calculations underestimate real requirements by 50%+ for multi-node services.
- Capacity plans must include the provisioning timeline as an output; "we'll need more capacity" without "by this date" is not actionable.
- Model for the business scenario with the highest traffic impact (major launch, seasonal peak) not the average scenario.

## Common Mistakes

- **Using average utilization** — averages hide peaks; p95 utilization is the correct baseline for capacity decisions.
- **No growth model** — "current capacity is sufficient" without a growth timeline misses the window to provision before SLO breach.
- **Forgetting stateful resource limits** — database connection limits, in-memory cache size, and file descriptor limits are hard caps that don't auto-scale; model them explicitly.
- **Not accounting for failover** — a system with 95% CPU utilization on three nodes has no capacity to absorb a node failure; always plan for N+1.

## When NOT to Use

- Serverless-first architectures where capacity is managed by the provider and auto-scales without provisioning decisions (though cost planning is still required)
- Prototype systems with no SLO and no production traffic
- Systems with unpredictable traffic spikes where horizontal auto-scaling is the capacity strategy (still need to set auto-scaling bounds correctly)
