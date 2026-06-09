---
name: design-zero-trust-architecture
description: Use when designing or migrating a security architecture to the zero trust model where no user, device, or network is trusted by default
source: 'NIST SP 800-207 "Zero Trust Architecture" (2020); John Kindervag Forrester "No More Chewy Centers" (2010); Google "BeyondCorp: A New Approach to Enterprise Security" (2014); CISA Zero Trust Maturity Model (2022)'
tags: [security, zero-trust, network, identity, access-control]
verified: true
---

# Design Zero Trust Architecture

Design a security architecture where every access request is authenticated, authorized, and continuously validated — trust is never assumed based on network location.

## Why This Is Best Practice

**Adopted by:** Google (BeyondCorp, production since 2011), US Federal Government (Executive Order 14028, 2021), Microsoft (internal zero trust migration), NIST as the recommended enterprise security model

**Impact:** Google's BeyondCorp eliminated VPN dependency for 80,000+ employees while improving security posture; CISA reports zero trust adoption reduces lateral movement blast radius by up to 90% in breach scenarios; Forrester Research found organizations with zero trust suffer average breach costs 40% lower than perimeter-only models

**Why best:** Traditional perimeter security assumes internal network traffic is safe — once an attacker breaches the perimeter (via phishing, credential theft, or insider threat), they move laterally with minimal friction. Zero trust treats every packet as potentially hostile, enforcing least-privilege access at every hop, reducing the blast radius of any compromise.

Sources: NIST SP 800-207 (2020); Kindervag, J. "No More Chewy Centers: Introducing the Zero Trust Model of Information Security" Forrester (2010); Ward, R. & Beyer, B. "BeyondCorp: A New Approach to Enterprise Security" Google (2014); CISA Zero Trust Maturity Model v2.0 (2022)

## Steps

1. **Inventory assets and data flows — map everything** — Enumerate all users, devices, applications, services, and data stores. Document every data flow between components. You cannot protect what you haven't mapped. Use automated discovery tools (e.g., network scanners, cloud asset inventory) to ensure completeness.

2. **Define protect surfaces — shrink the scope** — Instead of protecting the entire network, identify the minimal Protect Surfaces: data, assets, applications, and services (DAAS) that are critical. Each protect surface gets its own micro-perimeter enforced by a policy enforcement point (PEP).

3. **Design identity as the new perimeter — enforce strong authentication everywhere** — Implement a centralized Identity Provider (IdP) with multi-factor authentication (MFA) for all users and service accounts. Apply phishing-resistant MFA (FIDO2/passkeys) for privileged access. Every identity — human and machine — must be cryptographically verifiable.

4. **Implement device trust — validate health before access** — Require device compliance checks before granting access: OS patch level, endpoint detection and response (EDR) agent presence, disk encryption, and certificate validity. Use a Mobile Device Management (MDM) or Endpoint Management platform to continuously attest device health.

5. **Apply least-privilege access — grant minimum necessary permissions** — Replace broad network access (VPN, flat subnets) with application-level access grants scoped to specific resources and time windows. Use just-in-time (JIT) provisioning for privileged roles. Review and recertify access grants quarterly.

6. **Deploy policy enforcement points — intercept and evaluate every request** — Place PEPs (API gateways, reverse proxies, service meshes) in front of all resources. Each PEP consults a Policy Decision Point (PDP) that evaluates contextual signals: identity, device posture, location, time, and behavior before allowing or denying access.

7. **Encrypt all traffic — eliminate trusted internal channels** — Enforce mutual TLS (mTLS) for service-to-service communication. Encrypt data in transit everywhere, including east-west traffic within data centers. Use certificate management automation (e.g., cert-manager, AWS Certificate Manager) to prevent certificate expiry failures.

8. **Implement continuous monitoring and behavioral analytics** — Collect logs from every PEP, IdP, and endpoint into a SIEM. Apply User and Entity Behavior Analytics (UEBA) to detect anomalous access patterns. Set alerts for impossible travel, unusual data access volumes, or off-hours privileged activity.

9. **Establish an incident response playbook for zero trust failures** — Define what happens when a PEP blocks legitimate access (fast remediation path) and when anomalous behavior triggers investigation. Test playbooks via tabletop exercises at least twice per year.

10. **Iterate using a maturity model** — Use the CISA Zero Trust Maturity Model to benchmark current state across five pillars (Identity, Devices, Networks, Applications, Data) and plan incremental improvement. Avoid big-bang migrations; prioritize highest-risk protect surfaces first.

## Rules

- Never create implicit trust based on network location — source IP alone must never grant access
- Every access decision must be logged with full context (who, what device, when, from where, result)
- All secrets, certificates, and credentials must be rotated automatically — no static long-lived credentials
- Default deny all; explicitly allow only what is needed and verified
- Zero trust architecture must cover both external users and internal service-to-service calls

## Common Mistakes

- **Treating zero trust as a product** — Zero trust is an architecture strategy, not a single vendor solution. Buying a "zero trust" product without redesigning access policies and identity infrastructure achieves nothing.
- **Skipping device trust** — Implementing strong user MFA while ignoring device health creates a gap: a stolen compliant device can still be used by an attacker after initial authentication.
- **Ignoring east-west traffic** — Focusing only on north-south (ingress) traffic while leaving internal service calls unencrypted and unauthenticated defeats the purpose. Lateral movement happens east-west.
- **Over-permissive initial policies** — Teams often set broad allow rules during migration "just to keep things working" and never tighten them. Instrument and audit from day one.
- **Neglecting legacy systems** — Applications that cannot support modern auth protocols need a wrapper (reverse proxy or API gateway) — never exempted from zero trust controls.

## Examples

**Google BeyondCorp pattern:** Instead of a VPN, employees access internal apps through a cloud-hosted access proxy. The proxy validates identity (corporate SSO + hardware key), device certificate (issued to managed devices), and checks a dynamic access policy based on user role and resource sensitivity before proxying the request.

**Kubernetes zero trust with service mesh:** Deploy Istio or Linkerd as a service mesh. All pod-to-pod communication uses mTLS with certificates rotated by cert-manager. Authorization policies at the mesh level enforce that only the `orders-service` service account can call `payment-service`, and only on specific gRPC methods.

## When NOT to Use

- Small isolated systems with no external attack surface and no sensitive data (cost-benefit doesn't justify the complexity)
- Emergency break-glass scenarios where security controls are temporarily bypassed to restore critical services — document and audit immediately after
