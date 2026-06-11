---
name: prevent-redos
description: Use when writing or reviewing regular expressions that process user-controlled input — especially patterns with nested quantifiers, alternation with overlapping cases, or patterns applied to long strings.
source: 'OWASP Regular Expression Denial of Service (ReDoS) Defense Cheat Sheet (owasp.org/www-project-cheat-sheets); CWE-1333; OWASP Top 10 2021 A06'
tags: [security, owasp, redos, regex, denial-of-service, input-validation, developer]
---

# Prevent ReDoS

Eliminate catastrophic regex backtracking by using atomic groups, possessive quantifiers, or linear-time regex engines — preventing single-user input from consuming CPU for minutes or hours.

## Why This Is Best Practice

**Adopted by:** CWE-1333 (Inefficient Regular Expression Complexity) is a named weakness tracked by NVD with hundreds of CVEs. Node.js, Cloudflare, and Stack Overflow have all experienced ReDoS outages. OWASP Top 10 2021 A06 (Vulnerable and Outdated Components) includes ReDoS in third-party regex libraries. Google's `re2` library (used in Go's stdlib regex) is designed to guarantee linear-time matching. GitHub's CodeQL, Semgrep, and Snyk all include ReDoS detection rules.
**Impact:** Stack Overflow's 2016 outage was caused by a single malicious regex `^\s*(\w+\s*)+$` consuming 100% CPU on a crafted input string — taking down the site for 34 minutes. Cloudflare's 2019 global outage lasted 27 minutes due to a ReDoS in a WAF regex. Node.js has had multiple ReDoS CVEs in the `validator`, `marked`, and `moment` packages. A single well-crafted string can freeze a single-threaded server.
**Why best:** Input length limits alone are insufficient — catastrophic backtracking is exponential, so even 30-character inputs can cause seconds-long hangs. Rewriting to atomic/possessive patterns or using a linear-time engine eliminates the root cause.

Sources: OWASP ReDoS Cheat Sheet; Stack Overflow outage post-mortem (2016); Cloudflare outage post-mortem (2019); CWE-1333

## Steps

1. **Identify catastrophically vulnerable regex patterns** — look for:

   - Nested quantifiers: `(a+)+`, `(a*)*`, `([ab]+)+`
   - Alternation with overlapping options: `(a|aa)+`
   - Multiple optional groups applied to the same characters: `(a+b?)+`

   ```python
   # VULNERABLE — exponential backtracking on "aaaaaaaaaaaaaaaax"
   import re
   pattern = re.compile(r'^(\w+\s*)+$')           # nested quantifier
   pattern = re.compile(r'^(a+)+$')               # classic evil regex
   pattern = re.compile(r'^([a-zA-Z]+)*$')        # same structure

   # Test: time how long this takes (should be microseconds, not seconds)
   import time
   start = time.time()
   re.match(r'^(\w+\s*)+$', 'a' * 30 + '!')
   print(f"{time.time() - start:.3f}s")  # > 1s = vulnerable
   ```

2. **Use atomic groups or possessive quantifiers to prevent backtracking**:

   ```python
   # Python 3.11+ supports atomic groups (?>...) and possessive quantifiers
   import re

   # VULNERABLE
   pattern = re.compile(r'^(\w+\s*)+$')

   # FIXED — atomic group prevents backtracking into the group
   pattern = re.compile(r'^(?>\w+\s*)+$')   # Python 3.11+

   # Alternative for older Python — use regex module
   import regex
   pattern = regex.compile(r'^(\w+\s*)++$')  # possessive quantifier ++
   ```

3. **Use `re2` or linear-time engines for user-facing regex**:

   ```python
   # Python — google-re2 (linear time guarantee)
   pip install google-re2

   import re2
   pattern = re2.compile(r'^\w+(\.\w+)*$')
   match = pattern.match(user_input)
   # Note: re2 does not support backreferences or lookaheads
   ```

   ```javascript
   // Node.js — use 're2' package
   const RE2 = require('re2');
   const pattern = new RE2('^\\w+(\\.\\w+)*$');
   pattern.test(userInput);  // linear time, no catastrophic backtracking
   ```

4. **Set timeouts on regex execution** for patterns that cannot be rewritten:

   ```python
   import signal

   class RegexTimeout(Exception):
       pass

   def safe_match(pattern, text, timeout_ms=100):
       def handler(signum, frame):
           raise RegexTimeout("Regex timed out")
       signal.signal(signal.SIGALRM, handler)
       signal.setitimer(signal.ITIMER_REAL, timeout_ms / 1000)
       try:
           result = pattern.match(text)
           signal.setitimer(signal.ITIMER_REAL, 0)
           return result
       except RegexTimeout:
           return None
   ```

   Note: `signal.SIGALRM` is POSIX-only (not Windows).

5. **Limit input length before applying regex**:

   ```python
   MAX_INPUT_LEN = 1000  # adjust per expected input

   def validate_input(user_input, pattern):
       if len(user_input) > MAX_INPUT_LEN:
           raise ValueError("Input too long")
       return pattern.match(user_input)
   ```

   Length limits reduce the window of catastrophic backtracking but don't eliminate it — pair with pattern fixing.

6. **Use static analysis to detect vulnerable patterns at CI time**:

   ```bash
   # Semgrep ReDoS rules
   semgrep --config p/redos .

   # vuln-regex-detector (Node.js)
   npx vuln-regex-detector check '(a+)+'

   # safe-regex (npm)
   npx safe-regex '(\w+\s*)+'  # exits non-zero if unsafe
   ```

## Rules

- `re2` does not support backreferences (`\1`), lookaheads, or lookbehinds — evaluate compatibility before switching.
- Possessive quantifiers (`++`, `*+`, `?+`) in PCRE/Python `regex` module prevent backtracking into that group — they are safe but also change match semantics.
- Alternation with shared prefixes is safe — `(abc|abd)` — alternation with overlapping suffixes is dangerous — `(a|aa)+`.
- Testing a regex with crafted input in development is not sufficient — use static analysis tools.

## Common Mistakes

- **Assuming third-party library regexes are safe** — npm packages like `validator`, `path-to-regexp`, and `moment` have had ReDoS CVEs. Check `package.json` deps with `npx auditjs`.
- **Using complex regex for email validation** — the RFC-compliant email regex is catastrophically backtracking. Use a simple `\S+@\S+\.\S+` or a dedicated library.
- **Not testing with adversarial input** — a regex that works fine on valid inputs may still explode on malformed ones; always test with `'a' * 30 + '!'` style inputs.
- **Fixing ReDoS with length checks alone** — catastrophic backtracking is exponential; even short inputs (20–50 chars) can take seconds.
