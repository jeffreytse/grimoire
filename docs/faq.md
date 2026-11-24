# FAQ

Answers to common questions about grimoire.

← [Back to the README](../README.md)

**How is grimoire different from other skill package managers or agents' plugin marketplaces?**

They distribute skills: find, install, version and share them, and some evaluate them. That's a solved and useful layer, and grimoire does it too. Grimoire's focus is the next step: getting the agent to *apply* the right practice. It routes each request to the practices that fit, serves a 1,300+ skill library on demand without filling the agent's context, shows its decisions so you can redirect them, respects your project's declared standards, and checks compliance in CI. Grimoire's skills also work as plain `SKILL.md` folders, so they install through agents' native plugin commands as well.

**MCP mode or skills mode?**

MCP mode, wherever the agent supports it (Claude Code, Codex, Gemini CLI, Copilot CLI, OpenCode, Pi, Kimi Code, OMP, Antigravity). Agents put every installed skill's name and description into a fixed listing — about 8,000–16,000 characters in Claude Code and Codex — and silently truncate or drop what doesn't fit. grimoire-core's descriptions total ~320,000 characters, so a full skills-mode install leaves the agent able to discover only a few percent of it. MCP mode keeps the listing empty and finds skills by search when a task needs them, and adds grimoire's hooks and display. Skills mode is for agents without MCP support, or a deliberately small subset that fits the listing — about 30 skills. Skills mode stays fully supported, but new capabilities arrive in MCP mode first. Project standards, organization/user/project layers and search-based selection, for example, need MCP mode.

**Isn't this already in the model's training data?**

Yes — and no. Models know _about_ best practices. Skills make models _do_ them, reliably.

The difference:

| Without a skill                            | With a skill                                                          |
| ------------------------------------------ | --------------------------------------------------------------------- |
| Model improvises a version of the practice | Model follows the exact steps from the source institution             |
| Output varies every run                    | Same process, same structure, every time                              |
| Practice applied only if you know to ask   | Skill triggers automatically when the situation matches               |
| Generic advice                             | Specific: the right gate, the right question, the right output format |

For simple tasks (write a test, fix a bug), the skeptic is right — the model doesn't need a skill. For complex, multi-step workflows — an SLO design, a post-mortem, an incident response — skills measurably change what you get. The model knows Google's engineering review process exists. It does not reliably know which question to ask first, what the output format is, or when to stop. That's what a skill encodes.

**These are just textbook practices the model already knows. Why bother?**

Knowing a practice and reliably executing it are different things. Ask any model "I just had a production incident" — you'll get a generic write-up. Run `write-post-mortem` and you get: blameless framing, 5-whys, timeline, contributing factors, action items with owners, and a detection section. The model _knew_ all of that before the skill existed. The skill is what makes it happen consistently, in the right format, every time.

The "textbook" objection gets it backwards. Established practices are _ideal_ for skills precisely because they're falsifiable — you can verify whether the output matches what Google's SRE book, Amazon's mechanisms, or the WHO protocol actually prescribes. If you find a skill that adds nothing over a bare prompt, that's a quality failure. [File an issue.](https://github.com/jeffreytse/grimoire/issues)

**Some frameworks here are universally known — isn't the value already in the model?**

The question isn't whether the model knows the framework name. It's whether the skill encodes what practitioners who already know the name still get wrong.

A framework qualifies when the skill has substantial content beyond the acronym or label — the step most people skip, the failure mode they don't avoid, the discipline that separates expert application from surface-level application. SWOT, for example, is universally known, but most practitioners stop at the 4-quadrant list and never derive the TOWS cross-matrix (SO/WO/ST/WT strategies) — the step that converts a diagnosis into actionable options. The skill encodes that gap.

A framework doesn't qualify when the full implementation reduces to restating the framework name. If a skill's entire content would be "follow the acronym," it adds nothing — the model already knows the letters.

The test: _"What would this skill contain beyond the framework name?"_ If the answer is "the step practitioners skip + the failure mode they don't avoid" — it qualifies. If the answer is "the letters, explained" — it doesn't.

**Does grimoire conflict with my team's existing conventions?**

Skills describe what the world's top institutions do. Your team may do things differently — and be right to. Two ways to handle it:

**Pin your preference.** Tell grimoire which approach to follow when practices conflict:

```
User: We follow Google's engineering practices, not IBM's.
→ Claude pins this via `pin-best-practice-preference` — applies automatically from now on.
```

**Override or fork.** A skill is a starting point, not a mandate. Adapt any skill to your context, or ignore it entirely. The format is plain Markdown and the license is MIT.
