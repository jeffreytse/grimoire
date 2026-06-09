---
name: design-threat-model
description: Use when designing a new system, reviewing architecture for security, adding an external-facing feature, handling sensitive data, or before a security audit.
source: 'Microsoft STRIDE model (Adam Shostack, "Threat Modeling: Designing for Security", 2014), OWASP Threat Modeling Cheat Sheet'
tags: [threat-modeling, stride, security-design, attack-surface, security-engineer, developer]
verified: true
---

# Design Threat Model

Apply STRIDE to a system's data flow diagram to systematically enumerate threats, assess risk, and produce a prioritized mitigation list.

## Why This Is Best Practice

**Adopted by:** Microsoft has required threat modeling for all new product features since the Security Development Lifecycle (SDL) was mandated company-wide in 2004 following the "Trustworthy Computing" memo. Amazon, Google, and Meta all document threat modeling as part of their security review gates. OWASP's Threat Modeling Cheat Sheet (updated 2023) formalizes STRIDE as the recommended entry-point methodology for organizations without dedicated security teams.
**Impact:** Microsoft's SDL data (published by Adam Shostack in "Threat Modeling: Designing for Security", 2014) showed that finding and fixing a vulnerability during design costs ~$0.10 per defect vs. ~$100 in production — a 1000× cost difference. The 2020 Verizon Data Breach Investigations Report (DBIR) found that 86% of breaches involved threats that were categorized under STRIDE (primarily spoofing, information disclosure, and elevation of privilege) — all findable at design time.
**Why best:** Ad hoc security review ("is this secure?") produces inconsistent coverage and depends on reviewer expertise. STRIDE ensures every trust boundary is examined against six threat categories in a repeatable, auditable way. Compared to PASTA or attack trees, STRIDE is faster to learn and produces immediately actionable findings for developers without a security background.

Sources: Adam Shostack, "Threat Modeling: Designing for Security" (Wiley, 2014); Microsoft SDL (microsoft.com/en-us/securityengineering/sdl); OWASP Threat Modeling Cheat Sheet (owasp.org); Verizon DBIR 2020

## Steps

### 1. Draw the data flow diagram (DFD)

Sketch the system using four elements:

- **External entities** (rectangles) — users, external systems, third-party APIs
- **Processes** (circles/ovals) — code that transforms data
- **Data stores** (parallel lines) — databases, caches, files, queues
- **Data flows** (arrows) — data moving between elements, labeled with protocol and data type

Mark every **trust boundary** (dashed line) where data crosses privilege levels: browser → API, API → database, internal service → external service.

Minimum viable DFD: one page, covers all external-facing flows and trust boundaries.

### 2. Apply STRIDE to each element

For every process, data store, and data flow that crosses a trust boundary, ask each STRIDE question:

| Letter | Threat | Question to ask | Common controls |
|---|---|---|---|
| **S** | Spoofing | Can an attacker impersonate a legitimate identity? | Authentication, MFA, mutual TLS |
| **T** | Tampering | Can data be modified in transit or at rest? | HMAC, signing, integrity checks, TLS |
| **R** | Repudiation | Can an actor deny performing an action? | Audit logging, non-repudiation tokens |
| **I** | Information Disclosure | Can sensitive data be read by unauthorized parties? | Encryption at rest/transit, access control |
| **D** | Denial of Service | Can availability be degraded or eliminated? | Rate limiting, circuit breakers, quotas |
| **E** | Elevation of Privilege | Can an actor gain more permissions than granted? | Least privilege, input validation, RBAC |

For each threat identified, write one line:
```
[STRIDE-category] [Element] — [threat description] — [proposed control]
```

### 3. Score and prioritize threats

Use DREAD or a simple 3×3 matrix (Likelihood × Impact) for each threat. At minimum, classify each as High / Medium / Low.

Prioritize:
1. High-impact + easy to exploit (e.g., unauthenticated RCE)
2. High-impact + hard to exploit (fix before launch)
3. Low-impact + easy to exploit (fix in next sprint)
4. Low-impact + hard to exploit (accept or defer)

### 4. Define mitigations and assign owners

For each High or Medium threat, write a mitigation card:

```
Threat: Attacker spoofs API caller by replaying a stolen JWT
STRIDE: Spoofing
Severity: High
Mitigation: Set JWT expiry to 15 min; implement refresh token rotation;
            add jti (JWT ID) claim and blacklist on logout
Owner: @auth-team
Due: before launch
```

### 5. Document residual risk

For threats accepted without full mitigation, record:
- Why the residual risk is accepted
- What compensating control exists (if any)
- Who approved the acceptance

### 6. Review cycle

Re-run threat model when:
- A new external entity is added
- A trust boundary changes (new auth mechanism, new API surface)
- A data store changes scope (new PII field, new retention policy)
- A severity-1 incident reveals an unmodeled threat

## Rules

- Threat model the design, not the implemented code — design-time findings are 1000× cheaper to fix
- Every trust boundary must be covered — no exceptions for "internal only" flows
- STRIDE categories are not mutually exclusive — a single attack can hit multiple categories
- Mitigations must be verifiable — "we'll be careful" is not a mitigation
- Threat model document is a living artifact — version it alongside the architecture diagram
- Include developers in the session — security engineers alone miss implementation-level threats

## Examples

### Example threat finding

```
[S] Payment API /charge endpoint
Threat: Attacker replays a captured request to charge a card twice.
Control: Implement idempotency key (UUID per request); reject duplicate
         idempotency keys within 24h; return 200 with original response body.
Severity: High | Owner: payments-team | Due: Sprint 14
```

### Example trust boundary on DFD

```
[Browser] --HTTPS--> [API Gateway] --mTLS--> [Payment Service]
                     ^                       ^
              Trust boundary 1        Trust boundary 2
              (internet→DMZ)          (DMZ→internal)
```

Threat model each arrow crossing a trust boundary separately.

## Common Mistakes

- **Skipping the DFD** — jumping to STRIDE without drawing the system leads to missed trust boundaries and false coverage confidence.
- **Only modeling external attackers** — insider threats and compromised internal services are responsible for 34% of breaches (Verizon DBIR 2023). Model internal trust boundaries too.
- **Accepting "encryption solves I" for everything** — encryption does not prevent information disclosure from over-permissioned queries or verbose error messages.
- **One-time exercise** — threat models decay as systems evolve. A stale threat model is worse than none — it creates false confidence.
- **No ownership on mitigations** — findings without owners are never fixed. Every threat mitigation needs a name and a date.
