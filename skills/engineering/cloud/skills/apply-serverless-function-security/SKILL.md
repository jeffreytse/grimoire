---
name: apply-serverless-function-security
description: Use when writing serverless functions (AWS Lambda, Google Cloud Functions, Azure Functions) — validating event data from all trigger sources, preventing execution flow manipulation, and securing error handling.
source: 'OWASP Serverless Top 10 SLS-1 SLS-9 SLS-10 (owasp.org/www-project-serverless-top-10/); AWS Lambda security documentation; Palo Alto Unit 42 serverless research'
tags: [security, owasp, serverless, lambda, event-injection, cloud, developer]
---

# Apply Serverless Function Security

Validate all event data at function entry regardless of trigger source, prevent execution flow manipulation by isolating functions to single responsibilities, and return generic error messages — preventing injection attacks that exploit the expanded attack surface of serverless event sources.

## Why This Is Best Practice

**Adopted by:** OWASP Serverless Top 10 SLS-1 (Function Event Data Injection), SLS-9 (Serverless Function Execution Flow Manipulation), and SLS-10 (Improper Exception Handling). AWS Lambda security documentation mandates input validation for all event sources. Palo Alto Unit 42's "Serverless Security: A Deep Dive" (2021) identified event injection as the primary novel attack vector in serverless architectures. Netflix, Airbnb, and Capital One — all heavy Lambda users — treat Lambda event payloads as untrusted input from their first public serverless security guidance.
**Impact:** Unlike web applications where requests come through a single HTTP gateway, Lambda functions receive events from S3, SQS, SNS, DynamoDB Streams, API Gateway, EventBridge, and IoT — each with a different event schema and different injection surfaces. Palo Alto Unit 42 demonstrated SQL injection via S3 object key names that triggered Lambda functions. The 2022 Lacework serverless threat report found 35% of Lambda vulnerabilities were injection flaws introduced through non-HTTP event sources that developers didn't treat as attack vectors.
**Why best:** Traditional web input validation (request body + query params) misses serverless-specific injection paths: filenames in S3 events, message bodies in SQS, record data in DynamoDB Streams, and environment variable injection. Treating every event field as untrusted and validating schema at function entry closes these gaps vs. validating only API Gateway parameters.

Sources: OWASP Serverless Top 10 SLS-1, SLS-9, SLS-10; AWS Lambda security best practices; Palo Alto Unit 42 Serverless Security (2021); Lacework Cloud Security Report (2022)

## Steps

1. **Validate all event fields regardless of trigger source**:

   ```python
   import json
   import re
   from pydantic import BaseModel, validator

   # S3 trigger — validate object key before using in file operations
   class S3EventRecord(BaseModel):
       bucket_name: str
       object_key: str

       @validator("object_key")
       def no_path_traversal(cls, v):
           # S3 keys can contain ../ sequences — validate before using as filesystem path
           if ".." in v or v.startswith("/"):
               raise ValueError("Invalid object key")
           if not re.match(r'^[\w\-./]+$', v):
               raise ValueError("Object key contains invalid characters")
           return v

   def lambda_handler(event, context):
       for record in event.get("Records", []):
           if record.get("eventSource") == "aws:s3":
               try:
                   validated = S3EventRecord(
                       bucket_name=record["s3"]["bucket"]["name"],
                       object_key=record["s3"]["object"]["key"]
                   )
               except Exception as e:
                   # Log detail internally, return generic error externally
                   print(f"[ERROR] Invalid S3 event: {e}")
                   return {"statusCode": 400, "body": "Invalid event"}
               process_s3_object(validated.bucket_name, validated.object_key)
   ```

2. **Validate SQS message bodies as untrusted input**:

   ```python
   import json
   from pydantic import BaseModel, constr

   class OrderMessage(BaseModel):
       order_id: constr(pattern=r'^[0-9a-f-]{36}$')  # UUID format only
       amount: float
       currency: constr(max_length=3, pattern=r'^[A-Z]{3}$')

       # Explicit constraint — prevent numeric overflow
       @validator("amount")
       def amount_range(cls, v):
           if v <= 0 or v > 1_000_000:
               raise ValueError("Amount out of range")
           return v

   def process_sqs_event(event, context):
       for record in event["Records"]:
           try:
               body = json.loads(record["body"])
               order = OrderMessage(**body)
           except Exception as e:
               print(f"[ERROR] Malformed SQS message: {e}, MessageId: {record['messageId']}")
               # Don't raise — let SQS retry; log for investigation
               continue
           process_order(order)
   ```

3. **Prevent execution flow manipulation — one function, one responsibility**:

   ```python
   # BAD — function accepts a "mode" parameter to switch behavior
   # Attacker controls execution path
   def lambda_handler(event, context):
       mode = event.get("mode", "read")
       if mode == "admin":          # attacker sets mode=admin
           return admin_operation()
       elif mode == "read":
           return read_operation()

   # GOOD — separate Lambda functions per operation; API GW route determines which executes
   # Function: read-user-lambda (only reads)
   def lambda_handler(event, context):
       user_id = validate_uuid(event["pathParameters"]["userId"])
       return {"statusCode": 200, "body": json.dumps(get_user(user_id))}

   # Function: update-user-lambda (only updates, separate IAM role)
   def lambda_handler(event, context):
       user_id = validate_uuid(event["pathParameters"]["userId"])
       body = UserUpdateSchema(**json.loads(event["body"]))
       update_user(user_id, body)
       return {"statusCode": 200, "body": "{}"}
   ```

4. **Secure error handling — no stack traces or internals in responses**:

   ```python
   import uuid
   import traceback

   def lambda_handler(event, context):
       error_id = str(uuid.uuid4())
       try:
           return handle_request(event)
       except ValidationError as e:
           # User-facing: generic message
           return {
               "statusCode": 400,
               "body": json.dumps({"error": "Invalid request", "error_id": error_id})
           }
       except Exception as e:
           # Internal: full detail to CloudWatch Logs only
           print(f"[ERROR] Unhandled exception error_id={error_id}: {traceback.format_exc()}")
           return {
               "statusCode": 500,
               "body": json.dumps({"error": "Internal server error", "error_id": error_id})
           }
           # error_id lets support correlate logs without exposing internals
   ```

5. **Set function timeout and memory to prevent resource exhaustion**:

   ```yaml
   # SAM template / CloudFormation
   ProcessOrderFunction:
     Type: AWS::Serverless::Function
     Properties:
       Timeout: 30        # max 30s — prevent hanging on slow upstream
       MemorySize: 256    # sufficient for task; prevents runaway memory use
       ReservedConcurrencyLimit: 100  # cap concurrent executions
       Environment:
         Variables:
           LOG_LEVEL: INFO  # never DEBUG in production
   ```

## Rules

- Every Lambda trigger source (S3, SQS, SNS, DynamoDB Streams, EventBridge, IoT, Cognito) is an untrusted input surface — validate schema before processing.
- Never use event fields directly in shell commands, SQL queries, or file paths without validation — `subprocess.run(event["command"])` is direct command injection.
- Dead Letter Queues (DLQ) must be configured for async invocations — without DLQ, failed events are silently dropped.
- Function names, ARNs, and environment variable names must not be included in error responses — they aid enumeration.

## Common Mistakes

- **Validating only API Gateway events** — S3, SQS, and DynamoDB Stream events are equally injectable; treat all event sources as untrusted.
- **Catching `Exception` and returning `str(e)`** — Python exception messages often contain table names, file paths, or SQL snippets; log internally, return generic message externally.
- **No DLQ on async-triggered functions** — if an S3-triggered Lambda fails, the event is silently lost; configure DLQ to SQS for investigation.
- **Large monolithic functions with many code paths** — a single Lambda doing read/write/admin based on an event field is SLS-9 (execution flow manipulation); split into separate functions.
