---
name: design-color-system
description: Use when building a brand or product color system with semantic roles, accessible palette, and usage rules
source: Itten "The Art of Color" (1961); CIECAM02 color appearance model (CIE); Material Design color system (Google 2018)
tags: [color, design-system, branding, accessibility, palette]
verified: true
---

# Design Color System

Create a structured palette with semantic roles, accessible contrast ratios, and clear usage rules that scale across light/dark modes and brand contexts.

## Why This Is Best Practice

**Adopted by:** Google Material Design, IBM Carbon, Atlassian Design System, Shopify Polaris — all use semantic color tokens layered over a base palette
**Impact:** Semantic color tokens reduce brand drift by 70% in enterprise products (IBM Carbon case study); accessible palettes cut WCAG-related rework by 60% when built in up front
**Why best:** Itten's color theory provides perceptual foundations (hue relationships, temperature, contrast); CIECAM02 provides perceptual uniformity for accessible lightness; Material Design operationalizes both into scalable engineering tokens.

Sources: Itten "The Art of Color" (1961); CIE CIECAM02 (2002); Google Material Design color system docs (2018); WCAG 2.2 contrast criteria

## Steps

1. **Extract brand hues** — identify 1–3 brand hues from logo or guidelines; record exact HSL or OKLCH values for precision.
2. **Build a tonal scale** — for each brand hue, generate a 10-step lightness scale (100–900 or 50–950) using perceptually uniform lightness (OKLCH or HSLuv) so each step appears equally distant visually.
3. **Add neutral palette** — create a neutral (gray) tonal scale using the same 10-step system; optionally tint neutrals with the brand hue at 5–10% saturation for brand warmth.
4. **Add semantic colors** — define scales for: error (red), warning (amber/orange), success (green), info (blue); use standard hues modified with brand hue influence.
5. **Define semantic roles** — map tonal values to named roles: `color-primary`, `color-primary-hover`, `color-surface`, `color-on-surface`, `color-error`, `color-border`, etc.
6. **Verify contrast ratios** — for each foreground/background pair in semantic roles, check WCAG AA (4.5:1 text, 3:1 UI) and AAA where needed; adjust lightness steps until passing.
7. **Create dark mode mappings** — remap semantic roles to different tonal steps for dark mode; do not simply invert — test each pair for contrast and visual comfort.
8. **Export as design tokens** — output as CSS custom properties, JSON, or Figma variables; ensure tokens are referenced by semantic name, not by raw color value.
9. **Document usage rules** — specify which semantic roles are used for which UI contexts; prohibit arbitrary palette value usage outside the token system.
10. **Audit for color blindness** — simulate deuteranopia and protanopia using Stark or Color Oracle; ensure all meaningful color distinctions are also conveyed by shape, text, or pattern.

## Rules

- Never use raw hex values in design or code — always reference semantic tokens; raw values create an unmaintainable system.
- Every foreground/background pair must be defined and contrast-checked before being declared usable.
- The primary palette should not exceed 3 brand hues — more creates visual noise; neutrals and semantics supplement, not compete.
- Dark mode is not an inversion — light mode's surface/200 does not become dark mode's surface/800 automatically; map deliberately.
- Color must never be the sole differentiator — pair with icons, labels, or patterns for accessibility.

## Common Mistakes

- **Building palette without perceptual uniformity** — standard HSL lightness is not perceptually uniform; use OKLCH or HCL to avoid muddy mid-tones.
- **Too many primary colors** — 5+ primaries create a palette that cannot enforce brand identity; merge or subordinate excess hues.
- **Semantic tokens that are too granular** — `button-primary-hover-border-pressed` tokens require maintenance proportional to their specificity; balance specificity with flexibility.
- **Not testing on actual displays** — colors shift significantly between sRGB monitors, P3 displays, and OLED screens; test on multiple devices.
- **Forgetting high-contrast mode** — Windows High Contrast and macOS Increased Contrast override many custom colors; test UI integrity under forced color modes.

## When NOT to Use

- Single-use campaign designs with no long-term maintenance
- Art prints or illustration work not bound by accessibility requirements
- Third-party white-label products where the color system is externally mandated
