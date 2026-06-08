# Frequently Asked Questions

---

## Using grimoire

### How is grimoire different from a prompt library or AwesomePrompts?

Prompt libraries collect ad-hoc instructions — they tell an AI what to do, but don't verify that the instruction reflects what experts actually do. grimoire encodes practices adopted by top-tier organizations with cited evidence. Every skill must name specific adopters, cite measurable impact, and pass `review-best-practice-skill` before it enters the library.

The test: "Would this skill make an AI agent perform like a domain expert?" If not, it doesn't belong.

### Which AI agents does grimoire support?

| Agent | How it loads |
|-------|-------------|
| Claude Code | `/plugin marketplace add jeffreytse/grimoire` then `/plugin install grimoire@grimoire` |
| Codex | Same — uses `.codex-plugin/plugin.json` |
| Cursor | Same — uses `.cursor-plugin/plugin.json` |
| OpenCode | `.opencode/plugins/grimoire.js` ESM plugin |
| Gemini CLI | `gemini-extension.json` → `GEMINI.md` |
| Agents CLI | `AGENTS.md` at repo root |

### Can grimoire be used offline?

Skills are loaded into the agent's context at install time. Once installed, they're available locally — no network access is required during use.

### What if no grimoire skill covers my situation?

`suggest-best-practice` will say "No installed skills closely match this situation." Three options:

1. Install a domain that likely covers it: `suggest-best-practice` tells you what to install
2. Open a [new skill request](../.github/ISSUE_TEMPLATE/new-skill.md) — describe the practice and its evidence
3. Write the skill yourself — see [Authoring Skills](./authoring-skills.md)

### How do I find the right skill without knowing its name?

Invoke `suggest-best-practice` (or just describe your problem). It auto-classifies any situation and routes to the best match. You never need to know skill names in advance.

---

## Contributing skills

### What if I can't find a study or number for the Impact section?

Keep the causal argument and remove the specific number. A well-reasoned causal argument passes; "significantly improves" without data fails. The rule:

```
❌ "Significantly improves code quality"
❌ "Studies show 60% improvement"    ← which studies? this fails
✅ "Teams that skip problem definition produce solutions targeting the wrong problem,
    requiring full redesign (IDEO Human-Centered Design Field Guide, 2015)"
✅ "Reduces post-ship defects by ~50% (Google Engineering Practices documentation)"
```

If you're not certain a specific number is in a specific document, remove the number and keep the argument.

### Can I contribute practices from my company that aren't publicly documented?

If the practice is majority-adopted across top-tier organizations and you can name them, yes — cite the organizations. The standard is majority adoption at the top tier, not public documentation.

If it's your company's proprietary approach not adopted elsewhere, no. One organization's internal convention is not a best practice.

### What if a practice is contested among experts?

Two cases:

**Clear majority position**: encode it, acknowledge the minority argument. State which top-tier professionals hold each position and what evidence underlies the majority view.

**Genuine 50/50 split**: no consensus = no best practice. Skip it. If the split resolves over time (one side accumulates significantly more evidence), use `revise-best-practice-skill` to update.

### What happens if my PR doesn't pass review-best-practice-skill?

Maintainers apply `review-best-practice-skill` to every incoming skill PR. Three outcomes:

- **PASS** → merged
- **NEEDS-REVISION** → maintainer comments with specific findings; use `revise-best-practice-skill` to fix each one
- **REJECT** → a required field is missing or a core criterion fails; PR is closed with explanation; fix and re-open

Use `review-best-practice-skill` on your own file before opening a PR — it runs the same checklist. A skill that passes your self-review will pass the maintainer's review.

### How many skills do I need to write to add a new domain?

Minimum 2 seed skills. A domain with 1 skill signals an abandoned stub. Use `design-best-practice-domain` to set up the full domain structure, then `write-best-practice-skill` for each seed skill, then `audit-best-practice-domain` to verify all pass before opening the PR.

### Can I use grimoire for proprietary knowledge that can't be publicly cited?

No. Every skill requires a verifiable source — an institution, standard body, or named top-tier companies whose practices are at least publicly known (even if not fully documented). Skills that can't be verified by a reviewer fail the source criterion.

---

## Quality standards

### Why does every skill require "Why This Is Best Practice"?

Without it, there's no way for a reviewer — human or AI — to verify the skill belongs in grimoire. The section is the quality signal: it's how maintainers distinguish a proven best practice from an opinion, a company convention, or an interesting technique.

The section makes grimoire's quality claim self-enforcing: any agent that applies `review-best-practice-skill` checks the same criteria, producing a consistent verdict.

### How do I know if a skill is outdated?

Use the `deprecate-best-practice-skill` staleness criteria:

- The source institution revised its position → use `revise-best-practice-skill` (fix, not deprecate)
- The tool or API referenced no longer exists at scale → use `revise-best-practice-skill` (update reference)
- A newer practice has achieved majority top-tier adoption that directly supersedes this one → use `deprecate-best-practice-skill`
- The "Adopted by" claim is no longer accurate → use `revise-best-practice-skill` to update, or `deprecate-best-practice-skill` if the claim was always wrong

### What's the difference between revise-best-practice-skill and deprecate-best-practice-skill?

`revise-best-practice-skill`: the practice is still valid — something in the skill file is wrong or outdated. Fix specific findings without rewriting the whole skill.

`deprecate-best-practice-skill`: the practice itself has been superseded, or the skill never qualified. Mark for removal and point to the replacement.

### Why is the scope limit "one concept per skill"?

A skill covering two separable concepts produces ambiguous routing in `suggest-best-practice`, makes `review-best-practice-skill` harder to apply consistently, and prevents contributors from updating one concept independently. If you're unsure, split — two clean skills are always better than one overloaded skill.
