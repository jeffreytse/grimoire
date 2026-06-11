---
name: design-security-logging
description: Use when adding logging to authentication flows, authorization decisions, admin actions, or any security-sensitive operation — to enable detection, forensics, and incident response.
source: 'OWASP Logging Cheat Sheet (owasp.org/www-project-cheat-sheets); OWASP Top 10 2021 A09; NIST SP 800-92; CWE-778; PCI DSS v4.0 Requirement 10'
tags: [security, owasp, logging, monitoring, audit-trail, incident-response, developer, devops]
---

# Design Security Logging

Log security-relevant events with sufficient context for forensics while preventing log injection and avoiding logging sensitive data — enabling detection and investigation of attacks, breaches, and insider threats.

## Why This Is Best Practice

**Adopted by:** OWASP Top 10 2021 A09 (Security Logging and Monitoring Failures) is its own category. PCI DSS v4.0 Requirement 10 mandates audit logging for cardholder data environments with 12-month retention. NIST SP 800-92 (Guide to Computer Security Log Management) is the authoritative federal reference. SOC 2 Type II, ISO 27001, and HIPAA all require audit logging with tamper evidence. AWS, GCP, and Azure provide CloudTrail/Cloud Audit Logs/Activity Logs as core compliance features.
**Impact:** The Equifax breach (2017, 147M records) went undetected for 76 days partly due to insufficient logging — OWASP A09 was directly cited in post-mortem analysis. Capital One breach detection (2019) was triggered by an AWS CloudTrail anomaly alert, limiting the window. NIST estimates proper logging reduces mean time to detect (MTTD) breaches by 63 days on average (IBM Cost of a Data Breach Report 2023).
**Why best:** `print()` or unstructured log lines are the alternative — they're unsearchable at scale, don't aggregate, and lack the structured fields needed for SIEM correlation. Structured JSON logs with consistent field names enable automated alerting, SIEM ingestion (Splunk, DataDog, Elastic), and cross-service correlation of attack chains.

Sources: OWASP Logging Cheat Sheet; NIST SP 800-92; IBM Cost of Data Breach 2023; PCI DSS v4.0 Req 10; CWE-778

## Steps

1. **Log these security events at minimum**:

   | Event | Log level | Fields required |
   |-------|-----------|----------------|
   | Login success | INFO | user_id, ip, user_agent, timestamp |
   | Login failure | WARN | username_attempted, ip, reason, timestamp |
   | Logout | INFO | user_id, session_id, timestamp |
   | Password change | INFO | user_id, ip, changed_by, timestamp |
   | Permission change | WARN | user_id, changed_by, old_role, new_role, timestamp |
   | Admin action | WARN | admin_id, action, target, ip, timestamp |
   | Access denied (403) | WARN | user_id, resource, ip, timestamp |
   | Input validation failure | INFO | ip, endpoint, field, reason |

2. **Use structured logging with consistent fields**:

   ```python
   import structlog
   import logging

   logger = structlog.get_logger()

   def log_auth_event(event, user_id, ip, **kwargs):
       logger.info(
           "auth_event",
           event=event,
           user_id=user_id,
           ip=ip,
           timestamp=datetime.utcnow().isoformat() + 'Z',
           **kwargs
       )

   # Usage
   log_auth_event("login_success", user_id=42, ip=request.remote_addr)
   log_auth_event("login_failure", user_id=None, ip=request.remote_addr,
                  username_attempted=username[:50], reason="invalid_password")
   ```

3. **Never log sensitive data** — these must never appear in logs:

   - Passwords, password hashes
   - Full credit card numbers (PAN) — PCI DSS prohibits this
   - Session tokens, API keys, JWTs
   - Full SSN, full account numbers
   - Health data (HIPAA)

   ```python
   # BAD
   logger.info("login", password=request.form['password'])

   # GOOD — log outcome, not credential
   logger.info("login", user_id=user.id, success=True)

   # Mask sensitive fields if partial logging needed
   def mask(value, visible=4):
       return '*' * (len(value) - visible) + value[-visible:]
   logger.info("card", last4=mask(card_number))
   ```

4. **Prevent log injection** — sanitize user-controlled strings before including in log messages:

   ```python
   import re

   def sanitize_for_log(value, max_length=200):
       if not isinstance(value, str):
           value = str(value)
       # Remove newlines (log injection vector — splits one log entry into two)
       value = re.sub(r'[\r\n]', ' ', value)
       return value[:max_length]

   logger.info("user_action", username=sanitize_for_log(request.form['username']))
   ```

5. **Send logs to a centralized, tamper-resistant store** — logs on the same server as the app can be deleted by an attacker who compromises the server:

   - Cloud: CloudWatch (AWS), Cloud Logging (GCP), Azure Monitor
   - Self-hosted: Elastic (ELK), Splunk, Loki + Grafana
   - Configure log forwarding immediately on app start; don't rely on file-based log rotation for security logs

6. **Implement alerting for anomalous patterns**:

   ```yaml
   # Example: alert on brute-force (>5 login failures from one IP in 5 minutes)
   alert: brute_force_attempt
   condition: count(event=login_failure, ip=same) > 5 within 5m
   action: [notify_security_team, temp_block_ip]
   ```

   Key alerts: brute force, privilege escalation, mass data access, off-hours admin activity, access from new geolocation.

7. **Set retention policies** — PCI DSS requires 12 months (3 months immediately accessible); HIPAA requires 6 years; adjust per your compliance requirements. Archive to cold storage after the active window.

## Rules

- Log the who, what, when, where, and outcome — missing any one makes forensics incomplete.
- Use UTC timestamps with millisecond precision and include timezone offset — local time creates ambiguity during DST.
- Log IDs, not names — `user_id=42` is better than `username=alice` for joining with DB records during investigation.
- Failure to log is a security control failure — treat logging errors as security incidents.

## Common Mistakes

- **Logging full request bodies** — may contain passwords, tokens, or PII.
- **Using `print()` or unstructured strings** — unsearchable, doesn't aggregate, no structured fields for SIEM.
- **Only logging errors, not security events** — a successful brute force that guesses correctly generates no error; login success must also be logged.
- **Storing logs on the same host as the application** — an attacker who compromises the app can erase the evidence.
