---
name: prevent-open-redirect
description: Use when building any endpoint that redirects users based on a URL parameter — login redirects, logout returns, OAuth callback handling, or any redirect-after-action flow.
source: 'OWASP Unvalidated Redirects and Forwards Cheat Sheet (owasp.org/www-project-cheat-sheets); OWASP Top 10 2021 A01; CWE-601'
tags: [security, owasp, open-redirect, unvalidated-redirect, phishing, web, developer]
---

# Prevent Open Redirect

Validate redirect destination URLs against an allowlist of permitted domains or paths — never blindly redirecting to a URL parameter — to prevent phishing and OAuth token theft via open redirectors.

## Why This Is Best Practice

**Adopted by:** CWE-601 (URL Redirection to Untrusted Site) is a Common Weakness Enumeration entry flagged by every major SAST tool (Checkmarx, Veracode, Semgrep). OWASP Top 10 2021 A01 (Broken Access Control) includes open redirects. Google, Twitter, and PayPal have all had reported open redirect vulnerabilities. Google Bug Bounty pays for open redirectors used in OAuth flows. HackerOne reports open redirects are consistently in the top vulnerability classes submitted.
**Impact:** Open redirectors in OAuth callback endpoints are a primary exploitation technique — an attacker sends a crafted authorization URL where `redirect_uri` points to `https://legitimate-app.com/redirect?next=https://attacker.com`, and the authorization server (if redirect_uri validation is weak) sends the authorization code to the attacker's site via the chain of redirects. Non-OAuth open redirectors enable highly credible phishing: `https://your-bank.com/logout?next=https://attacker-bank.com` appears to originate from the bank's domain.
**Why best:** Sanitizing the URL by stripping protocols or checking for specific patterns fails against IDN homograph attacks, double encoding (`%2f%2f`), and browser-specific quirks. An allowlist of known-safe domains is the only reliable defense.

Sources: OWASP Unvalidated Redirects Cheat Sheet; CWE-601; Google Bug Bounty program; OWASP Top 10 2021

## Steps

1. **Allowlist permitted redirect destinations** — for same-app redirects, restrict to relative URLs or your own domain:

   ```python
   from urllib.parse import urlparse, urljoin

   ALLOWED_HOSTS = {'app.example.com', 'www.example.com'}

   def safe_redirect_url(next_url, fallback='/'):
       if not next_url:
           return fallback
       parsed = urlparse(next_url)
       # Allow relative URLs (no scheme or netloc)
       if not parsed.scheme and not parsed.netloc:
           return urljoin('/', next_url)  # normalize relative path
       # Allow absolute URLs only if host is in allowlist
       if parsed.scheme in ('http', 'https') and parsed.netloc in ALLOWED_HOSTS:
           return next_url
       # Reject all others — unknown host or protocol
       return fallback

   @app.route('/login', methods=['POST'])
   def login():
       # ... authenticate ...
       next_url = request.args.get('next', '/')
       return redirect(safe_redirect_url(next_url))
   ```

2. **For same-site redirects: restrict to relative paths only**:

   ```python
   import re

   def is_safe_relative_url(url):
       # Must not start with // (protocol-relative) or contain scheme
       if re.match(r'^(https?://|//)', url, re.IGNORECASE):
           return False
       # Must start with /
       if not url.startswith('/'):
           return False
       return True

   # Simpler but effective:
   def safe_local_redirect(next_url, fallback='/dashboard'):
       if next_url and is_safe_relative_url(next_url):
           return next_url
       return fallback
   ```

3. **Use indirect references instead of URL parameters** — map tokens to destinations server-side:

   ```python
   import uuid

   # On initial request: store destination, return a token
   def create_redirect_token(destination_url):
       if not is_internal_url(destination_url):
           raise ValueError("Invalid redirect destination")
       token = str(uuid.uuid4())
       redis.setex(f'redirect:{token}', 300, destination_url)  # 5 min TTL
       return token

   # On redirect: look up the token, don't accept raw URLs
   @app.route('/redirect/<token>')
   def perform_redirect(token):
       destination = redis.get(f'redirect:{token}')
       if not destination:
           return redirect('/dashboard')
       redis.delete(f'redirect:{token}')
       return redirect(destination)
   ```

4. **For OAuth redirect_uri: exact match only** (see also `apply-oauth2-security`):

   ```python
   REGISTERED_REDIRECT_URIS = {
       'https://app.example.com/oauth/callback',
       'https://staging.example.com/oauth/callback',
   }

   def validate_redirect_uri(requested_uri):
       if requested_uri not in REGISTERED_REDIRECT_URIS:
           raise OAuthError("invalid_redirect_uri")
   ```

5. **Warn users when redirecting to a different domain** — for intentional cross-domain redirects (e.g., external links):

   ```html
   <!-- Instead of immediate redirect, show intermediate page -->
   <p>You are leaving example.com and going to:</p>
   <p><strong>{{ destination_domain }}</strong></p>
   <a href="{{ safe_external_url }}">Continue</a>
   <a href="/dashboard">Go back</a>
   ```

## Rules

- `//attacker.com/path` is a protocol-relative URL treated as `https://attacker.com/path` — always check for `//` prefix in addition to `http://`.
- IDN homograph attacks (`аpple.com` with Cyrillic 'а') bypass domain string matching — normalize to punycode before comparing: `app.example.com` == `app.xn--example-abc.com`.
- `javascript:` URLs can appear in redirect parameters — block any non-http/https scheme.
- Logout redirects are a primary target — attackers send users to attacker-controlled sites after logout to re-capture credentials.

## Common Mistakes

- **Checking `startswith('https://app.example.com')` for the full URL** — `https://app.example.com.attacker.com/` passes this check.
- **Allowing redirects to `localhost` or `127.0.0.1`** — enables SSRF-via-redirect if the redirect is followed server-side.
- **Trusting the `Referer` header to validate redirects** — it can be spoofed and is often absent.
- **Forgetting about `forward` operations** — server-side internal forwards can also be redirected to unauthorized internal paths.
