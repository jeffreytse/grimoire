---
name: audit-infra-drift
description: Use when detecting, auditing, or remediating drift between declared infrastructure configuration and actual deployed state
source: GitOps Principles (opengitops.dev); HashiCorp Terraform documentation; AWS Config documentation
tags: [infrastructure, drift, terraform, gitops, aws-config, iac, devops]
verified: true
---

# Audit Infrastructure Drift

Detect and remediate differences between infrastructure-as-code declarations and the actual state of deployed resources.

## Why This Is Best Practice

**Adopted by:** WeaveWorks (GitOps origin), AWS (Config + CloudFormation drift detection), HashiCorp (Terraform plan as drift detection)
**Impact:** Puppet's State of DevOps report found that IaC adoption reduces change failure rate by 60%; undetected drift is the #1 cause of "works in staging, fails in production" incidents.

**Why best:** Drift accumulates when engineers make manual ("click-ops") changes to production that are not reflected in the IaC source of truth. Over time, the declared state diverges from reality, making changes unpredictable and compliance verification impossible. Regular drift audits restore confidence in IaC as the authoritative source.

## Steps

1. **Run a plan/diff against current state** — Execute `terraform plan` (Terraform), `kubectl diff` (Kubernetes), or `aws cloudformation detect-stack-drift` to surface all differences between declared and actual state.
2. **Classify each drift item** — Authorized drift (intentional manual override, documented), unauthorized drift (undocumented manual change), configuration skew (environment-specific differences).
3. **Trace drift origin** — Check CloudTrail / audit logs to identify who made the manual change and when. This determines remediation urgency and process improvement.
4. **Remediate by reconciling IaC** — Either: (a) update IaC to reflect the correct desired state and apply, or (b) revert the manual change via IaC apply. Never reconcile by making more manual changes.
5. **Enforce GitOps controls** — Remove manual IAM permissions for production resource modification where possible; all changes must flow through the IaC pipeline.
6. **Schedule continuous drift detection** — Run `terraform plan` in CI on a schedule (e.g., nightly); alert on non-empty plans. Use AWS Config rules for real-time compliance checks.
7. **Document exceptions** — For authorized drift (e.g., auto-scaling changes managed by AWS), configure Terraform `ignore_changes` or AWS Config suppressions to exclude from alerts.

## Rules

- A non-empty `terraform plan` on a production environment is an incident signal — investigate before the next deploy.
- All production changes must be traceable to a git commit; manual changes are a process violation, not a solution.
- Drift in security-sensitive resources (IAM policies, security groups, KMS) triggers immediate remediation regardless of authorization status.
- Use remote state locking (S3 + DynamoDB, Terraform Cloud) to prevent concurrent state modification.

## Examples

```bash
# Terraform drift detection
terraform plan -detailed-exitcode
# Exit code 2 = changes detected (drift present)

# AWS Config drift detection
aws cloudformation detect-stack-drift --stack-name production-stack
aws cloudformation describe-stack-resource-drifts --stack-name production-stack \
  --stack-resource-drift-status-filters MODIFIED DELETED
```

## Common Mistakes

- **Reconciling drift with more manual changes** — compounds the problem; always go through IaC.
- **`terraform apply -auto-approve` without reviewing the plan** — applies unintended drift removals that may cause outages.
- **No drift detection schedule** — drift discovered at deployment time is far more disruptive than drift caught nightly.
