---
name: prevent-mass-assignment
description: Use when binding HTTP request parameters or JSON bodies directly to model objects, database records, or data transfer objects — any place where user-supplied fields are mapped to internal properties.
source: 'OWASP Mass Assignment Cheat Sheet (owasp.org/www-project-cheat-sheets); OWASP API Security Top 10 2023 API3; CWE-915'
tags: [security, owasp, mass-assignment, parameter-binding, api, developer, input-validation]
---

# Prevent Mass Assignment

Allowlist only the fields users are permitted to set when binding request data to model objects — never allow users to set internal fields like `is_admin`, `role`, `account_balance`, or `created_at`.

## Why This Is Best Practice

**Adopted by:** OWASP API Security Top 10 2023 API3 (Broken Object Property Level Authorization) is entirely caused by mass assignment. Rails, Django, and Laravel have all had critical mass assignment CVEs: GitHub (2012, CVE-2012-2661, allowed pushing to any repo via mass assignment to ActiveRecord), Rails (CVE-2013-2615). Django REST Framework, FastAPI, Spring Boot, and ASP.NET Core all provide allowlist mechanisms. OWASP ranks this in the API Top 10 specifically because it's pervasive in modern auto-binding frameworks.
**Impact:** The GitHub 2012 mass assignment vulnerability (Egor Homakov) allowed any user to gain admin-level access by submitting `user[admin]=1` in a form — exploited publicly to demonstrate the issue. Every framework that auto-binds request parameters to model attributes is vulnerable by default unless explicitly configured. One missing allowlist in a user update endpoint can allow privilege escalation.
**Why best:** Denylist approaches (listing which fields to block) require knowing all dangerous fields in advance — new fields added to the model become vulnerable automatically. Allowlist (permitting only explicitly safe fields) is safe by default: new fields are blocked until explicitly permitted.

Sources: OWASP Mass Assignment Cheat Sheet; GitHub CVE-2012-2661; OWASP API Security Top 10 2023; CWE-915

## Steps

1. **Define explicit allowlists per operation** — different operations permit different fields:

   ```python
   # FastAPI / Pydantic — separate schemas for create vs update vs response
   from pydantic import BaseModel
   from typing import Optional

   class UserCreate(BaseModel):
       username: str
       email: str
       password: str
       # NOT included: id, role, is_admin, created_at, account_balance

   class UserUpdate(BaseModel):
       email: Optional[str] = None
       display_name: Optional[str] = None
       # More restrictive than create — can't change username

   class UserResponse(BaseModel):
       id: int
       username: str
       email: str
       display_name: Optional[str]
       # NOT included: password_hash, internal_notes, payment_info

   @app.post('/users')
   def create_user(data: UserCreate):  # only fields in UserCreate can be set
       user = User(**data.dict())
       db.save(user)
       return UserResponse.from_orm(user)
   ```

2. **Django — use `fields` on ModelForm and Serializer**:

   ```python
   # ModelForm: explicit fields
   class UserUpdateForm(forms.ModelForm):
       class Meta:
           model = User
           fields = ['email', 'first_name', 'last_name']
           # is_staff, is_superuser, groups NOT listed → cannot be set

   # Django REST Framework: explicit fields
   class UserUpdateSerializer(serializers.ModelSerializer):
       class Meta:
           model = User
           fields = ['email', 'first_name', 'last_name']
           read_only_fields = ['id', 'username', 'date_joined']
   ```

3. **Rails — use strong parameters (required since Rails 4)**:

   ```ruby
   def user_params
     params.require(:user).permit(:email, :first_name, :last_name)
     # role, admin, confirmed NOT permitted
   end

   def update
     @user.update(user_params)  # only permitted params applied
   end
   ```

4. **Spring Boot — use DTOs, not entities directly**:

   ```java
   // BAD — User entity has role, isAdmin, createdAt
   @PutMapping("/users/{id}")
   public User updateUser(@PathVariable Long id, @RequestBody User user) {
       return userRepo.save(user);  // attacker can set user.setAdmin(true)
   }

   // GOOD — DTO limits what can be changed
   public class UserUpdateDTO {
       private String email;
       private String displayName;
       // role, isAdmin, createdAt NOT present
   }

   @PutMapping("/users/{id}")
   public UserResponseDTO updateUser(@PathVariable Long id,
                                     @RequestBody @Valid UserUpdateDTO dto) {
       User user = userRepo.findById(id).orElseThrow();
       user.setEmail(dto.getEmail());
       user.setDisplayName(dto.getDisplayName());
       return toResponseDTO(userRepo.save(user));
   }
   ```

5. **Never use `update_attributes(params)` or equivalent with raw request data**:

   ```ruby
   # BAD — Rails, binds everything
   @user.update_attributes(params[:user])

   # GOOD — filter first
   @user.update_attributes(user_params)
   ```

   ```python
   # BAD — Django, binds everything
   User.objects.filter(pk=user_id).update(**request.data)

   # GOOD — use serializer with explicit fields
   serializer = UserUpdateSerializer(user, data=request.data, partial=True)
   serializer.is_valid(raise_exception=True)
   serializer.save()
   ```

6. **Use read-only fields for system-managed properties** — even if a field leaks into the allowlist, mark it non-writable:

   ```python
   class UserSerializer(serializers.ModelSerializer):
       class Meta:
           model = User
           fields = ['id', 'email', 'created_at']
           read_only_fields = ['id', 'created_at']  # API can return but not set
   ```

## Rules

- The allowlist must be per-operation, not per-model — an admin endpoint may permit more fields than a user self-update endpoint.
- Nested objects need their own allowlists — `user[address][attributes][admin]=true` bypasses flat parameter filters.
- Response serializers also need allowlists — don't accidentally return `password_hash`, `internal_token`, or `ssn` in API responses.
- Prefer separate DTO/schema classes over `exclude` lists — exclusion lists grow stale as models evolve.

## Common Mistakes

- **Using `exclude` instead of `fields`** — adding a new sensitive field to the model automatically exposes it. Always use `fields` (allowlist).
- **One schema for all operations** — a create schema and an update schema typically have different permitted fields; conflating them over-permits.
- **Trusting field-level validation to catch privilege escalation** — validation checks correctness of values, not permission to set them. They're different concerns.
- **Forgetting nested / related model binding** — `user.address = Address(**request.data['address'])` can mass-assign the nested object too.
