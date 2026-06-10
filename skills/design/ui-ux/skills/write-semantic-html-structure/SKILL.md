---
name: write-semantic-html-structure
description: Use when building or reviewing a web page's HTML structure — to ensure heading hierarchy, landmark regions, and semantic elements convey the correct document outline to screen readers and assistive technology.
source: W3C WCAG 2.1 SC 1.3.1 Info and Relationships (Level A), SC 2.4.6 Headings and Labels (Level AA); W3C HTML5 semantic element specification; ARIA Landmarks specification
tags: [accessibility, wcag, a11y, semantic-html, headings, landmarks, screen-reader, developer]
related: [design-accessibility-standards, apply-aria-roles, apply-skip-navigation]
---

# Write Semantic HTML Structure

Use heading hierarchy and landmark regions to give the page a logical structure that screen readers can navigate.

## Why This Is Best Practice

**Adopted by:** WCAG 2.1 SC 1.3.1 (Level A) and SC 2.4.6 (Level AA) are required by
all major accessibility standards (Section 508, EN 301 549, PSBAR 2018). The W3C HTML5
spec's semantic elements (`<nav>`, `<main>`, `<article>`, `<aside>`) are supported in
all modern browsers since 2014. WebAIM screen reader user surveys consistently show
heading navigation as the #1 navigation technique — 67% of screen reader users use
headings to find content.
**Impact:** WebAIM Million 2024 found that empty heading levels and missing landmark
regions appear in 43% of analyzed pages. Screen reader users navigate by landmarks
(`R` key in NVDA) and headings (`H` key) to skip to desired content — without these,
every page requires linear reading from top to bottom.
**Why best:** `<div>` and `<span>` have no semantic value. A page built entirely with
divs forces screen reader users to listen to everything sequentially. Semantic elements
provide free navigation, document structure, and AT compatibility with no additional
ARIA required.

Sources: W3C WCAG 2.1 SC 1.3.1, 2.4.6 (2018); WebAIM Screen Reader User Survey #10
(2024); W3C HTML5 semantic element specification

## Steps

### Step 1: Use one `<h1>` per page — match the page title

```html
<!-- Wrong — multiple h1s, skipped levels -->
<h1>Company Name</h1>
<h1>Products</h1>    <!-- second h1 — creates ambiguous top-level -->
<h3>Widget A</h3>    <!-- skipped h2 — broken outline -->

<!-- Right — one h1, sequential hierarchy -->
<h1>Products</h1>
  <h2>Widgets</h2>
    <h3>Widget A</h3>
    <h3>Widget B</h3>
  <h2>Services</h2>
```

Heading levels convey document outline — never choose a heading level for visual size.
Use CSS to style headings; use heading elements only for structure.

### Step 2: Use HTML5 landmark elements — not ARIA roles on divs

```html
<!-- Wrong — all divs, no structure -->
<div class="header">...</div>
<div class="nav">...</div>
<div class="content">...</div>
<div class="sidebar">...</div>
<div class="footer">...</div>

<!-- Right — native semantic elements -->
<header>...</header>
<nav aria-label="Main navigation">...</nav>
<main>
  <article>...</article>
  <aside>...</aside>
</main>
<footer>...</footer>
```

Each landmark element maps to an implicit ARIA role: `<main>` = `role="main"`,
`<nav>` = `role="navigation"`, `<header>` = `role="banner"`, `<footer>` = `role="contentinfo"`.
Native elements are always preferred over `role=` on divs.

### Step 3: Label multiple instances of the same landmark

```html
<!-- Wrong — two nav landmarks with identical names -->
<nav>...</nav>
<nav>...</nav>

<!-- Right — aria-label distinguishes them -->
<nav aria-label="Main navigation">...</nav>
<nav aria-label="Breadcrumb">...</nav>
```

Screen readers list all landmarks — duplicate names create confusion.
Use `aria-label` or `aria-labelledby` to distinguish repeated landmark types.

### Step 4: Use semantic elements for lists, tables, and emphasis

```html
<!-- Wrong — visual list with no semantics -->
<div class="list">
  <div>• Item one</div>
  <div>• Item two</div>
</div>

<!-- Right — screen reader announces "list, 2 items" -->
<ul>
  <li>Item one</li>
  <li>Item two</li>
</ul>

<!-- Tables: always include <caption> and <th scope="col/row"> -->
<table>
  <caption>Q3 Sales by Region</caption>
  <thead>
    <tr><th scope="col">Region</th><th scope="col">Revenue</th></tr>
  </thead>
  <tbody>
    <tr><td>North</td><td>$1.2M</td></tr>
  </tbody>
</table>
```

### Step 5: Verify the document outline with a heading visualizer

Use the browser extension "HeadingsMap" or the axe DevTools outline view to verify:
- One `<h1>` at the page level
- No skipped heading levels (h1 → h3 with no h2 is invalid)
- Headings describe the content of their sections

## When NOT to Use

- **Presentational layout elements** — `<header>` inside `<article>` scopes to the article, not the page. Use landmark elements deliberately, not as generic containers.
- **Email HTML** — most email clients strip semantic elements. Use `role=` attributes instead for accessible email.

## Common Mistakes

**Styling headings visually rather than structurally.** Using `<h3>` because it looks like the right size — not because it's the third level in the hierarchy — breaks the outline. Style with CSS; use heading levels for structure only.

**Multiple `<main>` elements.** A page must have exactly one `<main>`. Multiple `<main>` elements are invalid and confuse AT.

**`<section>` and `<article>` without headings.** Both elements are expected to have a heading that names them. A `<section>` without a heading provides no navigation benefit and is treated as a generic container.
