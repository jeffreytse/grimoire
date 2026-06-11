---
name: apply-iac-security
description: Use when writing Terraform, CloudFormation, Pulumi, or Bicep templates — to detect misconfigurations, enforce least privilege IAM, and prevent insecure defaults before infrastructure reaches production.
source: 'OWASP Infrastructure as Code Security Cheat Sheet (owasp.org/www-project-cheat-sheets); CIS AWS/GCP/Azure Benchmarks; Checkov documentation; tfsec documentation'
tags: [security, owasp, iac, terraform, cloudformation, cloud, devops, developer]
---

# Apply IaC Security

Scan Infrastructure as Code templates for misconfigurations, enforce least-privilege IAM, and validate security controls in CI — catching S3 public buckets, open security groups, and over-permissive roles before they reach production.

## Why This Is Best Practice

**Adopted by:** OWASP Infrastructure as Code Security Cheat Sheet defines the standard. CIS Benchmarks for AWS (v3.0), GCP (v2.0), and Azure (v2.0) are the authoritative hardening references. Bridgecrew (acquired by Palo Alto Networks), Checkov, and tfsec are used by Netflix, Shopify, and Stripe in their CI pipelines. AWS Security Hub and GCP Security Command Center both consume misconfig signals from IaC scanners.
**Impact:** Palo Alto Unit 42 (2021) found that 65% of cloud security incidents involved misconfigured IaC — public S3 buckets, overly permissive IAM roles, and unrestricted security groups. The 2019 Capital One breach (106M records) exploited an overly permissive IAM role on an EC2 instance. Checkov scans catch these classes in under 30 seconds before the `terraform apply` reaches production. CSPM tools report remediating IaC-caught misconfigs is 10× cheaper than post-deployment.
**Why best:** Runtime cloud security scanning (CSPM) finds misconfigurations after they're deployed — the IaC scanner finds them at PR time, before any resource exists. IaC scanning integrates with existing CI without requiring cloud credentials; CSPM requires a running environment and cloud API access. Combined approach (scan in CI + CSPM in production) provides defense-in-depth.

Sources: OWASP IaC Security Cheat Sheet; CIS AWS Benchmark v3.0; Bridgecrew State of Open Source Terraform Security Report (2022); Palo Alto Unit 42 Cloud Threat Report (2021)

## Steps

1. **Scan Terraform with Checkov in CI — fail the pipeline on HIGH/CRITICAL**:

   ```bash
   # Install
   pip install checkov

   # Scan a directory
   checkov -d ./infrastructure --framework terraform \
     --check CKV_AWS_18,CKV_AWS_19,CKV_AWS_20 \
     --compact --quiet

   # Fail on any HIGH or CRITICAL finding
   checkov -d ./infrastructure --soft-fail-on MEDIUM,LOW
   ```

   ```yaml
   # GitHub Actions
   - name: Checkov IaC scan
     uses: bridgecrewio/checkov-action@master
     with:
       directory: infrastructure/
       framework: terraform
       soft_fail: false
       output_format: github_failed_only
   ```

2. **Fix the most common Terraform misconfigurations**:

   ```hcl
   # S3: block public access + enable versioning + encrypt
   resource "aws_s3_bucket" "app_data" {
     bucket = "company-app-data"
   }

   resource "aws_s3_bucket_public_access_block" "app_data" {
     bucket                  = aws_s3_bucket.app_data.id
     block_public_acls       = true
     block_public_policy     = true
     ignore_public_acls      = true
     restrict_public_buckets = true
   }

   resource "aws_s3_bucket_server_side_encryption_configuration" "app_data" {
     bucket = aws_s3_bucket.app_data.id
     rule {
       apply_server_side_encryption_by_default {
         sse_algorithm = "aws:kms"
       }
     }
   }

   resource "aws_s3_bucket_versioning" "app_data" {
     bucket = aws_s3_bucket.app_data.id
     versioning_configuration { status = "Enabled" }
   }
   ```

3. **Enforce least-privilege IAM — no wildcard actions or resources**:

   ```hcl
   # BAD — wildcard action on all resources
   resource "aws_iam_policy" "bad_policy" {
     policy = jsonencode({
       Statement = [{
         Effect   = "Allow"
         Action   = "*"
         Resource = "*"
       }]
     })
   }

   # GOOD — specific actions on specific resource ARNs
   resource "aws_iam_policy" "app_policy" {
     policy = jsonencode({
       Version = "2012-10-17"
       Statement = [
         {
           Effect = "Allow"
           Action = [
             "s3:GetObject",
             "s3:PutObject"
           ]
           Resource = "${aws_s3_bucket.app_data.arn}/uploads/*"
         },
         {
           Effect = "Allow"
           Action = ["secretsmanager:GetSecretValue"]
           Resource = aws_secretsmanager_secret.db_password.arn
         }
       ]
     })
   }
   ```

4. **Lock down security groups — no `0.0.0.0/0` ingress except ports 80/443**:

   ```hcl
   resource "aws_security_group" "app" {
     name   = "app-sg"
     vpc_id = aws_vpc.main.id

     # Allow HTTPS from anywhere
     ingress {
       from_port   = 443
       to_port     = 443
       protocol    = "tcp"
       cidr_blocks = ["0.0.0.0/0"]
     }

     # Allow internal communication only
     ingress {
       from_port       = 8080
       to_port         = 8080
       protocol        = "tcp"
       security_groups = [aws_security_group.alb.id]
       # NOT: cidr_blocks = ["0.0.0.0/0"]
     }

     # Restrict egress to known destinations
     egress {
       from_port   = 443
       to_port     = 443
       protocol    = "tcp"
       cidr_blocks = ["0.0.0.0/0"]
     }
   }
   ```

5. **Never store secrets in IaC — use secret stores**:

   ```hcl
   # BAD — secret in tfvars or hardcoded
   resource "aws_db_instance" "main" {
     password = "supersecret123"  # visible in state file
   }

   # GOOD — fetch from Secrets Manager at runtime
   data "aws_secretsmanager_secret_version" "db_password" {
     secret_id = "/production/db/password"
   }

   resource "aws_db_instance" "main" {
     password = data.aws_secretsmanager_secret_version.db_password.secret_string
   }
   ```

   Terraform state also stores secrets — use remote state with encryption:
   ```hcl
   terraform {
     backend "s3" {
       bucket  = "terraform-state-company"
       key     = "production/terraform.tfstate"
       encrypt = true
       kms_key_id = "arn:aws:kms:..."
     }
   }
   ```

6. **Pin provider versions and module sources**:

   ```hcl
   terraform {
     required_providers {
       aws = {
         source  = "hashicorp/aws"
         version = "~> 5.0"  # no unpinned "latest"
       }
     }
   }

   # Pin modules to a specific tag — not a branch
   module "vpc" {
     source  = "terraform-aws-modules/vpc/aws"
     version = "5.1.2"  # not: "main" or "latest"
   }
   ```

## Rules

- Never commit `terraform.tfstate` or `*.tfvars` containing secrets — add both to `.gitignore`.
- Run `terraform plan` before `apply` and review the diff — automated `apply` without review risks unintended resource deletion.
- Use separate Terraform workspaces or directories for dev/staging/production — never use the same state file.
- Enable MFA delete on the S3 bucket storing Terraform state to prevent accidental or malicious deletion.

## Common Mistakes

- **`0.0.0.0/0` on port 22 or 3389** — SSH/RDP open to the internet is the #1 cloud misconfiguration; use Systems Manager Session Manager instead.
- **IAM roles with `sts:AssumeRole` for `*`** — allows any AWS principal to assume the role; always specify the trusted principal ARN.
- **Ignoring `checkov` violations with `#checkov:skip`** — each skip should have a documented justification and tracked exception.
- **Not encrypting RDS at rest** — `storage_encrypted = true` should be the default for any database with sensitive data.
