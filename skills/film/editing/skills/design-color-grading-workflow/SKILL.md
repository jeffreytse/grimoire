---
name: design-color-grading-workflow
description: Use when establishing a color grading workflow for a film project that needs a consistent visual look from on-set monitoring through final deliverables
source: SMPTE color standards; ASC CDL (Color Decision List) standard; Blackmagic Design DaVinci Resolve certification curriculum
tags: [color-grading, post-production, color-science, workflow, DaVinci-Resolve, deliverables]
verified: true
---

# Design Color Grading Workflow

Establish a technically sound, creatively consistent color grading workflow from on-set color management through final deliverables across multiple platforms.

## Why This Is Best Practice

**Adopted by:** SMPTE (Society of Motion Picture and Television Engineers), ASC Technology Committee (creators of the CDL standard), Netflix, Disney, and all major studio post-production pipelines.
**Impact:** Productions using a defined color workflow with ASC CDL from day one report 60% fewer color corrections required at the conform stage. Netflix's post-production specification requires an ACES or defined color-managed workflow as a condition of delivery acceptance.
**Why best:** Color grading without a defined workflow produces inconsistent looks between camera units, locations, and colorists. A structured workflow using color science standards (ASC CDL, ACES, or camera-native Log) ensures that creative intent established on set is preserved through edit, VFX, and final grade without accumulated error.

Sources: SMPTE ST 2067 (IMF), ST 2065 (ACES); ASC CDL specification v1.2; Blackmagic Design DaVinci Resolve certification curriculum; Alexis Van Hurkman "Color Correction Handbook" (2014)

## Steps

1. **Define the color science pipeline before the shoot** — Choose and document the color science approach: ACES (Academy Color Encoding System) for maximum deliverable flexibility, camera-native Log with a defined IDT/ODT, or a custom LUT-based workflow. The choice must be made before the DIT begins work on day one. Document the choice in a color pipeline document distributed to the DP, DIT, colorist, and VFX supervisor.

2. **Establish on-set color management with the DIT** — The DIT (Digital Imaging Technician) applies a technical display LUT to convert Log footage to a viewable image for the director and DP's monitor. This LUT must be documented and provided to the colorist — it defines the "on-set look" baseline. Confirm that the DIT's LUT is a technical viewing LUT, not a creative grade applied to the files.

3. **Apply ASC CDL values for on-set creative decisions** — When the director and DP make creative color decisions on set (a warmer look for a sunset scene, cooler shadows for a thriller sequence), record these as ASC CDL values: Slope, Offset, Power (SOP) per RGB channel, plus Saturation. CDL values are camera-agnostic and transferable to any grading system without quality loss.

4. **Organize the project structure in DaVinci Resolve** — Create the project with the correct timeline color space and gamma settings matching the pipeline document. Organize bins by shoot day and camera. Apply the technical IDT (Input Device Transform) to all clips in the media pool before the edit begins. Verify the signal chain with scopes (waveform, vectorscope, parade) against a known reference frame.

5. **Conduct a primary grade pass before secondary work** — Establish the primary grade (lift/gamma/gain, color balance, overall exposure) for every scene before any secondary work (qualifications, masks, windows). The primary pass creates the neutral baseline from which the look is built. Primaries must be measured against scopes, not by eye — the viewing monitor is an approximation.

6. **Build the look using a node-based structure** — Use a consistent node structure for every clip: (1) Input transform node, (2) Technical corrections node (exposure, white balance normalization), (3) Creative grade node (the look), (4) Output transform node. Label every node. This structure makes the grade auditable, transferable between colorists, and reversible at any stage.

7. **Create LUT exports for VFX and editorial reference** — Export a 33-point 3D LUT from the technical corrections + creative grade nodes for use by VFX artists and the editorial team. VFX plates must be processed through the same pipeline before integration — un-graded VFX integrated into a graded timeline is the most common source of visual inconsistency in post.

8. **Deliver against platform-specific specifications** — Each delivery platform has specific requirements: Theatrical DCP (P3-D65, 48 nits, SMPTE ST 428); Netflix (Rec.2020, HDR10, 1000 nits peak, -23 LUFS audio); broadcast (Rec.709, 100 nits, SMPTE C). Grade and QC each deliverable separately against its specification. Do not derive one from another without a documented transform.

## Rules

- Never apply a creative grade directly to camera originals — always work non-destructively on a copy or use a node-based system where the original is preserved as input.
- All ASC CDL values must be exported as a `.cdl` or `.edl` file and archived with the project — these values are the creative record of on-set decisions.
- The colorist must work on a calibrated monitor (minimum: DCI P3, D65 white point, 100 nits reference for Rec.709 deliverables); consumer monitors cannot evaluate color for deliverable compliance.
- VFX plates must be graded and QC'd before integration into the timeline, not after.
- Deliver a separate LUT package to the streaming platform or distributor as part of the deliverables package — it documents the creative intent for re-versioning.

## Common Mistakes

- **No color pipeline document distributed before the shoot** — DIT and colorist work in different color spaces; the grade built on set does not transfer to the DI suite without quality loss.
- **Creative grades applied as baked LUTs to camera originals** — Baking a LUT into the camera file destroys the dynamic range advantage of Log recording and makes any subsequent correction destructive.
- **Single deliverable grade stretched to all platforms** — A Rec.709 grade used for theatrical delivery and then converted to HDR for streaming produces incorrect HDR; each platform requires its own targeted grade.
- **No node labels or structure documentation** — An undocumented node tree cannot be handed to a second colorist or revisited six months later without reverse engineering.

## When NOT to Use

- Single-camera shoots with no post pipeline where a simple LUT applied in camera is sufficient for the delivery format.
- News and current affairs productions where turnaround prohibits color science workflow implementation.
- Archive restoration projects, which follow a separate color science methodology governed by the source material's original color space.
