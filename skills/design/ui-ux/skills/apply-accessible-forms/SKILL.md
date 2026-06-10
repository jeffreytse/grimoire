---
name: apply-accessible-forms
description: Use when building any form — login, checkout, search, settings — to ensure all inputs are correctly labeled, errors are programmatically associated, and users can complete the form with a keyboard and screen reader.
source: W3C WCAG 2.1 SC 1.3.5 Identify Input Purpose (Level AA), SC 3.3.1 Error Identification (Level A), SC 3.3.2 Labels or Instructions (Level A), SC 3.3.3 Error Suggestion (Level AA); WAI Forms Tutorial
tags: [accessibility, wcag, a11y, forms, labels, error-messages, autocomplete, screen-reader]
related: [design-accessibility-standards, apply-keyboard-accessibility, apply-aria-roles]
---

# Implement Accessible Forms

Label every input, associate error messages programmatically, mark required fields visually and semantically, and enable browser autocomplete.

## Why This Is Best Practice

**Adopted by:** WCAG 2.1 SCs 3.3.1 (Level A), 3.3.2 (Level A), and 3.3.3 (Level AA)
are required by all major accessibility laws. WAI's Forms Tutorial is the W3C-published
reference guide used by Google, Microsoft, and all major frameworks (React, Angular,
Vue form validation libraries implement these patterns). WebAIM Million 2024 found that
unlabeled form inputs is the #3 most common WCAG failure (45.4% of pages with forms).
**Impact:** Users with motor disabilities rely on autocomplete to avoid retyping personal
information; users with cognitive disabilities rely on clear labels and inline error
messages. Deque Systems research found that form errors not associated with inputs add
an average of 4–6 minutes of correction time for screen reader users per form
submission error.
**Why best:** Placeholder text disappears on input, leaving no label visible during
entry. `aria-label` labels are read only on focus, not while typing. A visible, persistent
`<label>` associated via `for`/`id` is the only mechanism that ensures the field purpose
is always announced.

Sources: W3C WCAG 2.1 SC 1.3.5, 3.3.1, 3.3.2, 3.3.3 (2018); WAI Forms Tutorial;
WebAIM Million 2024; Deque Systems

## Steps

### Step 1: Associate every input with a visible `<label>`

```html
<!-- Wrong — no label, placeholder only -->
<input type="email" placeholder="Email address">

<!-- Wrong — label not associated (screen reader won't link them) -->
<label>Email address</label>
<input type="email">

<!-- Right — for/id association -->
<label for="email">Email address</label>
<input type="email" id="email" name="email">

<!-- Right — wrapping label (implicit association) -->
<label>
  Email address
  <input type="email" name="email">
</label>
```

Every visible form control needs a `<label>`. Do not use `aria-label` or
`aria-labelledby` instead of a visible `<label>` unless the field has a persistent
visible label in the surrounding UI (e.g., a table column header).

### Step 2: Mark required fields both visually and with `aria-required`

```html
<!-- Right — asterisk with legend + aria-required -->
<form>
  <p>Fields marked with <span aria-hidden="true">*</span>
     <span class="sr-only">asterisk</span> are required.</p>

  <label for="name">
    Full name <span aria-hidden="true">*</span>
  </label>
  <input type="text" id="name" name="name"
         required aria-required="true">
</form>
```

`required` triggers browser validation; `aria-required="true"` announces the
requirement to screen readers before submission.

### Step 3: Associate error messages with `aria-describedby`

```html
<!-- Right — error linked to input -->
<label for="email">Email address</label>
<input type="email" id="email" name="email"
       aria-describedby="email-error"
       aria-invalid="true">
<span id="email-error" class="error-message" role="alert">
  Enter a valid email address (e.g., name@example.com)
</span>
```

Key points:
- `aria-invalid="true"` announces the field is in an error state
- `aria-describedby` links the error message to the input — screen reader reads both
- `role="alert"` announces the error immediately when injected into the DOM
- Error message must describe what's wrong AND suggest how to fix it (WCAG 3.3.3)

### Step 4: Group related fields with `<fieldset>` and `<legend>`

```html
<!-- Right — radio group with shared context -->
<fieldset>
  <legend>Preferred contact method</legend>
  <label><input type="radio" name="contact" value="email"> Email</label>
  <label><input type="radio" name="contact" value="phone"> Phone</label>
  <label><input type="radio" name="contact" value="sms"> SMS</label>
</fieldset>

<!-- Right — billing address group -->
<fieldset>
  <legend>Billing address</legend>
  <label for="street">Street</label>
  <input type="text" id="street" name="billing-street">
  <!-- more fields... -->
</fieldset>
```

`<fieldset>` + `<legend>` is required for radio/checkbox groups. Screen readers
announce the legend before each option ("Preferred contact method, Email radio button").

### Step 5: Enable browser autocomplete with correct `autocomplete` attributes

```html
<input type="text"   id="name"        autocomplete="name">
<input type="email"  id="email"       autocomplete="email">
<input type="tel"    id="phone"       autocomplete="tel">
<input type="text"   id="address"     autocomplete="street-address">
<input type="text"   id="city"        autocomplete="address-level2">
<input type="text"   id="postcode"    autocomplete="postal-code">
<input type="text"   id="cc-number"   autocomplete="cc-number">
<input type="month"  id="cc-expiry"   autocomplete="cc-exp">
```

WCAG SC 1.3.5 (Level AA) requires `autocomplete` attributes for personal data fields.
Autocomplete is essential for motor-disabled users who cannot type form data manually.

## When NOT to Use

- **OTP/security codes** — use `autocomplete="one-time-code"` and `autocomplete="off"` only for fields where pre-filling creates a genuine security risk (not for developer convenience).
- **Search fields without persistent labels** — a search input with a visible search button and surrounding context may use `aria-label="Search"` if no visible `<label>` is present.

## Common Mistakes

**Placeholder-as-label.** Placeholder disappears on focus or input. Screen readers may not reliably announce placeholder as an accessible name. Always use a persistent `<label>`.

**Error messages not in the DOM when announced.** Injecting `role="alert"` content into a pre-existing element doesn't trigger an announcement — the `role` must be present in the DOM before the content is added. Use `aria-live="assertive"` on a persistent container.

**`autocomplete="off"` on all fields.** Disabling autocomplete is a WCAG violation for personal data fields (SC 1.3.5). Reserve `off` for fields where pre-filling is genuinely harmful (security code, captcha).
