---
name: validate-external-input
description: Use when writing code that accepts data from outside the process boundary — HTTP requests, file uploads, database reads, CLI arguments, environment variables, or inter-service calls — to reject malformed, malicious, or out-of-range data before it reaches business logic.
source: OWASP Input Validation Cheat Sheet (2023); MITRE CWE Top 25; CERT Secure Coding Standards (SEI)
tags: [input-validation, security, injection, boundary, sanitization, developer, owasp, defect-reduction]
related: [apply-fail-fast, design-threat-model]
---

# Validate External Input

Reject invalid data at every process boundary before it reaches business logic.

## Why This Is Best Practice

**Adopted by:** OWASP (foundational to all 10 categories in the OWASP Top 10), Google
(mandatory in Google Security Guidelines), AWS (enforced via API Gateway validation),
and every major web framework (Django validators, Rails strong parameters, Spring
Validation, Express-validator).
**Impact:** Injection-class vulnerabilities — caused directly by insufficient input
validation — ranked #3 in OWASP Top 10 (2021) and appear in 40%+ of MITRE CVE records
(NVD data). The CERT Secure Coding Standard cites input validation as the single most
effective control against the widest class of software vulnerabilities.
**Why best:** Output-encoding-only approaches catch XSS at render time but miss SQL
injection, path traversal, overflow, and business-logic attacks. Validating at the
boundary is the only defense that covers all injection classes before the data is used.

Sources: OWASP Input Validation Cheat Sheet (2023); MITRE CWE-20; SEI CERT Coding
Standards; Google Security Engineering blog

## Steps

### Step 1: Identify all external boundaries

List every point where data enters the process:
- HTTP request body, query params, headers, cookies
- File and stream contents (uploads, config files, stdin)
- Environment variables and CLI arguments
- Database reads (treat as external if data was ever user-supplied)
- Messages from queues, webhooks, or other services

### Step 2: Validate at the boundary, not deep in the call chain

Place validation in the layer that first receives external data — controller, handler,
or adapter. Never rely on business logic or the database to catch it.

```python
# Wrong — validation buried in service layer
def process_order(order_id):
    order = db.get(order_id)
    if order_id < 0:  # too late — DB query already ran
        raise ValueError

# Right — validate at the handler
def handle_request(request):
    order_id = parse_int(request.params["order_id"], min=1)  # fail here
    process_order(order_id)
```

### Step 3: Validate in this order

1. **Type** — is it the expected type? (int, string, UUID, date)
2. **Format** — does it match the expected pattern? (email regex, ISO date, UUID format)
3. **Range / length** — within allowed bounds? (int 1–10000, string ≤255 chars)
4. **Business rules** — does it satisfy domain constraints? (order_id exists and belongs to this user)

Stop at the first failure. Do not run later checks on already-invalid data.

### Step 4: Whitelist, don't blacklist

Define what is valid and reject everything else. Never try to enumerate what is invalid.

```python
# Wrong — blacklist approach (incomplete, bypassable)
if "<script>" in value or "DROP TABLE" in value:
    reject()

# Right — whitelist approach
if not re.match(r'^[a-zA-Z0-9_\-]{1,64}$', value):
    reject()
```

### Step 5: Return a specific error; never echo unvalidated input

The error message must describe the constraint, not the input. Echoing unvalidated input
back to the caller is a reflected XSS vector.

```python
# Wrong — reflects input back to the caller
raise ValueError(f"Invalid input: {user_input}")

# Right — describes the constraint
raise ValueError("username must be 1–64 alphanumeric characters")
```

### Step 6: Validate all fields, not just the ones you currently use

Unused fields today may be used tomorrow. Validating them now prevents a future engineer
from trusting them without realizing they've never been checked.

## When NOT to Use

- **Internal function calls between components you own** — if the caller is internal code
  under your control, `apply-fail-fast` (assertions) is the right tool, not validation.
- **Re-validating already-validated data** — once data has passed boundary validation and
  is in your domain model, trust it. Re-validating on every use adds noise and hides the
  real boundary.

## Common Mistakes

**Validating at the wrong layer.** Validation in a service method or ORM model runs after
routing, logging, and sometimes partial processing. Move it to the entry handler.

**Sanitizing instead of rejecting.** Stripping bad characters from input silently masks
problems. If input doesn't match the expected pattern, reject it with an error.

**Trusting `Content-Type` headers.** Clients can send any Content-Type with any body.
Validate the actual content, not the declared type.

**Skipping validation on "internal" services.** Microservice A calling microservice B
is still an external boundary — both services should validate. Compromise of A is a
vector into B if B trusts its inputs.
