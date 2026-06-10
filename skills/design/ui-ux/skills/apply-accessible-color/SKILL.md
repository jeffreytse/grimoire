---
name: apply-accessible-color
description: Use when choosing colors for text, UI components, or data visualizations — to ensure sufficient contrast ratios and that color is never the sole means of conveying information.
source: W3C WCAG 2.1 SC 1.4.1 Use of Color (Level A), SC 1.4.3 Contrast Minimum (Level AA), SC 1.4.11 Non-text Contrast (Level AA); WebAIM Contrast Checker; Deque axe-core color rules
tags: [accessibility, wcag, a11y, color-contrast, color-blindness, developer, design, visual-accessibility]
related: [design-accessibility-standards, audit-accessibility, apply-text-alternatives]
---

# Implement Accessible Color

Ensure all text and UI components meet WCAG contrast ratios and that no information is conveyed by color alone.

## Why This Is Best Practice

**Adopted by:** WCAG 2.1 SC 1.4.3 (Level AA) is required by Section 508, EU EN 301
549, UK PSBAR 2018, and all major platform guidelines (Apple HIG, Material Design,
Microsoft Fluent). WebAIM Million 2024 found low contrast text is the #1 WCAG failure
on 81% of top homepages — the single most widespread accessibility violation on the web.
**Impact:** 8% of men and 0.5% of women have color vision deficiency (CVD); low
contrast affects all users in bright sunlight or on low-quality screens. Google's
Material Design team documented that meeting AA contrast ratios improved readability
scores across all user groups, not just users with disabilities.
**Why best:** Color-only information (red = error, green = success) is invisible to
colorblind users and screen reader users. Contrast ratios are mathematically verifiable
— the only objective, testable way to ensure text is readable across disability types
and display conditions.

Sources: W3C WCAG 2.1 SC 1.4.1, 1.4.3, 1.4.11 (2018); WebAIM Million 2024;
Colour Blind Awareness (colourblindawareness.org)

## Steps

### Step 1: Apply the correct contrast ratio for each text size

| Text type | Minimum ratio (AA) | Enhanced ratio (AAA) |
|-----------|--------------------|----------------------|
| Normal text (< 18pt / < 14pt bold) | **4.5:1** | 7:1 |
| Large text (≥ 18pt OR ≥ 14pt bold) | **3:1** | 4.5:1 |
| UI components, icons, focus indicators | **3:1** | — |
| Decorative text (no information) | none | — |
| Disabled controls | none | — |

Check ratios with: WebAIM Contrast Checker, Figma Contrast plugin, or browser DevTools
Accessibility panel.

### Step 2: Test contrast in code, not just design mockups

```css
/* Wrong — visually looks fine at full opacity, fails at 60% */
.helper-text {
  color: rgba(0, 0, 0, 0.6);   /* #666 on white = 5.74:1 ✅ but context-dependent */
}

/* Right — use design tokens that encode the verified contrast pair */
.helper-text {
  color: var(--color-text-secondary);  /* token defined as #595959 on #fff = 7.0:1 ✅ */
}
```

Test every `color`/`background-color` combination where text appears — including:
hover states, focus states, disabled states (exempt), placeholder text (not exempt).

### Step 3: Never use color as the sole means of conveying information

```html
<!-- Wrong — only color distinguishes error from success -->
<p style="color: red">Your password is too short</p>
<p style="color: green">Password saved successfully</p>

<!-- Right — pair color with text label and icon -->
<p class="error">
  <span aria-hidden="true">✗</span>
  <strong>Error:</strong> Your password is too short
</p>
<p class="success">
  <span aria-hidden="true">✓</span>
  <strong>Success:</strong> Password saved
</p>
```

Same principle applies to: charts and graphs (add patterns or direct labels), form
validation (add icons + text), link identification (add underline, not just color),
required field marking (add asterisk + legend, not just red color).

### Step 4: Check non-text contrast for UI components

WCAG 2.1 SC 1.4.11 (AA) requires 3:1 contrast for:
- Input borders against background
- Focus indicators against adjacent background
- Chart elements (lines, bars) against chart background
- Icon buttons (icon against button background)

```css
/* Wrong — light gray border on white background = 1.6:1 */
input { border: 1px solid #ccc; }

/* Right — 3:1 minimum */
input { border: 1px solid #767676; }  /* #767676 on #fff = 4.54:1 ✅ */
```

### Step 5: Test with a color blindness simulator

Use browser extensions (Colorblindly, Stark) or Figma's Color Blind plugin to simulate:
- Deuteranopia (red-green, most common)
- Protanopia (red-green)
- Tritanopia (blue-yellow)
- Achromatopsia (no color)

If the interface relies on color to distinguish items (e.g., chart lines), add patterns,
shapes, or direct labels that work in grayscale.

## When NOT to Use

- **Disabled UI controls** — WCAG explicitly exempts disabled (non-interactive) controls from contrast requirements.
- **Decorative text** (logotype text, purely ornamental) — no contrast requirement.
- **Placeholder text** — technically exempt at AA, but best practice is to meet 4.5:1 anyway as placeholder disappears on input.

## Common Mistakes

**Testing with mockup colors, not rendered colors.** Transparency, overlays, and CSS gradients change the effective contrast. Always measure the computed color in the browser DevTools.

**Passing contrast for default state but failing for hover/focus.** Every interactive state must independently meet the ratio. A button with passing text contrast but a low-contrast hover state fails SC 1.4.3.

**Exempting text over images.** Text rendered over a background image has no guaranteed contrast. Either avoid text over images, use a scrim/overlay with sufficient contrast, or use a solid text block.
