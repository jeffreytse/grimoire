---
name: apply-api-rate-limiting
description: Use when building any API endpoint — especially authentication endpoints, resource-creation endpoints, and any operation with a per-user or per-tenant cost — to protect against abuse, brute force, and resource exhaustion.
source: 'OWASP API Security Top 10 2023 API4 (Unrestricted Resource Consumption); OWASP API Security Cheat Sheet (owasp.org/www-project-cheat-sheets); CWE-770'
tags: [security, owasp, rate-limiting, throttling, api, dos-prevention, developer]
---

# Apply API Rate Limiting

Enforce per-client request quotas using token bucket or sliding window algorithms — protecting against brute force, scraping, and resource exhaustion attacks.

## Why This Is Best Practice

**Adopted by:** OWASP API Security Top 10 2023 API4 (Unrestricted Resource Consumption) is a top API vulnerability. AWS API Gateway, Google Cloud Endpoints, Azure APIM, Kong, Nginx, and Cloudflare all provide built-in rate limiting. Stripe, GitHub, and Twitter/X enforce rate limits on every API endpoint with well-documented headers. PCI DSS v4.0 Requirement 8.3.4 mandates account lockout after failed authentication attempts — which requires rate limiting.
**Impact:** Rate limiting prevents brute force credential attacks (see also `design-session-management`), credential stuffing (using credential lists), API scraping (competitive intelligence theft), and denial of service via resource exhaustion. Without authentication endpoint rate limits, attackers test millions of passwords in hours. Without resource endpoint limits, a single user can consume 100% of server capacity.
**Why best:** Client-side debouncing and application-level retries are the alternatives — they don't prevent malicious clients from ignoring them. Server-side rate limiting enforced at the network edge (or application layer with a distributed counter) is the only enforceable mechanism.

Sources: OWASP API Security Top 10 2023 API4; OWASP API Security Cheat Sheet; CWE-770; Stripe rate limiting design

## Steps

1. **Implement token bucket rate limiting** — smooth limit that allows short bursts:

   ```python
   import time
   import redis

   class TokenBucket:
       def __init__(self, redis_client, key, rate, capacity):
           self.redis = redis_client
           self.key = key
           self.rate = rate          # tokens added per second
           self.capacity = capacity  # max tokens (burst size)

       def consume(self, tokens=1):
           now = time.time()
           pipe = self.redis.pipeline()
           pipe.hgetall(self.key)
           results = pipe.execute()
           data = results[0]

           last_time = float(data.get(b'last_time', now))
           stored = float(data.get(b'tokens', self.capacity))

           # Add tokens for elapsed time
           elapsed = now - last_time
           stored = min(self.capacity, stored + elapsed * self.rate)

           if stored < tokens:
               return False  # rate limited

           stored -= tokens
           pipe = self.redis.pipeline()
           pipe.hset(self.key, mapping={'tokens': stored, 'last_time': now})
           pipe.expire(self.key, int(self.capacity / self.rate) + 1)
           pipe.execute()
           return True
   ```

2. **Apply different limits by endpoint sensitivity**:

   ```python
   RATE_LIMITS = {
       'auth_login':        {'rate': 5,    'capacity': 5,   'window': 60},   # 5/min
       'auth_password_reset': {'rate': 3,  'capacity': 3,   'window': 3600}, # 3/hr
       'api_read':          {'rate': 100,  'capacity': 200, 'window': 60},   # 100/min
       'api_write':         {'rate': 30,   'capacity': 60,  'window': 60},   # 30/min
       'api_export':        {'rate': 5,    'capacity': 10,  'window': 3600}, # 5/hr
   }

   def rate_limit_key(endpoint, identifier):
       # identifier = IP for unauthenticated, user_id for authenticated
       return f'ratelimit:{endpoint}:{identifier}'
   ```

3. **Return standard rate limit headers** (RFC 6585 + IETF RateLimit header draft):

   ```python
   from flask import g, request, jsonify

   @app.before_request
   def check_rate_limit():
       endpoint = request.endpoint
       identifier = current_user.id if current_user.is_authenticated else request.remote_addr
       key = rate_limit_key(endpoint, identifier)
       bucket = TokenBucket(redis_client, key, **RATE_LIMITS.get(endpoint, DEFAULT_LIMIT))

       if not bucket.consume():
           response = jsonify({'error': 'Rate limit exceeded. Retry after specified time.'})
           response.headers['Retry-After'] = str(bucket.retry_after_seconds())
           response.headers['X-RateLimit-Limit'] = str(bucket.capacity)
           response.headers['X-RateLimit-Remaining'] = '0'
           response.headers['X-RateLimit-Reset'] = str(int(time.time()) + bucket.retry_after_seconds())
           return response, 429
   ```

4. **Rate limit by multiple dimensions** — IP alone is insufficient (shared NAT, VPNs):

   ```python
   def get_rate_limit_identifiers(request, user):
       identifiers = [f'ip:{request.remote_addr}']
       if user and user.is_authenticated:
           identifiers.append(f'user:{user.id}')
           identifiers.append(f'tenant:{user.tenant_id}')
       return identifiers

   # Check ALL identifiers — fail if any is exceeded
   def is_rate_limited(request, user, endpoint):
       for identifier in get_rate_limit_identifiers(request, user):
           key = rate_limit_key(endpoint, identifier)
           if not TokenBucket(redis_client, key, **limits).consume():
               return True
       return False
   ```

5. **Use middleware/infrastructure-level rate limiting for edge protection**:

   ```nginx
   # Nginx — limit login endpoint to 5 requests/minute per IP
   limit_req_zone $binary_remote_addr zone=login:10m rate=5r/m;

   location /api/auth/login {
       limit_req zone=login burst=3 nodelay;
       limit_req_status 429;
       proxy_pass http://backend;
   }
   ```

   Combine with application-level limits for defense-in-depth.

6. **Implement exponential backoff for repeat offenders**:

   ```python
   def get_lockout_duration(violation_count):
       # Exponential: 1s, 2s, 4s, 8s ... up to 1 hour
       return min(2 ** (violation_count - 1), 3600)

   def record_violation(identifier):
       count = redis.incr(f'violations:{identifier}')
       redis.expire(f'violations:{identifier}', 86400)  # reset after 24h
       return count
   ```

## Rules

- Always rate limit authentication endpoints regardless of other protections — credential stuffing is fully automated.
- Rate limit by user ID (not just IP) for authenticated endpoints — a compromised account can still abuse from many IPs.
- Return 429 (Too Many Requests) with `Retry-After` — not 403 (which suggests permanent denial).
- Distributed systems need distributed rate limit counters (Redis, Memcached) — in-memory counters per instance are bypassable by distributing requests across instances.

## Common Mistakes

- **Only rate limiting unauthenticated endpoints** — authenticated users can also abuse resource-intensive endpoints.
- **Using fixed windows** (reset every minute) — vulnerable to bursting at the window boundary; use sliding window or token bucket.
- **Ignoring `X-Forwarded-For` header behind a load balancer** — all clients appear to come from the load balancer's IP; use the rightmost untrusted IP in the header chain.
- **Rate limiting login attempts but not password reset** — password reset is a brute force vector for account enumeration.
