---
name: apply-sensory-independence
description: Use when writing UI instructions or labels, and when building mobile or responsive interfaces — to ensure content is understandable without relying solely on shape, color, position, or orientation, and that display orientation is not unnecessarily locked.
source: W3C WCAG 2.1 SC 1.3.3 Sensory Characteristics (Level A), SC 1.3.4 Orientation (Level AA); WAI-ARIA Authoring Practices; UK GOV.UK Design System content guidelines
tags: [accessibility, wcag, a11y, sensory-characteristics, orientation, shape-independence, visual-independence, developer]
related: [apply-accessible-color, write-semantic-html-structure, design-accessibility-standards]
---

# Implement Sensory Independence

Write instructions that identify components by name or function — not only by shape, color,
size, or position — and avoid locking screen orientation unless the content genuinely requires it.

## Why This Is Best Practice

**Adopted by:** WCAG 2.1 SC 1.3.3 (Level A) and SC 1.3.4 (Level AA) are required by Section 508
(US), EU EN 301 549, and UK PSBAR 2018. GOV.UK Design System content guidelines explicitly prohibit
directional references ("the button on the left") as a content standard. Apple HIG and Android
accessibility guidelines both prohibit orientation locks except for inherently 2D content.
**Impact:** Screen reader users cannot perceive shape, color, or visual position — instructions
like "click the blue button" or "see the panel on the right" are meaningless when the page is
read linearly. For orientation: 1 in 13 users with mobility impairments mounts their phone in
a fixed orientation (portrait or landscape) because rotating the device is physically difficult
(WebAIM 2023). An app that forces landscape locks these users out entirely.
**Why best:** The alternative — relying on visual sensory properties for instructions — works only
for sighted users. Adding a name or label ("click **Submit**" instead of "click the blue button")
costs nothing and makes the instruction work for everyone. SC 1.3.4 orientation independence is
one CSS and JS line — the cost of compliance is trivially low.

Sources: W3C WCAG 2.1 SCs 1.3.3, 1.3.4 (2018); GOV.UK Design System content guidelines;
WebAIM Motor Disabilities Survey (2023); Apple HIG Adaptivity and Layout

## Steps

### Step 1: Identify UI components by name or function — not sensory properties (1.3.3)

Instructions must not rely solely on shape, color, size, visual location, or sound:

```html
<!-- Wrong — relies on color and position only -->
<p>Click the green button on the right to continue.</p>
<button style="color: green; float: right;">→</button>

<!-- Right — identifies by label/function -->
<p>Click <strong>Continue</strong> to proceed to the next step.</p>
<button>Continue</button>
```

```html
<!-- Wrong — relies on position -->
<p>Use the panel on the left to filter results.</p>

<!-- Right — identifies by name -->
<p>Use the <strong>Filters</strong> panel to narrow results.</p>
<aside aria-label="Filters"><!-- filter controls --></aside>
```

```html
<!-- Wrong — relies on shape -->
<p>Click the round icon to add to favourites.</p>

<!-- Right — identifies by label -->
<p>Click <strong>Add to favourites</strong> to save this item.</p>
<button aria-label="Add to favourites">♡</button>
```

Allowed in combination: "Click the **Submit** button (the green button on the right)" — using
a sensory property as an *additional* cue alongside a name is fine. The name alone must be
sufficient.

### Step 2: Pair icon-only indicators with accessible labels (1.3.3)

Status indicators and icons that convey meaning through shape or color alone need a text
equivalent:

```html
<!-- Wrong — status conveyed by icon shape only -->
<span class="status-icon">✓</span> Order confirmed
<span class="status-icon">✗</span> Payment failed

<!-- Right — icon is supplemental; text carries the meaning -->
<span aria-hidden="true">✓</span>
<span class="visually-hidden">Success:</span> Order confirmed

<span aria-hidden="true">✗</span>
<span class="visually-hidden">Error:</span> Payment failed
```

```css
/* Visually hidden but available to screen readers */
.visually-hidden {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}
```

### Step 3: Don't lock screen orientation (1.3.4)

Content must work in both portrait and landscape — don't restrict orientation in HTML or CSS:

```html
<!-- Wrong — forces landscape only -->
<meta name="viewport" content="width=device-width, initial-scale=1, orientation=landscape">
```

```css
/* Wrong — hides content in portrait mode */
@media (orientation: portrait) {
  body { display: none; }
}

/* Wrong — blocks rotation with CSS */
html {
  transform: rotate(90deg);
  overflow: hidden;
}
```

```javascript
// Wrong — locks to landscape programmatically
screen.orientation.lock('landscape');

// Right — only lock when the content genuinely requires it (see exceptions)
// Most apps: do nothing. Let the OS handle orientation.
```

If you use `screen.orientation.lock()` for genuinely essential content, provide an
unlock path:

```javascript
// If orientation lock is truly necessary — offer a way out
async function requestLandscape() {
  try {
    await screen.orientation.lock('landscape');
  } catch {
    // Lock failed (e.g. desktop browser) — show instruction instead
    showOrientationPrompt('Please rotate your device to landscape for best experience');
  }
}

document.getElementById('unlock-orientation').addEventListener('click', () => {
  screen.orientation.unlock();
});
```

### Step 4: Audit instructions in content and help text

Run a find-and-replace audit for sensory-only references in UI copy:

| Failing pattern | Fix |
|----------------|-----|
| "the button on the left/right/top/bottom" | "the **[Button Label]** button" |
| "the blue/red/green [element]" | "the **[Element Name]**" |
| "the large/small [element]" | "the **[Element Name]**" |
| "the round/square/circular icon" | "the **[Icon Name]** icon" |
| "click the ✓ mark" | "click **Confirm**" |
| "see the box on the right" | "see the **[Section Name]** section" |

## When NOT to Use

- **Orientation lock for inherently 2D content** — piano keyboard apps, landscape-only games,
  document editors that require wide format, some medical imaging viewers. SC 1.3.4 exempts
  content where a specific orientation is essential to function.
- **Supplemental sensory cues** — adding color or position as additional context alongside
  a name is fine and often helpful for sighted users. The rule requires sensory properties
  not be the *only* means of identification.

## Common Mistakes

**"Fill in the fields marked with a red asterisk."** Red color alone is a sensory characteristic.
Accompany with: "Fields marked with a red asterisk (*) are required" — the asterisk character
itself is non-color; with the text explanation it satisfies both 1.3.3 and 1.4.1.

**Rotation message instead of responsive layout.** Showing "Please rotate your device" in
portrait is a 1.3.4 failure — it restricts functionality to one orientation. Implement
a responsive layout that works in both orientations instead.

**`@media (orientation: portrait)` that hides navigation.** Using CSS to hide critical
navigation in portrait mode restricts content to landscape. If navigation must reflow,
show it differently (hamburger menu) — don't hide it entirely.
