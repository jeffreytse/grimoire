---
name: design-agentic-workflow
description: Use when building or prompting a multi-step agentic system — to structure execution so the agent plans before acting, creates verifiable checkpoints, and avoids irreversible mistakes.
source: Anthropic, "Building effective agents" (2024); Yao et al., "ReAct: Synergizing Reasoning and Acting in Language Models" (2022, Google Brain); Amazon Bedrock Agents documentation (2023); Wei et al., "Chain-of-Thought Prompting" (2022, Google)
tags: [agents, agentic, llm, planning, multi-step, ai-engineering, reliability, irreversible-actions]
related: [apply-evaluation-driven-development, write-eval-suite]
---

# Design Agentic Workflow

Structure multi-step agentic tasks as analyze → plan → validate → execute → verify. Never execute before planning.

## Why This Is Best Practice

**Adopted by:** Anthropic (documented in "Building effective agents", 2024), Google
(ReAct pattern — Yao et al. 2022, adopted in Gemini agents and AlphaCode), Amazon
(Bedrock agents architecture requires an explicit plan step before tool execution),
and Microsoft (AutoGen and Semantic Kernel agent patterns both implement
plan-then-execute as a core construct).
**Impact:** The ReAct pattern (Reason + Act) reduced task failure rates by 34% over
direct action on HotpotQA and Fever benchmarks (Yao et al., 2022). Anthropic's agent
safety research identifies unplanned execution as the primary source of irreversible
mistakes in production agents — errors that require human intervention to recover from.
**Why best:** Direct-action agents (receive task → immediately execute tools) have no
mechanism to catch ambiguous instructions, conflicting goals, or destructive paths
before damage is done. A plan-validate-execute structure creates a natural checkpoint:
ambiguity is caught at plan time (cheap to correct), not mid-execution (expensive or
unrecoverable).

Sources: Yao et al., "ReAct" (2022); Anthropic, "Building effective agents" (2024);
Amazon Bedrock Agents documentation (2023); Wei et al., "Chain-of-Thought Prompting" (2022)

## Steps

### Step 1: Analyze inputs and surface ambiguities before planning

Before generating a plan, identify every input that could be interpreted multiple ways
or that requires information not yet available:

```
Task: "Clean up the database"
Ambiguities to surface:
- Which database? (prod, staging, test?)
- What does "clean up" mean? (delete old records? vacuum? remove orphaned rows?)
- Is there a recovery plan if something goes wrong?
```

If the agent is user-facing, ask for clarification before proceeding. If running
autonomously, apply conservative defaults and log all assumptions explicitly.

### Step 2: Generate an explicit plan as structured output

Produce a list of discrete, verifiable steps before executing any of them:

```
Plan:
1. Connect to staging database (read-only first)
2. Count rows older than 90 days in events table
3. Display sample of rows targeted for deletion
[CHECKPOINT — approval required before step 4]
4. Delete rows in batches of 1000
5. Verify row count decreased by expected amount
6. Vacuum table
```

The plan must be externalized and inspectable. An agent that plans internally without
writing the plan as output cannot be validated before execution begins.

### Step 3: Insert checkpoints before irreversible or high-blast-radius actions

Any action that cannot be undone or that affects a large scope requires an explicit pause:

| Action type | Checkpoint requirement |
|-------------|----------------------|
| Deletes (files, DB rows, messages) | Always pause before executing |
| Writes to production systems | Always pause before executing |
| External API calls with side effects | Pause unless explicitly pre-authorized |
| Reads and reversible state changes | No checkpoint needed |

Checkpoints before irreversible actions are not optional. Autonomous convenience
does not outweigh recovery cost.

### Step 4: Execute step by step, verify after each step

Never batch-execute the full plan. After each step:
1. Check that the actual outcome matches the expected outcome
2. If mismatch: stop and surface the discrepancy before continuing
3. If match: proceed to the next step

```python
for step in plan:
    result = execute(step)
    if not matches(result, step.expected_outcome):
        raise AgentError(
            f"Step {step.id} diverged: expected {step.expected_outcome}, got {result}"
        )
    log(f"Step {step.id} complete: {result}")
```

Verifying after each step catches divergence early, before errors compound into
cascading failures.

### Step 5: Verify end state matches original intent

After all steps complete, compare the final system state against the original goal —
not just against whether the last step ran without error:

```
Goal: "Remove events older than 90 days"
End verification:
- Row count before: 4,820,311
- Rows deleted: 3,201,455
- Row count after: 1,618,856
- Oldest remaining event date: within 90 days ✅
- Vacuum complete ✅
- No errors in execution log ✅
```

A plan completing without error is not evidence the goal was achieved. Verify the state.

### Step 6: Prefer reversible actions; stage destructive operations

When two approaches achieve the same goal, prefer the reversible one:

| Destructive | Reversible alternative |
|-------------|----------------------|
| `DELETE FROM table WHERE ...` | Archive to separate table, then delete |
| Overwrite file in place | Write to `.tmp`, swap atomically on success |
| Drop database column | Rename to `_deprecated_col`, drop after validation period |

Stage then commit. The cost is one extra step; the benefit is recovery if something
goes wrong.

## When NOT to Use

- **Single-step tasks** — a task with exactly one tool call doesn't need a plan.
  Adding planning overhead to "look up this user's account balance" is waste.
- **High-frequency inner loops** — agents executing thousands of actions per second
  (simulation, game AI) cannot checkpoint every step architecturally. Use batch
  verification at interval boundaries instead.

## Common Mistakes

**Planning internally without externalizing the plan.** An agent that "reasons" in its
context but never writes the plan as structured output cannot be validated before
execution. The plan must be readable.

**Skipping checkpoints for speed.** "The user wants this done quickly" is not a reason
to bypass checkpoints on irreversible actions. Recovering from an unintended delete
takes far longer than a 3-second confirmation pause.

**Verifying only the final step.** If step 3 silently fails and step 4 continues, the
final state check may pass while the system is actually in a broken state. Verify
after every step.

**Treating plan completion as goal completion.** The plan is a means, not the end.
Always verify the real-world outcome against the original intent.
