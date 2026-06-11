---
name: apply-oauth2-security
description: Use when implementing OAuth 2.0 or OpenID Connect flows — building an authorization server, integrating a third-party identity provider, or securing API access delegation between services.
source: 'OWASP OAuth 2.0 Security Cheat Sheet (owasp.org/www-project-cheat-sheets); RFC 9700 OAuth 2.0 Security Best Current Practice; RFC 7636 PKCE; CWE-287'
tags: [security, owasp, oauth2, oidc, pkce, authorization, developer, api]
---

# Apply OAuth 2.0 Security

Implement OAuth 2.0 with PKCE for public clients, exact redirect URI matching, short-lived tokens, and state parameter CSRF protection — preventing authorization code interception, token theft, and open redirector attacks.

## Why This Is Best Practice

**Adopted by:** RFC 9700 "OAuth 2.0 Security Best Current Practice" (IETF, 2024) supersedes earlier guidance and mandates PKCE for all clients. Google, Microsoft, GitHub, Stripe, and Okta all implement PKCE for public clients. OAuth 2.0 is the universal standard for API authorization delegation — used by every major platform. NIST SP 800-63C references OAuth/OIDC for federated identity. The implicit flow and password grant are deprecated in RFC 9700.
**Impact:** The authorization code interception attack (PKCE was invented to prevent it) allows any app on a mobile device to steal authorization codes from a legitimate app's redirect URI, then exchange them for tokens. Slack, Facebook, and Uber have all had OAuth-related security vulnerabilities. Missing `state` parameter enables CSRF login attacks where an attacker initiates an OAuth flow and tricks a victim into completing it, binding the victim's account to the attacker's identity.
**Why best:** The implicit flow (returning tokens in URL fragments) was the original browser alternative — it's been deprecated since RFC 9700 because tokens in URLs appear in browser history, referrer headers, and server logs. Authorization code + PKCE provides equivalent functionality with no token exposure in the URL.

Sources: RFC 9700 OAuth 2.0 Security Best Current Practice; RFC 7636 (PKCE); OWASP OAuth Cheat Sheet; CWE-287

## Steps

1. **Use Authorization Code + PKCE for all public clients** (SPAs, mobile apps):

   ```python
   import secrets, hashlib, base64

   # Step 1: Generate PKCE code verifier and challenge
   code_verifier = secrets.token_urlsafe(64)  # 43-128 chars
   code_challenge = base64.urlsafe_b64encode(
       hashlib.sha256(code_verifier.encode()).digest()
   ).rstrip(b'=').decode()

   # Step 2: Authorization request
   auth_url = (
       f"{AUTH_SERVER}/authorize"
       f"?response_type=code"
       f"&client_id={CLIENT_ID}"
       f"&redirect_uri={REDIRECT_URI}"
       f"&scope=openid profile"
       f"&state={secrets.token_urlsafe(16)}"   # CSRF protection
       f"&code_challenge={code_challenge}"
       f"&code_challenge_method=S256"
   )

   # Step 3: Token exchange (include verifier, not challenge)
   token_response = requests.post(f"{AUTH_SERVER}/token", data={
       'grant_type': 'authorization_code',
       'code': authorization_code,
       'redirect_uri': REDIRECT_URI,
       'client_id': CLIENT_ID,
       'code_verifier': code_verifier,  # server verifies hash matches challenge
   })
   ```

2. **Validate `state` parameter to prevent CSRF**:

   ```python
   # Before redirecting to auth server: store state in session
   state = secrets.token_urlsafe(16)
   session['oauth_state'] = state
   session['oauth_code_verifier'] = code_verifier

   # After redirect back: verify state
   if request.args.get('state') != session.pop('oauth_state', None):
       raise SecurityError("OAuth state mismatch — possible CSRF")
   ```

3. **Use exact redirect URI matching** — register full URIs, not patterns or wildcards:

   ```
   # BAD — partial match allows attacker.example.com
   Allowed: https://example.com/*
   Allowed: https://*.example.com/callback

   # GOOD — exact match only
   Allowed: https://app.example.com/oauth/callback
   ```

   On the authorization server: reject any redirect_uri that doesn't exactly match a registered URI.

4. **Keep access tokens short-lived, use refresh token rotation**:

   ```python
   # Access token: 15 minutes
   # Refresh token: 7 days, rotate on each use (single-use)
   ACCESS_TOKEN_TTL = 15 * 60
   REFRESH_TOKEN_TTL = 7 * 24 * 3600

   def rotate_refresh_token(old_refresh_token):
       stored = db.get_refresh_token(old_refresh_token)
       if stored is None or stored.used:
           # Refresh token reuse detected — revoke entire family
           db.revoke_token_family(stored.family_id)
           raise SecurityError("Refresh token reuse detected")
       db.mark_used(old_refresh_token)
       return db.create_refresh_token(family_id=stored.family_id)
   ```

5. **Store tokens securely**:

   | Client type | Access token | Refresh token |
   |---|---|---|
   | SPA | Memory only (not localStorage) | HttpOnly cookie |
   | Mobile | OS Keychain/Keystore | OS Keychain/Keystore |
   | Server-side | In-memory or encrypted DB | Encrypted DB |

   Never store tokens in `localStorage` — XSS reads them. Never store in URL parameters — they appear in logs.

6. **Validate token audience (`aud`) on the resource server**:

   ```python
   # Resource server must verify the token was issued for this API
   payload = jwt.decode(token, PUBLIC_KEY, algorithms=['RS256'],
                        audience='https://api.example.com')
   # Rejects tokens issued for other services
   ```

7. **Use client authentication for confidential clients** — server-to-server flows:

   ```python
   # Client credentials flow (no user, service-to-service)
   # Use client_secret_jwt or private_key_jwt — not client_secret_post
   # client_secret_post sends secret in request body (logged by servers)
   token_response = requests.post(f"{AUTH_SERVER}/token", data={
       'grant_type': 'client_credentials',
       'client_assertion_type': 'urn:ietf:params:oauth:client-assertion-type:jwt-bearer',
       'client_assertion': create_client_jwt(CLIENT_ID, PRIVATE_KEY),
       'scope': 'api:read',
   })
   ```

## Rules

- The implicit flow (`response_type=token`) is deprecated — never implement it for new clients.
- The resource owner password credentials grant is deprecated — use authorization code flow instead.
- PKCE is required for public clients even when using a server-side authorization code exchange.
- `nonce` claim in OIDC is required to prevent ID token replay attacks in authorization code flows.

## Common Mistakes

- **Accepting any redirect URI with a matching domain** — allows open redirector injection where attacker registers `https://app.example.com.attacker.com/callback`.
- **Storing refresh tokens in localStorage** — XSS steals long-lived refresh tokens, enabling indefinite access.
- **Not validating `state`** — enables login CSRF where attacker tricks victim into completing the attacker's OAuth flow.
- **Using `response_type=code` without PKCE for SPAs** — without PKCE, authorization codes intercepted from the redirect can be exchanged by attackers.
