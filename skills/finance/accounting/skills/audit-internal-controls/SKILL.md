---
name: audit-internal-controls
description: Use when reviewing, designing, or stress-testing internal financial controls — e.g., "how do I prevent fraud?", "what internal controls should we have?", "SOX compliance", "reviewing controls before audit", "segregation of duties"
source: COSO Internal Control Framework (2013); Sarbanes-Oxley Act Section 404; AICPA auditing standards; Association of Certified Fraud Examiners (ACFE) "Report to the Nations" (2022)
tags: [finance, accounting, internal-controls, SOX, fraud-prevention, audit, compliance, segregation-of-duties]
verified: true
---

# Audit Internal Controls

Systematically review financial controls using the COSO framework to identify control gaps, prevent fraud, and ensure reliable financial reporting.

## Why This Is Best Practice

**Adopted by:** The COSO Internal Control Framework (Committee of Sponsoring Organizations) is the standard used by all public companies for SOX Section 404 compliance. The Big 4 accounting firms (Deloitte, PwC, KPMG, EY) use COSO as their audit methodology. The ACFE (Association of Certified Fraud Examiners) trains 90,000+ professionals using this framework.
**Impact:** ACFE "Report to the Nations" (2022) found that organizations lose 5% of revenue annually to fraud (median $125,000 per incident). Companies with strong internal controls detected fraud in a median of 12 months; those without detected it in 24 months — at double the cost. Public companies with material weaknesses in controls trade at a 5–10% discount to peers (Securities and Exchange Commission enforcement data).
**Why best:** Most fraud is prevented not by catching it after the fact, but by designing controls that make it difficult or impossible to commit. Segregation of duties — the single most powerful control — means no one person can both initiate and approve a transaction, eliminating the most common fraud vector. Control frameworks also prevent honest errors that cause restatements, the second-most-expensive accounting event after fraud.

## Steps

1. **Map the five COSO components** — Internal controls must address:
   - **Control Environment**: tone at the top, ethics culture, board oversight, HR policies.
   - **Risk Assessment**: what could go wrong? (financial misstatement, fraud, error, compliance breach)
   - **Control Activities**: the specific controls that prevent or detect errors/fraud.
   - **Information & Communication**: accurate financial data flows to decision makers.
   - **Monitoring**: ongoing assessment of whether controls are working.

2. **Inventory financial processes** — List all significant financial processes: revenue recognition, AR/AP, payroll, cash management, fixed assets, expense reporting, period-end close. Each is a risk area requiring control evaluation.

3. **Apply segregation of duties (SOD) analysis** — For each process, confirm no single person has all three of: (a) authorization, (b) custody, (c) record-keeping. Example: the person who approves vendor invoices should not also process payments AND should not reconcile the bank statement. Document who performs each role; flag any person with 2 or more of the three functions.

4. **Test key controls** — For each critical control:
   - **Preventive control** (stops errors before they occur): sample 25–50 transactions; verify control operated as designed every time.
   - **Detective control** (catches errors after the fact): verify the control fires when it should. Example: bank reconciliation should catch unauthorized payments.
   - **IT general controls**: access controls (who can enter/approve what), change management (who can modify financial systems), backup and recovery.

5. **Identify control gaps** — Rank gaps by risk (likelihood × financial impact):
   - **Material weakness**: significant deficiency that could lead to a material misstatement. Requires immediate remediation; must be disclosed in public company filings.
   - **Significant deficiency**: less severe; requires prompt attention.
   - **Control deficiency**: minor; remediate in normal course.

6. **Assess fraud risk** — Apply the Fraud Triangle: Opportunity (weak controls), Pressure (financial stress, quotas), Rationalization (ethical grey zones). High-risk combinations: cash handling + financial pressure + weak oversight = elevated fraud risk.

7. **Remediate prioritized gaps** — For each gap: define the control, assign ownership, set implementation deadline, test after remediation. Common remediations: add approval layer, add reconciliation step, restrict system access, implement dual-signature policy.

8. **Document for audit** — For each control: description, control owner, frequency (daily/monthly), evidence of operation (signed approval, system log, reconciliation file). Documentation is required for SOX; valuable for any audit.

## Rules

- Segregation of duties is the most cost-effective control — design processes around it before adding other controls.
- Every cash disbursement > $X (set your threshold, typically $5,000–$25,000) requires dual approval — no exceptions.
- Access controls in financial systems must be role-based: read-only for report viewers; entry-only for AP clerks; approval-only for managers. No admin access for operating staff.
- Controls must be tested, not assumed — a control on paper that no one follows is a documentation risk, not a real control.

## Examples

**Small company (50 employees) control gap assessment:**
Finding: the same AP clerk enters invoices, approves them, and processes the ACH payment. No one else reviews.
Risk: complete fraud vector — one person can create fictitious vendors and pay them without detection.
Remediation: AP clerk enters only; finance manager approves all payments > $500; controller reconciles bank weekly. Estimated fraud prevention value: eliminates 95% of AP fraud risk.
Implementation cost: 2 hours/week of manager time.

## Common Mistakes

- **Paper controls that aren't followed** — "The manager reviews all invoices" is not a control if the manager rubber-stamps everything in batch. Effective controls have evidence: signed documents, system approvals with timestamps.
- **IT access not reviewed after role changes** — Former employees with active system access are a major fraud risk; access must be revoked on the day of termination and reviewed quarterly.
- **Relying on trust instead of controls** — "We trust everyone here" is not an internal control. Most frauds are committed by long-tenured trusted employees (ACFE data: median fraud perpetrator has 7 years of tenure).

---

> **Finance disclaimer:** This skill encodes professional best practices for educational purposes. It is not financial advice. Consult a licensed financial advisor before making investment decisions.
