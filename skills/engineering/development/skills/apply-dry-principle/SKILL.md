---
name: apply-dry-principle
description: Use when you notice duplicated code, logic, data, or knowledge across a codebase and need to decide whether and how to consolidate it.
source: Hunt & Thomas, "The Pragmatic Programmer", 1999 (Addison-Wesley)
tags: [code-quality, duplication, refactoring, developer, maintainability, defect-reduction]
verified: true
---

# Apply DRY Principle

Every piece of knowledge must have a single, authoritative representation in the system.

## Why This Is Best Practice

**Adopted by:** Industry-universal. Codified in "The Pragmatic Programmer" (500k+
copies sold), cited in Google Engineering Practices, Microsoft internal style guides,
and Stripe's engineering handbook. No top-tier engineering org advocates for deliberate
duplication.
**Impact:** Every duplicated piece of knowledge is a latent defect — a future change
must be applied N times, and the probability of missing one instance scales with N.
DRY reduces the change surface proportionally. In large codebases, duplicated business
logic is the leading cause of behavioral divergence between services.
**Why best:** DRY addresses knowledge duplication, not just code copy-paste. A business
rule encoded in a database constraint, an API validator, and a frontend form is triplicated
knowledge — changing the rule requires three coordinated edits, each a defect opportunity.
WET (Write Everything Twice) produces inconsistency at scale.

Sources: Hunt & Thomas, "The Pragmatic Programmer" ch. 7 (The Evils of Duplication);
Google Engineering Practices; Martin Fowler, "Refactoring" (2018 ed.), ch. 3

## Steps

### 1. Classify the duplication type

| Type | Example | Risk |
|------|---------|------|
| Code | Same function copy-pasted | Low — caught by linters |
| Data | Magic number repeated | Medium — one instance gets updated |
| Logic | Business rule in 3 places | High — behavioral drift |
| Knowledge | Same concept with different names | Critical — silent divergence |

### 2. Find the authoritative location

Ask: "If this knowledge changes, what is the ONE place I'd update?"
That location should be the source of truth. Everything else should reference it.

### 3. Extract to the right abstraction level

Match the extraction scope to the duplication scope:

| Duplication scope | Extraction |
|-------------------|-----------|
| Within a function | Named constant or local variable |
| Within a module | Private function |
| Across modules in a service | Shared utility / helper |
| Across services | Shared library, schema registry, or API contract |
| Across frontend and backend | Code generation from a single schema (e.g., OpenAPI, Protobuf) |

### 4. Replace all instances

Replace every occurrence with a reference to the extracted source. Verify with grep
or your IDE's find-all-references — missing one instance recreates the problem.

### 5. Verify single-edit guarantee

Make a hypothetical change: "If this knowledge changes, how many edits are required?"
The answer should be exactly one. If it's more, the extraction is incomplete.

## Rules

- **Rule of Three:** Don't extract after the first or second occurrence — wait for the third. Premature extraction creates the wrong abstraction, which is harder to fix than duplication.
- **Coincidental similarity ≠ shared knowledge:** Two functions that look the same but represent different concepts should NOT be merged — they will diverge.
- **Test duplication is acceptable:** Tests often repeat setup for clarity. DRY tests can obscure what each test actually verifies. Prefer duplication in tests over abstraction that hides intent.
- **Don't DRY across service boundaries without a contract:** Sharing code between services via copy-paste is often correct; sharing via a live dependency creates coupling. Prefer generating both from a single schema.

## Common Mistakes

**Abstracting coincidental similarity.** Two functions that compute totals for orders
and invoices look the same today. Next quarter, order totals add tax; invoice totals
don't. Premature merging forces an ugly split or a flag parameter.

**DRY the wrong layer.** Deduplicating SQL queries by sharing a query builder
object couples unrelated features. The duplication was in SQL text, not in knowledge.

**Forgetting knowledge duplication across systems.** API documentation that duplicates
the API implementation is duplicated knowledge — use code generation or contract-first
design to keep one authoritative source.

## When NOT to Use

- When two similar things represent genuinely different concepts (coincidental similarity)
- When the abstraction would require a flag parameter to handle variants — this signals the wrong abstraction; keep them separate
- In tests where duplication aids readability and independence
- When deduplication would introduce a cross-service dependency where a schema/contract doesn't already exist
