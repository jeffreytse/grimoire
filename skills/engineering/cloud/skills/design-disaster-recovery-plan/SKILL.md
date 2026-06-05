---
name: design-disaster-recovery-plan
description: Use when establishing or auditing business continuity for cloud systems that must recover from catastrophic failures, data loss, or prolonged outages
source: AWS Well-Architected Reliability Pillar; NIST SP 800-34 Rev.1; ISO 22301:2019 Business Continuity Management
tags: [cloud, reliability, disaster-recovery, business-continuity]
verified: true
---

# Design Disaster Recovery Plan

Create a tested, documented plan that recovers critical systems within defined time and data-loss targets.

## Why This Is Best Practice

**Adopted by:** Financial institutions (required by FFIEC, BCBS239), healthcare (HIPAA §164.308), AWS Well-Architected framework mandates it for reliability pillar
**Impact:** Organizations with tested DRPs recover 3x faster (IBM 2022); average breach without DRP costs $4.7M more in downtime losses; NIST compliance requires documented RTOs
**Why best:** Untested assumptions about recovery are the leading cause of extended outages; a written, drilled plan eliminates ambiguity under pressure

Sources: AWS Well-Architected Reliability Pillar (2023); NIST SP 800-34 Rev.1 (2010); ISO 22301:2019

## Steps

1. **Conduct business impact analysis (BIA)** — Identify all systems, rank by business criticality, and quantify downtime cost ($/hour). Output: tiered system inventory with financial impact per outage hour.

2. **Define RTO and RPO per tier** — Recovery Time Objective (max acceptable downtime) and Recovery Point Objective (max acceptable data loss). Tier 1 (mission-critical): RTO < 1 hr, RPO < 15 min. Tier 2: RTO < 4 hr, RPO < 1 hr. Tier 3: RTO < 24 hr, RPO < 24 hr.

3. **Select DR strategy per tier** — Backup/restore (lowest cost, highest RTO); pilot light (minimal standby, minutes to scale); warm standby (reduced-capacity running system); active-active (highest cost, near-zero RTO). Match strategy to RTO/RPO; don't gold-plate low-criticality systems.

4. **Design data backup architecture** — Define backup frequency, retention period, and geographic distribution. Enforce 3-2-1 rule: 3 copies, 2 different media, 1 offsite. Test restores monthly; a backup never tested is not a backup.

5. **Document recovery procedures** — Write step-by-step runbooks for each failure scenario (region failure, data corruption, ransomware, service dependency outage). Include commands, escalation contacts, and decision trees. Store runbooks outside the failed system.

6. **Implement infrastructure as code** — All infrastructure must be reproducible from code (Terraform, CloudFormation, Pulumi). Manual click-ops environments cannot be reliably recreated under pressure. IaC is a prerequisite for sub-hour RTO.

7. **Define communication plan** — Establish incident commander role, internal escalation chain, and external stakeholder communication cadence. Draft status page templates in advance. Identify who authorizes failover decisions.

8. **Automate failover where possible** — Use health checks and automated DNS failover for tier 1 systems. Manual failover adds 15-60 minutes of human decision time. Automate the detection-to-failover path; require human sign-off only for irreversible actions.

9. **Conduct tabletop exercises** — Quarterly: walk through disaster scenarios with stakeholders without actually failing systems. Identify gaps in procedures, ownership ambiguity, and missing automation. Document findings and remediate.

10. **Run full DR drills** — Annually (or after major architecture changes): execute actual failover to DR environment. Measure actual RTO/RPO against targets. Treat drill failures as production incidents and fix root causes.

## Rules

- RTOs and RPOs are contracts, not aspirations; if you haven't measured actual recovery time in a drill, you don't have a DRP.
- Runbooks must be stored and accessible when the primary system is down — print copies, secondary wiki, offline storage.
- Every dependency (DNS provider, CDN, SaaS) must be included in the BIA; partial recovery that misses a dependency fails.

## Common Mistakes

- **Setting RTO/RPO without measuring cost** — arbitrary targets lead to over-engineering or under-engineering; always tie to business impact.
- **Never testing restores** — backup jobs succeed silently for months before a corrupt backup is discovered during an actual incident.
- **Manual infrastructure recreation** — clicking through consoles during a disaster is slow and error-prone; IaC is mandatory for realistic RTOs.
- **Single person knows the DR procedure** — bus-factor-one runbooks fail when that person is unavailable; at least two people must be trained per procedure.

## When NOT to Use

- Prototype or experimental systems with no production traffic
- Systems where downtime cost is less than the ongoing cost of DR infrastructure
- Greenfield projects in the first sprint — establish DR after architecture stabilizes
