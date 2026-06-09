---
name: design-data-governance-policy
description: Use when designing a data governance policy covering data ownership, classification, access control, lifecycle management, and regulatory compliance
source: 'DAMA International "DMBOK2: Data Management Body of Knowledge" 2nd ed. (2017); ISACA "COBIT 2019 Data Governance Framework"; GDPR Article 5 data processing principles (2018); NIST Privacy Framework v1.0 (2020)'
tags: [data, governance, gdpr, privacy, data-classification, access-control, compliance]
verified: true
---

# Design Data Governance Policy

Establish formal policies and operating procedures for how data is owned, classified, accessed, retained, and disposed of — enabling data utility while managing privacy, security, and regulatory risk.

## Why This Is Best Practice

**Adopted by:** All EU-regulated organizations (GDPR mandatory since 2018), US Federal agencies (NIST Privacy Framework), financial institutions (SOX, BCBS 239), healthcare organizations (HIPAA), and enterprise-scale organizations including JPMorgan Chase, Philips, and Unilever which have published their data governance frameworks publicly

**Impact:** GDPR enforcement actions totaled €4.5 billion in fines through 2023 (DLA Piper GDPR Fines Report), with most violations attributable to inadequate governance: unlawful access, excessive retention, and unclear ownership. Organizations with mature data governance (CMMI Data Management Level 3+) report 30% lower data breach costs (IBM Cost of Data Breach 2023) due to faster containment from known data inventories. Gartner found that organizations with formal data governance reduce time to trusted data for analytics from weeks to days.

**Why best:** Without governance, data proliferates without accountability: sensitive data lands in unsecured locations, access is over-provisioned and never revoked, retention exceeds necessity (increasing breach exposure), and regulatory obligations are unknown until an audit. Governance creates a systematic framework that makes data manageable, auditable, and compliant by design rather than by accident.

Sources: DAMA International "DMBOK2" (2017); ISACA "COBIT 2019" Data Governance guidance; GDPR Articles 5, 25, 30, 32 (2018); NIST Privacy Framework v1.0 (2020); IBM "Cost of a Data Breach Report" (2023); DLA Piper "GDPR Data Breach Survey" (2024)

## Steps

1. **Establish a data governance operating model — define who decides** — Create a Data Governance Council (executive sponsors + domain leads) that sets policy, and Data Stewards (domain-level owners) who implement and enforce policy. Document the RACI for data governance decisions: who is Responsible, Accountable, Consulted, and Informed for data policy changes, classification decisions, access approvals, and breach response. Without a clear operating model, governance policies are written and ignored.

2. **Inventory all data assets — you cannot govern what you don't know** — Conduct a data inventory across all systems: databases, data warehouses, data lakes, SaaS applications, and file stores. For each dataset document: system of record, data owner, data domains, presence of PII/regulated data, and approximate record count. Maintain the inventory in a data catalog. Automate discovery using catalog tools (Alation, DataHub, Collibra) with connectors to source systems.

3. **Define a data classification scheme — create a tiered sensitivity model** — Establish 3–4 classification tiers aligned to regulatory and business sensitivity:
   - **Public** — no restrictions, freely shareable
   - **Internal** — for employees only, no external disclosure
   - **Confidential** — sensitive business or customer data, need-to-know access
   - **Restricted** — PII, PHI, PCI, trade secrets; highest controls required

   Apply classification labels to every dataset in the inventory. Classification drives access control, encryption, retention, and breach notification requirements.

4. **Assign data owners — make accountability explicit** — Every dataset must have a named Data Owner (a business executive or senior manager, not an IT role) accountable for: approving access, certifying classification, approving retention periods, and signing off on data quality thresholds. Data Stewards support owners with operational tasks. Ownership must be tracked in the data catalog and reviewed annually.

5. **Implement access control policy — enforce least-privilege for data** — Define access control rules per classification tier: who may access, under what conditions, for how long, and with what controls (MFA, VPN, DLP monitoring). Implement attribute-based access control (ABAC) or role-based access control (RBAC) enforced at the data platform layer. Require access request workflows with owner approval for Confidential and Restricted data. Run quarterly access reviews (recertification campaigns).

6. **Define data retention and deletion policy — minimize data exposure over time** — For each data domain, define retention periods based on: business need, legal obligations (GDPR Article 5(e) storage limitation, SOX 7-year financial records), and regulatory requirements. Document retention schedules in the data catalog. Implement automated deletion or archival workflows triggered by retention expiry. Test deletion by auditing that expired records are actually removed.

7. **Apply privacy-by-design principles — embed privacy into data architecture** — Implement data minimization (collect only what is needed), purpose limitation (use data only for declared purposes), and storage limitation (retain only as long as necessary) per GDPR Article 5 and NIST Privacy Framework. Require a Data Protection Impact Assessment (DPIA) before any new processing of Restricted data or high-risk analytics.

8. **Establish a records of processing activities (RoPA) — meet Article 30 obligations** — Maintain a register of all data processing activities for Confidential and Restricted data: the purpose, legal basis, data categories, retention period, recipient categories, and international transfer mechanisms. The RoPA is required for GDPR compliance and serves as the master reference for breach response and regulatory inquiries.

9. **Define a data breach response procedure — prepare before incidents occur** — Document: breach detection triggers, containment steps, impact assessment methodology, notification SLAs (GDPR: 72 hours to supervisory authority, without undue delay to data subjects), communication templates, and post-incident review process. Conduct annual tabletop exercises. Ensure the breach response team includes legal, security, communications, and data governance.

10. **Measure governance maturity and improve continuously** — Use the DAMA DMBOK2 maturity model or CMMI Data Management maturity to assess governance capability annually. Track key metrics: % of datasets with assigned owners, % with current classification, average time to fulfill access requests, number of policy exceptions, and policy violation incidents. Set annual improvement targets for each metric.

## Rules

- Every dataset containing PII, PHI, or regulated data must have a documented owner, classification label, retention schedule, and access control policy before it enters production
- Access to Restricted data must require explicit owner approval and be time-limited — no standing access for human users to production Restricted data
- Data governance policy must be reviewed and updated at least annually, and immediately after any significant regulatory change or data incident
- Retention deletion must be verifiable via audit log — "we have a policy to delete" is insufficient; deletion must be evidenced
- Data governance decisions must be recorded — all classification changes, access approvals, and policy exceptions must be logged with rationale and approver identity

## Common Mistakes

- **Governance as a compliance checkbox** — Writing policies to satisfy an auditor without implementing controls creates paper compliance that fails on first inspection. Governance must change actual data handling behavior.
- **Assigning ownership to IT** — IT can implement controls but cannot make business decisions about what data is needed, how long to retain it, or who should access it. Ownership must sit with the business.
- **Over-classifying everything as Restricted** — When too much data is Restricted, the overhead of access controls blocks legitimate work and engineers route around the policy. Classification must reflect actual sensitivity.
- **No enforcement mechanism** — Publishing a governance policy without integrating it into access control systems, data platform configurations, and CI/CD pipelines means engineers simply don't follow it.
- **Ignoring unstructured and shadow data** — Governed databases are only part of the picture. Data in spreadsheets, email attachments, Slack exports, and personal drives is often the highest-risk. The inventory must cover all data stores, not just formal databases.

## Examples

**GDPR access request workflow:** A customer submits a Subject Access Request (SAR). The data governance portal triggers a workflow that queries the RoPA to identify all systems containing that customer's data, generates a report for Data Stewards in each system to confirm and export their data, and compiles a response within the 30-day GDPR deadline. Audit log records the full chain of handling.

**Classification-driven encryption enforcement:** A policy engine checks every new dataset registered in the data catalog. Restricted datasets without encryption at rest trigger an automated remediation task to enable encryption and alert the data owner. Confidential datasets missing an owner assignment block the dataset from being made queryable in the analytical warehouse until ownership is assigned.

## When NOT to Use

- Purely internal analytical experiments on synthetic or fully anonymized data with no connection to real individuals or regulated business data
- Early-stage startups with fewer than 10 employees where the overhead of a formal governance operating model exceeds the risk being managed (implement lightweight classification and ownership instead)
