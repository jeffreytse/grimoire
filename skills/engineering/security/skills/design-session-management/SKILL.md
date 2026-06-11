---
name: design-session-management
description: Use when implementing user authentication state — creating, storing, transmitting, and expiring session tokens or cookies in a web application.
source: 'OWASP Session Management Cheat Sheet (owasp.org/www-project-cheat-sheets); OWASP Top 10 2021 A07; CWE-613; NIST SP 800-63B'
tags: [security, owasp, sessions, cookies, authentication, web, developer]
---

# Design Session Management

Issue cryptographically random session tokens, transmit them exclusively over HTTPS with `HttpOnly`/`Secure`/`SameSite` cookie attributes, and enforce idle and absolute timeouts — eliminating fixation, hijacking, and token prediction attacks.

## Why This Is Best Practice

**Adopted by:** Codified in OWASP Top 10 2021 A07 (Identification and Authentication Failures) and NIST SP 800-63B (Digital Identity Guidelines). Django, Rails, Laravel, Spring Session, and ASP.NET Core all implement these defaults. PCI DSS v4.0 Requirement 8.2 mandates session controls. NIST 800-63B is the authoritative federal standard for session management in US government applications.
**Impact:** Session hijacking is among the most common post-authentication attack vectors — Verizon DBIR 2023 identifies stolen session tokens as a top initial access technique in web application breaches. Proper session invalidation on logout prevents 100% of fixation attacks. HttpOnly cookies prevent JavaScript from reading session tokens even when XSS occurs, limiting the blast radius.
**Why best:** JWTs stored in localStorage are the common alternative — they lack server-side revocation, persist after logout, and are readable by XSS. Server-side sessions with `HttpOnly` cookies give developers full control over validity and are immune to JavaScript-based theft.

Sources: OWASP Session Management Cheat Sheet; NIST SP 800-63B; Verizon DBIR 2023; CWE-613

## Steps

1. **Generate cryptographically random session IDs** — at least 128 bits of entropy from a CSPRNG. Never use sequential IDs, timestamps, or user-derived values.

   ```python
   import secrets
   session_id = secrets.token_urlsafe(32)  # 256-bit URL-safe token
   ```

   ```java
   // Java
   SecureRandom sr = new SecureRandom();
   byte[] token = new byte[32];
   sr.nextBytes(token);
   String sessionId = Base64.getUrlEncoder().withoutPadding().encodeToString(token);
   ```

2. **Set all four protective cookie attributes**:

   ```http
   Set-Cookie: session=<token>; HttpOnly; Secure; SameSite=Lax; Path=/
   ```

   - `HttpOnly` — JavaScript cannot read the cookie (blocks XSS-based token theft)
   - `Secure` — cookie only sent over HTTPS
   - `SameSite=Lax` — blocks CSRF from cross-site form POST (see `prevent-csrf`)
   - `Path=/` — restrict to your app path, not subpaths of shared hosts

3. **Rotate session ID on privilege escalation** — generate a NEW session ID after login, role changes, or permission grants to prevent session fixation:

   ```python
   # After successful login
   old_data = session.get_data()
   session.invalidate()
   new_session = session.create_new()
   new_session.set_data(old_data)
   new_session.set('user_id', authenticated_user.id)
   ```

4. **Enforce idle and absolute timeouts**:

   ```python
   IDLE_TIMEOUT = 30 * 60        # 30 minutes of inactivity
   ABSOLUTE_TIMEOUT = 8 * 3600  # 8 hours regardless of activity

   last_activity = session['last_activity']
   created_at = session['created_at']

   if time.time() - last_activity > IDLE_TIMEOUT:
       session.invalidate()
   if time.time() - created_at > ABSOLUTE_TIMEOUT:
       session.invalidate()
   ```

   PCI DSS requires ≤15 minutes idle timeout for cardholder data environments.

5. **Invalidate server-side on logout** — delete the session record from the server, not just the client cookie:

   ```python
   @app.route('/logout', methods=['POST'])
   def logout():
       session_id = request.cookies.get('session')
       session_store.delete(session_id)  # server-side deletion
       response = redirect('/login')
       response.delete_cookie('session')
       return response
   ```

   Client-only logout (clearing the cookie without server deletion) leaves the token valid for hijacking.

6. **Use framework session management** — don't implement from scratch. Use Django's `django.contrib.sessions`, Spring's `HttpSession`, or Rails' `ActionDispatch::Session`. They handle storage, rotation, and timeout correctly.

7. **For high-security actions, re-authenticate** — prompt for password before account deletion, fund transfers, or email/password changes. Don't rely solely on session validity for irreversible operations.

## Rules

- Never store session IDs in URLs (`?sessionid=...`) — they appear in logs, referrer headers, and browser history.
- Session tokens must not contain user data (user ID, role) — store a random ID that maps server-side to user state.
- Concurrent session control (invalidate old sessions on new login) prevents credential-sharing; implement where security requirements demand it.
- Session fixation: always rotate ID on login; never accept a session ID set by the client pre-login.

## Common Mistakes

- **Using userID + timestamp as session token** — predictable, brute-forceable.
- **Deleting the cookie but not the server-side session on logout** — the token remains valid if captured.
- **Setting long absolute timeouts for convenience** — 30-day sessions mean a stolen token is valid for 30 days.
- **HttpOnly cookies + JWT in localStorage** — mixing models: choose one. HttpOnly cookies handle both transport and XSS protection.
