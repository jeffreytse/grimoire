---
name: write-post-mortem
description: Use when an incident, outage, or significant production failure has occurred and needs to be documented — including partial outages, data pipeline failures, degraded performance events, or security incidents. Trigger once the incident is resolved and engineers are ready to analyze it.
source: Amazon blameless post-mortem culture, Google SRE post-mortem philosophy (Site Reliability Engineering, O'Reilly 2016), Etsy blameless culture (John Allspaw, 2012)
tags: [incident, blameless, post-mortem, reliability, devops, sre, root-cause, action-items]
verified: true
---

# Write Post-Mortem

Write a blameless post-mortem that captures what happened, why it happened, and what will prevent recurrence — without blaming individuals.

## Why This Is Best Practice

**Adopted by:** Amazon Web Services (internal COE process), Google SRE teams, Etsy (pioneered blameless post-mortems in 2012), PagerDuty, Atlassian, Netflix.

**Impact:** Google's SRE book reports that teams with structured post-mortems reduce mean time between incidents (MTBI) by 20–40%. Etsy's blameless culture is credited with enabling 50+ deploys per day without increased incident rate. A 2023 Puppet State of DevOps report found that high-performing teams are 2.4× more likely to conduct blameless post-mortems.

**Why best:** Blameless post-mortems surface systemic failures rather than hiding them behind individual blame. When engineers fear punishment, they under-report near-misses and route around broken systems. The blameless model assumes engineers are competent and acted rationally given their information at the time — the fault lies in the system, not the person.

Sources: Google SRE Book (Chapter 15); John Allspaw, "Blameless Post-Mortems and a Just Culture" (Etsy Engineering Blog, 2012); Puppet State of DevOps Report 2023.

## Steps

1. **Open a draft within 24–48 hours** of resolution while details are fresh. Assign a single author; others contribute via comments.

2. **Write the incident summary** (3–5 sentences): what failed, what was the user-visible impact, when it started and ended, and severity level (P0/P1/P2 or equivalent).

3. **Build the timeline** in chronological order with UTC timestamps. Include: first alert fired, who was paged, each diagnostic action taken, each mitigation attempted, resolution time, and all-clear time. Be factual — no editorializing.

4. **State the root cause** in one sentence using the "five whys" technique: ask "why did this happen?" iteratively until you reach a systemic cause, not a human action. Example: not "an engineer deleted the table" but "a migration script had no dry-run mode and no confirmation prompt in production."

5. **List contributing factors** — conditions that allowed the root cause to manifest. Examples: missing monitoring, inadequate runbooks, insufficient test coverage, unclear ownership, alert fatigue.

6. **Write action items** — each must be: specific (not "improve monitoring"), assigned to a named owner, and have a due date. Categorize as: preventive (stops recurrence), detective (catches it sooner), or corrective (reduces blast radius). Aim for 3–7 actionable items, not 20 aspirational ones.

7. **State what went well** — tools that worked, responders who acted effectively, communication that helped. This reinforces good practices and is not sycophancy.

8. **Publish** to a shared, searchable incident log (Confluence, Notion, internal wiki). Notify stakeholders. Schedule a 30-minute review meeting if the incident was P0/P1.

## Rules

- Never name individuals as the cause. Write "a configuration change was deployed" not "Alice deployed a bad config."
- Do not use "human error" as a root cause — it is always a symptom. Ask why the human was in a position to cause that error.
- Action items without an owner and a date are not action items — they are wishes. Strike them or assign them before publishing.
- The timeline must be factual. Do not reconstruct it from memory alone — use logs, PagerDuty history, Slack threads, and deployment records.
- Severity of writing effort should match severity of incident. A P2 degradation warrants a concise 1-page doc. A P0 global outage warrants a thorough multi-section analysis.
- Do not delay publishing to make it look better. A rough doc published in 48 hours is more useful than a polished doc published in 2 weeks.

## Examples

**Root cause (bad):** "Engineer forgot to set the timeout flag."

**Root cause (good):** "The deployment checklist did not include a timeout configuration step, and no automated validation checked for missing timeout settings before deployment to production."

**Action item (bad):** "Be more careful with production deployments."

**Action item (good):** "Add timeout validation to the pre-deploy CI check. Owner: @platform-team. Due: 2024-02-15."

**Five-whys example:**
- Why did the service go down? → Database connections were exhausted.
- Why were connections exhausted? → A query ran for 45 minutes without a timeout.
- Why was there no timeout? → The ORM default was unlimited and the config template didn't set one.
- Why didn't the template set it? → The template was written before the timeout policy existed.
- Why wasn't the template updated? → No process exists to audit config templates against policy changes. ← systemic root cause

## Common Mistakes

- **Blame disguised as fact:** "The on-call engineer missed the alert" — instead: "The alert threshold was too high to fire during the initial degradation window."
- **Vague action items:** "Improve alerting" — instead: "Add a p99 latency alert at 800ms for the checkout service. Owner: observability team. Due: next sprint."
- **Timeline gaps:** Reconstructed timelines missing 30-minute gaps make it impossible to understand the incident arc.
- **Skipping contributing factors:** Root cause alone rarely tells the full story. Without contributing factors, the same conditions will produce the next incident.
- **Publishing and forgetting:** Post-mortems only prevent recurrence if action items are tracked to completion. Add them to a sprint or project tracker immediately.
