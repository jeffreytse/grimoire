---
name: apply-typography-hierarchy
description: Use when designing a typographic system for a document, website, or brand that needs clear visual hierarchy and readable body text
source: Robert Bringhurst "The Elements of Typographic Style" (4th ed., 2012); Ellen Lupton "Thinking with Type" (2010)
tags: [typography, type-hierarchy, visual-design, layout, readability, type-scale]
verified: true
---

# Apply Typography Hierarchy

Build a type system with a coherent scale, contrast, and rhythm so readers navigate content without conscious effort.

## Why This Is Best Practice

**Adopted by:** Every major editorial, book, and web design discipline; Google Material Design type system; Apple Human Interface Guidelines typography specifications
**Impact:** Bringhurst's principles underpin professional typography in 50+ countries; eye-tracking research (Nielsen Norman Group, 2018) confirms users scan pages in F-patterns — visual hierarchy guides the eye through deliberate size and weight contrast

**Why best:** Good typography is invisible to the reader. Bad typography is a tax on every reading moment. A coherent hierarchy reduces cognitive load by encoding priority in visual weight before the reader processes a word.

## Steps

1. **Choose a type scale** — Use a modular scale with a ratio: Major Third (1.25), Perfect Fourth (1.333), or Golden Ratio (1.618). Apply consistently: base size × ratio^n. Tools: typescale.com. For web: set base at 16–18px.
2. **Select 2 typefaces maximum** — One serif or humanist sans-serif for body text; one contrasting face for headings. Pair by contrast (weight, width, or classification), not similarity. Avoid combining two serifs or two geometric sans-serifs.
3. **Set body text first** — Body text is the foundation. Set size, leading (line-height), and measure (line length) before touching headings. Optimal measure: 45–75 characters per line. Leading: 120–145% of font size.
4. **Establish heading levels** — Use the type scale to assign sizes for H1–H4. Each level should be visually distinct at a glance — at minimum a 1.25× size difference between adjacent levels. Supplement size with weight change (regular vs. bold) or color.
5. **Control spacing with a baseline grid** — All vertical spacing (margins, padding, line height) should align to an 8px or 4px grid. Consistent vertical rhythm makes the page feel ordered even subconsciously.
6. **Set color and contrast** — Body text: minimum 4.5:1 contrast ratio (WCAG AA). Headings at 18pt+: minimum 3:1. Decorative text labels: at least 3:1. Never gray-on-gray body text.
7. **Define hierarchy through weight, size, color, and case** — Use no more than two of these simultaneously per level. Three simultaneous signals (bold + large + all-caps) creates noise, not hierarchy.
8. **Test in context** — Apply the system to a real content sample. The hierarchy should communicate without labels: a reader should be able to rank content importance by visual weight alone.

## Rules

- Never set body text below 16px on screen or 10pt in print.
- Line length must not exceed 75 characters for sustained reading — wider lines cause readers to lose their place.
- Avoid justified text on screen without hyphenation; it creates irregular word spacing (rivers) that disrupts reading flow.
- All-caps is acceptable for short labels (3–5 words); never for body text or headings longer than 6 words.
- Font pairing must create contrast — if you squint and the two faces look the same, change one.

## Examples

A publication system using Inter (body) and Playfair Display (headings):
- H1: 48px / 700 weight / 1.1 leading
- H2: 36px / 600 weight / 1.2 leading
- H3: 24px / 600 weight / 1.3 leading
- Body: 18px / 400 weight / 1.6 leading / max 680px container width
- Caption: 14px / 400 weight / muted color

Consistent 8px grid spacing between all elements.

## Common Mistakes

- **Too many typefaces** — Three or more faces compete for attention; readers experience noise, not hierarchy.
- **Ignoring line length** — Full-width body text on a 1440px screen (150+ characters/line) is nearly unreadable; readers lose their place and skip lines.
- **Size without weight contrast** — Headings that are only slightly larger than body text fail to communicate hierarchy; add weight (bold) or letter-spacing to reinforce distinction.
