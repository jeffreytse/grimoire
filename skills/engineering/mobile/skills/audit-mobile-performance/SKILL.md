---
name: audit-mobile-performance
description: Use when diagnosing performance problems in a mobile app or conducting a systematic mobile performance review
source: "Android Vitals (Google Play Console); Apple MetricKit (Apple Developer Documentation); 'Mobile Performance' — Google Developers; RAIL model"
tags: [mobile, performance, android-vitals, metrickit, profiling, startup, frame-rate]
verified: true
---

# Audit Mobile Performance

Systematically identify and quantify performance bottlenecks using platform-native metrics before optimizing anything.

## Why This Is Best Practice

**Adopted by:** Google Play's Android Vitals program affects app store ranking; Apple App Store reviews MetricKit data for featured app eligibility; both platforms penalize apps that exceed thresholds in app store visibility algorithms

**Impact:** Google reports a 1-second improvement in mobile page speed improves conversions by up to 27% (Google/SOASTA research, 2017); Apple's MetricKit data shows apps with launch times above 2 seconds have 3× higher 1-day uninstall rates

**Why best:** Performance optimization without measurement is guesswork. Platform-provided metrics (Android Vitals, MetricKit) represent real-device, real-user data that synthetic benchmarks cannot replicate. Always measure first; optimize the metric with the largest user impact.

## Steps

1. **Collect baseline metrics** — pull Android Vitals from Play Console (ANR rate, crash rate, slow render rate, frozen frame rate) and MetricKit reports from Xcode Organizer (launch time, hang rate, memory); establish p50, p75, and p95 values
2. **Audit cold start time** — target: Android under 5 seconds cold start, iOS under 2 seconds (Apple's threshold for background task termination); profile with Android Profiler / Instruments Time Profiler
3. **Audit frame rendering** — target: 16ms per frame (60 fps) or 8ms (120 fps); identify dropped frames using Systrace (Android) or Instruments Core Animation (iOS); find methods executing on the main thread
4. **Audit memory usage** — profile heap allocations; identify memory leaks using LeakCanary (Android) or Instruments Allocations (iOS); target < 150MB resident memory for mid-range devices
5. **Audit network layer** — measure time-to-first-byte, payload sizes, and connection reuse; use Charles Proxy or Proxyman to inspect real traffic; compress payloads and implement request batching
6. **Audit battery impact** — use Android Battery Historian or Xcode Energy Organizer; identify wake lock abuse, excessive location polling, and background work exceeding platform budgets
7. **Prioritize by user impact** — rank findings by: percentage of users affected × severity × ease of fix; do not optimize a p99 cold start at the expense of an ANR affecting 5% of users

## Rules

- Always profile on real low-end devices, not simulators or flagship phones
- Never optimize before measuring: profiling reveals the actual bottleneck, which is rarely where engineers assume
- Android Vitals thresholds (crash rate >1.09%, ANR rate >0.47%) trigger Play Store ranking penalties — fix these first
- Memory leaks compound: one leak on a screen visited frequently will OOM crash after 10-20 visits

## Examples

**Finding:** Cold start p75 = 4.2s on Android. Profiler shows `SharedPreferences.getAll()` blocking main thread for 1.8s at startup. **Fix:** Move to async DataStore; cold start drops to 2.1s.

## Common Mistakes

- Optimizing only on flagship devices: performance issues surface on mid-range and low-end devices with slower CPUs and less RAM
- Measuring in the simulator: iOS Simulator uses host Mac CPU and memory; Instruments data from simulator is not representative of real-device behavior
- Fixing ANRs by moving work to background threads without checking the work is actually safe to parallelize

## When NOT to Use

- The app has not yet reached production or has fewer than 1,000 daily active users — insufficient real-device data makes platform metrics (Android Vitals, MetricKit) statistically unreliable for prioritization.
- The reported performance complaint stems from server-side latency exceeding 3 seconds — mobile profiling tools cannot diagnose backend bottlenecks and the audit effort should be redirected to backend observability.
- The development team is mid-sprint on a feature with a hard customer deadline and no existing baseline metrics — auditing without an established baseline produces findings that cannot be acted on without risking the release.
