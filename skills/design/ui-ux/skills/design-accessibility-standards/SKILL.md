---
name: design-accessibility-standards
description: Use when auditing or designing UI components to meet WCAG 2.2, Section 508, or ARIA accessibility requirements
source: WCAG 2.2 (W3C Web Content Accessibility Guidelines); Section 508 (US Rehabilitation Act); ARIA specification (W3C WAI-ARIA 1.2)
tags: [accessibility, wcag, aria, a11y, inclusive-design]
verified: true
---

# Design Accessibility Standards

Ensure digital interfaces are perceivable, operable, understandable, and robust for users with disabilities, meeting legal and ethical standards.

## Why This Is Best Practice

**Adopted by:** All US federal agencies (Section 508 mandate), EU (EN 301 549), UK (PSBAR 2018), Apple App Store accessibility review, Google Play accessibility guidelines
**Impact:** 1 in 4 US adults has a disability; accessible design improves usability for all users (curb-cut effect); ADA non-compliance lawsuits exceeded 4,000/year in 2022
**Why best:** WCAG 2.2 POUR framework (Perceivable, Operable, Understandable, Robust) provides a complete, testable taxonomy covering visual, motor, cognitive, and auditory needs.

Sources: W3C WCAG 2.2 (2023); W3C WAI-ARIA 1.2 (2023); US Section 508 (29 USC §794d)

## Steps

1. **Set conformance target** — define required level: WCAG 2.2 Level AA (legal minimum for most jurisdictions); Level AAA for healthcare or government.
2. **Audit color contrast** — all text must meet 4.5:1 ratio for normal text, 3:1 for large text (AA); use WebAIM Contrast Checker or Figma plugins; fix by adjusting foreground or background.
3. **Audit keyboard navigation** — tab through every interactive element; verify visible focus indicators, logical tab order, no keyboard traps; test skip navigation links.
4. **Audit screen reader compatibility** — test with NVDA+Chrome (Windows) and VoiceOver+Safari (macOS/iOS); verify all images have alt text, icons have aria-label, and live regions announce dynamic content.
5. **Audit form accessibility** — every input has an associated `<label>`; error messages are programmatically linked to inputs via `aria-describedby`; required fields marked both visually and with `aria-required`.
6. **Audit ARIA usage** — validate ARIA roles match element semantics; no ARIA attribute conflicts; verify landmarks (main, nav, aside, header, footer) are present and unique.
7. **Audit motion and animation** — implement `prefers-reduced-motion` for all animations; no content flashes more than 3 times per second (seizure threshold).
8. **Audit touch targets** — minimum 44×44px touch target for all interactive elements (iOS HIG) or 48×48dp (Material); spacing between targets ≥8px.
9. **Audit text and readability** — body text minimum 16px; line height 1.4–1.6; users can resize text 200% without loss of content or functionality.
10. **Run automated scan + manual review** — use axe-core or Deque Axe DevTools for automated coverage (~35% of issues); manually test the remaining 65% using keyboard and screen reader.

## Rules

- Automated tools catch only 30–40% of accessibility issues; manual keyboard and screen reader testing is non-negotiable.
- Color must never be the sole means of conveying information — pair with text labels, icons, or patterns.
- Focus indicators must be visible with at minimum a 3:1 contrast ratio against adjacent colors (WCAG 2.2 SC 2.4.11).
- All non-text content must have a text alternative; decorative images use `alt=""` to be ignored by screen readers.
- ARIA should supplement, not replace, semantic HTML — prefer native elements (`<button>`, `<nav>`) over ARIA roles on `<div>`.

## Common Mistakes

- **Removing focus outlines for aesthetics** — keyboard users become completely lost; style focus rings, do not remove them.
- **Using placeholder text as the only label** — placeholder disappears on input; screen readers may not announce it as a label.
- **Icon-only buttons without accessible names** — `<button><svg>...</svg></button>` is read as "button" by screen readers; add `aria-label`.
- **Modal dialogs without focus trapping** — keyboard focus escaping a modal creates an unusable experience for screen reader users.
- **Testing only with sighted keyboard** — screen reader users encounter different DOM ordering and ARIA announcement issues invisible to sighted keyboard testers.

## When NOT to Use

- Internal tools with 100% verified non-disabled user base (rare; still best practice)
- Prototypes explicitly not intended for user testing (wireframes)
- Content designed exclusively for print media
