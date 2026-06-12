---
name: audit-bundle-size
description: Use when investigating JavaScript bundle size, slow Time to Interactive, or preparing a performance budget for a web application
source: "web.dev 'Reduce JavaScript payloads' (Google); webpack-bundle-analyzer documentation; 'JavaScript Start-up Performance' (Addy Osmani, Google, 2017)"
tags: [performance, bundle-size, javascript, webpack, tree-shaking, code-splitting, web-vitals]
verified: true
---

# Audit Bundle Size

Identify and eliminate excess JavaScript to reduce parse/compile time and reach interactive state faster.

## Why This Is Best Practice

**Adopted by:** Google's web.dev performance curriculum; Vercel and Next.js built-in bundle analyzer; Create React App and Vite expose bundle analysis tooling by default

**Impact:** Addy Osmani (Google) documented that 1MB of JavaScript takes 8 seconds to parse on a median mobile device; Google's CrUX data shows JavaScript parse time is the leading cause of poor INP and TTI on mobile

**Why best:** JavaScript is the most expensive resource per byte on the web: it must be downloaded, parsed, compiled, and executed before it produces value. Images of equivalent size are far cheaper because they only require decoding. Every KB of JavaScript removed improves TTI on low-end devices disproportionately.

## Steps

1. **Generate a bundle report** — run `npx webpack-bundle-analyzer` or `npx vite-bundle-visualizer`; for Next.js, use `@next/bundle-analyzer`; identify the largest modules by parsed size (not gzip size)
2. **Audit third-party dependencies** — check each large dependency at bundlephobia.com for gzip size and tree-shakeability; identify libraries that can be replaced with smaller alternatives (e.g., `date-fns` instead of `moment`)
3. **Check tree shaking effectiveness** — import a single named export and verify in the bundle that only that export is included; barrel files (`index.js` re-exporting everything) defeat tree shaking
4. **Implement route-based code splitting** — ensure each route loads only its own code; in React: `React.lazy` + `Suspense`; in Next.js: automatic per-page splitting plus `dynamic()`
5. **Audit polyfills** — check `@babel/preset-env` targets; polyfills for ES2018+ features are unnecessary if browserslist excludes IE11; use `useBuiltIns: 'usage'` not `entry`
6. **Set and enforce performance budget** — add `bundlesize` or Lighthouse CI to the CI pipeline; fail builds where main bundle exceeds budget (e.g., 150KB gzip for initial JS)
7. **Measure real-world impact** — deploy and check TTI and INP in CrUX after 7-28 days; bundle size reduction alone does not guarantee improvement if the critical path is unchanged

## Rules

- Optimize parsed size, not only gzip size — gzip compression reduces transfer cost but not parse/compile cost
- Never import entire libraries when only one function is needed: `import _ from 'lodash'` includes 70KB; `import debounce from 'lodash/debounce'` includes 2KB
- Code splitting must split at route boundaries at minimum; component-level splitting is additive, not a substitute
- Performance budget must be part of CI — manual audits drift; automated enforcement holds the line

## Examples

**Finding:** `moment.js` appears at 67KB gzip in the bundle; only `moment().format()` is used.
**Fix:** Replace with `date-fns/format` (2KB gzip). Net saving: 65KB. On a median Android device, this saves ~2.5 seconds of JS execution.

## Common Mistakes

- Optimizing gzip size without checking parse time: a bundle that goes from 400KB to 350KB gzip is still 350KB for the JS engine to parse
- Splitting code without prefetching critical routes: code splitting without `<link rel="prefetch">` on likely-next routes can make navigation feel slower, not faster
- Running bundle analysis only once: bundle size regressions happen incrementally through dependency updates; automate the check

## When NOT to Use

- The application is a Node.js backend service or CLI tool — JavaScript parse/compile time on V8 in a server environment is not the bottleneck; bundle size analysis addresses browser delivery costs that do not apply.
- The performance problem is high Time to First Byte (TTFB > 800ms) — server response time must be fixed before client-side bundle optimization, because no amount of JS reduction compensates for a slow origin.
- The project uses a non-JavaScript frontend (e.g., server-rendered Rails, Django templates) — bundle analysis tooling and tree-shaking concepts do not apply to server-rendered HTML without a JS build pipeline.
