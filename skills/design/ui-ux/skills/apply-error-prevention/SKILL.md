---
name: apply-error-prevention
description: Use when building forms that submit legal, financial, or irreversible data — to ensure users can review, confirm, or reverse submissions, and are not asked to re-enter information they've already provided.
source: W3C WCAG 2.1 SC 3.3.4 Error Prevention Legal Financial Data (Level AA); W3C WCAG 2.2 SC 3.3.7 Redundant Entry (Level A); Nielsen Norman Group form usability research; UK GOV.UK Design System form patterns
tags: [accessibility, wcag, a11y, error-prevention, forms, confirmation, undo, cognitive-accessibility]
related: [apply-accessible-forms, apply-accessible-authentication, design-accessibility-standards]
---

# Implement Error Prevention

Give users the ability to review, confirm, or reverse submissions before data is lost or
committed — and never ask them to re-enter information they've already provided.

## Why This Is Best Practice

**Adopted by:** WCAG 2.1 SC 3.3.4 (Level AA) is required by Section 508 (US), EU EN 301 549, and
UK PSBAR 2018. WCAG 2.2 SC 3.3.7 (Level A) is required wherever WCAG 2.2 is mandated. GOV.UK
Design System, the US Web Design System, and checkout UX patterns at Amazon, Stripe, and Shopify
all implement review steps and session data pre-population as standard form patterns.
**Impact:** Nielsen Norman Group's e-commerce checkout research (2023) found that requiring users
to re-enter previously submitted data (billing address re-entered for shipping, password entered
twice on separate screens) increased form abandonment by 37% in users with cognitive disabilities.
Users with motor disabilities who make input errors in irreversible forms (deleted accounts, submitted
legal documents) face disproportionate impact — they cannot quickly re-complete the form.
**Why best:** Inline error messages (WCAG 3.3.1–3.3.3, `apply-accessible-forms`) catch errors
at the field level. Error prevention operates at the submission level: review steps, soft-delete with
restore windows, and session data pre-population are the only mechanisms that prevent already-submitted
errors from causing harm.

Sources: W3C WCAG 2.1 SC 3.3.4 (2018); WCAG 2.2 SC 3.3.7 (2023); NNGroup E-commerce UX (2023);
GOV.UK Design System Check Your Answers pattern

## Steps

### Step 1: Add a review step before irreversible submissions (3.3.4)

For legal agreements, financial transactions, account deletion, and any action that cannot be
undone, show a summary page before final submission:

```html
<!-- Step 1: Form input -->
<form action="/checkout/review" method="post">
  <!-- form fields -->
  <button type="submit">Review order</button>
</form>

<!-- Step 2: Review page — show all entered data before committing -->
<section aria-labelledby="review-heading">
  <h1 id="review-heading">Review your order</h1>

  <dl>
    <dt>Delivery address</dt>
    <dd>123 Main St, London, SW1A 1AA</dd>

    <dt>Payment</dt>
    <dd>Visa ending 4242</dd>

    <dt>Total</dt>
    <dd>£49.99</dd>
  </dl>

  <a href="/checkout">Change order details</a>

  <form action="/checkout/confirm" method="post">
    <input type="hidden" name="order_token" value="{{ order_token }}">
    <button type="submit">Place order</button>
  </form>
</section>
```

The review step satisfies "can be checked" — users read the summary. The "Change" link
satisfies "can be corrected." Final submit satisfies "can be confirmed."

### Step 2: Provide reversibility for destructive actions (3.3.4)

When an action cannot have a review step (single-button delete), provide a reversal window:

```javascript
// Soft-delete: mark as deleted, restore within 30 days
async function deleteAccount(userId) {
  await db.users.update(userId, {
    deleted_at: new Date(),
    restore_until: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000)
  });
  // Send restoration email with link
  await sendEmail(user.email, 'account-deletion-confirmation', {
    restoreUrl: `https://example.com/restore?token=${generateRestoreToken(userId)}`
  });
}
```

Alternatives to soft-delete:
- **Undo toast** — show "Deleted. [Undo]" immediately after action; user has 5–10 seconds to undo
- **Confirmation dialog** — "Are you sure? This cannot be undone" — satisfies "can be confirmed"
- **Trash/archive** — move to recoverable state instead of permanent deletion

```javascript
// Undo pattern — defer actual deletion
function deleteItem(itemId) {
  const item = removeFromUI(itemId);
  const toastTimeout = setTimeout(() => permanentlyDelete(itemId), 8000);

  showToast(`Item deleted. <button onclick="undoDelete()">Undo</button>`, 8000);

  window.__undoDelete = () => {
    clearTimeout(toastTimeout);
    restoreToUI(item);
  };
}
```

### Step 3: Pre-populate fields from previously entered session data (3.3.7)

Don't ask users to re-enter information already collected in the current session:

```javascript
// Checkout — pre-populate shipping from billing address
function initShippingForm() {
  const billing = getSessionData('billing_address');
  if (!billing) return;

  // Pre-fill shipping fields
  document.getElementById('ship-name').value    = billing.name;
  document.getElementById('ship-street').value  = billing.street;
  document.getElementById('ship-city').value    = billing.city;
  document.getElementById('ship-postcode').value = billing.postcode;

  // Let user override if different
  document.getElementById('same-as-billing').checked = true;
}
```

```html
<!-- Profile data available — pre-fill checkout contact details -->
<label for="email">Email</label>
<input type="email" id="email" name="email"
       value="{{ current_user.email }}"
       autocomplete="email">
<!-- User can change it; they shouldn't have to type it again -->
```

Exceptions: security re-entry (password confirmation on account deletion, PIN for payments)
is exempt — deliberate re-entry prevents accidental irreversible actions.

### Step 4: Allow correction at the review stage

The review step must offer a way to correct any field, not just confirm or cancel:

```html
<!-- Each section of the review has a change link -->
<section>
  <h2>Contact information
    <a href="/checkout/contact" aria-label="Change contact information">Change</a>
  </h2>
  <p>jane@example.com</p>
</section>

<section>
  <h2>Delivery address
    <a href="/checkout/address" aria-label="Change delivery address">Change</a>
  </h2>
  <p>123 Main St...</p>
</section>
```

"Change" links must return the user to the specific step, not the form beginning, and must
pre-populate the form with the values they entered.

## When NOT to Use

- **Security re-authentication** — requiring password entry before deleting an account or
  confirming a large transfer is a security control, not redundant entry. SC 3.3.7 explicitly
  exempts re-entry for security purposes.
- **Inherently confirmation-free processes** — real-time communications (sending a chat message,
  posting a comment) don't require a review step. Provide edit/delete functionality post-send instead.

## Common Mistakes

**Confirmation dialog with no way to review what will be deleted.** "Delete account — are you
sure?" without showing what data will be lost satisfies the letter of 3.3.4 but not the spirit.
Show what will be deleted ("This will delete your profile, 47 posts, and all order history").

**Multi-step form that doesn't preserve previous answers.** If a user goes back to correct step 2,
step 3 data must not be wiped. Pre-populate all steps from session state.

**Undo toast that auto-dismisses before screen reader announces it.** The undo toast must persist
long enough for screen reader announcement (apply-aria-live-regions). Use `role="status"` and
ensure the toast stays visible for at least 8 seconds.
