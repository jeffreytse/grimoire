---
name: apply-strangler-fig-pattern
description: Use when replacing a legacy system or component incrementally without a full rewrite, by routing traffic progressively from the old system to a new implementation.
source: Fowler "StranglerFigApplication" (martinfowler.com, 2004); Thoughtworks Technology Radar; Newman "Monolith to Microservices" (2019); AWS Migration Whitepaper
tags: [architecture, legacy, migration, refactoring, patterns, modernization, incremental]
verified: true
---

# Apply Strangler Fig Pattern

Replace a legacy system incrementally by intercepting calls, redirecting traffic to a new implementation one slice at a time, and decommissioning the old system only after the new one fully covers its functionality.

## Why This Is Best Practice

**Adopted by:** Netflix (streaming platform migration), AWS (recommended migration pattern in Well-Architected Framework), Thoughtworks (Technology Radar "Adopt" since 2016), Azure Architecture Center, Google Cloud migration guides, and the dominant pattern in Thoughtworks consulting engagements for enterprise legacy modernization.
**Impact:** The Standish Group CHAOS Report (2015) shows that "big bang" rewrites fail or are cancelled 70% of the time. Strangler Fig avoids this by keeping the legacy system in production throughout migration — the organization always has a working system. Netflix's incremental migration from monolith to microservices using this pattern is credited with enabling their 2× traffic growth in 2012–2015 with zero downtime.
**Why best:** A full rewrite requires freezing the legacy system's feature development for 12–24 months while the new system catches up — a period during which the business continues to change and the new system falls further behind. Strangler Fig allows feature development to continue in parallel with migration, delivers business value in each migration slice, and allows the team to learn the domain before committing to the full new architecture.

Sources: Fowler "StranglerFigApplication" (martinfowler.com 2004); Newman "Monolith to Microservices" (O'Reilly 2019) Ch. 3; AWS "Modernizing Legacy Applications" whitepaper; Thoughtworks Technology Radar vol. 25; Feathers "Working Effectively with Legacy Code" (2004)

## Steps

1. **Map the legacy system's capabilities** — Enumerate every capability the legacy system provides: endpoints, batch jobs, integrations, user-facing features. This becomes the migration roadmap. You cannot strangle what you have not mapped.
2. **Place an intercepting layer (facade)** — Deploy a proxy, API gateway, or facade in front of the legacy system that can route requests to either old or new implementation. All traffic passes through this layer. The facade must be deployed before migration begins — it is the control plane for the transition.
3. **Identify the first migration slice** — Choose the first capability to migrate based on: low risk (not business-critical), high standalone-ability (minimal dependencies on other legacy parts), and high value (demonstrates momentum). Do not start with the most complex or tightly coupled component.
4. **Build the new implementation for the slice** — Implement the selected capability in the new system. Do not replicate all edge cases from the legacy system — use this as an opportunity to eliminate known defects and unnecessary complexity. Validate with the same test cases as the legacy implementation.
5. **Route traffic to the new implementation** — Update the facade to route the migrated capability's traffic to the new implementation. Keep the legacy implementation running and accessible via the facade for rollback. Validate with real production traffic — do not rely solely on test environments.
6. **Verify and stabilize** — Monitor the new implementation for errors, latency, and behavioral differences from the legacy. Run both systems in parallel for a defined stabilization window (typically 1–2 weeks for non-critical paths, longer for critical). Compare outputs if the migration is high-risk.
7. **Decommission the migrated piece from legacy** — Only after stabilization, remove or disable the corresponding code in the legacy system. Update the facade to remove the legacy route. Do not leave dead code in the legacy system — it creates confusion about what has and has not been migrated.
8. **Repeat until the legacy system is empty** — Continue slicing, building, routing, verifying, and decommissioning until the legacy system handles no traffic. At that point, decommission the legacy system entirely and optionally remove the facade if no longer needed.

## Rules

- Never remove the legacy route before the new implementation has served production traffic without errors — the facade exists precisely for rollback capability.
- Migrate capabilities in dependency order — migrate dependencies before the components that depend on them; migrating a component before its dependencies forces you to call back into the legacy system, creating a circular dependency.
- The migration slice must be independently deployable — if a slice requires migrating 5 other components simultaneously to work, it is not a valid slice; decompose further.
- Track migration completeness explicitly — a list of "capabilities remaining in legacy" is the single metric that drives migration progress; without it, teams lose track of what has and has not been migrated.

## Common Mistakes

- **Migrating too large a slice** — Starting with a slice that spans multiple capabilities, has many integrations, and takes 6 months creates all the risks of a big bang rewrite inside a strangler wrapper. Slices should ship in 2–6 weeks.
- **Not deploying the facade before migration starts** — Some teams try to migrate without an intercepting layer, planning to "cut over" at the end. This eliminates the ability to route gradually and creates a binary switch that eliminates the pattern's risk-reduction benefit.
- **Running parallel systems indefinitely** — Leaving the legacy path live after the new implementation is stable incurs ongoing maintenance cost for both systems. Define explicit decommission criteria and enforce them.
- **Migrating in functional silos, not capability slices** — Migrating "all database logic" or "all authentication" across the system at once is not a strangler approach; it is a horizontal refactor that does not ship independently. Slice vertically by user-facing capability.

## Examples

**E-commerce platform:** Legacy monolith handles orders, inventory, accounts, payments. Facade: API gateway. Slice 1: product search (read-only, low risk). Slice 2: user account management. Slice 3: order history (read). Slice 4: checkout flow (last, highest risk). Each slice takes 3–6 weeks. Legacy decommissioned after 18 months.

**Internal reporting system:** Legacy batch report generation. Facade: message queue that routes report requests. New implementation handles financial summary reports first (highest volume, simplest logic). Legacy handles complex custom reports until migrated. Each report type migrated independently.

## When NOT to Use

- When the legacy system is well-understood, small, and can be safely replaced with a full rewrite in under 3 months — the overhead of a facade and incremental migration exceeds the risk reduction benefit.
- When the legacy system cannot be wrapped with a facade without a full refactor — some tightly coupled systems require architectural surgery before a strangler approach is feasible; use Feathers' "Seams" technique first.
- When regulatory or contractual requirements prohibit running two versions of a system simultaneously — some compliance contexts require point-in-time cutover; document the exception and plan a detailed cutover with rollback.
