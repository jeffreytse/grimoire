---
name: apply-evaluation-driven-development
description: Use when building or iterating on an LLM feature, prompt, or agentic system — to measure improvement objectively rather than speculating about what the model needs.
source: Anthropic, "Building effective agents" (2024); OpenAI Evals framework documentation; Liang et al., "HELM" Stanford CRFM (2022); Anthropic, Claude agent skills best practices (2024)
tags: [evals, llm, ai-engineering, iteration, baseline, regression, developer, measurement]
related: [write-eval-suite, design-prompt-template]
---

# Apply Evaluation-Driven Development

Build evals before writing extensive prompts. Measure first, then iterate.

## Why This Is Best Practice

**Adopted by:** Anthropic (documented practice for all production AI features), OpenAI
(eval-first methodology across all model iterations — public Evals framework), Google
DeepMind (HELM and BIG-Bench are eval-first by design), and Meta AI, Microsoft Research,
and DeepMind across production AI systems.
**Impact:** Teams iterating without evals consistently over-engineer prompts for scenarios
they haven't tested, adding complexity that doesn't improve measurable performance.
Anthropic's documented practice shows that baseline evals before any prompt engineering
reveal that minimal prompts often perform within 10–15% of heavily engineered versions —
eliminating weeks of speculative work. OpenAI's regression evals on GPT-4 Turbo caught
capability regressions not visible to internal red-teaming (OpenAI Evals framework
documentation).
**Why best:** The alternative — writing detailed prompts based on intuition, then testing
— produces unmeasurable progress. You cannot know if a change helped without a baseline.
Evaluation-driven development applies the same discipline as test-driven development:
define success criteria before building, measure against them, iterate to passing.

Sources: Anthropic, "Building effective agents" (2024); OpenAI Evals framework docs;
Liang et al., "HELM: Holistic Evaluation of Language Models", Stanford CRFM (2022);
Anthropic, Claude agent skills best practices (2024)

## Steps

### Step 1: Write evals before writing extensive prompts

Before adding detailed instructions, few-shot examples, or chain-of-thought scaffolding,
define at least 10–20 test cases that specify success:

```
input: "Summarize this support ticket in one sentence."
expected: concise factual summary, no filler phrases, identifies core issue
```

You cannot measure improvement without a starting point. Evals written after the fact
are biased toward the prompt already built.

### Step 2: Start with a minimal prompt — let evals reveal gaps

Use the simplest possible instruction that describes the task:

```
# Wrong — speculation about what the model needs before any measurement
You are a helpful assistant. When summarizing tickets, always identify the root cause,
use professional tone, avoid jargon, check for urgency markers, keep under 20 words...

# Right — minimal starting prompt
Summarize the support ticket in one sentence.
```

Run the minimal prompt against your evals. Record the baseline pass rate and failure
categories. Don't add complexity before you know where the model actually fails.

### Step 3: Measure and record the baseline

Run your eval suite against the minimal prompt. Record:
- Overall pass rate (e.g., 62% of test cases pass)
- Top failure categories (what type of input fails most often?)
- Surprising passes (edge cases that already work — don't break them)

Baseline numbers are the only honest signal you have. If you skip this step, every
subsequent change is speculation.

### Step 4: Iterate by fixing the largest failure category first

Each iteration follows one rule: one change, then measure.

1. Identify the single largest failure category from eval results
2. Add the minimum prompt change to address it
3. Re-run evals — did the target category improve? Did anything regress?
4. Accept only if target improved with no net regression

```python
while pass_rate < target:
    top_failure = eval_results.top_failure_category()
    candidate_prompt = minimal_fix_for(top_failure)
    new_results = run_evals(candidate_prompt)
    if new_results.pass_rate > current_pass_rate:
        accept(candidate_prompt)
        current_pass_rate = new_results.pass_rate
```

Changing multiple prompt elements in one iteration makes it impossible to isolate what
helped. One change per iteration.

### Step 5: Stop when evals show diminishing returns — not when the prompt looks complete

A prompt that "looks thorough" is not a success criterion. Stop adding complexity when:
- Pass rate is at or above target
- Each additional instruction improves fewer than 2–3 test cases
- You're writing rules for cases not in your eval set

Prompt complexity added beyond what evals justify is technical debt.

### Step 6: Lock evals in CI — catch regressions automatically

Once your feature passes the target bar, every future change (prompt edit, model
upgrade, context change) must pass the existing eval suite before shipping.

See `write-eval-suite` for how to structure the eval harness and integrate with CI/CD.

## When NOT to Use

- **One-off prompts** — a prompt that runs once with no ongoing quality requirement
  doesn't justify the overhead of a formal eval suite.
- **Pure human judgment tasks** — tasks where "correct" is entirely subjective and
  human raters disagree >30% of the time have no ground truth to eval against.

## Common Mistakes

**Writing extensive prompts first, then building evals.** The baseline is now
unavailable. You've invested effort in an unvalidated direction and can't measure
what you've gained.

**Too few eval cases.** 10 examples produce noisy, unreliable signal. Target 50+
before drawing conclusions about a failure category.

**Changing multiple prompt elements in one iteration.** You can't isolate which
change helped. One change per iteration.

**Removing eval cases that "seem outdated."** Shrinking the eval set masks regressions.
Only add cases; never remove them from an active suite.
