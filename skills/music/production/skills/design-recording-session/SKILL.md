---
name: design-recording-session
description: Use when planning and executing a professional recording session for any instrument or vocal, in a studio or home recording environment
source: Huber & Runstein "Modern Recording Techniques" (2017, 9th ed.); AES recording standards; Owsinski "The Recording Engineer's Handbook" (2017)
tags: [music-production, recording, studio-engineering, audio-engineering]
verified: true
---

# Design Recording Session

Plan and execute a professional recording session that captures high-quality audio efficiently while maintaining artist performance and technical standards.

## Why This Is Best Practice

**Adopted by:** AES (Audio Engineering Society, 12,000+ members) publishing standards for recording practice; Huber & Runstein is the most widely used recording techniques textbook in audio engineering programs (Full Sail, Berklee, SAE); Owsinski's Recording Engineer's Handbook is the professional reference for studio engineers.
**Impact:** Well-planned recording sessions achieve target sound in 30–50% fewer takes; proper gain staging and noise floor management prevent the 6dB signal-to-noise ratio degradation that occurs in 60%+ of home recordings; pre-session preparation reduces studio time costs by 25–40%.
**Why best:** Recording is the only phase of music production where the original signal is captured — errors in recording (noise, poor acoustics, clipping, phase issues) cannot be fully corrected in post-production. Planning prevents problems that are impossible to fix later.

Sources: Huber & Runstein "Modern Recording Techniques" 9th ed. (2017); Owsinski "The Recording Engineer's Handbook" 4th ed. (2017); AES standards for recording practice; Bartlett & Bartlett "Practical Recording Techniques" 7th ed. (2017).

## Steps

1. **Conduct pre-production planning** — before the session: confirm the track list (what will be recorded), decide on the arrangement (instruments, number of tracks), determine session format (live ensemble, overdub, hybrid), set the target tempo and key for all tracks, and create a session schedule with buffer time between takes.

2. **Prepare the recording environment** — treat the room acoustically: diffuse first reflections (absorb or scatter sound coming back from walls), control the bass buildup in corners (bass traps), ensure the recording space is isolated from external noise (air conditioning, traffic, HVAC). Test the noise floor: target −60 dBFS or lower for professional recording.

3. **Set up the signal chain** — establish the complete path: source (instrument/voice) → microphone → preamp → converter → DAW. Verify: phantom power (48V) for condenser mics, proper gain staging at each stage, cable integrity (no hum, buzz, or noise), and correct input routing in the DAW before the artist arrives.

4. **Apply gain staging correctly** — set levels to achieve: −18 dBFS average level (−18 dBu reference level, aligned to 0 VU on analog gear), −6 dBFS headroom for transient peaks. Never clip the converter (0 dBFS digital ceiling). A healthy signal at the preamp stage that clips the converter is unrecoverable.

5. **Select and position microphones** — choose mic based on source: large-diaphragm condenser (acoustic guitar, vocals, room ambience), dynamic (loud sources: drums, guitar amps), ribbon (smooth high-frequency response: brass, strings, guitar amps). Position: cardioid at 6–12 inches for close miking, experiment with distance for more room. Use the 3:1 rule for multi-mic setups to minimize phase cancellation.

6. **Record a reference level take** — before committing to tape (DAW), do a full run-through take to confirm: levels are set correctly, no clipping, no unwanted noise, headphone mix is comfortable for the performer, and the performance space is working acoustically. Adjust and then begin recording.

7. **Create a comfortable headphone (cue) mix** — the artist's headphone mix directly affects performance quality. Give performers: their own instrument/voice prominently, click track at a comfortable level, and relevant accompaniment tracks. Ask explicitly what they need; different performers have very different preferences. Use individual headphone mix capability when available.

8. **Record the performance** — use track comping workflow: record multiple complete takes rather than punching in on individual mistakes (especially for vocals and solo instruments). Create a comp (composite take) from the best sections after tracking. Keep all takes — you may need earlier material.

9. **Monitor phase coherence for multi-mic setups** — for drum kit or multi-mic acoustic recordings, check phase alignment: mono-sum the mix and listen for cancellation (loss of low frequencies signals phase problems). Use time alignment tools or physical mic repositioning to correct. Phase issues in recording cannot be fully corrected in mixing.

10. **Complete session documentation** — after recording, document: what was recorded, take numbers and notes (circle takes), file naming convention, tempo map and time signature, plugin settings used for any tracking effects, and storage location (local + backup). This information is essential for the mix engineer.

## Rules

- Record at 24-bit depth minimum; 32-bit float for complex sessions — bit depth determines dynamic range (24-bit = 144dB theoretical; far exceeds any acoustic environment).
- Set sample rate at 44.1kHz (delivery standard for music) or 48kHz (film/video standard); do not use higher rates unless mastering to require it (the difference is inaudible; the files are larger).
- Back up all recorded files immediately after each session — unrecoverable recording data loss is a professional failure.
- Never use audio compression as a tracking effect without also recording the dry signal in parallel — heavily compressed audio cannot be uncompressed.

## Common Mistakes

- **Gain staging errors** — recording too hot (clipping converter) or too low (noise floor dominates) are the most common and costly recording mistakes.
- **Ignoring headphone mix quality** — poor headphone mixes cause hesitant performances; investing 10 minutes in the cue mix often saves 30 minutes of additional takes.
- **Not acoustically treating the recording space** — reflective rooms add a "boxy" sound to close-mic recordings that cannot be removed in post.
- **Over-relying on "fix it in the mix"** — compression, EQ, and noise reduction can improve recordings but cannot fix fundamental acoustic problems (comb filtering from room reflections, excessive noise floor, off-axis mic placement).

## When NOT to Use

- When capturing live performance for archive rather than record production (different quality and time constraints apply).
- When sampling or working entirely with pre-recorded material (no tracking session needed; focus on arrangement and production).
- When producing entirely in the box with virtual instruments and no live recording (mix and production workflow applies instead).
