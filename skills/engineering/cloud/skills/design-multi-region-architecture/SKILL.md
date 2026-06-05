---
name: design-multi-region-architecture
description: Use when designing systems that must survive regional outages, minimize global latency, or meet geographic data residency requirements
source: AWS Well-Architected Framework (Reliability Pillar); Google Cloud Architecture Framework; NIST SP 800-34 contingency planning
tags: [cloud, reliability, architecture, disaster-recovery]
verified: true
---

# Design Multi-Region Architecture

Design cloud systems that remain available and performant across geographic failures.

## Why This Is Best Practice

**Adopted by:** Netflix, Amazon, Google, Stripe, financial institutions under DORA/BCBS239 regulation
**Impact:** Achieve 99.99%+ availability; reduce latency by 40-60% for globally distributed users; meet RTO < 15 min targets
**Why best:** Single-region deployments create single points of failure; multi-region distributes blast radius and puts compute near users

Sources: AWS Well-Architected Reliability Pillar (2023); Google Cloud Architecture Framework; NIST SP 800-34 Rev.1

## Steps

1. **Define availability targets** — Set RTO/RPO per service tier (tier 1: RTO < 15 min, RPO < 1 min; tier 2: RTO < 4 hr, RPO < 1 hr). These drive region count and replication strategy.

2. **Select regions** — Choose primary, secondary, and optional tertiary regions based on: user latency (<100 ms rule), data residency law (GDPR, PDPA), cloud provider availability, and inter-region latency (prefer <50 ms pairs).

3. **Choose a multi-region pattern** — Active-active: traffic served from all regions simultaneously (highest complexity, lowest RTO). Active-passive: standby region takes over on failure (simpler, higher RTO). Pilot-light: minimal standby scaled on failover.

4. **Design data replication** — Synchronous replication for zero RPO (adds latency); asynchronous for low-latency writes with non-zero RPO. Use CRDTs or last-write-wins for conflict resolution. Partition databases by geography where data residency requires it.

5. **Implement global load balancing** — Route traffic via anycast DNS (Route 53, Cloud DNS) or global load balancers with health checks. Use latency-based or geolocation routing. Set failover thresholds and TTLs (<30 s).

6. **Handle consistency trade-offs** — Apply CAP theorem: during partition, choose availability or consistency per service. Use eventual consistency for non-critical data; strong consistency only where required (financial transactions).

7. **Design stateless compute** — Ensure application servers carry no local state. Externalize sessions to distributed caches (Redis Cluster, Memorystore). Store ephemeral files in object storage, not local disk.

8. **Implement observability per region** — Deploy independent monitoring stacks per region. Aggregate metrics globally but alert regionally. Measure inter-region replication lag as a key metric.

9. **Test failover regularly** — Run game days quarterly; automate chaos experiments (inject region failure in staging). Measure actual RTO vs target; update runbooks.

10. **Document runbooks** — Write step-by-step regional failover procedures. Assign ownership. Test runbooks annually and after every architecture change.

## Rules

- Never share a single database across regions without replication — it becomes your single point of failure.
- Always test failover before you need it; untested failover plans fail when it counts.
- Match region count to RTO/RPO requirements, not aspirations — each additional region multiplies operational complexity.
- Treat inter-region bandwidth costs as a first-class architectural concern.

## Common Mistakes

- **Assuming active-active is always better** — it requires conflict resolution logic; active-passive is often sufficient and far simpler.
- **Ignoring replication lag** — asynchronous replication means data loss on failover; measure and communicate RPO honestly.
- **Forgetting third-party dependencies** — DNS providers, CDNs, and SaaS tools may be single-region; audit all external dependencies.
- **Underestimating data residency complexity** — GDPR and similar laws may prohibit cross-border replication; consult legal before designing replication topology.

## When NOT to Use

- Internal tools or low-criticality workloads where downtime cost is less than multi-region overhead
- Early-stage products where simplicity and iteration speed outweigh availability guarantees
- Applications with regulatory requirements that prohibit data leaving a single jurisdiction
