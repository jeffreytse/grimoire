---
name: apply-llm-output-sanitization
description: Use when rendering LLM-generated text in a browser, executing LLM-generated code, or inserting LLM output into a database or downstream system — any place LLM output leaves the AI layer and enters another execution context.
source: 'OWASP Top 10 for LLM Applications 2025 LLM02 (owasp.org/www-project-top-10-for-large-language-model-applications/); OWASP XSS Prevention Cheat Sheet; CWE-79; CWE-94'
tags: [security, owasp, llm, output-handling, xss, code-injection, ai-security, emerging, developer]
emerging: true
---

# Apply LLM Output Sanitization

Treat LLM-generated content as untrusted input to downstream systems — encoding for context before rendering HTML, sandboxing generated code before execution, and validating structure before inserting into databases.

## Why This Is Best Practice

**Adopted by:** OWASP Top 10 for LLM Applications 2025 LLM02 (Insecure Output Handling) is a dedicated category. Microsoft's AI Red Team and Google DeepMind's safety teams publish guidance on output handling. NIST AI RMF Govern 1.3 includes output validation controls. Production AI systems at Anthropic, OpenAI, and major cloud providers all implement output sandboxing for code execution.
**Status:** Emerging — the attack class has been well-documented since 2023, but standardized defense tooling is still maturing.
**Impact:** If an LLM generates HTML containing `<script>alert('xss')</script>` and it's rendered without escaping, XSS occurs. If an LLM generates shell commands and they're executed without sandboxing, RCE occurs. Bing Chat (2023) was demonstrated generating JavaScript that exfiltrated conversation history. ChatGPT plugins were shown capable of triggering SQL injection through LLM-generated database queries. The LLM is the attack vector; the browser/shell/database is the target.
**Why best:** Trusting LLM outputs because they came from your own model is the common approach — it ignores that LLMs can be manipulated via prompt injection (LLM01) to generate malicious content. Context-aware output encoding (the same defense as for user input XSS) prevents the malicious output from executing, even if generated.

Sources: OWASP LLM Top 10 2025 LLM02; OWASP XSS Prevention Cheat Sheet; CWE-79; Microsoft AI Red Team research

## Steps

1. **HTML-encode LLM output before rendering in a browser**:

   ```python
   import html
   from markupsafe import Markup, escape

   def render_ai_response(llm_output: str) -> str:
       # Escape HTML entities — treat LLM output as user input
       return str(escape(llm_output))

   # In templates (Jinja2, Django templates) — always use auto-escaping
   # {{ ai_response }}  ← auto-escaped (safe)
   # {{ ai_response|safe }}  ← NEVER use |safe on LLM output
   ```

2. **For Markdown rendering: sanitize after converting to HTML**:

   ```javascript
   import DOMPurify from 'dompurify';
   import { marked } from 'marked';

   function renderAIMarkdown(markdownOutput) {
     const rawHtml = marked.parse(markdownOutput);
     // Sanitize after markdown conversion — LLM may inject HTML in markdown
     return DOMPurify.sanitize(rawHtml, {
       ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'ul', 'ol', 'li', 'code', 'pre', 'h1', 'h2', 'h3'],
       ALLOWED_ATTR: [],  // no attributes — prevents href injection
     });
   }
   ```

3. **Sandbox LLM-generated code before execution**:

   ```python
   # Never exec() LLM output directly
   # BAD
   exec(llm_generated_code)

   # GOOD — use sandboxed execution environment
   # Option 1: Docker container with resource limits
   import subprocess
   result = subprocess.run(
       ['docker', 'run', '--rm', '--network=none',
        '--memory=128m', '--cpus=0.5',
        '--security-opt=no-new-privileges',
        'python:3.11-slim', 'python', '-c', llm_generated_code],
       capture_output=True, timeout=10
   )

   # Option 2: Restricted Python interpreter (RestrictedPython)
   from RestrictedPython import compile_restricted, safe_globals
   compiled = compile_restricted(llm_generated_code)
   exec(compiled, safe_globals)
   ```

4. **Validate LLM-generated SQL before execution**:

   ```python
   import sqlparse
   from sqlparse.sql import Statement
   from sqlparse.tokens import Keyword, DML

   ALLOWED_STATEMENT_TYPES = {'SELECT'}

   def validate_llm_sql(sql_string):
       parsed = sqlparse.parse(sql_string)
       for statement in parsed:
           stmt_type = statement.get_type()
           if stmt_type not in ALLOWED_STATEMENT_TYPES:
               raise SecurityError(f"LLM generated disallowed SQL type: {stmt_type}")
       # Also run through parameterization if any values are present
       return sql_string

   # Never pass LLM-generated SQL directly to execute() — use read-only DB user too
   ```

5. **For LLM-generated URLs: validate before rendering as links**:

   ```python
   from urllib.parse import urlparse

   def safe_llm_url(url_string):
       parsed = urlparse(url_string)
       if parsed.scheme not in ('http', 'https'):
           return '#'  # reject javascript:, data:, etc.
       return url_string
   ```

6. **Log LLM outputs for audit and anomaly detection**:

   ```python
   def process_llm_output(raw_output, user_id, session_id):
       # Log before sanitization — for forensics
       logger.info("llm_output", extra={
           'session_id': session_id,
           'output_hash': hashlib.sha256(raw_output.encode()).hexdigest(),
           'output_length': len(raw_output),
       })
       # Sanitize for context
       return sanitize_for_rendering(raw_output)
   ```

## Rules

- Every context where LLM output is consumed has its own encoding requirement — the same output needs HTML encoding for browsers, shell escaping for CLI, and parameterization for SQL.
- `|safe` or `innerHTML =` applied to LLM output requires DOMPurify sanitization first — no exceptions.
- Generated code execution requires sandboxing even when the prompt explicitly requested "safe" code — the LLM may have been injected.
- Structured output (JSON mode) reduces but does not eliminate the risk — a JSON string value can still contain XSS payloads.

## Common Mistakes

- **Trusting LLM output because it's your own model** — even your model can be manipulated via prompt injection in user-supplied context.
- **Rendering LLM markdown as raw HTML without sanitization** — markdown `[click here](javascript:evil())` becomes an XSS.
- **Using `eval()` for LLM-generated expressions** — even "math expressions" can escape the intended scope.
- **Sanitizing output once but using it in multiple contexts** — output sanitized for HTML may still be unsafe for SQL or shell.
