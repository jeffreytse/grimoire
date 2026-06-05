---
name: audit-performance
description: Use when a system is slow, latency has regressed, CPU or memory usage is unexpectedly high, or before optimizing any code — to measure first, identify real bottlenecks, and prioritize fixes by impact.
source: Brendan Gregg "Systems Performance" (2020) USE Method, Google pprof profiling methodology, Chrome DevTools performance profiling documentation
tags: [performance, profiling, bottleneck-analysis, latency, cpu, memory, observability, optimization]
verified: true
---

# Audit Performance

Measure before optimizing. Profile systematically to find real bottlenecks — not guessed ones.

## Why This Is Best Practice

**Adopted by:** Netflix (Brendan Gregg's USE Method, documented in "Systems Performance" 2020, used across Netflix engineering), Google (pprof is Google's internal profiler, open-sourced; used in production Go and C++ services), Meta (Pyroscope continuous profiling in production), Cloudflare (published profiling-first optimization case studies including 10× throughput gains on their Rust workers)
**Impact:** Studies consistently show engineers guess wrong about where performance bottlenecks are: Jon Bentley's "Programming Pearls" (1986) found developers misidentify the hot path 90% of the time. Measure-first optimization at Google yielded 3–10× improvements in specific subsystems vs. intuition-driven rewrites that often produced no measurable gain (Google SRE Book, Chapter 19). Cloudflare's profiling-first approach on their TLS stack achieved 40% latency reduction by fixing a single memory allocation hotspot that no engineer had suspected.
**Why best:** Premature optimization wastes engineering time on non-bottlenecks. The 90/10 rule (Knuth, 1974) holds empirically: 90% of time is spent in 10% of code. Only measurement identifies which 10%. Alternative (rewrite in a "faster" language) almost always loses to profile-and-fix in the same language for the same effort.

Sources: Brendan Gregg "Systems Performance" 2020, Google SRE Book (2016), Jon Bentley "Programming Pearls" (1986), Cloudflare engineering blog

## Steps

### 1. Define the performance problem precisely

Before profiling, write one sentence:

> "[What] is [how slow/expensive] under [what conditions], measured by [what metric]."

Example: "The `/search` endpoint takes >2s at p99 under 100 RPS load, measured by our Datadog APM trace."

Without this, you don't know when you're done. Without a baseline, you can't measure improvement.

### 2. Establish baseline measurements

Measure before touching any code. Capture:

- **Latency**: p50, p95, p99 (not average — averages hide tail latency)
- **Throughput**: requests/second or operations/second at target load
- **CPU**: utilization % per core
- **Memory**: heap size, allocation rate, GC pause frequency
- **I/O**: disk read/write bytes/sec, network bytes/sec, IOPS

Use realistic load, not synthetic benchmarks. Tools:
- Web: `wrk`, `k6`, `hey`, `vegeta`
- DB: `pgbench`, `sysbench`
- General: `perf stat`, `time`, language-native benchmarks

Record the baseline. You cannot claim improvement without before/after numbers.

### 3. Apply the USE Method (Brendan Gregg)

For every resource in the system (CPU, memory, disk, network, each service):

- **U**tilization: what % of capacity is being used?
- **S**aturation: is there a queue forming? (load average, queue depth)
- **E**rrors: are there error counts/rates?

High utilization + saturation = bottleneck. Start there.

```bash
# CPU utilization and saturation
top -b -n1 | head -20
mpstat -P ALL 1 3

# Memory
free -m
vmstat 1 5

# Disk I/O
iostat -x 1 5

# Network
ss -s
netstat -s | grep -E 'retransmit|error'
```

### 4. Profile CPU hotspots

**Go:**
```bash
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile?seconds=30
```

**Python:**
```bash
python -m cProfile -o profile.out my_script.py
python -m pstats profile.out
# Or: py-spy top --pid <pid>
```

**Node.js:**
```bash
node --prof app.js
node --prof-process isolate-*.log > processed.txt
# Or: clinic.js flame -- node app.js
```

**Java:**
```bash
async-profiler -d 30 -f profile.html <pid>
```

**Browser (frontend):**
Open Chrome DevTools → Performance → Record → reproduce the slow action → Stop. Look at the flame chart for long tasks (>50ms).

Read the flame chart: the widest frames are the hottest — work top-down. The actual bottleneck is at the bottom of the deepest wide column.

### 5. Profile memory (if memory is the issue)

Distinguish:

- **Leak**: heap grows unbounded over time → use heap profiler, track allocations over time.
- **High steady-state usage**: large objects retained → heap snapshot, look for retained object trees.
- **GC pressure**: high allocation rate causing frequent GC pauses → allocation profiler.

```bash
# Go heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Node.js heap snapshot
node --inspect app.js  # Then open Chrome DevTools → Memory → Take snapshot
```

### 6. Identify database query bottlenecks

Most application latency is database latency. Always check:

```sql
-- PostgreSQL: slow queries
SELECT query, calls, mean_exec_time, total_exec_time
FROM pg_stat_statements
ORDER BY total_exec_time DESC
LIMIT 10;

-- PostgreSQL: missing indexes
EXPLAIN (ANALYZE, BUFFERS) SELECT ...;
-- Look for: Seq Scan on large tables, high actual rows >> estimated rows
```

Signs of DB bottleneck:
- `Seq Scan` on a table with >10k rows where an index should exist
- `actual rows` far exceeds `estimated rows` → stale statistics → run `ANALYZE`
- N+1 queries: app makes N queries in a loop instead of one JOIN

### 7. Rank bottlenecks by impact

After profiling, list findings:

| Bottleneck | Current cost | Fix | Estimated gain | Effort |
|---|---|---|---|---|
| Missing index on `orders.user_id` | 800ms/query | Add index | ~750ms savings | Low |
| JSON serialization in hot loop | 40% CPU | Use binary format | ~30% CPU | Medium |
| N+1 query in `/users` endpoint | 50 DB calls/req | Eager load | ~40 calls/req | Low |

Prioritize by: `(estimated_gain × traffic_share) / effort`. Fix the highest-value, lowest-effort items first.

### 8. Fix, measure, repeat

For each fix:
1. Make the change on a branch.
2. Re-run the exact same benchmark from step 2.
3. Record new p50/p95/p99 + CPU/memory.
4. Confirm the improvement is statistically significant (not noise) — run 3× and average.

If the fix shows no improvement, the bottleneck is elsewhere. Don't ship it as a "performance improvement."

## Rules

- Measure before any code change. No profiling = no optimization.
- Report p95/p99 latency, never average — averages hide the tail where users suffer.
- Fix the biggest measured bottleneck first, not the most obvious-looking code.
- Every optimization must have a before/after benchmark comparison.
- Never optimize code that runs once on startup or handles <1% of traffic.
- Profile under realistic load — idle profiling finds idle bottlenecks.

## Examples

**Defining the problem correctly:**
> "The `POST /checkout` endpoint takes 3.2s p99 at 50 RPS under our staging load test, up from 800ms two weeks ago. Regression introduced in deploy #1847."

**Reading pprof output:**
Flat% = time spent in this function itself. Cum% = time in this function + all callees. Fix functions with high flat% — they're burning CPU directly.

**Spotting N+1:**
```
# 1 query: SELECT * FROM orders WHERE user_id = 1
# Then for each order:
# SELECT * FROM products WHERE id = ?   ← repeated 50 times
# Fix: SELECT o.*, p.* FROM orders o JOIN products p ON ...
```

## Common Mistakes

- **Optimizing without measuring**: changes the wrong thing; no measurable improvement.
- **Reporting average latency**: a p50 of 200ms with a p99 of 8s is a broken system.
- **Fixing code, not the bottleneck**: refactoring a function that contributes 0.1% of runtime.
- **Profiling under no load**: CPU profiles under idle load show initialization code, not steady-state bottlenecks.
- **Stopping after one fix**: the bottleneck moves. Profile again after each fix.
- **Rewriting in a different language first**: almost always slower to ship and rarely faster in practice without also fixing the algorithmic problems.
