---
name: design-auth-flow
description: Use when designing or reviewing authentication and authorization flows for web, mobile, or API systems
source: OAuth 2.0 RFC 6749; PKCE RFC 7636; OWASP Authentication Cheat Sheet (cheatsheetseries.owasp.org)
tags: [authentication, oauth2, pkce, security, authorization, owasp]
verified: true
---

# Design Auth Flow

Design secure authentication and authorization flows using OAuth 2.0 / OIDC standards with PKCE and OWASP best practices.

## Why This Is Best Practice

**Adopted by:** Google, Microsoft, GitHub, Okta, Auth0 — all implement OAuth 2.0 + PKCE as the baseline
**Impact:** OWASP reports that broken authentication is consistently in the top 3 web vulnerabilities; PKCE (RFC 7636) eliminates the authorization code interception attack that affected millions of mobile apps.

**Why best:** Rolling custom auth is the single highest-risk decision in software security. OAuth 2.0 + OIDC provides a peer-reviewed, widely-audited framework. PKCE extends it safely to public clients (SPAs, mobile apps) where client secrets cannot be stored securely.

## Steps

1. **Identify client type** — Confidential client (server-side, can store secrets) → Authorization Code Flow. Public client (SPA, mobile) → Authorization Code Flow + PKCE. Machine-to-machine → Client Credentials Flow.
2. **Never use Implicit Flow** — It is deprecated in OAuth 2.0 Security BCP (RFC 9700); use PKCE instead for all browser clients.
3. **Implement PKCE** — Generate a cryptographically random `code_verifier` (43-128 chars); hash it to `code_challenge` (S256); send challenge in auth request; send verifier in token exchange.
4. **Validate tokens correctly** — Verify JWT signature (RS256/ES256), `iss`, `aud`, `exp`, and `nbf` claims. Never trust unsigned tokens or skip expiry checks.
5. **Secure token storage** — Access tokens: memory only (not localStorage). Refresh tokens: httpOnly, Secure, SameSite=Strict cookies for web; secure enclave / keychain for mobile.
6. **Implement MFA** — Require TOTP or WebAuthn for sensitive operations; make MFA opt-in baseline for all users.
7. **Apply OWASP hardening** — Rate-limit login endpoints, implement account lockout, use constant-time comparison for secrets, log all auth events.

## Rules

- Never store plaintext passwords — use bcrypt (cost ≥12), scrypt, or Argon2id.
- Rotate refresh tokens on every use (refresh token rotation); invalidate the old token immediately.
- Short-lived access tokens (15 minutes), longer refresh tokens (7-30 days) with rotation.
- Log auth events (success, failure, MFA bypass) to a tamper-evident audit log.

## Examples

PKCE flow (SPA):
1. App generates `code_verifier = random(64)`, `code_challenge = base64url(sha256(verifier))`.
2. Redirect to `/authorize?response_type=code&code_challenge=...&code_challenge_method=S256`.
3. Exchange code + verifier at `/token` — server verifies `sha256(verifier) == challenge`.
4. Store access token in memory; refresh token in httpOnly cookie.

## Common Mistakes

- **Storing JWTs in localStorage** — XSS attacks steal tokens; use memory + httpOnly cookies.
- **Long-lived access tokens** — a stolen token is valid until expiry; keep them short-lived.
- **Skipping `state` parameter** — enables CSRF attacks on the OAuth callback.

## When NOT to Use

- When the system is an internal CLI tool or offline script used exclusively by a single authenticated OS user, OAuth 2.0 delegation flows add needless complexity; OS-level identity or a static API key with environment-scoped access is sufficient.
- When designing auth for a purely machine-to-machine backend service running in a trusted private network with no human login surface, PKCE and OIDC session management are irrelevant; focus on mTLS and service account credential rotation instead.
- When an existing identity provider (e.g., Auth0, Cognito, Okta) is already integrated and its default flows satisfy security requirements, redesigning the auth flow risks introducing vulnerabilities into a battle-tested configuration without added benefit.
