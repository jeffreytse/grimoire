---
name: apply-serverless-secrets
description: Use when managing secrets in serverless functions — replacing environment variable credentials with Secrets Manager or Parameter Store, auditing third-party Lambda layers, and preventing secret leakage in logs.
source: 'OWASP Serverless Top 10 SLS-6 SLS-7 (owasp.org/www-project-serverless-top-10/); AWS Secrets Manager documentation; Snyk State of Serverless Security 2022'
tags: [security, owasp, serverless, lambda, secrets, dependencies, cloud, developer]
---

# Apply Serverless Secrets

Replace environment variable secrets with AWS Secrets Manager or SSM Parameter Store, cache secrets with TTL to avoid per-invocation latency, audit Lambda layers for supply chain risks, and purge secrets from function logs — preventing credential exposure in configuration and execution artifacts.

## Why This Is Best Practice

**Adopted by:** OWASP Serverless Top 10 SLS-7 (Insecure Application Secrets Management) and SLS-6 (Insecure Third-Party Dependencies). AWS best practices documentation explicitly recommends against storing secrets in environment variables. HashiCorp's 2022 State of Security survey found 68% of organizations using Lambda store at least one database credential as an environment variable. Snyk's "State of Serverless Security" (2022) found 47% of serverless functions have secrets in environment variables, which are readable by any IAM principal with `lambda:GetFunctionConfiguration`.
**Impact:** Lambda environment variables are not encrypted by default at the configuration layer — they're visible in the AWS Console, `aws lambda get-function-configuration`, and any CloudFormation drift detection output. Anyone with `lambda:GetFunctionConfiguration` permission can read all environment variables including database passwords and API keys. In 2021, a misconfigured AWS Config rule at a financial institution exposed all Lambda environment variables to an internal logging system, compromising 43 production database credentials. Lambda layers from public repositories can include malicious code — 2022 saw multiple supply chain attacks via compromised npm packages published as Lambda layers.
**Why best:** Secrets Manager provides automatic rotation, fine-grained IAM access per secret (vs. reading all env vars), versioning, and an audit trail via CloudTrail. The alternative (KMS-encrypted environment variables) still stores ciphertext in the function configuration — better but requires KMS API calls to decrypt and doesn't support rotation.

Sources: OWASP Serverless Top 10 SLS-6, SLS-7; AWS Secrets Manager documentation; Snyk State of Serverless Security (2022); HashiCorp State of Security Survey (2022)

## Steps

1. **Replace environment variable secrets with Secrets Manager**:

   ```python
   import boto3
   import json
   import os
   import time

   # Cached secrets client — reuse across warm invocations
   _secrets_client = boto3.client("secretsmanager", region_name=os.environ["AWS_REGION"])
   _cache: dict = {}
   CACHE_TTL = 300  # 5 minutes

   def get_secret(secret_name: str) -> dict:
       cached = _cache.get(secret_name)
       if cached and time.time() - cached["fetched_at"] < CACHE_TTL:
           return cached["value"]

       response = _secrets_client.get_secret_value(SecretId=secret_name)
       value = json.loads(response["SecretString"])
       _cache[secret_name] = {"value": value, "fetched_at": time.time()}
       return value

   def lambda_handler(event, context):
       # Fetch at runtime, not from env var
       db_config = get_secret("/production/database")
       conn = connect_db(
           host=db_config["host"],
           password=db_config["password"]
       )
   ```

   Environment variable to keep: only non-secret config like `SECRET_NAME=/production/database` (the name, not the value).

2. **Use SSM Parameter Store for non-secret configuration**:

   ```python
   import boto3

   ssm = boto3.client("ssm")

   def get_parameter(name: str, decrypt: bool = True) -> str:
       response = ssm.get_parameter(Name=name, WithDecryption=decrypt)
       return response["Parameter"]["Value"]

   # For SecureString parameters (encrypted with KMS)
   api_endpoint = get_parameter("/production/payment-api/endpoint", decrypt=False)
   api_key = get_parameter("/production/payment-api/key", decrypt=True)
   ```

3. **Audit Lambda layers — only use trusted sources**:

   ```bash
   # List layers used by a function
   aws lambda get-function-configuration --function-name my-function \
     --query 'Layers[*].Arn'

   # Check layer source — verify publisher account ID
   LAYER_ARN="arn:aws:lambda:us-east-1:123456789:layer:my-layer:1"
   ACCOUNT=$(echo $LAYER_ARN | cut -d: -f5)
   echo "Layer owned by account: $ACCOUNT"
   # Verify this is your account or a known trusted publisher

   # Scan layer contents for vulnerabilities
   aws lambda get-layer-version-by-arn --arn $LAYER_ARN
   # Download and scan with Trivy or Snyk
   ```

   ```yaml
   # SAM/CloudFormation: pin layer versions — don't use "latest"
   Layers:
   - !Sub "arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:layer:my-layer:42"
   # NOT: arn:...:my-layer:$LATEST (mutable)
   ```

4. **Prevent secrets from appearing in logs**:

   ```python
   import logging
   import re

   # Custom formatter that scrubs sensitive patterns from log output
   class SecretScrubber(logging.Formatter):
       PATTERNS = [
           (re.compile(r'(password["\s:=]+)[^\s,}"\']+()', re.I), r'\1[REDACTED]\2'),
           (re.compile(r'(api[_-]?key["\s:=]+)[^\s,}"\']+', re.I), r'\1[REDACTED]'),
           (re.compile(r'(Bearer\s+)\S+'), r'\1[REDACTED]'),
           (re.compile(r'(Authorization["\s:=]+)[^\s,}"\']+', re.I), r'\1[REDACTED]'),
       ]

       def format(self, record):
           msg = super().format(record)
           for pattern, replacement in self.PATTERNS:
               msg = pattern.sub(replacement, msg)
           return msg

   logger = logging.getLogger()
   handler = logging.StreamHandler()
   handler.setFormatter(SecretScrubber())
   logger.addHandler(handler)
   logger.setLevel(logging.INFO)

   # Never log request bodies or response bodies that may contain credentials
   # logger.info(f"Processing request: {json.dumps(event)}")  # BAD
   logger.info(f"Processing request: method={event.get('httpMethod')}, path={event.get('path')}")
   ```

5. **Rotate secrets automatically — configure rotation in Secrets Manager**:

   ```python
   # Lambda rotation function (Secrets Manager calls this)
   def lambda_handler(event, context):
       step = event["Step"]
       secret_id = event["SecretId"]
       token = event["ClientRequestToken"]

       if step == "createSecret":
           # Generate new credential
           new_password = generate_secure_password()
           client.put_secret_value(
               SecretId=secret_id,
               ClientRequestToken=token,
               SecretString=json.dumps({"password": new_password}),
               VersionStages=["AWSPENDING"]
           )
       elif step == "setSecret":
           # Apply new credential to the service (e.g., update DB password)
           pending = client.get_secret_value(SecretId=secret_id, VersionStage="AWSPENDING")
           update_database_password(json.loads(pending["SecretString"])["password"])
       elif step == "finishSecret":
           # Promote AWSPENDING to AWSCURRENT
           client.update_secret_version_stage(
               SecretId=secret_id,
               VersionStage="AWSCURRENT",
               MoveToVersionId=token
           )
   ```

## Rules

- Never set database passwords, API keys, or tokens as Lambda environment variables — use Secrets Manager or Parameter Store.
- Cache secrets in the function's module-level scope (reused across warm invocations) but with a TTL — stale cache shouldn't persist rotated credentials indefinitely.
- Lambda layers used from public ARNs are third-party code running in your execution environment — audit them as you would any dependency.
- Set `POWERTOOLS_LOG_LEVEL=INFO` (not DEBUG) in production — debug logging often includes full event payloads which may contain tokens.

## Common Mistakes

- **Fetching secrets inside the handler on every invocation** — without caching, every Lambda call makes a Secrets Manager API call, adding 20–100ms latency and incurring API costs at scale; cache with TTL.
- **Logging `event` at INFO level** — API Gateway events include `Authorization` headers and request bodies; log only the fields needed for debugging.
- **Using Lambda layers from unknown public accounts** — verify the account ID in the layer ARN belongs to a known publisher (AWS, a major library vendor, or your own account).
- **Hardcoding secret names as strings** — store secret ARNs in environment variables so the same code works across dev/staging/prod with different secret paths.
