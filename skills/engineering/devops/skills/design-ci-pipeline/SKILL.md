---
name: design-ci-pipeline
description: Use when designing or improving a continuous integration pipeline for a software project
source: "Accelerate" — Forsgren, Humble & Kim (IT Revolution, 2018); DORA metrics (dora.dev); Trunk-Based Development (trunkbaseddevelopment.com)
tags: [ci, devops, dora, pipeline, automation, trunk-based-development]
verified: true
---

# Design CI Pipeline

Design a CI pipeline that gives fast, reliable feedback on every commit and enables trunk-based development at scale.

## Why This Is Best Practice

**Adopted by:** Google, Netflix, Amazon, Etsy — all DORA elite performers with documented CI/CD practices
**Impact:** DORA research (Accelerate, 2018) shows elite performers have 46× more frequent deployments and 2,555× faster lead time to change; fast CI (<10 min) is a key differentiator between elite and low performers.

**Why best:** A well-designed CI pipeline is the heartbeat of software delivery. It enforces quality gates, provides rapid feedback, and is the foundation that makes trunk-based development safe at scale. Slow or flaky CI is one of the highest-leverage engineering problems to fix.

## Steps

1. **Define pipeline stages** — Minimum viable: checkout → install → lint → test → build → security scan. Add: integration tests, performance tests, deploy-to-staging for pre-merge validation.
2. **Target <10 minutes to green** — Parallelize test stages; use test sharding; cache dependencies aggressively. DORA data: pipelines >10 min cut developer flow significantly.
3. **Fail fast** — Run the fastest checks first (lint, unit tests); run slower checks (integration, E2E) in parallel or after fast checks pass.
4. **Make the pipeline the gate** — No merge to main without green CI. Enforce via branch protection; never allow bypassing.
5. **Eliminate flakiness** — Track flaky tests (retry logic with reporting); quarantine and fix flaky tests within one sprint. Flaky CI erodes trust and causes developers to ignore failures.
6. **Cache and artifact management** — Cache: npm/pip/Maven deps, Docker build layers. Share build artifacts between stages; do not rebuild the same binary twice.
7. **Instrument the pipeline** — Measure: mean time to feedback, flaky test rate, pipeline duration trend. Alert when median duration exceeds SLO.

## Rules

- Every commit to main must trigger the full pipeline.
- Secrets must never be echoed in pipeline logs; mask all secrets.
- Pipelines must be defined as code (YAML) in the repository — no click-ops configuration.
- Differentiate PR pipelines (fast feedback) from post-merge pipelines (full validation + deploy).

## Examples

GitHub Actions structure:
```yaml
jobs:
  lint:      # 1 min
  unit-test: # 3 min, parallel matrix [node18, node20]
  build:     # 2 min, depends on unit-test
  security:  # 2 min, parallel with build
  integration: # 5 min, depends on build
```
Total wall time: ~8 minutes via parallelism.

## Common Mistakes

- **Sequential stages** — lint → test → build → scan run one-at-a-time; doubles pipeline time versus parallel.
- **No dependency caching** — `npm install` from scratch on every run adds 2-3 minutes; cache the node_modules layer.
- **Treating flaky tests as acceptable** — a 5% flaky rate means 1 in 20 PRs gets a false failure; developers start ignoring CI.
