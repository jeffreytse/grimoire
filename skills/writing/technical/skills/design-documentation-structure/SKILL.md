---
name: design-documentation-structure
description: Use when organizing or reorganizing documentation for a project, product, or knowledge base
source: "Diátaxis framework (Daniele Procida, 2017); Google Developer Documentation Style Guide"
tags: [documentation, information-architecture, diataxis, technical-writing, structure]
verified: true
---

# Design Documentation Structure

Organize documentation so readers always land in the right place for their current need.

## Why This Is Best Practice

**Adopted by:** Canonical (Ubuntu), Django, Gatsby, Cloudflare, and NumPy documentation teams

**Impact:** Canonical reported measurably reduced support requests after restructuring Ubuntu docs around Diátaxis; Django's docs are cited as a model for open-source project documentation

**Why best:** The Diátaxis framework distinguishes four fundamentally different reader needs — learning, doing, understanding, and looking up — each demanding a different writing style. Mixing these modes in one document forces readers to mentally filter content, increasing cognitive load and causing them to miss critical information.

## Steps

1. **Audit existing content** — list every page and classify it as tutorial, how-to guide, reference, or explanation; flag any page that tries to be two at once
2. **Map tutorials** — write learning-oriented content for complete beginners; focus on the journey and guarantee success, not completeness
3. **Map how-to guides** — write goal-oriented content for practitioners who know what they want; assume competence, focus on steps
4. **Map reference docs** — write information-oriented content structured around the system itself (API, CLI, config); accuracy and completeness over narrative
5. **Map explanation** — write understanding-oriented content exploring concepts, background, alternatives, and trade-offs; no procedures
6. **Design the navigation hierarchy** — surface all four quadrants from the top-level nav; do not bury tutorials under reference or vice versa
7. **Write an index page for each section** — one paragraph per section explaining what it contains and who it is for

## Rules

- Never mix tutorial and reference content in the same document
- How-to guides must be action-oriented: every heading should start with a verb
- Reference pages describe; they do not instruct or teach
- Explanation sections must not contain step-by-step procedures
- Every page must state its audience in the first paragraph

## Examples

**Misclassified:** A "Getting Started" page that also lists every configuration option — this is a tutorial mixed with reference.

**Correct:** Separate "Your First Deploy" (tutorial) from "Configuration Reference" (reference), linked from a clear top nav.

## Common Mistakes

- Creating an omnibus "Getting Started" page that tries to be all four types at once
- Hiding how-to guides inside tutorials, leaving intermediate users with no obvious path
- Writing explanation content inside reference pages, making them harder to scan

## When NOT to Use

- Do not apply this skill to single-page internal tools or scripts with fewer than five commands or options, where a single well-commented README section serves the reader better than a four-quadrant Diátaxis architecture.
- Do not use this skill when documentation is being written for a one-time deliverable such as a project handover or audit report, where the audience reads it once sequentially and navigational hierarchy adds overhead without benefit.
- Do not apply this skill to auto-generated API reference output from tools like Swagger UI or Javadoc, where the structure is dictated by the toolchain and manually imposing a custom information architecture breaks the generation pipeline.
