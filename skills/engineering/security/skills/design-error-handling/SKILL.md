---
name: design-error-handling
description: Use when designing how errors, exceptions, and failures are surfaced to users and logged internally — ensuring error messages don't leak sensitive implementation details while providing enough context for debugging.
source: 'OWASP Error Handling Cheat Sheet (owasp.org/www-project-cheat-sheets); OWASP Top 10 2021 A05; CWE-209; NIST SP 800-53 SI-11'
tags: [security, owasp, error-handling, information-disclosure, logging, developer]
---

# Design Error Handling

Return generic error messages to clients while logging detailed errors server-side — preventing stack traces, database schemas, and internal paths from leaking to attackers.

## Why This Is Best Practice

**Adopted by:** OWASP Top 10 2021 A05 (Security Misconfiguration) explicitly calls out verbose error messages as a configuration failure. CWE-209 (Generation of Error Message Containing Sensitive Information) is tracked across hundreds of CVEs annually. NIST SP 800-53 SI-11 mandates error handling that reveals minimal information. Django's `DEBUG = False` requirement, Rails' `config.consider_all_requests_local = false`, and Spring's production error configuration all implement this principle.
**Impact:** Error messages containing stack traces reveal framework versions (enabling targeted CVE exploitation), internal file paths (useful for path traversal), database table names and column names (enabling refined SQL injection), and internal IP addresses. The 2014 Heartbleed exploit was partly facilitated by verbose error output during discovery. Shodan scans routinely find live debug error pages that reveal full stack traces including database credentials in query parameters.
**Why best:** Catching and swallowing all exceptions silently (the alternative) creates debugging nightmares and hides operational issues. The correct design separates external-facing error messages (generic, safe) from internal logging (detailed, not user-facing) — giving both security and debuggability.

Sources: OWASP Error Handling Cheat Sheet; CWE-209; NIST SP 800-53 SI-11; Django/Rails production configuration guides

## Steps

1. **Define safe, generic error responses for clients**:

   ```python
   # Map exception types to safe user-facing messages
   ERROR_MESSAGES = {
       'not_found':      'The requested resource was not found.',
       'unauthorized':   'Authentication required.',
       'forbidden':      'You do not have permission to perform this action.',
       'validation':     'The request contains invalid data.',
       'server_error':   'An unexpected error occurred. Please try again.',
       'rate_limited':   'Too many requests. Please wait before retrying.',
   }

   @app.errorhandler(Exception)
   def handle_exception(e):
       error_id = generate_error_id()  # UUID for log correlation
       logger.exception("Unhandled exception [%s]", error_id, exc_info=e)

       if isinstance(e, NotFound):
           return jsonify({'error': ERROR_MESSAGES['not_found'],
                          'error_id': error_id}), 404
       if isinstance(e, Unauthorized):
           return jsonify({'error': ERROR_MESSAGES['unauthorized'],
                          'error_id': error_id}), 401
       # Default: generic 500
       return jsonify({'error': ERROR_MESSAGES['server_error'],
                      'error_id': error_id}), 500
   ```

   The `error_id` lets users report the ID to support, who can correlate to full server logs — without exposing the details publicly.

2. **Never return stack traces, SQL errors, or internal paths to clients**:

   ```python
   # BAD — exposes DB schema, internal path, Python version
   try:
       user = db.query(f"SELECT * FROM users WHERE id = {user_id}")
   except Exception as e:
       return jsonify({'error': str(e)}), 500
   # Returns: "psycopg2.errors.UndefinedColumn: column "admin_flag" of relation
   #          "users" does not exist\nLINE 1: SELECT admin_flag FROM users..."

   # GOOD — internal detail logged, generic message returned
   try:
       user = db.query("SELECT * FROM users WHERE id = %s", (user_id,))
   except DatabaseError as e:
       error_id = log_error(e)
       return jsonify({'error': 'Database error', 'error_id': error_id}), 500
   ```

3. **Disable debug error pages in production**:

   ```python
   # Flask
   app.config['PROPAGATE_EXCEPTIONS'] = False
   app.config['DEBUG'] = False  # never True in production

   # Django
   DEBUG = False  # settings.py — also set ALLOWED_HOSTS
   ```

   ```javascript
   // Express — never use this middleware in production
   if (process.env.NODE_ENV !== 'production') {
       app.use(require('errorhandler')());
   }

   // Production error handler — no stack trace to client
   app.use((err, req, res, next) => {
       const errorId = crypto.randomUUID();
       logger.error({ err, errorId }, 'Unhandled error');
       res.status(500).json({ error: 'Internal server error', errorId });
   });
   ```

4. **Log full exception details server-side with correlation ID**:

   ```python
   import uuid
   import traceback

   def log_error(exc, context=None):
       error_id = str(uuid.uuid4())
       logger.error(
           "error [%s]: %s\n%s",
           error_id,
           str(exc),
           traceback.format_exc(),
           extra={'context': context}
       )
       return error_id
   ```

5. **Handle specific exception types explicitly** — don't catch-all and swallow:

   ```python
   # BAD — swallows errors silently
   try:
       process_payment(amount)
   except Exception:
       pass  # user never knows, payment may have failed

   # GOOD — handle specifically, propagate unknown errors
   try:
       process_payment(amount)
   except PaymentDeclinedException as e:
       return jsonify({'error': 'Payment declined', 'code': e.decline_code}), 402
   except PaymentGatewayTimeout:
       return jsonify({'error': 'Payment service unavailable'}), 503
   # Let unknown exceptions propagate to the global handler
   ```

6. **Sanitize validation error messages** — field names and values are usually safe to return; internal system details are not:

   ```python
   # Safe to return: field names, user-supplied values, human-readable constraints
   {
     "error": "Validation failed",
     "fields": {
       "email": "Invalid email format",
       "age": "Must be between 18 and 120"
     }
   }

   # Never return: SQL constraint names, database column types, stack traces
   ```

7. **Set appropriate HTTP status codes** — attackers use status codes to enumerate behavior:

   | Scenario | Correct status |
   |----------|---------------|
   | Resource not found | 404 |
   | Not authenticated | 401 |
   | Authenticated but no permission | 403 |
   | Invalid input | 400 |
   | Rate limited | 429 |
   | Server error | 500 |
   | Auth endpoint: wrong credentials | 401 (not 403 — no enumeration of users vs passwords) |

## Rules

- `error_id` (a UUID) in the client response enables support correlation without exposing internals.
- HTTP 200 with `{"success": false, "error": "..."}` bodies break client error handling — use proper status codes.
- Error messages must not differ between "user not found" and "wrong password" — both return 401 with identical messaging to prevent user enumeration.
- Logging exceptions without re-raising them in middleware swallows errors; use `logger.exception()` which logs the full traceback.

## Common Mistakes

- **`except Exception as e: return str(e)`** — the single most common source of information leakage in web apps.
- **Different error messages for "user not found" vs "wrong password"** — enables user enumeration (attacker learns which usernames exist).
- **Logging PII in error context** — stack traces that include the request body may log passwords, SSNs, or payment data.
- **Catching BaseException** — catches `KeyboardInterrupt`, `SystemExit` — use `Exception` as the catch-all base class.
