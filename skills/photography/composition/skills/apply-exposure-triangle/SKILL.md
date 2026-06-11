---
name: apply-exposure-triangle
description: Use when setting camera exposure for any shot — to understand how aperture, shutter speed, and ISO interact so you can make deliberate creative choices rather than relying on full auto.
source: 'Peterson, Bryan. "Understanding Exposure" (4th ed., 2016); Cambridge in Colour, "Camera Exposure" (cambridgeincolour.com)'
tags: [exposure, aperture, shutter-speed, iso, camera-settings]
verified: true
---

# Apply the Exposure Triangle

Control aperture, shutter speed, and ISO as a unified system to achieve correct exposure while serving your creative intent.

## Why This Is Best Practice

**Adopted by:** Every photography education curriculum — from NYIP to CreativeLive to university programs. Bryan Peterson's *Understanding Exposure* has sold over 1 million copies and is the canonical self-teaching reference.

**Impact:** Photographers who understand the exposure triangle shift from reacting to light to controlling it. Full-auto mode optimizes for average exposure; manual control lets you choose shallow depth of field, freeze motion, or retain color in low light — outcomes auto cannot predict.

**Why best:** The three variables are coupled: changing one forces a trade-off in the others. Understanding the coupling is what separates creative exposure decisions from guesswork.

## Steps

### Step 1: Understand what each variable controls

**Aperture (f-stop)** — controls how wide the lens opening is.
- Wide aperture (f/1.8, f/2.8): more light in, shallow depth of field (background blur)
- Narrow aperture (f/11, f/16): less light in, deep depth of field (everything sharp)
- Counter-intuitive: smaller f-number = wider opening

**Shutter speed** — controls how long the sensor is exposed.
- Fast (1/1000s, 1/500s): freezes motion, less light
- Slow (1/30s, 1s): blurs motion, more light; requires tripod below ~1/focal-length

**ISO** — controls sensor sensitivity.
- Low ISO (100, 200): cleanest image, needs more light
- High ISO (3200, 6400+): usable in darkness, introduces grain/noise

### Step 2: Choose your priority based on creative intent

| Scene | Priority | Typical settings |
|-------|----------|-----------------|
| Portrait with blurred background | Aperture (wide) | f/1.8–f/2.8, adjust shutter/ISO |
| Sports, wildlife in motion | Shutter (fast) | 1/500s–1/2000s, adjust aperture/ISO |
| Landscape with full depth of field | Aperture (narrow) | f/8–f/16, tripod, low ISO |
| Handheld in low light | ISO first | Raise ISO until shutter fast enough |
| Silky waterfall/light trails | Shutter (slow) | 1s–30s, f/8–f/11, ISO 100, tripod |

Use Aperture Priority (Av), Shutter Priority (Tv), or Manual (M) mode — avoid Program (P) and Auto.

### Step 3: Apply the reciprocal relationship

Each stop of change doubles or halves exposure. To keep exposure constant while changing one variable, compensate in another:

```
Equivalent exposures (same brightness):
f/2.8  | 1/500s | ISO 400
f/4    | 1/250s | ISO 400   ← stopped down aperture, halved shutter speed
f/4    | 1/500s | ISO 800   ← stopped down aperture, doubled ISO instead
f/2.8  | 1/1000s| ISO 800   ← faster shutter, raised ISO to compensate
```

### Step 4: Use the metering system as a starting point

1. Point at the scene and half-press shutter to read meter
2. Adjust settings until meter reads 0 (or intentionally over/under for creative reasons)
3. Check histogram (see `review-histogram-exposure`) — not the LCD preview, which lies in bright light
4. Adjust if highlights are blown or shadows are blocked

### Step 5: Set a base exposure workflow

For most situations:
1. Set ISO to base (100 or 200) unless in low light
2. Set aperture for your desired depth of field
3. Adjust shutter speed until meter reads correctly
4. If shutter speed goes below safe handheld limit (1/focal-length rule), raise ISO

## Rules

- Never trust the LCD in bright sunlight — use the histogram
- Shutter speed minimum for handheld: 1/focal-length (e.g., 1/50s on a 50mm lens)
- ISO is the last resort: raise it when aperture and shutter can't get you the shot
- One stop = double or half the light — applies equally to all three variables
- f/8 is the sweet spot for most lenses: sharpest, good depth of field, no diffraction

## Common Mistakes

**Staying in Auto:** Camera meters for middle grey — silhouettes, snow, and backlit subjects all fool it. Learn manual or semi-manual.

**Raising ISO first:** Noise is permanent. Try widening aperture or slowing shutter first; raise ISO only when those options are exhausted.

**Ignoring depth of field:** Shooting environmental portraits at f/2.8 blurs the background so aggressively that context disappears. f/5.6 often serves the story better.

**Forgetting the reciprocal:** Stopping down aperture for sharpness without compensating produces an underexposed frame.

## When NOT to Use

- Documentary run-and-gun where speed matters more than precision: use Auto ISO with a minimum shutter speed set in camera
- Studio shooting with strobes: sync speed (~1/200s) constrains shutter; control exposure through aperture and strobe power only
