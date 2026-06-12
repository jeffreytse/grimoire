---
name: apply-eq-technique
description: Use when applying EQ to any audio track or bus to shape tone and reduce masking
source: David Gibson "The Art of Mixing" (2005), Mike Senior "Mixing Secrets for the Small Studio" (2011), Roey Izhaki "Mixing Audio" (2017)
tags: [music, production, eq, mixing, frequency, masking, subtractive-eq, additive-eq]
verified: true
---

# Apply EQ Technique

Apply EQ with intention — removing problem frequencies first, then adding character — to achieve clarity and separation in a mix.

## Why This Is Best Practice

**Adopted by:** Used by mixing engineers across all commercial genres; taught in every major audio engineering curriculum
**Impact:** David Gibson's frequency mapping in "The Art of Mixing" demonstrates that frequency masking between competing elements is the primary cause of muddy mixes; systematic EQ addresses this mechanically

**Why best:** EQ done wrong wastes headroom and creates masking. The standard professional approach is subtractive first (remove what hurts), additive second (add what helps). This order prevents boosting frequencies that are already crowded. Listening in context, not solo, reveals the actual masking problem rather than the track's isolated character.

## Steps

1. **High-pass filter first** — Apply a high-pass filter to every track that doesn't need sub-bass content. Typical HPF points: vocals 80-100Hz, guitars 100-120Hz, synths (context-dependent). Remove build-up that consumes headroom invisibly.
2. **Identify problem frequencies in context** — With the full mix playing, sweep a narrow boost (10-12dB, Q=8) slowly through the spectrum; the frequency that suddenly sounds harsh, muddy, or boxy is the problem frequency.
3. **Cut the problem frequency** — Switch to a moderate cut (-3 to -6dB) at the identified frequency with a moderate Q (3-5). The cut should fix the problem without creating a audible notch.
4. **Address frequency masking** — Identify competing elements in the same frequency range (e.g., kick and bass compete at 60-120Hz, vocals and guitars at 2-4kHz). Cut one to make room for the other rather than boosting both.
5. **Apply additive EQ for character** — Add air (12-16kHz shelf, +2-3dB) for presence; add warmth (200-300Hz, +2dB) if the track sounds thin. Use wide Q (1-2) for broad tone shaping.
6. **Check in mono** — Collapse the mix to mono and verify EQ decisions still work. Masking problems hidden by stereo separation become apparent in mono.
7. **A/B the EQ** — Bypass the entire EQ chain and compare to the processed version at the same loudness. EQ that doesn't improve the track in bypass comparison should be reconsidered.
8. **Check on reference speakers and headphones** — Verify EQ decisions translate across playback systems; a boost that sounds right on studio monitors may sound harsh on earbuds.

## Rules

- Never EQ in solo — masking problems only exist in the context of the full mix; solo EQ creates solutions to non-existent problems.
- Subtractive before additive — cuts remove real problems; boosts add character on top of existing problems.
- Cuts can be narrow (high Q); boosts should be wide (low Q) — narrow boosts sound unnatural and phasey.
- Never boost frequencies that another important element already occupies without first cutting that element.

## Examples

Bass guitar and kick drum masking: Kick has energy at 60Hz (punch) and 3-5kHz (click). Bass has energy at 80-100Hz (fundamental) and 250-400Hz (warmth). Solution: High-pass bass at 40Hz, cut bass at 60Hz by 4dB to give kick room; cut kick at 100Hz by 3dB to give bass fundamental space. Result: both elements audible and distinct simultaneously.

## Common Mistakes

- EQing in solo and then being surprised the mix doesn't improve — the masking relationship is the problem, not the track in isolation.
- Boosting muddy frequencies (200-400Hz) to add warmth — this frequency range is already overloaded in most mixes; cut rather than boost.
- Using the same EQ settings on every element of a type — every kick drum, bass, and vocal is different; EQ by ear, not by preset.
- Radical cuts (>12dB) as the first move — these create phase artifacts; identify the exact problem frequency precisely before cutting deeply.

## When NOT to Use

- When the recording has fundamental source problems (wrong microphone placement, room acoustics, instrument intonation) that EQ cannot fix and re-recording is the correct solution.
- When delivering unprocessed stems to a dedicated mix engineer who will perform all EQ decisions in the final mix context, not the production context.
- When working with already-mastered reference audio or sample packs where EQ changes will alter the material's calibrated tonal balance and licensing compliance.
