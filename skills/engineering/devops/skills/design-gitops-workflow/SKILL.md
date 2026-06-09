---
name: design-gitops-workflow
description: Use when designing a GitOps workflow where Git is the single source of truth for infrastructure and application configuration, with automated reconciliation to the desired state
source: 'Weaveworks "GitOps: Operations by Pull Request" Limoncelli (2017); CNCF OpenGitOps Principles v1.0 (2022); Argo CD Documentation (2023); Flux CD GitOps Toolkit patterns (2023)'
tags: [devops, gitops, kubernetes, infrastructure, ci-cd, argocd, flux]
verified: true
---

# Design GitOps Workflow

Use Git as the authoritative source of truth for all infrastructure and application state, with automated operators continuously reconciling live systems to the declared desired state.

## Why This Is Best Practice

**Adopted by:** Weaveworks (originators, production since 2017), Intuit (Argo CD, open-sourced 2018), CNCF member organizations via OpenGitOps Working Group, Shopify, Adobe, and hundreds of Kubernetes-native organizations

**Impact:** DORA 2023 research correlates GitOps practices with elite performance: teams using GitOps-style declarative configuration report 3x faster mean time to recovery (MTTR) and 50% fewer change failure rates compared to imperative CI/CD pipelines. Weaveworks reported 10x reduction in deployment lead time after full GitOps adoption. Argo CD manages 1M+ deployed applications across its user base.

**Why best:** Traditional CI/CD pushes changes imperativity — the pipeline is the source of truth and state drifts silently. GitOps inverts this: Git holds the desired state, and an operator continuously pulls and applies it. Drift is detected and corrected automatically. Every change has a full audit trail via Git history. Rollback is `git revert`. This eliminates entire categories of configuration drift, snowflake environments, and "who changed this?" incidents.

Sources: Limoncelli, T. "GitOps: Operations by Pull Request" Weaveworks (2017); CNCF OpenGitOps Principles v1.0 (2022); Argo CD docs argoproj.github.io; Flux CD docs fluxcd.io; DORA "Accelerate State of DevOps" (2023)

## Steps

1. **Define the four GitOps principles — align the team** — Before implementation, ensure the team understands and commits to the CNCF OpenGitOps principles: (1) desired state is declarative, (2) desired state is versioned and immutable in Git, (3) desired state is pulled automatically by software agents, (4) agents continuously reconcile and report divergence. These principles rule out "push" pipelines and out-of-band changes.

2. **Structure the Git repository topology — monorepo or polyrepo** — Choose a repository structure: a single config repo for all environments and apps (monorepo), or per-app repos with a separate environments/config repo. Recommended: separate application repos (containing source + Dockerfile) from a dedicated GitOps config repo (containing Kubernetes manifests or Helm values). This prevents application CI from triggering unnecessary environment reconciliation.

3. **Organize environments with directory or branch strategy** — Use directories per environment (`environments/dev/`, `environments/staging/`, `environments/prod/`) within the config repo, all on the main branch. Avoid long-lived environment branches — they create merge conflicts and obscure promotion history. Promotion is a pull request from `environments/staging/` changes into `environments/prod/`.

4. **Choose and deploy a GitOps operator — Argo CD or Flux CD** — Select an operator based on team needs. Argo CD offers a rich UI and Application CRD abstraction, better for teams wanting visual dashboards and application grouping. Flux CD is more modular and Kubernetes-native, better for automation-first teams. Deploy the operator into each cluster. Configure it to watch the config repository and target namespace.

5. **Define Applications or Kustomizations — point operator at desired state** — For Argo CD: define `Application` CRDs that map a Git path to a target cluster and namespace, with sync policy `automated` and `selfHeal: true`. For Flux CD: define `GitRepository` + `Kustomization` resources. Enable pruning so resources removed from Git are deleted from the cluster.

6. **Automate image promotion — update Git, not the cluster directly** — Configure image automation to detect new container image tags (Flux Image Automation, Argo CD Image Updater) and open pull requests or auto-commit version bumps to the config repo. Never update a running deployment via `kubectl set image` — all changes flow through Git.

7. **Enforce branch protection and PR review for production** — Require pull request review and CI status checks before merging to the path watched for production environments. Use CODEOWNERS to require infrastructure team approval for cluster-wide changes. Tag releases in the config repo to tie deployments to auditable events.

8. **Implement drift detection alerting — know when reality diverges** — Configure Argo CD or Flux alerts to notify (Slack, PagerDuty) when any Application is `OutOfSync` for more than N minutes. Treat persistent drift as an incident — it indicates someone made an out-of-band change or the operator is failing to reconcile.

9. **Define the emergency break-glass process** — Document when and how to make emergency changes (direct `kubectl` apply during an outage). Require that every break-glass action is immediately followed by a PR to bring Git back into alignment with what was applied. Log break-glass events in the incident management system.

10. **Progressively restrict direct cluster access** — Over time, remove `kubectl` write access from developers and CI pipelines for production namespaces. The GitOps operator becomes the only writer. This eliminates a class of unauthorized or accidental changes. Retain read access and break-glass access for on-call engineers via time-limited IAM escalation.

## Rules

- Git is the only source of truth — no manual `kubectl apply`, `helm upgrade`, or console-based changes in production
- All changes to desired state must go through a pull request, not direct commits to main/trunk
- The GitOps operator must run with reconciliation enabled — periodic sync intervals (≤5 minutes) are not sufficient for self-healing; enable continuous watch mode
- Secrets must never be stored in plaintext in the Git config repo — use Sealed Secrets, External Secrets Operator, or SOPS encryption
- Rollback is always `git revert` — never operator-level forced sync to a previous version without a corresponding Git state change

## Common Mistakes

- **Mixing application build and config in one repo** — When the same repo triggers both CI (build) and CD (deploy), every code commit causes a deployment evaluation. Separate concerns: source code repo triggers image build; config repo triggers deployment.
- **Using GitOps only for dev, not production** — Teams adopt GitOps for dev but keep manual processes for production "because it's sensitive." This defeats the audit trail and reproducibility benefits. Production is the highest-value target for GitOps controls.
- **Storing secrets in Git unencrypted** — A common early mistake. Even private repos are not appropriate for plaintext secrets. Integrate secret management from day one.
- **Enabling auto-sync without pruning** — Without pruning, resources removed from Git linger in the cluster as orphans. Enable `prune: true` and test deletion via Git.
- **Not testing the config repo with CI** — Manifests and Helm values in the config repo must be validated by CI (kubeval, helm lint, kustomize build) before merge to prevent broken desired state that blocks the operator.

## Examples

**Argo CD multi-environment promotion:** A `staging` Application watches `environments/staging/` and auto-syncs. After validating staging, a developer opens a PR copying the validated image tag into `environments/prod/`. A required reviewer approves. Merge triggers automatic sync of the `prod` Application. Full audit trail in Git.

**Flux CD with image automation:** Flux `ImageRepository` watches the container registry. `ImagePolicy` selects the latest semver tag matching `>=1.2.0`. `ImageUpdateAutomation` commits the new tag to the config repo path. A Flux `Kustomization` applies the updated manifest within 60 seconds, with a Slack notification on completion.

## When NOT to Use

- Non-Kubernetes infrastructure that lacks declarative state representations (e.g., legacy VM-based systems using imperative Ansible without idempotent modules)
- Very small teams with a single environment where the overhead of a GitOps operator exceeds the operational complexity being solved
