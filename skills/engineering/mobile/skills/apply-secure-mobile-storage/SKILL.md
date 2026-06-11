---
name: apply-secure-mobile-storage
description: Use when storing sensitive data on a mobile device — credentials, tokens, PII, health data, financial data, or any data that should not be accessible to other apps or readable from device backups.
source: 'OWASP Mobile Top 10 2024 M9 (owasp.org/www-project-mobile-top-10/); OWASP Mobile Application Security Verification Standard (MASVS); CWE-312; Apple Keychain documentation; Android Keystore documentation'
tags: [security, owasp, mobile, secure-storage, keychain, keystore, ios, android, developer]
---

# Apply Secure Mobile Storage

Store sensitive mobile app data using platform-provided secure enclaves (iOS Keychain, Android Keystore) and encrypted databases — never in plaintext files, SharedPreferences, or UserDefaults.

## Why This Is Best Practice

**Adopted by:** OWASP Mobile Top 10 2024 M9 (Insecure Data Storage) is a top mobile vulnerability. OWASP MASVS (Mobile Application Security Verification Standard) MSTG-STORAGE-1 through MSTG-STORAGE-14 specify secure storage requirements. Apple's iOS Security Guide mandates Keychain for sensitive data. Google's Android Security Best Practices require Keystore for cryptographic keys. PCI DSS v4.0 Requirement 3 and HIPAA both mandate encryption for stored sensitive data on mobile.
**Impact:** The OWASP Mobile Top 10 consistently ranks insecure data storage as a critical finding across mobile app security assessments. Common findings: app tokens stored in plaintext SQLite databases (accessible after physical device access or backup extraction), encryption keys stored in UserDefaults/SharedPreferences (readable without root on unencrypted devices), and PII cached in unencrypted temporary files. Physical device access, iCloud/adb backup extraction, and root/jailbreak all enable reading non-secure storage.
**Why best:** Encrypting files with a key stored in the same location as the encrypted data is the common alternative — it's equivalent to locking a safe and taping the combination to the door. Platform-provided secure hardware enclaves (Secure Enclave on iOS, TEE/StrongBox on Android) store keys in hardware that cannot be extracted even with root access.

Sources: OWASP Mobile Top 10 2024 M9; OWASP MASVS MSTG-STORAGE; Apple iOS Security Guide; Android Security Best Practices; CWE-312

## Steps

1. **iOS — use Keychain for all secrets and tokens**:

   ```swift
   import Security

   func saveToKeychain(key: String, data: Data) -> Bool {
       let query: [String: Any] = [
           kSecClass as String: kSecClassGenericPassword,
           kSecAttrAccount as String: key,
           kSecValueData as String: data,
           kSecAttrAccessible as String: kSecAttrAccessibleWhenUnlockedThisDeviceOnly
           // Never use kSecAttrAccessibleAlways — accessible even when locked
       ]
       SecItemDelete(query as CFDictionary)  // remove existing
       return SecItemAdd(query as CFDictionary, nil) == errSecSuccess
   }

   func loadFromKeychain(key: String) -> Data? {
       let query: [String: Any] = [
           kSecClass as String: kSecClassGenericPassword,
           kSecAttrAccount as String: key,
           kSecReturnData as String: true
       ]
       var result: AnyObject?
       SecItemCopyMatching(query as CFDictionary, &result)
       return result as? Data
   }
   ```

   Use `kSecAttrAccessibleWhenUnlockedThisDeviceOnly` for most secrets — requires device unlock, non-exportable.

2. **Android — use EncryptedSharedPreferences and Android Keystore**:

   ```kotlin
   import androidx.security.crypto.EncryptedSharedPreferences
   import androidx.security.crypto.MasterKey

   // Create master key backed by Android Keystore
   val masterKey = MasterKey.Builder(context)
       .setKeyScheme(MasterKey.KeyScheme.AES256_GCM)
       .build()

   // EncryptedSharedPreferences — transparently encrypted
   val securePrefs = EncryptedSharedPreferences.create(
       context,
       "secure_prefs",
       masterKey,
       EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
       EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
   )

   // Use like normal SharedPreferences
   securePrefs.edit().putString("auth_token", token).apply()
   val token = securePrefs.getString("auth_token", null)
   ```

3. **Never use these for sensitive data**:

   | Storage mechanism | Why it's insecure |
   |---|---|
   | `UserDefaults` (iOS) | Plaintext plist, readable from backup |
   | `SharedPreferences` without encryption | Plaintext XML, readable with root |
   | SQLite without encryption | Readable from backup or filesystem |
   | Files in Documents/ (iOS) | Included in iCloud backup by default |
   | External storage (Android) | World-readable before Android 10 |
   | Log output | Readable from ADB on debug builds |

4. **Use SQLCipher for encrypted databases**:

   ```kotlin
   // Android — SQLCipher
   import net.sqlcipher.database.SQLiteDatabase

   SQLiteDatabase.loadLibs(context)
   val db = SQLiteDatabase.openOrCreateDatabase(dbPath, getDatabaseKey(), null)
   // getDatabaseKey() should return key from Android Keystore

   // iOS — SQLCipher (via Swift Package Manager)
   let db = try Connection(dbPath)
   try db.key(getKeyFromKeychain())
   ```

5. **Exclude sensitive files from backups**:

   ```swift
   // iOS — mark file as excluded from iCloud backup
   var url = URL(fileURLWithPath: sensitiveFilePath)
   var values = URLResourceValues()
   values.isExcludedFromBackup = true
   try url.setResourceValues(values)
   ```

   ```xml
   <!-- Android — exclude files from Auto Backup (AndroidManifest.xml) -->
   <application
       android:allowBackup="true"
       android:fullBackupContent="@xml/backup_rules">

   <!-- res/xml/backup_rules.xml -->
   <full-backup-content>
       <exclude domain="database" path="sensitive.db"/>
       <exclude domain="sharedpref" path="secure_prefs.xml"/>
   </full-backup-content>
   ```

6. **Clear sensitive data from memory when no longer needed**:

   ```swift
   // Overwrite sensitive strings in memory
   var password = getPassword()
   defer {
       password = String(repeating: "\0", count: password.count)
   }
   authenticate(with: password)
   ```

## Rules

- Keys that decrypt sensitive data must live in the Keychain/Keystore — never in the database or alongside the data.
- Use `kSecAttrAccessibleWhenUnlockedThisDeviceOnly` / `setUserAuthenticationRequired(true)` for the highest-value secrets.
- On Android, `StrongBox`-backed keys are stored in dedicated hardware security module — use for the most sensitive keys on supported devices.
- Debug builds often log sensitive data — use `#if DEBUG` / `BuildConfig.DEBUG` guards to prevent logging tokens in release builds.

## Common Mistakes

- **Storing JWT access tokens in UserDefaults/SharedPreferences** — the most common finding in mobile security assessments.
- **Encrypting a database with a key derived only from the app's bundle ID** — predictable key defeats encryption.
- **Using the same Keychain item across all apps with the same bundle prefix** — set explicit `kSecAttrAccessGroup` to restrict access.
- **Logging sensitive data to NSLog/Log.d without stripping in release** — crashes and analytics platforms may capture logs.
