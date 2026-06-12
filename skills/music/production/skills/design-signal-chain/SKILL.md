---
name: design-signal-chain
description: Use when routing audio through processing plugins or hardware and need to determine the optimal processing order
source: 'Bob Katz "Mastering Audio: The Art and the Science" (2002), Mike Senior "Mixing Secrets for the Small Studio", Universal Audio signal chain methodology'
tags: [music, production, signal-chain, gain-staging, plugin-order, mastering, routing]
verified: true
---

# Design Signal Chain

Determine the correct order of processing tools in an audio signal chain to maximize headroom and achieve the intended sonic result.

## Why This Is Best Practice

**Adopted by:** Bob Katz's gain staging approach is used universally in mastering; UA Apollo console design reflects decades of hardware signal chain research
**Impact:** Bob Katz established the K-System metering standard adopted by the AES (Audio Engineering Society) precisely because signal chain gain staging errors are the most common cause of quality loss in digital audio

**Why best:** Processing order is not arbitrary — it determines what each subsequent processor "hears." A compressor after a distortion pedal hears a different signal than a compressor before it. Incorrect order causes the wrong problems to be compressed, the wrong frequencies to be EQ'd, and saturation to occur at the wrong stage. Getting the chain right once prevents hours of remedial work.

## Steps

1. **Start with gain staging** — Set input gain so the signal hits -18 dBFS RMS before any processing. Every subsequent processor inherits this headroom.
2. **Apply noise gates or expanders first** — Remove noise before any processing that would otherwise amplify it (compression, saturation). Gate threshold just above the noise floor.
3. **Apply corrective EQ second** — High-pass filter and corrective cuts before compression; compressor should hear the corrected tone, not the problem frequencies.
4. **Apply compression third** — Compressor now controls dynamics of the tonally corrected signal. Ratio, attack, and release tuned to the musical content.
5. **Apply saturation or harmonic enhancement fourth** — After dynamics are controlled, saturation adds consistent harmonic content without pumping from dynamic peaks.
6. **Apply tonal EQ (additive) fifth** — Final tone shaping after dynamics and saturation are set. This EQ shapes the character of the processed signal.
7. **Apply time-based effects last** — Reverb and delay are always last in the chain; sending a compressed/EQ'd/saturated signal to reverb is correct. Reversing this smears all processing into the reverb tail.
8. **Check gain at each stage** — Insert a metering plugin after each processor; output gain should not exceed -6 dBFS peak at any intermediate stage.

## Rules

- Gate before compress — compressor pumping on noise is irreversible without re-recording.
- Corrective EQ before compression — compressor responds to frequency imbalances; fix the tone first.
- Never put reverb before compression on a bus — the reverb tail will be compressed, killing depth perception.
- Saturation after compression produces consistent harmonic content; saturation before compression produces variable harmonic content that may compress unpredictably.

## Examples

Vocal chain: High-pass (100Hz, 12dB/oct) → De-esser (6-9kHz dynamic cut) → Compressor (4:1, fast attack, auto-release) → Saturation (subtle tape emulation, +0.5dB drive) → Additive EQ (air shelf +2dB at 12kHz, presence +1.5dB at 3kHz) → Send to Reverb bus. Gain at output: peaks -8 dBFS, RMS -14 dBFS.

## Common Mistakes

- Saturation before compression — produces unpredictable, variable distortion as dynamics hit the saturation stage at different levels.
- EQ boost before compression — boosted frequencies hit the compressor harder and trigger it more aggressively than intended.
- Reverb on the input of a bus chain — every processor after the reverb processes the reverb tail, destroying depth and muddying the mix.
- No gain check between stages — undetected clipping at an intermediate stage produces audible distortion that downstream processors cannot fix.

## When NOT to Use

- When deliberately designing a lo-fi or distortion-based aesthetic where conventional gain staging and processing order are intentionally violated for the desired sonic character.
- When working entirely within a hardware analog chain where plugin insert order is not applicable and physical signal flow is fixed by the console or patch bay routing.
- When the session is a live recording capture with no processing — designing a signal chain before the performance context is established adds constraints before any sonic problems are identified.
