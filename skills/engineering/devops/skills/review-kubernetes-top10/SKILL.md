---
name: review-kubernetes-top10
description: Use when auditing a Kubernetes cluster or reviewing a deployment for security — systematically checking all 10 OWASP Kubernetes Top 10 (2022) vulnerability classes with test commands and remediation guidance.
source: 'OWASP Kubernetes Top 10 2022 (owasp.org/www-project-kubernetes-top-ten/); CIS Kubernetes Benchmark v1.8; kube-bench; NSA/CISA Kubernetes Hardening Guidance'
tags: [security, owasp, kubernetes, audit, rbac, containers, devops]
---

# Review Kubernetes Top 10

Audit Kubernetes clusters against the OWASP Kubernetes Top 10 (2022) using kubectl commands, kube-bench, and configuration inspection — covering K01 through K10 with specific test procedures and remediation steps.

## Why This Is Best Practice

**Adopted by:** OWASP Kubernetes Top 10 (2022) is the authoritative vulnerability taxonomy for Kubernetes security. CIS Kubernetes Benchmark v1.8 provides detailed test procedures aligned with the OWASP categories. kube-bench (Aqua Security) automates CIS benchmark testing and is used by Google GKE, AWS EKS, and Azure AKS in their security assessments. NSA/CISA Kubernetes Hardening Guidance (2022) endorses the same vulnerability classes as the OWASP Top 10.
**Impact:** Palo Alto Unit 42 (2022) found that 65% of Kubernetes clusters had at least one K01–K10 class vulnerability. The 2021 Azurescape attack exploited a container escape (K08 equivalent) to gain cross-tenant host access in Azure Container Instances. Tesla's 2018 cryptojacking incident involved an exposed Kubernetes dashboard (K09 equivalent) with no authentication. A structured audit against K01–K10 provides coverage for the most frequently exploited Kubernetes attack vectors.
**Why best:** Ad-hoc Kubernetes reviews miss vulnerability classes because they focus on known issues rather than systematically checking all attack surfaces. The K01–K10 framework provides completeness — a reviewer who checks all 10 classes is significantly less likely to miss a critical misconfiguration than one using an unstructured checklist.

Sources: OWASP Kubernetes Top 10 (2022); CIS Kubernetes Benchmark v1.8; kube-bench documentation; NSA/CISA Kubernetes Hardening Guidance (August 2022)

## Steps

### K01 — Insecure Workload Configurations

```bash
# Find pods running as root or with privileged flag
kubectl get pods --all-namespaces -o json | \
  jq '.items[] | select(.spec.containers[].securityContext.runAsRoot == true or
      .spec.containers[].securityContext.privileged == true) |
      {name: .metadata.name, namespace: .metadata.namespace}'

# Find pods without resource limits
kubectl get pods --all-namespaces -o json | \
  jq '.items[] | select(.spec.containers[].resources.limits == null) |
      {name: .metadata.name, namespace: .metadata.namespace}'
```

**Fix:** Set `runAsNonRoot: true`, `allowPrivilegeEscalation: false`, resource limits, and `readOnlyRootFilesystem: true` in pod securityContext. Apply `restricted` Pod Security Standard to production namespaces.

### K02 — Supply Chain Vulnerabilities

```bash
# Check if image tags are pinned to digests
kubectl get pods --all-namespaces -o json | \
  jq '.items[].spec.containers[].image' | \
  grep -v '@sha256:'  # unpinned images appear here

# Verify Trivy operator is installed
kubectl get pods -n trivy-system
kubectl get vulnerabilityreports --all-namespaces | head -20
```

**Fix:** Pin images to digest (`image@sha256:...`), run Trivy operator for continuous scanning, use OPA Gatekeeper to enforce registry allowlist.

### K03 — Overly Permissive RBAC

```bash
# Find ClusterRoleBindings to cluster-admin
kubectl get clusterrolebindings -o json | \
  jq '.items[] | select(.roleRef.name == "cluster-admin") |
      {name: .metadata.name, subjects: .subjects}'

# Find wildcard permissions in roles
kubectl get clusterroles -o json | \
  jq '.items[] | select(.rules[].verbs[] == "*" or .rules[].resources[] == "*") |
      .metadata.name'

# Find service accounts with API access
kubectl get rolebindings,clusterrolebindings --all-namespaces -o json | \
  jq '.items[] | select(.subjects[]?.kind == "ServiceAccount")'
```

**Fix:** Remove `cluster-admin` bindings for workloads. Replace wildcard rules with specific resource/verb combinations. Disable `automountServiceAccountToken` on pods that don't use the K8s API.

### K04 — Lack of Centralized Policy Enforcement

```bash
# Check if OPA Gatekeeper or Kyverno is installed
kubectl get pods -n gatekeeper-system 2>/dev/null || echo "Gatekeeper not installed"
kubectl get pods -n kyverno 2>/dev/null || echo "Kyverno not installed"

# Check Pod Security Admission labels on namespaces
kubectl get namespaces -o json | \
  jq '.items[] | {name: .metadata.name, labels: .metadata.labels} |
      select(.labels | has("pod-security.kubernetes.io/enforce") | not)'
```

**Fix:** Deploy OPA Gatekeeper or Kyverno with deny policies. Apply Pod Security Admission labels to all namespaces.

### K05 — Inadequate Logging and Monitoring

```bash
# Check audit logging configuration on API server
kubectl -n kube-system get pod -l component=kube-apiserver -o json | \
  jq '.items[].spec.containers[].command[] | select(contains("audit-log"))'

# Check if Falco is running (runtime anomaly detection)
kubectl get pods --all-namespaces | grep falco
```

**Fix:** Enable API server audit logging with `--audit-log-path` and `--audit-policy-file`. Deploy Falco for runtime syscall monitoring. Ship audit logs to a SIEM.

### K06 — Broken Authentication Mechanisms

```bash
# Check if anonymous authentication is enabled (should be disabled)
kubectl -n kube-system get pod -l component=kube-apiserver -o json | \
  jq '.items[].spec.containers[].command[] | select(contains("anonymous-auth=true"))'

# Check if service account tokens have expiration
kubectl get serviceaccounts --all-namespaces -o json | \
  jq '.items[] | select(.metadata.annotations["kubernetes.io/service-account.uid"] != null) |
      .metadata.name'
```

**Fix:** Set `--anonymous-auth=false` on API server. Use short-lived tokens via bound service account tokens (`expirationSeconds: 3600`). Enable OIDC authentication for human users.

### K07 — Missing Network Segmentation

```bash
# Find namespaces without default-deny NetworkPolicy
for ns in $(kubectl get namespaces -o jsonpath='{.items[*].metadata.name}'); do
  count=$(kubectl get networkpolicies -n $ns 2>/dev/null | grep -c default-deny || echo 0)
  if [ "$count" -eq 0 ]; then
    echo "No default-deny NetworkPolicy in namespace: $ns"
  fi
done
```

**Fix:** Apply default-deny NetworkPolicy to all namespaces. Verify CNI plugin supports NetworkPolicy (Calico, Cilium, Weave). Use explicit allow policies for required communication paths.

### K08 — Secrets Management Failures

```bash
# Check if encryption at rest is configured for secrets
kubectl -n kube-system get pod -l component=kube-apiserver -o json | \
  jq '.items[].spec.containers[].command[] | select(contains("encryption-provider-config"))'

# Find secrets referenced as environment variables (riskier than volume mounts)
kubectl get pods --all-namespaces -o json | \
  jq '.items[] | select(.spec.containers[].env[]?.valueFrom.secretKeyRef != null) |
      {name: .metadata.name, namespace: .metadata.namespace}'
```

**Fix:** Enable `EncryptionConfiguration` with AES-CBC/GCM for etcd. Migrate to External Secrets Operator + Vault/Secrets Manager. Prefer volume mounts over env vars for secrets.

### K09 — Misconfigured Cluster Components

```bash
# Check if the Kubernetes Dashboard is exposed
kubectl get services --all-namespaces | grep kubernetes-dashboard
kubectl get ingress --all-namespaces | grep dashboard

# Check etcd access
kubectl -n kube-system get pod -l component=etcd -o json | \
  jq '.items[].spec.containers[].command[] | select(contains("listen-client-urls"))'
```

**Fix:** Do not expose Kubernetes Dashboard to the internet. Require RBAC + authentication for Dashboard. Ensure etcd listens only on localhost or internal network, not `0.0.0.0`.

### K10 — Outdated and Vulnerable Kubernetes Components

```bash
# Check cluster version
kubectl version --short

# Check node versions
kubectl get nodes -o json | jq '.items[] | {name: .metadata.name, kubelet: .status.nodeInfo.kubeletVersion}'

# Run kube-bench for CIS benchmark
kubectl apply -f https://raw.githubusercontent.com/aquasecurity/kube-bench/main/job.yaml
kubectl logs job/kube-bench | grep -E 'FAIL|WARN'
```

**Fix:** Keep Kubernetes within N-2 supported minor versions. Subscribe to kubernetes-security-announce mailing list. Apply security patches within 30 days of release.

## Rules

- Run this checklist against each cluster at least quarterly — Kubernetes evolves rapidly and new vulnerability classes emerge with each release.
- K01 and K03 (workload configs and RBAC) account for 70%+ of exploited Kubernetes vulnerabilities — prioritize these if time is limited.
- kube-bench automates K01, K03, K05, K06, K09, K10 checks with pass/fail output — run it as the baseline.

## Common Mistakes

- **Auditing only the `default` namespace** — production workloads are often in custom namespaces; audit all namespaces.
- **Treating kube-bench PASS as "secure"** — kube-bench covers CIS benchmarks, not all OWASP K01–K10 classes; both are needed.
- **Fixing findings in production without testing** — RBAC and NetworkPolicy changes can disrupt running applications; test in staging first.
