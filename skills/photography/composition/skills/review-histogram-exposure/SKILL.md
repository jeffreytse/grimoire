---
name: review-histogram-exposure
description: Use when evaluating a shot's exposure in-camera or in post-processing — to determine if highlights are blown, shadows are blocked, or exposure should be shifted before the scene changes.
source: 'Reichmann, Michael. "Exposing to the Right" (luminous-landscape.com, 2003); Cambridge in Colour, "Histograms" (cambridgeincolour.com)'
tags: [exposure, histogram, post-processing, raw, camera-settings]
verified: true
---

# Review Histogram for Exposure

Read the histogram rather than the LCD to make accurate exposure judgments and minimize noise in RAW files.

## Why This Is Best Practice

**Adopted by:** Universal in professional RAW workflow. Michael Reichmann's 2003 "Expose to the Right" article popularized histogram-based exposure for digital photography; the technique is now standard in Adobe Lightroom's documentation and all major RAW processors.

**Impact:** LCD screens lie — they look brighter or darker depending on ambient light. The histogram is always accurate. Photographers who check the histogram recover more shadow detail in post and avoid unrecoverable blown highlights on critical shots.

**Why best:** A RAW file stores more data in the brighter stops. Systematic histogram reading enables "ETTR" (expose to the right) — the practice of exposing as brightly as possible without clipping, then pulling down in post to maximize dynamic range and minimize noise.

## Steps

### Step 1: Enable the histogram in-camera

- Most cameras: **Info** or **Disp** button while reviewing an image cycles through overlay modes
- Enable **RGB histogram** (separate channels) if available — better than luminosity-only for detecting channel clipping
- Enable **highlight alert / "blinkies"** — overexposed areas flash in playback

### Step 2: Read the histogram shape

The x-axis = brightness (pure black left → pure white right). The y-axis = number of pixels at each brightness level.

```
Correct (normal scene):
     ▄▄▄
   ▄▄   ▄▄
  ▄        ▄
 ▄          ▄
|            |
0            255

Underexposed (left-heavy):
▄▄▄
▄   ▄
▄     ▄
▄       ▄
|            |
0            255

Overexposed / blown highlights (right clipped):
                 ▄▄▄▄▄▄
           ▄▄▄▄▄      |
      ▄▄▄▄▄             |
|                        |
0                       255  ← spike at right wall = clipped
```

### Step 3: Identify clipping

**Highlight clipping** (right wall): Data is permanently lost — pixels record as pure white with no color information. In RAW, slight clipping in 1–2 channels may be recoverable; all-3-channels blown is unrecoverable.

**Shadow clipping** (left wall): Pixels record as pure black. In RAW, shadows are recoverable up to 3–4 stops with noise penalty. Shadow clipping is less critical than highlight clipping.

Action on clipping:
- Right clipped on critical subject (skin, bright dress): reduce exposure by 1/3–1 stop, reshoot
- Left clipped on intentional dark background: acceptable
- Both walls clipped: scene has contrast beyond sensor range — use fill flash, HDR, or reframe

### Step 4: Apply ETTR (Expose to the Right)

ETTR = push exposure as far right as possible without clipping highlights. This technique maximizes signal-to-noise ratio in the RAW file because the brightest stop contains 50% of all digital values.

1. Shoot in RAW (not JPEG — JPEG clips before the RAW does)
2. Raise exposure until the histogram touches the right wall without overflowing it
3. In post, use the Exposure slider to bring brightness back to natural look
4. Result: noticeably cleaner shadows than if you had exposed for "correct" in-camera metering

**Do not apply ETTR for:** JPEGs, fast-moving scenes with no time to chimp, night sky (stars clip), deliberate high-contrast artistic choices.

### Step 5: Use RGB channels to detect color clipping

On an RGB histogram, if the Red channel clips while Green and Blue don't, you have blown reds (common with red flowers, warning lights). Adjust white balance or reduce exposure to recover.

Cameras that show only luminosity histogram can miss channel-specific clipping — a single-channel tool is less reliable.

## Rules

- Trust the histogram over the LCD in all lighting conditions
- Highlight clipping on the subject is a reshoot; clipping in unimportant sky or speculars is acceptable
- ETTR works only in RAW — never push JPEG exposure to the right
- One stop of underexposure in RAW adds ~2× the noise compared to correct exposure; ETTR reduces this
- Check the histogram immediately after the first frame of any new scene or lighting change

## Common Mistakes

**"Chimping" off the LCD:** The LCD misleads in bright sun (looks dark) and at night (looks bright). Histograms are environment-independent.

**Treating any right-leaning histogram as "overexposed":** A high-key scene (white dress on white background) naturally produces a right-heavy histogram. Context determines correct exposure, not a centered bell curve.

**Ignoring channel clipping:** A clean luminosity histogram can still have a blown red channel. Use RGB histograms or enable per-channel blinkies.

**Correcting underexposure in JPEG:** Pulling shadows up in a JPEG reveals banding and noise. The fix must happen before the shutter.

## When NOT to Use

- Studio with controlled, reproducible light: set exposure once with a grey card or incident meter, no need to chimp each frame
- Video production: waveform monitor is more precise than a histogram for exposure evaluation
