---
name: design-api-versioning-strategy
description: Use when designing public or partner APIs that must evolve without breaking existing integrations, or when API changes are causing client disruption
source: Microsoft REST API Guidelines versioning section; Stripe API versioning model (industry reference); Roy Fielding REST dissertation (2000)
tags: [architecture, api, versioning, rest]
verified: true
---

# Design API Versioning Strategy

Establish a versioning policy that allows APIs to evolve and add capabilities without breaking existing consumers.

## Why This Is Best Practice

**Adopted by:** Stripe (date-based versioning, industry gold standard), GitHub (version in URL path), Google APIs (version in URI), Twilio (URL versioning with LTS commitments)
**Impact:** Stripe's versioning strategy is cited as enabling 10+ years of API evolution without forced migrations; unversioned APIs that break consumers cause integration incidents that cost enterprise customers 40+ engineering hours per incident
**Why best:** APIs are contracts; breaking changes without versioning destroy trust and require coordinated cross-team upgrades; versioning decouples provider evolution from consumer upgrade timelines

Sources: Microsoft "REST API Guidelines" (2016, github.com/microsoft/api-guidelines); Stripe API versioning documentation; Fielding REST dissertation (2000)

## Steps

1. **Choose a versioning scheme** — URL path versioning (`/v1/orders`): most common, explicit, easy to route, easy to document. Date-based (`2024-01-01`): Stripe's model, enables precise semantics per version, ideal for continuous evolution. Header versioning (`Accept: application/vnd.api+json; version=1`): keeps URL clean, harder to test in browser. Choose URL path for new APIs unless you model after Stripe's release-date approach.

2. **Define breaking vs. non-breaking changes** — Non-breaking (never require a version bump): adding new optional fields to responses, adding new optional query parameters, adding new endpoints. Breaking (require a version bump): removing or renaming fields, changing field types, changing required parameters, altering authentication requirements. Enforce this distinction in API review.

3. **Establish a versioning policy** — Document: how long each version is supported (minimum 12 months after deprecation notice), how consumers are notified of deprecations (email, HTTP headers, changelog), and what "end-of-life" means (returns 410 Gone or is disabled). Publish this policy before releasing v1. Stripe's policy: versions are supported indefinitely unless the version has negligible usage.

4. **Implement version routing** — Route requests by version at the API gateway or load balancer layer, not in application code. This allows independent deployment of version implementations. If using URL path versioning, the path prefix (`/v1/`, `/v2/`) is the routing key. Version routing at the gateway enables zero-downtime version migrations.

5. **Use semantic versioning for major, minor, patch** — Major (v1 → v2): breaking changes. Minor (v1.1): additive, non-breaking changes. Patch: bug fixes. For REST APIs, expose only the major version in the URL; minor/patch changes are transparent to consumers. Document minor changes in the API changelog.

6. **Apply the tolerant reader pattern to consumers** — API consumers must ignore unknown fields (don't fail on extra fields in the response). Consumers that are strict about response schemas will break on non-breaking additive changes. Enforce this in integration tests: add extra fields to mock responses and verify consumers handle them.

7. **Deprecate before removing** — Before removing any endpoint or field: add a `Deprecation` header to responses (`Deprecation: true; rel="deprecation"; type="text/html"`), log a warning in the API docs, email all registered consumers, and set a minimum 12-month sunset timeline. Never remove without a deprecation period.

8. **Maintain a machine-readable changelog** — Use OpenAPI (Swagger) or GraphQL schema to document every version. Generate diffs between versions automatically (oasdiff for OpenAPI, graphql-inspector). Consumers can programmatically detect changes. This is the foundation for automated compatibility testing.

9. **Test version compatibility in CI** — Run contract tests (Pact or schema validation) between your API versions and consumer versions. The CI pipeline must catch breaking changes before they reach production. Use API linting tools (Spectral, Optic) to enforce non-breaking change policy on every PR.

10. **Sunset retired versions gracefully** — 6 months before EOL: increase deprecation warnings, reach out to high-usage consumers. 1 month before EOL: return `HTTP 301` redirects to the current version with documentation. On EOL date: return `HTTP 410 Gone` with a migration guide URL. Never return 404 for sunsetted versions; 410 signals intentional retirement.

## Rules

- Breaking changes require a version bump; no exceptions; this is a contractual commitment to consumers.
- The deprecation notice must be publicly documented and consumers must have been notified before any endpoint is removed.
- Never version a single endpoint (`/v2/orders/create`) while leaving others unversioned; version the API, not individual endpoints.
- Consumers that build with strict schema parsing are fragile; document the tolerant reader expectation in your API guide.

## Common Mistakes

- **Version in the query string** (`?version=1`) — query params are for filtering data, not API versions; hard to cache, easy to omit; URL path or header versioning is preferred.
- **Skipping deprecation for "internal" APIs** — internal consumers break on silent removal just like external consumers; all APIs need versioning and deprecation procedures.
- **Version inflation** — creating a new version for every additive change; additive changes are non-breaking and don't require a version bump; version only for breaking changes.
- **No consumer registry** — removing a version without knowing who uses it; maintain an API consumer registry so deprecation notices reach actual integrators.

## When NOT to Use

- Private APIs consumed by only one team where both sides can coordinate changes directly and simultaneously
- GraphQL APIs where schema evolution mechanisms (deprecated fields, union types) replace versioning
- Event-driven APIs where schema registry and forward/backward compatibility handle evolution
