---
name: apply-microservices-security
description: Use when building microservices that communicate with each other — implementing service-to-service mTLS, JWT propagation for user identity, and API gateway authentication to prevent lateral movement between services.
source: 'OWASP Microservices Security Cheat Sheet (owasp.org/www-project-cheat-sheets); NIST SP 800-204 (Microservices Security); Netflix Zuul/Istio documentation; Google BeyondProd whitepaper'
tags: [security, owasp, microservices, mtls, jwt, service-mesh, api-gateway, developer]
---

# Apply Microservices Security

Implement mutual TLS between services, propagate user identity via JWT without leaking service-to-service credentials to clients, and enforce per-service authorization — preventing lateral movement within the service mesh when any single service is compromised.

## Why This Is Best Practice

**Adopted by:** OWASP Microservices Security Cheat Sheet (2023) is the primary reference. NIST SP 800-204 (Security Strategies for Microservices) is the federal guidance. Google's BeyondProd (2019) describes their production microservices security model, which underpins Google Cloud's security design. Netflix, Uber, and Lyft all use Istio or Envoy-based service meshes with mTLS for service-to-service authentication. All major service mesh implementations (Istio, Linkerd, Consul Connect) provide mTLS as their primary security feature.
**Impact:** In a microservices architecture, services typically communicate over the internal network without authentication — any compromised container on the internal network can make arbitrary API calls to any service. The 2019 Capital One breach involved lateral movement from a compromised EC2 instance to S3 — the equivalent attack in microservices would be a compromised service accessing other services' data. Netflix's chaos engineering found that 80% of theoretical attack paths in their microservices architecture relied on unauthenticated internal service calls. mTLS eliminates this attack class by requiring cryptographic proof of service identity on every call.
**Why best:** Network segmentation (VPCs, security groups) provides perimeter security but doesn't authenticate services within the perimeter. mTLS provides cryptographic service identity — every service proves it is who it claims to be on every call. This enables zero-trust within the service mesh: even if an attacker compromises the network, they cannot make authenticated API calls without a valid service certificate.

Sources: OWASP Microservices Security Cheat Sheet; NIST SP 800-204; Google BeyondProd whitepaper (2019); Netflix Istio adoption guide

## Steps

1. **Enable mTLS between services with Istio**:

   ```yaml
   # Istio PeerAuthentication — require mTLS in production namespace
   apiVersion: security.istio.io/v1beta1
   kind: PeerAuthentication
   metadata:
     name: default
     namespace: production
   spec:
     mtls:
       mode: STRICT  # reject any plaintext connections
   ---
   # DestinationRule — enforce mTLS when calling payment service
   apiVersion: networking.istio.io/v1beta1
   kind: DestinationRule
   metadata:
     name: payment-service-mtls
     namespace: production
   spec:
     host: payment-service.production.svc.cluster.local
     trafficPolicy:
       tls:
         mode: ISTIO_MUTUAL  # use Istio-managed certs
   ```

2. **Propagate user identity via JWT without service credentials**:

   ```python
   # API Gateway validates external JWT, issues internal service token
   import jwt
   from datetime import datetime, timedelta

   INTERNAL_SECRET = get_secret("/production/internal-jwt-secret")
   EXTERNAL_JWKS_URL = "https://auth.company.com/.well-known/jwks.json"

   def gateway_auth_middleware(request):
       # Step 1: Validate external user JWT
       external_token = request.headers.get("Authorization", "").removeprefix("Bearer ")
       try:
           user_claims = verify_external_jwt(external_token, EXTERNAL_JWKS_URL)
       except Exception:
           return Response(401, "Unauthorized")

       # Step 2: Issue internal service token with user context
       # Internal token has shorter TTL and is not exposed to clients
       internal_token = jwt.encode({
           "sub": user_claims["sub"],      # user ID
           "email": user_claims["email"],
           "roles": user_claims["roles"],
           "iat": datetime.utcnow(),
           "exp": datetime.utcnow() + timedelta(minutes=5),  # short TTL
           "iss": "api-gateway",
           "aud": "internal",
       }, INTERNAL_SECRET, algorithm="HS256")

       # Forward to downstream services
       request.headers["X-Internal-Token"] = internal_token
       request.headers["X-User-ID"] = user_claims["sub"]
       # Remove original token — downstream services don't need external JWT
       del request.headers["Authorization"]
       return forward_to_service(request)
   ```

3. **Enforce per-service authorization with Istio AuthorizationPolicy**:

   ```yaml
   # Only order-service can call payment-service — no other service should
   apiVersion: security.istio.io/v1beta1
   kind: AuthorizationPolicy
   metadata:
     name: payment-service-policy
     namespace: production
   spec:
     selector:
       matchLabels:
         app: payment-service
     rules:
     - from:
       - source:
           principals:
           # Only order-service's service account
           - "cluster.local/ns/production/sa/order-service"
       to:
       - operation:
           methods: ["POST"]
           paths: ["/v1/payments/charge"]
     # Deny everything else — default deny
   ```

4. **Validate and extract user context in each service**:

   ```python
   # Each downstream service validates the internal token
   import jwt
   from functools import wraps

   def require_internal_auth(f):
       @wraps(f)
       def decorated(request, *args, **kwargs):
           token = request.headers.get("X-Internal-Token")
           if not token:
               return Response(401, "Missing internal auth token")

           try:
               claims = jwt.decode(
                   token,
                   INTERNAL_SECRET,
                   algorithms=["HS256"],
                   audience="internal",
                   issuer="api-gateway",
               )
           except jwt.ExpiredSignatureError:
               return Response(401, "Token expired")
           except jwt.InvalidTokenError as e:
               logger.warning(f"Invalid internal token: {e}")
               return Response(401, "Invalid token")

           # Attach user context to request for authorization checks
           request.user_id = claims["sub"]
           request.user_roles = claims["roles"]
           return f(request, *args, **kwargs)
       return decorated

   @require_internal_auth
   def get_order(request, order_id: str):
       order = Order.get(order_id)
       # Authorization: user can only access their own orders
       if order.user_id != request.user_id and "admin" not in request.user_roles:
           return Response(403, "Forbidden")
       return Response(200, order.to_dict())
   ```

5. **Secure service-to-service calls without user context (background jobs)**:

   ```python
   # Service-to-service token for internal calls without a user (cron jobs, event processing)
   def get_service_token(calling_service: str, target_service: str) -> str:
       return jwt.encode({
           "sub": calling_service,
           "aud": target_service,
           "iat": datetime.utcnow(),
           "exp": datetime.utcnow() + timedelta(minutes=5),
           "iss": "service-auth",
           "type": "service",  # distinguish from user tokens
       }, SERVICE_SECRET, algorithm="HS256")

   # Or: use workload identity / SPIFFE (via Istio) for automatic service tokens
   # No code needed — Istio provides X.509 SVIDs automatically to each workload
   ```

## Rules

- mTLS in PERMISSIVE mode (Istio default) allows plaintext — switch to STRICT before going to production.
- Internal tokens must have short TTLs (≤5 minutes) and must not be cached beyond that window — they're disposable credentials for a single request chain.
- Services must never accept the original client Authorization header as proof of internal identity — attackers who steal a user JWT can impersonate that user to all downstream services if there's no internal token.
- AuthorizationPolicy default action is ALLOW in Istio without a policy — explicitly deploy deny-by-default policies before adding allow rules.

## Common Mistakes

- **Sharing a single JWT signing key across all services** — compromise of any service exposes the key; use asymmetric keys (RS256) so services can verify tokens without being able to mint them.
- **Forwarding the client's JWT to downstream services directly** — the client JWT has a long TTL and broad scope; create short-lived internal tokens scoped to the specific downstream call.
- **No service-to-service authorization (only authentication)** — verifying that a call comes from payment-service doesn't mean payment-service should be able to access user-service's admin endpoints; add operation-level authorization.
- **Not rotating service certificates** — Istio automatically rotates certificates every 24 hours; if using manual mTLS, implement automatic rotation with ≤90 day certs.
