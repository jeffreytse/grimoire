---
name: audit-network-security
description: Use when assessing network security posture before a production launch, after a security incident, during compliance assessment, or on a periodic security review cycle
source: NIST SP 800-115 Technical Guide to Information Security Testing; CIS Controls v8; PCI DSS v4.0 assessment guidelines
tags: [networking, security, audit, compliance]
verified: true
---

# Audit Network Security

Systematically assess network security controls to identify misconfigurations, unauthorized access paths, and compliance gaps before attackers do.

## Why This Is Best Practice

**Adopted by:** PCI DSS requires quarterly external scans and annual penetration tests; SOC 2 requires periodic network security review; FedRAMP mandates continuous monitoring
**Impact:** Proactive audits find 3x more vulnerabilities than reactive incident investigation (Verizon DBIR 2023); average dwell time before detection is 204 days — audits compress this to zero
**Why best:** Network misconfigurations are the #2 cause of breaches (DBIR 2023); automated scanning cannot replace structured audit methodology that tests intent vs. implementation

Sources: NIST SP 800-115 (2008); CIS Controls v8 Control 12-13 (2021); PCI DSS v4.0 Requirement 11

## Steps

1. **Define audit scope** — Specify IP ranges, cloud accounts, VPCs, and services in scope. Document what is explicitly out of scope and why. Obtain written authorization before scanning; unauthorized scans violate cloud provider ToS and may trigger security alerts.

2. **Enumerate network topology** — Collect: network diagrams, routing tables, firewall rule sets, security group configurations, NACL policies. Compare documented architecture to actual configuration — discrepancies are immediate findings.

3. **Scan for open ports and services** — Run nmap against in-scope hosts: `nmap -sV -sC -p- --open <target>`. Identify services running on non-standard ports. Flag any management services (SSH port 22, RDP 3389, database ports) exposed to the internet.

4. **Review firewall and security group rules** — Audit every allow rule: source, destination, port, protocol. Flag rules with source `0.0.0.0/0` (any internet) on non-80/443 ports as critical findings. Check for redundant, overly broad, and stale rules (no traffic in 90 days).

5. **Test network segmentation** — From each segment, attempt to reach hosts in other segments that should be unreachable. Use tools like hping3 or netcat to test specific port connectivity. Verify database subnets are unreachable from internet-facing segments.

6. **Assess DNS configuration** — Check for DNS zone transfer (AXFR) enabled on authoritative servers. Review DNS records for dangling CNAME entries pointing to deprovisioned cloud resources (subdomain takeover risk). Verify DNSSEC where required.

7. **Review TLS/certificate posture** — Use testssl.sh or sslyze to assess all HTTPS endpoints: cipher suites (flag RC4, DES, 3DES, export ciphers), TLS version (require TLS 1.2+, flag SSLv3/TLS 1.0/1.1), certificate expiry, and HSTS header presence.

8. **Analyze flow logs for anomalies** — Review VPC Flow Logs / NSG Flow Logs for: unexpected cross-segment traffic, connections to known-malicious IPs (threat intel feed), port scanning patterns, and data exfiltration volumes (unusually high egress). Query 30-day window minimum.

9. **Check network device and instance configurations** — Verify: logging enabled on all network devices, NTP synchronized, management plane access restricted, default credentials changed, firmware/AMI patched within 30 days of CVE release for critical findings.

10. **Produce findings report** — Document each finding with: severity (Critical/High/Medium/Low), affected asset, description, evidence, remediation recommendation, and remediation timeline. Prioritize by risk: internet-exposed management services are Critical regardless of other controls.

## Rules

- Never scan without written authorization; include scan windows and excluded systems in the authorization document.
- Findings without evidence (screenshots, tool output, log snippets) are not findings — they will be disputed.
- Prioritize internet-exposed findings above all others; internal findings matter but external exposure is immediate risk.
- Retest every finding after remediation; remediation without verification is incomplete.

## Common Mistakes

- **Stopping at automated scan output** — vulnerability scanners produce false positives and miss logic flaws; always validate findings manually.
- **Skipping flow log analysis** — static configuration review misses active malicious traffic; live traffic analysis reveals what configuration reviews cannot.
- **Auditing only the prod account** — attackers pivot from dev/staging to prod via shared credentials, VPC peering, or CI/CD pipelines; audit all environments.
- **No remediation tracking** — findings without SLA-tracked remediation reappear in every subsequent audit; integrate findings into the team's issue tracker.

## When NOT to Use

- Production systems during peak traffic hours where nmap scans cause noticeable load (schedule during maintenance windows)
- Systems where penetration testing is legally prohibited by the data processing agreement with a third-party vendor
