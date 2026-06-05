---
name: run-chaos-engineering
description: Use when validating system reliability before a major launch, when SLO breach investigations reveal unknown failure modes, or when establishing a proactive reliability practice
source: Rosenthal et al. "Chaos Engineering" (O'Reilly 2020); Nygard "Release It!" (2018); Netflix Chaos Monkey (Basiri et al. 2016)
tags: [reliability, chaos-engineering, sre, testing]
verified: true
---

# Run Chaos Engineering

Deliberately inject failures into a system in a controlled way to discover weaknesses before they cause unplanned outages.

## Why This Is Best Practice

**Adopted by:** Netflix (Chaos Monkey, Simian Army), Amazon (GameDay), Google (DiRT — Disaster Recovery Testing), LinkedIn, Microsoft Azure
**Impact:** Netflix found that Chaos Monkey reduced the impact of actual cloud failures by 85% after adoption; Basiri et al. (2016): systems with regular chaos testing have 50% fewer unplanned outages; proactive discovery is 10-100x cheaper to fix than production incidents
**Why best:** Systems fail in unexpected ways under real conditions; testing assumptions in a controlled environment is the only way to discover failure modes before users do

Sources: Rosenthal, Jones, et al. "Chaos Engineering" O'Reilly (2020); Basiri et al. "Chaos Engineering" IEEE Software (2016); Nygard "Release It!" Pragmatic Programmers (2018)

## Steps

1. **Define steady-state behavior** — Before injecting chaos, define what normal looks like: p99 latency < 200 ms, error rate < 0.1%, throughput > 1000 RPS. These are your steady-state hypotheses. Chaos engineering tests whether the system maintains steady state under failure conditions. Undefined steady state means undefined success criteria.

2. **Formulate a hypothesis** — "We believe that if a single availability zone fails, the system will fail over within 60 seconds with no more than 0.5% of requests failing." The hypothesis is specific, measurable, and falsifiable. Vague hypotheses ("the system will handle failure") are untestable.

3. **Choose the minimum blast radius** — Start with the smallest scope that validates the hypothesis: one container restart, not a full AZ failure. Chaos engineering is not about maximum destruction; it is about minimum sufficient injection to test the hypothesis. Escalate scope as confidence in the system's resilience grows.

4. **Run in non-production first** — Execute chaos experiments in staging with production-like traffic (via traffic mirroring or synthetic load). Graduating to production requires: successful non-prod results, automated abort conditions, stakeholder awareness, and a defined rollback procedure.

5. **Define abort conditions** — Before starting any experiment, define automatic and manual abort triggers: if error rate exceeds 1%, abort immediately and restore steady state. If on-call is paged during the experiment, abort. If the experiment is still running at the defined end time, abort. Abort conditions are non-negotiable; define them before you inject.

6. **Inject failure** — Use chaos engineering tools: AWS Fault Injection Service (FIS) for cloud failure injection, Chaos Monkey for instance termination, Pumba for Docker container chaos, Litmus or Chaos Mesh for Kubernetes, Gremlin for agent-based injection. Common injection types: instance termination, CPU pressure, memory pressure, network latency, network packet loss, disk I/O stress, dependency failure (kill a downstream service).

7. **Monitor continuously during injection** — Watch your observability stack in real time during the experiment. Track: error rate, latency, service health, and auto-recovery metrics. Compare against steady-state baseline. The goal is to observe what happens, not just whether the system recovers.

8. **Analyze results and compare to hypothesis** — Did the system maintain steady state? If yes: the hypothesis is validated; document the result and increase scope or move to the next failure mode. If no: the hypothesis is falsified; you've found a real reliability weakness. Document the failure mode and create a remediation backlog item.

9. **Automate recurring experiments** — After a successful experiment, automate it to run regularly in staging (weekly or on every deployment). Automated chaos prevents regressions: a code change that inadvertently breaks failure handling is caught before production. Use CI/CD integration to block deployments that fail chaos gates.

10. **Document and share results** — Publish experiment results (hypothesis, injection, observations, conclusions, remediations) to the engineering team. A culture of visible failure discovery reduces fear of chaos engineering and accelerates organizational learning. Every fixed weakness is a prevented outage.

## Rules

- Never run chaos experiments in production without a defined abort procedure, monitoring coverage, and explicit stakeholder awareness.
- Start with the smallest possible blast radius; escalate scope incrementally based on validated confidence.
- An experiment that doesn't test a meaningful hypothesis is not chaos engineering; it is random destruction.
- Fix weaknesses found in experiments before running new experiments; an unfixed weakness compounding with a new injection creates uncontrolled chaos.

## Common Mistakes

- **Running chaos without monitoring** — injecting failures without observing the system's response produces no learning; observability infrastructure must precede chaos engineering.
- **Skipping the hypothesis** — running random failures without a hypothesis tests unknown things against unknown criteria; you can't know if you learned anything.
- **Chaos in production before non-production** — validating experiments in staging is a prerequisite; production chaos without non-production validation multiplies risk.
- **No abort conditions** — an experiment that spirals into an unplanned outage is an incident, not a chaos experiment; abort conditions are mandatory.

## When NOT to Use

- Systems without adequate observability infrastructure — chaos engineering without monitoring is random destruction
- Systems under active incident response — never inject failures into a system that is already degraded
- Critical financial processing systems during peak trading hours — schedule experiments during low-traffic windows
