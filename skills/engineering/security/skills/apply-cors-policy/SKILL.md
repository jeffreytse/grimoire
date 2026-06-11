---
name: apply-cors-policy
description: Use when building an API or web service that needs to be called from browser JavaScript on a different origin — including public APIs, microservices accessed from SPAs, and third-party integrations.
source: 'OWASP Cross-Origin Resource Sharing Cheat Sheet (owasp.org/www-project-cheat-sheets); W3C CORS specification; OWASP Top 10 2021 A01; CWE-942'
tags: [security, owasp, cors, cross-origin, api, web, developer]
---

# Apply CORS Policy

Configure CORS headers with an explicit allowlist of trusted origins — never reflecting arbitrary origins or using `Access-Control-Allow-Origin: *` with credentials — to prevent cross-origin data theft.

## Why This Is Best Practice

**Adopted by:** W3C CORS specification is implemented in all modern browsers. OWASP Top 10 2021 A01 (Broken Access Control) includes CORS misconfiguration. AWS API Gateway, Azure APIM, and Cloudflare Workers all provide explicit CORS configuration with origin allowlisting. Major APIs (Stripe, Twilio, SendGrid) require explicit CORS origin configuration.
**Impact:** CWE-942 (Overly Permissive Cross-origin Policy) is one of the most common API security misconfigurations. A misconfigured CORS that reflects any origin with credentials enables cross-origin attacks: any website can make authenticated requests to your API and read the response. Portswigger Web Security Academy documents CORS misconfigurations leading to account takeover in real-world bug bounty reports. Netflix, Coinbase, and numerous fintechs have paid bug bounties for CORS misconfigurations exposing authenticated API data.
**Why best:** The browser's same-origin policy already blocks cross-origin reads — CORS is the controlled exception. `Access-Control-Allow-Origin: *` is safe for public, unauthenticated APIs (public CDN assets, public data APIs) but catastrophic for authenticated APIs. The distinction: `*` cannot be combined with `Access-Control-Allow-Credentials: true` in the spec, but reflection attacks bypass this by echoing back the requesting origin instead of using `*`.

Sources: OWASP CORS Cheat Sheet; W3C CORS spec (fetch.spec.whatwg.org); CWE-942; Portswigger CORS research

## Steps

1. **Define an explicit allowlist of permitted origins**:

   ```python
   ALLOWED_ORIGINS = {
       'https://app.example.com',
       'https://admin.example.com',
       'https://partner.example.com',
   }

   def get_cors_origin(request_origin):
       if request_origin in ALLOWED_ORIGINS:
           return request_origin
       return None  # do not set header for disallowed origins
   ```

2. **Set CORS headers on every response, not just the preflight**:

   ```python
   from flask import request, make_response

   @app.after_request
   def add_cors_headers(response):
       origin = request.headers.get('Origin')
       allowed = get_cors_origin(origin)
       if allowed:
           response.headers['Access-Control-Allow-Origin'] = allowed
           response.headers['Access-Control-Allow-Credentials'] = 'true'
           response.headers['Vary'] = 'Origin'  # required when echoing origin
       return response

   @app.route('/api/resource', methods=['OPTIONS'])
   def preflight():
       response = make_response()
       origin = request.headers.get('Origin')
       allowed = get_cors_origin(origin)
       if allowed:
           response.headers['Access-Control-Allow-Origin'] = allowed
           response.headers['Access-Control-Allow-Methods'] = 'GET, POST, PUT, DELETE'
           response.headers['Access-Control-Allow-Headers'] = 'Content-Type, Authorization'
           response.headers['Access-Control-Max-Age'] = '86400'
           response.headers['Access-Control-Allow-Credentials'] = 'true'
       return response, 204
   ```

3. **Use `Access-Control-Allow-Origin: *` only for public, unauthenticated endpoints**:

   ```http
   # Safe: public CDN assets, open data API (no cookies/auth)
   Access-Control-Allow-Origin: *

   # Never combine * with credentials:
   Access-Control-Allow-Origin: *
   Access-Control-Allow-Credentials: true   ← browsers block this per spec, but see reflection attack
   ```

4. **Always set `Vary: Origin`** when reflecting origins dynamically — without this, CDNs and proxies cache one origin's CORS response and serve it to others:

   ```http
   Access-Control-Allow-Origin: https://app.example.com
   Vary: Origin
   ```

5. **Restrict allowed headers and methods to what your API actually uses**:

   ```http
   Access-Control-Allow-Methods: GET, POST
   Access-Control-Allow-Headers: Content-Type, Authorization, X-CSRF-Token
   ```

   Do not use `Access-Control-Allow-Headers: *` for credentialed requests (unsupported in Safari and creates unnecessary exposure).

6. **For development environments** — use environment variables for allowed origins so local development doesn't bleed into production:

   ```python
   # development: allow localhost
   # production: allow only production domains
   ALLOWED_ORIGINS = set(os.environ.get('CORS_ALLOWED_ORIGINS', '').split(','))
   ```

7. **Use framework middleware rather than manual headers**:

   ```python
   # Flask-CORS
   from flask_cors import CORS
   CORS(app, origins=['https://app.example.com'], supports_credentials=True)
   ```

   ```javascript
   // Express cors package
   const cors = require('cors');
   app.use(cors({
     origin: ['https://app.example.com'],
     credentials: true,
   }));
   ```

   ```typescript
   // Next.js App Router — app/api/route.ts
   const ALLOWED_ORIGINS = ['https://app.example.com'];

   export async function OPTIONS(request: Request) {
     const origin = request.headers.get('origin') ?? '';
     const allowed = ALLOWED_ORIGINS.includes(origin) ? origin : '';
     return new Response(null, {
       status: 204,
       headers: {
         'Access-Control-Allow-Origin': allowed,
         'Access-Control-Allow-Methods': 'GET, POST',
         'Access-Control-Allow-Headers': 'Content-Type, Authorization',
         'Access-Control-Max-Age': '86400',
         'Vary': 'Origin',
       },
     });
   }
   ```

   ```go
   // Go net/http — manual middleware
   var allowedOrigins = map[string]bool{
       "https://app.example.com": true,
   }

   func corsMiddleware(next http.Handler) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           origin := r.Header.Get("Origin")
           if allowedOrigins[origin] {
               w.Header().Set("Access-Control-Allow-Origin", origin)
               w.Header().Set("Vary", "Origin")
               if r.Method == http.MethodOptions {
                   w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
                   w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
                   w.Header().Set("Access-Control-Max-Age", "86400")
                   w.WriteHeader(http.StatusNoContent)
                   return
               }
           }
           next.ServeHTTP(w, r)
       })
   }
   ```

   ```nginx
   # Nginx — origin allowlist via map
   map $http_origin $cors_origin {
       default                      "";
       "https://app.example.com"    $http_origin;
       "https://admin.example.com"  $http_origin;
   }

   server {
       location /api/ {
           if ($request_method = OPTIONS) {
               add_header Access-Control-Allow-Origin  $cors_origin always;
               add_header Access-Control-Allow-Methods "GET, POST" always;
               add_header Access-Control-Allow-Headers "Content-Type, Authorization" always;
               add_header Access-Control-Max-Age       86400 always;
               add_header Vary                         Origin always;
               return 204;
           }
           add_header Access-Control-Allow-Origin $cors_origin always;
           add_header Vary Origin always;
           proxy_pass http://backend;
       }
   }
   ```

   ```yaml
   # AWS API Gateway — REST API (serverless.yml / SAM)
   Cors:
     AllowOrigins:
       - https://app.example.com
     AllowHeaders:
       - Content-Type
       - Authorization
     AllowMethods:
       - GET
       - POST
     AllowCredentials: true
   ```

## Rules

- Never reflect `Origin` without checking against an allowlist — `if (origin) response.set('Access-Control-Allow-Origin', origin)` is the CORS misconfiguration pattern.
- `Access-Control-Allow-Origin: null` is dangerous — sandbox iframes and local files send `Origin: null`, which `null` allows.
- CORS is not authentication — it controls browser-to-server data visibility, not server-to-server access. API keys and auth tokens are still required.
- CORS does not prevent CSRF — a credentialed cross-origin request is blocked from reading the response but the request itself is still sent and processed.

## Common Mistakes

- **Reflecting `Origin` header without validation** — `origin = request.headers['Origin']; response['ACAO'] = origin` allows any attacker site to call the API with credentials.
- **Testing CORS with Postman and concluding it's secure** — Postman doesn't enforce the same-origin policy; only browser-based tests reveal CORS behavior.
- **Forgetting `Vary: Origin`** — a cached CORS response for `https://app.example.com` gets served to `https://attacker.com` if `Vary` is missing.
- **Using regex with substring matching** — `if 'example.com' in origin` allows `https://evil-example.com` through.
