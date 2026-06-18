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

  <a href="https://github.com/jeffreytse/grimoire-skills">
    <img src="https://img.shields.io/badge/skills-900%2B-blue"
      alt="1000+ Skills" />
  </a>
</p>

<div align="center">
  <h4>
    <a href="#-why-grimoire">Why</a> |
    <a href="#-skills-in-action">Features</a> |
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

> "We're launching our SaaS in 48 hours. I'm terrified something will break. What do we do?"

![grimoire demo — pre-launch protocol: apply-premortem → design-slo → plan-incident-response → run-game-day](./assets/demo.gif)

Not an engineer? See [Skills in Action](#-skills-in-action) for law, finance, cooking, and health examples — same idea, different domain.

## 🤔 Why Grimoire?

> Your AI knows everything — Grimoire makes it practice it.

Books gave everyone knowledge. Google gave everyone access. AI gave everyone comprehension. None of them gave everyone *practice*. Grimoire does.

The model knows SOLID, DDD, Google SRE, and the ABA Model Rules. Without explicit guidance, it enforces none of them. Grimoire is the enforcement layer — **skills**: named, citable, executable units of expert practice, one concept, one source, one set of steps.

- 🔍 **You don't know what you don't know.** Grimoire surfaces the governing standard you didn't know applied.
- 🤖 **Knowing is not enough. Grimoire enforces.** Specific steps, verifiable criteria, repeatable results. Not summaries — verdicts.
- 🌍 **The world's best practices belong to everyone.** McKinsey charges $1M. Senior lawyers bill $800/hr. The practices they follow are not proprietary. Grimoire makes them free — as executable steps, cited and verified.
- 🧱 **Every profession. Every AI.** 1000+ skills across 27 domains. Works with Claude, Codex, Cursor, Gemini, OpenCode, and OpenClaw.

If you've spent 10,000 hours mastering something, your practice belongs here.

→ [Read the full philosophy](https://github.com/jeffreytse/grimoire-skills/blob/main/PHILOSOPHY.md)

## ✨ Skills in Action

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

**Or describe the situation in plain language** — you don't need to know which skill applies:

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

## ⚒️ Install

**macOS / Linux:**

```bash
curl -fsSL https://raw.githubusercontent.com/jeffreytse/grimoire/main/scripts/install.sh | bash
```

**Windows (PowerShell):**

```powershell
Invoke-WebRequest -Uri "https://github.com/jeffreytse/grimoire/releases/latest/download/grimoire-windows-amd64.exe" -OutFile "$env:USERPROFILE\bin\grimoire.exe"
```

Add `%USERPROFILE%\bin` to PATH if not already set. Run as Administrator for symlink support.

**Go:**

```bash
go install github.com/jeffreytse/grimoire@latest
```

**Then install skills (all platforms):**

```bash
grimoire update    # clone skill library → ~/.grimoire
grimoire install   # auto-detect & install to all AI agents
```

Auto-detects Claude Code, Codex, Gemini CLI, OpenClaw, and OpenCode.

**Native plugin shortcuts (Claude Code):**

```bash
# Step 1: add the marketplace
/plugin marketplace add jeffreytse/grimoire

# Step 2: install (skills are namespaced, e.g. /grimoire-engineering:propose-conventional-commit)
/plugin install grimoire@grimoire                   # all domains (latest)
/plugin install grimoire-engineering@grimoire       # one domain

# For subdomain-level installs, use grimoire
```

**Granular installs:**

```bash
grimoire install                                            # interactive TUI
grimoire install --domain engineering
grimoire install --domain engineering --subdomain development
grimoire install --skill engineering/development/propose-conventional-commit
grimoire install --target all      # install to all agents, even if not detected
grimoire update                    # pull latest (choose stable or unstable channel)
grimoire doctor                    # health check: git repo, symlinks, config
grimoire version                   # version info with commit and date
grimoire list                      # list available domains and skills
```

**Gemini CLI:**

```bash
gemini extensions install https://github.com/jeffreytse/grimoire-skills          # latest
gemini extensions install https://github.com/jeffreytse/grimoire-skills@v1.0.0   # pin to a release
gemini extensions update grimoire                                         # update later
```

**Cursor:**

```bash
grimoire install --target cursor
```

**OpenCode:**

```bash
grimoire install --target opencode
```

Or via plugin in `opencode.json`:
```json
{ "plugin": ["grimoire@git+https://github.com/jeffreytse/grimoire-skills.git"] }
```

**OpenClaw:** see [`.openclaw/INSTALL.md`](./.openclaw/INSTALL.md) or run `grimoire install --target openclaw`.

## 🤖 Agent Support

| Agent | Plugin install | Script install |
| ----- | -------------- | -------------- |
| Claude Code | `/plugin marketplace add jeffreytse/grimoire-skills` then `/plugin install grimoire@grimoire-skills` | `grimoire install --target claude` |
| GitHub Copilot CLI | `copilot plugin marketplace add jeffreytse/grimoire-skills` then `copilot plugin install grimoire@grimoire-skills` | `grimoire install --target all` |
| Gemini CLI | `gemini extensions install https://github.com/jeffreytse/grimoire-skills` | `grimoire install --target gemini` |
| OpenCode | See [`.opencode/INSTALL.md`](./.opencode/INSTALL.md) | `grimoire install --target opencode` |
| OpenClaw | See [`.openclaw/INSTALL.md`](./.openclaw/INSTALL.md) | `grimoire install --target openclaw` |
| Codex CLI | `AGENTS.md` auto-loaded; browse `/plugins` in CLI | `grimoire install --target codex` |
| Cursor | `AGENTS.md` context injection | `grimoire install --target cursor` |

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

For CI enforcement, initialize once per project and gate PRs with the `grimoire` CLI:

```bash
grimoire init    # creates .grimoire/settings.toml with auto-detected profile
grimoire check   # reads compliance report, exits 0/1/2
```

**New to grimoire?** Start with `/suggest-best-practice`. Describe any professional or life situation — it reads your context and routes you to the matching skill, or tells you exactly what to install if the skill isn't in your library yet.

## 🎯 Workflows

| Your situation | Start here |
|----------------|------------|
| Know exactly which skill you need | `/skill-name` directly |
| Have a problem, unsure which skill | `/suggest-best-practice` |
| Already have a plan, want gaps checked | `/review-best-practice-fit` |
| Needs 2+ practices coordinated — within one domain or across many — and sub-problems are identifiable upfront | `/plan-best-practice-solution` |
| Complex problem where sub-problems are opaque and emerge through execution | `/apply-best-practice-tree` |
| Don't know what practices exist for a topic | `/discover-best-practices` |
| About to start a task — want to catch gaps before you begin | `/start-best-practice` |
| Problem isn't clear yet — need to define it before solving | `/analyze-best-practice-problem` |
| Activate a paradigm's best practices (OOP, TDD, etc.) | `/apply-best-practice-profile` |
| Align any project or artifact to stated best practice preferences (BPDD) | `/apply-best-practice-driven-development` |
| Check if any artifact aligns with stated best practice preferences | `/check-best-practice-compliance` |
| Have a specific compliance finding to fix | `/fix-best-practice-finding` |
| Two practices exist — want side-by-side comparison | `/compare-best-practices` |
| Two practices conflict — want to reason through which fits | `/resolve-best-practice-conflict` |
| Resolved a conflict — want to save the decision for future sessions | `/pin-best-practice-preference` |

→ [BPDD guide](./docs/bpdd.md) — cycle, compliance linter, LSP output, suppression, CI integration

### Common paths

**Describe any problem → `suggest-best-practice`**

Describe any situation. `suggest-best-practice` reads it, finds the best-matching skill(s), and routes you there.

```
User: My senior engineer just quit and the team is losing confidence.

Claude: Situation matches: design-onboarding-program + write-leadership-principles
        → /design-onboarding-program to rebuild team structure
        → /write-leadership-principles to establish decision-making clarity
```

**Multiple practices needed → `plan-best-practice-solution`**

When a problem needs 2+ practices coordinated in sequence — within one domain or across many — `plan-best-practice-solution` decomposes it with MECE methodology, sequences skills by dependency, and recursively plans complex sub-problems.

```
User: I want to leave my job and launch a SaaS startup.

Claude: Situation spans 4 domains. Solution plan:
        1. design-business-model    — validate the idea before building
        2. calculate-unit-economics — confirm revenue model is viable
        3. review-saas-contract     — legal structure before hiring
        4. design-go-to-market      — growth strategy before launch

        Apply in order? I'll pause for confirmation after each step.
```

### Learning workflows

| Want to... | Start here |
|------------|------------|
| Understand what a skill does and why it exists | `/explain-best-practice` |
| Adapt a skill to a different context or constraint | `/adapt-best-practice` |
| Teach a practice to someone else (structured walkthrough) | `/teach-best-practice` |

### Contributor workflows

**Setup**

```
/install-grimoire              # first-time setup
/configure-grimoire            # customize behavior for your context
```

**Adding a skill**

```
/write-best-practice-skill    # author the skill
/review-best-practice-skill   # validate against STANDARD.md (5 criteria)
/revise-best-practice-skill   # fix any review findings
→ open PR
```

**Adding a profile (named practice set)**

```
/write-best-practice-profile   # author a named profile (e.g., "oop", "tdd")
/review-best-practice-profile  # validate it
/share-best-practice-profile   # publish to others
```

**Maintaining a domain**

```
/audit-best-practice-domain        # batch health check — surfaces outdated or weak-sourced skills
/audit-applied-best-practices      # audit what practices have actually been applied to a project
/design-best-practice-domain       # design a new domain from scratch
/revise-best-practice-skill        # update stale or under-sourced skills
/deprecate-best-practice-skill     # retire skills superseded by newer practices
```

## ⚙️ Settings

When top-tier orgs disagree (Google vs IBM, startup vs enterprise), pin which approach you prefer once — grimoire applies it automatically from then on.

```
User: I follow Google's engineering practices when there's a conflict.

Claude: Pinning preference: Google Engineering Practices → engineering domain.
        Save to: session / project / global?

User: global

Claude: Saved to ~/.config/grimoire/settings.toml.
        Future engineering conflicts resolve toward Google's approach automatically.
```

**Override priority** — grimoire checks preferences in this order, first match wins:

| Priority | Level | Stored in | Scope |
|----------|-------|-----------|-------|
| 1st | Session | In-memory | Current session only — resets when session ends |
| 2nd | Project personal | `.grimoire/settings.local.toml` | Current project, gitignored |
| 3rd | Project shared | `.grimoire/settings.toml` | Current project, committed to repo |
| 4th | Global | `~/.config/grimoire/settings.toml` | All projects on this machine |

**Configure manually** — edit the settings files directly without going through the AI:

```toml
# ~/.config/grimoire/settings.toml  (global — applies everywhere)
profiles = ["oop"]                 # activate all skills tagged "oop"
# profiles = ["clean-architecture", "tdd"]  # multiple — first entry wins conflicts

[engineering]
practices = ["Google Engineering Practices"]

[finance]
practices = ["CFA Institute standards"]
```

```toml
# <project-root>/.grimoire/settings.toml  (project — overrides global for this repo)
[engineering.architecture]
practices = [
  "SOLID principles: production code",
  "KISS: prototypes, scripts"
]
fallback = "ask"
```

Project settings override global. Session pins override both. Teams can share a global standard while individual projects deviate where needed.

**`practices = ["OOP"]` vs `profiles = ["oop"]`** — both signal OOP intent, but differently. `practices = ["OOP"]` in a domain section is a loose hint — the AI leans toward OOP conventions from its training. `profiles = ["oop"]` at the top level activates specific installed skills (exact steps, validated sources). Use `profiles` for precision; `practices` for domain-level style preference. → [Full comparison](./docs/profiles.md#profiles-vs-practices)

**Guided settings management:** Use `/configure-grimoire` to view, edit, or validate settings without touching TOML directly. Use `/apply-best-practice-profile` to activate a full paradigm (OOP, TDD, clean architecture) in one command. Use `/resolve-best-practice-conflict` to resolve contradictions between two installed skills and record the priority automatically. Use `/apply-best-practice-driven-development` to run the full BPDD cycle — or see the [BPDD guide](./docs/bpdd.md) for the linter, LSP output, and CI integration details.

→ [Full settings reference](./docs/settings.md) — all keys, override hierarchy, TOML examples

## 🎭 Profiles

Activate a named set of skills in one line — no list to maintain, no file to create.

```toml
# .grimoire/settings.toml
profiles = ["oop"]   # activates every installed skill tagged "oop"
```

Grimoire resolves the name in this order, first match wins:

1. `.grimoire/profiles/<name>.toml` — project-level file
2. `~/.grimoire/profiles/<name>.toml` — user-level file
3. `.grimoire/profiles/default.toml` — project-level fallback
4. `~/.grimoire/profiles/default.toml` — user-level fallback
5. Tag query — all installed skills where `tags` contains the name

If no file exists, the tag query fires automatically. `profiles = ["oop"]` works without creating any file.

**Multiple profiles** — combine paradigms; first entry wins conflicts, duplicates are deduplicated:

```toml
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

**`profiles` vs `practices`** — `profiles` activates skill bundles globally; `practices` is a domain-scoped explicit list. → [Full comparison](./docs/profiles.md#profiles-vs-practices)

→ [Full profiles guide](./docs/profiles.md) — resolution order, conflict handling, sharing profiles

## 📏 BPDD — Best Practice Driven Development

Grimoire doubles as a best-practice linter. Encode your quality criteria once in `settings.toml`, then run `/check-best-practice-compliance` against any artifact — a codebase, a legal contract, a business plan, a training program. Same criteria every run. Gaps that survive human review get caught by the check.

**The cycle** — same inversion as TDD: declare what "good" looks like first, then bring the artifact into alignment.

```
1. Red      — run compliance check; identify which practices FAIL or are PARTIAL
2. Green    — invoke the relevant grimoire skill; fix until the check passes
3. Refactor — clean up while keeping the check green
4. Commit   — record progress; repeat for next gap
```

Run `/apply-best-practice-driven-development` to drive the full cycle. Run `/check-best-practice-compliance` for a one-off check.

**Output** — dual format, always written to `.grimoire/reports/`:

| File | Format | Use |
|------|--------|-----|
| `compliance-<timestamp>.json` | LSP-compatible JSON | editors, CI pipelines, LSP servers, dashboards |
| `compliance-<timestamp>.html` | HTML | browser or CI artifact upload |

The JSON follows the LSP Diagnostic schema — `uri` + `range` locate any finding in any text artifact, not just code. Gate CI with the `grimoire` CLI (no jq required):

```bash
# Initialize once
grimoire init

# After running /check-best-practice-compliance in your AI session:
grimoire check          # exits 0 (pass) or 1 (fail)
grimoire check --json   # machine-readable output
```

Install: `curl -fsSL https://raw.githubusercontent.com/jeffreytse/grimoire/main/scripts/install.sh | bash`

**Coverage thresholds** — set in `settings.toml`, enforced on every check:

```toml
[engineering]
compliance-threshold = 80        # fail if overall criteria coverage < 80%
compliance-threshold-error = 0   # fail if any error-severity violations remain
```

Use `/fix-best-practice-finding` to fix one specific compliance finding — targeted, location-aware, verified. Use `/apply-best-practice-driven-development` to fix everything systematically.

→ [Full BPDD guide](./docs/bpdd.md) — cycle, linter, LSP schema, false positive suppression, incremental mode

## 🌟 Featured Skills

| Skill | Domain | Source methodology | Verified |
|-------|--------|--------------------|----------|
| [`review-saas-contract`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/law/contracts/skills/review-saas-contract/) | law/contracts | ABA model SaaS agreements | ✓ |
| [`calculate-fire-number`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/finance/personal-finance/skills/calculate-fire-number/) | finance/personal-finance | Bengen (1994) / Trinity Study | ✓ |
| [`negotiate-salary`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/finance/personal-finance/skills/negotiate-salary/) | finance/personal-finance | Fisher & Ury "Getting to Yes" / BLS data | ✓ |
| [`design-sleep-protocol`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/health/sleep/skills/design-sleep-protocol/) | health/sleep | Matthew Walker "Why We Sleep" / AASM | ✓ |
| [`apply-mise-en-place`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/cooking/techniques/skills/apply-mise-en-place/) | cooking/techniques | Culinary Institute of America | ✓ |
| [`apply-five-whys`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/reliability/skills/apply-five-whys/) | engineering/reliability | Toyota Production System / Google SRE | ✓ |
| [`design-training-program`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/health/fitness/skills/design-training-program/) | health/fitness | NSCA CSCS curriculum | ✓ |
| [`apply-acceptance-commitment-therapy`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/psychology/cognitive/skills/apply-acceptance-commitment-therapy/) | psychology/cognitive | Hayes / ACBS meta-analyses | ✓ |
| [`write-value-proposition`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/writing/copywriting/skills/write-value-proposition/) | writing/copywriting | Osterwalder "Value Proposition Design" | ✓ |
| [`design-training-periodization-plan`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/sports/training/skills/design-training-periodization-plan/) | sports/training | Bompa "Periodization" / NSCA | ✓ |

→ [Browse all skills by domain](https://github.com/jeffreytse/grimoire-skills/blob/main/SKILLS.md)

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

→ [Read the full standard](https://github.com/jeffreytse/grimoire-skills/blob/main/STANDARD.md) · [Adopt this standard](https://github.com/jeffreytse/grimoire-skills/blob/main/STANDARD.md#adopting-this-standard)

See [CONTRIBUTING.md](https://github.com/jeffreytse/grimoire-skills/blob/main/CONTRIBUTING.md) to submit a skill.

## 🗺️ Domains

grimoire is a framework + reference skills. The domain structure is ready — contribute to fill your domain.

📦 **Skills library:** [github.com/jeffreytse/grimoire-skills](https://github.com/jeffreytse/grimoire-skills)

| Domain | Sub-domains |
| ------ | ----------- |
| [grimoire](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/) | **Setup:** [install-grimoire](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/install-grimoire/) · [configure-grimoire](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/configure-grimoire/) · **Problem analysis:** [analyze-best-practice-problem](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/analyze-best-practice-problem/) · [discover-best-practices](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/discover-best-practices/) · **Routing:** [suggest-best-practice](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/suggest-best-practice/) · [start-best-practice](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/start-best-practice/) · **Solution planning:** [plan-best-practice-solution](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/plan-best-practice-solution/) · [apply-best-practice-tree](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/apply-best-practice-tree/) · **Practice evaluation:** [review-best-practice-fit](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/review-best-practice-fit/) · [compare-best-practices](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/compare-best-practices/) · [audit-applied-best-practices](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/audit-applied-best-practices/) · **Practice understanding:** [explain-best-practice](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/explain-best-practice/) · [adapt-best-practice](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/adapt-best-practice/) · [teach-best-practice](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/teach-best-practice/) · **Preferences:** [pin-best-practice-preference](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/pin-best-practice-preference/) · [resolve-best-practice-conflict](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/resolve-best-practice-conflict/) · [apply-best-practice-profile](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/apply-best-practice-profile/) · [write-best-practice-profile](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/write-best-practice-profile/) · [review-best-practice-profile](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/review-best-practice-profile/) · [share-best-practice-profile](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/share-best-practice-profile/) · **Compliance:** [apply-best-practice-driven-development](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/apply-best-practice-driven-development/) · [check-best-practice-compliance](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/check-best-practice-compliance/) · **Contributors:** [write-best-practice-skill](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/write-best-practice-skill/) · [review-best-practice-skill](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/review-best-practice-skill/) · [revise-best-practice-skill](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/revise-best-practice-skill/) · [audit-best-practice-domain](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/audit-best-practice-domain/) · [deprecate-best-practice-skill](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/deprecate-best-practice-skill/) · [design-best-practice-domain](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/design-best-practice-domain/) |
| [engineering](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/) | [development](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/development/skills/), [frontend](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/frontend/skills/), [architecture](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/architecture/skills/), [testing](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/testing/skills/), [reliability](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/reliability/skills/), [devops](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/devops/skills/), [cloud](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/cloud/skills/), [networking](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/networking/skills/), [security](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/security/skills/), [data](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/data/skills/), [ai](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/ai/skills/), [hardware](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/hardware/skills/), [mobile](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/mobile/skills/), [performance](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/performance/skills/), [project-management](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/project-management/skills/), [product](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/product/skills/), [documentation](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/engineering/documentation/skills/) |
| [writing](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/writing/) | [creative](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/writing/creative/skills/), [technical](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/writing/technical/skills/), [copywriting](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/writing/copywriting/skills/), [academic](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/writing/academic/skills/), [journalism](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/writing/journalism/skills/) |
| [design](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/design/) | [ui-ux](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/design/ui-ux/skills/), [graphic](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/design/graphic/skills/), [branding](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/design/branding/skills/), [motion](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/design/motion/skills/), [product](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/design/product/skills/) |
| [business](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/business/) | [strategy](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/business/strategy/skills/), [operations](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/business/operations/skills/), [leadership](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/business/leadership/skills/), [entrepreneurship](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/business/entrepreneurship/skills/), [hr](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/business/hr/skills/) |
| [science](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/science/) | [biology](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/science/biology/skills/), [physics](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/science/physics/skills/), [chemistry](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/science/chemistry/skills/), [mathematics](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/science/mathematics/skills/), [earth-science](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/science/earth-science/skills/), [astronomy](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/science/astronomy/skills/) |
| [marketing](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/marketing/) | [seo](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/marketing/seo/skills/), [content](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/marketing/content/skills/), [social-media](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/marketing/social-media/skills/), [paid-ads](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/marketing/paid-ads/skills/), [growth](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/marketing/growth/skills/), [analytics](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/marketing/analytics/skills/) |
| [health](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/health/) | [fitness](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/health/fitness/skills/), [nutrition](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/health/nutrition/skills/), [mental-health](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/health/mental-health/skills/), [sleep](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/health/sleep/skills/), [medicine](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/health/medicine/skills/) |
| [finance](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/finance/) | [personal-finance](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/finance/personal-finance/skills/), [investing](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/finance/investing/skills/), [accounting](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/finance/accounting/skills/), [real-estate](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/finance/real-estate/skills/), [corporate](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/finance/corporate/skills/) |
| [education](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/education/) | [curriculum](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/education/curriculum/skills/), [teaching](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/education/teaching/skills/), [e-learning](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/education/e-learning/skills/), [assessment](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/education/assessment/skills/), [learning-science](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/education/learning-science/skills/) |
| [film](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/film/) | [cinematography](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/film/cinematography/skills/), [directing](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/film/directing/skills/), [editing](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/film/editing/skills/), [screenwriting](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/film/screenwriting/skills/), [production](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/film/production/skills/) |
| [law](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/law/) | [contracts](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/law/contracts/skills/), [ip](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/law/ip/skills/), [employment](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/law/employment/skills/), [privacy](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/law/privacy/skills/), [corporate](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/law/corporate/skills/) |
| [photography](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/photography/) | [composition](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/photography/composition/skills/), [lighting](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/photography/lighting/skills/), [editing](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/photography/editing/skills/), [genres](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/photography/genres/skills/) |
| [music](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/music/) | [composition](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/music/composition/skills/), [production](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/music/production/skills/), [mixing](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/music/mixing/skills/), [theory](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/music/theory/skills/), [performance](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/music/performance/skills/) |
| [cooking](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/cooking/) | [techniques](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/cooking/techniques/skills/), [baking](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/cooking/baking/skills/), [flavor](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/cooking/flavor/skills/), [nutrition](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/cooking/nutrition/skills/), [world-cuisine](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/cooking/world-cuisine/skills/) |
| [language](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/language/) | [learning](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/language/learning/skills/), [linguistics](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/language/linguistics/skills/), [translation](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/language/translation/skills/), [communication](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/language/communication/skills/) |
| [art](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/art/) | [drawing](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/art/drawing/skills/), [painting](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/art/painting/skills/), [digital-art](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/art/digital-art/skills/), [illustration](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/art/illustration/skills/), [color-theory](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/art/color-theory/skills/) |
| [sports](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/sports/) | [training](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/sports/training/skills/), [coaching](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/sports/coaching/skills/), [nutrition](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/sports/nutrition/skills/), [tactics](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/sports/tactics/skills/), [recovery](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/sports/recovery/skills/) |
| [productivity](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/productivity/) | [time-management](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/productivity/time-management/skills/), [habits](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/productivity/habits/skills/), [focus](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/productivity/focus/skills/), [goals](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/productivity/goals/skills/), [tools](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/productivity/tools/skills/) |
| [travel](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/travel/) | [planning](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/travel/planning/skills/), [budgeting](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/travel/budgeting/skills/), [cultural](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/travel/cultural/skills/), [adventure](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/travel/adventure/skills/) |
| [psychology](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/psychology/) | [cognitive](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/psychology/cognitive/skills/), [behavioral](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/psychology/behavioral/skills/), [social](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/psychology/social/skills/), [clinical](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/psychology/clinical/skills/), [positive](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/psychology/positive/skills/) |
| [home](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/home/) | [renovation](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/home/renovation/skills/), [interior-design](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/home/interior-design/skills/), [gardening](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/home/gardening/skills/), [organization](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/home/organization/skills/), [smart-home](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/home/smart-home/skills/) |
| [environment](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/environment/) | [sustainability](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/environment/sustainability/skills/), [ecology](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/environment/ecology/skills/), [climate](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/environment/climate/skills/), [energy](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/environment/energy/skills/), [policy](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/environment/policy/skills/) |
| [pets](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/pets/) | [dogs](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/pets/dogs/skills/), [cats](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/pets/cats/skills/), [training](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/pets/training/skills/), [nutrition](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/pets/nutrition/skills/), [health](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/pets/health/skills/) |
| [fashion](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/fashion/) | [styling](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/fashion/styling/skills/), [wardrobe](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/fashion/wardrobe/skills/), [design](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/fashion/design/skills/), [sustainability](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/fashion/sustainability/skills/), [accessories](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/fashion/accessories/skills/) |
| [parenting](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/parenting/) | [infant](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/parenting/infant/skills/), [toddler](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/parenting/toddler/skills/), [school-age](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/parenting/school-age/skills/), [teen](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/parenting/teen/skills/) |
| [automotive](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/automotive/) | [maintenance](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/automotive/maintenance/skills/), [troubleshooting](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/automotive/troubleshooting/skills/), [buying](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/automotive/buying/skills/), [modifications](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/automotive/modifications/skills/), [ev](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/automotive/ev/skills/) |

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

**Some frameworks here are universally known — isn't the value already in the model?**

The question isn't whether the model knows the framework name. It's whether the skill encodes what practitioners who already know the name still get wrong.

A framework qualifies when the skill has substantial content beyond the acronym or label — the step most people skip, the failure mode they don't avoid, the discipline that separates expert application from surface-level application. SWOT, for example, is universally known, but most practitioners stop at the 4-quadrant list and never derive the TOWS cross-matrix (SO/WO/ST/WT strategies) — the step that converts a diagnosis into actionable options. The skill encodes that gap.

A framework doesn't qualify when the full implementation reduces to restating the framework name. If a skill's entire content would be "follow the acronym," it adds nothing — the model already knows the letters.

The test: *"What would this skill contain beyond the framework name?"* If the answer is "the step practitioners skip + the failure mode they don't avoid" — it qualifies. If the answer is "the letters, explained" — it doesn't.

**Does grimoire conflict with my team's existing conventions?**

Skills describe what the world's top institutions do. Your team may do things differently — and be right to. Two ways to handle it:

**Pin your preference.** Tell grimoire which approach to follow when practices conflict:

```
User: We follow Google's engineering practices, not IBM's.
→ Claude pins this via `pin-best-practice-preference` — applies automatically from now on.
```

**Override or fork.** A skill is a starting point, not a mandate. Adapt any skill to your context, or ignore it entirely. The format is plain Markdown and the license is MIT.

## 🤝 Contributing

**grimoire has many best practices. It needs more. Pick a domain.**

Every domain has empty sub-domains waiting for skills. If you know a field — engineering, law, finance, music, cooking, anything — add the practices you've seen work at the highest level.

**Your first skill in ~30 minutes:**
1. Pick a practice you've used at the highest level in your field
2. Run `/write-best-practice-skill` — it guides you through the format step by step
3. Open a PR — `/review-best-practice-skill` runs automatically and flags any gaps
4. Merge after review passes

Not sure where to start? Browse [open issues](https://github.com/jeffreytse/grimoire/issues) for requested skills, or pick any empty sub-domain from the table below.

Skills must pass [`review-best-practice-skill`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/review-best-practice-skill/) before merge.
The meta skills guide the full contribution workflow:

| Task | Skill |
|------|-------|
| Write a new skill | [`write-best-practice-skill`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/write-best-practice-skill/) |
| Review a skill PR | [`review-best-practice-skill`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/review-best-practice-skill/) |
| Fix review findings | [`revise-best-practice-skill`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/revise-best-practice-skill/) |
| Add a new domain | [`design-best-practice-domain`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/design-best-practice-domain/) |
| Audit a domain's health | [`audit-best-practice-domain`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/audit-best-practice-domain/) |
| Retire an outdated skill | [`deprecate-best-practice-skill`](https://github.com/jeffreytse/grimoire-skills/tree/main/skills/meta/skills/deprecate-best-practice-skill/) |

See [CONTRIBUTING.md](https://github.com/jeffreytse/grimoire-skills/blob/main/CONTRIBUTING.md) for the full standard and [GOVERNANCE.md](https://github.com/jeffreytse/grimoire-skills/blob/main/GOVERNANCE.md) for how the project and standard evolve.

## ❤️ Support

grimoire is free. It replaces $500/hr lawyers, $300 doctor visits, and $1M McKinsey
engagements — at zero cost, forever.

If it saved you time, money, or a bad decision:

- **[⭐ Star this repo](https://github.com/jeffreytse/grimoire)** — takes 2 seconds, helps thousands of people find it
- **[💖 Sponsor on GitHub](https://github.com/sponsors/jeffreytse)** — keeps the maintainer funded to add more skills across more domains
- **[☕ Ko-fi](https://ko-fi.com/jeffreytse)** · **[Patreon](https://patreon.com/jeffreytse)** · **[Liberapay](https://liberapay.com/jeffreytse)** — one-time or recurring

Every star makes grimoire more visible. Every sponsorship funds one more domain.

## 📄 License

This project is licensed under the [MIT license](https://opensource.org/licenses/mit-license.php) © Jeffrey Tse.
