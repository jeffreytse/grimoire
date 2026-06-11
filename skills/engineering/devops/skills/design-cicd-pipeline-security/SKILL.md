---
name: design-cicd-pipeline-security
description: Use when designing or hardening CI/CD pipelines — securing secret injection, restricting pipeline permissions, signing artifacts, and preventing supply chain attacks in GitHub Actions, GitLab CI, or Jenkins.
source: 'OWASP CI/CD Security Cheat Sheet (owasp.org/www-project-cheat-sheets); CNCF Software Supply Chain Best Practices; SLSA Framework (slsa.dev); NIST SP 800-218 (Secure Software Development Framework)'
tags: [security, owasp, cicd, supply-chain, github-actions, secrets, devops, developer]
---

# Design CI/CD Pipeline Security

Secure CI/CD pipelines by restricting pipeline permissions to least privilege, injecting secrets from vaults (not env vars), signing artifacts with Sigstore/cosign, and pinning actions to commit SHAs — preventing secret exfiltration and supply chain attacks.

## Why This Is Best Practice

**Adopted by:** OWASP CI/CD Security Cheat Sheet (2023) and CNCF Software Supply Chain Best Practices define the standard. The SLSA Framework (slsa.dev, backed by Google, Microsoft, NIST) provides a supply chain integrity maturity model adopted by GitHub, GitLab, and npm. NIST SP 800-218 (SSDF) mandates supply chain controls for US federal software. GitHub, Google, and Shopify all use artifact signing and pinned action versions in production pipelines.
**Impact:** The 2021 Codecov breach compromised CI environment variables and exfiltrated credentials from thousands of pipelines. The 2020 SolarWinds attack injected malicious code at the build stage — directly what SLSA framework prevents. GitHub's own research (2022) found 95% of GitHub Actions workflows use at least one unpinned action, enabling tag mutation attacks. CircleCI's 2023 breach was caused by compromised CI credentials exposed in pipeline logs.
**Why best:** Pipeline security vs. application security: application security controls (SAST, SCA) catch code vulnerabilities, but pipeline security prevents the pipeline itself from becoming an attack vector. An attacker with pipeline access can inject arbitrary code, exfiltrate all secrets, and publish malicious artifacts — bypassing all code review and security scanning.

Sources: OWASP CI/CD Security Cheat Sheet; CNCF Software Supply Chain Best Practices Guide; SLSA framework specification; NIST SP 800-218; GitHub Security Lab research (2022)

## Steps

1. **Restrict GitHub Actions pipeline permissions to minimum required**:

   ```yaml
   # At top of workflow — restrict all permissions
   permissions:
     contents: read   # only what's needed

   # Or per-job
   jobs:
     build:
       permissions:
         contents: read
         packages: write  # only for publish job
   ```

   ```yaml
   # Never use: permissions: write-all
   # Default GitHub token has read-all — explicitly restrict
   ```

2. **Pin actions to commit SHAs, not tags — tags are mutable**:

   ```yaml
   # BAD — tag can be changed to point at malicious code
   - uses: actions/checkout@v4
   - uses: actions/setup-node@main

   # GOOD — SHA is immutable
   - uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683  # v4.2.2
   - uses: actions/setup-node@39370e3970a6d050c480ffad4ff0ed4d3fdee5af  # v4.1.0
   ```

   Use Dependabot to auto-update pinned SHAs:
   ```yaml
   # .github/dependabot.yml
   version: 2
   updates:
   - package-ecosystem: github-actions
     directory: "/"
     schedule:
       interval: weekly
   ```

3. **Never hardcode secrets — use GitHub Secrets or a vault**:

   ```yaml
   # BAD — secret visible in logs if run echoes it
   env:
     DATABASE_URL: postgresql://user:password@host/db

   # GOOD — reference from GitHub Secrets
   env:
     DATABASE_URL: ${{ secrets.DATABASE_URL }}

   # BETTER — fetch from HashiCorp Vault or AWS Secrets Manager at runtime
   - name: Import secrets from Vault
     uses: hashicorp/vault-action@d1720f055e0635fd932a1d2a48f87a666a57906c  # v3.0.0
     with:
       url: https://vault.company.com
       method: jwt
       role: github-actions-role
       secrets: |
         secret/data/production/db password | DATABASE_PASSWORD
   ```

4. **Sign container images and binaries with Sigstore/cosign**:

   ```yaml
   - name: Sign image with cosign (keyless — uses OIDC)
     run: |
       cosign sign --yes \
         registry.company.com/myapp:${{ github.sha }}
   ```

   ```bash
   # Verify signature before deploying
   cosign verify \
     --certificate-identity "https://github.com/company/myapp/.github/workflows/release.yml@refs/heads/main" \
     --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
     registry.company.com/myapp:abc123
   ```

5. **Generate and verify SBOM (Software Bill of Materials)**:

   ```yaml
   - name: Generate SBOM
     uses: anchore/sbom-action@e8d2a6937ecead383dfe75190d104edd1f9c5666  # v0.15.0
     with:
       image: registry.company.com/myapp:${{ github.sha }}
       format: spdx-json
       output-file: sbom.spdx.json

   - name: Attest SBOM to image
     run: |
       cosign attest --yes \
         --predicate sbom.spdx.json \
         --type spdxjson \
         registry.company.com/myapp:${{ github.sha }}
   ```

6. **Isolate pipeline environments — separate jobs for build vs. deploy**:

   ```yaml
   jobs:
     build:
       runs-on: ubuntu-latest
       permissions:
         contents: read
         packages: write
       # No production credentials here

     deploy:
       needs: build
       runs-on: ubuntu-latest
       environment: production  # requires manual approval
       permissions:
         id-token: write  # for OIDC auth to cloud
       steps:
       - name: Configure AWS credentials via OIDC (no long-lived keys)
         uses: aws-actions/configure-aws-credentials@e3dd6a429d7300a6a4c196c26e071d42e0343502  # v4.0.2
         with:
           role-to-assume: arn:aws:iam::123456789:role/github-actions-deploy
           aws-region: us-east-1
   ```

## Rules

- Never print secrets to pipeline logs — `echo $SECRET` is visible in logs even if the secret is masked by GitHub.
- Use OIDC federated identity (not long-lived IAM keys) for cloud authentication — keys can be rotated, stolen keys are not auditable.
- Set branch protection rules: require PR reviews and passing status checks before merge to main.
- Audit third-party actions before using them — treat them as code dependencies with the same security scrutiny.

## Common Mistakes

- **`pull_request_target` with checkout of PR head** — `pull_request_target` runs with write permissions; checking out untrusted PR code gives attackers write access.
- **Self-hosted runners shared across security boundaries** — a compromised job on a self-hosted runner can affect all subsequent jobs on that runner via filesystem persistence.
- **Caching secrets inadvertently** — `actions/cache` may persist files containing secrets across runs; exclude sensitive directories.
- **No approval gate for production deployments** — automated deploy to production without manual approval means any merged PR can trigger a production deployment.
