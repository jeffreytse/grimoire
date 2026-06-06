# Glossary

Terms used throughout grimoire with specific meanings. When these terms appear in STANDARD.md, meta skills, or contribution docs, this is what they mean.

---

## A

**Adopted by**
The field in the "Why This Is Best Practice" section that names which organizations or institutions use the practice. Must name at least two specific organizations or one named institution. Claims must state majority adoption, not just notable examples. "Google, Netflix, and most software companies" passes; "some companies" fails.

---

## B

**Best practice**
A technique or approach that is: (1) used by MOST top-tier organizations or credentialed professionals in the domain, (2) has measurable, cited impact, and (3) has no genuine top-tier consensus against it. An interesting technique, a company convention, or an emerging trend does not qualify. See also: *majority adoption*, *top-tier*.

---

## C

**Checklist verification**
One of the two required testing passes for a skill. Runs `review-best-practice-skill` against the SKILL.md to verify all required format criteria are met. The other pass is *behavioral testing*.

**Contested practice**
A practice where credentialed top-tier professionals are split approximately 50/50. No grimoire skill can encode a contested practice as a best practice. If a clear majority position exists (e.g., 70/30), encode the majority position and acknowledge the minority argument. See also: *majority adoption*.

**Contributor**
Anyone opening a PR to add or revise a skill. Contributors write skills and self-review with `review-best-practice-skill` but cannot merge their own PRs. See also: *maintainer*.

**Credentialed professionals**
Practitioners with formal certification or accreditation in their domain — e.g., board-certified physicians, licensed attorneys, CFA charterholders, NSCA-CSCS certified coaches. Used as the reference population for "majority adoption" in regulated domains where top-tier companies aren't the relevant unit.

---

## D

**Deprecated skill**
A skill marked for removal because the practice has been superseded or never qualified. Deprecation is a two-step process: (1) add deprecation notice and `deprecated: true` in marketplace.json, (2) remove after one release cycle. Use `deprecate-best-practice-skill` to guide the process.

**Domain**
A top-level category in the grimoire skill library — e.g., `engineering`, `health`, `law`, `finance`. A domain requires at least 3 separable sub-domains before it can be added as a top-level domain. Each domain has its own `plugin.json`. See also: *subdomain*.

---

## E

**Evidence tier**
A tag applied to factual claims in health/medicine skills indicating how strong the supporting evidence is. Required on every factual claim in health skills. The four tiers:

| Tag | Meaning |
|-----|---------|
| `[RCT]` | Randomized controlled trial |
| `[SR]` | Systematic review or meta-analysis |
| `[Consensus]` | Expert consensus or clinical guideline |
| `[Practical]` | Widely used practice without strong RCT support |

`[RCT]` and `[SR]` are strongest. Unlabeled factual claims in health skills are a REJECT finding.

---

## F

**Four axes**
The four required dimensions that a skill's tags must cover: *problem keyword*, *tool/method*, *role/context*, and *outcome*. Tags missing any axis weaken `suggest-best-practice` routing. See [tags.md](./tags.md) for accepted values per axis.

**Frontmatter**
The YAML block at the top of every SKILL.md, enclosed in `---`. Contains four required fields: `name`, `description`, `source`, `tags`. Missing or malformed frontmatter is a REJECT finding.

---

## I

**Impact**
The field in "Why This Is Best Practice" that states the measurable outcome of applying the practice. Must contain a number (%, ratio, time unit) OR a named study. "Significantly improves" without data fails. If a specific number can't be verified, keep the causal argument without the number rather than citing an uncertain source.

---

## J

**Jurisdiction**
Required in law skills: the legal system or geographic scope the skill applies to. Must be stated explicitly in the first paragraph — e.g., "US law (federal and Delaware state)" or "EU GDPR." "International" or "global" is not a valid jurisdiction; law is jurisdiction-specific. Missing jurisdiction is a REJECT finding in law skills.

---

## M

**Maintainer**
A project contributor with merge access. Maintainers review PRs using `review-best-practice-skill`, run periodic domain audits, and execute deprecation removals. Contributors write skills; maintainers gate quality and merge.

**Majority adoption**
The threshold a practice must clear to qualify for  used by MOST top-tier organizations or credentialed professionals in the domain, not just a few notable ones. "Most" means you could name 5+ organizations without straining. See also: *top-tier*.

**MECE**
Mutually Exclusive, Collectively Exhaustive. A decomposition method where sub-problems don't overlap and together cover the whole problem. Used by `plan-best-practice-solution` to break complex multi-domain problems into skill-sized pieces.

**Meta skill**
A skill that operates on the grimoire framework itself rather than on a domain subject. There are 9 meta skills in `skills/grimoire/`: three user-facing (suggest-best-practice, plan-best-practice-solution, review-best-practice-fit) and six contributor/maintainer skills (write-best-practice-skill, review-best-practice-skill, revise-best-practice-skill, audit-best-practice-domain, deprecate-best-practice-skill, design-best-practice-domain). See [grimoire-skills.md](./grimoire-skills.md).

---

## N

**NEEDS-REVISION**
One of three `review-best-practice-skill` verdicts. The skill has specific fixable findings but doesn't require a full rewrite. Each finding identifies the criterion, what's wrong, and how to fix it. Use `revise-best-practice-skill` to address NEEDS-REVISION findings. See also: *PASS*, *REJECT*.

---

## O

**Outcome**
One of the four required tag axes. Tags the result the user achieves by applying the skill — e.g., `defect-prevention`, `consistency`, `domain-release`. The most commonly missing axis.

---

## P

**PASS**
One of three `review-best-practice-skill` verdicts. All criteria met — the skill is ready to open a PR. See also: *NEEDS-REVISION*, *REJECT*.

**Plugin**
A configuration file that tells an AI agent how to load grimoire's skills. Each supported agent has its own plugin format (`.claude-plugin/plugin.json` for Claude Code, `gemini-extension.json` for Gemini CLI, etc.). See [agents.md](./agents.md).

**Problem keyword**
One of the four required tag axes. Tags what situation or problem the skill addresses — e.g., `commit-messages`, `gap-analysis`, `skill-review`. The most intuitive axis; contributors rarely miss it.

---

## R

**Regulated domain**
A domain where skills can influence decisions with real consequences — health/medicine, law, finance/investing, psychology/mental health. Skills in regulated domains require additional safeguards: evidence tier tags (health), jurisdiction statements (law), performance disclaimers (finance), and exact required disclaimers at the end. See [domain-safety.md](./domain-safety.md).

**REJECT**
One of three `review-best-practice-skill` verdicts. A required field is missing or a core criterion fails — the skill cannot be routed correctly by `suggest-best-practice` and fails the basic contract of the skill format. Fix the blocking issue before addressing anything else. REJECT on frontmatter blocks further evaluation. See also: *PASS*, *NEEDS-REVISION*.

**Role/context**
One of the four required tag axes. Tags who uses the skill — e.g., `developer`, `maintainer`, `reviewer`. Must be specific; `practitioner` is a fallback for broad-audience skills only.

---

## S

**Safety footer**
A required disclaimer at the end of every skill in a regulated domain. Exact text required — not paraphrased. Missing safety footer is a REJECT finding. See [domain-safety.md](./domain-safety.md) for the exact text per domain.

**Seed skill**
A skill included when a new domain is first added to grimoire. Minimum 2 seed skills required for any new domain PR. Seed skills demonstrate the domain's scope and prove the domain structure works before merge. All seed skills must pass `review-best-practice-skill` before the domain PR opens.

**Skill**
A SKILL.md file encoding a single best practice. A skill covers one concept, has a specific audience, names specific adopters, cites measurable impact, and passes `review-best-practice-skill`. Not a prompt, not a template, not a guideline — an evidence-backed instruction.

**Skill scope**
The boundary of what a skill covers. One concept per skill. If two concepts are separable — meaning you could apply one without the other — they belong in two skills. "One concept" test: can you describe the skill's trigger in one clause? If not, it may cover two concepts.

**Source**
A required frontmatter field listing the institution, standard body, or named top-tier organizations whose practices the skill encodes. Must be verifiable — a reviewer must be able to look it up. Personal blogs, anonymous "studies," and fabricated citations are REJECT findings.

**Subdomain**
A category within a domain. Engineering has subdomains: development, testing, devops, security, etc. A subdomain has its own `plugin.json` and `marketplace.json` entry. Three or more separable subdomains are required before a top-level domain qualifies. See also: *domain*.

---

## T

**Tool/method**
One of the four required tag axes. Tags the technique or approach the skill applies — e.g., `checklist`, `mece`, `situation-analysis`. Sometimes implicit in the problem; the axis is still required.

**Top-tier**
The reference population for "majority adoption." In engineering: companies like Google, Netflix, Meta, Stripe, Amazon — organizations operating at significant scale with documented engineering practices. In regulated domains: credentialed professionals (board-certified physicians, licensed attorneys, etc.) rather than companies. "Top-tier" is not "well-known" or "innovative" — it means organizations whose practices set the standard others follow.

**Triggering condition**
The situation that causes a skill to activate. Defined in the skill's `description` frontmatter field. Must describe WHEN to use the skill, not WHAT the skill does. `suggest-best-practice` uses triggering conditions to route situations to the right skill.

---

## V

**Verdict**
The output of `review-best-practice-skill`: PASS, NEEDS-REVISION, or REJECT, plus a per-criterion findings table. Verdicts are deterministic — the same skill file produces the same verdict regardless of who runs the review.

---

## W

**Why best**
The field in "Why This Is Best Practice" that explains why this approach wins over the named alternative. Must name at least one alternative and give the causal reason this approach beats it. "This is the best approach" without naming an alternative fails.

**Why This Is Best Practice**
The required section that proves a skill belongs in grimoire. Contains three required fields: *Adopted by* (who uses it), *Impact* (measurable outcome with citation), *Why best* (comparison to named alternative). This section is the quality signal that distinguishes grimoire from a prompt library.
