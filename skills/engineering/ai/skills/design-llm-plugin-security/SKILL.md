---
name: design-llm-plugin-security
description: Use when building tools, plugins, or function-calling integrations for LLM systems — defining what actions the model can invoke, how those invocations are validated, and how to prevent misuse of plugin capabilities.
source: 'OWASP Top 10 for LLM Applications 2025 LLM07 (owasp.org/www-project-top-10-for-large-language-model-applications/); NIST AI RMF 1.0 Govern 1.3; Anthropic tool use documentation; OpenAI function calling documentation'
tags: [security, owasp, llm, plugins, tool-use, function-calling, ai-security, emerging, developer]
emerging: true
---

# Design LLM Plugin Security

Build LLM tool and plugin interfaces with explicit input validation, minimal permissions, human confirmation for consequential actions, and sandboxed execution — preventing prompt injection from triggering unauthorized tool invocations.

## Why This Is Best Practice

**Adopted by:** OWASP Top 10 for LLM Applications 2025 LLM07 (Insecure Plugin Design). OpenAI, Anthropic, and Google's function calling documentation all include security guidance for tool definitions. The Model Context Protocol (MCP) specification from Anthropic includes authorization requirements for tool invocations. NIST AI RMF 1.0 Govern 1.3 requires human oversight mechanisms for automated AI systems.
**Status:** Emerging — function calling / tool use became mainstream in 2023; security standards are still being developed.
**Impact:** LLM plugins that invoke APIs, execute code, or modify data are the primary escalation path from prompt injection (information disclosure) to action execution (data deletion, unauthorized transactions, system compromise). Demonstrated attacks: manipulated AI email assistants sent emails to contacts, AI code assistants triggered malicious package installs via injected code comments, AI financial assistants initiated unauthorized transfers via forged instructions in documents.
**Why best:** Trusting the LLM to only call tools when appropriate is the common approach — it fails whenever the LLM is manipulated via prompt injection. Explicit authorization checks, input validation, and confirmation for high-impact actions provide defense-in-depth that holds even when the LLM is injected.

Sources: OWASP LLM Top 10 2025 LLM07; NIST AI RMF 1.0; MCP specification; OpenAI function calling security guide

## Steps

1. **Define strict, narrow tool schemas** — minimize what each tool can accept:

   ```python
   tools = [
       {
           "name": "search_documents",
           "description": "Search the user's own documents. Only returns documents owned by the authenticated user.",
           "input_schema": {
               "type": "object",
               "properties": {
                   "query": {
                       "type": "string",
                       "maxLength": 500,  # hard limit
                       "description": "Search query"
                   }
               },
               "required": ["query"],
               "additionalProperties": False  # reject unexpected params
           }
       }
   ]
   ```

2. **Validate all tool inputs before execution** — never pass LLM-generated parameters directly to tools:

   ```python
   from pydantic import BaseModel, validator, constr

   class SearchDocumentsInput(BaseModel):
       query: constr(min_length=1, max_length=500)

       @validator('query')
       def no_injection_patterns(cls, v):
           # Reject obvious injection patterns
           suspicious = ['ignore previous', 'system prompt', 'act as root']
           for pattern in suspicious:
               if pattern.lower() in v.lower():
                   raise ValueError(f"Suspicious query rejected")
           return v

   def execute_tool(tool_name: str, raw_params: dict, user_id: str):
       # Parse and validate against schema — never use raw_params directly
       if tool_name == 'search_documents':
           params = SearchDocumentsInput(**raw_params)
           return search_documents(params.query, owner_id=user_id)  # enforce ownership
       raise ValueError(f"Unknown tool: {tool_name}")
   ```

3. **Enforce user context in every tool call** — tools must operate within the caller's permissions:

   ```python
   class ToolExecutor:
       def __init__(self, user_id: str, user_permissions: set):
           self.user_id = user_id
           self.user_permissions = user_permissions

       def execute(self, tool_name: str, params: dict) -> dict:
           required_permission = TOOL_PERMISSIONS[tool_name]
           if required_permission not in self.user_permissions:
               raise PermissionError(f"Tool {tool_name} requires {required_permission}")

           # All tools receive user_id — they cannot access other users' data
           return TOOL_HANDLERS[tool_name](params, user_id=self.user_id)

   TOOL_PERMISSIONS = {
       'search_documents': 'documents:read',
       'delete_document': 'documents:delete',
       'send_email': 'email:send',
   }
   ```

4. **Require human confirmation for consequential or irreversible actions**:

   ```python
   HIGH_IMPACT_TOOLS = {'delete_document', 'send_email', 'make_payment', 'create_user'}

   def execute_tool_with_confirmation(tool_name: str, params: dict,
                                      user_id: str, session_id: str):
       if tool_name in HIGH_IMPACT_TOOLS:
           # Store pending action for user approval
           action_id = store_pending_action(tool_name, params, user_id, session_id)
           return {
               'status': 'awaiting_confirmation',
               'action_id': action_id,
               'description': describe_action(tool_name, params),
               'confirm_url': f'/confirm/{action_id}',
           }
       return execute_tool(tool_name, params, user_id)
   ```

5. **Log all tool invocations with full context**:

   ```python
   def log_tool_call(tool_name: str, params: dict, user_id: str,
                     result: dict, triggered_by: str):
       logger.info("tool_invocation", extra={
           'tool': tool_name,
           'user_id': user_id,
           'params_hash': hashlib.sha256(json.dumps(params, sort_keys=True).encode()).hexdigest(),
           'triggered_by': triggered_by,  # 'user_direct' or 'llm_agent'
           'success': 'error' not in result,
       })
   ```

6. **Limit tool exposure in the context** — only expose tools relevant to the current task:

   ```python
   def get_relevant_tools(task_type: str) -> list:
       TOOL_SETS = {
           'document_qa': ['search_documents', 'get_document_excerpt'],
           'email_draft': ['search_documents', 'get_email_template'],
           # NOT: send_email, delete_document — not needed for drafting
       }
       return [tools[t] for t in TOOL_SETS.get(task_type, [])]
   ```

## Rules

- Never pass raw LLM tool parameters to backend functions — always parse through a typed schema first.
- Tool descriptions in the schema influence what the LLM calls them for — be precise: "Search the authenticated user's documents" not "Search documents".
- Sandboxed code execution tools (Python interpreter, shell) are the highest-risk tool type — require the most restrictive containment (no network, no filesystem outside /tmp, resource limits).
- Plugin chains (tool A calls tool B) multiply the attack surface — audit each tool independently and require authorization at each hop.

## Common Mistakes

- **Exposing all tools for all tasks** — the LLM sees what it can call; fewer tools = smaller injection target.
- **Not validating that tool output is within the user's access scope** — a search tool may return documents from other users if ownership filtering is in the DB query, not the tool schema.
- **Treating tool call parameters as safe because they came from the LLM** — they came from user-influenced context; treat them as untrusted.
- **No audit log for tool invocations** — agentic systems need forensic trails; "the agent did it" is not sufficient incident response.
