---
name: apply-jwt-security
description: Use when implementing JWT-based authentication or authorization — covering token generation, algorithm selection, validation, expiry, and revocation for stateless APIs.
source: 'OWASP JSON Web Token Cheat Sheet (owasp.org/www-project-cheat-sheets); RFC 7519 (JWT); RFC 8725 (JWT Best Current Practices); CWE-347'
tags: [security, owasp, jwt, authentication, authorization, api, developer]
---

# Apply JWT Security

Issue JWTs with explicit algorithm pinning, short expiry, and server-side validation — protecting against algorithm confusion attacks, token forgery, and indefinite token validity.

## Why This Is Best Practice

**Adopted by:** RFC 8725 "JWT Best Current Practices" (IETF, 2020) is the definitive standard. Auth0, Okta, AWS Cognito, and Google Identity Platform all enforce strict JWT validation. OWASP API Security Top 10 2023 (API2:Broken Authentication) cites JWT misuse. PCI DSS v4.0 Requirement 8.6 requires authentication token expiration controls.
**Impact:** The `alg:none` attack (CVE-2015-9235 in jwt.js) allowed any JWT to be accepted as valid by setting the algorithm to "none" — effectively bypassing all authentication. The RS256/HS256 algorithm confusion attack (exploited at Accenture, SecureWorks, and in CTF competitions) allows attackers to forge tokens using the public key as the HMAC secret. Auth0 paid $1,500–$5,000 bounties for JWT misuse findings. Improper JWT validation is a top finding in API security audits.
**Why best:** Client-side storage of authentication state (JWT) vs. server-side sessions — JWTs eliminate database lookups but require strict validation. The tradeoff: JWTs cannot be revoked without a denylist, so short expiry is mandatory. Opaque session tokens (see `design-session-management`) are better for applications needing instant revocation; JWTs are better for stateless, distributed systems.

Sources: RFC 7519 (JWT); RFC 8725 (JWT Best Current Practices); IETF; CWE-347; Auth0 security blog

## Steps

1. **Pin the algorithm explicitly on verification** — never accept the algorithm from the token header:

   ```python
   import jwt  # PyJWT

   # BAD — accepts whatever alg the token claims
   payload = jwt.decode(token, key, algorithms=None)

   # GOOD — pin to one algorithm
   payload = jwt.decode(token, SECRET_KEY, algorithms=['HS256'])

   # For RS256 (asymmetric):
   payload = jwt.decode(token, PUBLIC_KEY, algorithms=['RS256'])
   ```

   ```javascript
   // Node.js — jsonwebtoken
   // BAD
   jwt.verify(token, key);

   // GOOD
   jwt.verify(token, key, { algorithms: ['RS256'] });
   ```

2. **Use RS256 (asymmetric) for public APIs, HS256 for internal services**:

   - **RS256**: Sign with private key, verify with public key. Public key can be shared without risk.
   - **HS256**: Single shared secret — anyone with the secret can forge tokens. Use only for internal, controlled environments.
   - Never use RS256 public key as HS256 secret — this is the algorithm confusion attack.

3. **Set short `exp` (expiry) claims — minutes for access tokens, hours for refresh tokens**:

   ```python
   from datetime import datetime, timedelta, timezone
   import jwt

   payload = {
       'sub': str(user_id),
       'iat': datetime.now(timezone.utc),
       'exp': datetime.now(timezone.utc) + timedelta(minutes=15),  # access token
       'jti': str(uuid.uuid4()),  # unique token ID for revocation
   }
   token = jwt.encode(payload, PRIVATE_KEY, algorithm='RS256')
   ```

   Typical values: access token 15 min, refresh token 7 days. Never issue tokens without `exp`.

4. **Validate all required claims** — `exp`, `iss` (issuer), `aud` (audience), `nbf` (not before):

   ```python
   payload = jwt.decode(
       token,
       PUBLIC_KEY,
       algorithms=['RS256'],
       options={
           'require': ['exp', 'iss', 'sub'],
           'verify_exp': True,
           'verify_iss': True,
       },
       issuer='https://auth.example.com',
       audience='https://api.example.com',
   )
   ```

5. **Implement token revocation with a denylist for logout and high-security actions**:

   ```python
   # Store revoked JTIs in Redis with TTL = token remaining expiry
   def revoke_token(jti, ttl_seconds):
       redis.setex(f'revoked:{jti}', ttl_seconds, '1')

   def is_revoked(jti):
       return redis.exists(f'revoked:{jti}')

   # On every request:
   payload = jwt.decode(token, ...)
   if is_revoked(payload['jti']):
       raise AuthError("Token revoked")
   ```

6. **Store JWTs securely on the client side**:

   - `HttpOnly` cookie: protected from XSS, vulnerable to CSRF (mitigate with SameSite=Lax + CSRF token).
   - `localStorage`: accessible to JavaScript, vulnerable to XSS token theft.
   - For most web apps: `HttpOnly` + `SameSite=Lax` cookie is the right choice.

7. **Rotate signing keys periodically** — support multiple valid keys during rotation:

   ```python
   # Publish a JWKS endpoint for key discovery
   # kid (key ID) in JWT header identifies which key was used for signing
   # Keep old key active for existing tokens' remaining lifetime during rotation
   ```

## Rules

- Never decode JWT claims without verifying the signature first — `base64.decode(token.split('.')[1])` is not verification.
- `exp` must be validated server-side — a library that doesn't check `exp` by default is a security risk.
- Do not store sensitive data (PII, passwords) in JWT payloads — they are base64-encoded, not encrypted. Encrypt the JWT (JWE) if confidentiality is needed.
- Algorithm `none` must always be rejected — any library that accepts it should be replaced or wrapped.

## Common Mistakes

- **Long-lived tokens without revocation** — a 30-day JWT stolen once is valid for 30 days. Use short expiry + refresh token rotation.
- **Trusting `alg` from the token header** — the foundation of algorithm confusion attacks; always pin explicitly.
- **Storing JWTs in localStorage for convenience** — XSS on any page in your app steals the token.
- **Not including `aud` claim** — a token issued for service A accepted by service B because `aud` isn't checked is a privilege escalation.
