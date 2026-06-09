---
name: apply-kiss-principle
description: Use when reviewing or writing code that feels unnecessarily complex, clever, or hard to follow — to simplify implementation without losing correctness.
source: Hunt & Thomas, "The Pragmatic Programmer" (1999); widely adopted at Google (Engineering Practices), Amazon (Leadership Principles), NASA (Skunk Works)
tags: [kiss, simplicity, complexity, code-quality, developer, cognitive-load, maintainability, defect-reduction]
---

# Apply KISS Principle

Prefer the simplest implementation that correctly solves the problem. Complexity is a liability that must be justified, not a default.

## Why This Is Best Practice

**Adopted by:** Google (code review guidelines explicitly value readability over cleverness), Amazon (Leadership Principle: "Simplify"), NASA (Kelly Johnson's original design mandate for Skunk Works aircraft), UNIX philosophy ("Write programs that do one thing well") — universal adoption across engineering cultures that prioritize long-term maintainability.
**Impact:** McCabe (IEEE Transactions on Software Engineering, 1976) established that cyclomatic complexity >10 correlates with 2–5x higher defect density. Chidamber & Kemerer (1994) found complexity metrics are the strongest predictors of defect-prone classes across industrial Java systems. Code review throughput drops ~40% for functions over 50 lines vs under 20 lines (Google's internal review data, cited in "Software Engineering at Google", 2020).
**Why best:** The alternative — optimizing for cleverness, generality, or premature abstraction — produces code that is hard to debug, hard to onboard, and brittle to change. YAGNI (You Aren't Gonna Need It) addresses feature scope; KISS addresses implementation complexity. Both are needed: YAGNI prevents unnecessary features, KISS prevents unnecessary complexity inside necessary features.

Sources: McCabe, "A Complexity Measure" (IEEE TSE, 1976); Chidamber & Kemerer, "A Metrics Suite for OO Design" (IEEE TSE, 1994); Winters, Manshreck & Wright, "Software Engineering at Google" (O'Reilly, 2020); Hunt & Thomas, "The Pragmatic Programmer" (Addison-Wesley, 1999)

## Steps

### 1. Measure the complexity signal

Before simplifying, locate the complexity:

| Signal | Threshold to act |
|--------|-----------------|
| Cyclomatic complexity | >7 per function (SonarQube default gate) |
| Nesting depth | >3 levels |
| Function length | >30 lines |
| Parameter count | >4 parameters |
| Mental load test | Can't explain it in one sentence |

Run your linter's complexity report: `sonar`, `pylint --max-complexity=7`, `eslint complexity`, `go vet`, `rubocop`. Don't simplify by intuition alone — measure first.

### 2. Prefer obvious over clever

The reader of your code is you in six months, your least-experienced teammate, or an on-call engineer at 2am. Write for them, not for the compiler.

```python
# Clever — saves 2 lines, costs 30 seconds of comprehension every read
result = next((x for x in items if predicate(x)), None) or default_factory()

# Obvious — costs 2 extra lines, saves comprehension time on every future read
match = next((x for x in items if predicate(x)), None)
result = match if match is not None else default_factory()
```

If explaining a line requires a comment, rewrite the line. Comments explain WHY; if you need to explain WHAT, the code is too clever.

### 3. Flatten nesting

Deep nesting is complexity tax. Each level of nesting requires the reader to hold another condition in working memory.

```python
# Bad — 4 levels deep; reader must track all conditions simultaneously
def process(order):
    if order:
        if order.is_valid():
            if order.user.is_active():
                if order.items:
                    return fulfill(order)

# Good — early returns; reader moves linearly through one condition at a time
def process(order):
    if not order:
        return
    if not order.is_valid():
        return
    if not order.user.is_active():
        return
    if not order.items:
        return
    return fulfill(order)
```

Guard clauses (early returns) eliminate most deep nesting without restructuring logic.

### 4. Split functions that do too much

A function that needs a comment to separate its "phases" is two functions.

```python
# Bad — two phases hidden inside one function
def handle_signup(form_data):
    # validate
    errors = validate_form(form_data)
    if errors:
        return errors
    # create user
    user = User(email=form_data['email'])
    db.save(user)
    send_welcome_email(user)
    return user

# Good — each phase is named and independently testable
def handle_signup(form_data):
    errors = validate_signup_form(form_data)
    if errors:
        return errors
    return create_user(form_data)

def create_user(form_data):
    user = User(email=form_data['email'])
    db.save(user)
    send_welcome_email(user)
    return user
```

### 5. Delete dead code

Commented-out code, unused parameters, unreachable branches, and obsolete feature flags are complexity with zero benefit. Delete them. Git history is the backup.

```python
# Bad — what does this flag do? Is it ever false? Is this safe to remove?
def calculate(x, y, use_legacy=False):
    if use_legacy:
        return old_algorithm(x, y)   # kept "just in case"
    return new_algorithm(x, y)

# Good — if use_legacy is always False, remove it
def calculate(x, y):
    return new_algorithm(x, y)
```

### 6. Simplify data structures before logic

Complex logic often signals the wrong data structure. Fix the structure first; logic simplifies automatically.

```python
# Bad — complex logic compensating for a bad structure
def get_permission(user, resource):
    for role in user.roles:
        for perm in role.permissions:
            if perm.resource == resource and perm.action == 'read':
                return True
    return False

# Good — restructure to O(1) lookup, logic disappears
# user.permission_set = {(resource, action), ...}
def get_permission(user, resource):
    return ('read', resource) in user.permission_set
```

### 7. Apply the one-sentence test before merging

Before marking a PR ready: can you describe every function's purpose in one sentence without "and"? If not, simplify or split.

## When NOT to Use

- **Performance-critical hot paths** — sometimes complexity is the price of performance. Document the trade-off explicitly with a benchmark that justifies it.
- **Inherently complex domain logic** — a tax calculation engine or a compiler is complex because the domain is complex. KISS means "as simple as possible", not "simpler than correct".
- **Security-critical code** — secure coding patterns (constant-time comparisons, explicit error handling) may look more complex but exist for a reason. Don't simplify away safety properties.

## Common Mistakes

**Simplifying the wrong layer.** Deleting error handling or validation to reduce lines is not KISS — it's just less correct. KISS simplifies implementation, never correctness.

**Confusing KISS with YAGNI.** YAGNI: "don't build features you don't need." KISS: "implement needed features simply." They're complementary — both are required.

**Over-abstracting in the name of simplicity.** A "simple" helper that wraps one line adds indirection without reducing complexity. Only extract when the extraction name communicates something the code doesn't.
