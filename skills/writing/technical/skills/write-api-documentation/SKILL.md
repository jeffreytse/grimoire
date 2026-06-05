---
name: write-api-documentation
description: Use when writing or reviewing API reference documentation, endpoint descriptions, or developer-facing SDK documentation
source: Google Developer Documentation Style Guide; OpenAPI Specification (Swagger) documentation standards; Stripe API documentation (industry benchmark)
tags: [technical-writing, api-documentation, developer-experience, openapi]
verified: true
---

# Write API Documentation

Produce accurate, developer-friendly API documentation that enables independent integration without support escalation.

## Why This Is Best Practice

**Adopted by:** Stripe (industry benchmark for API docs), Twilio, GitHub, AWS, and all major API-first companies; OpenAPI Specification is the de facto standard with 10,000+ public APIs documented; Google Developer Documentation Style Guide is used by 1,000+ Google products and widely adopted externally.
**Impact:** Well-documented APIs reduce developer time-to-first-call by 60–80%; Stripe's docs are credited as a key factor in developer adoption that drove the company from $0 to $95B valuation.
**Why best:** APIs are products; documentation is the UI. Bad docs mean developers integrate incorrectly, generate support tickets, or abandon the API.

Sources: Google Developer Documentation Style Guide; OpenAPI Specification 3.1 (OAI); Stripe API Documentation (reference benchmark); Redocly API documentation standards.

## Steps

1. **Document every endpoint with OpenAPI 3.1** — write a machine-readable OpenAPI spec first. Each endpoint needs: HTTP method, path, operationId, summary (one line), description, parameters, request body schema, response schemas for every status code.

2. **Write a getting-started guide** — before any reference docs, provide: authentication setup (with working code example), the simplest possible working request, expected response, and a "what next" path. Target time-to-first-successful-call under 15 minutes.

3. **Document authentication explicitly** — specify the auth scheme (API key, OAuth 2.0, JWT), where credentials go (header, query param, body), how to obtain credentials, and how to rotate or revoke them. Include a working curl example.

4. **Write parameter descriptions that explain, not restate** — bad: `user_id: The user's ID`. Good: `user_id: Unique identifier of the user whose data to retrieve. Find this in the Users list or account settings.`

5. **Provide working code examples in multiple languages** — include curl, Python, JavaScript/Node, and Ruby at minimum. Examples must be copy-paste runnable. Use real, syntactically correct values — not `<your-value-here>` placeholders.

6. **Document all error responses** — for every endpoint, list every non-200 status code with: HTTP status, error code string, human-readable message, what caused it, and how to fix it. This is the most under-documented aspect of most APIs.

7. **Explain rate limits, pagination, and versioning** — document: requests/second or requests/minute limits, pagination strategy (cursor, offset, page), API versioning scheme and deprecation policy, and how to request higher limits.

8. **Add a changelog** — every API change (new endpoint, new parameter, deprecation, breaking change) must appear in a dated changelog with migration guidance for breaking changes.

9. **Test every code example** — run every example in the documentation against the production API before publishing. Broken examples destroy developer trust immediately.

10. **Implement feedback mechanisms** — add inline "Was this helpful?" ratings, a GitHub link to suggest edits, and a support path. Review low-rated pages monthly.

## Rules

- Never document a field as "optional" without specifying what happens when it's omitted.
- Every breaking change requires a migration guide, not just a changelog entry.
- Code examples must use realistic values, not generic placeholders.
- Maintain documentation in the same repository as the API code; outdated docs are worse than none.

## Common Mistakes

- **Reference-only documentation** — providing only endpoint specs without tutorials or use-case guides leaves developers unable to understand how to accomplish real tasks.
- **Undocumented error codes** — returning `error: "invalid_request"` without explanation forces developers to contact support for every integration issue.
- **Inconsistent naming** — mixing `user_id`, `userId`, and `userID` in different endpoints creates integration errors and erodes trust.
- **No sandbox/test environment** — requiring production credentials to test an API increases the barrier to integration and risk of data corruption.

## When NOT to Use

- When documenting internal-only APIs with no external consumers (lightweight interface comments may suffice).
- When the API is still in active design and will change significantly (document after the interface stabilizes).
- When SDK auto-generation from the OpenAPI spec will cover developer needs without supplementary prose.
