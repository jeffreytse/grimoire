---
name: write-design-system-token
description: Use when defining, naming, or documenting design tokens for a design system to ensure consistent visual decisions across platforms
source: W3C Design Tokens Community Group specification (2023); Salesforce Lightning Design System; Figma Token Studio documentation
tags: [design-tokens, design-system, tokens, theming, css-variables, component-design]
verified: true
---

# Write Design System Token

Define design tokens with clear naming, semantic layers, and platform-ready output so design decisions are consistent and themeable.

## Why This Is Best Practice

**Adopted by:** Salesforce Lightning, IBM Carbon, Adobe Spectrum, GitHub Primer, Atlassian Design System
**Impact:** Spotify reported a 60% reduction in design inconsistency across platforms after adopting a structured token layer; the W3C Design Tokens format (`.tokens.json`) is becoming the cross-tool standard for Figma→code pipelines

**Why best:** Tokens transform design decisions from hardcoded values scattered across codebases into a single source of truth. Without a semantic layer, changing a brand color requires touching thousands of individual components instead of one token value.

## Steps

1. **Define the primitive layer** — Create raw value tokens with no semantic meaning. These are named for what they are: `color-blue-500: #2563EB`, `spacing-4: 16px`, `font-size-base: 16px`. Never use these tokens directly in components.
2. **Define the semantic layer** — Create tokens that reference primitives and carry meaning: `color-action-primary: {color-blue-500}`, `space-component-padding: {spacing-4}`. Components consume semantic tokens, not primitives.
3. **Define component tokens (optional)** — For large systems, add a third layer: `button-background-default: {color-action-primary}`. This allows component-level overrides without breaking the semantic layer.
4. **Name tokens using the pattern** — `[category]-[concept]-[variant]-[state]`: `color-text-primary`, `color-text-secondary`, `color-border-focus`, `space-inset-md`. Avoid color names in semantic tokens (`color-blue` not allowed at semantic level).
5. **Document the decision** — For each token, record: the problem it solves, which components use it, and what must change to update it. A token without documentation is a mystery to the next designer.
6. **Output platform formats** — Use Style Dictionary or Token Studio to transform tokens into CSS custom properties, iOS Swift constants, Android XML, and JSON. One source, many targets.
7. **Version and publish** — Store tokens in a shared repository. Publish changes with semantic versioning. Breaking changes (renames, deletions) are major version bumps; additions are minor.

## Rules

- Components must only consume semantic tokens, never primitive tokens directly.
- Token names must describe usage (purpose), not appearance (value). `color-danger` not `color-red`.
- Every token value must be traceable back to a primitive via a reference chain — no hardcoded values in semantic or component layers.
- Rename tokens with a deprecation period, not immediate deletion — consumers need migration time.
- Tokens and components must be versioned together; a component shipped without token updates is a breaking change.

## Examples

Salesforce Lightning tokens:
- Primitive: `color-palette-brand-dark: #0176D3`
- Semantic: `color-brand: {color-palette-brand-dark}`
- Component: `button-color-brand: {color-brand}`

A brand update changes `color-palette-brand-dark` to `#0154A5` — every component consuming `button-color-brand` updates automatically.

## Common Mistakes

- **Semantic tokens named after values** — `color-blue-primary` breaks when the brand changes to green; semantic tokens must describe intent, not appearance.
- **Skipping the semantic layer** — Components consuming primitives directly means a brand refresh requires updating every component individually.
- **No tooling for output** — Hand-maintaining CSS variables and iOS Swift files from the same source diverges within weeks; automation is not optional.
