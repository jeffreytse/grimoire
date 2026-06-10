---
name: apply-timed-content-controls
description: Use when building carousels, auto-advancing slideshows, session management, or any content that moves, blinks, or updates automatically — to ensure users can pause, stop, or extend time-sensitive content.
source: W3C WCAG 2.1 SC 2.2.1 Timing Adjustable (Level A), SC 2.2.2 Pause Stop Hide (Level A), SC 2.2.6 Timeouts (Level AA); WAI-ARIA Carousel Design Pattern; NIST SP 800-63B session management guidance
tags: [accessibility, wcag, a11y, timing, carousels, session-timeout, auto-advancing, pause-controls]
related: [apply-keyboard-accessibility, apply-reduced-motion, design-accessibility-standards]
---

# Implement Timed Content Controls

Provide user controls to pause, stop, or extend moving content and session timelines — so users with motor or cognitive disabilities are not blocked by content that advances without their control.

## Why This Is Best Practice

**Adopted by:** WCAG 2.1 SCs 2.2.1 and 2.2.2 (Level A) are required by Section 508 (US),
EU EN 301 549, and UK PSBAR 2018. SC 2.2.6 (Level AA) is required by jurisdictions mandating
WCAG 2.1 AA. The WAI-ARIA Carousel Pattern provides the reference implementation for accessible
auto-advancing carousels. UK GOV.UK Design System and the US Web Design System both require pause
controls on animated content.
**Impact:** Users with motor disabilities who have slow response times cannot interact with
auto-advancing content before it advances. Users with cognitive disabilities who read slowly cannot
process content before it disappears. WHO reports 16% of adults have a significant motor or cognitive
disability affecting interaction speed. Auto-advancing carousels — used on 72% of major e-commerce
homepages — fail WCAG 2.2.2 without a pause button.
**Why best:** CSS `prefers-reduced-motion` respects a system-level preference but does not give users
control over specific auto-advancing elements. WCAG 2.2.2 requires an explicit user-controlled pause
mechanism — not just a media query. Session timeouts without warning (failing 2.2.6) are a separate
problem: users with cognitive disabilities completing forms slowly lose unsaved data.

Sources: W3C WCAG 2.1 SCs 2.2.1, 2.2.2, 2.2.6 (2018); WAI-ARIA Carousel Pattern; GOV.UK Design System;
NIST SP 800-63B

## Steps

### Step 1: Add pause/play controls to all auto-advancing content (2.2.2)

Any content that moves, blinks, scrolls, or auto-updates for more than 5 seconds must have
a mechanism to pause, stop, or hide it:

```html
<section aria-label="Featured products" aria-roledescription="carousel">
  <button id="carousel-pause" aria-label="Pause carousel" type="button">⏸ Pause</button>

  <div id="carousel-slides" aria-live="off">
    <div role="group" aria-label="Slide 1 of 3" aria-roledescription="slide">
      <!-- slide content -->
    </div>
  </div>
</section>
```

```javascript
let isPlaying = true;
let intervalId = setInterval(advanceSlide, 5000);
const pauseBtn = document.getElementById('carousel-pause');
const slides = document.getElementById('carousel-slides');

pauseBtn.addEventListener('click', () => {
  if (isPlaying) {
    clearInterval(intervalId);
    isPlaying = false;
    pauseBtn.textContent = '▶ Play';
    pauseBtn.setAttribute('aria-label', 'Play carousel');
    slides.setAttribute('aria-live', 'off');
  } else {
    intervalId = setInterval(advanceSlide, 5000);
    isPlaying = true;
    pauseBtn.textContent = '⏸ Pause';
    pauseBtn.setAttribute('aria-label', 'Pause carousel');
  }
});

// Also pause when keyboard focus enters the carousel
const carousel = document.querySelector('[aria-roledescription="carousel"]');
carousel.addEventListener('focusin',  () => clearInterval(intervalId));
carousel.addEventListener('focusout', () => { if (isPlaying) intervalId = setInterval(advanceSlide, 5000); });
carousel.addEventListener('mouseenter', () => clearInterval(intervalId));
carousel.addEventListener('mouseleave', () => { if (isPlaying) intervalId = setInterval(advanceSlide, 5000); });
```

### Step 2: Warn and allow extension before session timeout (2.2.1 + 2.2.6)

If a session timeout will cause data loss, warn the user at least 20 seconds before expiry
and offer extension:

```html
<!-- Warning dialog — present in DOM at page load, shown on timeout warning -->
<div id="timeout-dialog" role="alertdialog" aria-modal="true"
     aria-labelledby="timeout-title" aria-describedby="timeout-desc" hidden>
  <h2 id="timeout-title">Session expiring soon</h2>
  <p id="timeout-desc">
    Your session will expire in <span id="countdown">2:00</span>.
    Unsaved changes will be lost.
  </p>
  <button id="extend-btn" onclick="extendSession()">Extend session</button>
  <button onclick="saveAndLogOut()">Save and log out</button>
</div>

<!-- Screen reader announcement region -->
<div id="timeout-live" role="status" aria-live="polite" aria-atomic="true"></div>
```

```javascript
const SESSION_MS  = 20 * 60 * 1000; // 20 minutes
const WARNING_MS  = 2  * 60 * 1000; // warn 2 minutes before

let sessionTimer  = setTimeout(expireSession,     SESSION_MS);
let warningTimer  = setTimeout(showWarning,        SESSION_MS - WARNING_MS);

function showWarning() {
  document.getElementById('timeout-dialog').removeAttribute('hidden');
  document.getElementById('extend-btn').focus();
  document.getElementById('timeout-live').textContent =
    'Session expires in 2 minutes. Choose Extend session to continue.';
  startCountdown();
}

function extendSession() {
  clearTimeout(sessionTimer);
  clearTimeout(warningTimer);
  document.getElementById('timeout-dialog').setAttribute('hidden', '');
  document.getElementById('timeout-live').textContent = 'Session extended by 20 minutes.';
  sessionTimer = setTimeout(expireSession, SESSION_MS);
  warningTimer = setTimeout(showWarning,   SESSION_MS - WARNING_MS);
}
```

### Step 3: Handle auto-refreshing data content (2.2.2)

Auto-updating tables or dashboards must offer a pause/disable control:

```html
<div class="live-data-header">
  <h2>Live stock prices</h2>
  <label>
    <input type="checkbox" id="auto-refresh" checked>
    Auto-refresh every 30 seconds
  </label>
</div>
<div role="status" aria-live="polite" id="refresh-status"></div>
<table id="stock-table"><!-- ... --></table>
```

```javascript
let refreshInterval = setInterval(refreshData, 30000);

document.getElementById('auto-refresh').addEventListener('change', (e) => {
  if (!e.target.checked) {
    clearInterval(refreshInterval);
    document.getElementById('refresh-status').textContent = 'Auto-refresh paused.';
  } else {
    refreshData();
    refreshInterval = setInterval(refreshData, 30000);
    document.getElementById('refresh-status').textContent = 'Auto-refresh enabled.';
  }
});
```

## When NOT to Use

- **Emergency alerts** — WCAG 2.2.2 exempts moving content that conveys an emergency. Don't add
  pause controls to evacuation notices.
- **Animations under 5 seconds** — content that blinks or moves for 5 seconds or less is exempt
  from 2.2.2. Use `prefers-reduced-motion` instead (`apply-reduced-motion`).
- **Real-time event countdowns** — live auction timers and broadcast clocks are exempt from 2.2.1
  (the time limit is essential). Still follow 2.2.6: warn before data loss.

## Common Mistakes

**Carousel pauses visually but `aria-live` keeps announcing.** When paused, set `aria-live="off"`
on the slides container so screen readers stop announcing slide changes.

**Timeout dialog doesn't move focus.** The timeout warning is a dialog — move focus to the "Extend
session" button when it appears (`apply-focus-management`). Restore focus to the user's previous
element when dismissed.

**`setInterval` drift in countdown display.** `setInterval` drifts over long sessions. Track
remaining time using `Date.now()` deltas rather than decrementing a counter.
