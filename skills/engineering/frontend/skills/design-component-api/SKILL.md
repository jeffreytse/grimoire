---
name: design-component-api
description: Use when designing the props interface, composition model, or public API for a reusable React or UI component
source: "Atomic Design (Brad Frost, 2013); 'Inversion of Control' pattern (Kent C. Dodds, 2020); React docs — Thinking in React"
tags: [frontend, react, component-api, props, design-system, atomic-design, inversion-of-control]
verified: true
---

# Design Component API

Design component interfaces that are easy to use correctly, hard to use incorrectly, and flexible enough to outlast their initial use case.

## Why This Is Best Practice

**Adopted by:** Radix UI, Headless UI, Chakra UI, and shadcn/ui — the most widely adopted React component libraries all use these patterns; Kent C. Dodds's compound component and inversion of control patterns are curriculum in React training programs globally

**Impact:** Radix UI's accessibility-first, behavior-composable API is now the standard reference for accessible component design; libraries using inversion of control report dramatically fewer "escape hatch" issues where users need to patch around the library

**Why best:** A component API is a contract. Once published and used, breaking changes cost the entire consuming codebase. Designing for composition, not configuration, produces components that can adapt to unforeseen requirements without API changes — extending the useful life of the interface.

## Steps

1. **Define the component's single responsibility** — a component that does two distinct things should be two components; name it by what it renders, not what page it appears on
2. **Apply Atomic Design classification** — classify as atom (no dependencies, pure HTML + style), molecule (composes atoms), organism (business logic + molecules), template, or page; this governs where it lives and who owns it
3. **Prefer composition over configuration** — instead of `<Modal showFooter showCloseButton footerContent={...}>`, use `<Modal><Modal.Header /><Modal.Body /><Modal.Footer /></Modal>`; each compound component controls its own area
4. **Apply inversion of control for behavior** — instead of `onSort={(data) => sort(data)}` (library controls behavior), expose `renderItem` or `getFilteredItems` (consumer controls behavior); use this when built-in behavior will inevitably be wrong for some use case
5. **Type every prop explicitly** — use TypeScript; avoid `any`; prefer union types over booleans for mutually exclusive states (`variant: 'primary' | 'secondary'` not `isPrimary: boolean`)
6. **Forward refs and pass through HTML attributes** — every component that renders a DOM element must accept `ref` and spread remaining props onto the root element; this enables consumers to attach event listeners and accessibility attributes
7. **Test the API with two real use cases before publishing** — if implementing both requires duplicating props or workarounds, the API is wrong; revise before shipping

## Rules

- Do not add props for every possible customization — add a `className` or `style` prop and let consumers handle visual edge cases
- Boolean prop names must describe the positive state: `isDisabled` not `notEnabled`; `isOpen` not `closed`
- Avoid prop names that duplicate HTML attribute names with different semantics: do not use `onClick` to mean "on card click" when the component is not the target element
- A component must never silently ignore a passed prop — if a prop is unsupported, throw a development-only warning

## Examples

**Configuration anti-pattern:** `<DataTable sortable filterable paginationPosition="bottom" emptyStateText="No results" />`
**Composition pattern:** `<DataTable><DataTable.Toolbar><DataTable.Filter /></DataTable.Toolbar><DataTable.Body /><DataTable.Pagination /></DataTable>` — each child is independently replaceable.

## Common Mistakes

- Designing for the first use case only: a `UserCard` with hardcoded user-domain props cannot become a `ProductCard` without a rewrite; design for the rendering pattern, not the data domain
- Boolean prop explosion: when a component accumulates 12 boolean props, it has 4096 possible states — most untested and some contradictory; use variant enums instead
- Not forwarding refs: consuming code that needs to measure a DOM node or integrate with a third-party library will be blocked and forced to wrap the component unnecessarily

## When NOT to Use

- When building a one-off page-specific component that will never be reused outside its current context, investing in a composable compound-component API adds complexity with no return.
- When prototyping or validating a feature idea in a throwaway spike, formal API design slows iteration; defer until the design is confirmed and the component graduates to the shared library.
- When wrapping a third-party component purely for local theming, designing a full props interface duplicates an already-stable upstream API; prefer thin wrappers or CSS overrides instead.
