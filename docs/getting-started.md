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

## Find the right skill: suggest-practice

You don't need to know which skill applies. Describe your situation:

```
User: /suggest-practice
```

Then describe the problem in natural language. `suggest-practice` auto-classifies the situation and routes to the best match:

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

If a skill isn't installed, `suggest-practice` tells you exactly what to install. It never pretends a skill exists when it doesn't.

---

## Multi-domain problems: plan-solution

Some problems span multiple disciplines. `plan-solution` decomposes them using MECE (Mutually Exclusive, Collectively Exhaustive) analysis and sequences the applicable skills in dependency order.

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

`suggest-practice` routes to `plan-solution` automatically when it detects a multi-domain problem. You can also invoke it directly.

---

## Evaluate your existing solution: review-practice-fit

If you already have a design, plan, or process and want to know how well it aligns with best practices:

```
User: I have an API design: REST endpoints, JWT auth, PostgreSQL, no rate limiting,
      logs to console only.

Claude: You have an existing solution. Applying review-practice-fit...

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

`review-practice-fit` evaluates your solution against each applicable best practice — ALIGNED, PARTIAL, or MISSING — and produces a prioritized fix list.

`suggest-practice` routes to `review-practice-fit` automatically when it detects you're describing an existing solution ("is this good?", "what am I missing?").

---

## When a skill isn't installed

`suggest-practice` always tells you what to install:

```
Claude: Situation matches: review-saas-contract (law/contracts)
        Not installed — get it with:
        /plugins add github:jeffreytse/grimoire/skills/law/contracts
```

After installing, invoke `suggest-practice` again or call the skill directly.

If no skill in grimoire covers your situation yet, `suggest-practice` says so and asks a clarifying question to narrow the domain. You can then request the skill via a [GitHub issue](../.github/ISSUE_TEMPLATE/new-skill.md).

---

## Next steps

- **Browse all domains**: see the [README domains table](../README.md#domains)
- **Contribute a skill**: read [Authoring Skills](./authoring-skills.md)
- **Request a missing skill**: open a [new skill request](../.github/ISSUE_TEMPLATE/new-skill.md)
- **Report a skill issue**: open a [skill revision request](../.github/ISSUE_TEMPLATE/skill-revision.md)
