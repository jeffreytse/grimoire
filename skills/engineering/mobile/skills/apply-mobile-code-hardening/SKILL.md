---
name: apply-mobile-code-hardening
description: Use when preparing a mobile app for production release — especially apps handling financial data, DRM content, or anti-cheat systems — to raise the bar against reverse engineering, tampering, and runtime manipulation.
source: 'OWASP Mobile Top 10 2024 M7 (owasp.org/www-project-mobile-top-10/); OWASP MASVS MSTG-RESILIENCE; CWE-693; Google Play Integrity API; Apple DeviceCheck documentation'
tags: [security, owasp, mobile, code-hardening, obfuscation, anti-tamper, reverse-engineering, ios, android, developer]
---

# Apply Mobile Code Hardening

Apply code obfuscation, root/jailbreak detection, debug detection, and integrity verification to raise the cost of reverse engineering and runtime manipulation — critical for apps with high-value secrets, DRM, or anti-fraud requirements.

## Why This Is Best Practice

**Adopted by:** OWASP Mobile Top 10 2024 M7 (Insufficient Binary Protections). OWASP MASVS MSTG-RESILIENCE defines the standard for mobile binary protection. Google Play Integrity API and Apple DeviceCheck are platform-provided integrity attestation services. Financial services apps (banking, payment), gaming apps (anti-cheat), and DRM apps (streaming) universally apply these controls — required by PCI DSS v4.0 Requirement 6.3 for payment apps.
**Impact:** Mobile apps are distributed to untrusted devices — the reverse engineering tooling (jadx for Android, Hopper/Ghidra for iOS, Frida for both) is freely available and capable. Unprotected banking apps have had API keys extracted to build unauthorized clients. Gaming apps without anti-cheat have been reverse-engineered to enable automated bots. DRM content has been extracted from streaming apps. Code hardening raises the cost from minutes to weeks or months.
**Why best:** Code hardening is an arms race — determined attackers with resources will eventually break any protection. The goal is to raise the cost of attack beyond the attacker's expected return, not to make it impossible. Combined with server-side validation and monitoring, hardening reduces the class of opportunistic attackers effectively.

Sources: OWASP Mobile Top 10 2024 M7; OWASP MASVS MSTG-RESILIENCE; Google Play Integrity API documentation; Apple DeviceCheck documentation; CWE-693

## Steps

1. **Enable ProGuard/R8 code obfuscation on Android** (minimum baseline):

   ```groovy
   // build.gradle
   android {
       buildTypes {
           release {
               minifyEnabled true         // enables R8 code shrinking + obfuscation
               shrinkResources true       // removes unused resources
               proguardFiles getDefaultProguardFile('proguard-android-optimize.txt'),
                            'proguard-rules.pro'
           }
       }
   }
   ```

   ```
   # proguard-rules.pro — keep only what's required for reflection/serialization
   -keep class com.example.api.model.** { *; }  # keep data models for JSON parsing
   # Everything else: obfuscated by default
   ```

2. **Use Google Play Integrity API to verify app authenticity at runtime**:

   ```kotlin
   // Request integrity verdict before sensitive operations
   val integrityManager = IntegrityManagerFactory.create(context)

   suspend fun verifyAppIntegrity(): Boolean {
       val nonce = generateNonce()  // send to server for verification

       val request = StandardIntegrityManager.StandardIntegrityTokenRequest.builder()
           .setRequestHash(hashOfRequest(nonce))
           .build()

       val tokenResponse = integrityManager.requestIntegrityToken(request).await()
       val token = tokenResponse.token()

       // Send token to server — server decrypts and verifies with Google
       return backendApi.verifyIntegrity(token, nonce)
   }
   ```

3. **Detect root (Android) and jailbreak (iOS)**:

   ```kotlin
   // Android — use RootBeer library (community-maintained root detection)
   val rootBeer = RootBeer(context)
   if (rootBeer.isRooted) {
       // Log the detection — don't silently fail
       analytics.log("rooted_device_detected")
       // Policy decision: warn user, restrict features, or refuse operation
       showRootedDeviceWarning()
   }
   ```

   ```swift
   // iOS — jailbreak detection heuristics
   func isJailbroken() -> Bool {
       // Check for common jailbreak artifacts
       let jailbreakPaths = [
           "/Applications/Cydia.app",
           "/usr/sbin/sshd",
           "/etc/apt",
           "/private/var/lib/apt/",
       ]
       for path in jailbreakPaths {
           if FileManager.default.fileExists(atPath: path) { return true }
       }
       // Attempt to write outside sandbox
       let testPath = "/private/test_jailbreak_\(UUID().uuidString)"
       do {
           try "test".write(toFile: testPath, atomically: true, encoding: .utf8)
           try FileManager.default.removeItem(atPath: testPath)
           return true  // write succeeded — not possible on non-jailbroken device
       } catch { return false }
   }
   ```

4. **Detect debuggers and instrumentation frameworks at runtime**:

   ```swift
   // iOS — detect debugger attachment
   func isBeingDebugged() -> Bool {
       var info = kinfo_proc()
       var mib: [Int32] = [CTL_KERN, KERN_PROC, KERN_PROC_PID, getpid()]
       var size = MemoryLayout<kinfo_proc>.stride
       sysctl(&mib, UInt32(mib.count), &info, &size, nil, 0)
       return (info.kp_proc.p_flag & P_TRACED) != 0
   }
   ```

   ```kotlin
   // Android — detect Frida (dynamic instrumentation)
   fun isFridaDetected(): Boolean {
       return try {
           // Frida opens a TCP port on the device
           Socket("127.0.0.1", 27042).use { true }
       } catch (e: Exception) { false }
   }
   ```

5. **Use Apple DeviceCheck for persistent device-level fraud prevention**:

   ```swift
   import DeviceCheck

   func setDeviceFlag(bit: Bool, completion: @escaping (Bool) -> Void) {
       DCDevice.current.generateToken { token, error in
           guard let token = token else {
               completion(false)
               return
           }
           // Send token to your server — server uses DeviceCheck API
           // to set per-device bit (persists across reinstalls)
           backendAPI.recordDeviceToken(token.base64EncodedString()) { success in
               completion(success)
           }
       }
   }
   ```

6. **Verify app signature at runtime**:

   ```kotlin
   // Android — verify APK signing certificate hasn't changed
   fun verifyAppSignature(context: Context): Boolean {
       val expectedSignatureHash = "SHA256_OF_YOUR_RELEASE_CERT"
       val packageInfo = context.packageManager.getPackageInfo(
           context.packageName,
           PackageManager.GET_SIGNATURES
       )
       val signature = packageInfo.signatures[0]
       val certHash = MessageDigest.getInstance("SHA-256")
           .digest(signature.toByteArray())
           .toHexString()
       return certHash == expectedSignatureHash
   }
   ```

## Rules

- Code hardening is defense-in-depth, not a complete defense — all sensitive operations must still be validated server-side.
- Root/jailbreak detection should log and optionally degrade functionality — avoid outright blocking as it may affect legitimate users with custom ROMs.
- Detection checks must be spread throughout the code, not concentrated in one place — removing a single check should not bypass all protections.
- Obfuscation does not protect secrets hardcoded in the binary — no hardcoded keys, regardless of obfuscation level.

## Common Mistakes

- **Applying code hardening only to the release build without testing it** — obfuscation can break reflection-based code; test before releasing.
- **Treating root detection as a security boundary** — it reduces attack surface but is bypassable; server-side validation is the authoritative check.
- **Hardcoding API keys in the binary and relying on obfuscation to protect them** — obfuscated keys are still extractable by a motivated attacker; use certificate pinning and server-side key validation instead.
- **One-time integrity check at launch only** — move checks to multiple points in the code so patching the launch check doesn't disable all protection.
