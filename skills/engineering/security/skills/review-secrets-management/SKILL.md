---
name: review-secrets-management
description: Use when reviewing how secrets, credentials, API keys, or certificates are stored, rotated, and accessed in a system
source: OWASP Secrets Management Cheat Sheet (cheatsheetseries.owasp.org); HashiCorp Vault documentation; AWS Secrets Manager docs
tags: [security, secrets, credentials, vault, aws-secrets-manager, owasp]
verified: true
---

# Review Secrets Management

Audit how secrets are stored, distributed, and rotated to eliminate hardcoded credentials and reduce blast radius from leaks.

## Why This Is Best Practice

**Adopted by:** HashiCorp (Vault used by 70%+ of Fortune 500), AWS Secrets Manager, GitHub Advanced Security (secret scanning)
**Impact:** GitGuardian detected 10 million secrets exposed on GitHub in 2023; hardcoded credentials are the #1 cause of cloud breach initial access (Verizon DBIR 2023).

**Why best:** Secrets in code, environment files, or CI logs are a single git clone away from exposure. Centralized secrets management provides: access control, audit logging, automatic rotation, and short-lived credential issuance — none of which are possible with hardcoded or manually managed secrets.

## Steps

1. **Scan for hardcoded secrets** — Run `gitleaks`, `trufflehog`, or `detect-secrets` across the entire git history (not just HEAD). Enable GitHub Advanced Security secret scanning.
2. **Inventory all secrets** — Catalog every secret the system uses: DB passwords, API keys, TLS certificates, signing keys, third-party tokens. Document owner, rotation frequency, and current storage location.
3. **Classify by sensitivity** — Tier 1 (DB root, signing keys, payment credentials): most rigorous controls. Tier 2 (internal service tokens). Tier 3 (non-sensitive API keys for public APIs).
4. **Migrate to a secrets manager** — Move Tier 1/2 secrets to HashiCorp Vault, AWS Secrets Manager, GCP Secret Manager, or Azure Key Vault. Never store in environment variable files committed to git.
5. **Implement dynamic secrets** — Use Vault's dynamic secrets for DB credentials — issue short-lived credentials per connection; old credentials expire automatically, eliminating rotation burden.
6. **Enforce least privilege access** — Each service/role should access only the secrets it needs; use Vault policies or IAM conditions to enforce.
7. **Automate rotation** — Configure automatic rotation for all long-lived credentials; validate rotation works in staging before enabling in production.

## Rules

- Zero secrets in source code, ever — including test fixtures and example configs.
- Secrets must never appear in logs, error messages, or HTTP responses.
- All secret access must be audited — who accessed what, when.
- Short-lived credentials (TTL <24h) are always preferable to long-lived static keys.

## Examples

Bad: `DB_PASSWORD=s3cr3t` in `.env` committed to git.
Good: Application calls `vault read database/creds/my-role` at startup; receives a unique username/password valid for 1 hour; Vault auto-revokes when TTL expires.

CI/CD: Secrets injected at runtime via Vault agent sidecar or AWS IAM role; never stored in pipeline YAML.

## Common Mistakes

- **`.env` files in git** — even if removed in a later commit, they exist in git history; `gitleaks --no-git` will not find them; scan full history.
- **Long-lived service account keys** — GCP/AWS keys that never rotate are indefinite blast radius; prefer Workload Identity / IRSA.
- **Secrets in CI environment variables visible in logs** — always mask secrets in CI output; audit pipeline logs for leakage.

## When NOT to Use

- When reviewing a purely static frontend application that has no server-side secrets and communicates only with public, unauthenticated APIs, a secrets management audit has no applicable surface area.
- When the system is an air-gapped, offline embedded device with no network connectivity and no dynamic credential issuance, centralized secrets manager patterns designed for cloud or networked systems do not apply.
- When the scope of the review is limited to a single algorithm or business logic function with no I/O, credentials, or external service calls, a secrets management audit is out of scope and should not block that review.
