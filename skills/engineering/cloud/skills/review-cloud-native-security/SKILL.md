---
name: review-cloud-native-security
description: Use when auditing a cloud-native application or infrastructure — systematically checking all 10 OWASP Cloud-Native Application Security Top 10 vulnerability classes with AWS CLI, kubectl, and Terraform commands.
source: 'OWASP Cloud-Native Application Security Top 10 (owasp.org/www-project-cloud-native-application-security-top-10/); CIS AWS/GCP/Azure Benchmarks; AWS Security Hub; Prowler documentation'
tags: [security, owasp, cloud, audit, kubernetes, aws, terraform, developer]
---

# Review Cloud-Native Security

Audit cloud-native applications against the OWASP Cloud-Native Application Security Top 10 using Prowler, AWS Security Hub, and kubectl — covering C1 through C10 with specific test commands and remediation steps.

## Why This Is Best Practice

**Adopted by:** OWASP Cloud-Native Application Security Top 10 is the authoritative vulnerability taxonomy for cloud-native workloads. Prowler (open-source, used by AWS, Stripe, and security consulting firms) automates checks for C1–C10 equivalent categories. AWS Security Hub consolidates findings from GuardDuty, Inspector, Macie, and third-party tools against the AWS Foundational Security Best Practices standard. CIS Benchmarks for AWS (v3.0), GCP (v2.0), and Azure (v2.0) provide detailed test procedures for each vulnerability class.
**Impact:** Cloud Security Alliance's 2022 "Top Threats to Cloud Computing" maps directly to C1–C10. Palo Alto Unit 42's "2023 Cloud-Native Security Report" found 76% of organizations have at least one C1–C10 class vulnerability in production. Organizations using structured cloud security review checklists remediate critical findings 40% faster than those using ad-hoc review (Gartner Cloud Security report, 2022). The OWASP Cloud-Native Top 10 provides completeness guarantees that neither tool-based scanning nor manual review alone achieves.
**Why best:** Cloud-native environments combine container, orchestration, and cloud service vulnerabilities — a reviewer focused on Kubernetes misconfigurations may miss S3 bucket exposure, IAM over-permissioning, or logging gaps. C1–C10 provides a complete cross-layer checklist covering all three tiers of the cloud-native stack.

Sources: OWASP Cloud-Native Application Security Top 10; CIS AWS Foundations Benchmark v3.0; Prowler documentation; Palo Alto Unit 42 Cloud-Native Security Report (2023)

## Steps

### Pre-Audit Setup

```bash
# Install Prowler (comprehensive AWS security scanner)
pip install prowler

# Run Prowler against all CIS and AWS Foundational Security checks
prowler aws --checks cis_level2_1_1 --output-formats json,html

# Or target specific services
prowler aws --services ec2 iam s3 guardduty cloudtrail

# For GCP
prowler gcp --project-id my-project
```

### C1 — Insecure Cloud, Container, and Orchestration Configuration

```bash
# AWS: check for public S3 buckets
aws s3api list-buckets --query 'Buckets[*].Name' --output text | \
  xargs -I{} aws s3api get-bucket-acl --bucket {} 2>/dev/null | \
  grep -i "AllUsers\|AuthenticatedUsers"

# Check Security Group open to internet
aws ec2 describe-security-groups \
  --filters "Name=ip-permission.cidr,Values=0.0.0.0/0" \
  --query 'SecurityGroups[?IpPermissions[?IpRanges[?CidrIp==`0.0.0.0/0`]] && IpPermissions[?FromPort!=`443` && FromPort!=`80`]].GroupId'

# Kubernetes: check Pod Security Standards
kubectl get namespaces -o json | \
  jq '.items[] | {name: .metadata.name, labels: .metadata.labels} |
      select(.labels | has("pod-security.kubernetes.io/enforce") | not) |
      .name'
```

**Fix:** Block public S3 access, restrict security groups, apply `restricted` Pod Security Standards.

### C2 — Injection Flaws (App Layer, Cloud Events, Cloud Services)

```bash
# Check Lambda functions for environment variable injection risk
aws lambda list-functions --query 'Functions[*].FunctionName' --output text | \
  xargs -I{} aws lambda get-function-configuration --function-name {} \
  --query 'Environment.Variables' 2>/dev/null | grep -v "^null"
```

**Check:**
- [ ] Lambda event inputs validated (prevent SLS-1 equivalent)
- [ ] No user input interpolated into AWS CLI commands or SDK calls
- [ ] SNS/SQS message bodies treated as untrusted input

### C3 — Improper Authentication and Authorization

```bash
# Check for Lambda functions with public URLs (no auth)
aws lambda list-function-url-configs \
  --query 'FunctionUrlConfigs[?AuthType==`NONE`].FunctionArn'

# Check for API Gateway stages without auth
aws apigateway get-rest-apis --query 'items[*].id' --output text | \
  xargs -I{} aws apigateway get-resources --rest-api-id {} \
  --query 'items[?resourceMethods.GET.authorizationType==`NONE`].path'
```

**Fix:** Require IAM or Cognito auth on all Lambda URLs and API Gateway endpoints.

### C4 — CI/CD Pipeline and Software Supply Chain Flaws

```bash
# Check for GitHub Actions using unpinned action versions
grep -r "uses: " .github/workflows/ | grep -v "@[a-f0-9]\{40\}" | grep -v "^#"
# Any result = unpinned action (should be @sha256)

# Check for Terraform providers without version pins
grep -r "required_providers" . -A 20 | grep -v "version"
```

**Fix:** Pin all GitHub Actions to commit SHA. Pin Terraform provider versions. Enable code signing.

### C5 — Insecure Secrets Storage

```bash
# Check for secrets in environment variables
aws lambda list-functions --output text --query 'Functions[*].FunctionName' | \
  xargs -I{} aws lambda get-function-configuration --function-name {} \
  --query 'Environment.Variables' | \
  python3 -c "import sys,json; d=json.load(sys.stdin); [print(k) for k in d if any(s in k.upper() for s in ['PASSWORD','SECRET','KEY','TOKEN'])]"

# Check Kubernetes secrets for plaintext credentials
kubectl get secrets --all-namespaces -o json | \
  jq '.items[] | .data // {} | to_entries[] | .value | @base64d' 2>/dev/null | \
  grep -E "password|secret|key" -i
```

**Fix:** Move secrets to Secrets Manager/SSM Parameter Store. Use IRSA instead of credential env vars.

### C6 — Over-Permissioned Resources or Insecure IAM

```bash
# Find IAM roles with admin access
aws iam list-roles --query 'Roles[*].RoleName' --output text | \
  xargs -I{} aws iam list-attached-role-policies --role-name {} \
  --query 'AttachedPolicies[?PolicyArn==`arn:aws:iam::aws:policy/AdministratorAccess`]'

# Find IAM policies with wildcards
aws iam list-policies --scope Local --query 'Policies[*].Arn' --output text | \
  xargs -I{} aws iam get-policy-version --policy-arn {} --version-id v1 \
  --query 'PolicyVersion.Document.Statement[?Effect==`Allow` && Action==`*`]'

# Run IAM Access Analyzer
aws accessanalyzer list-findings --analyzer-arn $(aws accessanalyzer list-analyzers \
  --query 'analyzers[0].arn' --output text) --filter '{"resourceType": {"eq": ["AWS::IAM::Role"]}}'
```

**Fix:** Replace AdministratorAccess with least-privilege policies. Remove wildcard actions.

### C7 — Insecure Network Controls

```bash
# Check VPC flow logs disabled
aws ec2 describe-vpcs --query 'Vpcs[*].VpcId' --output text | \
  xargs -I{} aws ec2 describe-flow-logs \
  --filter "Name=resource-id,Values={}" \
  --query 'FlowLogs[0].FlowLogStatus' 2>/dev/null | grep -v "ACTIVE"

# Check for default VPC in use
aws ec2 describe-vpcs \
  --filters "Name=isDefault,Values=true" \
  --query 'Vpcs[?State==`available`].VpcId'
```

**Fix:** Enable VPC flow logs. Delete default VPC. Implement three-tier subnet architecture.

### C8 — Using Components with Known Vulnerabilities

```bash
# AWS Inspector v2 — scan EC2 and Lambda for CVEs
aws inspector2 enable --resource-types EC2 ECR LAMBDA

aws inspector2 list-findings \
  --filter-criteria '{"findingSeverity": [{"comparison": "EQUALS", "value": "CRITICAL"}]}' \
  --query 'findings[*].title'
```

**Fix:** Enable Inspector. Patch HIGH/CRITICAL CVEs within 30 days per org policy.

### C9 — Improper Assets Management

```bash
# Find untagged EC2 instances (untracked resources)
aws ec2 describe-instances \
  --query 'Reservations[*].Instances[?!not_null(Tags[?Key==`Environment`])].[InstanceId]' \
  --output text

# Find old/unused Lambda functions (last invocation > 90 days)
aws lambda list-functions --query 'Functions[*].FunctionName' --output text | \
  xargs -I{} aws cloudwatch get-metric-statistics \
  --namespace AWS/Lambda --metric-name Invocations \
  --dimensions Name=FunctionName,Value={} \
  --start-time $(date -d "90 days ago" +%Y-%m-%dT%H:%M:%S) \
  --end-time $(date +%Y-%m-%dT%H:%M:%S) \
  --period 7776000 --statistics Sum
```

### C10 — Insufficient Logging and Monitoring

```bash
# Check CloudTrail enabled in all regions
aws cloudtrail describe-trails --include-shadow-trails \
  --query 'trailList[?IsMultiRegionTrail==`true` && HomeRegion==`us-east-1`].Name'

# Check GuardDuty enabled
aws guardduty list-detectors --query 'DetectorIds'

# Check for CloudTrail log file validation
aws cloudtrail describe-trails \
  --query 'trailList[?LogFileValidationEnabled!=`true`].Name'
```

**Fix:** Enable multi-region CloudTrail with log file validation. Enable GuardDuty in all regions.

### Audit Summary Checklist

```
□ C1  Configuration — no public S3, no open SGs, restricted Pod Security
□ C2  Injection — cloud event inputs validated, no input in SDK calls
□ C3  Authentication — no public Lambda URLs/APIs without auth, IRSA used
□ C4  Supply chain — GitHub Actions pinned to SHA, provider versions locked
□ C5  Secrets — no credentials in env vars, Secrets Manager used
□ C6  IAM — no wildcard actions, no AdministratorAccess on workloads
□ C7  Network — VPC flow logs enabled, no default VPC, three-tier subnets
□ C8  Vulnerabilities — Inspector enabled, no unpatched CRITICAL CVEs
□ C9  Asset management — all resources tagged, unused resources cleaned
□ C10 Logging — multi-region CloudTrail, GuardDuty, Security Hub enabled
```

## Rules

- Run Prowler quarterly and after major infrastructure changes — configuration drift is continuous.
- C6 (IAM) and C5 (secrets) findings are the highest-priority classes — exploit them first in a penetration test, remediate them first in a security review.
- C10 (logging) gaps mean other finding classes are invisible — fix logging before investigating other findings.

## Common Mistakes

- **Running audit against dev account, not production** — production and dev accounts often have different compliance states; always audit production directly.
- **Closing findings without root cause analysis** — "we enabled GuardDuty" without investigating why it was disabled doesn't prevent recurrence; find and fix the cause.
- **Treating Prowler PASS as secure** — Prowler covers CIS benchmarks but not business-logic authorization or application-layer injection; both tool and manual review are needed.
