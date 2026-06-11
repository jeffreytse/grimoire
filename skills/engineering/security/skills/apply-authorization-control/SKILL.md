---
name: apply-authorization-control
description: Use when implementing access control decisions — determining which users can perform which actions on which resources — in any API, web application, or service.
source: 'OWASP Authorization Cheat Sheet (owasp.org/www-project-cheat-sheets); OWASP API Security Top 10 2023 API1/API3/API5; NIST SP 800-162 (ABAC); CWE-285'
tags: [security, owasp, authorization, rbac, abac, access-control, developer]
---

# Apply Authorization Control

Enforce explicit, deny-by-default authorization at every resource access point using Role-Based (RBAC) or Attribute-Based (ABAC) access control — never relying on UI hiding or obscurity as security.

## Why This Is Best Practice

**Adopted by:** OWASP API Security Top 10 2023 lists Broken Object Level Authorization (API1), Broken Object Property Level Authorization (API3), and Broken Function Level Authorization (API5) as the top API vulnerability classes — all caused by missing or incomplete authorization checks. NIST SP 800-162 defines ABAC as the federal standard for fine-grained access control. AWS IAM, Google Cloud IAM, and Azure RBAC all implement these patterns at cloud scale. PCI DSS v4.0 Requirement 7 mandates least-privilege access control.
**Impact:** OWASP reports that broken access control was the #1 web vulnerability in 2021 (A01). Broken Object Level Authorization is the top API vulnerability in 2023 — it allows any authenticated user to access any other user's data by changing an ID in the request. Dropbox's 2012 breach and Instagram's 2019 BOLA vulnerability both stemmed from missing object-level authorization. Proper authorization prevents entire classes of data exposure.
**Why best:** UI-level hiding (not rendering a button) and security through obscurity (unlinked URLs) are the common alternatives — they provide zero protection against direct API calls. Explicit server-side authorization checks enforced on every request are the only reliable defense.

Sources: OWASP Authorization Cheat Sheet; OWASP API Security Top 10 2023; NIST SP 800-162; CWE-285

## Steps

1. **Default deny: every resource access must pass an explicit authorization check**:

   ```python
   def get_document(doc_id, current_user):
       doc = db.get(doc_id)
       if doc is None:
           raise NotFound()
       # Explicit ownership/permission check — never skip this
       if doc.owner_id != current_user.id and not current_user.has_role('admin'):
           raise Forbidden()
       return doc
   ```

   Never assume that authenticated = authorized. Every object-level fetch needs a check.

2. **Implement Role-Based Access Control (RBAC)** for function-level permissions:

   ```python
   # Define roles and their permitted actions
   ROLE_PERMISSIONS = {
       'viewer':    {'document:read'},
       'editor':    {'document:read', 'document:write'},
       'admin':     {'document:read', 'document:write', 'document:delete', 'user:manage'},
   }

   def require_permission(permission):
       def decorator(func):
           def wrapper(*args, **kwargs):
               if permission not in ROLE_PERMISSIONS.get(current_user.role, set()):
                   raise Forbidden(f"Requires {permission}")
               return func(*args, **kwargs)
           return wrapper
       return decorator

   @require_permission('document:delete')
   def delete_document(doc_id):
       ...
   ```

3. **For multi-tenancy: always scope queries by owner/tenant** — never return all records and filter in application code:

   ```python
   # BAD — fetches all, filters in Python (BOLA if filter fails)
   all_docs = db.query("SELECT * FROM documents")
   user_docs = [d for d in all_docs if d.owner_id == current_user.id]

   # GOOD — database enforces ownership
   user_docs = db.query(
       "SELECT * FROM documents WHERE owner_id = %s",
       (current_user.id,)
   )
   ```

4. **Validate field-level access for sensitive attributes** (Broken Object Property Level Authorization):

   ```python
   # Define which fields each role can read/write
   FIELD_READ_PERMISSIONS = {
       'viewer': {'id', 'title', 'content', 'created_at'},
       'admin':  {'id', 'title', 'content', 'created_at', 'owner_id', 'internal_notes'},
   }

   def serialize_document(doc, user_role):
       allowed = FIELD_READ_PERMISSIONS.get(user_role, set())
       return {k: v for k, v in doc.items() if k in allowed}
   ```

5. **Apply ABAC for fine-grained, context-aware decisions**:

   ```python
   def can_edit(user, document, context):
       # Attribute-based: combine user attributes, resource attributes, context
       if user.department != document.department:
           return False
       if context.time_of_day not in document.allowed_hours:
           return False
       if document.classification_level > user.clearance_level:
           return False
       return True
   ```

6. **Use a centralized authorization library** — don't scatter `if user.role == 'admin'` checks:

   - Python: `casbin` (RBAC/ABAC policy engine)
   - Node.js: `casl`, `node-casbin`
   - Java: Spring Security's method security (`@PreAuthorize`)
   - .NET: ASP.NET Core policy-based authorization

   ```python
   from casbin import Enforcer
   e = Enforcer('model.conf', 'policy.csv')

   if not e.enforce(user.id, resource, action):
       raise Forbidden()
   ```

7. **Log all authorization failures** — repeated failures indicate probing:

   ```python
   if not authorized:
       logger.warning("authz_denied",
           user_id=current_user.id,
           resource=resource_id,
           action=action,
           ip=request.remote_addr)
       raise Forbidden()
   ```

## Rules

- Never check authorization only in the UI — always enforce server-side on every API call.
- Object-level authorization (BOLA) check must happen after fetching the object, not before — you need the object's owner/attributes to check.
- Admin routes must still check authorization — "internal" or "admin" endpoints are primary targets.
- Indirect object references (UUIDs instead of sequential IDs) reduce guessability but are NOT authorization — still enforce ownership checks.

## Common Mistakes

- **Checking role but not ownership** — an admin check passes for all admins, but a user should only access their own data.
- **Filtering in application code after fetching all records** — if the filter code has a bug, all records are exposed.
- **Relying on HTTP method to infer permission** — attackers change GET to POST or DELETE. Check permission explicitly per action.
- **Missing authorization on bulk/batch endpoints** — `PATCH /api/documents` with an array body needs per-document ownership checks, not just route-level RBAC.
