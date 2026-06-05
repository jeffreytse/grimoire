---
name: design-photo-retouching-workflow
description: Use when establishing a systematic photo retouching workflow for portrait, product, or commercial image post-production
source: NAPP (National Association of Photoshop Professionals) retouching standards; Scott Kelby "The Adobe Photoshop Book for Digital Photographers" (2020); ASMP digital imaging guidelines
tags: [retouching, photoshop, post-production, editing, workflow]
verified: true
---

# Design Photo Retouching Workflow

Build a systematic, non-destructive retouching workflow that produces consistent, commercially acceptable results efficiently across a volume of images.

## Why This Is Best Practice

**Adopted by:** NAPP Photoshop certification curriculum, ASMP professional digital imaging standards, all major commercial photography studios
**Impact:** Non-destructive workflows reduce re-do time by 60% when client revision requests arrive; systematic retouching reduces per-image time by 40% through habit and layer organization (NAPP efficiency studies)
**Why best:** Kelby's layered retouching approach and ASMP non-destructive principles ensure that every decision remains reversible, the original file is never altered, and the retouch can be dialed back or handed to another retoucher without rebuilding from scratch.

Sources: Kelby "The Adobe Photoshop Book for Digital Photographers" (2020); ASMP "Digital Imaging Guidelines" (2022); NAPP "Photoshop Retouching Certification Standards"

## Steps

1. **Organize source files** — import RAW files into Lightroom Classic or Capture One; apply camera calibration and lens correction profiles; back up originals to minimum two locations before retouching begins.
2. **Global corrections first (Lightroom/Capture One)** — correct: white balance, exposure, highlights, shadows, whites, blacks; apply noise reduction and sharpening at global level; these global corrections apply to all selected images simultaneously via sync.
3. **Export layered master to Photoshop** — open as Smart Object to preserve RAW edit re-access; duplicate background layer immediately (Cmd/Ctrl+J); never retouch on the background layer.
4. **Frequency separation for skin** — create high-frequency (texture) and low-frequency (tone/color) layers; retouch skin color and tone on low-frequency without altering pores and texture; heal texture irregularities on high-frequency layer independently.
5. **Dodge and burn for dimension** — create two 50% gray layers (mode: Overlay); paint white (dodge) on areas to brighten, black (burn) to deepen; this shapes and sculpts the face or product surface without color contamination.
6. **Liquify for shaping** — use Photoshop Liquify as a Smart Filter on a duplicate merged layer; make only subtle adjustments; before/after comparison at 50% view is the quality check.
7. **Color and tone refinement** — use Curves adjustment layer for global contrast; Hue/Saturation for selective color correction (skin, sky, product color); add a graduated or radial filter adjustment for vignette and spotlight.
8. **Clean up distractions** — remove: sensor dust (heal on empty layer), background distractions (clone, patch tool), unwanted objects; work on separate empty layers above source.
9. **Sharpen for output** — add output sharpening last as a Smart Filter (Unsharp Mask or High Pass); sharpening amount depends on output: web (moderate), large print (aggressive), small print (light).
10. **Flatten and export** — save layered PSD as master file; flatten a copy for export; export JPEG at 100% quality for print delivery or 80% for web; verify exported file dimensions and color profile match client specifications.

## Rules

- Never retouch on the background/original layer — every destructive edit must be on a duplicate or empty layer above.
- Save the layered PSD — the JPEG export is a delivery artifact; the PSD is the asset; never overwrite the PSD with a flatten.
- Frequency separation must be used for skin — healing directly on a texture layer creates blurring artifacts that look worse than the original skin.
- Retouching must remain believable — over-retouched skin (no pores, no texture, plastic appearance) is a quality failure in commercial and beauty standards; use the "blink test" (before/after blink) to check for over-processing.
- All adjustments must be delivered at the specified color space — CMYK for print, sRGB for web; wrong color space is a professional error.

## Common Mistakes

- **Healing directly on the background layer** — any healing on the background layer makes the change permanent and irreversible.
- **Sharpening before global corrections** — sharpening applied before tone correction is compounded by later contrast adjustments, creating oversharpened artifacts.
- **Over-smoothing skin** — aggressive frequency separation with heavy Gaussian blur produces skin that reads as CGI rather than photography.
- **Not calibrating the display** — retouching on an uncalibrated monitor produces files that look wrong on calibrated displays used by print vendors.
- **Forgetting to check at 100% zoom** — dust, healing artifacts, and edge contamination are only visible at 100% (1:1 pixel view); always inspect before delivery.

## When NOT to Use

- Documentary or photojournalism work where image manipulation is ethically prohibited
- Casual social photography where professional retouching is out of scope
- RAW-to-JPEG workflows for speed-priority delivery (sports, news) where retouching is impractical
