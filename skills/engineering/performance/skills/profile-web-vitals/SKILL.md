---
name: profile-web-vitals
description: Use when diagnosing web performance problems or preparing a performance audit for a website or web application
source: "Google web.dev Core Web Vitals; Chrome DevTools Performance panel documentation; CrUX (Chrome User Experience Report)"
tags: [performance, web-vitals, lcp, cls, inp, chrome-devtools, lighthouse, profiling]
verified: true
---

# Profile Web Vitals

Measure LCP, CLS, and INP on real user data, then trace each metric to its root cause using DevTools before optimizing.

## Why This Is Best Practice

**Adopted by:** Google uses Core Web Vitals as a ranking signal in Google Search (confirmed May 2021); Chrome, Edge, and Firefox report CWV metrics; Cloudflare, Vercel, and Netlify expose CWV dashboards in their analytics products

**Impact:** Google's CrUX data shows pages passing all Core Web Vitals thresholds have 24% fewer abandonment events; BBC found a 10% improvement in LCP correlated with a 1% increase in additional page views

**Why best:** Core Web Vitals measure what users experience, not what the server reports. A fast server response with a slow LCP means users still perceive the page as slow. Lab measurements (Lighthouse) show potential; field measurements (CrUX, RUM) show reality. Always diagnose with field data first.

## Steps

1. **Pull field data first** — check CrUX via PageSpeed Insights (pagespeed.web.dev) or Search Console Core Web Vitals report; look at p75 values for LCP, CLS, and INP segmented by device type (mobile typically 2-3× worse)
2. **Audit LCP (Largest Contentful Paint)** — identify the LCP element using DevTools Performance panel → LCP marker; diagnose render-blocking resources, slow server response (TTFB), and unoptimized LCP image (no `fetchpriority="high"`, no preload)
3. **Audit CLS (Cumulative Layout Shift)** — record a DevTools performance trace with screenshots; find layout shift clusters; common causes: images without explicit dimensions, ads without reserved space, web fonts causing FOUT
4. **Audit INP (Interaction to Next Paint)** — use Chrome DevTools Performance panel → Interactions; INP > 200ms requires investigation; find long tasks blocking the main thread during input handling
5. **Run Lighthouse for opportunity scoring** — use Lighthouse in DevTools or CI to get scored recommendations; prioritize by estimated savings × user segment size
6. **Implement and measure in field** — deploy one change at a time; wait 7-28 days for CrUX data to update; confirm p75 improvement before attributing causation
7. **Set performance budgets in CI** — use Lighthouse CI or web-vitals library to fail builds that regress LCP > 2.5s, CLS > 0.1, or INP > 200ms

## Rules

- Never optimize based on Lighthouse score alone — a 100 Lighthouse score does not guarantee good field CWV
- Measure on real mobile hardware with throttled CPU, not desktop DevTools device emulation
- Fix TTFB before optimizing rendering — a slow server response negates all client-side optimization
- CLS must be measured across full page lifecycle including late-loading ads and embeds

## Examples

**LCP diagnosis:** LCP element is a hero image. Waterfall shows it is discovered late (not in HTML, injected by JS). Fix: move image to HTML `<img>` with `fetchpriority="high"`. LCP drops from 4.2s to 1.8s.

## Common Mistakes

- Optimizing desktop Lighthouse score while ignoring mobile field data: 60%+ of web traffic is mobile, where CWV scores are dramatically worse
- Attributing CLS to one element without checking the full page scroll: CLS accumulates across all shifts in the viewport lifetime
- Deploying multiple performance changes simultaneously: makes it impossible to attribute which change caused the improvement

## When NOT to Use

- The site receives fewer than 1,000 monthly sessions — CrUX requires a minimum traffic threshold to populate field data, making the entire field-first diagnostic approach impossible.
- The performance complaint is about a React or Vue rendering bottleneck causing UI freezes unrelated to page load — profiling Web Vitals metrics (LCP, CLS, INP at page level) will not isolate in-component re-render loops; use the React Profiler or Vue DevTools instead.
- The website is a server-rendered admin dashboard used exclusively on corporate LAN by internal users — Core Web Vitals are optimized for public internet traffic patterns, and optimizing them for this audience yields no business value.
