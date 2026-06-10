---
name: apply-predictable-navigation
description: Use when building multi-page applications, SPAs with multiple routes, or any UI with repeated navigation elements — to ensure focus and input changes don't trigger unexpected page transitions, and navigation appears consistently across views.
source: W3C WCAG 2.1 SC 3.2.1 On Focus (Level A), SC 3.2.2 On Input (Level A), SC 3.2.3 Consistent Navigation (Level AA), SC 3.2.4 Consistent Identification (Level AA); Nielsen Norman Group predictability research; UK GOV.UK Design System navigation principles
tags: [accessibility, wcag, a11y, navigation, predictability, consistent-ui, keyboard-navigation, developer]
related: [apply-focus-management, apply-keyboard-accessibility, design-accessibility-standards]
---

# Implement Predictable Navigation

Ensure focus events and input changes don't trigger unexpected page changes, and repeated UI
elements appear consistently across pages — so keyboard and screen reader users can build an
accurate mental model of the interface.

## Why This Is Best Practice

**Adopted by:** WCAG 2.1 SCs 3.2.1 and 3.2.2 (Level A) are required universally. SCs 3.2.3 and 3.2.4
(Level AA) are required by Section 508 (US), EU EN 301 549, and UK PSBAR 2018. Nielsen Norman Group's
research on predictability (2000–2023) consistently identifies unexpected context changes as the
#1 cause of user disorientation. GOV.UK Design System mandates consistent navigation as a core principle.
**Impact:** Screen reader users build a spatial model of a page as they navigate. When focus alone
changes the page context (opens a modal, submits a form, navigates away), the model is invalidated
without warning. NNGroup (2023) found that unexpected context changes increase task failure rate by 35%
for screen reader users. Speech input users who tab through fields rely on focus not triggering side
effects — unexpected navigation breaks their flow invisibly.
**Why best:** The alternative — triggering actions on focus for "convenience" (auto-submit on field
selection, auto-navigate on radio select) — creates side effects that keyboard users cannot predict or
undo. A submit button that fires on focus is inaccessible by definition; a select that navigates on
`change` without warning removes user agency.

Sources: W3C WCAG 2.1 SCs 3.2.1–3.2.4 (2018); NNGroup Predictability research (2023); GOV.UK
Design System navigation guidance

## Steps

### Step 1: Never trigger context changes on focus alone (3.2.1)

Focus must not submit a form, open a new window, navigate to a new page, or cause any significant
UI change:

```javascript
// Wrong — submits if a value is present when user tabs into field
searchInput.addEventListener('focus', () => {
  showAutocomplete();
  if (searchInput.value) submitSearch(); // fires whenever field receives focus
});

// Right — autocomplete shows on focus; user must explicitly submit
searchInput.addEventListener('focus', showAutocomplete);
searchInput.addEventListener('blur',  hideAutocomplete);
searchForm.addEventListener('submit', handleSearch);
```

Permitted on focus: showing a tooltip, revealing help text, opening a dropdown that lists
options without selecting one. The effect must be reversible and must not move the page or
submit data.

### Step 2: Don't trigger context changes on input without warning (3.2.2)

Changing a form control value must not navigate or submit unless the user was explicitly informed:

```html
<!-- Wrong — selecting a country navigates immediately -->
<select id="country" onchange="window.location = '/store/' + this.value">
  <option value="us">United States</option>
  <option value="uk">United Kingdom</option>
</select>

<!-- Right — select + explicit submit button -->
<label for="country">Country</label>
<select id="country">
  <option value="us">United States</option>
  <option value="uk">United Kingdom</option>
</select>
<button type="submit">Go</button>
```

If you must act on `change` (live filter, currency switcher), inform the user with visible
associated text:

```html
<label for="currency">Currency</label>
<select id="currency" aria-describedby="currency-note">
  <option value="USD">USD</option>
  <option value="EUR">EUR</option>
</select>
<p id="currency-note">Prices update immediately when you change currency.</p>
```

### Step 3: Keep repeated navigation in the same relative order (3.2.3)

Navigation components that appear on multiple pages must appear in the same relative position
and DOM order across all pages:

```jsx
// React — shared nav component guarantees consistent order
function MainNav() {
  return (
    <nav aria-label="Main navigation">
      <NavLink to="/">Home</NavLink>
      <NavLink to="/products">Products</NavLink>
      <NavLink to="/contact">Contact</NavLink>
    </nav>
  );
}
// Import MainNav in every page layout — never inline nav per-page
```

"Same relative order" allows adding new items (insert "About" between "Products" and "Contact"
on all pages). Reversing order between pages or making per-page nav variations fails 3.2.3.

### Step 4: Identify components that serve the same function consistently (3.2.4)

Components serving the same function must use the same label, icon, and accessible name
across all pages:

```html
<!-- Wrong: edit action uses different labels on different pages -->
<!-- product page: -->
<button aria-label="Edit product">✏️</button>
<!-- profile page: -->
<button aria-label="Modify profile">✏️</button>

<!-- Right: same label everywhere for same function -->
<button aria-label="Edit">✏️</button>
```

Audit for inconsistency:
- "Edit" vs "Modify" vs "Change" for the same action
- "Delete" vs "Remove" vs "Clear" for the same action
- Icon-only buttons with different `aria-label` values on different pages
- Different tooltip text for the same icon across views

## When NOT to Use

- **Context-sensitive components** — a "Save" button in a text editor and a "Save" button in
  settings both serve the same function — that's consistent. If the function genuinely differs,
  different labels are correct.
- **Intentional `change`-driven navigation with explicit warning** — a language selector that
  reloads the page is acceptable if visible text states "Changing language will reload the page."
  Provide a submit button as the preferable alternative.

## Common Mistakes

**React `onChange` treated as native `change`.** `onChange` fires on every keystroke in React
controlled inputs (unlike native `change`, which fires on blur). Submitting on `onChange` in a
text field fails 3.2.2. Use `onBlur` for validation; explicit submit buttons for form submission.

**Different mobile and desktop nav order.** CSS `order` or `flex-direction: row-reverse` visually
reorders nav items. Screen readers follow DOM order, not visual order. If DOM order differs between
breakpoints (hidden mobile nav with separate desktop nav), screen reader users see a different
sequence than sighted users. Keep one nav in the DOM; use CSS for visual layout only.

**Breadcrumb in `<nav>` on some pages, `<div>` on others.** Screen reader users navigate by
landmarks. Inconsistent landmark use for the same component breaks 3.2.4 and `apply-skip-navigation`.
