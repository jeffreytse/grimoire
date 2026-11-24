# Practice Profiles

Grimoire maintains skills and standards. Profiles are yours to own and share.

← [Back to the README](../README.md)

A profile is a named set of skills. Set `profiles = ["oop"]` in grimoire.toml and an AI assistant activates every skill tagged `oop` — no list to maintain, no file to create.

---

## How It Works

```toml
# grimoire.toml
[standards]
profiles = ["oop"]
```

When `apply-best-practice-profile` sees this, it resolves each profile name in this order:

1. `.grimoire/profiles/<name>.toml` — project-level file (highest priority)
2. `~/.grimoire/profiles/<name>.toml` — package-level file
3. `grimoire.toml` `[profiles.<name>]` — inline definition (see below)
4. `.grimoire/profiles/default.toml` — project-level fallback
5. `~/.grimoire/profiles/default.toml` — package-level fallback
6. Tag query — all installed skills where `tags` contains `<name>`

If no file matches, tags activate the skills. `profiles = ["oop"]` works immediately without creating any file.

Once a file is found, the profile itself applies **four composition layers** in order:

```
extends → tags → [[skills]] → exclude
```

Each layer is optional. A minimal profile with only `[[skills]]` entries works exactly as before.

---

## Multiple Profiles

Combine paradigms by listing multiple profiles:

```toml
[standards]
profiles = ["clean-architecture", "tdd"]
```

Skills from all profiles are merged. When the same skill appears in two profiles, the `priority` field on `[[skills]]` determines which wins — lower number = higher priority (default 50). If both have the same priority, the first profile in the array wins.

Duplicate skills are silently deduplicated — activated once regardless of how many profiles include them.

---

## User-Defined Profile Format

Create a file to curate exactly which skills activate. Profiles compose with four layers — each is optional:

```toml
# .grimoire/profiles/my-team.toml
name = "my-team"
description = "Our backend team's default practices"

# Layer 1: inherit all skills from other profiles (recursive)
extends = ["oop", "tdd"]

# Layer 2: bulk-activate all skills matching any of these tags
tags = ["go", "backend"]

# Layer 3: explicit additions — always included regardless of tags
[[skills]]
name = "apply-internal-convention"
priority = 10  # lower = higher priority; default = 50

# Layer 4: remove specific skills after all inclusions
exclude = ["apply-law-of-demeter"]
```

Resolution order within a profile: `extends` → `tags` → `[[skills]]` → `exclude`.
Explicit `[[skills]]` entries can override the priority of inherited skills by name.

Profiles are flat files — no subdirs. `extends` handles hierarchy.

### Fields

| Field | Type | Default | Description |
|---|---|---|---|
| `name` | string | filename stem | Profile identifier |
| `description` | string | `""` | One sentence, shown in `grimoire profile show` |
| `extends` | `string[]` | `[]` | Inherit skills from named profiles (resolved depth-first, cycle-safe). Resolves profile files only — use `tags` for tag-based bulk activation. |
| `tags` | `string[]` | `[]` | Bulk-activate all installed skills matching any tag. Works without a profile file. |
| `[[skills]].name` | string | required | Explicit skill addition |
| `[[skills]].priority` | int | `50` | Conflict weight — lower wins; overrides inherited priority for same skill |
| `exclude` | `string[]` | `[]` | Remove skills by name after all inclusions |

### Minimal profile (backward compatible)

```toml
# .grimoire/profiles/oop.toml
name = "oop"
description = "Object-oriented practices"

[[skills]]
name = "apply-solid-principles"

[[skills]]
name = "apply-domain-driven-design"
```

All new fields are optional. Existing profiles with only `[[skills]]` continue to work unchanged.

---

## Inline Profiles in grimoire.toml

Define a profile directly inside `grimoire.toml` — no separate file needed.
Use the `[profiles.<name>]` section. It supports the same fields as a profile TOML file,
plus optional compliance settings.

```toml
# grimoire.toml

[standards]
profiles = ["myteam"]

[profiles.myteam]
name = "myteam"  # optional — defaults to section key
description = "Our backend standards"
extends = ["oop"]
tags = ["backend", "go"]
compliance-threshold = 90
compliance-threshold-error = 0

[[profiles.myteam.skills]]
name = "apply-internal-convention"
priority = 1

exclude = ["apply-law-of-demeter"]
```

This is useful for:
- **All-in-one config** — single file defines the profile AND compliance settings
- **Shareable bundles** — commit `grimoire.toml` with profile definitions for teammates
- **Package profiles** — ship `grimoire.toml` alongside `profiles/` in a package bundle

When both a file and an inline definition exist for the same name, the file wins.

---

## The Ecosystem

Any git repo is a valid grimoire package bundle. Publish one and others install it with
`grimoire package add <name> <url>`.

### Bundle layout

```
my-standards/
  skills/          # optional — custom skill SKILL.md files
  profiles/        # optional — profile TOML files  (e.g. backend.toml, frontend.toml)
  grimoire.toml    # optional — compliance defaults + [profiles.*] inline definitions
```

All three directories are optional. A profiles-only bundle (`profiles/` + `grimoire.toml`,
no `skills/`) is fully supported.

### grimoire.toml in a bundle

The `grimoire.toml` in a package bundle serves two purposes:

1. **Compliance defaults** — `[standards.*]` sections set recommended thresholds shown as
   wizard defaults during `grimoire init`:
   ```toml
   [standards]
   compliance-threshold = 80
   compliance-threshold-error = 0

   [standards.finance]
   compliance-threshold = 90
   ```

2. **Inline profiles** — `[profiles.*]` sections define profiles without separate TOML files.
   Users get them all in one import:
   ```toml
   [profiles.backend]
   description = "Backend team defaults"
   tags = ["engineering"]
   compliance-threshold = 85

   [[profiles.backend.skills]]
   name = "apply-solid-principles"
   priority = 1
   ```

### How to publish

```bash
grimoire package validate   # check structure before pushing
grimoire package publish    # step-by-step checklist + opens github.com/new
```

Naming convention: `grimoire-package-<yourname>`. Add GitHub topic `grimoire-package` to make it discoverable.

After publishing, users install with one command:

```bash
grimoire package add myorg https://github.com/myorg/grimoire-package-myorg.git
```

Profiles from your `profiles/` and `grimoire.toml` then appear in `grimoire init` picker automatically.

### Project-level sharing

- **Commit `.grimoire/profiles/`** — teammates clone and get the profiles automatically
- **Commit `grimoire.toml`** — teammates get inline profiles and compliance settings
- **Use `default.toml`** — activates whenever an unknown profile is requested

This mirrors how `eslint-config-airbnb` or VS Code profiles work: the tool ships primitives,
the community ships bundles.

---

## Profiles and Compliance

The key insight: **the profile IS the context**. You never write `when:` conditions or trigger rules — activating a profile is the declaration that these skills apply to this work.

`grimoire check` reads the active profile, resolves its full skill list, and uses each skill's body as a compliance rubric:

```
active profile
  → skill list (via extends + tags + [[skills]] − exclude)
    → each SKILL.md body read as rubric
      → AI evaluates changed files against all rubrics
        → compliance result
```

No separate criteria schema. The skill body itself is the spec — `grimoire check` passes it directly to the AI as the evaluation standard.

### Profile detection order

`grimoire check` (and `grimoire watch`) resolves the active profile in this order:

1. `[standards] profiles` in `grimoire.toml` — explicit project declaration
2. `[standards] profiles` in `~/.config/grimoire/grimoire.toml` — user global
3. File-signal heuristics — `go.mod` → `go-engineering`, `package.json` → `node-engineering`, etc.

Run `grimoire status` to confirm which profile is active and how many skills it resolved.

### Compliance thresholds per profile

Profiles can ship compliance thresholds in their TOML file or inline settings:

```toml
# .grimoire/profiles/strict.toml
name = "strict"
description = "Production-grade engineering standards"

[[skills]]
name = "apply-solid-principles"

[[skills]]
name = "apply-domain-driven-design"
```

Thresholds live in `[standards.*]` settings, not the profile file — keeping skill selection and enforcement level separate. Set them in `grimoire.toml` to override package-provided defaults.

---

## Profiles vs Practices

Two separate keys in grimoire.toml — complementary, not redundant:

| | `profiles` | `practices` |
|---|---|---|
| Where | Top-level key | Inside a `[domain]` section |
| Scope | Global — activates skills across all domains | Scoped — applies to one domain only |
| Granularity | Bulk — all skills matching the name/tag | Explicit ordered list of specific skills |
| Use when | Declaring paradigm intent for the project | Fine-tuning skill priority for a specific domain |

```toml
# grimoire.toml

[standards]
profiles = ["oop"]  # load all oop-tagged skills globally

[standards.engineering.architecture]
practices = [       # override: in architecture, prefer these specific skills in this order
  "apply-solid-principles",
  "apply-law-of-demeter"
]
```

`profiles` sets the baseline. `practices` overrides per domain.

**Can I use `practices = ["OOP"]` instead of `profiles`?**

Yes — it works as a loose hint. The AI will lean toward OOP conventions based on its training. But it's less precise than `profiles = ["oop"]`:

- No guaranteed skill steps — depends on the AI's interpretation of "OOP"
- No priority rules when two practices conflict (profiles resolve that per skill)
- Not scoped to specific installed SKILL.md files

For best results, use both:

```toml
[standards]
profiles = ["oop"]   # activate oop skill bundle (specific skills)

[standards.engineering.architecture]
practices = ["OOP"]  # also pin as explicit preference for this domain
```

---

## CLI Reference

```bash
# Profile inspection
grimoire profile list                  # list profile files + active profiles from settings
grimoire profile show <name>           # resolve a profile and display its skill list
grimoire profile init <name>           # scaffold .grimoire/profiles/<name>.toml
grimoire config show                  # resolved config with per-key source
grimoire profile show <name>          # a profile's resolved skills

# Profile-driven compliance
grimoire check                         # run AI compliance check using active profile
grimoire check --via claude            # force a specific local AI agent
grimoire check --from-report           # read pre-generated report instead of running AI
grimoire check --ci                    # GitHub Actions annotations (auto on GITHUB_ACTIONS=true)
grimoire watch                         # re-run compliance on every file save (uses active profile)

# Project health
grimoire status                        # show active profile, skill count, last check result, staleness
```

`grimoire profile show` resolves the full composition chain (extends + tags + exclude) using installed skill sources. Skills from tag queries are shown even when no profile file exists.

---

## Examples

### Settings examples

| What you write in `[standards]` | What happens |
|---|---|
| `profiles = ["oop"]` | Tag query → all skills tagged `oop` |
| `profiles = ["functional"]` | Tag query → all skills tagged `functional` |
| `profiles = ["clean-architecture", "tdd"]` | Union of both; priority field resolves conflicts |
| `profiles = ["my-team"]` | Reads `profiles/my-team.toml` → applies 4-layer composition |
| `profiles = ["unknown"]` | Falls through to `default.toml` if present, else tag query (likely empty) |

### Profile file examples

**Tag-based profile** (no explicit skill list):
```toml
# .grimoire/profiles/backend-go.toml
name = "backend-go"
tags = ["go", "backend", "oop"]     # activate all matching skills automatically
exclude = ["apply-law-of-demeter"]  # too strict for this team
```

**Composition profile** (builds on shared profiles):
```toml
# .grimoire/profiles/my-team.toml
name = "my-team"
extends = ["oop", "tdd"]  # inherit from two profiles
tags = ["backend"]        # also pull in backend-tagged skills

[[skills]]
name = "apply-internal-convention" # team-specific addition
priority = 1              # activate before any inherited skills

exclude = ["apply-law-of-demeter"]
```


## Quick reference

Activate a named set of skills in one line — no list to maintain, no file to create.

```toml
# grimoire.toml
[standards]
profiles = ["oop"]  # activates every installed skill tagged "oop"
```

Grimoire resolves the name in this order, first match wins:

1. `.grimoire/profiles/<name>.toml` — project-level file
2. `~/.grimoire/profiles/<name>.toml` — user-level file
3. `.grimoire/profiles/default.toml` — project-level fallback
4. `~/.grimoire/profiles/default.toml` — user-level fallback
5. `[profiles.<name>]` in `grimoire.toml` — inline definition (no separate file needed)
6. Tag query — all installed skills where `tags` contains the name

If no file or inline definition exists, the tag query fires automatically. `profiles = ["oop"]` works without creating any file. Profiles go under `[standards]`; a top-level `profiles = [...]` is not read.

**Multiple profiles** — combine paradigms; first entry wins conflicts, duplicates are deduplicated:

```toml
[standards]
profiles = ["clean-architecture", "tdd"]  # clean-architecture wins if both include the same skill
```

**Custom profile** — curate a team-specific subset when the tag set is too broad:

```toml
# .grimoire/profiles/my-team.toml
name = "my-team"
description = "Our backend team's default practices"

[[skills]]
name = "apply-solid-principles"

[[skills]]
name = "apply-domain-driven-design"
```

Commit `.grimoire/profiles/` to share standards across the team. Publish as a gist or repo (`grimoire-profile-<name>`) for the community.

**Inline profile** — define directly in `grimoire.toml`, no separate file:

```toml
# grimoire.toml
[standards]
profiles = ["my-team"]

[profiles.my-team]
description = "Our backend team's default practices"
extends = ["oop"]

[[profiles.my-team.skills]]
name = "apply-solid-principles"

[[profiles.my-team.skills]]
name = "apply-domain-driven-design"

exclude = ["apply-law-of-demeter"]
```

If a `profiles/my-team.toml` file also exists for the same name, the file wins.

**`profiles` vs `practices`** — `profiles` activates skill bundles globally; `practices` is a domain-scoped explicit list. → [Full comparison](./profiles.md#profiles-vs-practices)

→ [Full profiles guide](./profiles.md) — resolution order, conflict handling, sharing profiles
