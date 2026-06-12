---
name: write-push-notification
description: Use when writing copy for push notifications or designing a push notification permission and delivery strategy
source: "Apple Human Interface Guidelines — Notifications; Google Material Design — Notifications; 'Mobile Push Notification Best Practices' (Urban Airship/Airship, 2023)"
tags: [mobile, push-notification, copywriting, ux, permissions, engagement, ios, android]
verified: true
---

# Write Push Notification

Write notification copy and design delivery timing so users find notifications useful rather than silencing them permanently.

## Why This Is Best Practice

**Adopted by:** Apple HIG mandates notification design standards for App Store approval; Google Material Design guidelines govern Android notification behavior; Airship (formerly Urban Airship) processes 15 billion+ notifications monthly and publishes research on engagement benchmarks

**Impact:** Airship research (2023) shows personalized notifications achieve 4-7× higher click-through than generic broadcast notifications; 60% of users who disable notifications never re-enable them — the first impression is permanent

**Why best:** Push notification permission is among the most fragile trust signals in mobile UX. Users who grant permission and then receive irrelevant notifications revoke it immediately and permanently. Every notification must pass a single test: does this require the user's attention right now, and does it give them enough information to decide whether to open the app?

## Steps

1. **Design permission priming before the system prompt** — show an in-app screen explaining exactly what notifications the user will receive and why, before triggering the OS permission dialog; this increases grant rates by 30-50% (Airship)
2. **Write the notification title** — 40 characters or fewer; must communicate the subject without the body; do not use the app name (it appears automatically)
3. **Write the notification body** — 100 characters or fewer on iOS; one complete thought; answer "what happened and what should I do?" without opening the app
4. **Include a deep link** — every notification must open directly to the relevant content, not the app home screen; a notification about an order must open that order
5. **Choose timing based on user behavior** — send during user's active hours using per-user send time optimization; never send between midnight and 6 AM local time without explicit user opt-in for time-sensitive alerts
6. **Segment and personalize** — use user attributes, behavior, and preferences to target; a "Your cart is waiting" notification is 5× more actionable than "Don't forget to shop"
7. **Measure opt-out rate per notification type** — if a notification type generates >5% opt-outs, pause and revise before resuming

## Rules

- Never ask for notification permission on app first open — establish value first
- Do not send more than one marketing notification per day per user
- Notification copy must not be deceptive — do not use urgency language ("Act now!") for non-urgent events
- Badge counts must reflect real unread items; inflating badges to drive opens is a deceptive pattern

## Examples

**Weak:** "Check out what's new in the app!"
**Strong:** "Your order #4821 has shipped — arrives Thursday. Track it."

**Permission priming:** "We'll notify you when your orders ship and when prices drop on items you've saved. You can control which alerts you receive in Settings."

## Common Mistakes

- Sending the OS permission prompt on the first screen of the first launch: users have no context for why they should grant permission
- Notification body that duplicates the title: both lines should add information, not repeat it
- Generic broadcast notifications with no personalization: they train users that notifications from your app are noise

## When NOT to Use

- The app is a consumer-facing tool where the user has not yet experienced core value — sending notifications before the user has completed onboarding will increase uninstall rates, not retention.
- The notification is purely informational with no required action and the same information is visible next time the user opens the app — interrupting the user adds no value and erodes permission trust.
- The team lacks backend infrastructure to segment users by behavior or attributes — sending the same blast notification to all users is an anti-pattern this skill explicitly prohibits, so apply this skill only after segmentation capability is in place.
