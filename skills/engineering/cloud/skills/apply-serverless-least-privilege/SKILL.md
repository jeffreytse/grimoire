---
name: apply-serverless-least-privilege
description: Use when writing serverless functions — assigning per-function IAM execution roles with minimal permissions, disabling unused triggers, and securing function URLs and API authentication.
source: 'OWASP Serverless Top 10 SLS-2 SLS-3 SLS-4 (owasp.org/www-project-serverless-top-10/); AWS IAM best practices; CIS AWS Foundations Benchmark; Palo Alto Unit 42 serverless research'
tags: [security, owasp, serverless, lambda, iam, least-privilege, cloud, developer]
---

# Apply Serverless Least Privilege

Assign per-function IAM roles with resource-level permissions, disable unused event source mappings, require authentication on function URLs, and restrict deployment configuration — preventing lateral movement when any single function is compromised.

## Why This Is Best Practice

**Adopted by:** OWASP Serverless Top 10 SLS-4 (Over-Privileged Function Permissions), SLS-3 (Insecure Serverless Deployment Configuration), and SLS-2 (Broken Authentication). CIS AWS Foundations Benchmark v3.0 requires per-function execution roles and denies wildcard IAM actions. AWS Well-Architected Security Pillar mandates least-privilege IAM for all Lambda functions. Netflix's "Chaos Engineering for Serverless" (2022) found over-privileged Lambda roles as the most common critical finding in their Lambda security reviews.
**Impact:** Palo Alto Unit 42 (2022) found 70% of Lambda functions in production have excessive IAM permissions — most commonly `s3:*` on `*` resources when the function only needs `s3:GetObject` on one bucket. The 2022 Codecov incident (pipeline compromise, $hundreds of millions in remediation across affected organizations) occurred because CI functions had write access to all repos. A single Lambda with `iam:*` or `sts:AssumeRole` on `*` allows complete account takeover from that one function's compromise.
**Why best:** Shared execution roles across multiple functions are operationally convenient but mean that compromising any function exposes all resources accessible to that role. Per-function roles limit blast radius — a compromised order-processing Lambda can only access the order database, not payment records or user credentials.

Sources: OWASP Serverless Top 10 SLS-2, SLS-3, SLS-4; CIS AWS Foundations Benchmark v3.0; AWS Lambda security documentation; Netflix serverless security guidance (2022)

## Steps

1. **Create per-function execution roles — never share roles across functions**:

   ```yaml
   # SAM / CloudFormation: separate role per function
   ReadOrdersFunction:
     Type: AWS::Serverless::Function
     Properties:
       Role: !GetAtt ReadOrdersRole.Arn

   ReadOrdersRole:
     Type: AWS::IAM::Role
     Properties:
       AssumeRolePolicyDocument:
         Statement:
         - Effect: Allow
           Principal:
             Service: lambda.amazonaws.com
           Action: sts:AssumeRole
       Policies:
       - PolicyDocument:
           Statement:
           # Specific table, specific actions only
           - Effect: Allow
             Action:
             - dynamodb:GetItem
             - dynamodb:Query
             Resource: !GetAtt OrdersTable.Arn
           # CloudWatch Logs — required for all functions
           - Effect: Allow
             Action:
             - logs:CreateLogGroup
             - logs:CreateLogStream
             - logs:PutLogEvents
             Resource: !Sub "arn:aws:logs:${AWS::Region}:${AWS::AccountId}:log-group:/aws/lambda/ReadOrders*"
   ```

2. **Use resource-level permissions — never `Resource: "*"`**:

   ```python
   # Terraform example — specific ARN permissions
   resource "aws_iam_role_policy" "process_payment" {
     role = aws_iam_role.process_payment_lambda.id
     policy = jsonencode({
       Version = "2012-10-17"
       Statement = [
         {
           Effect = "Allow"
           Action = ["secretsmanager:GetSecretValue"]
           Resource = aws_secretsmanager_secret.payment_api_key.arn
           # NOT: Resource = "*"
         },
         {
           Effect = "Allow"
           Action = ["sqs:ReceiveMessage", "sqs:DeleteMessage", "sqs:GetQueueAttributes"]
           Resource = aws_sqs_queue.payment_queue.arn
           # Specific queue, not all SQS
         },
         {
           Effect = "Allow"
           Action = ["kms:Decrypt"]
           Resource = aws_kms_key.payment_key.arn
           # Only the payment KMS key
         }
       ]
     })
   }
   ```

3. **Require authentication on all invocation paths**:

   ```yaml
   # API Gateway: require Cognito or IAM auth on all routes
   # BAD — no auth
   Events:
     ApiEvent:
       Type: Api
       Properties:
         Path: /orders
         Method: GET
         # Auth: NONE — anyone can call

   # GOOD — require JWT authorization
   Events:
     ApiEvent:
       Type: HttpApi
       Properties:
         Path: /orders
         Method: GET
         Auth:
           Authorizer: CognitoAuthorizer

   # For Lambda Function URLs — require IAM auth (not NONE)
   FunctionUrlConfig:
     AuthType: AWS_IAM  # NOT: AuthType: NONE
     Cors:
       AllowOrigins:
       - "https://app.company.com"  # specific origin, not "*"
   ```

4. **Disable unused event source mappings and triggers**:

   ```bash
   # List all event source mappings for a function
   aws lambda list-event-source-mappings --function-name my-function

   # Disable unused mappings (don't delete — may be needed again)
   aws lambda update-event-source-mapping \
     --uuid <mapping-uuid> \
     --enabled false

   # Review Lambda resource-based policies — remove stale invoke permissions
   aws lambda get-policy --function-name my-function
   aws lambda remove-permission \
     --function-name my-function \
     --statement-id old-s3-trigger
   ```

5. **Restrict deployment: prevent unauthorized function code updates**:

   ```yaml
   # IAM policy for deployment role — only CI/CD service account has lambda:UpdateFunctionCode
   DeploymentPolicy:
     Statement:
     - Effect: Allow
       Action:
       - lambda:UpdateFunctionCode
       - lambda:UpdateFunctionConfiguration
       Resource: "arn:aws:lambda:*:*:function:*"
       Condition:
         StringEquals:
           aws:PrincipalTag/Role: "cicd-deployer"
   ```

   ```bash
   # Enable Lambda code signing — only signed deployments accepted
   aws lambda create-code-signing-config \
     --description "production-signing" \
     --allowed-publishers "SigningProfileVersionArns=arn:aws:signer:..."
     --code-signing-policies "UntrustedArtifactOnDeployment=Enforce"

   aws lambda put-function-code-signing-config \
     --function-name my-function \
     --code-signing-config-arn arn:aws:lambda:...
   ```

6. **Use IAM Access Analyzer to detect over-permissive policies**:

   ```bash
   # Scan Lambda execution roles for unused permissions
   aws accessanalyzer validate-policy \
     --policy-document file://lambda-role-policy.json \
     --policy-type IDENTITY_POLICY

   # Generate least-privilege policy from CloudTrail activity
   aws iam generate-service-last-accessed-details \
     --arn arn:aws:iam::123456789:role/MyLambdaRole

   # Use this output to prune permissions not used in last 90 days
   ```

## Rules

- Never attach `AdministratorAccess` or `AWSLambdaFullAccess` managed policies to a Lambda execution role.
- Lambda Function URLs with `AuthType: NONE` are publicly callable by anyone on the internet — require IAM or implement custom authorizer.
- Remove `lambda:InvokeFunction` resource-based policy statements for services no longer triggering the function — stale principals can still invoke.
- Enable AWS CloudTrail for all regions to audit Lambda invocations and IAM role assumptions.

## Common Mistakes

- **One shared execution role for all Lambda functions** — common in starter templates (`LambdaExecutionRole` used everywhere); creates shared blast radius.
- **`s3:*` on `*` because "Lambda needs S3 access"** — specify the exact bucket ARN and limit to the actions the function actually performs (GetObject, PutObject, etc.).
- **`logs:*` on `*`** — Lambda only needs logs permissions for its own log group; scope to `/aws/lambda/FunctionName*`.
- **Not monitoring for IAM privilege escalation** — a Lambda with `iam:AttachRolePolicy` or `iam:PutRolePolicy` on `*` can grant itself admin; use IAM Access Analyzer to detect these.
