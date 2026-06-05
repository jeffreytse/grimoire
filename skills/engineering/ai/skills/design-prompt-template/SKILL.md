---
name: design-prompt-template
description: Use when building LLM-powered features that require consistent, reliable outputs across varied inputs, or when prompt quality is causing inconsistent application behavior
source: OpenAI "Prompt Engineering Guide" best practices; Anthropic model card guidance; Wei et al. "Chain-of-Thought Prompting" NeurIPS (2022)
tags: [ai, llm, prompt-engineering, product]
verified: true
---

# Design Prompt Template

Create structured, reusable prompt templates that produce reliable, consistent LLM outputs across diverse inputs in production systems.

## Why This Is Best Practice

**Adopted by:** Anthropic, OpenAI, Google (all publish prompt engineering guides); every production LLM product uses templated prompts, not ad-hoc strings
**Impact:** Structured prompts reduce output variance by 40-70% (OpenAI cookbook benchmarks); chain-of-thought prompting improves reasoning accuracy by 18-57% on benchmark tasks (Wei et al. 2022)
**Why best:** LLMs are sensitive to exact wording; unstructured prompts produce inconsistent outputs; templates enforce structure, enable testing, and make prompt evolution manageable

Sources: Wei et al. "Chain-of-Thought Prompting Elicits Reasoning in Large Language Models" NeurIPS (2022); OpenAI "Prompt Engineering" docs; Anthropic "Claude's Constitution" guidance

## Steps

1. **Define the task precisely** — Write a one-sentence task definition: "Given X, produce Y in format Z." Vague tasks produce vague outputs. If you cannot state the task in one sentence, decompose it. The task definition becomes the core of your system prompt.

2. **Design the system prompt** — Assign a clear role: "You are a [role] that [does task]." Specify: output format (JSON schema, markdown, numbered list), tone and style, constraints ("respond only in English," "do not speculate"), and what to do when input is ambiguous or out of scope.

3. **Structure the user prompt template** — Use clear delimiters to separate variable content from instructions. XML tags work well for Claude: `<document>{{content}}</document>`. Markdown code fences work for code. Place variable content after instructions to prevent prompt injection via user data.

4. **Apply chain-of-thought for reasoning tasks** — For tasks requiring multi-step reasoning (analysis, math, classification), add: "Think step by step before giving your final answer" or use `<thinking>` tags to elicit reasoning before the answer. Chain-of-thought consistently improves accuracy on complex tasks.

5. **Use few-shot examples for format enforcement** — Include 2-5 representative input-output examples when output format is complex or unusual. Examples are more reliable than format instructions alone. Choose examples that cover edge cases and represent the expected distribution of inputs.

6. **Specify output format explicitly** — If returning structured data: provide a JSON schema example or a template with placeholder values. Request XML or JSON over free prose when downstream code parses the output. Validate output format in your application and retry on parse failure.

7. **Handle edge cases in the prompt** — Explicitly instruct the model on what to do with: ambiguous inputs, out-of-scope requests, insufficient information, and adversarial inputs. "If the input is not a valid X, respond with {error: 'invalid_input', reason: '...'}." Models without explicit edge case handling produce unpredictable outputs on edge inputs.

8. **Parameterize and version the template** — Store prompt templates as versioned files or in a prompt registry (not hardcoded strings). Use a templating system (Jinja2, Handlebars) for variable substitution. Treat prompt changes as code changes requiring review and testing.

9. **Evaluate with a test set** — Build a labeled test set of 50-200 input-output pairs covering normal cases, edge cases, and adversarial inputs. Run evaluations before deploying prompt changes. Define pass/fail criteria (ROUGE score, JSON validity, human rating threshold). Never ship prompt changes without evaluation.

10. **Monitor output quality in production** — Log a sample of production inputs and outputs. Implement output classifiers or LLM-as-judge to detect failures (format violations, refusals, off-topic responses). Alert on quality degradation. Model updates can break prompts without notice.

## Rules

- Never concatenate user input directly into prompts without delimiters — this enables prompt injection; always wrap user content in explicit tags.
- Prompts are code: version them, review them, test them before shipping.
- System prompts should be comprehensive enough that the model never encounters an unspecified situation — all edge cases should be handled explicitly.
- One prompt does one task; compound tasks are better served by chaining single-purpose prompts.

## Common Mistakes

- **Vague role assignment** — "You are a helpful assistant" provides no task context; "You are a JSON extraction agent that extracts named entities from text" is specific and constraining.
- **No output format specification** — models produce varied formats across calls; always specify format and validate it.
- **Ignoring token limits** — prompts with long few-shot examples may leave insufficient context for the actual input; measure prompt token length and reserve budget for input + output.
- **No evaluation set** — shipping prompt changes without evaluation is shipping code without tests; regressions are discovered in production.

## When NOT to Use

- Simple classification tasks where a fine-tuned model is more cost-effective and reliable than prompting a large model
- Tasks where deterministic code (regex, parsing, rule engine) is more reliable than probabilistic LLM output
- High-stakes decisions where model reasoning cannot be audited or explained to affected parties
