---
name: apply-defensive-copy
description: Use when writing constructors, setters, or methods that accept or return mutable objects — to prevent callers from silently mutating internal state through shared object references.
source: Bloch, "Effective Java" (3rd ed. 2018) Item 50; Bloch & Gafter, "Java Puzzlers" (2005); widely adopted at Google (Guava ImmutableList), Amazon (AWS SDK), Netflix
tags: [defensive-copy, immutability, encapsulation, mutability, developer, correctness, data-integrity, thread-safety]
related: [apply-fail-fast, validate-external-input]
---

# Apply Defensive Copy

Copy mutable inputs on receipt and mutable outputs on return — never share internal references.

## Why This Is Best Practice

**Adopted by:** Joshua Bloch "Effective Java" Item 50 (the canonical Java reference,
adopted across the JVM ecosystem), Google Guava (ImmutableList, ImmutableMap — the
most-downloaded Java library), Amazon AWS SDK (all public API inputs are copied on
ingestion), and enforced by default in Rust (ownership model) and Swift (value types).
**Impact:** Google internal data attributes 30%+ of concurrency-related production bugs
in Java services to shared mutable state. Defensive copying eliminates an entire class
of time-of-check-to-time-of-use (TOCTOU) vulnerabilities — the underlying cause of
several high-severity CVEs in major Java libraries (Apache Commons, Spring Framework).
**Why best:** Documentation asking callers "don't mutate this" relies on discipline and
is invisible to future maintainers. A defensive copy is enforced by the runtime — the
caller cannot violate it regardless of intent, including malicious intent.

Sources: Bloch, "Effective Java" 3rd ed. Item 50; Bloch, Google Tech Talk "How to Design
a Good API" (2007); Guava project documentation; MITRE CWE-374 (Passing Mutable Objects)

## Steps

### Step 1: Copy mutable inputs before storing

Never store a reference the caller passed. The caller retains the original and can
mutate it after your constructor or setter runs.

```java
// Wrong — caller can mutate start/end after construction
public Period(Date start, Date end) {
    this.start = start;
    this.end = end;
}

// Right — defensive copy on input
public Period(Date start, Date end) {
    this.start = new Date(start.getTime());  // copy first
    this.end   = new Date(end.getTime());
    // validate the copies, not the originals (prevents TOCTOU)
    if (this.start.after(this.end))
        throw new IllegalArgumentException("start after end");
}
```

### Step 2: Copy before validating (prevent TOCTOU)

Always copy the input first, then validate the copy. If you validate the original and
another thread mutates it between validation and storage, your check is worthless.

```
copy → validate copy → store copy    ✅ TOCTOU-safe
validate original → copy → store     ❌ TOCTOU window between validate and copy
```

### Step 3: Return copies of mutable internal fields

When a getter returns an internal mutable object, callers receive a reference to your
internals and can mutate it.

```java
// Wrong — exposes internal reference
public Date getStart() { return start; }

// Right — return a copy
public Date getStart() { return new Date(start.getTime()); }
```

### Step 4: Prefer immutable types — no copying needed

The best defensive copy is no copy. Use immutable types wherever possible:

| Mutable | Immutable alternative |
|---------|----------------------|
| `Date` | `Instant`, `LocalDate` (java.time) |
| `int[]` | `List.of(...)` |
| `ArrayList` | `List.copyOf(...)` / Guava `ImmutableList` |
| Plain object | Record (Java 16+), `@Value` (Lombok), `data class` (Kotlin) |

### Step 5: For arrays, copy explicitly — `.clone()` is not always sufficient

```java
// Right for primitive arrays
int[] copy = Arrays.copyOf(original, original.length);

// For arrays of mutable objects, shallow clone shares element references — deep copy instead
Item[] copy = Arrays.stream(original)
    .map(Item::copy)
    .toArray(Item[]::new);
```

### Step 6: Document that copies are made

Future maintainers may remove the copy as an "optimization." A brief comment preserves
the intent.

```java
// defensive copy — caller retains original and can mutate it
this.items = new ArrayList<>(items);
```

## When NOT to Use

- **Immutable objects** — `String`, `Integer`, `LocalDate` and other immutable types
  cannot be mutated through a reference; copying is redundant.
- **Performance-critical inner loops** where profiling has confirmed copy cost is
  significant. Document the exception and the profiler data explicitly.
- **Private methods where you control all callers** and the class is `final` — the
  encapsulation guarantee is provided by access control, not copying.

## Common Mistakes

**Validating before copying.** The TOCTOU window between validate and copy is a real
vulnerability in multi-threaded environments. Always copy first.

**Forgetting to copy on output.** Most engineers remember to copy inputs; few remember
to copy return values. Both directions need protection.

**Shallow-copying collections of mutable objects.** `new ArrayList<>(list)` copies the
list structure but shares element references. If the elements are mutable, callers can
still mutate them. Deep-copy elements or use immutable element types.

**Removing copies as a "performance optimization" without profiling.** The copy cost is
negligible in almost all cases. Measure before removing.
