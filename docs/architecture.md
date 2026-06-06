# grimoire Architecture

grimoire has three layers: **skills** (content), **meta skills** (framework logic), and the **plugin system** (delivery to AI agents). This document describes how each layer works and how they fit together.

---

## Overview

```
┌─────────────────────────────────────────────────┐
│                   AI Agent                      │
│   (Claude Code / Codex / Cursor / Gemini / ...)  │
└──────────────────────┬──────────────────────────┘
                       │ plugin install
┌──────────────────────▼──────────────────────────┐
│               Plugin System                     │
│   .claude-plugin/ · .codex-plugin/ · gemini-    │
│   extension.json · marketplace.json             │
└──────────────────────┬──────────────────────────┘
                       │ loads
┌──────────────────────▼──────────────────────────┐
│                  Skills                         │
│   skills/<domain>/<subdomain>/skills/           │
│   <skill-name>/SKILL.md                         │
└──────────────────────┬──────────────────────────┘
                       │ managed by
┌──────────────────────▼──────────────────────────┐
│               Meta Skills                       │
│   suggest-best-practice · review-best-practice-skill · write-best-practice-skill │
│   plan-best-practice-solution · review-best-practice-fit · ...     │
└─────────────────────────────────────────────────┘
```

---

## Skill format: SKILL.md

Every skill is a single `SKILL.md` file. The format is fixed:

```markdown
---
name: verb-first-kebab-case
description: Use when <triggering conditions>.
source: <institution, standard body, or top-tier companies>
tags: [problem-keyword, tool-method, role-context, outcome]
---

# Skill Title

One-sentence purpose.

## Why This Is Best Practice

**Adopted by:** [specific companies or institutions]
**Impact:** [measurable outcome — a number or named study]
**Why best:** [why this over the named alternative]

Sources: [verifiable citations]

## Steps

### 1. Step name
[concrete, immediately executable instruction]

...
```

**Frontmatter fields** are used by `suggest-best-practice` for automatic routing. The `tags` field is the primary matching signal — it must cover all four axes: problem keyword, tool/method, role/context, and outcome.

**"Why This Is Best Practice"** is grimoire's quality signal. Every skill must prove it belongs: named adopters, measurable impact with evidence, and an explicit comparison to alternatives. Vague claims ("many companies use this", "significantly improves quality") are rejected at review.

**Size target: 50–300 lines.** Under 50 is too shallow. Over 300 is two skills.

Full format specification: [STANDARD.md](../STANDARD.md).

---

## Directory layout

```
skills/
  <domain>/
    .claude-plugin/
      plugin.json               ← domain-level manifest
    <subdomain>/
      .claude-plugin/
        plugin.json             ← subdomain-level manifest
      skills/
        <skill-name>/
          SKILL.md
```

Example — the grimoire meta domain:

```
skills/
│   ├── grimoire/                         # Meta skills (the framework's nervous system)
│   │   └── skills/
│   │       ├── analyze-problem/
│   │       ├── discover-best-practices/
│   │       ├── intercept-best-practice/
│   │       ├── suggest-best-practice/
│   │       ├── plan-best-practice-solution/
│   │       ├── apply-best-practice-tree/
│   │       ├── review-best-practice-fit/
│   │       ├── compare-best-practices/
│   │       ├── audit-applied-best-practices/
│   │       ├── explain-best-practice/
│   │       ├── adapt-best-practice/
│   │       ├── teach-best-practice/
│   │       ├── pin-best-practice-preference/
│   │       ├── write-best-practice-skill/
│   │       ├── review-best-practice-skill/  # Quality gate
│   │       ├── revise-best-practice-skill/
│   │       ├── audit-best-practice-domain/
│   │       ├── deprecate-best-practice-skill/
│   │       └── design-best-practice-domain/
```

Example — the reference skill:

```
skills/
  engineering/
    .claude-plugin/plugin.json
    development/
      .claude-plugin/plugin.json
      skills/
        propose-conventional-commit/
          SKILL.md
```

**Domain plugin.json** (`skills/<domain>/.claude-plugin/plugin.json`):

```json
{
  "name": "grimoire-engineering",
  "description": "Engineering skills: development, testing, architecture, ...",
  "version": "0.1.0",
  "author": { "name": "Jeffrey Tse", "email": "jeffreytse.mail@gmail.com" },
  "homepage": "https://github.com/jeffreytse/grimoire",
  "repository": "https://github.com/jeffreytse/grimoire",
  "license": "MIT",
  "skills": [
    "./development/skills",
    "./testing/skills",
    "./architecture/skills"
  ]
}
```

`skills` is an array of paths — each pointing to a subdomain's `skills/` directory. The agent auto-discovers all `SKILL.md` files under each path.

**Subdomain plugin.json** (`skills/<domain>/<subdomain>/.claude-plugin/plugin.json`):

```json
{
  "name": "grimoire-engineering-development",
  "description": "Development skills: coding, implementation, code review, debugging.",
  "version": "0.1.0",
  "author": { "name": "Jeffrey Tse", "email": "jeffreytse.mail@gmail.com" },
  "homepage": "https://github.com/jeffreytse/grimoire",
  "repository": "https://github.com/jeffreytse/grimoire",
  "license": "MIT",
  "skills": "./skills"
}
```

`skills` is a string (not an array) pointing to the `skills/` directory. The agent auto-discovers all `SKILL.md` files in that directory.

---

## Plugin system

### Claude Code

Claude Code reads `.claude-plugin/plugin.json` when you run `/plugins add`. The `skills` field tells it where to look for `SKILL.md` files. Skills are loaded into context and become available as slash commands (if the skill name matches a command pattern) or as background knowledge the agent applies when relevant.

Install commands:
```bash
/plugins add github:jeffreytse/grimoire                              # all domains
/plugins add github:jeffreytse/grimoire/skills/engineering           # one domain
/plugins add github:jeffreytse/grimoire/skills/engineering/development  # one subdomain
```

### Marketplace

`.claude-plugin/marketplace.json` is the registry of all installable units. Each entry has a `name`, `source` (with GitHub path), and `description`. This is what powers the install command resolution.

```json
{
  "name": "grimoire-engineering-development",
  "source": {
    "source": "github",
    "repo": "jeffreytse/grimoire",
    "path": "skills/engineering/development"
  },
  "description": "Software development: coding, implementation, review, debugging"
}
```

---

## Multi-agent support

grimoire ships plugin configurations for five AI agents:

| Agent | Config file | Loading mechanism |
|-------|-------------|-------------------|
| Claude Code | `.claude-plugin/plugin.json` | `/plugins add` reads plugin.json, auto-discovers SKILL.md |
| Codex | `.codex-plugin/plugin.json` | Same structure as Claude Code |
| Cursor | `.cursor-plugin/plugin.json` | Same structure |
| OpenCode | `.opencode/plugins/grimoire.js` | ESM module — injects AGENTS.md into first user message via transform hook; registers skills paths via config hook |
| Gemini CLI | `gemini-extension.json` | Points to `GEMINI.md`, which is loaded as context |

All agent-facing content (CLAUDE.md, AGENTS.md, GEMINI.md) describes grimoire's available skills and how to invoke them. Agent-specific docs live at the repo root.

---

## Meta skills: the framework's nervous system

grimoire is self-managing. The meta skills in `skills/grimoire/` run the framework itself:

**User-facing meta skills** — help users find and apply practices:

| Skill | What it does |
|-------|-------------|
| `analyze-problem` | Clarifies an ill-defined problem through structured questioning, then maps the problem space and surfaces possible routes |
| `discover-best-practices` | Surfaces available practices for a domain before the user has a specific problem — grouped by subdomain, framed as gaps |
| `intercept-best-practice` | Proactively fires before a task starts — matches the most relevant practice and offers to apply it before gaps emerge |
| `suggest-best-practice` | Universal entry point — classifies any situation and routes to the matching skill or install command |
| `plan-best-practice-solution` | Decomposes multi-domain problems into sequenced skill applications (MECE methodology) |
| `apply-best-practice-tree` | Recursively decomposes a complex single-domain problem into sub-problems, matching each to the best installed skill |
| `review-best-practice-fit` | Evaluates an existing solution against best practices — ALIGNED/PARTIAL/MISSING per practice, prioritized fix list |
| `compare-best-practices` | Side-by-side comparison when multiple practices apply — produces a table and a clear recommendation |
| `audit-applied-best-practices` | Audits existing work for which practices were applied, which are missing, and what gaps to close |
| `explain-best-practice` | Educational deep-dive: problem → origin → evidence → mechanism → failure modes → misconceptions |
| `adapt-best-practice` | Adapts a practice for specific constraints, classifying each step as Core/Adjustable/Optional |
| `teach-best-practice` | Produces audience-tailored talking points, brief, or slide outline for sharing a practice with others |
| `pin-best-practice-preference` | Saves a practice preference at session, project, or global level for a domain or subdomain |

**Contributor meta skills** — manage the library itself:

| Skill | What it does |
|-------|-------------|
| `write-best-practice-skill` | Guides authoring a new SKILL.md from scratch |
| `review-best-practice-skill` | Evaluates a SKILL.md — PASS/NEEDS-REVISION/REJECT with specific findings |
| `revise-best-practice-skill` | Applies review findings to an existing skill without touching passing sections |
| `audit-best-practice-domain` | Batch-evaluates all skills in a domain, calls review-best-practice-skill per file |
| `deprecate-best-practice-skill` | Marks an outdated skill for removal with migration guidance |
| `design-best-practice-domain` | Architects a new domain: directory structure, plugin files, marketplace entries, seed skills |

**Why this matters:** `review-best-practice-skill` runs the exact same criteria as `STANDARD.md`. As long as both are maintained together, the quality standard is self-enforcing — it cannot drift between the written standard and what's actually checked. Any agent using grimoire's meta skills applies the same bar every time.

---

## Quality gate

The contribution pipeline is:

```
write-best-practice-skill  →  review-best-practice-skill  →  revise-best-practice-skill  →  review-best-practice-skill  →  PR merge
                    ↑                                    |
                    └────────────── loop ────────────────┘
```

For new domains:

```
design-best-practice-domain  →  write-best-practice-skill (seed skills)  →  audit-best-practice-domain  →  PR merge
```

`review-best-practice-skill` uses the Fagan Inspection method (IBM 1976) — a deterministic checklist that produces consistent verdicts regardless of which agent applies it. PASS requires all criteria to be ✅. A REJECT on any frontmatter field blocks merge.

The full criteria and checklist are in [STANDARD.md](../STANDARD.md).
