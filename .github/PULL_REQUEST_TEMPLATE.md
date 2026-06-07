## Summary

<!-- What skill(s) does this PR add, fix, or deprecate? One sentence per skill. -->

## Self-check

Run `review-best-practice-skill` on every SKILL.md in this PR, then check each box.

### Deduplication

- [ ] No near-duplicate: ran `suggest-best-practice`, top similarity < 0.7 — or added `duplicate-reviewed: true` to frontmatter with justification below

### Frontmatter

- [ ] `name` passes naming standard: verb-first, specific subject, 2–4 words, no rejected verbs (`handle-`, `manage-`, `improve-`, `get-`, `use-`, `help-`)
- [ ] `description` starts with "Use when"
- [ ] `description` describes triggering conditions only — no workflow summary
- [ ] `source` names a specific institution, standard body, or top-tier adopter list (not `grimoire STANDARD.md`)
- [ ] `tags` has 3–8 tags covering all 4 axes: problem keyword, tool/method, role/context, outcome

### Content

- [ ] `## Why This Is Best Practice` section present
- [ ] **Adopted by**: names specific companies or institutions (not "many companies")
- [ ] **Impact**: contains a % number or named study (not "significantly improves")
- [ ] **Why best**: explicitly names an alternative approach and explains why this one wins
- [ ] **Sources**: external, verifiable (not self-referential)
- [ ] Steps are immediately executable — commands, not advice
- [ ] Scoped to one concept (no "and" in skill name or scope)
- [ ] Covers failure modes and edge cases

### Safety (if applicable)

- [ ] Health/medicine: evidence tier tags (`[RCT]`, `[SR]`, `[Consensus]`); healthcare provider footer
- [ ] Law: jurisdiction scope stated; legal advice disclaimer footer
- [ ] Finance: past-performance disclaimer; financial advice footer
- [ ] Mental health: no diagnosis or medication recommendations; mental health professional footer

### Size

- [ ] 50–300 lines

### Deprecations (if any)

<!-- List any existing skills this PR supersedes. Include the path of the skill to deprecate. -->
