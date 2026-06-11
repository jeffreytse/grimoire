---
name: apply-security-headers
description: Use when configuring any web server or application that serves HTTP responses to browsers — covers the full set of security response headers that harden against common attacks.
source: 'OWASP HTTP Security Response Headers Cheat Sheet (owasp.org/www-project-cheat-sheets); OWASP Top 10 2021 A05; securityheaders.com scoring criteria; MDN Web Docs Security headers'
tags: [security, owasp, http-headers, hsts, csp, web, developer, devops]
---

# Apply Security Headers

Configure the full set of HTTP security response headers to eliminate browser-level attack vectors — clickjacking, MIME sniffing, information disclosure, and downgrade attacks.

## Why This Is Best Practice

**Adopted by:** OWASP Top 10 2021 A05 (Security Misconfiguration) cites missing security headers as a primary finding. PCI DSS v4.0 Requirement 6.2.4 requires protection against injection and misconfiguration. Mozilla Observatory, securityheaders.com, and Google's PageSpeed all include security header grading. GitHub, Google, Stripe, and all major cloud providers configure these headers by default in their platforms.
**Impact:** Scott Helme's analysis of the Alexa Top 1 Million sites (2023) found 72% lack HSTS, 80% lack CSP, and 60% lack X-Frame-Options — each representing an exploitable gap. HSTS alone eliminates SSL-stripping attacks (responsible for credential theft in public WiFi scenarios). Missing `X-Content-Type-Options` enables MIME-confusion attacks where text files execute as scripts.
**Why best:** Per-header configuration over a WAF or CDN setting — because WAFs can be bypassed and CDN configs can be overridden. Setting headers at the application layer ensures they're always present regardless of routing or proxy configuration.

Sources: OWASP HTTP Security Response Headers Cheat Sheet; Scott Helme Observatory data (2023); MDN Web Security; CWE-16

## Steps

1. **Set `Strict-Transport-Security` (HSTS)** — force all connections to use HTTPS, including subdomains. Prevents SSL stripping.

   ```http
   Strict-Transport-Security: max-age=63072000; includeSubDomains; preload
   ```

   - Start with `max-age=300` to test, then increase to `63072000` (2 years) before submitting to the HSTS preload list.
   - Submit at `hstspreload.org` to be baked into browsers permanently.
   - Only set on HTTPS — HTTP servers must not send HSTS.

2. **Set `X-Content-Type-Options: nosniff`** — prevents browsers from MIME-sniffing responses away from the declared `Content-Type`.

   ```http
   X-Content-Type-Options: nosniff
   ```

   Required when serving user-uploaded files — without it, a `.jpg` containing HTML may execute as `text/html`.

3. **Set `X-Frame-Options: DENY` or `SAMEORIGIN`** — prevents your pages from being embedded in iframes on attacker sites (clickjacking).

   ```http
   X-Frame-Options: DENY
   ```

   Use `SAMEORIGIN` if your app embeds its own pages in iframes. Superseded by CSP `frame-ancestors` but keep both for legacy browser support.

4. **Set `Referrer-Policy`** — control how much referrer information is sent with cross-origin requests.

   ```http
   Referrer-Policy: strict-origin-when-cross-origin
   ```

   Prevents leaking URL paths (which may contain session tokens or sensitive parameters) to third-party sites.

5. **Set `Permissions-Policy`** — restrict browser features (camera, microphone, geolocation) to only what your app needs.

   ```http
   Permissions-Policy: camera=(), microphone=(), geolocation=(), payment=()
   ```

   Prevents malicious scripts or third-party embeds from accessing hardware APIs.

6. **Set `Content-Security-Policy`** — see `apply-content-security-policy` for the full guide. Minimum viable CSP:

   ```http
   Content-Security-Policy: default-src 'self'; object-src 'none'; base-uri 'self'
   ```

7. **Remove information-leaking headers** — strip headers that reveal stack details:

   ```http
   # Remove these (configure in web server):
   Server: Apache/2.4.51    →  Server: (remove or set to generic value)
   X-Powered-By: PHP/8.1    →  remove entirely
   X-AspNet-Version: 4.x    →  remove entirely
   ```

   Nginx: `server_tokens off;`
   Apache: `ServerTokens Prod; ServerSignature Off`
   Express: `app.disable('x-powered-by')`

8. **Verify with securityheaders.com** — paste your URL to get an A–F grade and per-header findings.

## Rules

- HSTS `preload` is irreversible — once submitted, the domain is HTTPS-only in all browsers permanently. Test thoroughly before preloading.
- `X-Frame-Options` is deprecated in favor of CSP `frame-ancestors` but should still be set for IE11 and old Safari.
- `Permissions-Policy` syntax changed (from `Feature-Policy`) — use the newer form shown above.
- Security headers must be set on every response, not just HTML — API responses and file downloads also need `X-Content-Type-Options`.

## Common Mistakes

- **Setting headers only on HTML responses** — JavaScript and JSON responses also need `nosniff`.
- **Using `X-Frame-Options: ALLOWALL`** — this is not a valid value and is silently ignored.
- **Setting HSTS on HTTP responses** — browsers ignore HSTS on non-HTTPS responses; this does nothing.
- **Adding `unsafe-inline` to CSP thinking it helps** — it removes CSP's XSS protection entirely.
