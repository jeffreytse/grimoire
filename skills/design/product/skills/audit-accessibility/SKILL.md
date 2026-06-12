---
name: audit-accessibility
description: Use when auditing a product or component for accessibility compliance or when preparing for WCAG 2.1 AA certification
source: W3C WCAG 2.1 (2018); WAI-ARIA Authoring Practices 1.2; Deque axe-core open-source ruleset
tags: [accessibility, wcag, a11y, screen-reader, aria, inclusive-design, audit]
verified: true
---

# Audit Accessibility

Systematically test a product against WCAG 2.1 AA to find and prioritize accessibility barriers before they reach users with disabilities.

## Why This Is Best Practice

**Adopted by:** UK Government Digital Service (mandatory WCAG 2.1 AA), US Section 508, EU EN 301 549, and all major tech companies under legal obligation or voluntary commitment
**Impact:** WebAIM Million report (2024) found 95.9% of top 1M homepages had detectable WCAG failures; automated tools catch ~30–40% of issues — manual and assistive technology testing is required for full coverage

**Why best:** Accessibility is a legal requirement in most jurisdictions and expands the addressable market by ~1.3 billion people globally with disabilities. Retrofitting accessibility after launch costs 3–5× more than building it in.

## Steps

1. **Run automated scanning** — Use axe-core (browser extension or CI integration) on every page/view. Deque Axe DevTools or Lighthouse accessibility audit catch structural violations: missing alt text, insufficient color contrast, unlabeled form fields, missing landmark roles.
2. **Audit color contrast** — Check all text against WCAG 2.1 AA: normal text ≥ 4.5:1 ratio, large text (≥18pt or 14pt bold) ≥ 3:1. Use the Colour Contrast Analyser or browser DevTools.
3. **Test keyboard navigation** — Tab through the entire interface using only a keyboard. Verify: visible focus indicator on all interactive elements, logical tab order, no keyboard traps, all interactive elements reachable.
4. **Test with a screen reader** — Use NVDA+Chrome (Windows) or VoiceOver+Safari (macOS/iOS). Test: page title announced correctly, headings navigable in sequence, form labels associated, dynamic content announced (live regions), modal focus management.
5. **Audit heading structure** — Every page should have one `<h1>`. Headings must not skip levels (h1→h3 with no h2 is a violation). Use the browser's accessibility tree or a heading visualizer extension.
6. **Check images and media** — All informative images need descriptive `alt` text. Decorative images need `alt=""`. Videos need captions; audio needs transcripts.
7. **Audit interactive components** — For each custom widget (modal, dropdown, date picker, slider): verify ARIA roles, states, and properties match WAI-ARIA Authoring Practices patterns. Keyboard interactions must match the expected pattern.
8. **Document findings in a severity matrix** — Rate each issue: Critical (blocks access), Serious (significantly limits access), Moderate (inconvenient), Minor (cosmetic). Fix critical and serious before launch.

## Rules

- Automated tools alone are never sufficient — manual screen reader testing is mandatory for any shipped product.
- Color must not be the only means of conveying information (e.g., red = error must also have a text label or icon).
- All interactive elements must have a programmatically determinable name — not just a visible label.
- Focus must never be trapped except in modals, and modals must return focus to the trigger element when closed.
- Audit with real assistive technology users at least once per major release cycle.

## Examples

A modal audit reveals the close button has `aria-label="close"` but pressing Escape does not close the modal and focus does not return to the trigger. Fix: add `keydown` Escape handler and restore focus on close. This matches the WAI-ARIA Dialog pattern exactly.

## Common Mistakes

- **Audit = run Lighthouse** — Lighthouse catches ~30% of WCAG violations; teams that stop there ship inaccessible products confidently.
- **ARIA without semantics** — Adding `role="button"` to a `<div>` without also adding `tabindex="0"` and keyboard handlers creates the appearance of accessibility without the reality.
- **Fixing automation findings only** — Critical barriers like incorrect screen reader announcements for dynamic content are invisible to automated tools.

## When NOT to Use

- Do not substitute a WCAG 2.1 AA audit for a higher standard when the product serves users with cognitive disabilities or when the organization has committed to WCAG 2.2 or AAA — a passing AA audit does not guarantee the product is usable by all disability groups.
- Do not run an accessibility audit as a one-time pre-launch gate rather than a recurring practice — new components, content updates, and third-party script changes introduce violations continuously, so a single audit provides only a point-in-time snapshot.
- Do not treat an audit as complete if no real assistive technology users were involved in testing — automated tools and developer screen reader checks consistently miss interaction patterns that diverge from how actual users navigate with specialized hardware or custom AT configurations.
