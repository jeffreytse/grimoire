---
name: design-serverless-cost-protection
description: Use when deploying serverless functions at scale — setting reserved concurrency, SQS visibility timeouts, recursion detection, and budget alerts to prevent denial-of-service via financial resource exhaustion.
source: 'OWASP Serverless Top 10 SLS-8 (owasp.org/www-project-serverless-top-10/); AWS Lambda reserved concurrency documentation; AWS re:Invent "Serverless Security" (2022); Datadog State of Serverless 2023'
tags: [security, owasp, serverless, lambda, cost-protection, dos, concurrency, cloud, developer]
---

# Design Serverless Cost Protection

Set reserved concurrency limits, configure SQS dead-letter queues with backoff, enable Lambda Recursion Detection, and configure AWS Budgets alerts — preventing runaway functions and financial denial-of-service where an attacker or bug can generate millions of invocations and thousands of dollars in minutes.

## Why This Is Best Practice

**Adopted by:** OWASP Serverless Top 10 SLS-8 (Denial of Service and Financial Resource Exhaustion). AWS Lambda reserved concurrency is documented as the primary defense against both DoS and runaway cost. AWS re:Invent 2022 "Serverless Security" session lists SLS-8 as the most financially damaging class of serverless incident. Datadog's "State of Serverless" (2023) found that 12% of organizations reported unexpected Lambda cost spikes of >300% month-over-month, with most caused by event loop bugs or abuse of public-facing endpoints.
**Impact:** Unlike EC2 (fixed cost per instance), Lambda scales automatically to handle any traffic volume — a single misconfigured SQS trigger in an infinite retry loop was reported by a startup founder in 2019 to have generated $72,000 in Lambda and DynamoDB charges in one weekend. AWS documented a case where a Lambda recursive loop caused 100,000+ invocations per second for 3 hours before the account limit was hit. Public-facing API Gateway endpoints without rate limits are directly exploitable — a botnet sending 10,000 requests/second at $0.0000002/invocation costs $2,000/hour.
**Why best:** Reserved concurrency provides a hard cap — functions cannot scale beyond the reservation, preventing both cost exhaustion and cascade failures where one function consumes all account concurrency. The alternative (account-level concurrency limit only) allows a single function to starve all others in the account. Financial alerts provide a safety net when limits are misconfigured.

Sources: OWASP Serverless Top 10 SLS-8; AWS Lambda documentation on reserved concurrency; AWS re:Invent SVS401 "Serverless Security" (2022); Datadog State of Serverless (2023)

## Steps

1. **Set reserved concurrency on all production functions**:

   ```yaml
   # SAM template — reserved concurrency per function
   ProcessOrderFunction:
     Type: AWS::Serverless::Function
     Properties:
       ReservedConcurrencyLimit: 100  # hard cap: max 100 concurrent executions
       # Never leave at -1 (unlimited) for functions with external triggers

   PublicApiFunction:
     Type: AWS::Serverless::Function
     Properties:
       ReservedConcurrencyLimit: 50   # public endpoint: lower cap
   ```

   ```bash
   # Set via CLI
   aws lambda put-function-concurrency \
     --function-name my-function \
     --reserved-concurrent-executions 100
   ```

   Calculate appropriate limit: `max_rps × avg_duration_seconds = required_concurrency`
   Add 20% headroom: `required_concurrency × 1.2 = reserved_limit`

2. **Prevent SQS-triggered Lambda infinite loops with DLQ and backoff**:

   ```yaml
   # SQS queue with DLQ — prevents infinite retry loops
   OrderQueue:
     Type: AWS::SQS::Queue
     Properties:
       VisibilityTimeout: 90  # must be ≥ Lambda timeout × 6
       RedrivePolicy:
         deadLetterTargetArn: !GetAtt OrderDLQ.Arn
         maxReceiveCount: 3   # after 3 failures, move to DLQ — stop retrying

   OrderDLQ:
     Type: AWS::SQS::Queue
     Properties:
       MessageRetentionPeriod: 1209600  # 14 days — for investigation

   # Lambda event source mapping — limit batch size
   OrderProcessorMapping:
     Type: AWS::Lambda::EventSourceMapping
     Properties:
       FunctionName: !GetAtt ProcessOrderFunction.Arn
       EventSourceArn: !GetAtt OrderQueue.Arn
       BatchSize: 10          # process 10 at a time, not 10,000
       MaximumBatchingWindowInSeconds: 5
       FunctionResponseTypes:
       - ReportBatchItemFailures  # partial batch failure — don't re-process successes
   ```

3. **Enable Lambda Recursion Detection (prevents Lambda-to-Lambda infinite loops)**:

   ```bash
   # Enable recursive loop detection on all Lambda functions (AWS default since 2023)
   aws lambda put-function-recursion-config \
     --function-name my-function \
     --recursive-loop Terminate

   # Verify
   aws lambda get-function-recursion-config --function-name my-function
   # Expected: {"RecursiveLoop": "Terminate"}
   ```

   ```python
   # Also implement application-level recursion guard
   import os

   MAX_RECURSION_DEPTH = 10

   def lambda_handler(event, context):
       # Check for recursion via custom header/attribute
       depth = int(event.get("_recursion_depth", 0))
       if depth >= MAX_RECURSION_DEPTH:
           raise RuntimeError(f"Maximum recursion depth {MAX_RECURSION_DEPTH} exceeded")

       # Pass depth to downstream invocations
       invoke_next_function(event={**event, "_recursion_depth": depth + 1})
   ```

4. **Rate limit public API Gateway endpoints**:

   ```yaml
   # API Gateway usage plan — rate + burst limits per API key or per stage
   ApiUsagePlan:
     Type: AWS::ApiGateway::UsagePlan
     Properties:
       UsagePlanName: standard
       Throttle:
         RateLimit: 1000   # requests per second
         BurstLimit: 2000  # burst capacity
       Quota:
         Limit: 100000     # daily request quota
         Period: DAY

   # Stage-level throttling (applies to all endpoints)
   ApiStage:
     Type: AWS::ApiGateway::Stage
     Properties:
       DefaultRouteSettings:
         ThrottlingRateLimit: 1000
         ThrottlingBurstLimit: 500
   ```

5. **Configure AWS Budgets alerts for Lambda cost anomalies**:

   ```python
   # CDK / Python: programmatic budget creation
   import boto3

   budgets = boto3.client("budgets")

   budgets.create_budget(
       AccountId="123456789012",
       Budget={
           "BudgetName": "lambda-daily-limit",
           "BudgetType": "COST",
           "TimeUnit": "DAILY",
           "BudgetLimit": {"Amount": "100", "Unit": "USD"},
           "CostFilters": {"Service": ["AWS Lambda"]},
       },
       NotificationsWithSubscribers=[
           {
               "Notification": {
                   "NotificationType": "ACTUAL",
                   "ComparisonOperator": "GREATER_THAN",
                   "Threshold": 80,  # alert at 80% of daily limit
               },
               "Subscribers": [
                   {"SubscriptionType": "SNS", "Address": "arn:aws:sns:...alert-topic"},
               ],
           },
           {
               "Notification": {
                   "NotificationType": "FORECASTED",
                   "ComparisonOperator": "GREATER_THAN",
                   "Threshold": 100,  # alert when forecasted to exceed limit
               },
               "Subscribers": [
                   {"SubscriptionType": "SNS", "Address": "arn:aws:sns:...alert-topic"},
               ],
           },
       ],
   )
   ```

## Rules

- Every function with an external trigger (API Gateway, SQS, SNS, S3) must have a reserved concurrency limit — unlimited scaling is a financial risk, not just a performance feature.
- SQS → Lambda event source mappings must have a DLQ with `maxReceiveCount` ≤ 5 — without DLQ, a poisoned message triggers infinite retries at full concurrency.
- Account-level Lambda concurrency limit (default 1000) is shared across all functions — reserve concurrency per critical function so one runaway function can't starve all others.
- Monitor the SQS DLQ depth — a growing DLQ indicates a processing failure that should trigger an alert, not just silent message accumulation.

## Common Mistakes

- **`ReservedConcurrencyLimit: -1` (unlimited) on all functions** — the SAM default; acceptable during development but must be set before production deployment.
- **SQS without DLQ** — a message that consistently fails to process retries indefinitely until the retention period expires (4 days default), generating continuous Lambda invocations.
- **Not accounting for downstream throttling** — a Lambda with reserved concurrency of 1000 calling DynamoDB at 1000 RPS will hit DynamoDB throughput limits; size the entire call chain, not just Lambda.
- **Budget alerts sent to email only** — email alerts can be missed; route to SNS → PagerDuty/OpsGenie for financial anomalies that represent potential security incidents.
