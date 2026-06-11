---
name: apply-api-client-security
description: Use when your application consumes external or third-party APIs — configuring HTTP clients, validating TLS certificates, handling secrets, and sanitizing responses from upstream services.
source: 'OWASP API Security Top 10 2023 API10 (Unsafe Consumption of APIs); OWASP Transport Layer Security Cheat Sheet; CWE-295; CWE-918'
tags: [security, owasp, api-client, http-client, certificate-validation, upstream-trust, developer]
---

# Apply API Client Security

Configure outbound HTTP clients with strict TLS verification, timeouts, response size limits, and response sanitization — treating data from third-party APIs as untrusted until validated.

## Why This Is Best Practice

**Adopted by:** OWASP API Security Top 10 2023 API10 (Unsafe Consumption of APIs) is a dedicated category for client-side API security failures. AWS SDK, Google Cloud Client Libraries, and Azure SDK all enforce TLS certificate validation by default. NIST SP 800-52 Rev 2 requires certificate validation for all TLS connections. PCI DSS v4.0 Requirement 6.3.3 mandates secure configurations for components that connect to payment APIs.
**Impact:** Disabling TLS certificate verification (`verify=False` in Python requests, `rejectUnauthorized: false` in Node.js) enables man-in-the-middle attacks on all API calls — attackers on the same network can intercept credentials and response data. OWASP documents cases where third-party APIs injected malicious data into responses that was then rendered to users without sanitization, enabling XSS via upstream compromise. Unlimited response sizes allow external APIs to send gigabyte responses that exhaust server memory.
**Why best:** Trusting upstream APIs completely is the common approach — it's fast to implement but creates a dependency chain attack surface. Treating third-party API responses as potentially hostile (validating schema, sanitizing content, enforcing limits) provides defense-in-depth against upstream compromise.

Sources: OWASP API Security Top 10 2023 API10; CWE-295; CWE-918; NIST SP 800-52 Rev 2

## Steps

1. **Never disable TLS certificate verification**:

   ```python
   import requests

   # BAD — disables all TLS security
   response = requests.get('https://api.example.com', verify=False)

   # GOOD — verify enabled by default, specify CA bundle if needed
   response = requests.get('https://api.example.com')  # verify=True default
   response = requests.get('https://api.example.com', verify='/path/to/custom-ca.pem')
   ```

   ```javascript
   // BAD — Node.js
   const https = require('https');
   const agent = new https.Agent({ rejectUnauthorized: false });
   fetch('https://api.example.com', { agent });

   // GOOD
   fetch('https://api.example.com');  // no custom agent needed
   ```

2. **Set connection and read timeouts** — never wait indefinitely on external APIs:

   ```python
   import requests

   # Set both connect timeout and read timeout (seconds)
   response = requests.get(
       'https://api.example.com/data',
       timeout=(5, 30)  # (connect_timeout, read_timeout)
   )

   # For long-running requests: stream and enforce size limits
   response = requests.get(url, stream=True, timeout=(5, 30))
   content = b''
   MAX_BYTES = 10 * 1024 * 1024  # 10 MB limit
   for chunk in response.iter_content(8192):
       content += chunk
       if len(content) > MAX_BYTES:
           response.close()
           raise ValueError("Response too large")
   ```

3. **Validate response schema before using data**:

   ```python
   from pydantic import BaseModel, ValidationError
   from typing import List, Optional

   class ExternalUserResponse(BaseModel):
       id: str
       email: str
       name: Optional[str] = None
       # NOT accepting arbitrary fields

   def fetch_user(user_id):
       response = requests.get(f'{EXTERNAL_API}/users/{user_id}', timeout=(5, 10))
       response.raise_for_status()
       try:
           return ExternalUserResponse(**response.json())
       except (ValidationError, ValueError) as e:
           logger.error("Invalid response from upstream API: %s", e)
           raise ServiceError("Upstream API returned unexpected data")
   ```

4. **Sanitize response content before rendering** — treat upstream API data as untrusted user input:

   ```python
   import html
   from markupsafe import Markup

   def render_product_description(description_from_api):
       # Escape HTML — the upstream API may be compromised
       safe = html.escape(description_from_api)
       return Markup(safe)  # mark as safe only after escaping
   ```

   If the API returns HTML for rich text: use DOMPurify (browser) or bleach (Python) to sanitize.

5. **Store API credentials securely, rotate them** — never hardcode:

   ```python
   import os
   from functools import lru_cache

   @lru_cache(maxsize=1)
   def get_api_key():
       return os.environ['EXTERNAL_API_KEY']  # from secret manager or env

   # Use per-request headers, don't log them
   headers = {
       'Authorization': f'Bearer {get_api_key()}',
       'User-Agent': 'MyApp/1.0',
   }
   ```

6. **Validate redirect chains** — HTTP clients following redirects may end up at attacker-controlled endpoints:

   ```python
   class SafeHTTPAdapter(requests.adapters.HTTPAdapter):
       ALLOWED_HOSTS = {'api.example.com', 'cdn.example.com'}

       def send(self, request, **kwargs):
           from urllib.parse import urlparse
           host = urlparse(request.url).hostname
           if host not in self.ALLOWED_HOSTS:
               raise SecurityError(f"Redirect to disallowed host: {host}")
           return super().send(request, **kwargs)

   session = requests.Session()
   session.mount('https://', SafeHTTPAdapter())
   ```

7. **Log upstream API errors with context, not credentials**:

   ```python
   def call_external_api(endpoint, params):
       try:
           response = requests.get(endpoint, params=params, timeout=(5, 30))
           response.raise_for_status()
           return response.json()
       except requests.Timeout:
           logger.error("API timeout: %s", endpoint)
           raise ServiceUnavailable()
       except requests.HTTPError as e:
           logger.error("API HTTP error: %s %s", e.response.status_code, endpoint)
           raise
       # Never log: headers (may contain API keys), full request body (may contain credentials)
   ```

## Rules

- `verify=False` in any language is never acceptable in production — use custom CA bundles if your infrastructure uses internal CAs.
- Always validate HTTP status codes before parsing response bodies — a 200 does not mean the body matches the expected schema.
- Retry logic must include jitter and exponential backoff — thundering herds after upstream recovery cause secondary outages.
- Third-party API responses that contain user-controlled data (e.g., social login display names) must be treated as untrusted user input.

## Common Mistakes

- **Propagating upstream 500 errors directly to clients** — leaks implementation details; translate to generic errors (see `design-error-handling`).
- **Using `requests.get(url)` without timeout in a web server** — one slow upstream API call blocks the request worker indefinitely.
- **Logging full API responses** — may contain PII or secrets from the upstream service.
- **Assuming HTTPS means the server is legitimate** — HTTPS only proves the domain is valid; the domain itself may be compromised. Validate the response content too.
