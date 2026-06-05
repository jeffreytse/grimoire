---
name: design-supply-chain-security
description: Use when designing software supply chain security controls including SBOM generation, artifact signing, dependency management, and build pipeline integrity
source: NIST SP 800-161r1 "Cybersecurity Supply Chain Risk Management" (2022); SLSA Framework v1.0 Google (2023); CISA "Software Bill of Materials" guidance (2023); US Executive Order 14028 on Improving the Nation's Cybersecurity (2021)
tags: [security, supply-chain, sbom, signing, dependencies, slsa]
verified: true
---

# Design Supply Chain Security

Establish end-to-end integrity controls for software from source code through deployment so that every artifact's provenance is verified and tampering is detectable.

## Why This Is Best Practice

**Adopted by:** US Federal agencies (mandated by EO 14028 and OMB M-22-18), Google (internal SLSA framework, open-sourced 2021), GitHub (artifact attestations GA 2024), Linux Foundation (OpenSSF Scorecard, sigstore project)

**Impact:** The SolarWinds attack (2020) compromised 18,000 organizations through a single build pipeline; the Log4Shell vulnerability demonstrated how transitive dependencies become attack vectors at scale. CISA found that organizations with mature SBOM practices reduced mean time to identify vulnerable components by 72% following a zero-day disclosure. SLSA level 3 provenance makes build tampering detectable without access to the original build system.

**Why best:** Modern software is composed — 80%+ of code in a typical application comes from open-source dependencies and third-party tooling. Attackers target the supply chain because compromising one upstream package yields access to thousands of downstream consumers. Systematic controls (signing, SBOM, provenance) make attacks detectable before they reach production.

Sources: NIST SP 800-161r1 (2022); Google SLSA Framework v1.0 slsa.dev (2023); CISA "Framing Software Component Transparency: Establishing a Common Software Bill of Materials" (2023); Executive Order 14028 (May 2021); OpenSSF Scorecard (2023)

## Steps

1. **Audit current dependency graph — know what you include** — Run dependency analysis across all repositories to enumerate direct and transitive dependencies with pinned versions. Use tools like `syft`, `cyclonedx-cli`, or language-native tools (`pip-audit`, `npm audit`, `gradle dependencies`). Document findings as the baseline SBOM for each artifact.

2. **Adopt a standard SBOM format — choose CycloneDX or SPDX** — Generate machine-readable SBOMs in CycloneDX 1.4+ or SPDX 2.3+ format for every release. Include: package name, version, supplier, license, and cryptographic hash. Automate SBOM generation in CI so it is always current and never manually maintained.

3. **Pin all dependencies and enforce hash verification** — Replace version ranges with exact pinned versions and checksum verification in lock files (`package-lock.json`, `poetry.lock`, `go.sum`). Enable hash checking mode (`pip install --require-hashes`). Never fetch dependencies at build time without pre-verified checksums.

4. **Implement artifact signing — make provenance cryptographically verifiable** — Sign all build outputs (container images, binaries, packages) using Sigstore/cosign or a hardware-backed code signing certificate. Use keyless signing with OIDC-based identity tied to the CI/CD pipeline identity (not a developer's personal key). Store signatures in a transparency log (Rekor).

5. **Achieve SLSA provenance — attest how artifacts were built** — Implement SLSA Supply-chain Levels for Software Artifacts. Target SLSA Level 2 minimum (hosted build, signed provenance) for all production artifacts; Level 3 (hardened build, non-falsifiable provenance) for critical components. Generate and attach SLSA provenance attestations to every artifact.

6. **Harden the build environment — isolate and lock down CI/CD** — Use ephemeral build environments (new runner per build). Restrict outbound network access during builds to known package registries only. Pin CI action versions to full commit SHAs, not mutable tags (e.g., `actions/checkout@abc123def` not `@v3`). Audit all secrets accessible in build pipelines.

7. **Operate a private artifact registry — control what enters the build** — Mirror approved packages in an internal registry (Nexus, Artifactory, AWS CodeArtifact). Enforce that builds only pull from the internal mirror, never directly from public registries. Scan every package entering the mirror against CVE databases and license policy.

8. **Implement continuous vulnerability scanning — track exposure in real time** — Integrate SBOM-driven vulnerability scanning (Grype, Trivy, or Snyk) into CI gates and production monitoring. Alert on new CVEs affecting components in deployed SBOMs within 24 hours of disclosure. Define an SLA for patching based on CVSS severity (Critical ≤24h, High ≤7d).

9. **Define a vendor risk process for third-party software** — For commercial software and managed services, request SBOMs and attestations during procurement. Assess vendor SLSA maturity. Include supply chain security requirements in contracts. For critical vendors, require notification of any build system changes.

10. **Test the controls — verify detection capability** — Run tabletop exercises simulating a compromised dependency (e.g., a typosquatting package). Verify that SBOM comparison between releases detects unexpected new components. Test that signature verification rejects tampered artifacts before deployment.

## Rules

- Every artifact deployed to production must have a signed SBOM and verifiable build provenance
- Dependency pinning with hash verification is mandatory — floating version ranges are not permitted in production
- CI pipeline secrets must never be accessible to untrusted code (e.g., from forked PRs without review)
- All signing keys must be stored in HSMs or use keyless OIDC-based signing — no plaintext private keys in repositories or environment variables
- SBOM must be updated automatically on every build — manual SBOM maintenance is not acceptable

## Common Mistakes

- **SBOM as a one-time snapshot** — Generating an SBOM at release time only and never updating it means the SBOM is stale within days as dependencies change. Automate in CI.
- **Pinning direct dependencies but not transitive** — Lock files that capture only direct dependencies miss the long tail of transitive packages that are frequently the attack vector (e.g., event-stream incident in npm).
- **Signing the wrong thing** — Signing a container image tag (mutable) instead of the image digest (immutable) provides false assurance. Always sign and verify by digest.
- **Treating SLSA as a checkbox** — Achieving SLSA Level 2 paperwork without actually enforcing provenance verification at deployment time means the attestations are never checked and provide no real security.
- **Ignoring the build tool chain** — Pinning application dependencies while using unpinned build tools (compilers, linters, base images) leaves a gap — the SolarWinds attack targeted the build toolchain, not the application code.

## Examples

**GitHub Actions with Sigstore:** Use `sigstore/cosign-installer` in CI and `cosign sign` with OIDC keyless signing after pushing a container image. Consumers verify with `cosign verify --certificate-identity-regexp "github.com/org/repo" --certificate-oidc-issuer "https://token.actions.githubusercontent.com"`.

**SBOM-driven patch SLA:** After the Log4Shell disclosure, an organization with a CycloneDX SBOM for every deployed service ran `grype sbom:./sbom.json` across all service SBOMs in 20 minutes, producing a ranked list of affected services with CVSS scores, enabling prioritized patching within hours rather than the industry average of days.

## When NOT to Use

- Internal tooling or throwaway scripts with no external distribution and access to no sensitive systems (overhead exceeds risk)
- Proof-of-concept projects that will be rewritten before any production deployment
