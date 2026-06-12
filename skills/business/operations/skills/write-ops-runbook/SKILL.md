---
name: write-ops-runbook
description: Use when documenting a repeatable operational procedure that on-call engineers or operators must execute under pressure
source: Google SRE Book (Beyer et al., 2016), Chapter 12 "Effective Troubleshooting"; PagerDuty Runbook Guide
tags: [runbook, operations, on-call, documentation, sre, reliability]
verified: true
---

# Write Runbook

Create an operational runbook that any qualified engineer can execute correctly under pressure without prior context.

## Why This Is Best Practice

**Adopted by:** Google SRE teams, AWS operations, Cloudflare, PagerDuty — all treat runbooks as living operational documentation
**Impact:** Google SRE reports that documented runbooks reduce mean time to resolution (MTTR) by 40–60% for known failure modes; they also reduce on-call burnout by eliminating repeated tribal-knowledge lookups

**Why best:** A good runbook transfers expert knowledge into a form that a competent non-expert can execute at 2 AM during an incident. Every step of ambiguity in a runbook multiplies the risk of error under stress.

## Steps

1. **Write the header block** — Include: Alert name or trigger condition, severity level, service owner, last reviewed date, and link to the service dashboard. Someone opening this runbook needs orientation in under 10 seconds.
2. **State the symptoms** — List the exact alert text, error messages, or observable behaviors that bring someone to this runbook. Include screenshot or log snippet if available.
3. **Define the impact** — Who is affected? What is degraded? What is the customer-visible symptom? One paragraph, no jargon.
4. **List prerequisites** — Tools, permissions, VPN access, environment variables, or credentials required before starting. Do not assume the reader has these set up.
5. **Write numbered steps** — Each step is one action. Include the exact command to run (copy-pasteable), the expected output, and what to do if the output differs. Steps must be in order; no "see also" references mid-procedure.
6. **Add decision branches** — At key points, add "If X, go to Step 8. If Y, go to Step 12." Avoid branching more than twice in a single runbook; split into separate runbooks instead.
7. **Document escalation** — At what point does the responder escalate? To whom? With what information? Include Slack channel, PagerDuty escalation policy, and on-call contact.
8. **Include rollback steps** — For any action that modifies state (restart, deploy, config change), provide the exact rollback command. Never leave the operator without an undo.

## Rules

- Every command must be copy-pasteable with no placeholder editing required during execution.
- Runbooks are reviewed and tested at least quarterly; stale runbooks are worse than no runbooks.
- Each runbook covers exactly one failure mode — do not bundle multiple procedures.
- Avoid prose explanations mid-procedure; move background to an appendix so operators skip to steps.
- Link to the postmortem that created this runbook so the reader understands the failure context.

## Examples

A Redis memory runbook might include:
- **Step 3:** Check memory usage: `redis-cli -h $REDIS_HOST info memory | grep used_memory_human`
- Expected output: `used_memory_human:1.50G`
- If output shows >4G → proceed to Step 4 (eviction policy check). If Redis is unreachable → go to Step 9 (failover).

## Common Mistakes

- **Prose runbooks** — Paragraphs of explanation without numbered steps require the operator to parse meaning under stress; errors increase dramatically.
- **Stale commands** — Commands referencing deprecated flags or old hostnames fail silently; operators lose time debugging the runbook instead of the incident.
- **Missing rollback** — An operator who cannot undo a step will hesitate or skip it; hesitation extends incidents.

## When NOT to Use

- Do not write a runbook for a novel incident type that has never recurred — invest that effort in a postmortem and write the runbook only after the second occurrence confirms it is a repeatable failure mode.
- Do not use a runbook to document architectural decisions or system design rationale; those belong in an ADR or design doc, and mixing them into runbooks forces operators to parse explanation when they need action.
- Do not create runbooks for procedures that require real-time expert judgment on ambiguous system state — if the steps cannot be deterministic, a runbook creates false confidence and should be replaced with an escalation path to the subject matter expert.
