---
name: apply-kubernetes-security
description: Use when deploying workloads to Kubernetes — configuring RBAC, network policies, pod security standards, secrets management, and admission controls to prevent privilege escalation and lateral movement.
source: 'OWASP Kubernetes Security Cheat Sheet (owasp.org/www-project-cheat-sheets); CIS Kubernetes Benchmark; NIST SP 800-190; Kubernetes Pod Security Standards documentation'
tags: [security, owasp, kubernetes, rbac, containers, devops, developer]
---

# Apply Kubernetes Security

Harden Kubernetes clusters by enforcing RBAC least privilege, network policies, pod security standards, and encrypted secrets — preventing privilege escalation, lateral movement, and container escape.

## Why This Is Best Practice

**Adopted by:** OWASP Kubernetes Security Cheat Sheet and CIS Kubernetes Benchmark v1.8 are the authoritative references. NSA/CISA released "Kubernetes Hardening Guidance" (2022) as official federal guidance. Google GKE Autopilot, AWS EKS default node groups, and Azure AKS enforce pod security admission and RBAC by default. Kubernetes Pod Security Standards (PSS) replaced PodSecurityPolicy in K8s 1.25 and are now the built-in hardening mechanism.
**Impact:** The 2022 Tesla cryptojacking incident and the 2020 SolarWinds supply chain attack both involved misconfigured Kubernetes clusters with over-permissive RBAC. Palo Alto Unit 42 (2022) found 65% of Kubernetes clusters in production had at least one container running as root. NSA/CISA guidance found that default Kubernetes configurations grant excessive permissions that enable cluster-wide lateral movement from a single compromised pod.
**Why best:** Default Kubernetes configuration prioritizes compatibility over security — default service accounts have API access, pods can communicate freely, and secrets are stored as base64 (not encrypted). Applying RBAC least privilege, network policies, and pod security standards systematically closes these gaps vs. relying on application-level controls alone.

Sources: OWASP Kubernetes Security Cheat Sheet; CIS Kubernetes Benchmark v1.8; NSA/CISA Kubernetes Hardening Guidance (2022); Kubernetes Pod Security Standards documentation

## Steps

1. **Enforce Pod Security Standards at the namespace level**:

   ```bash
   # Apply restricted policy to production namespace
   kubectl label namespace production \
     pod-security.kubernetes.io/enforce=restricted \
     pod-security.kubernetes.io/enforce-version=latest \
     pod-security.kubernetes.io/warn=restricted \
     pod-security.kubernetes.io/audit=restricted
   ```

   ```yaml
   # Pod spec that passes restricted standard
   apiVersion: v1
   kind: Pod
   spec:
     securityContext:
       runAsNonRoot: true
       runAsUser: 1000
       seccompProfile:
         type: RuntimeDefault
     containers:
     - name: app
       securityContext:
         allowPrivilegeEscalation: false
         readOnlyRootFilesystem: true
         capabilities:
           drop: ["ALL"]
       resources:
         limits:
           cpu: "500m"
           memory: "256Mi"
         requests:
           cpu: "100m"
           memory: "128Mi"
   ```

2. **Apply RBAC least privilege — never use `cluster-admin` for workloads**:

   ```yaml
   # Role: namespace-scoped, specific resources and verbs only
   apiVersion: rbac.authorization.k8s.io/v1
   kind: Role
   metadata:
     name: app-reader
     namespace: production
   rules:
   - apiGroups: [""]
     resources: ["configmaps"]
     verbs: ["get", "list"]
   # NOT: resources: ["*"], verbs: ["*"]
   ---
   apiVersion: rbac.authorization.k8s.io/v1
   kind: RoleBinding
   metadata:
     name: app-reader-binding
     namespace: production
   subjects:
   - kind: ServiceAccount
     name: myapp-sa
     namespace: production
   roleRef:
     kind: Role
     name: app-reader
     apiGroup: rbac.authorization.k8s.io
   ```

   ```yaml
   # Disable default service account auto-mount in pod spec
   spec:
     automountServiceAccountToken: false
   ```

3. **Implement network policies — deny by default, allow by exception**:

   ```yaml
   # Default deny all ingress and egress in namespace
   apiVersion: networking.k8s.io/v1
   kind: NetworkPolicy
   metadata:
     name: default-deny-all
     namespace: production
   spec:
     podSelector: {}
     policyTypes:
     - Ingress
     - Egress
   ---
   # Allow only frontend → backend on port 8080
   apiVersion: networking.k8s.io/v1
   kind: NetworkPolicy
   metadata:
     name: allow-frontend-to-backend
     namespace: production
   spec:
     podSelector:
       matchLabels:
         app: backend
     ingress:
     - from:
       - podSelector:
           matchLabels:
             app: frontend
       ports:
       - protocol: TCP
         port: 8080
   ```

4. **Use Kubernetes Secrets with encryption at rest — never store credentials in ConfigMaps or pod env from literals**:

   ```bash
   # Enable encryption at rest in kube-apiserver (add to --encryption-provider-config)
   # encryption-config.yaml:
   apiVersion: apiserver.config.k8s.io/v1
   kind: EncryptionConfiguration
   resources:
   - resources:
     - secrets
     providers:
     - aescbc:
         keys:
         - name: key1
           secret: <base64-encoded-32-byte-key>
     - identity: {}
   ```

   ```yaml
   # Reference secrets as volumes, not env vars (avoids exposure in pod spec)
   spec:
     volumes:
     - name: db-creds
       secret:
         secretName: database-credentials
     containers:
     - name: app
       volumeMounts:
       - name: db-creds
         mountPath: "/etc/secrets"
         readOnly: true
   ```

   For production: use External Secrets Operator with AWS Secrets Manager, HashiCorp Vault, or GCP Secret Manager instead of native K8s Secrets.

5. **Scan images in CI and enforce admission with an OPA/Gatekeeper policy**:

   ```bash
   # Trivy operator — scans images automatically within cluster
   kubectl apply -f https://raw.githubusercontent.com/aquasecurity/trivy-operator/main/deploy/static/trivy-operator.yaml

   # Check scan results
   kubectl get vulnerabilityreports -n production
   ```

   ```yaml
   # OPA Gatekeeper constraint: require specific image registry
   apiVersion: constraints.gatekeeper.sh/v1beta1
   kind: K8sAllowedRepos
   metadata:
     name: require-internal-registry
   spec:
     match:
       kinds:
       - apiGroups: [""]
         kinds: ["Pod"]
     parameters:
       repos:
       - "registry.internal.company.com/"
   ```

6. **Audit cluster configuration with kube-bench**:

   ```bash
   # Run CIS benchmark against current node
   kubectl apply -f https://raw.githubusercontent.com/aquasecurity/kube-bench/main/job.yaml
   kubectl logs job/kube-bench

   # Or locally
   kube-bench run --targets node,master,etcd,policies
   ```

## Rules

- Never bind `cluster-admin` to a service account used by an application workload.
- Disable `automountServiceAccountToken: false` on every pod that doesn't call the Kubernetes API.
- Network policies require a CNI plugin that enforces them (Calico, Cilium, Weave) — vanilla kubenet ignores NetworkPolicy objects.
- etcd must be encrypted at rest and accessible only to the API server — direct etcd access bypasses all RBAC.

## Common Mistakes

- **`securityContext` at pod level vs container level** — `runAsNonRoot` must be set at both levels; container-level overrides pod-level.
- **Storing secrets as environment variables** — environment variables appear in `kubectl describe pod` and process listings; use volume mounts.
- **No resource limits** — a pod without limits can starve the node; resource limits are also required for `restricted` pod security standard.
- **Using `hostNetwork: true` or `hostPID: true`** — gives the container access to the host network stack and process namespace, defeating isolation.
