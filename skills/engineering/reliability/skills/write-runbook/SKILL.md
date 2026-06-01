---
name: write-runbook
description: Use when writing or updating an on-call runbook, incident response playbook, or alert-to-action guide
source: "The Site Reliability Workbook (Google, 2018) Ch. 8; PagerDuty Incident Response Documentation; The Art of SRE (Niall Murphy et al.)"
tags: [sre, runbook, on-call, incident-response, reliability, operations, playbook]
---

# Write Runbook

Write a runbook that an on-call engineer who has never seen this service can follow at 3 AM to diagnose and resolve an incident.

## Why This Is Best Practice

**Adopted by:** Google SRE teams; PagerDuty's incident response model; AWS operational excellence pillar in the Well-Architected Framework

**Impact:** PagerDuty reports that teams with documented runbooks reduce mean time to resolution (MTTR) by 40-60%; AWS Well-Architected reviews flag missing runbooks as a reliability risk requiring remediation

Runbooks fail when they assume knowledge the on-call engineer does not have at 3 AM under stress. The baseline assumption must be: the reader knows nothing about this specific service, is tired, and has five minutes before an executive asks for an update.

## Steps

1. **Write the alert header** — link the runbook to the exact alert name; include: alert severity, SLO at risk, expected MTTR, and escalation contact
2. **Describe the symptom in user-facing terms** — "Checkout requests failing" not "HTTP 502 on pod replica set"; state what users are experiencing
3. **Write the triage checklist** — ordered list of diagnostic commands with expected output; annotate what "normal" looks like so the engineer can confirm the hypothesis
4. **Provide branching resolution paths** — cover the top 3-5 root causes with distinct remediation steps for each; label each branch with its trigger condition
5. **Include rollback steps** — every runbook must have a "make it stop now" option: rollback deployment, flip feature flag, shed load; this comes before the root-cause fix
6. **Add escalation criteria** — state explicitly when to page the service owner, when to declare a major incident, and who to contact in what order
7. **End with post-incident steps** — file the ticket, update the runbook if it was wrong, and link to the postmortem template

## Rules

- Every command must be copy-paste executable — no placeholders like `<your-cluster-name>` without a pointer to where to find the value
- Runbooks must be tested: run through them in a staging incident before publishing
- Keep runbooks under 2 pages; if longer, split by symptom into multiple runbooks
- Update runbooks within 48 hours of any incident where the runbook was used or found lacking

## Examples

**Alert:** `CheckoutLatencyHigh` (P1, SLO: 99.5% of requests under 800ms, MTTR target: 30 min)
**Triage command:** `kubectl top pods -n checkout-service | sort -k3 -rn | head -5` — expected: no single pod above 80% CPU; if one pod is at 100%, proceed to Branch A (pod restart).

## Common Mistakes

- Runbook describes symptoms but not commands: a reader should never have to decide what to check — the runbook decides for them
- Missing rollback path: root cause fixes take time; users need relief now; always provide an immediate mitigation first
- Runbook last updated when the service was first deployed: stale runbooks are worse than no runbook because they build false confidence
