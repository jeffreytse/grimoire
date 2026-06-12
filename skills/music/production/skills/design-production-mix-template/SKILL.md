---
name: design-production-mix-template
description: Use when setting up a DAW session template for mixing or starting a new production project
source: Mike Senior "Mixing Secrets for the Small Studio" (2011), Pensado's Place mixing workflows, Pro Tools/Ableton session organization standards
tags: [music, production, mixing, daw, template, gain-staging, bus-routing, organization]
verified: true
---

# Design Mix Template

Build a reusable DAW session template that enforces correct gain staging, bus architecture, and track organization before a single audio file is loaded.

## Why This Is Best Practice

**Adopted by:** Professional mixing engineers including Michael Brauer (multi-bus technique), Chris Lord-Alge, and Tchad Blake publish their template approaches
**Impact:** Mike Senior's research in "Mixing Secrets for the Small Studio" shows that session organization and gain staging decisions made before mixing account for 60% of final mix quality in small-room contexts

**Why best:** A template eliminates setup time and enforces discipline. Correct gain staging from the start prevents headroom loss that forces remedial compression later. Bus routing enables parallel processing and recall-able bus compression — the core of the "mix bus glue" approach used on virtually every commercial release.

## Steps

1. **Define the track category groups** — Create groups for: Drums, Bass, Guitars, Keys/Synths, Vocals (Lead + BG), FX/Samples, and Orchestral/Other. Each group feeds a corresponding bus.
2. **Create the bus architecture** — Each group feeds a stereo sub-bus. All sub-buses feed a master bus. Optionally add a Parallel Compression bus and a Reverb/Delay return bus.
3. **Set gain staging targets** — Individual tracks: average -18 to -12 dBFS RMS with peaks no higher than -6 dBFS. Sub-buses: -10 to -6 dBFS. Master bus: at least 6dB headroom before limiting.
4. **Add utility processing to each bus** — Insert a gain plugin (for trim without moving fader), a spectrum analyzer, and a placeholder for bus compression (bypassed until needed).
5. **Configure the master bus chain** — EQ (high-pass at 20Hz, gentle air shelf) → bus compressor (ratio 2:1, slow attack/release) → limiter (ceiling -0.3 dBFS, bypassed until mastering prep). Leave all bypassed by default.
6. **Create track color coding** — Assign consistent colors: drums = blue, bass = green, guitars = orange, keys = yellow, vocals = red, FX = purple. Color coding survives session transfers.
7. **Add reference track channel** — Import a commercially mastered reference track. Set its volume to match the mix bus level. Use for periodic A/B comparison.
8. **Save as a locked template** — Export the session as a template (Logic: Template; Ableton: .als default set; Pro Tools: session template). Never modify the template in-session — use Save As.

## Rules

- Never place processing on the master bus until the mix is 80% complete — premature master bus limiting hides dynamic problems in individual tracks.
- All tracks must be gain-staged before any EQ or compression is applied — EQ into clipping is not EQ, it's distortion.
- The reference track channel must be loudness-matched to the mix, not raw file volume.
- Template sub-buses must be stereo even if the source is mono — summing mono into stereo is reversible; the reverse is not.

## Examples

Pop/R&B template: 8 drum tracks (kick, snare, hats, OH L/R, room L/R, perc) → Drums Bus. 2 bass tracks (DI + amp) → Bass Bus. 4 guitar tracks → Guitar Bus. 6 synth tracks → Keys Bus. 8 vocal tracks (lead, doubles, harmonies, BGs) → Vocal Bus. All buses → Master Bus. Parallel Compression bus fed by sends from Drums and Bass buses. Total: 28 tracks pre-loaded, color-coded, gain-staged at -18 dBFS RMS.

## Common Mistakes

- No bus architecture — mixing directly to the master bus loses the ability to control section balance independently.
- Gain staging after processing — EQ and compression behavior changes dramatically at different input levels; stage gain first.
- Using absolute fader positions as gain reference — faders should be at unity (0 dB) after gain staging; track levels set with gain trim plugin.
- Skipping the reference track — ears adapt to the mix; periodic A/B against a professional reference catches tonal drift.

## When NOT to Use

- When working on a one-off mix for a single-instrument recording (e.g., solo piano or voice) where a full bus architecture adds complexity without any corresponding benefit.
- When the DAW session has already been built by a tracking engineer with an established routing structure that would require destructive reorganization to replace.
- When the deliverable is stems or individual rendered tracks for a sync library, where the final mix will be assembled externally and no master bus chain is needed.
