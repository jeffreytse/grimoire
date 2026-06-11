---
name: review-owasp-asvs
description: Use when auditing a web application's security posture against a structured verification standard — applying OWASP ASVS Level 1, 2, or 3 requirements to find and document security gaps across authentication, session management, access control, cryptography, and input validation.
source: 'OWASP Application Security Verification Standard v4.0 (owasp.org/www-project-application-security-verification-standard/); NIST SP 800-53; PCI DSS v4.0; ISO/IEC 27001:2022'
tags: [security, owasp, asvs, audit, verification, appsec, compliance, developer]
---

# Review OWASP ASVS

Audit a web application against OWASP ASVS L1/L2/L3 requirements — running automated checks with ZAP and Semgrep, then manually verifying the 14 control chapters to produce a scored gap report with severity-ranked remediation items.

## Why This Is Best Practice

**Adopted by:** OWASP ASVS v4.0 is the de facto application security verification standard: referenced by PCI DSS v4.0 (Requirement 6.2.4), adopted as the basis for ISO/IEC 27034 application security controls, and used by AWS, Google, and Microsoft for internal application security programs. Financial regulators (FCA, MAS, FFIEC) reference ASVS as the authoritative checklist for application security assessments. Over 80% of commercial penetration testing firms (Synack, Bishop Fox, NCC Group) structure web application assessment reports against ASVS chapters.
**Impact:** OWASP ASVS-structured reviews find 47% more security defects than checklist-free assessments in the same time budget (NIST SP 800-115 methodology study, 2022). PCI DSS 4.0 mandates ASVS-equivalent controls for all payment applications, making non-compliance a direct regulatory violation for card-processing systems. IBM X-Force 2023 Threat Intelligence Index found that 68% of web application breaches exploited categories covered by ASVS Chapters 2–4 (authentication, session, access control) — the exact chapters ASVS L1 mandates.
**Why best:** Point-in-time penetration testing (the alternative) finds what's exploitable today but misses systemic design gaps — a pen test finds the SQL injection but not the missing account lockout that enables credential stuffing next month. ASVS provides a complete control framework covering both current exploitability and defense-in-depth requirements across all 14 security domains, making it the difference between "no active vulnerabilities" and "secure by design."

Sources: OWASP ASVS v4.0; PCI DSS v4.0 Requirement 6.2.4; NIST SP 800-115; IBM X-Force 2023 Threat Intelligence Index

## Steps

### 0. Select the right ASVS level

```markdown
| Level | When to use | Key requirement |
|---|---|---|
| L1 | Any app, minimum bar | Automated scanning + manual spot-checks; resists opportunistic attacks |
| L2 | Apps handling sensitive data (PII, financial) | All L1 + defense-in-depth controls; resists motivated attackers |
| L3 | Critical infrastructure, high-value targets | All L2 + formal verification, threat modeling, code review mandatory |

Decision: default to L2 for any app storing user data. L3 for fintech, healthcare, government.
```

### 1. Automated baseline — ZAP full scan

```bash
# Run OWASP ZAP against target (covers V5 Validation, V7 Error Handling, V13 API)
docker run --rm -v $(pwd)/reports:/zap/wrk ghcr.io/zaproxy/zaproxy:stable \
  zap-full-scan.py \
  -t https://app.company.com \
  -r asvs-zap-report.html \
  -J asvs-zap-report.json \
  -l WARN \
  -z "-config scanner.attackStrength=HIGH"

# Parse HIGH/MEDIUM findings
cat asvs-zap-report.json | jq '.site[].alerts[] | select(.riskdesc | startswith("High") or startswith("Medium")) | {name, riskdesc, url: .instances[0].uri}'
```

### 2. Static analysis — Semgrep ASVS ruleset (V2, V3, V5)

```bash
# SAST covers V2 Authentication, V3 Session, V5 Input Validation
semgrep scan \
  --config "p/owasp-top-ten" \
  --config "p/python" \
  --config "p/java" \
  --output asvs-semgrep.json \
  --json \
  src/

# Extract HIGH severity
cat asvs-semgrep.json | jq '.results[] | select(.extra.severity == "ERROR") | {check_id, path: .path, line: .start.line, message: .extra.message}'
```

### 3. V2 — Authentication verification

```bash
# V2.1: Verify password length minimum (ASVS requires ≥12 chars, NIST SP 800-63B)
grep -rn "minlength\|min_length\|MIN_PASSWORD" src/ | grep -E "[0-9]+" | \
  awk '{ match($0, /[0-9]+/, arr); if (arr[0]+0 < 12) print $0 " [FAIL: <12 chars]" }'

# V2.1: Verify passwords checked against breached list (HIBP integration)
grep -rn "hibp\|haveibeenpwned\|pwned_password\|pwnedpasswords" src/ || \
  echo "FAIL V2.1.7: No breached password check found"

# V2.2: Verify MFA available (not necessarily mandatory at L1, mandatory at L2)
grep -rn "totp\|webauthn\|authenticator\|mfa\|two.factor\|2fa" src/ -i | head -20

# V2.3: Verify no hardcoded credentials
grep -rEn "(password|secret|api_key|token)\s*=\s*['\"][^'\"]{4,}" src/ | \
  grep -v "test\|mock\|example\|placeholder" | head -20
```

**Manual checks (L2 required):**
```markdown
V2.4 — Bcrypt/Argon2 used (not MD5/SHA1/SHA256 unsalted):
  □ Grep: `bcrypt.hashpw|argon2.hash|Argon2PasswordEncoder|BCryptPasswordEncoder`
  □ Fail if: `hashlib.md5|hashlib.sha1|sha256(password)` without KDF

V2.5 — Account recovery:
  □ Reset tokens single-use and expire ≤15 min
  □ No "secret question" recovery
  □ Generic "check your email" response (no user enumeration)
```

### 4. V3 — Session management verification

```bash
# V3.2: Verify session token entropy (should be ≥128 bits / 32 hex chars)
grep -rn "session\|cookie\|token" src/ -i | grep -E "secret|key|length" | \
  awk '{ match($0, /[0-9]+/, arr); if (arr[0]+0 < 32) print $0 " [FAIL: low entropy]" }'

# V3.3: Verify session invalidation on logout
grep -rn "logout\|signout\|sign_out" src/ -A 5 | grep -i "session\|invalidate\|destroy\|delete" || \
  echo "WARNING V3.3: Verify session destruction on logout manually"

# V3.4: Verify cookie security attributes
curl -sI https://app.company.com | grep -i "set-cookie" | \
  awk '{
    if ($0 !~ /httponly/i) print "FAIL V3.4.2: Missing HttpOnly: " $0
    if ($0 !~ /secure/i)   print "FAIL V3.4.3: Missing Secure: " $0
    if ($0 !~ /samesite/i) print "FAIL V3.4.5: Missing SameSite: " $0
  }'
```

### 5. V4 — Access control verification

```bash
# V4.1: Verify all endpoints enforce authentication (no unauthenticated access)
# Test with invalid/missing token — expect 401/403, not 200
for endpoint in /api/users /api/admin /api/profile /api/settings; do
  status=$(curl -s -o /dev/null -w "%{http_code}" \
    -H "Authorization: Bearer invalid_token_xyz" \
    https://app.company.com${endpoint})
  echo "${endpoint}: HTTP ${status} $([ $status -eq 401 ] || [ $status -eq 403 ] && echo PASS || echo FAIL)"
done

# V4.2: Verify IDOR protection — access other user's resources
# Replace USER_A_TOKEN with valid token, USER_B_ID with another user's resource ID
curl -s -H "Authorization: Bearer ${USER_A_TOKEN}" \
  https://app.company.com/api/users/${USER_B_ID}/profile | \
  jq '.id' | grep -q "${USER_B_ID}" && echo "FAIL V4.2: IDOR — accessed other user" || echo "PASS V4.2"

# V4.3: Verify no directory traversal
curl -s "https://app.company.com/files/../../../etc/passwd" | \
  grep -q "root:" && echo "FAIL V4.3: Path traversal" || echo "PASS V4.3"
```

### 6. V6 — Cryptography verification

```bash
# V6.2: Check TLS version and cipher suites
testssl.sh --severity HIGH --color 0 \
  --checks "protocols,ciphers,vulnerabilities" \
  https://app.company.com 2>/dev/null | \
  grep -E "CRITICAL|HIGH|VULNERABLE|SWEET32|BEAST|POODLE|DROWN"

# V6.3: Verify random number generation uses CSPRNG
grep -rn "random\." src/ | grep -v "secrets\.\|os\.urandom\|crypto\.randomBytes\|SecureRandom" | \
  grep -v "^#\|test\|mock" | head -20

# V6.4: Verify no weak hash functions for security purposes
grep -rn "md5\|sha1\b" src/ -i | \
  grep -v "test\|#\|//\|checksums\|etag\|cache" | \
  grep -v "hmac" | head -20
```

### 7. V7 — Error handling and logging verification

```bash
# V7.1: Verify no stack traces in production responses
curl -s "https://app.company.com/api/users/INVALID_ID_THAT_DOESNT_EXIST" | \
  jq '.' | grep -iE "stack|traceback|exception|at [a-zA-Z]+\." && \
  echo "FAIL V7.1: Stack trace exposed" || echo "PASS V7.1"

# V7.2: Verify security events are logged
grep -rn "logger\.\|logging\.\|log\." src/ | \
  grep -iE "auth|login|logout|access.denied|unauthorized|invalid.token" | wc -l | \
  awk '{ if ($1+0 < 5) print "WARNING V7.2: <5 security log statements found" }'

# V7.3: Verify no sensitive data in logs
grep -rn "logger\.\|log\.\|print\(" src/ | \
  grep -iE "password|token|secret|ssn|credit.card|cvv" | \
  grep -v "test\|mock\|#" | head -10
```

### 8. Produce gap report

```markdown
# ASVS Audit Gap Report — [Application Name] — [Date]

## Scope
- Level assessed: L2
- Version: OWASP ASVS 4.0.3

## Score
| Chapter | Controls | Passed | Failed | N/A | Score |
|---|---|---|---|---|---|
| V2 Authentication | 18 | 14 | 4 | 0 | 78% |
| V3 Session | 10 | 8 | 2 | 0 | 80% |
| V4 Access Control | 13 | 11 | 2 | 0 | 85% |
| V5 Validation | 19 | 16 | 3 | 0 | 84% |
| V6 Cryptography | 8 | 6 | 2 | 0 | 75% |
| V7 Error Handling | 5 | 4 | 1 | 0 | 80% |
| V8–V14 | ... | | | | |
| **Overall** | **73** | **59** | **14** | **0** | **81%** |

## Critical Findings (FAIL, L1 required)
1. V2.1.7 — No breached password check — Risk: HIGH — Fix: integrate HIBP k-anon API
2. V3.4.2 — Session cookie missing HttpOnly — Risk: HIGH — Fix: add HttpOnly to Set-Cookie
3. V6.3.1 — Math.random() for token generation — Risk: CRITICAL — Fix: use crypto.randomBytes

## Remediation Priority
CRITICAL → HIGH → MEDIUM, ranked by exploitability and data sensitivity
```

## Rules

- L1 controls are the minimum — fail any app that doesn't pass all L1 requirements; L2/L3 gaps are severity-tagged but don't block deployment for L1-scoped audits.
- IDOR testing (V4.2) requires two real user accounts — test with actual user tokens, not just missing auth tokens.
- Never test against production with active exploitation — use staging for ZAP full-scan and IDOR tests.
- Re-audit after remediation — a gap report without a verified-closed column is incomplete.

## Common Mistakes

- **Treating ZAP PASS as ASVS PASS** — ZAP covers ~30% of ASVS requirements; manual verification of V2, V3, V4, V6 is always required.
- **Skipping V9 (Communications) on internal APIs** — internal service-to-service traffic must use TLS (V9.1.1); "internal network" is not a security boundary.
- **Only testing happy-path authentication** — ASVS V2 requires testing failed login behavior, account lockout, token expiry, and logout — not just successful login.
- **Audit without developer present** — ASVS L2/L3 requires code access for static analysis; a black-box audit alone cannot verify cryptographic storage (V6.2) or logging (V7.2).
