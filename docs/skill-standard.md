# The Grimoire Skill Standard

What every official skill must meet.

← [Back to the README](../README.md)

grimoire-core maintains an open standard for AI agent skill quality — freely adoptable by any skill library. Community packages are encouraged to follow it; grimoire installs any git repo regardless. This quality gate applies to grimoire-core contributions only.

Every skill must pass `review-best-practice-skill` before merge:

| Criterion      | Requirement                           | Rejection example                |
| -------------- | ------------------------------------- | -------------------------------- |
| **Adopted by** | Named organizations or institutions   | "Many top companies"             |
| **Impact**     | Cited study or % number               | "Significantly improves quality" |
| **Steps**      | Immediately executable                | Abstract theory or advice        |
| **Scope**      | One concept per skill                 | "Nutrition and training program" |
| **Source**     | External institution or standard body | Internal opinion                 |

**`grimoire validate`** is the reference implementation of the standard — a static linter for SKILL.md files:

```bash
grimoire validate                        # validate all skills in current directory
grimoire validate skills/                # validate a specific skills tree
grimoire validate skills/engineering/    # validate one domain
grimoire validate path/to/SKILL.md       # validate a single skill
grimoire validate skills/ --fix          # fix errors via AI (requires local agent or API key)
grimoire validate skills/ --strict       # treat warnings as errors (CI gate)
grimoire validate skills/ --no-duplicates  # skip semantic near-duplicate detection
grimoire validate --test-schema schema/tests  # run conformance test suite
```

Checks run on every skill: required frontmatter fields, tag format, description quality, name uniqueness, and semantic near-duplicate detection across the library.

Skills that load or score *other* skills (the meta workflows) must also work when grimoire serves skills over MCP: the standard's [Skill Access](https://github.com/jeffreytse/grimoire-core/blob/main/STANDARD.md#skill-access-conditional) section tells them to use the Skill tool in skills mode and `grimoire_get_skill` / `grimoire_search_skills` in MCP mode.

→ [Read the full standard](https://github.com/jeffreytse/grimoire-core/blob/main/STANDARD.md) · [Adopt this standard](https://github.com/jeffreytse/grimoire-core/blob/main/STANDARD.md#adopting-this-standard)

See [CONTRIBUTING.md](https://github.com/jeffreytse/grimoire-core/blob/main/CONTRIBUTING.md) to submit a skill.
