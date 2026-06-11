---
name: apply-hyrums-law
description: Use when designing APIs, evolving existing interfaces, or planning a change to system behavior — to account for the fact that all observable behavior will be depended upon by someone, regardless of what the documented contract says.
source: 'Hyrum Wright "Hyrum\'s Law" (hyrumslaw.com); Winters, Manshreck & Wright "Software Engineering at Google" (O\'Reilly, 2020) ch.1; Rich Hickey "Spec-ulation Keynote" (Clojure/conj 2016); Martin Fowler "Tolerant Reader" (martinfowler.com, 2011)'
tags: [api-design, compatibility, contracts, evolution, software-engineering, systems-thinking, dependencies]
verified: true
---

# Apply Hyrum's Law

Design APIs with the smallest possible observable surface area, and before changing any behavior — documented or not — detect who depends on it, because with sufficient users every observable behavior becomes an implicit contract.

## Why This Is Best Practice

**Origin:** Hyrum Wright at Google observed this pattern after maintaining widely-used internal APIs: "With a sufficient number of users of an API, it does not matter what you promise in the contract: all observable behaviors of your system will be depended upon by somebody." The aphorism emerged from real incidents at Google where technically non-breaking changes broke production systems at scale.

**Adopted by:** *Software Engineering at Google* (Winters, Manshreck & Wright, 2020) opens with Hyrum's Law as a foundational principle of sustainable software engineering. Google uses it to justify the investment in automated change detection (Rosie, LSC infrastructure) before any cross-codebase API change. Teams that internalize Hyrum's Law design smaller surfaces, write better deprecation notices, and have fewer surprise breakages.

**Impact:** The law explains why "non-breaking" changes break things. Response body field ordering (unspecified by JSON) — users parse positionally. Error message text (undocumented) — users string-match it. Sort stability (implementation detail) — users relied on CPython's incidentally stable sort before it was formally guaranteed. Hash ordering (explicitly randomized in Python 3.3) — broke code depending on dict iteration order. Each of these was a change to undocumented behavior that produced real failures.

**Why best:** Most API designers think the contract is the spec. Hyrum's Law says the contract is everything users can observe — timing, error messages, field ordering, side effects, response sizes, retry behavior. Designing with this in mind produces smaller, more opaque interfaces with fewer implicit contracts, and a change process that detects real dependents before a change ships.

Sources: Hyrum Wright, hyrumslaw.com; *Software Engineering at Google* ch.1; Hickey "Spec-ulation" 2016; Fowler "Tolerant Reader"

## Steps

### Step 1: Map the full observable surface — not just the documented interface

For every API you own, list what callers can observe beyond the documented contract:

| Category | Examples |
|----------|---------|
| Response structure | Field ordering in JSON, presence of undocumented fields, envelope shape |
| Error behavior | Exact error message text, error code values, HTTP status for each failure mode |
| Performance | p50/p99 latency, response size, rate of retries |
| Ordering | Sort order of list results (even if "unspecified") |
| Side effects | Email sent, audit log written, cache invalidated |
| Timing | Debounce windows, async delivery lag, TTL values |
| Transitive behavior | Which downstream service is called, in what order |

Everything in this table is a de-facto contract for at least one of your users.

### Step 2: Minimize observable surface area at design time

The smaller the observable surface, the fewer implicit contracts. Apply at design time:

**Use opaque types for implementation details:**
```python
# BAD — exposes internal structure; users will depend on field names and ordering
{"result": {"db_row_id": 42, "cache_key": "prod:42", "name": "Widget"}}

# GOOD — opaque ID; internal representation is hidden
{"id": "prod_01HXYZ", "name": "Widget"}
```

**Return only what the contract requires:**
```python
# BAD — extra fields become implicit contracts
def get_user(user_id):
    user = db.query(user_id)
    return user.__dict__  # exposes internal_score, legacy_flag, temp_field

# GOOD — explicit projection
def get_user(user_id):
    user = db.query(user_id)
    return {"id": user.id, "name": user.name, "email": user.email}
```

**Use error codes, not error messages:**
```python
# BAD — error message text is now a contract; you can never reword it
raise ValueError("User not found in database (id=42)")

# GOOD — stable code; message can change without breaking callers
raise ApiError(code="USER_NOT_FOUND", message="User not found")
```

**Paginate and sort explicitly:**
```python
# BAD — "unspecified" ordering that users will depend on
SELECT * FROM products WHERE active = true

# GOOD — explicit ordering is a contract you control
SELECT * FROM products WHERE active = true ORDER BY created_at DESC, id ASC
```

### Step 3: Document the non-contract explicitly

State what you do NOT guarantee, not just what you do. This reduces the number of users who accidentally depend on implementation details:

```yaml
# OpenAPI spec
/products:
  get:
    description: |
      Returns active products.

      **Ordering:** Results are sorted by created_at descending. This ordering is guaranteed.

      **Response fields:** Only fields documented here are part of the contract.
      Additional fields may appear and should be ignored (tolerant reader pattern).

      **Error messages:** The `message` field is for human display only and may change.
      Use `code` for programmatic handling.
```

### Step 4: Before changing behavior, detect who depends on it

Never change observable behavior — even undocumented behavior — without first discovering who relies on it.

**Telemetry approach:** Log the behavior you intend to change, then query for callers who observe it:
```python
# Before changing error message text:
logger.info("error_message_text", text=str(e), caller=request.user_agent)
# Query: how many unique callers received this text in the last 30 days?
```

**Canary approach:** Shadow the new behavior alongside the old, compare outcomes:
```python
# Return old behavior; log what new behavior would have returned
if feature_flag("new_sort_order"):
    new_result = query_new_order()
    if new_result != old_result:
        log_divergence(old_result, new_result, caller_id)
return old_result  # still serve old behavior during canary
```

**Grep/static analysis for string-matched error messages:**
```bash
# Find callers who hardcode your error text
grep -r "User not found in database" ../consumer-services/
grep -r "rate limit exceeded" ../consumer-services/
```

### Step 5: Evolve with explicit signals and a deprecation window

Once you know who depends on the behavior, communicate the change with a timeline:

1. **Announce:** Add a deprecation header or log warning when the old behavior is invoked
2. **Wait:** Give consumers a deprecation window (minimum: one release cycle for internal; 6–12 months for public APIs)
3. **Migrate:** Help dependents move to the new behavior (provide a migration guide, offer a compatibility shim if needed)
4. **Remove:** Delete the old behavior only after confirmed migration

```http
# Deprecation headers (RFC 8594)
Deprecation: Sat, 01 Feb 2025 00:00:00 GMT
Sunset: Sat, 01 Aug 2025 00:00:00 GMT
Link: <https://api.example.com/migration/v3>; rel="successor-version"
```

### Step 6: Write behavioral contract tests, not just interface tests

Tests that only verify the documented interface don't protect against Hyrum's Law breakages. Add tests for observable behaviors you want to preserve:

```python
def test_product_list_sort_order_is_stable():
    """Sort order by created_at desc is a documented contract — verify it."""
    products = api.list_products()
    dates = [p.created_at for p in products]
    assert dates == sorted(dates, reverse=True)

def test_error_code_not_message_is_stable():
    """Error code USER_NOT_FOUND is a contract; message text is not."""
    response = api.get_user("nonexistent")
    assert response.error.code == "USER_NOT_FOUND"
    # deliberately do NOT assert on response.error.message
```

## Rules

- Every observable behavior is a contract for at least one user. Treat undocumented behaviors as contracts-in-waiting.
- Minimize observable surface area: opaque IDs, explicit projections, error codes over error messages, explicit sort orders.
- Never change behavior — documented or not — without first finding who depends on it.
- Document non-contracts explicitly: "this field may change", "message text is non-binding", "additional fields may appear".
- Deprecation window is non-negotiable — the size depends on reach (internal: days to weeks; public: months to years).

## Common Mistakes

**"It's not in the spec, so it's not a contract":** The spec describes your intent; Hyrum's Law describes reality. Undocumented behavior is observed and depended upon regardless.

**Changing error messages without notice:** Error text is the most commonly hard-coded implicit contract. Users parse it to classify failures. Any reword is a breaking change to those users.

**Assuming "unspecified" ordering means no one depends on it:** If your list endpoint has returned results in the same incidental order for 6 months, multiple callers are relying on that order in production.

**Skipping telemetry before the change:** "No one will depend on this" is always wrong at sufficient scale. Measure first.

**Tight deprecation windows:** A 2-week window for a behavior with 50 downstream callers is not a deprecation — it's a breaking change with a short warning.

## Examples

**Python dict ordering:** Pre-3.7, dict iteration order was implementation-defined but incidentally consistent in CPython. Users depended on it. Python 3.7 made it officially ordered — by observing the existing de-facto contract and formalizing it rather than changing it.

**Google sort stability:** Python's `list.sort()` was documented as "not stable" but was stable in CPython. Hyrum's Law: users depended on stability. The fix was to make stability official, not to introduce instability.

**HTTP response field injection:** A team adds a new `_debug` field to JSON responses in staging. It leaks to production. Users start parsing it. Now it's a contract — removing it breaks callers. The fix: never return fields you're not prepared to support forever.
