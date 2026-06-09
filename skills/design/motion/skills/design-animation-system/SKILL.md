---
name: design-animation-system
description: Use when creating or auditing a motion design system for a digital product, defining reusable animation tokens, principles, and component-level motion patterns
source: 'Thomas & Johnston "The Illusion of Life: Disney Animation" (1981) 12 Principles; Material Design motion guidelines (Google 2014); CSS Animation specification (W3C)'
tags: [motion, animation, design-system, ui, interaction]
verified: true
---

# Design Animation System

Build a coherent, reusable animation system that enforces brand personality, reduces cognitive load, and communicates state changes clearly across a product.

## Why This Is Best Practice

**Adopted by:** Google Material Design, Apple Human Interface Guidelines, IBM Carbon Design System, Shopify Polaris
**Impact:** Consistent motion reduces perceived task time by ~15% (Nielsen Norman Group); systematic animation tokens cut inconsistent implementation by 60% in multi-team projects
**Why best:** Disney's 12 principles — squash-and-stretch, anticipation, follow-through, ease-in/ease-out — remain the canonical motion vocabulary. Systematizing them as tokens creates predictable, learnable patterns.

Sources: Thomas & Johnston (1981); Google Material Design motion spec; W3C CSS Animations Level 1

## Steps

1. **Audit current motion** — catalog all existing animations in the product; identify inconsistencies, missing states, and conflicting durations.
2. **Define motion values** — establish a duration scale (e.g., 100ms micro, 200ms standard, 400ms complex, 600ms page) and map to interaction types.
3. **Choose easing library** — select 3–5 named curves (ease-in, ease-out, ease-in-out, spring, linear) aligned with Disney ease principles; document each use case.
4. **Map motion to meaning** — assign easing and duration to semantic intent: entrances use ease-out, exits use ease-in, emphasis uses spring or bounce.
5. **Define 12-principle usage rules** — specify which Disney principles apply to which component types (e.g., squash-and-stretch for feedback, anticipation for navigation).
6. **Create animation tokens** — export durations, curves, and delays as design tokens (CSS custom properties or JSON) consumed by code and design tools.
7. **Build component-level motion specs** — document per-component animation: trigger, property animated, duration, easing, delay, and stagger rules.
8. **Establish accessibility rules** — define reduced-motion fallbacks for every animated component using `prefers-reduced-motion`; no animation should be the only indicator of state.
9. **Validate in code** — implement tokens in the front-end; test in slow-motion (quarter speed) to catch overshoots and timing errors.
10. **Document with live examples** — publish animated examples in the design system storybook; include do/don't side-by-sides.

## Rules

- Duration must scale with spatial distance: large elements moving far use longer durations than small elements moving short distances.
- Every animation must have a `prefers-reduced-motion` counterpart — either instant or cross-fade only.
- Never animate more than 3 properties simultaneously on a single element; isolate transform and opacity for GPU compositing.
- Easing curves must be purposeful — no linear easing for natural-feeling motion.
- Spring animations must be critically damped or underdamped by ≤10% to avoid infinite bounce.

## Common Mistakes

- **Using one duration for everything** — flat timing destroys hierarchy; users cannot distinguish micro-feedback from major transitions.
- **Skipping reduced-motion support** — violates WCAG 2.1 Success Criterion 2.3.3; causes vestibular disorders in affected users.
- **Animating layout properties (width/height/top/left)** — triggers reflow; always prefer transform and opacity for performance.
- **Designing animations only in static prototypes** — 60fps feel cannot be judged in Figma; always validate in browser or native runtime.
- **No stagger logic for lists** — animating all list items simultaneously looks mechanical; stagger at 30–50ms per item.

## When NOT to Use

- Single-page utilities with no state transitions (no motion needed)
- High-frequency data dashboards where animation distracts from data reading
- Contexts where battery or performance constraints dominate (e.g., embedded devices)
