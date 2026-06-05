---
name: apply-feature-toggles
description: Use when you need to separate code deployment from feature release, enable A/B testing, implement gradual rollouts, or allow rapid disabling of features in production without a rollback deployment.
source: Hodgson "Feature Toggles (aka Feature Flags)" (martinfowler.com, 2017); Fowler endorsement and publication; DORA "Accelerate" (2018); LaunchDarkly State of Feature Management Report (2023)
tags: [devops, deployment, feature-flags, continuous-delivery, trunk-based-development, release-management, experimentation]
verified: true
---

# Apply Feature Toggles

Use conditional flags in code to control which features are active in production, decoupling the act of deploying code from the act of releasing a feature — enabling gradual rollouts, instant kill switches, A/B tests, and trunk-based development at scale.

## Why This Is Best Practice

**Adopted by:** Google (Abseil flags library, used across all products), Facebook (Gatekeeper system, controls every feature release), Etsy (pioneered continuous deployment with feature flags), Amazon, Netflix, and documented in DORA's "Accelerate" (Forsgren et al. 2018) as a key practice of elite software delivery. LaunchDarkly reports 80% of Fortune 500 engineering teams use feature management platforms.
**Impact:** DORA research (2023) shows teams using feature flags as part of continuous delivery practices deploy 208× more frequently and have 2,604× faster recovery times than low performers. Facebook's Gatekeeper system enables thousands of experiments simultaneously — without feature flags, each experiment would require a separate deployment. Google's research shows feature flags reduce mean time to recover (MTTR) by 90% compared to rollback deployments for feature-level incidents.
**Why best:** Traditional release gating (branch-per-feature, release trains) creates integration complexity that grows exponentially with team size. Feature toggles allow all code to ship to production continuously on trunk while features remain dark until ready. When a feature causes issues in production, the fix is a configuration change (seconds) rather than a rollback deployment (minutes to hours). This is the mechanism that makes trunk-based development viable at scale.

Sources: Hodgson "Feature Toggles (aka Feature Flags)" (martinfowler.com/articles/feature-toggles.html 2017); Forsgren, Humble & Kim "Accelerate" (2018) Ch. 4; LaunchDarkly "State of Feature Management" (2023); Hammant "Feature Branching is Evil" (trunkbaseddevelopment.com); Fowler "FeatureToggle" (martinfowler.com/bliki)

## Steps

1. **Identify the toggle type** — Choose the right toggle category for the use case: (a) **Release toggle** — hides incomplete features from users, removed after full release; (b) **Experiment toggle** — enables A/B testing, statistical analysis required; (c) **Ops toggle** — circuit breaker for performance-sensitive paths, may be permanent; (d) **Permission toggle** — controls access by user segment or entitlement. Each type has a different lifecycle and cleanup obligation.
2. **Define the toggle decision point** — Identify where in the code the toggle check belongs. Decision points should be at the entry to a feature's behavior, not scattered throughout implementation. Avoid toggle checks deep inside business logic — keep them at boundaries (controllers, service entry points, UI components).
3. **Implement the toggle check** — Use a centralized toggle service or library (LaunchDarkly, Unleash, AWS AppConfig, Split.io, or home-built). The toggle check takes a flag name and an evaluation context (user ID, account, environment) and returns a boolean or variant. Never hardcode toggle state — always read from the toggle system.
4. **Define toggle context** — Specify what context the toggle evaluates against: user ID (personal targeting), cohort (% rollout), account (B2B tier), environment (production vs. staging), or server region. Rich context enables precise targeting and progressive delivery.
5. **Configure the rollout strategy** — For release toggles: start at 0% → internal users → 1% → 10% → 25% → 50% → 100%. For experiment toggles: A/B split with statistical power calculation (minimum 95% confidence, calculated sample size). For ops toggles: maintain a documented runbook for when to flip.
6. **Instrument and monitor** — Log toggle evaluations with context. Track: flag evaluation count, variant distribution, error rate per variant, latency per variant. A feature flag causing a 2× latency increase at 10% rollout is caught before 100%.
7. **Set a cleanup date at toggle creation** — Release toggles are temporary. Every toggle created should have: (a) the condition under which it will be removed (full rollout, experiment conclusion), and (b) a target cleanup date added to the backlog immediately. A toggle without a cleanup date becomes permanent technical debt.
8. **Remove the toggle after release** — After a release toggle reaches 100% and is stable, remove the toggle check, remove the old code path, and delete the flag from the toggle system. Keeping dead toggle code increases cognitive load and risks confusion if the flag is accidentally re-enabled.

## Rules

- Toggles are temporary by default — only Ops toggles may be permanent; all other toggle types must have an explicit cleanup plan.
- Toggle logic must not contain business logic — the toggle decides which path to take; business logic lives in the paths themselves. A toggle check that contains conditional logic beyond "which variant?" is a design smell.
- Test both sides of every toggle — every code path behind a toggle must have test coverage for both the enabled and disabled state; a toggle that disables a path also removes test signal for that path if not explicitly tested.
- Toggle evaluation must be fast — toggle checks happen on every request; a toggle system that adds >1ms per evaluation at scale will degrade p99 latency. Use local evaluation (flag rules cached in-process) not per-request remote calls.

## Common Mistakes

- **Toggle debt accumulation** — Adding toggles without removing them produces a codebase where every feature path has 2–5 conditional branches and developers cannot reason about which code is actually active in production. Schedule cleanup aggressively.
- **Scattering toggle checks throughout implementation** — Checking the same toggle flag in 12 different methods creates a maintenance nightmare when removing the toggle. One toggle, one decision point at the feature boundary.
- **Using toggles as a substitute for automated testing** — "We'll test it with a 1% rollout" is not a testing strategy. Code entering production via a toggle must have test coverage in CI before any traffic reaches it.
- **Not monitoring per-variant metrics** — Releasing to 10% and not measuring error rate, latency, and business conversion separately for the two variants removes the risk-reduction benefit of gradual rollout.

## Examples

**New checkout flow:** Release toggle "new-checkout-v2" evaluated per user ID. Internal team: 100% → QA cohort: 100% → 1% of users → 10% → 50% → 100% over 3 weeks. Error rate and conversion monitored per variant. Toggle removed 2 weeks after full rollout.

**Payment provider failover:** Ops toggle "use-stripe-primary" defaults true. When Stripe has an incident, ops team flips toggle to false in 30 seconds, routing all payments to backup provider. Average deployment rollback: 8 minutes. Toggle flip: 30 seconds.

**Pricing experiment:** Experiment toggle "annual-pricing-upsell-v2" routes 50% of users to new pricing page variant, 50% to control. Measures 30-day conversion to annual plan. Statistical significance reached at 14 days. Winner deployed at 100%, toggle removed.

## When NOT to Use

- When a change is not user-facing and does not need gradual rollout — internal refactors, infrastructure changes, and dependency upgrades do not require feature toggles; they need good test coverage and canary deployments instead.
- When the toggle would persist indefinitely without a clear removal condition — if you cannot define "this toggle is removed when X," it is not a toggle, it is a configuration value; design it as such with proper configuration management.
- When the team has no monitoring infrastructure — a toggle without monitoring of per-variant behavior is just a switch that makes the codebase more complex; implement basic observability before adopting feature toggles at scale.
