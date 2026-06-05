---
name: calculate-cloud-costs
description: Use when estimating, forecasting, or auditing cloud infrastructure costs before provisioning, during architecture review, or at budget planning time
source: FinOps Foundation "Cloud FinOps" (Storment & Fuller 2019); AWS Cost Explorer methodology; Google Cloud pricing calculator best practices
tags: [cloud, finops, cost, budgeting]
verified: true
---

# Calculate Cloud Costs

Produce accurate, defensible cloud cost estimates and forecasts to inform architecture decisions and budget planning.

## Why This Is Best Practice

**Adopted by:** FinOps Foundation member companies (Atlassian, Spotify, Nike, Goldman Sachs); required by enterprise procurement and budget governance processes
**Impact:** Organizations using structured cost estimation reduce budget overruns by 40% (FinOps Foundation 2023 State of FinOps); unchecked cloud spend grows 23% YoY without governance
**Why best:** Cloud pricing is complex (compute, storage, network egress, API calls, support tiers all compound); structured calculation prevents surprise bills and enables architecture trade-off decisions

Sources: Storment & Fuller "Cloud FinOps" O'Reilly (2019); FinOps Foundation State of FinOps Report (2023); AWS Well-Architected Cost Optimization Pillar

## Steps

1. **Inventory all billable components** — List every resource category your architecture requires: compute (instance type, count, hours/month), storage (GB, IOPS, class), network (egress GB, inter-region, CDN), managed services (RDS, Kafka, ML APIs), and support plan. Missing a category is the most common source of estimate error.

2. **Gather usage baselines** — For existing systems: pull 90-day actuals from Cost Explorer / Cloud Billing. For new systems: estimate peak and average RPS, data volumes, and user counts. Never estimate on peak alone; use p95 load for sizing.

3. **Use provider pricing calculators** — AWS Pricing Calculator, GCP Pricing Calculator, Azure Calculator. Input specific SKUs, not generic "server" estimates. Capture on-demand prices first, then apply commitment discounts separately.

4. **Apply commitment discounts** — Reserved Instances (AWS) or Committed Use Discounts (GCP) reduce compute costs 30-60% for stable workloads. Savings Plans (AWS) offer 20-66% off on-demand. Only apply commitments to stable baseline; keep variable workloads on-demand or spot.

5. **Model data transfer costs** — Egress costs are the most underestimated line item. Model: inter-region replication, CDN origin pull, API gateway egress, and end-user downloads. AWS charges $0.09/GB egress to internet; at scale this dominates.

6. **Estimate managed service overhead** — RDS doubles the cost of equivalent EC2+storage. EKS, MSK, ElastiCache, Elasticsearch Service add 20-40% overhead vs self-managed. Quantify the operational time saved vs cost delta to justify managed services.

7. **Build a cost model spreadsheet** — Structure: resource → unit price → quantity → monthly cost → annual cost. Separate fixed costs (reserved capacity) from variable (per-request, egress). Add a contingency buffer (15-20%) for usage variance.

8. **Allocate costs by service/team** — Implement tagging strategy (team, environment, service, cost-center) before deployment. Tag coverage < 80% makes chargeback impossible. Define tagging enforcement via SCPs (AWS) or Organization Policies (GCP).

9. **Set budgets and alerts** — Configure billing alerts at 50%, 80%, and 100% of budget. Use AWS Budgets or GCP Budget alerts. Alert the service owner, not just finance. Treat a 100% alert as an incident.

10. **Reconcile monthly against forecast** — Compare actuals to estimate. Investigate variances > 10%. Update model with actual usage data. Cost estimation accuracy improves only through feedback loops.

## Rules

- Always model data egress separately — it is never included in compute estimates and routinely exceeds $10K/month at scale.
- Commitment discounts require 1- or 3-year lock-in; only commit to stable baseline capacity, never to anticipated peak.
- Tag everything before launch; retroactive tagging is 10x harder and leaves gaps in cost allocation forever.

## Common Mistakes

- **Using list price for compute** — failing to apply Reserved Instance or Savings Plan discounts overstates costs by 30-60%; stakeholders reject over-priced architectures.
- **Ignoring support plan costs** — AWS Business Support is 10% of monthly bill (minimum $100); Enterprise Support is 3%; at $100K/month this is $3,000/month overlooked.
- **Estimating on peak traffic only** — compute sized for peak but billed for actual; misrepresents cost if auto-scaling is in place.
- **Forgetting snapshot and backup storage** — daily snapshots of a 10 TB database accumulate to 300 TB of snapshot data per month if not lifecycle-managed.

## When NOT to Use

- Proof-of-concept workloads running for < 30 days where rough order of magnitude is sufficient
- Cost calculations for systems where pricing is entirely fixed (negotiated enterprise agreements with flat fees)
