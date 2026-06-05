---
name: design-api
description: Use when designing a new REST API, adding endpoints to an existing API, or reviewing an API design for correctness, consistency, and long-term maintainability.
source: Fielding "REST Architectural Constraints" (2000 dissertation); Google API Design Guide (cloud.google.com/apis/design); Stripe API design; OpenAPI Specification (openapis.org)
tags: [api-design, rest, http, developer-experience, versioning, contracts]
verified: true
---

# Design API

Design a clean, consistent, and evolvable REST API — with correct resource naming, HTTP semantics, status codes, versioning, pagination, and error format.

## Why This Is Best Practice

**Adopted by:** Stripe (widely cited as the gold standard for developer experience — their API design decisions are documented in their engineering blog and referenced across the industry), Google (Google API Design Guide is public and applies to all Google Cloud APIs, covering 200+ APIs), Twilio (developer-centric API design has been their primary competitive differentiator since 2010).

**Impact:** Stripe's API quality directly drives revenue — their developer NPS of 60+ (vs. industry average of 30) is attributed primarily to API consistency and documentation. The Postman State of the API Report (2023, n=40,000 developers) found that 52% of developers cite inconsistent APIs as the primary reason they abandon an integration. A well-designed API reduces integration time by 30–50% according to Twilio's internal developer experience metrics.

**Why best:** Ad-hoc API design produces inconsistency across endpoints (some use `user_id`, others use `userId`, others use `id`) that forces clients to write special-case code for each endpoint. REST conventions (Fielding 2000) are the most widely understood API contract in software — following them means developers can correctly guess your API's behavior before reading the docs.

Sources: Fielding REST dissertation (2000), Google API Design Guide (cloud.google.com/apis/design), Stripe API reference (stripe.com/docs/api), Postman State of the API 2023

## Steps

### 1. Define resources first, endpoints second

REST APIs are organized around resources (nouns), not actions (verbs). Identify your resources before choosing URLs.

Resources map to domain entities: `user`, `order`, `payment`, `product`, `invoice`.

Each resource has two URL patterns:
```
/users          → collection
/users/{id}     → individual resource
```

Never put verbs in URLs:
- Bad: `/getUser`, `/createOrder`, `/deletePayment`
- Good: `GET /users/{id}`, `POST /orders`, `DELETE /payments/{id}`

Exception: use action sub-resources for operations that do not map to CRUD — sparingly:
- `POST /payments/{id}/capture` — captures an authorized payment
- `POST /accounts/{id}/deactivate` — transitions account state

### 2. Use HTTP methods with correct semantics

| Method | Semantics | Idempotent | Safe |
|--------|-----------|-----------|------|
| GET | Read resource or collection | Yes | Yes |
| POST | Create resource or trigger action | No | No |
| PUT | Replace resource entirely | Yes | No |
| PATCH | Partial update (send only changed fields) | No* | No |
| DELETE | Remove resource | Yes | No |

*PATCH idempotency depends on implementation — prefer patch formats (JSON Merge Patch RFC 7396) that make it idempotent.

Never use GET with a request body — proxies and caches may discard it.

### 3. Choose a consistent URL convention and never mix

Pick one and enforce it across the entire API:
- Use **plural nouns** for collections: `/users`, `/orders`, not `/user`, `/order`
- Use **kebab-case** for multi-word paths: `/payment-methods`, not `/paymentMethods` or `/payment_methods`
- Use **snake_case** for JSON fields: `{ "first_name": "Ada" }` — consistent with most backend languages and less error-prone than camelCase in case-sensitive languages
- Nest resources only when the child cannot exist without the parent: `/orders/{id}/line-items` is correct; `/users/{id}/invoices` is debatable (invoices can exist without a user context)

### 4. Return correct HTTP status codes

Use status codes that match HTTP semantics — clients use these for error handling logic, not just logging.

**2xx Success:**
- `200 OK` — GET, PUT, PATCH response with body
- `201 Created` — POST that created a resource; include `Location: /resources/{new-id}` header
- `204 No Content` — DELETE success, or PUT/PATCH with no body
- `202 Accepted` — async operation started; return a job/operation ID

**4xx Client Error:**
- `400 Bad Request` — malformed request, validation failure
- `401 Unauthorized` — missing or invalid authentication
- `403 Forbidden` — authenticated but not authorized for this resource
- `404 Not Found` — resource does not exist
- `409 Conflict` — resource state conflict (e.g., duplicate creation, optimistic lock failure)
- `422 Unprocessable Entity` — syntactically valid but semantically invalid (prefer over 400 for validation errors per RFC 9110)
- `429 Too Many Requests` — rate limit exceeded; include `Retry-After` header

**5xx Server Error:**
- `500 Internal Server Error` — unexpected server failure
- `503 Service Unavailable` — server overloaded or in maintenance; include `Retry-After`

### 5. Design a consistent error response format

Every error must return the same JSON structure — clients should be able to handle all errors with one code path.

```json
{
  "error": {
    "code": "validation_failed",
    "message": "The request body contains invalid fields.",
    "request_id": "req_01HX7G5F8K4PBNJD4W3TYZV9AB",
    "details": [
      {
        "field": "email",
        "code": "invalid_format",
        "message": "Must be a valid email address."
      }
    ]
  }
}
```

Rules for errors:
- `code` is a stable machine-readable string — clients branch on this, not on `message`
- `message` is human-readable — it can change without breaking clients
- `request_id` is present on every response — essential for support and debugging
- `details` is present for validation errors — one entry per invalid field

### 6. Version the API

Plan for breaking changes before they happen. The two dominant approaches:

**URL versioning** (recommended for public APIs): `/v1/users`, `/v2/users`
- Pros: immediately visible, cacheable, no header negotiation
- Cons: clients must change URLs on upgrade
- Used by: Stripe, Twilio, GitHub

**Header versioning**: `Accept: application/vnd.api+json;version=2`
- Pros: URL is stable; version is a negotiation concern
- Cons: harder to test in browser, less visible in logs
- Used by: GitHub (for media type versioning)

Rules:
- Never change the meaning of an existing field — add new fields instead
- Never remove a field without a deprecation period (minimum 6 months for public APIs)
- A version bump is required for any breaking change: removing fields, changing field types, changing error codes clients branch on

### 7. Design pagination for all collection endpoints

Never return unbounded collections. Every `GET /resources` must be paginated.

**Cursor-based pagination** (preferred for large or frequently updated datasets):
```json
{
  "data": [...],
  "pagination": {
    "next_cursor": "cursor_01HX7G5F8K4PBNJD4W3",
    "has_more": true
  }
}
```
Request: `GET /orders?after=cursor_01HX7G5F8K4PBNJD4W3&limit=20`

**Offset-based pagination** (acceptable for small, stable datasets):
```json
{
  "data": [...],
  "pagination": {
    "total": 1847,
    "page": 3,
    "per_page": 20
  }
}
```

Cursor-based is preferred because offset pagination breaks when items are inserted or deleted during pagination.

### 8. Document with OpenAPI before implementing

Write an OpenAPI 3.1 spec before writing any handler code. The spec is the contract:
- Forces you to think through every field, every status code, every error before the implementation is wired up
- Enables mock servers so frontend/client teams can start integration immediately
- Generates client SDKs and interactive docs (Swagger UI, Redoc)

Minimum per endpoint: summary, all parameters, all request body fields with types, all response codes with bodies, and at least one example.

## Rules

- Every endpoint must handle `request_id` — generate one at the edge if the client does not provide it
- All timestamps are ISO 8601 in UTC: `"2024-03-15T14:30:00Z"` — never Unix epoch in JSON
- All money amounts are integers in the smallest currency unit (cents) — never floats; `"amount": 2999` = $29.99
- Boolean fields must not use integer 0/1 — use `true`/`false`
- Nullable fields must be declared explicitly — absence of a field and `null` are different states
- Rate limits must be documented and their headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`, `Retry-After`) must be returned on every response

## Examples

**Correct resource design:**
```
POST   /orders                     → create order
GET    /orders                     → list orders (paginated)
GET    /orders/{id}                → get one order
PATCH  /orders/{id}                → update order fields
DELETE /orders/{id}                → cancel order
POST   /orders/{id}/fulfillments   → fulfill an order (action sub-resource)
```

**Correct error for a validation failure (422):**
```json
{
  "error": {
    "code": "validation_failed",
    "message": "Invalid request parameters.",
    "request_id": "req_abc123",
    "details": [
      { "field": "quantity", "code": "must_be_positive", "message": "Quantity must be greater than 0." }
    ]
  }
}
```

## Common Mistakes

- Verbs in URLs (`/getUser`, `/createOrder`) — breaks REST resource model; clients cannot predict URL structure
- Returning `200 OK` with `{ "success": false }` in the body — clients cannot use HTTP status for error detection
- Returning all data in one request with no pagination — breaks at scale; client memory and latency blow up
- Using `400` for all errors — clients cannot distinguish auth failures from validation failures
- Storing money as a float — floating-point arithmetic errors cause real financial discrepancies; always use integer cents
- No `request_id` — support tickets cannot be correlated to server logs
