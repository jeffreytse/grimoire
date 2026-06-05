---
name: write-database-migration
description: Use when changing a database schema — adding or removing columns, creating tables, adding indexes, or altering constraints.
source: Widely adopted at Rails, Django, Laravel, Flyway, Liquibase — codified in Fowler & Sadalage, "Evolutionary Database Design", 2016 (martinfowler.com/articles/evodb.html)
tags: [database, schema-evolution, migrations, developer, deployment-safety, rollback]
verified: true
---

# Write Database Migration

Version every schema change as a migration file committed alongside the code that requires it.

## Why This Is Best Practice

**Adopted by:** Every major web framework ships a migration tool: Rails ActiveRecord
Migrations, Django Migrations, Laravel Migrations. Flyway is used at Google, Zalando,
Wix; Liquibase is used at 10,000+ enterprises. Schema-as-code is the universal default
in 2024 — ad-hoc ALTER scripts are considered a smell in professional engineering teams.
**Impact:** Teams using versioned migrations report near-zero schema drift incidents
between environments (dev/staging/prod). Automated rollback capability reduces mean
time to recovery from hours (restore from backup) to minutes (down migration). Every
change is code-reviewed, auditable via git log, and exactly reproducible across all
environments.
**Why best:** Ad-hoc ALTER scripts run manually have no audit trail, can't be replayed,
diverge between environments, and can't be rolled back predictably. A migration file
is committed in the same PR as the application code that requires it — impossible to
deploy code that references a column that doesn't exist.

Sources: Fowler & Sadalage, "Evolutionary Database Design" (martinfowler.com/articles/evodb.html);
Rails Active Record Migrations guide; Flyway documentation (flyway.org)

## Steps

### 1. One migration = one schema change

Never bundle unrelated alterations. Each migration should be independently deployable
and independently rollbackable:

```
# Bad — one migration does too much
20240601_setup_users_and_orders.sql

# Good — separate concerns
20240601120000_create_users_table.sql
20240601120001_create_orders_table.sql
20240602090000_add_email_index_to_users.sql
```

### 2. Name with a timestamp prefix

Format: `YYYYMMDDHHMMSS_<what_changed>.sql` (or the equivalent in your framework).

The timestamp ensures deterministic ordering across branches and parallel development:

```
20240601120000_add_email_to_users.sql
20240601120001_create_sessions_table.sql
20240602090000_add_index_on_users_email.sql
```

### 3. Write both up and down

Every migration must have a reversal path:

```sql
-- up
ALTER TABLE users ADD COLUMN email VARCHAR(255);

-- down
ALTER TABLE users DROP COLUMN email;
```

If the rollback is destructive (data loss), document it explicitly rather than omitting it:

```sql
-- down (WARNING: drops all email data — backup required before rollback)
ALTER TABLE users DROP COLUMN email;
```

### 4. Test on production-like data before merge

Run the migration against a recent anonymized copy of the production database.
Check:
- Migration completes without errors
- Application code works correctly after migration
- Down migration reverses cleanly
- Query performance is acceptable with production-scale data

### 5. Handle large tables safely (>1M rows)

Standard `ALTER TABLE` locks the table. On large tables, use non-locking patterns:

**Adding a NOT NULL column with a default:**
```sql
-- Step 1: Add nullable (no lock)
ALTER TABLE orders ADD COLUMN status VARCHAR(50);

-- Step 2: Backfill in batches (separate migration or script)
UPDATE orders SET status = 'pending' WHERE status IS NULL AND id BETWEEN 1 AND 10000;
-- ... repeat in batches

-- Step 3: Add NOT NULL constraint after backfill (fast — no data scan needed if all filled)
ALTER TABLE orders ALTER COLUMN status SET NOT NULL;
```

**Adding an index:** Use `CREATE INDEX CONCURRENTLY` (PostgreSQL) or equivalent.
Never create indexes without `CONCURRENTLY` on live tables with significant traffic.

### 6. Never edit a merged migration

A merged migration has run in staging or production. Editing it creates a divergence
between the recorded migration history and the actual database state.

For corrections: add a new migration that fixes the error.

```
# Bad — editing merged migration
git commit -m "fix typo in 20240601_add_email_to_users"  # already ran in prod

# Good — new corrective migration
20240603_rename_email_typo_to_correct_column.sql
```

## Rules

- Migrations live in the same PR as the code that requires the schema change — never merge application code before its migration
- Never apply raw SQL directly to production databases — always run through the migration tool to maintain history
- Down migrations are required — "we'll never roll back" is an assumption that fails during incidents
- Seed data (test fixtures, reference data) belongs in seed files, not migrations — migrations model structure, not content
- All migrations must be idempotent when possible: `CREATE TABLE IF NOT EXISTS`, `ADD COLUMN IF NOT EXISTS` (PostgreSQL 9.6+)

## Common Mistakes

**Skipping the down migration.** The first time you need to roll back at 2am, an absent
down migration forces a manual database restore.

**Bundling schema changes with data migrations.** Schema changes (DDL) and data
transformations (DML) have different risk profiles and should be separate migrations.

**Locking production tables.** A plain `ALTER TABLE` on a 50M-row table can lock writes
for minutes. Always check whether your database's ALTER operation is non-locking.

**Deploying code before its migration.** Application code referencing `users.email`
will crash if the migration hasn't run. Always run migrations before deploying the
application code that requires them.

## When NOT to Use

- **Schema-less or document databases** (MongoDB, DynamoDB): structure evolution is handled differently — use schema validation at the application layer and versioned document schemas
- **One-off data corrections on production**: use a documented, reviewed SQL script with a rollback plan — not a permanent migration that will run forever
- **Seed/fixture data**: use dedicated seed commands, not migration files
