---
name: apply-branch-by-abstraction
description: Use when replacing a large or deeply integrated component of a system incrementally without creating a long-lived feature branch, by introducing an abstraction layer that allows old and new implementations to coexist.
source: Fowler "BranchByAbstraction" (martinfowler.com/bliki/BranchByAbstraction.html, 2014); Hammant (trunkbaseddevelopment.com); Humble & Farley "Continuous Delivery" (2010)
tags: [development, refactoring, continuous-delivery, trunk-based-development, patterns, abstraction, migration]
verified: true
---

# Apply Branch by Abstraction

Replace a large integrated component by introducing an abstraction layer that both old and new implementations satisfy, allowing incremental migration on trunk without a long-lived feature branch.

## Why This Is Best Practice

**Adopted by:** Google (primary technique for large-scale refactoring across monorepo), Thoughtworks (standard practice for trunk-based development), documented in Humble & Farley "Continuous Delivery" (2010) as the enabling technique for replacing large components without branching, and implemented by teams using Trunk-Based Development at scale.
**Impact:** Long-lived feature branches are one of the top predictors of integration failure — teams using branches longer than 1 day have 3× higher merge conflict rates (DORA research 2018). Branch by Abstraction eliminates long-lived branches for component replacement work by keeping the change on trunk throughout. Google uses this technique to refactor components used by thousands of callers across their monorepo without any team experiencing a broken trunk.
**Why best:** When replacing a large component (database library, HTTP client, payment provider, authentication system), the naive approaches both fail: (a) a feature branch becomes a multi-month integration nightmare as trunk diverges, (b) an in-place replacement breaks all callers simultaneously. Branch by Abstraction provides a third path — introduce an interface, migrate callers gradually, swap the implementation, remove the abstraction — all on trunk, all shippable at every commit.

Sources: Fowler "BranchByAbstraction" (martinfowler.com/bliki/BranchByAbstraction.html 2014); Humble & Farley "Continuous Delivery" (Addison-Wesley 2010) Ch. 14; Hammant "Branch by Abstraction" (trunkbaseddevelopment.com); Feathers "Working Effectively with Legacy Code" (2004) on Seams

## Steps

1. **Identify the component to replace** — Specify exactly what will be replaced: a library, a service client, a database access layer, a third-party integration. The component must have a clear boundary — a set of callers and a set of responsibilities. Ambiguous boundaries produce ambiguous abstractions.
2. **Introduce an abstraction layer** — Create an interface (or abstract class, or protocol, depending on language) that captures the component's contract — the methods and behaviors its callers depend on. This interface becomes the "branch" — all callers will depend on the interface, not the implementation. Keep the interface thin: only the behaviors that callers actually use.
3. **Place the existing implementation behind the abstraction** — Wrap or adapt the current component to implement the new interface. No caller behavior changes — the interface matches current behavior. Run all tests; they must pass with zero changes. This step verifies the abstraction captures the existing contract correctly.
4. **Build the new implementation behind the abstraction** — Implement the new component to satisfy the same interface. The new implementation can be built incrementally; it does not need to be complete before migration begins. Only the behaviors already implemented in the new version are migrated in each step.
5. **Migrate callers incrementally** — Update callers one by one (or one service at a time) to use the new implementation, routing through the abstraction. Use a feature toggle or a factory/configuration to control which implementation each caller gets. Both implementations remain active; you can ship every migration step.
6. **Run both implementations in parallel if needed** — For high-risk components, run old and new implementations in parallel: send each request to both, compare outputs, log discrepancies. This "shadow mode" validates the new implementation against production traffic without exposing users to any failures.
7. **Complete migration and verify** — Once all callers use the new implementation and the old implementation has no active callers, verify with production monitoring that no traffic reaches the old path.
8. **Remove the abstraction** — After full migration, assess whether the abstraction is worth retaining (useful for future changes, testing) or should be removed. If removing: inline the new implementation directly into callers and delete the interface. Remove the old implementation. This is the cleanup step that prevents abstraction accumulation.

## Rules

- The abstraction must match the existing contract exactly before any migration starts — tests must pass with the old implementation behind the abstraction before building the new one.
- Migrate callers one at a time, shipping after each — never migrate 50 callers simultaneously; if the new implementation has a defect, 50 callers are affected. Incremental migration limits blast radius.
- The abstraction must be thin — it should capture only the behaviors callers depend on, not every method the component exposes. A bloated interface defeats the purpose.
- Do not retain the abstraction "just in case" unless it has ongoing value — every interface adds indirection cost and cognitive overhead; remove it after migration unless it serves a clear ongoing purpose.

## Common Mistakes

- **Designing the abstraction around the new implementation's API** — The abstraction must satisfy existing callers first; the new implementation conforms to the abstraction, not the other way around. Designing for the new implementation forces all callers to change simultaneously.
- **Not testing both implementations** — After introducing the abstraction, tests that only run with the old implementation behind it will not catch regressions in the new implementation until callers migrate. Run the test suite with both implementations.
- **Keeping the old implementation alive too long** — After migration is complete, the old implementation creates confusion about which path is active and maintenance cost for dead code. Remove it promptly.
- **Applying this to a component with no clear boundary** — If the "component" is actually business logic scattered across 40 files with no clear interface, the first step is refactoring to create a seam (Feathers), not introducing an abstraction.

## Examples

**Replacing HTTP client:** System uses legacy HttpClient with 200 callers. Abstraction: `HttpGateway` interface with `get()`, `post()`, `delete()`. Wrap legacy client behind `HttpGateway`. Build new async client implementing `HttpGateway`. Migrate 20 callers/sprint. 10 sprints later: all callers use new client, legacy removed, interface retained for testability.

**Replacing payment provider:** 15 order flow methods call Stripe SDK directly. Abstraction: `PaymentGateway` interface with `charge()`, `refund()`, `retrievePayment()`. Wrap Stripe behind interface. Build Adyen implementation. Shadow mode: 30 days comparing outputs. Migrate 3 methods/week. Complete in 5 weeks.

## When NOT to Use

- When the component being replaced is small (under 5 callers, under 500 lines) — the overhead of an abstraction layer exceeds the benefit; do an in-place replacement with a good test suite and a single well-prepared PR.
- When the component's callers are already behind an abstraction — if callers already use an interface, just provide a new implementation; no new abstraction is needed.
- When the migration timeline is so short that trunk divergence is not a risk — a 3-day replacement does not need branch-by-abstraction; the technique is designed for multi-week or multi-month migrations.
