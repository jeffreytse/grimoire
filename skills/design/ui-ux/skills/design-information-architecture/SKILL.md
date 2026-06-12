---
name: design-information-architecture
description: Use when organizing content, navigation, or features in a product so users can find what they need without confusion
source: Peter Morville & Louis Rosenfeld "Information Architecture" (4th ed., 2015); Nielsen Norman Group IA Guidelines
tags: [information-architecture, navigation, ux, card-sorting, tree-testing, content-strategy]
verified: true
---

# Design Information Architecture

Organize product content and navigation so users can find information quickly using mental models, not memorization.

## Why This Is Best Practice

**Adopted by:** NN Group IA practice, GOV.UK digital service standard, Shopify design system team
**Impact:** Nielsen Norman Group studies show that users abandon tasks 40% more often due to poor navigation than due to visual design problems; tree testing can identify IA problems before any visual design is built

**Why best:** Information architecture fails invisibly — users blame themselves when they can't find things, not the product. Structural problems discovered in tree testing cost 10× less to fix than those found post-launch.

## Steps

1. **Inventory the content** — List every piece of content, feature, or page that needs a home. Group into a spreadsheet: content item, format, audience, frequency of access.
2. **Research user mental models** — Conduct 5–8 user interviews asking: "When you think about [domain], what categories or groups come to mind?" Look for natural language and groupings users already use.
3. **Run a card sort** — Give participants 30–50 cards (one content item each) and ask them to group them and name each group. Use open card sort for discovery, closed card sort to validate existing categories. Optimal Workshop or Maze work well remotely.
4. **Analyze card sort results** — Use a similarity matrix to identify content items users consistently group together and items with high disagreement. High-agreement clusters become navigation categories.
5. **Draft the IA structure** — Build a hierarchy: top-level categories → sub-categories → content items. Aim for breadth over depth (7±2 categories, no more than 3 levels deep for most products).
6. **Run a tree test** — Present the hierarchy as text only (no visual design) and give participants realistic tasks: "Where would you go to update your billing address?" Measure success rate and time per task.
7. **Iterate on failures** — Any task with <70% success rate indicates an IA problem. Rename labels, re-categorize items, or split/merge categories. Repeat tree testing.

## Rules

- Navigation labels must use users' language, not product or organizational jargon.
- Avoid "Other" or "Miscellaneous" categories — they are cognitive trash cans that hide things.
- No content item should live in more than one primary location; use cross-links for secondary access.
- Tree test before doing any visual navigation design — separating structure from aesthetics produces cleaner tests.
- Validate with a second round of tree testing after each major revision.

## Examples

GOV.UK redesigned its navigation from department-based categories ("Home Office," "HMRC") to task-based categories ("Visas and immigration," "Tax") after card sorting revealed citizens couldn't map government departments to their actual needs. Task completion rates rose from 43% to 76%.

## Common Mistakes

- **Org-chart IA** — Structuring navigation by internal department makes sense internally but fails for external users who don't know how the organization is structured.
- **Skipping tree testing** — Designing navigation visually first embeds the IA into mockups before it's been validated; teams then resist changes due to sunk cost.
- **Too many levels deep** — Users lose context below three levels; they forget where they are and cannot backtrack effectively.

## When NOT to Use

- Do not apply a full IA design process to single-purpose tools with fewer than 10 distinct content items — the overhead of card sorting and tree testing exceeds the complexity of the problem.
- Do not redesign information architecture when poor findability is caused by search quality or content gaps rather than structural organization — fixing navigation will not help users who cannot articulate what they are looking for.
- Do not run card sorts or tree tests as a substitute for understanding actual user goals — if the content inventory itself is wrong or incomplete, structural validation produces a well-organized system for the wrong content.
