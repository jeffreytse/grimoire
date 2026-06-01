# grimoire

> The world's professional best practices of every field — your expert consultant, anywhere, anytime.

[![GitHub Stars](https://img.shields.io/github/stars/jeffreytse/grimoire?style=social)](https://github.com/jeffreytse/grimoire)
[![Contributors](https://img.shields.io/github/contributors/jeffreytse/grimoire)](https://github.com/jeffreytse/grimoire/graphs/contributors)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Works with](https://img.shields.io/badge/works%20with-Claude%20·%20Codex%20·%20Cursor%20·%20Gemini%20·%20OpenCode-blue)](#agent-support)
[![Skill Validation](https://github.com/jeffreytse/grimoire/actions/workflows/validate.yml/badge.svg)](https://github.com/jeffreytse/grimoire/actions/workflows/validate.yml)

A lawyer charges $500/hr. A doctor charges $300 per visit. A financial advisor, a software architect, a structural engineer — each expert is expensive, scarce, and unavailable at 2am when you need them.

grimoire encodes what they know: the professional best practices of every field, adopted by top-tier institutions, with cited evidence. Install once. Your expert consultant in every discipline — available anywhere, anytime, for life.

## Philosophy

Grimoire encodes the world's best practitioners' actual methods — proven at scale, cited, and actionable. Every skill is a verb you can use right now. Built by anyone who has mastered their craft. Free, open source, for every field, forever.

- **The world's best practices belong to everyone.**
- **Skills are verbs.**
- **Every claim must be proven.**
- **Consensus is the floor.**
- **Anyone who has mastered their craft can contribute.**
- **The skill outlasts the AI.**

→ [Read the full philosophy](./PHILOSOPHY.md)

---

## What a skill looks like

```
User: /propose-conventional-commit

Claude: Inspecting staged changes...

  feat(auth): add OAuth2 PKCE flow for mobile clients

  Replaces implicit grant with PKCE to address RFC 6749 §10.16
  implicit flow vulnerabilities in native apps.

  Approve? [y/n]
```

The skill runs `git diff --cached`, extracts scope and intent, drafts a Conventional Commits message, and waits for your approval. No freeform output, no guessing.

Skill source: `skills/engineering/development/skills/propose-conventional-commit/`

---

## Or describe your situation

You don't need to know which skill applies. Describe the problem:

```
User: My PRs always get rejected in code review.

Claude: Situation matches: review-pull-request (engineering/development)
        Not installed — get it with:
        /plugins add github:jeffreytse/grimoire/skills/engineering/development
```

```
User: I have an existing API design — is it following best practices?

Claude: You have an existing solution. Applying review-practice-fit...

        ### design-api-architecture — PARTIAL
        ✓ REST endpoints, stateless auth
        ✗ No versioning strategy, no pagination standard
        → Fix: Add /v1/ prefix and cursor-based pagination before next deploy

        🔴 Critical: No rate limiting → DoS exposure
```

`suggest-practice` auto-classifies any situation, routes to the matching skill,
or tells you exactly what to install if the skill isn't in your library yet.

---

## Install

**All skills (Claude Code):**

```bash
/plugins add github:jeffreytse/grimoire
```

**One domain:**

```bash
/plugins add github:jeffreytse/grimoire/skills/engineering
/plugins add github:jeffreytse/grimoire/skills/writing
```

**One sub-domain:**

```bash
/plugins add github:jeffreytse/grimoire/skills/engineering/development
```

**Via script (Claude Code, Codex, Gemini CLI, Copilot):**

```bash
curl -fsSL https://raw.githubusercontent.com/jeffreytse/grimoire/main/scripts/install.sh | bash
```

**Granular installs:**

```bash
./scripts/install.sh --domain engineering
./scripts/install.sh --domain engineering --subdomain development
./scripts/install.sh --skill engineering/development/propose-conventional-commit
./scripts/install.sh --domain writing --target all
```

---

## Featured skills

| Skill | Domain | Source methodology | Verified |
|-------|--------|--------------------|----------|
| [`propose-conventional-commit`](./skills/engineering/development/skills/propose-conventional-commit/) | engineering/development | Angular/Google Conventional Commits | ✓ |
| [`suggest-practice`](./skills/meta/skills/suggest-practice/) | meta | grimoire meta-skill | ✓ |
| [`plan-solution`](./skills/meta/skills/plan-solution/) | meta | McKinsey MECE decomposition | ✓ |
| [`review-pull-request`](./skills/engineering/development/skills/review-pull-request/) | engineering/development | Google Engineering Practices | ✓ |
| [`write-post-mortem`](./skills/engineering/devops/skills/write-post-mortem/) | engineering/devops | Amazon blameless post-mortem | ✓ |

---

## The Grimoire Skill Standard

grimoire maintains an open standard for AI agent skill quality — freely adoptable by any skill library.

Every skill must pass `review-skill` before merge:

| Criterion | Requirement | Rejection example |
|-----------|-------------|-------------------|
| **Adopted by** | Named organizations or institutions | "Many top companies" |
| **Impact** | Cited study or % number | "Significantly improves quality" |
| **Steps** | Immediately executable | Abstract theory or advice |
| **Scope** | One concept per skill | "Nutrition and training program" |
| **Source** | External institution or standard body | `grimoire STANDARD.md` |

→ [Read the full standard](./STANDARD.md) · [Adopt this standard](./STANDARD.md#adopting-this-standard)

See [CONTRIBUTING.md](./CONTRIBUTING.md) to submit a skill.

---

## Domains

grimoire is a framework + reference skills. The domain structure is ready — contribute to fill your domain.

| Domain | Sub-domains |
| ------ | ----------- |
| [meta](./skills/meta/) | **User-facing:** [suggest-practice](./skills/meta/skills/suggest-practice/) · [plan-solution](./skills/meta/skills/plan-solution/) · [review-practice-fit](./skills/meta/skills/review-practice-fit/) · **Contributors:** [write-skill](./skills/meta/skills/write-skill/) · [review-skill](./skills/meta/skills/review-skill/) · [revise-skill](./skills/meta/skills/revise-skill/) · [audit-domain](./skills/meta/skills/audit-domain/) · [deprecate-skill](./skills/meta/skills/deprecate-skill/) · [design-domain](./skills/meta/skills/design-domain/) |
| [engineering](./skills/engineering/) | [development](./skills/engineering/development/skills/), [frontend](./skills/engineering/frontend/skills/), [architecture](./skills/engineering/architecture/skills/), [testing](./skills/engineering/testing/skills/), [reliability](./skills/engineering/reliability/skills/), [devops](./skills/engineering/devops/skills/), [cloud](./skills/engineering/cloud/skills/), [networking](./skills/engineering/networking/skills/), [security](./skills/engineering/security/skills/), [data](./skills/engineering/data/skills/), [ai](./skills/engineering/ai/skills/), [hardware](./skills/engineering/hardware/skills/), [mobile](./skills/engineering/mobile/skills/), [performance](./skills/engineering/performance/skills/), [project-management](./skills/engineering/project-management/skills/), [product](./skills/engineering/product/skills/), [documentation](./skills/engineering/documentation/skills/) |
| [writing](./skills/writing/) | [creative](./skills/writing/creative/skills/), [technical](./skills/writing/technical/skills/), [copywriting](./skills/writing/copywriting/skills/), [academic](./skills/writing/academic/skills/), [journalism](./skills/writing/journalism/skills/) |
| [design](./skills/design/) | [ui-ux](./skills/design/ui-ux/skills/), [graphic](./skills/design/graphic/skills/), [branding](./skills/design/branding/skills/), [motion](./skills/design/motion/skills/), [product](./skills/design/product/skills/) |
| [business](./skills/business/) | [strategy](./skills/business/strategy/skills/), [operations](./skills/business/operations/skills/), [leadership](./skills/business/leadership/skills/), [entrepreneurship](./skills/business/entrepreneurship/skills/), [hr](./skills/business/hr/skills/) |
| [science](./skills/science/) | [biology](./skills/science/biology/skills/), [physics](./skills/science/physics/skills/), [chemistry](./skills/science/chemistry/skills/), [mathematics](./skills/science/mathematics/skills/), [earth-science](./skills/science/earth-science/skills/), [astronomy](./skills/science/astronomy/skills/) |
| [marketing](./skills/marketing/) | [seo](./skills/marketing/seo/skills/), [content](./skills/marketing/content/skills/), [social-media](./skills/marketing/social-media/skills/), [paid-ads](./skills/marketing/paid-ads/skills/), [growth](./skills/marketing/growth/skills/), [analytics](./skills/marketing/analytics/skills/) |
| [health](./skills/health/) | [fitness](./skills/health/fitness/skills/), [nutrition](./skills/health/nutrition/skills/), [mental-health](./skills/health/mental-health/skills/), [sleep](./skills/health/sleep/skills/), [medicine](./skills/health/medicine/skills/) |
| [finance](./skills/finance/) | [personal-finance](./skills/finance/personal-finance/skills/), [investing](./skills/finance/investing/skills/), [accounting](./skills/finance/accounting/skills/), [real-estate](./skills/finance/real-estate/skills/), [corporate](./skills/finance/corporate/skills/) |
| [education](./skills/education/) | [curriculum](./skills/education/curriculum/skills/), [teaching](./skills/education/teaching/skills/), [e-learning](./skills/education/e-learning/skills/), [assessment](./skills/education/assessment/skills/), [learning-science](./skills/education/learning-science/skills/) |
| [film](./skills/film/) | [cinematography](./skills/film/cinematography/skills/), [directing](./skills/film/directing/skills/), [editing](./skills/film/editing/skills/), [screenwriting](./skills/film/screenwriting/skills/), [production](./skills/film/production/skills/) |
| [law](./skills/law/) | [contracts](./skills/law/contracts/skills/), [ip](./skills/law/ip/skills/), [employment](./skills/law/employment/skills/), [privacy](./skills/law/privacy/skills/), [corporate](./skills/law/corporate/skills/) |
| [photography](./skills/photography/) | [composition](./skills/photography/composition/skills/), [lighting](./skills/photography/lighting/skills/), [editing](./skills/photography/editing/skills/), [genres](./skills/photography/genres/skills/) |
| [music](./skills/music/) | [composition](./skills/music/composition/skills/), [production](./skills/music/production/skills/), [mixing](./skills/music/mixing/skills/), [theory](./skills/music/theory/skills/), [performance](./skills/music/performance/skills/) |
| [cooking](./skills/cooking/) | [techniques](./skills/cooking/techniques/skills/), [baking](./skills/cooking/baking/skills/), [flavor](./skills/cooking/flavor/skills/), [nutrition](./skills/cooking/nutrition/skills/), [world-cuisine](./skills/cooking/world-cuisine/skills/) |
| [language](./skills/language/) | [learning](./skills/language/learning/skills/), [linguistics](./skills/language/linguistics/skills/), [translation](./skills/language/translation/skills/), [communication](./skills/language/communication/skills/) |
| [art](./skills/art/) | [drawing](./skills/art/drawing/skills/), [painting](./skills/art/painting/skills/), [digital-art](./skills/art/digital-art/skills/), [illustration](./skills/art/illustration/skills/), [color-theory](./skills/art/color-theory/skills/) |
| [sports](./skills/sports/) | [training](./skills/sports/training/skills/), [coaching](./skills/sports/coaching/skills/), [nutrition](./skills/sports/nutrition/skills/), [tactics](./skills/sports/tactics/skills/), [recovery](./skills/sports/recovery/skills/) |
| [productivity](./skills/productivity/) | [time-management](./skills/productivity/time-management/skills/), [habits](./skills/productivity/habits/skills/), [focus](./skills/productivity/focus/skills/), [goals](./skills/productivity/goals/skills/), [tools](./skills/productivity/tools/skills/) |
| [travel](./skills/travel/) | [planning](./skills/travel/planning/skills/), [budgeting](./skills/travel/budgeting/skills/), [cultural](./skills/travel/cultural/skills/), [adventure](./skills/travel/adventure/skills/) |
| [psychology](./skills/psychology/) | [cognitive](./skills/psychology/cognitive/skills/), [behavioral](./skills/psychology/behavioral/skills/), [social](./skills/psychology/social/skills/), [clinical](./skills/psychology/clinical/skills/), [positive](./skills/psychology/positive/skills/) |
| [home](./skills/home/) | [renovation](./skills/home/renovation/skills/), [interior-design](./skills/home/interior-design/skills/), [gardening](./skills/home/gardening/skills/), [organization](./skills/home/organization/skills/), [smart-home](./skills/home/smart-home/skills/) |
| [environment](./skills/environment/) | [sustainability](./skills/environment/sustainability/skills/), [ecology](./skills/environment/ecology/skills/), [climate](./skills/environment/climate/skills/), [energy](./skills/environment/energy/skills/), [policy](./skills/environment/policy/skills/) |
| [pets](./skills/pets/) | [dogs](./skills/pets/dogs/skills/), [cats](./skills/pets/cats/skills/), [training](./skills/pets/training/skills/), [nutrition](./skills/pets/nutrition/skills/), [health](./skills/pets/health/skills/) |
| [fashion](./skills/fashion/) | [styling](./skills/fashion/styling/skills/), [wardrobe](./skills/fashion/wardrobe/skills/), [design](./skills/fashion/design/skills/), [sustainability](./skills/fashion/sustainability/skills/), [accessories](./skills/fashion/accessories/skills/) |
| [parenting](./skills/parenting/) | [infant](./skills/parenting/infant/skills/), [toddler](./skills/parenting/toddler/skills/), [school-age](./skills/parenting/school-age/skills/), [teen](./skills/parenting/teen/skills/) |
| [automotive](./skills/automotive/) | [maintenance](./skills/automotive/maintenance/skills/), [troubleshooting](./skills/automotive/troubleshooting/skills/), [buying](./skills/automotive/buying/skills/), [modifications](./skills/automotive/modifications/skills/), [ev](./skills/automotive/ev/skills/) |

---

## Agent support

| Agent | Plugin install | Script install |
| ----- | -------------- | -------------- |
| Claude Code | `/plugins add github:jeffreytse/grimoire` | `--target claude` |
| Codex | `/plugins add github:jeffreytse/grimoire` | `--target codex` |
| Cursor | `/plugins add github:jeffreytse/grimoire` | `--target all` |
| OpenCode | See [`.opencode/INSTALL.md`](./.opencode/INSTALL.md) | `--target all` |
| Gemini CLI | — | `--target gemini` |

---

## Contributing

**grimoire has 32 skills. It needs 500. Pick a domain.**

Every domain has empty sub-domains waiting for skills. If you know a field — engineering, law, finance, music, cooking, anything — add the practices you've seen work at the highest level.

Skills must pass [`review-skill`](./skills/meta/skills/review-skill/) before merge.
The meta skills guide the full contribution workflow:

| Task | Skill |
|------|-------|
| Write a new skill | [`write-skill`](./skills/meta/skills/write-skill/) |
| Review a skill PR | [`review-skill`](./skills/meta/skills/review-skill/) |
| Fix review findings | [`revise-skill`](./skills/meta/skills/revise-skill/) |
| Add a new domain | [`design-domain`](./skills/meta/skills/design-domain/) |
| Audit a domain's health | [`audit-domain`](./skills/meta/skills/audit-domain/) |
| Retire an outdated skill | [`deprecate-skill`](./skills/meta/skills/deprecate-skill/) |

See [CONTRIBUTING.md](./CONTRIBUTING.md) for the full standard and [GOVERNANCE.md](./GOVERNANCE.md) for how the project and standard evolve.

## License

MIT
