---
name: plan-incident-response
description: Use when a production incident is declared or suspected, when building an incident response runbook, or when running a post-incident review that needs a structured process.
source: Google Site Reliability Engineering book (Beyer et al., O'Reilly 2016), PagerDuty Incident Response Guide (response.pagerduty.com)
tags: [incident-response, sre, on-call, postmortem, devops, reliability, production]
verified: true
---

# Plan Incident Response

Execute a structured live incident: declare severity, assign roles, mitigate, communicate, resolve, and run a blameless postmortem.

## Why This Is Best Practice

**Adopted by:** Google's SRE practice (documented in the SRE book, O'Reilly 2016) is the canonical reference, followed by Atlassian, PagerDuty, Slack, and Stripe who have all published their incident response processes derived from the same model. The PagerDuty Incident Response Guide (response.pagerduty.com) is used by thousands of engineering organizations and is openly licensed.
**Impact:** Google SRE data (SRE book, Ch. 14) shows that unstructured incident response — where responders simultaneously debug, communicate, and make decisions without role separation — extends mean time to resolution (MTTR) by 2–4× compared to role-separated response. PagerDuty's State of Digital Operations report (2022) found that organizations with a documented incident response process resolve incidents 23% faster and have 30% lower escalation rates.
**Why best:** Under stress, cognitive load spills from the Incident Commander to debuggers, degrading both. Role separation (IC, responder, communicator) is borrowed from aviation crew resource management (CRM) — the same principle that reduced aviation fatal accidents by 65% between 1978 and 2000 after CRM was mandated (FAA, 1995). Blameless postmortems (from the SRE book) produce systemic fixes rather than scapegoating, which Google attributes to their ability to publish post-mortems publicly and drive industry-wide learning.

Sources: Beyer, Jones, Petoff, Murphy, "Site Reliability Engineering" (O'Reilly 2016, Ch. 14–15); PagerDuty Incident Response Guide (response.pagerduty.com); FAA Advisory Circular 120-51E (Crew Resource Management Training, 2004)

## Steps

### 1. Declare and classify severity

Anyone can declare an incident. Do not wait for certainty — declare early and downgrade later if wrong. The cost of a false positive is minutes; the cost of a delayed declaration is compounded user impact.

| Severity | Definition | Response time | War room |
|---|---|---|---|
| SEV-1 | Complete outage or data loss affecting all users or production data | Immediate (< 5 min) | Required |
| SEV-2 | Significant degradation: major feature down, >20% of users affected | < 15 min | Required |
| SEV-3 | Partial degradation: minor feature, <20% of users, workaround exists | < 1 hour | Optional |
| SEV-4 | Non-urgent: cosmetic, low-traffic path, no user impact reported | Next business day | No |

Post severity in the incident channel immediately:
```
🚨 SEV-2 DECLARED — [timestamp UTC]
Summary: Payment checkout returning 500 errors for ~30% of users
IC: @alice | Responder: @bob | Comms: @carol
War room: #incident-2024-11-14 | Status page: updating
```

### 2. Assign roles — do not skip this under pressure

**Incident Commander (IC):**
- Owns the incident from declaration to close
- Makes final decisions; everyone else executes
- Does NOT debug — delegates all technical tasks
- Runs 10-min check-ins to re-assess severity and progress

**Responder(s):**
- Debug, mitigate, and fix
- Report findings to IC; do not act on mitigations without IC approval during SEV-1/2
- One responder per workstream (DB, API, infra) to avoid collision

**Communications Lead (Comms):**
- Drafts and posts all external communications (status page, customer email)
- Updates internal stakeholders (VP, support, sales) every 30 min during SEV-1/2
- Never speculates about root cause externally — use: "We are investigating"

### 3. Establish the incident channel and war room

Create a dedicated Slack/Teams channel: `#inc-YYYY-MM-DD-short-description`.

Pin to the channel:
- Incident timeline (running log, append-only)
- Runbook link (if one exists)
- Dashboard/monitoring link
- Current severity and IC name

Start a video call for SEV-1/2. Mute everyone except active speakers. Use channel for async updates so the video call stays uncluttered.

### 4. Investigate — timeline before hypotheses

Before proposing a fix, the IC demands a **timeline of changes**:

```bash
# What changed in the last 2 hours?
git log --oneline --since="2 hours ago"
kubectl rollout history deployment/payment-api
# Check infrastructure changes
terraform plan -out=tfplan && terraform show tfplan  # or check change log
```

Post findings to the incident channel with timestamps:
```
13:42 UTC — payment-api v2.14.1 deployed (contains DB schema migration)
13:47 UTC — first 500 errors appear in Datadog
13:52 UTC — error rate climbs to 32%
```

The timeline almost always points to root cause faster than log spelunking.

### 5. Mitigate first, fix later

Mitigation = stopping the bleeding. Fix = preventing recurrence.

Common mitigations (fastest options first):
1. **Rollback** — revert the last deployment. Fastest path if a deploy triggered it.
2. **Feature flag** — disable the broken feature without deploying.
3. **Traffic shift** — route traffic to a known-good instance or region.
4. **Rate limit / circuit breaker** — protect downstream services from cascading failure.
5. **Scale up** — if the cause is resource exhaustion, buy time while the fix is built.

IC approves mitigation before execution. Responder announces: "About to rollback payment-api to v2.13.9 — standing by for IC go."

### 6. Communicate externally

Update the status page within 10 minutes of declaration. Use this language:

- **Investigating:** "We are aware of an issue affecting [service]. Our team is investigating. We will provide an update in 30 minutes."
- **Identified:** "We have identified the cause of [issue] and are working on a fix."
- **Monitoring:** "A fix has been applied. We are monitoring to confirm resolution."
- **Resolved:** "This incident has been resolved. A post-incident review will be published within 5 business days."

Never write: "Due to [technical cause], users experienced..." — this is speculation before root cause is confirmed. Never promise an ETA unless you have high confidence.

### 7. Declare resolution

Resolution requires:
- Error rate at baseline for ≥ 10 minutes (not just a single green data point)
- IC confirms with responders that the fix is stable
- Status page updated to "Resolved"
- All stakeholders notified

IC announces in channel:
```
✅ RESOLVED — [timestamp UTC]
Duration: 47 min
Impact: ~30% checkout failure rate for 47 min
Mitigation: Rolled back payment-api to v2.13.9
Next: Postmortem scheduled for [date], owner: @alice
```

### 8. Run a blameless postmortem

Schedule within 5 business days for SEV-1/2. Within 2 weeks for SEV-3.

Postmortem document structure:
1. **Summary** — 2-3 sentence overview (what broke, impact, duration)
2. **Impact** — quantified user and business impact (errors, revenue, SLA breach)
3. **Timeline** — chronological log of events, detection, actions, resolution
4. **Root cause** — the specific technical condition that caused the incident (not "human error")
5. **Contributing factors** — systemic conditions that made the incident possible or worse
6. **What went well** — concrete things the team did right (detection, communication, rollback speed)
7. **Action items** — specific, owne, time-bound tasks to prevent recurrence. No action item is "be more careful."

Blameless means: the goal is systemic improvement, not assigning fault. People follow the processes and tools available to them. If a person made a mistake, ask why the system made the mistake easy to make.

## Rules

- IC does not debug — if the IC is in the terminal, they are not managing the incident
- One IC at a time — leadership confusion extends incidents. Incoming managers join as observers unless IC hands off explicitly
- The incident channel is append-only during the incident — no edits to past messages, timestamped log only
- Never speculate on root cause in external communications before it is confirmed
- Mitigation takes priority over root cause investigation during active user impact
- No action item in a postmortem is "we need to be more careful" — only systemic, verifiable changes
- SEV-1/2 postmortems are published internally to all of engineering — learning is a shared good

## Common Mistakes

- **No IC declared** — everyone debugs, no one coordinates, mitigations conflict, MTTR doubles.
- **Skipping rollback to find root cause** — prioritizing understanding over stopping user impact. Mitigate first; debug after users are unaffected.
- **Premature resolution** — declaring resolved after one green data point. Monitor for ≥ 10 minutes at baseline before declaring done.
- **Postmortem theater** — scheduling the postmortem, producing a document, and filing no action items. The postmortem is only valuable if it produces systemic change.
- **Blame in postmortem** — "The engineer should have tested this" is not a root cause. The root cause is that the system allowed untested code to reach production.
- **No status page update** — customers discover incidents from social media before the status page. Update the status page within 10 minutes, even with "we are investigating."
