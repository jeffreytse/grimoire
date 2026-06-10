---
name: apply-focus-management
description: Use when building modals, drawers, route changes in SPAs, dynamic content insertion, or any interaction that moves or removes focus — to ensure keyboard focus moves predictably and is never lost.
source: W3C WCAG 2.1 SC 2.4.3 Focus Order (Level A), SC 2.4.7 Focus Visible (Level AA), SC 2.4.11 Focus Appearance (Level AA, WCAG 2.2); WAI-ARIA Dialog Pattern; WCAG 2.2
tags: [accessibility, wcag, a11y, focus-management, modals, focus-trap, developer, keyboard-navigation]
related: [apply-keyboard-accessibility, apply-skip-navigation, design-accessibility-standards]
---

# Implement Focus Management

Move keyboard focus intentionally when UI changes — so focus is never lost, never trapped outside a dialog, and always lands in a logical position after dynamic updates.

## Why This Is Best Practice

**Adopted by:** WCAG 2.1 SC 2.4.3 (Level A) and 2.4.7 (Level AA) are required by all
major accessibility laws. The WAI-ARIA Dialog Pattern is the W3C specification for
modal focus management — adopted in every major component library (React Aria, Headless
UI, Radix UI, Angular CDK). WCAG 2.2 added SC 2.4.11 Focus Appearance (AA) requiring
focus indicators to be sufficiently visible.
**Impact:** Lost focus — where Tab focus disappears after an action — is the #1 modal
accessibility defect per Deque Systems research. When a modal opens and focus is not
moved inside it, a screen reader user continues interacting with the content behind
the modal, creating confusion and incorrect actions. When a modal closes and focus is
not restored, the user loses their position in the page.
**Why best:** Browsers do not automatically manage focus across dynamic UI changes. SPA
route changes, modals, drawers, and AJAX-injected content all move or remove DOM
elements without browser focus recovery. Manual focus management is the only mechanism
that ensures users retain orientation.

Sources: W3C WCAG 2.1 SC 2.4.3, 2.4.7 (2018); WCAG 2.2 SC 2.4.11 (2023);
WAI-ARIA Dialog Pattern; Deque Systems accessibility audit data

## Steps

### Step 1: Move focus into a modal when it opens

```javascript
function openModal(modal) {
  modal.removeAttribute('hidden');
  modal.removeAttribute('aria-hidden');

  // Move focus to modal container or first focusable element
  const firstFocusable = modal.querySelector(
    'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
  );
  (firstFocusable || modal).focus();

  // Store the element that triggered the modal — needed for restoration
  modal._triggerElement = document.activeElement;
}
```

Prefer focusing the first interactive element inside the modal. If the modal has
a heading that introduces the content, focus the heading (`tabindex="-1"`) first
to announce context before the first control.

### Step 2: Trap focus inside the modal while open

```javascript
function trapFocus(modal) {
  const focusableSelector =
    'button:not([disabled]), [href], input:not([disabled]), select:not([disabled]), ' +
    'textarea:not([disabled]), [tabindex]:not([tabindex="-1"])';

  modal.addEventListener('keydown', (e) => {
    if (e.key !== 'Tab') return;

    const focusable = [...modal.querySelectorAll(focusableSelector)];
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
}
```

### Step 3: Restore focus when modal closes

```javascript
function closeModal(modal) {
  modal.setAttribute('hidden', '');
  modal.setAttribute('aria-hidden', 'true');

  // Return focus to the element that opened the modal
  if (modal._triggerElement) {
    modal._triggerElement.focus();
    modal._triggerElement = null;
  }
}
```

If the trigger element is no longer in the DOM (e.g., deleted), move focus to the
nearest logical parent or to the `<main>` element.

### Step 4: Manage focus on SPA route changes

```javascript
// React example — move focus to page heading after route change
function RouteChangeHandler() {
  const location = useLocation();

  useEffect(() => {
    // Update document title
    document.title = getPageTitle(location.pathname);

    // Move focus to main heading (tabindex="-1" allows programmatic focus)
    const h1 = document.querySelector('main h1');
    if (h1) h1.focus();
  }, [location.pathname]);

  return null;
}
```

Without this, screen reader users hear nothing after a route change — the page
appears unchanged from their perspective.

### Step 5: Style focus indicators — never remove them

```css
/* Wrong — removes all focus styling */
:focus { outline: none; }
*:focus { outline: 0; }

/* Right — replace with custom visible indicator */
:focus-visible {
  outline: 3px solid #0066cc;
  outline-offset: 2px;
  border-radius: 2px;
}

/* WCAG 2.2 SC 2.4.11 minimum: 2px solid, 3:1 contrast against adjacent colors */
```

Use `:focus-visible` (not `:focus`) to show focus rings only for keyboard users,
not mouse clicks — this resolves the design/accessibility conflict.

## When NOT to Use

- **Non-interactive dynamic content** — adding a banner or notification that doesn't require user action does not need focus management. Use `role="status"` or `role="alert"` to announce it instead (see `apply-aria-roles`).

## Common Mistakes

**Moving focus before the element is visible.** Calling `.focus()` on a `display: none` element silently fails. Ensure the element is visible (`display` is not `none`, `visibility` is not `hidden`) before calling `.focus()`.

**Focus trap that blocks Escape.** The modal must close and restore focus on Escape key. A focus trap without an Escape handler permanently traps keyboard users.

**Focus ring removed globally.** `* { outline: none }` in a CSS reset removes keyboard accessibility for every element on the page. Replace with a custom `:focus-visible` style.
