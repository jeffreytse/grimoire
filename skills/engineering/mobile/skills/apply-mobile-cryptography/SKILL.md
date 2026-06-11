---
name: apply-mobile-cryptography
description: Use when implementing cryptographic operations in a mobile app — encrypting data, hashing passwords, generating keys, signing data, or establishing secure communication channels between the app and backend.
source: 'OWASP Mobile Top 10 2024 M10 (owasp.org/www-project-mobile-top-10/); OWASP MASVS MSTG-CRYPTO; NIST SP 800-175B; Apple CryptoKit documentation; Android Keystore documentation'
tags: [security, owasp, mobile, cryptography, encryption, key-management, ios, android, developer]
---

# Apply Mobile Cryptography

Use platform-provided cryptographic APIs with approved algorithms (AES-256-GCM, ChaCha20-Poly1305, RSA-2048+) and hardware-backed key storage — never rolling custom crypto or hardcoding keys.

## Why This Is Best Practice

**Adopted by:** OWASP Mobile Top 10 2024 M10 (Insufficient Cryptography). OWASP MASVS MSTG-CRYPTO-1 through MSTG-CRYPTO-6 specify mobile cryptography requirements. Apple's CryptoKit and Android's Keystore API are the official platform cryptography interfaces. NIST SP 800-175B defines approved algorithms. PCI DSS v4.0 Requirement 3.5 mandates strong cryptography for stored cardholder data.
**Impact:** OWASP's Mobile Security Testing Guide (MSTG) consistently finds weak cryptography in the top 3 mobile app vulnerabilities during assessments. Common findings: AES-ECB mode encryption (identical blocks → visible patterns), hardcoded keys in source code (recovered via static analysis), MD5/SHA-1 for security purposes (collision-vulnerable), and custom XOR encryption (trivially reversible). Mobile app binary analysis tools (jadx, Hopper, frida) make hardcoded keys trivially recoverable.
**Why best:** Custom cryptography implementations (rolling your own cipher, key derivation, or protocol) are the alternative — they are routinely broken because subtle implementation errors invalidate mathematical security proofs. Platform-provided cryptographic primitives are reviewed, maintained, and hardware-accelerated.

Sources: OWASP Mobile Top 10 2024 M10; OWASP MASVS MSTG-CRYPTO; NIST SP 800-175B; Apple CryptoKit documentation; Android Keystore system documentation

## Steps

1. **iOS — use Apple CryptoKit for all cryptographic operations**:

   ```swift
   import CryptoKit

   // Symmetric encryption: AES-GCM
   func encrypt(data: Data, key: SymmetricKey) throws -> Data {
       let sealedBox = try AES.GCM.seal(data, using: key)
       return sealedBox.combined!  // nonce + ciphertext + tag
   }

   func decrypt(combined: Data, key: SymmetricKey) throws -> Data {
       let sealedBox = try AES.GCM.SealedBox(combined: combined)
       return try AES.GCM.open(sealedBox, using: key)
   }

   // Key generation: 256-bit key
   let key = SymmetricKey(size: .bits256)

   // Store key in Keychain (see apply-secure-mobile-storage)
   // Never hardcode or persist in UserDefaults
   ```

2. **Android — use Android Keystore for key generation and storage**:

   ```kotlin
   import android.security.keystore.KeyGenParameterSpec
   import android.security.keystore.KeyProperties
   import javax.crypto.Cipher
   import javax.crypto.KeyGenerator
   import javax.crypto.SecretKey

   fun generateKeyInKeystore(alias: String): SecretKey {
       val keyGenerator = KeyGenerator.getInstance(
           KeyProperties.KEY_ALGORITHM_AES, "AndroidKeyStore"
       )
       keyGenerator.init(
           KeyGenParameterSpec.Builder(alias,
               KeyProperties.PURPOSE_ENCRYPT or KeyProperties.PURPOSE_DECRYPT)
               .setBlockModes(KeyProperties.BLOCK_MODE_GCM)
               .setEncryptionPaddings(KeyProperties.ENCRYPTION_PADDING_NONE)
               .setKeySize(256)
               .setUserAuthenticationRequired(true)  // require biometrics/PIN
               .build()
       )
       return keyGenerator.generateKey()
   }

   fun encrypt(plaintext: ByteArray, key: SecretKey): Pair<ByteArray, ByteArray> {
       val cipher = Cipher.getInstance("AES/GCM/NoPadding")
       cipher.init(Cipher.ENCRYPT_MODE, key)
       val ciphertext = cipher.doFinal(plaintext)
       return Pair(cipher.iv, ciphertext)  // store IV alongside ciphertext
   }
   ```

3. **Algorithm selection guide for mobile**:

   | Use case | Recommended | Avoid |
   |---|---|---|
   | Symmetric encryption | AES-256-GCM, ChaCha20-Poly1305 | AES-ECB, AES-CBC without MAC, 3DES |
   | Key exchange | ECDH P-256 or P-384 | RSA key exchange (performance), DH <2048 |
   | Digital signatures | ECDSA P-256, Ed25519 | RSA <2048, SHA-1 |
   | Hashing (integrity) | SHA-256, SHA-3 | MD5, SHA-1 |
   | Password hashing | Argon2id, bcrypt | SHA-256(password), MD5 |
   | Random numbers | `SecRandomCopyBytes` (iOS), `KeyGenerator` (Android) | `Math.random()`, `java.util.Random` |

4. **Never hardcode cryptographic keys**:

   ```swift
   // BAD — key visible in binary analysis
   let hardcodedKey = "0123456789abcdef0123456789abcdef"

   // GOOD — generate per-device key in Secure Enclave
   func getOrCreateKey() -> SecureEnclave.P256.Signing.PrivateKey {
       if let existing = loadKeyFromKeychain(tag: "app.signing.key") {
           return existing
       }
       let newKey = try! SecureEnclave.P256.Signing.PrivateKey()
       saveToKeychain(key: newKey, tag: "app.signing.key")
       return newKey
   }
   ```

5. **Use platform-provided TLS, not custom certificate validation**:

   ```swift
   // iOS — do NOT override certificate validation
   // BAD
   func urlSession(_ session: URLSession, didReceive challenge: URLAuthenticationChallenge,
                   completionHandler: @escaping (URLSession.AuthChallengeDisposition, URLCredential?) -> Void) {
       completionHandler(.useCredential, URLCredential(trust: challenge.protectionSpace.serverTrust!))
       // Accepts any certificate — MitM vulnerability
   }

   // GOOD — use default URLSession (validates certificates automatically)
   let session = URLSession.shared
   ```

6. **Generate random values with platform CSPRNG**:

   ```swift
   // iOS
   import CryptoKit
   let randomBytes = SymmetricKey(size: .bits256).withUnsafeBytes { Data($0) }

   // Or using SecRandomCopyBytes
   var bytes = [UInt8](repeating: 0, count: 32)
   SecRandomCopyBytes(kSecRandomDefault, 32, &bytes)
   ```

   ```kotlin
   // Android
   import java.security.SecureRandom
   val random = SecureRandom()
   val bytes = ByteArray(32)
   random.nextBytes(bytes)
   ```

## Rules

- AES-ECB mode is never acceptable for any data longer than 16 bytes — identical input blocks produce identical ciphertext, revealing data patterns.
- IVs/nonces must be unique per encryption operation — reusing an IV with GCM mode completely breaks confidentiality.
- Keys derived from `SecRandomCopyBytes` or `SecureRandom` are CSPRNG-generated; keys derived from `Date().timeIntervalSince1970` or device IMEI are not.
- The Secure Enclave (iOS) and StrongBox (Android) cannot export keys — operations happen inside the hardware module. Use them for the highest-value keys.

## Common Mistakes

- **Using `AES/ECB/PKCS5Padding` because it's in the Android documentation examples** — older documentation used ECB; always use GCM.
- **Deriving an encryption key from the user's password with SHA-256** — SHA-256 is a hash, not a KDF; use Argon2id, scrypt, or PBKDF2.
- **Catching certificate validation exceptions and proceeding** — "just for testing" certificates in development code often reach production.
- **Storing IVs in a predictable location** — IVs can be public, but must be random and stored alongside the ciphertext.
