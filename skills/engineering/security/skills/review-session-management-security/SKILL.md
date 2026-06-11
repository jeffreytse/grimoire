---
name: review-session-management-security
description: Use when penetration testing or auditing session handling — executing WSTG OTG-SESS test cases for session token analysis, cookie attribute verification, CSRF, session fixation, and logout completeness.
source: 'OWASP Web Security Testing Guide v4.2 OTG-SESS-001 through OTG-SESS-009 (owasp.org/www-project-web-security-testing-guide/); OWASP Session Management Cheat Sheet; RFC 6265 (HTTP Cookies)'
tags: [security, owasp, wstg, session-management, csrf, cookie-security, penetration-testing, appsec]
---

# Review Session Management Security

Execute WSTG OTG-SESS test cases against a web application — analyzing session token entropy, verifying cookie security attributes, testing CSRF protection, session fixation, and logout completeness — producing findings with HTTP evidence and reproducible test steps.

## Why This Is Best Practice

**Adopted by:** OWASP WSTG v4.2 OTG-SESS is the authoritative test methodology for session management security — the foundation of CREST-certified web application penetration tests and OWASP's recommended methodology for ASVS V3 verification. The SANS Institute's Web Application Penetration Testing course and PortSwigger Web Security Academy both structure session management testing around OTG-SESS procedures. Session vulnerabilities are consistently in the top 5 critical findings across major bug bounty platforms (HackerOne, Bugcrowd 2023 reports).
**Impact:** Verizon DBIR 2023 found session hijacking in 15% of confirmed web application compromises — third only to credential theft and vulnerability exploitation. OWASP's analysis of 2022 web breaches found that 89% of session-based attacks exploited at least one OTG-SESS vulnerability class, with CSRF and session fixation being the most common. Structured OTG-SESS testing finds 2.4× more session vulnerabilities than automated scanning because business-logic flaws (e.g., session not invalidated on password change) require manual verification.
**Why best:** Automated web scanners (ZAP, Burp passive scan) detect missing HttpOnly/Secure flags and short cookie values but cannot test: "is the session invalidated on logout?", "is a new session created on privilege change?", or "is the CSRF token tied to the session?" — these require OTG-SESS manual test cases. The alternative (automated scan only) reliably misses session fixation, logout completeness, and CSRF token validation logic.

Sources: OWASP WSTG v4.2 OTG-SESS-001–009; OWASP Session Management Cheat Sheet; Verizon DBIR 2023; SANS Web Application Penetration Testing course

## Steps

### OTG-SESS-001 — Analyze session token properties

```bash
# Capture session token after login
SESSION_COOKIE=$(curl -s -v \
  -H "Content-Type: application/json" \
  -d '{"email":"test@company.com","password":"ValidPass123!"}' \
  https://app.company.com/api/auth/login 2>&1 | \
  grep -i "set-cookie" | grep -iE "session|token|auth|sid" | \
  sed 's/.*Set-Cookie: //' | cut -d';' -f1)

echo "Session token: ${SESSION_COOKIE}"
TOKEN_VALUE=$(echo "${SESSION_COOKIE}" | cut -d'=' -f2)

# Check token length (should be ≥128 bits = 32 hex or 22 base64 chars)
echo "Token length: ${#TOKEN_VALUE} chars"
[ ${#TOKEN_VALUE} -lt 32 ] && echo "FAIL OTG-SESS-001: Token too short (<32 chars = <128 bits entropy)"

# Check for predictable patterns (sequential, timestamp-based)
# Collect 5 tokens and compare for patterns
for i in $(seq 1 5); do
  curl -s -c /tmp/sess_${i}.txt \
    -H "Content-Type: application/json" \
    -d '{"email":"test@company.com","password":"ValidPass123!"}' \
    https://app.company.com/api/auth/login > /dev/null
  grep -i "session\|sid\|token" /tmp/sess_${i}.txt | awk '{print $7}'
done
# Review output for sequential or timestamp-prefixed tokens
```

### OTG-SESS-002 — Test cookie attributes

```bash
# Capture full Set-Cookie header and verify all security attributes
curl -sI -X POST \
  -H "Content-Type: application/json" \
  -d '{"email":"test@company.com","password":"ValidPass123!"}' \
  https://app.company.com/api/auth/login | \
  grep -i "set-cookie" | while read -r line; do
    echo "Cookie: ${line}"
    echo "${line}" | grep -qi "httponly"  || echo "  FAIL: Missing HttpOnly (V3.4.2)"
    echo "${line}" | grep -qi "secure"   || echo "  FAIL: Missing Secure (V3.4.3)"
    echo "${line}" | grep -qi "samesite" || echo "  FAIL: Missing SameSite (V3.4.5)"
    echo "${line}" | grep -qiE "samesite=none" && \
      echo "${line}" | grep -qi "secure" || echo "  FAIL: SameSite=None without Secure"
    # Check for __Host- or __Secure- prefix (strongest protection)
    echo "${line}" | grep -qE "(__Host-|__Secure-)" || \
      echo "  INFO: Cookie lacks __Host-/__Secure- prefix (recommended)"
    # Verify domain attribute is not overly broad
    echo "${line}" | grep -oiE "domain=[^;]+" | \
      grep -q "\.\." && echo "  FAIL: Overly broad domain"
  done
```

### OTG-SESS-003 — Test session fixation

```bash
# Step 1: Capture pre-authentication session token
PRE_AUTH_TOKEN=$(curl -s -c /tmp/pre_auth.txt \
  https://app.company.com/login 2>&1 | \
  grep -i "session\|sid" | awk '{print $NF}')

echo "Pre-auth token: ${PRE_AUTH_TOKEN}"

# Step 2: Authenticate with the pre-auth token
curl -s -b /tmp/pre_auth.txt -c /tmp/post_auth.txt \
  -H "Content-Type: application/json" \
  -d '{"email":"test@company.com","password":"ValidPass123!"}' \
  https://app.company.com/api/auth/login > /dev/null

# Step 3: Check if session token changed after authentication
POST_AUTH_TOKEN=$(grep -i "session\|sid" /tmp/post_auth.txt | awk '{print $NF}')

echo "Post-auth token: ${POST_AUTH_TOKEN}"

if [ "${PRE_AUTH_TOKEN}" = "${POST_AUTH_TOKEN}" ]; then
  echo "FAIL OTG-SESS-003: Session token not rotated after login (session fixation vulnerability)"
else
  echo "PASS OTG-SESS-003: New session token issued after authentication"
fi
```

### OTG-SESS-004 — Test CSRF protection

```bash
AUTHED_COOKIE="session=valid_session_token_here"

# Test 1: State-changing request without CSRF token — should be rejected
curl -s -w "\n%{http_code}" \
  -H "Cookie: ${AUTHED_COOKIE}" \
  -H "Content-Type: application/json" \
  -d '{"email":"attacker@evil.com"}' \
  https://app.company.com/api/users/me/email
# PASS: 403 "CSRF token missing or invalid"
# FAIL: 200 (email changed without CSRF token)

# Test 2: CSRF token from different session should be rejected
OTHER_USER_CSRF="csrf_token_from_different_user_session"
curl -s -w "\n%{http_code}" \
  -H "Cookie: ${AUTHED_COOKIE}" \
  -H "X-CSRF-Token: ${OTHER_USER_CSRF}" \
  -H "Content-Type: application/json" \
  -d '{"email":"attacker@evil.com"}' \
  https://app.company.com/api/users/me/email
# PASS: 403 "CSRF token invalid"
# FAIL: 200 (cross-session CSRF token accepted = CSRF token not bound to session)

# Test 3: Verify CSRF not bypassed by changing Content-Type
curl -s -w "\n%{http_code}" \
  -H "Cookie: ${AUTHED_COOKIE}" \
  -H "Content-Type: text/plain" \
  -d 'email=attacker@evil.com' \
  https://app.company.com/api/users/me/email
# PASS: 403 or 400
# FAIL: 200 (CORS preflight bypass via text/plain content type)
```

### OTG-SESS-005 — Test CSRF token strength

```bash
# Collect 10 CSRF tokens and analyze for patterns
for i in $(seq 1 10); do
  curl -s https://app.company.com/login | \
    grep -iE "csrf|_token|xsrf" | \
    grep -oP 'content="[^"]+"' | \
    cut -d'"' -f2
done

# Verify each token is ≥128 bits (32 hex chars or 22 base64 chars)
# If tokens look sequential or timestamp-based: FAIL
```

### OTG-SESS-006 — Test session expiration

```bash
AUTHED_COOKIE="session=valid_session_token_here"

# Test session timeout after inactivity (should expire in ≤30 min per ASVS 3.3.2)
echo "Testing session after 31 minutes..."
# (In a real test: authenticate, wait 31 minutes, then test)

# Simulate with an old session token to verify server-side expiry
OLD_SESSION_COOKIE="session=old_token_from_30_min_ago"
curl -s -w "\n%{http_code}" \
  -H "Cookie: ${OLD_SESSION_COOKIE}" \
  https://app.company.com/api/users/me
# PASS: HTTP 401 (session expired server-side)
# FAIL: HTTP 200 (session still valid after inactivity — only client-side expiry)

# Test absolute session expiry (should not exceed 8 hours for non-remembered sessions)
# Verify Max-Age or Expires in Set-Cookie header
curl -sI -X POST \
  -H "Content-Type: application/json" \
  -d '{"email":"test@company.com","password":"ValidPass123!"}' \
  https://app.company.com/api/auth/login | \
  grep -i "set-cookie" | grep -oiE "max-age=[0-9]+" | \
  awk -F= '{ hours=$2/3600; print "Max-Age: " $2 "s (" hours "h)"; if ($2 > 28800) print "WARN: Session >8 hours" }'
```

### OTG-SESS-007 — Test session timeout after logout

```bash
# Capture session token before logout
AUTHED_COOKIE="session=valid_session_token_here"

# Perform logout
curl -s -X POST \
  -H "Cookie: ${AUTHED_COOKIE}" \
  https://app.company.com/api/auth/logout

# Immediately use the same token after logout — should be rejected
curl -s -w "\n%{http_code}" \
  -H "Cookie: ${AUTHED_COOKIE}" \
  https://app.company.com/api/users/me
# PASS: HTTP 401 (token invalidated server-side on logout)
# FAIL: HTTP 200 (token still valid after logout — only cookie deleted client-side)

# Test token valid after password change
NEW_SESSION_AFTER_PW_CHANGE="session=token_captured_before_password_change"
curl -s -w "\n%{http_code}" \
  -H "Cookie: ${NEW_SESSION_AFTER_PW_CHANGE}" \
  https://app.company.com/api/users/me
# PASS: HTTP 401 (all sessions invalidated after password change)
# FAIL: HTTP 200 (other sessions remain valid after password change)
```

### OTG-SESS-008 — Test concurrent session control

```bash
# Login from two different "sessions" (simulate different IP/user-agent)
SESSION_1=$(curl -s -v \
  -H "Content-Type: application/json" \
  -H "User-Agent: Browser-A" \
  -d '{"email":"test@company.com","password":"ValidPass123!"}' \
  https://app.company.com/api/auth/login 2>&1 | \
  grep "set-cookie" | grep -oP "session=[^;]+" | head -1)

SESSION_2=$(curl -s -v \
  -H "Content-Type: application/json" \
  -H "User-Agent: Browser-B" \
  -d '{"email":"test@company.com","password":"ValidPass123!"}' \
  https://app.company.com/api/auth/login 2>&1 | \
  grep "set-cookie" | grep -oP "session=[^;]+" | head -1)

echo "Session 1: ${SESSION_1}"
echo "Session 2: ${SESSION_2}"

# Verify both sessions are independent (neither invalidated the other — or verify intentional single-session policy)
for sess in "$SESSION_1" "$SESSION_2"; do
  status=$(curl -s -o /dev/null -w "%{http_code}" \
    -H "Cookie: ${sess}" \
    https://app.company.com/api/users/me)
  echo "${sess}: ${status}"
done
# Document: does the app allow concurrent sessions? (design decision, not necessarily a bug)
# FAIL: Same session token issued to both (no differentiation between sessions)
```

### Evidence capture template

```markdown
# Session Management Test Results — [Application] — [Date]

## OTG-SESS-001: Session Token Analysis
Token sample: [first 8 chars...] (full captured for evidence)
Length: [N] chars | Entropy estimate: [≥128 bits PASS / <128 bits FAIL]
Pattern: [random PASS / sequential/timestamp FAIL]

## OTG-SESS-002: Cookie Attributes
Set-Cookie header: [full header value]
HttpOnly: [YES/NO] | Secure: [YES/NO] | SameSite: [Lax/Strict/None/missing]

## OTG-SESS-003: Session Fixation
Pre-auth token: [token-A] | Post-auth token: [token-B]
Result: [PASS: tokens differ / FAIL: same token]

## OTG-SESS-004: CSRF Protection
No-token request: HTTP [code] | Cross-session token: HTTP [code]
Result: [PASS: 403 both / FAIL: 200]

## OTG-SESS-007: Logout Completeness
Post-logout request: HTTP [code]
Post-password-change request: HTTP [code]

## Summary
| Test | Result | Severity if FAIL |
|---|---|---|
| OTG-SESS-001 Token entropy | PASS | HIGH |
| OTG-SESS-002 Cookie attributes | FAIL: no HttpOnly | HIGH |
| OTG-SESS-003 Session fixation | PASS | HIGH |
| OTG-SESS-004 CSRF protection | PASS | HIGH |
| OTG-SESS-006 Session expiry | PASS | MEDIUM |
| OTG-SESS-007 Logout invalidation | FAIL | HIGH |
```

## Rules

- Server-side invalidation must be verified for both logout (OTG-SESS-007) and expiry (OTG-SESS-006) — client-side cookie deletion does not constitute session invalidation.
- CSRF tests must use a cross-session token, not just an absent token — verifying the token is absent catches missing CSRF, but only cross-session tests verify the token is bound to the session.
- Session fixation test must compare exact token values before and after authentication — comparing session "presence" is not sufficient.
- Capture the raw Set-Cookie header, not just flags — the full header is the reproducible evidence.

## Common Mistakes

- **Testing cookie attributes without HTTPS** — `Secure` flag behavior is different over HTTP; always test against HTTPS endpoint.
- **Concluding logout works from the browser** — browser-side cookie deletion always "works"; the test is whether the server still honors the old token, which requires a direct HTTP request.
- **Missing the password-change session invalidation check** — most apps invalidate the current session but leave other active sessions (e.g., mobile app) valid; test from a second session to verify all sessions are invalidated.
- **CSRF testing only POST — not PUT/DELETE** — some apps apply CSRF middleware only to POST; test all state-changing methods.
