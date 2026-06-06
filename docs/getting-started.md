# Getting Started with grimoire

grimoire delivers professional best practices from every field directly into your AI agent. This guide walks through installation, your first skill, and the three meta-skills that help you find and apply the right practice for any situation.

---

## Install

### Claude Code

```bash
# All skills across all domains
/plugins add github:jeffreytse/grimoire

# One domain
/plugins add github:jeffreytse/grimoire/skills/engineering

# One sub-domain
/plugins add github:jeffreytse/grimoire/skills/engineering/development
```

### Codex / Cursor

```bash
/plugins add github:jeffreytse/grimoire
```

### OpenCode

See [`.opencode/INSTALL.md`](../.opencode/INSTALL.md) for plugin configuration.

### Gemini CLI

```bash
curl -fsSL https://raw.githubusercontent.com/jeffreytse/grimoire/main/scripts/install.sh | bash --target gemini
```

### Script install (all agents)

```bash
curl -fsSL https://raw.githubusercontent.com/jeffreytse/grimoire/main/scripts/install.sh | bash

# Options
./scripts/install.sh --domain engineering
./scripts/install.sh --domain engineering --subdomain development
./scripts/install.sh --skill engineering/development/propose-conventional-commit
./scripts/install.sh --target all    # install for all detected agents
```

---

## Your first skill

Once the engineering domain is installed, try the reference skill:

```
User: /propose-conventional-commit

Claude: Inspecting staged changes...

  feat(auth): add OAuth2 PKCE flow for mobile clients

  Replaces implicit grant with PKCE to address RFC 6749 §10.16
  implicit flow vulnerabilities in native apps.

  Approve? [y/n]
```

The skill reads `git diff --cached`, extracts scope and intent from the diff, and drafts a Conventional Commits message. You approve or reject — nothing is committed without your confirmation.

---

## Find the right skill: suggest-best-practice

You don't need to know which skill applies. Describe your situation:

```
User: /suggest-best-practice
```

Then describe the problem in natural language. `suggest-best-practice` auto-classifies the situation and routes to the best match:

```
User: My pull requests always take 3 days to get reviewed and then get rejected.

Claude: Situation matches: review-pull-request (engineering/development)
        Not installed — get it with:
        /plugins add github:jeffreytse/grimoire/skills/engineering/development
```

**Browse mode** — if you want to see what's available without applying anything:

```
User: What practices exist for database performance?

Claude: Best practices for: database performance

  1. optimize-query-latency — eliminate slow queries before they reach production
     Domain: engineering/data  |  Install: /plugins add ...

  2. design-schema-indexing — index strategy for high-read workloads
     Domain: engineering/data  |  Install: /plugins add ...
```

If a skill isn't installed, `suggest-best-practice` tells you exactly what to install. It never pretends a skill exists when it doesn't.

---

## Multi-domain problems: plan-best-practice-solution

Some problems span multiple disciplines. `plan-best-practice-solution` decomposes them using MECE (Mutually Exclusive, Collectively Exhaustive) analysis and sequences the applicable skills in dependency order.

```
User: I want to launch a SaaS startup.

Claude: Your situation spans 5 domains. Here is the solution plan:

  1. design-business-model — validate revenue model before building
     Domain: business/strategy

  2. calculate-unit-economics — confirm LTV/CAC before committing to go-to-market
     Domain: finance/corporate

  3. review-saas-contract — legal structure before signing vendors or hiring
     Domain: law/contracts

  4. design-go-to-market — customer acquisition after unit economics confirmed
     Domain: marketing/growth

  5. plan-hiring — team structure after role and revenue model are clear
     Domain: business/hr

  Apply in order? I'll pause for your confirmation after each step.
```

Skills run one at a time, with your confirmation between each. The plan adapts if an earlier step reveals new constraints.

`suggest-best-practice` routes to `plan-best-practice-solution` automatically when it detects a multi-domain problem. You can also invoke it directly.

---

## Evaluate your existing solution: review-best-practice-fit

If you already have a design, plan, or process and want to know how well it aligns with best practices:

```
User: I have an API design: REST endpoints, JWT auth, PostgreSQL, no rate limiting,
      logs to console only.

Claude: You have an existing solution. Applying review-best-practice-fit...

        ### design-api-architecture — PARTIAL
        ✓ REST endpoints, stateless auth
        ✗ No versioning strategy (/v1/ prefix)
        ✗ No pagination standard
        → Fix: Add /v1/ prefix and cursor-based pagination before next release

        ### review-security-posture — MISSING
        ✗ No rate limiting
        ✗ No input validation mentioned
        → Fix: Add rate limiter (100 req/min per IP) at gateway before next deploy

        🔴 Critical: No rate limiting → DoS exposure
        🔴 Critical: Console logs only → incidents uninvestigable

        Verdict: NEEDS WORK — address rate limiting and structured logging first.
```

`review-best-practice-fit` evaluates your solution against each applicable best practice — ALIGNED, PARTIAL, or MISSING — and produces a prioritized fix list.

`suggest-best-practice` routes to `review-best-practice-fit` automatically when it detects you're describing an existing solution ("is this good?", "what am I missing?").

---

## Browse what exists first: discover-best-practices

**When to use:** You don't have a specific problem yet but want to know what best practices exist for a domain or field.

```bash
/discover-best-practices "software architecture"
```

**What it does:**
1. Scans installed skills for the detected domain
2. Groups them by subdomain with gap framing ("what goes wrong without this")
3. Highlights 1–3 practices most commonly discovered too late
4. Offers to apply any or learn more about one

## Catch gaps before you start: intercept-best-practice

**When to use:** You're about to begin a task — writing, coding, designing, planning — and want to apply the right practice before starting rather than after discovering a gap.

```bash
/intercept-best-practice "I'm about to refactor our auth module"
```

**What it does:**
1. Matches the most relevant practice for the task you're about to do
2. Surfaces the specific gaps it would prevent
3. Offers to apply it before you begin — or lets you continue without

## Define the problem first: analyze-problem

**When to use:** Your problem isn't clearly defined yet — you know something is wrong but not what exactly to solve.

```bash
/analyze-problem "My team isn't performing well"
```

**What it does:**
1. Asks clarifying questions one at a time until the problem boundary is clear
2. Produces a structured problem space map (statement / scope / constraints / root cause)
3. Surfaces 2–4 possible solution routes with trade-offs
4. Routes to `suggest-best-practice` or `plan-best-practice-solution` based on your choice

## When a skill isn't installed

`suggest-best-practice` always tells you what to install:

```
Claude: Situation matches: review-saas-contract (law/contracts)
        Not installed — get it with:
        /plugins add github:jeffreytse/grimoire/skills/law/contracts
```

After installing, invoke `suggest-best-practice` again or call the skill directly.

If no skill in grimoire covers your situation yet, `suggest-best-practice` says so and asks a clarifying question to narrow the domain. You can then request the skill via a [GitHub issue](../.github/ISSUE_TEMPLATE/new-skill.md).

---

## Next steps

- **Learn a new practice:** Use `/suggest-best-practice` for any question
- **Browse a domain:** Use `/discover-best-practices` to see what skills exist before you have a problem
- **Intercept before starting:** Use `/intercept-best-practice` before any major task
- **Define your problem:** Use `/analyze-problem` when the problem isn't clear yet
- **Contribute:** Create a new skill with `/write-best-practice-skill`
- **Review a contribution:** Audit a skill with `/review-best-practice-skill`
- **Audit a domain:** Review all skills in a domain with `/audit-best-practice-domain`
