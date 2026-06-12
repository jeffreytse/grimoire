---
name: design-incident-response
description: Use when creating or improving an incident response process for a production system or operational team
source: PagerDuty Incident Response Guide (2023); NIST SP 800-61r2 "Computer Security Incident Handling Guide"
tags: [incident-response, operations, reliability, on-call, postmortem, sre]
verified: true
---

# Design Incident Response

Build a structured incident response process that minimizes customer impact, coordinates responders, and prevents recurrence.

## Why This Is Best Practice

**Adopted by:** Google SRE, PagerDuty, Atlassian, AWS, and organizations following NIST SP 800-61 as the federal standard for incident handling
**Impact:** Organizations with a documented incident response process resolve incidents 28% faster and have 50% lower repeat-incident rates than those without (PagerDuty State of Digital Operations, 2022)

**Why best:** Without a defined process, incidents devolve into heroics — one person firefighting while everyone else watches. A defined incident command structure, clear roles, and a blameless postmortem loop convert incidents into organizational learning.

## Steps

1. **Define severity levels** — Create 3–4 tiers (SEV1–SEV3 or P0–P3) with explicit, measurable criteria: "SEV1 = customer-facing service unavailable for >1% of users or data loss." Avoid vague language like "significant impact."
2. **Assign incident command roles** — Incident Commander (IC): owns the incident and decisions. Communications Lead: updates stakeholders. Subject Matter Experts (SMEs): do the technical work. No one plays two roles during a SEV1.
3. **Create the alert-to-acknowledge SLA** — Define time limits per severity: SEV1 = 5 min page, 15 min acknowledge; SEV2 = 30 min. Document escalation path if the primary on-call does not respond.
4. **Build the incident channel template** — On-call tooling should auto-create a dedicated channel with: severity, incident summary, IC assigned, customer impact, and links to runbooks. Standardize the format so responders orient instantly.
5. **Define the update cadence** — SEV1: status update every 15 minutes to stakeholders. SEV2: every 30 minutes. Updates follow the format: Current state / Actions taken / Next action / ETA. No open-ended updates.
6. **Conduct blameless postmortem within 48 hours** — Document timeline, contributing factors (no "root cause" — systems have multiple causes), and action items with owners and due dates. Publish internally.
7. **Track postmortem action items** — Assign each item to a team with a sprint commitment. Review completion in the next postmortem. Untracked items repeat the incident.

## Rules

- SEV criteria must be objective and measurable — not "feels bad" or "leadership is watching."
- The Incident Commander does not perform technical work during the incident; their job is coordination only.
- Postmortems are blameless by policy — name systems and decisions, never individuals.
- Publish all postmortems to the full engineering organization within one week.
- Every runbook referenced during an incident must be updated within 24 hours if it was wrong or missing steps.

## Examples

PagerDuty's own incident command model assigns an "External Communications" role during SEV1s who owns the status page and customer email — keeping the IC focused purely on resolution. This separation reduced customer complaint volume by 40% by keeping communications proactive rather than reactive.

## Common Mistakes

- **IC also coding** — The Incident Commander context-switching into technical work means no one is tracking the big picture; incidents take longer to resolve.
- **Postmortem theater** — Writing a postmortem but not tracking action items means the same incident recurs; the postmortem becomes a ritual with no learning loop.
- **Vague severity definitions** — When SEV1 means "whatever the on-call thinks is bad," escalation and response are inconsistent across teams.

## When NOT to Use

- Do not apply a full incident command structure to routine maintenance windows or planned deployments that have a pre-approved rollback plan and zero unexpected customer impact.
- Do not use this process for one-person teams or solo operators where assigning separate Incident Commander and SME roles is structurally impossible — instead document a simplified checklist tailored to solo response.
- Do not run a blameless postmortem for every minor bug fix or low-severity alert; reserve the process for SEV1/SEV2 events where the effort of a full postmortem is proportional to the potential for systemic learning.
