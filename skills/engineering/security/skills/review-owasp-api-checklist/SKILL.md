---
name: review-owasp-api-checklist
description: Use when auditing an API for security vulnerabilities, conducting a pre-launch security review, performing a penetration test on API endpoints, or assessing API security posture against the OWASP API Security Top 10 (2023).
source: 'OWASP API Security Top 10 2023 (owasp.org/www-project-api-security/); OWASP API Security Testing Guide; CWE-285; CWE-770'
tags: [security, owasp, api-security, audit, review, penetration-testing, developer]
---

# Review OWASP API Security Checklist

Systematically check an API against all 10 OWASP API Security Top 10 (2023) vulnerability classes — producing a finding list with severity, evidence, and remediation references.

## Why This Is Best Practice

**Adopted by:** OWASP API Security Top 10 2023 is the authoritative reference for API security, used by security teams at Google, Microsoft, Netflix, and financial institutions globally. The Salt Security API Security Report (2023) found that 94% of organizations experienced API security incidents in the past year, and 59% had API-related breaches. Gartner predicts API abuse will be the most frequent attack vector by 2025. PCI DSS v4.0 Requirement 6.2.4 requires protection against all OWASP Top 10 vulnerability classes including API-specific ones.
**Impact:** Broken Object Level Authorization (API1) is the most common API vulnerability, responsible for the Peloton data exposure (2021), Experian API leak (2021), and countless undisclosed incidents. The entire OWASP API Top 10 represents the vulnerabilities most likely to be exploited in real-world attacks. A structured checklist review catches 80%+ of API security issues before production — at a fraction of the cost of post-breach remediation.
**Why best:** Ad-hoc API security review misses systematic vulnerability classes. This checklist-based approach — parallel to the OWASP Top 10 web checklist (`review-owasp-checklist`) — ensures every API vulnerability class is explicitly checked, not just the ones the reviewer happens to remember.

Sources: OWASP API Security Top 10 2023; Salt Security API Security Report 2023; Gartner API Security research; CWE-285

## Steps

Work through each API Top 10 item. For each: describe what to look for, how to test, and what evidence to record.

---

### API1:2023 — Broken Object Level Authorization (BOLA)

**What:** API endpoints that accept object IDs (user IDs, document IDs, order IDs) without verifying the requester owns or is authorized to access that object.

**Test:**
1. Authenticate as User A, note any object IDs in responses (`/api/orders/1234`).
2. Authenticate as User B, attempt to access User A's object IDs directly.
3. Try sequential IDs: if User B accessed `/orders/1234`, try `/orders/1233`, `/orders/1235`.

**Finding:** If User B can read, modify, or delete User A's objects, BOLA is confirmed.
**Fix reference:** `apply-authorization-control` — object-level ownership checks.

---

### API2:2023 — Broken Authentication

**What:** Weak authentication mechanisms — missing token validation, JWT `alg:none`, accepting expired tokens, weak credential policies.

**Test:**
1. Send requests without `Authorization` header to authenticated endpoints.
2. Modify JWT payload (change `user_id`) without updating signature — if accepted, signature not validated.
3. Use an expired token — if accepted, `exp` not validated.
4. Test brute force on login endpoint — no lockout after N attempts.

**Fix references:** `apply-jwt-security`, `design-session-management`, `apply-api-rate-limiting`.

---

### API3:2023 — Broken Object Property Level Authorization

**What:** APIs returning more data than the caller should see, or accepting writes to protected fields (mass assignment).

**Test:**
1. Inspect API responses for fields the user shouldn't see: `password_hash`, `internal_notes`, `credit_card_number`.
2. On update endpoints: try adding unexpected fields: `{"email": "new@test.com", "role": "admin", "is_verified": true}`.
3. Check if the unexpected fields are silently accepted.

**Fix references:** `prevent-mass-assignment`, `apply-authorization-control` (field-level).

---

### API4:2023 — Unrestricted Resource Consumption

**What:** No rate limits on expensive operations — allows DoS, brute force, and resource exhaustion.

**Test:**
1. Send 100+ requests to login/registration endpoints in rapid succession — check for 429 responses.
2. Test endpoints that trigger background jobs, emails, or external API calls — can they be spammed?
3. Test for unlimited pagination: `?page=999999&per_page=9999`.
4. Upload a large file — is there a size limit?

**Fix reference:** `apply-api-rate-limiting`.

---

### API5:2023 — Broken Function Level Authorization

**What:** User-level callers can access admin/elevated-privilege functions by guessing endpoint names.

**Test:**
1. As a regular user, try admin endpoints: `/api/admin/users`, `/api/admin/stats`, `/api/internal/config`.
2. Try HTTP method switching: GET `/api/users/123` works — try DELETE `/api/users/123` without admin role.
3. Try accessing undocumented endpoints by fuzzing: `/api/v1/users/export`, `/api/v1/users/bulk-delete`.

**Fix reference:** `apply-authorization-control` — function-level permission checks.

---

### API6:2023 — Unrestricted Access to Sensitive Business Flows

**What:** APIs that can be abused to perform business actions at unrealistic volume — ticket scalping, coupon stacking, fake reviews, account creation at scale.

**Test:**
1. Identify business-critical flows: checkout, referral bonus, coupon redemption.
2. Attempt to automate the flow at high volume — can 1000 referral bonuses be claimed from one account?
3. Check if there are bot detection measures (CAPTCHAs, device fingerprinting).

**Fix:** Implement business-logic rate limits, CAPTCHA, device fingerprinting, and anomaly detection on critical flows.

---

### API7:2023 — Server Side Request Forgery (SSRF)

**What:** API endpoints that fetch URLs or make HTTP requests based on user-supplied input.

**Test:**
1. Find any endpoint that accepts a URL: webhook URLs, image import URLs, link preview.
2. Supply internal URLs: `http://169.254.169.254/latest/meta-data/`, `http://localhost:8080/admin`.
3. Check if responses contain metadata or internal service data.

**Fix reference:** `prevent-ssrf`.

---

### API8:2023 — Security Misconfiguration

**What:** Default credentials, unnecessary features enabled, missing security headers, verbose error messages exposing internals.

**Test:**
1. Check response headers for security headers (`X-Content-Type-Options`, `Content-Security-Policy`).
2. Trigger a 500 error — does the response include a stack trace, file path, or framework version?
3. Check for default credentials on admin interfaces.
4. Check HTTP methods: does `OPTIONS` reveal unexpected allowed methods?

**Fix references:** `apply-security-headers`, `design-error-handling`.

---

### API9:2023 — Improper Inventory Management

**What:** Outdated or undocumented API versions still exposed; debug/development endpoints in production.

**Test:**
1. Try deprecated versions: `/api/v1/`, `/api/v2/` alongside the current `/api/v3/`.
2. Test if old versions have weaker security than current (missing auth, missing validation).
3. Check for developer/test endpoints: `/api/test`, `/api/debug`, `/api/echo`.

**Fix:** Maintain an API version inventory; actively retire old versions; remove test endpoints.

---

### API10:2023 — Unsafe Consumption of APIs

**What:** Your API trusts data from third-party APIs without validation — passing it to users, databases, or internal systems.

**Test:**
1. Trace what data your API fetches from upstream services.
2. Check if that data is validated/sanitized before use.
3. Check if upstream API credentials are hardcoded or logged.

**Fix reference:** `apply-api-client-security`.

---

## Rules

- Document findings with: vulnerability class, endpoint, test payload, response evidence, CVSS score, fix reference.
- Re-test after fixes — BOLA fixes often leave related endpoints unpatched.
- API authentication bypass (API2) blocks all further testing — fix it first before testing authorization logic.
- Use both authenticated (as different user roles) and unauthenticated sessions during testing.

## Common Mistakes

- **Only testing the documented API** — BOLA and function-level authorization issues often exist on undocumented or internal endpoints.
- **Testing only happy paths** — security testing requires sending unexpected inputs, wrong IDs, wrong roles.
- **Marking as "low risk" because it requires authentication** — BOLA only requires being authenticated as any user, not as the target user.
