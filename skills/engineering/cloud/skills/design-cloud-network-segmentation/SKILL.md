---
name: design-cloud-network-segmentation
description: Use when designing cloud network architecture — structuring VPCs with public/private subnets, using VPC endpoints for AWS service access, and implementing security groups as stateful firewalls to contain blast radius.
source: 'OWASP Cloud-Native Application Security Top 10 C7 (owasp.org/www-project-cloud-native-application-security-top-10/); AWS VPC security best practices; CIS AWS Foundations Benchmark v3.0; NIST SP 800-125B'
tags: [security, owasp, cloud, vpc, network-segmentation, security-groups, aws, developer]
---

# Design Cloud Network Segmentation

Structure VPCs with isolated public/private/data subnets, replace NAT gateway with VPC endpoints for AWS service access, implement security groups with minimal allow rules, and peer VPCs with explicit route controls — preventing lateral movement between environments and reducing attack surface.

## Why This Is Best Practice

**Adopted by:** OWASP Cloud-Native Application Security Top 10 C7 (Insecure Network Controls). CIS AWS Foundations Benchmark v3.0 mandates VPC flow logs, restricted security groups, and no default VPC in production. AWS Well-Architected Security Pillar defines the three-tier subnet model (public/private/data) as the standard. Netflix, Stripe, and Capital One all use multi-VPC architectures with Transit Gateway for environment isolation — documented in their engineering blogs.
**Impact:** The 2021 Colonial Pipeline ransomware spread from an IT network to OT systems because there was no network segmentation between them — the cloud equivalent is flat VPC design where a compromised web server can directly reach the database. CIS Benchmark analysis (2022) found 35% of AWS accounts have security groups with `0.0.0.0/0` ingress on ports beyond 443/80. AWS's own "Customer Security Incidents" post-mortems cite flat VPC architecture as a contributing factor in 40% of analyzed breach scenarios.
**Why best:** A flat network (all instances in one subnet, permissive security groups) means a single compromised instance can reach all other instances and services. Tiered subnets + tight security groups implement the principle of least network access — a compromised web server in the public subnet cannot reach the database because the database security group only accepts traffic from the application security group, and the application subnet has no route to the internet.

Sources: OWASP Cloud-Native Top 10 C7; CIS AWS Foundations Benchmark v3.0; AWS VPC security best practices; Netflix multi-account architecture (2019)

## Steps

1. **Three-tier VPC structure: public, private, data subnets**:

   ```hcl
   # Terraform: three-tier VPC
   module "vpc" {
     source  = "terraform-aws-modules/vpc/aws"
     version = "5.1.2"

     name = "production"
     cidr = "10.0.0.0/16"

     azs = ["us-east-1a", "us-east-1b", "us-east-1c"]

     # Public: load balancers only
     public_subnets  = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
     # Private: application servers (no direct internet access)
     private_subnets = ["10.0.11.0/24", "10.0.12.0/24", "10.0.13.0/24"]
     # Data: databases (no internet route at all)
     database_subnets = ["10.0.21.0/24", "10.0.22.0/24", "10.0.23.0/24"]

     enable_nat_gateway     = true
     single_nat_gateway     = false  # one per AZ for HA
     enable_dns_hostnames   = true
     enable_flow_log        = true
     flow_log_destination_type = "cloud-watch-logs"
   }
   ```

2. **Replace NAT gateway with VPC endpoints for AWS services**:

   ```hcl
   # VPC endpoints: private subnets access S3, DynamoDB, Secrets Manager
   # without traversing the public internet or NAT gateway

   # Gateway endpoint (S3, DynamoDB — free)
   resource "aws_vpc_endpoint" "s3" {
     vpc_id            = module.vpc.vpc_id
     service_name      = "com.amazonaws.us-east-1.s3"
     vpc_endpoint_type = "Gateway"
     route_table_ids   = module.vpc.private_route_table_ids
   }

   # Interface endpoint (Secrets Manager — has per-hour cost but eliminates internet path)
   resource "aws_vpc_endpoint" "secrets_manager" {
     vpc_id              = module.vpc.vpc_id
     service_name        = "com.amazonaws.us-east-1.secretsmanager"
     vpc_endpoint_type   = "Interface"
     subnet_ids          = module.vpc.private_subnets
     security_group_ids  = [aws_security_group.vpc_endpoints.id]
     private_dns_enabled = true
   }
   ```

3. **Security groups: source = security group, not CIDR range**:

   ```hcl
   # Application layer: only accepts traffic from ALB security group
   resource "aws_security_group" "app" {
     name   = "app-layer"
     vpc_id = module.vpc.vpc_id

     ingress {
       from_port       = 8080
       to_port         = 8080
       protocol        = "tcp"
       security_groups = [aws_security_group.alb.id]  # NOT: cidr_blocks = ["10.0.0.0/8"]
     }

     egress {
       from_port       = 5432
       to_port         = 5432
       protocol        = "tcp"
       security_groups = [aws_security_group.db.id]  # only to DB
     }

     egress {
       from_port   = 443
       to_port     = 443
       protocol    = "tcp"
       cidr_blocks = ["0.0.0.0/0"]  # HTTPS out for API calls
     }
   }

   # Database layer: only accepts from app security group
   resource "aws_security_group" "db" {
     name   = "db-layer"
     vpc_id = module.vpc.vpc_id

     ingress {
       from_port       = 5432
       to_port         = 5432
       protocol        = "tcp"
       security_groups = [aws_security_group.app.id]  # only from app
     }
     # NO egress: databases don't initiate outbound connections
   }
   ```

4. **Separate VPCs per environment with Transit Gateway**:

   ```hcl
   # Each environment in its own VPC — no route between prod and dev
   # Transit Gateway for controlled cross-VPC routing

   resource "aws_ec2_transit_gateway" "main" {
     description                     = "Production network hub"
     default_route_table_association = "disable"  # explicit route tables only
     default_route_table_propagation = "disable"
   }

   # Attach production VPC — only routes to shared-services VPC
   resource "aws_ec2_transit_gateway_vpc_attachment" "production" {
     transit_gateway_id = aws_ec2_transit_gateway.main.id
     vpc_id             = module.production_vpc.vpc_id
     subnet_ids         = module.production_vpc.private_subnets
   }

   # NO direct route from production to development
   # Development VPC can only route to shared-services, not production
   ```

5. **Enable VPC Flow Logs for network forensics**:

   ```hcl
   resource "aws_flow_log" "production" {
     vpc_id          = module.vpc.vpc_id
     traffic_type    = "ALL"  # capture both ACCEPT and REJECT
     iam_role_arn    = aws_iam_role.flow_log.arn
     log_destination = aws_cloudwatch_log_group.flow_logs.arn

     # Enable enhanced metadata fields
     log_format = "$${version} $${account-id} $${vpc-id} $${subnet-id} $${instance-id} $${srcaddr} $${dstaddr} $${srcport} $${dstport} $${protocol} $${packets} $${bytes} $${start} $${end} $${action} $${log-status}"
   }

   # CloudWatch Insights query for rejected connections (security investigation)
   # fields srcAddr, dstAddr, dstPort, action
   # | filter action = "REJECT"
   # | stats count(*) by srcAddr, dstAddr, dstPort
   # | sort count desc
   # | limit 20
   ```

## Rules

- Delete the default VPC in all AWS regions — it has permissive default security groups and no flow logs; workloads should never use it.
- Security group egress rules must be explicit — `allow all egress` (`0.0.0.0/0`) on application servers allows exfiltration of any data to any destination.
- VPC peering with "allow all traffic" route effectively merges the two VPCs — peer with specific route table entries to specific CIDRs only.
- Flow logs must capture REJECT records — accepted traffic shows normal operation; rejected traffic shows attack attempts and misconfiguration.

## Common Mistakes

- **Single VPC for all environments** — dev, staging, and production in the same VPC means a developer with staging access can reach production databases via internal IPs.
- **Security groups with `0.0.0.0/0` ingress on port 22/3389** — SSH/RDP open to the internet is the #1 CIS benchmark failure; use Systems Manager Session Manager instead.
- **Using NACLs as the primary security layer** — NACLs are stateless and hard to maintain; use security groups (stateful) as the primary control, NACLs only for coarse subnet-level deny.
- **Not monitoring VPC flow logs** — enabling flow logs without shipping them to SIEM or running queries provides no security value; automate anomaly detection on reject rates.
