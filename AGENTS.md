# grimoire

Multi-domain skill collection for AI assistants. Skills work across Claude Code, Codex,
Copilot, Gemini CLI, OpenCode, and more.

## Skill Path Convention

All skill paths follow this pattern:
```
<domain>/<subdomain>/skills/<skill-name>/SKILL.md
```

If a task matches a skill, invoke it. Never implement directly when a skill applies.

## Domains

| Domain | Sub-domains |
|--------|-------------|
| `meta` | *(flat — no sub-domains)* |
| `engineering` | development, frontend, architecture, testing, reliability, devops, cloud, networking, security, data, ai, hardware, mobile, performance, project-management, product, documentation |
| `business` | strategy, operations, leadership, entrepreneurship, hr |
| `science` | biology, physics, chemistry, mathematics, earth-science, astronomy |
| `health` | fitness, nutrition, mental-health, sleep, medicine |
| `writing` | creative, technical, copywriting, academic, journalism |
| `design` | ui-ux, graphic, branding, motion, product |
| `marketing` | seo, content, social-media, paid-ads, growth, analytics |
| `finance` | personal-finance, investing, accounting, real-estate, corporate |
| `law` | contracts, ip, employment, privacy, corporate |
| `film` | cinematography, directing, editing, screenwriting, production |
| `music` | composition, production, mixing, theory, performance |
| `art` | drawing, painting, digital-art, illustration, color-theory |
| `sports` | training, coaching, nutrition, tactics, recovery |
| `psychology` | cognitive, behavioral, social, clinical, positive |
| `language` | learning, linguistics, translation, communication |
| `education` | curriculum, teaching, e-learning, assessment, learning-science |
| `cooking` | techniques, baking, flavor, nutrition, world-cuisine |
| `travel` | planning, budgeting, cultural, adventure |
| `home` | renovation, interior-design, gardening, organization, smart-home |
| `parenting` | infant, toddler, school-age, teen |
| `automotive` | maintenance, troubleshooting, buying, modifications, ev |
| `pets` | dogs, cats, training, nutrition, health |
| `fashion` | styling, wardrobe, design, sustainability, accessories |
| `environment` | sustainability, ecology, climate, energy, policy |
| `photography` | composition, lighting, editing, genres |
| `productivity` | time-management, habits, focus, goals, tools |

## Installation

```bash
# All skills
./scripts/install.sh

# One domain
./scripts/install.sh --domain engineering

# One sub-domain
./scripts/install.sh --domain engineering --subdomain development

# One skill
./scripts/install.sh --skill engineering/development/propose-commit

# All agents
./scripts/install.sh --target all
```

Or via Claude Code plugin system:
```bash
/plugins add github:jeffreytse/grimoire                         # all
/plugins add github:jeffreytse/grimoire/engineering             # one domain
/plugins add github:jeffreytse/grimoire/engineering/development # one sub-domain
```

## Skill Format

```markdown
---
name: skill-name
description: Use when <triggering conditions>
---

# Skill Name

One-sentence purpose.

## Steps
...
```

See [STANDARD.md](./STANDARD.md) for the full quality standard.
