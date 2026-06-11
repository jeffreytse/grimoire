---
name: design-serverless-observability
description: Use when operating serverless functions in production — implementing structured logging, distributed tracing, anomaly-based alerting, and cold start monitoring to detect security incidents and performance degradation.
source: 'OWASP Serverless Top 10 SLS-5 (owasp.org/www-project-serverless-top-10/); AWS Lambda Powertools documentation; AWS Well-Architected Operational Excellence Pillar; Datadog State of Serverless 2023'
tags: [security, owasp, serverless, lambda, observability, logging, tracing, cloud, developer]
---

# Design Serverless Observability

Implement structured JSON logs with correlation IDs, distributed tracing with AWS X-Ray, anomaly-based CloudWatch alarms, and security-relevant event logging — enabling incident detection and forensic investigation in serverless architectures where traditional server logs don't exist.

## Why This Is Best Practice

**Adopted by:** OWASP Serverless Top 10 SLS-5 (Inadequate Function Monitoring and Logging). AWS Lambda Powertools (used by thousands of teams for production Lambda observability) is the reference implementation. AWS Well-Architected Operational Excellence Pillar mandates distributed tracing for all serverless workloads. Datadog's "State of Serverless" (2023) found that teams using structured logging and distributed tracing resolve incidents 3× faster than those using unstructured logs.
**Impact:** Unlike EC2/ECS, serverless functions produce logs only in CloudWatch — there's no server to SSH into, no persistent process to attach a debugger to, and no network-level visibility without explicit instrumentation. The 2022 Lacework serverless threat research found that 60% of Lambda security incidents were detected only through billing anomalies (sudden cost spikes from unexpected invocations) rather than from application logs — because most teams lack structured Lambda security logging. Without correlation IDs, a multi-function execution chain is nearly impossible to trace after an incident.
**Why best:** Unstructured print/console logs are unsearchable at scale and can't be correlated across functions. Structured JSON logs with consistent fields (request_id, user_id, function_name, duration) enable CloudWatch Insights queries, automated anomaly detection, and SIEM ingestion. Distributed tracing with X-Ray links a user request through API Gateway → Lambda A → DynamoDB → Lambda B — critical for security forensics in event-driven architectures.

Sources: OWASP Serverless Top 10 SLS-5; AWS Lambda Powertools documentation; AWS Well-Architected Operational Excellence Pillar; Datadog State of Serverless (2023)

## Steps

1. **Use AWS Lambda Powertools for structured logging**:

   ```python
   from aws_lambda_powertools import Logger, Tracer, Metrics
   from aws_lambda_powertools.metrics import MetricUnit
   from aws_lambda_powertools.utilities.typing import LambdaContext

   logger = Logger(service="order-processor")
   tracer = Tracer(service="order-processor")
   metrics = Metrics(namespace="OrderService", service="order-processor")

   @logger.inject_lambda_context(correlation_id_path="requestContext.requestId", log_event=False)
   @tracer.capture_lambda_handler
   @metrics.log_metrics(capture_cold_start_metric=True)
   def lambda_handler(event: dict, context: LambdaContext) -> dict:
       order_id = event["pathParameters"]["orderId"]

       # Structured log with consistent fields
       logger.info("Processing order", extra={"order_id": order_id, "user_id": event["requestContext"]["authorizer"]["userId"]})

       try:
           result = process_order(order_id)
           metrics.add_metric(name="OrdersProcessed", unit=MetricUnit.Count, value=1)
           return {"statusCode": 200, "body": json.dumps(result)}
       except OrderNotFoundError:
           logger.warning("Order not found", extra={"order_id": order_id})
           return {"statusCode": 404, "body": '{"error": "Order not found"}'}
   ```

2. **Log security-relevant events explicitly**:

   ```python
   def log_auth_event(user_id: str, action: str, resource: str, granted: bool, context):
       logger.info("authorization_decision", extra={
           "event_type": "authorization",
           "user_id": user_id,
           "action": action,
           "resource": resource,
           "granted": granted,
           "source_ip": context.get("sourceIp"),
           "user_agent": context.get("userAgent"),
       })

   def log_data_access(user_id: str, data_type: str, record_count: int):
       logger.info("data_access", extra={
           "event_type": "data_access",
           "user_id": user_id,
           "data_type": data_type,
           "record_count": record_count,
       })
       # Alert if record_count is unusually high (possible data exfiltration)
       if record_count > 1000:
           logger.warning("large_data_access", extra={
               "event_type": "security_alert",
               "user_id": user_id,
               "record_count": record_count,
           })
   ```

3. **Configure CloudWatch alarms for anomaly detection**:

   ```yaml
   # CloudFormation — alarms on function error rate and invocation count
   ErrorRateAlarm:
     Type: AWS::CloudWatch::Alarm
     Properties:
       AlarmName: !Sub "${FunctionName}-error-rate"
       MetricName: Errors
       Namespace: AWS/Lambda
       Dimensions:
       - Name: FunctionName
         Value: !Ref ProcessOrderFunction
       Statistic: Sum
       Period: 60
       EvaluationPeriods: 5
       Threshold: 10
       ComparisonOperator: GreaterThanThreshold
       AlarmActions:
       - !Ref SecurityAlertTopic

   # Anomaly detection — alert when invocation count deviates from historical norm
   InvocationAnomalyAlarm:
     Type: AWS::CloudWatch::Alarm
     Properties:
       AlarmName: !Sub "${FunctionName}-invocation-anomaly"
       Metrics:
       - Id: invocations
         MetricStat:
           Metric:
             Namespace: AWS/Lambda
             MetricName: Invocations
             Dimensions:
             - Name: FunctionName
               Value: !Ref ProcessOrderFunction
           Stat: Sum
           Period: 300
       - Id: anomaly_band
         Expression: "ANOMALY_DETECTION_BAND(invocations, 2)"
       ComparisonOperator: GreaterThanUpperThreshold
       ThresholdMetricId: anomaly_band
       AlarmActions:
       - !Ref SecurityAlertTopic
   ```

4. **Enable X-Ray distributed tracing across function chains**:

   ```python
   from aws_lambda_powertools import Tracer

   tracer = Tracer()

   @tracer.capture_method
   def query_database(order_id: str) -> dict:
       # X-Ray automatically traces DynamoDB calls within this segment
       response = table.get_item(Key={"order_id": order_id})
       return response.get("Item")

   @tracer.capture_method
   def call_payment_api(order: dict) -> dict:
       # X-Ray traces the HTTP call to payment service
       tracer.put_annotation(key="order_id", value=order["id"])
       tracer.put_metadata(key="payment_amount", value=order["amount"])
       return payment_client.charge(order)
   ```

   ```yaml
   # SAM: enable X-Ray tracing
   Globals:
     Function:
       Tracing: Active  # all functions
   ```

5. **Monitor cost as a security signal — unexpected cost spikes indicate compromise**:

   ```bash
   # AWS Budgets — alert on Lambda cost anomaly
   aws budgets create-budget \
     --account-id 123456789 \
     --budget '{
       "BudgetName": "lambda-cost-alert",
       "BudgetType": "COST",
       "TimeUnit": "DAILY",
       "BudgetLimit": {"Amount": "50", "Unit": "USD"},
       "CostFilters": {"Service": ["AWS Lambda"]}
     }' \
     --notifications-with-subscribers '[{
       "Notification": {
         "NotificationType": "ACTUAL",
         "ComparisonOperator": "GREATER_THAN",
         "Threshold": 80
       },
       "Subscribers": [{"SubscriptionType": "EMAIL", "Address": "security@company.com"}]
     }]'
   ```

## Rules

- Always inject `correlation_id` into logs — without it, tracing a request across multiple functions in a chain requires guessing from timestamps.
- Never log event payloads at INFO level — use DEBUG level for payloads and ensure DEBUG is disabled in production (`LOG_LEVEL: INFO`).
- Security events (authentication, authorization decisions, data access) must be logged even when no error occurs — security logging is separate from error logging.
- CloudWatch Log Groups must have retention policies — without them, logs accumulate indefinitely and are never purged (cost and compliance risk).

## Common Mistakes

- **Using `print()` instead of structured logger** — print output is a single string in CloudWatch; can't be queried by field, can't be shipped to SIEM without parsing.
- **Not sampling X-Ray in high-traffic functions** — X-Ray at 100% sampling adds latency and cost at scale; use reservoir + fixed rate sampling in production.
- **Separate alarms per function instead of service-level dashboards** — Lambda functions are often orchestrated; alert on service-level error rates, not individual functions.
- **No log retention policy on CloudWatch Log Groups** — Lambda automatically creates log groups with no expiry; set retention to 90 days to balance cost and forensic investigation window.
