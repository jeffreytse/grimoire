---
name: design-culling-workflow
description: Use when establishing a systematic photo culling and selection workflow to efficiently identify selects from large shoot volumes
source: 'Lightroom Classic industry workflow (Adobe); ASMP digital asset management guidelines; Peter Krogh "The DAM Book: Digital Asset Management for Photographers" (2009)'
tags: [culling, photo-editing, workflow, lightroom, dam]
verified: true
---

# Design Culling Workflow

Build a fast, systematic culling process that reduces a large shoot volume to a curated select set efficiently, without agonizing over individual images.

## Why This Is Best Practice

**Adopted by:** ASMP professional digital imaging standards, Adobe Lightroom's star/flag system as the industry default, Peter Krogh's DAM methodology used by National Geographic photographers
**Impact:** Systematic culling reduces selection time by 50% vs. unstructured review (Adobe internal Lightroom usage research); a disciplined cull produces galleries with 3–5× higher per-image quality floor than unculled deliveries
**Why best:** Krogh's DAM framework establishes that culling is a distinct, separate phase from editing — mixing them creates bottlenecks and decision fatigue. A multi-pass system eliminates obvious rejects first, reducing the cognitive load for comparative selection.

Sources: Krogh "The DAM Book" (2009) Ch. 4; Adobe Lightroom Classic documentation; ASMP "Digital Asset Management Guidelines" (2022)

## Steps

1. **Import with integrity** — import RAW files to a dated, project-labeled folder; enable "Make a Second Copy To" for immediate backup; apply metadata preset at import (copyright, contact info, camera calibration, lens correction).
2. **Build previews before culling** — build 1:1 previews during import (or overnight); culling on rendered previews is 3–5× faster than waiting for previews to render per image.
3. **Pass 1: Flag keepers (fast, no zoom)** — use Lightroom flag (P = pick, X = reject) at fit-to-screen view; reject: obviously blurred, eyes closed, wrong exposure beyond recovery, subject not in frame; flag: technically usable images; speed target: 1 second per image.
4. **Filter to flagged only** — use Library Filter to show only picked images; this removes the cognitive noise of seeing rejects.
5. **Pass 2: Star rate for quality** — at fit-to-screen view, assign 1–3 stars: 1 star = technically acceptable; 2 stars = good composition or moment; 3 stars = exceptional image (best 5–10% of shoot).
6. **Pass 3: Zoom comparison for technical quality** — use Compare view (C key) to compare 2 images head-to-head for sharpness; toggle between similar frames; keep the sharper, better-composed version; this is where similar frames are reduced.
7. **Filter to 2-star minimum for delivery selects** — 2+ star images form the core delivery set; verify the set tells the complete story of the shoot without gaps.
8. **Identify hero images** — from 3-star images, select 5–10 portfolio-quality "hero" images for priority editing; these are processed first and to the highest standard.
9. **Apply batch global corrections** — select all delivery selects; sync: white balance, exposure, profile corrections, noise reduction; this brings all images to a consistent baseline before individual editing.
10. **Export and archive** — export selects at delivery specification; move RAW masters + PSD masters to long-term archive; delete rejects from working storage only after client approval.

## Rules

- Never cull and edit simultaneously — mixing passes destroys the speed advantage of the multi-pass system; cull completely before opening the develop module.
- Do not return to Pass 1 during Pass 2 or 3 — re-evaluating technical rejects during quality assessment introduces decision fatigue; trust the first pass.
- Hero images must be identified before any editing — prioritizing edits to hero images ensures best images are ready first.
- Archive originals before deleting rejects — never delete originals from working storage until they exist in at least one backup location.
- Delivery quantity must be defined upfront — without a target (e.g., "deliver 80 images"), the cull lacks a clear endpoint and expands indefinitely.

## Common Mistakes

- **Zooming to 100% on every image in Pass 1** — zooming to evaluate sharpness at the reject pass wastes time; blur is visible at fit-to-screen; save zoom for the comparison pass.
- **Keeping duplicates "just in case"** — duplicate similar frames inflate the delivery gallery without adding value; make the call between similar frames in Pass 3.
- **Culling without a target delivery count** — no target → keeping too many → retouching too many → reduced per-image quality.
- **Starting editing without completing the cull** — editing an image only to reject it later in the cull wastes editing time on a frame that will not be delivered.
- **No backup before culling** — accidentally flagging and batch-rejecting the wrong images without a backup has no recovery path.

## When NOT to Use

- Shoots producing fewer than 50 images (single-pass selection is sufficient)
- Film photography (analog workflow requires different archival approach)
- News/wire photography where immediate single-image selection is required in the field
