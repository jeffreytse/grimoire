---
name: apply-domain-based-naming
description: Use when naming variables, functions, classes, modules, or directories — especially when new team members struggle to connect code to domain concepts, or when domain experts and engineers use different terms for the same thing.
source: Dan North, "CUPID—for joyful coding", https://dannorth.net/2022/02/10/cupid-for-joyful-coding/, 2022; Evans, "Domain-Driven Design", Addison-Wesley, 2003
emerging: true
tags: [naming, domain-language, ubiquitous-language, readability, developer, onboarding, maintainability]
related: [apply-composable-design, apply-unix-philosophy, apply-predictable-code, apply-idiomatic-code, apply-domain-driven-design]
---

# Apply Domain-Based Naming

Name code elements using the vocabulary of the problem domain so that code reads like the business it models.

## Why This Is Best Practice

**Adopted by:** Evans' Domain-Driven Design (adopted at Amazon, Netflix, Spotify, and virtually every large-scale service-oriented architecture); Google's internal style guides; Microsoft's .NET naming guidelines (domain vocabulary preferred over technical vocabulary).
**Impact:** When code uses domain language, domain experts can participate in code review, specification errors are caught before they become bugs, and new engineers who know the domain can read the code without a translator. The cognitive translation tax — "ah, `record` means `order`" — compounds across every file read, every PR review, every incident diagnosis. Evans (2003) cites vocabulary mismatch as the primary cause of requirement misunderstanding between domain experts and engineering teams.
**Why best:** Technical naming (suffixes like `Manager`, `Handler`, `Data`, `DTO`; directories named `utils/`, `helpers/`, `common/`) describes implementation artifacts, not domain concepts. When the code changes its implementation (e.g., switches from a relational to a document store), technically-named code becomes misleading. Domain-named code remains accurate as long as the business concept holds — which is far longer than any implementation detail.

Sources: Evans, "Domain-Driven Design" (Addison-Wesley, 2003); Dan North, "CUPID—for joyful coding" (2022); Google, "Google Style Guides"; Feathers, "Working Effectively with Legacy Code" (Prentice Hall, 2004)

**Status:** Emerging — domain naming is a well-understood practice within DDD circles; its framing as a standalone CUPID property by Dan North (2022) is gaining adoption but not yet majority-recognized independently of DDD.

## Steps

### 1. Identify domain vocabulary first

Before naming anything, establish the terms the domain uses. If a glossary does not exist, create one.

- Interview domain experts: what do they call the entities, actions, and states in their work?
- Capture the terms in a glossary (`docs/glossary.md` or similar) alongside the code
- Resolve conflicts: if engineering says "transaction" and finance says "payment", pick one and use it everywhere

---

### 2. Name entities as domain nouns

Classes, structs, tables, and types represent domain concepts. Name them with domain nouns, not technical patterns.

| Technical (avoid) | Domain-based (prefer) |
|---|---|
| `UserRecord` | `Customer` |
| `OrderDTO` | `Order` |
| `PaymentManager` | `PaymentProcessor` |
| `AbstractFactory` | (name the thing it makes) |
| `DataHelper` | (name what data it operates on) |

```python
# Bad — technical suffix adds no meaning
class UserDataRecord:
    user_id: str
    email_address: str

# Good — the domain calls this a Customer
class Customer:
    id: str
    email: str
```

---

### 3. Name operations as domain verbs

Functions and methods name actions. Use the verb the domain uses, not a generic technical verb.

| Technical (avoid) | Domain-based (prefer) |
|---|---|
| `processPayment()` | `chargeCard()`, `authorizePayment()` |
| `handleOrder()` | `placeOrder()`, `cancelOrder()`, `fulfillOrder()` |
| `updateUser()` | `changeEmail()`, `promoteToAdmin()`, `deactivateAccount()` |
| `getData()` | `fetchInvoice()`, `lookupCustomer()` |

```typescript
// Bad — "process" hides what actually happens
async processPayment(payment: Payment): Promise<void>

// Good — domain knows the difference between authorize and capture
async authorizePayment(payment: Payment): Promise<Authorization>
async capturePayment(auth: Authorization): Promise<Receipt>
```

---

### 4. Structure directories by domain, not by technical layer

File organization is naming at the module level. Technical layers (`models/`, `services/`, `controllers/`, `utils/`) create a structure that mirrors the framework, not the domain. Domain-based structure mirrors the business.

```
# Bad — technical layers
src/
  models/user.py, order.py, product.py
  services/user_service.py, order_service.py
  controllers/user_controller.py, order_controller.py
  utils/

# Good — domain structure
src/
  customers/customer.py, customer_store.py
  orders/order.py, order_processor.py, order_store.py
  catalog/product.py, catalog.py
```

Domain-structured code puts everything related to a concept in one place. When an `Order` changes, you look in `orders/`.

---

### 5. Maintain consistency across the entire stack

The same concept must have the same name in every layer: API, database, logs, error messages, documentation. Inconsistency forces every reader to maintain a mental translation table.

- API returns `customer_id` → database column is `customer_id` → log message says `customer_id` → not `user_id` in one place and `account_id` in another
- Run a grep across the codebase: how many terms refer to the same concept? Normalize to one.

## Common Mistakes

**Domain prefix on every class.** `InvoiceInvoiceProcessor`, `OrderOrderRepository` — don't repeat the domain prefix when the file path already provides context. Inside `orders/`, the class is just `Processor` or `Repository`.

**Naming for the current implementation.** `SQLUserRepository` encodes the storage technology in the name. When you switch to MongoDB, the name is wrong. Name what the unit does, not how: `UserRepository` or `UserStore`.

**Over-precise names that don't exist in the domain.** `PendingUnshippedDomesticOrderWithDiscountApplied` is technically accurate but no domain expert uses that phrase. If the domain has a name for the concept, use it; if not, use the simplest name that is accurate.

**utils/ and helpers/.** These are wastebaskets, not domain concepts. When code ends up in `utils/`, it means the domain concept it belongs to has not been named yet. Name it and move the code there.

## When NOT to Use

- **Library/framework internals** — library code often uses technical vocabulary by necessity (e.g., `middleware`, `handler`, `resolver`); these are infrastructure concepts, not domain concepts
- **Generic utilities with no domain meaning** — a `retry` function or a `parse_iso_date` helper has no domain concept to be named after; technical names are correct here
- **Early prototyping** — when the domain is not yet understood, technical placeholders are fine; apply domain naming once the concepts stabilize
