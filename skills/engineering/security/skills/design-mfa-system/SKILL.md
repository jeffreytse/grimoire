---
name: design-mfa-system
description: Use when adding multi-factor authentication to an application — implementing TOTP, WebAuthn/passkeys, recovery codes, and phishing-resistant options for high-value accounts.
source: 'OWASP MFA Cheat Sheet (owasp.org/www-project-cheat-sheets); NIST SP 800-63B (Digital Identity Guidelines); WebAuthn Level 2 specification (W3C); FIDO Alliance implementation guidance'
tags: [security, owasp, mfa, totp, webauthn, authentication, identity, developer]
---

# Design MFA System

Implement TOTP as baseline MFA, WebAuthn/passkeys as phishing-resistant option, secure recovery codes, and account lockout — protecting accounts even when passwords are compromised.

## Why This Is Best Practice

**Adopted by:** OWASP MFA Cheat Sheet and NIST SP 800-63B (AAL2/AAL3 requirements) are the authoritative references. Google, GitHub, Stripe, and all major financial institutions mandate MFA for administrative accounts. WebAuthn (FIDO2) is supported by all major browsers and operating systems since 2019. NIST SP 800-63B mandates MFA for any system at Authenticator Assurance Level 2 or 3, which includes most financial, healthcare, and government systems.
**Impact:** Microsoft's 2019 analysis of 1.2 million compromised accounts found that 99.9% could have been prevented by MFA. Verizon DBIR 2023 found that stolen credentials account for 49% of breaches — MFA prevents credential-stuffing and password spray attacks from succeeding even with valid credentials. The 2022 Twilio breach compromised SMS MFA codes via phishing — illustrating why TOTP and WebAuthn are more secure than SMS, which is susceptible to SIM swapping and phishing relay attacks.
**Why best:** SMS OTP (the alternative) is susceptible to SIM swapping attacks (real-time account takeover) and phishing relay attacks where an attacker proxies a live login. TOTP is significantly more resistant than SMS. WebAuthn/passkeys bind the authenticator to the origin (domain), making phishing impossible — the credential only works on the legitimate site. Combining TOTP baseline with WebAuthn option provides both broad coverage and maximum security for high-value accounts.

Sources: OWASP MFA Cheat Sheet; NIST SP 800-63B; WebAuthn Level 2 W3C specification; FIDO Alliance implementation guide; Microsoft Identity Security (2019)

## Steps

1. **Implement TOTP (Time-based One-Time Password) as baseline MFA**:

   ```python
   import pyotp
   import qrcode
   import io
   import base64
   from cryptography.fernet import Fernet

   class TOTPManager:
       def __init__(self, encryption_key: bytes):
           self.fernet = Fernet(encryption_key)

       def generate_totp_secret(self) -> str:
           # 160-bit secret (base32 encoded) — NIST SP 800-63B requirement
           return pyotp.random_base32()

       def get_provisioning_qr(self, user_email: str, secret: str, issuer: str) -> str:
           totp = pyotp.TOTP(secret)
           uri = totp.provisioning_uri(name=user_email, issuer_name=issuer)
           # Return QR code as base64 PNG
           qr = qrcode.make(uri)
           buffer = io.BytesIO()
           qr.save(buffer, format="PNG")
           return base64.b64encode(buffer.getvalue()).decode()

       def verify_totp(self, secret: str, token: str) -> bool:
           totp = pyotp.TOTP(secret)
           # valid_window=1 allows ±30 seconds clock drift
           return totp.verify(token, valid_window=1)

       def encrypt_secret_for_storage(self, secret: str) -> bytes:
           return self.fernet.encrypt(secret.encode())

       def decrypt_secret(self, encrypted: bytes) -> str:
           return self.fernet.decrypt(encrypted).decode()

   # Store encrypted secret per user
   # Never store raw TOTP secret in database
   ```

2. **Generate and store recovery codes**:

   ```python
   import secrets
   import bcrypt

   def generate_recovery_codes(count: int = 10) -> tuple[list[str], list[bytes]]:
       codes = []
       hashes = []
       for _ in range(count):
           # 10-digit code with hyphen for readability: XXXXX-XXXXX
           code = f"{secrets.randbelow(100000):05d}-{secrets.randbelow(100000):05d}"
           codes.append(code)
           hashes.append(bcrypt.hashpw(code.encode(), bcrypt.gensalt(rounds=10)))
       return codes, hashes
       # Return plaintext codes to user ONCE (display at setup, never again)
       # Store hashes in database

   def use_recovery_code(user_id: str, submitted_code: str, stored_hashes: list[bytes]) -> bool:
       code_clean = submitted_code.replace("-", "").replace(" ", "")
       for i, stored_hash in enumerate(stored_hashes):
           submitted_normalized = f"{code_clean[:5]}-{code_clean[5:]}"
           if bcrypt.checkpw(submitted_normalized.encode(), stored_hash):
               # Invalidate used code — recovery codes are single-use
               delete_recovery_code(user_id, i)
               return True
       return False
   ```

3. **Implement WebAuthn/passkeys for phishing-resistant MFA**:

   ```python
   from webauthn import (
       generate_registration_options,
       verify_registration_response,
       generate_authentication_options,
       verify_authentication_response,
   )
   from webauthn.helpers.structs import (
       AuthenticatorSelectionCriteria,
       UserVerificationRequirement,
       ResidentKeyRequirement,
   )

   RP_ID = "app.company.com"
   RP_NAME = "Company App"

   def begin_registration(user_id: str, user_email: str):
       options = generate_registration_options(
           rp_id=RP_ID,
           rp_name=RP_NAME,
           user_id=user_id.encode(),
           user_name=user_email,
           authenticator_selection=AuthenticatorSelectionCriteria(
               user_verification=UserVerificationRequirement.REQUIRED,
               resident_key=ResidentKeyRequirement.PREFERRED,
           ),
       )
       # Store challenge in session for verification
       session["webauthn_challenge"] = options.challenge
       return options

   def complete_registration(user_id: str, credential_response: dict):
       verification = verify_registration_response(
           credential=credential_response,
           expected_challenge=session["webauthn_challenge"],
           expected_rp_id=RP_ID,
           expected_origin=f"https://{RP_ID}",
           require_user_verification=True,
       )
       # Store credential for user — use for future authentication
       store_credential(user_id, {
           "id": verification.credential_id,
           "public_key": verification.credential_public_key,
           "sign_count": verification.sign_count,
       })
   ```

4. **Rate limit MFA verification attempts**:

   ```python
   from datetime import datetime, timedelta
   import redis

   redis_client = redis.Redis()

   class MFARateLimiter:
       MAX_ATTEMPTS = 5
       LOCKOUT_MINUTES = 15

       def check_and_record(self, user_id: str) -> bool:
           key = f"mfa_attempts:{user_id}"
           current = redis_client.get(key)

           if current and int(current) >= self.MAX_ATTEMPTS:
               return False  # locked out

           pipe = redis_client.pipeline()
           pipe.incr(key)
           pipe.expire(key, int(timedelta(minutes=self.LOCKOUT_MINUTES).total_seconds()))
           pipe.execute()
           return True

       def clear_on_success(self, user_id: str):
           redis_client.delete(f"mfa_attempts:{user_id}")

   def verify_mfa(user_id: str, token: str) -> dict:
       limiter = MFARateLimiter()
       if not limiter.check_and_record(user_id):
           return {"success": False, "error": "Too many attempts. Try again in 15 minutes."}

       user = get_user(user_id)
       if verify_totp(user.totp_secret, token):
           limiter.clear_on_success(user_id)
           return {"success": True}

       return {"success": False, "error": "Invalid code"}
   ```

5. **Enforce MFA at session creation — don't allow bypassing**:

   ```python
   # Multi-step login: password → MFA → session creation
   def login_step1(email: str, password: str) -> dict:
       user = authenticate_password(email, password)
       if not user:
           return {"success": False, "error": "Invalid credentials"}

       # Create partial session — only MFA verification completes it
       partial_token = create_partial_session_token(user.id, expires_minutes=5)
       return {
           "success": True,
           "requires_mfa": user.mfa_enabled,
           "mfa_types": user.mfa_types,  # totp, webauthn, recovery
           "partial_token": partial_token,
       }

   def login_step2(partial_token: str, mfa_code: str) -> dict:
       user_id = verify_partial_session(partial_token)
       if not user_id:
           return {"success": False, "error": "Session expired"}

       if not verify_mfa(user_id, mfa_code)["success"]:
           return {"success": False, "error": "Invalid MFA code"}

       # Only now create full authenticated session
       session_token = create_full_session(user_id)
       return {"success": True, "session_token": session_token}
   ```

## Rules

- SMS OTP must not be used as sole MFA factor — it's vulnerable to SIM swapping; offer TOTP or WebAuthn as alternatives.
- Recovery codes must be single-use and hashed in storage — displaying them again after initial setup allows an attacker who gains database access to use them.
- MFA bypass tokens (device-remember cookies) must be scoped to a specific device, expire in ≤30 days, and be invalidated on password change.
- WebAuthn challenges must be cryptographically random, at least 16 bytes, and expire within 5 minutes — stale challenges allow replay attacks.

## Common Mistakes

- **Implementing MFA as client-side validation only** — MFA verification must happen server-side; client-side JavaScript checks are trivially bypassed.
- **TOTP with valid_window > 2** — allowing more than ±60 seconds of clock drift significantly widens the brute-force window.
- **Storing recovery codes in plaintext** — recovery codes are as sensitive as passwords; hash with bcrypt before storage.
- **Not invalidating existing sessions after MFA enrollment** — a session created before MFA was enabled should require re-authentication.
