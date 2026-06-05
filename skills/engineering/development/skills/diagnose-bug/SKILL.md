---
name: diagnose-bug
description: Use when a bug, test failure, crash, or unexpected behavior is reported — before proposing any fix. Use when behavior diverges from specification and the root cause is unknown.
source: Systematic debugging methodology — rubber duck debugging (Pragmatic Programmer), git bisect (Linus Torvalds), scientific method applied to software (Mozilla Engineering, Google SRE)
tags: [debugging, root-cause-analysis, bug-isolation, testing, reliability, developer]
verified: true
---

# Diagnose Bug

Isolate the root cause of a bug through systematic hypothesis-and-test cycles before writing any fix.

## Why This Is Best Practice

**Adopted by:** Mozilla (engineering blog documents bisect-driven debugging as standard practice), Google SRE (Site Reliability Engineering book, O'Reilly 2016 — "debugging is hypothesis testing"), Linux kernel development (git bisect used as primary regression isolation tool since 2005).

**Impact:** Hunt (2019, "The Pragmatic Programmer" 20th anniversary) estimates that developers spend 35–50% of time debugging. Studies at Microsoft Research (Begel & Simon, 2008) found that programmers who skip reproduction steps and jump to fixes re-open the same bug 2.4× more often than those who first establish a minimal reproduction. `git bisect` can isolate a regression among 1,000 commits in 10 binary search steps.

**Why best:** "Fix first, understand later" leads to patches that mask symptoms rather than causes — the bug resurfaces under different conditions. Systematic isolation — reproduce, shrink, hypothesize, test — is the only approach that guarantees the fix targets the actual cause. It is strictly superior to random trial-and-error or reading code top-down hoping to spot the issue.

Sources: Pragmatic Programmer (Hunt & Thomas), Google SRE Book (Beyer et al., O'Reilly 2016), Mozilla Engineering Blog, Zeller "Why Programs Fail" (2009)

## Steps

### 1. Reproduce the bug reliably

Before touching code, establish a reproducible test case:
- Record the exact inputs, environment (OS, runtime version, config), and steps
- Confirm you can trigger the bug on demand; if not, the bug is intermittent — treat this as a separate class (see intermittent bugs below)
- Write the reproduction as a failing test or script — this becomes your regression guard

If you cannot reproduce it, do not proceed to fixing. Ask for more information.

### 2. Establish the smallest reproduction

Shrink the reproduction to the minimum that still triggers the bug:
- Remove unrelated inputs, dependencies, and code paths one at a time
- If the bug disappears, add the last removed piece back
- Goal: a 10-line reproduction is better than a 1,000-line one; it eliminates false hypotheses

### 3. Characterize the bug precisely

Answer these questions before hypothesizing:
- What is the *actual* output or behavior?
- What is the *expected* output or behavior?
- Where in the stack is the divergence first observable?
- What changed recently? (check git log, dependency updates, environment changes)

### 4. Isolate with binary search (bisect)

If the bug is a regression (worked before, broken now):
```bash
git bisect start
git bisect bad                  # current commit is broken
git bisect good <known-good>    # last known-good commit or tag
# git bisect runs binary search; test each commit, then:
git bisect good   # or: git bisect bad
# Repeat until bisect identifies the introducing commit
git bisect reset
```

If git bisect is not applicable, apply the binary search principle manually: comment out half the code path, confirm which half contains the bug, repeat.

### 5. Form a single hypothesis

State a falsifiable hypothesis: "The bug occurs because *X* when *Y*."

Example: "The pagination breaks because `offset` is computed before the filter is applied, so it counts unfiltered rows."

Do not form multiple hypotheses at once — test one at a time.

### 6. Test the hypothesis

Add targeted logging, assertions, or a unit test that would fail if your hypothesis is correct:
- If the test fails as predicted, the hypothesis is supported — proceed to fix
- If the test does not fail, the hypothesis is wrong — discard it and return to step 5
- Never modify the fix and the test simultaneously; change one variable at a time

### 7. Identify root cause vs. symptom

Before fixing, distinguish:
- **Symptom**: the observable failure (e.g., NullPointerException on line 42)
- **Root cause**: the condition that made the failure possible (e.g., the factory returns null for unregistered types, and no caller checks)

Fix the root cause. Fixing only symptoms leaves the system brittle.

### 8. Write the fix and verify

- Make the minimal change that eliminates the root cause
- Run the reproduction test from step 1 — it must now pass
- Run the full test suite — confirm no regressions
- Add the reproduction test as a permanent regression test if it is not already covered

### Intermittent bugs

For bugs that do not reproduce deterministically:
- Look for concurrency issues: missing locks, shared mutable state, unsynchronized counters
- Look for timing dependencies: sleeps, timeouts, event ordering
- Add logging around the suspected area and run under load
- Use race detectors if available (`go test -race`, ThreadSanitizer, Valgrind/Helgrind)

## Rules

- Never modify production code before establishing a reproduction
- Never fix two bugs in one commit — isolation breaks; each fix is independently verifiable
- If the root cause is in a dependency you do not control, file a bug upstream AND add a workaround with a comment citing the upstream issue
- Document the root cause in the commit message body — not just what changed, but why

## Examples

**Good hypothesis:** "The session expires prematurely because the token TTL is read from config at startup rather than per-request, so config reloads do not take effect."

**Bad hypothesis:** "Something is wrong with the auth code."

**Bisect finding root cause:**
```
git bisect identified commit a3f9c12 as the first bad commit.
Commit message: "refactor: centralize config loading"
Root cause: Config is now cached at startup; per-request refresh was silently removed.
```

## Common Mistakes

- Jumping to the fix without a reproduction — the fix may target the wrong code path
- Forming multiple hypotheses and testing them all at once — confounds results
- Fixing the symptom (catching the exception) instead of the cause (why the exception occurs)
- Skipping the regression test — the bug will silently re-enter the codebase
- Stopping at "it works now" without explaining *why* it works — leaves the team unable to prevent recurrence
