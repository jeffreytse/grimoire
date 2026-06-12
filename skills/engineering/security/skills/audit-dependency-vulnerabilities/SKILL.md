---
name: audit-dependency-vulnerabilities
description: Use when auditing third-party dependencies for known vulnerabilities, license issues, or supply chain risks
source: OWASP A06:2021 – Vulnerable and Outdated Components; NIST NVD (nvd.nist.gov); Snyk (snyk.io); GitHub Advisory Database
tags: [security, dependencies, vulnerabilities, owasp, supply-chain, cve, snyk]
verified: true
---

# Audit Dependency Vulnerabilities

Systematically scan, triage, and remediate known vulnerabilities in third-party dependencies to reduce supply chain risk.

## Why This Is Best Practice

**Adopted by:** GitHub (Dependabot), Snyk (used by Google, Salesforce, Adobe), OWASP top 10 mandates
**Impact:** OWASP A06:2021 moved Vulnerable Components from #9 to #6; the Log4Shell vulnerability (CVE-2021-44228) affected 93% of enterprise cloud environments and cost organizations an average of $4.6M to remediate (IBM 2022).

**Why best:** Third-party code is the majority of most applications. Vulnerabilities in transitive dependencies are as exploitable as direct ones — Log4Shell was a transitive dependency in most affected systems. Automated scanning must be continuous, not a one-time audit.

## Steps

1. **Generate a Software Bill of Materials (SBOM)** — Produce a dependency manifest: `npm list --all`, `pip freeze`, `mvn dependency:tree`, or use `syft` for container images.
2. **Run automated scanning** — Execute `npm audit`, `pip-audit`, `snyk test`, or `trivy` against the SBOM. Integrate into CI so every PR is scanned.
3. **Triage by CVSS score and exploitability** — Critical (CVSS 9+): fix immediately. High (7-8.9): fix within 7 days. Medium: next sprint. Low: backlog. Check if the vulnerable code path is actually reachable in your usage.
4. **Check transitive dependencies** — Direct deps rarely expose the full risk; audit the full dependency tree including transitive deps.
5. **Apply fixes** — Update to patched versions; if no patch exists, assess: remove the dependency, replace it, or apply a virtual patch (WAF rule, input validation).
6. **Enable automated PRs** — Configure Dependabot or Renovate to auto-open PRs for patch updates; auto-merge patch-level updates that pass CI.
7. **Monitor continuously** — Subscribe to GitHub Advisory Database or Snyk alerts for your dependency set; new CVEs emerge after your last scan.

## Rules

- Never dismiss a critical CVE without documented justification (the vulnerable function is not called, etc.).
- Treat `npm audit --production` and `--all` separately — dev dependency vulnerabilities are lower risk but not zero risk.
- Pin dependency versions in production; use lockfiles (package-lock.json, Pipfile.lock); verify checksums.
- Audit container base images separately — OS-level vulnerabilities are not caught by language package scanners.

## Examples

```bash
# Quick audit
npm audit --audit-level=high
snyk test --severity-threshold=high

# SBOM generation
syft . -o spdx-json > sbom.json
grype sbom:sbom.json
```

Triage note: `lodash@4.17.20` — CVE-2021-23337 (prototype pollution, CVSS 7.2). Assessment: `lodash.template` not used in this codebase → Medium priority; update in next sprint.

## Common Mistakes

- **Auditing only direct dependencies** — most exploited CVEs are in transitive deps.
- **Ignoring "low" severity** — chained with other vulnerabilities, low-severity issues enable escalation.
- **No CI integration** — weekly manual audits miss vulnerabilities introduced between scans.
