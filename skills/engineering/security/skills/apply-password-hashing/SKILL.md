---
name: apply-password-hashing
description: Use when storing user passwords or credentials in a database — covers algorithm selection, work factor configuration, and migration from legacy hashing schemes.
source: 'OWASP Password Storage Cheat Sheet (owasp.org/www-project-cheat-sheets); NIST SP 800-63B Section 5.1.1.2; CWE-916; OWASP Top 10 2021 A07'
tags: [security, owasp, passwords, bcrypt, argon2, hashing, authentication, developer]
---

# Apply Password Hashing

Store passwords with a memory-hard, salted, slow hash function (Argon2id, bcrypt, or scrypt) — never MD5, SHA-1, SHA-256, or unsalted hashes.

## Why This Is Best Practice

**Adopted by:** NIST SP 800-63B (authoritative US federal standard for digital identity) mandates slow hashing for memorized secrets. OWASP Top 10 2021 A07 (Authentication Failures) cites weak password storage as a primary failure mode. Django defaults to Argon2, ASP.NET Core Identity defaults to PBKDF2, Spring Security defaults to bcrypt. PCI DSS v4.0 Requirement 8.3.2 mandates strong hashing for stored passwords.
**Impact:** Every major password breach demonstrates the cost of weak hashing — LinkedIn (2012, SHA-1, 117M cracked in days), RockYou (2009, plaintext, 32M passwords), Adobe (2013, 3DES ECB, trivially reversible). bcrypt with cost factor 12 requires ~300ms per attempt; modern GPU rigs crack MD5 at 100 billion hashes/second. Proper hashing makes offline cracking economically infeasible for strong passwords.
**Why best:** Fast hashes (MD5, SHA-256) are designed for speed and destroy security — a single GPU cracks billions per second. Argon2id wins over bcrypt by resisting GPU attacks via memory-hardness; it wins over PBKDF2 by resisting ASICs. bcrypt is the safe fallback when Argon2 is unavailable.

Sources: OWASP Password Storage Cheat Sheet; NIST SP 800-63B; HaveIBeenPwned breach analysis; CWE-916

## Steps

1. **Choose the right algorithm by priority**:

   | Priority | Algorithm | Use when |
   |----------|-----------|----------|
   | 1st | Argon2id | Modern apps — best resistance to GPU/ASIC cracking |
   | 2nd | bcrypt | Argon2 unavailable; widely supported; proven 25+ years |
   | 3rd | scrypt | Alternative memory-hard option |
   | 4th | PBKDF2-HMAC-SHA256 | FIPS compliance required |
   | Never | MD5, SHA-*, plain | Not password hashing algorithms |

2. **Configure Argon2id with secure parameters** — OWASP minimum (2023):

   ```python
   from argon2 import PasswordHasher
   ph = PasswordHasher(
       time_cost=2,          # iterations
       memory_cost=19456,    # 19 MB memory
       parallelism=1,
       hash_len=32,
       salt_len=16
   )
   hash = ph.hash(password)
   ```

   Tune `memory_cost` and `time_cost` so hashing takes ≥500ms on your hardware.

3. **Configure bcrypt with work factor ≥12**:

   ```python
   import bcrypt
   salt = bcrypt.gensalt(rounds=12)     # never use rounds < 10
   hashed = bcrypt.hashpw(password.encode(), salt)
   ```

   ```javascript
   const bcrypt = require('bcrypt');
   const ROUNDS = 12;
   const hash = await bcrypt.hash(password, ROUNDS);
   ```

   Benchmark: increase `rounds` until hashing takes ~300–500ms on your hardware.

4. **Never manage salts manually** — Argon2, bcrypt, and scrypt generate and embed a random salt automatically. The output string includes the salt. Do NOT add a separate global "pepper" without understanding the implications.

   ```python
   # bcrypt output embeds the salt — store the full string:
   # $2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj9bIyZzAQm.
   ```

5. **Verify correctly — never re-hash and compare**:

   ```python
   # Argon2
   try:
       ph.verify(stored_hash, provided_password)
       if ph.check_needs_rehash(stored_hash):
           # rehash with updated parameters on next login
           new_hash = ph.hash(provided_password)
           update_hash_in_db(user_id, new_hash)
   except argon2.exceptions.VerifyMismatchError:
       raise AuthenticationError()
   ```

6. **Migrate legacy hashes without forcing password resets** — wrap old hashes:

   ```python
   # On login: detect old bcrypt hash, rehash with Argon2, store new
   if is_legacy_bcrypt(stored_hash):
       if bcrypt.checkpw(password, stored_hash):
           new_hash = ph.hash(password)
           update_to_argon2(user_id, new_hash)
   ```

7. **Enforce a maximum password length of 64–128 characters** — bcrypt silently truncates at 72 bytes; Argon2 does not, but extremely long passwords can cause DoS. Validate before hashing:

   ```python
   if len(password) > 128:
       raise ValueError("Password too long")
   ```

## Rules

- Always use the library's built-in salt generation — manual salting is error-prone.
- Store only the full hash string (which includes algorithm, parameters, and salt) — never store passwords or salts separately.
- Increase work factors annually — hardware gets faster; a 2015 bcrypt cost 10 is a 2025 cost 9 equivalent.
- Do NOT pre-hash passwords before bcrypt (e.g., `bcrypt(sha256(password))`) — introduces new vulnerabilities without benefit.

## Common Mistakes

- **Using `sha256(password + salt)` and calling it "secure"** — fast hashing with a salt is still fast; bcrypt with no salt is still better than SHA-256 with a salt.
- **Rolling your own hashing** — every custom implementation has missed bcrypt's cost factor, truncation handling, or encoding edge cases.
- **Storing the hash encoding as `hex` instead of the full output string** — loses the embedded parameters and salt.
- **Using rounds < 10 for bcrypt** — cost factors below 10 are crackable with consumer hardware.
