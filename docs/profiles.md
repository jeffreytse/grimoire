# Practice Profiles

Grimoire maintains skills and standards. Profiles are yours to own and share.

A profile is a named set of skills. Set `profiles = ["oop"]` in settings.toml and an AI assistant activates every skill tagged `oop` — no list to maintain, no file to create.

---

## How It Works

```toml
# .grimoire/settings.toml
profiles = ["oop"]
```

When `apply-best-practice-profile` sees this, it resolves each profile name in this order:

1. `.grimoire/profiles/<name>.toml` — project-level file (highest priority)
2. `~/.grimoire/profiles/<name>.toml` — user-level file
3. `.grimoire/profiles/default.toml` — project-level fallback
4. `~/.grimoire/profiles/default.toml` — user-level fallback
5. Tag query — all installed skills where `tags` contains `<name>`

If no file matches, tags activate the skills. `profiles = ["oop"]` works immediately without creating any file.

---

## Multiple Profiles

Combine paradigms by listing multiple profiles:

```toml
profiles = ["clean-architecture", "tdd"]  # clean-architecture wins conflicts
```

**Array order = conflict priority.** When two profiles include skills with conflicting guidance, the first profile in the array wins.

**Duplicate skills** (same skill name in two profiles) are silently deduplicated — activated once. The activation summary notes which profiles both requested it.

---

## User-Defined Profile Format

Create a file to curate exactly which skills activate — useful when the tag set is too broad, or you want a team-specific subset:

```toml
# .grimoire/profiles/my-team.toml
name = "my-team"
description = "Our backend team's default practices"

[[skills]]
name = "apply-solid-principles"

[[skills]]
name = "apply-domain-driven-design"
```

Flat files only. No `extends`, no subdirs.

### Fields

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | Profile identifier |
| `description` | string | no | One sentence, when to use this profile |
| `skills[].name` | string | yes | Skill name (matches frontmatter `name:`) |

---

## The Ecosystem

Grimoire ships no built-in profiles — the tag system covers the common cases. But you can:

- **Share your team profile** — commit `.grimoire/profiles/` to your repo
- **Publish a community profile** — open-source as a gist or repo (`grimoire-profile-oop`, etc.)
- **Use `default.toml`** — put team-wide defaults there; activates whenever any unknown profile is requested

This mirrors how `eslint-config-airbnb` or VS Code profiles work: the tool ships primitives, the community ships bundles.

---

## Profiles vs Practices

Two separate keys in settings.toml — complementary, not redundant:

| | `profiles` | `practices` |
|---|---|---|
| Where | Top-level key | Inside a `[domain]` section |
| Scope | Global — activates skills across all domains | Scoped — applies to one domain only |
| Granularity | Bulk — all skills matching the name/tag | Explicit ordered list of specific skills |
| Use when | Declaring paradigm intent for the project | Fine-tuning skill priority for a specific domain |

```toml
# .grimoire/settings.toml

profiles = ["oop"]              # load all oop-tagged skills globally

[engineering.architecture]
practices = [                   # override: in architecture, prefer these specific skills in this order
  "apply-solid-principles",
  "apply-law-of-demeter"
]
```

`profiles` sets the baseline. `practices` overrides per domain.

**Can I use `practices = ["OOP"]` instead of `profiles`?**

Yes — it works as a loose hint. The AI will lean toward OOP conventions based on its training. But it's less precise than `profiles = ["oop"]`:

- No guaranteed skill steps — depends on the AI's interpretation of "OOP"
- No dedup or conflict logic between skills
- Not scoped to specific installed SKILL.md files

For best results, use both:

```toml
profiles = ["oop"]           # activate oop skill bundle (specific skills)

[engineering.architecture]
practices = ["OOP"]          # also pin as explicit preference for this domain
```

---

## Examples

| What you write | What happens |
|---|---|
| `profiles = ["oop"]` | Tag query → all skills tagged `oop` |
| `profiles = ["functional"]` | Tag query → all skills tagged `functional` |
| `profiles = ["clean-architecture", "tdd"]` | Union of both; `clean-architecture` wins conflicts |
| `profiles = ["my-team"]` | Reads `profiles/my-team.toml` → activates listed skills |
| `profiles = ["unknown"]` | Falls through to `default.toml` if present, else tag query (likely empty) |
