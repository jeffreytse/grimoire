---
name: review-authentication-security
description: Use when penetration testing or auditing authentication mechanisms — executing WSTG OTG-AUTHN test cases for username enumeration, brute force, default credentials, weak lockout, MFA bypass, and password reset flaws.
source: 'OWASP Web Security Testing Guide v4.2 OTG-AUTHN-001 through OTG-AUTHN-010 (owasp.org/www-project-web-security-testing-guide/); OWASP Authentication Cheat Sheet; NIST SP 800-63B'
tags: [security, owasp, wstg, authentication, penetration-testing, enumeration, brute-force, appsec]
---

# Review Authentication Security

Execute WSTG OTG-AUTHN test cases against a web application — testing for username enumeration, credential brute force, default credentials, weak lockout policy, MFA bypass, and password reset vulnerabilities — producing findings with HTTP evidence.

## Why This Is Best Practice

**Adopted by:** OWASP WSTG v4.2 OTG-AUTHN is the authoritative test methodology for web authentication review — used by CREST-certified penetration testers, OSCP practitioners, and application security assessors globally. CREST (Council of Registered Ethical Security Testers) recommends WSTG as the primary reference methodology for web application testing engagements. HackerOne's 2023 Hacker Report found authentication vulnerabilities in the top 3 most common critical findings in bug bounty programs, with account takeover being the highest-paying bug class.
**Impact:** Verizon DBIR 2023 found that 74% of breaches involve a human element, with credential-based attacks accounting for 49% of all breach vectors — authentication flaws are the single highest-ROI target for attackers. OWASP WSTG-structured authentication review catches the full class of auth vulnerabilities (not just SQLi-in-login), including account lockout gaps, enumeration leaks, and MFA bypasses that automated scanners miss. Organizations that perform structured WSTG authentication testing before launch find 3× more authentication vulnerabilities than those relying on automated scanning alone (SANS Institute Application Security Survey, 2022).
**Why best:** Automated scanners test for SQLi and XSS in login forms but cannot test business-logic flaws: "does this app lock out after 10 failed attempts?", "does the password reset token expire?", "does the app reveal whether an email exists?" — these require WSTG-style manual test cases executed with HTTP clients. The alternative (automated scan only) misses the majority of authentication vulnerabilities that lead to account takeover.

Sources: OWASP WSTG v4.2 OTG-AUTHN-001–010; CREST penetration testing methodology; Verizon DBIR 2023; SANS Application Security Survey 2022

## Steps

### OTG-AUTHN-001 — Test for credentials transported over encrypted channel

```bash
# Verify login does not send credentials over HTTP
curl -v -s -o /dev/null http://app.company.com/login 2>&1 | grep -i "location\|https"
# PASS: 301/302 redirect to HTTPS before credentials transmitted
# FAIL: 200 response with login form on HTTP

# Verify no credentials in URL (GET parameters)
# Look for ?password= or ?token= in access logs or referer header
curl -sI https://app.company.com/login | grep -i "referer\|referrer-policy"
# PASS: Referrer-Policy: strict-origin-when-cross-origin (or stricter)
```

### OTG-AUTHN-002 — Test for default credentials

```bash
# Test common default admin credentials
CREDENTIALS=(
  "admin:admin" "admin:password" "admin:admin123"
  "administrator:administrator" "root:root" "test:test"
  "user:user" "guest:guest" "demo:demo"
)

for cred in "${CREDENTIALS[@]}"; do
  user="${cred%%:*}"
  pass="${cred##*:}"
  status=$(curl -s -o /dev/null -w "%{http_code}" \
    -d "username=${user}&password=${pass}" \
    https://app.company.com/api/auth/login)
  redirect=$(curl -s -o /dev/null -w "%{redirect_url}" \
    -d "username=${user}&password=${pass}" \
    -L https://app.company.com/login)
  echo "${cred}: HTTP ${status} | redirect=${redirect}"
  # FAIL: HTTP 200 with session cookie, or redirect to dashboard
done
```

### OTG-AUTHN-003 — Test account lockout mechanism

```bash
TARGET="victim@company.com"
WRONG_PASS="definitely_wrong_password_xyz"
LOCKOUT_TEST_URL="https://app.company.com/api/auth/login"

# Send 20 failed attempts and monitor when lockout occurs
for i in $(seq 1 20); do
  response=$(curl -s -w "\n%{http_code}" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"${TARGET}\",\"password\":\"${WRONG_PASS}_${i}\"}" \
    "${LOCKOUT_TEST_URL}")
  http_code=$(echo "$response" | tail -1)
  body=$(echo "$response" | head -1)
  echo "Attempt ${i}: HTTP ${http_code} | $(echo $body | jq -r '.message // .error // "no message"')"
done

# PASS: 429 or account_locked response before attempt 10
# FAIL: All 20 attempts return same error (no lockout), or HTTP 200 on correct pass after all attempts
```

### OTG-AUTHN-004 — Test for username/email enumeration

```bash
# Compare responses for existing vs. non-existing user
EXISTING_USER="known_user@company.com"
FAKE_USER="definitely_not_registered_xyzabc@company.com"
WRONG_PASS="wrong_password_123"

response_existing=$(curl -s -w "\n%{http_code}" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${EXISTING_USER}\",\"password\":\"${WRONG_PASS}\"}" \
  https://app.company.com/api/auth/login)

response_fake=$(curl -s -w "\n%{http_code}" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${FAKE_USER}\",\"password\":\"${WRONG_PASS}\"}" \
  https://app.company.com/api/auth/login)

echo "Existing user response:"
echo "$response_existing"
echo "Non-existing user response:"
echo "$response_fake"

# PASS: Identical response body and HTTP code (e.g., both: "Invalid credentials")
# FAIL: Different messages ("Password incorrect" vs "Account not found"), timing difference >50ms, different HTTP codes
```

**Timing-based enumeration check:**
```bash
# Measure response time for existing vs. non-existing user
time_existing=$(curl -s -o /dev/null -w "%{time_total}" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${EXISTING_USER}\",\"password\":\"${WRONG_PASS}\"}" \
  https://app.company.com/api/auth/login)

time_fake=$(curl -s -o /dev/null -w "%{time_total}" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${FAKE_USER}\",\"password\":\"${WRONG_PASS}\"}" \
  https://app.company.com/api/auth/login)

echo "Existing: ${time_existing}s | Non-existing: ${time_fake}s"
# FAIL: >100ms consistent difference (bcrypt not applied to fake users = timing side-channel)
```

### OTG-AUTHN-005 — Test for vulnerable password reset

```bash
# Step 1: Request password reset and capture token format from email
# (Observe token in reset URL: length, charset, predictability)
curl -s -H "Content-Type: application/json" \
  -d "{\"email\":\"${EXISTING_USER}\"}" \
  https://app.company.com/api/auth/forgot-password

# Observe token: short tokens (<128 bits entropy) or sequential tokens FAIL

# Step 2: Verify token expires after 15 minutes
# (Request token, wait, attempt use after expiry)
EXPIRED_TOKEN="expired_token_captured_earlier"
curl -s -H "Content-Type: application/json" \
  -d "{\"token\":\"${EXPIRED_TOKEN}\",\"password\":\"NewPassword123!\"}" \
  https://app.company.com/api/auth/reset-password
# PASS: {"error": "token expired/invalid"}
# FAIL: {"success": true} (expired token accepted)

# Step 3: Verify token is single-use (cannot be reused after successful reset)
VALID_TOKEN="valid_token_from_email"
# Use token once:
curl -s -d "{\"token\":\"${VALID_TOKEN}\",\"password\":\"NewPassword123!\"}" \
  https://app.company.com/api/auth/reset-password
# Use same token again:
curl -s -d "{\"token\":\"${VALID_TOKEN}\",\"password\":\"AnotherPassword456!\"}" \
  https://app.company.com/api/auth/reset-password
# PASS: Second use returns error
# FAIL: Second use returns success (token reuse allowed)

# Step 4: Test Host header injection in reset email
curl -s -H "Content-Type: application/json" \
  -H "Host: attacker.com" \
  -d "{\"email\":\"${EXISTING_USER}\"}" \
  https://app.company.com/api/auth/forgot-password
# PASS: Reset link uses hardcoded domain, not Host header
# FAIL: Reset link uses attacker.com (Host header injection)
```

### OTG-AUTHN-006 — Test for weak password policy

```bash
# Test password policy enforcement
WEAK_PASSWORDS=(
  "password" "12345678" "qwerty123" "letmein1"
  "abc" "a"  # too short
  "aaaaaaaaaa"  # no complexity
)

for pass in "${WEAK_PASSWORDS[@]}"; do
  status=$(curl -s -w "%{http_code}" -o /dev/null \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"test@test.com\",\"password\":\"${pass}\",\"password_confirm\":\"${pass}\"}" \
    https://app.company.com/api/auth/register)
  echo "Password '${pass}': HTTP ${status}"
  # PASS (for weak passwords): HTTP 400 "Password too weak"
  # FAIL: HTTP 200/201 — weak password accepted
done
```

### OTG-AUTHN-007 — Test MFA bypass

```bash
# Verify MFA is enforced after password authentication
# Step 1: Authenticate with valid credentials (will get MFA challenge)
SESSION=$(curl -s -c session.txt \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${EXISTING_USER}\",\"password\":\"ValidPass123!\"}" \
  https://app.company.com/api/auth/login | jq -r '.partial_token // .session_token')

# Step 2: Try accessing protected resource BEFORE completing MFA
curl -s -H "Authorization: Bearer ${SESSION}" \
  https://app.company.com/api/users/me
# PASS: 401/403 "MFA required"
# FAIL: 200 with user data (MFA step can be skipped)

# Step 3: Verify MFA codes expire (TOTP window)
OLD_TOTP_CODE="123456"  # expired code from >30 seconds ago
curl -s -H "Authorization: Bearer ${SESSION}" \
  -d "{\"mfa_code\":\"${OLD_TOTP_CODE}\"}" \
  https://app.company.com/api/auth/mfa/verify
# PASS: 400/401 "Invalid or expired code"
# FAIL: 200 (expired code accepted — no time window enforcement)
```

### OTG-AUTHN-008 — Test remember-me functionality

```bash
# Capture remember-me token from login
REMEMBER_ME_COOKIE=$(curl -s -v \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${EXISTING_USER}\",\"password\":\"ValidPass123!\",\"remember_me\":true}" \
  https://app.company.com/api/auth/login 2>&1 | \
  grep -i "set-cookie.*remember\|set-cookie.*persistent" | \
  awk '{print $3}')

echo "Remember-me cookie: ${REMEMBER_ME_COOKIE}"

# Verify token is not predictable (should be ≥128-bit random)
# Length check: hex token should be ≥32 chars, base64 ≥22 chars
echo "${REMEMBER_ME_COOKIE}" | awk -F= '{print length($2)}' | \
  awk '{ if ($1+0 < 32) print "FAIL OTG-AUTHN-008: Short remember-me token (" $1 " chars)" }'
```

### Evidence capture template

```markdown
# Authentication Test Results — [Application] — [Date]

## OTG-AUTHN-001: Credentials over TLS
Status: PASS / FAIL
Evidence: [HTTP response headers showing 301 redirect or lack thereof]

## OTG-AUTHN-003: Account Lockout
Status: PASS / FAIL
Evidence: Attempt 1-10 HTTP 401; Attempt 11+ HTTP 429 (or FAIL: all 401)
Lock threshold: [N] attempts | Lockout duration: [N] minutes

## OTG-AUTHN-004: Username Enumeration
Status: PASS / FAIL
Evidence:
  Existing user: HTTP [code] — "[message]" — [Xms]
  Non-existing user: HTTP [code] — "[message]" — [Xms]

## OTG-AUTHN-005: Password Reset
Status: PASS / FAIL
Evidence: Token length [N] chars | Token expires: [Y/N] | Single-use: [Y/N]

## Summary Table
| Test | Result | Severity if FAIL |
|---|---|---|
| OTG-AUTHN-001 Encrypted transport | PASS | HIGH |
| OTG-AUTHN-002 Default credentials | PASS | CRITICAL |
| OTG-AUTHN-003 Account lockout | FAIL | HIGH |
| OTG-AUTHN-004 Username enumeration | FAIL | MEDIUM |
| OTG-AUTHN-005 Password reset | PASS | HIGH |
| OTG-AUTHN-006 Password policy | PASS | MEDIUM |
| OTG-AUTHN-007 MFA bypass | PASS | CRITICAL |
| OTG-AUTHN-008 Remember-me | PASS | MEDIUM |
```

## Rules

- Always test with two accounts (existing + non-existing) for enumeration tests — single-account testing cannot detect timing or response differences.
- Capture raw HTTP request/response as evidence — "login fails" is not a finding; the HTTP code and body diff is the finding.
- Test lockout with exact count — "appears to lock out" is not a result; establish the exact threshold and lockout duration.
- Do not test with real user credentials from production — create dedicated test accounts with consent for all brute-force tests.

## Common Mistakes

- **Testing only the POST body, not GET parameters** — some reset flows embed tokens in GET params, which appear in server logs and referer headers.
- **Accepting "I see a lockout message" as PASS** — verify the lockout is enforced server-side; some apps show a lockout message while still processing the authentication attempt.
- **Missing MFA bypass via password reset** — resetting a password can bypass MFA if the app doesn't require MFA after password reset; this is a separate test from OTG-AUTHN-007.
- **Not checking email-based enumeration** — forgot-password flow often reveals whether an email exists even when login does not; test both flows.
