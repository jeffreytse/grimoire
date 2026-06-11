---
name: audit-cors-configuration
description: Use when you need to verify that a CORS policy is actually enforced correctly — after configuring CORS, before shipping an API, or during a security review. Tests origin reflection, null origin, subdomain bypass, wildcard-with-credentials, and preflight correctness.
source: 'Portswigger Web Security Academy — CORS vulnerabilities (portswigger.net/web-security/cors); OWASP Testing Guide OTG-CLIENT-007; CWE-942; James Kettle "Exploiting CORS misconfigurations for Bitcoins and bounties" (AppSec EU 2017)'
tags: [security, cors, owasp, audit, testing, api, web, penetration-testing]
verified: true
---

# Audit CORS Configuration

Probe your CORS policy with targeted `curl` and browser tests to confirm it blocks what it should — misconfigured CORS is common because developers add the policy but never verify it.

## Why This Is Best Practice

**Adopted by:** OWASP Testing Guide includes CORS as a required audit item (OTG-CLIENT-007). Portswigger Web Security Academy documents CORS misconfiguration as a systematic exploitable class. Bug bounty programmes at Netflix, Coinbase, and dozens of fintechs have paid critical/high bounties for CORS issues — almost all were configurations that looked correct but were not verified.

**Impact:** James Kettle's 2017 research found CORS misconfigurations on hundreds of the Alexa top 1000 sites, including one that would have allowed theft of Bitcoin from an exchange's authenticated API. The pattern was always the same: CORS was *configured* but never *tested*. `curl` does not enforce the same-origin policy — tests using curl or Postman will show the header is present but will not reveal whether the policy actually restricts cross-origin reads.

**Why best:** Configuring CORS and auditing CORS are different activities. A middleware can be correctly installed but misconfigured. A regex can match too broadly. An environment variable can be set to the wrong value in production. Auditing closes the gap between intent and reality.

Sources: Portswigger CORS research; OWASP Testing Guide OTG-CLIENT-007; CWE-942; Kettle, AppSec EU 2017

## Steps

### Step 1: Check for origin reflection without allowlist validation

Send a clearly untrusted origin. A correct implementation returns no `Access-Control-Allow-Origin` header. A misconfigured one reflects it back.

```bash
curl -s -I \
  -H "Origin: https://evil.com" \
  https://api.example.com/v1/user

# FAIL — server reflects arbitrary origin:
# Access-Control-Allow-Origin: https://evil.com

# PASS — no ACAO header returned, or returns only for trusted origins
```

If the endpoint requires authentication, include a valid session cookie or token — the dangerous case is reflection on *authenticated* endpoints:

```bash
curl -s -I \
  -H "Origin: https://evil.com" \
  -H "Authorization: Bearer <valid_token>" \
  https://api.example.com/v1/user
```

### Step 2: Check for null origin allowance

Sandboxed iframes, `data:` URIs, and redirected cross-origin requests send `Origin: null`. Allowing `null` permits any page to embed your API in a sandboxed iframe and read responses.

```bash
curl -s -I \
  -H "Origin: null" \
  https://api.example.com/v1/user

# FAIL:
# Access-Control-Allow-Origin: null

# PASS: no ACAO header, or no match
```

### Step 3: Check for substring and regex bypass

If your allowlist uses `in` / `contains` / a substring regex, test variants that embed your domain inside an attacker-controlled domain:

```bash
# Tests for: if 'example.com' in origin — matches both
curl -s -I -H "Origin: https://evil-example.com" https://api.example.com/v1/user
curl -s -I -H "Origin: https://notexample.com" https://api.example.com/v1/user
curl -s -I -H "Origin: https://example.com.evil.io" https://api.example.com/v1/user

# All three should return no ACAO header
```

### Step 4: Check wildcard combined with credentials

`Access-Control-Allow-Origin: *` with `Access-Control-Allow-Credentials: true` is blocked by browsers per spec — but some servers produce this combination, which means they will work with non-browser clients and also signals a confused implementation.

```bash
curl -s -I https://api.example.com/v1/public

# FAIL — dangerous combination:
# Access-Control-Allow-Origin: *
# Access-Control-Allow-Credentials: true

# PASS for public endpoint:
# Access-Control-Allow-Origin: *
# (no credentials header, or credentials: false)
```

### Step 5: Verify the preflight response is correctly restricted

OPTIONS preflights must restrict methods and headers to only what your API uses. An overly permissive preflight defeats the purpose of restricting actual requests.

```bash
curl -s -I \
  -X OPTIONS \
  -H "Origin: https://app.example.com" \
  -H "Access-Control-Request-Method: DELETE" \
  -H "Access-Control-Request-Headers: X-Custom-Dangerous-Header" \
  https://api.example.com/v1/user

# FAIL:
# Access-Control-Allow-Methods: GET, POST, PUT, DELETE, PATCH, OPTIONS
# Access-Control-Allow-Headers: *

# PASS:
# Access-Control-Allow-Methods: GET, POST
# Access-Control-Allow-Headers: Content-Type, Authorization
# (method/header not in list = browser will block the actual request)
```

### Step 6: Confirm `Vary: Origin` is set when reflecting dynamically

Without `Vary: Origin`, a CDN or proxy caches one origin's CORS response and serves it to all others — including attacker.com.

```bash
curl -s -I \
  -H "Origin: https://app.example.com" \
  https://api.example.com/v1/resource

# Required when ACAO reflects the request origin:
# Vary: Origin
```

### Step 7: Run a browser-based cross-origin fetch test

`curl` and Postman bypass the same-origin policy — they cannot test whether the browser actually allows or blocks the request. Use a browser fetch from a different origin.

Option A — test from browser console on a trusted page:
```javascript
// Paste into DevTools console on https://app.example.com
fetch('https://api.example.com/v1/user', { credentials: 'include' })
  .then(r => r.json())
  .then(console.log)
  .catch(console.error);
// Expected: data returned (allowed origin)
```

Option B — test from a controlled test page on a different origin:
```html
<!-- host this on https://test.yourlab.com, not on example.com -->
<script>
fetch('https://api.example.com/v1/user', { credentials: 'include' })
  .then(r => console.log('FAIL — cross-origin read succeeded:', r.status))
  .catch(e => console.log('PASS — blocked:', e.message));
</script>
```

### Step 8: Audit all endpoints, not just the ones you think are public

CORS misconfigurations commonly exist on endpoints that developers considered internal or low-risk. Enumerate all routes and test each:

```bash
# Collect all routes (example for Express/Node)
grep -r "app\.\(get\|post\|put\|delete\|patch\)" src/ | grep -oP "['\"](/[^'\"]+)" | sort -u

# Or check your OpenAPI/Swagger spec
cat openapi.yaml | grep -oP "(?<=^  /)[\w/{}]+"
```

Run Steps 1–3 against each endpoint. Pay particular attention to any endpoint that returns user data, session tokens, or account information.

## Rules

- Never conclude CORS is secure from `curl` output alone — `curl` has no same-origin policy and will always receive the response regardless of CORS headers.
- Test authenticated endpoints with a valid token — unauthenticated endpoints returning `*` are often intentional and safe; the dangerous case is always authenticated endpoints.
- Re-audit after any change to origin allowlist logic, middleware version upgrades, or CDN configuration changes — CORS is often broken by infrastructure changes, not just code changes.
- Check CORS on error responses as well as success responses — some frameworks apply CORS middleware only to successful paths, leaving error responses un-headered (which causes confusing CORS errors) or incorrectly headered.

## Common Mistakes

**Trusting that the CORS library handles it correctly without verifying the configuration passed to it:** `CORS(app, origins='*')` is one argument away from `CORS(app, origins=['https://app.example.com'])`. Library is correct; configuration is wrong.

**Only testing from the same origin as the server during development:** `localhost:3000` calling `localhost:8080` does not trigger same-origin policy in all environments. Test from a genuinely different origin.

**Concluding CORS is blocking requests because the browser shows a CORS error:** A CORS error means the browser blocked the *read* — the request was still sent and processed by the server. CORS is not a firewall; it controls what the browser exposes to JavaScript.

**Not testing with credentials:** `fetch(url)` (no credentials) has different CORS rules than `fetch(url, { credentials: 'include' })`. Test both, since the dangerous case always involves credentials.

## Examples

**Quick audit of a production API before launch:**
```bash
BASE=https://api.example.com

for path in /v1/user /v1/account /v1/payment/methods; do
  echo "=== $path ==="
  curl -s -I -H "Origin: https://evil.com" -H "Authorization: Bearer $TOKEN" "$BASE$path" \
    | grep -i "access-control\|vary"
  echo ""
done
```

**Automated CI check (fail if arbitrary origin reflected):**
```bash
#!/usr/bin/env bash
RESPONSE=$(curl -s -I -H "Origin: https://evil-corp.com" https://api.example.com/v1/user)
if echo "$RESPONSE" | grep -qi "access-control-allow-origin: https://evil-corp.com"; then
  echo "FAIL: CORS reflects arbitrary origin"
  exit 1
fi
echo "PASS"
```
