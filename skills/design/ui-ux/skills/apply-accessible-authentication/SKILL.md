---
name: apply-accessible-authentication
description: Use when building login, registration, or any authentication flow — to ensure users with cognitive disabilities can authenticate without being blocked by memory tests, cognitive puzzles, or disabled paste functionality.
source: W3C WCAG 2.2 SC 3.3.8 Accessible Authentication Minimum (Level AA), SC 3.3.9 Accessible Authentication Enhanced (Level AAA); WebAuthn W3C specification; NIST SP 800-63B digital identity guidelines
tags: [accessibility, wcag, a11y, authentication, password-manager, captcha, cognitive-accessibility, developer]
related: [apply-accessible-forms, design-accessibility-standards, audit-accessibility]
---

# Implement Accessible Authentication

Allow users to authenticate without relying on memory or cognitive tests — by supporting password managers, copy/paste, and providing CAPTCHA alternatives.

## Why This Is Best Practice

**Adopted by:** WCAG 2.2 SC 3.3.8 (Level AA, 2023) is now required wherever WCAG 2.2
is mandated — UK PSBAR updated to WCAG 2.2 AA in 2024; US Section 508 refresh is
expected to follow. NIST SP 800-63B (the US federal digital identity standard) explicitly
recommends allowing paste in password fields and supporting password managers as security
best practice — not just an accessibility requirement. Apple's Sign in with Apple and
Google's One Tap fulfill SC 3.3.8 by design.
**Impact:** 1 in 7 people have a cognitive disability (WHO). Blocking paste, disabling
password managers, or requiring complex CAPTCHA solutions disproportionately blocks
users with dyslexia, memory impairments, and ADHD. Troy Hunt (HaveIBeenPwned) documented
that blocking paste in password fields reduces password quality and increases reuse —
harming security AND accessibility.
**Why best:** Blocking paste to "prevent phishing" is a security myth disproved by NCSC
(National Cyber Security Centre, UK) and NIST. Password managers are a security
improvement. An authentication method that blocks assistive tools is simultaneously
less secure and less accessible.

Sources: W3C WCAG 2.2 SC 3.3.8 (2023); NIST SP 800-63B (2020); NCSC Password
Guidance (2021); Troy Hunt "The "Cobra Effect" of Blocking Paste on Password Fields"

## Steps

### Step 1: Never disable paste in password fields

```html
<!-- Wrong — blocks password managers, forces manual typing -->
<input type="password"
       onpaste="return false"
       oncopy="return false"
       oncut="return false">

<!-- Right — allow paste; support password manager autofill -->
<input type="password"
       id="password"
       name="password"
       autocomplete="current-password">
```

Disabling paste is a WCAG 2.2 SC 3.3.8 failure and a security anti-pattern.
No legitimate security justification exists for blocking paste.

### Step 2: Add correct `autocomplete` attributes to all auth fields

```html
<!-- Login form -->
<input type="text"     name="username" autocomplete="username">
<input type="password" name="password" autocomplete="current-password">

<!-- Registration form -->
<input type="email"    name="email"    autocomplete="email">
<input type="password" name="new-pw"   autocomplete="new-password">
<input type="password" name="confirm"  autocomplete="new-password">

<!-- Two-factor / OTP -->
<input type="text" name="otp" autocomplete="one-time-code"
       inputmode="numeric" pattern="[0-9]{6}">
```

Password managers and browser autofill depend on `autocomplete` values. Without them,
autofill is unreliable across browsers.

### Step 3: Provide a CAPTCHA alternative or replace CAPTCHA

CAPTCHA that requires image recognition, distorted text reading, or complex drag
interactions fails SC 3.3.8 for users with visual and cognitive disabilities.

Alternatives (from least to most user-friendly):
1. **Audio CAPTCHA alternative** — required if visual CAPTCHA is used (WCAG 1.1.1)
2. **Honeypot fields** — invisible to users; traps bots without user interaction
3. **Rate limiting + progressive challenge** — no CAPTCHA for normal use; only after
   suspicious behavior
4. **Cloudflare Turnstile / hCaptcha accessibility mode** — invisible challenge for
   verified humans, visible only for suspected bots
5. **Passkey / WebAuthn** — no CAPTCHA needed; authentication is cryptographic

```html
<!-- Honeypot example: hidden from users, bots fill it -->
<div style="display: none" aria-hidden="true">
  <label for="company">Company (leave blank)</label>
  <input type="text" id="company" name="company" tabindex="-1" autocomplete="off">
</div>
```

### Step 4: Support passkeys / WebAuthn as a CAPTCHA-free, memory-free option

WebAuthn (passkeys) satisfies SC 3.3.8 natively — no password to remember, no CAPTCHA.

```javascript
// Register a passkey
const credential = await navigator.credentials.create({
  publicKey: {
    challenge: serverChallenge,
    rp: { name: "Acme Corp", id: "acme.com" },
    user: { id: userId, name: userEmail, displayName: userName },
    pubKeyCredParams: [{ alg: -7, type: "public-key" }],
    authenticatorSelection: { userVerification: "preferred" }
  }
});

// Authenticate with passkey (no password required)
const assertion = await navigator.credentials.get({
  publicKey: { challenge: serverChallenge, rpId: "acme.com" }
});
```

Passkeys are supported in Chrome, Safari, Firefox (2023+) and on iOS/Android.

### Step 5: For required cognitive tests (image match, etc.) — provide a text alternative

If a cognitive function test cannot be removed, SC 3.3.8 requires at least one of:
- An alternative authentication method (email link, passkey)
- An object recognition alternative (select all images containing a bicycle — OK if
  an audio alternative is also available for each image)

The test must not require the user to memorize or transcribe text.

## When NOT to Use

- **Government high-assurance identity verification** — NIST AAL3 requires phishing-resistant hardware tokens. Password manager support still applies, but the overall authentication process may require in-person verification.
- **Financial transaction re-authentication** — re-auth for sensitive actions (wire transfers) may require additional friction. Still support paste and password managers for the credential entry step.

## Common Mistakes

**`autocomplete="off"` on login fields.** This suppresses password manager autofill. Banks and healthcare portals frequently do this citing security — it is both ineffective and a WCAG 2.2 AA violation.

**Image-only CAPTCHA without audio alternative.** A CAPTCHA that only shows distorted images with no audio or text alternative fails SC 1.1.1 (Level A) in addition to SC 3.3.8.

**Requiring CAPTCHA on every login attempt.** Rate limiting + progressive challenge is equally effective at blocking bots while eliminating CAPTCHA for legitimate users. Evaluate the real threat model before adding CAPTCHA friction.
