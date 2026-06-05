---
name: design-network-segmentation
description: Use when designing network architecture for systems that handle sensitive data, require regulatory compliance, or need to contain blast radius from breaches
source: NIST SP 800-125B network security guidelines; CIS Controls v8 Control 12; PCI DSS v4.0 network segmentation requirements
tags: [networking, security, architecture, compliance]
verified: true
---

# Design Network Segmentation

Partition networks into isolated zones so that a compromised segment cannot pivot laterally to compromise others.

## Why This Is Best Practice

**Adopted by:** PCI DSS mandates it for cardholder data environments; HIPAA requires it for ePHI isolation; DoD mandates zero-trust segmentation; every major financial institution
**Impact:** Reduces breach blast radius by 60-80% (IBM Cost of a Data Breach 2023); PCI DSS compliance requires segmentation to reduce audit scope; lateral movement limited to compromised segment
**Why best:** Flat networks allow attackers who breach any host to reach all hosts; segmentation enforces least-privilege at the network layer and contains compromise

Sources: NIST SP 800-125B (2016); CIS Controls v8 Control 12 (2021); PCI DSS v4.0 (2022)

## Steps

1. **Identify and classify assets** — Inventory all systems and classify by data sensitivity (public, internal, confidential, restricted) and function (web tier, app tier, database, management). This classification drives segment boundaries.

2. **Define security zones** — Establish standard zones: DMZ (public-facing services), application tier, database tier, management/bastion, and DevOps CI/CD. Each zone has a defined trust level and allowed traffic flows.

3. **Apply the principle of least-privilege to network flows** — Define an allowlist of required traffic flows between zones. Default-deny all other traffic. Document every allowed flow with business justification. Deny-by-default is the foundational rule.

4. **Implement VPC/VNET segmentation** — Use cloud-native constructs: separate VPCs per environment (dev, staging, prod) and per security domain. Use VPC peering or Transit Gateway only for explicitly needed cross-VPC flows. Never peer dev VPC to prod.

5. **Deploy security groups and NACLs** — Security groups: stateful, applied at instance level, source/destination specific. NACLs: stateless, applied at subnet level, for subnet-boundary enforcement. Layer both; NACLs catch security group misconfigurations.

6. **Isolate database tier** — Database subnets must have no route to the internet (no NAT gateway, no IGW). Access only from application tier security group. Deny direct admin access; route through bastion host with session recording.

7. **Implement a dedicated management zone** — All administrative access (SSH, RDP, kubectl, cloud console) flows through a hardened bastion/jump host or VPN. The management zone has no internet access. Log all sessions with tools like AWS Systems Manager Session Manager.

8. **Design micro-segmentation for critical workloads** — For PCI or highly sensitive data: implement host-based firewall rules (iptables, Windows Firewall, Calico for Kubernetes) to restrict east-west traffic within a subnet. Service mesh (Istio) provides mTLS between services.

9. **Validate with traffic flow analysis** — Use VPC Flow Logs (AWS), Cloud Flow Logs (GCP), or NSG Flow Logs (Azure) to verify only allowed traffic flows exist. Scan for unexpected cross-segment flows. Run this analysis after every architecture change.

10. **Document and review segment topology** — Maintain a network diagram showing all segments, routing, and firewall rules. Review quarterly and after major changes. Treat undocumented traffic flows as security findings.

## Rules

- Default-deny is non-negotiable; allowlist only, never blocklist-only for inter-zone traffic.
- Production and non-production networks must never be peered without a documented, reviewed exception.
- The database tier must have zero direct internet access; any route to the internet is a critical finding.
- All traffic crossing a zone boundary must traverse a firewall or security group; routing around security controls defeats segmentation.

## Common Mistakes

- **Overly broad security group rules** — `0.0.0.0/0` on port 443 in internal segments; specify exact source security groups, not CIDR ranges.
- **Flat VPC with subnets as only isolation** — subnets in the same VPC can route to each other unless NACLs or security groups block them; subnets alone are not segmentation.
- **Forgetting VPC endpoint security** — S3 VPC endpoints require endpoint policies; without them, all instances in the VPC can access any S3 bucket including those in other accounts.
- **Developer access to production databases** — direct prod DB access from dev workstations bypasses application-layer controls; always route through the application tier.

## When NOT to Use

- Single-tier applications with no sensitive data and no compliance requirements (simple static sites)
- Local development environments where segmentation overhead adds no security value
- Proof-of-concept systems with no sensitive data and a defined decommission date
