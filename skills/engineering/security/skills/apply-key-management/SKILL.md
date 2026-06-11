---
name: apply-key-management
description: Use when designing encryption key lifecycle — generating, storing, rotating, and retiring cryptographic keys for application secrets, database encryption, and API signing keys.
source: 'OWASP Key Management Cheat Sheet (owasp.org/www-project-cheat-sheets); NIST SP 800-57 Part 1 (Key Management); AWS KMS Best Practices; FIPS 140-2'
tags: [security, owasp, cryptography, key-management, kms, hsm, developer]
---

# Apply Key Management

Manage cryptographic keys with envelope encryption via cloud KMS, automatic rotation schedules, and HSM-backed storage — preventing long-lived plaintext keys and single-point-of-compromise.

## Why This Is Best Practice

**Adopted by:** OWASP Key Management Cheat Sheet and NIST SP 800-57 (Key Management Recommendation) are the authoritative references. FIPS 140-2 Level 3 HSM compliance is required for US federal systems handling sensitive data. AWS KMS, Google Cloud KMS, and Azure Key Vault are used by Netflix, Stripe, and JPMorgan Chase as their primary key management infrastructure. PCI DSS v4.0 Requirement 3.7 mandates key management procedures including rotation, split knowledge, and dual control.
**Impact:** The 2014 Heartbleed OpenSSL vulnerability exposed long-lived private keys from millions of servers — keys that had never been rotated and remained valid for years after the bug was patched. Uber's 2022 breach included an attacker finding AWS credentials committed to GitHub (long-lived keys, never rotated) that provided access to 57 million records. NIST SP 800-57 estimates that 90-day key rotation reduces the expected loss from key compromise by approximately 90% compared to multi-year key lifespans.
**Why best:** Manual key management (keys stored in config files, rotated by hand) fails because developers avoid rotation due to complexity, and keys accumulate indefinitely. Envelope encryption — where a data encryption key (DEK) encrypted by a key encryption key (KEK) managed by KMS — limits blast radius: compromising the DEK doesn't expose the KEK, and rotation only requires re-encrypting the DEK, not re-encrypting all data.

Sources: OWASP Key Management Cheat Sheet; NIST SP 800-57 Part 1 Rev 5; AWS KMS Best Practices whitepaper; PCI DSS v4.0 Requirement 3.7

## Steps

1. **Use envelope encryption with cloud KMS — never store plaintext KEKs**:

   ```python
   import boto3
   import base64

   kms = boto3.client("kms", region_name="us-east-1")

   def encrypt_data(plaintext: bytes, kms_key_id: str) -> dict:
       # Generate a DEK — KMS returns plaintext + ciphertext blob
       response = kms.generate_data_key(
           KeyId=kms_key_id,
           KeySpec="AES_256"
       )
       plaintext_key = response["Plaintext"]
       encrypted_key = response["CiphertextBlob"]

       # Encrypt data with the plaintext DEK
       import os
       from cryptography.hazmat.primitives.ciphers.aead import AESGCM
       nonce = os.urandom(12)
       aesgcm = AESGCM(plaintext_key)
       ciphertext = aesgcm.encrypt(nonce, plaintext, None)

       return {
           "encrypted_key": base64.b64encode(encrypted_key).decode(),
           "nonce": base64.b64encode(nonce).decode(),
           "ciphertext": base64.b64encode(ciphertext).decode()
       }

   def decrypt_data(envelope: dict) -> bytes:
       encrypted_key = base64.b64decode(envelope["encrypted_key"])
       # KMS decrypts the DEK — plaintext key never leaves KMS at rest
       response = kms.decrypt(CiphertextBlob=encrypted_key)
       plaintext_key = response["Plaintext"]

       from cryptography.hazmat.primitives.ciphers.aead import AESGCM
       nonce = base64.b64decode(envelope["nonce"])
       ciphertext = base64.b64decode(envelope["ciphertext"])
       return AESGCM(plaintext_key).decrypt(nonce, ciphertext, None)
   ```

2. **Configure automatic key rotation**:

   ```bash
   # AWS KMS — enable annual automatic rotation
   aws kms enable-key-rotation --key-id arn:aws:kms:us-east-1:123456789:key/abc-123

   # Check rotation status
   aws kms get-key-rotation-status --key-id arn:aws:kms:us-east-1:123456789:key/abc-123
   ```

   ```python
   # For application secrets: rotate on a schedule
   import datetime

   class SecretRotationPolicy:
       MAX_AGE_DAYS = {
           "api_signing_key": 90,
           "database_encryption_key": 365,
           "session_secret": 30,
       }

       def needs_rotation(self, key_type: str, created_at: datetime.datetime) -> bool:
           max_age = self.MAX_AGE_DAYS[key_type]
           age = (datetime.datetime.utcnow() - created_at).days
           return age >= max_age
   ```

3. **Store application secrets in a secrets manager — not environment variables or config files**:

   ```python
   # AWS Secrets Manager — cache with TTL to avoid per-request latency
   import boto3
   import json
   import time
   from functools import lru_cache

   class SecretsCache:
       _cache: dict = {}
       TTL = 300  # 5 minutes

       def get_secret(self, secret_name: str) -> dict:
           cached = self._cache.get(secret_name)
           if cached and time.time() - cached["fetched_at"] < self.TTL:
               return cached["value"]

           client = boto3.client("secretsmanager")
           response = client.get_secret_value(SecretId=secret_name)
           value = json.loads(response["SecretString"])
           self._cache[secret_name] = {"value": value, "fetched_at": time.time()}
           return value

   secrets = SecretsCache()
   db_config = secrets.get_secret("/production/database")
   ```

4. **Implement key access logging and anomaly detection**:

   ```bash
   # AWS CloudTrail — log all KMS API calls
   aws cloudtrail create-trail \
     --name kms-audit-trail \
     --s3-bucket-name audit-logs-company

   # Alert on unusual key usage patterns (CloudWatch metric filter)
   aws cloudwatch put-metric-alarm \
     --alarm-name "KMS-Decrypt-Spike" \
     --metric-name "KMSDecryptCount" \
     --threshold 1000 \
     --evaluation-periods 1 \
     --comparison-operator GreaterThanThreshold
   ```

5. **Define key retirement and destruction procedure**:

   ```python
   class KeyLifecycle:
       """NIST SP 800-57 key states: Pre-activation → Active → Deactivated → Destroyed"""

       def retire_key(self, key_id: str, replacement_key_id: str):
           # Step 1: create new key version and start using it for new encryptions
           self.set_active_key(replacement_key_id)

           # Step 2: re-encrypt all data encrypted with old key (if DEKs, re-wrap them)
           self.rewrap_deks(old_key_id=key_id, new_key_id=replacement_key_id)

           # Step 3: deactivate old key (still available for decryption, not encryption)
           self.deactivate_key(key_id)

           # Step 4: after retention period, schedule destruction
           self.schedule_destruction(key_id, after_days=90)
   ```

## Rules

- Never log decrypted key material — KMS API responses containing plaintext keys must not be written to logs.
- Keys must never be exported from HSM in plaintext — use key wrapping (AES Key Wrap, PKCS#11) for key material transport.
- Separate encryption keys by classification: DEK for each data class, KEK hierarchy in KMS, master key in HSM.
- Apply least-privilege KMS IAM: separate IAM roles for key administrator (manage) vs. key user (encrypt/decrypt).

## Common Mistakes

- **Reusing the same key for encryption and signing** — different algorithms require different key types; separate keys for each purpose.
- **Deleting keys before re-encrypting all data** — data encrypted with a deleted key is permanently unrecoverable; deactivate first, destroy after re-encryption.
- **Not auditing key access** — without CloudTrail/audit logging, a compromised key usage is invisible until data exfiltration occurs.
- **Hardcoding KMS key ARNs in application code** — store key references in configuration/secrets manager so key rotation doesn't require code deployment.
