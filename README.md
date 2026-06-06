<div align="center">
  <a href="https://github.com/jeffreytse/grimoire">
    <img alt="grimoire" src="./assets/banner.svg" width="700">
  </a>

  <p>The world's knowledge is in your AI. The world's practice is not.<br>Most people don't know which best practice applies — and AI won't enforce one unless guided. Grimoire closes both gaps.</p>

  <br><h1>📖 Grimoire 📖</h1>

</div>

<h4 align="center">
  Promotes the expert standard you didn't know applied. Guides your <a href="#-agent-support">AI</a> through applying it — automatically, across every field.
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

  <a href="https://github.com/jeffreytse/grimoire/releases">
    <img src="https://img.shields.io/github/v/release/jeffreytse/grimoire?color=brightgreen"
      alt="Release Version" />
  </a>

  <a href="https://github.com/jeffreytse/grimoire/graphs/contributors">
    <img src="https://img.shields.io/github/contributors/jeffreytse/grimoire?color=brightgreen"
      alt="Contributors" />
  </a>

  <a href="./LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-greygreen.svg"
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
    <img height="20" src="https://ko-fi.com/img/githubbutton_sm.svg"
      alt="Donate (Ko-fi)" />
  </a>

  <a href="#-agent-support">
    <img src="https://img.shields.io/badge/works%20with-Claude%20%C2%B7%20Codex%20%C2%B7%20Cursor%20%C2%B7%20Gemini%20%C2%B7%20OpenCode-blue"
      alt="Works with" />
  </a>

  <a href="./skills/">
    <img src="https://img.shields.io/badge/skills-500%2B-blue"
      alt="505 Skills" />
  </a>
</p>

<div align="center">
  <h4>
    <a href="#-why-grimoire">Why</a> |
    <a href="#-what-a-skill-looks-like">Features</a> |
    <a href="#%EF%B8%8F-install">Install</a> |
    <a href="#-quick-start">Quick Start</a> |
    <a href="#%EF%B8%8F-domains">Domains</a> |
    <a href="#-contributing">Contributing</a> |
    <a href="https://github.com/jeffreytse/grimoire/releases">Changelog</a> |
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

## 🎬 Demo

> "I'm 42, AI just took my job, I have a mortgage. What do I do?"

![grimoire demo — natural language problem solved with grimoire best practices](./assets/demo.gif)

---

## 🤔 Why Grimoire?

AI assistants have ingested every textbook, every paper, every article ever written. They
understand fields. They do not practice them. Practice is what happens after 10,000 hours.
Practice is what a senior surgeon does without thinking. Practice is what a staff engineer
knows not to do. Practice is what grimoire encodes.

- 🔍 **Most people don't know the right practice exists.** When you face a problem, you search for a solution — not for the standard that governs it. The ISO certification process, the ABA clause audit, the NSCA periodization model — these exist. Most people solving those problems have never heard of them. Grimoire closes the discovery gap.

- 🤖 **LLMs know the practice. They won't apply it.** Ask an AI to help with a contract and it gives general advice. Ask it to review an architecture and it summarizes what it sees. The model knows ISO, ABA, Google SRE — but without explicit guidance, it won't enforce any of them. Grimoire provides that guidance.

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
        Applying review-best-practice-fit...

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

`suggest-best-practice` auto-classifies any situation, routes to the matching skill,
or tells you exactly what to install if the skill isn't in your library yet.

---

## ⚒️ Install

**One command. Every AI agent on your system.**

```bash
curl -fsSL https://raw.githubusercontent.com/jeffreytse/grimoire/main/scripts/install.sh | bash
```

Auto-detects Claude Code, Codex, and Gemini CLI. Installs to every agent found. No flags needed.

---

**Windows (PowerShell):**

```powershell
Invoke-WebRequest https://raw.githubusercontent.com/jeffreytse/grimoire/main/scripts/install.ps1 -OutFile install.ps1; .\install.ps1
```

---

**Native plugin shortcuts (Claude Code / Codex):**

```bash
/plugins add github:jeffreytse/grimoire                         # all skills (latest)
/plugins add github:jeffreytse/grimoire@v1.0.0                  # pin to a release
/plugins add github:jeffreytse/grimoire/skills/engineering      # one domain
/plugins add github:jeffreytse/grimoire/skills/engineering/development  # one sub-domain
```

---

**Granular script installs:**

```bash
./scripts/install.sh --domain engineering
./scripts/install.sh --domain engineering --subdomain development
./scripts/install.sh --skill engineering/development/propose-conventional-commit
./scripts/install.sh --target all    # force install to all agents, even if not detected
```

**Gemini CLI:**

```bash
gemini extensions install https://github.com/jeffreytse/grimoire          # latest
gemini extensions install https://github.com/jeffreytse/grimoire@v1.0.0   # pin to a release
gemini extensions update grimoire                                         # update later
```

---

**Cursor** — in Agent chat:

```
/add-plugin grimoire
```

Or search "grimoire" in the plugin marketplace.

---

**OpenCode:** add to `opencode.json`:
```json
{ "plugins": ["grimoire@git+https://github.com/jeffreytse/grimoire.git"] }
```

---

## 🚀 Quick Start

**After install, describe any problem in plain language:**

```
User: I need to raise a Series A but don't know how to pitch investors.

Claude: Situation matches: write-value-proposition + design-go-to-market + apply-pyramid-principle
        Applying suggest-best-practice...
        → Start with your value prop. /write-value-proposition
```

Or invoke a skill directly:

```bash
/suggest-best-practice     # describe any problem — auto-routes to the right skill
/review-pull-request       # engineering code review
/calculate-fire-number     # how much do I need to retire?
/review-saas-contract      # flag dangerous clauses before signing
/design-training-program   # build a training program
```

**New to grimoire?** Start with `/suggest-best-practice`. Describe any professional or life situation — it reads your context and routes you to the matching skill, or tells you exactly what to install if the skill isn't in your library yet.

---

## 🔄 Workflows

### User workflows

**1. Know the skill — invoke directly**

```
/propose-conventional-commit
/calculate-fire-number
/review-saas-contract
```

**2. Have a problem — use `suggest-best-practice`**

Describe any situation. `suggest-best-practice` reads it, finds the best-matching skill(s), and routes you there.

```
User: My senior engineer just quit and the team is losing confidence.

Claude: Situation matches: design-onboarding-program + write-leadership-principles
        → /design-onboarding-program to rebuild team structure
        → /write-leadership-principles to establish decision-making clarity
```

**3. Have a solution — use `review-best-practice-fit`**

Already have a plan or approach? `review-best-practice-fit` checks it against best practices and surfaces gaps.

```
User: I'm planning to run weekly 1:1s with each engineer, 30 min, no agenda.
      Is this a good approach?

Claude: Applying review-best-practice-fit (engineering/project-management)...

        🟡 CAUTION No agenda → likely drifts to status updates
        → Best practice: structured 1:1 (Lencioni format) with 3 fixed topics
        🟢 PASS Frequency and duration align with Google re:Work research
```

**4. Problem spans multiple domains — use `plan-best-practice-solution`**

When a problem involves 3+ independent fields (legal, financial, technical, marketing), `plan-best-practice-solution` decomposes it with MECE methodology and sequences skills by dependency.

```
User: I want to leave my job and launch a SaaS startup.

Claude: Situation spans 4 domains. Solution plan:
        1. design-business-model    — validate the idea before building
        2. calculate-unit-economics — confirm revenue model is viable
        3. review-saas-contract     — legal structure before hiring
        4. design-go-to-market      — growth strategy before launch

        Apply in order? I'll pause for confirmation after each step.
```

**5. Complex problem, one domain — use `apply-best-practice-tree`**

When a problem is deeply complex but stays within one field, `apply-best-practice-tree` uses skills as the decomposition engine — recursively matching sub-problems to best practices until each is covered.

```
User: Our production incident response is a mess — we have no process.

Claude: Top match: design-incident-response-process (confidence 0.81). Applying...

        Sub-problems surfaced:
          A. Detection    → apply-alerting-strategy
          B. Escalation   → design-on-call-rotation
          C. Communication → write-status-page-update
          D. Post-incident → apply-five-whys

        Ready to apply A: apply-alerting-strategy. Continue?
```

> `plan-best-practice-solution` handles **breadth** (multiple domains, sequenced). `apply-best-practice-tree` handles **depth** (one domain, recursive).

---

### Decision guide

| Your situation | Start here |
|----------------|------------|
| Know exactly which skill you need | `/skill-name` directly |
| Have a problem, unsure which skill | `/suggest-best-practice` |
| Already have a plan, want gaps checked | `/review-best-practice-fit` |
| Problem spans 3+ independent domains | `/plan-best-practice-solution` |
| Complex problem within one domain | `/apply-best-practice-tree` |
| Don't know what practices exist for a topic | `/discover-best-practices` |
| About to start a task — want to catch gaps before you begin | `/intercept-best-practice` |
| Problem isn't clear yet — need to define it before solving | `/analyze-problem` |
| Two practices conflict (e.g., Google vs IBM) | `/pin-best-practice-preference` |

**6. Don't know what practices exist — use `discover-best-practices`**

Browse what's available for a domain before you have a specific problem.

```
User: What best practices exist for software architecture?

Claude: Best practices available for software architecture:

  Design
    design-api-architecture — most APIs are built without versioning until a
      breaking change forces it; this closes that gap
    design-observability — teams discover they can't debug incidents after the
      first production outage

  Security
    review-security-posture — security reviews done after a breach cost 10× more than before

Most commonly discovered too late:
  ★ review-security-posture — applied after a breach, not before
  ★ design-observability — discovered when the first incident is uninvestigable

Want to apply any of these now, or learn more about one?
```

**7. About to start a task — use `intercept-best-practice`**

Proactively catches applicable practices before you begin, not after problems emerge.

```
User: I'm going to refactor our authentication module.

Claude: Before you start: design-api-security applies to authentication refactors.

  Common gaps caught at this stage:
  • Token expiry and refresh strategy (most teams omit this until a session bug)
  • Rate limiting on auth endpoints
  • Audit logging for authentication events

  Apply design-api-security now, or continue without? [y/n]
```

**8. Problem needs defining first — use `analyze-problem`**

When the problem is vague, define it before routing to a solution skill.

```
User: My team isn't performing well.

Claude: Before finding the right practice, I need to understand the problem.
  What outcome are you trying to achieve — better output, faster delivery, or fewer conflicts?

User: Better output, but also morale is low.

Claude: Problem statement: Team output is below expectations and morale is low.

  Scope: In — team performance and culture. Out — individual skill gaps (separate problem).

  Possible routes:
  A. Fix root cause — run a structured retrospective to surface what's blocking the team
  B. Workaround — redistribute workload while root cause is investigated
  C. Prevent recurrence — establish team health rituals to catch issues earlier

  Which route fits best? I can find the applicable best practices.
```

**Prefer a specific org's practices — use `pin-best-practice-preference`**

When top-tier orgs disagree (Google vs IBM, startup vs enterprise), pin which approach you prefer once — grimoire applies it automatically from then on.

```
User: I follow Google's engineering practices when there's a conflict.

Claude: Pinning preference: Google Engineering Practices → engineering domain.
        Save to: session / project / global?

User: global

Claude: Saved to ~/.grimoire/preferences.md.
        Future engineering conflicts resolve toward Google's approach automatically.
```

**Override priority** — grimoire checks preferences in this order, first match wins:

| Priority | Level | Stored in | Scope |
|----------|-------|-----------|-------|
| 1st | Session | In-memory | Current session only — resets when session ends |
| 2nd | Project | `.grimoire/preferences.md` in project root | Current project only |
| 3rd | Global | `~/.grimoire/preferences.md` | All projects on this machine |
| 4th | Legacy | `CLAUDE.md` under `## Grimoire Preferences` | Fallback |

**Configure manually** — edit the preference files directly without going through the AI:

```markdown
# ~/.grimoire/preferences.md  (global — applies everywhere)
engineering: Google Engineering Practices
finance: CFA Institute standards
law: ABA Model Rules
```

```markdown
# <project-root>/.grimoire/preferences.md  (project — overrides global for this repo)
engineering: startup  # this project moves fast; override the global Google preference
```

Project preferences override global. Session pins override both. Teams can share a global standard while individual projects deviate where needed.

---

### Contributor workflows

**9. Adding a skill**

```
/write-best-practice-skill    # author the skill
/review-best-practice-skill   # validate against STANDARD.md (5 criteria)
/revise-best-practice-skill   # fix any review findings
→ open PR
```

**10. Maintaining a domain**

```
/audit-best-practice-domain    # batch health check — surfaces outdated or weak-sourced skills
/revise-best-practice-skill    # update stale or under-sourced skills
/deprecate-best-practice-skill # retire skills superseded by newer practices
```

---

## 🌟 Featured Skills

| Skill | Domain | Source methodology | Verified |
|-------|--------|--------------------|----------|
| [`apply-five-whys`](./skills/engineering/reliability/skills/apply-five-whys/) | engineering/reliability | Toyota Production System / Google SRE | ✓ |
| [`design-go-to-market`](./skills/business/strategy/skills/design-go-to-market/) | business/strategy | Moore "Crossing the Chasm" | ✓ |
| [`audit-gdpr-compliance`](./skills/law/privacy/skills/audit-gdpr-compliance/) | law/privacy | GDPR / EDPB guidelines | ✓ |
| [`calculate-fire-number`](./skills/finance/personal-finance/skills/calculate-fire-number/) | finance/personal-finance | Bengen (1994) / Trinity Study | ✓ |
| [`design-training-program`](./skills/health/fitness/skills/design-training-program/) | health/fitness | NSCA CSCS curriculum | ✓ |
| [`apply-mise-en-place`](./skills/cooking/techniques/skills/apply-mise-en-place/) | cooking/techniques | Culinary Institute of America | ✓ |
| [`apply-acceptance-commitment-therapy`](./skills/psychology/cognitive/skills/apply-acceptance-commitment-therapy/) | psychology/cognitive | Hayes / ACBS meta-analyses | ✓ |
| [`apply-spaced-repetition`](./skills/education/curriculum/skills/apply-spaced-repetition/) | education/curriculum | Ebbinghaus / Roediger & Karpicke | ✓ |
| [`write-value-proposition`](./skills/writing/copywriting/skills/write-value-proposition/) | writing/copywriting | Osterwalder "Value Proposition Design" | ✓ |
| [`design-training-periodization-plan`](./skills/sports/training/skills/design-training-periodization-plan/) | sports/training | Bompa "Periodization" / NSCA | ✓ |

→ [Browse all 500+ skills by domain](./SKILLS.md)

---

## 📐 The Grimoire Skill Standard

grimoire maintains an open standard for AI agent skill quality — freely adoptable by any skill library.

Every skill must pass `review-best-practice-skill` before merge:

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
| [grimoire](./skills/grimoire/) | **Problem analysis:** [analyze-problem](./skills/grimoire/skills/analyze-problem/) · [discover-best-practices](./skills/grimoire/skills/discover-best-practices/) · **Routing:** [suggest-best-practice](./skills/grimoire/skills/suggest-best-practice/) · [intercept-best-practice](./skills/grimoire/skills/intercept-best-practice/) · **Solution planning:** [plan-best-practice-solution](./skills/grimoire/skills/plan-best-practice-solution/) · [apply-best-practice-tree](./skills/grimoire/skills/apply-best-practice-tree/) · **Practice evaluation:** [review-best-practice-fit](./skills/grimoire/skills/review-best-practice-fit/) · [compare-best-practices](./skills/grimoire/skills/compare-best-practices/) · [audit-applied-best-practices](./skills/grimoire/skills/audit-applied-best-practices/) · **Practice understanding:** [explain-best-practice](./skills/grimoire/skills/explain-best-practice/) · [adapt-best-practice](./skills/grimoire/skills/adapt-best-practice/) · [teach-best-practice](./skills/grimoire/skills/teach-best-practice/) · **Preferences:** [pin-best-practice-preference](./skills/grimoire/skills/pin-best-practice-preference/) · **Contributors:** [write-best-practice-skill](./skills/grimoire/skills/write-best-practice-skill/) · [review-best-practice-skill](./skills/grimoire/skills/review-best-practice-skill/) · [revise-best-practice-skill](./skills/grimoire/skills/revise-best-practice-skill/) · [audit-best-practice-domain](./skills/grimoire/skills/audit-best-practice-domain/) · [deprecate-best-practice-skill](./skills/grimoire/skills/deprecate-best-practice-skill/) · [design-best-practice-domain](./skills/grimoire/skills/design-best-practice-domain/) |
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
| Cursor | `/add-plugin grimoire` (in Agent chat) | `--target all` |
| OpenCode | See [`.opencode/INSTALL.md`](./.opencode/INSTALL.md) | `--target all` |
| Gemini CLI | `gemini extensions install https://github.com/jeffreytse/grimoire` | `--target gemini` |

---

## ❓ FAQ

**Isn't this already in the model's training data?**

Yes — and no. Models know *about* best practices. Skills make models *do* them, reliably.

The difference:

| Without a skill | With a skill |
|-----------------|--------------|
| Model improvises a version of the practice | Model follows the exact steps from the source institution |
| Output varies every run | Same process, same structure, every time |
| Practice applied only if you know to ask | Skill triggers automatically when the situation matches |
| Generic advice | Specific: the right gate, the right question, the right output format |

For simple tasks (write a test, fix a bug), the skeptic is right — the model doesn't need a skill. For complex, multi-step workflows — an SLO design, a post-mortem, an incident response — skills measurably change what you get. The model knows Google's engineering review process exists. It does not reliably know which question to ask first, what the output format is, or when to stop. That's what a skill encodes.

**These are just textbook practices the model already knows. Why bother?**

Knowing a practice and reliably executing it are different things. Ask any model "I just had a production incident" — you'll get a generic write-up. Run `write-post-mortem` and you get: blameless framing, 5-whys, timeline, contributing factors, action items with owners, and a detection section. The model *knew* all of that before the skill existed. The skill is what makes it happen consistently, in the right format, every time.

The "textbook" objection gets it backwards. Established practices are *ideal* for skills precisely because they're falsifiable — you can verify whether the output matches what Google's SRE book, Amazon's mechanisms, or the WHO protocol actually prescribes. If you find a skill that adds nothing over a bare prompt, that's a quality failure. [File an issue.](https://github.com/jeffreytse/grimoire/issues)

**Does grimoire conflict with my team's existing conventions?**

Skills describe what the world's top institutions do. Your team may do things differently — and be right to. Two ways to handle it:

**Pin your preference.** Tell grimoire which approach to follow when practices conflict:

```
User: We follow Google's engineering practices, not IBM's.
→ Claude pins this via `pin-best-practice-preference` — applies automatically from now on.
```

**Override or fork.** A skill is a starting point, not a mandate. Adapt any skill to your context, or ignore it entirely. The format is plain Markdown and the license is MIT.

---

## 🤝 Contributing

**grimoire has 500+ skills. It needs 1000. Pick a domain.**

Every domain has empty sub-domains waiting for skills. If you know a field — engineering, law, finance, music, cooking, anything — add the practices you've seen work at the highest level.

**Your first skill in ~30 minutes:**
1. Pick a practice you've used at the highest level in your field
2. Run `/write-best-practice-skill` — it guides you through the format step by step
3. Open a PR — `/review-best-practice-skill` runs automatically and flags any gaps
4. Merge after review passes

Not sure where to start? Browse [open issues](https://github.com/jeffreytse/grimoire/issues) for requested skills, or pick any empty sub-domain from the table below.

Skills must pass [`review-best-practice-skill`](./skills/grimoire/skills/review-best-practice-skill/) before merge.
The meta skills guide the full contribution workflow:

| Task | Skill |
|------|-------|
| Write a new skill | [`write-best-practice-skill`](./skills/grimoire/skills/write-best-practice-skill/) |
| Review a skill PR | [`review-best-practice-skill`](./skills/grimoire/skills/review-best-practice-skill/) |
| Fix review findings | [`revise-best-practice-skill`](./skills/grimoire/skills/revise-best-practice-skill/) |
| Add a new domain | [`design-best-practice-domain`](./skills/grimoire/skills/design-best-practice-domain/) |
| Audit a domain's health | [`audit-best-practice-domain`](./skills/grimoire/skills/audit-best-practice-domain/) |
| Retire an outdated skill | [`deprecate-best-practice-skill`](./skills/grimoire/skills/deprecate-best-practice-skill/) |

See [CONTRIBUTING.md](./CONTRIBUTING.md) for the full standard and [GOVERNANCE.md](./GOVERNANCE.md) for how the project and standard evolve.

## ❤️ Support

grimoire is free. It replaces $500/hr lawyers, $300 doctor visits, and $1M McKinsey
engagements — at zero cost, forever.

If it saved you time, money, or a bad decision:

- **[⭐ Star this repo](https://github.com/jeffreytse/grimoire)** — takes 2 seconds, helps thousands of people find it
- **[💖 Sponsor on GitHub](https://github.com/sponsors/jeffreytse)** — keeps the maintainer funded to add more skills across more domains
- **[☕ Ko-fi](https://ko-fi.com/jeffreytse)** · **[Patreon](https://patreon.com/jeffreytse)** · **[Liberapay](https://liberapay.com/jeffreytse)** — one-time or recurring

Every star makes grimoire more visible. Every sponsorship funds one more domain.

---

## 📄 License

This project is licensed under the [MIT license](https://opensource.org/licenses/mit-license.php) © Jeffrey Tse.
