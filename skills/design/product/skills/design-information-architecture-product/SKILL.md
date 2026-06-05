---
name: design-information-architecture-product
description: Use when structuring content, navigation, or feature hierarchies for a product so users can find information predictably and efficiently
source: Morville & Rosenfeld "Information Architecture for the Web and Beyond" (2015, 4th ed.); IA Institute principles; Card sorting research (Spencer "Card Sorting" 2009)
tags: [ia, navigation, ux, taxonomy, content-structure]
verified: true
---

# Design Information Architecture for Products

Organize content and features into a structure that matches users' mental models so they can find and understand information without effort.

## Why This Is Best Practice

**Adopted by:** Nielsen Norman Group IA methodology, GovUK design system, Salesforce Lightning IA, Apple App Store taxonomy
**Impact:** Good IA reduces user navigation time by 30–50% (Nielsen Norman Group findability research); poor IA is the #1 cause of "lost" users on content-heavy sites
**Why best:** Morville & Rosenfeld's three circles (users, content, context) remain the definitive IA framework — ignoring any one dimension produces structures that look logical but fail in use.

Sources: Morville & Rosenfeld (2015) Ch. 1–4; Spencer "Card Sorting" (2009); IA Institute Core Concepts

## Steps

1. **Inventory existing content** — catalog all content types, features, and data objects; record volume, ownership, and update frequency.
2. **Define user goals and mental models** — from research (interviews, analytics, search logs) identify the top 10 tasks users arrive to complete; note the vocabulary they use.
3. **Run open card sort** — give 30–50 representative content items to 15–20 users; ask them to group and name groups; use similarity matrix to reveal natural clusters.
4. **Analyze sort results** — use Optimal Workshop or manual similarity scoring; identify agreement clusters (>60% agreement = strong grouping) and outliers.
5. **Draft candidate taxonomy** — create 2–3 structural options based on card sort data; explore top-down (organization's logic) vs. bottom-up (user groupings) hybrids.
6. **Run tree test** — present the candidate taxonomy as a text-only tree; ask users to find specific items; measure directness, success rate, and time on task.
7. **Refine structure** — fix navigation paths where <70% of participants succeed; rename labels using users' own vocabulary from card sort and search logs.
8. **Design navigation system** — define global nav, local nav, contextual links, search, and breadcrumbs; document each system's role and interaction patterns.
9. **Create wireframe flows** — sketch key user journeys through the IA; validate depth (prefer ≤3 clicks to any primary content) vs. breadth trade-offs.
10. **Validate with usability test** — run task-based sessions on the implemented navigation; iterate until directness >80% on top-5 tasks.

## Rules

- Navigation labels must use the users' vocabulary, not the organization's internal terminology.
- Maximum meaningful depth is 3 levels for most products; deeper hierarchies require breadcrumbs and cross-links.
- Every item must belong to exactly one primary location; cross-links are permitted but must not replace clear taxonomy.
- Search is a supplement, not a substitute, for clear navigation — do not rely on search to compensate for poor IA.
- IA must be validated with real users via tree testing before front-end build begins.

## Common Mistakes

- **Mirroring org chart in navigation** — internal structure rarely matches user mental models; always validate with card sorting.
- **Too many top-level categories** — more than 7±2 primary nav items overloads working memory (Miller 1956).
- **Inconsistent naming patterns** — mixing nouns and verbs at the same level (e.g., "Reports" and "Create Dashboard") creates confusion.
- **No cross-linking for adjacent content** — users follow scent trails; missing contextual links strand them in silos.
- **Skipping tree testing** — card sorting reveals groupings but not findability; tree testing validates the final structure.

## When NOT to Use

- Single-page applications with fewer than 10 distinct content types
- Real-time tools (chat, games) where IA is irrelevant to interaction model
- When content volume is too small to require taxonomy (fewer than 20 items)
