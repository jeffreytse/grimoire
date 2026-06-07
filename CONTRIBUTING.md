# Contributing to grimoire

## Quality Standard

All skills must meet the [Skill Standards](./STANDARD.md). Read it before writing your first skill.
Use [SKILL_TEMPLATE.md](./SKILL_TEMPLATE.md) as your starting point.

The bar: would this skill make an AI agent perform like a domain expert?
If not, it's not ready.

---

## Adding a Skill

Skills live at: `skills/<domain>/<subdomain>/skills/<skill-name>/SKILL.md`

```
skills/engineering/development/skills/code-review/SKILL.md
skills/health/fitness/skills/design-training-program/SKILL.md
skills/finance/investing/skills/dcf-valuation/SKILL.md
skills/law/contracts/skills/review-saas-agreement/SKILL.md
```

Steps:
0. Run `suggest-best-practice` with your topic — if any result scores ≥ 0.7, extend that skill instead of creating a new one.
1. Copy `SKILL_TEMPLATE.md` → `skills/<domain>/<subdomain>/skills/<skill-name>/SKILL.md`
2. Fill in all frontmatter fields: `name`, `description`, `source`, `tags`
3. Write content following [STANDARD.md](./STANDARD.md)
4. Run the review checklist in STANDARD.md
5. Test with at least one AI agent before submitting

---

## Adding a Sub-domain

1. Create `<domain>/<subdomain>/.claude-plugin/plugin.json`:
   ```json
   {
     "name": "grimoire-<domain>-<subdomain>",
     "description": "<Subdomain> skills: <brief description>.",
     "version": "0.1.0",
     "author": { "name": "Your Name", "email": "you@example.com" },
     "homepage": "https://github.com/jeffreytse/grimoire",
     "repository": "https://github.com/jeffreytse/grimoire",
     "license": "MIT",
     "skills": "./skills"
   }
   ```
2. Add `"./<subdomain>/skills"` to the domain's `.claude-plugin/plugin.json` `skills` array
3. Add entry to `.claude-plugin/marketplace.json`
4. Add sub-domain link to the domain's row in `README.md`

---

## Adding a Domain

1. Create `<domain>/.claude-plugin/plugin.json` listing sub-domains
2. Add entry to `.claude-plugin/marketplace.json`
3. Add `"./<domain>"` to root `.claude-plugin/plugin.json` `domains` array
4. Add row to `README.md` domain table
5. Add domain to `AGENTS.md`

---

## Pull Request Guidelines

- One skill per PR when possible — easier to review
- Skill name follows naming standard: kebab-case, verb-first **required**, specific subject, 2–4 words. See `STANDARD.md` `name` section for approved verbs and examples.
- `description` starts with "Use when" and describes triggering conditions only
- No vendor lock-in unless the skill is explicitly for that tool
- Include the STANDARD.md review checklist as a PR description checklist

---

## Self-Check Before Submitting

Run this checklist against your SKILL.md before opening a PR:

### Deduplication
- [ ] No near-duplicate: ran `suggest-best-practice`, top similarity < 0.7 — or added `duplicate-reviewed: true` to frontmatter with justification in PR body

### Frontmatter
- [ ] `name` passes naming standard: verb-first, specific subject, 2–4 words, no rejected verbs (`handle-`, `manage-`, `improve-`, `get-`, `use-`, `help-`)
- [ ] `description` starts with "Use when" and describes triggering conditions only (not what the skill does)
- [ ] `source` names a specific institution, standard body, or top-tier companies (not "widely known" or "industry standard")
- [ ] `tags` has 3–8 tags covering all 4 axes: problem keyword, tool/method, role/context, outcome

### Content
- [ ] `## Why This Is Best Practice` section present
- [ ] **Adopted by**: names specific companies or institutions — not "many companies"
- [ ] **Impact**: contains a number (%, ratio, time unit) or named study
- [ ] **Why best**: mentions at least one alternative and explains why this approach wins
- [ ] Steps are concrete and immediately executable — not theory
- [ ] Scoped to one concept (if you're wondering whether to split it, split it)
- [ ] Covers edge cases and failure modes — not just the happy path

### Safety (domain-specific)
- [ ] Health/medicine: evidence tier tags present; healthcare provider footer
- [ ] Law: jurisdiction stated; "not legal advice" footer
- [ ] Finance: risk disclosure; "not financial advice" footer
- [ ] Psychology: no diagnosis/prescription; mental health professional footer

### Size
- [ ] 50–300 lines (under 50 = too shallow, over 300 = split it)

---

## License

By contributing, you agree your contributions are licensed under [MIT](./LICENSE).
