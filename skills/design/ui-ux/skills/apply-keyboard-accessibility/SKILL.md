---
name: apply-keyboard-accessibility
description: Use when building interactive UI components — to ensure every action achievable with a mouse is also achievable with a keyboard alone, with no traps or dead ends.
source: W3C WCAG 2.1 SC 2.1.1 Keyboard (Level A), SC 2.1.2 No Keyboard Trap (Level A); WAI-ARIA Authoring Practices 1.2 (keyboard interaction patterns); MDN Web Docs keyboard accessibility guide
tags: [accessibility, wcag, a11y, keyboard-navigation, tab-order, focus, developer, interactive-components]
related: [apply-focus-management, apply-skip-navigation, design-accessibility-standards]
---

# Implement Keyboard Accessibility

Ensure every interactive element is reachable and operable by keyboard — Tab to reach it, Enter/Space to activate it, Escape to dismiss it, and arrow keys for internal navigation.

## Why This Is Best Practice

**Adopted by:** WCAG 2.1 SC 2.1.1 and 2.1.2 are Level A — the absolute baseline
required by every accessibility law worldwide (Section 508, EU EN 301 549, UK PSBAR
2018, Australian DDA). Motor disability affects 2.6% of the US population; keyboard-only
operation is essential for users who cannot use a mouse (tremors, paralysis, switch
access). Power users also rely on keyboard for speed.
**Impact:** WebAIM survey 2024 found keyboard navigation is the primary input method for
21% of screen reader users and an essential fallback for many more. Custom widgets
(datepickers, dropdowns, sliders) fail keyboard access in the majority of implementations
without deliberate effort — the WAI-ARIA Authoring Practices exist specifically to
document correct keyboard patterns.
**Why best:** Native HTML elements (`<button>`, `<input>`, `<select>`, `<a href>`) are
keyboard accessible by default. Every custom widget built on `<div>` or `<span>`
requires explicit keyboard implementation. The cost of retrofitting is 3–5× the cost
of building it correctly.

Sources: W3C WCAG 2.1 SC 2.1.1, 2.1.2 (2018); WAI-ARIA Authoring Practices 1.2;
WebAIM Screen Reader User Survey #10 (2024)

## Steps

### Step 1: Use native elements — they're keyboard accessible by default

```html
<!-- Wrong — div requires manual keyboard implementation -->
<div class="btn" onclick="submit()">Submit</div>

<!-- Right — button is keyboard accessible, focusable, and announced correctly -->
<button type="submit">Submit</button>
```

Use: `<button>` for actions, `<a href>` for navigation, `<input>` for data entry,
`<select>` for dropdowns of options, `<details>` for disclosure widgets.
Build custom widgets only when native elements cannot meet the design requirement.

### Step 2: Implement the correct keyboard pattern for each widget type

Per WAI-ARIA Authoring Practices:

| Widget | Tab behavior | Internal navigation |
|--------|-------------|---------------------|
| Button | Focusable | Enter/Space = activate |
| Link | Focusable | Enter = follow |
| Checkbox | Focusable | Space = toggle |
| Radio group | One tab stop | Arrow keys = change selection |
| Listbox/select | One tab stop | Arrow keys = navigate, Enter = select |
| Menu | One tab stop | Arrow keys = navigate, Enter/Space = select, Escape = close |
| Dialog | Trap focus inside | Escape = close, Tab cycles within |
| Slider | Focusable | Arrow keys = adjust value |
| Tabs | One tab stop per tablist | Arrow keys = change tab |

```javascript
// Example: arrow key navigation for a custom listbox
listbox.addEventListener('keydown', (e) => {
  switch (e.key) {
    case 'ArrowDown':
      e.preventDefault();
      focusNextOption();
      break;
    case 'ArrowUp':
      e.preventDefault();
      focusPrevOption();
      break;
    case 'Enter':
    case ' ':
      e.preventDefault();
      selectFocusedOption();
      break;
    case 'Escape':
      closeListbox();
      break;
  }
});
```

### Step 3: Never use `tabindex > 0`

```html
<!-- Wrong — positive tabindex creates an unpredictable tab order -->
<div tabindex="3">Third in tab order? Probably not.</div>

<!-- Right — tabindex="0" adds to natural DOM order -->
<div role="button" tabindex="0" onclick="...">Custom button</div>

<!-- Right — tabindex="-1" makes focusable programmatically, not via Tab -->
<div id="modal-content" tabindex="-1">Dialog content</div>
```

`tabindex > 0` creates a separate tab sequence that runs before the natural order —
producing a confusing jump that most users do not expect.

### Step 4: Ensure no keyboard trap exists

```javascript
// Wrong — modal that doesn't handle Tab creates a trap
// (keyboard focus leaves the modal, never to return)

// Right — confine Tab to modal contents
modal.addEventListener('keydown', (e) => {
  if (e.key !== 'Tab') return;

  const focusable = modal.querySelectorAll(
    'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
  );
  const first = focusable[0];
  const last = focusable[focusable.length - 1];

  if (e.shiftKey && document.activeElement === first) {
    e.preventDefault();
    last.focus();
  } else if (!e.shiftKey && document.activeElement === last) {
    e.preventDefault();
    first.focus();
  }
});
```

See `apply-focus-management` for the full modal focus trap pattern.

### Step 5: Test keyboard-only before every release

Manual keyboard test checklist:
- [ ] Tab visits every interactive element in logical DOM order
- [ ] Shift+Tab moves backward through the same order
- [ ] Enter/Space activates every button and link
- [ ] Escape closes modals, dropdowns, and tooltips
- [ ] Arrow keys navigate within composite widgets (menus, sliders, tabs)
- [ ] No focus is lost after a dynamic update (AJAX load, route change)

## When NOT to Use

- **Non-interactive elements** — purely presentational elements need no keyboard support. Don't add `tabindex="0"` to static text, images without actions, or decorative elements.

## Common Mistakes

**`onclick` without `onkeydown`.** A `<div>` with `onclick` responds to mouse click but not Enter or Space. Add `role="button"`, `tabindex="0"`, and keyboard event handling — or switch to `<button>`.

**Arrow keys that also scroll the page.** Arrow key handlers must call `e.preventDefault()` to prevent the page from scrolling while navigating a widget.

**Tab order that doesn't match visual order.** If the DOM order doesn't match visual layout (due to CSS flexbox/grid reordering), Tab order diverges from visual order and confuses sighted keyboard users. Fix the DOM order; use CSS only for visual reordering.
