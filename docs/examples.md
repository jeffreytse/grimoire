# Examples

Skills in action, and ways to start.

← [Back to the README](../README.md)

## Skills in action

**Engineering** — propose a commit message from staged diff:

```
User: grimoire:propose-conventional-commit

Claude: Inspecting staged changes...

  feat(auth): add OAuth2 PKCE flow for mobile clients

  Replaces implicit grant with PKCE to address RFC 6749 §10.16
  implicit flow vulnerabilities in native apps.

  Approve? [y/n]
```

Every skill is structured: triggering condition, verified source, numbered steps, no freeform guessing.

**Or describe the situation in plain language** — you don't need to know which skill applies:

<pre>
❯ My team keeps missing sprint goals and I don't know why.

grimoire ▸ finding a practice for missed sprint goals — a problem without a chosen approach

<b>Options</b>
1. apply-systems-iceberg — trace the repeated misses past the surface events to the patterns
                           and structures causing them ★ recommended
2. plan-retrospective — run a retro with the team to find blockers and agree on improvements
                        with owners
3. plan-sprint — check how you set goals, estimate capacity and choose stories, in case the
                 commitment is the problem
4. run-team-health-check — measure clarity, workload, trust and psychological safety for
                           causes that aren't about process

Reply 1–4, "all", or "skip".
</pre>

`suggest-best-practice` auto-classifies any situation, routes to the matching skill,
or tells you exactly what to install if the skill isn't in your library yet.

→ [See skills across all 26 domains — law, finance, health, cooking, and more](https://github.com/jeffreytse/grimoire-core)

## Ways to start

**New to grimoire? Run the interactive wizard:**

```bash
grimoire wizard
```

The wizard asks how agents should get skills (MCP mode first, recommended), then the scope and agents — no manual config editing. It detects your installed AI agents, sets up MCP mode for those that support it, and links skills into the rest.

---

**After install, describe any problem in plain language:**

<pre>
❯ I need to raise a Series A but don't know how to pitch investors.

grimoire ▸ finding a practice for pitching a Series A — goal stated, no approach chosen yet

<b>Options</b>
1. design-fundraising-strategy — plan the Series A: how much to raise, acceptable dilution,
                                 which investors to target, how to run the process ★
                                 recommended
2. design-presentation-structure — build the pitch story and the order of your evidence so
                                   investors move toward a yes
3. calculate-tam-sam-som — size your market so it holds up when investors check your numbers
4. write-business-plan — write the formal document behind the pitch, for investors who ask
                         for one

Reply 1–4, "all", or "skip".
</pre>

Or invoke a skill directly. In MCP mode, write `grimoire:<skill-name> [task]` in any agent, or `/grimoire:skill <skill-name>` in Claude Code, e.g. `grimoire:suggest-best-practice I need to raise a Series A`. In skills mode, each skill is a slash command:

```
/suggest-best-practice     # describe any problem — auto-routes to the right skill
/review-pull-request       # engineering code review
/calculate-fire-number     # how much do I need to retire?
/review-saas-contract      # flag dangerous clauses before signing
/design-training-program   # build a training program
```

For CI enforcement, initialize once per project and gate PRs with the `grimoire` CLI:

```bash
grimoire init    # creates grimoire.toml with auto-detected profile
grimoire check   # exits 0 (pass) or 1 (fail)
```

**New to grimoire?** Start with `/suggest-best-practice`. Describe any professional or life situation — it reads your context and routes you to the matching skill, or tells you exactly what to install if the skill isn't in your library yet.
