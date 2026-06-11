---
name: design-llm-action-boundaries
description: Use when building LLM agents that can take actions in the world — sending emails, modifying files, making API calls, or executing code — to prevent the model from taking actions beyond its intended scope.
source: 'OWASP Top 10 for LLM Applications 2025 LLM08 (owasp.org/www-project-top-10-for-large-language-model-applications/); NIST AI RMF 1.0 Govern 1.3; EU AI Act Article 14; CWE-284'
tags: [security, owasp, llm, agents, excessive-agency, action-boundaries, ai-safety, emerging, developer]
emerging: true
---

# Design LLM Action Boundaries

Constrain LLM agents to the minimum necessary capabilities, enforce human-in-the-loop for irreversible actions, and implement hard permission limits that cannot be overridden by prompt — preventing agents from acting outside their intended scope.

## Why This Is Best Practice

**Adopted by:** OWASP Top 10 for LLM Applications 2025 LLM08 (Excessive Agency). NIST AI RMF 1.0 (2023) Govern 1.3 requires human oversight for autonomous AI systems. EU AI Act Article 14 (2024) mandates human oversight mechanisms for high-risk AI. Anthropic's Model Specification, OpenAI's GPT policies, and Google's Responsible AI Practices all include capability limitation guidance for agentic systems.
**Status:** Emerging — agentic AI (AI taking actions autonomously) became mainstream in 2023; safety standards are still maturing. This is one of the most rapidly evolving areas in AI safety.
**Impact:** AutoGPT and BabyAGI (2023) were demonstrated deleting files, sending emails, and making API calls based on prompt injections from web pages. AI coding assistants have been demonstrated installing malicious packages via injected code comments. An AI customer service agent with CRM write access and no approval requirement was manipulated into issuing fraudulent refunds. The more actions an agent can take, the higher the potential blast radius of a successful injection.
**Why best:** Reactive monitoring (detecting unauthorized actions after they occur) is the alternative — it cannot undo irreversible actions like sent emails, deleted data, or executed financial transactions. Proactive capability limitation prevents the actions from occurring in the first place.

Sources: OWASP LLM Top 10 2025 LLM08; NIST AI RMF 1.0 Govern 1.3; EU AI Act Article 14; AutoGPT security research (2023)

## Steps

1. **Apply strict least-privilege to agent capabilities**:

   ```python
   from enum import Enum

   class AgentCapability(Enum):
       READ_DOCS = "read_docs"
       WRITE_DOCS = "write_docs"
       SEND_EMAIL = "send_email"
       EXECUTE_CODE = "execute_code"
       CALL_EXTERNAL_API = "call_external_api"

   # Task-specific capability grants — grant only what the task needs
   TASK_CAPABILITIES = {
       'document_summary': {AgentCapability.READ_DOCS},
       'draft_email': {AgentCapability.READ_DOCS},  # NOT send_email
       'send_approved_email': {AgentCapability.SEND_EMAIL},
       'code_review': {AgentCapability.READ_DOCS, AgentCapability.EXECUTE_CODE},
   }

   class BoundedAgent:
       def __init__(self, task_type: str, user_id: str):
           self.capabilities = TASK_CAPABILITIES.get(task_type, set())
           self.user_id = user_id

       def can_do(self, capability: AgentCapability) -> bool:
           return capability in self.capabilities
   ```

2. **Separate planning from execution — require approval for action plans**:

   ```python
   def two_phase_execution(agent, task: str):
       # Phase 1: Generate action plan (read-only)
       plan = agent.plan(task)  # LLM generates list of proposed actions

       # Phase 2: Human reviews and approves plan
       approved_actions = present_plan_for_approval(plan)

       if approved_actions is None:
           return {'status': 'rejected', 'reason': 'User rejected plan'}

       # Phase 3: Execute only approved actions
       results = []
       for action in approved_actions:
           result = agent.execute_single_action(action)
           results.append(result)

       return {'status': 'completed', 'results': results}
   ```

3. **Hard-code irreversible action gates in non-LLM code**:

   ```python
   class ActionGate:
       """Hard gates that cannot be bypassed by prompt manipulation."""

       IRREVERSIBLE_ACTIONS = {
           'delete_file', 'send_email', 'make_payment',
           'drop_database', 'terminate_instance'
       }

       @staticmethod
       def require_confirmation(action: str, params: dict, user_id: str) -> bool:
           if action not in ActionGate.IRREVERSIBLE_ACTIONS:
               return True  # no confirmation needed

           # Confirmation must come from user, not from LLM
           confirmation_id = create_pending_confirmation(action, params, user_id)
           # Send to user via out-of-band channel (SMS, push notification)
           send_confirmation_request(user_id, confirmation_id, action, params)

           # Wait for explicit confirmation (max 5 minutes)
           return wait_for_confirmation(confirmation_id, timeout=300)
   ```

4. **Scope filesystem and network access to a sandbox**:

   ```python
   import os
   import subprocess

   AGENT_SANDBOX_DIR = f'/tmp/agent_sandbox/{session_id}'

   def execute_agent_code(code: str, sandbox_dir: str) -> dict:
       os.makedirs(sandbox_dir, exist_ok=True)
       result = subprocess.run(
           ['docker', 'run', '--rm',
            f'--volume={sandbox_dir}:/workspace:rw',  # limited filesystem access
            '--network=none',                            # no network access
            '--memory=256m', '--cpus=0.5',
            '--security-opt=no-new-privileges',
            '--read-only',
            'python:3.11-slim',
            'python', '-c', code],
           capture_output=True,
           timeout=30
       )
       return {'stdout': result.stdout.decode(), 'returncode': result.returncode}
   ```

5. **Implement scope boundaries on data access** — agents should not access data beyond the current task:

   ```python
   class ScopedDataAccess:
       def __init__(self, user_id: str, session_id: str, allowed_doc_ids: list):
           self.user_id = user_id
           self.session_id = session_id
           self.allowed_doc_ids = set(allowed_doc_ids)  # explicit list for this task

       def read_document(self, doc_id: str) -> str:
           if doc_id not in self.allowed_doc_ids:
               raise PermissionError(f"Agent not authorized to read {doc_id} in this session")
           return db.get_document(doc_id, owner_id=self.user_id)  # ownership enforced too
   ```

6. **Log all agent decisions and actions for audit**:

   ```python
   def log_agent_action(session_id, action, params, approved_by, outcome):
       logger.info("agent_action", extra={
           'session_id': session_id,
           'action': action,
           'approved_by': approved_by,  # 'auto' or user_id
           'outcome': outcome,
           'timestamp': datetime.utcnow().isoformat(),
       })
   ```

## Rules

- The LLM should plan; deterministic code should execute — never let the LLM execute actions directly without a typed execution layer.
- Irreversible action gates must be implemented outside the LLM's influence — hard-coded in application logic, not enforced via system prompt.
- An agent's capabilities should be tied to the specific task, not the user's maximum permissions — a document-summarization task doesn't need email sending capability even if the user has that permission.
- The more autonomous the agent, the more restrictive the capability set should be — fully autonomous agents need the tightest constraints.

## Common Mistakes

- **Giving agents the same permissions as the user** — users have permissions for interactive tasks; agents acting autonomously should have a subset.
- **Using the system prompt as the only capability gate** — "only read files, never delete" in a system prompt is bypassed by injection; hard code the gate.
- **No human-in-the-loop for multi-step agentic tasks** — each step can compound errors or malicious actions; checkpoints limit blast radius.
- **Logging only errors, not successful actions** — forensics requires knowing what the agent did, not just when it failed.
