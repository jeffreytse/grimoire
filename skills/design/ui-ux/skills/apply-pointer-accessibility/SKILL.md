---
name: apply-pointer-accessibility
description: Use when building touch interfaces, drag-and-drop interactions, gesture-based controls, or any feature activated on mouse or touch events — to ensure users with motor disabilities can operate all functionality with a single pointer without complex gestures.
source: W3C WCAG 2.1 SC 2.5.1 Pointer Gestures (Level A), SC 2.5.2 Pointer Cancellation (Level A), SC 2.5.3 Label in Name (Level A), SC 2.5.4 Motion Actuation (Level A); Apple HIG Touch Interaction guidelines; Google Material Design interaction guidelines
tags: [accessibility, wcag, a11y, pointer, gestures, motor-accessibility, speech-input, developer]
related: [apply-keyboard-accessibility, apply-hover-content, design-accessibility-standards]
---

# Implement Pointer Accessibility

Ensure all functionality works with a single pointer and a simple tap — by providing alternatives
to multi-point gestures, activating on pointer-up not pointer-down, matching visible labels to
accessible names, and offering alternatives to device motion.

## Why This Is Best Practice

**Adopted by:** WCAG 2.1 SCs 2.5.1–2.5.4 (Level A) are required by Section 508 (US), EU EN 301 549,
and UK PSBAR 2018. Apple's Human Interface Guidelines and Google Material Design both specify
pointer-up activation and single-touch alternatives as baseline interaction requirements for
all platforms.
**Impact:** 26 million people in the US have mobility impairments (CDC, 2023). Users with motor
disabilities may use a single switch, head pointer, or adaptive touch devices that cannot perform
multi-finger gestures or path-based swipes. Speech input users (Dragon NaturallySpeaking, iOS Voice
Control) activate controls by saying the visible label — if the visible label doesn't match the
accessible name, speech commands fail silently.
**Why best:** Browser native controls (buttons, links, select) satisfy 2.5.2 by default because
`click` fires on pointer-up. Custom components built on `mousedown`/`touchstart` fail it. Providing
single-tap alternatives to gestures doesn't remove the gesture — it adds an equivalent path for
users who need it.

Sources: W3C WCAG 2.1 SCs 2.5.1–2.5.4 (2018); CDC Disability Statistics (2023); Apple HIG
Touch Interaction; Google Material Design Interaction guidelines

## Steps

### Step 1: Provide single-pointer alternatives to multi-point gestures (2.5.1)

Every function requiring a multi-point gesture (pinch, two-finger swipe) or a path-based gesture
(draw a shape, swipe arc) must also work with a single tap or click:

```html
<!-- Map with pinch-to-zoom — add + / - buttons as single-pointer alternative -->
<div id="map" aria-label="Interactive map"></div>
<div class="map-controls" role="group" aria-label="Map zoom controls">
  <button id="zoom-in"  aria-label="Zoom in"  type="button">+</button>
  <button id="zoom-out" aria-label="Zoom out" type="button">−</button>
</div>
```

```javascript
document.getElementById('zoom-in').addEventListener('click',
  () => map.setZoom(map.getZoom() + 1));
document.getElementById('zoom-out').addEventListener('click',
  () => map.setZoom(map.getZoom() - 1));
// Keep the multi-point gesture — just ADD the single-pointer path
```

Swipeable carousels must have previous/next buttons. Drag-to-sort lists must have
up/down buttons or a keyboard sort mechanism (`apply-keyboard-accessibility`).

### Step 2: Activate on pointer-up, not pointer-down (2.5.2)

Use `click` (fires on up) — never `mousedown` or `touchstart` for actions:

```javascript
// Wrong — activates immediately; user cannot cancel by moving pointer away
button.addEventListener('mousedown', submitOrder);
button.addEventListener('touchstart', submitOrder);

// Right — fires on release; user cancels by moving pointer off element before release
button.addEventListener('click', submitOrder);
```

For drag operations where pointer-down is required to begin tracking, abort the action
when pointer is released outside a valid drop target:

```javascript
let dragging = false;
let dragSource = null;

draggable.addEventListener('pointerdown', (e) => {
  dragging = true;
  dragSource = e.currentTarget;
});

document.addEventListener('pointerup', (e) => {
  if (!dragging) return;
  const dropTarget = e.target.closest('[data-drop-zone]');
  if (!dropTarget) {
    cancelDrag(dragSource); // pointer released outside valid target — abort
  } else {
    completeDrop(dragSource, dropTarget);
  }
  dragging = false;
});
```

Workaround for 300ms tap delay (old iOS): use `touch-action: manipulation` in CSS instead
of `touchstart` — eliminates the delay without breaking pointer cancellation.

```css
button, a, [role="button"] {
  touch-action: manipulation; /* removes 300ms delay, keeps pointer-up activation */
}
```

### Step 3: Ensure visible label text is in the accessible name (2.5.3)

Speech input users activate buttons by saying the visible label text. The accessible name
must contain the visible text:

```html
<!-- Wrong — "Submit form data" doesn't contain visible text "Submit" -->
<button aria-label="Submit form data">Submit</button>

<!-- Right — accessible name contains visible text -->
<button aria-label="Submit form">Submit</button>
<!-- Or simply — use visible text directly -->
<button>Submit</button>

<!-- Wrong — aria-label doesn't match visible tooltip text -->
<button aria-label="Delete item" title="Remove from list">🗑</button>
<!-- User says "click Remove from list" — fails to match "Delete item" -->

<!-- Right — accessible name matches what user would say -->
<button aria-label="Remove from list">🗑</button>
```

Rule: if a button has visible text, `aria-label` must **contain** that visible text
(case-insensitive). Prefer no `aria-label` on buttons with visible text labels.

### Step 4: Provide alternatives to device motion (2.5.4)

If a feature is triggered by device motion (shake to undo, tilt to navigate), provide
an equivalent UI control and allow users to disable the motion trigger:

```javascript
// Feature triggered by shake — also provide a UI alternative
if (window.DeviceMotionEvent) {
  window.addEventListener('devicemotion', handleShake);
}
document.getElementById('undo-button').addEventListener('click', undoLastAction);

// Provide a setting to disable motion triggering
const motionToggle = document.getElementById('disable-motion-toggle');
motionToggle.addEventListener('change', (e) => {
  if (e.target.checked) {
    window.removeEventListener('devicemotion', handleShake);
  } else {
    window.addEventListener('devicemotion', handleShake);
  }
});

// Persist preference across sessions
motionToggle.checked = localStorage.getItem('motion-disabled') === 'true';
motionToggle.addEventListener('change', (e) =>
  localStorage.setItem('motion-disabled', e.target.checked));
```

## When NOT to Use

- **Essential path-based gestures** — if the path of movement is essential to the function
  (handwriting recognition, drawing app, signature pad), 2.5.1 does not require an alternative.
  The gesture IS the content.
- **Browser-native viewport behavior** — browser scroll and pinch-to-zoom at the viewport level
  are exempt. Don't disable them with `user-scalable=no` on the `<meta viewport>` tag.

## Common Mistakes

**`aria-label` overwrites visible text, breaking speech input.** `<button aria-label="Close dialog">✕ Cancel</button>` — speech user says "click Cancel" but accessible name is "Close dialog". Either remove `aria-label` and fix visible text, or set `aria-label="Cancel"` to match.

**Drag-and-drop with no keyboard alternative.** Drag-and-drop fails 2.5.1 and also WCAG 2.1.1
(keyboard). Provide a keyboard sort mechanism or move-to buttons alongside the drag UI.

**`touchstart` for tap delay optimization.** Some developers use `touchstart` over `click` for
~300ms latency on old iOS. This violates 2.5.2. Use `touch-action: manipulation` in CSS
instead — eliminates the delay without breaking pointer cancellation.
