---
name: design-mobile-privacy-controls
description: Use when building a mobile app that collects, processes, or transmits user data — implementing data minimization, purpose limitation, consent management, and user data controls required by app store policies and privacy regulations.
source: 'OWASP Mobile Top 10 2024 M6 (owasp.org/www-project-mobile-top-10/); Apple App Store Review Guidelines 5.1; Google Play Data Safety; GDPR Article 5; CCPA; OWASP MASVS MSTG-STORAGE-11'
tags: [security, owasp, mobile, privacy, data-minimization, gdpr, consent, ios, android, developer]
---

# Design Mobile Privacy Controls

Implement data minimization, granular permission requests, user-facing privacy controls, and transparent data practices — meeting Apple App Store, Google Play, GDPR, and CCPA requirements while reducing data breach impact.

## Why This Is Best Practice

**Adopted by:** OWASP Mobile Top 10 2024 M6 (Inadequate Privacy Controls). Apple App Store Guidelines Section 5.1 and Google Play Data Safety policy both require privacy nutrition labels and consent for data collection. GDPR (EU, 2018) and CCPA (California, 2020) mandate data minimization, purpose limitation, and user data rights. Fines: GDPR fines have reached €1.2B (Meta, 2023); CCPA enforcement has resulted in millions in settlements.
**Impact:** Mobile apps are primary vectors for unauthorized data collection. Research by AppCensus (2022) found 70% of top Android apps share data with advertising networks beyond what users expect. Without privacy controls: fines under GDPR (up to 4% of global annual revenue), App Store/Play Store removal, and reputational damage. Apple's App Tracking Transparency (ATT) framework (2021) caused a $10B revenue impact on Facebook due to users opting out of tracking.
**Why best:** Collecting all possible data and deciding later what to use ("data lake" approach) creates regulatory liability and maximizes breach impact. Data minimization — collecting only what's necessary for the stated purpose — reduces both compliance burden and breach damage.

Sources: OWASP Mobile Top 10 2024 M6; Apple App Store Guidelines 5.1; Google Play Data Safety policy; GDPR Articles 5-7; CCPA Section 1798.100

## Steps

1. **Request permissions at the point of need, not at app launch**:

   ```swift
   // iOS — request location permission when the feature is triggered
   func userTappedFindNearbyButton() {
       // DON'T request location in AppDelegate.didFinishLaunching
       // DO request when the user triggers the feature
       locationManager.requestWhenInUseAuthorization()
   }

   // Provide usage description that explains WHY before the system dialog
   // Info.plist:
   // NSLocationWhenInUseUsageDescription: "We show nearby stores when you tap Find Nearby."
   ```

   ```kotlin
   // Android — request permission in context
   fun onFindNearbyClicked() {
       if (checkSelfPermission(Manifest.permission.ACCESS_FINE_LOCATION) != PERMISSION_GRANTED) {
           // Show rationale BEFORE requesting if previously denied
           if (shouldShowRequestPermissionRationale(Manifest.permission.ACCESS_FINE_LOCATION)) {
               showLocationRationaleDialog()
           } else {
               requestPermissions(arrayOf(Manifest.permission.ACCESS_FINE_LOCATION), REQUEST_CODE)
           }
       }
   }
   ```

2. **Implement data minimization** — collect only what the feature requires:

   ```swift
   // Location: use lowest precision required
   locationManager.desiredAccuracy = kCLLocationAccuracyKilometer  // not kCLLocationAccuracyBest
   // unless GPS-level precision is functionally necessary

   // For analytics: anonymize before sending
   struct AnalyticsEvent {
       let eventName: String
       let timestamp: Date
       // NOT: userId, deviceId, IP address
   }
   ```

3. **Implement user-facing privacy controls** (required by GDPR Article 17, CCPA Section 1798.105):

   ```swift
   // Privacy settings screen — expose all data controls
   class PrivacySettingsView {
       func showPrivacySettings() -> [PrivacyControl] {
           return [
               PrivacyControl(title: "Analytics",
                              description: "Helps us improve the app",
                              isEnabled: UserDefaults.analyticsEnabled,
                              onToggle: { enabled in
                                  Analytics.setEnabled(enabled)
                              }),
               PrivacyControl(title: "Delete My Data",
                              description: "Permanently deletes all your data",
                              action: { self.requestDataDeletion() }),
               PrivacyControl(title: "Export My Data",
                              description: "Download everything we have about you",
                              action: { self.requestDataExport() }),
           ]
       }
   }
   ```

4. **Complete Apple's App Privacy Nutrition Label accurately** (App Store requirement):

   ```
   Checklist for App Privacy disclosure:
   □ List all data types collected (location, contacts, identifiers, usage data, etc.)
   □ For each type: does it link to identity? Used for tracking?
   □ List all third-party SDKs and their data collection (advertising SDKs are major source)
   □ Remove analytics/advertising SDKs that collect data not needed for core functionality
   ```

5. **Audit third-party SDK data collection**:

   ```bash
   # iOS — use Emerge Tools or Privacy Manifest audit
   # Check PrivacyInfo.xcprivacy in each SDK

   # Android — review DATA_SAFETY declarations in play console
   # Use Mobile Security Framework (MobSF) to scan for tracking SDKs
   mobsf analyze --apk app.apk
   ```

   Remove or replace SDKs that collect data beyond what you disclose.

6. **Minimize data in crash reports and analytics**:

   ```swift
   // Configure crash reporter to exclude PII
   Crashlytics.crashlytics().setCrashlyticsCollectionEnabled(!isEUUser())

   // Scrub PII from logs before sending
   func sanitizeForAnalytics(_ event: [String: Any]) -> [String: Any] {
       var sanitized = event
       sanitized.removeValue(forKey: "email")
       sanitized.removeValue(forKey: "userId")  // use anonymous session ID instead
       return sanitized
   }
   ```

7. **Implement data retention limits** — delete data that's no longer needed:

   ```swift
   func enforceDataRetentionPolicy() {
       let retentionPeriod: TimeInterval = 90 * 24 * 3600  // 90 days
       let cutoff = Date().addingTimeInterval(-retentionPeriod)
       CoreDataStack.shared.deleteEvents(before: cutoff)
   }
   ```

## Rules

- Never request permissions not listed in your App Privacy disclosure — Apple and Google audit this.
- `ACCESS_FINE_LOCATION` (GPS-level) vs `ACCESS_COARSE_LOCATION` (city-level) — request the least precise permission that meets the feature requirement.
- Third-party advertising and analytics SDKs often collect more data than documented — audit each SDK's privacy manifest before including.
- GDPR consent must be freely given, specific, informed, and unambiguous — pre-ticked boxes and bundled consent are not valid.

## Common Mistakes

- **Requesting all permissions at app launch** — users deny permissions they don't understand, and app stores flag this behavior.
- **Not disclosing third-party SDK data collection** — the SDK's data collection is your responsibility to disclose.
- **Treating GDPR as EU-only** — CCPA, PIPEDA (Canada), LGPD (Brazil), and PDPA (Thailand/Singapore) all have similar requirements; design for the most stringent.
- **Collecting precise location when city-level suffices** — precise location is PII; coarse location reduces privacy impact significantly.
