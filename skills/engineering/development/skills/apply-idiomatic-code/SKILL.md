---
name: apply-idiomatic-code
description: Use when writing or reviewing code in a specific language or ecosystem — especially when contributors from different language backgrounds write in conflicting styles, or when code compiles but feels foreign to native speakers of the language.
source: Dan North, "CUPID—for joyful coding", https://dannorth.net/2022/02/10/cupid-for-joyful-coding/, 2022
emerging: true
tags: [idioms, code-quality, readability, conventions, developer, onboarding, maintainability]
related: [apply-composable-design, apply-unix-philosophy, apply-domain-based-naming, apply-predictable-code, review-code-style]
---

# Apply Idiomatic Code

Write code that feels natural to experienced practitioners of the language — using its conventions, patterns, and standard library as intended.

## Why This Is Best Practice

**Adopted by:** All major language communities maintain idiomatic guides: Go (Effective Go), Python (PEP 8 + Pythonic idioms), Rust (Rust API Guidelines), JavaScript/TypeScript (Airbnb, Google style guides), Ruby (community idioms enforced by RuboCop). Language designers invest significant effort making certain patterns the "natural" path — idiomatic code walks that path.
**Impact:** Idiomatic code is understood instantly by anyone experienced in the language. Non-idiomatic code requires readers to translate: "this is Java-style code in Python" or "this is procedural code in a functional language." That translation tax compounds across every review, every debugging session, and every onboarding. The Go team found in user research that non-idiomatic code is the most common source of confusion for new contributors to an existing Go codebase.
**Why best:** Idiomatic code has an additional advantage: it works *with* the language's tooling. Linters, formatters, refactoring tools, and documentation generators are designed around idiomatic patterns. Non-idiomatic code resists these tools — it generates false positives, confuses autocompletion, and cannot be safely refactored automatically. This is distinct from code style (`review-code-style`): style governs formatting; idioms govern which language features and patterns to choose.

Sources: Dan North, "CUPID—for joyful coding" (2022); Pike, "Effective Go" (golang.org); van Rossum, Warsaw & Coghlan, "PEP 8" (python.org); Blandy & Orendorff, "Programming Rust" (O'Reilly, 2021)

**Status:** Emerging — idiomatic coding is widely practiced but not yet codified as a named, standalone best practice across industry; CUPID's framing (2022) is consolidating the concept.

## Steps

### 1. Read the language's official idiom guide

Every major language has one. Read it before writing production code in a new language.

| Language | Primary idiomatic reference |
|---|---|
| Go | Effective Go (golang.org/doc/effective_go) |
| Python | PEP 8 + "The Hitchhiker's Guide to Python" |
| Rust | Rust API Guidelines (rust-lang.github.io/api-guidelines) |
| JavaScript | MDN + language spec; Airbnb or Google style guides for conventions |
| Ruby | Ruby Style Guide (rubystyle.guide) |
| Kotlin | Kotlin Coding Conventions (kotlinlang.org/docs/coding-conventions) |

---

### 2. Use language-native patterns, not patterns from your previous language

The most common source of non-idiomatic code is a developer bringing patterns from a language they know well into one they are learning.

```python
# Bad — Java-style loop, written by someone coming from Java
i = 0
while i < len(items):
    process(items[i])
    i += 1

# Good — Python idiom
for item in items:
    process(item)

# Also good — when index is needed
for i, item in enumerate(items):
    process(i, item)
```

```go
// Bad — error handling style from Java (exception-like chaining)
result, _ := fetchUser(id)  // silently ignored error

// Good — Go idiom: always handle the error
result, err := fetchUser(id)
if err != nil {
    return fmt.Errorf("fetching user %s: %w", id, err)
}
```

---

### 3. Prefer standard library over custom implementations

The standard library is the most idiomatic code in any language — it sets the patterns everything else follows. A custom implementation of something in the standard library is almost always worse.

- Check the standard library before reaching for an external package
- Check external packages before writing custom implementations
- The order: stdlib → well-known package → custom

```python
# Bad — custom date parsing when stdlib handles it
def parse_date(s):
    parts = s.split("-")
    return {"year": int(parts[0]), "month": int(parts[1]), "day": int(parts[2])}

# Good — stdlib
from datetime import date
d = date.fromisoformat("2024-03-15")
```

---

### 4. Use the language's error and null handling patterns

Every language has a preferred way to signal absence and failure. Using a different pattern (e.g., returning `null` in a language that prefers `Optional`, or throwing exceptions in a language that returns errors) makes code harder to compose with the rest of the ecosystem.

```rust
// Bad — panic instead of returning Result; non-idiomatic for library code
fn parse_config(path: &str) -> Config {
    let content = std::fs::read_to_string(path).unwrap(); // panics on error
    toml::from_str(&content).unwrap()
}

// Good — Rust idiom: return Result, let caller decide how to handle
fn parse_config(path: &str) -> Result<Config, Box<dyn Error>> {
    let content = std::fs::read_to_string(path)?;
    Ok(toml::from_str(&content)?)
}
```

---

### 5. Let the linter and formatter enforce idioms automatically

Manual idiom enforcement in code review is slow and inconsistent. Configure tooling to enforce idioms at commit time.

| Language | Idiom enforcement tool |
|---|---|
| Go | `go vet`, `staticcheck`, `golangci-lint` |
| Python | `ruff`, `pylint`, `mypy` |
| Rust | `clippy` |
| TypeScript | `eslint` with language-specific rule sets |
| Ruby | `rubocop` |

Run these in CI. Fail the build on violations. This removes idiom discussions from code review and into tooling.

## Common Mistakes

**Idiomatic ≠ correct.** Idiomatic code can still be wrong. Idioms govern style; correctness governs behavior. Don't trade correctness for idiom.

**Imposing idioms from a different version.** Python 2 idioms are wrong in Python 3. Go 1.18 generics idioms don't exist in Go 1.16. Know which version's idioms apply.

**Confusing idioms with style.** Code style (spacing, bracket placement) is enforced by formatters (`gofmt`, `black`, `prettier`). Idioms are patterns — which features to use, how to handle errors, how to structure modules. Both matter; they are different.

**Fighting the type system.** A common non-idiomatic pattern is working around the type system (excessive casting, `any` everywhere in TypeScript, `interface{}` in Go) instead of using it. The type system is idiomatic — use it.

## When NOT to Use

- **Adapting to an existing non-idiomatic codebase** — introducing one idiomatic module into a codebase full of non-idiomatic code creates inconsistency, which is worse than consistent non-idioms. Migrate incrementally or all-at-once; don't mix.
- **Following an established team convention that differs from the language default** — if the team has a documented reason for a non-default pattern (e.g., explicit `null` over `Optional` for performance), follow the team convention over the language default. Team consistency beats language purity.
- **Research or prototype code** — idiom compliance adds a review overhead that is not justified for throwaway code
