---
name: design-offline-first
description: Use when designing a mobile or web application that must function without a reliable network connection
source: "Offline First (offlinefirst.org); CouchDB Guide (Apache); 'Designing Offline-First Web Applications' (A List Apart, Alex Feyerke, 2013)"
tags: [mobile, offline-first, sync, conflict-resolution, service-worker, pwa, couchdb]
verified: true
---

# Design Offline-First

Design data and sync architecture so the app works fully offline and merges changes correctly when connectivity returns.

## Why This Is Best Practice

**Adopted by:** Salesforce mobile, Google Docs, Notion, Linear, and Figma — all implement offline-first sync; CouchDB/PouchDB pattern is the canonical reference implementation for offline sync

**Impact:** Google research shows 53% of users abandon mobile sites that take over 3 seconds to load; apps that work offline eliminate the entire class of "no connection" failures that drive uninstalls; Figma reported offline support as a top enterprise adoption driver

**Why best:** Treating the network as unreliable by default — rather than an exception — produces more resilient applications. Offline-first means writes go to local storage immediately and sync to the server eventually, rather than blocking on network availability for every user action.

## Steps

1. **Identify offline scope** — decide which features must work offline (core write operations) vs. which can degrade gracefully (analytics, search, real-time collaboration)
2. **Choose a local storage strategy** — SQLite for relational data on mobile, IndexedDB or PouchDB for web, Core Data for iOS; select based on query complexity and sync library support
3. **Design the conflict resolution model** — choose a strategy: last-write-wins (simple, lossy), operational transforms (complex, lossless), or CRDTs (eventual consistency, no central coordinator)
4. **Implement optimistic UI** — write to local store immediately and reflect the change in the UI; do not wait for server confirmation; queue the sync operation
5. **Build the sync queue** — persist all local mutations to a durable queue; replay the queue on reconnect in the order they were created; mark each operation with a logical timestamp
6. **Handle sync conflicts explicitly** — detect conflicts on server merge; apply your resolution strategy; surface unresolvable conflicts to the user rather than silently discarding data
7. **Test under network conditions** — use Chrome DevTools throttling, Android emulator network profiles, and XCTest network conditioning to verify behavior under 0%, 2G, and lossy connections

## Rules

- Never make a write operation dependent on network availability — local write must always succeed
- Sync errors must not surface as user-facing errors unless they are unresolvable; queue silently, resolve in background
- Use logical clocks (Lamport timestamps or vector clocks) rather than wall-clock time for ordering conflicting operations
- Offline scope must be documented and communicated; users must know which features require connectivity

## Examples

**Optimistic write flow:** User checks off a to-do → UI updates immediately → write queued to IndexedDB → background sync pushes to server → server confirms → queue entry removed. If offline for 3 days: all checkmarks persist locally, sync on reconnect, server applies all in order.

## Common Mistakes

- Designing offline as an afterthought: retrofitting sync onto an online-first data model is an order-of-magnitude harder than designing for it initially
- Last-write-wins without user awareness: silently discarding a user's edits because a teammate saved 1 second later will cause data loss complaints
- Not persisting the sync queue across app restarts: a queue only in memory is lost on crash, causing silent data loss

## When NOT to Use

- The application is inherently real-time and requires live data to be useful (e.g., stock trading, live video streaming, multiplayer gaming) — stale local state would produce incorrect or dangerous outcomes.
- The data is highly sensitive with strict regulatory requirements (e.g., HIPAA, PCI-DSS) where persisting data to the local device storage creates an unacceptable compliance risk without additional encryption and audit controls being designed first.
- The team is rebuilding an existing online-only product under a 2-week deadline — offline-first sync architecture cannot be safely retrofitted in this timeframe and the attempt will introduce data corruption bugs worse than the original connectivity errors.
