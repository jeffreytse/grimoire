---
name: prevent-clickjacking
description: Use when building any web page that performs state-changing actions on click — login forms, payment buttons, delete confirmations, or settings toggles that could be exploited if framed by an attacker.
source: 'OWASP Clickjacking Defense Cheat Sheet (owasp.org/www-project-cheat-sheets); OWASP Top 10 2021 A05; CWE-1021; W3C CSP Level 3 frame-ancestors'
tags: [security, owasp, clickjacking, iframe, x-frame-options, csp, web, developer]
---

# Prevent Clickjacking

Block your pages from being embedded in attacker-controlled iframes using `X-Frame-Options` and CSP `frame-ancestors` — eliminating UI redress attacks where invisible overlays trick users into unintended clicks.

## Why This Is Best Practice

**Adopted by:** Recommended by OWASP Top 10 2021 A05 (Security Misconfiguration) and NIST SP 800-53 SC-18. Twitter fell victim to a major clickjacking attack ("likejacking") in 2009 before deploying frame-busting. Facebook, Google, and GitHub all deploy `X-Frame-Options: DENY` or equivalent CSP. Mozilla Observatory and securityheaders.com penalize missing frame protection.
**Impact:** Clickjacking enables attacks ranging from social engineering (forced likes, follows) to high-impact account compromise (changing email/password via one-click on a framed settings page). The 2010 "Filejacking" attack used clickjacking to steal Adobe Flash local storage via framed permissions dialogs. Proper frame protection eliminates the entire attack class at zero performance cost.
**Why best:** JavaScript frame-busting code (checking `window.top !== window.self`) is the legacy alternative — it's defeated by the `sandbox="allow-scripts allow-same-origin"` iframe attribute, introduced in HTML5. Header-based protection (`X-Frame-Options`, `frame-ancestors`) operates at the browser level and cannot be bypassed by iframe sandbox attributes.

Sources: OWASP Clickjacking Defense Cheat Sheet; Marcus Niemietz (Ruhr-Universität Bochum) — "Clickjacking Revisited" (2012); CWE-1021

## Steps

1. **Set `X-Frame-Options` header** — works in all browsers including IE:

   ```http
   X-Frame-Options: DENY
   ```

   Use `SAMEORIGIN` only if your app legitimately embeds its own pages in iframes (e.g., an internal dashboard embedding sub-pages):

   ```http
   X-Frame-Options: SAMEORIGIN
   ```

   Never use `ALLOW-FROM` — it is not supported in Chrome/Firefox and is silently ignored.

2. **Set CSP `frame-ancestors` (preferred, supersedes X-Frame-Options)**:

   ```http
   Content-Security-Policy: frame-ancestors 'none'
   ```

   `frame-ancestors 'none'` is equivalent to `X-Frame-Options: DENY`.
   `frame-ancestors 'self'` is equivalent to `X-Frame-Options: SAMEORIGIN`.

   For specific trusted origins:
   ```http
   Content-Security-Policy: frame-ancestors 'self' https://trusted-parent.example.com
   ```

3. **Set both headers for maximum browser coverage** — `X-Frame-Options` covers older browsers; `frame-ancestors` covers modern ones and takes precedence when both are present:

   ```nginx
   add_header X-Frame-Options "DENY" always;
   add_header Content-Security-Policy "frame-ancestors 'none'" always;
   ```

4. **Apply to all HTML responses** — not just login pages. Clickjacking targets any page with interactive elements. Configure at the server/framework level, not per-route:

   **Express:**
   ```javascript
   const helmet = require('helmet');
   app.use(helmet.frameguard({ action: 'deny' }));
   ```

   **Django:**
   ```python
   # settings.py
   MIDDLEWARE = ['django.middleware.clickjacking.XFrameOptionsMiddleware']
   X_FRAME_OPTIONS = 'DENY'
   ```

   **Spring Security:**
   ```java
   http.headers().frameOptions().deny();
   ```

5. **If you must allow framing (widgets, embeds)** — use `frame-ancestors` with an explicit allowlist and validate on the server that the embedding origin is trusted. Never use `X-Frame-Options: ALLOW-FROM` (unsupported in modern browsers):

   ```http
   Content-Security-Policy: frame-ancestors 'self' https://partner.example.com
   ```

## Rules

- Remove any JavaScript frame-busting code — it's defeated by `sandbox` and creates false confidence.
- `X-Frame-Options: ALLOW-FROM` is not supported in Chrome or Firefox — use CSP `frame-ancestors` with specific origins instead.
- The `frame-ancestors` directive in CSP overrides `X-Frame-Options` in modern browsers — set both for full coverage.
- These headers have zero effect on non-HTML responses (images, JSON, PDFs) — they're only relevant for HTML page framing.

## Common Mistakes

- **Only protecting the login page** — attackers frame password-change or payment confirmation pages more often than login pages.
- **Using JavaScript frame-busting** — `if (top !== self) top.location = self.location` is bypassed by `<iframe sandbox="allow-scripts allow-same-origin">`.
- **Setting `SAMEORIGIN` when `DENY` would work** — if no legitimate use case for framing exists, `DENY` is stricter and correct.
