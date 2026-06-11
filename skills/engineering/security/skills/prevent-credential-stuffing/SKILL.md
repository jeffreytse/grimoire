---
name: prevent-credential-stuffing
description: Use when building login endpoints — detecting and blocking automated credential stuffing attacks using breached password detection, device fingerprinting, and bot detection.
source: 'OWASP Credential Stuffing Prevention Cheat Sheet (owasp.org/www-project-cheat-sheets); NIST SP 800-63B Section 5.1.1.2; HaveIBeenPwned API documentation; Shape Security research'
tags: [security, owasp, credential-stuffing, authentication, bot-detection, rate-limiting, developer]
---

# Prevent Credential Stuffing

Detect and block credential stuffing attacks by checking submitted passwords against breached credential databases, implementing multi-signal bot detection, and using device fingerprinting — distinguishing automated attacks from legitimate users.

## Why This Is Best Practice

**Adopted by:** OWASP Credential Stuffing Prevention Cheat Sheet (2023) is the primary reference. NIST SP 800-63B Section 5.1.1.2 mandates checking new passwords against lists of commonly used, expected, or compromised values. HaveIBeenPwned (HIBP) API is used by Microsoft, 1Password, Google Chrome, and Mozilla Firefox for real-time breached password checking. Shape Security (now F5) research documents credential stuffing as responsible for 90% of login traffic on major e-commerce sites.
**Impact:** Verizon DBIR 2023 attributes 49% of breaches to stolen credentials — credential stuffing is the primary mechanism for exploiting those credentials at scale. The 2022 Rockstar Games breach used credential stuffing to access employee accounts. The 2021 Coinbase credential stuffing attack compromised 6,000+ accounts by exploiting weak SMS MFA. Shape Security estimates credential stuffing attacks generate 2–3 billion login attempts per day across the internet. A site with 1 million users and 1% password reuse rate has 10,000 accounts vulnerable to stuffing from any major data breach.
**Why best:** IP-based rate limiting alone fails because attackers use botnets with thousands of IPs (e.g., 1 attempt per IP per day is below most limits but allows millions of attempts). Multi-signal detection (failed login velocity, device fingerprint novelty, IP reputation, user-agent analysis) distinguishes automated attacks without degrading legitimate user experience. Breached password detection prevents accounts from using compromised credentials in the first place.

Sources: OWASP Credential Stuffing Prevention Cheat Sheet; NIST SP 800-63B Section 5.1.1.2; Shape Security "State of Credential Stuffing" (2022); HaveIBeenPwned API documentation

## Steps

1. **Check passwords against HaveIBeenPwned on login and registration**:

   ```python
   import hashlib
   import httpx

   def is_password_breached(password: str) -> bool:
       """
       HIBP k-anonymity API: send first 5 chars of SHA-1 hash.
       Server returns all hashes with that prefix. Never sends full hash.
       """
       sha1 = hashlib.sha1(password.encode()).hexdigest().upper()
       prefix = sha1[:5]
       suffix = sha1[5:]

       response = httpx.get(
           f"https://api.pwnedpasswords.com/range/{prefix}",
           timeout=3.0,
           headers={"Add-Padding": "true"},  # prevents traffic analysis
       )
       response.raise_for_status()

       for line in response.text.splitlines():
           hash_suffix, count = line.split(":")
           if hash_suffix == suffix:
               return True  # password found in {count} breaches
       return False

   def check_at_login(email: str, password: str) -> dict:
       user = authenticate(email, password)
       if not user:
           return {"success": False}

       if is_password_breached(password):
           # Don't block login — force immediate password reset
           return {
               "success": True,
               "force_password_reset": True,
               "reason": "This password was found in a data breach. Please change it."
           }
       return {"success": True}
   ```

2. **Implement multi-dimensional rate limiting — IP + account + global**:

   ```python
   import redis
   from datetime import timedelta

   redis_client = redis.Redis()

   class CredentialStuffingDetector:
       LIMITS = {
           "per_ip": (20, timedelta(minutes=10)),       # 20 attempts per IP per 10 min
           "per_account": (10, timedelta(minutes=15)),   # 10 attempts per account per 15 min
           "global_failure_rate": (1000, timedelta(minutes=1)),  # circuit breaker
       }

       def check_and_record_attempt(
           self, ip: str, email: str, success: bool
       ) -> tuple[bool, str]:
           if success:
               self._clear_on_success(ip, email)
               return True, ""

           # Check all dimensions
           checks = [
               (f"login_ip:{ip}", "per_ip", f"Too many attempts from this IP"),
               (f"login_acct:{email}", "per_account", "Account temporarily locked"),
               ("login_global_failures", "global_failure_rate", "Service temporarily unavailable"),
           ]

           for key, limit_name, message in checks:
               limit, window = self.LIMITS[limit_name]
               current = int(redis_client.get(key) or 0)
               if current >= limit:
                   return False, message

           # Record the failed attempt
           pipe = redis_client.pipeline()
           for key, limit_name, _ in checks:
               _, window = self.LIMITS[limit_name]
               pipe.incr(key)
               pipe.expire(key, int(window.total_seconds()))
           pipe.execute()

           return True, ""

       def _clear_on_success(self, ip: str, email: str):
           redis_client.delete(f"login_ip:{ip}", f"login_acct:{email}")
   ```

3. **Fingerprint devices to detect credential stuffing bots**:

   ```python
   import hashlib

   def compute_device_fingerprint(request) -> str:
       """
       Combine signals that are stable for a legitimate user but vary for bots.
       Never use fingerprint alone for blocking — use as a signal.
       """
       signals = [
           request.headers.get("User-Agent", ""),
           request.headers.get("Accept-Language", ""),
           request.headers.get("Accept-Encoding", ""),
           # JavaScript-collected signals (via frontend):
           request.json.get("screen_resolution", ""),
           request.json.get("timezone_offset", ""),
           request.json.get("canvas_hash", ""),  # canvas fingerprint
       ]
       return hashlib.sha256("|".join(signals).encode()).hexdigest()

   def is_suspicious_login(user_id: str, device_fp: str, ip: str) -> bool:
       known_devices = get_user_known_devices(user_id)
       known_ips = get_user_known_ips(user_id)

       # New device + new IP = higher risk
       if device_fp not in known_devices and ip not in known_ips:
           return True
       return False

   def handle_suspicious_login(user_id: str, request_context: dict):
       # Step-up auth: require MFA or email verification
       send_login_notification_email(user_id, request_context)
       return {"requires_step_up": True}
   ```

4. **Detect automated behavior patterns**:

   ```python
   import time

   class BotBehaviorDetector:
       def analyze_timing(self, session_start: float, form_submit: float) -> bool:
           """
           Humans take >2 seconds to fill a form.
           Bots often submit in <500ms.
           """
           fill_time = form_submit - session_start
           return fill_time < 0.5  # likely bot

       def check_honeypot(self, form_data: dict) -> bool:
           """
           Hidden field bots fill, humans leave empty.
           """
           return bool(form_data.get("_hp_field", ""))

       def check_interaction_signals(self, interaction_data: dict) -> bool:
           """
           Legitimate users move mouse, click fields, type with variable speed.
           Bots often submit without any interaction.
           """
           return (
               interaction_data.get("mouse_events", 0) == 0 and
               interaction_data.get("key_events", 0) == 0
           )
   ```

5. **Progressive response — avoid blocking legitimate users**:

   ```python
   from enum import Enum

   class ThreatLevel(Enum):
       LOW = "low"
       MEDIUM = "medium"
       HIGH = "high"
       CRITICAL = "critical"

   RESPONSES = {
       ThreatLevel.LOW: {"action": "allow"},
       ThreatLevel.MEDIUM: {"action": "captcha", "captcha_type": "recaptcha_v3"},
       ThreatLevel.HIGH: {"action": "step_up", "require": "email_verification"},
       ThreatLevel.CRITICAL: {"action": "block", "duration_minutes": 60},
   }

   def get_threat_level(ip: str, email: str, fingerprint: str, bot_signals: dict) -> ThreatLevel:
       score = 0
       if is_ip_on_blocklist(ip): score += 40
       if failed_attempts_in_window(ip) > 5: score += 20
       if failed_attempts_in_window(email) > 3: score += 20
       if bot_signals.get("timing_suspicious"): score += 10
       if bot_signals.get("honeypot_filled"): score += 30
       if bot_signals.get("no_interaction"): score += 15

       if score >= 60: return ThreatLevel.CRITICAL
       if score >= 40: return ThreatLevel.HIGH
       if score >= 20: return ThreatLevel.MEDIUM
       return ThreatLevel.LOW
   ```

## Rules

- Never block logins solely based on IP — botnets have millions of IPs; per-IP limits slow bots but don't stop them; combine with account-level and global signals.
- Breached password checks must use k-anonymity (HIBP range API) — sending full password hashes to a third party would create a credential database.
- Response timing must be identical for valid accounts and invalid accounts — timing differences allow username enumeration that feeds credential stuffing lists.
- Account lockout must use exponential backoff, not permanent lockout — permanent lockout enables DoS against a specific user.

## Common Mistakes

- **Returning different HTTP status codes for valid vs invalid usernames** — 404 for missing user vs 401 for wrong password enables account enumeration; always return 401.
- **IP allowlisting that bypasses MFA or rate limits** — corporate NAT IPs may appear in blocklists; investigate before permanent allowlisting.
- **Not logging failed login attempts with IP and user-agent** — without this data, credential stuffing campaigns are invisible until account takeovers are reported.
- **CAPTCHA as the only defense** — CAPTCHA solving services ($1–2 per 1000 solves) make CAPTCHA-only defenses economically viable to bypass; use CAPTCHA as one layer among many.
