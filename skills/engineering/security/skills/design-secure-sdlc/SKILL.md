---
name: design-secure-sdlc
description: Use when establishing or improving a software development lifecycle — embedding security requirements, threat modeling, and security testing at each phase to find vulnerabilities when they cost 10× less to fix.
source: 'OWASP Top 10 Proactive Controls C04 (owasp.org/www-project-top-ten/); OWASP SAMM v2.0; Microsoft Security Development Lifecycle (SDL); NIST SP 800-218 (SSDF)'
tags: [security, owasp, sdlc, threat-modeling, security-requirements, appsec, developer]
---

# Design Secure SDLC

Integrate security requirements at design, threat modeling in sprint planning, SAST/DAST in CI, and security acceptance criteria in definition of done — finding vulnerabilities at 10× lower cost than post-deployment remediation.

## Why This Is Best Practice

**Adopted by:** OWASP Top 10 Proactive Controls C04 (Address Security from the Start) and OWASP SAMM v2.0 are the primary references. Microsoft's Security Development Lifecycle (SDL) has been mandatory for all Microsoft products since 2004 and documented a 50% reduction in security vulnerabilities post-SDL. NIST SP 800-218 (Secure Software Development Framework) is the US federal mandate for software supply chain security. Google, AWS, and GitHub all publicly describe security-integrated SDLC as core to their engineering practice.
**Impact:** IBM Systems Science Institute's research (confirmed by NIST) found that fixing a security defect in requirements costs $0.1×, in design $0.3×, in development $1×, in testing $3×, and in production $30×. Gartner (2022) estimates that organizations integrating security into SDLC reduce security incidents by 70% compared to bolt-on security testing. Microsoft's SDL reduced critical vulnerabilities by 50% in Windows Vista vs. Windows XP across the same codebase by adding security requirements, threat modeling, and mandatory security review.
**Why best:** Penetration testing at the end of a release cycle (the alternative) finds vulnerabilities too late — they've already been designed, coded, and deployed. Retrospective security testing cannot change architecture decisions that create entire classes of vulnerabilities. Shifting security left means decisions like "use parameterized queries" are made in sprint 1, not discovered as a critical finding post-launch.

Sources: OWASP Top 10 Proactive Controls C04; OWASP SAMM v2.0; Microsoft Security Development Lifecycle documentation; NIST SP 800-218; IBM Systems Science Institute (defect cost research)

## Steps

1. **Define security requirements at project start using abuse cases**:

   ```markdown
   # Security Requirements Template

   For each user story, add corresponding abuse cases:

   User story: "As a user, I can reset my password via email link"

   Abuse cases:
   - Attacker requests password reset for victim's email (enumeration)
   - Attacker reuses an expired reset link
   - Attacker brute-forces the reset token
   - Attacker intercepts the reset email (phishing)

   Security acceptance criteria:
   - Reset tokens must be cryptographically random (128-bit entropy)
   - Reset tokens expire after 15 minutes
   - Reset tokens are single-use (invalidated after first use)
   - Old token invalidated when new reset requested
   - Password reset doesn't confirm whether email exists (generic message)
   - Rate limit: 3 reset requests per email per hour
   ```

2. **Run threat modeling in sprint planning for new features**:

   ```markdown
   # Sprint Threat Model (STRIDE per component)

   Feature: Add payment processing
   Session: 45 minutes, security champion + lead developer

   ## Data Flow Diagram
   Browser → API Gateway → Payment Service → Stripe API
                                          → Payments DB

   ## STRIDE Analysis
   | Component | Threat | Mitigation |
   |---|---|---|
   | API Gateway | Spoofing: forged user ID | JWT validation, IRSA |
   | Payment Service | Tampering: amount manipulation | Server-side amount from order, not request |
   | Stripe API | Info Disclosure: card data in logs | Never log card numbers; use Stripe tokens |
   | Payments DB | Privilege Escalation | Separate DB user with SELECT/INSERT only |

   ## Security Stories (added to sprint backlog)
   - [ ] Validate payment amount server-side from order record (not request body)
   - [ ] Implement idempotency key for duplicate payment prevention
   - [ ] Add PCI-required logging (exclude card numbers)
   ```

3. **Integrate SAST and dependency scanning in CI**:

   ```yaml
   # GitHub Actions: security gates in CI
   security:
     runs-on: ubuntu-latest
     steps:
     - uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683

     # SAST: static analysis
     - name: Run Semgrep SAST
       uses: semgrep/semgrep-action@713efdd345f3035192eaa63f56867b88e63e4e5d
       with:
         config: "p/owasp-top-ten p/security-audit"
         fail_on: error

     # SCA: dependency vulnerabilities
     - name: Audit dependencies
       run: |
         pip install pip-audit
         pip-audit --fail-on-vuln --severity high

     # Secrets detection
     - name: Detect secrets
       uses: gitleaks/gitleaks-action@cb7149a9b57195b609c63e8518d2c6ef96efa375

     # DAST: dynamic testing on staging (separate job)
     - name: OWASP ZAP baseline scan
       if: github.ref == 'refs/heads/main'
       uses: zaproxy/action-baseline@73d8a505fb288291e668f9d9930ae05cb7b2da21
       with:
         target: https://staging.app.company.com
   ```

4. **Security in Definition of Done**:

   ```markdown
   # Definition of Done — Security Checklist

   Before marking any story complete:

   Code:
   □ No new HIGH/CRITICAL SAST findings (Semgrep, Bandit, ESLint-security)
   □ No new vulnerable dependencies (pip-audit / npm audit)
   □ No hardcoded secrets (Gitleaks clean)
   □ Security acceptance criteria from abuse cases met

   Review:
   □ Security champion reviewed changes touching auth, payments, PII, or admin
   □ Threat model updated if new external integrations added

   Testing:
   □ Security test cases written for each abuse case
   □ OWASP Top 10 relevant checks verified (SQLi, XSS, IDOR, etc.)

   Documentation:
   □ Sensitive data flows documented
   □ New permissions/access changes documented
   ```

5. **Designate security champions per team**:

   ```markdown
   # Security Champion Program

   Role: One developer per squad acts as security liaison
   Time commitment: ~10% of sprint time

   Responsibilities:
   - Facilitate threat modeling sessions for new features
   - Review security-sensitive PRs (auth, payments, PII handling)
   - Triage SAST/SCA findings and create backlog items
   - Attend monthly security champion sync with AppSec team
   - Spread security knowledge within their squad

   Training: OWASP Top 10, Secure Code Review, Threat Modeling
   (AppSec team provides training; security champions provide scale)
   ```

## Rules

- Threat modeling must happen BEFORE coding — threat models done after implementation are review artifacts, not design tools.
- SAST findings with HIGH/CRITICAL severity must block PR merge — introducing known vulnerabilities knowingly requires sign-off from security.
- Abuse cases must be written for every feature that handles authentication, authorization, PII, or financial data — not optional for these categories.
- Security acceptance criteria belong in the story, not a separate security story — security requirements for a feature die when separated from the feature.

## Common Mistakes

- **Threat modeling as a one-time annual exercise** — architecture changes constantly; threat models must be updated with each significant feature or architecture change.
- **SAST with no agreed triage process** — SAST generates findings that need prioritization; without a triage process, engineers suppress or ignore findings rather than fix them.
- **Security champion with no time allocation** — security review is invisible work that gets deprioritized; explicitly protect 10% of the champion's sprint time or it won't happen.
- **Security requirements as a security team document** — security requirements owned by the security team become a bottleneck; embed security champions to scale requirement creation across squads.
