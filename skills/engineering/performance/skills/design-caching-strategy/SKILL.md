---
name: design-caching-strategy
description: Use when designing or auditing HTTP caching, CDN configuration, or application-level cache policies for a web service
source: "RFC 7234 HTTP/1.1 Caching; MDN Cache-Control documentation; Fastly Caching Best Practices; 'Caching at Scale' (Varnish Software)"
tags: [performance, caching, cache-control, cdn, http, stale-while-revalidate, headers]
verified: true
---

# Design Caching Strategy

Configure cache-control directives and CDN rules to serve assets at maximum speed while keeping content fresh enough to be correct.

## Why This Is Best Practice

**Adopted by:** Fastly, Cloudflare, and Akamai — serving trillions of requests using these patterns; RFC 7234 is the HTTP/1.1 standard implemented by every browser and CDN globally

**Impact:** Cloudflare reports CDN cache hits serving responses in <10ms vs. 100-400ms for origin fetches; Facebook's static asset caching strategy (immutable + content-addressed URLs) eliminates cache invalidation for static content entirely

**Why best:** Caching decisions are permanent performance wins: unlike code optimizations, a correctly cached response costs nothing to serve regardless of origin server load. The strategy must match the content's mutability — aggressive caching of mutable content breaks correctness; conservative caching of immutable content wastes the performance budget.

## Steps

1. **Classify content by mutability** — immutable (versioned JS/CSS/images), semi-mutable (API responses, user-generated content), and mutable (personalized pages, real-time data); each class needs a different strategy
2. **Set immutable asset caching** — for content-addressed files (e.g., `app.a3f2c.js`): `Cache-Control: public, max-age=31536000, immutable`; the URL changes on content change, so max TTL is safe
3. **Set API response caching** — for frequently read, infrequently changed data: `Cache-Control: public, s-maxage=300, stale-while-revalidate=60`; CDN serves stale content while fetching fresh in background, eliminating user-visible latency on revalidation
4. **Set personalized content headers** — for user-specific responses: `Cache-Control: private, no-store` or `Vary: Cookie, Authorization` to prevent CDN from serving one user's data to another
5. **Configure CDN cache keys** — strip unnecessary query parameters from cache keys to improve hit ratio; normalize URL variants (trailing slashes, case) before caching
6. **Implement cache invalidation** — use surrogate keys (Fastly) or cache tags (Cloudflare) to invalidate groups of cached objects by content type; do not rely on TTL expiry for time-sensitive invalidation
7. **Measure hit ratio and TTFB** — target CDN cache hit ratio > 85% for static assets; monitor origin TTFB separately from edge TTFB; a rising origin TTFB hidden by cache is a latent incident

## Rules

- Never cache `Set-Cookie` responses at the CDN layer — session cookies cached and served to wrong users is a security incident
- `stale-while-revalidate` requires a CDN or service worker — browsers alone do not implement background revalidation
- Do not use `no-cache` when you mean `no-store` — `no-cache` still stores and revalidates; `no-store` never stores
- Always set explicit `Cache-Control` headers; never rely on CDN defaults, which vary by provider

## Examples

**Static asset:** `/assets/logo.a1b2c3.png` → `Cache-Control: public, max-age=31536000, immutable`
**API response:** `GET /api/products` → `Cache-Control: public, s-maxage=60, stale-while-revalidate=30`
**User dashboard:** `GET /dashboard` → `Cache-Control: private, no-store`

## Common Mistakes

- Caching by file extension instead of content type: a `.json` endpoint with user data should not be cached public regardless of extension
- Setting `max-age=0` thinking it disables caching: it sets TTL to zero, meaning the resource is immediately stale but still cacheable and revalidatable
- Forgetting to vary on `Accept-Encoding`: serving a gzip response to a client that sent no `Accept-Encoding` breaks parsing

## When NOT to Use

- The endpoint returns financial transaction records, healthcare data, or other content where serving a stale response even briefly could cause incorrect decisions — correctness requirements override caching benefits and `no-store` must be used with no exceptions.
- The service does not yet have a CDN in front of it and the team has no near-term plans to add one — HTTP caching strategy design for `s-maxage` and `stale-while-revalidate` requires CDN infrastructure to take effect; without it, only browser-level caching applies.
- The performance bottleneck has been identified as database query latency on cache misses, not high request volume — application-level query result caching (Redis, Memcached) should be designed instead of HTTP caching headers, which do not reduce database load on the first uncached request.
