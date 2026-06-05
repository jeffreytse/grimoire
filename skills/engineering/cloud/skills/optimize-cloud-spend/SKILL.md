---
name: optimize-cloud-spend
description: Use when cloud costs are growing faster than usage, when preparing for budget reviews, or when seeking to reduce infrastructure spend without degrading performance
source: FinOps Foundation "Cloud FinOps" (Storment & Fuller 2019); CNCF FinOps whitepaper; Gartner cloud optimization research (2022)
tags: [cloud, finops, cost-optimization, infrastructure]
verified: true
---

# Optimize Cloud Spend

Systematically identify and eliminate cloud waste to reduce costs 20-40% without compromising reliability or performance.

## Why This Is Best Practice

**Adopted by:** Lyft (reduced spend 50%), Spotify (30% reduction), Pinterest — documented FinOps case studies; FinOps Foundation standard methodology
**Impact:** Average organization wastes 32% of cloud spend (Gartner 2022); structured optimization programs achieve 20-40% reduction within 90 days
**Why best:** Cloud billing complexity obscures waste; systematic analysis reveals idle resources, over-provisioning, and commitment discount gaps that accumulate silently

Sources: Storment & Fuller "Cloud FinOps" O'Reilly (2019); FinOps Foundation State of FinOps (2023); Gartner "Optimize Cloud Spending" (2022)

## Steps

1. **Establish cost visibility** — Enable Cost Explorer (AWS) or Cloud Billing reports (GCP/Azure). Enforce a tagging policy (service, team, environment, cost-center). Set tag coverage target of ≥ 95%. Without attribution, optimization is guesswork.

2. **Identify idle and unattached resources** — Run automated scans for: EC2/GCE instances with CPU < 5% for 14+ days, unattached EBS/PD volumes, unused load balancers, idle NAT gateways, and orphaned snapshots older than 90 days. These are immediate zero-risk savings.

3. **Right-size compute** — Compare current instance type utilization (CPU p95, memory p95) over 30 days to instance specs. Downsize when p95 CPU < 40% and p95 memory < 60%. Use AWS Compute Optimizer or GCP Recommender to automate analysis. Right-sizing typically saves 20-30%.

4. **Apply commitment discounts** — Analyze on-demand spend for stable workloads (running > 80% of the time). Purchase Reserved Instances (1-year, no upfront) for 30-40% savings. Apply Savings Plans for flexibility across instance families. Commit only to observed baseline, not projected peak.

5. **Use Spot/Preemptible instances for fault-tolerant workloads** — Batch jobs, CI/CD workers, stateless web tier, ML training. Spot pricing is 60-90% below on-demand. Implement interruption handling; design workloads to checkpoint and resume.

6. **Optimize storage tiers** — Move infrequently accessed data to cold storage (S3 Glacier, GCS Nearline). Implement S3 Intelligent-Tiering for unpredictable access patterns. Delete logs and temporary files after retention period. Add lifecycle policies before launching new systems.

7. **Reduce data transfer costs** — Deploy CloudFront/Cloud CDN to cache frequently accessed content (reduces origin egress 60-80%). Co-locate services in the same region/AZ to eliminate inter-region transfer. Use VPC endpoints for S3/DynamoDB to avoid NAT gateway charges.

8. **Eliminate over-provisioned managed services** — Audit RDS instance sizes, ElastiCache node counts, and MSK broker sizes against actual utilization. Enable auto-scaling for DynamoDB, Aurora Serverless. Consider migrating low-traffic services to serverless.

9. **Implement a FinOps review cadence** — Weekly: review top 10 cost changes. Monthly: right-sizing sweep, commitment coverage review. Quarterly: architecture review for cost-optimization opportunities. Assign cost ownership to service teams, not centralized ops.

10. **Track unit economics** — Define cost per unit (cost per API call, cost per active user, cost per GB processed). Rising unit costs signal architectural inefficiency. Falling unit costs confirm optimization is working. Report to engineering leadership monthly.

## Rules

- Never optimize without measurement; establish baselines before making changes or you cannot quantify savings.
- Commit discounts before optimizing utilization — buying RIs for over-provisioned instances locks in waste for a year.
- Right-size in staging first; validate performance before applying to production.
- Treat idle resource cleanup as a recurring process, not a one-time event; new idle resources accumulate continuously.

## Common Mistakes

- **Optimizing infrequently** — a one-time cleanup without governance reverts to baseline waste within 6 months as teams provision new resources.
- **Purchasing Reserved Instances prematurely** — committing before right-sizing locks in over-provisioned capacity at a discount, reducing savings.
- **Ignoring data transfer in optimization** — compute optimization saves 20-30%; combined with data transfer optimization the saving is 35-50%.
- **No team-level cost accountability** — centralized cost teams cannot keep up with decentralized provisioning; embed FinOps into service team workflows.

## When NOT to Use

- Systems under active architectural redesign — optimize the new design, not the interim state
- Workloads with less than 3 months of stable usage data (insufficient baseline for commitment decisions)
- Crisis response: never optimize during an active incident; stability first, cost second
