---
name: diagnose-memory-leak
description: Use when a process shows continuously growing memory usage that doesn't stabilize, when OOM kills are occurring in production, or when heap dumps are needed to identify object accumulation
source: Brendan Gregg "Systems Performance" (2020) memory chapter; Java/JVM GC analysis (Oracle docs); Gregg "BPF Performance Tools" (2019)
tags: [performance, debugging, memory, reliability]
verified: true
---

# Diagnose Memory Leak

Systematically identify the source of unbounded memory growth in a process by profiling allocations, analyzing heap contents, and tracing the code path responsible.

## Why This Is Best Practice

**Adopted by:** Netflix (Java heap analysis for streaming services), Google (memory profiling as standard performance practice), all production engineering teams that run long-lived processes
**Impact:** Memory leaks are the #3 cause of production incidents (after hardware and deployment failures) in long-running services; undiagnosed leaks cause 2-4 AM OOM kills that take 20+ minutes to recover; structured diagnosis reduces MTTI from hours to 30 minutes
**Why best:** Memory leaks have many causes (event listener accumulation, cache without eviction, retained closures, circular references, native memory leaks); only systematic profiling pinpoints the specific cause

Sources: Gregg "Systems Performance" 2nd ed. (2020) Ch. 7; Gregg "BPF Performance Tools" Pearson (2019); Oracle JVM GC documentation; Node.js memory profiling guide

## Steps

1. **Confirm the leak exists** — Monitor process memory over time (at least 24-48 hours): `ps -o pid,rss,vsz,comm -p <pid>` every minute, or use Prometheus process_resident_memory_bytes. A genuine leak shows monotonically increasing RSS that doesn't return to baseline after load decreases. Natural heap growth that stabilizes is not a leak.

2. **Determine the memory region** — Distinguish heap vs. non-heap memory growth. JVM: heap (`-Xmx` bounded) vs. off-heap (Metaspace, direct buffers, native libraries). Node.js: V8 heap vs. native addons. A growing heap is a managed code leak; growing non-heap is a native memory or off-heap allocation leak. Diagnosis tools differ by region.

3. **Enable detailed GC logging (JVM)** — Add JVM flags: `-Xlog:gc*:file=/tmp/gc.log:time,uptime:filecount=10,filesize=10m`. Analyze with GCEasy.io or IBM GC analyzer. A rising old-gen baseline after full GC indicates objects surviving collection that should be collected — classic leak signature. GC logs reveal memory pressure before OOM.

4. **Capture heap snapshots at intervals** — JVM: `jmap -dump:format=b,file=/tmp/heap1.hprof <pid>` at two points 30 minutes apart. Node.js: `v8.writeHeapSnapshot()` or `--inspect` + Chrome DevTools heap snapshot. Compare the two snapshots to identify which object types grew. Objects that accumulate between snapshots without being released are leak candidates.

5. **Analyze heap with a profiler** — JVM: open heap dump in Eclipse MAT (Memory Analyzer Tool). Use "Leak Suspects" report — MAT identifies objects with high retained heap. Look for: collections (HashMap, ArrayList) growing without bound, listener registries, thread-local variables. Node.js: Chrome DevTools "Comparison" between two heap snapshots shows which object types increased.

6. **Identify the retaining reference chain** — In Eclipse MAT: right-click the leaking object → "Path to GC Roots" → "Exclude soft/weak/phantom references". This shows the reference chain keeping the object alive. The root of the chain (a static field, a thread-local, a framework registry) is the leak source. In Node.js DevTools: "Retainers" panel shows the same chain.

7. **Profile allocations in real time** — JVM: use async-profiler (`./profiler.sh -e alloc -d 60 -f alloc.svg <pid>`) to generate an allocation flame graph. Shows which call stacks allocate the most objects — high-frequency allocation of objects that should be short-lived but are retained points to the leak site. Gregg's BPF tools: `memleak` from bcc-tools traces native memory allocations.

8. **Check common leak sources by runtime** — JVM: static collections used as caches without eviction, ClassLoader leaks in hot-reload environments, ThreadLocal variables not removed, listeners registered but never deregistered. Node.js: event emitters with unlimited listeners (`emitter.setMaxListeners(0)` hides the warning), global caches with no TTL, closures capturing large objects. Go: goroutine leaks (goroutines blocked forever), map entries never deleted.

9. **Reproduce in a controlled environment** — Write a test that reproduces the leak: run the suspected code path in a loop and monitor heap. Confirmed if heap grows without bound. This test becomes the regression test — the fix must make this test pass (stable heap under load). Reproducing the leak is a prerequisite for verifying the fix.

10. **Implement and verify the fix** — Apply the fix: remove the static reference, add eviction to the cache, deregister the listener on cleanup, use weak references where appropriate. Run the reproduction test for 30 minutes. Verify heap stabilizes at a flat baseline after warmup. Deploy to staging with extended monitoring before promoting to production.

## Rules

- Heap dumps from production processes containing user data must be handled as sensitive data (PII may be in memory); apply the same access controls as production databases.
- Never tune heap size (-Xmx) as the first response to a memory leak; you are delaying the OOM, not fixing the leak.
- A fix that reduces growth rate but doesn't eliminate it is not a fix; verify heap is flat (zero net growth) after warmup.
- Always reproduce the leak in a controlled test before claiming a fix; "it seems better" after a production deploy is confirmation bias, not validation.

## Common Mistakes

- **Increasing heap size as the fix** — this delays the OOM crash by hours or days; the leak recurs and the system still crashes; find and fix the root cause.
- **Analyzing heap during high load** — heap snapshots during load show normal temporary allocations mixed with leaked objects; capture snapshots after load decreases to isolate what should have been collected.
- **Ignoring non-heap memory** — a Java process where heap is stable but RSS is growing has an off-heap leak (DirectByteBuffer, JNI, Metaspace); heap analysis tools miss this.
- **Single heap snapshot** — a single snapshot shows the heap state but cannot distinguish what grew; always compare two snapshots taken at different points in time to isolate growth.

## When NOT to Use

- Memory growth that stabilizes after warmup — this is natural heap sizing, not a leak; JVMs fill heap to improve GC efficiency
- Memory growth proportional to data volume — storing user data in memory that grows with users is architectural, not a leak; evaluate caching strategy instead
- OOM during initialization before the service is serving traffic — this is a configuration issue (insufficient heap for the workload), not a leak
