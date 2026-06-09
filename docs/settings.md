# Settings

`settings.toml` is where you pin preferences, activate profiles, and set compliance thresholds. Grimoire merges multiple files — global, project-shared, project-personal, and session — so individuals can override team defaults without touching shared config.

---

## Override Hierarchy

Grimoire checks these sources in order. First match per key wins.

| Priority | Level | File | Scope | Committed? |
|----------|-------|------|-------|------------|
| 1st | Session | in-memory | current session only | — |
| 2nd | Project personal | `.grimoire/settings.local.toml` | this project, this machine | gitignored |
| 3rd | Project shared | `.grimoire/settings.toml` | this project, all team members | yes |
| 4th | Global | `~/.config/grimoire/settings.toml` | all projects on this machine | — |

Each key resolves independently. Project-shared can set `practices` while global sets `profiles` — both apply without conflict. A key set at level 2 does not suppress unrelated keys at levels 3 or 4.

**Setting up gitignore:**

```bash
echo '.grimoire/settings.local.toml' >> .gitignore
```

Or add to an existing `.gitignore`:

```
# grimoire personal settings
.grimoire/settings.local.toml
```

---

## Key Reference

| Key | Where | Type | Description |
|-----|-------|------|-------------|
| `profiles` | top-level | `string[]` | Activate skill bundles by name or tag — resolved in [profile lookup order](./profiles.md#how-it-works) |
| `practices` | `[domain]` or `[domain.subdomain]` | `string[]` | Ordered skill/practice hints for that domain — first entry wins conflicts |
| `fallback` | `[domain]` or `[domain.subdomain]` | `"ask"` \| `"skip"` | Behavior when no matching skill is found: `"ask"` prompts for clarification, `"skip"` silently moves on |
| `compliance-threshold` | `[domain]` or `[domain.subdomain]` | number | Minimum overall criteria coverage % to pass `check-best-practice-compliance` |
| `compliance-threshold-error` | `[domain]` or `[domain.subdomain]` | number | Maximum allowed error-severity violations — `0` means none allowed |

Domain sections can be top-level (`[engineering]`) or subdomain (`[engineering.architecture]`). More specific section wins — `[engineering.architecture]` overrides `[engineering]` for architecture work.

---

## Examples

### Global defaults

Applies to every project on this machine. Good for personal style preferences that span all work.

```toml
# ~/.config/grimoire/settings.toml
profiles = ["oop"]          # OOP skills active everywhere by default

[engineering]
practices = ["Google Engineering Practices"]

[engineering.testing]
practices = ["apply-tdd"]

[finance]
practices = ["CFA Institute standards"]

[law]
practices = ["ABA Model Rules"]
fallback = "ask"            # always ask before applying a law skill — don't guess jurisdiction
```

---

### Project shared config

Committed to the repo. Every team member gets the same baseline.

```toml
# .grimoire/settings.toml
profiles = ["clean-architecture", "tdd"]   # two profiles — clean-architecture wins conflicts

[engineering]
practices = ["Google Engineering Practices"]
compliance-threshold = 75           # 75% overall coverage to pass CI

[engineering.architecture]
practices = [
  "apply-solid-principles",
  "apply-law-of-demeter",
  "apply-domain-driven-design"
]
compliance-threshold = 80           # architecture held to a higher bar
compliance-threshold-error = 0      # zero tolerance for error-severity findings

[engineering.testing]
practices = ["apply-tdd"]
fallback = "skip"                   # if no exact test skill matches, skip silently
```

---

### Personal override

Gitignored. Use to deviate from team settings on your machine without touching shared config.

```toml
# .grimoire/settings.local.toml
[engineering]
practices = ["Facebook Engineering Standards"]   # personal preference overrides team default
fallback = "ask"                                 # prefer to be prompted rather than auto-routed
```

---

### Multi-domain project (startup)

A project that spans engineering, legal, and finance — each domain has its own standard:

```toml
# .grimoire/settings.toml
profiles = ["lean-startup"]

[engineering]
practices = ["Google Engineering Practices"]
compliance-threshold = 70           # early-stage: 70% is the bar

[law.contracts]
practices = ["review-saas-contract", "audit-gdpr-compliance"]
compliance-threshold-error = 0      # legal errors are always blocking

[finance.personal-finance]
practices = ["CFA Institute standards"]

[writing.copywriting]
practices = ["apply-aida-framework", "write-value-proposition"]
```

---

### Enforcement-only project (legacy codebase)

Existing codebase with tech debt — set a lower threshold now, raise it over time:

```toml
# .grimoire/settings.toml
[engineering]
practices = [
  "apply-solid-principles",
  "apply-kiss-principle"
]
compliance-threshold = 40           # current baseline — raise 5% per sprint
compliance-threshold-error = 0      # errors always block regardless of overall score
```

---

### Session pin (via AI — no TOML edit needed)

Grimoire can pin a preference for the current session without writing any file:

```
User: For this session, use Facebook's engineering practices instead of Google's.

Claude: Pinned for this session: Facebook Engineering Standards → [engineering].
        Reverts when session ends. Save permanently? [y] global  [p] project
```

The session pin overrides all files. Nothing is written to disk until you confirm.

---

## How Layers Compose

Given this setup:

```
~/.config/grimoire/settings.toml          profiles = ["oop"]
                                          [engineering] practices = ["Google"]

.grimoire/settings.toml                   [engineering.architecture] practices = ["SOLID", "LoD"]
                                          [engineering.architecture] compliance-threshold = 80

.grimoire/settings.local.toml             [engineering] fallback = "ask"
```

The resolved spec for `engineering.architecture` work:

| Key | Value | Source |
|-----|-------|--------|
| `profiles` | `["oop"]` | global |
| `[engineering].practices` | `["Google"]` | project-shared (no local override for this key) |
| `[engineering.architecture].practices` | `["SOLID", "LoD"]` | project-shared |
| `[engineering.architecture].compliance-threshold` | `80` | project-shared |
| `[engineering].fallback` | `"ask"` | project-personal |

`[engineering.architecture]` practices do not fall back to `[engineering]` practices — they replace them for architecture-scoped work. Other keys not set at the subdomain level do inherit from the domain level.

---

## Guided Management

Use skills instead of editing TOML directly:

- `/configure-grimoire` — view the resolved spec for any domain, edit and validate settings interactively
- `/pin-best-practice-preference` — pin a preference when two practices conflict; choose session, project, or global scope
- `/resolve-best-practice-conflict` — resolve contradictions between two installed skills and record the priority automatically

---

## Related

- [profiles.md](./profiles.md) — how `profiles` resolves, custom profile files, `profiles` vs `practices`
- [bpdd.md](./bpdd.md) — compliance thresholds in the BPDD cycle, CI gating, incremental mode
