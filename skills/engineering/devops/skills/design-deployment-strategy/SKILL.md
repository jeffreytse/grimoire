---
name: design-deployment-strategy
description: Use when selecting or designing a deployment strategy for releasing software to production safely
source: AWS Well-Architected Framework (Reliability Pillar); Netflix Tech Blog (canary deployments); Martin Fowler "BlueGreenDeployment" (martinfowler.com)
tags: [deployment, blue-green, canary, rolling, devops, reliability, aws]
verified: true
---

# Design Deployment Strategy

Select and implement the deployment strategy that best matches the system's risk tolerance, rollback requirements, and infrastructure capabilities.

## Why This Is Best Practice

**Adopted by:** Netflix (canary), AWS (blue-green with CodeDeploy), Google (gradual rollouts in GKE), Facebook (progressive push)
**Impact:** Netflix's canary deployments catch ~95% of production issues before they affect all users; blue-green deployments reduce mean time to recover (MTTR) from hours to minutes via instant rollback.

**Why best:** The choice of deployment strategy directly determines blast radius when something goes wrong. A rolling deploy with no traffic control can expose 100% of users to a bad release in minutes; canary deploys can limit exposure to 1% while metrics are evaluated.

## Steps

1. **Assess risk and rollback requirements** — High-risk changes (DB migrations, major refactors): blue-green or canary. Low-risk patches: rolling. Zero-downtime requirement: any strategy except in-place restart.
2. **Evaluate infrastructure capability** — Blue-green requires double capacity temporarily. Canary requires traffic splitting (service mesh, load balancer weights, or feature flags). Rolling requires health check support.
3. **Design blue-green for instant rollback** — Maintain two identical production environments; route traffic via DNS/load balancer switch; keep blue alive for 24h after green proves healthy.
4. **Design canary for gradual exposure** — Route 1-5% of traffic to new version; monitor error rate, latency, and business metrics for a defined window; auto-rollback if thresholds breach.
5. **Design rolling for stateless services** — Replace instances in batches (e.g., 25% at a time); configure health checks to gate each batch; set `maxUnavailable=0` for zero-downtime.
6. **Automate rollback triggers** — Define rollback criteria (error rate >1%, P99 latency >2×, specific log patterns) and automate rollback execution; do not rely on manual intervention.
7. **Test the rollback path** — Practice rollback in staging; a rollback procedure that has never been tested will fail when needed.

## Rules

- Never deploy and walk away — monitor for at least 15 minutes post-deploy before declaring success.
- Database migrations must be backward-compatible with the current production code before deployment begins.
- Feature flags decouple deployment from release — deploy dark, release with a flag flip.
- Document the rollback procedure in a runbook; link it from the deployment pipeline.

## Examples

Canary with Kubernetes and Argo Rollouts:
- Deploy v2 to 5% of pods; Argo monitors Prometheus error rate.
- If error rate <0.5% for 10 minutes → promote to 50% → then 100%.
- If error rate >1% at any step → automatic rollback to v1.

## Common Mistakes

- **No rollback plan** — "we'll just redeploy" is not a plan; migration side effects may not be reversible.
- **Canary without meaningful metrics** — routing 5% of traffic but not measuring business-specific error rates makes the canary blind.
- **Blue-green without database schema alignment** — running two app versions against a single DB with incompatible schemas causes immediate data corruption.

## When NOT to Use

- The service is a stateful database engine or distributed storage cluster — deployment strategies designed for stateless application tiers do not apply; database version upgrades follow vendor-specific rolling upgrade procedures that must not be overridden with generic canary or blue-green patterns.
- The deployment is a hotfix for a critical security vulnerability actively being exploited in production — gradual canary exposure intentionally delays full rollout, which prolongs user exposure to the vulnerability; deploy to 100% immediately with monitoring and accept the higher blast radius.
- The infrastructure has no health check mechanism and no traffic splitting capability — canary and rolling strategies both depend on automated health signal to gate promotion; without these primitives, the strategy cannot be executed safely and infrastructure prerequisites must be built first.
