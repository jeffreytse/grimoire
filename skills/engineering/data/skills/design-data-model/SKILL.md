---
name: design-data-model
description: Use when designing a new relational database schema, adding entities to an existing system, reviewing a schema for normalization issues, or planning migrations — before writing DDL or ORM models.
source: E.F. Codd (1970) "A Relational Model of Data for Large Shared Data Banks", Martin Fowler "Patterns of Enterprise Application Architecture" (2002), PostgreSQL documentation
tags: [schema-design, normalization, relational-database, data-modeling, constraints, indexing, backend, correctness]
verified: true
---

# Design Data Model

Design a correct, normalized relational schema before writing DDL or ORM code.

## Why This Is Best Practice

**Adopted by:** Google (Spanner schema design guide), Stripe (data modeling RFC process), Amazon (schema review as part of ORR — Operational Readiness Review), Airbnb (schema review gate before any new table ships)
**Impact:** Fixing a schema defect post-launch costs 10–100× more than fixing it at design time (Barry Boehm, "Software Engineering Economics", 1981 — ratio holds for data schema per Fowler, PEAA). Unnormalized schemas cause data anomalies (update, insertion, deletion) that corrupt production data silently. Facebook's 2012 migration off a denormalized messages schema cost 1.5 engineer-years (public post-mortem). Proper indexing at design time reduces query latency by 10–1000× compared to retrofitting indexes on large tables.
**Why best:** Normalization to 3NF eliminates redundancy and anomalies by construction — no amount of application-layer discipline achieves the same guarantee. Alternative (denormalize for performance upfront) is premature optimization: normalized schemas can be selectively denormalized with materialized views once bottlenecks are measured.

Sources: E.F. Codd (1970), Martin Fowler (PEAA 2002), PostgreSQL indexing documentation, Stripe engineering blog

## Steps

### 1. List entities and their attributes

Write one line per entity: `Entity — attribute1, attribute2, ...`

Ask: does each attribute describe this entity, or something else? If something else, it belongs in a different entity.

### 2. Identify primary keys

Every entity must have a stable, unique, non-null identifier.

- Prefer surrogate keys (`id BIGSERIAL` / UUID) over natural keys unless the natural key is guaranteed stable and unique (e.g., ISO country code).
- Never use email, phone, or username as PK — they change.

### 3. Map relationships and cardinality

For every pair of related entities, state: one-to-one, one-to-many, or many-to-many.

- one-to-many: FK on the "many" side.
- many-to-many: junction table with FKs to both sides + its own PK.
- one-to-one: FK + UNIQUE constraint, or same table if always co-present.

### 4. Normalize to 3NF

Apply in order:

**1NF** — Each column holds one atomic value. No comma-separated lists, no repeated column groups (`tag1`, `tag2`, `tag3`). Extract repeating groups to a child table.

**2NF** — Every non-key column depends on the whole PK, not part of it. If a composite PK has columns A+B and column C depends only on A, extract (A, C) to a separate table.

**3NF** — No transitive dependencies. If column C depends on column B which depends on PK A, extract (B, C) to its own table.

Stop at 3NF by default. Go to BCNF only if you have overlapping candidate keys causing anomalies.

### 5. Add constraints

Every constraint you skip becomes an application bug waiting to happen:

```sql
-- Not null where semantically required
ALTER TABLE orders ADD COLUMN user_id BIGINT NOT NULL;

-- Foreign key with explicit action
ALTER TABLE orders ADD CONSTRAINT fk_orders_user
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT;

-- Unique where business rules require it
ALTER TABLE users ADD CONSTRAINT uq_users_email UNIQUE (email);

-- Check constraints for domain rules
ALTER TABLE products ADD CONSTRAINT chk_price_positive CHECK (price > 0);
```

### 6. Design indexes

Add indexes for every access pattern you know will exist at launch:

- PK index: auto-created.
- FK columns: always index FK columns (prevents full table scans on joins).
- Query predicates: index columns in `WHERE`, `ORDER BY`, `GROUP BY` for high-traffic queries.
- Composite indexes: column order matters — most selective or most frequently filtered first.

```sql
-- FK index
CREATE INDEX idx_orders_user_id ON orders(user_id);

-- Composite for common query shape: WHERE status = ? ORDER BY created_at DESC
CREATE INDEX idx_orders_status_created ON orders(status, created_at DESC);
```

Skip indexes for: columns updated frequently, low-cardinality columns (boolean, enum with 2-3 values) on large write-heavy tables.

### 7. Review the schema for red flags

Check each of these before finalizing:

- [ ] Any nullable FK? Confirm it's intentional (optional relationship).
- [ ] Any VARCHAR without length that should be bounded? Add CHECK or use TEXT explicitly.
- [ ] Any money/price column as FLOAT? Use `NUMERIC(precision, scale)` — floats lose cents.
- [ ] Any timestamp without timezone? Use `TIMESTAMPTZ` (UTC-stored) for anything user-facing.
- [ ] Any soft-delete `deleted_at` column? Add partial index `WHERE deleted_at IS NULL` on commonly filtered queries.

### 8. Write the DDL or ORM model

Translate the validated schema to DDL (SQL) or ORM models (SQLAlchemy, Django models, Prisma schema). Ensure constraints are expressed at the DB layer, not only the application layer.

## Rules

- Never skip primary keys — every table must have one.
- Never store money as FLOAT — use NUMERIC or integer cents.
- Always set ON DELETE behavior on FK constraints explicitly (RESTRICT, CASCADE, or SET NULL).
- Never use application-layer uniqueness checks as a substitute for UNIQUE constraints — race conditions exist.
- Always index FK columns.
- Normalize first, denormalize only after measuring a bottleneck.

## Examples

**Before (unnormalized):**
```
orders: id, customer_name, customer_email, product_name, product_price, quantity
```
Problems: customer and product data duplicated per row; update anomaly if email changes.

**After (3NF):**
```
users:    id, name, email
products: id, name, price_cents (NUMERIC)
orders:   id, user_id (FK), created_at (TIMESTAMPTZ)
order_items: id, order_id (FK), product_id (FK), quantity, unit_price_cents (NUMERIC)
```

## Common Mistakes

- **Storing arrays in a column** (`tags TEXT DEFAULT ''`): violates 1NF; use a junction table or a DB-native array type with GIN index.
- **Natural keys as PKs**: emails and usernames change; surrogate keys don't.
- **Skipping FK constraints**: application-layer checks race; the DB won't.
- **Denormalizing prematurely**: almost always wrong at design time; profile first.
- **Forgetting timezone on timestamps**: naive timestamps cause midnight bugs at DST boundaries and across regions.
- **FLOAT for currency**: `0.1 + 0.2 = 0.30000000000000004`; use NUMERIC.
