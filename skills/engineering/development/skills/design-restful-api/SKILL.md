---
name: design-restful-api
description: Use when designing HTTP API endpoints — choosing resource names, HTTP methods, status codes, versioning strategy, and response shapes.
source: Roy Fielding, "Architectural Styles and the Design of Network-based Software Architectures", 2000; IETF RFC 9110 (HTTP Semantics, 2022)
tags: [api-design, rest, http, developer, interoperability, client-compatibility]
verified: true
---

# Design RESTful API

Model resources as nouns, use HTTP verbs for actions, and return consistent status codes.

## Why This Is Best Practice

**Adopted by:** GitHub, Stripe, Twilio, Google (Cloud APIs), AWS, Shopify, Salesforce
— REST is the dominant public API style. Stripe's API (launched 2011, unchanged v1
surface) is the industry benchmark for developer experience.
**Impact:** REST APIs outlive their implementations. Stripe's v1 API is still compatible
13+ years later because HTTP semantics don't change. Teams with consistent REST
conventions resolve API design disputes in minutes instead of days — the RFC is the
tie-breaker.
**Why best:** GraphQL requires schema introspection tooling and shifts query complexity
to clients. RPC (gRPC, JSON-RPC) couples clients to server method signatures.
REST uses HTTP semantics that every HTTP client already understands — no SDK required
to consume a well-designed REST API.

Sources: Fielding dissertation (ics.uci.edu/~fielding/pubs/dissertation/), IETF RFC 9110,
Stripe API reference, Google Cloud API Design Guide

## Steps

### 1. Model resources as plural nouns

Resources are things, not actions:

```
# Bad — verb in URL
GET /getUser/123
POST /createOrder
DELETE /deleteItem?id=5

# Good — noun resources
GET /users/123
POST /orders
DELETE /items/5
```

Nest sub-resources when they only exist in the context of a parent:

```
GET /orders/42/items        # items belonging to order 42
POST /orders/42/items       # add item to order 42
DELETE /orders/42/items/7   # remove specific item
```

Don't nest deeper than 2 levels — URLs become unreadable and brittle.

### 2. Map actions to HTTP verbs

| Verb | Semantics | Idempotent? | Body? |
|------|-----------|-------------|-------|
| GET | Read — never mutate | Yes | No |
| POST | Create / non-idempotent action | No | Yes |
| PUT | Full replace | Yes | Yes |
| PATCH | Partial update | No | Yes |
| DELETE | Remove | Yes | No |

Use POST for actions that don't fit CRUD: `/orders/42/cancel`, `/users/me/verify-email`.

### 3. Return correct status codes

Always return the semantically correct code — clients branch on these:

| Code | When |
|------|------|
| 200 | Successful GET, PUT, PATCH |
| 201 | Successful POST that created a resource — include `Location` header |
| 204 | Successful DELETE or action with no response body |
| 400 | Malformed request (syntax, missing required field) |
| 401 | Not authenticated |
| 403 | Authenticated but not authorized |
| 404 | Resource not found |
| 409 | Conflict (duplicate create, stale update) |
| 422 | Validation failed (well-formed but semantically invalid) |
| 429 | Rate limited — include `Retry-After` header |
| 500 | Server error — never leak stack traces |

### 4. Version via URL prefix

```
/v1/users
/v2/users
```

Never version via header — it's invisible in browser logs, curl output, and URLs shared
in Slack. Maintain old versions until all clients have migrated; announce deprecation
with a sunset date in the `Sunset` response header (RFC 8594).

### 5. Use a consistent response envelope

Pick one shape and never deviate:

```json
// Success (single resource)
{ "data": { "id": "usr_123", "email": "..." } }

// Success (collection)
{ "data": [...], "meta": { "total": 500, "page": 2, "per_page": 20 } }

// Error
{ "error": { "code": "validation_failed", "message": "...", "fields": {...} } }
```

### 6. Paginate all collections

Never return unbounded lists — a table with 10M rows will break clients and servers.

```
GET /users?page=2&per_page=50          # offset pagination — simple, allows jumping
GET /users?cursor=eyJpZCI6MTAwfQ==     # cursor pagination — consistent under writes
```

Cursor pagination is preferred for large datasets or frequently-changing collections.

## Rules

- GET must never mutate state — caches and proxies will replay GET requests
- POST is not idempotent — use an idempotency key header (`Idempotency-Key`) for payment or critical creation endpoints
- Don't put actions in the path as verbs unless using POST: `/cancel` is acceptable as a POST target, not as a GET
- Error responses must always include a machine-readable `code` field — human-readable `message` alone forces string parsing
- Never expose database IDs directly in public APIs — use opaque string IDs (`usr_abc123`) that allow backend changes

## Common Mistakes

**Using 200 for errors.** `{ "success": false }` with HTTP 200 breaks every HTTP client
library's error handling. Use the right 4xx/5xx code.

**Returning different shapes for success and error.** Clients handle one envelope, not two.

**Skipping versioning.** "We'll add versioning later" means the first breaking change
is a crisis. Start with `/v1/` on day one.

**Deeply nested routes.** `/companies/1/departments/2/teams/3/members/4` — use
`/members/4` and include parent IDs in the response body.

## When NOT to Use

- **Real-time bidirectional communication:** Use WebSocket or SSE — REST is request/response only
- **Bulk operations crossing many resources atomically:** GraphQL mutations or a custom RPC endpoint is cleaner than a REST workaround
- **Internal service-to-service calls at very high frequency:** gRPC (binary, multiplexed) outperforms REST at 10k+ RPS with strict latency requirements
- **File streaming:** Use pre-signed URLs (S3-style) rather than proxying through REST
