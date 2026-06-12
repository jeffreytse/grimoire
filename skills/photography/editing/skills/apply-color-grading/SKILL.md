---
name: apply-color-grading
description: Use when color grading photographs to achieve a consistent look, mood, or brand aesthetic
source: Adobe Color Methodology, VSCO Film emulation research, color theory application in photography (Itten "The Art of Color")
tags: [photography, editing, color-grading, hsl, tone-curve, lightroom, color-theory, mood]
verified: true
---

# Apply Color Grading

Apply color grading using HSL adjustments, tone curves, and color wheels to achieve a specific mood, consistency, or brand aesthetic.

## Why This Is Best Practice

**Adopted by:** VSCO's film emulation methodology is used by 200M+ users; Adobe's color grading tools are the industry standard; LACMA and major galleries use consistent color grading in documentary photography exhibitions
**Impact:** Color psychology research (Itten's color theory, validated by subsequent neuroaesthetics studies) shows that warm-shifted images are rated as more emotionally positive and trustworthy; consistent color grading increases brand recognition by 33% (Pantone/Forbes study)

**Why best:** Color grading transforms technically correct images into emotionally specific images. It is the difference between "accurate" and "intentional." Inconsistent color grading across a session or brand portfolio signals amateur work even when individual images are technically excellent. The craft is knowing which color relationships to reinforce or suppress to serve the emotional goal.

## Steps

1. **Establish the mood target** — Define the intended emotional quality: warm/intimate, cool/editorial, faded/nostalgic, high-contrast/dramatic. Every subsequent decision must serve this goal.
2. **Correct white balance first** — Color grading is meaningless on images with uncorrected color casts. Ensure WB is neutral before grading; all creative color shifts are added on top of a corrected base.
3. **Set the tone curve** — S-curve for contrast. Lift the blacks slightly (shadow fade/matte look) or crush them (deep blacks). Pull highlights down slightly for a film-like rolloff. The tone curve controls luminosity and sets the grading foundation.
4. **Adjust the HSL panel** — Hue shifts each color range; Saturation boosts or mutes; Luminance brightens or darkens. For skin tones: check orange and red hue and saturation first. Shift orange hue toward red (+10-15) to warm skin; reduce orange saturation slightly to avoid oversaturation.
5. **Apply the color grader (color wheels)** — Use Lightroom's Color Grading (or equivalent): add a warm tone (orange-yellow) to the Shadows wheel for a faded-film look; add a complementary cool (teal) to the Highlights wheel for a cinematic split tone. Midtones wheel sets the overall key.
6. **Check skin tones on a reference image** — Open a close portrait; verify skin tones look healthy not orange, gray, or green. Use the HSL orange and red channels to correct without affecting the overall grade.
7. **Assess global saturation** — Pull back global saturation if the grade feels oversaturated. Vibrance is preferable to Saturation for photographic images because it protects skin tones from oversaturation.
8. **Apply and sync to the session** — Apply the grade to the reference image, then sync to all images from the same lighting setup. Review the synced images as a slideshow to confirm consistency; adjust outliers individually.

## Rules

- Always grade from a corrected white balance — grading over a color cast produces inconsistent results across different shots with varying WB.
- Skin tones are the reference for any grade that includes people — if the grade makes skin look unhealthy, the grade must change.
- Vibrance, not Saturation, for portraits — Saturation clips skin tones; Vibrance protects them while boosting less-saturated colors.
- A consistent grade across a session is more important than a "perfect" grade on a single image — consistency is what makes a portfolio look professional.

## Examples

Travel/lifestyle brand grade: Tone curve — slight S-curve, shadows lifted to RGB value 15 (matte look), highlights rolled off at 230. HSL — aquas/teals: hue -10 (shift toward teal), saturation +15. Oranges: hue +8 (warmer skin), saturation -8. Color grading wheels — Shadows: orange-amber tint at 15% strength. Highlights: teal-blue at 10% strength. Result: warm, sun-kissed look with cinematic teal-orange split tone consistent with travel brand aesthetic.

## Common Mistakes

- Grading before correcting white balance — the grade builds on a flawed foundation; images with different WB will never look consistent.
- Over-saturating to make images "pop" — oversaturation looks artificial in print and on calibrated monitors; use Vibrance at +15-25 maximum.
- Applying the same grade preset to all image types — outdoor daylight and indoor mixed-light images require different base corrections before the same grade can be applied consistently.
- Never reviewing the grade on a calibrated display — grades applied on uncalibrated monitors produce images that look correct in the editing environment but wrong in delivery.

## When NOT to Use

- When delivering images for scientific, forensic, or medical documentation purposes where color accuracy must be preserved and creative interpretation would compromise evidentiary integrity.
- When the client brief specifies straight color-correct processing (e.g., product catalog photography requiring exact color matching to physical samples).
- When images are destined for press/offset printing with a strict ICC profile workflow managed by the print production team, where creative grading applied before handoff will be overridden by print calibration.
