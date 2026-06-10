---
name: apply-responsive-text
description: Use when styling text on any web page — to ensure content remains readable and functional when users resize text to 200%, zoom to 400%, or override text spacing with their own stylesheets.
source: W3C WCAG 2.1 SC 1.4.4 Resize Text (Level AA), SC 1.4.10 Reflow (Level AA), SC 1.4.12 Text Spacing (Level AA); MDN responsive typography guide; UK Government Design System typography
tags: [accessibility, wcag, a11y, responsive-text, reflow, text-spacing, zoom, developer]
related: [design-accessibility-standards, apply-accessible-color, audit-accessibility]
---

# Implement Responsive Text

Use relative units and flexible layouts so text remains readable when users increase font size, zoom to 400%, or override text spacing — without horizontal scrolling or lost content.

## Why This Is Best Practice

**Adopted by:** WCAG 2.1 SCs 1.4.4, 1.4.10, and 1.4.12 are all Level AA — required by
Section 508 (US), EU EN 301 549, and UK PSBAR 2018. The UK Government Design System,
GOV.AU Design System, and US Web Design System all mandate responsive typography as a
core accessibility requirement. Apple HIG and Material Design both specify relative font
units as baseline for accessible typography.
**Impact:** 13% of people over 65 use browser zoom as their primary accessibility tool
(Nielsen Norman Group, 2023). Low-vision users who zoom to 400% will hit horizontal
scrollbars on fixed-width layouts — requiring them to scroll two dimensions, which is a
WCAG 2.1 AA failure. WCAG 1.4.12 was added after research showed dyslexia software
overriding text spacing consistently broke layouts that used fixed line heights.
**Why best:** `px` units in font-size and line-height override user browser font
preferences and break zoom/reflow. `rem`/`em` units inherit from and respect user
browser settings. Responsive text is a CSS architecture decision that must be made at
the foundation — retrofitting is expensive.

Sources: W3C WCAG 2.1 SC 1.4.4, 1.4.10, 1.4.12 (2018); Nielsen Norman Group Zoom
study (2023); UK Government Design System typography guidance

## Steps

### Step 1: Use `rem` for font sizes — never `px` alone

```css
/* Wrong — fixed px overrides user's browser font preference */
body { font-size: 16px; }
h1   { font-size: 32px; }
p    { font-size: 14px; }

/* Right — rem respects browser default (typically 16px) */
:root { font-size: 100%; }   /* inherits browser default, not override */
body  { font-size: 1rem; }   /* = browser default */
h1    { font-size: 2rem; }   /* = 32px at default, scales with user preference */
p     { font-size: 0.875rem; }
```

Setting `html { font-size: 62.5%; }` (10px base) undermines user preferences — avoid.
Use `font-size: 100%` on `:root` and `rem` everywhere else.

### Step 2: Use relative units for line height and spacing

```css
/* Wrong — px line-height breaks WCAG 1.4.12 text spacing override */
p {
  font-size: 1rem;
  line-height: 24px;   /* fixed — breaks when font-size changes */
}

/* Right — unitless or em line-height scales with font size */
p {
  font-size: 1rem;
  line-height: 1.5;    /* 1.5× the font size — scales automatically */
}

/* WCAG 1.4.12 minimum text spacing (must not break layout when applied): */
/* line-height: ≥ 1.5× font-size */
/* letter-spacing: ≥ 0.12em */
/* word-spacing: ≥ 0.16em */
/* paragraph spacing: ≥ 2× font-size */
```

### Step 3: Test WCAG 1.4.12 by applying the text spacing bookmarklet

Inject this CSS and verify no content is lost or clipped:

```css
/* Text Spacing Test — apply as bookmarklet or browser devtools */
* {
  line-height: 1.5 !important;
  letter-spacing: 0.12em !important;
  word-spacing: 0.16em !important;
}
p { margin-bottom: 2em !important; }
```

Everything that was visible must still be visible. Truncated text (`overflow: hidden`
on a fixed-height container) is the most common failure.

### Step 4: Make layouts reflow at 320px viewport (WCAG 1.4.10)

400% zoom on a 1280px desktop = a 320px viewport. Content must reflow to a single
column without horizontal scrolling:

```css
/* Wrong — fixed widths cause overflow at 320px */
.container { width: 800px; }
.sidebar   { width: 300px; float: left; }

/* Right — fluid layouts with max-width */
.container { width: 100%; max-width: 1200px; }

/* Right — CSS Grid/Flexbox that wraps naturally */
.layout {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}
.main    { flex: 1 1 300px; }
.sidebar { flex: 1 1 200px; }
```

Test with: Chrome DevTools → device toolbar → set width to 320px. Or zoom browser to
400% and verify no horizontal scrollbar appears.

Exception: content that requires 2D layout to function (data tables, maps, complex
diagrams) is exempt from WCAG 1.4.10 — but must still be scrollable.

### Step 5: Verify text resize at 200% browser zoom (WCAG 1.4.4)

```
1. Set browser font size to maximum (Chrome: Settings → Appearance → Font Size → Very Large)
2. OR zoom to 200% (Ctrl/Cmd +)
3. Verify: all text is still readable
4. Verify: no text is clipped, truncated, or overlapping
5. Verify: all interactive elements remain operable
```

Common failures:
- Fixed-height containers clip overflow text: use `min-height` not `height`
- Tooltips and popups overflow viewport: use responsive positioning
- Navigation wraps to overlap content: test with overflow visible

## When NOT to Use

- **Canvas-rendered text** — text drawn on `<canvas>` is not affected by CSS zoom. If canvas text must be accessible, provide a text alternative outside the canvas.
- **Data tables** — WCAG 1.4.10 exempts content that requires 2D navigation (tables). Horizontal scrolling is acceptable for data tables; wrap in a scrollable container with `role="region"` and `aria-label`.

## Common Mistakes

**`overflow: hidden` on fixed-height containers with text.** When a user increases text spacing or font size, text overflows the container and is clipped. Use `min-height` instead of `height`, or `overflow: auto`.

**`px` in media query breakpoints.** Browser zoom affects viewport width — a 320px media query breakpoint in `px` doesn't adjust with zoom. Use `em`-based breakpoints: `@media (max-width: 20em)`.

**Testing only at the default browser font size.** Increase the browser's default font size in settings (not just zoom) to catch layout breaks that only appear when the base `rem` value increases.
