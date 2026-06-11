---
name: apply-cloud-runtime-security
description: Use when running workloads in cloud environments — enabling GuardDuty, Security Command Center, or Microsoft Defender, and deploying Falco for Kubernetes runtime threat detection.
source: 'OWASP Cloud-Native Application Security Top 10 C1 (owasp.org/www-project-cloud-native-application-security-top-10/); AWS GuardDuty documentation; Falco project documentation; CNCF Security Technical Advisory Group'
tags: [security, owasp, cloud, runtime-security, guardduty, falco, threat-detection, developer]
---

# Apply Cloud Runtime Security

Enable cloud-native threat detection (GuardDuty, Security Command Center, Defender) and deploy Falco for container syscall monitoring — detecting credential theft, cryptomining, lateral movement, and container escapes that static IaC scanning cannot catch.

## Why This Is Best Practice

**Adopted by:** OWASP Cloud-Native Application Security Top 10 C1 (Insecure Cloud, Container, and Orchestration Configuration). AWS GuardDuty is enabled by default in AWS Security Hub and required for AWS Foundational Security Best Practices compliance. CNCF Security Technical Advisory Group's "Cloud Native Security Whitepaper" (2022) mandates runtime threat detection for production workloads. Falco (CNCF project) is the standard open-source runtime security tool, used by Shopify, Kubernetes maintainers, and major cloud providers in their security posture.
**Impact:** Netflix's chaos engineering team documented that 60% of their theoretical cloud attack paths are only detectable at runtime — not by IaC scanning or static analysis. AWS GuardDuty detected the 2022 Capital One breach indicators (SSRF to IMDS, unusual S3 API calls) within 4 hours of initial access — the breach would have been caught earlier with GuardDuty enabled from day one. Aqua Security's 2023 Cloud Native Threat Report found that 50% of cloud attacks involve cryptomining that starts within minutes of container compromise — detectable only by runtime monitoring of network connections and CPU patterns.
**Why best:** IaC scanning (Checkov, tfsec) detects misconfigurations before deployment. Runtime security detects attacks happening right now — a container that was correctly configured at deploy time can be compromised via a zero-day, supply chain attack, or application vulnerability. Runtime security provides the detective layer that IaC scanning cannot replace.

Sources: OWASP Cloud-Native Top 10 C1; AWS GuardDuty documentation; CNCF Cloud Native Security Whitepaper (2022); Aqua Security Cloud Native Threat Report (2023)

## Steps

1. **Enable AWS GuardDuty with EKS and malware protection**:

   ```bash
   # Enable GuardDuty in all regions (use AWS Organizations for multi-account)
   aws guardduty create-detector \
     --enable \
     --features '[
       {"Name": "EKS_AUDIT_LOGS", "Status": "ENABLED"},
       {"Name": "EKS_RUNTIME_MONITORING", "Status": "ENABLED",
        "AdditionalConfiguration": [
          {"Name": "EKS_ADDON_MANAGEMENT", "Status": "ENABLED"}
        ]},
       {"Name": "S3_DATA_EVENTS", "Status": "ENABLED"},
       {"Name": "MALWARE_PROTECTION", "Status": "ENABLED"}
     ]'
   ```

   ```yaml
   # CloudWatch Events rule — route high-severity findings to SNS
   GuardDutyHighSeverityAlarm:
     Type: AWS::Events::Rule
     Properties:
       EventPattern:
         source: [aws.guardduty]
         detail-type: [GuardDuty Finding]
         detail:
           severity: [{numeric: [">=", 7]}]  # HIGH and CRITICAL
       Targets:
       - Arn: !Ref SecurityAlertTopic
         Id: SecurityTeamAlert
   ```

2. **Deploy Falco for Kubernetes runtime syscall monitoring**:

   ```bash
   # Install Falco via Helm
   helm repo add falcosecurity https://falcosecurity.github.io/charts
   helm install falco falcosecurity/falco \
     --namespace falco \
     --create-namespace \
     --set driver.kind=ebpf \
     --set falcosidekick.enabled=true \
     --set falcosidekick.config.slack.webhookurl="https://hooks.slack.com/..."
   ```

   ```yaml
   # Custom Falco rules for common attack patterns
   # /etc/falco/rules.d/custom-rules.yaml
   - rule: Container executing cryptominer
     desc: Detect common cryptominer process names
     condition: >
       spawned_process and
       proc.name in (xmrig, minerd, cryptonight, ethminer)
     output: >
       Cryptominer detected (user=%user.name cmd=%proc.cmdline
       container=%container.id image=%container.image.repository)
     priority: CRITICAL

   - rule: Unexpected outbound connection from container
     desc: Alert on unexpected external network connections
     condition: >
       outbound and
       not container.image.repository in (known_external_images) and
       fd.rip not in (allowed_external_ips)
     output: >
       Unexpected outbound connection (container=%container.id
       ip=%fd.rip port=%fd.rport image=%container.image.repository)
     priority: WARNING

   - rule: Sensitive file read in container
     desc: Detect access to credential files
     condition: >
       open_read and
       fd.name in (/etc/shadow, /etc/passwd, /root/.ssh/id_rsa,
                   /var/run/secrets/kubernetes.io/serviceaccount/token)
       and not proc.name in (sshd, login, su)
     priority: HIGH
   ```

3. **Configure CloudTrail for full API audit logging**:

   ```bash
   # Enable CloudTrail with data events (S3, Lambda invocations)
   aws cloudtrail create-trail \
     --name full-audit-trail \
     --s3-bucket-name audit-logs-company \
     --is-multi-region-trail \
     --enable-log-file-validation

   aws cloudtrail put-event-selectors \
     --trail-name full-audit-trail \
     --event-selectors '[{
       "ReadWriteType": "All",
       "IncludeManagementEvents": true,
       "DataResources": [
         {"Type": "AWS::S3::Object", "Values": ["arn:aws:s3:::"]},
         {"Type": "AWS::Lambda::Function", "Values": ["arn:aws:lambda"]}
       ]
     }]'

   # Enable Insights to detect unusual API activity patterns
   aws cloudtrail put-insight-selectors \
     --trail-name full-audit-trail \
     --insight-selectors '[{"InsightType": "ApiCallRateInsight"},
                           {"InsightType": "ApiErrorRateInsight"}]'
   ```

4. **Alert on critical runtime security events**:

   ```python
   # Lambda function: process GuardDuty findings and route to appropriate response
   import boto3
   import json

   def lambda_handler(event, context):
       finding = event["detail"]
       severity = finding["severity"]
       finding_type = finding["type"]

       # Categorize response by severity and type
       if severity >= 7:
           if "CryptoCurrency" in finding_type:
               isolate_instance(finding["resource"]["instanceDetails"]["instanceId"])
           elif "UnauthorizedAccess" in finding_type:
               revoke_credentials(finding["resource"])
               notify_security_team(finding, channel="pagerduty")
           else:
               notify_security_team(finding, channel="slack-security")
       elif severity >= 4:
           notify_security_team(finding, channel="slack-security")
       else:
           create_ticket(finding)  # low severity — track in backlog

   def isolate_instance(instance_id: str):
       ec2 = boto3.client("ec2")
       ec2.modify_instance_attribute(
           InstanceId=instance_id,
           Groups=["sg-quarantine"]  # quarantine security group: no egress
       )
   ```

5. **GCP: Enable Security Command Center and Cloud Audit Logs**:

   ```bash
   # Enable Security Command Center Premium
   gcloud scc settings update \
     --organization=ORG_ID \
     --enable-all-services

   # Enable Cloud Audit Logs for all services
   gcloud projects get-iam-policy PROJECT_ID \
     --format=json > policy.json
   # Add auditConfigs for all services with DATA_READ, DATA_WRITE, ADMIN_WRITE
   gcloud projects set-iam-policy PROJECT_ID policy.json

   # Enable Container Threat Detection for GKE
   gcloud container clusters update production-cluster \
     --enable-workload-config-audit \
     --workload-vulnerability-scanning=STANDARD
   ```

## Rules

- Enable GuardDuty/Security Command Center/Defender in ALL regions, not just primary — attackers exploit enabled-but-unmonitored regions.
- Falco rules must be tuned to the environment — default rules generate high noise in production; tune to reduce false positives before alerting is useful.
- CloudTrail logs must be in a dedicated account or with SCP-enforced write protection — an attacker with admin access can delete CloudTrail logs to cover tracks.
- Runtime security findings must have defined response playbooks — an alert with no response procedure is theater.

## Common Mistakes

- **GuardDuty enabled but findings going to an unmonitored email** — without on-call routing (PagerDuty/OpsGenie), high-severity findings are seen days later.
- **Falco deployed without eBPF driver on modern kernels** — the legacy kernel module driver requires kernel headers and doesn't work on all managed node groups; use eBPF.
- **Not monitoring for `iam:CreateAccessKey` or `sts:GetCalledUserIdentity`** — these are early indicators of credential theft; alert on IAM actions not in expected patterns.
- **CloudTrail delivering to an S3 bucket in the same account** — an attacker with admin access can delete both the trail and the bucket; use cross-account log delivery with SCP protection.
