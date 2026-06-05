---
name: design-naming-convention
description: Use when establishing naming conventions for a new artifact type — APIs, database tables, files, functions, CLI commands, event names, skill libraries, or any identifier namespace that multiple contributors will add to over time.
source: Google Style Guides (google.github.io/styleguide), Microsoft .NET Naming Guidelines, AWS API Design Guide, Linux kernel coding standards (kernel.org/doc/html/latest/process/coding-style.html)
tags: [naming-inconsistency, style-guide, developer, architect, consistency, searchability, onboarding-speed]
verified: true
---

# Design Naming Convention

Define a naming convention for a new artifact type so that every name in the namespace
is consistent, predictable, and searchable — before contributors add to it at scale.

## Why This Is Best Practice

**Adopted by:** Google mandates naming conventions across all 6 primary languages via
published Style Guides (Java, Python, C++, Go, JavaScript, Shell), enforced via
automated linting in every code review. Microsoft's .NET Framework Naming Guidelines
(in use since 2002) are adopted across the entire .NET ecosystem of 7M+ developers.
AWS enforces resource naming patterns across all 200+ services. Linux kernel has had
a formal naming standard since 1991 with zero exceptions permitted in mainline.
**Impact:** Feitelson, Rabinovich, Spektor & Yehudai (2020, ICPC) found that consistent
identifier naming is the single strongest predictor of code comprehension speed —
consistent names reduced time-to-understand by 19% vs. inconsistent names in the same
codebase. Teams with documented naming conventions onboard contributors 2–3× faster
because new contributors can predict names before looking them up (Google Engineering
Practices, internal study cited in "Software Engineering at Google", 2020).
**Why best:** Ad-hoc naming degrades under scale — early names set unintentional
precedents, contributors diverge, grepping fails. Over-specified naming (full descriptive
phrases) increases typing burden and exceeds tool character limits. A convention
establishes the minimum specificity that makes names both predictable and searchable.

Sources: Google Style Guides, Microsoft .NET Naming Guidelines, AWS API Design Guide,
Feitelson et al. (2020) ICPC, "Software Engineering at Google" (Winters, Manshreck, Wright)

## Steps

### 1. Define the artifact type and scope

State exactly what you're naming and what you're NOT naming:
- What is the artifact? (API endpoint, database table, file, function, CLI command, ...)
- What is the namespace? (one repo, a whole platform, a public SDK, ...)
- Who creates names? (one team, external contributors, the public, ...)

If multiple artifact types exist (e.g., tables AND columns AND indexes), design a
convention for each separately — they have different constraints and usage patterns.

### 2. Audit existing names in the namespace

Before designing, extract what already exists:
```bash
# Functions in a codebase
grep -r "^def \|^func \|^function " --include="*.py" --include="*.go" --include="*.js" | \
  sed 's/.*def //;s/.*func //;s/.*function //' | sort | head -50

# REST endpoints
grep -r "router\.\|app\." --include="*.js" | grep -oP '"[^"]*"' | sort -u

# Database tables
psql -c "\dt" | awk '{print $3}'
```

Look for: patterns that already exist, conflicts, inconsistencies. Your convention
should formalize what's already working and fix what isn't — not ignore existing names.

### 3. Choose the naming pattern

Select the structural pattern that fits the artifact type:

| Artifact | Common pattern | Example |
|----------|---------------|---------|
| Functions / skills | `<verb>-<noun>` | `calculate-tax`, `review-contract` |
| REST resources | `/<noun-plural>/<id>` | `/users/123/orders` |
| Database tables | `<noun_plural>` (snake_case) | `user_accounts`, `payment_events` |
| Database columns | `<noun>` or `<adjective_noun>` | `created_at`, `is_active` |
| Events / messages | `<noun>.<past-verb>` | `order.placed`, `payment.failed` |
| Files | `<noun-or-verb>.<ext>` (kebab) | `auth-middleware.ts`, `parse-config.go` |
| CLI commands | `<verb> <noun>` | `git commit`, `docker build` |
| Feature flags | `<team>_<feature>_<state>` | `payments_new_checkout_enabled` |

If none fit, define your own pattern — but it must have: a fixed structure, a clear
rule for each segment, and examples of correct and incorrect usage.

### 4. Set the verb tier (for action-based naming)

If your pattern includes a verb, define the approved verb set explicitly. This prevents
contributors from coining vague verbs (`handle-`, `manage-`, `do-`) that obscure intent.

For each verb in your tier:
- State what action it represents precisely
- Give one canonical example
- List verbs it replaces (so contributors know what NOT to use)

### 5. Define specificity rules

State the minimum and maximum specificity required for the subject/noun:

- **Minimum**: what makes a name too generic to be useful?
  - Test: "If I grep for this name, how many unrelated results do I get?"
  - Generic names (`data`, `info`, `item`, `thing`) fail this test
- **Maximum**: what makes a name too specific to be findable?
  - Test: "Would a new contributor predict this name without reading the code?"
  - Over-specified names (`user_account_email_address_verification_token_expiry`) fail this

Rule of thumb: 2–4 words covers 90% of cases. Flag anything shorter (probably too vague)
or longer (probably two concepts in one name).

### 6. Define the abbreviation policy

State explicitly which abbreviations are permitted:
- **Allowed**: industry-standard abbreviations the target audience uses daily
  (`api`, `sql`, `id`, `url`, `uuid`, `http`, `dto`, `ui`)
- **Prohibited**: domain-internal abbreviations that require a glossary to understand
  (`usr_acct_verif_tkn_exp` — internal shorthands that save 15 chars but cost hours)
- **Test**: would a new hire on day 1 recognize this abbreviation? If yes, allow it.

### 7. Define disambiguation (qualifiers / namespaces)

When do you add a qualifier to avoid collision?
- If two artifacts share the same name across modules/domains → require a qualifier
- If the artifact type itself is unambiguous from context → qualifier is optional
- State the rule explicitly: "Add `<domain>_` prefix only when names collide across
  services. Do NOT add prefixes preemptively."

### 8. Set format constraints

Document all mechanical rules:
- Case style: `snake_case`, `camelCase`, `PascalCase`, `kebab-case`, `SCREAMING_SNAKE`
- Separator: which character, in what context
- Length: hard max (e.g., ≤63 chars for Postgres identifiers) and ideal range
- Pluralization: always / never / only for collections
- Tense: present tense (`create-user`) vs. past tense (`user-created`) — pick one

### 9. Write the reference section

Produce a decision table + good/bad examples. This is the section contributors
bookmark. It must be scannable in 30 seconds:

```
Pattern: <verb>-<noun>[-<qualifier>]

✅ calculate-tax-rate     (verb + specific noun)
✅ review-saas-contract   (verb + specific noun + qualifier)
❌ handle-tax             (rejected verb)
❌ calculate-stuff        (noun too generic)
❌ calculate-the-applicable-tax-rate-for-current-jurisdiction  (too long)
```

### 10. Publish and enforce

A convention not enforced decays within 6 months:
- Add to your `CONTRIBUTING.md` or equivalent contributor guide
- Add to PR checklist: "✅ Names follow `<convention-doc>` pattern"
- Automate if possible: linter rule, pre-commit hook, CI check

If you can't automate, designate a reviewer whose job includes convention compliance.

## Rules

- Design for the artifact type that exists NOW — don't abstract for hypothetical types
- Audit before designing — never ignore what's already there
- If a convention already exists (Google Style Guide applies to your language), adopt
  it instead of designing your own
- Every convention must have a decision table with good/bad examples before publishing
- A convention without enforcement is aspirational, not operational — always pair with a checklist item or linter

## Common Mistakes

**Designing too late**: waiting until 200 names exist creates a migration cost. Design
before the third contributor joins.

**Over-constraining verbs**: a 3-verb tier (`create`, `get`, `delete`) doesn't cover
the real action space — contributors invent workarounds. Start with 8–12 verbs.

**No abbreviation policy**: "use common sense" produces half the team writing `usr` and
the other half writing `user`. Enumerate the allowed list.

**Convention without examples**: abstract rules without concrete right/wrong examples
are ignored. The decision table is not optional.

**Retrofitting inconsistently**: when fixing existing names to comply, either fix all
of them or none — partial fixes create a two-tier system harder to learn than the
original inconsistency.
