---
name: design-cloud-workload-identity
description: Use when Kubernetes pods or cloud services need to access cloud APIs — replacing long-lived IAM keys with IRSA (AWS), Workload Identity (GCP), or Managed Identity (Azure) to eliminate credential leakage risks.
source: 'OWASP Cloud-Native Application Security Top 10 C3 (owasp.org/www-project-cloud-native-application-security-top-10/); AWS IRSA documentation; Google Workload Identity Federation documentation; NIST SP 800-190'
tags: [security, owasp, cloud, kubernetes, iam, irsa, workload-identity, k8s, developer]
---

# Design Cloud Workload Identity

Replace long-lived IAM access keys in Kubernetes pods with IRSA (AWS), Workload Identity (GCP), or Pod Identity (Azure) — cryptographically binding cloud permissions to the pod's identity without credentials that can be stolen from environment variables or files.

## Why This Is Best Practice

**Adopted by:** OWASP Cloud-Native Application Security Top 10 C3 (Improper Authentication and Authorization). AWS IRSA (IAM Roles for Service Accounts) is the recommended approach in the AWS EKS documentation and AWS Well-Architected Framework. Google Cloud's Workload Identity is required for GKE security compliance. Azure's Managed Identity for AKS is the Microsoft security standard. Netflix, Spotify, and Shopify use workload identity federation exclusively — no static IAM keys in production Kubernetes.
**Impact:** The 2022 Microsoft Azure customer breach involved Kubernetes workloads using long-lived managed identity credentials stored as Kubernetes secrets — visible to anyone with `kubectl get secret`. The 2021 CircleCI breach exposed environment variables including AWS_ACCESS_KEY_ID credentials stored in pods. HashiCorp's 2022 security research found AWS access keys in Kubernetes ConfigMaps, Secrets, and environment variables in 43% of audited clusters. Static credentials in pods are exposed in CloudFormation/Terraform state, Kubernetes secret stores, container image layers, and process environments.
**Why best:** OIDC-based workload identity provides temporary, automatically-rotating credentials (valid ~1 hour) that are tied to the specific pod's service account — not exportable as static keys. A stolen OIDC token from a compromised pod is valid only while the pod runs and can't be used from outside the cluster. Static IAM keys (the alternative) are valid indefinitely until manually rotated and can be used from any IP address.

Sources: OWASP Cloud-Native Top 10 C3; AWS IRSA documentation; Google Workload Identity docs; Azure Managed Identity for AKS; Netflix cloud security architecture (2022)

## Steps

1. **AWS: configure IRSA (IAM Roles for Service Accounts)**:

   ```bash
   # Step 1: Create OIDC provider for the EKS cluster
   eksctl utils associate-iam-oidc-provider \
     --cluster production-cluster \
     --approve

   # Step 2: Create IAM role with trust policy scoped to specific service account
   ```

   ```json
   {
     "Version": "2012-10-17",
     "Statement": [{
       "Effect": "Allow",
       "Principal": {
         "Federated": "arn:aws:iam::123456789:oidc-provider/oidc.eks.us-east-1.amazonaws.com/id/EXAMPLED539D4633E53DE1B71EXAMPLE"
       },
       "Action": "sts:AssumeRoleWithWebIdentity",
       "Condition": {
         "StringEquals": {
           "oidc.eks.us-east-1.amazonaws.com/id/EXAMPLED:sub":
             "system:serviceaccount:production:order-service"
         }
       }
     }]
   }
   ```

   ```yaml
   # Step 3: Annotate the Kubernetes service account
   apiVersion: v1
   kind: ServiceAccount
   metadata:
     name: order-service
     namespace: production
     annotations:
       eks.amazonaws.com/role-arn: "arn:aws:iam::123456789:role/order-service-role"
   ---
   # Step 4: Use the service account in the pod — no credentials needed
   apiVersion: v1
   kind: Pod
   spec:
     serviceAccountName: order-service
     containers:
     - name: app
       # boto3 automatically uses IRSA credentials via token projection
       # NO: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY env vars
   ```

2. **GCP: configure Workload Identity for GKE**:

   ```bash
   # Step 1: Enable Workload Identity on the cluster
   gcloud container clusters update production-cluster \
     --workload-pool=PROJECT_ID.svc.id.goog

   # Step 2: Create GCP service account
   gcloud iam service-accounts create order-service-sa \
     --project=PROJECT_ID

   # Step 3: Bind Kubernetes service account to GCP service account
   gcloud iam service-accounts add-iam-policy-binding \
     order-service-sa@PROJECT_ID.iam.gserviceaccount.com \
     --role=roles/iam.workloadIdentityUser \
     --member="serviceAccount:PROJECT_ID.svc.id.goog[production/order-service]"
   ```

   ```yaml
   # Step 4: Annotate Kubernetes service account
   apiVersion: v1
   kind: ServiceAccount
   metadata:
     name: order-service
     namespace: production
     annotations:
       iam.gke.io/gcp-service-account: "order-service-sa@PROJECT_ID.iam.gserviceaccount.com"
   ```

3. **Azure: configure Pod Identity (Workload Identity for AKS)**:

   ```bash
   # Step 1: Enable OIDC issuer and workload identity on cluster
   az aks update -g myResourceGroup -n production-cluster \
     --enable-oidc-issuer --enable-workload-identity

   # Step 2: Create Managed Identity
   az identity create -g myResourceGroup -n order-service-identity

   # Step 3: Get OIDC issuer URL
   export OIDC_ISSUER=$(az aks show -g myResourceGroup -n production-cluster \
     --query "oidcIssuerProfile.issuerUrl" -o tsv)

   # Step 4: Create federated credential
   az identity federated-credential create \
     --name order-service-federated \
     --identity-name order-service-identity \
     --resource-group myResourceGroup \
     --issuer "${OIDC_ISSUER}" \
     --subject "system:serviceaccount:production:order-service"
   ```

4. **Verify no static credentials exist in the cluster**:

   ```bash
   # Scan for AWS keys in Kubernetes secrets
   kubectl get secrets --all-namespaces -o json | \
     jq '.items[] | select(.data != null) |
         .data | to_entries[] |
         select(.value != null) |
         {key: .key, value: (.value | @base64d)}' | \
     grep -E "AKIA|ASIA"  # AWS access key prefixes

   # Scan for GCP service account JSON in secrets
   kubectl get secrets --all-namespaces -o json | \
     jq '.items[].data // {}' | \
     grep -i "private_key_id"

   # Scan pod specs for credential environment variables
   kubectl get pods --all-namespaces -o json | \
     jq '.items[].spec.containers[].env[]? |
         select(.name | test("KEY|SECRET|PASSWORD|TOKEN|CREDENTIAL"; "i"))'
   ```

5. **Restrict instance metadata service access to prevent credential theft**:

   ```yaml
   # Block pods from accessing EC2 instance metadata (which has node IAM role)
   # Using Kubernetes NetworkPolicy or IMDSv2 hop limit
   apiVersion: networking.k8s.io/v1
   kind: NetworkPolicy
   metadata:
     name: block-imds
     namespace: production
   spec:
     podSelector: {}
     policyTypes:
     - Egress
     egress:
     # Allow all egress except to IMDS endpoint (169.254.169.254)
     - to:
       - ipBlock:
           cidr: 0.0.0.0/0
           except:
           - 169.254.169.254/32
   ```

   ```bash
   # AWS: enforce IMDSv2 (hop limit 1 prevents pod access to node credentials)
   aws ec2 modify-instance-metadata-options \
     --instance-id i-xxx \
     --http-tokens required \
     --http-put-response-hop-limit 1
   ```

## Rules

- Never create IAM users with long-lived access keys for Kubernetes workloads — use IRSA/Workload Identity exclusively.
- Each pod service account must have its own IAM role with only the permissions that specific service needs.
- Regularly audit IAM roles for unused permissions — AWS IAM Access Analyzer can show which permissions have not been used in the last 90 days.
- Disabling `automountServiceAccountToken` on pods that don't call the Kubernetes API also prevents OIDC token projection — verify the pod actually needs the service account token.

## Common Mistakes

- **Using node IAM role for pod access** — the node's IAM role (attached to EC2 instance) is accessible to all pods via IMDS without IRSA; use IRSA to scope permissions per pod.
- **Single IAM role for all pods in a namespace** — all pods then have the union of all permissions; per-pod roles enforce least privilege.
- **Not setting `automountServiceAccountToken: false` on non-IRSA pods** — pods that don't need cloud API access still get a projected service account token that could be misused.
- **Trust policy scoped to the namespace but not the service account** — `system:serviceaccount:production:*` allows any service account in the namespace to assume the role; always scope to the specific service account.
