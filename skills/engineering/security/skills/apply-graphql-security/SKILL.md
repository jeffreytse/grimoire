---
name: apply-graphql-security
description: Use when building or securing a GraphQL API — limiting query depth and complexity, disabling introspection in production, enforcing field-level authorization, and preventing batching abuse.
source: 'OWASP GraphQL Security Cheat Sheet (owasp.org/www-project-cheat-sheets); OWASP API Security Top 10 2023; HackerOne GraphQL security research; Apollo Server security documentation'
tags: [security, owasp, graphql, api, authorization, rate-limiting, developer]
---

# Apply GraphQL Security

Secure GraphQL APIs by enforcing query depth limits, disabling introspection in production, implementing field-level authorization, and rate limiting by query complexity — preventing DoS via deeply nested queries and unauthorized data access via field enumeration.

## Why This Is Best Practice

**Adopted by:** OWASP GraphQL Security Cheat Sheet (2023) is the primary reference. HackerOne's GraphQL security research documents the most common GraphQL vulnerability classes found across bug bounty programs. Apollo Server (used by Airbnb, The New York Times, Expedia) provides built-in query depth limiting and complexity analysis. GitHub's GraphQL API enforces complexity limits (5,000 points per query) and requires authentication for all requests including introspection.
**Impact:** HackerOne's 2022 "Hacking GraphQL for Fun and Profit" found introspection enabled in production in 40% of GraphQL APIs tested — allowing attackers to map entire schemas. Deeply nested GraphQL queries (10+ levels) can cause O(n^k) database queries from a single HTTP request, enabling DoS with a single request. Shopify paid out multiple GraphQL-related bug bounties for missing field-level authorization allowing access to other merchants' data — direct object reference without per-field auth checks.
**Why best:** REST APIs use URL-based authorization (middleware checks each endpoint); GraphQL resolvers require per-field authorization because a single endpoint serves all queries. A REST security model applied to GraphQL (middleware on `/graphql`) allows a user with read access to `viewer { name }` to also query `viewer { paymentMethods { cardNumber } }` if field-level checks are missing.

Sources: OWASP GraphQL Security Cheat Sheet; HackerOne "Hacking GraphQL for Fun and Profit" (2022); Apollo Server documentation; GitHub GraphQL API documentation

## Steps

1. **Limit query depth and complexity — prevent DoS via nested queries**:

   ```python
   # Graphene (Python) with graphene-django
   from graphql import parse, validate
   from graphql.validation.rules import NoSchemaIntrospectionCustomRule

   # Install: pip install graphql-core graphene
   # Depth limit middleware
   class QueryDepthLimiter:
       MAX_DEPTH = 7

       def resolve(self, next, root, info, **args):
           if root is None:  # top-level only
               depth = self._get_query_depth(info.operation)
               if depth > self.MAX_DEPTH:
                   raise Exception(f"Query depth {depth} exceeds maximum {self.MAX_DEPTH}")
           return next(root, info, **args)

       def _get_query_depth(self, node, depth=0):
           max_depth = depth
           if hasattr(node, 'selection_set') and node.selection_set:
               for selection in node.selection_set.selections:
                   child_depth = self._get_query_depth(selection, depth + 1)
                   max_depth = max(max_depth, child_depth)
           return max_depth
   ```

   ```javascript
   // Apollo Server (Node.js) with graphql-depth-limit
   import depthLimit from 'graphql-depth-limit';
   import { createComplexityLimitRule } from 'graphql-validation-complexity';

   const server = new ApolloServer({
     schema,
     validationRules: [
       depthLimit(7),
       createComplexityLimitRule(1000, {
         onCost: (cost) => console.log('Query complexity:', cost),
       }),
     ],
   });
   ```

2. **Disable introspection in production**:

   ```javascript
   // Apollo Server
   const server = new ApolloServer({
     schema,
     introspection: process.env.NODE_ENV !== 'production',
   });
   ```

   ```python
   # Graphene — disable introspection via validation rule
   from graphql import build_ast_schema, validate
   from graphql.validation import NoSchemaIntrospectionCustomRule

   def execute_query(query_string, schema, context):
       document = parse(query_string)
       errors = validate(schema, document, rules=[NoSchemaIntrospectionCustomRule])
       if errors:
           return {"errors": [str(e) for e in errors]}
       # execute...
   ```

3. **Implement field-level authorization in resolvers**:

   ```python
   # BAD — authorization only at query level
   @login_required
   def resolve_user(root, info, id):
       return User.objects.get(id=id)  # returns all fields including email, SSN

   # GOOD — field-level checks
   class UserType(graphene.ObjectType):
       email = graphene.String()
       ssn_last4 = graphene.String()

       def resolve_email(root, info):
           # Only return email if viewer is the user or an admin
           viewer = info.context.user
           if viewer.id != root.id and not viewer.is_admin:
               return None
           return root.email

       def resolve_ssn_last4(root, info):
           viewer = info.context.user
           if viewer.id != root.id:
               raise PermissionError("Cannot access SSN of another user")
           return root.ssn_last4
   ```

4. **Rate limit by query complexity, not just request count**:

   ```python
   import time
   from collections import defaultdict

   class ComplexityRateLimiter:
       WINDOW_SECONDS = 60
       MAX_COMPLEXITY_PER_WINDOW = 5000

       def __init__(self):
           self._usage: dict = defaultdict(list)

       def check_and_record(self, user_id: str, complexity: int) -> bool:
           now = time.time()
           window_start = now - self.WINDOW_SECONDS

           # Remove expired entries
           self._usage[user_id] = [
               (ts, c) for ts, c in self._usage[user_id] if ts > window_start
           ]

           current_usage = sum(c for _, c in self._usage[user_id])
           if current_usage + complexity > self.MAX_COMPLEXITY_PER_WINDOW:
               return False

           self._usage[user_id].append((now, complexity))
           return True
   ```

5. **Prevent batching abuse — limit query batching**:

   ```javascript
   // Apollo Server — disable query batching or set a batch size limit
   const server = new ApolloServer({
     schema,
     allowBatchedHttpRequests: false,  // disable batching entirely
   });

   // Or limit batch size via middleware
   app.use('/graphql', (req, res, next) => {
     if (Array.isArray(req.body) && req.body.length > 10) {
       return res.status(400).json({ error: 'Batch size exceeds limit of 10' });
     }
     next();
   });
   ```

6. **Validate and sanitize all input arguments**:

   ```python
   import graphene
   from graphene import String, Int
   import re

   class SearchUsers(graphene.ObjectType):
       search_users = graphene.List(
           UserType,
           query=String(required=True),
           limit=Int(default_value=20)
       )

       def resolve_search_users(root, info, query, limit):
           # Validate inputs
           if len(query) > 100:
               raise ValueError("Search query too long")
           if limit > 100:
               raise ValueError("Limit exceeds maximum of 100")
           # Use ORM — never string interpolation into raw SQL
           return User.objects.filter(name__icontains=query)[:limit]
   ```

## Rules

- Disable introspection in production — if schema discovery is needed (for internal tools), require authentication.
- Every resolver that returns sensitive data must check authorization — not just the top-level query.
- Query aliases allow executing the same field multiple times; depth and complexity limits must account for aliases.
- Persisted queries (allowlist of approved query hashes) eliminate arbitrary query execution — consider for production APIs.

## Common Mistakes

- **`__typename` bypassing introspection blocks** — `__typename` is a meta-field that leaks type names even when introspection is disabled; consider whether it needs to be blocked too.
- **Using DataLoader only for performance** — DataLoader (batched loading) also mitigates N+1 query amplification attacks; skipping it allows a single query to cause thousands of DB round trips.
- **Missing authorization on mutation arguments** — checking `canUpdateUser(viewer, userId)` but not `canSetRole(viewer, newRole)` allows privilege escalation through mutation arguments.
- **Error messages exposing schema details** — default GraphQL error messages often include resolver and field names; customize error formatting to remove internal details.
