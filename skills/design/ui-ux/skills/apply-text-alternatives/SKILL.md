---
name: apply-text-alternatives
description: Use when adding images, icons, SVGs, canvas elements, or any non-text content to a page — to ensure users who cannot see the content receive an equivalent text description.
source: W3C WCAG 2.1 SC 1.1.1 Non-text Content (Level A); WAI-ARIA Authoring Practices 1.2; WebAIM Alt Text guide
tags: [accessibility, wcag, a11y, alt-text, images, aria-label, developer, screen-reader]
related: [design-accessibility-standards, audit-accessibility, apply-aria-roles]
---

# Implement Text Alternatives

Provide a text equivalent for every non-text element so screen reader users receive the same information as sighted users.

## Why This Is Best Practice

**Adopted by:** WCAG 2.1 SC 1.1.1 is Level A — the lowest conformance level, meaning
it is the baseline required by every accessibility standard worldwide: Section 508 (US),
EN 301 549 (EU), PSBAR 2018 (UK), and the Australian DDA. Every major tech company's
accessibility guidelines (Apple, Google, Microsoft) begin here.
**Impact:** WebAIM Million 2024 found missing alternative text on images is the #2 most
common WCAG failure (35.2% of homepages), making it the single most impactful quick
fix. Screen readers announce image filenames or "image" when alt text is absent —
producing a meaningless or disruptive experience for blind users.
**Why best:** The alternative — relying on surrounding text to convey image meaning — is
insufficient because screen readers read images in DOM order, which may not match the
surrounding context. Explicit alt text guarantees the description is attached to the
content, not assumed from context.

Sources: W3C WCAG 2.1 SC 1.1.1 (2018); WebAIM Million 2024; Google Lighthouse
accessibility documentation

## Steps

### Step 1: Classify the image — informative, decorative, or functional

| Type | Correct treatment |
|------|------------------|
| Informative — conveys meaning | Descriptive `alt="..."` text |
| Decorative — visual only, no meaning | `alt=""` (empty string — screen reader skips it) |
| Functional — a link or button | `alt` describes the destination or action, not appearance |
| Complex — chart, graph, diagram | Short `alt` + longer description in `aria-describedby` or adjacent text |

### Step 2: Write informative alt text

```html
<!-- Wrong — filename, "image of", or no context -->
<img src="hero.jpg" alt="hero.jpg">
<img src="chart.png" alt="image of a chart">

<!-- Right — describes the meaning, not the appearance -->
<img src="revenue-q3.png" alt="Q3 revenue: $4.2M, up 18% from Q2">
<img src="warning.svg" alt="Warning: this action cannot be undone">
```

Rules for writing good alt text:
- Describe the purpose, not the visual appearance ("chart showing upward trend" not "blue bar chart")
- Keep it under 150 characters; use `aria-describedby` for longer descriptions
- Do not start with "image of" or "picture of" — screen readers already announce "image"
- For functional images (linked), describe the destination: `alt="View order #4821"`

### Step 3: Mark decorative images with empty alt

```html
<!-- Empty alt — screen reader skips entirely -->
<img src="divider-line.svg" alt="">

<!-- CSS background images are already invisible to screen readers -->
<div style="background-image: url(decorative-bg.jpg)" aria-hidden="true"></div>
```

Never omit `alt` — a missing `alt` causes the screen reader to announce the filename.
Empty `alt=""` is the correct decoration signal.

### Step 4: Label icon-only buttons and SVGs

```html
<!-- Wrong — screen reader announces "button" with no name -->
<button><svg>...</svg></button>

<!-- Right — aria-label provides the accessible name -->
<button aria-label="Close dialog">
  <svg aria-hidden="true" focusable="false">...</svg>
</button>

<!-- Right — title inside SVG (also sets tooltip) -->
<svg role="img" aria-labelledby="chart-title">
  <title id="chart-title">Monthly revenue for Q3 2026</title>
  ...
</svg>
```

### Step 5: Handle complex images with extended descriptions

```html
<figure>
  <img src="org-chart.png"
       alt="Company organizational chart"
       aria-describedby="org-chart-desc">
  <figcaption id="org-chart-desc">
    CEO reports to Board. Three VPs (Engineering, Sales, Finance) report to CEO.
    Each VP has 4–6 direct reports.
  </figcaption>
</figure>
```

## When NOT to Use

- **CSS background images used purely for decoration** — they are already inaccessible to AT by default; no action needed.
- **Redundant images** — an image of a "Submit" button next to a visible "Submit" text label should use `alt=""` to avoid announcing "Submit" twice.

## Common Mistakes

**Omitting `alt` entirely.** A missing `alt` attribute causes screen readers to announce the image filename (e.g., "img_2847_final_v3.jpg"). Always include `alt` — empty if decorative.

**Using `alt` to stuff keywords.** Alt text is not an SEO field. Keyword-stuffed alt text ("cheap flights cheap hotels cheap car rental image") is disruptive to screen reader users.

**Setting `aria-hidden="true"` on informative images.** This hides the image from all AT, including screen readers. Only use `aria-hidden` on decorative content.
