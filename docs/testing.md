# Testing Skills

CONTRIBUTING.md says "test with at least one AI agent before opening a PR." This guide explains what that means and how to do it.

---

## Two types of testing

A skill needs two passes before it's ready for PR review:

| Type | What it checks | How to run |
|------|---------------|------------|
| **Checklist verification** | Does the file satisfy all required format criteria? | Invoke `review-skill` on the file |
| **Behavioral testing** | Does the skill actually change agent behavior? | Invoke the skill in a live agent session and test it |

Both are required. A skill that passes `review-skill` but produces wrong agent behavior isn't ready. A skill that an agent follows correctly but fails `review-skill` can't be merged.

---

## Checklist verification: review-skill

Invoke `review-skill` on your file before opening a PR:

```
Applying review-skill to skills/<domain>/<subdomain>/skills/<name>/SKILL.md
```

`review-skill` runs the full STANDARD.md checklist and produces a structured verdict:

```
## Skill Review: <skill-name>

### Overall Verdict: PASS | NEEDS-REVISION | REJECT

| Criterion     | Result | Finding |
|---------------|--------|---------|
| name          | ✅     | —       |
| description   | ✅     | —       |
| source        | ✅     | —       |
| tags          | ⚠️     | Missing outcome axis |
| Adopted by    | ✅     | —       |
| Impact        | ⚠️     | No number or named study |
| Why best      | ✅     | —       |
| Steps         | ✅     | —       |
```

**What each verdict means:**

- **PASS** — file is ready to open a PR
- **NEEDS-REVISION** — specific findings flagged; use `revise-skill` to fix each one, then re-run `review-skill`
- **REJECT** — a required field is missing or a core criterion fails; fix the blocking issue before anything else

A REJECT on any required frontmatter field (name, description, source, tags) means `suggest-practice` can't route the skill correctly. Fix it before testing behavior.

Fix every finding. Re-run `review-skill` after each round of fixes. Don't open a PR until you get PASS.

---

## Behavioral testing: does the skill work?

Checklist verification tells you the file is well-formed. Behavioral testing tells you whether an agent actually follows the skill correctly when invoked.

### What behavioral testing means

A skill is a set of instructions loaded into an agent's context. Testing behavior means:

1. Starting a fresh agent session
2. Installing or loading the skill
3. Creating a scenario that should trigger the skill
4. Observing whether the agent follows the skill's steps correctly
5. Checking the output against the skill's intended outcome

### How to test in Claude Code

```bash
# If your skill is already in an installed plugin:
# Just start a new claude session — the skill loads automatically

# If testing a new skill not yet in a plugin:
/plugins add <path-to-skill-directory>
```

Then trigger the skill naturally:

```
User: "I want to commit the current changes"
# → should trigger propose-conventional-commit if installed
# → should not hallucinate if no staging exists
```

### Designing a test scenario

A good test scenario has three parts:

**1. A triggering condition that matches the skill's description**

The skill's `description` field defines when it should activate. Your scenario must match those conditions exactly.

```yaml
# Skill description:
description: Use when the user asks to commit, wants a commit message, or invokes /propose-commit.

# Correct test trigger:
"Can you propose a commit message for these changes?"

# Wrong test trigger (too indirect):
"What should I do with these staged files?"
```

**2. An expected output you can verify**

Know what the skill is supposed to produce before you test it. If the Steps section says "present a draft commit message for approval before committing," your test should verify the agent does that — not that it silently commits.

**3. A failure condition**

Know what wrong behavior looks like. Document what the agent did wrong so you can fix the skill.

### What to check

After triggering the skill, verify:

- **Routing**: Did the agent recognize this as the right situation for the skill?
- **Steps followed**: Did the agent execute the steps in the correct order?
- **Required behaviors**: Did any "must" or "required" instructions get followed?
- **Scope**: Did the agent stay within the skill's defined scope, or drift into adjacent behavior?
- **Output format**: Does the output match what the skill specifies?

### Testing regulated domain skills

Skills in health, law, finance, and psychology domains need an extra behavioral check: verify that the required disclaimer appears in the output.

```
# Health skill test
User: "What's the best approach to progressive overload?"
Expected: skill content + footer "For personal health decisions, consult a qualified healthcare provider."

# Verify:
- Footer is present
- Footer is exact text (not paraphrased)
- No specific medication doses or diagnoses appeared in the output
- All factual claims have evidence tier tags [RCT], [SR], [Consensus], or [Practical]
```

If the footer is missing or paraphrased, the skill's Steps section is incomplete — add an explicit step directing the agent to include it.

---

## Common testing gaps

| Gap | What goes wrong | Fix |
|-----|----------------|-----|
| Testing only the happy path | Agent follows the skill for the expected trigger but breaks on edge cases | Test at least one edge case (e.g., no staged files for a commit skill) |
| Not verifying scope | Agent applies the skill too broadly or not at all | Check description triggers — are they specific enough? |
| Skipping regulated domain footer check | Footer appears in the file but agent omits it from output | Add an explicit step in the skill: "End your response with: [exact footer text]" |
| Testing in the same session you wrote the skill | You remember the skill's intent and interpret the agent charitably | Start a fresh session with no context from authoring |
| Not testing the REJECT case | You know the skill works for valid input but never verify invalid handling | If the skill has Rules or Common Mistakes, trigger one of them and verify the agent catches it |

---

## Testing checklist

Before opening a PR:

- [ ] `review-skill` produces PASS verdict
- [ ] Skill tested in a fresh agent session (not the session where you wrote it)
- [ ] Triggering condition produces correct routing
- [ ] All Steps executed in order
- [ ] Output matches skill's intended outcome
- [ ] Regulated domain: required footer present in agent output (exact text)
- [ ] At least one edge case tested
- [ ] No NEEDS-REVISION or REJECT findings remain

A skill that passes both checklist verification and behavioral testing is ready to open a PR.
