<div align="center">
  <a href="https://github.com/jeffreytse/grimoire">
    <img alt="grimoire" src="./docs/banner.svg" width="700">
  </a>

  <p>📖 The world's professional best practices — your expert consultant, anywhere, anytime.</p>

  <br><h1>📖 Grimoire 📖</h1>

</div>

<h4 align="center">
  Multi-domain skill collection for <a href="#-agent-support">AI assistants</a>.
</h4>

<p align="center">
  <a href="https://github.com/jeffreytse/grimoire/actions/workflows/validate.yml">
    <img src="https://github.com/jeffreytse/grimoire/actions/workflows/validate.yml/badge.svg"
      alt="Skill Validation" />
  </a>

  <a href="https://github.com/sponsors/jeffreytse">
    <img src="https://img.shields.io/static/v1?label=sponsor&message=%E2%9D%A4&logo=GitHub&link=&color=greygreen"
      alt="Donate (GitHub Sponsor)" />
  </a>

  <a href="https://github.com/jeffreytse/zsh-vi-mode/releases">
    <img src="https://img.shields.io/github/v/release/jeffreytse/grimoire?color=brightgreen"
      alt="Release Version" />
  </a>

  <a href="https://github.com/jeffreytse/grimoire/graphs/contributors">
    <img src="https://img.shields.io/github/contributors/jeffreytse/grimoire"
      alt="Contributors" />
  </a>

  <a href="./LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg"
      alt="License: MIT" />
  </a>

  <a href="https://liberapay.com/jeffreytse">
    <img src="http://img.shields.io/liberapay/goal/jeffreytse.svg?logo=liberapay"
      alt="Donate (Liberapay)" />
  </a>

  <a href="https://patreon.com/jeffreytse">
    <img src="https://img.shields.io/badge/support-patreon-F96854.svg?style=flat-square"
      alt="Donate (Patreon)" />
  </a>

  <a href="https://ko-fi.com/jeffreytse">
    <img height="20" src="https://www.ko-fi.com/img/githubbutton_sm.svg"
      alt="Donate (Ko-fi)" />
  </a>

  <a href="#-agent-support">
    <img src="https://img.shields.io/badge/works%20with-Claude%20%C2%B7%20Codex%20%C2%B7%20Cursor%20%C2%B7%20Gemini%20%C2%B7%20OpenCode-blue"
      alt="Works with" />
  </a>

  <a href="https://github.com/jeffreytse/grimoire/stargazers">
    <img src="https://img.shields.io/github/stars/jeffreytse/grimoire?style=social"
      alt="GitHub Stars" />
  </a>
</p>

<div align="center">
  <h4>
    <a href="#-why-grimoire">Why</a> |
    <a href="#-what-a-skill-looks-like">Features</a> |
    <a href="#%EF%B8%8F-install">Install</a> |
    <a href="#%EF%B8%8F-domains">Domains</a> |
    <a href="#-contributing">Contributing</a> |
    <a href="#-license">License</a>
  </h4>
</div>

<div align="center">
  <sub>Built with ❤︎ by
  <a href="https://jeffreytse.net">jeffreytse</a> and
  <a href="https://github.com/jeffreytse/grimoire/graphs/contributors">contributors</a>
  </sub>
</div>
<br>

## 🤔 Why Grimoire?

> The world's knowledge is in your AI. The world's practice is not.

AI assistants have ingested every textbook, every paper, every article ever written. They
understand fields. They do not practice them. Practice is what happens after 10,000 hours.
Practice is what a senior surgeon does without thinking. Practice is what a staff engineer
knows not to do. Practice is what grimoire encodes.

- 🌍 **The world's best practices belong to everyone.** A McKinsey engagement costs $1M. A senior lawyer bills $800/hr. A structural engineer isn't available at 2am. The practices they follow — proven at the highest levels — are not proprietary. They belong to the world. Grimoire makes them free.

- ⚡ **Skills are verbs.** Not descriptions of what experts know — the exact steps they take, in the exact situation they face, proven at scale. If you can't act on it in the next five minutes, it isn't a skill.

- 🔬 **Every claim must be proven.** One skill. One concept. Adopted by most top-tier institutions in the field, with measurable impact and a named source. If you can't prove it, you can't ship it.

- 🏔️ **Consensus is the floor.** If the world's best professionals are split, grimoire acknowledges the debate — and encodes the majority position. When there is no consensus, there is no best practice to ship.

- 🤝 **Anyone who has mastered their craft can contribute.** A nurse. A jazz musician. A securities lawyer. A structural engineer. Grimoire is not a developer project. It is a project for everyone who has spent 10,000 hours in a field and has something to say about how it's really done.

- ♾️ **The skill outlasts the AI.** Plain Markdown. No lock-in. No proprietary format. These skills will outlive every AI assistant currently running.

→ [Read the full philosophy](./PHILOSOPHY.md)

---

## ✨ What a Skill Looks Like

**Engineering** — propose a commit message from staged diff:

```
User: /propose-conventional-commit

Claude: Inspecting staged changes...

  feat(auth): add OAuth2 PKCE flow for mobile clients

  Replaces implicit grant with PKCE to address RFC 6749 §10.16
  implicit flow vulnerabilities in native apps.

  Approve? [y/n]
```

**Law** — audit a SaaS contract before signing:

```
User: /review-saas-contract

Claude: Red flags found.

  🔴 CRITICAL  Liability cap is $500 — your contract value is $240,000
  🔴 CRITICAL  No data processing agreement — GDPR Art.28 requires one
               before processing any EU customer data
  🟡 WARNING   Auto-renewal clause: 90-day cancellation window in §14.3
  🟡 WARNING   IP ownership of custom integrations is unassigned
  ✓  PASS      SLA 99.9% with defined service credits

  Do not sign. Renegotiate: liability cap, DPA, IP clause.
  Source: ABA model SaaS agreements
```

**Finance** — calculate how much you need to retire:

```
User: /calculate-fire-number

Claude: Applying Trinity Study model (Bengen 1994, 95% 30-year success rate)

  Annual expenses:       $60,000
  Safe withdrawal rate:  4.0%
  ─────────────────────────────
  FIRE number:           $1,500,000

  At $4,000/month invested + 7% real return → retire in 20.4 years
  Cut annual spend by $12,000               → retire in 16.8 years

  Next: /design-portfolio-allocation
```

**Health** — design a training program:

```
User: /design-training-program

Claude: Hypertrophy block — intermediate, 4 days/week

  UPPER (Mon/Thu)   Bench 4×8 · Row 4×8 · OHP 3×10 · Face pull 3×15
  LOWER (Tue/Fri)   Squat 4×6 · RDL 3×8 · Leg press 3×12

  Progressive overload: +2.5 kg when all reps completed clean
  Week 5: deload — 40% volume, same intensity

  Source: NSCA CSCS + Israetel volume landmarks (2019)
```

Every skill is structured: triggering condition, verified source, numbered steps, no freeform guessing.

---

## 🎯 Or Describe Your Situation

You don't need to know which skill applies. Just describe the problem:

```
User: I signed a vendor contract and now they're saying they own my data.

Claude: Situation matches: review-saas-contract (law/contracts)
        Applying review-practice-fit...

        🔴 CRITICAL Data ownership clause absent — defaults to vendor jurisdiction
        → Escalate to legal counsel immediately. Request amendment to §8.2.
```

```
User: My team keeps missing sprint goals and I don't know why.

Claude: Situation matches: plan-retrospective (engineering/project-management)

        Run a blameless retrospective (Amazon format):
        1. Timeline reconstruction — what actually happened
        2. Five whys — surface root causes, not symptoms
        3. Action items with owners and deadlines
        → /plan-retrospective to run the session now
```

`suggest-practice` auto-classifies any situation, routes to the matching skill,
or tells you exactly what to install if the skill isn't in your library yet.

---

## ⚒️ Install

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

## 🌟 Featured Skills

| Skill | Domain | Source methodology | Verified |
|-------|--------|--------------------|----------|
| [`propose-conventional-commit`](./skills/engineering/development/skills/propose-conventional-commit/) | engineering/development | Angular/Google Conventional Commits | ✓ |
| [`suggest-practice`](./skills/meta/skills/suggest-practice/) | meta | grimoire meta-skill | ✓ |
| [`design-slo`](./skills/engineering/reliability/skills/design-slo/) | engineering/reliability | Google SRE Book | ✓ |
| [`review-pull-request`](./skills/engineering/development/skills/review-pull-request/) | engineering/development | Google Engineering Practices | ✓ |
| [`write-post-mortem`](./skills/engineering/devops/skills/write-post-mortem/) | engineering/devops | Amazon blameless post-mortem | ✓ |
| [`audit-gdpr-compliance`](./skills/law/privacy/skills/audit-gdpr-compliance/) | law/privacy | GDPR / EDPB guidelines | ✓ |
| [`calculate-dcf`](./skills/finance/investing/skills/calculate-dcf/) | finance/investing | CFA Institute / Damodaran | ✓ |
| [`design-training-program`](./skills/health/fitness/skills/design-training-program/) | health/fitness | NSCA CSCS curriculum | ✓ |
| [`apply-first-principles`](./skills/productivity/focus/skills/apply-first-principles/) | productivity/focus | Aristotle / Descartes / SpaceX | ✓ |
| [`design-go-to-market`](./skills/business/strategy/skills/design-go-to-market/) | business/strategy | Moore "Crossing the Chasm" | ✓ |

---

## 📐 The Grimoire Skill Standard

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

## 🗺️ Domains

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

## 🤖 Agent Support

| Agent | Plugin install | Script install |
| ----- | -------------- | -------------- |
| Claude Code | `/plugins add github:jeffreytse/grimoire` | `--target claude` |
| Codex | `/plugins add github:jeffreytse/grimoire` | `--target codex` |
| Cursor | `/plugins add github:jeffreytse/grimoire` | `--target all` |
| OpenCode | See [`.opencode/INSTALL.md`](./.opencode/INSTALL.md) | `--target all` |
| Gemini CLI | — | `--target gemini` |

---

## 🤝 Contributing

**grimoire has 152 skills. It needs 500. Pick a domain.**

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

## ❤️ Support

grimoire is free. It replaces $500/hr lawyers, $300 doctor visits, and $1M McKinsey
engagements — at zero cost, forever.

If it saved you time, money, or a bad decision:

- **[⭐ Star this repo](https://github.com/jeffreytse/grimoire)** — takes 2 seconds, helps thousands of people find it
- **[💖 Sponsor on GitHub](https://github.com/sponsors/jeffreytse)** — keeps the maintainer funded to add more skills across more domains
- **[☕ Ko-fi](https://ko-fi.com/jeffreytse)** · **[Patreon](https://patreon.com/jeffreytse)** · **[Liberapay](https://liberapay.com/jeffreytse)** — one-time or recurring

Every star makes grimoire more visible. Every sponsorship funds one more domain.

[![Star History Chart](https://api.star-history.com/svg?repos=jeffreytse/grimoire&type=Date)](https://star-history.com/#jeffreytse/grimoire&Date)

---

## 📄 License

This project is licensed under the [MIT license](https://opensource.org/licenses/mit-license.php) © Jeffrey Tse.
