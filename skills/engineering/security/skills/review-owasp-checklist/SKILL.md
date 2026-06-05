---
name: review-owasp-checklist
description: Use when reviewing a new codebase for security, before a production launch, after a security incident, or during a compliance audit requiring OWASP Top 10 coverage.
source: OWASP Top 10 2021 (owasp.org/Top10), NIST SP 800-53, OWASP Application Security Verification Standard (ASVS) v4.0
tags: [owasp, security-review, vulnerability, injection, authentication, developer, security-engineer]
verified: true
---

# Review OWASP Checklist

Systematically check a codebase or design against the OWASP Top 10 (2021) to identify the highest-impact vulnerability classes before they reach production.

## Why This Is Best Practice

**Adopted by:** The OWASP Top 10 is referenced by PCI DSS (v4.0, Requirement 6.2.4), NIST SP 800-53 (SA-11), ISO 27001, SOC 2 Type II audits, and is mandated by FedRAMP for web application security testing. Google, Microsoft, and Amazon all include OWASP Top 10 coverage in their internal application security review programs and vendor security questionnaires.
**Impact:** The 2021 OWASP Top 10 categories collectively account for the majority of real-world application breaches. Verizon DBIR 2023 reports that web application attacks (primarily injection, broken access control, and misconfigurations from OWASP Top 10) represent 26% of all breaches. IBM Cost of a Data Breach Report 2023 found the average cost of a breach caused by an exploited vulnerability is $4.45M — preventable with a structured review.
**Why best:** Ad hoc security review produces inconsistent coverage. OWASP Top 10 is community-validated against real-world breach data, updated every 3–4 years, and has tooling support (OWASP ZAP, Semgrep OWASP rulesets, Snyk). It covers the attack surface that attackers actually exploit, not theoretical vulnerabilities.

Sources: OWASP Top 10 2021 (owasp.org/Top10); OWASP ASVS v4.0; Verizon DBIR 2023; IBM Cost of a Data Breach Report 2023; PCI DSS v4.0

## Steps

### 1. Scope the review

Define what is in scope: which services, which endpoints, which data flows. Exclude only what is explicitly out of scope (e.g., third-party SaaS). Document the scope decision.

### 2. Check A01 — Broken Access Control

- Every endpoint: is there an authorization check before returning or mutating data?
- Check for IDOR (Insecure Direct Object Reference): can a user change `?id=123` to `?id=124` and access another user's data?
- Verify CORS policy: `Access-Control-Allow-Origin: *` on authenticated endpoints is a fail.
- Check for path traversal: user-supplied filenames used in file operations without sanitization.
- Admin/privileged routes: are they protected by role checks, not just hidden from the UI?

```bash
# Grep for common IDOR patterns (Python/JS examples)
grep -r "request.params\|req.params\|request.GET" --include="*.py" --include="*.js" | grep -v "auth\|permission\|authorize"
```

### 3. Check A02 — Cryptographic Failures

- Is sensitive data (PII, credentials, payment data) encrypted at rest and in transit?
- Are passwords hashed with bcrypt, argon2, or scrypt? Fail if MD5, SHA-1, or SHA-256 without salt.
- Are TLS 1.0/1.1 disabled? Check server config.
- Are API keys, tokens, or secrets stored in source code or environment-variable files committed to git?

```bash
# Detect secrets in git history
git log --all --full-history -- '*.env' '*.key' '*.pem'
grep -r "password\s*=\s*['\"]" --include="*.py" --include="*.js" --include="*.ts"
```

### 4. Check A03 — Injection

SQL, NoSQL, command, LDAP, and template injection all qualify.

- Is every external input parameterized or escaped before use in a query?
- ORM queries with raw string concatenation are a fail: `f"SELECT * FROM users WHERE name='{name}'"`.
- Shell commands using user input without `shlex.quote` or equivalent are a fail.
- Template engines: is user input rendered with `{{ }}` (auto-escaped) or `{% raw %}` / `|safe` (unescaped)?

```bash
# Grep for raw SQL concatenation
grep -rn "execute.*%" --include="*.py" | grep -v "#"
grep -rn "query.*\+\|query.*\`" --include="*.js" --include="*.ts"
```

### 5. Check A04 — Insecure Design

This is a design-level check, not code-level:
- Is there a threat model for this system? (See design-threat-model skill.)
- Are rate limits in place on authentication, password reset, and payment endpoints?
- Are business logic limits enforced server-side (e.g., cannot order negative quantity, cannot apply discount twice)?

### 6. Check A05 — Security Misconfiguration

- Default credentials changed? Admin panels, database users, cloud storage buckets.
- Debug mode disabled in production? (`DEBUG=True` in Django, `app.run(debug=True)` in Flask.)
- Directory listing disabled on web server.
- Error messages: do they expose stack traces, SQL schema, or internal paths to users?
- HTTP security headers present: `Content-Security-Policy`, `X-Frame-Options`, `Strict-Transport-Security`, `X-Content-Type-Options`.

```bash
curl -I https://yourdomain.com | grep -E "X-Frame|Content-Security|Strict-Transport|X-Content-Type"
```

### 7. Check A06 — Vulnerable and Outdated Components

```bash
# Check for known CVEs in dependencies
npm audit                  # Node.js
pip-audit                  # Python
mvn dependency-check:check # Java (OWASP Dependency Check)
bundle audit               # Ruby
```

Any Critical or High CVE with a fixed version available = immediate action required.

### 8. Check A07 — Identification and Authentication Failures

- Passwords: minimum 8 chars, check against breach databases (e.g., HaveIBeenPwned API).
- Session tokens: cryptographically random, ≥128 bits, invalidated on logout.
- MFA available on all privileged accounts?
- Password reset flow: does it use time-limited, single-use tokens? Reusable reset links are a fail.
- Account lockout or CAPTCHA after N failed login attempts.

### 9. Check A08 — Software and Data Integrity Failures

- CI/CD pipeline: are build artifacts signed? Are dependencies pinned to exact versions with hash verification?
- Deserialization: is user-supplied serialized data (pickle, Java serialization, XML) deserialized without validation? This is a critical fail.
- Auto-update mechanisms: do they verify signatures before applying updates?

### 10. Check A09 — Security Logging and Monitoring Failures

- Are authentication events (login, logout, failed login) logged with timestamp, user ID, and IP?
- Are authorization failures (403) logged?
- Are logs shipped to a tamper-resistant store (not the same host)?
- Is there an alert on > N failed logins within M minutes?

### 11. Check A10 — Server-Side Request Forgery (SSRF)

- Any endpoint that accepts a URL and fetches it server-side?
- Validate that the resolved IP is not in RFC 1918 space (10.x, 172.16.x, 192.168.x) or 169.254.x (AWS metadata).
- Use an allowlist of permitted domains, not a blocklist.

### 12. Produce a findings report

For each finding:

```
Category: A01 Broken Access Control
Location: /api/v1/orders/{id} (GET)
Severity: High
Finding: No authorization check — any authenticated user can retrieve any order by ID.
Proof: GET /api/v1/orders/9999 with user token for order owner of order 1001 returns 200.
Remediation: Add owner check: assert order.user_id == current_user.id before returning.
Owner: @backend-team
```

## Rules

- Findings without a reproduction step are not findings — they are hypotheses
- Severity must be assigned (Critical / High / Medium / Low / Informational) — never just "vulnerability found"
- Every Critical and High finding blocks release until remediated or formally risk-accepted by a named owner
- Do not rely solely on automated scanners — A04 (Insecure Design) and A01 (Broken Access Control) require manual review
- Recheck mitigations before sign-off — verify the fix in the actual codebase, not just the PR description

## Common Mistakes

- **Treating scanner output as complete** — automated tools catch A03/A06/A05 well but miss A01, A04, and most A07 issues entirely.
- **Checking only the happy path** — attackers probe edge cases. Test with negative IDs, empty strings, Unicode, oversized inputs.
- **Marking "low severity" findings as done** — low-severity issues combine. An information disclosure (A09) + a weak token (A07) can chain into account takeover.
- **No retesting** — verify fixes. Developers sometimes misunderstand findings and apply incomplete patches.
- **Ignoring third-party integrations** — OAuth callbacks, webhook receivers, and payment processor redirects are frequent A01/A08 failure points.
