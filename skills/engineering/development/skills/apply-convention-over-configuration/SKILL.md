---
name: apply-convention-over-configuration
description: Use when designing a project structure, framework, or tooling setup and deciding how much explicit configuration to require from developers.
source: DHH/Rails, "Convention over Configuration", 2005 — Widely adopted at Spring Boot, Django, Laravel, Angular CLI, Next.js, NestJS
tags: [configuration, project-structure, developer, boilerplate-reduction, onboarding-speed]
verified: true
---

# Apply Convention over Configuration

Encode the common case as a default so developers only configure deviations.

## Why This Is Best Practice

**Adopted by:** Every major web framework post-2005 — Rails, Spring Boot, Django,
Laravel, Angular CLI, Next.js, NestJS — all ship opinionated defaults rather than
requiring upfront configuration.
**Impact:** Rails vs. J2EE (2005) demonstrated a 60–80% reduction in boilerplate
configuration files for equivalent functionality. New developer onboarding time drops
from days (explicit config systems) to hours (convention-based) because the project
structure is already familiar.
**Why best:** Explicit configuration requires each team to re-decide settled questions
(where do models live? how are routes named?) — producing N divergent answers across
projects. Conventions encode collective wisdom once; every subsequent project inherits
it for free. Deviations are immediately visible because they require explicit config.

Sources: DHH, "Rails Doctrine" (rubyonrails.org/doctrine), Spring Boot reference docs,
Django design philosophies (djangoproject.com/design-philosophies)

## Steps

### 1. Identify the decisions to conventionalize

List the recurring structural decisions in your domain:
- File/directory locations (where do models, controllers, tests live?)
- Naming patterns (how are files, classes, routes named relative to each other?)
- Wiring rules (how does a file in location X get auto-discovered?)

### 2. Encode defaults

Make the common case require zero configuration:
- File named `UserController` → automatically handles `/users` routes
- File in `models/` → automatically loaded as a model
- Test file named `*_test.rb` → automatically discovered by the test runner

If it requires a config entry to do the standard thing, the convention is missing.

### 3. Config only for deviations

When a project departs from the convention, require an explicit declaration at the
deviation point — not a global config file that re-describes the default:

```ruby
# Bad — re-states the convention explicitly
class User < ApplicationRecord
  self.table_name = "users"  # this IS the convention; remove it
end

# Good — only configure when deviating
class User < ApplicationRecord
  self.table_name = "people"  # explicit deviation
end
```

### 4. Document the conventions, not the config

Write a single reference listing all conventions. New developers read this once.
Don't document configuration options for cases the convention already handles.

### 5. Validate discoverability

Test: give a new team member a task that requires finding an existing file.
If they need a guided tour or a map, the convention isn't strong enough.

## Rules

- One convention per decision — conflicting conventions force developers to pick, which defeats the purpose
- Conventions must be consistent: if `UserController` lives in `controllers/`, `OrderController` must also live in `controllers/` — partial conventions are worse than none
- Don't invent conventions that conflict with the framework you're building on — two competing convention layers create two mental models

## Common Mistakes

**Documenting defaults as if they were config options.** If a developer needs to write
`table_name = "users"` to get the default behavior, the convention isn't working.

**Convention without discoverability.** A convention that can't be discovered by reading
the code is magic — add a visible hook or naming pattern that makes the wiring obvious.

**Over-conventionalizing.** Genuine one-off decisions don't need a convention.
Only standardize what recurs across every project in the domain.

## When NOT to Use

- Novel system with no established community conventions — inventing conventions before patterns emerge creates the wrong defaults
- Small scripts or one-off tools where structure doesn't pay off
- Teams with strong existing conventions that differ — forcing a new convention without buy-in creates resistance and inconsistency
