---
name: apply-llm-resource-limits
description: Use when building applications that make LLM API calls, run inference locally, or build agentic systems — to prevent runaway costs, infinite loops, and denial of service via unrestricted model consumption.
source: 'OWASP Top 10 for LLM Applications 2025 LLM04 (owasp.org/www-project-top-10-for-large-language-model-applications/); OpenAI usage policies; Anthropic responsible scaling policy; CWE-770'
tags: [security, owasp, llm, resource-limits, dos-prevention, cost-control, ai-security, emerging, developer]
emerging: true
---

# Apply LLM Resource Limits

Enforce token budgets, per-user quotas, request timeouts, and loop detection in LLM applications — preventing runaway inference costs, agent infinite loops, and user-driven denial of service.

## Why This Is Best Practice

**Adopted by:** OWASP Top 10 for LLM Applications 2025 LLM04 (Model Denial of Service). OpenAI, Anthropic, and Google all enforce per-account and per-minute rate limits on their APIs. AWS Bedrock, Azure OpenAI, and Google Vertex AI all provide quota management controls. The AI engineering community (LangChain, LlamaIndex, AutoGen) all include resource limit configurations in their framework defaults.
**Status:** Emerging — the attack class is well-understood, but defense tooling and best practices are still being standardized in 2024-2025.
**Impact:** A single user submitting a prompt with a 100,000-token context window at $0.01/1K tokens generates $1 per request — at 100 concurrent requests, $100/minute in inference costs. Agentic systems (AutoGPT-style) with no iteration limits have been observed consuming $50–$500 in a single runaway session. Without limits, a single malicious or buggy user can drain a monthly budget in minutes. Recursive tool-calling loops in multi-agent systems can saturate compute indefinitely.
**Why best:** Monitoring costs after the fact is the common approach — it detects abuse after damage occurs. Pre-emptive per-request and per-user limits cap the blast radius at known thresholds.

Sources: OWASP LLM Top 10 2025 LLM04; CWE-770; OpenAI rate limit documentation; Anthropic responsible scaling policy

## Steps

1. **Set per-request token budgets**:

   ```python
   MAX_INPUT_TOKENS = 4096     # max context size per request
   MAX_OUTPUT_TOKENS = 2048    # max generated tokens per response

   def call_llm(prompt: str, system: str = "") -> str:
       # Estimate tokens before calling (approx 4 chars per token)
       estimated_input = len(prompt + system) // 4
       if estimated_input > MAX_INPUT_TOKENS:
           raise ValueError(f"Input too long: ~{estimated_input} tokens")

       response = client.messages.create(
           model="claude-sonnet-4-6",
           max_tokens=MAX_OUTPUT_TOKENS,  # hard cap on output
           system=system,
           messages=[{"role": "user", "content": prompt}]
       )
       return response.content[0].text
   ```

2. **Implement per-user and per-tenant token quotas**:

   ```python
   import redis
   from datetime import datetime

   DAILY_TOKEN_LIMIT = 100_000   # per user per day
   MONTHLY_TOKEN_LIMIT = 2_000_000  # per user per month

   def check_and_consume_quota(user_id: str, estimated_tokens: int):
       day_key = f'tokens:{user_id}:{datetime.utcnow().date()}'
       month_key = f'tokens:{user_id}:{datetime.utcnow().strftime("%Y-%m")}'

       pipe = redis.pipeline()
       pipe.incrby(day_key, estimated_tokens)
       pipe.expire(day_key, 86400 * 2)  # 2-day TTL
       pipe.incrby(month_key, estimated_tokens)
       pipe.expire(month_key, 86400 * 35)
       day_total, _, month_total, _ = pipe.execute()

       if day_total > DAILY_TOKEN_LIMIT:
           raise QuotaExceeded(f"Daily token limit reached. Resets at midnight UTC.")
       if month_total > MONTHLY_TOKEN_LIMIT:
           raise QuotaExceeded(f"Monthly token limit reached.")
   ```

3. **Limit agent iteration counts and recursion depth**:

   ```python
   class SafeAgent:
       MAX_ITERATIONS = 20
       MAX_TOOL_CALLS_PER_ITER = 5

       def run(self, task: str) -> str:
           iteration = 0
           total_tool_calls = 0

           while iteration < self.MAX_ITERATIONS:
               iteration += 1
               response = self.llm.complete(task, tools=self.tools)

               if not response.tool_calls:
                   return response.content  # task complete

               if len(response.tool_calls) > self.MAX_TOOL_CALLS_PER_ITER:
                   raise AgentError("Too many tool calls in single step")

               total_tool_calls += len(response.tool_calls)
               self.execute_tool_calls(response.tool_calls)

           raise AgentError(f"Agent exceeded {self.MAX_ITERATIONS} iterations without completion")
   ```

4. **Set wall-clock timeouts on LLM calls**:

   ```python
   import asyncio

   async def call_llm_with_timeout(prompt: str, timeout_seconds: float = 30.0) -> str:
       try:
           response = await asyncio.wait_for(
               llm_client.acompletions(prompt),
               timeout=timeout_seconds
           )
           return response.text
       except asyncio.TimeoutError:
           logger.warning("LLM call timed out after %ss", timeout_seconds)
           raise ServiceTimeout("AI response took too long")
   ```

5. **Monitor and alert on cost anomalies**:

   ```python
   def track_llm_cost(user_id: str, input_tokens: int, output_tokens: int,
                      model: str):
       # Calculate cost (example rates — check current pricing)
       cost_usd = (input_tokens * 0.003 + output_tokens * 0.015) / 1000

       metrics.increment('llm.tokens.input', input_tokens, tags={'user': user_id})
       metrics.increment('llm.tokens.output', output_tokens, tags={'user': user_id})
       metrics.gauge('llm.cost.request', cost_usd, tags={'user': user_id, 'model': model})

       # Alert if single request exceeds threshold
       if cost_usd > 1.00:  # $1 per request
           alert_ops(f"High-cost LLM request: ${cost_usd:.2f} from user {user_id}")
   ```

6. **Validate and truncate user-supplied context before including in prompts**:

   ```python
   import tiktoken

   def truncate_to_token_limit(text: str, max_tokens: int, model: str = "gpt-4") -> str:
       enc = tiktoken.encoding_for_model(model)
       tokens = enc.encode(text)
       if len(tokens) <= max_tokens:
           return text
       # Truncate and add indicator
       truncated = enc.decode(tokens[:max_tokens])
       return truncated + "\n[Content truncated due to length]"
   ```

## Rules

- Set `max_tokens` on every LLM API call — never rely on default limits.
- Agent loop limits must be enforced in code, not just in the system prompt ("stop after 5 steps" in a prompt is not a hard limit).
- Cost limits should fail-open (allow the request, log the anomaly) at low thresholds and fail-closed at absolute limits.
- Monitor token consumption by user, tenant, and endpoint — not just total usage.

## Common Mistakes

- **No iteration limit on agentic loops** — the most common source of runaway costs in production AI systems.
- **Trusting user-supplied document sizes** — users can submit 1MB text files claiming to be 1KB; always measure actual token count.
- **Setting limits only at the API gateway** — application-level limits are needed too for per-user quota enforcement.
- **Not monitoring until the month-end invoice** — real-time alerting on anomalous usage is required; invoices arrive too late.
