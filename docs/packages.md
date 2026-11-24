# Skills as packages

How grimoire installs, pins and shares skill libraries.

← [Back to the README](../README.md)

Every AI coding tool ships its own config format — `CLAUDE.md`, `.cursorrules`, `AGENTS.md`. Developers copy-paste and drift. Grimoire solves this with a package manager model: declare the practices you need in `grimoire.toml`, pin versions in the ref (`@0.1.0`), and serve them to every agent. `grimoire.lock` records exactly what was installed.

Any git repo is a valid package. grimoire-core is the official one — curated, cited, peer-reviewed. Beyond distributing skills, grimoire closes the loop: your agent applies the declared practices, and `grimoire check` verifies the work against them.

```toml
# grimoire.toml — commit this to your repo
[package]
name    = "my-project"
version = "0.1.0"

[dependencies]
"jeffreytse/grimoire-core"             = "*"   # official package — all 1,300+ skills
"jeffreytse/grimoire-core:engineering" = "*"   # one domain only
"jeffreytse/grimoire-core@0.1.0"       = "*"   # pinned to a specific release
"mycompany/internal-skills"            = "*"   # private company package
```

**In MCP mode** (recommended), nothing is linked into the agent. The server serves your libraries directly:

```bash
grimoire package add my-team mycompany/internal-skills         # add a library for yourself (user config)
grimoire update                                                # fetch libraries, incl. a project's [dependencies]
grimoire package list                                          # show your libraries
```

A team makes a library its standard by declaring it under `[dependencies]` in the project's `grimoire.toml` (above). Teammates run `grimoire update` to fetch it, and the assistant prefers its practices.

**In skills mode** (agents without MCP), `grimoire install` also links the matched skills into each agent:

```bash
grimoire install jeffreytse/grimoire-core                      # all skills from official package
grimoire install "jeffreytse/grimoire-core:engineering"        # engineering domain only
grimoire install "jeffreytse/grimoire-core@0.1.0"              # pinned to a specific release
grimoire install "jeffreytse/grimoire-core@0.1.0:health"       # pinned release + health domain
grimoire install "gitlab.com/mycompany/internal-skills"        # any git repo as a package
grimoire install                                               # (re-)install all from grimoire.toml
grimoire update                                                # update all packages to latest
grimoire list                                                  # show installed packages and skills
grimoire uninstall mycompany/internal-skills                   # remove a package
```

**Ref syntax:** `owner/repo[@tag][:path]` — the `:path` suffix is a Standard Glob pattern (doublestar) matched against each skill's domain path. `**` matches any depth; `*` matches within one segment. Examples: `engineering/**` (all engineering skills), `health/sleep/**` (one subdomain only), `**/development/**` (development subdomain in any domain).

| Package                                                        | Type                | Quality gate                     |
| -------------------------------------------------------------- | ------------------- | -------------------------------- |
| [`grimoire-core`](https://github.com/jeffreytse/grimoire-core) | Official            | STANDARD.md peer review required |
| Any git repo                                                   | Community / private | Owner's discretion               |

grimoire-core is the official package — curated, cited, and reviewed against STANDARD.md. But grimoire manages _any_ package without restriction. Mix and match.

**Community packages** install identically to grimoire-core. We encourage following the [grimoire skill standard](https://github.com/jeffreytse/grimoire-core/blob/main/STANDARD.md) for quality and interoperability — but any git repo works. Publish as `grimoire-<name>` for discoverability:

| Package example             | Focus                             |
| --------------------------- | --------------------------------- |
| `yourorg/grimoire-fintech`  | fintech-specific practices        |
| `yourorg/grimoire-medical`  | clinical protocols, FDA workflows |
| `yourorg/grimoire-legal-us` | US jurisdiction guides            |
| `yourorg/internal-skills`   | your team's private practices     |

The practices in your installed packages define what `grimoire check` validates against — declare the package, get the linter rules for free.
