---
name: apply-yagni-principle
description: Use when deciding whether to add a feature, abstraction, configuration option, or generalization that isn't required by a current, concrete requirement.
source: Beck & Fowler, "Planning Extreme Programming", Addison-Wesley, 2000
tags: [yagni, over-engineering, scope-creep, developer, feature-planning, waste-reduction, maintainability, complexity-reduction]
---

# Apply YAGNI Principle

Don't build it until you need it. Every line of code has a lifetime maintenance cost; code you don't write has none.

## Why This Is Best Practice

**Adopted by:** Google (Engineering Practices: "avoid premature generalization"), Basecamp (Shape Up methodology explicitly bans imagined requirements), Spotify (engineering blog advocates YAGNI as default), Extreme Programming (Kent Beck) — unanimous adoption across agile and lean engineering cultures.
**Impact:** Fowler cites that unused features account for 64% of software features in a typical enterprise application (Standish CHAOS Report baseline). Each speculative feature costs an estimated 3x its initial build cost in total lifecycle expense (planning + build + test + maintain + eventually delete). Google's "Software Engineering at Google" (2020) documents that premature abstraction is among the top causes of unnecessary complexity in large codebases.
**Why best:** The alternative — building for anticipated future needs — requires predicting requirements correctly. Requirements change: Standish Group finds 45% of features in shipped software are never used, and 19% are used rarely. Building for a future that arrives differently than predicted produces dead code that must be maintained, abstracted, and eventually deleted. KISS addresses how to implement things simply; YAGNI addresses whether to implement them at all.

Sources: Beck & Fowler, "Planning Extreme Programming" (Addison-Wesley, 2000); Fowler, "Refactoring" 2nd ed. (2018); Winters, Manshreck & Wright, "Software Engineering at Google" (O'Reilly, 2020); Standish Group CHAOS Report (2015)

## Steps

### 1. Identify the concrete requirement driving the change

Before writing a line of code, state the specific, current requirement in one sentence. If you can't, stop.

- ✅ "User story #42 requires filtering orders by date range."
- ❌ "We might need to filter by other fields someday."
- ❌ "This abstraction will make future features easier."

If the requirement is "might", "someday", or "when we need it", that's YAGNI territory — don't build it.

### 2. Challenge generalizations and abstractions

Every proposed abstraction, config option, or extension point is a YAGNI candidate unless there are two or more concrete current uses.

Ask before adding:

| Proposal | YAGNI challenge |
|----------|----------------|
| Plugin system | "Name two current plugins that require this." |
| Config flag | "Who specifically needs this toggle today?" |
| Abstract base class | "What are the two concrete subclasses right now?" |
| Generic data pipeline | "Which three current pipelines justify a framework?" |
| Versioned API | "Is there a current client that requires v1 and v2 simultaneously?" |

If the answer is "not yet" or "probably someone later", don't build it.

### 3. Apply the Rule of Three before abstracting

Don't generalize on the first or second occurrence of a pattern. Wait for the third.

```python
# First occurrence — implement directly
def send_order_confirmation(order):
    send_email(order.user.email, "Order confirmed", render_order_template(order))

# Second occurrence — still implement directly, note the similarity
def send_shipping_notification(shipment):
    send_email(shipment.user.email, "Shipped", render_shipment_template(shipment))

# Third occurrence — NOW extract the abstraction
def send_notification(recipient_email, subject, body):
    send_email(recipient_email, subject, body)
```

Premature abstraction on the first occurrence encodes today's assumptions into the interface. By the third occurrence, the real shape of the abstraction is visible.

### 4. Delete speculative code in code review

Flag these patterns in PR review as YAGNI violations:

```python
# YAGNI violation — "just in case" parameter
def create_user(name, email, role="user", legacy_id=None, migration_source=None):
    # legacy_id and migration_source never used in any current caller
    ...

# YAGNI violation — unused extension point
class PaymentProcessor:
    def process(self, payment):
        return self._process_impl(payment)

    def _process_impl(self, payment):   # indirection with no current purpose
        ...

# YAGNI violation — dead configuration branch
ENABLE_NEW_CHECKOUT = os.getenv("ENABLE_NEW_CHECKOUT", "false")
if ENABLE_NEW_CHECKOUT == "true":
    # this path was never enabled in production
    ...
```

Request removal, not "we'll use it later."

### 5. Distinguish YAGNI from leaving known problems unfixed

YAGNI is not an excuse to ignore real current problems:

| Situation | Correct response |
|-----------|-----------------|
| Known bug that will definitely be hit | Fix it — that's a current requirement |
| Scaling problem that is 2 weeks away at current growth | Fix it — it's current |
| Abstraction needed for a feature in the current sprint | Build it — it's current |
| "We might need multi-tenancy in 18 months" | YAGNI — don't build it |
| "This will be easier to test if we add an interface" | Build it — testability is a current requirement |

### 6. When you get a "but what if" in code review

"But what if we need to support multiple currencies later?" is the most common form of YAGNI pressure. Respond with:

1. **Cost of adding later**: Is this genuinely hard to retrofit? If not, defer.
2. **Probability**: Is there a concrete roadmap item, or is this speculation?
3. **Cost of the abstraction now**: How much complexity does this add to every reader of this code forever?

Most "what if" abstractions are retrofittable in a day when actually needed, but cost weeks of complexity over their lifetime if added speculatively.

## When NOT to Use

- **Security** — don't defer security controls to "when we need them." Retrofitting auth, encryption, or audit logging is expensive and error-prone. Build them correctly from the start.
- **Foundational schema decisions** — database column types, API field names, and protocol choices are expensive to change. Spend extra time getting these right upfront.
- **Known performance ceilings** — if current load testing shows you'll hit a bottleneck in 4 weeks, design for it now.
- **Regulatory requirements** — compliance controls (GDPR, SOC2, HIPAA) must be built when the system is built, not retrofitted.

## Common Mistakes

**Using YAGNI to justify shortcuts.** "We don't need error handling yet" is not YAGNI — error handling is a current requirement of correct software. YAGNI prevents speculative features, not required correctness.

**Deferring until it's too late.** YAGNI works when you have confident refactorability. If the architecture makes change expensive, YAGNI breaks down. Keep code clean so you can add things when actually needed.

**Confusing YAGNI with KISS.** KISS: implement the current requirement simply. YAGNI: don't implement requirements that don't exist yet. Both apply independently to every feature decision.
