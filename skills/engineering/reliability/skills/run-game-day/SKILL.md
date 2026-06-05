---
name: run-game-day
description: Use when validating incident response capabilities, onboarding new engineers to on-call, preparing for a major launch, or after significant architecture changes
source: Google SRE Workbook disaster role playing chapter; Amazon Game Day methodology (Vogels 2011); Rosenthal "Chaos Engineering" failure injection
tags: [reliability, sre, incident-response, training]
verified: true
---

# Run Game Day

Simulate realistic failure scenarios in a controlled exercise to validate incident response procedures, find gaps in runbooks, and build muscle memory before real incidents occur.

## Why This Is Best Practice

**Adopted by:** Amazon (Werner Vogels mandated game days after the 2011 AWS outage), Google (DiRT — Disaster Recovery Testing), Stripe, PagerDuty, Netflix
**Impact:** Amazon: teams that run game days quarterly reduce mean time to recover (MTTR) by 50% (Vogels 2011); Google DiRT program: found 70% of DR procedures had at least one critical gap on first execution; teams with game day practice resolve incidents 3x faster
**Why best:** Incident response is a perishable skill; runbooks written months ago are stale; procedures not practiced are procedures that fail under stress

Sources: Murphy et al. "Site Reliability Workbook" O'Reilly (2018) Ch. 9; Vogels "Working Backwards" (2021); Rosenthal "Chaos Engineering" O'Reilly (2020)

## Steps

1. **Define the game day scenario** — Choose a realistic failure scenario relevant to your system: regional cloud outage, database primary failure, cascading microservice failures, data corruption, DDoS attack, or a critical dependency outage. Base scenarios on past incidents or known high-risk failure modes. Document the scenario in detail but keep it hidden from participants.

2. **Set objectives** — Define what you want to test: How fast does the team detect the failure? Do the runbooks have all the steps needed? Can engineers restore service within the RTO? Are communication channels working? Objectives determine which aspects to measure and what success looks like.

3. **Assemble roles** — Incident commander (runs the exercise, injects events), observers (silently note gaps without helping), and responders (the on-call team treating it as a real incident). Observers must not intervene; their job is to document, not rescue. Use actual on-call engineers, not volunteers.

4. **Brief participants (partially)** — Tell responders: a game day exercise is starting, it involves a failure scenario, and they should respond exactly as they would during a real incident (use the runbook, escalate, declare incident, communicate to stakeholders). Do not reveal the scenario. The exercise must feel real to produce realistic responses.

5. **Inject the failure** — At the designated time, inject the simulated failure: silently remove a database host, block a network route, inject error responses from a dependency, or simply announce a fictional alert matching the scenario. Use your monitoring and chaos engineering tooling for realistic injection; paper-based tabletop exercises miss infrastructure behavior.

6. **Monitor and document in real time** — Observers track: time to first alert, time to first responder acknowledgment, which runbooks were consulted, which steps failed or were missing, communication quality, and decision points where responders were uncertain. Use a shared document with timestamps.

7. **Maintain the fiction** — Do not help responders who are struggling, unless they are about to take an action that could cause real harm. The gaps revealed by struggling are the game day's most valuable output. The discomfort of a difficult exercise prevents real incident discomfort.

8. **Define a kill switch** — If a responder is about to take an action that could impact real production systems, the incident commander calls "exercise stop, exercise stop, exercise stop" (NATO convention) and explains the boundary. Participants immediately stop and do not take the planned action. Real harm prevention overrides exercise realism.

9. **Conduct an immediate debrief (hot wash)** — Within 30 minutes of the exercise end: gather all participants for a 45-60 minute debrief. Structure: What was the scenario? What did we observe? What worked well? What failed or was missing? What actions would have prevented faster recovery? Capture every action item. The hot wash while the exercise is fresh produces the most valuable findings.

10. **Publish findings and track remediation** — Write a game day report: scenario, objectives, timeline, findings (positive and negative), and action items with owners and due dates. Track action items in your issue tracker with a game-day label. Review completion of action items at the start of the next game day. Unresolved action items are known gaps in your reliability posture.

## Rules

- Game days must be scheduled at least 2 weeks in advance; surprise exercises cause real outages when responders take irreversible actions under time pressure.
- The kill switch protocol must be established and understood by all participants before the exercise begins; real harm prevention is non-negotiable.
- Observers must not provide hints or help; intervention eliminates the learning value of the exercise.
- Every action item from a game day must have an owner and due date; unowned action items are never completed.

## Common Mistakes

- **Tabletop only (no infrastructure injection)** — talking through responses without actually experiencing the monitoring alerts, the runbook steps, and the tooling reveals far fewer gaps than a live exercise.
- **Rescuing responders who are struggling** — the gaps revealed by struggling are the most valuable findings; resist the urge to help.
- **No debrief** — an exercise without a structured debrief captures no learnings; the debrief is where the value is extracted.
- **Running game days without tracking action items** — a game day that produces 10 findings but no follow-through is an expensive exercise with no reliability improvement.

## When NOT to Use

- Systems under active incidents — never run a game day on a system experiencing real degradation
- Teams with no runbooks or incident response process — establish the basics first; game days validate procedures, not create them from scratch
- Immediately before a critical business event (product launch, peak season) — game days sometimes reveal issues that require emergency remediation; plan them 4-6 weeks before critical events
