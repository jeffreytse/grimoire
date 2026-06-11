---
name: configure-cors-for-oauth
description: Use when building an OAuth 2.0 authorization server or configuring CORS on OAuth endpoints — token, userinfo, revocation, and introspection endpoints each have different CORS requirements that differ from ordinary API CORS.
source: 'RFC 9700 OAuth 2.0 Security Best Current Practice (2024); draft-ietf-oauth-browser-based-apps (OAuth 2.0 for Browser-Based Apps); RFC 8252 OAuth 2.0 for Native Apps; OWASP OAuth Cheat Sheet'
tags: [security, cors, oauth, pkce, authorization-server, api, browser-based-apps, tokens]
verified: true
---

# Configure CORS for OAuth

Apply CORS selectively per OAuth endpoint — token and userinfo endpoints need tight origin allowlists when called from browser-based apps, while authorize, introspect, and most server-to-server endpoints should have no CORS at all.

## Why This Is Best Practice

**Adopted by:** RFC 9700 (OAuth 2.0 Security Best Current Practice, 2024) explicitly addresses CORS on token endpoints for browser-based apps. The OAuth 2.0 for Browser-Based Apps draft (IETF) specifies which endpoints require CORS and which must not expose it. Auth0, Okta, and Keycloak each configure CORS per-endpoint, not globally.

**Impact:** Applying a single global CORS policy to an authorization server is a common misconfiguration. The token endpoint returns access tokens and refresh tokens — an overly permissive CORS policy on it allows any attacker website to initiate flows and potentially harvest tokens. Conversely, missing CORS on the token endpoint for a PKCE flow causes legitimate browser apps to fail at the last step (exchanging the code for tokens), which developers then "fix" by opening CORS to `*`.

**Why best:** Each OAuth endpoint has a different caller profile. The authorization endpoint (`/authorize`) is never called via XHR — it's a browser redirect. The token endpoint (`/token`) is called server-to-server in confidential client flows, but browser-to-server in PKCE flows. The userinfo endpoint is always called from browser JS. Treating all endpoints identically breaks one or opens the other.

Sources: RFC 9700 §4.11; draft-ietf-oauth-browser-based-apps §8; RFC 8252; OWASP OAuth Cheat Sheet

## Steps

### Step 1: Map each OAuth endpoint to its caller

Before configuring CORS, identify who calls each endpoint:

| Endpoint | Called by | Needs CORS? |
|----------|-----------|-------------|
| `/oauth/authorize` | Browser redirect (not XHR) | No — never via `fetch` |
| `/oauth/token` | Server (confidential clients) | No |
| `/oauth/token` | Browser JS (PKCE / public clients) | Yes — tight allowlist |
| `/oauth/userinfo` | Browser JS after token grant | Yes — tight allowlist |
| `/oauth/revoke` | Server or browser | Yes if browser clients use it |
| `/oauth/introspect` | Server-to-server only | No — block browser access |
| `/.well-known/openid-configuration` | Browser JS (OIDC discovery) | Yes — `*` is safe here |
| `/oauth/jwks` | Browser JS (token verification) | Yes — `*` is safe here |

### Step 2: Configure the token endpoint for PKCE flows

The token endpoint must only accept origins that match your registered clients. Never use `*` — the response contains access and refresh tokens.

```python
# Python / Flask example
REGISTERED_CLIENT_ORIGINS = {
    'https://app.example.com',
    'https://mobile.example.com',
}

@app.route('/oauth/token', methods=['POST', 'OPTIONS'])
def token_endpoint():
    origin = request.headers.get('Origin')

    if request.method == 'OPTIONS':
        # Only respond to preflights from registered origins
        if origin not in REGISTERED_CLIENT_ORIGINS:
            return '', 403
        response = make_response('', 204)
        response.headers['Access-Control-Allow-Origin'] = origin
        response.headers['Access-Control-Allow-Methods'] = 'POST'
        response.headers['Access-Control-Allow-Headers'] = 'Content-Type'
        response.headers['Access-Control-Max-Age'] = '86400'
        response.headers['Vary'] = 'Origin'
        return response

    # Process token request...
    response = make_response(jsonify(issue_tokens(request)))
    if origin in REGISTERED_CLIENT_ORIGINS:
        response.headers['Access-Control-Allow-Origin'] = origin
        response.headers['Vary'] = 'Origin'
    return response
```

```javascript
// Node.js / Express example
const REGISTERED_CLIENT_ORIGINS = new Set([
  'https://app.example.com',
  'https://mobile.example.com',
]);

app.use('/oauth/token', (req, res, next) => {
  const origin = req.headers.origin;
  if (origin && REGISTERED_CLIENT_ORIGINS.has(origin)) {
    res.setHeader('Access-Control-Allow-Origin', origin);
    res.setHeader('Vary', 'Origin');
    if (req.method === 'OPTIONS') {
      res.setHeader('Access-Control-Allow-Methods', 'POST');
      res.setHeader('Access-Control-Allow-Headers', 'Content-Type');
      res.setHeader('Access-Control-Max-Age', '86400');
      return res.sendStatus(204);
    }
  } else if (req.method === 'OPTIONS') {
    return res.sendStatus(403);
  }
  next();
});
```

### Step 3: Configure the userinfo endpoint

The userinfo endpoint is always called from browser JS (it's the point of OIDC — the client fetches user claims after receiving tokens). Apply the same registered-origin allowlist.

```python
@app.route('/oauth/userinfo', methods=['GET', 'OPTIONS'])
def userinfo_endpoint():
    origin = request.headers.get('Origin')
    allowed = origin if origin in REGISTERED_CLIENT_ORIGINS else None

    if request.method == 'OPTIONS':
        if not allowed:
            return '', 403
        response = make_response('', 204)
        response.headers['Access-Control-Allow-Origin'] = allowed
        response.headers['Access-Control-Allow-Methods'] = 'GET'
        response.headers['Access-Control-Allow-Headers'] = 'Authorization'
        response.headers['Vary'] = 'Origin'
        return response

    # Validate Bearer token, return claims...
    response = make_response(jsonify(get_claims(request)))
    if allowed:
        response.headers['Access-Control-Allow-Origin'] = allowed
        response.headers['Vary'] = 'Origin'
    return response
```

### Step 4: Open discovery and JWKS endpoints to `*`

These endpoints contain only public metadata — no secrets, no tokens. `*` is correct and required here: any client needs to discover your authorization server's configuration without pre-registration.

```python
@app.route('/.well-known/openid-configuration')
def openid_configuration():
    response = make_response(jsonify(build_discovery_document()))
    response.headers['Access-Control-Allow-Origin'] = '*'
    return response

@app.route('/.well-known/jwks.json')
def jwks():
    response = make_response(jsonify(get_public_keys()))
    response.headers['Access-Control-Allow-Origin'] = '*'
    return response
```

### Step 5: Block CORS on introspection and server-to-server endpoints

The introspection endpoint (`/oauth/introspect`) verifies tokens — it's a resource server-to-authorization server call. No browser client should ever call it directly. Return no CORS headers; if a browser sends an OPTIONS preflight, return 403.

```python
@app.route('/oauth/introspect', methods=['POST', 'OPTIONS'])
def introspect_endpoint():
    if request.headers.get('Origin'):
        # Any cross-origin request (browser) is rejected
        return '', 403
    # Process introspection for server clients...
```

### Step 6: Sync CORS allowlist with client registration

The origin allowlist for the token endpoint must match the `redirect_uri` origins in your client registry. When a new client is registered with `redirect_uri: https://newapp.example.com/callback`, add `https://newapp.example.com` to the CORS allowlist at the same time.

```python
# On client registration:
def register_client(client_data):
    redirect_origins = {
        urlparse(uri).scheme + '://' + urlparse(uri).netloc
        for uri in client_data['redirect_uris']
    }
    client_data['allowed_cors_origins'] = list(redirect_origins)
    save_client(client_data)
    # Update in-memory CORS allowlist
    REGISTERED_CLIENT_ORIGINS.update(redirect_origins)
```

## Rules

- Never use `Access-Control-Allow-Origin: *` on the token endpoint — the response contains credentials (access tokens, refresh tokens, id tokens).
- Do not apply a global CORS middleware to the entire authorization server — configure per-endpoint.
- The `/oauth/authorize` endpoint is never called via `fetch` or `XMLHttpRequest` in a correct OAuth flow — it is a browser redirect. If a client is calling it via XHR, the client is misconfigured, not the server.
- `Vary: Origin` is required on all endpoints that reflect the origin dynamically, including the token endpoint.
- CORS on the token endpoint does not replace PKCE — both are required for browser-based apps. CORS controls what the browser exposes to JavaScript; PKCE prevents authorization code interception.

## Common Mistakes

**Global CORS middleware applied to all routes:** `app.use(cors({ origin: allowedOrigins }))` before route definitions applies the same policy to `/oauth/introspect` as to `/oauth/userinfo`. Always configure CORS per-route on an authorization server.

**Adding `localhost` origins to the production allowlist:** During development, adding `http://localhost:3000` to the production CORS allowlist is a common shortcut. Any attacker running a local server can then call your production token endpoint. Use environment-specific allowlists.

**Missing CORS on the token endpoint, then "fixing" it with `*`:** When a PKCE flow fails because the token endpoint has no CORS headers, the path of least resistance is `Access-Control-Allow-Origin: *`. This exposes tokens to any origin. Use a registered-client allowlist instead.

**Assuming confidential client flows need CORS:** Client credentials grant, authorization code grant with a backend server exchanging the code, and all server-to-server flows do not involve browser XHR. Adding CORS to these flows increases attack surface without benefit.

## Examples

**PKCE flow: what needs CORS at each step**
```
1. Redirect to /oauth/authorize?...       ← browser redirect, no XHR, no CORS needed
2. User authenticates on AS              ← same-origin to AS, no CORS needed
3. AS redirects back with ?code=...      ← browser redirect, no CORS needed
4. Client JS POST /oauth/token {code}    ← XHR, CORS required on token endpoint
5. Client JS GET /oauth/userinfo         ← XHR, CORS required on userinfo endpoint
```

**Keycloak: configure per-client origins (not global)**
```
Realm Settings → Clients → [your client] → Web Origins
Add: https://app.example.com
(Keycloak uses this list for token endpoint CORS, separate from redirect_uri list)
```

**Auth0: configure allowed origins per application**
```
Dashboard → Applications → [your app] → Settings → Allowed Web Origins
Add: https://app.example.com
(distinct from Allowed Callback URLs — controls token endpoint CORS)
```
