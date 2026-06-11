---
name: apply-content-security-policy
description: Use when deploying a web application that serves HTML — especially one with user-generated content, third-party scripts, or inline JavaScript — to restrict what resources the browser may load or execute.
source: 'OWASP Content Security Policy Cheat Sheet (owasp.org/www-project-cheat-sheets); W3C CSP Level 3 spec; OWASP Top 10 2021 A05; CWE-693'
tags: [security, owasp, csp, content-security-policy, xss-mitigation, web, developer]
---

# Apply Content Security Policy

Reduce the impact of XSS and injection attacks by declaring a strict Content Security Policy that whitelists trusted script, style, and resource origins.

## Why This Is Best Practice

**Adopted by:** Mandated by PCI DSS v4.0 Requirement 6.2.4 and referenced in OWASP Top 10 2021 (A05:Security Misconfiguration). Google, GitHub, Stripe, and Twitter deploy strict CSP headers. Google's CSP Evaluator is used internally and publicly to audit policies. W3C Content Security Policy Level 3 is implemented across all major browsers (Chrome, Firefox, Safari, Edge).
**Impact:** CSP reduces the exploitability of XSS vulnerabilities — Google's security team (2016 study, "CSP Is Dead, Long Live CSP!") found that 94.72% of deployed CSP policies are bypassable, but a strict nonce-based or hash-based CSP eliminates inline script injection entirely. GitHub's adoption of strict-dynamic CSP reduced XSS-related bug bounty payouts significantly. Without CSP, any injected script executes with full page privileges.
**Why best:** `unsafe-inline` and `unsafe-eval` policies look like CSP but provide nearly zero protection — they pass scanners while blocking nothing meaningful. Nonce-based or hash-based strict CSP (using `strict-dynamic`) is the modern approach that actually blocks injected scripts while supporting legitimate inline scripts via one-time tokens.

Sources: OWASP CSP Cheat Sheet; W3C CSP Level 3; Google Security Blog "CSP Is Dead, Long Live CSP!" (2016); CWE-693

## Steps

1. **Start with a report-only policy to avoid breaking production** — set `Content-Security-Policy-Report-Only` first to observe violations without enforcement.

   ```http
   Content-Security-Policy-Report-Only: default-src 'self'; report-uri /csp-report
   ```

   Collect violation reports for 1–2 weeks before switching to enforcing mode.

2. **Build a strict nonce-based policy** — generate a cryptographically random nonce per request, attach it to all `<script>` and `<style>` tags, and reference it in the header.

   ```python
   import secrets
   nonce = secrets.token_urlsafe(16)  # 128-bit random value
   # Attach to template context
   ```

   ```html
   <script nonce="{{ nonce }}">
     // legitimate inline script
   </script>
   ```

   ```http
   Content-Security-Policy:
     default-src 'self';
     script-src 'nonce-{nonce}' 'strict-dynamic';
     style-src 'nonce-{nonce}' 'self';
     object-src 'none';
     base-uri 'self';
     require-trusted-types-for 'script'
   ```

3. **Use `strict-dynamic` to allow dynamically loaded scripts** — `strict-dynamic` propagates trust from a nonced script to scripts it dynamically loads, eliminating the need to whitelist CDN URLs.

   ```http
   script-src 'nonce-{nonce}' 'strict-dynamic'
   ```

   Do NOT use `script-src 'unsafe-inline' https://cdn.example.com` — this is bypassable via any CDN-hosted JSONP or redirect endpoint.

4. **Lock down `object-src`, `base-uri`, and `form-action`** — these are frequently forgotten and allow bypass:

   ```http
   object-src 'none';      # blocks Flash and plugin injection
   base-uri 'self';        # blocks <base href> hijacking
   form-action 'self';     # blocks form exfiltration to attacker domains
   ```

5. **For hash-based CSP (alternative to nonces)** — compute SHA-256 of each inline script's exact content and include the hash:

   ```http
   script-src 'sha256-{hash_of_inline_script}' 'strict-dynamic'
   ```

   Useful for static pages where scripts don't change. Fragile if any whitespace changes.

6. **Set up violation reporting** — use the `report-to` directive (CSP Level 3) to collect real violations in production:

   ```http
   Content-Security-Policy: ...; report-to csp-endpoint
   Report-To: {"group":"csp-endpoint","max_age":10886400,"endpoints":[{"url":"/csp-report"}]}
   ```

7. **Test with Google's CSP Evaluator** — paste your policy at `csp-evaluator.withgoogle.com` to detect bypasses before deployment.

## Rules

- Never use `unsafe-inline` or `unsafe-eval` in a production CSP — they negate most protection.
- Nonces must be regenerated per request — a static nonce is equivalent to no nonce.
- `strict-dynamic` requires a nonce or hash to anchor trust — it's not standalone.
- CSP is defense-in-depth against XSS, not a replacement for output encoding (see `prevent-xss`).

## Common Mistakes

- **Allowlisting entire CDN domains** (`script-src https://cdn.jsdelivr.net`) — any JSONP endpoint or open redirect on that CDN bypasses the policy. Use nonces or hashes instead.
- **Forgetting `object-src 'none'`** — `default-src 'self'` does NOT cover `object-src` in all browsers.
- **Deploying CSP Level 1 allowlist policies** — they are nearly universally bypassable per Google's research. Nonce-based is the minimum effective policy.
